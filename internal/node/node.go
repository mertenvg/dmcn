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
	"github.com/mertenvg/logr/v2"
	"github.com/multiformats/go-multiaddr"

	"github.com/mertenvg/dmcn/internal/keystore"
	"github.com/mertenvg/dmcn/internal/registry"
	"github.com/mertenvg/dmcn/internal/relay"
)

// Config holds configuration for a DMCN node.
type Config struct {
	ListenAddr     string   // multiaddr string, e.g. "/ip4/0.0.0.0/tcp/7400"
	BootstrapPeers []string // multiaddr strings of bootstrap peers
	OrgPeers       []string // multiaddr strings of organizational peers (relay fallbacks)
	KeystorePath   string   // path to encrypted keystore file
	Passphrase     string   // passphrase for keystore encryption
}

// Node is a combined DMCN development node running DHT registry and relay.
type Node struct {
	host     host.Host
	registry *registry.Registry
	relay    *relay.Relay
	keystore *keystore.Keystore
	orgPeers []string
	log      logr.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates and starts a new DMCN node.
func New(ctx context.Context, cfg Config, log ...logr.Logger) (*Node, error) {
	var l logr.Logger
	if len(log) > 0 {
		l = log[0]
	} else {
		l = logr.With(logr.M("component", "node"))
	}

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
	rl := relay.New(h, reg.Lookup, relay.WithOrgPeers(cfg.OrgPeers))
	rl.Start()

	// If org peers provided but no bootstrap peers, use org peers for DHT bootstrap.
	bootstrapPeers := cfg.BootstrapPeers
	if len(bootstrapPeers) == 0 && len(cfg.OrgPeers) > 0 {
		bootstrapPeers = cfg.OrgPeers
	}

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
		orgPeers: cfg.OrgPeers,
		log:      l,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Connect to bootstrap peers
	for _, peerAddr := range bootstrapPeers {
		if err := n.ConnectPeer(peerAddr); err != nil {
			// Non-fatal: log but continue
			l.Warnf("failed to connect to bootstrap peer %s: %v", peerAddr, err)
		}
	}

	// Discover additional org peers from connected org peers.
	if len(cfg.OrgPeers) > 0 {
		n.discoverOrgPeers(ctx, cfg.OrgPeers)
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

// RelayHints returns the node's own addresses plus org peers, suitable for
// populating IdentityRecord.RelayHints.
func (n *Node) RelayHints() []string {
	return append(n.Addrs(), n.orgPeers...)
}

// ParseRelayHint parses a relay hint multiaddr string into peer.AddrInfo.
func ParseRelayHint(hint string) (*peer.AddrInfo, error) {
	ma, err := multiaddr.NewMultiaddr(hint)
	if err != nil {
		return nil, fmt.Errorf("invalid relay hint: %w", err)
	}
	info, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return nil, fmt.Errorf("parse relay hint peer info: %w", err)
	}
	return info, nil
}

// discoverOrgPeers queries connected org peers for the full cluster list
// and connects to any newly discovered peers.
func (n *Node) discoverOrgPeers(ctx context.Context, initialPeers []string) {
	known := make(map[string]bool)
	for _, p := range initialPeers {
		known[p] = true
	}

	for _, peerAddr := range initialPeers {
		info, err := ParseRelayHint(peerAddr)
		if err != nil {
			continue
		}

		discovered, err := n.relay.ClientOrgPeers(ctx, info.ID)
		if err != nil {
			n.log.Warnf("failed to discover org peers from %s: %v", peerAddr, err)
			continue
		}

		for _, dp := range discovered {
			if known[dp] {
				continue
			}
			known[dp] = true
			n.orgPeers = append(n.orgPeers, dp)
			if err := n.ConnectPeer(dp); err != nil {
				n.log.Warnf("failed to connect to discovered org peer %s: %v", dp, err)
			}
		}
	}
}

// Close shuts down the node, stopping the relay and registry.
func (n *Node) Close() error {
	n.relay.Stop()
	n.registry.Close()
	n.host.Close()
	n.cancel()
	return nil
}
