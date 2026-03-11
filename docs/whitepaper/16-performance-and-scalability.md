## 16. Performance and Scalability Analysis

This section provides quantitative estimates of the DMCN's performance and scalability characteristics under realistic operating conditions. The estimates are derived from first-principles analysis of the proposed architecture, benchmarks of comparable systems, and published performance data for the cryptographic primitives involved. They are presented with explicit assumptions and uncertainty ranges, not as guaranteed specifications.

The purpose of this analysis is twofold: to demonstrate that the proposed architecture is viable at global email scale, and to identify the components that present the greatest engineering challenge and will require the most careful optimisation in prototype development.

> **Methodology**
> *All estimates use order-of-magnitude reasoning with stated assumptions. Where published benchmarks for comparable systems exist, they are cited. Estimates should be treated as planning figures rather than performance guarantees. Prototype benchmarking will be required to validate or revise these figures.*

---

### 16.1 Scale Targets

The DMCN must be capable of supporting global email-scale usage to be a credible replacement for SMTP. The following figures define the scale targets against which the architecture is assessed:

| Metric | Current Global Email Scale | DMCN Target (Year 5) |
|---|---|---|
| Active users | ~4 billion | 50–500 million (1–12% adoption) |
| Messages sent per day | ~350 billion (including spam) | 5–50 billion (spam-free) |
| Messages per second (peak) | ~4 million | 60,000–600,000 |
| Average message size | ~75KB (with attachments) | ~50KB (encrypted envelope) |
| Identity registry entries | N/A | 50–500 million |
| Relay nodes | N/A | 10,000–100,000 |

The Year 5 target represents a realistic early-adoption scenario — comparable to Signal's growth trajectory in its first five years — rather than full global deployment. Full global-scale deployment is a longer-horizon target that the architecture must support in principle but need not be optimised for in the prototype phase.

---

### 16.2 Cryptographic Operation Latency

Every message in the DMCN requires a fixed set of cryptographic operations at the sender, at each relay node, and at the recipient. The latency contribution of these operations is the irreducible floor below which no optimisation can reduce message latency.

The following benchmarks are for modern mid-range hardware (Apple M2, AMD Ryzen 7, comparable server CPUs). Mobile device benchmarks are approximately 3–5× slower for the same operations.

| Operation | Location | Benchmark | Notes |
|---|---|---|---|
| Ed25519 key generation | Sender (once, at account creation) | ~50 µs | One-time cost; hardware-accelerated on modern devices |
| Ed25519 sign | Sender (per message) | ~20 µs | Signing the plaintext_message |
| X25519 key exchange + HKDF | Sender (per message) | ~30 µs | Deriving symmetric key from ephemeral pair |
| AES-256-GCM encrypt (50KB) | Sender (per message) | ~150 µs | Hardware AES-NI; ~3.3 GB/s throughput |
| Ed25519 verify | Relay node (per message) | ~40 µs | Verifying sender signature against registry |
| AES-256-GCM decrypt (50KB) | Recipient (per message) | ~150 µs | Hardware AES-NI |
| Ed25519 verify | Recipient (per message) | ~40 µs | Verifying sender signature post-decryption |
| **Total cryptographic latency (sender)** | — | **~200 µs** | Dominated by encryption |
| **Total cryptographic latency (recipient)** | — | **~190 µs** | Dominated by decryption |
| **Total cryptographic latency (relay node, per message)** | — | **~40 µs** | Signature verification only |

Cryptographic latency is negligible relative to network latency for typical internet paths (10–100ms). It is not a bottleneck at any realistic message volume.

---

### 16.3 Identity Registry Performance

The identity registry is the component most likely to be a bottleneck at global scale, because every new message to an unknown recipient requires a registry lookup, and the registry must support consistent reads across a globally distributed DHT.

#### 16.3.1 Lookup Latency

A Kademlia DHT with N nodes converges in O(log₂ N) hops. For a registry with 100 million entries distributed across 100,000 nodes:

- log₂(100,000) ≈ 17 hops
- Each hop: one network round trip, estimated 20–50ms for geographically distributed nodes
- **Estimated lookup latency: 340–850ms** for a cold lookup (no cache)

This is the worst case. In practice, two factors reduce effective lookup latency significantly:

