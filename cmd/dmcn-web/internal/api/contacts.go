package api

import (
	"encoding/json"
	"net/http"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
)

// ContactHandler handles contact-related payload update requests.
type ContactHandler struct {
	users    *store.UserStore
	sessions *store.SessionStore
	log      logr.Logger
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(
	users *store.UserStore,
	sessions *store.SessionStore,
	log logr.Logger,
) *ContactHandler {
	return &ContactHandler{
		users:    users,
		sessions: sessions,
		log:      log,
	}
}

// updatePayloadRequest is the JSON body for HandleUpdatePayload.
type updatePayloadRequest struct {
	EncryptedPayload store.EncryptedPayload `json:"encrypted_payload"`
}

// HandleUpdatePayload updates the encrypted payload for the authenticated user.
func (h *ContactHandler) HandleUpdatePayload(w http.ResponseWriter, r *http.Request) {
	address := store.AddressFromContext(r.Context())
	if address == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req updatePayloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.users.Load(address)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	user.Version++
	user.EncryptedPayload = req.EncryptedPayload

	if err := h.users.Save(user); err != nil {
		h.log.Error("failed to save user", logr.M("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to update payload")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"version": user.Version})
}
