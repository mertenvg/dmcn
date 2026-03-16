package bridge

import (
	"crypto/ed25519"
	"time"
)

// Classify determines the trust tier for an inbound legacy email based
// on its authentication results.
//
// Classification rules:
//   - DKIM pass AND DMARC pass → TrustTierVerifiedLegacy
//   - Any DKIM or DMARC failure → TrustTierSuspicious
//   - Otherwise (missing/neutral checks) → TrustTierUnverifiedLegacy
func Classify(result *AuthResult) BridgeTrustTier {
	if result.DKIM == DKIMFail || result.DMARC == DMARCFail {
		return TrustTierSuspicious
	}
	if result.DKIM == DKIMPass && result.DMARC == DMARCPass {
		return TrustTierVerifiedLegacy
	}
	return TrustTierUnverifiedLegacy
}

// NewClassificationRecord constructs a BridgeClassificationRecord from the
// authentication result and trust tier. The record is unsigned; call Sign()
// before attaching to a message.
func NewClassificationRecord(
	bridgeAddr string,
	bridgePubKey ed25519.PublicKey,
	smtpFrom string,
	authResult *AuthResult,
	tier BridgeTrustTier,
) *BridgeClassificationRecord {
	return &BridgeClassificationRecord{
		BridgeAddress:   bridgeAddr,
		BridgePublicKey: bridgePubKey,
		SMTPFrom:        smtpFrom,
		SMTPSenderIP:    authResult.SenderIP,
		SPFResult:       authResult.SPF,
		DKIMResult:      authResult.DKIM,
		DMARCResult:     authResult.DMARC,
		ReputationScore: 0, // stubbed for PoC
		TrustTier:       tier,
		ClassifiedAt:    time.Now().UTC(),
	}
}
