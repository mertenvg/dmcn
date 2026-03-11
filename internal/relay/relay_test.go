package relay

import (
	"testing"
	"time"

	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
)

func TestMessageStoreBasic(t *testing.T) {
	store := NewMessageStore()

	kp, _ := identity.GenerateIdentityKeyPair()
	msg, _ := message.NewPlaintextMessage(
		"alice@localhost", "bob@localhost", "Test", "Hello",
		kp.Ed25519Public,
	)
	sm := &message.SignedMessage{Plaintext: *msg}
	sm.Sign(kp.Ed25519Private)

	recipKP, _ := identity.GenerateIdentityKeyPair()
	env, _ := message.Encrypt(sm, []message.RecipientInfo{{
		DeviceID: recipKP.DeviceID, X25519Pub: recipKP.X25519Public,
	}})

	hash := [32]byte{1, 2, 3}
	store.Store("bob@localhost", env, hash)

	// Count
	if c := store.Count(); c != 1 {
		t.Errorf("count = %d, want 1", c)
	}

	// Fetch
	envs, hashes := store.Fetch("bob@localhost")
	if len(envs) != 1 {
		t.Fatalf("fetched %d, want 1", len(envs))
	}
	if hashes[0] != hash {
		t.Error("hash mismatch")
	}

	// Delivery status
	status, err := store.DeliveryStatusOf(hash)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if status != Pending {
		t.Errorf("status = %d, want Pending", status)
	}

	// Ack
	if err := store.Ack(hash); err != nil {
		t.Fatalf("ack: %v", err)
	}

	status, _ = store.DeliveryStatusOf(hash)
	if status != Delivered {
		t.Errorf("status after ack = %d, want Delivered", status)
	}

	// Fetch should return empty after ack (only pending)
	envs, _ = store.Fetch("bob@localhost")
	if len(envs) != 0 {
		t.Errorf("fetched %d after ack, want 0", len(envs))
	}
}

func TestMessageStoreNotFound(t *testing.T) {
	store := NewMessageStore()

	// Ack non-existent
	err := store.Ack([32]byte{99})
	if err != ErrEnvelopeNotFound {
		t.Errorf("ack non-existent: got %v, want ErrEnvelopeNotFound", err)
	}

	// Status non-existent
	_, err = store.DeliveryStatusOf([32]byte{99})
	if err != ErrEnvelopeNotFound {
		t.Errorf("status non-existent: got %v, want ErrEnvelopeNotFound", err)
	}
}

func TestMessageStoreMultiple(t *testing.T) {
	store := NewMessageStore()

	kp, _ := identity.GenerateIdentityKeyPair()
	recipKP, _ := identity.GenerateIdentityKeyPair()

	for i := 0; i < 5; i++ {
		msg, _ := message.NewPlaintextMessage(
			"alice@localhost", "bob@localhost", "", "msg",
			kp.Ed25519Public,
		)
		sm := &message.SignedMessage{Plaintext: *msg}
		sm.Sign(kp.Ed25519Private)
		env, _ := message.Encrypt(sm, []message.RecipientInfo{{
			DeviceID: recipKP.DeviceID, X25519Pub: recipKP.X25519Public,
		}})

		hash := [32]byte{byte(i)}
		store.Store("bob@localhost", env, hash)
	}

	if c := store.Count(); c != 5 {
		t.Errorf("count = %d, want 5", c)
	}

	envs, _ := store.Fetch("bob@localhost")
	if len(envs) != 5 {
		t.Errorf("fetched %d, want 5", len(envs))
	}

	// Fetch for different address returns empty
	envs, _ = store.Fetch("charlie@localhost")
	if len(envs) != 0 {
		t.Errorf("fetched %d for charlie, want 0", len(envs))
	}
}

func TestRateLimiterBasic(t *testing.T) {
	rl := NewRateLimiter(3)

	// First 3 should be allowed
	for i := 0; i < 3; i++ {
		if !rl.Allow("alice@localhost") {
			t.Errorf("attempt %d should be allowed", i)
		}
	}

	// 4th should be denied
	if rl.Allow("alice@localhost") {
		t.Error("4th attempt should be denied")
	}

	// Different sender should be allowed
	if !rl.Allow("bob@localhost") {
		t.Error("bob should be allowed (separate limit)")
	}
}

func TestRateLimiterExpiry(t *testing.T) {
	rl := NewRateLimiter(2)

	// Manually set timestamps in the past
	past := time.Now().Add(-2 * time.Hour)
	rl.timestamps["alice@localhost"] = []time.Time{past, past}

	// Should be allowed since old timestamps are pruned
	if !rl.Allow("alice@localhost") {
		t.Error("should be allowed after expiry")
	}
}

func TestRateLimiterWindow(t *testing.T) {
	rl := NewRateLimiter(2)

	// Override nowFunc to control time
	now := time.Now()
	rl.nowFunc = func() time.Time { return now }

	rl.Allow("alice@localhost")
	rl.Allow("alice@localhost")

	// At current time, 3rd should be denied
	if rl.Allow("alice@localhost") {
		t.Error("3rd should be denied within window")
	}

	// Move time forward past the window
	rl.nowFunc = func() time.Time { return now.Add(61 * time.Minute) }

	// Should be allowed now
	if !rl.Allow("alice@localhost") {
		t.Error("should be allowed after window expires")
	}
}
