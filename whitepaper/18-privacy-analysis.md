## 18. Privacy Analysis

This section addresses a question distinct from the threat model in Section 13: not whether the DMCN can be attacked, but what the system *inherently reveals* during normal, correct operation. A communication network can be cryptographically secure against active attackers while still exposing significant information about its users through the ordinary mechanics of message routing, identity discovery, and protocol operation.

The privacy analysis is structured around four areas: metadata exposure at the network layer, the identity registry as a surveillance surface, bridge node privacy, and regulatory compliance in a decentralised architecture. Each area is assessed against a baseline of what the current SMTP email ecosystem reveals, so that the comparison is grounded rather than abstract.

> **Scope**
> *This analysis addresses privacy in the technical sense — what information is exposed to which parties as a structural consequence of the protocol — rather than the policy sense of what operators choose to do with data. Operator conduct is a governance and regulatory matter addressed in Section 19.3 and Section 18.4.*

---

### 18.1 Baseline: What SMTP Reveals

Before assessing the DMCN, it is worth being precise about the privacy properties of the system it proposes to replace. SMTP email, as deployed by major providers, exposes the following to varying parties:

**To the email provider (Gmail, Outlook, etc.):** Full message content, subject lines, sender and recipient addresses, timestamps, device metadata, IP addresses, and — through scanning for features like Smart Reply and spam classification — inferred behavioural patterns and social graphs. Major providers operate under privacy policies that permit substantial use of this data for advertising and product improvement, subject to jurisdiction-specific regulatory constraints.

**To relay infrastructure in transit:** Historically, SMTP transmitted message content in plaintext. Opportunistic TLS between mail transfer agents, now widely deployed, encrypts content in transit between servers — but each server in the relay chain can read content, and TLS is not universally enforced. Message headers, including sender, recipient, routing path, and timestamps, are structurally visible to all relay infrastructure.

**To passive network observers:** On links where TLS is not enforced, full message content is visible. Even with TLS, connection metadata — which servers are communicating, at what times, with what data volumes — is observable at the network layer.

**To recipients:** Full message content, the sender's email address, and whatever metadata the sending client and relay chain have appended to message headers.

This is the baseline against which the DMCN's privacy properties should be measured. The bar is not high.

---

### 18.2 Metadata Exposure at the Network Layer

End-to-end encryption protects message *content* from relay nodes — no node in the DMCN transport layer can read what Alice is saying to Bob. What encryption does not protect is *metadata*: the fact that Alice is communicating with Bob, the frequency of their exchanges, the timing of messages, and approximate message sizes. This metadata can be as revealing as content in many threat models.

#### 18.2.1 What Relay Nodes Can Observe

A DMCN relay node handling a message in transit can observe the following:

- The sender's public key (as the message's cryptographic identifier)
- The recipient's public key (as the routing address)
- The approximate size of the encrypted payload
- The timestamp of receipt and forwarding
- The IP address of the upstream node that delivered the message, and the IP address of the downstream node to which it is forwarded

A relay node cannot observe the message content, subject, or any human-readable metadata. It also cannot — in a correctly implemented onion routing scheme — observe both the originating sender's IP address and the ultimate recipient's IP address simultaneously. Each relay node sees only the previous and next hop in the delivery chain.

This is a material improvement over SMTP, where relay nodes can read full message content and headers. However, it is not equivalent to anonymity. A relay node that handles a high volume of traffic for a small number of users can build a detailed picture of communication patterns between pseudonymous identities (public keys) even without reading content.

#### 18.2.2 What a Global Passive Observer Can Infer

The most sophisticated metadata threat is a global passive adversary — an entity capable of observing a significant fraction of all network traffic simultaneously. This is the same threat that onion routing networks such as Tor are known to be vulnerable to through traffic correlation attacks.

By observing that a message-sized packet left Alice's IP address at time T, and that an equivalently-sized packet arrived at Bob's relay node shortly after T, a global passive observer can probabilistically correlate the two events and infer that Alice sent a message to Bob — even without reading either packet's content.

The DMCN's onion routing layer partially mitigates this by introducing multiple hops and variable timing, increasing the difficulty of correlation. It does not eliminate it. For users whose threat model includes nation-state-level traffic analysis — journalists communicating with sources, activists in authoritarian jurisdictions, legal counsel in sensitive matters — the DMCN should be understood as providing strong content privacy with meaningful but imperfect metadata privacy. Users with these threat models should be directed to Tor-over-DMCN configurations or equivalent overlay networks for the transport layer.

