## 15. Protocol Specification Outline

This section provides a structured technical outline of the DMCN protocol. It is not a complete specification — a production-ready protocol specification would be published as a series of formal documents analogous to IETF RFCs — but it defines the principal data structures, message formats, and protocol flows with sufficient precision to guide prototype implementation and to invite technical critique.

The outline is organised into five layers: identity, addressing, message format, transport, and storage and delivery. A sixth subsection covers the bridge protocol interface. Each layer is described in terms of its data structures, the operations it supports, and its interface with adjacent layers.

> **Status**
> *This outline represents the current design intent as of Version 0.2. Field names, encodings, and parameter values are indicative and subject to revision through the prototype and community review process. Where open questions remain, they are explicitly noted.*

---

### 15.1 Encoding and Serialisation Conventions

All DMCN protocol messages use Protocol Buffers (protobuf) version 3 as the canonical wire encoding, chosen for its compact binary representation, language-neutral schema definitions, and forward-compatibility properties. JSON representations of the same schemas are defined for debugging, human-readable export, and bridge protocol use.

All binary fields (keys, signatures, hashes, nonces) are encoded as raw bytes in protobuf and as base64url (RFC 4648 §5, no padding) in JSON representations.

All timestamps are Unix epoch seconds as a 64-bit unsigned integer.

String fields use UTF-8 encoding. Address strings follow the `local@domain` format defined in Section 15.2.

Protocol version negotiation uses a single `uint32 version` field present in all top-level message types. The current protocol version is `1`. Nodes must reject messages with version numbers they do not support and return a `VERSION_NOT_SUPPORTED` error code.

---

### 15.2 Identity Layer

#### 15.2.1 Key Pair Specification

Each DMCN identity is represented by an elliptic curve key pair using **Curve25519** for key exchange (X25519) and **Ed25519** for signatures. These two curves are mathematically related (both defined over the same field) and are used in combination throughout the Signal Protocol and modern TLS 1.3.

```
identity_keypair {
    ed25519_public_key:   bytes[32]   // signing public key
    x25519_public_key:    bytes[32]   // key exchange public key
    created_at:           uint64      // Unix timestamp of key generation
    device_id:            bytes[16]   // random UUID identifying the generating device
}
```

The private keys corresponding to `ed25519_public_key` and `x25519_public_key` never leave the device's secure enclave and are not represented in any protocol message.

A **fingerprint** is defined as the first 20 bytes of the SHA-256 hash of the concatenation of `ed25519_public_key` and `x25519_public_key`, encoded as a 40-character uppercase hex string for display purposes (e.g. `A3F2...B901`). Fingerprints are used for out-of-band identity verification.

#### 15.2.2 Identity Record

An identity record is the unit of data published to the distributed identity registry. It binds a human-readable address to a key pair and is signed by the identity's own Ed25519 private key, making the binding self-certifying.

```
identity_record {
    version:              uint32      // protocol version
    address:              string      // e.g. "alice@example.com"
    ed25519_public_key:   bytes[32]
    x25519_public_key:    bytes[32]
    created_at:           uint64
    expires_at:           uint64      // 0 = no expiry; positive = Unix timestamp
    relay_hints:          repeated string  // suggested relay node addresses
    verification_tier:    enum { UNVERIFIED, PROVIDER_HOSTED, DOMAIN_DNS, DANE }
    attestations:         repeated attestation_record  // optional web-of-trust
    self_signature:       bytes[64]   // Ed25519 signature over all preceding fields
}
```

The `self_signature` is computed over the canonical protobuf serialisation of all fields except `self_signature` itself. Any node receiving an identity record must verify this signature before storing or forwarding the record.

#### 15.2.3 Attestation Record

An attestation is a signed statement by one identity vouching for another. Attestations are optional and may be published publicly or retained privately by the attesting party.

```
attestation_record {
    attester_address:     string      // address of the attesting identity
    attester_pubkey:      bytes[32]   // Ed25519 public key of attester
    subject_address:      string      // address being attested
    subject_pubkey:       bytes[32]   // Ed25519 public key of subject
    attestation_type:     enum { IN_PERSON, FINGERPRINT, NETWORK, ORGANISATIONAL }
    attested_at:          uint64
    expires_at:           uint64
    signature:            bytes[64]   // Ed25519 signature by attester over all preceding fields
}
```

#### 15.2.4 Identity Registry Operations

The distributed identity registry exposes four operations:

