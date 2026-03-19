package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/api"
	"github.com/mertenvg/dmcn/internal/core/identity"
)

func TestHandleRelayHints(t *testing.T) {
	expectedHints := []string{"/ip4/1.2.3.4/tcp/7400/p2p/QmTest1", "/ip4/5.6.7.8/tcp/7400/p2p/QmTest2"}

	h := api.NewIdentityHandler(
		func(ctx context.Context, address string) (*identity.IdentityRecord, error) {
			return nil, nil
		},
		func() []string { return expectedHints },
		logr.With(logr.M("test", true)),
	)

	req := httptest.NewRequest("GET", "/api/v1/relay-hints", nil)
	rr := httptest.NewRecorder()
	h.HandleRelayHints(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	hints, ok := resp["relay_hints"].([]interface{})
	if !ok {
		t.Fatal("expected relay_hints array")
	}
	if len(hints) != 2 {
		t.Fatalf("expected 2 hints, got %d", len(hints))
	}
	if hints[0].(string) != expectedHints[0] {
		t.Errorf("hint[0] = %q, want %q", hints[0], expectedHints[0])
	}
}

func TestHandleRelayHints_Empty(t *testing.T) {
	h := api.NewIdentityHandler(
		func(ctx context.Context, address string) (*identity.IdentityRecord, error) {
			return nil, nil
		},
		func() []string { return nil },
		logr.With(logr.M("test", true)),
	)

	req := httptest.NewRequest("GET", "/api/v1/relay-hints", nil)
	rr := httptest.NewRecorder()
	h.HandleRelayHints(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	hints := resp["relay_hints"].([]interface{})
	if len(hints) != 0 {
		t.Fatalf("expected 0 hints, got %d", len(hints))
	}
}
