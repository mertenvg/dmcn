// Binary dmcn-web is the DMCN web client backend. It serves the embedded
// frontend, exposes REST and WebSocket APIs for messaging, and connects to
// a running dmcn-node for DHT registry and relay access.
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/api"
	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/server"
	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/ws"
	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/node"
	"github.com/mertenvg/dmcn/internal/relay"
)

//go:embed web/dist
var frontendFS embed.FS

var (
	version = "dev"
	log     logr.Logger
)

// nodeRelayRouter implements api.RelayRouter by connecting to peers and
// storing pre-signed envelopes via the node's relay client methods.
type nodeRelayRouter struct {
	node *node.Node
}

func (r *nodeRelayRouter) ConnectPeer(addr string) error {
	return r.node.ConnectPeer(addr)
}

func (r *nodeRelayRouter) StorePreSignedOnPeer(ctx context.Context, hint string, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error) {
	info, err := node.ParseRelayHint(hint)
	if err != nil {
		return [32]byte{}, err
	}
	return r.node.Relay().ClientStorePreSigned(ctx, info.ID, senderAddr, signature, env)
}

// relayAdapter bridges the node's relay Client* methods to the RelayProxy
// interface expected by the WebSocket ConnManager.
type relayAdapter struct {
	relay          *relay.Relay
	node           *node.Node
	registryLookup func(ctx context.Context, address string) (*identity.IdentityRecord, error)
}

func (ra *relayAdapter) FetchChallenge(ctx context.Context, address string) ([]byte, network.Stream, error) {
	// Look up the user's identity record to find relay hints.
	if ra.registryLookup != nil {
		rec, err := ra.registryLookup(ctx, address)
		if err == nil && len(rec.RelayHints) > 0 {
			// Try relay hints in order.
			for _, hint := range rec.RelayHints {
				info, parseErr := node.ParseRelayHint(hint)
				if parseErr != nil {
					continue
				}
				if connectErr := ra.node.ConnectPeer(hint); connectErr != nil {
					continue
				}
				nonce, stream, fetchErr := ra.relay.ClientFetchChallenge(ctx, info.ID, address)
				if fetchErr == nil {
					return nonce, stream, nil
				}
			}
			return nil, nil, fmt.Errorf("all relay hints failed for %s", address)
		}
	}
	// Fallback to local node.
	return ra.relay.ClientFetchChallenge(ctx, ra.node.PeerID(), address)
}