| Operation | Input | Output | Notes |
|---|---|---|---|
| `REGISTER` | `identity_record` | `ack` or `error` | Idempotent; re-registration updates the record if self-signature is valid |
| `LOOKUP` | `address: string` | `identity_record` or `not_found` | Rate-limited per source; see Section 18.3.1 |
| `REVOKE` | `address`, `revocation_signature` | `ack` or `error` | Revocation is permanent; revoked keys cannot be re-registered |
| `UPDATE` | `identity_record` | `ack` or `error` | For key rotation; triggers key-change notifications to whitelisted contacts |

Registry nodes maintain a Kademlia-style DHT keyed on the SHA-256 hash of the identity address string. Lookup queries converge in O(log N) hops where N is the number of **participating registry nodes**, not the number of registered identities. The number of identities affects how much data each node stores; it does not affect routing hop count. This distinction is significant for scalability: lookup latency grows logarithmically with the size of the node network, which is expected to remain orders of magnitude smaller than the identity population. A registry of 100,000 nodes serving 500 million identities converges in approximately log₂(100,000) ≈ 17 hops regardless of identity count.

**Key design consequence — address search is not supported.** Because registry entries are keyed on the SHA-256 hash of the address string, the DHT supports only exact-match lookups: a client must know the precise address `alice@example.com` in order to retrieve Alice's identity record. Prefix search, wildcard lookup, and domain-level enumeration (e.g. "all addresses at example.com") are not supported by the DHT structure, as hashing destroys address ordering and grouping. This is a deliberate design choice — it prevents bulk harvesting of registered identities — but it means that address discovery and autocomplete functionality must be implemented through a separate, opt-in directory service outside the core DHT. Clients that wish to offer contact search must either maintain a local address book populated through direct exchange, or query a supplementary directory operated by organisations that choose to publish their member addresses.

#### 15.2.5 Device Sub-Key Record

As described in Section 7.5, each device on which a user activates their DMCN account generates its own sub-key pair. Sub-keys are registered in the identity registry as children of the primary identity record and returned alongside it in `LOOKUP` responses.

```
device_subkey_record {
    version:                uint32
    primary_address:        string      // address of the owning identity
    primary_ed25519_pubkey: bytes[32]   // Ed25519 public key of the primary identity
    device_id:              bytes[16]   // random UUID identifying the device
    device_label:           string      // optional human-readable label (e.g. "iPhone 15")
    sub_ed25519_pubkey:     bytes[32]   // Ed25519 signing key for this device
    sub_x25519_pubkey:      bytes[32]   // X25519 encryption key for this device
    created_at:             uint64
    expires_at:             uint64      // 0 = no expiry; positive = Unix timestamp
    primary_signature:      bytes[64]   // Ed25519 signature by primary key over all preceding fields
    device_self_signature:  bytes[64]   // Ed25519 signature by sub key over all preceding fields
}
```

Both signatures must be present and valid for a sub-key record to be accepted by the registry. The `primary_signature` proves the sub-key is authorised by the identity owner. The `device_self_signature` proves the registering device holds the corresponding private key and prevents a primary-key holder from registering phantom sub-keys for devices they do not control.

The `device_label` field is optional and intended solely for the user's own client UI — to allow them to identify and revoke specific devices by name. It is not used in any routing or trust decision and should not be relied upon by external parties.

**Sub-key registry operations** extend the four base operations defined in Section 15.2.4:

| Operation | Input | Output | Notes |
|---|---|---|---|
| `SUBKEY_REGISTER` | `device_subkey_record` | `ack` or `error` | Adds sub-key to the primary identity's sub-key set |
| `SUBKEY_REVOKE` | `primary_address`, `device_id`, `revocation_signature` | `ack` or `error` | Revocation is permanent; signed by either primary key or the sub-key itself |
| `SUBKEY_LIST` | `primary_address` | `repeated device_subkey_record` or `not_found` | Returns all active sub-keys; included in `LOOKUP` response |

`SUBKEY_REVOKE` may be signed by the primary key (owner-initiated, e.g. for a lost device) or by the sub-key itself (device-initiated, e.g. on a clean logout or decommission). Both are valid. A primary-key-signed revocation takes effect immediately; a sub-key-signed revocation is also immediate but additionally signals to contacts that the device performed a clean deactivation rather than being forcibly removed.

**Encryption to multiple sub-keys.** When a sender looks up a recipient who has multiple active sub-keys, the client encrypts the message independently to each active `sub_x25519_pubkey`. Each encrypted copy is addressed to its sub-key and delivered through the transport layer. Relay nodes and storage nodes treat each copy as a separate message and deliver them independently. The recipient's client on whichever device retrieves the message first decrypts its copy; the copies on other devices are retrieved when those devices next connect. This is the same multi-device delivery model used by the Signal protocol.

