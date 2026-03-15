# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DMCN (Decentralized Mesh Communication Network) is a peer-to-peer, end-to-end encrypted messaging protocol designed to replace email. It uses cryptographic identity (Ed25519 signing + X25519 key exchange) instead of SMTP-style trust. This is a proof-of-concept implementation in Go.

The whitepaper is in `docs/whitepaper/`. Code references whitepaper sections in comments (e.g., "See whitepaper Section 15.2.1").

## Build & Development Commands

```bash
make build          # Build binary to bin/dmcn-node
make test           # Run all tests (120s timeout)
make test-cover     # Run tests with per-package coverage
make lint           # Run buf lint + go vet
make proto          # Regenerate protobuf Go code (requires buf CLI)
make clean          # Remove coverage files and binary
```

Run a single package's tests:
```bash
go test ./internal/core/crypto/...
go test ./internal/relay/... -timeout 120s
```

The integration test in `internal/node/` spins up real libp2p nodes and requires the 120s timeout.

## Architecture

### Three-Layer Message Model (whitepaper Section 15.3)

1. **PlaintextMessage** (`internal/core/message/message.go`) - composed message with sender/recipient, body, attachments
2. **SignedMessage** - wraps PlaintextMessage with Ed25519 sender signature
3. **EncryptedEnvelope** (`internal/core/message/encrypt.go`) - hybrid encryption using per-message CEK (AES-256-GCM) wrapped per-recipient via X25519 ECDH + HKDF. Payloads are padded to size-class buckets for traffic analysis resistance.

### Key Packages

- `internal/core/crypto/` - Thin wrappers around Go stdlib crypto (Ed25519, X25519, AES-256-GCM, HKDF-SHA256). No custom crypto implementations. Uses `randReader` var that tests can override.
- `internal/core/identity/` - Identity key pairs (Ed25519 + X25519) and self-certifying IdentityRecords. Address format is `local@domain`.
- `internal/core/message/` - Message composition, signing, and hybrid encryption/decryption.
- `internal/registry/` - DHT-based identity registry using libp2p Kademlia. Records are keyed on `SHA256(address)` under the `/dmcn/` namespace. Includes a `record.Validator` implementation for DHT record validation.
- `internal/relay/` - Message relay service over libp2p streams (protocol `/dmcn/relay/1.0.0`). Length-prefixed protobuf wire protocol. Supports STORE (with sender signature verification + rate limiting), FETCH (with challenge-response auth), ACK, and PING operations. In-memory message store (PoC only).
- `internal/keystore/` - Encrypted on-disk key storage (AES-256-GCM with HKDF-derived key from passphrase). JSON format.
- `internal/node/` - Combined development node that runs DHT registry + relay in a single process.
- `cmd/dmcn-node/` - CLI entrypoint with subcommands: `start`, `identity generate|register|lookup`, `message send|fetch`.

### Protobuf

Proto definitions are in `proto/` (identity.proto, message.proto, relay.proto). Generated Go code goes to `internal/proto/dmcnpb/`. Uses buf v2 for generation and linting.

### Serialization Convention

All signable data is serialized using deterministic protobuf marshaling (`proto.MarshalOptions{Deterministic: true}`). The `protoMarshal` var in identity and message packages can be overridden in tests.

### Logging

Uses `github.com/mertenvg/logr/v2` for structured logging. Key conventions:
- The CLI (`cmd/dmcn-node/main.go`) initializes logr with `logr.AddWriter(os.Stderr, ...)` and `logr.Verbose` filter (Panic, Error, Warning, Info, Success — no Debug).
- Internal packages (`node`, `relay`) create component-scoped loggers via `logr.With(logr.M("component", "..."))`.
- `node.New()` accepts an optional `logr.Logger` as a variadic parameter for caller-provided loggers.
- Log levels: `logr.Error`/`logr.Warn` for failures and rejections, `logr.Info` for operational events, `logr.Debug` for protocol-level detail (STORE/FETCH traces), `logr.Success` for completed user actions.
- Private key material must never appear in log output.
- Call `logr.Wait()` before `os.Exit()` to flush pending log messages.

## Key Patterns

- Functional options pattern used for `relay.New()` and `registry.New()` configuration.
- Relay has both server-side stream handlers and `Client*` methods for sending requests to remote nodes.
- Error sentinel values are used throughout (e.g., `registry.ErrNotFound`, `relay.ErrRateLimited`). Check with `errors.Is()`.
- The relay stores envelopes indexed by hex-encoded recipient X25519 public key, not by email address.
