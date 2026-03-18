package api_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mertenvg/logr/v2"
	"google.golang.org/protobuf/proto"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/api"
	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/core/identity"
)

func newTestAuthHandler(t *testing.T) (*api.AuthHandler, *store.UserStore, *store.SessionStore) {
	t.Helper()
	us, err := store.NewUserStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ss := store.NewSessionStore(time.Hour)
	register := func(ctx context.Context, rec *identity.IdentityRecord) error {
		return nil // stub
	}
	h := api.NewAuthHandler(us, ss, register, logr.With(logr.M("test", true)))
	return h, us, ss
}

func createSignedIdentityRecord(t *testing.T, address string) (*identity.IdentityKeyPair, *identity.IdentityRecord, []byte) {
	t.Helper()
	kp, err := identity.GenerateIdentityKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	rec, err := identity.NewIdentityRecord(address, kp)
	if err != nil {
		t.Fatal(err)
	}
	if err := rec.Sign(kp); err != nil {
		t.Fatal(err)
	}
	pb := rec.ToProto()
	data, err := proto.MarshalOptions{Deterministic: true}.Marshal(pb)
	if err != nil {
		t.Fatal(err)
	}
	return kp, rec, data
}

func TestHandleRegister_Success(t *testing.T) {
	h, _, _ := newTestAuthHandler(t)
	kp, _, recBytes := createSignedIdentityRecord(t, "alice@dmcn.me")

	body := map[string]interface{}{
		"address":          "alice@dmcn.me",
		"ed25519_pub":      base64.StdEncoding.EncodeToString(kp.Ed25519Public),
		"x25519_pub":       base64.StdEncoding.EncodeToString(kp.X25519Public[:]),
		"encrypted_payload": map[string]string{"salt": "s", "nonce": "n", "ciphertext": "c", "tag": "t"},
		"identity_record":  base64.StdEncoding.EncodeToString(recBytes),
		"self_signature":   base64.StdEncoding.EncodeToString([]byte("unused")),
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/register", strings.NewReader(string(b)))
	rr := httptest.NewRecorder()
	h.HandleRegister(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["session_token"] == "" {
		t.Fatal("expected session_token in response")
	}
}

func TestHandleRegister_Duplicate(t *testing.T) {
	h, us, _ := newTestAuthHandler(t)

	_ = us.Save(&store.UserRecord{Address: "alice@dmcn.me"})

	body := `{"address":"alice@dmcn.me","ed25519_pub":"AA==","x25519_pub":"BB==","identity_record":"","self_signature":""}`
	req := httptest.NewRequest("POST", "/api/v1/register", strings.NewReader(body))
	rr := httptest.NewRecorder()
	h.HandleRegister(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rr.Code)
	}
}

func TestHandleRegister_MissingFields(t *testing.T) {
	h, _, _ := newTestAuthHandler(t)

	req := httptest.NewRequest("POST", "/api/v1/register", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	h.HandleRegister(rr, req)

	// Empty identity_record should fail decoding.
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHandleLogin_AndVerify(t *testing.T) {
	h, us, _ := newTestAuthHandler(t)
	kp, _, _ := createSignedIdentityRecord(t, "bob@dmcn.me")

	// Save user.
	_ = us.Save(&store.UserRecord{
		Version:    1,
		Address:    "bob@dmcn.me",
		Ed25519Pub: base64.StdEncoding.EncodeToString(kp.Ed25519Public),
		X25519Pub:  base64.StdEncoding.EncodeToString(kp.X25519Public[:]),
	})

	// Step 1: Login.
	loginBody, _ := json.Marshal(map[string]string{"address": "bob@dmcn.me"})
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginBody)))
	rr := httptest.NewRecorder()
	h.HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var loginResp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&loginResp)

	nonceB64, ok := loginResp["challenge_nonce"].(string)
	if !ok || nonceB64 == "" {
		t.Fatal("expected challenge_nonce in login response")
	}

	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the nonce.
	sig, err := crypto.Sign(kp.Ed25519Private, nonce)
	if err != nil {
		t.Fatal(err)
	}

	// Step 2: Verify.
	verifyBody, _ := json.Marshal(map[string]string{
		"address":             "bob@dmcn.me",
		"challenge_signature": base64.StdEncoding.EncodeToString(sig),
		"challenge_nonce":     nonceB64,
	})
	req2 := httptest.NewRequest("POST", "/api/v1/login/verify", strings.NewReader(string(verifyBody)))
	rr2 := httptest.NewRecorder()
	h.HandleLoginVerify(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("verify: expected 200, got %d: %s", rr2.Code, rr2.Body.String())
	}

	var verifyResp map[string]string
	json.NewDecoder(rr2.Body).Decode(&verifyResp)
	if verifyResp["session_token"] == "" {
		t.Fatal("expected session_token in verify response")
	}
}

func TestHandleLoginVerify_BadSignature(t *testing.T) {
	h, us, _ := newTestAuthHandler(t)
	kp, _, _ := createSignedIdentityRecord(t, "carol@dmcn.me")

	_ = us.Save(&store.UserRecord{
		Version:    1,
		Address:    "carol@dmcn.me",
		Ed25519Pub: base64.StdEncoding.EncodeToString(kp.Ed25519Public),
	})

	// Login to get challenge.
	loginBody, _ := json.Marshal(map[string]string{"address": "carol@dmcn.me"})
	req := httptest.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginBody)))
	rr := httptest.NewRecorder()
	h.HandleLogin(rr, req)

	var loginResp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&loginResp)
	nonceB64 := loginResp["challenge_nonce"].(string)

	// Send bad signature.
	verifyBody, _ := json.Marshal(map[string]string{
		"address":             "carol@dmcn.me",
		"challenge_signature": base64.StdEncoding.EncodeToString(make([]byte, 64)),
		"challenge_nonce":     nonceB64,
	})
	req2 := httptest.NewRequest("POST", "/api/v1/login/verify", strings.NewReader(string(verifyBody)))
	rr2 := httptest.NewRecorder()
	h.HandleLoginVerify(rr2, req2)

	if rr2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rr2.Code, rr2.Body.String())
	}
}