**Interaction with the primary key.** The primary `identity_record` (Section 15.2.2) retains its own `x25519_public_key` field, which is used when the recipient has no active sub-keys registered — for example, immediately after account creation before any device sub-key has been issued, or during account recovery before a new device sub-key is established. Senders should prefer sub-keys when available and fall back to the primary key only when no active sub-keys are present.

---

### 15.3 Message Format

#### 15.3.1 Plaintext Message

Before encryption, a DMCN message has the following structure:

```
plaintext_message {
    version:          uint32
    message_id:       bytes[16]      // random UUID, globally unique
    thread_id:        bytes[16]      // UUID linking messages in a conversation thread
    sender_address:   string
    sender_pubkey:    bytes[32]      // Ed25519 public key (for recipient verification)
    recipient_address: string
    sent_at:          uint64         // Unix timestamp
    subject:          string         // optional; empty string if no subject
    body:             message_body
    attachments:      repeated attachment_record
    reply_to_id:      bytes[16]      // optional; message_id of the message being replied to
}

message_body {
    content_type:     string         // MIME type, e.g. "text/plain" or "text/html"
    content:          bytes          // UTF-8 encoded body text
}

attachment_record {
    attachment_id:    bytes[16]
    filename:         string
    content_type:     string
    size_bytes:       uint64
    content_hash:     bytes[32]      // SHA-256 of the plaintext attachment content
    content:          bytes          // encrypted separately; see Section 15.3.3
}
```

#### 15.3.2 Signed Message

Before encryption, the plaintext message is signed by the sender's Ed25519 private key. The signature covers the canonical protobuf serialisation of the `plaintext_message`.

```
signed_message {
    plaintext:        plaintext_message
    sender_signature: bytes[64]      // Ed25519 signature by sender over plaintext
}
```

Recipients must verify `sender_signature` after decryption. A message with an invalid or missing sender signature must be rejected and must not be displayed to the user.

#### 15.3.3 Encrypted Envelope

The DMCN uses a **Key Encapsulation Mechanism (KEM)** pattern for message encryption. This separates the encryption of the message payload (which happens once, regardless of how many devices the recipient has enrolled) from the distribution of the decryption key (which is wrapped individually for each intended recipient key). The result is that large payloads and attachments appear on the wire exactly once, with only a small per-recipient overhead for the wrapped key material.

**Step 1 — Generate a content key.** The sender generates a random 256-bit symmetric content key (CEK, Content Encryption Key) for the message. This key is used to encrypt the `signed_message` payload once using AES-256-GCM.

**Step 2 — Wrap the CEK for each recipient key.** For each device sub-key (or primary key, if no sub-keys are active) of the intended recipient, the sender performs an X25519 key exchange between a freshly generated ephemeral private key and the recipient key's `x25519_public_key`. The resulting shared secret is passed through HKDF-SHA256 to derive a 256-bit key-wrapping key (KWK), which is used to encrypt the CEK using AES-256-GCM. Each such wrapped CEK, together with the ephemeral public key used to produce it, forms a `recipient_record`. The ephemeral key pair is discarded after wrapping; a distinct ephemeral key is generated per recipient key.

**Step 3 — Assemble the envelope.** The encrypted payload and the set of recipient records are assembled into a single `encrypted_envelope`:

```
recipient_record {
    recipient_pubkey:     bytes[32]   // X25519 public key this record is wrapped for
    ephemeral_pubkey:     bytes[32]   // ephemeral X25519 public key used for this wrapping
    wrapped_cek:          bytes[32]   // AES-256-GCM ciphertext of the 256-bit CEK
    wrap_aead_tag:        bytes[16]   // GCM authentication tag for the CEK wrapping
    wrap_nonce:           bytes[12]   // 96-bit random nonce for the CEK wrapping
}

encrypted_envelope {
    version:              uint32
    message_id:           bytes[16]               // matches plaintext_message.message_id
    recipients:           repeated recipient_record // one entry per enrolled device key
    encrypted_payload:    bytes                    // AES-256-GCM ciphertext of signed_message
    payload_aead_tag:     bytes[16]               // GCM authentication tag for the payload
    payload_nonce:        bytes[12]               // 96-bit random nonce for payload encryption
    payload_size_class:   uint32                  // padded size class (see Section 18.2.3)
    created_at:           uint64
}
```

**Decryption.** A recipient device locates the `recipient_record` whose `recipient_pubkey` matches its own device sub-key. It performs X25519 key exchange between its private sub-key and the `ephemeral_pubkey` in that record, derives the KWK via HKDF-SHA256, and uses it to unwrap the CEK. It then uses the CEK to decrypt the `encrypted_payload`. No other device's `recipient_record` is needed or accessed.