For the vast majority of users whose threat model does not include a global passive adversary, the DMCN's metadata privacy properties represent a substantial improvement over SMTP.

#### 18.2.3 Message Size as a Side Channel

Encrypted message sizes are observable by relay nodes and passive network observers even when content is not. In some contexts, message size is itself informative — a 50KB encrypted message is more likely to contain a document attachment than a brief reply. The DMCN transport layer should implement padding to normalise message sizes into a small number of size classes, reducing the inferential value of size observation. This is a standard technique in privacy-preserving transport protocols and is a recommended design requirement, though its implementation detail is deferred to the protocol specification.

---

### 18.3 The Identity Registry as a Surveillance Surface

The DMCN's distributed identity registry is public by architectural necessity. For the system to function — for any sender to be able to look up a recipient's public key and send them an encrypted message — the registry must be queryable by anyone. This openness is a deliberate design choice, but it creates a surveillance surface that requires explicit attention.

#### 18.3.1 Account Existence Confirmation

Anyone with access to the identity registry can query it to determine whether a given email address has been registered as a DMCN identity. This means:

- An advertiser, stalker, data broker, or intelligence agency can compile lists of DMCN users by querying the registry against known email addresses.
- The registry reveals not just public keys but the fact of account existence, which may itself be sensitive in some contexts.
- Bulk enumeration of the registry — attempting to discover all registered identities — is a privacy concern if not rate-limited.

**Mitigation:** The registry should implement rate limiting on lookups per source IP or per authenticated identity, and should not support bulk enumeration queries. Lookups should return only the specific queried identity, not adjacencies or related entries. The registry design should also consider whether to support unlisted identities — accounts that are reachable by existing contacts but do not appear in registry searches initiated by unknown parties.

#### 18.3.2 Social Graph Inference from the Registry

Because registry entries include the identity's public key and any cross-signatures from the web of trust, a determined observer who accumulates registry data over time can begin to map social relationships: if Alice's key is cross-signed by Bob's and Carol's keys, it is inferable that Alice, Bob, and Carol have a trust relationship. At scale, the web-of-trust attestation data constitutes a partial social graph that is structurally visible without reading any message content.

**Mitigation:** Web-of-trust attestations should be opt-in for public visibility. Users should be able to maintain private attestations — stored locally or exchanged out-of-band — that inform their own trust decisions without being published to the global registry. The registry specification should distinguish between public attestations (visible to all) and private attestations (held locally or shared only with specific contacts).

#### 18.3.3 Timing and Correlation via Registry Lookups

When Alice's client performs a registry lookup for Bob's public key, that lookup is itself observable as a network event. A network observer monitoring Alice's traffic can infer that Alice is about to send a message to Bob — before the message is even sent — simply by observing the registry query.

**Mitigation:** The client should implement a registry prefetching strategy, maintaining a local cache of public keys for recent and likely contacts and refreshing it on a schedule rather than on demand. This decouples the timing of registry lookups from the timing of message composition, reducing the inferential value of lookup timing.

---

### 18.4 Bridge Node Privacy

The SMTP-DMCN bridge architecture, addressed in Section 11 from a security perspective, has distinct privacy implications that require separate treatment.

#### 18.4.1 Outbound Path: What the Bridge Operator Sees

When a DMCN user sends a message to a legacy email address, the message must be decrypted at the bridge node to be re-encoded as SMTP. This is an unavoidable consequence of protocol translation, disclosed in Section 11.2.2. The privacy implication is explicit: the bridge operator has technical access to:

- The full content of every outbound message sent to legacy email recipients
- The sender's DMCN identity
- The recipient's legacy email address
- Timestamps and message sizes

This is structurally equivalent to the trust placed in a conventional email service provider such as Gmail or Outlook, which has identical access to message content. The difference is that users choosing the DMCN are typically doing so with an expectation of enhanced privacy — and the bridge's necessary content access may conflict with that expectation if not clearly disclosed at onboarding.

**Disclosure requirement:** The DMCN client must present a clear, non-technical disclosure at the point where a user first sends a message to a legacy email recipient, explaining that the bridge operator can read the content of messages sent to non-DMCN addresses. This disclosure should be persistent — not a one-time consent flow that users will click through without reading — and the privacy policy of the chosen bridge operator should be linked and surfaced in the client UI.

**Mitigation through operator choice:** Because the bridge architecture is federated (Section 11.5), users can choose bridge operators with strong privacy commitments, including operators that commit to zero message logging and are subject to independent audit. Organisations with strong confidentiality requirements can operate their own bridge nodes, eliminating third-party access entirely.

