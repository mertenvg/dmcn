package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
)

// --- Context ---

func TestAddressFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), store.ContextKeyAddress, "alice@example.com")
	if got := store.AddressFromContext(ctx); got != "alice@example.com" {
		t.Fatalf("expected alice@example.com, got %q", got)
	}
}

func TestAddressFromContext_Missing(t *testing.T) {
	if got := store.AddressFromContext(context.Background()); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestAddressFromContext_WrongKeyType(t *testing.T) {
	// Using a plain string key should not match the typed contextKey.
	ctx := context.WithValue(context.Background(), "address", "alice@example.com")
	if got := store.AddressFromContext(ctx); got != "" {
		t.Fatalf("expected empty string for plain string key, got %q", got)
	}
}

// --- UserStore ---

func TestUserStore_SaveAndLoad(t *testing.T) {
	us, err := store.NewUserStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	rec := &store.UserRecord{
		Version:    1,
		Address:    "alice@dmcn.me",
		Ed25519Pub: "AAAA",
		X25519Pub:  "BBBB",
	}
	if err := us.Save(rec); err != nil {
		t.Fatal(err)
	}

	loaded, err := us.Load("alice@dmcn.me")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Address != rec.Address || loaded.Version != rec.Version {
		t.Fatalf("loaded record mismatch: %+v", loaded)
	}
}

func TestUserStore_Exists(t *testing.T) {
	us, err := store.NewUserStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if us.Exists("nobody@dmcn.me") {
		t.Fatal("expected Exists to return false for non-existent user")
	}

	_ = us.Save(&store.UserRecord{Address: "bob@dmcn.me"})
	if !us.Exists("bob@dmcn.me") {
		t.Fatal("expected Exists to return true after save")
	}
}

func TestUserStore_LoadNonexistent(t *testing.T) {
	us, err := store.NewUserStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	_, err = us.Load("missing@dmcn.me")
	if err == nil {
		t.Fatal("expected error loading nonexistent user")
	}
}

// --- SessionStore ---

func TestSessionStore_CreateAndValidate(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	token, err := ss.Create("alice@dmcn.me")
	if err != nil {
		t.Fatal(err)
	}
	if len(token) != 64 {
		t.Fatalf("expected 64-char token, got %d chars", len(token))
	}

	addr, err := ss.Validate(token)
	if err != nil {
		t.Fatal(err)
	}
	if addr != "alice@dmcn.me" {
		t.Fatalf("expected alice@dmcn.me, got %q", addr)
	}
}

func TestSessionStore_ValidateExpired(t *testing.T) {
	ss := store.NewSessionStore(time.Millisecond)
	token, _ := ss.Create("alice@dmcn.me")
	time.Sleep(5 * time.Millisecond)

	_, err := ss.Validate(token)
	if err == nil {
		t.Fatal("expected error for expired session")
	}
}

func TestSessionStore_Delete(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	token, _ := ss.Create("alice@dmcn.me")
	ss.Delete(token)

	_, err := ss.Validate(token)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestSessionStore_ValidateUnknown(t *testing.T) {
	ss := store.NewSessionStore(time.Hour)
	_, err := ss.Validate("nonexistenttoken")
	if err == nil {
		t.Fatal("expected error for unknown token")
	}
}

// --- EnvelopeStore ---

func TestEnvelopeStore_StoreAndList(t *testing.T) {
	es, err := store.NewEnvelopeStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	hash := [32]byte{1, 2, 3}
	data := []byte("envelope-data")

	if err := es.Store("alice@dmcn.me", hash, data); err != nil {
		t.Fatal(err)
	}

	dataList, hashes, err := es.List("alice@dmcn.me")
	if err != nil {
		t.Fatal(err)
	}
	if len(dataList) != 1 {
		t.Fatalf("expected 1 envelope, got %d", len(dataList))
	}
	if string(dataList[0]) != "envelope-data" {
		t.Fatalf("unexpected data: %q", dataList[0])
	}
	if hashes[0] != hash {
		t.Fatalf("hash mismatch")
	}
}

func TestEnvelopeStore_Delete(t *testing.T) {
	es, err := store.NewEnvelopeStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	hash := [32]byte{4, 5, 6}
	_ = es.Store("bob@dmcn.me", hash, []byte("data"))
	if err := es.Delete("bob@dmcn.me", hash); err != nil {
		t.Fatal(err)
	}

	dataList, _, err := es.List("bob@dmcn.me")
	if err != nil {
		t.Fatal(err)
	}
	if len(dataList) != 0 {
		t.Fatalf("expected 0 envelopes after delete, got %d", len(dataList))
	}
}

func TestEnvelopeStore_ListEmpty(t *testing.T) {
	es, err := store.NewEnvelopeStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	dataList, hashes, err := es.List("nobody@dmcn.me")
	if err != nil {
		t.Fatal(err)
	}
	if len(dataList) != 0 || len(hashes) != 0 {
		t.Fatal("expected empty results for unknown address")
	}
}
