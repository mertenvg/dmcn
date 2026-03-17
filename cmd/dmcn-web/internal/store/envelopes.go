package store

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// EnvelopeStore manages file-based envelope persistence. Envelopes are stored
// as protobuf binary files at <dataDir>/envelopes/<address>/<hex-hash>.bin.
type EnvelopeStore struct {
	dataDir string
	mu      sync.RWMutex
}

// NewEnvelopeStore creates an EnvelopeStore backed by the given data directory.
// It creates the <dataDir>/envelopes/ directory if it does not exist.
func NewEnvelopeStore(dataDir string) (*EnvelopeStore, error) {
	dir := filepath.Join(dataDir, "envelopes")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("store: create envelopes dir: %w", err)
	}
	return &EnvelopeStore{dataDir: dataDir}, nil
}

// Store writes envelope data to disk for the given address and hash.
// It creates the per-address subdirectory if needed.
func (s *EnvelopeStore) Store(address string, hash [32]byte, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Join(s.dataDir, "envelopes", address)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("store: create envelope dir: %w", err)
	}

	filename := hex.EncodeToString(hash[:]) + ".bin"
	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("store: write envelope file: %w", err)
	}

	return nil
}

// List returns all stored envelopes for the given address. It returns the
// raw data and corresponding hashes for each envelope file.
func (s *EnvelopeStore) List(address string) ([][]byte, [][32]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dir := filepath.Join(s.dataDir, "envelopes", address)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("store: read envelope dir: %w", err)
	}

	var dataList [][]byte
	var hashes [][32]byte

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".bin") {
			continue
		}

		hashHex := strings.TrimSuffix(entry.Name(), ".bin")
		hashBytes, err := hex.DecodeString(hashHex)
		if err != nil || len(hashBytes) != 32 {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("store: read envelope file: %w", err)
		}

		var hash [32]byte
		copy(hash[:], hashBytes)

		dataList = append(dataList, data)
		hashes = append(hashes, hash)
	}

	return dataList, hashes, nil
}

// Delete removes an envelope file for the given address and hash.
func (s *EnvelopeStore) Delete(address string, hash [32]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filename := hex.EncodeToString(hash[:]) + ".bin"
	path := filepath.Join(s.dataDir, "envelopes", address, filename)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("store: delete envelope file: %w", err)
	}

	return nil
}
