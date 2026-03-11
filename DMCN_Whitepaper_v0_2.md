---
title: "Decentralized Mesh Communication Network"
subtitle: "Rebuilding Digital Communication on a Foundation of Trust"
version: "0.2 — Integrated Draft"
date: "March 2026"
status: "CONFIDENTIAL — Draft for Review"
---

# Decentralized Mesh Communication Network

**Rebuilding Digital Communication on a Foundation of Trust**

Version 0.2 — Integrated Draft · March 2026 · *CONFIDENTIAL — Draft for Review*

---

## Table of Contents

- [Abstract](#abstract)
- [Executive Summary](#executive-summary)
- [1. The Problem with Email](#1-the-problem-with-email)
- [2. Why Existing Solutions Have Failed](#2-why-existing-solutions-have-failed)
- [3. Why This Hasn't Been Done Before](#3-why-this-hasnt-been-done-before)
- [4. The Competitive Landscape](#4-the-competitive-landscape)
- [5. The Full Problem Space: Beyond Spam](#5-the-full-problem-space-beyond-spam)
- [6. Proposed Architecture: Decentralized Mesh Communication Network](#6-proposed-architecture)
- [7. Cryptographic Identity and Key Management](#7-cryptographic-identity-and-key-management)
- [8. Spam Elimination at the Protocol Level](#8-spam-elimination-at-the-protocol-level)
- [9. User Experience: Hiding Complexity Without Sacrificing Security](#9-user-experience)
- [10. Transition Strategy: Coexistence with Legacy Email](#10-transition-strategy)
- [11. The SMTP-DMCN Bridge Architecture](#11-the-smtp-dmcn-bridge-architecture)
- [12. Bringing Existing Email Addresses to the DMCN](#12-bringing-existing-email-addresses-to-the-dmcn)
- [13. Trust Management: Allowlists, Pending Queue, and Blocklists](#13-trust-management)
- [14. Threat Model](#14-threat-model)
- [15. Open Challenges and Research Questions](#15-open-challenges-and-research-questions)
- [16. Conclusion](#16-conclusion)
- [17. Privacy Analysis](#17-privacy-analysis)
- [18. Protocol Specification Outline](#18-protocol-specification-outline)
- [19. Performance and Scalability Analysis](#19-performance-and-scalability-analysis)
- [Glossary](#glossary)
- [References](#references)

---


---

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

## Executive Summary

Email is broken. What began as a simple message-passing protocol for a
small trusted academic network has become the foundational communication
layer of global civilization — yet its architecture has not
fundamentally changed in over five decades. The result is a system where
the majority of all traffic is spam, where identity is trivially forged,
and where trust is an afterthought bolted on through a patchwork of
incomplete standards.

This whitepaper investigates the feasibility and design of a
Decentralized Mesh Communication Network (DMCN) — a next-generation
messaging infrastructure built from the ground up on cryptographic
identity, peer-to-peer routing, and public/private key verification. The
primary objectives of this investigation are:

- To document the structural failures of the existing email system that make spam, phishing, and identity fraud endemic rather than exceptional.

- To survey the landscape of existing solutions and articulate clearly why each has failed to achieve mainstream adoption.

- To propose a technical and architectural framework for a decentralized mesh network that eliminates spam at the protocol level through cryptographic sender verification.

- To address the key challenge of user experience — demonstrating that cryptographic identity need not create barriers to entry for mainstream users.

- To position this proposal relative to existing efforts such as Dmail Network, and articulate its distinct advantages.

The central thesis is straightforward: spam and email fraud are not
filtering problems — they are identity problems. Any solution that
addresses spam at the filtering layer rather than the identity layer is
treating symptoms rather than causes. A mesh network where every node is
cryptographically identified and every message is cryptographically
signed eliminates the conditions under which spam is economically
viable.


> **Core Proposition**
> *A message from an unverified sender is not delivered. A message from
> a verified sender cannot be forged. These two properties, enforced at
> the protocol level, make spam structurally impossible rather than
> merely inconvenient.*

---


---

## 1. The Problem with Email


### 1.1 A Protocol Designed for a Different World


The Simple Mail Transfer Protocol (SMTP), first defined in RFC 821 in
1982, was engineered for a network of a few hundred trusted nodes ---
universities, government research labs, and military institutions. In
that context, identity verification was unnecessary. Everyone on the
network was known. The protocol's defining characteristic — openness
--- was a feature, not a vulnerability.

Today, SMTP underpins communication for over 4 billion email users
globally. The network it was designed for no longer exists. The trust
assumptions it was built on are not merely strained — they are
completely inverted. The openness that made email powerful is precisely
what makes it exploitable.


### 1.2 The Scale of the Spam Problem


Spam is not a minor nuisance. It is the dominant form of email
communication on the planet by volume. Industry estimates consistently
place spam at between 45% and 85% of all global email traffic, with some
peak periods substantially higher. On any given day, hundreds of
billions of spam messages are transmitted across global email
infrastructure.

The consequences extend well beyond annoyance. Email-based phishing is
among the most costly forms of cybercrime. Business Email Compromise
(BEC) — a form of fraud in which attackers impersonate executives or
trusted partners — costs organizations billions of dollars annually.
The FBI's Internet Crime Complaint Center consistently ranks BEC among
the highest-impact cybercrime categories by financial loss.


> **Scale of the Problem**
> *Spam accounts for the majority of all email ever sent. Email fraud
> costs the global economy billions annually. These are not edge cases
> — they are the normal operating conditions of the current system.*


### 1.3 The Structural Root Cause


The spam problem is an identity problem. SMTP provides no mechanism for
a receiving server to verify that the sending server is who it claims to
be, and no mechanism for a recipient to verify that a message actually
came from the stated sender. Sending an email that appears to come from
any address — including your bank, your employer, or a government
agency — requires no special access, no credentials, and no technical
sophistication beyond a basic SMTP client.

This means that spam and phishing are not aberrations that a better
filter can eliminate. They are rational economic behaviors enabled by
the protocol itself. As long as sending a message claiming to be from
any sender costs essentially nothing and carries no verifiable identity,
the conditions for spam will persist regardless of how sophisticated
filtering becomes.


### 1.4 Existing Mitigations and Their Limitations


The email ecosystem has accumulated several layers of mitigation over
the decades, none of which address the root cause:

- Sender Policy Framework (SPF) — allows domain owners to specify which IP addresses are authorized to send email on their behalf. Widely adopted but easily circumvented through compromised authorized servers, and provides no per-message signing.

- DomainKeys Identified Mail (DKIM) — adds a cryptographic signature to outgoing messages, allowing receivers to verify that content was not altered in transit. Addresses integrity but not spam; a spammer controlling their own domain can produce valid DKIM signatures.

- Domain-based Message Authentication, Reporting and Conformance (DMARC) — a policy layer on top of SPF and DKIM. Adoption is inconsistent and enforcement is frequently weak.

- AI-based spam filtering — the approach used by major providers like Google and Microsoft. Highly effective at classifying known spam patterns, but reactive by nature, computationally expensive, and produces significant false positive rates that affect legitimate communication.

Each of these mitigations is a layer of additional complexity applied to
a fundamentally trust-less protocol. They reduce spam volumes in
practice, but they cannot eliminate spam in principle, because they do
not address sender identity at the protocol level.



---

## 2. Why Existing Solutions Have Failed


### 2.1 PGP and S/MIME — Technically Sound, Practically Abandoned


Pretty Good Privacy (PGP), created by Phil Zimmermann in 1991,
represented the first serious attempt to bring cryptographic identity to
email. Using a public/private key model and a decentralized web-of-trust
for key distribution, PGP was — and remains — technically capable of
providing exactly the sender verification and message encryption that
the email system lacks.

Yet after more than three decades, PGP adoption among general users
remains negligible. Key management required users to understand and
interact with cryptographic concepts that were entirely foreign to them;
there was no intuitive discovery mechanism for finding another person's
public key; and the user experience of most PGP implementations was
designed for technically sophisticated users. Crucially, both parties
needed to set it up — a network effect problem that was never
overcome.

S/MIME addressed some of these issues by integrating with certificate
authorities and building key management into email clients, but
reintroduced centralization through its dependence on a PKI hierarchy,
and added cost barriers through certificate fees. Enterprise adoption is
moderate; consumer adoption is essentially absent.


### 2.2 Protocol-Level Patches — Necessary but Insufficient


SPF, DKIM, and DMARC represent meaningful progress and have been widely
deployed. They have materially reduced certain categories of spoofing
and given domain owners tools for asserting sending policies. However,
they address domain-level authentication rather than individual sender
identity, and they are layered on top of a fundamentally unauthenticated
protocol rather than replacing it.


### 2.3 Blockchain-Based Approaches — Right Problem, Wrong Architecture


A number of projects have attempted to use blockchain technology as the
basis for a decentralized email identity and messaging layer. These
projects share the correct insight that decentralized cryptographic
identity is the solution to the trust problem — but their execution
has been constrained by dependence on blockchain infrastructure.

Blockchain-based systems introduce several structural disadvantages for
a general-purpose communication protocol: transaction latency makes
real-time messaging difficult; gas fees or transaction costs create
economic friction for every message; the requirement that users possess
and manage cryptocurrency wallets creates a significant barrier for
non-technical users; and the Web3 ecosystem orientation limits appeal to
a small subset of the potential user base.


### 2.4 The Network Effect Trap


Perhaps the most consistent failure mode across all alternative email
systems is the network effect problem. A secure communication system
provides zero value to its first user and limited value to any user
unless the people they want to communicate with are also on the system.
Solutions that require wholesale replacement of email face a
near-insurmountable adoption barrier. Any viable successor must provide
a credible migration path that allows incremental adoption.


> **Key Insight**
> *Previous solutions failed not because the cryptographic technology
> was inadequate, but because they did not solve the user experience
> problem, the key discovery problem, or the migration problem. A
> successful proposal must address all three.*


---


---

## 3. Why This Hasn't Been Done Before

The core technical proposition of the DMCN — cryptographic identity at the protocol level, peer-to-peer routing, spam rejection at the network boundary rather than the inbox — has been understood as the correct class of solution since at least the early 1990s. Phil Zimmermann understood it in 1991. The question a technically sophisticated reader will immediately ask is: if this is the right answer, why hasn't it been built?

The honest response is not that nobody thought of it. It is that five compounding obstacles have, until recently, made a successful execution implausible. Understanding those obstacles is important context for evaluating this proposal — and for understanding why the conditions have changed enough to make a serious attempt viable now.

### 3.1 The Network Effect Is Brutal

Email's value derives almost entirely from the fact that everyone is already on it. A replacement network starts at zero users, which means zero value, which means no rational user switches, which means it stays at zero users. This is not a technical problem — it is a cold-start problem of the kind that has ended every serious attempt at email replacement.

The only communication platforms that have broken through this barrier in recent years — Signal, WhatsApp, Telegram — did so by targeting a different use case (mobile messaging) rather than attempting a direct email replacement, and by benefiting from specific external catalysts: the Snowden revelations, Facebook's acquisition of WhatsApp, and government attempts to ban encrypted messaging. No equivalent catalyst has yet destabilised email incumbency at scale.

Any credible proposal for an email replacement must have a specific, credible answer to the cold-start problem. The DMCN's answer — the SMTP bridge architecture and address portability — is the design element that makes or breaks the proposal's real-world viability. The technical architecture is the easier part.

### 3.2 PGP Poisoned the Well

PGP did the cryptographically correct thing and failed so visibly, for so long, that it created a durable belief in the technical community that secure email is a structural impossibility. That belief suppressed serious engineering investment and venture capital attention for three decades.

The diagnosis was wrong. PGP failed on user experience and key discovery — not on cryptography. The underlying model was sound. But "if PGP couldn't do it, nothing can" became a default assumption that was rarely examined. The result is that the problem space has been underfunded and understaffed relative to its importance, and most serious cryptographic engineering talent has been directed toward problems perceived as more tractable.

### 3.3 Incumbent Platforms Have No Incentive to Fix It

Google and Microsoft collectively handle the majority of global email. Their business models are built on the email ecosystem as it currently exists. Gmail's advertising revenue depends on the ability to analyse message content. A cryptographically private email system with end-to-end encryption would directly undermine that model.

Both companies have invested substantially in spam filtering — which keeps users satisfied enough to prevent churn — but have no commercial incentive to solve the identity problem at the protocol level. Solving it would also solve their surveillance capability. The result is a rational corporate decision to treat spam as a manageable nuisance rather than a solvable structural problem, and to direct engineering resources accordingly.

This dynamic means that a solution to the email identity problem is unlikely to come from the incumbent platforms. It must come from outside.

### 3.4 The Transition Cost Is Asymmetric

Every prior attempt at secure or decentralised email forced users to make a binary choice: adopt the new system and abandon your existing address, or stay on SMTP. For most users and organisations, abandoning an established email address is not a realistic option — it is embedded in contracts, business cards, institutional systems, and years of correspondence. The switching cost is effectively prohibitive.

This asymmetry has killed technically sound proposals repeatedly. A system can be cryptographically superior in every measurable way and still achieve zero adoption if the migration path requires users to start over. The DMCN's address portability and bridge architecture are specifically designed to eliminate this barrier — but it is worth being clear that this design challenge is harder than the cryptographic design challenge, and that no previous attempt has solved it convincingly.

### 3.5 The UX Precondition Was Not Yet Met

The DMCN's central UX claim — that public/private key cryptography can be made entirely invisible to mainstream users — has only become demonstrably true in the last two to three years. The mass deployment of passkeys by Apple, Google, and Microsoft from 2022 onwards established, at a scale of hundreds of millions of users, that people can use elliptic curve cryptography daily without any awareness that they are doing so.

Before that precedent existed, the honest answer to "can you hide cryptographic key management from non-technical users at scale" was "we believe so, but we have not proven it." That uncertainty was a legitimate objection to the entire approach. It is no longer a legitimate objection. The passkey deployment provides a real-world existence proof at the exact scale the DMCN requires.

This is the most important change in the conditions since PGP's failure, and it is recent. It is the primary reason why the timing of this proposal differs from prior attempts.

### 3.6 Blockchain Absorbed the Decentralisation Impulse

From approximately 2017 onwards, the engineering energy and venture capital that might otherwise have been directed at decentralised communication infrastructure was largely absorbed by the blockchain and Web3 ecosystem. Projects that correctly identified the decentralised identity problem attempted to solve it using blockchain infrastructure — which introduced transaction latency, economic friction, and a requirement for cryptocurrency wallet ownership that made mainstream adoption structurally impossible.

The result was a decade in which the correct problem was identified by a large number of well-funded teams, and the wrong architectural choice was consistently made. The DMCN's explicit design decision to achieve decentralisation through peer-to-peer networking rather than blockchain infrastructure is a direct response to this pattern.

### 3.7 The Conditions Have Changed

The five obstacles above have not disappeared. The network effect problem remains the hardest unsolved challenge, and this proposal does not claim to have a guaranteed solution to it. What has changed is:

- The UX precondition is now met, demonstrably, at scale
- The cryptographic primitives required (Curve25519, Ed25519, AES-256-GCM) are mature, hardware-accelerated, and universally available
- The blockchain detour has produced useful lessons about what decentralised identity infrastructure should not be
- Regulatory pressure on email security and data privacy (GDPR, NIS2, increasing BEC enforcement) is creating institutional demand for a more trustworthy email substrate
- The distributed systems infrastructure required to run a global peer-to-peer network — cloud compute, global CDN, open-source DHT implementations — is commoditised in a way it was not in 2005 or even 2015

This is not an argument that the DMCN will succeed where others have failed. It is an argument that the conditions under which a well-designed attempt could succeed are better now than they have been at any prior point. The proposal stands or falls on the quality of its execution — particularly the migration strategy — not on the novelty of its core insight.

> **Historical Assessment**
> *The right technical answer to the email identity problem has been known for thirty years. The barriers to implementation have been predominantly economic, social, and experiential rather than cryptographic. The conditions that sustained those barriers are weaker now than at any previous point. That is the case for attempting this now, stated plainly.*




---

## 4. The Competitive Landscape


### 4.1 Dmail Network


Dmail Network (dmail.ai) is currently the most prominent active project
in the decentralized encrypted email space. It is an AI-enhanced
decentralized communication platform that combines blockchain-anchored
identity with encrypted messaging, cross-chain notifications, and Web3
productivity tools. The platform reports over 2 million monthly active
users and positions itself as the communication backbone of the Web3
ecosystem.

Dmail's technical approach uses Decentralized Identifiers (DIDs), NFT
domains, and cryptocurrency wallet addresses as the basis for user
identity. Messages are end-to-end encrypted and stored in decentralized
storage. The platform bridges to traditional email protocols, allowing
some interoperability with legacy systems.

While Dmail represents meaningful progress, its orientation toward the
Web3 and cryptocurrency ecosystem limits its potential for mainstream
adoption. Its identity model requires users to engage with blockchain
infrastructure, and its broader feature set — NFT domain trading,
token incentives, on-chain marketing tools — presupposes a user
already embedded in the crypto ecosystem. It is not designed as a
general replacement for email.


### 4.2 ProtonMail and Tutanota


ProtonMail and Tutanota represent the encrypted email approach within
the existing SMTP infrastructure. Both offer end-to-end encrypted email
between users on their own platform, with varying levels of encryption
when communicating with external addresses. They have achieved
meaningful consumer adoption, particularly among privacy-conscious
users. However, both are centralized services and neither addresses the
spam problem at the protocol level.


### 4.3 Signal and Matrix


Signal and Matrix (Element) demonstrate that decentralized,
cryptographically secure messaging is technically viable at scale.
Signal has achieved significant mainstream adoption while providing
state-of-the-art end-to-end encryption. Matrix provides a federated,
decentralized protocol with open-source infrastructure. Neither is
designed as an email replacement, but both provide important technical
and UX precedents that a DMCN design should draw upon.


### 4.4 Summary Comparison


| **Solution** | **Decentralized** | **Spam-Free by Design** | **No Crypto Required** | **Mainstream UX** | **Email Compatible** |
|---|---|---|---|---|---|
| PGP / S/MIME | Partial | No | No | No | Yes |
| Dmail Network | Yes | Partial | No | No | Partial |
| ProtonMail | No | No | Yes | Yes | Yes |
| Signal / Matrix | Yes | N/A | Yes | Yes | No |
| Proposed DMCN | Yes | Yes | Yes | Yes | Yes |


---


---

## 5. The Full Problem Space: Beyond Spam

The case for the DMCN presented so far has focused primarily on spam and phishing — the most visible symptoms of SMTP's identity problem. This framing, while accurate, understates the proposal's commercial and institutional relevance. SMTP's architectural failures generate a cluster of distinct, costly pain points that affect different constituencies in different ways. Understanding the full problem space serves two purposes: it strengthens the case for the infrastructure investment required to build the DMCN, and it identifies the multiple adoption wedges available to a deployment strategy that does not depend on convincing the entire global email user base simultaneously.

This section maps the complete landscape of SMTP pain points addressable by the DMCN, grouped by the constituency that bears the cost and has the motivation to act.

---

### 5.1 Email Deliverability: The Invisible Tax on Legitimate Senders

The whitepaper's framing of spam as a problem for *recipients* obscures an equally significant problem for *senders*. Legitimate organisations — SaaS companies, e-commerce platforms, financial services providers, healthcare systems — depend on transactional email for core business operations: account verification, password reset, order confirmation, appointment reminders, invoice delivery. These messages must reach their recipients reliably to function.

They frequently do not. Spam filters trained on the behaviour of billions of messages cannot reliably distinguish a legitimate transactional email from a spam campaign using the same infrastructure. The result is an entire industry — email deliverability consulting, reputation monitoring services, IP warm-up infrastructure, dedicated sending platforms such as SendGrid, Mailgun, and Amazon SES — that exists entirely to work around the fundamental untrustworthiness of SMTP identity.

The financial scale of this problem is substantial. Deliverability platforms represent a multi-billion dollar market. Engineering time spent on deliverability configuration, monitoring, and incident response is a significant cost centre for any organisation that depends on transactional email. A system in which sender identity is cryptographically verified at the protocol level makes deliverability infrastructure unnecessary: a verified sender's message is either explicitly accepted by the recipient or it is not — there is no probabilistic classification, no reputation score to manage, no warm-up period to endure.

For organisations that send transactional email at scale, the DMCN's deliverability guarantee is a direct and quantifiable cost saving, independent of any spam or phishing benefit. This constituency is already paying to manage a problem the DMCN eliminates.

> **Commercial Implication**
> *The deliverability market represents a paying customer base that already has budget allocated to solving a problem the DMCN solves structurally. Framing the DMCN as a deliverability solution — not just a security solution — significantly expands the commercial addressable market.*

---

### 5.2 Message Authenticity, Non-Repudiation, and Legal Admissibility

Every message sent through the DMCN bears a cryptographic signature that makes it tamper-evident and non-repudiable. The sender cannot credibly deny having sent a signed message; a recipient cannot credibly claim a message was altered after delivery. This property, which emerges as a side effect of the DMCN's identity layer, has significant standalone value in legal, financial, and compliance contexts.

#### 5.2.1 Litigation and E-Discovery

Email is the primary documentary record of business decision-making and is routinely produced in litigation, regulatory investigations, and employment disputes. Under current SMTP infrastructure, proving that a produced email has not been altered since it was sent requires a chain of custody argument supported by server logs, metadata analysis, and forensic examination — a complex, expensive, and imperfect process.

DMCN-signed messages are self-evidently authentic. The cryptographic signature is verifiable by any party with access to the sender's public key, without reference to server logs or forensic infrastructure. This reduces the cost and complexity of e-discovery authentication, and may be relevant to the admissibility standards for electronic evidence in jurisdictions that require authentication of digital records.

#### 5.2.2 Financial Services and Contractual Communications

Financial services firms operating under MiFID II, FINRA, and equivalent frameworks are required to record and archive communications related to investment advice and transactions. The authenticity of archived communications is a regulatory requirement, not an optional feature. DMCN signatures provide a technically robust authentication mechanism that satisfies this requirement structurally, rather than through the fragile chain-of-custody approaches currently employed.

Similarly, the increasing use of email for contractual communications — terms acceptance, instruction confirmation, amendment agreements — creates a need for message authenticity that current SMTP infrastructure cannot provide. A DMCN-signed instruction to transfer funds or execute a trade is cryptographically attributable to the sender's identity in a way that a standard email is not.

#### 5.2.3 The Non-Repudiation Premium

Non-repudiation — the inability of a sender to credibly deny having sent a message — is a property that organisations in regulated industries are effectively required to demonstrate but cannot currently achieve at the protocol level. This creates a compliance gap that the DMCN closes as a structural feature, not as an add-on product. Regulated industries represent early adopters with both the motivation and the budget to pay for this property.

---

### 5.3 Secure Document and Data Exchange

A significant fraction of business communication involves the exchange of sensitive documents: contracts, financial statements, medical records, legal filings, intellectual property. The current infrastructure for this exchange is a patchwork of inadequate solutions: email with password-protected attachments (security theatre), dedicated secure portals (friction-heavy, require recipient registration), consumer file-sharing services (insufficient audit trails), and specialised platforms such as DocuSign or Citrix ShareFile (purpose-built but siloed).

This fragmentation exists because email is fundamentally untrustworthy as a secure document transport layer. The DMCN changes this. End-to-end encryption and verified sender identity make DMCN a credible substrate for sensitive document exchange without requiring separate infrastructure. A document sent through the DMCN is encrypted to the recipient's public key (accessible only to them), signed by the sender's private key (authenticity guaranteed), and delivered through a network that provides no plaintext access to intermediate nodes.

For legal firms, healthcare providers, financial advisors, and any organisation that regularly exchanges sensitive documents with external parties, the DMCN eliminates the need to choose between the convenience of email and the security of a dedicated portal. This is a use case with clear commercial value and a constituency — professional services and regulated industries — that is already paying for inferior solutions.

---

### 5.4 Machine-to-Machine and Automated Communication

A large and growing fraction of email traffic is not human-generated. System alerts, CI/CD pipeline notifications, invoice processing, EDI (Electronic Data Interchange), B2B API event notifications, IoT device reporting — all of these use email or email-adjacent protocols as a communication substrate. This traffic class has no inherent need for a human inbox; what it requires is reliable, authenticated, machine-readable message delivery with a strong audit trail.

The DMCN's cryptographic identity layer is in many respects a better fit for machine identities than for human ones. Machines do not struggle with key management UX. A server that generates a key pair at provisioning time and registers it in the DMCN identity layer has, as a structural consequence, a cryptographic identity that can be used to sign and encrypt all outbound communications with zero ongoing management overhead.

This creates a compelling enterprise adoption path that does not depend on consumer behaviour at all. An organisation that deploys DMCN for its automated B2B communication — invoice delivery, order confirmation, API event notification — immediately solves its deliverability, authentication, and non-repudiation problems for that traffic class. It benefits from the DMCN's properties from day one, with no dependency on its counterparties adopting DMCN natively, because the bridge handles the translation to SMTP for recipients who have not yet adopted.

As counterparties adopt DMCN natively, the communication path upgrades automatically to fully encrypted and verified. The human-facing email experience improves as a trailing consequence of an adoption decision initially made for machine-to-machine efficiency.

---

### 5.5 Phishing Resistance as Cyber Insurance Risk Reduction

The insurance market for cyber risk has undergone significant hardening since 2020. Premiums have risen substantially, coverage terms have tightened, and underwriters are increasingly requiring evidence of specific security controls as a condition of coverage. Phishing resistance — specifically, the ability to demonstrate that email-based credential theft and BEC are structurally mitigated — is directly relevant to cyber insurance underwriting.

An organisation that has deployed DMCN for its internal and B2B communication has a technically defensible claim that a significant class of phishing attack is structurally impossible against its DMCN-protected identities. This is not a claim that can be made about spam filtering or security awareness training, both of which are probabilistic defences. It is a structural argument that an underwriter can evaluate and price.

The cyber insurance market is large, growing, and actively seeking ways to differentiate between organisations that have meaningfully reduced their risk exposure and those that have only deployed conventional defences. DMCN adoption may become a premium-reduction lever for organisations with significant cyber insurance exposure — creating a financial incentive for adoption that operates independently of the email experience itself.

---

### 5.6 Archival Integrity and Regulatory Compliance

Organisations subject to records retention requirements — which includes virtually every regulated business globally — maintain email archives as a matter of legal obligation. The integrity of those archives — the ability to demonstrate that archived messages have not been altered since receipt — is both a compliance requirement and a practical necessity for their use as evidence.

Current email archiving solutions rely on hash-based integrity verification of messages as they enter the archive system. This approach protects against post-archival modification but cannot authenticate the original message at the point of sending. A message that was altered before archiving, or that was forged in the first place, passes integrity checks if it enters the archive cleanly.

DMCN signatures provide origin authentication that archive integrity checks cannot. A signed message in the DMCN is verifiably attributable to its sender's cryptographic identity at the point of composition, not merely at the point of archiving. For organisations whose archives are subject to regulatory scrutiny or legal production, this is a qualitative improvement in the evidentiary value of their records.

---

### 5.7 Calendar, Scheduling, and Meeting Authenticity

Meeting invitations are a primary and growing phishing vector. Attackers send calendar invitations impersonating executives, financial counterparties, or IT support personnel to induce recipients to join fraudulent calls, reveal credentials, or authorise transactions. The current calendar and scheduling infrastructure — iCalendar, CalDAV — inherits all of SMTP's identity weaknesses. There is no reliable mechanism for verifying that a meeting invitation from an unknown sender is legitimate.

The DMCN's identity layer extends naturally to calendar and scheduling. A meeting invitation sent through the DMCN carries the same cryptographic identity guarantees as a message: it is signed by the sender's private key and verifiable against their registered public key. A verified meeting invitation from a known organisational identity is a qualitatively different artefact from the current model.

This is not a core feature of the DMCN specification but is a natural extension of the identity infrastructure, achievable without significant additional protocol work. Its inclusion in the DMCN's extension roadmap strengthens the value proposition for enterprise adoption.

---

### 5.8 Identity Infrastructure Beyond Email

The DMCN's cryptographic identity layer — a distributed registry of public keys anchored to human-readable addresses — is, at a structural level, general-purpose identity infrastructure. Email is the first application built on top of it, but the same infrastructure supports:

- **Website and service authentication** — a domain's DMCN identity can serve as a verifiable credential for web authentication, complementing or replacing certificate authority-based TLS in some contexts
- **Software supply chain signing** — packages, binaries, and configuration files signed with a DMCN identity provide the same non-repudiation guarantees as signed messages
- **API authentication** — service-to-service API calls authenticated with DMCN identities eliminate the management overhead of rotating API keys and secrets
- **Organisational identity attestation** — a DMCN organisational identity can attest to employee identities, creating a verifiable credential chain from organisation to individual without dependency on a central certificate authority

This reframes the DMCN from an email replacement to an identity infrastructure project with email as its initial and most visible application. The investment in deploying DMCN identity infrastructure yields returns across multiple use cases simultaneously, not just in the email context. For enterprise buyers, this changes the cost-benefit calculation significantly.

---

### 5.9 The Constituency Matrix

The pain points above map to distinct buying constituencies, each with independent motivation to adopt:

| Pain Point | Primary Constituency | Current Annual Cost Proxy | DMCN Treatment |
|---|---|---|---|
| Spam and phishing | All organisations | Billions in BEC losses; filtering infrastructure costs | Eliminated at protocol level |
| Deliverability | SaaS, e-commerce, financial services | Multi-billion dollar deliverability market | Eliminated — verified senders guaranteed delivery |
| Message authenticity | Legal, financial, regulated industries | E-discovery costs; compliance infrastructure | Structural — all messages cryptographically signed |
| Secure document exchange | Professional services, healthcare, finance | Dedicated portal market; portal management overhead | Replaced by native DMCN transport |
| Machine-to-machine comms | Enterprises with B2B automation | Engineering time on deliverability and auth | Structurally solved by machine identity layer |
| Cyber insurance risk | Any organisation with cyber coverage | Premium costs; underwriting requirements | Defensible structural risk reduction |
| Archival integrity | All regulated businesses | Archive infrastructure; forensic authentication costs | Structural — origin authentication at composition |
| Calendar phishing | All organisations | Incident response; fraud losses | Addressable through identity layer extension |
| General identity infrastructure | Enterprises | PKI management; API key rotation; supply chain signing | DMCN identity layer as shared substrate |

No single constituency needs to solve all of these problems to justify adoption. Each constituency needs to solve one of them — and solving one of them contributes to the network density that makes the next constituency's adoption more valuable.

This is the structural argument against the network effect objection. The DMCN does not require a single mass adoption event. It has nine distinct adoption wedges, each with a paying constituency, each contributing incrementally to network value. Enterprise adoption for machine-to-machine communication seeds the identity registry. Regulated industry adoption for compliance and archival integrity expands it. SaaS adoption for deliverability expands it further. Consumer adoption follows the density, rather than preceding it.

> **Strategic Implication**
> *The DMCN's go-to-market strategy should not be framed as "replace email." It should be framed as "deploy cryptographic identity infrastructure with email as the first application." Each constituency that deploys for their specific pain point strengthens the network for every other constituency. The network effect becomes an accelerant rather than a barrier once the first adoption wedges are established.*




---

## 6. Proposed Architecture: Decentralized Mesh Communication Network


### 6.1 Design Principles


The Decentralized Mesh Communication Network (DMCN) is designed around a
set of foundational principles derived from the failure analysis of
prior approaches. These principles are architectural constraints that
shape every design decision.

- Identity is cryptographic and self-sovereign. Every participant in the network has a unique identity derived from a public/private key pair. Identity is not assigned by any central authority and cannot be revoked by any third party.

- Sender verification is mandatory and protocol-enforced. A message without a valid cryptographic signature from a registered identity cannot enter the network. Verification is not optional, not opt-in, and not a filter applied after the fact — it is a gate at the point of transmission.

- The network is peer-to-peer with no central routing authority. Messages are relayed through a distributed mesh of nodes. No single entity controls routing, storage, or delivery.

- Complexity is hidden from end users. The cryptographic machinery that makes the network trustworthy operates entirely below the user experience layer. Users interact with human-readable identities and familiar communication patterns.

- Legacy email interoperability is a first-class requirement. The network must be capable of sending to and receiving from legacy SMTP addresses during a transition period.


### 6.2 Network Architecture


The DMCN consists of three logical layers, each with distinct
responsibilities:


#### 6.2.1 Identity Layer


The Identity Layer is responsible for the creation, storage, and
discovery of cryptographic identities. Each user identity is represented
by an elliptic curve key pair. The public key serves as the user's
network identifier and is registered in a distributed identity registry
--- a peer-to-peer data structure analogous to a distributed hash table
(DHT) — that allows any node to resolve a human-readable address to
its corresponding public key.

Human-readable addresses follow a format similar to traditional email
(user@domain) but resolve not to a mail server IP address but to a
public key and a set of authorized relay nodes. This means users can be
addressed in a familiar way while the underlying identity is
cryptographic and decentralized.


#### 6.2.2 Transport Layer


The Transport Layer is responsible for routing messages through the mesh
network. Messages are addressed to the recipient's public key,
encrypted with that public key, signed with the sender's private key,
and relayed through a network of nodes using an onion-routing-inspired
protocol that provides metadata privacy in addition to content privacy.

Relay nodes can verify message signatures against the identity layer,
which is the mechanism by which spam is rejected at the network level
--- a node that receives a message signed by an identity not registered
in the identity layer drops the message without delivery.


#### 6.2.3 Storage and Delivery Layer


Unlike real-time messaging systems, email is inherently asynchronous.
The Storage and Delivery Layer provides distributed, encrypted message
storage that holds messages until the recipient's client connects to
retrieve them. Messages are stored encrypted with the recipient's
public key; relay nodes providing storage cannot read message content.


### 6.3 Message Lifecycle


A message in the DMCN follows this lifecycle:

- The sender's client composes a message and signs it with the sender's private key.

- The client encrypts the signed message with the recipient's public key, retrieved from the Identity Layer.

- The encrypted, signed message is submitted to the transport layer with the recipient's public key as the address.

- Relay nodes verify the sender's signature against the Identity Layer. Messages with invalid or absent signatures are dropped.

- The message is routed through the mesh to relay nodes serving the recipient's address, where it is held in encrypted storage.

- When the recipient's client connects, it retrieves and decrypts messages using the recipient's private key.

- The recipient's client verifies the sender's signature, confirming the message genuinely originated from the stated sender.



---

## 7. Cryptographic Identity and Key Management


### 7.1 Key Generation and Storage


Each user account in the DMCN is associated with an elliptic curve key
pair, generated at account creation using well-established cryptographic
standards (specifically the Curve25519 or secp256k1 curve families,
which underpin both modern TLS and cryptocurrency wallets respectively).
The private key is generated on the user's device and never transmitted
to any server or relay node in plaintext form.

Private keys are stored in a hardware-backed secure enclave on devices
that support it (the Secure Enclave on Apple devices, the Trusted
Execution Environment on Android devices, TPM on Windows machines). On
devices without hardware security modules, keys are stored in encrypted
form protected by the user's authentication credential (biometric or
PIN). This approach is identical to how Apple Passkeys, Google Passkeys,
and hardware security keys manage private key material today.


### 7.2 Public Key Distribution and Discovery


The public key for each user identity is published to the DMCN's
distributed identity registry at account creation. The registry is a
peer-to-peer distributed data structure — conceptually similar to the
Distributed Hash Table used by the BitTorrent protocol — in which each
node stores a portion of the registry and can resolve any query through
a bounded number of hops.

When a user wishes to send a message to an address they have not
previously contacted, the client performs a registry lookup to retrieve
the recipient's public key. This lookup is cryptographically verifiable
--- the registry entry for each identity is itself signed by the private
key of the identity owner, allowing any node to verify that the public
key they receive genuinely belongs to the claimed identity.


### 7.3 The Key Management UX Problem — And Its Solution


The failure of PGP, despite its technical soundness, is primarily
attributable to the burden it placed on users to manage cryptographic
keys. The DMCN takes a fundamentally different approach, drawing on the
model established by passkeys and mobile device security:

- Key generation is automatic and invisible. When a user creates an account, keys are generated on their device without any user-facing step involving cryptographic concepts.

- Private keys are never shown to the user. Unlike cryptocurrency wallets, users are not presented with seed phrases or private key strings.

- Key backup is automatic and encrypted. Private keys are backed up using the device's existing encrypted cloud backup infrastructure (iCloud Keychain, Google Password Manager, or an equivalent DMCN-native encrypted key backup service).

- Multi-device access is handled through the primary/sub-key model described in Section 7.5. Each device holds its own sub-key; the user's primary key is the stable identity anchor. Losing or replacing a device requires only revoking that device's sub-key, not the entire identity.

- Account recovery is possible through a social recovery mechanism. Users designate trusted contacts who hold encrypted shards of a recovery key. A threshold of contacts (for example, 3 of 5) must participate to restore access.


### 7.4 Identity Verification and the Trust Graph


The DMCN's identity model provides cryptographic certainty that a
message came from the private key associated with a given public key.
The system also provides a voluntary trust graph through which users can
cross-sign each other's identities, creating a web of trust analogous
to PGP's original model but with modern UX — implemented as a simple
QR code exchange in the mobile app.


> **Trust Without Central Authority**
> *The DMCN does not require a central certificate authority to issue
> or revoke identity credentials. Trust is established through direct
> cryptographic verification and through the voluntary web of
> attestations that users build over time — exactly as trust is
> established in the physical world.*


### 7.5 Primary Keys and Device Sub-Keys

The DMCN uses a two-level key hierarchy to separate the question of *who a user is* from the question of *which device they are currently using*. This distinction is essential for multi-device support, graceful key rotation, and per-device revocation without identity loss.

#### 7.5.1 The Primary Key

Each DMCN identity has exactly one **primary key pair** at any point in time. The primary key is the canonical representation of the identity: it is what appears in the identity registry, what the trust graph is anchored to, what allowlists are bound to, and what contacts see when they look up an address. The primary key's public half is the stable long-term identifier for the address.

The primary key is generated on the user's first device at account creation and is backed up through the encrypted key backup infrastructure described in Section 7.3. It is never generated or held by any server or relay node.

#### 7.5.2 Device Sub-Keys

Each device on which the user activates their DMCN account generates its own **device sub-key pair** locally. The device sub-key is signed by the primary key, creating a cryptographically verifiable delegation: any party who trusts the primary key can verify that the sub-key is a legitimately authorised device for that identity.

Device sub-keys serve two functions. First, they are the keys actually used for per-message encryption and signing on that device — the primary key is not used for routine message operations, reducing its exposure. Second, they allow per-device revocation: if a device is lost, stolen, or decommissioned, only that device's sub-key need be revoked. The primary key, the identity's trust relationships, and all other active devices are unaffected.

Sub-keys are registered in the identity registry as children of the primary key record and are returned alongside the primary key in registry lookups. When sending to a multi-device user, the message payload is encrypted exactly once using a single randomly generated symmetric content key. That content key is then individually wrapped (encrypted) for each active device sub-key. Any device that holds the corresponding private sub-key can unwrap the content key and decrypt the payload. This approach — a Key Encapsulation Mechanism (KEM) pattern — ensures that message and attachment payload bytes appear on the wire only once regardless of how many devices the recipient has enrolled, eliminating the per-recipient payload duplication that would otherwise result. The protocol structure for this is defined in Section 15.3.3.

#### 7.5.3 Sub-Key Lifecycle

Sub-keys are created when a new device is activated and revoked when a device is decommissioned. Revocation is published to the identity registry and propagated through the same mechanism as primary key revocation. A revoked sub-key cannot be re-activated; a replacement device generates a new sub-key and re-registers it under the same primary key.

Sub-keys carry an optional `expires_at` field, enabling organisations to enforce periodic rotation of device credentials — relevant in managed domain contexts where the domain authority (Section 13) may require regular key refresh as a compliance control.

#### 7.5.4 Key Rotation

Periodic rotation of the primary key — whether on a schedule, following a suspected compromise, or as a policy requirement — is handled by publishing a new primary key to the registry via the `UPDATE` operation (Section 15.2.4), signed by both the old and new primary keys to prove continuity of control. This dual-signature rotation is self-authenticating: only the legitimate key holder could have produced a valid signature from the old private key. Contacts who have the rotating identity on their allowlist will accept the new key silently, with a non-blocking notification, rather than being required to re-verify — because the cryptographic chain of custody is intact. The full handling logic for signed versus unsigned rotations is specified in Section 14.1.2.

Following a signed rotation, the old key is retained in the registry in parallel with the new key for a seven-day window (configurable by domain authorities). This retention window provides a recovery path if the owner later discovers the old key was stolen: they can publish a `COMPROMISE` declaration against the old key during this window, which flags any rotation signed by it for re-verification by contacts. This is possible because a key theft is a copy, not a removal — the legitimate owner still holds the key and can sign the declaration. The `COMPROMISE` operation and its association rules are defined in Section 15.2.4.

In cases where the old private key is unavailable — device loss, destruction, or social recovery — the rotation cannot carry an old-key signature and will trigger re-verification prompts for allowlisted contacts. This is the correct behaviour: an unsigned rotation is indistinguishable from a registry substitution attack without additional out-of-band confirmation.

All device sub-keys must be re-issued under the new primary key following a primary key rotation, as existing sub-key signatures reference the old primary key fingerprint.

#### 7.5.5 Separation of Concerns

The two-level hierarchy also supports an intentional use case: a user who wishes to maintain distinct cryptographic personas for different communication contexts — for example, a public-facing address for general correspondence and a private address for a trusted inner circle — can register these as separate DMCN identities, each with its own primary key and sub-key tree. This is architecturally identical to having two email addresses, with the added property that the separation is cryptographically enforced rather than merely conventional. The DMCN client can manage multiple identities within a single account interface, presenting them as distinct inboxes.


The DMCN's identity model provides cryptographic certainty that a
message came from the private key associated with a given public key.
The system also provides a voluntary trust graph through which users can
cross-sign each other's identities, creating a web of trust analogous
to PGP's original model but with modern UX — implemented as a simple
QR code exchange in the mobile app.


> **Trust Without Central Authority**
> *The DMCN does not require a central certificate authority to issue
> or revoke identity credentials. Trust is established through direct
> cryptographic verification and through the voluntary web of
> attestations that users build over time — exactly as trust is
> established in the physical world.*



---

## 8. Spam Elimination at the Protocol Level


### 8.1 Why Cryptographic Identity Eliminates Spam


Spam exists because sending email is effectively free and sender
identity is effectively unverifiable. A spammer can send a billion
messages per day, claiming to be from any sender, at a cost measured in
fractions of a cent per message. The DMCN eliminates the conditions that
make spam possible rather than trying to detect and filter spam after it
has entered the network:

- Every sender must possess a valid registered private key. Creating a DMCN identity requires an account creation process — while frictionless for legitimate users, it is not free in the way that sending an SMTP email is free.

- Every message must bear a valid cryptographic signature from a registered identity. Relay nodes verify this signature before accepting a message for routing. A message without a valid signature is dropped at the first relay node.

- Sender identity is non-repudiable. Because messages are signed with the sender's private key, it is cryptographically impossible to forge a message that appears to come from a registered identity.

- Identity reputation is persistent and portable. An identity that sends unwanted messages can be blocked, and that block persists across sessions and devices.


### 8.2 Consent-Based Communication


The DMCN introduces consent-based message acceptance as a first-class
protocol feature. By default, a user's inbox accepts messages only from
identities that meet one of the following criteria:

- The sender is in the recipient's contact list.

- The sender shares a verified organizational identity with the recipient.

- The sender's identity has been vouched for by a trusted contact in the recipient's web of trust.

- The recipient has explicitly opted in to receiving messages from unknown senders (for public figures or customer-facing businesses).

Before any of these criteria are evaluated, the sender's cryptographic identity is checked against the recipient's personal blocklist and any shared reputation feeds the recipient has subscribed to. A sender present on either is silently dropped at the relay node before the message reaches the recipient's device — no delivery failure is returned to the sender, and no notification is shown to the recipient. Blocklist evaluation is a prerequisite gate, not a fallback; it applies regardless of whether the sender would otherwise have qualified for inbox or pending queue delivery.

Messages from senders that clear the blocklist check but do not meet any of the allowlist criteria are placed in a pending queue where the recipient can review them. These messages still bear valid cryptographic signatures — the sender's identity is known — allowing the recipient to make an informed decision.


### 8.3 Economic Disincentives for Spam


Beyond protocol-level barriers, the DMCN can implement optional economic
mechanisms that further disincentivize spam. A micro-payment channel
system allows senders to unknown recipients to attach a small,
refundable deposit to their message. If the recipient accepts, the
deposit is returned. If rejected, the deposit is forfeited. This imposes
no cost on messages between known contacts, and only trivial cost on
legitimate outreach — but makes mass spam campaigns economically
prohibitive.



---

## 9. User Experience: Hiding Complexity Without Sacrificing Security


### 9.1 The Fundamental Principle


The history of secure communication tools is, in large part, a history
of UX failures. The DMCN's user experience layer is designed around a
single principle: the security model must be invisible to the user in
normal operation. Users should experience the DMCN as a familiar
messaging application — with inboxes, contacts, compose windows, and
threads — and should never encounter the words 'key', 'signature',
'certificate', or 'encryption' in the normal flow of using the
product.


### 9.2 Familiar Addressing


Users are addressed using a format that mirrors traditional email: a
local identifier and a domain, separated by the @ symbol. Internally,
this address resolves to a public key — but from the user's
perspective, it is simply their address, just as a phone number is
simply a phone number without any awareness of the SS7 routing protocol
underneath.


### 9.3 Onboarding Flow


Account creation in the DMCN is designed to be comparable in friction to
creating a Gmail account. The user provides a chosen identifier,
authenticates with biometric or PIN, and the application generates their
key pair silently in the background. The entire process takes under two
minutes. There is no seed phrase, no key download, no certificate
request, and no cryptographic terminology.


### 9.4 Contact Discovery


Finding another user on the DMCN requires only their address. The
application resolves the address to a public key through the distributed
identity registry, and the contact appears in the user's contact list.
All messages to that contact are automatically encrypted and signed. The
user does not need to take any additional steps to enable security ---
it is on by default and cannot be turned off.


### 9.5 Trust Indicators


The application surfaces trust information in intuitive, non-technical
ways. Verified organizational identities display a checkmark alongside
the sender's domain. Messages from unknown senders appear in a separate
pending section. A simple trust indicator shows whether a contact's
identity has been verified by mutual connections in the user's network.



---

## 10. Transition Strategy: Coexistence with Legacy Email


### 10.1 The Migration Problem


No communication platform has ever achieved mainstream adoption by
requiring users to abandon their existing communication infrastructure.
The transition strategy for the DMCN is built on the principle of
graceful degradation — the system provides maximum value and security
to users communicating with each other on the DMCN, while maintaining
the ability to communicate with legacy email users at reduced security
levels during the transition period.


### 10.2 DMCN-to-DMCN Communication


When both sender and recipient are on the DMCN, messages are fully
encrypted, cryptographically signed, peer-to-peer routed, and spam-free
by protocol. This is the target state of the system and the experience
that should be promoted as the default.


### 10.3 DMCN-to-Email Communication


When a DMCN user sends a message to a legacy email address, the message
passes through a gateway node that translates it to SMTP format for
delivery. The message can include a footer inviting the recipient to
join the DMCN. Security properties are reduced in this path — message
content must be decrypted at the gateway for SMTP delivery — but
sender identity remains verifiable at the gateway level.


### 10.4 Email-to-DMCN Communication


Receiving a message from a legacy email sender requires a verified
gateway address system, where legacy emails pass through a gateway that
performs basic spam filtering and sender verification at the SMTP level
before delivering to the DMCN inbox. Users may also maintain a connected
legacy email address displayed in a separate, clearly labeled section of
their DMCN client.


### 10.5 Progressive Migration Incentives


The transition strategy includes mechanisms that actively incentivize
migration to native DMCN communication: visible trust indicators that
distinguish DMCN-verified messages from legacy email; organizational
compliance features requiring DMCN for sensitive internal
correspondence; and developer APIs allowing third-party applications to
integrate DMCN identity as a communication primitive.


### 10.6 IMAP and POP3 Compatibility — The Local Bridge Model

The transition strategy addresses SMTP interoperability through the bridge architecture in Section 11, but does not yet address the retrieval side of the legacy email stack. POP3 and IMAP are the protocols by which mail clients — Outlook, Apple Mail, Thunderbird, and their enterprise equivalents — pull messages from a mail server for display. Organisations and users deeply embedded in IMAP-based workflows represent a significant adoption barrier if the only path to DMCN participation requires replacing their mail client entirely.

A local IMAP bridge addresses this directly. Analogous to ProtonMail Bridge, it is a lightweight process that runs on the user's own device, exposes a localhost IMAP interface to the existing mail client, and handles DMCN message retrieval and decryption transparently in the background. The user's mail client connects to localhost rather than a remote server; the bridge fetches encrypted messages from DMCN relay nodes, decrypts them using the private key stored in the device's secure enclave, and presents them as ordinary IMAP mailbox folders. From the mail client's perspective, nothing has changed.

The security properties of this model are meaningfully different from a server-side IMAP bridge. Because decryption occurs on the user's device and the private key never leaves the secure enclave, end-to-end encryption is preserved end-to-user — the cleartext is only ever present in the local bridge process on the user's own hardware. A server-side IMAP bridge, by contrast, would decrypt messages on a third-party server, reducing security to the level of a conventional hosted email provider. The DMCN specification should explicitly support the local bridge model as a first-class transition mechanism and should discourage server-side IMAP bridge deployments except where the operator and user have the same trust relationship as a self-hosted mail server.

The local IMAP bridge does not affect the outbound path. Messages composed in the legacy mail client are intercepted by the bridge, signed with the user's private key, encrypted to the recipient's public key, and submitted to the DMCN transport layer — again transparently. For recipients on legacy email, the existing SMTP bridge handles delivery as normal.

This mechanism allows an organisation to deploy DMCN for its entire user base without requiring any change to end-user mail clients, devices, or workflows. It is the most friction-free enterprise adoption path available and should be a priority deliverable in the initial implementation roadmap.



---

## 11. The SMTP-DMCN Bridge Architecture


A Decentralized Mesh Communication Network that cannot communicate with
the existing email ecosystem is, for practical purposes, a closed
system. The SMTP-DMCN Bridge is the infrastructure component that makes
gradual, real-world adoption possible — allowing DMCN users to send
and receive messages with the 4 billion people on legacy email, without
compromising the security properties of native DMCN communication.

This section provides a detailed architectural treatment of the bridge,
covering both directions of message flow, the trust model the bridge
operates under, its internal components, and the mechanisms by which it
preserves as much security as possible within the constraints of
protocol translation.


### 11.1 Architectural Role and Placement


Bridge nodes occupy a distinct position in the DMCN topology. They are
DMCN-native nodes — registered in the identity layer with their own
cryptographic identity — that additionally operate SMTP listener and
sender services on the legacy internet. A bridge node is, in effect,
bilingual: it speaks DMCN natively and SMTP at the boundary.

Bridge nodes can be operated by the DMCN foundation, by commercial
providers, or by organizations that wish to run their own bridge
infrastructure for compliance or privacy reasons. The bridge is not a
single central server — multiple independent bridge operators can
coexist, and DMCN users can choose which bridge operator handles their
legacy email traffic, much as organizations today choose their own mail
exchanger (MX) records.


> **Design Principle**
> *The bridge is a temporary necessity during the transition period,
> not a permanent fixture of the architecture. Every design decision
> prioritizes making native DMCN adoption more attractive, so that
> dependence on the bridge diminishes naturally over time as the
> network grows.*


### 11.2 Outbound Path: DMCN to SMTP


When a DMCN user addresses a message to a legacy email address, the
message follows a modified delivery path through a bridge node. This is
the simpler of the two directions, and the security properties are
well-defined.


#### 11.2.1 Message Flow


- The sender's DMCN client composes and signs the message with the sender's private key, as in a standard DMCN message.

- The client detects that the recipient address resolves to a legacy email address (no DMCN public key found in the identity registry) and routes the message to a bridge node rather than the standard transport layer.

- The bridge node receives the encrypted, signed DMCN message, verifies the sender's signature against the identity registry, and decrypts the message content using the bridge's own private key (the message is re-encrypted to the bridge's key rather than a non-existent recipient key).

- The bridge constructs a standard SMTP message from the decrypted content, applying DKIM signing using the bridge operator's domain key. The `From` header is set to the sender's human-readable DMCN address (e.g. `alice@mycompany.com`), so that the recipient's email client displays the sender's familiar identity. The `Sender` header is set to the bridge's own address (e.g. `bridge@bridge.dmcn.net`), identifying the mailbox actually responsible for transmission per RFC 5322. The `Reply-To` header is set to the sender's DMCN bridge receive address, so that replies from legacy clients are routed back through the bridge correctly. This combination — `From` showing the human author, `Sender` showing the transmitting agent — is the same pattern used by legitimate bulk email providers (such as Mailchimp and SendGrid) when sending on behalf of their customers. SPF and DKIM pass against the bridge's own domain; the `From` address is not required to be an authorised sender for those checks under standard DMARC evaluation when the `Sender` header is present.

- The bridge delivers the SMTP message to the recipient's mail server using standard MX lookup and SMTP relay.

- A standardized footer is appended to the message body, displaying the sender's verified DMCN identity and an invitation link for the recipient to join the DMCN and receive future messages natively.


#### 11.2.2 Trust Properties


The outbound path involves one unavoidable trust compromise: the bridge
node decrypts message content in order to re-encode it as SMTP. This
means the bridge operator has technical access to message content in
transit. This is an honest limitation that must be clearly disclosed to
users and is analogous to the trust placed in any email service provider
today.

This limitation is mitigated by several factors: bridge operators are
registered DMCN identities with cryptographic accountability; users can
choose their bridge operator and can migrate between operators;
organizations with strong confidentiality requirements can operate their
own bridge nodes; and the limitation affects only messages sent to
legacy email recipients, not native DMCN-to-DMCN communication.


### 11.3 Inbound Path: SMTP to DMCN


Receiving messages from legacy email senders is the more complex
direction, because the sender has no DMCN identity and no cryptographic
signing capability. The bridge must make a trust determination about an
unverified sender and communicate that determination clearly to the
recipient.


#### 11.3.1 Addressing


DMCN users who wish to receive legacy email obtain a bridge address ---
a standard email address managed by the bridge operator (e.g.,
username@bridge.dmcn.net). This address is registered as an MX record
pointing to the bridge node's SMTP listener. Users can publish this
address on websites, business cards, and legacy systems as their contact
point for people who have not yet adopted DMCN.


#### 11.3.2 Message Flow


- The bridge's SMTP listener receives an inbound message addressed to the user's bridge address.

- The bridge performs a full suite of legacy authentication checks: SPF validation, DKIM signature verification, DMARC policy evaluation, and reverse DNS lookup on the sending mail server.

- The bridge queries distributed reputation databases (analogous to existing RBL/DNSBL services) for the sending IP address and domain.

- Messages that fail hard authentication checks (invalid DKIM, SPF failure with DMARC reject policy) are dropped with a delivery failure response to the sending server.

- Messages that pass or partially pass authentication are classified into trust tiers: Verified Legacy Sender (valid DKIM + positive reputation), Unverified Legacy Sender (no DKIM or neutral reputation), and Suspicious Legacy Sender (reputation flags present).

- The bridge wraps the classified message in a DMCN envelope, signed by the bridge's own private key as an attestation of the classification outcome. The DMCN envelope includes the full authentication result metadata, the trust tier assignment, and the original SMTP headers.

- The wrapped message is encrypted to the recipient's DMCN public key and delivered through the standard DMCN transport layer to the recipient's inbox.


#### 11.3.3 Recipient Experience


Inbound legacy messages appear in a clearly distinguished section of the
recipient's DMCN client, visually separated from native DMCN messages.
Each message displays its trust tier — Verified Legacy, Unverified
Legacy, or Suspicious — with a plain-language explanation of what the
classification means. The recipient can promote a legacy sender to their
contact list (which triggers the DMCN client to send the sender an
invitation to join the network natively) or block the sender
permanently.


### 11.4 Bridge Node Security Model


Bridge nodes are high-value infrastructure components and require a
rigorous security model. Several specific threats must be addressed:

- Bridge impersonation — a malicious actor operating a fraudulent bridge that misrepresents message authentication results. Mitigated by requiring bridge operators to register their identity in the DMCN identity registry, publish their security practices, and submit to periodic independent audits.

- Content interception — a bridge operator reading or modifying message content in transit on the outbound path. Mitigated by end-to-end message signing (recipients can verify the sender's signature even after bridge translation), audit logging, and regulatory accountability for commercial operators.

- SMTP relay abuse — spammers attempting to use the inbound bridge path to inject spam into DMCN inboxes. Mitigated by the authentication classification system and by rate limiting on the bridge's SMTP listener per sending domain and IP.

- Bridge node compromise — an attacker gaining control of a bridge node. Mitigated by the decentralized bridge model (no single bridge handles all traffic), key rotation protocols, and the ability for users to revoke trust in a specific bridge operator.


### 11.5 Federated Bridge Architecture


To avoid reintroducing centralization through the bridge layer, the DMCN
specification defines an open Bridge Operator Protocol (BOP) that any
qualified operator can implement. Bridge operators are discoverable
through the DMCN identity registry, and DMCN clients implement automatic
bridge selection based on operator reputation, geographic proximity,
organizational policy, and user preference.

This federated model mirrors the existing email MX record system — the
delivery path is determined by the recipient's published preferences,
not by any central routing authority — while adding the cryptographic
accountability layer that SMTP lacks. An organization can run its own
bridge, use a commercial bridge provider, or designate different bridges
for inbound and outbound traffic.


### 11.6 Precedents and Comparable Systems


The bridge architecture is a well-established pattern in communication
infrastructure. Several analogous systems demonstrate that protocol
translation at scale is an engineering challenge with proven solutions:

- Matrix bridges (Synapse) — the Matrix protocol has operated production bridges to Slack, Discord, WhatsApp, Telegram, SMS, and IRC simultaneously, with hundreds of thousands of active users. The Matrix bridge architecture provides a detailed precedent for federated, bidirectional protocol translation.

- Email-to-SMS gateways — carriers have operated SMTP-to-SMS translation services for decades, handling billions of messages. The protocol mismatch between email and SMS is in many respects more severe than the mismatch between SMTP and DMCN.

- SIP-PSTN gateways in VoIP — Voice over IP systems routinely bridge between SIP-native networks and the legacy public switched telephone network, including cryptographic signaling translation. The architecture for handling identity and trust across this boundary is directly analogous to the SMTP-DMCN bridge problem.


> **Engineering Assessment**
> *The SMTP-DMCN bridge is a significant but tractable engineering
> project. It does not require novel cryptographic research — it
> requires careful implementation of known patterns applied to a new
> protocol boundary. A working prototype bridge is a realistic
> deliverable for an initial proof-of-concept phase.*



---

## 12. Bringing Existing Email Addresses to the DMCN


One of the most significant friction points in any transition away from
legacy email is the email address itself. For most people and
organizations, an email address is not merely a routing string — it is
a persistent identity, published on business cards, embedded in
contracts, known to colleagues and clients accumulated over years or
decades. Any system that requires users to abandon their existing
address in order to participate faces an adoption barrier that is
largely insurmountable. The DMCN's address portability feature directly
addresses this by allowing users to bring their existing email addresses
into the DMCN identity layer without abandoning them.


### 12.1 The Principle of Address Portability


Address portability in the DMCN means that an existing email address can
be registered as a DMCN identity anchor — a verified link between a
known email address and a cryptographic public key. Once this link is
established and published in the DMCN identity registry, the address
functions simultaneously as a legacy SMTP address and as a native DMCN
identity. Senders who are on the DMCN will automatically discover the
cryptographic key and send natively; senders on legacy email will
continue to reach the address through conventional SMTP delivery.

This dual-mode operation is the bridge between the old world and the
new. It requires no change of address, no notification to existing
contacts, and no interruption of legacy email delivery. The upgrade is
invisible to legacy senders and automatic for DMCN senders.


### 12.2 Verification Mechanisms


The strength of the link between an email address and a DMCN identity
depends on the level of control the user has over the address and its
underlying domain. The DMCN supports three verification tiers, each with
distinct trust properties:


#### 12.2.1 Provider-Hosted Address Verification (e.g., Gmail, Outlook)


For addresses hosted by third-party providers — the most common case
for individual users — verification proceeds through an email
confirmation flow analogous to standard account verification practices:

- The user claims ownership of their existing address (e.g., alice@gmail.com) within the DMCN client.

- The DMCN identity service sends a time-limited, single-use verification code to that address via the bridge's outbound SMTP path.

- The user retrieves the code from their legacy inbox and enters it in the DMCN client, confirming they control the address.

- The DMCN identity registry publishes a signed binding record linking the address to the user's public key.

This tier provides practical ownership verification — proof that the
user can receive mail at the address — but does not constitute
cryptographic domain ownership. The binding is valid as long as the user
retains control of the legacy account.


#### 12.2.2 Custom Domain Address Verification


For users and organizations who control their own domain (e.g.,
alice@mycompany.com), a stronger verification path is available that
mirrors the DNS-based ownership proof used by DKIM and SPF:

- The user requests a DMCN verification token for their domain from the DMCN identity service.

- The user publishes the token as a DNS TXT record at a standardized subdomain (e.g., \_dmcn.mycompany.com), alongside their public key fingerprint.

- The DMCN identity service performs a DNS lookup to confirm the record is present and correctly formatted.

- The binding is published in the identity registry with a Domain-Verified status, indicating that the link is backed by DNS control rather than merely inbox access.

Domain-verified bindings are significantly more robust. They can be
updated by the domain owner at any time through DNS changes, they do not
depend on the policies of any email provider, and they are resistant to
account suspension or provider-side interference. For organizations,
this is the recommended verification path.

Domain ownership verification establishes *that* an organisation controls a domain. It does not address *how* the organisation governs which addresses are authorised under that domain, or how it provisions and deprovisions staff identities. Those operational requirements are addressed by the Domain Authority Record model in Section 13.


#### 12.2.3 DANE-Style Cryptographic Domain Binding


For domains that have enabled DNSSEC — the cryptographic extension to
DNS that provides tamper-evident records — a third verification tier
is available that provides the highest level of assurance. In this
model, the domain owner publishes the DMCN public key directly in a
DNSSEC-signed record, creating a chain of cryptographic trust from the
DNS root through the domain to the individual identity. This approach is
analogous to DANE (DNS-based Authentication of Named Entities), which is
already used in some contexts to bind TLS certificates to domain names
without relying on certificate authorities.


### 12.3 Trust Implications of Each Tier


| **Verification Tier** | **Proof of Control** | **Resistant to Provider Action** | **Recommended For** |
|---|---|---|---|
| Provider-Hosted (Gmail, Outlook) | Inbox access only | No — provider can suspend | Individual users during transition |
| Custom Domain DNS Verification | DNS record control | Yes — domain owner controls DNS | Professionals, small businesses |
| DNSSEC / DANE Cryptographic Binding | Cryptographic DNS chain | Yes — highest assurance | Enterprises, regulated industries |


### 12.4 The Honest Limitation: Ownership vs. Control


Address portability is a powerful adoption mechanism, but it requires an
honest disclosure that the whitepaper must not obscure: bringing a
provider-hosted address to the DMCN does not give the user cryptographic
ownership of that address at the domain level. Google still controls
\@gmail.com. Microsoft still controls \@outlook.com. If a user's Google
account is suspended, terminated, or if Google elects to block
DMCN-related traffic, the legacy delivery path for that address breaks
--- though the user's DMCN identity and their cryptographic key persist
independently.

This distinction has practical consequences that users should understand
at onboarding. Provider-hosted address linking is a convenience feature
for the transition period. For users who want long-term,
provider-independent identity, the DMCN should actively encourage
migration to a custom domain. The client application can surface this
recommendation appropriately — not as a barrier, but as a path to
greater identity sovereignty over time.


> **Identity Sovereignty Principle**
> *A provider-hosted address gives you a key to a house you rent. A
> custom domain address gives you a key to a house you own. The DMCN
> supports both — but only one provides true long-term identity
> independence.*


### 12.5 Address Portability and the Spam Problem


Address portability introduces one additional consideration for the spam
model: a user who verifies ownership of an existing email address brings
with them the reputation — positive or negative — associated with
that address in legacy spam databases. The DMCN identity layer should
initialize the reputation of a newly verified address using available
legacy reputation signals as a starting point, rather than treating
every ported address as a clean slate.

Conversely, address portability is a meaningful barrier to spam identity
laundering. A spammer who wishes to port a known-good email address to a
DMCN identity in order to inherit its reputation must actually control
that address — they cannot simply claim it. This is a meaningfully
higher bar than the trivially low cost of sending SMTP mail from an
arbitrary claimed address.


### 12.6 Precedents


The address portability model draws on several well-established
precedents in both identity verification and email infrastructure:

- Keybase — demonstrated the viability of linking a cryptographic identity to multiple existing identities (email, Twitter, GitHub, domain) through a system of cryptographic proofs. The DMCN's verification model is conceptually similar but integrated at the protocol level rather than as a third-party overlay.

- Google Workspace and Microsoft 365 custom domain onboarding — both services use DNS TXT record verification to prove domain ownership before allowing custom domain email. The DMCN's Domain Verification tier follows this exact pattern, which is already familiar to IT administrators globally.

- DKIM public key DNS records — the practice of publishing cryptographic keys in DNS records is already standard email infrastructure. The DMCN's DANE-style binding extends this established pattern.

- Number portability in mobile telephony — the telecommunications industry solved an analogous problem when it allowed consumers to bring their phone numbers between carriers. The lesson from that transition is directly applicable: portability dramatically lowers switching costs and accelerates adoption of superior infrastructure.

Address portability answers the question of how an existing address enters the DMCN identity layer. For organisations deploying DMCN across a workforce, a further layer is required: ongoing governance of who is authorised to operate under the domain, how new addresses are provisioned, and how access is revoked when people leave. This organisational layer is specified in Section 13.



---

## 13. Domain Authority and Organisational Address Management

Individual cryptographic identity — a person generating a key pair and registering an address — is the foundation of the DMCN model. But organisations are not collections of independent individuals. A company that deploys DMCN across its workforce needs to control which identities are authorised to operate under its domain, provision addresses for new staff, revoke access when people leave, enforce policy across its namespace, and delegate administrative authority within large or complex structures. None of this is addressed by the individual identity layer described in Section 7 or the address portability mechanisms in Section 12.

This section defines the **Domain Authority Record** (DAR) — an organisational layer that sits above individual identity records in the registry and gives domain owners administrative control over their address namespace without compromising the cryptographic self-sovereignty of individual users.

> **Design Principle**
> *The DAR gives organisations control over their namespace while preserving individual key sovereignty. A domain authority can provision and deprovision addresses, but it cannot read messages, impersonate users, or extract private keys. Administrative authority over an address is distinct from cryptographic ownership of the identity.*

---

### 13.1 The Domain Authority Record

A Domain Authority Record is a signed registry entry that declares a domain's administrative policy and the public key of the authority responsible for enforcing it. It is published by the domain owner and referenced by all individual identity records under that domain.

```
domain_authority_record {
    version:                  uint32
    domain:                   string          // e.g. "company.com"
    authority_ed25519_pubkey:  bytes[32]       // Ed25519 public key of the domain authority
    authority_x25519_pubkey:   bytes[32]       // X25519 public key of the domain authority
    policy_flags:             uint32          // bitmask; see Section 13.1.1
    sub_authorities:          repeated sub_authority_record
    verification_tier:        enum { DOMAIN_DNS, DNSSEC_DANE }
    created_at:               uint64
    expires_at:               uint64          // 0 = no expiry
    dns_verification_token:   string          // published at _dmcn-authority.<domain> TXT record
    self_signature:           bytes[64]       // Ed25519 signature over all preceding fields
}

sub_authority_record {
    sub_authority_address:    string          // DMCN address of the delegated administrator
    sub_authority_pubkey:     bytes[32]       // Ed25519 public key of the delegated administrator
    scope:                    string          // namespace scope, e.g. "engineering.company.com"
    permissions:              uint32          // bitmask of delegated permissions
    granted_at:               uint64
    expires_at:               uint64
    grant_signature:          bytes[64]       // Ed25519 signature by domain authority
}
```

The DAR is published to the identity registry using the same DHT infrastructure as individual identity records, keyed on the SHA-256 hash of the domain string. Any node can verify the DAR's authenticity by checking the `self_signature` against the `authority_ed25519_pubkey`, and can verify the domain owner's control of the domain by checking that the `dns_verification_token` appears in the DNS TXT record at `_dmcn-authority.<domain>`.

#### 13.1.1 Policy Flags

The `policy_flags` bitmask allows domain authorities to declare administrative policies that relay nodes and clients enforce automatically:

| Flag | Value | Meaning |
|---|---|---|
| `REQUIRE_DOMAIN_COUNTERSIG` | 0x01 | Identity records under this domain must carry a valid domain countersignature to be treated as verified |
| `REJECT_UNMANAGED` | 0x02 | Relay nodes must reject messages from addresses under this domain that do not carry a valid domain countersignature |
| `ARCHIVE_REQUIRED` | 0x04 | All messages to/from this domain must be routed through an approved archive bridge (see Section 13.5) |
| `MANAGED_DEVICES_ONLY` | 0x08 | Identities under this domain may only be registered from device identifiers pre-approved by the domain authority |
| `SUBDOMAIN_DELEGATION` | 0x10 | Sub-authority delegation is permitted; sub-authority records in this DAR are valid |

---

### 13.2 Identity Provisioning Under a Managed Domain

When `REQUIRE_DOMAIN_COUNTERSIG` is set, an individual identity record under the domain is only treated as fully verified if it carries a countersignature from the domain authority in addition to the individual's own self-signature.

The provisioning flow is:

1. The new user generates their key pair on their device, as in standard DMCN account creation (Section 7.1). The private key never leaves their device.
2. The user's DMCN client submits the unsigned identity record to the domain authority's provisioning endpoint — a service operated by the organisation's IT administrator.
3. The provisioning service verifies that the requested address is within the domain's namespace, that it has not already been allocated, and that the request is authorised (for example, by cross-referencing an HR system or directory service).
4. The domain authority signs the identity record with its Ed25519 private key, producing a `domain_countersignature`, and returns it to the user's client.
5. The client publishes the identity record to the registry with both the individual self-signature and the domain countersignature. Relay nodes and peers see a fully verified managed identity.

This flow ensures that only authorised addresses are registered under the domain, that the domain authority has a record of every active identity in its namespace, and that the individual's private key is generated and held exclusively on their own device throughout. The domain authority signs the address binding — it does not generate or hold the private key.

---

### 13.3 Deprovisioning and Revocation

Deprovisioning is the critical operational requirement that the individual identity model cannot satisfy alone. When an employee leaves an organisation, the company needs to invalidate their `alice@company.com` identity immediately, regardless of whether Alice cooperates or whether her device is available.

The DAR model provides two complementary revocation mechanisms:

#### 13.3.1 Domain Authority Revocation

The domain authority can publish a signed revocation notice for any address under its domain. This is distinct from the individual's own revocation capability described in Section 15.2.4: it does not revoke Alice's underlying key pair, but it withdraws the domain countersignature that makes her address valid under `company.com`.

```
domain_revocation_record {
    domain:                   string
    revoked_address:          string
    revoked_pubkey:           bytes[32]
    revocation_reason:        enum { OFFBOARDING, POLICY_VIOLATION, KEY_COMPROMISE, ADMINISTRATIVE }
    revoked_at:               uint64
    authority_signature:      bytes[64]      // Ed25519 signature by domain authority
}
```

Once a domain revocation record is published to the registry, relay nodes enforcing `REQUIRE_DOMAIN_COUNTERSIG` or `REJECT_UNMANAGED` will reject messages from the revoked address. The revocation propagates through the DHT within the same timeframe as any other registry update.

The revoked user's underlying DMCN identity — their key pair — is unaffected. They retain their cryptographic identity and can register a new address under a personal domain or provider-hosted address. What they lose is the authorised binding to the organisation's domain.

#### 13.3.2 Key-Change Notification on Revocation

When a domain authority revokes an address, the allowlists of users who had that address as a trusted contact must be updated. The key-change notification mechanism described in Section 14.1.2 fires automatically on revocation: contacts are alerted that the identity previously known as `alice@company.com` is no longer active under that address, and are prompted to re-verify before resuming communication. This prevents a revoked identity from continuing to receive messages intended for the legitimate address holder.

---

### 13.4 Namespace Management and Conflict Resolution

Without a domain authority, nothing prevents two people from both attempting to register as `cfo@company.com`. The first registration wins in the DHT, leaving the domain with no recourse and creating a potential for namespace squatting or impersonation within the domain.

When a DAR is published with `REQUIRE_DOMAIN_COUNTERSIG`, this problem is structurally eliminated. A registration attempt for an address under a managed domain without the required countersignature is rejected by the registry as unverified. Only addresses that have passed through the domain authority's provisioning flow carry valid countersignatures and are accepted as verified identities under that domain.

Relay nodes enforce this at message receipt: a message purporting to come from `cfo@company.com` without a valid domain countersignature from `company.com`'s registered DAR is treated as unverified and routed to the recipient's pending queue regardless of the individual self-signature. The display name and address string are the same; the verification indicator is absent, which is the signal the recipient's client surfaces.

---

### 13.5 Compliance Archiving Under the DAR Model

Regulated industries — financial services, healthcare, legal — are required to archive communications and make them available to regulators on demand. The `ARCHIVE_REQUIRED` policy flag in the DAR enables this without compromising end-to-end encryption for unregulated communication paths.

When `ARCHIVE_REQUIRED` is set, the domain authority registers one or more approved archive bridge identities in its DAR. The DMCN client automatically BCC-encrypts a copy of every outbound message to the archive bridge's public key before sending. The archive bridge is a registered DMCN identity that stores received messages in encrypted form, indexed by sender, recipient, and timestamp. Because the copy is encrypted to the archive bridge's public key using the standard DMCN encryption scheme, the archive is end-to-end encrypted between the sender's client and the archive service — relay nodes and third parties cannot read archived content.

The domain authority, operating or contracting the archive bridge, holds the archive bridge's private key and can decrypt archived messages in response to a regulatory request or legal process. Individual users cannot disable archiving while connected to a domain that sets `ARCHIVE_REQUIRED`. The archive copy is sent transparently by the client, analogous to a BCC to a compliance mailbox — a pattern already familiar in regulated enterprise email deployments.

This approach satisfies the compliance tension identified in Section 19.3 for managed domain users, without requiring any change to the encryption model for native DMCN-to-DMCN communication between unmanaged users.

---

### 13.6 Delegated Administration

Large organisations require delegated administrative authority. A global enterprise cannot manage all address provisioning from a single root authority; subsidiaries, business units, and departments need their own administrative scope.

The `sub_authority_record` within the DAR supports this through a two-level delegation model:

- The root domain authority for `company.com` can delegate administrative scope over `engineering.company.com` to a designated sub-authority identity.
- The sub-authority can provision and deprovision addresses within its delegated scope without involving the root authority.
- The root authority retains override capability — it can revoke any address under the full domain regardless of which sub-authority provisioned it.
- Sub-authority delegation is itself time-bounded and can be revoked by the root authority by publishing an updated DAR.

Sub-authorities are identified by their own DMCN identity (address and public key), making their administrative actions cryptographically attributable. An audit log of provisioning and revocation events signed by the responsible authority's key provides a tamper-evident record of all address lifecycle events under the domain.

The current specification supports one level of delegation. Deeper hierarchies (sub-sub-authorities) are deferred as an extension — the `SUBDOMAIN_DELEGATION` policy flag will govern whether sub-authorities may themselves delegate, when that extension is specified.

---

### 13.7 Threat Considerations for the Domain Authority Model

The DAR introduces a concentration of authority that has no equivalent in the individual identity model and creates specific threats that must be addressed.

**Domain authority key compromise.** The domain authority's Ed25519 private key is the root of trust for all managed identities under the domain. Compromise of this key allows an attacker to provision fraudulent identities under the domain, revoke legitimate identities, and modify domain policy. This is analogous to the compromise of an enterprise certificate authority root key, and should be treated with equivalent seriousness. The domain authority key should be held in hardware security modules (HSMs), subject to multi-party signing requirements for sensitive operations, and rotated on a defined schedule. The DAR's `expires_at` field enforces periodic rotation.

**Insider threat and administrative abuse.** Domain authorities and sub-authorities have the power to revoke users and modify namespace policy. This power can be abused — a malicious administrator could revoke a whistleblower's identity or provision a shadow identity to conduct surveillance. The audit log of all provisioning and revocation events, signed by the acting authority, provides accountability after the fact. Real-time alerting to affected users on key-change events (Section 14.1.2) limits the window of undetected abuse. The DMCN specification should require that domain revocation events be visible to the revoked user — silent revocation should not be possible.

**Domain authority as a surveillance surface.** A domain authority knows every address registered under its domain, can observe the timing of provisioning and deprovisioning events, and — if operating its own relay or archive infrastructure — has additional visibility into communication patterns. This is a materially greater surveillance capability than any individual relay node. It is, however, structurally equivalent to the existing capability of corporate email administrators under conventional hosted email. The DAR model does not increase this capability relative to the current state; it structures and makes it explicit.

**Interaction with Section 17 threat categories.** Domain authority key compromise has elements of both the key compromise threat (Section 17.7) and the infrastructure attack threat (Section 17.3). The mitigations for both apply: hardware key storage, multi-party signing, and the revocation and re-registration capability of the DHT. The specific additional mitigation is key rotation enforced by the `expires_at` field and the requirement that rotated DAR keys trigger re-verification of all managed identities.

---

### 13.8 Privacy Implications of the DAR Model

The DAR gives domain authorities structural visibility into their address namespace that individual relay nodes do not have. Specifically, a domain authority operating under the DAR model can observe:

- Every address registered under its domain, including registration and revocation timestamps
- The public keys associated with each address (not the private keys, and not message content)
- The set of archive-bridge-routed messages if `ARCHIVE_REQUIRED` is set and the authority operates the archive bridge

This visibility is appropriate for the organisational context — an employer has legitimate interest in knowing which employees have active communication identities under the company domain — but it should be explicitly disclosed at onboarding and reflected in the organisation's privacy policy and employment terms. The privacy analysis in Section 18 should be understood as applying to the unmanaged identity case; managed domain users operate under a modified privacy model in which the domain authority has the administrative visibility described above.

The DMCN client should display a clear, persistent indicator when a user's identity is operating under a managed domain, so that the user understands they are not in the fully self-sovereign privacy model. This is analogous to the way corporate device management (MDM) indicators work on mobile platforms.

---

### 13.9 Interaction with Existing Sections

The DAR model intersects with several other sections of this whitepaper:

**Section 12 (Address Portability)** — The DNS verification tiers in Section 12.2 establish domain ownership but do not address namespace governance. Domain authorities building on the `DOMAIN_DNS` or `DNSSEC_DANE` verification tiers in Section 12.2.2 and 12.2.3 should publish a DAR alongside their DNS verification record to activate managed domain behaviour. The two mechanisms are complementary: DNS verification proves domain ownership to the registry; the DAR declares the administrative policy the domain owner wishes to enforce.

**Section 14 (Trust Management)** — Contacts managed under a DAR-governed domain are subject to automatic key-change notification (Section 14.1.2) on both individual key rotation and domain authority revocation. Allowlisted contacts of a managed identity should be treated as requiring re-verification on deprovisioning events, not merely on key changes initiated by the individual.

**Section 19.3 (Regulatory Compliance)** — The `ARCHIVE_REQUIRED` flag and the archive bridge model described in Section 13.5 directly address the compliance gap identified in Section 19.3. This does not fully close the challenge — edge cases around cross-jurisdiction archiving, legal holds, and eDiscovery production remain open — but it provides the structural mechanism on which compliance tooling can be built.

**Section 18 (Privacy Analysis)** — The DAR's administrative visibility is a distinct privacy surface not covered in Section 17's bridge-focused analysis. A future revision of Section 18 should add a subsection addressing managed domain privacy, covering the domain authority's visibility, the distinction between managed and unmanaged users, and the client disclosure requirement described in Section 13.8.

**Section 15.2.4 (Identity Registry Operations)** — The `REGISTER` and `REVOKE` operations defined in the protocol specification apply to both individual identity records and Domain Authority Records. The `UPDATE` operation for a DAR triggers re-verification of all managed identities under the domain, which the protocol specification should make explicit in a future revision.


---

## 14. Trust Management: Allowlists, Pending Queue, and Blocklists


Cryptographic identity verification is the foundation of the DMCN's
trust model — it answers the question of whether a message genuinely
came from a claimed sender. But verification alone does not answer a
second, equally important question: whether the user actually wants to
hear from that sender. Trust management is the user-facing system that
sits on top of cryptographic verification and allows each participant to
define, on their own terms, who they trust, who they are uncertain
about, and who they actively reject.

The DMCN's trust management system operates at three tiers ---
allowlist, pending queue, and blocklist — each with distinct delivery
semantics, key storage implications, and sharing properties. Together
they form a layered defence that is more powerful than anything
available in legacy email, precisely because the identities being
managed are cryptographic and persistent rather than superficial and
easily spoofed.


### 14.1 The Allowlist: Confirmed Trusted Senders


The allowlist is the user's registry of confirmed trusted contacts. It
is not merely an address book — it is a cryptographically anchored
record that binds a human identity to a specific public key, with a
record of how and when that binding was established. Every entry in the
allowlist carries a trust provenance — the mechanism by which the user
confirmed the contact's identity — which is surfaced in the client UI
to help users understand the strength of each trust relationship.


#### 14.1.1 Trust Establishment Mechanisms


The DMCN supports multiple mechanisms for adding a contact to the
allowlist, ranked here in descending order of trust strength:

- Direct key exchange — the user and contact are physically present and exchange public keys via a QR code scan in the DMCN mobile application. This establishes an in-person verified binding with the highest possible assurance. The resulting allowlist entry is marked Verified In-Person.

- Fingerprint verification — the user retrieves a contact's public key from the identity registry and then verifies the key fingerprint through an independent channel (a phone call, a video call, a previously trusted communication method). The user confirms that the fingerprint read aloud by the contact matches the one in the registry. Marked Fingerprint Verified.

- Web-of-trust promotion — the contact is already allowlisted by two or more of the user's existing Verified contacts. The user can choose to extend trust on the basis of their network's endorsement, with a clear indication of which mutual contacts vouch for the new addition. Marked Network Vouched.

- Organisational verification — the contact holds a DMCN identity attested by an organisation the user has already verified (e.g., a colleague whose identity is attested by a shared employer domain). Marked Organisationally Verified.

- First-message confirmation — the user receives a message from an unknown sender and actively chooses to approve and allowlist them. This is the weakest trust mechanism — the user has reviewed the message and chosen to accept the sender, but has not independently verified the key. Marked User Approved.

Trust provenance is preserved indefinitely in the allowlist record and
is visible to the user at any time. A contact marked Verified In-Person
carries a fundamentally different assurance than one marked User
Approved, and the client communicates this distinction without requiring
the user to understand the underlying cryptography.


#### 14.1.2 Key Binding and Update Handling


Because allowlist entries are bound to specific public keys rather than
addresses alone, the DMCN client must handle the case where a contact's
key changes — for example, when they migrate to a new device, perform
a key rotation, or recover their account through the social recovery
mechanism.

The client distinguishes between two categories of key change, applying different handling to each:

**Signed rotation — silent update permitted.** When a new key is published to the registry accompanied by a valid signature from the previous private key, the rotation is self-authenticating: only the legitimate key holder could have produced that signature. In this case the client updates the allowlist binding silently, without prompting the user to re-verify. A brief non-blocking notification is surfaced — "Alice updated her key" — so the user is aware the change occurred, but no action is required and message delivery continues uninterrupted. This covers the common cases of routine key rotation, device migration, and scheduled credential refresh.

**Unsigned rotation — re-verification required.** When a new key is published without a valid signature from the previous key — because the old key was lost, the device was destroyed, or the account was recovered through the social recovery mechanism — the client cannot cryptographically distinguish a legitimate recovery from an attacker who has substituted a key in the registry. In this case the client suspends delivery from that identity, alerts the user that the contact's key has changed without a verifiable chain of custody, and requires explicit re-verification before the allowlist binding is updated. This preserves the original protection against registry substitution attacks for precisely the cases where that protection is needed.

This distinction means that the friction of re-verification is reserved for genuinely ambiguous key changes, rather than being imposed on every routine rotation. The social recovery case — where re-confirmation is unavoidable because the old key no longer exists — is handled honestly: recovery is disclosed as an unsigned transition, and contacts are informed accordingly.

**Rotation retention window.** When a signed rotation is published, the old key is not immediately retired. It remains active in the registry in parallel with the new key for a configurable retention window — default seven days. This window serves a specific security purpose: if the legitimate owner discovers after rotating that their old private key was stolen, they can use the old key to publish a `COMPROMISE` declaration against it (see below) during this window. Because the attacker may have used the stolen old key to push the new key, a compromise declaration on the old key within the retention window automatically flags the new key as tainted — requiring re-verification by all contacts before the new key is trusted. The retention window therefore gives the legitimate owner a recovery path in the scenario where an attacker rotates the identity before the owner detects the theft.

After the retention window expires, the old key is retired and the new key becomes the sole active primary key. Domain authorities managing a DMCN domain (Section 13) may configure a longer retention window — up to 30 days — for environments where delayed detection of key theft is a realistic concern.

**Compromise declarations.** A `COMPROMISE` declaration is a signed registry operation distinct from a standard `REVOKE`. It carries the semantic that the declared key should be treated as having been in hostile hands, not merely retired from use. The declaration is signed by the key being declared compromised — possible because the legitimate owner still holds that key; an attacker who copied it did not remove it. Registry nodes that receive a `COMPROMISE` declaration propagate it with higher urgency than a standard revocation, and clients receiving it immediately suspend trust in any message signed by the compromised key, regardless of when that message was sent.

The association rule is deliberately conservative: a `COMPROMISE` on the old key during the retention window flags the new key for re-verification, but does not automatically revoke it. This allows a legitimate owner who rotated cleanly and later discovered the old key was also stolen to recover the situation — they can publish a fresh signed rotation from the new key they hold, which they control and the attacker does not, resolving the re-verification flag. An attacker who pushed the new key cannot perform this recovery because they do not hold the new private key.

The same notification mechanism fires when a contact's address is deprovisioned by a domain authority — for example, when an employee leaves an organisation and their `@company.com` identity is revoked. Domain authority revocations are always treated as unsigned transitions requiring re-verification, since the individual's old key is no longer valid as an authorising signature. The domain authority revocation model that triggers this behaviour is specified in Section 13.3.


> **Key Change Handling**
> *A rotation signed by the old private key is accepted silently — only the legitimate key holder could have signed it. A rotation without that signature suspends delivery and requires the user to re-verify. If the old key is later found to have been stolen, the owner can declare it compromised during the seven-day retention window, flagging any rotation it signed for re-verification. Friction is applied where it is genuinely warranted.*


#### 14.1.3 Allowlist Portability and Backup


The allowlist is an asset of significant personal value — it
represents years of accumulated trust relationships. It is therefore
backed up as part of the user's encrypted key material and can be
exported in a standardised, encrypted format for migration between
clients or for archival. The export format includes not only the public
keys and addresses but also the full trust provenance record for each
entry, so that the history of how trust was established is preserved
across migrations.


### 14.2 The Pending Queue: Unknown but Unblocked Senders


The pending queue is not a list that senders are added to — it is where messages arrive by default when the sender is on neither the allowlist nor the blocklist. A sender reaches the pending queue simply by being unknown to the recipient: their cryptographic identity is registered and their message signature is valid, but the recipient has made no prior decision about them in either direction.


#### 14.2.1 Pending Queue Delivery Semantics


Messages held in the pending queue are presented in a section of the client visually distinct from the primary inbox. The client displays the sender's verified DMCN identity, the trust provenance of the sender in the network (whether any of the user's contacts have vouched for them, and if so how many), and a summary of the message sufficient to make an informed accept-or-reject decision — without requiring the user to fully open and read the message first.

From the pending queue the user has four options for each pending
message: Accept and allowlist the sender (promoting all future messages
to the primary inbox), Accept this message only (delivering the message
without allowlisting the sender), Reject and ignore (discarding the
message without any notification to the sender), or Reject and blocklist
(discarding the message and adding the sender to the blocklist to
prevent future delivery attempts).


#### 14.2.2 Pending Queue Auto-Resolution Rules


To reduce the burden of manual review, the client supports
configurable auto-resolution rules that can automatically promote or
demote senders before they reach the pending queue, based on network signals:

- Auto-promote if vouched by N or more allowlist contacts — configurable threshold, default of 3.

- Auto-promote if sender holds a verified organisational identity matching a domain the user has previously allowlisted.

- Auto-promote if the sender's identity has a reputation score above a configurable threshold in the user's chosen shared reputation feed.

- Auto-reject if the sender's identity appears on any blocklist feed the user has subscribed to.

These rules run at delivery time, before the message reaches the pending
queue, and are fully configurable. Users who want complete manual
control can disable all auto-resolution rules. Users who want a more
automated experience can enable conservative defaults that handle the
common cases without requiring intervention.


### 14.3 The Blocklist: Blocking Known Bad Actors


The blocklist is the user's registry of explicitly rejected senders.
Unlike a legacy email block — which can be trivially circumvented by
creating a new address — a DMCN blocklist entry is bound to a
cryptographic identity. A blocklisted sender cannot reach the user by
creating a new address, because their underlying key pair is what is
blocked, not the surface-level address string. This is a fundamentally
stronger guarantee than any blocking mechanism available in legacy
email.


#### 14.3.1 Personal Blocklist


The personal blocklist is private to the user and is never shared
externally. Adding a sender to the personal blocklist causes the DMCN
relay nodes handling the user's incoming messages to silently drop any
message signed by that identity before it reaches the user's device ---
the sender receives no delivery failure notification and no indication
that they have been blocked. This is consistent with the behaviour of
email blocking in major clients today and prevents the blocked sender
from using delivery failures as a signal to probe for workarounds.

Personal blocklist entries include the blocked identity's public key,
the address at which they were known, the date of blocking, and an
optional private note from the user recording their reason for the
block. This note is stored encrypted with the user's private key and is
never transmitted.


#### 14.3.2 Shared Reputation Feeds


Beyond the personal blocklist, the DMCN supports an opt-in shared
reputation feed system — a decentralised, community-maintained
registry of known bad actor public keys. This is the cryptographic
equivalent of the DNS-based blocklists (RBLs/DNSBLs) that legacy email
infrastructure has relied upon for decades, but with a critical
structural advantage: the identities being listed are cryptographic and
persistent.

In legacy email, a spammer who is listed on a blocklist can rotate to a
new IP address or domain within hours, effectively resetting their
reputation. In the DMCN, a public key that has been reported and listed
carries that reputation permanently. The key cannot be changed without
creating an entirely new identity — which requires going through the
account creation process again, imposing the same friction cost that
limits Sybil attacks. This asymmetry fundamentally favours the
defenders.


#### 14.3.3 Reputation Feed Architecture


Shared reputation feeds are operated independently of the DMCN core
protocol, allowing multiple competing feed operators to exist ---
analogous to how multiple DNS blocklist operators (Spamhaus, SORBS,
Barracuda) coexist today. Each feed operator maintains a signed,
distributed registry of reported public keys, with associated metadata
including the category of reported behaviour (spam, phishing,
harassment, malware distribution), the number of independent reports,
and the date of first and most recent report.

Users subscribe to one or more feeds through their DMCN client settings.
Feed data is retrieved and cached locally, so delivery decisions do not
require a real-time network lookup for every incoming message. Feed
operators publish their listing criteria, their dispute resolution
process, and their removal policy — users choose feeds whose policies
align with their needs.


#### 14.3.4 Reporting and Feed Contribution


Any user can submit a report against a sender's public key to a feed
they are subscribed to. The report is signed with the reporting user's
private key, providing cryptographic accountability for the report ---
false or malicious reports can be traced back to the reporter's
identity. This accountability mechanism is important: it discourages
coordinated campaigns to falsely blocklist legitimate identities,
because the reporters themselves are identifiable.

Feed operators implement their own thresholds and policies for when a
reported key is listed. A conservative operator might require reports
from twenty or more independent, verified identities before listing a
key. A more aggressive operator might list on fewer reports with greater
weighting for reports from highly trusted identities. Users select feed
operators whose threshold policies match their tolerance for false
positives versus false negatives.


#### 14.3.5 The Persistence Advantage


The most significant property of a cryptographic blocklist relative to
its legacy equivalents deserves explicit emphasis. When a DMCN identity
is reported and listed across multiple feeds, that listing is
effectively permanent for that key pair. The spammer's investment in
building a sending reputation — any messages that passed through
the pending queue, any contacts who user-approved them, any network vouching
they accumulated — is entirely lost. Starting over requires a new
identity, new account creation friction, and the same uphill
reputation-building process from scratch.

This is the economic property that makes spam structurally unviable in a
mature DMCN ecosystem. In legacy email, spamming is profitable because
the cost of reputation loss is near zero — a new domain or IP address
restores sending capability immediately. In the DMCN, reputation loss is
permanent per identity, and the cost of new identities, while low, is
non-zero and cumulative. At scale, this shifts the economics of spam
from profitable to unprofitable.


> **The Economic Argument Against Spam**
> *Legacy email: lose a sending address, acquire a new one in minutes
> at zero cost. DMCN: lose a cryptographic identity permanently,
> acquire a new one at non-zero cost with zero inherited reputation.
> Repeat at the scale required for spam economics to work, and the
> model collapses.*


### 14.4 Trust Tier Interaction Summary


| **Tier** | **Sender Type** | **Delivery Destination** | **Key Bound?** | **Shareable?** |
|---|---|---|---|---|
| Allowlist | Verified trusted contact | Primary inbox, immediate delivery | Yes — with provenance | Exportable (private) |
| Pending Queue | Verified but unknown sender | Pending queue, user review | No — state not stored | No |
| Personal Blocklist | Explicitly rejected sender | Silently dropped at relay | Yes — key blocked | No (private) |
| Shared Reputation Feed | Community-reported bad actor | Dropped per feed policy | Yes — persistent listing | Yes — community opt-in |

---


---

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

The distributed identity registry exposes five operations:

| Operation | Input | Output | Notes |
|---|---|---|---|
| `REGISTER` | `identity_record` | `ack` or `error` | Idempotent; re-registration updates the record if self-signature is valid |
| `LOOKUP` | `address: string` | `identity_record` or `not_found` | Rate-limited per source; see Section 18.3.1 |
| `REVOKE` | `address`, `revocation_signature` | `ack` or `error` | Permanent retirement of a key; revoked keys cannot be re-registered |
| `UPDATE` | `identity_record`, `old_key_signature` | `ack` or `error` | Key rotation; old key is retained for `retention_window_days` alongside new key; triggers allowlist notifications |
| `COMPROMISE` | `compromise_record` | `ack` or `error` | Declares a key as having been in hostile hands; propagated with high urgency; see below |

The `UPDATE` operation includes a `retention_window_days` field (default: 7, domain-authority-configurable maximum: 30). During the retention window, both the old and new keys are returned in `LOOKUP` responses with their respective roles clearly tagged. After the window expires, the old key is retired automatically.

The `COMPROMISE` operation uses the following record structure:

```
compromise_record {
    version:                  uint32
    address:                  string      // address of the identity
    compromised_pubkey:       bytes[32]   // Ed25519 public key being declared compromised
    declaration_timestamp:    uint64      // Unix timestamp of the declaration
    compromise_signature:     bytes[64]   // Ed25519 signature by the compromised key itself
    narrative:                string      // optional human-readable note; not displayed to third parties
}
```

The `compromise_signature` is produced by the key being declared compromised. This is possible because a key theft is a copy — the legitimate owner retains their copy and can sign the declaration. Registry nodes verify this signature before accepting the declaration.

**Propagation.** `COMPROMISE` declarations are propagated through the DHT with the same urgency as security-critical updates, rather than being batched with routine registry synchronisation. Clients that receive a compromise notification for a key they have stored — whether in an allowlist, a pending message, or a cached registry entry — must immediately surface an alert and suspend trust in that key.

**Retention window association rule.** If a `COMPROMISE` declaration is received for a key during its retention window — that is, the key was declared compromised after being used to sign a rotation — the newly rotated key is automatically flagged as requiring re-verification by all contacts. It is not automatically revoked, because the legitimate owner may have rotated legitimately and later discovered the old key was also compromised. The owner can resolve the re-verification flag by publishing a fresh `UPDATE` signed by the new key, demonstrating they hold it. An attacker who pushed the new key cannot perform this resolution.

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

## 17. Threat Model


This section provides a structured analysis of the threat landscape
facing the Decentralized Mesh Communication Network. For each identified
threat category, we document the nature of the attack, how it manifests
in the existing SMTP email ecosystem, how the DMCN architecture changes
the threat surface, and an honest assessment of whether the DMCN
improves, maintains, or introduces new risk relative to the status quo.


> **Methodology**
> *This threat model follows a structured adversarial analysis
> approach: for each threat, we identify the adversary's goal, the
> attack vector available under SMTP, the attack vector (if any)
> available under DMCN, and the net change in risk. Threats are grouped
> by adversary type: mass senders, targeted attackers, infrastructure
> attackers, and state-level actors.*


### 17.1 Threat Category 1: Spam and Bulk Unsolicited Messaging


#### 17.1.1 Nature of the Threat


Spam is the dominant form of abuse in the current email ecosystem. A
spam operator's goal is to deliver promotional, fraudulent, or
malicious messages to as many recipients as possible at the lowest
possible cost per delivery. The economics are straightforward: even at a
conversion rate of a fraction of a percent, the near-zero marginal cost
of sending email makes spam campaigns profitable.


#### 17.1.2 How SMTP Enables This Threat


SMTP imposes no cost and no identity requirement on senders. Any actor
with access to an SMTP server — whether legitimately operated,
compromised, or rented from a bulletproof hosting provider — can send
messages claiming to originate from any address. Existing mitigations
(SPF, DKIM, DMARC) have meaningfully reduced spoofing of established
domains, but have not eliminated the underlying problem: a spammer who
controls their own domain can produce fully authenticated messages that
pass all verification checks.

- No mandatory sender identity at the protocol level

- Negligible marginal cost per message at scale

- Address fabrication is trivial and carries no consequence

- Blocklisting based on IP or domain is easily circumvented by rotating infrastructure


#### 17.1.3 How DMCN Changes the Threat Surface


The DMCN eliminates the conditions that make spam economically viable at
the protocol level. Every sender must possess a registered cryptographic
identity, and every message must bear a valid signature from that
identity. Relay nodes reject unsigned or unverifiably signed messages
before they enter the network. This imposes a non-trivial, though low,
cost on account creation — and critically, each identity's reputation
is permanent and non-transferable.

A spam operator who wishes to send at scale must create a large number
of registered identities. Each identity that is reported and blocklisted
is permanently lost — there is no equivalent of rotating to a new IP
address. The mathematical relationship between spam volume and identity
cost shifts the economics of spam from profitable to uneconomical at
scale.


#### 17.1.4 Residual Risk and Honest Limitations


The DMCN does not make spam creation infinitely expensive — it makes
it non-zero in cost and permanently cumulative in consequence. A
determined, well-resourced spam operation could potentially automate the
account creation process (a Sybil attack), creating large numbers of
identities before they are reported. Section 17.5 addresses Sybil
resistance specifically. The consent-based inbox model (Section 8.2) provides a secondary layer: even a registered
identity cannot reach a user's primary inbox without meeting one of the
allowlisting criteria.


> **Net Assessment**
> *Spam: Significantly mitigated. The protocol-level economics of spam
> are fundamentally changed. Residual risk exists through Sybil attacks
> and requires robust account creation friction to fully close.*


### 17.2 Threat Category 2: Phishing and Identity Spoofing


#### 17.2.1 Nature of the Threat


Phishing attacks exploit the inability of email recipients to reliably
verify sender identity. An attacker impersonates a trusted entity — a
bank, an employer, a government agency — to induce the recipient to
reveal credentials, transfer funds, or install malware. Business Email
Compromise (BEC) is a sophisticated variant in which attackers
impersonate executives or financial officers within an organisation to
authorise fraudulent wire transfers. BEC alone causes billions of
dollars in losses annually.


#### 17.2.2 How SMTP Enables This Threat


Sender identity in SMTP is determined by the From header field, which is
a free-text string with no cryptographic binding to the actual sending
infrastructure. While DKIM signs the message content and SPF restricts
authorised sending IPs, neither prevents a spoofed display name, a
lookalike domain (e.g., paypa1.com instead of paypal.com), or a
compromised legitimate account from being used for phishing. The human
perception of trustworthiness is based entirely on the displayed From
field, which is trivially forged or manipulated.

- Display name spoofing requires no technical capability

- Lookalike domain registration costs a few dollars

- Compromised email accounts produce messages that pass all authentication checks

- No mechanism for the recipient to verify the message was written by a known human contact


#### 17.2.3 How DMCN Changes the Threat Surface


In the DMCN, every message carries a cryptographic signature tied to the
sender's registered key pair. It is mathematically impossible to forge
a message that appears to originate from a registered identity without
possessing that identity's private key. The recipient's client
verifies this signature automatically, and surfaces a clear,
non-technical trust indicator based on the sender's position in the
recipient's trust graph.

A phishing attempt from an unregistered identity cannot enter the
network. A phishing attempt from a registered identity that is not in
the recipient's allowlist lands in the pending queue, visibly
distinguished from trusted messages. The only viable phishing vector in
the DMCN is a fully registered identity that the recipient has
explicitly trusted — which requires the attacker to have established a
prior trust relationship.


#### 17.2.4 Account Compromise and the Key Binding Problem


The DMCN does not eliminate the threat of account compromise — it
changes its character. In SMTP, a compromised email account allows an
attacker to send messages that are indistinguishable from legitimate
messages from that sender. In the DMCN, a compromised account requires
the attacker to have stolen the private key itself, not merely the login
credentials. Private keys stored in hardware-backed secure enclaves (as
specified in Section 7.1) cannot be extracted even if the device's
operating system is compromised.

This represents a meaningful improvement over SMTP account compromise,
but it introduces a new concern: if a private key is stolen (e.g., from
a device without hardware security support), the attacker gains the full
trust relationships of that identity with no visible indicator to
contacts. The allowlist key-change notification system (Section 14.1.2)
partially mitigates this: if the attacker uses a new device, contacts
will be alerted that the key has changed and prompted to re-verify.


> **Net Assessment**
> *Phishing: Substantially mitigated for protocol-level spoofing.
> Account compromise risk shifts from credential theft to key theft ---
> a meaningfully higher bar, but not zero. Hardware key storage is
> essential for this property to hold.*


### 17.3 Threat Category 3: Infrastructure Attacks


#### 17.3.1 Denial of Service Against the Network


Any distributed network is a potential target for denial-of-service
attacks. In the DMCN, the primary infrastructure targets are relay
nodes, the distributed identity registry, and bridge nodes. The goal of
such an attack is to prevent message delivery, disrupt identity lookups,
or degrade the network to the point of unusability.


#### 17.3.2 Comparison with SMTP


SMTP email infrastructure is frequently the target of distributed
denial-of-service attacks. Major email providers defend against these
attacks through substantial investment in distributed infrastructure,
traffic scrubbing, and rate limiting. Smaller operators are more
vulnerable. The centralised nature of many email providers means that a
successful attack on a major provider affects a large fraction of global
email.

The DMCN's peer-to-peer architecture distributes this attack surface.
There is no single point of failure equivalent to a major email
provider's infrastructure. An attacker wishing to disrupt the network
must simultaneously target a significant fraction of all relay nodes ---
a substantially harder target than attacking a centralised mail server.


#### 17.3.3 New Infrastructure Risks Introduced by DMCN


The distributed identity registry represents a novel attack surface with
no direct SMTP equivalent. A successful attack that corrupts or makes
unavailable a significant portion of the identity registry could prevent
new message delivery (recipients' public keys cannot be resolved) or,
in a more sophisticated attack, allow injection of false key mappings.
The cryptographic verification of registry entries (each entry is signed
by the identity owner's private key) provides strong resistance to the
latter — a false key cannot be injected without the private key of the
identity being spoofed.

Bridge nodes represent a concentration of trust and traffic that may be
attractive targets. A successful attack on a widely-used bridge node
disrupts both inbound and outbound legacy email communication for its
users. The federated bridge architecture (Section 11.5) distributes this
risk, but organisations using a single bridge provider remain exposed to
single-point-of-failure risk.


> **Net Assessment**
> *Infrastructure DoS: Comparable to SMTP for distributed attacks;
> improved for centralised attack scenarios due to peer-to-peer
> architecture. The identity registry is a novel attack surface
> requiring careful design. Bridge nodes reintroduce some
> centralisation risk during the transition period.*


### 17.4 Threat Category 4: Relay Node Misbehaviour


#### 17.4.1 Nature of the Threat


Unlike centralised email providers, DMCN relay nodes can be operated by
any party. This raises the question of what happens when a relay node
operator acts maliciously or negligently. Potential misbehaviours
include selectively dropping messages, logging message metadata (even
though content is encrypted), injecting false routing information, or
colluding with other nodes to deanonymise communication patterns.


#### 17.4.2 Comparison with SMTP


In the existing email ecosystem, routing trust is placed in a chain of
mail transfer agents (MTAs) that may or may not be operated by the
ultimate sender and recipient. SMTP messages in transit are visible to
any MTA in the delivery chain — historically in plaintext, though TLS
adoption has significantly improved transport encryption. However, TLS
between MTAs does not guarantee end-to-end confidentiality; each MTA can
read message content.

The DMCN improves on this substantially: relay nodes cannot read message
content because messages are encrypted with the recipient's public key,
which no relay node possesses. What relay nodes can observe is message
metadata — the pseudonymous identities of sender and recipient (as
public keys), message size, and timing. This is a meaningful improvement
over SMTP, where message content is accessible to routing
infrastructure.


#### 17.4.3 Metadata Privacy and the Onion Routing Layer


The proposed onion-routing-inspired transport protocol (Section 6.2.2)
is specifically designed to limit the metadata visibility of individual
relay nodes. In an onion routing scheme, each relay node knows only the
previous hop and the next hop — it does not know both the originating
sender and the ultimate recipient. This prevents a single malicious
relay node from observing a complete communication relationship.

However, onion routing is not a complete solution. A global passive
adversary — one that can observe a significant fraction of all network
traffic — may be able to correlate message timing and sizes across
multiple hops to deanonymise communication patterns. This is a known
limitation of onion routing schemes and is the same threat faced by
networks such as Tor. For most threat models relevant to a
general-purpose communication network, this level of sophistication is
beyond the realistic adversary.


> **Net Assessment**
> *Relay node misbehaviour: Improved relative to SMTP. Content
> confidentiality is guaranteed regardless of relay node behaviour.
> Metadata privacy is meaningfully improved through onion routing but
> not absolute. Global passive adversaries represent a residual risk
> for high-sensitivity use cases.*


### 17.5 Threat Category 5: Sybil Attacks


#### 17.5.1 Nature of the Threat


A Sybil attack occurs when a malicious actor creates a large number of
fake identities to subvert a trust-based system. In the context of the
DMCN, the primary Sybil attack scenarios are: creating large numbers of
identities to conduct spam campaigns before they are blocklisted;
creating fake identities to inflate web-of-trust vouching for a
malicious identity; and creating fake identities to manipulate shared
reputation feeds.


#### 17.5.2 Comparison with SMTP


SMTP is essentially infinitely susceptible to Sybil attacks — there is
no meaningful identity system to attack, and the cost of registering a
new sending domain is a few dollars. The DMCN's identity model is
inherently more resistant because it requires account creation friction,
and because blocklisted identities cannot be recovered. However, the
DMCN is not immune, and Sybil resistance is one of the most significant
open design challenges.


#### 17.5.3 Proposed Mitigations


Several mechanisms can be combined to raise the cost of Sybil attacks to
uneconomical levels:

- Account creation friction: requiring a verified phone number, email address, or proof-of-work computation during account creation raises the per-identity cost above zero

- Rate limiting on new identity behaviour: newly created identities are subject to stricter pending queue treatment and lower throughput limits until they have established a reputation history

- Web-of-trust bootstrapping requirements: the consent-based inbox model means that a new identity must earn its way into recipients' allowlists; a Sybil farm of identities with no trust relationships has no delivery capability

- Reputation feed correlation: feed operators can flag clusters of identities with similar creation timing, device fingerprints, or behaviour patterns as likely Sybil farms


> **Net Assessment**
> *Sybil attacks: Improved relative to SMTP, but not fully solved. The
> per-identity cost is non-zero and cumulative in consequence, raising
> the economics of Sybil attacks. Full resistance requires careful
> design of account creation friction — a balance between security
> and accessibility that requires user research and iteration.*


### 17.6 Threat Category 6: State-Level Surveillance and Censorship


#### 17.6.1 Nature of the Threat


Nation-state actors represent the most sophisticated and well-resourced
adversaries in the threat landscape. Their objectives may include mass
surveillance of communications content, targeted surveillance of
specific individuals, censorship of communications between specific
parties, or disruption of communication infrastructure for geopolitical
purposes.


#### 17.6.2 How SMTP Enables State Surveillance


SMTP email is extraordinarily accessible to state surveillance. In most
jurisdictions, intelligence agencies have legal authority to compel
email providers to hand over message content and metadata. Beyond legal
compulsion, the centralised infrastructure of major email providers
represents a small number of high-value interception points.
Transport-level encryption (TLS) protects messages in transit between
MTAs, but messages are stored in plaintext at the provider and are
accessible through legal process, compromise of provider systems, or
provider cooperation.

Countries that control significant internet routing infrastructure can
also conduct passive surveillance of SMTP traffic at the network level,
particularly for messages between providers that do not enforce
opportunistic TLS.


#### 17.6.3 How DMCN Changes the Threat Surface


The DMCN substantially increases the difficulty and cost of mass
surveillance. Messages are encrypted end-to-end with the recipient's
public key; there is no centralised service provider that can be
compelled to hand over message content in plaintext. A state actor
wishing to read DMCN messages must either obtain the recipient's
private key (through device seizure, legal compulsion, or compromise of
the device's secure enclave) or conduct a cryptanalytic attack against
the underlying elliptic curve cryptography — which is currently
considered computationally infeasible.

Metadata surveillance is more feasible. Even with onion routing, a state
actor controlling significant network infrastructure can observe which
IP addresses are communicating with DMCN relay nodes and can attempt
traffic correlation attacks. The DMCN's pseudonymous identity model
(public keys rather than real names) provides some protection, but is
not equivalent to anonymity.


#### 17.6.4 Censorship and Network Disruption


A state that wishes to prevent DMCN communication within its
jurisdiction has several options: blocking IP addresses of known relay
nodes (analogous to how Tor exit nodes are blocked in some countries);
requiring ISPs to perform deep packet inspection to identify and block
DMCN traffic; or compelling local device manufacturers to remove or
disable DMCN client applications. These are the same techniques used
against other encrypted communication platforms (Signal, WhatsApp) and
represent a cat-and-mouse dynamic between the network and censoring
states rather than a fundamental architectural vulnerability.

The peer-to-peer architecture of the DMCN provides more censorship
resistance than centralised email providers, because there is no single
point that can be blocked or compelled. However, the bridge architecture
--- which provides interoperability with legacy email — does introduce
centralisation points during the transition period that could be
targeted for censorship or legal compulsion.


> **Net Assessment**
> *State-level surveillance: Substantially improved relative to SMTP
> for content privacy. Mass content surveillance is no longer feasible
> against a correctly implemented DMCN. Metadata surveillance and
> censorship remain possible for well-resourced state actors but
> require significantly more effort and sophistication than SMTP
> interception. Bridge nodes represent a transitional vulnerability.*


### 17.7 Threat Category 7: Key Compromise and Recovery Attacks


#### 17.7.1 Nature of the Threat


The security of the entire DMCN identity model rests on the secrecy of
each user's private key. If a private key is compromised, the attacker
gains the ability to read all future messages sent to that identity, to
send messages that appear to come from that identity, and to modify or
revoke that identity's trust relationships. Key compromise is therefore
the most serious category of attack specific to a cryptographic identity
system.


#### 17.7.2 Comparison with SMTP


SMTP account security is typically based on password authentication.
Password-based accounts are vulnerable to credential stuffing, phishing
of login credentials, database breaches at the email provider, and
SIM-swapping attacks against phone-based multi-factor authentication. In
practice, large-scale email account compromise is common and inexpensive
for attackers. The consequences are significant — a compromised email
account often provides access to password reset flows for other services
--- but the compromised account can be recovered through
provider-managed account recovery.

The DMCN's private key model changes this threat in important ways.
Hardware-backed key storage (Secure Enclave, TPM) substantially raises
the bar for key extraction — private keys in secure hardware cannot be
exported even with full device access. However, the consequences of a
key compromise that does occur are more severe: there is no centralised
provider who can reset an account.


#### 17.7.3 The Social Recovery Attack Surface


The social recovery mechanism (Section 7.3) — in which trusted
contacts hold encrypted shards of a recovery key — introduces a new
attack surface. An attacker who wishes to compromise an account could
target the recovery contacts rather than the primary user, attempting to
compromise enough contacts to meet the recovery threshold. This is a
form of social engineering attack that is specific to threshold recovery
systems.

Mitigations include: requiring that recovery contacts independently
verify the identity of the person requesting recovery (e.g., via a video
call) before releasing their shard; implementing time delays on recovery
requests to allow the legitimate user to be notified and object; and
limiting the recovery mechanism to account access restoration rather
than providing a path to re-issue the underlying key pair.


#### 17.7.4 Key Revocation and Forward Secrecy


The whitepaper's current architecture does not fully specify key
revocation and forward secrecy mechanisms. These are important open
questions. If a key is compromised, the affected user must be able to
publish a revocation that is propagated through the identity registry,
invalidating future messages signed by the old key. Forward secrecy ---
the property that compromise of a long-term key does not expose
historical messages — requires additional protocol design, such as
ephemeral session keys derived from the long-term identity key for each
conversation.


> **Net Assessment**
> *Key compromise: The bar for key theft is substantially higher than
> password theft, particularly with hardware-backed key storage.
> However, the consequences of compromise are also higher, and recovery
> mechanisms introduce new attack surfaces. Forward secrecy and key
> revocation are open design questions that must be resolved before a
> production deployment.*


### 17.8 Threat Category 8: Bridge Node Attacks


#### 17.8.1 Nature of the Threat


The SMTP-DMCN bridge architecture (Section 10) is a necessary component
of any viable transition strategy, but it reintroduces several trust and
security challenges that the native DMCN architecture otherwise
eliminates. Bridge nodes represent the interface between the trustless
SMTP world and the cryptographically verified DMCN, and they must make
trust determinations about SMTP senders that have no direct equivalent
in the native protocol.


#### 17.8.2 Bridge-Specific Attack Vectors


The following attacks are specific to the bridge architecture and have
no equivalent in the native DMCN:

- Content interception on the outbound path: bridge nodes must decrypt outbound messages to re-encode them as SMTP. A malicious or compromised bridge operator gains access to message content in transit. This is disclosed in Section 11.2.2 and is an unavoidable consequence of protocol translation.

- False trust classification: a malicious bridge could misrepresent the trust tier of an inbound SMTP message — for example, classifying a spam message as 'Verified Legacy Sender' to bypass the recipient's filters. The bridge's classification is signed with the bridge's own DMCN key, creating accountability, but this only helps if users verify which bridge they are trusting.

- Legacy spam injection: spammers may target the bridge's SMTP listener as a pathway into DMCN inboxes, attempting to exploit weaknesses in the bridge's SMTP authentication checks to inject messages that would be rejected if sent natively.

- Bridge impersonation: an attacker could operate a bridge that presents itself as trustworthy in the identity registry but maliciously handles messages. Mitigated by requiring bridge operators to publish their security practices and undergo periodic audits.


#### 17.8.3 Bridge Risk as a Transitional Concern


It is important to contextualise bridge risks appropriately. Bridge
nodes handle only the traffic that crosses between SMTP and DMCN; native
DMCN-to-DMCN communication is not affected by bridge security
properties. As the DMCN user base grows and a larger fraction of
communication is native, the fraction of traffic passing through bridge
nodes diminishes. The bridge architecture is explicitly designed as a
transitional mechanism, not a permanent feature.


> **Net Assessment**
> *Bridge attacks: The bridge architecture necessarily reintroduces
> SMTP-era trust challenges for legacy communication paths. These risks
> are bounded, disclosed, and diminish as native DMCN adoption grows.
> The bridge is a transitional vulnerability, not a permanent
> architectural weakness.*


### 17.9 Threat Model Summary


The table below summarises each threat category, the current severity in
SMTP, the treatment under DMCN, and the net outcome for each:

| **Threat Category** | **SMTP Severity** | **DMCN Treatment** | **Net Outcome** |
|---|---|---|---|
| Spam / Bulk Messaging | Critical — protocol endemic | Protocol-level identity cost eliminates economic viability | **Significantly Reduced** |
| Phishing / Spoofing | Critical — trivially executed | Cryptographic signing makes spoofing mathematically infeasible | **Significantly Reduced** |
| Infrastructure DoS | High — centralised targets | Distributed architecture reduces single-point risk | **Partially Mitigated** |
| Relay Node Misbehaviour | High — plaintext in transit | End-to-end encryption limits relay visibility to metadata | **Partially Mitigated** |
| Sybil Attacks | N/A — no identity system | Non-zero identity cost; permanent reputation loss | **Partially Mitigated** |
| State Surveillance | Critical — provider access | End-to-end encryption; no centralised interception point | **Significantly Reduced** |
| Key Compromise | High — passwords weak | Hardware keys raise bar; recovery introduces new surface | **Partially Mitigated** |
| Bridge Attacks | N/A — DMCN-specific | Bounded to legacy traffic; diminishes with adoption | **Partially Mitigated** |


> **Overall Assessment**
> *The DMCN architecture represents a meaningful and substantial
> improvement over SMTP across every threat category where comparison
> is possible. The threats that remain partially mitigated rather than
> eliminated — Sybil resistance, relay metadata leakage, key recovery
> attacks, and bridge node risks — are well-understood engineering
> challenges with known mitigation approaches, rather than fundamental
> architectural flaws. None of these residual risks represents a
> regression relative to the current SMTP-based email ecosystem.*
>
> *Organisations deploying the Domain Authority Record model (Section 13) introduce an additional threat surface: compromise of the domain authority key. This threat is analogous to enterprise CA root key compromise and should be treated with equivalent operational rigour. Mitigations are specified in Section 13.7.*

---


---

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

The DMCN client should implement a full data export function that produces a portable, encrypted archive of the user's message history, allowlist, blocklist, and trust attestations in a documented, open format. This export serves both the regulatory compliance function and the practical function of enabling migration between DMCN client applications without loss of data.

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

## 19. Open Challenges and Research Questions


This whitepaper represents a preliminary investigation into the design
space of a Decentralized Mesh Communication Network. Several significant
challenges remain open and will require further research, prototyping,
and community input.


### 19.1 Scale and Performance


The distributed identity registry and peer-to-peer routing architecture
must be demonstrated to perform adequately at the scale of global email
--- billions of users, hundreds of billions of messages per day. The
performance characteristics of the proposed architecture under realistic
load conditions must be validated through simulation and prototype
deployment.


### 19.2 Key Recovery Without Central Authority


The social recovery model proposed in Section 7 is promising, but its UX
and security properties require careful design and user research. The
threshold for recovery must balance security against the practical
reality that trusted contacts may be unavailable, may themselves lose
access to their accounts, or may be compromised. Alternative recovery
mechanisms should be investigated and compared.


### 19.3 Regulatory and Legal Compliance


End-to-end encrypted communication creates compliance challenges for
regulated industries — financial services, healthcare, law — that
are required to maintain records of communications and provide them to
regulators on demand. The DMCN architecture must provide mechanisms that
allow regulated entities to meet their compliance obligations without
compromising the security properties of the system for other users.

The `ARCHIVE_REQUIRED` policy flag in the Domain Authority Record model (Section 13.5) provides the structural mechanism for compliance archiving under managed domains: outbound messages are BCC-encrypted to an approved archive bridge at send time, giving the domain authority access to a decryptable record without exposing message content to relay infrastructure. Edge cases around legal holds, cross-jurisdiction archiving, and eDiscovery production workflows remain open and require further engagement with compliance counsel and regulators in target jurisdictions.


### 19.4 Governance


A truly decentralized network requires a governance model that allows
the protocol to evolve — to address security vulnerabilities,
incorporate technical improvements, and respond to regulatory changes
--- without any central authority having unilateral control. The
governance model for the DMCN is a critical design question with
significant implications for the network's long-term resilience and
trustworthiness.


### 19.5 Sybil Resistance


While the DMCN's identity model prevents spam from unregistered
senders, it must also resist Sybil attacks — scenarios in which a
malicious actor creates a large number of registered identities to
overwhelm spam defenses. The account creation process must impose
sufficient cost or friction to make large-scale Sybil attacks
uneconomical, without imposing unacceptable burden on legitimate users.


### 19.6 Route Selection Algorithm

The onion routing transport layer specified in Section 15.4.4 requires each client to select three relay nodes for every outbound message path. The hard constraints on this selection — operator diversity, ASN diversity, subnet diversity — are well-defined. The weighted scoring algorithm that balances latency, geographic distribution, reputation, and node capacity within those constraints is not yet empirically validated.

Several open questions require prototype deployment and measurement to resolve:

**Guard node policy.** Section 15.4.4 recommends a stable set of preferred entry nodes rotated every 30 days, following Tor's guard node model. The security benefit of guard nodes — preventing an adversary from gaining statistical entry position through repeated random selection — must be weighed against the privacy cost: a stable entry node that is compromised or surveilled provides persistent observation of the sender's traffic patterns over the rotation period. The optimal rotation interval and guard set size for the DMCN's threat model require empirical analysis.

**Latency versus diversity trade-off.** Geographic diversity in route selection increases the number of distinct network vantage points an adversary requires for traffic correlation, but also increases cumulative path latency. For users in geographically peripheral locations — where the nearest relay nodes are clustered in a small number of regions — enforcing strict continental diversity may impose unacceptable latency penalties. The protocol needs a principled approach to relaxing diversity constraints when strict application would degrade the user experience below acceptable thresholds.

**Relay node directory freshness.** The client's local relay node directory is refreshed on a schedule from the identity registry. A stale directory may include nodes that have gone offline, been delisted, or had their reputation updated since the last refresh. The interaction between directory staleness, route construction failures, and fallback selection behaviour requires careful specification to avoid creating exploitable patterns — for example, an adversary who can predict which nodes a client will fall back to when its preferred nodes are unavailable.

**Adversarial node concentration.** An adversary who operates a large number of relay nodes can increase their probability of appearing in a given path even with operator diversity constraints, by registering nodes under many different operator identities. The extent to which ASN and subnet diversity constraints mitigate this — and whether additional heuristics such as maximum per-operator traffic share are warranted — requires analysis of realistic adversary node deployment strategies.

These questions are research problems as much as engineering problems, and their resolution should draw on the academic literature on anonymous communication network design. Community engagement with the anonymity research community is a recommended early step before the route selection algorithm is finalised.

---


---

## 20. Conclusion


Email is the foundational communication layer of the digital world, and
it is broken in ways that incremental fixes cannot repair. Spam,
phishing, and identity fraud are not aberrations — they are structural
consequences of a protocol designed for a world of trusted academic
networks, deployed in a world of adversarial global actors. The
mitigations that have accumulated over five decades have managed the
symptoms without addressing the cause.

The cause is the absence of cryptographic identity at the protocol
level. In a network where anyone can claim to be anyone and sending a
message costs nothing, spam is not a problem to be solved — it is a
rational economic behavior to be expected. No filtering system, however
sophisticated, can permanently resolve a problem that is embedded in the
architecture.

The Decentralized Mesh Communication Network proposed in this whitepaper
addresses the root cause. By making cryptographic identity mandatory ---
not optional, not recommended, not a premium feature — and by
enforcing sender verification at the protocol level rather than at the
filtering layer, the DMCN eliminates the conditions under which spam is
possible rather than trying to detect spam after it has been generated.

The key insight that distinguishes this proposal from prior attempts is
the recognition that cryptographic complexity and user simplicity are
not in conflict. The same elliptic curve cryptography that underpins
cryptocurrency wallets and modern TLS connections also underpins Apple
Passkeys — an authentication technology that hundreds of millions of
users interact with daily, without any awareness that they are using
public/private key cryptography. The technology is not the barrier. The
user experience design is the barrier, and it is a solvable engineering
problem.

The DMCN is not proposed as a finished design — it is proposed as a
research agenda and a design direction. The open challenges documented
in Section 19 are real and significant. The competitive landscape
documented in Section 4 demonstrates that the market is beginning to
recognize the problem space, even if existing solutions have not yet
solved it effectively.

The opportunity exists for a solution that combines the cryptographic
soundness of PGP, the decentralized architecture of peer-to-peer
networks, the user experience standards of modern consumer applications,
and a credible migration path from the existing email ecosystem. The
DMCN is a proposal toward that solution.


> **Next Steps**
> *This whitepaper represents Version 0.2 of an ongoing investigation.
> Subsequent versions will incorporate technical prototyping results,
> user research findings, and engagement with the broader cryptographic
> and communications research community. Feedback, critique, and
> collaboration are actively solicited.*

---

*End of Document*

---


---

## Glossary

Terms are listed alphabetically. Where a term has a common abbreviation used in this document, the abbreviation is shown in parentheses.

---

**Blocklist**
A user-maintained or community-maintained registry of cryptographic identities that are explicitly blocked from delivering messages. In the DMCN, blocklist entries are bound to public keys rather than surface addresses, making them impossible to circumvent by simply creating a new address string. See also: *Pending Queue*, *Allowlist*, *Shared Reputation Feed*.

---

**Bridge Node**
A DMCN-native network node that additionally speaks SMTP, allowing messages to pass between the DMCN and the legacy email ecosystem. Bridge nodes handle protocol translation in both directions — outbound (DMCN to SMTP) and inbound (SMTP to DMCN) — and are registered in the DMCN identity layer with their own cryptographic identity. Bridge nodes are a transitional component; their role diminishes as native DMCN adoption grows.

---

**Business Email Compromise (BEC)**
A category of email fraud in which an attacker impersonates a trusted party — typically an executive, financial officer, or business partner — to authorise fraudulent wire transfers or obtain sensitive information. BEC is consistently ranked by the FBI among the highest-impact cybercrime categories by financial loss. It exploits the absence of cryptographic sender verification in SMTP.

---

**Curve25519**
An elliptic curve used for public-key cryptography, widely regarded as one of the most secure and performant curves available. Curve25519 is the basis for the X25519 key exchange protocol used in modern TLS, Signal, and many other security systems. It is one of the candidate curve families for the DMCN identity layer.

---

**Compromise Declaration**
A signed registry operation by which a DMCN identity owner declares that their private key has been stolen or is otherwise in hostile hands. Unlike a standard revocation (which merely retires a key from active use), a compromise declaration carries the semantic that the key should be treated as having been controlled by an attacker. The declaration is signed by the compromised key itself — possible because key theft is a copy, not a removal, so the legitimate owner retains their copy. Registry nodes propagate compromise declarations with high urgency. If a compromise declaration is issued during the key rotation retention window, any rotation signed by the compromised key is automatically flagged for contact re-verification. See also: *Key Rotation Retention Window*, *Revocation*.

---

**Key Rotation Retention Window**
A period — default seven days, domain-authority-configurable up to 30 — during which both the old and new primary keys coexist in the identity registry following a signed key rotation. The retention window exists to give the legitimate owner time to detect and respond to a key theft: if they discover the old key was stolen after rotating, they can publish a Compromise Declaration against it during this window, which flags the rotated-to key for re-verification by contacts. After the window expires, the old key is automatically retired. See also: *Compromise Declaration*, *Key Rotation*.

---

**DANE (DNS-based Authentication of Named Entities)**
An internet standard that allows domain owners to publish cryptographic key material directly in DNS records secured by DNSSEC, creating a chain of trust from the DNS root to a specific certificate or key. The DMCN's highest-assurance address verification tier uses a DANE-style model to bind email addresses to public keys without relying on a certificate authority.

---

**Decentralized Identifier (DID)**
A type of identifier defined by the W3C that enables verifiable, self-sovereign digital identity without dependence on any central registry or authority. DIDs are used by some blockchain-based identity systems (including Dmail Network) as the basis for user identity. The DMCN's identity model shares the self-sovereign property of DIDs but does not require blockchain infrastructure.

---

**Distributed Hash Table (DHT)**
A peer-to-peer data structure that distributes storage and lookup of key-value pairs across a network of nodes, with no central coordinator. Each node stores a portion of the data and can resolve any query through a bounded number of hops. The BitTorrent protocol uses a DHT for peer discovery. The DMCN uses a DHT-like structure for its distributed identity registry.

---

**DKIM (DomainKeys Identified Mail)**
An email authentication standard that allows a sending mail server to attach a cryptographic signature to outgoing messages, enabling receivers to verify that the message content was not altered in transit and that it was sent by a server authorised by the domain owner. DKIM addresses message integrity but not spam: a spammer controlling their own domain can produce fully valid DKIM signatures.

---

**DMARC (Domain-based Message Authentication, Reporting and Conformance)**
A policy framework built on top of SPF and DKIM that allows domain owners to specify how receiving servers should handle messages that fail authentication checks, and to receive reports on authentication outcomes. DMARC has improved domain-level spoofing resistance but does not address individual sender identity or eliminate spam from authenticated domains.

---

**Domain Authority Record (DAR)**
A signed registry entry published by a domain owner that declares the domain's administrative policy and the public key of the authority responsible for enforcing it. The DAR enables organisations to control which identities are authorised under their domain, provision and deprovision staff addresses, enforce compliance archiving, and delegate administrative authority to sub-administrators. Identity records under a managed domain carry a domain countersignature issued by the domain authority alongside the individual's own self-signature. See Section 13.

---

**Domain Countersignature**
A cryptographic signature applied to an individual identity record by a domain authority, certifying that the address has been provisioned through the organisation's authorised process. Relay nodes and clients check for a valid domain countersignature when the domain's Domain Authority Record sets the `REQUIRE_DOMAIN_COUNTERSIG` policy flag. An identity record without a valid countersignature under a managed domain is treated as unverified regardless of the individual's own self-signature. The domain authority can withdraw the countersignature at any time through a domain revocation record, immediately depro­visioning the address. See Section 13.2 and Section 13.3.

---

**DNSBL (DNS-based Blocklist)**
A system that publishes lists of IP addresses or domains known to be sources of spam or other abuse, in a format queryable via DNS. Mail servers use DNSBLs to reject or flag messages from listed senders. Well-known DNSBLs include those operated by Spamhaus and Barracuda. In the DMCN context, shared reputation feeds serve an analogous function but list cryptographic public keys rather than IP addresses.

---

**DNSSEC (Domain Name System Security Extensions)**
A suite of extensions to DNS that adds cryptographic signatures to DNS records, allowing resolvers to verify that records have not been tampered with in transit. DNSSEC is required for the highest-assurance DMCN address verification tier (DANE-style binding), as it provides a tamper-evident chain of trust from the DNS root to the domain's published key material.

---

**Elliptic Curve Cryptography**
A form of public-key cryptography based on the algebraic structure of elliptic curves over finite fields. Elliptic curve key pairs provide equivalent security to RSA at significantly smaller key sizes, making them well-suited to constrained environments such as mobile devices. The DMCN identity layer uses elliptic curve key pairs (specifically the Curve25519 or secp256k1 families) for all user identities.

---

**End-to-End Encryption (E2EE)**
An encryption model in which messages are encrypted by the sender and can only be decrypted by the intended recipient. Intermediate infrastructure — relay nodes, bridge nodes, service providers — cannot read the message content. The DMCN provides end-to-end encryption for all native DMCN-to-DMCN messages by encrypting each message to the recipient's public key before it enters the transport layer.

---

**Forward Secrecy**
A cryptographic property whereby compromise of a long-term private key does not expose previously recorded encrypted communications. Forward secrecy is achieved by using ephemeral session keys — short-lived keys derived for each conversation or message — so that each session's encryption is independent. Forward secrecy is an open design question in the current DMCN specification and is noted as a required feature before production deployment.

---

**Pending Queue**
The default destination in the DMCN client for messages from senders who appear on neither the recipient's allowlist nor their blocklist. The pending queue is not a curated list — senders arrive there by the absence of any prior decision about them, not by being added to anything. Messages in the pending queue are held for user review, displayed with the sender's verified cryptographic identity and any network trust signals available. The recipient can allowlist the sender, accept the individual message, ignore it, or blocklist the sender. See also: *Allowlist*, *Blocklist*.

---

**Key Pair**
A matched pair of cryptographic keys — a public key and a private key — generated together using a mathematical relationship such that data encrypted with the public key can only be decrypted with the corresponding private key, and data signed with the private key can be verified with the public key. In the DMCN, every user identity is represented by a key pair; the public key is published in the identity registry, while the private key never leaves the user's device.

---

**Key Encapsulation Mechanism (KEM)**
A cryptographic pattern in which a message payload is encrypted once with a randomly generated symmetric content key (CEK), and the CEK is then individually wrapped (encrypted) for each intended recipient using that recipient's public key. Any recipient who holds the corresponding private key can unwrap the CEK and decrypt the payload. The KEM pattern ensures that large payloads are transmitted exactly once regardless of how many recipients or enrolled devices are involved, with only small per-recipient overhead for the wrapped key material. The DMCN uses a KEM pattern for all message and attachment encryption (Section 15.3.3). The approach is standardised in RFC 9180 (HPKE).

---

**Mesh Network**
A network topology in which each node can relay data for the network, with no single central hub through which all traffic must pass. Mesh networks are resilient — the failure of any individual node does not disrupt overall connectivity — and resistant to centralised censorship or surveillance. The DMCN uses a mesh routing topology for message delivery.

---

**MX Record (Mail Exchanger Record)**
A type of DNS record that specifies the mail server responsible for accepting email for a given domain. When a sending server wants to deliver email to user@example.com, it looks up the MX record for example.com to find the destination server. In the DMCN bridge architecture, bridge nodes are published via MX records so that legacy SMTP senders can reach DMCN users.

---

**NFT (Non-Fungible Token)**
A type of blockchain-based digital asset representing unique ownership of a specific item. Dmail Network uses NFT domains as one form of user identity. The DMCN does not use NFTs or blockchain infrastructure; this term appears in the document only in the competitive analysis of Dmail.

---

**Onion Routing**
A technique for anonymous communication over a network in which messages are encrypted in multiple layers and routed through a sequence of relay nodes, each of which decrypts one layer to learn only the next hop — no single node knows both the origin and the destination. Tor is the most well-known onion routing implementation. The DMCN transport layer uses an onion-routing-inspired protocol to provide metadata privacy in addition to content encryption.

---

**Peer-to-Peer (P2P)**
A network architecture in which participants communicate directly with each other rather than through a central server. Each peer can act as both client and server. P2P networks are resilient and decentralised by design. The DMCN is a peer-to-peer network: messages are routed through a distributed mesh of nodes with no central routing authority.

---

**PKI (Public Key Infrastructure)**
A framework for managing public-key cryptography at scale, typically involving certificate authorities (CAs) that issue signed certificates binding public keys to identities. S/MIME email encryption relies on PKI. PKI introduces centralisation — a compromised or malicious CA can issue fraudulent certificates — which the DMCN avoids through its decentralised identity registry.

---

**Primary Key**
The canonical cryptographic key pair representing a DMCN identity. There is exactly one active primary key per address at any point in time. It is published in the identity registry, anchors the trust graph, and is the key against which allowlist bindings are made. The primary key is not used for routine per-message operations; that role is delegated to device sub-keys. See also: *Device Sub-Key*, *Key Pair*.

---

**Device Sub-Key**
A subordinate key pair generated on a specific device and signed by the identity's primary key. Sub-keys are the keys used for day-to-day message signing and decryption on a given device. Multiple active sub-keys may exist simultaneously — one per enrolled device. A sub-key can be revoked independently when a device is lost, decommissioned, or rotated, without affecting the primary key or other devices. Senders encrypt messages to all active sub-keys so that the recipient can decrypt on whichever device they first open the message. See also: *Primary Key*, *Key Pair*.

---

**Private Key**
The secret half of a cryptographic key pair, kept exclusively on the user's device and never transmitted. Private keys are used to sign outgoing messages (proving they came from the key's owner) and to decrypt incoming messages (which were encrypted to the corresponding public key). The security of the entire DMCN identity model rests on the secrecy of private keys; they are stored in hardware-backed secure enclaves wherever the device supports it.

---

**Public Key**
The shareable half of a cryptographic key pair, published in the DMCN identity registry so that anyone wishing to send a message can encrypt it to the recipient and verify the sender's signatures. Knowing someone's public key does not allow an attacker to impersonate them or decrypt their messages — only the corresponding private key can do those things.

---

**RBL (Real-time Blocklist)**
A DNS-based list of IP addresses known to originate spam, queried in real time by mail servers to make delivery decisions. RBLs are a long-established component of legacy email spam defence. Their limitation is that listed IP addresses can be rotated; in the DMCN, shared reputation feeds serve an analogous purpose but list persistent cryptographic identities rather than easily-changed IP addresses.

---

**Relay Node**
A node in the DMCN mesh network that participates in routing and (optionally) storing messages. Relay nodes verify sender signatures against the identity registry before forwarding messages, which is the mechanism by which unsigned or fraudulently signed messages are rejected at the network level. Because messages are encrypted to the recipient's public key, relay nodes cannot read message content.

---

**secp256k1**
An elliptic curve used in the Bitcoin and Ethereum cryptographic systems, and widely deployed in cryptocurrency wallets and Web3 infrastructure. It is one of the candidate curve families for the DMCN identity layer, alongside Curve25519. Both provide strong security properties; the choice between them involves trade-offs in performance, ecosystem compatibility, and implementation maturity.

---

**Secure Enclave**
A hardware-isolated execution environment within a device's processor, designed to protect sensitive operations (such as private key generation and signing) from the rest of the operating system. Even if the device's OS is compromised, private keys stored in the Secure Enclave cannot be extracted. Apple devices use a dedicated Secure Enclave chip; Android devices use a Trusted Execution Environment (TEE); Windows machines use a Trusted Platform Module (TPM). The DMCN stores private keys in hardware-backed secure enclaves wherever supported.

---

**Shared Reputation Feed**
An opt-in, community-maintained registry of cryptographic public keys that have been reported for abuse (spam, phishing, harassment, malware distribution). Analogous to DNS blocklists in legacy email, but with a critical structural advantage: because entries are bound to persistent cryptographic identities rather than easily-rotated IP addresses or domains, a listed identity's reputation loss is permanent. Users subscribe to one or more feeds through their DMCN client.

---

**SMTP (Simple Mail Transfer Protocol)**
The protocol that has underpinned internet email since its standardisation in RFC 821 in 1982. SMTP was designed for a small network of trusted academic nodes and provides no mechanism for cryptographic sender verification. This architectural absence is the root cause of spam, phishing, and email identity fraud. The DMCN is designed to eventually replace SMTP as the primary internet messaging substrate, while maintaining interoperability during a transition period.

---

**SPF (Sender Policy Framework)**
An email authentication standard that allows domain owners to publish, via DNS, a list of IP addresses authorised to send email on their behalf. Receiving servers can check whether an incoming message originated from an authorised IP. SPF addresses domain-level sending authorisation but does not provide per-message cryptographic signing, and cannot prevent spoofing from compromised authorised servers.

---

**Sybil Attack**
An attack on a trust-based network in which a malicious actor creates a large number of fake identities to gain disproportionate influence or to overwhelm defences. In the DMCN context, the primary Sybil attack scenario involves creating many registered identities to conduct spam campaigns before they are blocklisted. The DMCN mitigates this through account creation friction and permanent reputation consequences, but full Sybil resistance is an open design challenge.

---

**TPM (Trusted Platform Module)**
A hardware security chip found in most modern Windows PCs and many enterprise laptops, designed to store cryptographic keys and perform security-sensitive operations in isolation from the main processor. The TPM is the Windows equivalent of the Secure Enclave on Apple devices. The DMCN uses the TPM for private key storage on Windows machines that support it.

---

**Web of Trust**
A decentralised model for establishing cryptographic identity assurance, in which users cross-sign each other's public keys to attest that a given key genuinely belongs to the claimed identity. Rather than relying on a central certificate authority, trust is built through a network of individual attestations. PGP pioneered the web-of-trust model in 1991; the DMCN implements a modernised version with improved UX, using QR code exchange for in-person key verification.

---

**Web3**
A broad term for a vision of a decentralised internet built on blockchain infrastructure, typically involving cryptocurrency, NFTs, and decentralised applications. Several competing email and messaging projects (including Dmail Network) are oriented toward the Web3 ecosystem. The DMCN is explicitly designed as a general-purpose email replacement accessible to mainstream users, and does not require engagement with Web3 or cryptocurrency infrastructure.

---

**Allowlist**
In the DMCN trust model, the allowlist is the user's registry of confirmed trusted contacts. Unlike a simple address book, allowlist entries in the DMCN are cryptographically bound to specific public keys and carry a record of how trust was established (in-person verification, fingerprint check, network vouching, etc.). Messages from allowlisted contacts are delivered directly to the primary inbox without passing through the pending queue.



---


---

## References

References are grouped thematically. Internet standards and RFCs are listed first, followed by academic and research literature, then industry and project sources. All URLs were verified as of March 2026.

---

### Internet Standards and RFCs

**[RFC821]**
Postel, J. (1982). *Simple Mail Transfer Protocol*. RFC 821. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc821

**[RFC5321]**
Klensin, J. (2008). *Simple Mail Transfer Protocol*. RFC 5321. Internet Engineering Task Force. (Supersedes RFC 2821 and RFC 821.)
https://www.rfc-editor.org/rfc/rfc5321

**[RFC7208]**
Kitterman, S. (2014). *Sender Policy Framework (SPF) for Authorizing Use of Domains in Email, Version 1*. RFC 7208. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc7208

**[RFC6376]**
Crocker, D., Hansen, T., & Kucherawy, M. (2011). *DomainKeys Identified Mail (DKIM) Signatures*. RFC 6376. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc6376

**[RFC7489]**
Kucherawy, M., & Zwicky, E. (2015). *Domain-based Message Authentication, Reporting, and Conformance (DMARC)*. RFC 7489. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc7489

**[RFC7671]**
Dukhovni, V., & Hardaker, W. (2015). *The DANE Protocol: Updates and Operational Guidance*. RFC 7671. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc7671

**[RFC4033]**
Arends, R., Austein, R., Larson, M., Massey, D., & Rose, S. (2005). *DNS Security Introduction and Requirements* (DNSSEC). RFC 4033. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc4033

**[RFC4880]**
Callas, J., Donnerhacke, L., Finney, H., Shaw, D., & Thayer, R. (2007). *OpenPGP Message Format*. RFC 4880. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc4880

**[RFC5652]**
Housley, R. (2009). *Cryptographic Message Syntax (CMS)* — the basis for S/MIME. RFC 5652. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc5652

**[RFC7519]**
Jones, M., Bradley, J., & Sakimura, N. (2015). *JSON Web Token (JWT)*. RFC 7519. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc7519

**[RFC9180]**
Barnes, R., Bhargavan, K., Lipp, B., & Wood, C. (2022). *Hybrid Public Key Encryption (HPKE)*. RFC 9180. Internet Engineering Task Force. (Standardises the KEM + KDF + AEAD pattern used in the DMCN encrypted envelope design.)
https://www.rfc-editor.org/rfc/rfc9180

---

### Cryptographic Foundations

**[Bernstein2006]**
Bernstein, D. J. (2006). *Curve25519: New Diffie-Hellman Speed Records*. In: Yung, M., Dodis, Y., Kiayias, A., Malkin, T. (eds), *Public Key Cryptography — PKC 2006*, Lecture Notes in Computer Science, vol 3958. Springer, Berlin, Heidelberg.
https://cr.yp.to/ecdh/curve25519-20060209.pdf

**[Certicom2000]**
Certicom Research. (2000). *SEC 2: Recommended Elliptic Curve Domain Parameters*. Standards for Efficient Cryptography Group. (Defines secp256k1.)
https://www.secg.org/sec2-v2.pdf

**[Bernstein2012]**
Bernstein, D. J., & Lange, T. (2012). *SafeCurves: Choosing Safe Curves for Elliptic-Curve Cryptography.*
https://safecurves.cr.yp.to

**[Dingledine2004]**
Dingledine, R., Mathewson, N., & Syverson, P. (2004). *Tor: The Second-Generation Onion Router*. Proceedings of the 13th USENIX Security Symposium.
https://svn.torproject.org/svn/projects/design-paper/tor-design.pdf

**[Maymounkov2002]**
Maymounkov, P., & Mazières, D. (2002). *Kademlia: A Peer-to-Peer Information System Based on the XOR Metric*. Proceedings of the 1st International Workshop on Peer-to-Peer Systems (IPTPS 2002).
https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf

**[Marlinspike2016]**
Marlinspike, M., & Perrin, T. (2016). *The Double Ratchet Algorithm*. Signal Foundation.
https://signal.org/docs/specifications/doubleratchet/

**[Shamir1979]**
Shamir, A. (1979). *How to Share a Secret*. Communications of the ACM, 22(11), 612–613. (Foundational paper for threshold secret sharing, the basis for the DMCN social recovery mechanism.)

---

### Decentralised Identity and Key Management

**[W3C-DID]**
Sporny, M., Guy, A., Sabadello, M., & Reed, D. (2022). *Decentralized Identifiers (DIDs) v1.0*. W3C Recommendation. World Wide Web Consortium.
https://www.w3.org/TR/did-core/

**[W3C-VC]**
Sporny, M., Longley, D., & Chadwick, D. (2022). *Verifiable Credentials Data Model v1.1*. W3C Recommendation. World Wide Web Consortium.
https://www.w3.org/TR/vc-data-model/

**[FIDO-Passkeys]**
FIDO Alliance. (2023). *Passkeys: FIDO Authentication Overview*. FIDO Alliance White Paper.
https://fidoalliance.org/passkeys/

**[Apple-Passkeys]**
Apple Inc. (2023). *About the Secure Enclave*. Apple Platform Security Guide.
https://support.apple.com/guide/security/secure-enclave-sec59b0b31ff/web

**[Zimmermann1995]**
Zimmermann, P. R. (1995). *The Official PGP User's Guide*. MIT Press, Cambridge, MA. (PGP's original design rationale and web-of-trust model.)

---

### Spam, Phishing, and Email Abuse

**[FBI-IC3-2023]**
Federal Bureau of Investigation, Internet Crime Complaint Center. (2023). *Internet Crime Report 2023*. U.S. Department of Justice. (Annual report documenting BEC as a leading cybercrime category by financial loss.)
https://www.ic3.gov/Media/PDF/AnnualReport/2023_IC3Report.pdf

**[Symantec-ISTR]**
Broadcom / Symantec. (2023). *Internet Security Threat Report*. (Annual industry report on spam volumes, phishing trends, and email-based malware delivery.)
https://symantec-enterprise-blogs.security.com/blogs/threat-intelligence/istr-2023

**[Enisa2023]**
European Union Agency for Cybersecurity (ENISA). (2023). *ENISA Threat Landscape 2023*. (European threat intelligence report covering phishing and BEC trends.)
https://www.enisa.europa.eu/publications/enisa-threat-landscape-2023

**[Levchenko2011]**
Levchenko, K., Pitsillidis, A., Chachra, N., Enright, B., Félegyházi, M., Grier, C., Halvorson, T., Kanich, C., Kreibich, C., Liu, H., McCoy, D., Weaver, N., Paxson, V., Voelker, G. M., & Savage, S. (2011). *Click Trajectories: End-to-End Analysis of the Spam Value Chain*. Proceedings of the IEEE Symposium on Security and Privacy. (Landmark study of spam economics; foundational for the DMCN's economic argument against spam.)

---

### Decentralised Messaging Systems and Precedents

**[Matrix-Spec]**
Hodgson, M., et al. (2024). *The Matrix Specification v1.9*. Matrix.org Foundation.
https://spec.matrix.org/latest/

**[Matrix-Bridges]**
Matrix.org Foundation. (2023). *Bridging*. Matrix Specification Documentation. (Documents the federated bridge architecture that the DMCN bridge design draws upon.)
https://matrix.org/docs/matrix-concepts/bridges/

**[Signal-Protocol]**
Marlinspike, M. (2023). *Signal Protocol Documentation*. Signal Foundation.
https://signal.org/docs/

**[BitTorrent-DHT]**
BitTorrent Inc. (2008). *BitTorrent Enhancement Proposal 5: DHT Protocol*.
http://www.bittorrent.org/beps/bep_0005.html

**[Keybase-Proofs]**
Keybase Inc. (2019). *Keybase Identity Proofs*. Keybase Documentation. (Documents the cryptographic proof-of-ownership model that the DMCN address portability feature is modelled on.)
https://book.keybase.io/docs/server/proofs

---

### Competing and Related Projects

**[Dmail-Whitepaper]**
Dmail Network. (2023). *Dmail Network: Decentralized AI Communication Platform*. Dmail Network Technical Documentation.
https://dmail.ai/whitepaper

**[ProtonMail-Security]**
Proton AG. (2023). *Proton Mail Security Features*. Proton Mail Documentation.
https://proton.me/support/proton-mail-encryption-explained

**[Tutanota-Security]**
Tutanota GmbH. (2023). *Tutanota Security Model*. Tutanota Documentation.
https://tuta.com/support/articles/blog/security

---

### Regulatory and Compliance Context

**[GDPR]**
European Parliament and Council. (2016). *Regulation (EU) 2016/679 on the protection of natural persons with regard to the processing of personal data (General Data Protection Regulation)*. Official Journal of the European Union.
https://eur-lex.europa.eu/eli/reg/2016/679/oj

**[HIPAA]**
U.S. Department of Health and Human Services. (1996, as amended). *Health Insurance Portability and Accountability Act of 1996 (HIPAA)*. (Relevant to the DMCN compliance discussion for regulated industries.)
https://www.hhs.gov/hipaa/index.html

**[NIS2]**
European Parliament and Council. (2022). *Directive (EU) 2022/2555 on measures for a high common level of cybersecurity across the Union (NIS2 Directive)*. Official Journal of the European Union.
https://eur-lex.europa.eu/eli/dir/2022/2555

---

### Domain Name and PKI Infrastructure

**[DNS-RFC1034]**
Mockapetris, P. (1987). *Domain Names — Concepts and Facilities*. RFC 1034. Internet Engineering Task Force.
https://www.rfc-editor.org/rfc/rfc1034

**[CAB-Forum]**
CA/Browser Forum. (2023). *Baseline Requirements for the Issuance and Management of Publicly-Trusted Certificates, Version 2.0.0*. (Defines the PKI hierarchy that S/MIME relies upon, and which the DMCN is designed to avoid.)
https://cabforum.org/baseline-requirements-documents/



---