#### 18.4.2 Inbound Path: Legacy Sender Metadata

When a legacy email sender sends a message to a DMCN user's bridge address, the bridge operator observes the full SMTP headers of the inbound message: sender address, sending server IP, timestamps, and routing path. This metadata is used to perform the authentication classification described in Section 11.3.2 and is necessarily retained for that purpose.

The DMCN specification should define minimum and maximum retention periods for bridge-held metadata, consistent with applicable data protection law, and should require bridge operators to publish their metadata retention policies.

#### 18.4.3 Bridge Operator as Data Controller

Under the EU General Data Protection Regulation and equivalent frameworks, any entity that determines the purposes and means of processing personal data is a data controller with obligations to data subjects. A bridge operator processing message content and metadata on behalf of DMCN users is a data controller for that processing.

This has practical implications: bridge operators must have a lawful basis for processing, must respond to data subject access requests, must implement appropriate technical and organisational security measures, and must be located in or have adequate data transfer arrangements with the jurisdictions in which their users reside. The DMCN Bridge Operator Protocol (BOP) should incorporate these requirements as conditions of registry participation, so that users can trust that any registered bridge operator meets a minimum compliance baseline.

#### 18.4.4 IMAP and POP3 Bridge Privacy

The local IMAP bridge model described in Section 10.6 introduces a privacy surface distinct from both the native DMCN path and the SMTP bridge path, and its privacy properties depend critically on where the bridge process runs.

