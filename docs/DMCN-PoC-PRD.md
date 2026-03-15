---
title: "DMCN — Proof of Concept Product Requirements Document"
version: "1.0"
status: "Draft — For Development Use"
date: "March 2026"
author: "Merten van Gerven"
language: "Go"
whitepaper: "DMCN Whitepaper v0.2"
confidentiality: "CONFIDENTIAL"
---

# DMCN — Proof of Concept PRD

**Decentralized Mesh Communication Network**
Version 1.0 · March 2026 · *CONFIDENTIAL — Draft for Development Use*

---

## Table of Contents

- [1. Purpose and Scope](#1-purpose-and-scope)
- [2. Background](#2-background)
- [3. Repository Structure](#3-repository-structure)
- [4. Milestone 1 — Cryptographic Core](#4-milestone-1--cryptographic-core)
- [5. Milestone 2 — Node and Registry](#5-milestone-2--node-and-registry)
- [6. Milestone 3 — Bridge Node](#6-milestone-3--bridge-node)
- [7. Non-Functional Requirements](#7-non-functional-requirements)
- [8. Milestone Summary](#8-milestone-summary)
- [9. Explicitly Deferred (Post-PoC)](#9-explicitly-deferred-post-poc)
- [10. Reference](#10-reference)

---

## 1. Purpose and Scope

This document defines the requirements for the initial proof-of-concept (PoC) implementation of the Decentralized Mesh Communication Network (DMCN). It is intended as a direct input to Claude Code and serves as the authoritative specification for what must be built, in what order, and to what standard, during the bootstrapping phase of the project.

The PoC has two sequential goals:

1. Validate the core cryptographic and protocol design by implementing the identity and message format layers in isolation, before any networking is introduced.
2. Demonstrate end-to-end message delivery between two nodes, including identity registration, lookup, message signing, encryption, routing, and decryption — sufficient to support technical review and stakeholder demonstration.

> **Out of scope for PoC:** Onion routing, shared reputation feeds, social key recovery, the full trust management stack, mobile clients, and production hardening. These are addressed in subsequent milestones.

---

## 2. Background

DMCN is a next-generation messaging protocol designed to replace SMTP by enforcing cryptographic sender identity at the protocol level. Spam and email fraud are treated as identity problems, not filtering problems. Every participant holds a public/private key pair; every message must carry a valid cryptographic signature from a registered identity; relay nodes reject unsigned messages at the network boundary.

The full protocol design is specified in the DMCN Whitepaper v0.2. This PRD covers only the PoC scope. All section references below (e.g. Section 15.2) refer to the whitepaper.

The implementation language is Go. This choice is appropriate for all server-side services given Go's concurrency model, mature cryptographic standard library, strong Protocol Buffers support, and suitability for I/O-bound distributed systems work.

---

## 3. Repository Structure

The repository should be structured as a Go module with clearly separated packages from the outset. Future services import from the core package; the core package has no dependency on any service.

```
dmcn/
  cmd/
    dmcn-node/        # combined dev node binary (Milestone 2)
    dmcn-bridge/      # bridge node binary (Milestone 3)
  internal/
    core/             # cryptographic primitives and data structures
      identity/       # key generation, identity records, registry ops
      message/        # plaintext, signed, encrypted envelope types
      crypto/         # low-level crypto wrappers (Ed25519, X25519, AES-GCM, HKDF)
    registry/         # DHT node and lookup client
    relay/            # relay node message handling
    bridge/           # SMTP-DMCN bridge protocol
  proto/              # Protocol Buffer definitions (.proto files)
  docs/               # architecture notes, protocol decisions
  go.mod
```

> **Convention:** All packages under `internal/` are unexported. Public API surface is minimal and intentional. Test files live alongside their packages.

---

## 4. Milestone 1 — Cryptographic Core

The cryptographic core is the foundation of every other component. It must be correct, well-tested, and reviewed before anything is built on top of it. This milestone produces no binary — it is a Go package with comprehensive test coverage.

### 4.1 Cryptographic Primitives (`internal/core/crypto`)

Implement thin, well-documented wrappers around Go standard library and `golang.org/x/crypto` primitives. Do not implement any cryptographic algorithm from scratch.

| Operation | Implementation |
|---|---|
| Ed25519 key generation | `crypto/ed25519` — standard library |
| Ed25519 sign | `ed25519.Sign` |
| Ed25519 verify | `ed25519.Verify` |
| X25519 key generation | `golang.org/x/crypto/curve25519` |
| X25519 key exchange | `curve25519.X25519` |
| HKDF-SHA256 | `golang.org/x/crypto/hkdf` + `crypto/sha256` |
| AES-256-GCM encrypt | `crypto/aes` + `crypto/cipher` (GCM mode) |
| AES-256-GCM decrypt | `crypto/aes` + `crypto/cipher` (GCM mode) |
| SHA-256 hash | `crypto/sha256` — standard library |
| Secure random bytes | `crypto/rand.Read` |

Each wrapper function must carry a doc comment that cites the relevant whitepaper section and states the exact input/output contract. All wrappers return explicit errors; no panics on invalid input.

### 4.2 Identity Layer (`internal/core/identity`)

Implement the data structures and operations defined in whitepaper Section 15.2.

#### 4.2.1 Key Pair

The `IdentityKeyPair` struct holds both the Ed25519 signing pair and the X25519 key exchange pair for a single identity, generated together at account creation.

```go
type IdentityKeyPair struct {
    Ed25519Public   ed25519.PublicKey    // 32 bytes — signing
    Ed25519Private  ed25519.PrivateKey   // 64 bytes — signing
    X25519Public    [32]byte             // 32 bytes — key exchange
    X25519Private   [32]byte             // 32 bytes — key exchange
    CreatedAt       time.Time
    DeviceID        [16]byte             // random UUID
}
```

`GenerateIdentityKeyPair()` generates both key pairs in a single call. Private key material is never logged.

#### 4.2.2 Identity Record

The `IdentityRecord` struct maps a human-readable address to a key pair. It is self-certifying: the `SelfSignature` field covers all other fields and is produced by the identity's own Ed25519 private key.

```go
type VerificationTier int

const (
    TierUnverified    VerificationTier = iota
    TierProviderHosted
    TierDomainDNS
    TierDANE
)

type IdentityRecord struct {
    Version          uint32
    Address          string           // local@domain
    Ed25519Public    ed25519.PublicKey
    X25519Public     [32]byte
    CreatedAt        time.Time
    ExpiresAt        time.Time        // zero = no expiry
    RelayHints       []string
    VerificationTier VerificationTier
    SelfSignature    [64]byte         // Ed25519 sig over all preceding fields
}
```

`Sign(keypair)` computes and sets the `SelfSignature`. `Verify()` validates the `SelfSignature`. The signed byte sequence is the canonical protobuf serialisation of all fields except `SelfSignature`. Both methods must be unit-tested with known vectors.

#### 4.2.3 Fingerprint

`Fingerprint()` returns the first 20 bytes of `SHA-256(Ed25519Public || X25519Public)`, encoded as a 40-character uppercase hex string. Used for out-of-band identity verification.

### 4.3 Message Format (`internal/core/message`)

Implement the three-layer message structure defined in whitepaper Section 15.3: `PlaintextMessage`, `SignedMessage`, and `EncryptedEnvelope`.

#### 4.3.1 PlaintextMessage

Represents a composed message before signing or encryption.

```go
type PlaintextMessage struct {
    Version          uint32
    MessageID        [16]byte         // random UUID — required; see Section 15.3.5
    ThreadID         [16]byte         // UUID linking conversation thread
    SenderAddress    string
    SenderPublicKey  ed25519.PublicKey
    RecipientAddress string
    SentAt           time.Time
    Subject          string
    Body             MessageBody
    Attachments      []AttachmentRecord
    ReplyToID        [16]byte         // zero = not a reply
}

type MessageBody struct {
    ContentType string               // e.g. "text/plain"
    Content     []byte
}
```

#### 4.3.2 SignedMessage

The sender signs the canonical serialisation of `PlaintextMessage` with their Ed25519 private key.

```go
type SignedMessage struct {
    Plaintext       PlaintextMessage
    SenderSignature [64]byte          // Ed25519 sig over Plaintext
}
```

`Sign(senderPrivKey)` computes and sets `SenderSignature`. `Verify(senderPubKey)` validates it. A `SignedMessage` with an invalid signature must never be displayed to a user — this invariant must be enforced at the API boundary, not left to callers.

#### 4.3.3 EncryptedEnvelope and KEM Pattern

Implements the hybrid encryption scheme from whitepaper Section 15.3.3. The message is encrypted once with a randomly generated Content Encryption Key (CEK). The CEK is wrapped individually for each recipient device using X25519 key exchange + HKDF-SHA256 to derive a wrapping key, then AES-256-GCM to encrypt the CEK.

```go
type RecipientRecord struct {
    DeviceID         [16]byte
    RecipientXPub    [32]byte         // X25519 public key of recipient device
    EphemeralXPub    [32]byte         // per-recipient ephemeral X25519 public key
    WrappedCEK       []byte           // AES-256-GCM ciphertext of CEK
    CEKNonce         [12]byte         // 96-bit nonce for CEK wrapping
    CEKTag           [16]byte         // GCM auth tag for CEK wrapping
}

type EncryptedEnvelope struct {
    Version          uint32
    MessageID        [16]byte
    Recipients       []RecipientRecord
    EncryptedPayload []byte           // AES-256-GCM ciphertext of SignedMessage
    PayloadNonce     [12]byte         // 96-bit nonce for payload
    PayloadTag       [16]byte         // GCM auth tag for payload
    PayloadSizeClass uint32           // padded size bucket
    CreatedAt        time.Time
    RatchetPubKey    [32]byte         // reserved; zero in protocol v1
}
```

`Encrypt(msg SignedMessage, recipients []RecipientInfo)` produces an `EncryptedEnvelope`. `Decrypt(envelope EncryptedEnvelope, recipientPrivKey [32]byte, deviceID [16]byte)` returns the `SignedMessage`. Both functions must validate all cryptographic material and return typed errors on failure.

> **RatchetPubKey:** This field is reserved for the Double Ratchet upgrade path (protocol v2). It must be present in the struct and serialised as a zero-valued 32-byte field in v1. Relay nodes and recipients must silently ignore it. Do not implement ratchet mechanics at this stage.

### 4.4 Protocol Buffer Definitions (`proto/`)

Define `.proto` (protobuf v3) schemas that correspond 1:1 to the Go structs above. The proto definitions are the canonical wire format; the Go structs are generated from them where practical, or manually maintained in parallel with explicit serialisation tests.

- `identity.proto` — `IdentityRecord`, `AttestationRecord`, `VerificationTier` enum
- `message.proto` — `PlaintextMessage`, `SignedMessage`, `EncryptedEnvelope`, `RecipientRecord`, `MessageBody`, `AttachmentRecord`

All binary fields (keys, signatures, nonces, tags) are `bytes` in proto and `[]byte` in Go. All timestamps are `int64` Unix seconds in proto and `time.Time` in Go, with explicit conversion functions. String fields are UTF-8.

### 4.5 Test Requirements for Milestone 1

Milestone 1 is not complete until all of the following tests pass:

1. **Round-trip:** `GenerateIdentityKeyPair` → `Sign` `IdentityRecord` → `Verify` succeeds.
2. **Tamper:** mutating any field of a signed `IdentityRecord` causes `Verify` to return an error.
3. **Round-trip:** compose `PlaintextMessage` → `Sign` → `Encrypt` to one recipient → `Decrypt` → `Verify` signature succeeds, plaintext matches.
4. **Multi-device:** `Encrypt` to three recipients → each can independently `Decrypt` and verify.
5. **Wrong key:** `Decrypt` with a key not in `Recipients` returns a typed error, not a panic.
6. **Tamper:** mutating `EncryptedPayload` causes `Decrypt` to return an authentication error.
7. **RatchetPubKey:** serialise and deserialise an `EncryptedEnvelope` and confirm `RatchetPubKey` is present as 32 zero bytes.
8. **Fingerprint:** `Fingerprint()` returns a 40-character uppercase hex string; two different key pairs produce different fingerprints.
9. **MessageID:** every `PlaintextMessage` produced by the constructor has a non-zero, unique `message_id` UUID. Two independently constructed messages must not share a `message_id`.

> **Coverage target:** `internal/core/crypto`, `internal/core/identity`, and `internal/core/message` must each reach 90% line coverage. Use `go test -cover` to measure.

---

## 5. Milestone 2 — Node and Registry

Milestone 2 introduces the network layer. The goal is a single binary (`dmcn-node`) that can register an identity to a local DHT instance and exchange a signed, encrypted message with another running instance of itself. Onion routing is explicitly excluded at this stage — messages are delivered directly between nodes.

### 5.1 DHT and Identity Registry (`internal/registry`)

Use libp2p (`github.com/libp2p/go-libp2p`) as the DHT foundation rather than implementing Kademlia from scratch. libp2p provides production-grade DHT infrastructure used by IPFS and Ethereum, and allows the PoC to focus on DMCN-specific protocol logic rather than DHT plumbing.

The registry package wraps libp2p and exposes the four operations defined in whitepaper Section 15.2.4:

| Operation | Behaviour |
|---|---|
| `REGISTER` | Store a signed `IdentityRecord` in the DHT, keyed on `SHA-256(address)`. Validates `SelfSignature` before storing. Idempotent. |
| `LOOKUP` | Retrieve an `IdentityRecord` by exact address string. Returns not-found if absent. Validates `SelfSignature` on retrieval. |
| `REVOKE` | Mark an identity as revoked. Revocation is permanent; revoked keys cannot be re-registered. |
| `UPDATE` | Replace an `IdentityRecord` with a new one. The new record must be signed by both the old and new private keys to prove continuity of control. |

> **Scope note:** For the PoC, the DHT can be operated as a local in-process instance or a small cluster of two or three nodes on localhost. Global DHT deployment is a post-PoC concern.

### 5.2 Relay Node (`internal/relay`)

Implement a minimal relay node that can receive, verify, store, and forward `EncryptedEnvelope`s. At this stage, direct delivery is used — the full onion routing transport from Section 15.4 is deferred to a later milestone.

The relay node must implement the following operations over TLS 1.3 connections:

| Operation | Behaviour |
|---|---|
| `STORE` | Receive an `EncryptedEnvelope`. Verify the sender's identity exists in the registry. Store the envelope indexed by recipient public key. Reject envelopes from unregistered senders. |
| `FETCH` | Recipient authenticates by signing a challenge nonce with their Ed25519 private key. Relay verifies the signature against the registry and returns pending envelopes for that identity. |
| `ACK` | Recipient confirms successful decryption. Relay marks the envelope as delivered. |
| `PING` | Liveness check. Returns node metadata: version, capacity, uptime. |

The relay node must enforce: sender identity must be registered in the DHT before a `STORE` is accepted. An envelope from an unregistered sender is dropped with a typed error response — this is the mechanism by which spam is rejected at the network boundary.

> **Rate limiting:** Implement basic per-sender rate limiting in the PoC: maximum 100 `STORE` operations per hour per registered identity. This is intentionally conservative; production limits are defined in Section 15.4.3 of the whitepaper.

### 5.3 The `dmcn-node` Binary

The `dmcn-node` binary is a combined development node that runs a DHT registry node and a relay node in a single process. It is not a production architecture — it exists to make end-to-end testing possible without running multiple separate services.

The binary exposes a minimal CLI for development use:

```
dmcn-node start --listen 0.0.0.0:7400 --dht-port 7401
dmcn-node identity generate --address alice@localhost
dmcn-node identity register --address alice@localhost
dmcn-node identity lookup --address bob@localhost
dmcn-node message send --from alice@localhost --to bob@localhost --body "hello"
dmcn-node message fetch --address alice@localhost
```

Identity key material is stored on disk in an encrypted keystore file, protected by a passphrase supplied at startup. The keystore format is documented and versioned.

### 5.4 End-to-End Test Scenario

Milestone 2 is complete when the following scenario executes successfully in an automated integration test:

1. Start two `dmcn-node` instances (node-A on port 7400, node-B on port 7402) sharing a local DHT.
2. Generate identity `alice@localhost` on node-A. Register it.
3. Generate identity `bob@localhost` on node-B. Register it.
4. node-A looks up `bob@localhost` and retrieves his `IdentityRecord`. Signature validates.
5. node-A composes a `PlaintextMessage`, signs it, encrypts it to bob's X25519 public key, and `STORE`s it on node-B's relay.
6. node-B authenticates and `FETCH`es its pending envelopes. Decrypts the envelope. Verifies the sender signature. Plaintext matches original.
7. node-B sends `ACK`. node-A confirms delivery.

> **Rejection test:** An additional test must confirm that a `STORE` from an unregistered identity is rejected by the relay node before the envelope enters storage.

---

## 6. Milestone 3 — Bridge Node

The bridge node is the first component that makes the PoC demonstrable to a non-technical audience: it allows a legacy email client (Gmail, Outlook, Apple Mail) to send a message that arrives in a DMCN inbox, and allows a DMCN user to reply to a legacy email address. This milestone produces the `dmcn-bridge` binary.

### 6.1 Inbound Path (SMTP to DMCN)

The bridge runs an SMTP listener that receives messages from legacy senders addressed to the DMCN user's bridge address (e.g. `alice@bridge.localhost` in development).

For each inbound SMTP message:

1. Perform SPF, DKIM, and DMARC verification on the inbound message.
2. Classify the sender into a trust tier: Verified Legacy (valid DKIM + pass), Unverified Legacy (no DKIM or neutral), Suspicious (failures or reputation flags).
3. Construct a `BridgeClassificationRecord` (signed by the bridge's own Ed25519 key) attesting to the classification outcome. This corresponds to the `bridge_classification_record` structure in whitepaper Section 15.6.2.
4. Look up the DMCN recipient's `IdentityRecord` from the registry.
5. Wrap the original message content and classification record in a `PlaintextMessage`, sign as bridge, encrypt to recipient's public key, and `STORE` on the recipient's relay node.

> **PoC scope:** For the PoC, reputation database lookups (RBL/DNSBL) may be stubbed. The classification logic and signed record format must be correctly implemented.

### 6.2 Outbound Path (DMCN to SMTP)

When a DMCN user sends to a legacy email address, the encrypted envelope is addressed to the bridge node's public key. The bridge:

1. Receives the `EncryptedEnvelope` addressed to its own public key.
2. Decrypts using its own X25519 private key.
3. Verifies the sender's Ed25519 signature against the registry.
4. Constructs a standard SMTP message. Sets the `From` address to a bridge-scoped representation of the sender's DMCN address.
5. DKIM-signs the outbound message using the bridge operator's domain key.
6. Delivers via SMTP to the recipient's MX.
7. Returns a signed `BridgeDeliveryReceipt` to the DMCN sender.

> **Trust disclosure:** The bridge operator has technical access to outbound message content. This is an unavoidable consequence of protocol translation, disclosed in whitepaper Section 11.2.2. The PoC must log a clear warning in the bridge server output whenever a message is decrypted for outbound delivery.

### 6.3 Bridge Node Identity

The bridge node is a registered DMCN identity with a `bridge_capability` flag in its `IdentityRecord`. It must be registered in the DHT before it can accept or deliver messages. The `dmcn-bridge` binary handles its own key generation and registration on first run.

### 6.4 End-to-End Test Scenario

Milestone 3 is complete when the following scenario executes successfully:

1. Start `dmcn-node` (local DHT + relay) and `dmcn-bridge`.
2. Register `alice@dmcn.localhost` as a DMCN identity.
3. Send an SMTP message to `alice@bridge.localhost` from an external mail client or smtp-test tool.
4. Confirm the message arrives in alice's DMCN fetch queue with a `BridgeClassificationRecord` attached.
5. `alice@dmcn.localhost` composes a reply addressed to the legacy sender's email address.
6. Confirm the reply is delivered as SMTP to the legacy address, DKIM-signed by the bridge.

---

## 7. Non-Functional Requirements

### 7.1 Error Handling

All functions that can fail must return an error as the final return value. Errors must be typed (using `errors.As` / `errors.Is` patterns) and must carry enough context to identify the failing operation without exposing private key material. Panics are not acceptable in production code paths — only in `init()` validation of compile-time constants.

### 7.2 Logging

Use a structured logger (`log/slog` from Go 1.21 standard library). Log levels: DEBUG for protocol trace, INFO for normal operation events, WARN for degraded conditions, ERROR for failures. Private key material must never appear in log output at any level. Message content must never appear in log output above DEBUG level, and DEBUG logging must be disabled by default.

### 7.3 Configuration

All node configuration is supplied via a YAML config file and/or environment variables. No hardcoded addresses, ports, or key material. The config schema is documented and validated at startup with clear error messages for missing or invalid values.

### 7.4 Testing

Unit tests live alongside their packages (`foo_test.go`). Integration tests live in a top-level `tests/` directory. The integration test suite must be runnable with a single command (`make test-integration` or equivalent) that starts all required local services and tears them down after. CI must run both unit and integration tests on every commit.

### 7.5 Dependencies

Minimise external dependencies. The approved external dependencies for the PoC are:

| Package | Purpose |
|---|---|
| `golang.org/x/crypto` | X25519, HKDF, additional crypto primitives |
| `google.golang.org/protobuf` | Protocol Buffer serialisation |
| `github.com/libp2p/go-libp2p` | DHT and peer-to-peer networking foundation |
| `gopkg.in/yaml.v3` | Configuration file parsing |
| `log/slog` | Structured logging (Go >= 1.21 standard library) |

Any additional dependency requires a brief written justification in the PR that introduces it. Dependencies with C bindings (cgo) are discouraged unless no pure-Go alternative exists.

### 7.6 Security Constraints

- Private key material must never be written to disk in plaintext. The keystore must be AES-256-GCM encrypted with a key derived from the user's passphrase via PBKDF2 or Argon2id.
- TLS 1.3 is required for all inter-node communication. TLS 1.2 must not be accepted.
- All received `IdentityRecord`s and `SignedMessage`s must have their signatures verified before any processing occurs. Skipping signature verification is not acceptable under any code path.
- The relay node must never store an envelope from an unregistered sender. This check must occur before any disk write.

---

## 8. Milestone Summary

| # | Milestone | Deliverable | Done When |
|---|---|---|---|
| M1 | Cryptographic Core | `internal/core` packages + proto definitions | All 9 specified tests pass at 90% coverage |
| M2 | Node and Registry | `internal/registry`, `internal/relay`, `cmd/dmcn-node` | End-to-end alice→bob message test passes |
| M3 | Bridge Node | `internal/bridge`, `cmd/dmcn-bridge` | Legacy SMTP↔DMCN round-trip test passes |

Milestones are sequential. M2 must not begin until M1's test suite is complete and passing. M3 must not begin until M2's integration test is complete and passing.

---

## 9. Explicitly Deferred (Post-PoC)

The following capabilities are confirmed out of scope for the PoC and must not be partially implemented or stubbed in ways that create technical debt:

- **Onion routing** (Section 15.4.1) — direct delivery is used throughout the PoC.
- **Double Ratchet forward secrecy** (whitepaper Section 15.7, protocol v2) — the `RatchetPubKey` field is reserved as zero bytes only.
- **Device state synchronisation** (whitepaper Section 15.3.5, `SyncEnvelope`) — the data structure is defined in the protocol spec; implementation is deferred. The PoC must populate `message_id` correctly on all `PlaintextMessage` instances to ensure no structural changes are required when sync is implemented.
- **Social key recovery** (Section 7.3) — not implemented; key material is passphrase-protected on disk.
- **Shared reputation feeds** (Section 14.3.2) — not implemented.
- **Full trust management stack** (allowlist, pending queue, blocklist) — not implemented in the PoC node; the data structures may be defined but no UI or enforcement logic.
- **Address portability / domain verification tiers** — not implemented; all identities are `TierUnverified` in the PoC.
- **Native desktop or mobile client** — the CLI is the only client interface in the PoC.
- **Production DHT deployment** — localhost only.

---

## 10. Reference

| Document | Description |
|---|---|
| DMCN Whitepaper v0.2 | Full protocol design and rationale. All section references in this PRD refer to this document. |
| Section 15.2 | Identity layer data structures and registry operations. |
| Section 15.3 | Message format: `PlaintextMessage`, `SignedMessage`, `EncryptedEnvelope`, KEM pattern. |
| Section 15.3.5 | `SyncEnvelope`: device state synchronisation structure. Defined in spec; implementation deferred to post-PoC. |
| Section 15.4 | Transport layer and relay node protocol (onion routing deferred to post-PoC). |
| Section 15.5 | Storage and delivery layer. |
| Section 15.6 | Bridge protocol interface and `BridgeClassificationRecord`. |
| Section 11 | SMTP-DMCN bridge architecture, inbound and outbound paths. |
| Section 17.5 | Sybil resistance — rationale for new identity rate limiting. |
| Section 19.7 | Standards process context — why prototype precedes IETF engagement. |
