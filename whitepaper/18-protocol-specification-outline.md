## 18. Protocol Specification Outline

This section provides a structured technical outline of the DMCN protocol. It is not a complete specification — a production-ready protocol specification would be published as a series of formal documents analogous to IETF RFCs — but it defines the principal data structures, message formats, and protocol flows with sufficient precision to guide prototype implementation and to invite technical critique.

The outline is organised into five layers: identity, addressing, message format, transport, and storage and delivery. A sixth subsection covers the bridge protocol interface. Each layer is described in terms of its data structures, the operations it supports, and its interface with adjacent layers.

> **Status**
> *This outline represents the current design intent as of Version 0.2. Field names, encodings, and parameter values are indicative and subject to revision through the prototype and community review process. Where open questions remain, they are explicitly noted.*

---

### 18.1 Encoding and Serialisation Conventions

All DMCN protocol messages use Protocol Buffers (protobuf) version 3 as the canonical wire encoding, chosen for its compact binary representation, language-neutral schema definitions, and forward-compatibility properties. JSON representations of the same schemas are defined for debugging, human-readable export, and bridge protocol use.

All binary fields (keys, signatures, hashes, nonces) are encoded as raw bytes in protobuf and as base64url (RFC 4648 §5, no padding) in JSON representations.

All timestamps are Unix epoch seconds as a 64-bit unsigned integer.

String fields use UTF-8 encoding. Address strings follow the `local@domain` format defined in Section 18.2.

Protocol version negotiation uses a single `uint32 version` field present in all top-level message types. The current protocol version is `1`. Nodes must reject messages with version numbers they do not support and return a `VERSION_NOT_SUPPORTED` error code.

---

### 18.2 Identity Layer

#### 18.2.1 Key Pair Specification

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

#### 18.2.2 Identity Record

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

#### 18.2.3 Attestation Record

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

#### 18.2.4 Identity Registry Operations

The distributed identity registry exposes four operations:

| Operation | Input | Output | Notes |
|---|---|---|---|
| `REGISTER` | `identity_record` | `ack` or `error` | Idempotent; re-registration updates the record if self-signature is valid |
| `LOOKUP` | `address: string` | `identity_record` or `not_found` | Rate-limited per source; see Section 17.3.1 |
| `REVOKE` | `address`, `revocation_signature` | `ack` or `error` | Revocation is permanent; revoked keys cannot be re-registered |
| `UPDATE` | `identity_record` | `ack` or `error` | For key rotation; triggers key-change notifications to whitelisted contacts |

Registry nodes maintain a Kademlia-style DHT keyed on the SHA-256 hash of the identity address string. Lookup queries converge in O(log N) hops where N is the number of **participating registry nodes**, not the number of registered identities. The number of identities affects how much data each node stores; it does not affect routing hop count. This distinction is significant for scalability: lookup latency grows logarithmically with the size of the node network, which is expected to remain orders of magnitude smaller than the identity population. A registry of 100,000 nodes serving 500 million identities converges in approximately log₂(100,000) ≈ 17 hops regardless of identity count.

**Key design consequence — address search is not supported.** Because registry entries are keyed on the SHA-256 hash of the address string, the DHT supports only exact-match lookups: a client must know the precise address `alice@example.com` in order to retrieve Alice's identity record. Prefix search, wildcard lookup, and domain-level enumeration (e.g. "all addresses at example.com") are not supported by the DHT structure, as hashing destroys address ordering and grouping. This is a deliberate design choice — it prevents bulk harvesting of registered identities — but it means that address discovery and autocomplete functionality must be implemented through a separate, opt-in directory service outside the core DHT. Clients that wish to offer contact search must either maintain a local address book populated through direct exchange, or query a supplementary directory operated by organisations that choose to publish their member addresses.

---

### 18.3 Message Format

