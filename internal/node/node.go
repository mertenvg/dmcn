// Package node provides a combined DMCN node that runs a DHT registry and
// relay service in a single process. This is the PoC development node
// described in PRD Section 5.3.
package node

import (
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/mertenvg/dmcn/internal/keystore"
	"github.com/mertenvg/dmcn/internal/registry"
	"github.com/mertenvg/dmcn/internal/relay"
)

// Config holds configuration for a DMCN node.
type Config struct {
	ListenAddr     string   // multiaddr string, e.g. "/ip4/0.0.0.0/tcp/7400"
	BootstrapPeers []string // multiaddr strings of bootstrap peers
	KeystorePath   string   // path to encrypted keystore file
	Passphrase     string   // passphrase for keystore encryption
}

// Node is a combined DMCN development node running DHT registry and relay.
type Node struct {
	host     host.Host
	registry *registry.Registry
	relay    *relay.Relay
	keystore *keystore.Keystore
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates and starts a new DMCN node.
func New(ctx context.Context, cfg Config) (*Node, error) {
	ctx, cancel := context.WithCancel(ctx)

	listenAddr, err := multiaddr.NewMultiaddr(cfg.ListenAddr)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("node: invalid listen address: %w", err)
	}

	h, err := libp2p.New(
		libp2p.ListenAddrs(listenAddr),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("node: create libp2p host: %w", err)
	}

	// Create registry
	reg, err := registry.New(ctx, h,
		registry.WithDHTMode(dht.ModeServer),
	)
	if err != nil {
		h.Close()
		cancel()
		return nil, fmt.Errorf("node: create registry: %w", err)
	}

	// Create relay with registry lookup
	rl := relay.New(h, reg.Lookup)
	rl.Start()

	// Create keystore
	var ks *keystore.Keystore
	if cfg.KeystorePath != "" {
		ks = keystore.New(cfg.KeystorePath, cfg.Passphrase)
	}

	n := &Node{
		host:     h,
		registry: reg,
		relay:    rl,
		keystore: ks,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Connect to bootstrap peers
	for _, peerAddr := range cfg.BootstrapPeers {
		if err := n.ConnectPeer(peerAddr); err != nil {
			// Non-fatal: log but continue
			fmt.Printf("node: warning: failed to connect to bootstrap peer %s: %v\n", peerAddr, err)
		}
	}

	return n, nil
}

// Host returns the underlying libp2p host.
func (n *Node) Host() host.Host {
	return n.host
}

// Registry returns the DHT identity registry.
func (n *Node) Registry() *registry.Registry {
	return n.registry
}

// Relay returns the relay service.
func (n *Node) Relay() *relay.Relay {
	return n.relay
}

// Keystore returns the encrypted keystore. May be nil if no keystore path
// was configured.
func (n *Node) Keystore() *keystore.Keystore {
	return n.keystore
}

// PeerID returns the node's libp2p peer ID.
func (n *Node) PeerID() peer.ID {
	return n.host.ID()
}

// Addrs returns the node's listen multiaddrs with peer ID included.
func (n *Node) Addrs() []string {
	hostAddr := n.host.Addrs()
	peerInfo := peer.AddrInfo{
		ID:    n.host.ID(),
		Addrs: hostAddr,
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return nil
	}
	result := make([]string, len(addrs))
	for i, a := range addrs {
		result[i] = a.String()
	}
	return result
}

// ConnectPeer connects to a peer by multiaddr string.
func (n *Node) ConnectPeer(addr string) error {
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return fmt.Errorf("node: invalid peer address: %w", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return fmt.Errorf("node: parse peer info: %w", err)
	}

	if err := n.host.Connect(n.ctx, *peerInfo); err != nil {
		return fmt.Errorf("node: connect: %w", err)
	}

	return nil
}

// Close shuts down the node, stopping the relay and registry.
func (n *Node) Close() error {
	n.relay.Stop()
	n.registry.Close()
	n.host.Close()
	n.cancel()
	return nil
}
