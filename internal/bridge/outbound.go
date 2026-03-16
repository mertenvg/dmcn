package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
)

// OutboundHandler processes DMCN messages addressed to legacy email
// recipients and delivers them via SMTP.
type OutboundHandler struct {
	bridgeKP     *identity.IdentityKeyPair
	bridgeAddr   string
	deliverer    SMTPDeliverer
	lookup       LookupFunc
	bridgeDomain string
	dmcnDomain   string
	log          logr.Logger
}

// OutboundConfig configures the outbound handler.
type OutboundConfig struct {
	BridgeKP     *identity.IdentityKeyPair
	BridgeAddr   string
	Deliverer    SMTPDeliverer
	Lookup       LookupFunc
	BridgeDomain string
	DMCNDomain   string
	Log          logr.Logger
}

// NewOutboundHandler creates a new outbound message handler.
func NewOutboundHandler(cfg OutboundConfig) *OutboundHandler {
	return &OutboundHandler{
		bridgeKP:     cfg.BridgeKP,
		bridgeAddr:   cfg.BridgeAddr,
		deliverer:    cfg.Deliverer,
		lookup:       cfg.Lookup,
		bridgeDomain: cfg.BridgeDomain,
		dmcnDomain:   cfg.DMCNDomain,
		log:          cfg.Log,
	}
}

// HandleEnvelope decrypts a DMCN envelope addressed to the bridge,
// verifies the sender, delivers the message via SMTP, and returns a
// signed delivery receipt.
func (h *OutboundHandler) HandleEnvelope(ctx context.Context, env *message.EncryptedEnvelope) (*BridgeDeliveryReceipt, error) {
	// 1. Decrypt
	sm, err := message.Decrypt(env, h.bridgeKP.X25519Private, h.bridgeKP.X25519Public)
	if err != nil {
		return nil, fmt.Errorf("bridge: decrypt: %w", err)
	}

	// 2. Verify sender signature
	if err := sm.Verify(); err != nil {
		return nil, fmt.Errorf("bridge: verify sender: %w", err)
	}

	// 3. Log warning — PRD requirement: bridge must log when decrypting
	// message content for outbound delivery.
	h.log.Warnf("TRUST DISCLOSURE: decrypting message from %s for outbound SMTP delivery to %s",
		sm.Plaintext.SenderAddress, sm.Plaintext.RecipientAddress)

	// 4. Verify sender exists in registry
	senderRec, err := h.lookup(ctx, sm.Plaintext.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrSenderNotFound, sm.Plaintext.SenderAddress, err)
	}
	_ = senderRec // verified existence; public key already validated via signature

	// 5. Check recipient is a legacy address
	recipientAddr := sm.Plaintext.RecipientAddress
	if !IsLegacyAddress(recipientAddr, h.bridgeDomain, h.dmcnDomain) {
		return nil, fmt.Errorf("%w: %s", ErrNotLegacyAddress, recipientAddr)
	}

	// 6. Deliver via SMTP
	smtpFrom := DMCNToSMTPFrom(sm.Plaintext.SenderAddress, h.bridgeDomain)
	deliverErr := h.deliverer.Deliver(ctx, smtpFrom, recipientAddr, sm.Plaintext.Subject, string(sm.Plaintext.Body.Content))

	// 7. Construct delivery receipt
	receipt := &BridgeDeliveryReceipt{
		OriginalMessageID: sm.Plaintext.MessageID,
		RecipientEmail:    recipientAddr,
		BridgeAddress:     h.bridgeAddr,
		DeliveredAt:       time.Now().UTC(),
		Success:           deliverErr == nil,
	}
	if deliverErr != nil {
		receipt.ErrorDetail = deliverErr.Error()
		h.log.Warnf("outbound delivery failed to %s: %v", recipientAddr, deliverErr)
	} else {
		h.log.Infof("outbound message delivered from %s to %s via SMTP", sm.Plaintext.SenderAddress, recipientAddr)
	}

	if err := receipt.Sign(h.bridgeKP.Ed25519Private); err != nil {
		return nil, fmt.Errorf("bridge: sign receipt: %w", err)
	}

	return receipt, deliverErr
}
