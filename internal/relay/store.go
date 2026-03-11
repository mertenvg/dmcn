// Package relay implements the DMCN relay node protocol for message
// storage and delivery. See whitepaper Section 15.4.2 and PRD Section 5.2.
package relay

import (
	"errors"
	"sync"

	"github.com/mertenvg/dmcn/internal/core/message"
)

// DeliveryStatus represents the state of a stored envelope.
type DeliveryStatus int

const (
	// Pending means the envelope has not been fetched yet.
	Pending DeliveryStatus = iota
	// Delivered means the recipient has acknowledged receipt.
	Delivered
)

var (
	// ErrEnvelopeNotFound is returned when an envelope hash is not in the store.
	ErrEnvelopeNotFound = errors.New("relay: envelope not found")
)

// storedEnvelope holds an envelope along with its delivery metadata.
type storedEnvelope struct {
	Envelope *message.EncryptedEnvelope
	Hash     [32]byte
	Status   DeliveryStatus
}

// MessageStore is an in-memory store for encrypted envelopes, indexed by
// recipient address. Sufficient for PoC; persistent storage is post-PoC.
type MessageStore struct {
	mu        sync.RWMutex
	byAddr    map[string][]*storedEnvelope    // recipient address → envelopes
	byHash    map[[32]byte]*storedEnvelope    // envelope hash → envelope
}

// NewMessageStore creates an empty in-memory message store.
func NewMessageStore() *MessageStore {
	return &MessageStore{
		byAddr: make(map[string][]*storedEnvelope),
		byHash: make(map[[32]byte]*storedEnvelope),
	}
}

// Store adds an encrypted envelope to the store, indexed by recipient address.
func (s *MessageStore) Store(recipientAddr string, env *message.EncryptedEnvelope, hash [32]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	se := &storedEnvelope{
		Envelope: env,
		Hash:     hash,
		Status:   Pending,
	}

	s.byAddr[recipientAddr] = append(s.byAddr[recipientAddr], se)
	s.byHash[hash] = se
}

// Fetch returns all pending envelopes for a recipient address along with
// their hashes. Does not remove them from the store.
func (s *MessageStore) Fetch(recipientAddr string) ([]*message.EncryptedEnvelope, [][32]byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored := s.byAddr[recipientAddr]
	var envs []*message.EncryptedEnvelope
	var hashes [][32]byte
	for _, se := range stored {
		if se.Status == Pending {
			envs = append(envs, se.Envelope)
			hashes = append(hashes, se.Hash)
		}
	}
	return envs, hashes
}

// Ack marks an envelope as delivered by its hash.
func (s *MessageStore) Ack(hash [32]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	se, ok := s.byHash[hash]
	if !ok {
		return ErrEnvelopeNotFound
	}
	se.Status = Delivered
	return nil
}

// DeliveryStatusOf returns the delivery status of an envelope by its hash.
func (s *MessageStore) DeliveryStatusOf(hash [32]byte) (DeliveryStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	se, ok := s.byHash[hash]
	if !ok {
		return 0, ErrEnvelopeNotFound
	}
	return se.Status, nil
}

// Count returns the total number of stored envelopes.
func (s *MessageStore) Count() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return uint32(len(s.byHash))
}
