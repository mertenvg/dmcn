package api

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/mertenvg/logr/v2"
	"google.golang.org/protobuf/proto"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/proto/dmcnpb"
)

// RelayRouter provides the ability to connect to peers and store envelopes on
// remote relays identified by their peer ID.
type RelayRouter interface {
	ConnectPeer(addr string) error
	StorePreSignedOnPeer(ctx context.Context, peerID string, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error)
}

// MessageHandler handles message send, list, and ack requests.
type MessageHandler struct {
	users          *store.UserStore
	sessions       *store.SessionStore
	envelopes      *store.EnvelopeStore
	storePreSigned func(ctx context.Context, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error)
	registryLookup func(ctx context.Context, address string) (*identity.IdentityRecord, error)
	relayRouter    RelayRouter
	log            logr.Logger
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(
	users *store.UserStore,
	sessions *store.SessionStore,
	envelopes *store.EnvelopeStore,
	storePreSigned func(ctx context.Context, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error),
	registryLookup func(ctx context.Context, address string) (*identity.IdentityRecord, error),
	relayRouter RelayRouter,
	log logr.Logger,
) *MessageHandler {
	return &MessageHandler{
		users:          users,
		sessions:       sessions,
		envelopes:      envelopes,
		storePreSigned: storePreSigned,
		registryLookup: registryLookup,
		relayRouter:    relayRouter,
		log:            log,
	}
}

// sendRequest is the JSON body for HandleSend.
type sendRequest struct {
	SenderAddress    string `json:"sender_address"`
	SenderSignature  string `json:"sender_signature"`
	Envelope         string `json:"envelope"`
	RecipientAddress string `json:"recipient_address"`
}

// HandleSend handles sending a message.
func (h *MessageHandler) HandleSend(w http.ResponseWriter, r *http.Request) {
	address := store.AddressFromContext(r.Context())
	if address == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Verify sender address matches session.
	if req.SenderAddress != address {
		writeError(w, http.StatusForbidden, "sender address does not match session")
		return
	}

	// Decode envelope protobuf from base64.
	envBytes, err := base64.StdEncoding.DecodeString(req.Envelope)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid envelope encoding")
		return
	}

	var pbEnv dmcnpb.EncryptedEnvelope
	if err := proto.Unmarshal(envBytes, &pbEnv); err != nil {
		writeError(w, http.StatusBadRequest, "invalid envelope protobuf: "+err.Error())
		return
	}

	env, err := message.EncryptedEnvelopeFromProto(&pbEnv)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid envelope: "+err.Error())
		return
	}

	// Decode signature from base64.
	sigBytes, err := base64.StdEncoding.DecodeString(req.SenderSignature)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid signature encoding")
		return
	}

	// Route STORE to recipient's relay hints if recipient_address is provided.
	var hash [32]byte
	if req.RecipientAddress != "" && h.registryLookup != nil && h.relayRouter != nil {
		recipientRec, lookupErr := h.registryLookup(r.Context(), req.RecipientAddress)
		if lookupErr != nil {
			writeError(w, http.StatusBadRequest, "recipient not found")
			return
		}
		if len(recipientRec.RelayHints) == 0 {
			writeError(w, http.StatusBadRequest, "recipient has no relay hints")
			return
		}

		var lastErr error
		for _, hint := range recipientRec.RelayHints {
			if connectErr := h.relayRouter.ConnectPeer(hint); connectErr != nil {
				lastErr = connectErr
				continue
			}
			hash, lastErr = h.relayRouter.StorePreSignedOnPeer(r.Context(), hint, req.SenderAddress, sigBytes, env)
			if lastErr == nil {
				break
			}
		}
		if lastErr != nil {
			h.log.Error("failed to store envelope on any relay", logr.M("error", lastErr.Error()))
			writeError(w, http.StatusBadGateway, "failed to deliver message to recipient relays")
			return
		}
	} else {
		// Fallback to default relay (legacy behavior).
		var storeErr error
		hash, storeErr = h.storePreSigned(r.Context(), req.SenderAddress, sigBytes, env)
		if storeErr != nil {
			h.log.Error("failed to store envelope", logr.M("error", storeErr.Error()))
			writeError(w, http.StatusInternalServerError, "failed to store message")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"envelope_hash": hex.EncodeToString(hash[:])})
}

// envelopeEntry is a single envelope in the list response.
type envelopeEntry struct {
	Hash string `json:"hash"`
	Data string `json:"data"`
}

// HandleList handles listing envelopes for the authenticated user.
func (h *MessageHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	address := store.AddressFromContext(r.Context())
	if address == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	dataList, hashes, err := h.envelopes.List(address)
	if err != nil {
		h.log.Error("failed to list envelopes", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}

	entries := make([]envelopeEntry, len(dataList))
	for i := range dataList {
		entries[i] = envelopeEntry{
			Hash: hex.EncodeToString(hashes[i][:]),
			Data: base64.StdEncoding.EncodeToString(dataList[i]),
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"envelopes": entries})
}

// ackRequest is the JSON body for HandleAck.
type ackRequest struct {
	EnvelopeHash string `json:"envelope_hash"`
}

// HandleAck handles acknowledging (deleting) an envelope.
func (h *MessageHandler) HandleAck(w http.ResponseWriter, r *http.Request) {
	address := store.AddressFromContext(r.Context())
	if address == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req ackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	hashBytes, err := hex.DecodeString(req.EnvelopeHash)
	if err != nil || len(hashBytes) != 32 {
		writeError(w, http.StatusBadRequest, "invalid envelope hash")
		return
	}

	var hash [32]byte
	copy(hash[:], hashBytes)

	if err := h.envelopes.Delete(address, hash); err != nil {
		h.log.Error("failed to delete envelope", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to acknowledge message")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
