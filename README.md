# Decentralized Mesh Communication Network (DMCN)

**Rebuilding digital communication on a foundation of trust.**

DMCN is a proposed next-generation messaging infrastructure that replaces email's broken identity model with cryptographic sender verification at the protocol level. Every participant has a cryptographic identity. Every message carries a verifiable signature. Unsigned or unverifiably signed messages are rejected at the network boundary — not filtered after the fact.

The result: spam and phishing become structurally impossible rather than statistically unlikely.

---

## The Problem

Email was designed in 1982 for a network of a few hundred trusted academic nodes. It has not been architecturally updated since. The consequences are systemic:

- Spam accounts for the majority of all global email traffic by volume
- Business Email Compromise (BEC) costs organisations billions annually
- Phishing remains the leading vector for ransomware and credential theft

Five decades of mitigations — SPF, DKIM, DMARC, AI-based filtering — have treated these as filtering problems. They are identity problems. No filtering system can permanently resolve a problem embedded in the protocol architecture.

---

## The Solution

DMCN is built on three interlocking commitments:

**Cryptographic identity.** Every participant has a public/private key pair generated on their device. Identity is self-sovereign — not assigned by any central authority and not revocable by any third party.

**Mandatory sender verification.** Every message must carry a valid cryptographic signature from a registered identity before any relay node will accept it. Verification is a gate at transmission, not a filter at the inbox.

**Peer-to-peer routing.** No central routing authority. No single point of failure. No centralised interception point. Messages are relayed through a distributed mesh with onion-routing-inspired metadata privacy.

All cryptographic complexity is invisible to end users. Onboarding is comparable in friction to creating a Gmail account. Users interact with familiar `user@domain` addresses and never encounter the words "key", "signature", or "certificate" in normal operation. The precedent for this is already established at scale: Apple and Google Passkeys have demonstrated that hundreds of millions of people can use elliptic curve cryptography daily without knowing they are doing so.

A bridge architecture maintains full interoperability with the 4 billion users on legacy SMTP email throughout the transition period.

---

## Status

> **Early proof-of-concept in progress.** The whitepaper is complete. Implementation has begun.

| Milestone | Description | Status |
|---|---|---|
| **Whitepaper v0.2** | Full architecture specification | ✅ Complete |
| **M1 — Cryptographic Core** | Identity layer, message format, KEM encryption (`internal/core`) | 🔨 In progress |
| **M2 — Node and Registry** | libp2p-based DHT, relay node, `dmcn-node` binary | ⏳ Planned |
| **M3 — Bridge Node** | Inbound SMTP→DMCN and outbound DMCN→SMTP paths | ⏳ Planned |

The implementation is being developed in Go. The cryptographic core is the starting point, deliberately prior to any networking layer.

---

## Documentation

The full technical whitepaper is the primary design document for the project.

- [**Read the whitepaper**](docs/whitepaper/README.md) — architecture, protocol specification, threat model, privacy analysis, and open challenges
- [**Download the PDF**](docs/whitepaper/DMCN-Whitepaper-v0.2.pdf)

The whitepaper covers:

- Structural analysis of SMTP's identity failure and why five decades of mitigations have not resolved it
- Survey of prior approaches (PGP, S/MIME, blockchain-based messaging) and why each failed
- Full protocol specification: identity layer, KEM message format, onion routing transport, storage and delivery
- SMTP-DMCN bridge architecture for legacy interoperability
- Domain Authority Record model for organisational address management
- Trust management framework: allowlists, pending queue, blocklists, shared reputation feeds, guardian controls
- Threat model across eight adversary categories
- Privacy analysis and regulatory compliance
- Performance and scalability analysis at global email scale
- Open challenges and research questions

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                   DMCN Network                       │
│                                                      │
│  ┌──────────────┐     ┌──────────────┐              │
│  │ Identity     │     │ Transport    │              │
│  │ Layer (DHT)  │────▶│ Layer        │              │
│  │              │     │ (Onion       │              │
│  │ Ed25519 keys │     │  Routing)    │              │
│  │ X25519 KEM   │     │              │              │
│  └──────────────┘     └──────┬───────┘              │
│                              │                       │
│                       ┌──────▼───────┐              │
│                       │ Storage &    │              │
│                       │ Delivery     │              │
│                       │ (Encrypted   │              │
│                       │  at rest)    │              │
│                       └──────────────┘              │
└─────────────────────────┬───────────────────────────┘
                          │
                   ┌──────▼───────┐
                   │  SMTP Bridge │  ◀── Legacy email
                   │  (Transition │       interop
                   │   layer)     │
                   └──────────────┘
```

Key cryptographic primitives: **Ed25519** (signing), **X25519 + HKDF-SHA256** (key exchange), **AES-256-GCM** (encryption), **KEM pattern per RFC 9180** (multi-device envelope).

---

## Getting Started

The proof-of-concept implementation is in early development. The cryptographic core (`internal/core`) is the current focus and the right place to start for contributors.

**Prerequisites**

- Go 1.22 or later

**Clone the repository**

```bash
git clone https://github.com/mertenvg/dmcn.git
cd dmcn
```

Further setup instructions will be added as the M1 milestone progresses.

---

## Contributing

Feedback, critique, and contributions are actively welcomed — including from cryptographers who spot weaknesses in the design, engineers who want to contribute to the implementation, and practitioners in regulated industries with compliance requirements the architecture should address.

The best way to engage:

- **Open an issue** for design questions, bug reports, or proposals
- **Start a discussion** in GitHub Discussions for broader architectural questions
- **Submit a pull request** for implementation contributions once M1 is underway

Please read the whitepaper before opening a design-level issue — most architectural decisions are documented with their rationale, and understanding the reasoning makes for more productive discussion.

---

## Licence

This project is licensed under the **Apache License 2.0** — see [LICENSE](LICENSE) for details.

Apache 2.0 was chosen deliberately: it includes an explicit patent non-aggression clause, meaning contributors grant users a licence to any patents covering their contributions. This is consistent with the project's intent to remain an open protocol.

---

## Citation

If you reference this work, please cite it as:

```
van Gerven, M. (2026). Decentralized Mesh Communication Network:
Rebuilding Digital Communication on a Foundation of Trust (Version 0.2).
https://github.com/mertenvg/dmcn
```

---

*DMCN is a research agenda and design direction, not a finished specification. The open challenges are real and documented. Feedback and collaboration are what move it forward.*