**Local caching:** The client caches public keys for all recent and frequent contacts. For a typical user who communicates with a stable set of contacts, the majority of lookups are served from cache. Cache hit rates of 90%+ are realistic for established users, reducing the average lookup latency to tens of milliseconds.

**Relay node caching:** Relay nodes cache recently looked-up keys for their served users. A relay node serving 10,000 users will see significant lookup overlap and can serve cached keys for the majority of inter-user messages without a DHT query.

**Effective average lookup latency estimate: 30–100ms** accounting for realistic cache hit rates.

#### 16.3.2 Registry Throughput

A single DHT node handling registry lookups must process:

- Lookup requests from relay nodes and clients in its geographic region
- Routing table maintenance messages (Kademlia PING, FIND_NODE)
- Registry updates (new registrations, key rotations, revocations)

For a 100,000-node registry with uniform load distribution and 500 million registered identities, each node is responsible for approximately 5,000 identity records. At peak global messaging load (600,000 messages/second), assuming a 10% cache miss rate, the DHT must process approximately 60,000 lookups/second across all nodes — approximately 0.6 lookups/second per node on average, with significant geographic skew toward high-density regions.

This is well within the throughput capacity of modern server hardware. Kademlia DHT implementations routinely handle thousands of operations per second per node. Registry throughput is **not a scalability bottleneck** under the target load.

#### 16.3.3 Registry Storage

Each identity record is approximately 500 bytes (public keys, address string, metadata, signature). For 500 million registered identities:

- Total registry data: ~250GB
- Per node (100,000 nodes, with 3× replication): ~7.5MB

This is negligible. Even with 10× replication for reliability, the per-node storage requirement is under 75MB. Registry storage is **not a scalability bottleneck**.

---

### 16.4 Message Relay Throughput

#### 16.4.1 Per-Node Throughput

A relay node's primary work per message is:

1. Receive and parse the onion packet (~10 µs CPU)
2. Verify the layer signature (~40 µs CPU)
3. Decrypt the onion layer (~10 µs CPU)
4. Forward to next hop or store (~network I/O)

Total CPU time per relayed message: approximately **60 µs**, or roughly **16,000 messages/second** on a single CPU core. A relay node running on a 4-core server can process approximately **64,000 messages/second** on CPU alone, with the practical limit being network I/O bandwidth.

At 50KB per message (the target average encrypted envelope size), 64,000 messages/second requires approximately **3.2 GB/s** of network throughput — within reach of a 40GbE network interface, but requiring dedicated hardware for a heavily loaded node.

In practice, relay nodes will not be uniformly loaded. A network of 10,000 relay nodes handling 600,000 messages/second distributes to an average of **60 messages/second per node** — well within the capacity of commodity hardware. Peak load on the highest-traffic nodes (in major metropolitan areas, handling traffic for large concentrations of users) may be 10–50× average, or 600–3,000 messages/second per node, still comfortably within the capacity of a single server.

**Relay throughput is not a scalability bottleneck** at Year 5 target scale.

#### 16.4.2 End-to-End Message Latency

The end-to-end latency for a DMCN message from send to delivery for an online recipient is the sum of:

| Component | Estimate | Notes |
|---|---|---|
| Sender cryptographic operations | ~200 µs | See Section 16.2 |
| Hop 1 network latency | 20–100ms | Sender to first relay node |
| Hop 1 relay processing | ~60 µs | Signature verify + route |
| Hop 2 network latency | 20–100ms | Relay to relay |
| Hop 2 relay processing | ~60 µs | |
| Hop 3 network latency | 20–100ms | Relay to delivery node |
| Hop 3 relay processing | ~60 µs | |
| Recipient fetch + decrypt | ~200 µs + network | Recipient polling interval dependent |
| **Total (online recipient, optimistic)** | **~100ms** | Three nearby hops |
| **Total (online recipient, pessimistic)** | **~500ms** | Three geographically dispersed hops |

For offline recipients, delivery latency is determined by the recipient's client polling interval. A default polling interval of 60 seconds gives a maximum additional latency of 60 seconds for offline delivery, comparable to standard email delivery expectations.

This latency profile is **comparable to existing encrypted messaging systems** (Signal, WhatsApp) and substantially better than traditional email, which has no guaranteed delivery time.

