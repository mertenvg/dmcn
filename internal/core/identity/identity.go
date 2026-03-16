// Package identity implements the DMCN identity layer data structures
// and operations defined in whitepaper Section 15.2.
//
// An identity consists of an Ed25519 signing key pair and an X25519
// key exchange pair, bound together in a self-certifying IdentityRecord.
package identity

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/proto/dmcnpb"
	"google.golang.org/protobuf/proto"
)

// protoMarshal is the protobuf marshaling function, overridable for testing.
var protoMarshal = func(m proto.Message) ([]byte, error) {
	return proto.MarshalOptions{Deterministic: true}.Marshal(m)
}

var (
	// ErrInvalidAddress is returned when an address is malformed.
	ErrInvalidAddress = errors.New("identity: invalid address format")
	// ErrInvalidSignature is returned when a self-signature is invalid.
	ErrInvalidSignature = errors.New("identity: invalid self-signature")
	// ErrExpired is returned when an identity record has expired.
	ErrExpired = errors.New("identity: record expired")
)

// VerificationTier represents the level of identity verification.
// See whitepaper Section 15.2.2.
type VerificationTier int

const (
	TierUnverified    VerificationTier = iota // No verification
	TierProviderHosted                        // Provider-hosted key
	TierDomainDNS                             // Domain DNS verification
	TierDANE                                  // DANE (TLSA) verification
)

// IdentityKeyPair holds both the Ed25519 signing pair and the X25519 key
// exchange pair for a single identity, generated together at account creation.
//
// See whitepaper Section 15.2.1.
type IdentityKeyPair struct {
	Ed25519Public  ed25519.PublicKey
	Ed25519Private ed25519.PrivateKey
	X25519Public   [32]byte
	X25519Private  [32]byte
	CreatedAt      time.Time
	DeviceID       [16]byte
}

// GenerateIdentityKeyPair generates both key pairs in a single call.
// Private key material is never logged.
//
// See whitepaper Section 15.2.1.
func GenerateIdentityKeyPair() (*IdentityKeyPair, error) {
	edPub, edPriv, err := crypto.GenerateEd25519KeyPair()
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	xPub, xPriv, err := crypto.GenerateX25519KeyPair()
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	deviceID, err := crypto.RandomUUID()
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	return &IdentityKeyPair{
		Ed25519Public:  edPub,
		Ed25519Private: edPriv,
		X25519Public:   xPub,
		X25519Private:  xPriv,
		CreatedAt:      time.Now().UTC(),
		DeviceID:       deviceID,
	}, nil
}

// IdentityRecord maps a human-readable address to a key pair.
// It is self-certifying: the SelfSignature field covers all other fields
// and is produced by the identity's own Ed25519 private key.
//
// See whitepaper Section 15.2.2.
type IdentityRecord struct {
	Version          uint32
	Address          string // local@domain
	Ed25519Public    ed25519.PublicKey
	X25519Public     [32]byte
	CreatedAt        time.Time
	ExpiresAt        time.Time // zero = no expiry
	RelayHints       []string
	VerificationTier VerificationTier
	BridgeCapability bool
	SelfSignature    [64]byte
}

// NewIdentityRecord creates a new unsigned IdentityRecord from a key pair
// and address.
func NewIdentityRecord(address string, kp *IdentityKeyPair) (*IdentityRecord, error) {
	if err := validateAddress(address); err != nil {
		return nil, err
	}

	return &IdentityRecord{
		Version:          1,
		Address:          address,
		Ed25519Public:    kp.Ed25519Public,
		X25519Public:     kp.X25519Public,
		CreatedAt:        kp.CreatedAt,
		VerificationTier: TierUnverified,
	}, nil
}

// Sign computes and sets the SelfSignature. The signed byte sequence is
// the canonical protobuf serialisation of all fields except SelfSignature.
//
// See whitepaper Section 15.2.2.
func (r *IdentityRecord) Sign(kp *IdentityKeyPair) error {
	data, err := r.signableBytes()
	if err != nil {
		return fmt.Errorf("identity: sign: %w", err)
	}

	sig, err := crypto.Sign(kp.Ed25519Private, data)
	if err != nil {
		return fmt.Errorf("identity: sign: %w", err)
	}

	copy(r.SelfSignature[:], sig)
	return nil
}

