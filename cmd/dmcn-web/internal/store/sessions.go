package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrSessionNotFound is returned when a session token is not in the store.
	ErrSessionNotFound = errors.New("store: session not found")
	// ErrSessionExpired is returned when a session has exceeded its expiry duration.
	ErrSessionExpired = errors.New("store: session expired")
)

// session holds per-session state.
type session struct {
	Address   string
	CreatedAt time.Time
}

// SessionStore manages in-memory session tokens mapped to user addresses.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*session
	expiry   time.Duration
}

// NewSessionStore creates a SessionStore with the given expiry duration.
// If expiry is zero, a default of 24 hours is used.
func NewSessionStore(expiry time.Duration) *SessionStore {
	if expiry == 0 {
		expiry = 24 * time.Hour
	}
	return &SessionStore{
		sessions: make(map[string]*session),
		expiry:   expiry,
	}
}

// Create generates a new session token for the given address and returns it.
// The token is 32 random bytes, hex-encoded to a 64-character string.
func (s *SessionStore) Create(address string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("store: generate session token: %w", err)
	}

	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[token] = &session{
		Address:   address,
		CreatedAt: time.Now(),
	}
	s.mu.Unlock()

	return token, nil
}

// Validate checks that the token exists and has not expired. It returns the
// associated address on success.
func (s *SessionStore) Validate(token string) (string, error) {
	s.mu.RLock()
	sess, ok := s.sessions[token]
	s.mu.RUnlock()

	if !ok {
		return "", ErrSessionNotFound
	}

	if time.Since(sess.CreatedAt) > s.expiry {
		s.mu.Lock()
		delete(s.sessions, token)
		s.mu.Unlock()
		return "", ErrSessionExpired
	}

	return sess.Address, nil
}

// Delete removes a session by token.
func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}
