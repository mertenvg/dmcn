package bridge_test

import (
	"context"
	"testing"
	"time"

	"github.com/mertenvg/dmcn/internal/bridge"
	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/node"
)

// createConnectedNodes creates two connected DMCN nodes and waits for
// DHT bootstrap. Returns the nodes and a cleanup function.
func createConnectedNodes(t *testing.T, ctx context.Context) (*node.Node, *node.Node) {
	t.Helper()

	nodeA, err := node.New(ctx, node.Config{
		ListenAddr: "/ip4/127.0.0.1/tcp/0",
	})
	if err != nil {
		t.Fatalf("create node-A: %v", err)
	}

	nodeB, err := node.New(ctx, node.Config{
		ListenAddr: "/ip4/127.0.0.1/tcp/0",
	})
	if err != nil {
		nodeA.Close()
		t.Fatalf("create node-B: %v", err)
	}

	// Connect nodes
	addrs := nodeB.Addrs()
	if len(addrs) == 0 {
		nodeA.Close()
		nodeB.Close()
		t.Fatal("node-B has no addresses")
	}
	if err := nodeA.ConnectPeer(addrs[0]); err != nil {
		nodeA.Close()
		nodeB.Close()
		t.Fatalf("connect A→B: %v", err)
	}

	time.Sleep(500 * time.Millisecond) // DHT bootstrap

	return nodeA, nodeB
}

func registerIdentity(t *testing.T, ctx context.Context, n *node.Node, address string) *identity.IdentityKeyPair {
	t.Helper()

	kp, err := identity.GenerateIdentityKeyPair()
	if err != nil {
		t.Fatalf("generate keys for %s: %v", address, err)
	}
	rec, err := identity.NewIdentityRecord(address, kp)
	if err != nil {
		t.Fatalf("create record for %s: %v", address, err)
	}
	if err := rec.Sign(kp); err != nil {
		t.Fatalf("sign %s: %v", address, err)
	}
	if err := n.Registry().Register(ctx, rec); err != nil {
		t.Fatalf("register %s: %v", address, err)
	}
	return kp
}

