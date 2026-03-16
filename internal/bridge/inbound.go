package bridge

import (
	"context"
	"fmt"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
)

// StoreFunc stores an encrypted envelope for a recipient. The recipientAddr
// is the hex-encoded X25519 public key of the recipient. Returns the
// envelope hash.
type StoreFunc func(recipientAddr string, env *message.EncryptedEnvelope, hash [32]byte)

// LookupFunc looks up an identity record by address from the registry.
type LookupFunc func(ctx context.Context, address string) (*identity.IdentityRecord, error)

// InboundHandler processes inbound SMTP messages and delivers them as
// encrypted DMCN envelopes to the recipient's relay node.
type InboundHandler struct {
	bridgeKP     *identity.IdentityKeyPair
	bridgeAddr   string
	authVerifier AuthVerifier
	lookup       LookupFunc
	store        StoreFunc
	bridgeDomain string
	dmcnDomain   string
	log          logr.Logger
}

// InboundConfig configures the inbound handler.
type InboundConfig struct {
	BridgeKP     *identity.IdentityKeyPair
	BridgeAddr   string
	AuthVerifier AuthVerifier
	Lookup       LookupFunc
	Store        StoreFunc
	BridgeDomain string
	DMCNDomain   string
	Log          logr.Logger
}

// NewInboundHandler creates a new inbound message handler.
func NewInboundHandler(cfg InboundConfig) *InboundHandler {
	return &InboundHandler{
		bridgeKP:     cfg.BridgeKP,
		bridgeAddr:   cfg.BridgeAddr,
		authVerifier: cfg.AuthVerifier,
		lookup:       cfg.Lookup,
		store:        cfg.Store,
		bridgeDomain: cfg.BridgeDomain,
		dmcnDomain:   cfg.DMCNDomain,
		log:          cfg.Log,
	}
}

// HandleMessage processes an inbound SMTP message, classifies it, wraps it
// in a DMCN envelope, and stores it on the relay.
func (h *InboundHandler) HandleMessage(ctx context.Context, senderIP, from, to string, rawMsg []byte) error {
	// 1. Verify authentication
	authResult, err := h.authVerifier.Verify(ctx, senderIP, from, rawMsg)
	if err != nil {
		return fmt.Errorf("bridge: auth verify: %w", err)
	}

	// 2. Classify
	tier := Classify(authResult)
	h.log.Debugf("classified %s from %s as tier %d", to, from, tier)

	// 3. Construct and sign classification record
	classRec := NewClassificationRecord(h.bridgeAddr, h.bridgeKP.Ed25519Public, from, authResult, tier)
	if err := classRec.Sign(h.bridgeKP.Ed25519Private); err != nil {
		return fmt.Errorf("bridge: sign classification: %w", err)
	}

	classBytes, err := classRec.Marshal()
	if err != nil {
		return fmt.Errorf("bridge: marshal classification: %w", err)
	}

	// 4. Map bridge address to DMCN address
	dmcnAddr := SMTPToDMCN(to, h.bridgeDomain, h.dmcnDomain)

	// 5. Look up recipient
	recipientRec, err := h.lookup(ctx, dmcnAddr)
	if err != nil {
		return fmt.Errorf("%w: %s: %v", ErrRecipientNotFound, dmcnAddr, err)
	}

	// 6. Create PlaintextMessage with classification attachment
	msg, err := message.NewPlaintextMessage(
		h.bridgeAddr,
		dmcnAddr,
		fmt.Sprintf("Bridged message from %s", from),
		string(rawMsg),
		h.bridgeKP.Ed25519Public,
	)
	if err != nil {
		return fmt.Errorf("bridge: compose message: %w", err)
	}

	classHash := crypto.SHA256Hash(classBytes)
	attID, err := crypto.RandomUUID()
	if err != nil {
		return fmt.Errorf("bridge: generate attachment ID: %w", err)
	}
	msg.Attachments = append(msg.Attachments, message.AttachmentRecord{
		AttachmentID: attID,
		Filename:     "classification.bin",
		ContentType:  ClassificationContentType,
		SizeBytes:    uint64(len(classBytes)),
		ContentHash:  classHash,
		Content:      classBytes,
	})

	// 7. Sign as bridge
	sm := &message.SignedMessage{Plaintext: *msg}
	if err := sm.Sign(h.bridgeKP.Ed25519Private); err != nil {
		return fmt.Errorf("bridge: sign message: %w", err)
	}

	// 8. Encrypt to recipient
	recipients := []message.RecipientInfo{{
		DeviceID:  h.bridgeKP.DeviceID,
		X25519Pub: recipientRec.X25519Public,
	}}
	env, err := message.Encrypt(sm, recipients)
	if err != nil {
		return fmt.Errorf("bridge: encrypt: %w", err)
	}

	// 9. Store on relay — compute hash and store directly
	envHash := computeEnvelopeHash(env)
	for _, rec := range env.Recipients {
		addr := fmt.Sprintf("%x", rec.RecipientXPub[:])
		h.store(addr, env, envHash)
	}

	h.log.Infof("inbound message from %s to %s stored, hash: %x", from, dmcnAddr, envHash)
	return nil
}

// computeEnvelopeHash computes the SHA-256 hash of an envelope's proto bytes.
func computeEnvelopeHash(env *message.EncryptedEnvelope) [32]byte {
	pb := env.ToProto()
	data, err := protoMarshal(pb)
	if err != nil {
		// This should not happen with valid envelopes
		return [32]byte{}
	}
	return crypto.SHA256Hash(data)
}