func (ra *relayAdapter) FetchComplete(stream network.Stream, address string, nonce, signature []byte) ([]*message.EncryptedEnvelope, [][32]byte, error) {
	return ra.relay.ClientFetchComplete(stream, address, nonce, signature)
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		fmt.Println("dmcn-web", version)
		return
	}

	logr.AddWriter(os.Stderr, logr.WithFormatter(logr.FormatWithColours), logr.WithFilter(logr.Verbose))
	log = logr.With(logr.M("component", "web"))

	// Read configuration from environment variables.
	listenAddr := envOrDefault("DMCN_WEB_LISTEN", ":8443")
	domain := envOrDefault("DMCN_WEB_DOMAIN", "dmcn.me")
	nodeAddr := os.Getenv("DMCN_WEB_NODE")
	dataDir := envOrDefault("DMCN_WEB_DATA_DIR", "data")
	tlsCert := os.Getenv("DMCN_WEB_TLS_CERT")
	tlsKey := os.Getenv("DMCN_WEB_TLS_KEY")
	devMode := os.Getenv("DMCN_WEB_DEV") == "true" || os.Getenv("DMCN_WEB_DEV") == "1"
	pollIntervalStr := envOrDefault("DMCN_WEB_POLL_INTERVAL", "10s")
	orgPeersEnv := os.Getenv("DMCN_WEB_ORG_PEERS")

	if nodeAddr == "" {
		fmt.Fprintln(os.Stderr, "DMCN_WEB_NODE is required (multiaddr of running dmcn-node)")
		logr.Wait()
		os.Exit(1)
	}

	pollInterval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid DMCN_WEB_POLL_INTERVAL: %v\n", err)
		logr.Wait()
		os.Exit(1)
	}

	// In dev mode, default to local TLS certificates.
	if devMode {
		if tlsCert == "" {
			tlsCert = "certs/localhost.crt"
		}
		if tlsKey == "" {
			tlsKey = "certs/localhost.key"
		}
	}

	// Create context with cancellation for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse org peers from environment.
	var orgPeers []string
	if orgPeersEnv != "" {
		orgPeers = strings.Split(orgPeersEnv, ",")
	}

	// Create DMCN node connecting to the specified bootstrap peer.
	n, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{nodeAddr},
		OrgPeers:       orgPeers,
	})
	if err != nil {
		log.Errorf("failed to create node: %v", err)
		logr.Wait()
		os.Exit(1)
	}
	defer n.Close()

	log.Infof("connected to DMCN network, local peer ID: %s", n.PeerID())

	// Create stores.
	userStore, err := store.NewUserStore(dataDir)
	if err != nil {
		log.Errorf("failed to create user store: %v", err)
		logr.Wait()
		os.Exit(1)
	}

	envelopeStore, err := store.NewEnvelopeStore(dataDir)
	if err != nil {
		log.Errorf("failed to create envelope store: %v", err)
		logr.Wait()
		os.Exit(1)
	}

	sessionStore := store.NewSessionStore(24 * time.Hour)

	// Build closures that the API handlers need.
	registryRegister := func(ctx context.Context, rec *identity.IdentityRecord) error {
		return n.Registry().Register(ctx, rec)
	}

	storePreSigned := func(ctx context.Context, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error) {
		return n.Relay().ClientStorePreSigned(ctx, n.PeerID(), senderAddr, signature, env)
	}

	registryLookup := func(ctx context.Context, address string) (*identity.IdentityRecord, error) {
		return n.Registry().Lookup(ctx, address)
	}

	relayHints := func() []string {
		return append([]string{nodeAddr}, orgPeers...)
	}

	// Create API handlers.
	authHandler := api.NewAuthHandler(userStore, sessionStore, registryRegister, log)
	msgHandler := api.NewMessageHandler(userStore, sessionStore, envelopeStore, storePreSigned, registryLookup, &nodeRelayRouter{node: n}, log)
	identHandler := api.NewIdentityHandler(registryLookup, relayHints, log)
	contactHandler := api.NewContactHandler(userStore, sessionStore, log)

	// Create WebSocket connection manager.
	relayProxy := &relayAdapter{
		relay:          n.Relay(),
		node:           n,
		registryLookup: registryLookup,
	}
	connManager := ws.NewConnManager(sessionStore, relayProxy, envelopeStore, pollInterval, log)

	// Create server.
	srv := server.New(server.Config{
		ListenAddr: listenAddr,
		Domain:     domain,
		TLSCert:    tlsCert,
		TLSKey:     tlsKey,
		DevMode:    devMode,
		DataDir:    dataDir,
	}, log)

	// Build sub-FS for embedded frontend.
	subFS, err := fs.Sub(frontendFS, "web/dist")
	if err != nil {
		log.Errorf("failed to create frontend sub-FS: %v", err)
		logr.Wait()
		os.Exit(1)
	}

	// Register API routes, WebSocket handler, and frontend.
	authMiddleware := server.AuthMiddleware(sessionStore)
	srv.RegisterAPI(authHandler, msgHandler, identHandler, contactHandler, connManager.HandleUpgrade, authMiddleware, subFS)

	// Start server in a goroutine.
	go func() {
		var err error
		if tlsCert == "" && tlsKey == "" && !devMode {
			err = srv.StartAutocert(domain, filepath.Join(dataDir, "certs"))
		} else {
			err = srv.Start(tlsCert, tlsKey)
		}
		if err != nil {
			log.Errorf("server error: %v", err)
			cancel()
		}
	}()

	log.Infof("DMCN web client listening on %s", listenAddr)

	// Wait for interrupt signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Infof("received signal %s, shutting down...", sig)
	case <-ctx.Done():
	}

	// Graceful shutdown with timeout.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorf("server shutdown error: %v", err)
	}

	log.Info("DMCN web client stopped")
	logr.Wait()
}

// envOrDefault returns the value of the environment variable named by key,
// or defaultVal if the variable is empty or unset.
func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