**Local bridge (runs on user's device):** The private key remains in the device's secure enclave. Decryption occurs in the local bridge process; cleartext is present only in memory on the user's own hardware and is never transmitted to any third party. The IMAP connection is to localhost only. This model preserves end-to-end encryption end-to-user — the privacy properties are equivalent to the native DMCN client for inbound messages. The local bridge process itself is a potential attack surface (a malicious or compromised bridge binary could exfiltrate decrypted content), which makes the integrity of the bridge software supply chain a security requirement. Bridge software should be open source, reproducibly built, and distributed through verifiable channels.

**Server-side bridge (runs on a third-party server):** Decryption necessarily occurs on the server, meaning the bridge operator has access to message plaintext. The privacy model degrades to that of a conventional hosted email provider — better than unauthenticated SMTP but short of true end-to-end encryption. Server-side IMAP bridges should be treated as equivalent to SMTP bridge nodes for privacy disclosure purposes: users must be clearly informed that the bridge operator can read their messages, and the operator must meet the same data controller obligations described in Section 18.4.3.

**POP3 considerations:** POP3 is a simpler protocol than IMAP — it downloads messages to the client and (by default) deletes them from the server. A local POP3 bridge is viable under the same security model as the local IMAP bridge. However, POP3's lack of server-side folder synchronisation makes it unsuitable for multi-device use, which limits its practical relevance as a transition mechanism. IMAP is the recommended bridge protocol for enterprise deployments.

**Recommendation:** The DMCN specification should mandate that any local bridge implementation stores no cleartext on disk — decrypted message content should be held in memory only and written to local storage exclusively in re-encrypted form using a key derived from the user's authentication credential. This prevents cleartext message recovery from disk forensics on a seized or stolen device.

---

### 18.5 Regulatory Privacy Compliance in a Decentralised Architecture

Decentralisation creates genuine tension with data protection frameworks that were designed around the assumption of an identifiable data controller. The DMCN's privacy architecture must honestly engage with this tension rather than treating decentralisation as a compliance shield.

#### 18.5.1 The Data Controller Problem

In the DMCN's native peer-to-peer layer, there is no central operator. Messages are stored encrypted on relay nodes, routed through the mesh, and held until the recipient retrieves them. No single entity controls the processing of any given user's messages. This makes it difficult to identify a data controller in the GDPR sense — and without a data controller, data subjects' rights (access, rectification, erasure, portability) become difficult to exercise.

**Practical positions:**

- For the core protocol layer, the user themselves may be considered the data controller for their own encrypted data, since only they hold the decryption key. This is analogous to the position taken by some self-hosted encrypted services.
- For relay nodes storing encrypted messages, the relay node operator may be considered a data processor acting on behalf of the user-controller, with a data processing agreement required.
- For bridge nodes, as discussed in Section 18.4.3, the operator is a data controller in their own right for the content they can access.
- For the identity registry, the distributed architecture means there is no single controller; each node operator is a processor of the subset of registry data they hold.

These positions are not fully settled in law and will require engagement with data protection authorities in relevant jurisdictions as the DMCN matures. The governance framework (Section 19.4) should include a dedicated working group on regulatory compliance.

#### 18.5.2 The Right to Erasure

GDPR Article 17 grants data subjects the right to erasure of their personal data. In the DMCN, this creates a specific challenge: encrypted messages stored on relay nodes cannot be deleted by the user on demand, because the user has no direct administrative relationship with relay node operators.

**Partial mitigation:** Messages stored on relay nodes are encrypted with the recipient's public key. If the recipient deletes their private key — or if the relay node's retention policy expires the message — the encrypted data becomes permanently inaccessible even if the bytes persist on disk. This achieves functional erasure (the data is unrecoverable) even without literal deletion of the stored bytes.

The DMCN specification should define a maximum message retention period for relay nodes, after which stored messages are deleted regardless of whether they have been retrieved. A default of 30 days with user-configurable extension is a reasonable starting point, consistent with practices in existing encrypted messaging systems.

#### 18.5.3 Data Portability

GDPR Article 20 grants data subjects the right to receive their personal data in a structured, machine-readable format and to transmit it to another controller. In practice, for a messaging system, this means the user's message history, contact list, and trust relationships.

The DMCN client should implement a full data export function that produces a portable, encrypted archive of the user's message history, whitelist, greylist, blacklist, and trust attestations in a documented, open format. This export serves both the regulatory compliance function and the practical function of enabling migration between DMCN client applications without loss of data.

---

### 18.6 Privacy Properties Summary

The table below summarises the DMCN's privacy properties across key dimensions, compared to the SMTP baseline.

| Privacy Dimension | SMTP (Gmail/Outlook) | DMCN Native | DMCN via SMTP Bridge | DMCN via Local IMAP Bridge |
|---|---|---|---|---|
| Message content visible to provider | Yes — always | No — E2EE | Yes — to bridge operator | No — decryption on device |
| Message content visible to relay infrastructure | Yes — historically; TLS in transit only | No — E2EE throughout | Partial — decrypted at bridge | No — E2EE to local bridge |
| Sender/recipient visible to relay nodes | Yes — full headers | Pseudonymous public keys only | Yes — to bridge operator | Pseudonymous public keys only |
| Metadata visible to passive network observer | Yes — sender/recipient, timing, size | Timing and size only (onion routing limits more) | Timing and size | Timing and size |
| Social graph inferable from infrastructure | Yes — from provider data | Partially — from registry attestations | Partially | Partially |
| Account existence discoverable | Yes — MX lookup | Yes — registry query | Yes | Yes |
| Data subject rights (GDPR etc.) | Provider is data controller; rights exercisable | Distributed; complex controller picture | Bridge operator is data controller | User is controller; local storage only |
| Message retention | Provider-controlled; typically indefinite | Relay node retention policy; finite | Bridge operator retention policy | User-controlled local storage |

---

### 18.7 Design Recommendations

The privacy analysis above yields the following concrete design recommendations for the DMCN specification, in priority order:

**Message size normalisation** — implement payload padding in the transport layer to reduce the inferential value of size observation. This is a low-cost, high-value privacy improvement.

**Registry rate limiting and unlisted identity support** — prevent bulk enumeration of the identity registry and allow users to opt out of public discoverability while remaining reachable by existing contacts.

**Private web-of-trust attestations** — allow trust attestations to be held locally rather than published to the global registry, preserving the utility of the web of trust without exposing social graph data.

**Registry lookup prefetching** — decouple registry lookups from message composition timing to reduce the inferential value of lookup timing to network observers.

**Bridge operator disclosure** — require persistent, prominent disclosure in the DMCN client when messages will be processed by a bridge operator, including a link to the operator's privacy policy.

**Relay node retention limits** — specify a maximum message retention period in the protocol, ensuring that unread encrypted messages do not persist indefinitely on relay infrastructure.

**Data export function** — implement a full, portable, encrypted data export in the DMCN client to satisfy data portability obligations and enable client migration.

**Regulatory working group** — establish a dedicated working group within the DMCN governance structure to engage with data protection authorities and develop jurisdiction-specific compliance guidance.

**Local IMAP bridge cleartext policy** — mandate that local IMAP and POP3 bridge implementations store no decrypted message content on disk; cleartext must be held in memory only and written to local storage exclusively in re-encrypted form, preventing recovery from disk forensics on a seized or stolen device.



---


---