**Wire overhead.** Each `recipient_record` is 108 bytes (32 + 32 + 32 + 16 + 12 + 4 bytes padding alignment). For a user with five enrolled devices, the total per-recipient overhead is 540 bytes — negligible relative to even a minimal message payload. The payload itself, regardless of size, is encrypted and transmitted exactly once.

The `payload_size_class` field records the size bucket into which the payload has been padded (e.g. 1KB, 4KB, 16KB, 64KB, 256KB, 1MB), not the actual payload size. Relay nodes and passive observers can observe only the size class, not the precise message size.

#### 15.3.4 Attachment Handling

Attachments use the same KEM pattern as the message envelope, but are encrypted and stored independently of it. Each attachment is encrypted with its own randomly generated CEK, and that CEK is wrapped for each recipient device sub-key exactly as in Section 15.3.3. The resulting `attachment_envelope` has the same `recipients` / `encrypted_payload` structure as the message envelope, with the attachment ciphertext as the payload.

The `attachment_record` embedded in the `plaintext_message` contains the `content_hash` of the plaintext attachment for integrity verification after decryption, but the attachment ciphertext itself is stored as a separately addressed blob in the storage layer, referenced by `attachment_id`. This separation allows large attachments to be stored and retrieved independently of the message envelope, reducing storage requirements at relay nodes that buffer messages for offline recipients, and allowing recipients to defer download of large attachments until they choose to open them.

Because each attachment has its own CEK, a recipient who receives a message with three attachments can decrypt the message body immediately using the message CEK, and decrypt each attachment independently as they open it — without re-fetching or re-processing the message envelope.

---

### 15.4 Transport Layer

#### 15.4.1 Onion Routing Packet Format

Messages in the DMCN transport layer are wrapped in an onion routing structure with a fixed number of layers (default: 3 hops). Each layer is encrypted to the relay node at that position in the route, and contains the routing instruction for that hop.

```
onion_packet {
    version:          uint32
    layer:            onion_layer     // the layer for the current node
    next_payload:     bytes           // encrypted blob for the next hop (or final delivery)
}

onion_layer {
    next_hop:         string          // address of the next relay node, or "DELIVER" for final hop
    layer_signature:  bytes[64]       // Ed25519 signature by the originating sender
    created_at:       uint64
    ttl_seconds:      uint32          // time-to-live; nodes drop expired packets
}
```

The sender constructs the onion packet by layering encryptions from the innermost hop outward. Each relay node decrypts its layer, reads the `next_hop` instruction, and forwards the `next_payload` to the specified next node. No relay node can determine both the origin and the destination of the packet.

Route selection is performed by the sender's client, which queries the identity registry for relay node candidates and selects a path based on geographic distribution, node reputation, and latency estimates. The route is not disclosed to relay nodes — each knows only its predecessor and successor.

#### 15.4.2 Relay Node Protocol

Relay nodes communicate over persistent TLS 1.3 connections using a simple request-response protocol. The primary relay node operations are:

| Operation | Initiator | Description |
|---|---|---|
| `RELAY` | Sender node | Forward an onion packet to the next hop |
| `STORE` | Previous relay | Store an encrypted envelope for an offline recipient |
| `FETCH` | Recipient client | Retrieve stored envelopes for the authenticated identity |
| `ACK` | Any | Acknowledge receipt of a relayed or stored message |
| `PING` | Any | Liveness check; used for routing table maintenance |
| `NODE_INFO` | Any | Retrieve relay node metadata (capacity, supported versions, reputation) |

Relay nodes authenticate to each other and to clients using their registered DMCN identities. A relay node that presents an identity not found in the identity registry, or whose self-signature is invalid, must be rejected.

#### 15.4.3 Flow Control and Rate Limiting

Relay nodes implement per-sender rate limiting based on the sender's registered identity. New identities (registered within the past 30 days) are subject to stricter throughput limits than established identities, implementing the reputation bootstrapping behaviour described in Section 17.5.

Rate limits are expressed as:
- Maximum messages per hour per sending identity: default 500 (new identity), 5000 (established)
- Maximum total payload bytes per hour per sending identity: default 10MB (new), 100MB (established)
- Maximum recipient identities per hour per sending identity: default 50 (new), 500 (established)

These defaults are configurable by relay node operators and represent recommended baseline values for the reference implementation.

---

### 15.5 Storage and Delivery Layer

#### 15.5.1 Message Store

