// Package store provides file-based and in-memory storage for the DMCN web client.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// EncryptedPayload holds the encrypted key material for a user, with all
// fields base64-encoded.
type EncryptedPayload struct {
	Salt       string `json:"salt"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
	Tag        string `json:"tag"`
}

// UserRecord is the JSON structure persisted per user.
type UserRecord struct {
	Version          int              `json:"version"`
	Address          string           `json:"address"`
	Ed25519Pub       string           `json:"ed25519_pub"`
	X25519Pub        string           `json:"x25519_pub"`
	EncryptedPayload EncryptedPayload `json:"encrypted_payload"`
}

// UserStore manages file-based user storage. Each user is stored as a
// separate JSON file at <dataDir>/users/<address>.json.
type UserStore struct {
	dataDir string
	mu      sync.RWMutex
}

// NewUserStore creates a UserStore backed by the given data directory.
// It creates the <dataDir>/users/ directory if it does not exist.
func NewUserStore(dataDir string) (*UserStore, error) {
	dir := filepath.Join(dataDir, "users")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("store: create users dir: %w", err)
	}
	return &UserStore{dataDir: dataDir}, nil
}

// Save writes a UserRecord to disk as a JSON file with 0600 permissions.
func (s *UserStore) Save(user *UserRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("store: marshal user: %w", err)
	}

	path := s.userPath(user.Address)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("store: write user file: %w", err)
	}

	return nil
}

// Load reads a UserRecord from disk by address. Returns an error if the
// file does not exist or cannot be parsed.
func (s *UserStore) Load(address string) (*UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.userPath(address)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("store: read user file: %w", err)
	}

	var user UserRecord
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("store: unmarshal user: %w", err)
	}

	return &user, nil
}

// Exists returns true if a user file exists for the given address.
func (s *UserStore) Exists(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := os.Stat(s.userPath(address))
	return err == nil
}

// userPath constructs the file path for a user record.
func (s *UserStore) userPath(address string) string {
	return filepath.Join(s.dataDir, "users", address+".json")
}
