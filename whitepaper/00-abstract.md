## Abstract

Email is the foundational communication layer of the digital world. Over
four billion people rely on it daily for personal correspondence,
commercial transaction, and institutional communication. Yet the
protocol underpinning it — the Simple Mail Transfer Protocol, first
standardised in 1982 — was designed for a network of a few hundred
trusted academic nodes and has not been architecturally updated to
reflect the adversarial reality of the modern internet.

The consequences are systemic. Spam accounts for the majority of all
email traffic by volume. Business Email Compromise costs organisations
billions of dollars annually. Phishing remains the most common vector
for credential theft and ransomware delivery. These are not edge cases
or unsolved technical puzzles — they are structural outcomes of a
protocol that provides no mechanism for verifying who is actually
sending a message.

Decades of mitigation work — Sender Policy Framework (SPF), DomainKeys
Identified Mail (DKIM), Domain-based Message Authentication Reporting
and Conformance (DMARC), and increasingly sophisticated AI-based
filtering — have reduced the visible impact of spam without
eliminating its root cause. That cause is the absence of cryptographic
identity at the protocol level. Any system that treats spam as a
filtering problem rather than an identity problem is treating symptoms,
not the disease.

This whitepaper proposes the Decentralized Mesh Communication Network
(DMCN): a next-generation messaging infrastructure designed from the
ground up to resolve the identity problem at the protocol level. The
DMCN's central proposition is that a message from an unverified sender
should not be deliverable, and a message from a verified sender should
not be forgeable. These two properties, enforced cryptographically at
the point of transmission rather than filtered probabilistically after
the fact, make spam and phishing structurally impossible rather than
merely statistically unlikely.

The DMCN achieves this through three interlocking architectural
commitments. First, every participant in the network is assigned a
cryptographic identity — a public/private key pair generated on their
device — that serves as their immutable network identifier. Second,
every message must carry a valid cryptographic signature from a
registered identity before it is accepted by any relay node; unsigned or
unverifiably signed messages are dropped at the network boundary, not
filtered at the inbox. Third, the network is peer-to-peer with no
central routing authority, eliminating both the single points of failure
that make centralised systems vulnerable to infrastructure attacks and
the centralised interception points that make them accessible to mass
surveillance.

The primary contributions of this investigation are:

- A structured analysis of the structural failures of SMTP that make spam, phishing, and identity fraud endemic, and a critique of why five decades of mitigation work have not resolved them

- A survey of the existing landscape of alternative approaches — including PGP, S/MIME, blockchain-based messaging, and federated encrypted email — and an articulation of the specific failure modes that have prevented each from achieving mainstream adoption

- A proposed technical architecture for the DMCN, covering the identity layer, transport layer, and storage and delivery layer, with particular attention to the message lifecycle and the mechanism by which spam is rejected at the protocol level

- A detailed treatment of the SMTP-DMCN bridge architecture that enables gradual, real-world adoption by maintaining interoperability with the four billion users on legacy email infrastructure during the transition period

- An analysis of address portability — the mechanism by which existing email addresses can be brought into the DMCN identity layer without requiring users to abandon their established address — and its implications for adoption and spam resistance

- A layered trust management framework covering allowlists, the pending queue, shared reputation feeds, and the cryptographic blocklisting model that makes identity reputation permanent and non-transferable

- A threat model covering eight adversary categories, comparing DMCN's threat surface against SMTP and providing an honest assessment of which threats are eliminated, which are mitigated, and which remain open challenges

A recurring theme across these contributions is the distinction between
technical capability and user experience. The cryptographic technology
required to implement the DMCN — elliptic curve key pairs, distributed
hash tables, onion routing — is mature and well-understood. The
barrier to adoption of prior cryptographic email systems was not the
technology; it was the imposition of cryptographic complexity on users
who had no interest in managing keys, certificates, or seed phrases. The
DMCN is designed explicitly to hide this complexity below the user
experience layer. The security model operates invisibly; users interact
with familiar addressing formats, contact lists, and message threads.

This design philosophy draws on a significant and underappreciated
precedent: Apple Passkeys and Google Passkeys have demonstrated that
hundreds of millions of people can use public/private key cryptography
daily without knowing they are doing so. The DMCN applies this same
principle to messaging — cryptographic identity as infrastructure, not
as a user-facing feature.

This whitepaper is presented as Version 0.2 of an ongoing investigation.
It is a research agenda and design direction, not a finished
specification. Several significant open challenges are documented ---
including Sybil resistance, key recovery without central authority,
regulatory compliance for encrypted communications, and network
governance — that will require further research, prototyping, and
community engagement before a production system can be specified.

The opportunity the DMCN addresses is real and urgent. The competitive
landscape is beginning to recognise the problem space, but no existing
solution combines the cryptographic soundness of PGP, the decentralised
architecture of peer-to-peer networks, the user experience standards of
modern consumer applications, and a credible migration path from the
existing email ecosystem. This whitepaper is a proposal toward that
combination.

> **Central Thesis**
> *Spam and email fraud are not filtering problems — they are
> identity problems. A mesh network where every node is
> cryptographically identified and every message is cryptographically
> signed eliminates the conditions under which spam is economically
> viable, rather than attempting to detect and discard it after the
> fact. This is the only class of solution that addresses the root
> cause.*

---


---

