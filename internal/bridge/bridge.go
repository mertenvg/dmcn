package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/keystore"
	"github.com/mertenvg/dmcn/internal/node"
)

// Config holds configuration for a bridge node.
type Config struct {
	NodeAddr       string // multiaddr of running dmcn-node to connect to
	SMTPListenAddr string // SMTP listen address (default ":2525")
	LibP2PAddr     string // libp2p listen address (default "/ip4/127.0.0.1/tcp/0")
	BridgeDomain   string // domain for bridge email addresses
	DMCNDomain     string // domain for DMCN addresses
	BridgeAddress  string // bridge's own DMCN address
	KeystorePath   string
	Passphrase     string
	PollInterval   time.Duration // how often to poll relay for outbound messages
	AuthVerifier   AuthVerifier  // nil = use stub
	Deliverer      SMTPDeliverer // nil = use stub
}

// Bridge is the SMTP-DMCN bridge node.
type Bridge struct {
	node     *node.Node
	bridgeKP *identity.IdentityKeyPair
	inbound  *InboundHandler
	outbound *OutboundHandler
	smtp     *SMTPServer
	poll     time.Duration
	log      logr.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates and initializes a bridge node. It connects to the DMCN
// network, generates or loads bridge keys, and registers the bridge identity.
func New(ctx context.Context, cfg Config, log ...logr.Logger) (*Bridge, error) {
	var l logr.Logger
	if len(log) > 0 {
		l = log[0]
	} else {
		l = logr.With(logr.M("component", "bridge"))
	}

	ctx, cancel := context.WithCancel(ctx)

	// Defaults
	if cfg.SMTPListenAddr == "" {
		cfg.SMTPListenAddr = ":2525"
	}
	if cfg.LibP2PAddr == "" {
		cfg.LibP2PAddr = "/ip4/127.0.0.1/tcp/0"
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 5 * time.Second
	}
	if cfg.AuthVerifier == nil {
		cfg.AuthVerifier = &StubAuthVerifier{
			DefaultSPF:   SPFNone,
			DefaultDKIM:  DKIMNone,
			DefaultDMARC: DMARCNone,
		}
	}
	if cfg.Deliverer == nil {
		cfg.Deliverer = &StubSMTPDeliverer{}
	}

	// Create DMCN node
	n, err := node.New(ctx, node.Config{
		ListenAddr:     cfg.LibP2PAddr,
		BootstrapPeers: []string{cfg.NodeAddr},
		KeystorePath:   cfg.KeystorePath,
		Passphrase:     cfg.Passphrase,
	}, l)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("bridge: create node: %w", err)
	}

	// Load or generate bridge keys
	bridgeKP, err := loadOrGenerateBridgeKeys(cfg, l)
	if err != nil {
		n.Close()
		cancel()
		return nil, fmt.Errorf("bridge: keys: %w", err)
	}

	// Register bridge identity (retry to allow DHT bootstrap to complete)
	rec, err := identity.NewIdentityRecord(cfg.BridgeAddress, bridgeKP)
	if err != nil {
		n.Close()
		cancel()
		return nil, fmt.Errorf("bridge: create identity record: %w", err)
	}
	rec.BridgeCapability = true
	if err := rec.Sign(bridgeKP); err != nil {
		n.Close()
		cancel()
		return nil, fmt.Errorf("bridge: sign identity: %w", err)
	}
	var regErr error
	for i := 0; i < 10; i++ {
		if regErr = n.Registry().Register(ctx, rec); regErr == nil {
			break
		}
		l.Debugf("bridge registration attempt %d failed: %v", i+1, regErr)
		select {
		case <-ctx.Done():
			n.Close()
			cancel()
			return nil, fmt.Errorf("bridge: register identity: %w", ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
	if regErr != nil {
		n.Close()
		cancel()
		return nil, fmt.Errorf("bridge: register identity: %w", regErr)
	}
	l.Successf("bridge identity registered: %s", cfg.BridgeAddress)

	// Create handlers — use direct store into the relay's message store
	// (the bridge is the relay, so we bypass the stream protocol)
	relayStore := n.Relay().Store()
	inbound := NewInboundHandler(InboundConfig{
		BridgeKP:     bridgeKP,
		BridgeAddr:   cfg.BridgeAddress,
		AuthVerifier: cfg.AuthVerifier,
		Lookup:       n.Registry().Lookup,
		Store:        relayStore.Store,
		BridgeDomain: cfg.BridgeDomain,
		DMCNDomain:   cfg.DMCNDomain,
		Log:          l,
	})

	outbound := NewOutboundHandler(OutboundConfig{
		BridgeKP:     bridgeKP,
		BridgeAddr:   cfg.BridgeAddress,
		Deliverer:    cfg.Deliverer,
		Lookup:       n.Registry().Lookup,
		BridgeDomain: cfg.BridgeDomain,
		DMCNDomain:   cfg.DMCNDomain,
		Log:          l,
	})

	smtpSrv := NewSMTPServer(ctx, cfg.SMTPListenAddr, inbound, cfg.BridgeDomain, l)

	return &Bridge{
		node:     n,
		bridgeKP: bridgeKP,
		inbound:  inbound,
		outbound: outbound,
		smtp:     smtpSrv,
		poll:     cfg.PollInterval,
		log:      l,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start begins the SMTP server and outbound relay polling.
func (b *Bridge) Start() error {
	if err := b.smtp.Start(); err != nil {
		return fmt.Errorf("bridge: start SMTP: %w", err)
	}

	go b.pollLoop()

	b.log.Info("bridge started")
	return nil
}

// pollLoop periodically fetches envelopes from the relay addressed to
// the bridge and processes them for outbound SMTP delivery.
func (b *Bridge) pollLoop() {
	ticker := time.NewTicker(b.poll)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			b.processPending()
		}
	}
}

func (b *Bridge) processPending() {
	// Read directly from the relay's message store (no stream dial needed)
	store := b.node.Relay().Store()
	bridgeXPubHex := fmt.Sprintf("%x", b.bridgeKP.X25519Public[:])
	envs, hashes := store.Fetch(bridgeXPubHex)

	for i, env := range envs {
		receipt, deliverErr := b.outbound.HandleEnvelope(b.ctx, env)
		if deliverErr != nil && receipt == nil {
			b.log.Warnf("outbound handling failed: %v", deliverErr)
			continue
		}

		// ACK the envelope directly in the store
		if err := store.Ack(hashes[i]); err != nil {
			b.log.Warnf("ack failed for %x: %v", hashes[i], err)
		}

		// Send delivery receipt back to sender
		if receipt != nil {
			b.sendReceipt(b.ctx, env, receipt)
		}
	}
}

func (b *Bridge) sendReceipt(ctx context.Context, originalEnv *message.EncryptedEnvelope, receipt *BridgeDeliveryReceipt) {
	// Decrypt original to get sender address for the receipt
	sm, err := message.Decrypt(originalEnv, b.bridgeKP.X25519Private, b.bridgeKP.X25519Public)
	if err != nil {
		b.log.Warnf("cannot decrypt for receipt: %v", err)
		return
	}

	// Look up sender to encrypt receipt to them
	senderRec, err := b.node.Registry().Lookup(ctx, sm.Plaintext.SenderAddress)
	if err != nil {
		b.log.Warnf("cannot look up sender for receipt: %v", err)
		return
	}

	receiptBytes, err := receipt.Marshal()
	if err != nil {
		b.log.Warnf("marshal receipt: %v", err)
		return
	}

	msg, err := message.NewPlaintextMessage(
		b.inbound.bridgeAddr,
		sm.Plaintext.SenderAddress,
		"Delivery Receipt",
		"Message delivery receipt attached.",
		b.bridgeKP.Ed25519Public,
	)
	if err != nil {
		b.log.Warnf("compose receipt message: %v", err)
		return
	}

	msg.Attachments = append(msg.Attachments, message.AttachmentRecord{
		Filename:    "receipt.bin",
		ContentType: ReceiptContentType,
		SizeBytes:   uint64(len(receiptBytes)),
		Content:     receiptBytes,
	})

	signed := &message.SignedMessage{Plaintext: *msg}
	if err := signed.Sign(b.bridgeKP.Ed25519Private); err != nil {
		b.log.Warnf("sign receipt: %v", err)
		return
	}

	recipients := []message.RecipientInfo{{
		DeviceID:  b.bridgeKP.DeviceID,
		X25519Pub: senderRec.X25519Public,
	}}
	env, err := message.Encrypt(signed, recipients)
	if err != nil {
		b.log.Warnf("encrypt receipt: %v", err)
		return
	}

	// Store directly into the relay's message store
	envHash := computeEnvelopeHash(env)
	store := b.node.Relay().Store()
	for _, rec := range env.Recipients {
		addr := fmt.Sprintf("%x", rec.RecipientXPub[:])
		store.Store(addr, env, envHash)
	}

	b.log.Debugf("delivery receipt sent to %s", sm.Plaintext.SenderAddress)
}

// Node returns the underlying DMCN node.
func (b *Bridge) Node() *node.Node {
	return b.node
}

// BridgeKeyPair returns the bridge's identity key pair.
func (b *Bridge) BridgeKeyPair() *identity.IdentityKeyPair {
	return b.bridgeKP
}

// Inbound returns the inbound handler for direct testing.
func (b *Bridge) Inbound() *InboundHandler {
	return b.inbound
}

// Outbound returns the outbound handler for direct testing.
func (b *Bridge) Outbound() *OutboundHandler {
	return b.outbound
}

// SMTPAddr returns the SMTP server's listen address.
func (b *Bridge) SMTPAddr() string {
	return b.smtp.Addr()
}

// Stop shuts down the bridge.
func (b *Bridge) Stop() error {
	b.cancel()
	b.smtp.Stop()
	b.node.Close()
	b.log.Info("bridge stopped")
	return nil
}

func loadOrGenerateBridgeKeys(cfg Config, log logr.Logger) (*identity.IdentityKeyPair, error) {
	if cfg.KeystorePath == "" {
		// No keystore — generate ephemeral keys
		kp, err := identity.GenerateIdentityKeyPair()
		if err != nil {
			return nil, err
		}
		log.Info("generated ephemeral bridge keys (no keystore configured)")
		return kp, nil
	}

	ks := keystore.New(cfg.KeystorePath, cfg.Passphrase)

	// Try to load existing keys
	kp, err := ks.Load(cfg.BridgeAddress)
	if err == nil {
		log.Infof("loaded bridge keys from %s", cfg.KeystorePath)
		return kp, nil
	}

	// Generate new keys and store them
	kp, err = identity.GenerateIdentityKeyPair()
	if err != nil {
		return nil, err
	}

	if err := ks.Store(cfg.BridgeAddress, kp); err != nil {
		return nil, fmt.Errorf("store bridge keys: %w", err)
	}

	log.Successf("generated and stored bridge keys in %s", cfg.KeystorePath)
	return kp, nil
}