Relay nodes providing storage services maintain an encrypted message store indexed by recipient public key. The store is content-addressed: each stored item is identified by the SHA-256 hash of its `encrypted_envelope`, allowing deduplication across relay nodes that may both receive the same message.

```
stored_message_record {
    envelope_hash:        bytes[32]      // SHA-256 of encrypted_envelope
    recipient_pubkey:     bytes[32]      // X25519 public key of recipient
    stored_at:            uint64
    expires_at:           uint64         // relay node retention expiry
    size_class:           uint32         // from encrypted_envelope.payload_size_class
    delivery_status:      enum { PENDING, DELIVERED, EXPIRED }
}
```

The `encrypted_envelope` itself is stored as an opaque blob. Relay nodes cannot read its contents.

#### 15.5.2 Recipient Fetch Protocol

When a recipient client connects to retrieve messages, it authenticates by signing a challenge nonce with its Ed25519 private key. The relay node verifies the signature against the identity registry and returns all `stored_message_record` headers for messages addressed to that identity. The client then fetches the full `encrypted_envelope` for each message it wishes to retrieve.

This two-phase fetch (headers first, then content on demand) allows clients to make informed decisions about large attachments before downloading, and supports efficient operation on constrained network connections.

#### 15.5.3 Delivery Receipts

The DMCN supports optional end-to-end delivery receipts. When the recipient client successfully decrypts and verifies a message, it may send a signed receipt back to the sender through the transport layer.

```
delivery_receipt {
    message_id:       bytes[16]      // matches plaintext_message.message_id
    recipient_address: string
    delivered_at:     uint64
    receipt_signature: bytes[64]     // Ed25519 signature by recipient
}
```

Delivery receipts are encrypted to the sender's public key and routed through the standard transport layer. They are optional — recipients may disable receipt sending as a privacy measure.

---

### 15.6 Bridge Protocol Interface

The SMTP-DMCN Bridge Operator Protocol (BOP) defines the interface between bridge nodes and the core DMCN network. Bridge nodes are registered DMCN identities with an additional `bridge_capability` flag set in their identity record.

#### 15.6.1 Outbound Bridge Message

When a DMCN client sends a message to a legacy email address, the encrypted envelope is addressed to the bridge node's public key rather than a recipient public key. The bridge node decrypts, reconstructs the SMTP message, and delivers it. The bridge attaches a standardised footer and DKIM-signs the outbound SMTP message using its registered domain key.

The outbound bridge flow is:
1. Client looks up bridge node identity from registry
2. Client encrypts signed message to bridge node's X25519 public key
3. Client routes encrypted envelope through transport layer to bridge node
4. Bridge node decrypts, verifies sender signature, constructs SMTP message
5. Bridge node delivers via standard SMTP with DKIM signing
6. Bridge node sends a signed `bridge_delivery_receipt` back to sender

#### 15.6.2 Inbound Bridge Classification Record

For inbound messages from SMTP senders, the bridge node wraps the classified message in a DMCN envelope and attaches a signed classification record:

```
bridge_classification_record {
    bridge_address:       string      // DMCN address of the bridge operator
    bridge_pubkey:        bytes[32]   // Ed25519 public key of bridge operator
    smtp_from:            string      // original SMTP From address
    smtp_sender_ip:       string      // sending server IP
    spf_result:           enum { PASS, FAIL, SOFTFAIL, NEUTRAL, NONE }
    dkim_result:          enum { PASS, FAIL, NONE }
    dmarc_result:         enum { PASS, FAIL, NONE }
    reputation_score:     int32       // -100 to +100; 0 = neutral
    trust_tier:           enum { VERIFIED_LEGACY, UNVERIFIED_LEGACY, SUSPICIOUS }
    classified_at:        uint64
    bridge_signature:     bytes[64]   // Ed25519 signature over all preceding fields
}
```

Recipients can verify the `bridge_signature` to confirm the classification was produced by a registered, trusted bridge operator.

---

### 15.7 Protocol Extension Mechanism

The DMCN protocol is designed to be extensible without breaking backward compatibility. Each top-level message type includes a `repeated extension_field extensions` field using protobuf's extension mechanism. Nodes that do not understand a given extension field must ignore it and must not reject the message on that basis.

Proposed extensions are published as numbered extension specifications (analogous to IETF Internet-Drafts) and progress through a community review process before being assigned stable extension numbers and included in the reference implementation.

Planned first-generation extensions include: group messaging (multi-recipient encrypted envelopes), message expiry (sender-specified deletion after a time period), read receipts distinct from delivery receipts, and rich text body support beyond the base `text/plain` and `text/html` content types.



---


---