// Verify validates the SelfSignature against the record's Ed25519 public key.
// Returns nil if valid, ErrInvalidSignature if not.
//
// See whitepaper Section 15.2.2.
func (r *IdentityRecord) Verify() error {
	data, err := r.signableBytes()
	if err != nil {
		return fmt.Errorf("identity: verify: %w", err)
	}

	if err := crypto.Verify(r.Ed25519Public, data, r.SelfSignature[:]); err != nil {
		return ErrInvalidSignature
	}
	return nil
}

// Fingerprint returns the first 20 bytes of SHA-256(Ed25519Public || X25519Public),
// encoded as a 40-character uppercase hex string.
//
// Used for out-of-band identity verification.
// See whitepaper Section 15.2.1.
func (r *IdentityRecord) Fingerprint() string {
	data := make([]byte, 0, len(r.Ed25519Public)+len(r.X25519Public))
	data = append(data, r.Ed25519Public...)
	data = append(data, r.X25519Public[:]...)
	hash := crypto.SHA256Hash(data)
	return strings.ToUpper(hex.EncodeToString(hash[:20]))
}

// signableBytes returns the canonical protobuf serialisation of the record
// with all fields except SelfSignature. This is the byte sequence over
// which the signature is computed.
func (r *IdentityRecord) signableBytes() ([]byte, error) {
	pb := &dmcnpb.IdentityRecord{
		Version:          r.Version,
		Address:          r.Address,
		Ed25519PublicKey: r.Ed25519Public,
		X25519PublicKey:  r.X25519Public[:],
		CreatedAt:        r.CreatedAt.Unix(),
		ExpiresAt:        r.ExpiresAt.Unix(),
		RelayHints:       r.RelayHints,
		VerificationTier: dmcnpb.VerificationTier(r.VerificationTier),
		BridgeCapability: r.BridgeCapability,
		// SelfSignature intentionally omitted — this is what we sign over
	}

	data, err := protoMarshal(pb)
	if err != nil {
		return nil, fmt.Errorf("protobuf marshal: %w", err)
	}
	return data, nil
}

// ToProto converts the IdentityRecord to its protobuf representation.
func (r *IdentityRecord) ToProto() *dmcnpb.IdentityRecord {
	return &dmcnpb.IdentityRecord{
		Version:          r.Version,
		Address:          r.Address,
		Ed25519PublicKey: r.Ed25519Public,
		X25519PublicKey:  r.X25519Public[:],
		CreatedAt:        r.CreatedAt.Unix(),
		ExpiresAt:        r.ExpiresAt.Unix(),
		RelayHints:       r.RelayHints,
		VerificationTier: dmcnpb.VerificationTier(r.VerificationTier),
		BridgeCapability: r.BridgeCapability,
		SelfSignature:    r.SelfSignature[:],
	}
}

// IdentityRecordFromProto creates an IdentityRecord from its protobuf
// representation.
func IdentityRecordFromProto(pb *dmcnpb.IdentityRecord) (*IdentityRecord, error) {
	if pb == nil {
		return nil, errors.New("identity: nil protobuf record")
	}

	var x25519Pub [32]byte
	copy(x25519Pub[:], pb.X25519PublicKey)

	var selfSig [64]byte
	copy(selfSig[:], pb.SelfSignature)

	var expiresAt time.Time
	if pb.ExpiresAt != 0 {
		expiresAt = time.Unix(pb.ExpiresAt, 0).UTC()
	}

	return &IdentityRecord{
		Version:          pb.Version,
		Address:          pb.Address,
		Ed25519Public:    pb.Ed25519PublicKey,
		X25519Public:     x25519Pub,
		CreatedAt:        time.Unix(pb.CreatedAt, 0).UTC(),
		ExpiresAt:        expiresAt,
		RelayHints:       pb.RelayHints,
		VerificationTier: VerificationTier(pb.VerificationTier),
		BridgeCapability: pb.BridgeCapability,
		SelfSignature:    selfSig,
	}, nil
}

// IdentityRecordFromProtoBytes deserializes an IdentityRecord from raw
// protobuf bytes. This is a convenience wrapper used by the registry package.
func IdentityRecordFromProtoBytes(data []byte) (*IdentityRecord, error) {
	pb := &dmcnpb.IdentityRecord{}
	if err := proto.Unmarshal(data, pb); err != nil {
		return nil, fmt.Errorf("identity: unmarshal: %w", err)
	}
	return IdentityRecordFromProto(pb)
}

// validateAddress performs basic validation on an address string.
// Addresses must be in local@domain format.
func validateAddress(address string) error {
	parts := strings.SplitN(address, "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("%w: %q", ErrInvalidAddress, address)
	}
	return nil
}
