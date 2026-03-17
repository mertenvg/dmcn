// Package api implements HTTP handlers for the DMCN web client backend.
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/core/identity"
)

// pendingChallenge holds a login challenge nonce and its expiry.
type pendingChallenge struct {
	nonce     []byte
	expiresAt time.Time
}

// AuthHandler handles registration, login, and logout requests.
type AuthHandler struct {
	users             *store.UserStore
	sessions          *store.SessionStore
	registryRegister  func(ctx context.Context, rec *identity.IdentityRecord) error
	log               logr.Logger
	pendingChallenges sync.Map // address -> pendingChallenge
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(
	users *store.UserStore,
	sessions *store.SessionStore,
	registryRegister func(ctx context.Context, rec *identity.IdentityRecord) error,
	log logr.Logger,
) *AuthHandler {
	return &AuthHandler{
		users:            users,
		sessions:         sessions,
		registryRegister: registryRegister,
		log:              log,
	}
}

// registerRequest is the JSON body for HandleRegister.
type registerRequest struct {
	Address          string                 `json:"address"`
	Ed25519Pub       string                 `json:"ed25519_pub"`
	X25519Pub        string                 `json:"x25519_pub"`
	EncryptedPayload store.EncryptedPayload `json:"encrypted_payload"`
	IdentityRecord   string                 `json:"identity_record"`
	SelfSignature    string                 `json:"self_signature"`
}

// HandleRegister handles user registration.
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if h.users.Exists(req.Address) {
		writeError(w, http.StatusConflict, "address already registered")
		return
	}

	// Decode Ed25519 public key from base64.
	ed25519PubBytes, err := base64.StdEncoding.DecodeString(req.Ed25519Pub)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid ed25519 public key encoding")
		return
	}

	// Decode identity record protobuf from base64.
	identityRecordBytes, err := base64.StdEncoding.DecodeString(req.IdentityRecord)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identity record encoding")
		return
	}

	rec, err := identity.IdentityRecordFromProtoBytes(identityRecordBytes)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identity record: "+err.Error())
		return
	}

	// Verify the self-signature.
	if err := rec.Verify(); err != nil {
		writeError(w, http.StatusBadRequest, "identity record signature verification failed")
		return
	}

	// Save user record.
	user := &store.UserRecord{
		Version:          1,
		Address:          req.Address,
		Ed25519Pub:       base64.StdEncoding.EncodeToString(ed25519PubBytes),
		X25519Pub:        req.X25519Pub,
		EncryptedPayload: req.EncryptedPayload,
	}
	if err := h.users.Save(user); err != nil {
		h.log.Error("failed to save user", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to save user")
		return
	}

	// Register identity in DHT.
	if err := h.registryRegister(r.Context(), rec); err != nil {
		h.log.Error("failed to register identity in DHT", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to register identity")
		return
	}

	// Create session.
	token, err := h.sessions.Create(req.Address)
	if err != nil {
		h.log.Error("failed to create session", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"session_token": token})
}

// loginRequest is the JSON body for HandleLogin.
type loginRequest struct {
	Address string `json:"address"`
}

// HandleLogin handles the first step of login: returns user info and a challenge nonce.
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.users.Load(req.Address)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// Generate 32-byte challenge nonce.
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		h.log.Error("failed to generate challenge nonce", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to generate challenge")
		return
	}

	// Store the pending challenge with 60s expiry.
	h.pendingChallenges.Store(req.Address, pendingChallenge{
		nonce:     nonce,
		expiresAt: time.Now().Add(60 * time.Second),
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"version":           user.Version,
		"ed25519_pub":       user.Ed25519Pub,
		"encrypted_payload": user.EncryptedPayload,
		"challenge_nonce":   base64.StdEncoding.EncodeToString(nonce),
	})
}

// loginVerifyRequest is the JSON body for HandleLoginVerify.
type loginVerifyRequest struct {
	Address            string `json:"address"`
	ChallengeSignature string `json:"challenge_signature"`
	ChallengeNonce     string `json:"challenge_nonce"`
}

// HandleLoginVerify handles the second step of login: verifies the signed challenge.
func (h *AuthHandler) HandleLoginVerify(w http.ResponseWriter, r *http.Request) {
	var req loginVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Load and validate pending challenge.
	val, ok := h.pendingChallenges.Load(req.Address)
	if !ok {
		writeError(w, http.StatusBadRequest, "no pending challenge for address")
		return
	}
	challenge := val.(pendingChallenge)

	if time.Now().After(challenge.expiresAt) {
		h.pendingChallenges.Delete(req.Address)
		writeError(w, http.StatusBadRequest, "challenge expired")
		return
	}

	// Load user to get Ed25519 public key.
	user, err := h.users.Load(req.Address)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	ed25519PubBytes, err := base64.StdEncoding.DecodeString(user.Ed25519Pub)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid stored public key")
		return
	}

	// Decode the signature.
	sigBytes, err := base64.StdEncoding.DecodeString(req.ChallengeSignature)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid signature encoding")
		return
	}

	// Verify the signature of the nonce.
	if err := crypto.Verify(ed25519PubBytes, challenge.nonce, sigBytes); err != nil {
		writeError(w, http.StatusUnauthorized, "challenge signature verification failed")
		return
	}

	// Delete pending challenge.
	h.pendingChallenges.Delete(req.Address)

	// Create session.
	token, err := h.sessions.Create(req.Address)
	if err != nil {
		h.log.Error("failed to create session", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"session_token": token})
}

// HandleLogout handles session logout.
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		writeError(w, http.StatusBadRequest, "missing bearer token")
		return
	}

	h.sessions.Delete(token)
	w.WriteHeader(http.StatusNoContent)
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