#### 18.3.1 Plaintext Message

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
    content:          bytes          // encrypted separately; see Section 18.3.3
}
```

#### 18.3.2 Signed Message

Before encryption, the plaintext message is signed by the sender's Ed25519 private key. The signature covers the canonical protobuf serialisation of the `plaintext_message`.

```
signed_message {
    plaintext:        plaintext_message
    sender_signature: bytes[64]      // Ed25519 signature by sender over plaintext
}
```

Recipients must verify `sender_signature` after decryption. A message with an invalid or missing sender signature must be rejected and must not be displayed to the user.

#### 18.3.3 Encrypted Envelope

The `signed_message` is encrypted using a hybrid encryption scheme: an ephemeral X25519 key pair is generated for each message, a shared secret is derived via X25519 key exchange between the ephemeral private key and the recipient's `x25519_public_key`, and the shared secret is used to derive a symmetric key via HKDF-SHA256 for AES-256-GCM encryption of the message content.

```
encrypted_envelope {
    version:              uint32
    message_id:           bytes[16]      // matches plaintext_message.message_id
    recipient_pubkey:     bytes[32]      // X25519 public key of intended recipient
    ephemeral_pubkey:     bytes[32]      // ephemeral X25519 public key
    encrypted_payload:    bytes          // AES-256-GCM ciphertext of signed_message
    aead_tag:             bytes[16]      // GCM authentication tag
    nonce:                bytes[12]      // 96-bit random nonce for AES-GCM
    payload_size_class:   uint32         // padded size class (see Section 17.2.3)
    created_at:           uint64
}
```

The `payload_size_class` field records the size bucket into which the payload has been padded (e.g. 1KB, 4KB, 16KB, 64KB, 256KB, 1MB), not the actual payload size. Relay nodes and passive observers can observe only the size class, not the precise message size.

#### 18.3.4 Attachment Handling

Attachments are encrypted separately from the message body using the same hybrid scheme, with a separate ephemeral key pair per attachment. The `attachment_record` in the `plaintext_message` contains the `content_hash` of the plaintext attachment for integrity verification after decryption, but the attachment content itself is stored as a separately addressed blob in the storage layer, referenced by `attachment_id`.

This separation allows large attachments to be stored and retrieved independently of the message envelope, reducing storage requirements at relay nodes that buffer messages for offline recipients.

---

### 18.4 Transport Layer

#### 18.4.1 Onion Routing Packet Format

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

#### 18.4.2 Relay Node Protocol

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

#### 18.4.3 Flow Control and Rate Limiting

Relay nodes implement per-sender rate limiting based on the sender's registered identity. New identities (registered within the past 30 days) are subject to stricter throughput limits than established identities, implementing the reputation bootstrapping behaviour described in Section 14.5.

Rate limits are expressed as:
- Maximum messages per hour per sending identity: default 500 (new identity), 5000 (established)
- Maximum total payload bytes per hour per sending identity: default 10MB (new), 100MB (established)
- Maximum recipient identities per hour per sending identity: default 50 (new), 500 (established)

These defaults are configurable by relay node operators and represent recommended baseline values for the reference implementation.

---

### 18.5 Storage and Delivery Layer

#### 18.5.1 Message Store

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

#### 18.5.2 Recipient Fetch Protocol

When a recipient client connects to retrieve messages, it authenticates by signing a challenge nonce with its Ed25519 private key. The relay node verifies the signature against the identity registry and returns all `stored_message_record` headers for messages addressed to that identity. The client then fetches the full `encrypted_envelope` for each message it wishes to retrieve.

This two-phase fetch (headers first, then content on demand) allows clients to make informed decisions about large attachments before downloading, and supports efficient operation on constrained network connections.

#### 18.5.3 Delivery Receipts

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

### 18.6 Bridge Protocol Interface

The SMTP-DMCN Bridge Operator Protocol (BOP) defines the interface between bridge nodes and the core DMCN network. Bridge nodes are registered DMCN identities with an additional `bridge_capability` flag set in their identity record.

#### 18.6.1 Outbound Bridge Message

When a DMCN client sends a message to a legacy email address, the encrypted envelope is addressed to the bridge node's public key rather than a recipient public key. The bridge node decrypts, reconstructs the SMTP message, and delivers it. The bridge attaches a standardised footer and DKIM-signs the outbound SMTP message using its registered domain key.

The outbound bridge flow is:
1. Client looks up bridge node identity from registry
2. Client encrypts signed message to bridge node's X25519 public key
3. Client routes encrypted envelope through transport layer to bridge node
4. Bridge node decrypts, verifies sender signature, constructs SMTP message
5. Bridge node delivers via standard SMTP with DKIM signing
6. Bridge node sends a signed `bridge_delivery_receipt` back to sender

#### 18.6.2 Inbound Bridge Classification Record

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

### 18.7 Protocol Extension Mechanism

The DMCN protocol is designed to be extensible without breaking backward compatibility. Each top-level message type includes a `repeated extension_field extensions` field using protobuf's extension mechanism. Nodes that do not understand a given extension field must ignore it and must not reject the message on that basis.

Proposed extensions are published as numbered extension specifications (analogous to IETF Internet-Drafts) and progress through a community review process before being assigned stable extension numbers and included in the reference implementation.

Planned first-generation extensions include: group messaging (multi-recipient encrypted envelopes), message expiry (sender-specified deletion after a time period), read receipts distinct from delivery receipts, and rich text body support beyond the base `text/plain` and `text/html` content types.



---
