// Package registry implements a DHT-based identity registry using libp2p
// Kademlia DHT. Identity records are stored keyed on SHA-256(address).
//
// See whitepaper Section 15.2.4 and PRD Section 5.1.
package registry

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/core/host"
	"google.golang.org/protobuf/proto"

	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/core/identity"
)

var (
	// ErrNotFound is returned when an identity is not in the registry.
	ErrNotFound = errors.New("registry: identity not found")
	// ErrInvalidRecord is returned when a stored identity record fails validation.
	ErrInvalidRecord = errors.New("registry: invalid identity record")
)

// dhtKeyPrefix is the namespace prefix for identity records in the DHT.
// Must match the namespace registered with NamespacedValidator.
const dhtKeyPrefix = "/dmcn/"

// Registry wraps a libp2p Kademlia DHT for identity record storage.
type Registry struct {
	dht  *dht.IpfsDHT
	host host.Host
}

// Option configures a Registry.
type Option func(*options)

type options struct {
	bootstrapPeers []string
	dhtMode        dht.ModeOpt
}

// WithBootstrapPeers sets the bootstrap peer addresses.
func WithBootstrapPeers(peers []string) Option {
	return func(o *options) {
		o.bootstrapPeers = peers
	}
}

// WithDHTMode sets the DHT mode (client or server).
func WithDHTMode(mode dht.ModeOpt) Option {
	return func(o *options) {
		o.dhtMode = mode
	}
}

// New creates a new Registry backed by a Kademlia DHT on the given libp2p host.
func New(ctx context.Context, h host.Host, opts ...Option) (*Registry, error) {
	cfg := &options{
		dhtMode: dht.ModeServer,
	}
	for _, o := range opts {
		o(cfg)
	}

	validator := &identityValidator{}

	d, err := dht.New(ctx, h,
		dht.Mode(cfg.dhtMode),
		dht.NamespacedValidator(dhtNamespace, validator),
		dht.ProtocolPrefix("/dmcn"),
	)
	if err != nil {
		return nil, fmt.Errorf("registry: create DHT: %w", err)
	}

	if err := d.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("registry: bootstrap DHT: %w", err)
	}

	return &Registry{dht: d, host: h}, nil
}

// Register stores a signed IdentityRecord in the DHT.
// The record's self-signature is validated before storage.
func (r *Registry) Register(ctx context.Context, rec *identity.IdentityRecord) error {
	if err := rec.Verify(); err != nil {
		return fmt.Errorf("registry: register: %w", err)
	}

	pb := rec.ToProto()
	data, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("registry: register: marshal: %w", err)
	}

	key := addressToKey(rec.Address)
	if err := r.dht.PutValue(ctx, key, data); err != nil {
		return fmt.Errorf("registry: register: put: %w", err)
	}

	return nil
}

// Lookup retrieves an IdentityRecord by exact address string.
// Returns ErrNotFound if absent. Validates the self-signature on retrieval.
func (r *Registry) Lookup(ctx context.Context, address string) (*identity.IdentityRecord, error) {
	key := addressToKey(address)

	data, err := r.dht.GetValue(ctx, key)
	if err != nil {
		// libp2p DHT returns routing.ErrNotFound when key is absent
		return nil, fmt.Errorf("%w: %s", ErrNotFound, address)
	}

	rec, err := unmarshalRecord(data)
	if err != nil {
		return nil, fmt.Errorf("registry: lookup: %w", err)
	}

	if err := rec.Verify(); err != nil {
		return nil, fmt.Errorf("%w: signature verification failed", ErrInvalidRecord)
	}

	return rec, nil
}

// Close shuts down the DHT.
func (r *Registry) Close() error {
	return r.dht.Close()
}

// addressToKey converts an address string to a DHT key.
func addressToKey(address string) string {
	hash := crypto.SHA256Hash([]byte(address))
	return dhtKeyPrefix + hex.EncodeToString(hash[:])
}

// unmarshalRecord deserializes a protobuf-encoded IdentityRecord.
func unmarshalRecord(data []byte) (*identity.IdentityRecord, error) {
	rec, err := identity.IdentityRecordFromProtoBytes(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRecord, err)
	}
	return rec, nil
}

// dhtNamespace is the DHT validator namespace matching our key prefix.
const dhtNamespace = "dmcn"

// identityValidator implements record.Validator for identity records in the DHT.
type identityValidator struct{}

var _ record.Validator = (*identityValidator)(nil)

// Validate checks that a DHT value is a valid signed IdentityRecord.
func (v *identityValidator) Validate(key string, value []byte) error {
	rec, err := unmarshalRecord(value)
	if err != nil {
		return err
	}
	if err := rec.Verify(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidRecord, err)
	}
	return nil
}

// Select chooses between two records. We prefer the newer record (by CreatedAt).
func (v *identityValidator) Select(key string, vals [][]byte) (int, error) {
	if len(vals) == 0 {
		return 0, errors.New("no values to select from")
	}
	bestIdx := 0
	bestRec, err := unmarshalRecord(vals[0])
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(vals); i++ {
		rec, err := unmarshalRecord(vals[i])
		if err != nil {
			continue
		}
		if rec.CreatedAt.After(bestRec.CreatedAt) {
			bestIdx = i
			bestRec = rec
		}
	}
	return bestIdx, nil
}
