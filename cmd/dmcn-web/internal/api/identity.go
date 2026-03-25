package api

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/internal/core/identity"
)

// IdentityHandler handles DHT identity lookup requests.
type IdentityHandler struct {
	lookup     func(ctx context.Context, address string) (*identity.IdentityRecord, error)
	relayHints func() []string
	log        logr.Logger
}

// NewIdentityHandler creates a new IdentityHandler.
func NewIdentityHandler(
	lookup func(ctx context.Context, address string) (*identity.IdentityRecord, error),
	relayHints func() []string,
	log logr.Logger,
) *IdentityHandler {
	return &IdentityHandler{
		lookup:     lookup,
		relayHints: relayHints,
		log:        log,
	}
}

// HandleRelayHints returns the relay hints for this web backend's node.
func (h *IdentityHandler) HandleRelayHints(w http.ResponseWriter, r *http.Request) {
	hints := h.relayHints()
	if hints == nil {
		hints = []string{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"relay_hints": hints})
}

// HandleLookup handles an identity lookup by address query parameter.
func (h *IdentityHandler) HandleLookup(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		writeError(w, http.StatusBadRequest, "missing address query parameter")
		return
	}

	rec, err := h.lookup(r.Context(), address)
	if err != nil {
		h.log.Error("identity lookup failed", logr.M("error", err.Error()), logr.M("address", address))
		writeError(w, http.StatusNotFound, "identity not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address":           rec.Address,
		"ed25519_pub":       base64.StdEncoding.EncodeToString(rec.Ed25519Public),
		"x25519_pub":        base64.StdEncoding.EncodeToString(rec.X25519Public[:]),
		"fingerprint":       rec.Fingerprint(),
		"verification_tier": int(rec.VerificationTier),
	})
}