---

### 16.5 Storage Requirements

#### 16.5.1 Relay Node Message Storage

Relay nodes buffer messages for offline recipients until they are fetched or until the retention period expires (default: 30 days). The storage requirement per relay node depends on the number of users served and their average message volume.

Assumptions:
- Average user receives 50 messages/day at 50KB each = 2.5MB/day
- 30-day retention = 75MB per user
- A relay node serving 10,000 users = **750GB** of message storage

This is a substantial but entirely manageable storage requirement for a dedicated server. Consumer NVMe storage at 750GB costs under $100; cloud object storage at equivalent capacity is approximately $15/month at current prices. Relay node operators serving larger user populations will need proportionally more storage, but the per-user cost remains low.

#### 16.5.2 Client Storage

Client-side message storage is bounded by the user's device storage and retention preferences. The DMCN client should implement configurable local retention with automatic archival to encrypted cloud backup, consistent with the behaviour of modern email clients.

---

### 16.6 Network Bandwidth

#### 16.6.1 Onion Routing Overhead

Each message traverses 3 relay hops rather than the 1–2 hops typical in SMTP delivery. The bandwidth cost of each additional hop is one additional transmission of the encrypted message across the network. For a 50KB message traversing 3 hops, the total network bandwidth consumed is approximately **150KB** (3 × 50KB), compared to approximately **50–100KB** for a typical SMTP delivery.

The onion routing overhead therefore increases total network bandwidth consumption by approximately **1.5–3×** relative to direct delivery. This is the privacy cost of the onion routing layer and is the correct trade-off given the privacy benefits described in Section 18.2.

At Year 5 target scale (50 billion messages/day at 50KB each with 3× onion overhead), the total daily network bandwidth consumption of the DMCN is approximately **7.5 petabytes/day**. This is a large but entirely tractable figure — the global internet carries approximately **500 exabytes/day** of traffic, and global email traffic already accounts for a significant fraction of that.

#### 16.6.2 Size Class Padding Overhead

Message size class padding (Section 15.3.3) adds up to 3× overhead in the worst case (a 1KB message padded to the 4KB size class). For the average 50KB message, padding to the nearest size class (64KB) adds approximately 28% overhead. Across the full message volume, padding overhead is estimated at **15–30%** of total payload bandwidth.

This is a worthwhile privacy cost: size normalisation substantially reduces the inferential value of traffic analysis as described in Section 18.2.3.

---

### 16.7 Scalability Bottleneck Summary

| Component | Bottleneck Risk | Assessment | Primary Mitigation |
|---|---|---|---|
| Identity registry lookup latency | Medium | Cold lookups ~340–850ms; acceptable with caching | Client and relay node key caching |
| Identity registry throughput | Low | 0.6 lookups/sec/node average; well within capacity | Standard Kademlia implementation |
| Identity registry storage | Low | ~7.5MB per node at 500M users; negligible | Standard DHT replication |
| Relay node CPU throughput | Low | 64K msg/sec per 4-core node; average load ~60 msg/sec | Horizontal scaling |
| Relay node network I/O | Medium | 3.2 GB/s at max single-node load; requires 40GbE | High-traffic nodes require dedicated hardware |
| Relay node storage | Medium | 750GB per 10K users at 30-day retention | Tiered storage; shorter retention for high-volume nodes |
| End-to-end latency | Low | 100–500ms for online recipients; acceptable | Geographic relay node distribution |
| Total network bandwidth | Low | 7.5 PB/day at Year 5; tractable at global scale | Standard CDN and transit infrastructure |

The overall assessment is that the DMCN architecture is **viable at Year 5 adoption scale** without requiring novel infrastructure. The components that present the most engineering attention are identity registry lookup latency (addressed by caching strategy) and relay node storage management (addressed by retention policy and tiered storage). Neither represents a fundamental architectural obstacle.

At full global-scale deployment (4 billion users, 350 billion messages/day), relay node storage and network bandwidth requirements increase by approximately two orders of magnitude and would require infrastructure investment comparable to that of a major cloud provider. This is a longer-horizon engineering challenge that the architecture supports in principle but that is explicitly deferred as outside the scope of the prototype phase.

---


---