// TestBridgeEndToEnd is the PRD Section 6.4 end-to-end integration test.
//
// Scenario:
// 1. Start two connected dmcn-nodes (DHT + relay).
// 2. Register alice@dmcn.localhost on nodeA.
// 3. Create bridge connected to nodeB with stub auth verifier + stub SMTP deliverer.
// 4. Inbound: simulate email from external@gmail.com to alice@bridge.localhost.
// 5. Alice fetches, decrypts, verifies classification record.
// 6. Outbound: Alice replies to external@gmail.com via bridge.
// 7. Verify stub SMTP deliverer captured the outbound message.
// 8. Alice fetches delivery receipt.
func TestBridgeEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 1: Start two connected nodes
	nodeA, nodeB := createConnectedNodes(t, ctx)
	defer nodeA.Close()
	defer nodeB.Close()

	// Step 2: Register Alice on nodeA
	aliceKP := registerIdentity(t, ctx, nodeA, "alice@dmcn.localhost")

	// Step 3: Create bridge connected to nodeB
	stubDeliverer := &bridge.StubSMTPDeliverer{}
	nodeBAddrs := nodeB.Addrs()

	b, err := bridge.New(ctx, bridge.Config{
		NodeAddr:       nodeBAddrs[0],
		SMTPListenAddr: "127.0.0.1:0",
		LibP2PAddr:     "/ip4/127.0.0.1/tcp/0",
		BridgeDomain:   "bridge.localhost",
		DMCNDomain:     "dmcn.localhost",
		BridgeAddress:  "bridge@bridge.localhost",
		PollInterval:   100 * time.Millisecond,
		AuthVerifier: &bridge.StubAuthVerifier{
			DefaultSPF:   bridge.SPFPass,
			DefaultDKIM:  bridge.DKIMPass,
			DefaultDMARC: bridge.DMARCPass,
		},
		Deliverer: stubDeliverer,
	})
	if err != nil {
		t.Fatalf("create bridge: %v", err)
	}
	defer b.Stop()

	if err := b.Start(); err != nil {
		t.Fatalf("start bridge: %v", err)
	}

	// Give DHT time to propagate bridge identity
	time.Sleep(1 * time.Second)

	// Create Alice's client node connected to the bridge
	bridgeNode := b.Node()
	bridgeAddrs := bridgeNode.Addrs()
	if len(bridgeAddrs) == 0 {
		t.Fatal("bridge node has no addresses")
	}
	aliceNode, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{bridgeAddrs[0]},
	})
	if err != nil {
		t.Fatalf("create alice node: %v", err)
	}
	defer aliceNode.Close()
	time.Sleep(500 * time.Millisecond)

	// Step 4: Inbound — simulate email from external@gmail.com to alice@bridge.localhost
	err = b.Inbound().HandleMessage(ctx, "1.2.3.4", "external@gmail.com", "alice@bridge.localhost", []byte("Hello Alice from legacy email!"))
	if err != nil {
		t.Fatalf("inbound handle: %v", err)
	}

	// Step 5: Alice fetches from the bridge's relay
	envs, hashes, err := aliceNode.Relay().ClientFetch(ctx, bridgeNode.PeerID(), aliceKP, "alice@dmcn.localhost")
	if err != nil {
		t.Fatalf("alice fetch: %v", err)
	}
	if len(envs) == 0 {
		t.Fatal("alice expected at least 1 envelope, got 0")
	}

	// Decrypt and verify
	sm, err := message.Decrypt(envs[0], aliceKP.X25519Private, aliceKP.X25519Public)
	if err != nil {
		t.Fatalf("alice decrypt: %v", err)
	}
	if err := sm.Verify(); err != nil {
		t.Fatalf("alice verify sender: %v", err)
	}

	// Verify message content
	if string(sm.Plaintext.Body.Content) != "Hello Alice from legacy email!" {
		t.Fatalf("unexpected body: %q", string(sm.Plaintext.Body.Content))
	}

	// Check classification record attachment
	if len(sm.Plaintext.Attachments) == 0 {
		t.Fatal("expected classification record attachment")
	}
	classAtt := sm.Plaintext.Attachments[0]
	if classAtt.ContentType != bridge.ClassificationContentType {
		t.Fatalf("unexpected attachment type: %s", classAtt.ContentType)
	}
	classRec, err := bridge.UnmarshalClassificationRecord(classAtt.Content)
	if err != nil {
		t.Fatalf("unmarshal classification: %v", err)
	}
	if err := classRec.Verify(); err != nil {
		t.Fatalf("verify classification signature: %v", err)
	}
	if classRec.TrustTier != bridge.TrustTierVerifiedLegacy {
		t.Fatalf("expected TrustTierVerifiedLegacy, got %d", classRec.TrustTier)
	}
	if classRec.SMTPFrom != "external@gmail.com" {
		t.Fatalf("expected smtp_from 'external@gmail.com', got %q", classRec.SMTPFrom)
	}

	// ACK
	if err := aliceNode.Relay().ClientAck(ctx, bridgeNode.PeerID(), hashes[0]); err != nil {
		t.Fatalf("alice ack: %v", err)
	}

	// Step 6: Outbound — Alice replies to external@gmail.com via bridge
	replyMsg, err := message.NewPlaintextMessage(
		"alice@dmcn.localhost",
		"external@gmail.com",
		"Re: Hello",
		"Hello from DMCN!",
		aliceKP.Ed25519Public,
	)
	if err != nil {
		t.Fatalf("compose reply: %v", err)
	}

	replySigned := &message.SignedMessage{Plaintext: *replyMsg}
	if err := replySigned.Sign(aliceKP.Ed25519Private); err != nil {
		t.Fatalf("sign reply: %v", err)
	}

	// Encrypt to bridge's public key
	bridgeKP := b.BridgeKeyPair()
	recipients := []message.RecipientInfo{{
		DeviceID:  aliceKP.DeviceID,
		X25519Pub: bridgeKP.X25519Public,
	}}
	replyEnv, err := message.Encrypt(replySigned, recipients)
	if err != nil {
		t.Fatalf("encrypt reply: %v", err)
	}

	// Store on the bridge's relay for the bridge to pick up
	_, err = aliceNode.Relay().ClientStoreWithAddress(ctx, bridgeNode.PeerID(), "alice@dmcn.localhost", aliceKP, replyEnv)
	if err != nil {
		t.Fatalf("store reply: %v", err)
	}

	// Step 7: Wait for bridge poll to process the outbound message
	time.Sleep(2 * time.Second)

	// Verify stub deliverer captured the message
	if len(stubDeliverer.Messages) == 0 {
		t.Fatal("expected outbound SMTP delivery, got none")
	}
	delivered := stubDeliverer.Messages[0]
	if delivered.To != "external@gmail.com" {
		t.Fatalf("expected To 'external@gmail.com', got %q", delivered.To)
	}
	if delivered.From != "alice@bridge.localhost" {
		t.Fatalf("expected From 'alice@bridge.localhost', got %q", delivered.From)
	}
	if delivered.Body != "Hello from DMCN!" {
		t.Fatalf("unexpected body: %q", delivered.Body)
	}

	// Step 8: Alice fetches delivery receipt
	time.Sleep(1 * time.Second)
	envs2, _, err := aliceNode.Relay().ClientFetch(ctx, bridgeNode.PeerID(), aliceKP, "alice@dmcn.localhost")
	if err != nil {
		t.Fatalf("alice fetch receipt: %v", err)
	}
	if len(envs2) == 0 {
		t.Fatal("expected delivery receipt envelope, got 0")
	}

	sm2, err := message.Decrypt(envs2[0], aliceKP.X25519Private, aliceKP.X25519Public)
	if err != nil {
		t.Fatalf("decrypt receipt: %v", err)
	}
	if err := sm2.Verify(); err != nil {
		t.Fatalf("verify receipt message: %v", err)
	}
	if len(sm2.Plaintext.Attachments) == 0 {
		t.Fatal("expected receipt attachment")
	}
	receiptAtt := sm2.Plaintext.Attachments[0]
	if receiptAtt.ContentType != bridge.ReceiptContentType {
		t.Fatalf("unexpected receipt attachment type: %s", receiptAtt.ContentType)
	}
	receipt, err := bridge.UnmarshalDeliveryReceipt(receiptAtt.Content)
	if err != nil {
		t.Fatalf("unmarshal receipt: %v", err)
	}
	if !receipt.Success {
		t.Fatalf("receipt indicates failure: %s", receipt.ErrorDetail)
	}
	if receipt.RecipientEmail != "external@gmail.com" {
		t.Fatalf("receipt recipient: got %q, want 'external@gmail.com'", receipt.RecipientEmail)
	}
}

// TestBridgeInboundClassificationTiers tests that different auth results
// produce different trust tiers in the classification record.
func TestBridgeInboundClassificationTiers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Start two connected nodes
	nodeA, nodeB := createConnectedNodes(t, ctx)
	defer nodeA.Close()
	defer nodeB.Close()

	// Register recipient on nodeA
	aliceKP := registerIdentity(t, ctx, nodeA, "alice@dmcn.localhost")

	tests := []struct {
		name     string
		spf      bridge.SPFResult
		dkim     bridge.DKIMResult
		dmarc    bridge.DMARCResult
		wantTier bridge.BridgeTrustTier
	}{
		{"verified", bridge.SPFPass, bridge.DKIMPass, bridge.DMARCPass, bridge.TrustTierVerifiedLegacy},
		{"unverified", bridge.SPFNone, bridge.DKIMNone, bridge.DMARCNone, bridge.TrustTierUnverifiedLegacy},
		{"suspicious", bridge.SPFPass, bridge.DKIMFail, bridge.DMARCFail, bridge.TrustTierSuspicious},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeBAddrs := nodeB.Addrs()
			b, err := bridge.New(ctx, bridge.Config{
				NodeAddr:       nodeBAddrs[0],
				SMTPListenAddr: "127.0.0.1:0",
				LibP2PAddr:     "/ip4/127.0.0.1/tcp/0",
				BridgeDomain:   "bridge.localhost",
				DMCNDomain:     "dmcn.localhost",
				BridgeAddress:  "bridge@bridge.localhost",
				AuthVerifier: &bridge.StubAuthVerifier{
					DefaultSPF:   tt.spf,
					DefaultDKIM:  tt.dkim,
					DefaultDMARC: tt.dmarc,
				},
				Deliverer: &bridge.StubSMTPDeliverer{},
			})
			if err != nil {
				t.Fatalf("create bridge: %v", err)
			}
			defer b.Stop()

			// Create Alice's client node connected to bridge
			bridgeNode := b.Node()
			bridgeAddrs := bridgeNode.Addrs()
			aliceClient, err := node.New(ctx, node.Config{
				ListenAddr:     "/ip4/127.0.0.1/tcp/0",
				BootstrapPeers: []string{bridgeAddrs[0]},
			})
			if err != nil {
				t.Fatalf("create alice client: %v", err)
			}
			defer aliceClient.Close()

			time.Sleep(500 * time.Millisecond) // DHT propagation

			err = b.Inbound().HandleMessage(ctx, "1.2.3.4", "sender@test.com", "alice@bridge.localhost", []byte("test"))
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			envs, _, err := aliceClient.Relay().ClientFetch(ctx, bridgeNode.PeerID(), aliceKP, "alice@dmcn.localhost")
			if err != nil {
				t.Fatalf("fetch: %v", err)
			}
			if len(envs) == 0 {
				t.Fatal("no envelopes")
			}

			sm, err := message.Decrypt(envs[0], aliceKP.X25519Private, aliceKP.X25519Public)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if len(sm.Plaintext.Attachments) == 0 {
				t.Fatal("no classification attachment")
			}

			rec, err := bridge.UnmarshalClassificationRecord(sm.Plaintext.Attachments[0].Content)
			if err != nil {
				t.Fatal(err)
			}
			if rec.TrustTier != tt.wantTier {
				t.Errorf("trust tier: got %d, want %d", rec.TrustTier, tt.wantTier)
			}
		})
	}
}
