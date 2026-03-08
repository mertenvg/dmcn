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
- [3. The Competitive Landscape](#3-the-competitive-landscape)
- [4. Proposed Architecture: Decentralized Mesh Communication Network](#4-proposed-architecture)
- [5. Cryptographic Identity and Key Management](#5-cryptographic-identity-and-key-management)
- [6. Spam Elimination at the Protocol Level](#6-spam-elimination-at-the-protocol-level)
- [7. User Experience: Hiding Complexity Without Sacrificing Security](#7-user-experience)
- [8. Transition Strategy: Coexistence with Legacy Email](#8-transition-strategy)
- [9. The SMTP-DMCN Bridge Architecture](#9-the-smtp-dmcn-bridge-architecture)
- [10. Bringing Existing Email Addresses to the DMCN](#10-bringing-existing-email-addresses-to-the-dmcn)
- [11. Trust Management: Whitelists, Greylists, and Blacklists](#11-trust-management)
- [12. Threat Model](#12-threat-model)
- [13. Open Challenges and Research Questions](#13-open-challenges-and-research-questions)
- [14. Conclusion](#14-conclusion)
- [Glossary](#glossary)
- [15. Privacy Analysis](#15-privacy-analysis)
- [References](#references)

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

- A layered trust management framework covering whitelists, greylists, shared reputation feeds, and the cryptographic blacklisting model that makes identity reputation permanent and non-transferable

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

This whitepaper is presented as Version 0.1 of an ongoing investigation.
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
> *Spam and email fraud are not filtering problems --- they are
> identity problems. A mesh network where every node is
> cryptographically identified and every message is cryptographically
> signed eliminates the conditions under which spam is economically
> viable, rather than attempting to detect and discard it after the
> fact. This is the only class of solution that addresses the root
> cause.*

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
> --- they are the normal operating conditions of the current system.*


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


## 3. The Competitive Landscape


### 3.1 Dmail Network


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


### 3.2 ProtonMail and Tutanota


ProtonMail and Tutanota represent the encrypted email approach within
the existing SMTP infrastructure. Both offer end-to-end encrypted email
between users on their own platform, with varying levels of encryption
when communicating with external addresses. They have achieved
meaningful consumer adoption, particularly among privacy-conscious
users. However, both are centralized services and neither addresses the
spam problem at the protocol level.


### 3.3 Signal and Matrix


Signal and Matrix (Element) demonstrate that decentralized,
cryptographically secure messaging is technically viable at scale.
Signal has achieved significant mainstream adoption while providing
state-of-the-art end-to-end encryption. Matrix provides a federated,
decentralized protocol with open-source infrastructure. Neither is
designed as an email replacement, but both provide important technical
and UX precedents that a DMCN design should draw upon.


### 3.4 Summary Comparison


  ---------------- ------------------- ------------- ------------ -------------- --------------
  **Solution**     **Decentralized**   **Spam-Free   **No Crypto  **Mainstream   **Email
                                       by Design**   Required**   UX**           Compatible**

  PGP / S/MIME     Partial             No            No           No             Yes

  Dmail Network    Yes                 Partial       No           No             Partial

  ProtonMail       No                  No            Yes          Yes            Yes

  Signal / Matrix  Yes                 N/A           Yes          Yes            No

  Proposed DMCN    Yes                 Yes           Yes          Yes            Yes
  ---------------- ------------------- ------------- ------------ -------------- --------------


## 4. Proposed Architecture: Decentralized Mesh Communication Network


### 4.1 Design Principles


The Decentralized Mesh Communication Network (DMCN) is designed around a
set of foundational principles derived from the failure analysis of
prior approaches. These principles are architectural constraints that
shape every design decision.

- Identity is cryptographic and self-sovereign. Every participant in the network has a unique identity derived from a public/private key pair. Identity is not assigned by any central authority and cannot be revoked by any third party.

- Sender verification is mandatory and protocol-enforced. A message without a valid cryptographic signature from a registered identity cannot enter the network. Verification is not optional, not opt-in, and not a filter applied after the fact — it is a gate at the point of transmission.

- The network is peer-to-peer with no central routing authority. Messages are relayed through a distributed mesh of nodes. No single entity controls routing, storage, or delivery.

- Complexity is hidden from end users. The cryptographic machinery that makes the network trustworthy operates entirely below the user experience layer. Users interact with human-readable identities and familiar communication patterns.

- Legacy email interoperability is a first-class requirement. The network must be capable of sending to and receiving from legacy SMTP addresses during a transition period.


### 4.2 Network Architecture


The DMCN consists of three logical layers, each with distinct
responsibilities:


#### 4.2.1 Identity Layer


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


#### 4.2.2 Transport Layer


The Transport Layer is responsible for routing messages through the mesh
network. Messages are addressed to the recipient's public key,
encrypted with that public key, signed with the sender's private key,
and relayed through a network of nodes using an onion-routing-inspired
protocol that provides metadata privacy in addition to content privacy.

Relay nodes can verify message signatures against the identity layer,
which is the mechanism by which spam is rejected at the network level
--- a node that receives a message signed by an identity not registered
in the identity layer drops the message without delivery.


#### 4.2.3 Storage and Delivery Layer


Unlike real-time messaging systems, email is inherently asynchronous.
The Storage and Delivery Layer provides distributed, encrypted message
storage that holds messages until the recipient's client connects to
retrieve them. Messages are stored encrypted with the recipient's
public key; relay nodes providing storage cannot read message content.


### 4.3 Message Lifecycle


A message in the DMCN follows this lifecycle:

- The sender's client composes a message and signs it with the sender's private key.

- The client encrypts the signed message with the recipient's public key, retrieved from the Identity Layer.

- The encrypted, signed message is submitted to the transport layer with the recipient's public key as the address.

- Relay nodes verify the sender's signature against the Identity Layer. Messages with invalid or absent signatures are dropped.

- The message is routed through the mesh to relay nodes serving the recipient's address, where it is held in encrypted storage.

- When the recipient's client connects, it retrieves and decrypts messages using the recipient's private key.

- The recipient's client verifies the sender's signature, confirming the message genuinely originated from the stated sender.


## 5. Cryptographic Identity and Key Management


### 5.1 Key Generation and Storage


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


### 5.2 Public Key Distribution and Discovery


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


### 5.3 The Key Management UX Problem — And Its Solution


The failure of PGP, despite its technical soundness, is primarily
attributable to the burden it placed on users to manage cryptographic
keys. The DMCN takes a fundamentally different approach, drawing on the
model established by passkeys and mobile device security:

- Key generation is automatic and invisible. When a user creates an account, keys are generated on their device without any user-facing step involving cryptographic concepts.

- Private keys are never shown to the user. Unlike cryptocurrency wallets, users are not presented with seed phrases or private key strings.

- Key backup is automatic and encrypted. Private keys are backed up using the device's existing encrypted cloud backup infrastructure (iCloud Keychain, Google Password Manager, or an equivalent DMCN-native encrypted key backup service).

- Multi-device access is handled through secure key synchronization — a flow identical to the device migration process used by Signal.

- Account recovery is possible through a social recovery mechanism. Users designate trusted contacts who hold encrypted shards of a recovery key. A threshold of contacts (for example, 3 of 5) must participate to restore access.


### 5.4 Identity Verification and the Trust Graph


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
> attestations that users build over time --- exactly as trust is
> established in the physical world.*


## 6. Spam Elimination at the Protocol Level


### 6.1 Why Cryptographic Identity Eliminates Spam


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


### 6.2 Consent-Based Communication


The DMCN introduces consent-based message acceptance as a first-class
protocol feature. By default, a user's inbox accepts messages only from
identities that meet one of the following criteria:

- The sender is in the recipient's contact list.

- The sender shares a verified organizational identity with the recipient.

- The sender's identity has been vouched for by a trusted contact in the recipient's web of trust.

- The recipient has explicitly opted in to receiving messages from unknown senders (for public figures or customer-facing businesses).

Messages from senders that do not meet any of these criteria are placed
in a pending queue where the recipient can review them. These messages
still bear valid cryptographic signatures — the sender's identity is
known — allowing the recipient to make an informed decision.


### 6.3 Economic Disincentives for Spam


Beyond protocol-level barriers, the DMCN can implement optional economic
mechanisms that further disincentivize spam. A micro-payment channel
system allows senders to unknown recipients to attach a small,
refundable deposit to their message. If the recipient accepts, the
deposit is returned. If rejected, the deposit is forfeited. This imposes
no cost on messages between known contacts, and only trivial cost on
legitimate outreach — but makes mass spam campaigns economically
prohibitive.


## 7. User Experience: Hiding Complexity Without Sacrificing Security


### 7.1 The Fundamental Principle


The history of secure communication tools is, in large part, a history
of UX failures. The DMCN's user experience layer is designed around a
single principle: the security model must be invisible to the user in
normal operation. Users should experience the DMCN as a familiar
messaging application — with inboxes, contacts, compose windows, and
threads — and should never encounter the words 'key', 'signature',
'certificate', or 'encryption' in the normal flow of using the
product.


### 7.2 Familiar Addressing


Users are addressed using a format that mirrors traditional email: a
local identifier and a domain, separated by the @ symbol. Internally,
this address resolves to a public key — but from the user's
perspective, it is simply their address, just as a phone number is
simply a phone number without any awareness of the SS7 routing protocol
underneath.


### 7.3 Onboarding Flow


Account creation in the DMCN is designed to be comparable in friction to
creating a Gmail account. The user provides a chosen identifier,
authenticates with biometric or PIN, and the application generates their
key pair silently in the background. The entire process takes under two
minutes. There is no seed phrase, no key download, no certificate
request, and no cryptographic terminology.


### 7.4 Contact Discovery


Finding another user on the DMCN requires only their address. The
application resolves the address to a public key through the distributed
identity registry, and the contact appears in the user's contact list.
All messages to that contact are automatically encrypted and signed. The
user does not need to take any additional steps to enable security ---
it is on by default and cannot be turned off.


### 7.5 Trust Indicators


The application surfaces trust information in intuitive, non-technical
ways. Verified organizational identities display a checkmark alongside
the sender's domain. Messages from unknown senders appear in a separate
pending section. A simple trust indicator shows whether a contact's
identity has been verified by mutual connections in the user's network.


## 8. Transition Strategy: Coexistence with Legacy Email


### 8.1 The Migration Problem


No communication platform has ever achieved mainstream adoption by
requiring users to abandon their existing communication infrastructure.
The transition strategy for the DMCN is built on the principle of
graceful degradation — the system provides maximum value and security
to users communicating with each other on the DMCN, while maintaining
the ability to communicate with legacy email users at reduced security
levels during the transition period.


### 8.2 DMCN-to-DMCN Communication


When both sender and recipient are on the DMCN, messages are fully
encrypted, cryptographically signed, peer-to-peer routed, and spam-free
by protocol. This is the target state of the system and the experience
that should be promoted as the default.


### 8.3 DMCN-to-Email Communication


When a DMCN user sends a message to a legacy email address, the message
passes through a gateway node that translates it to SMTP format for
delivery. The message can include a footer inviting the recipient to
join the DMCN. Security properties are reduced in this path — message
content must be decrypted at the gateway for SMTP delivery — but
sender identity remains verifiable at the gateway level.


### 8.4 Email-to-DMCN Communication


Receiving a message from a legacy email sender requires a verified
gateway address system, where legacy emails pass through a gateway that
performs basic spam filtering and sender verification at the SMTP level
before delivering to the DMCN inbox. Users may also maintain a connected
legacy email address displayed in a separate, clearly labeled section of
their DMCN client.


### 8.5 Progressive Migration Incentives


The transition strategy includes mechanisms that actively incentivize
migration to native DMCN communication: visible trust indicators that
distinguish DMCN-verified messages from legacy email; organizational
compliance features requiring DMCN for sensitive internal
correspondence; and developer APIs allowing third-party applications to
integrate DMCN identity as a communication primitive.


## 9. The SMTP-DMCN Bridge Architecture


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


### 9.1 Architectural Role and Placement


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


### 9.2 Outbound Path: DMCN to SMTP


When a DMCN user addresses a message to a legacy email address, the
message follows a modified delivery path through a bridge node. This is
the simpler of the two directions, and the security properties are
well-defined.


#### 9.2.1 Message Flow


- The sender's DMCN client composes and signs the message with the sender's private key, as in a standard DMCN message.

- The client detects that the recipient address resolves to a legacy email address (no DMCN public key found in the identity registry) and routes the message to a bridge node rather than the standard transport layer.

- The bridge node receives the encrypted, signed DMCN message, verifies the sender's signature against the identity registry, and decrypts the message content using the bridge's own private key (the message is re-encrypted to the bridge's key rather than a non-existent recipient key).

- The bridge constructs a standard SMTP message from the decrypted content, applying DKIM signing using the bridge operator's domain key. The From address is set to a bridge-scoped representation of the sender's DMCN address (e.g., username=dmcn.net@bridge.dmcn.net), preserving sender identity in a form that legacy email clients can display.

- The bridge delivers the SMTP message to the recipient's mail server using standard MX lookup and SMTP relay.

- A standardized footer is appended to the message body, displaying the sender's verified DMCN identity and an invitation link for the recipient to join the DMCN and receive future messages natively.


#### 9.2.2 Trust Properties


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


### 9.3 Inbound Path: SMTP to DMCN


Receiving messages from legacy email senders is the more complex
direction, because the sender has no DMCN identity and no cryptographic
signing capability. The bridge must make a trust determination about an
unverified sender and communicate that determination clearly to the
recipient.


#### 9.3.1 Addressing


DMCN users who wish to receive legacy email obtain a bridge address ---
a standard email address managed by the bridge operator (e.g.,
username@bridge.dmcn.net). This address is registered as an MX record
pointing to the bridge node's SMTP listener. Users can publish this
address on websites, business cards, and legacy systems as their contact
point for people who have not yet adopted DMCN.


#### 9.3.2 Message Flow


- The bridge's SMTP listener receives an inbound message addressed to the user's bridge address.

- The bridge performs a full suite of legacy authentication checks: SPF validation, DKIM signature verification, DMARC policy evaluation, and reverse DNS lookup on the sending mail server.

- The bridge queries distributed reputation databases (analogous to existing RBL/DNSBL services) for the sending IP address and domain.

- Messages that fail hard authentication checks (invalid DKIM, SPF failure with DMARC reject policy) are dropped with a delivery failure response to the sending server.

- Messages that pass or partially pass authentication are classified into trust tiers: Verified Legacy Sender (valid DKIM + positive reputation), Unverified Legacy Sender (no DKIM or neutral reputation), and Suspicious Legacy Sender (reputation flags present).

- The bridge wraps the classified message in a DMCN envelope, signed by the bridge's own private key as an attestation of the classification outcome. The DMCN envelope includes the full authentication result metadata, the trust tier assignment, and the original SMTP headers.

- The wrapped message is encrypted to the recipient's DMCN public key and delivered through the standard DMCN transport layer to the recipient's inbox.


#### 9.3.3 Recipient Experience


Inbound legacy messages appear in a clearly distinguished section of the
recipient's DMCN client, visually separated from native DMCN messages.
Each message displays its trust tier — Verified Legacy, Unverified
Legacy, or Suspicious — with a plain-language explanation of what the
classification means. The recipient can promote a legacy sender to their
contact list (which triggers the DMCN client to send the sender an
invitation to join the network natively) or block the sender
permanently.


### 9.4 Bridge Node Security Model


Bridge nodes are high-value infrastructure components and require a
rigorous security model. Several specific threats must be addressed:

- Bridge impersonation — a malicious actor operating a fraudulent bridge that misrepresents message authentication results. Mitigated by requiring bridge operators to register their identity in the DMCN identity registry, publish their security practices, and submit to periodic independent audits.

- Content interception — a bridge operator reading or modifying message content in transit on the outbound path. Mitigated by end-to-end message signing (recipients can verify the sender's signature even after bridge translation), audit logging, and regulatory accountability for commercial operators.

- SMTP relay abuse — spammers attempting to use the inbound bridge path to inject spam into DMCN inboxes. Mitigated by the authentication classification system and by rate limiting on the bridge's SMTP listener per sending domain and IP.

- Bridge node compromise — an attacker gaining control of a bridge node. Mitigated by the decentralized bridge model (no single bridge handles all traffic), key rotation protocols, and the ability for users to revoke trust in a specific bridge operator.


### 9.5 Federated Bridge Architecture


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


### 9.6 Precedents and Comparable Systems


The bridge architecture is a well-established pattern in communication
infrastructure. Several analogous systems demonstrate that protocol
translation at scale is an engineering challenge with proven solutions:

- Matrix bridges (Synapse) — the Matrix protocol has operated production bridges to Slack, Discord, WhatsApp, Telegram, SMS, and IRC simultaneously, with hundreds of thousands of active users. The Matrix bridge architecture provides a detailed precedent for federated, bidirectional protocol translation.

- Email-to-SMS gateways — carriers have operated SMTP-to-SMS translation services for decades, handling billions of messages. The protocol mismatch between email and SMS is in many respects more severe than the mismatch between SMTP and DMCN.

- SIP-PSTN gateways in VoIP — Voice over IP systems routinely bridge between SIP-native networks and the legacy public switched telephone network, including cryptographic signaling translation. The architecture for handling identity and trust across this boundary is directly analogous to the SMTP-DMCN bridge problem.


> **Engineering Assessment**
> *The SMTP-DMCN bridge is a significant but tractable engineering
> project. It does not require novel cryptographic research --- it
> requires careful implementation of known patterns applied to a new
> protocol boundary. A working prototype bridge is a realistic
> deliverable for an initial proof-of-concept phase.*


## 10. Bringing Existing Email Addresses to the DMCN


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


### 10.1 The Principle of Address Portability


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


### 10.2 Verification Mechanisms


The strength of the link between an email address and a DMCN identity
depends on the level of control the user has over the address and its
underlying domain. The DMCN supports three verification tiers, each with
distinct trust properties:


#### 10.2.1 Provider-Hosted Address Verification (e.g., Gmail, Outlook)


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


#### 10.2.2 Custom Domain Address Verification


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


#### 10.2.3 DANE-Style Cryptographic Domain Binding


For domains that have enabled DNSSEC — the cryptographic extension to
DNS that provides tamper-evident records — a third verification tier
is available that provides the highest level of assurance. In this
model, the domain owner publishes the DMCN public key directly in a
DNSSEC-signed record, creating a chain of cryptographic trust from the
DNS root through the domain to the individual identity. This approach is
analogous to DANE (DNS-based Authentication of Named Entities), which is
already used in some contexts to bind TLS certificates to domain names
without relying on certificate authorities.


### 10.3 Trust Implications of Each Tier


  ------------------- ---------------- ---------------- -----------------
  **Verification      **Proof of       **Resistant to   **Recommended
  Tier**              Control**        Provider         For**
                                       Action**         

  Provider-Hosted     Inbox access     No — provider  Individual users
  (Gmail, Outlook)    only             can suspend      during transition

  Custom Domain DNS   DNS record       Yes — domain   Professionals,
  Verification        control          owner controls   small businesses
                                       DNS              

  DNSSEC / DANE       Cryptographic    Yes — highest  Enterprises,
  Cryptographic       DNS chain        assurance        regulated
  Binding                                               industries
  ------------------- ---------------- ---------------- -----------------


### 10.4 The Honest Limitation: Ownership vs. Control


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
> supports both --- but only one provides true long-term identity
> independence.*


### 10.5 Address Portability and the Spam Problem


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


### 10.6 Precedents


The address portability model draws on several well-established
precedents in both identity verification and email infrastructure:

- Keybase — demonstrated the viability of linking a cryptographic identity to multiple existing identities (email, Twitter, GitHub, domain) through a system of cryptographic proofs. The DMCN's verification model is conceptually similar but integrated at the protocol level rather than as a third-party overlay.

- Google Workspace and Microsoft 365 custom domain onboarding — both services use DNS TXT record verification to prove domain ownership before allowing custom domain email. The DMCN's Domain Verification tier follows this exact pattern, which is already familiar to IT administrators globally.

- DKIM public key DNS records — the practice of publishing cryptographic keys in DNS records is already standard email infrastructure. The DMCN's DANE-style binding extends this established pattern.

- Number portability in mobile telephony — the telecommunications industry solved an analogous problem when it allowed consumers to bring their phone numbers between carriers. The lesson from that transition is directly applicable: portability dramatically lowers switching costs and accelerates adoption of superior infrastructure.


## 11. Trust Management: Whitelists, Greylists, and Blacklists


Cryptographic identity verification is the foundation of the DMCN's
trust model — it answers the question of whether a message genuinely
came from a claimed sender. But verification alone does not answer a
second, equally important question: whether the user actually wants to
hear from that sender. Trust management is the user-facing system that
sits on top of cryptographic verification and allows each participant to
define, on their own terms, who they trust, who they are uncertain
about, and who they actively reject.

The DMCN's trust management system operates at three tiers ---
whitelist, greylist, and blacklist — each with distinct delivery
semantics, key storage implications, and sharing properties. Together
they form a layered defence that is more powerful than anything
available in legacy email, precisely because the identities being
managed are cryptographic and persistent rather than superficial and
easily spoofed.


### 11.1 The Whitelist: Confirmed Trusted Senders


The whitelist is the user's registry of confirmed trusted contacts. It
is not merely an address book — it is a cryptographically anchored
record that binds a human identity to a specific public key, with a
record of how and when that binding was established. Every entry in the
whitelist carries a trust provenance — the mechanism by which the user
confirmed the contact's identity — which is surfaced in the client UI
to help users understand the strength of each trust relationship.


#### 11.1.1 Trust Establishment Mechanisms


The DMCN supports multiple mechanisms for adding a contact to the
whitelist, ranked here in descending order of trust strength:

- Direct key exchange — the user and contact are physically present and exchange public keys via a QR code scan in the DMCN mobile application. This establishes an in-person verified binding with the highest possible assurance. The resulting whitelist entry is marked Verified In-Person.

- Fingerprint verification — the user retrieves a contact's public key from the identity registry and then verifies the key fingerprint through an independent channel (a phone call, a video call, a previously trusted communication method). The user confirms that the fingerprint read aloud by the contact matches the one in the registry. Marked Fingerprint Verified.

- Web-of-trust promotion — the contact is already whitelisted by two or more of the user's existing Verified contacts. The user can choose to extend trust on the basis of their network's endorsement, with a clear indication of which mutual contacts vouch for the new addition. Marked Network Vouched.

- Organisational verification — the contact holds a DMCN identity attested by an organisation the user has already verified (e.g., a colleague whose identity is attested by a shared employer domain). Marked Organisationally Verified.

- First-message confirmation — the user receives a message from an unknown sender and actively chooses to approve and whitelist them. This is the weakest trust mechanism — the user has reviewed the message and chosen to accept the sender, but has not independently verified the key. Marked User Approved.

Trust provenance is preserved indefinitely in the whitelist record and
is visible to the user at any time. A contact marked Verified In-Person
carries a fundamentally different assurance than one marked User
Approved, and the client communicates this distinction without requiring
the user to understand the underlying cryptography.


#### 11.1.2 Key Binding and Update Handling


Because whitelist entries are bound to specific public keys rather than
addresses alone, the DMCN client must handle the case where a contact's
key changes — for example, when they migrate to a new device, perform
a key rotation, or recover their account through the social recovery
mechanism.

When a contact's public key changes, the client presents an explicit
notification to the user: the previous key is no longer active, a new
key has been published, and the user must re-verify the contact before
the whitelist binding is updated. Automatic silent key updates are not
permitted for whitelist entries — the user must consciously re-confirm
the relationship. This prevents a class of attack in which an adversary
replaces a contact's key in the identity registry and silently
intercepts subsequent communication.


> **Key Change Alert**
> *When a whitelisted contact\'s public key changes, the DMCN client
> suspends delivery from that identity and alerts the user. No message
> is delivered under an unconfirmed new key until the user explicitly
> re-verifies. This is a deliberate friction point --- it is the
> correct response to a high-assurance security event.*


#### 11.1.3 Whitelist Portability and Backup


The whitelist is an asset of significant personal value — it
represents years of accumulated trust relationships. It is therefore
backed up as part of the user's encrypted key material and can be
exported in a standardised, encrypted format for migration between
clients or for archival. The export format includes not only the public
keys and addresses but also the full trust provenance record for each
entry, so that the history of how trust was established is preserved
across migrations.


### 11.2 The Greylist: Unknown but Unblocked Senders


The greylist occupies the space between explicit trust and explicit
rejection. It is the default destination for messages from DMCN-verified
senders who are not yet on the user's whitelist — verified in the
cryptographic sense, meaning their signature is valid and their identity
is registered, but not yet confirmed as trusted by the user.


#### 11.2.1 Greylist Delivery Semantics


Messages arriving from greylist senders are held in a pending queue,
visually distinct from the primary inbox. The client displays the
sender's verified DMCN identity, the trust provenance of the sender in
the network (whether any of the user's contacts have vouched for them,
and if so how many), and a summary of the message sufficient to make an
informed accept-or-reject decision — without requiring the user to
fully open and read the message first.

From the greylist queue the user has four options for each pending
message: Accept and whitelist the sender (promoting all future messages
to the primary inbox), Accept this message only (delivering the message
without whitelisting the sender), Reject and ignore (discarding the
message without any notification to the sender), or Reject and blacklist
(discarding the message and adding the sender to the blacklist to
prevent future delivery attempts).


#### 11.2.2 Greylist Auto-Resolution Rules


To reduce the burden of manual greylist management, the client supports
configurable auto-resolution rules that can automatically promote or
demote senders based on network signals:

- Auto-promote if vouched by N or more whitelist contacts — configurable threshold, default of 3.

- Auto-promote if sender holds a verified organisational identity matching a domain the user has previously whitelisted.

- Auto-promote if the sender's identity has a reputation score above a configurable threshold in the user's chosen shared reputation feed.

- Auto-reject if the sender's identity appears on any blacklist feed the user has subscribed to.

These rules run at delivery time, before the message reaches the pending
queue, and are fully configurable. Users who want complete manual
control can disable all auto-resolution rules. Users who want a more
automated experience can enable conservative defaults that handle the
common cases without requiring intervention.


### 11.3 The Blacklist: Blocking Known Bad Actors


The blacklist is the user's registry of explicitly rejected senders.
Unlike a legacy email block — which can be trivially circumvented by
creating a new address — a DMCN blacklist entry is bound to a
cryptographic identity. A blacklisted sender cannot reach the user by
creating a new address, because their underlying key pair is what is
blocked, not the surface-level address string. This is a fundamentally
stronger guarantee than any blocking mechanism available in legacy
email.


#### 11.3.1 Personal Blacklist


The personal blacklist is private to the user and is never shared
externally. Adding a sender to the personal blacklist causes the DMCN
relay nodes handling the user's incoming messages to silently drop any
message signed by that identity before it reaches the user's device ---
the sender receives no delivery failure notification and no indication
that they have been blocked. This is consistent with the behaviour of
email blocking in major clients today and prevents the blocked sender
from using delivery failures as a signal to probe for workarounds.

Personal blacklist entries include the blocked identity's public key,
the address at which they were known, the date of blocking, and an
optional private note from the user recording their reason for the
block. This note is stored encrypted with the user's private key and is
never transmitted.


#### 11.3.2 Shared Reputation Feeds


Beyond the personal blacklist, the DMCN supports an opt-in shared
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


#### 11.3.3 Reputation Feed Architecture


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


#### 11.3.4 Reporting and Feed Contribution


Any user can submit a report against a sender's public key to a feed
they are subscribed to. The report is signed with the reporting user's
private key, providing cryptographic accountability for the report ---
false or malicious reports can be traced back to the reporter's
identity. This accountability mechanism is important: it discourages
coordinated campaigns to falsely blacklist legitimate identities,
because the reporters themselves are identifiable.

Feed operators implement their own thresholds and policies for when a
reported key is listed. A conservative operator might require reports
from twenty or more independent, verified identities before listing a
key. A more aggressive operator might list on fewer reports with greater
weighting for reports from highly trusted identities. Users select feed
operators whose threshold policies match their tolerance for false
positives versus false negatives.


#### 11.3.5 The Persistence Advantage


The most significant property of a cryptographic blacklist relative to
its legacy equivalents deserves explicit emphasis. When a DMCN identity
is reported and listed across multiple feeds, that listing is
effectively permanent for that key pair. The spammer's investment in
building a sending reputation — any messages that passed through
greylists, any contacts who user-approved them, any network vouching
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


### 11.4 Trust Tier Interaction Summary


  ------------- -------------------- --------------- -------------- ----------------
  **Tier**      **Sender Type**      **Delivery      **Key Bound?** **Shareable?**
                                     Destination**                  

  Whitelist     Verified trusted     Primary inbox,  Yes — with   Exportable
                contact              immediate       provenance     (private)
                                     delivery                       

  Greylist      Verified but unknown Pending queue,  Yes —        No
                sender               user review     identity       
                                                     displayed      

  Personal      Explicitly rejected  Silently        Yes — key    No (private)
  Blacklist     sender               dropped at      blocked        
                                     relay                          

  Shared        Community-reported   Dropped per     Yes —        Yes ---
  Reputation    bad actor            feed policy     persistent     community opt-in
  Feed                                               listing        
  ------------- -------------------- --------------- -------------- ----------------

---

## 12. Threat Model


This section provides a structured analysis of the threat landscape
facing the Decentralized Mesh Communication Network. For each identified
threat category, we document the nature of the attack, how it manifests
in the existing SMTP email ecosystem, how the DMCN architecture changes
the threat surface, and an honest assessment of whether the DMCN
improves, maintains, or introduces new risk relative to the status quo.


> **Methodology**
> *This threat model follows a structured adversarial analysis
> approach: for each threat, we identify the adversary\'s goal, the
> attack vector available under SMTP, the attack vector (if any)
> available under DMCN, and the net change in risk. Threats are grouped
> by adversary type: mass senders, targeted attackers, infrastructure
> attackers, and state-level actors.*


### 12.1 Threat Category 1: Spam and Bulk Unsolicited Messaging


#### 12.1.1 Nature of the Threat


Spam is the dominant form of abuse in the current email ecosystem. A
spam operator's goal is to deliver promotional, fraudulent, or
malicious messages to as many recipients as possible at the lowest
possible cost per delivery. The economics are straightforward: even at a
conversion rate of a fraction of a percent, the near-zero marginal cost
of sending email makes spam campaigns profitable.


#### 12.1.2 How SMTP Enables This Threat


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


#### 12.1.3 How DMCN Changes the Threat Surface


The DMCN eliminates the conditions that make spam economically viable at
the protocol level. Every sender must possess a registered cryptographic
identity, and every message must bear a valid signature from that
identity. Relay nodes reject unsigned or unverifiably signed messages
before they enter the network. This imposes a non-trivial, though low,
cost on account creation — and critically, each identity's reputation
is permanent and non-transferable.

A spam operator who wishes to send at scale must create a large number
of registered identities. Each identity that is reported and blacklisted
is permanently lost — there is no equivalent of rotating to a new IP
address. The mathematical relationship between spam volume and identity
cost shifts the economics of spam from profitable to uneconomical at
scale.


#### 12.1.4 Residual Risk and Honest Limitations


The DMCN does not make spam creation infinitely expensive — it makes
it non-zero in cost and permanently cumulative in consequence. A
determined, well-resourced spam operation could potentially automate the
account creation process (a Sybil attack), creating large numbers of
identities before they are reported. Section 12.5 addresses Sybil
resistance specifically. The consent-based inbox model (Section 6.2 of
the main whitepaper) provides a secondary layer: even a registered
identity cannot reach a user's primary inbox without meeting one of the
whitelisting criteria.


> **Net Assessment**
> *Spam: Significantly mitigated. The protocol-level economics of spam
> are fundamentally changed. Residual risk exists through Sybil attacks
> and requires robust account creation friction to fully close.*


### 12.2 Threat Category 2: Phishing and Identity Spoofing


#### 12.2.1 Nature of the Threat


Phishing attacks exploit the inability of email recipients to reliably
verify sender identity. An attacker impersonates a trusted entity — a
bank, an employer, a government agency — to induce the recipient to
reveal credentials, transfer funds, or install malware. Business Email
Compromise (BEC) is a sophisticated variant in which attackers
impersonate executives or financial officers within an organisation to
authorise fraudulent wire transfers. BEC alone causes billions of
dollars in losses annually.


#### 12.2.2 How SMTP Enables This Threat


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


#### 12.2.3 How DMCN Changes the Threat Surface


In the DMCN, every message carries a cryptographic signature tied to the
sender's registered key pair. It is mathematically impossible to forge
a message that appears to originate from a registered identity without
possessing that identity's private key. The recipient's client
verifies this signature automatically, and surfaces a clear,
non-technical trust indicator based on the sender's position in the
recipient's trust graph.

A phishing attempt from an unregistered identity cannot enter the
network. A phishing attempt from a registered identity that is not in
the recipient's whitelist lands in the greylist pending queue, visibly
distinguished from trusted messages. The only viable phishing vector in
the DMCN is a fully registered identity that the recipient has
explicitly trusted — which requires the attacker to have established a
prior trust relationship.


#### 12.2.4 Account Compromise and the Key Binding Problem


The DMCN does not eliminate the threat of account compromise — it
changes its character. In SMTP, a compromised email account allows an
attacker to send messages that are indistinguishable from legitimate
messages from that sender. In the DMCN, a compromised account requires
the attacker to have stolen the private key itself, not merely the login
credentials. Private keys stored in hardware-backed secure enclaves (as
specified in Section 5.1) cannot be extracted even if the device's
operating system is compromised.

This represents a meaningful improvement over SMTP account compromise,
but it introduces a new concern: if a private key is stolen (e.g., from
a device without hardware security support), the attacker gains the full
trust relationships of that identity with no visible indicator to
contacts. The whitelist key-change notification system (Section 11.1.2)
partially mitigates this: if the attacker uses a new device, contacts
will be alerted that the key has changed and prompted to re-verify.


> **Net Assessment**
> *Phishing: Substantially mitigated for protocol-level spoofing.
> Account compromise risk shifts from credential theft to key theft ---
> a meaningfully higher bar, but not zero. Hardware key storage is
> essential for this property to hold.*


### 12.3 Threat Category 3: Infrastructure Attacks


#### 12.3.1 Denial of Service Against the Network


Any distributed network is a potential target for denial-of-service
attacks. In the DMCN, the primary infrastructure targets are relay
nodes, the distributed identity registry, and bridge nodes. The goal of
such an attack is to prevent message delivery, disrupt identity lookups,
or degrade the network to the point of unusability.


#### 12.3.2 Comparison with SMTP


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


#### 12.3.3 New Infrastructure Risks Introduced by DMCN


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
users. The federated bridge architecture (Section 9.5) distributes this
risk, but organisations using a single bridge provider remain exposed to
single-point-of-failure risk.


> **Net Assessment**
> *Infrastructure DoS: Comparable to SMTP for distributed attacks;
> improved for centralised attack scenarios due to peer-to-peer
> architecture. The identity registry is a novel attack surface
> requiring careful design. Bridge nodes reintroduce some
> centralisation risk during the transition period.*


### 12.4 Threat Category 4: Relay Node Misbehaviour


#### 12.4.1 Nature of the Threat


Unlike centralised email providers, DMCN relay nodes can be operated by
any party. This raises the question of what happens when a relay node
operator acts maliciously or negligently. Potential misbehaviours
include selectively dropping messages, logging message metadata (even
though content is encrypted), injecting false routing information, or
colluding with other nodes to deanonymise communication patterns.


#### 12.4.2 Comparison with SMTP


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


#### 12.4.3 Metadata Privacy and the Onion Routing Layer


The proposed onion-routing-inspired transport protocol (Section 4.2.2)
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


### 12.5 Threat Category 5: Sybil Attacks


#### 12.5.1 Nature of the Threat


A Sybil attack occurs when a malicious actor creates a large number of
fake identities to subvert a trust-based system. In the context of the
DMCN, the primary Sybil attack scenarios are: creating large numbers of
identities to conduct spam campaigns before they are blacklisted;
creating fake identities to inflate web-of-trust vouching for a
malicious identity; and creating fake identities to manipulate shared
reputation feeds.


#### 12.5.2 Comparison with SMTP


SMTP is essentially infinitely susceptible to Sybil attacks — there is
no meaningful identity system to attack, and the cost of registering a
new sending domain is a few dollars. The DMCN's identity model is
inherently more resistant because it requires account creation friction,
and because blacklisted identities cannot be recovered. However, the
DMCN is not immune, and Sybil resistance is one of the most significant
open design challenges.


#### 12.5.3 Proposed Mitigations


Several mechanisms can be combined to raise the cost of Sybil attacks to
uneconomical levels:

- Account creation friction: requiring a verified phone number, email address, or proof-of-work computation during account creation raises the per-identity cost above zero

- Rate limiting on new identity behaviour: newly created identities are subject to stricter greylist treatment and lower throughput limits until they have established a reputation history

- Web-of-trust bootstrapping requirements: the consent-based inbox model means that a new identity must earn its way into recipients' whitelists; a Sybil farm of identities with no trust relationships has no delivery capability

- Reputation feed correlation: feed operators can flag clusters of identities with similar creation timing, device fingerprints, or behaviour patterns as likely Sybil farms


> **Net Assessment**
> *Sybil attacks: Improved relative to SMTP, but not fully solved. The
> per-identity cost is non-zero and cumulative in consequence, raising
> the economics of Sybil attacks. Full resistance requires careful
> design of account creation friction --- a balance between security
> and accessibility that requires user research and iteration.*


### 12.6 Threat Category 6: State-Level Surveillance and Censorship


#### 12.6.1 Nature of the Threat


Nation-state actors represent the most sophisticated and well-resourced
adversaries in the threat landscape. Their objectives may include mass
surveillance of communications content, targeted surveillance of
specific individuals, censorship of communications between specific
parties, or disruption of communication infrastructure for geopolitical
purposes.


#### 12.6.2 How SMTP Enables State Surveillance


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


#### 12.6.3 How DMCN Changes the Threat Surface


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


#### 12.6.4 Censorship and Network Disruption


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


### 12.7 Threat Category 7: Key Compromise and Recovery Attacks


#### 12.7.1 Nature of the Threat


The security of the entire DMCN identity model rests on the secrecy of
each user's private key. If a private key is compromised, the attacker
gains the ability to read all future messages sent to that identity, to
send messages that appear to come from that identity, and to modify or
revoke that identity's trust relationships. Key compromise is therefore
the most serious category of attack specific to a cryptographic identity
system.


#### 12.7.2 Comparison with SMTP


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


#### 12.7.3 The Social Recovery Attack Surface


The social recovery mechanism (Section 5.3) — in which trusted
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


#### 12.7.4 Key Revocation and Forward Secrecy


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


### 12.8 Threat Category 8: Bridge Node Attacks


#### 12.8.1 Nature of the Threat


The SMTP-DMCN bridge architecture (Section 9) is a necessary component
of any viable transition strategy, but it reintroduces several trust and
security challenges that the native DMCN architecture otherwise
eliminates. Bridge nodes represent the interface between the trustless
SMTP world and the cryptographically verified DMCN, and they must make
trust determinations about SMTP senders that have no direct equivalent
in the native protocol.


#### 12.8.2 Bridge-Specific Attack Vectors


The following attacks are specific to the bridge architecture and have
no equivalent in the native DMCN:

- Content interception on the outbound path: bridge nodes must decrypt outbound messages to re-encode them as SMTP. A malicious or compromised bridge operator gains access to message content in transit. This is disclosed in Section 9.2.2 and is an unavoidable consequence of protocol translation.

- False trust classification: a malicious bridge could misrepresent the trust tier of an inbound SMTP message — for example, classifying a spam message as 'Verified Legacy Sender' to bypass the recipient's filters. The bridge's classification is signed with the bridge's own DMCN key, creating accountability, but this only helps if users verify which bridge they are trusting.

- Legacy spam injection: spammers may target the bridge's SMTP listener as a pathway into DMCN inboxes, attempting to exploit weaknesses in the bridge's SMTP authentication checks to inject messages that would be rejected if sent natively.

- Bridge impersonation: an attacker could operate a bridge that presents itself as trustworthy in the identity registry but maliciously handles messages. Mitigated by requiring bridge operators to publish their security practices and undergo periodic audits.


#### 12.8.3 Bridge Risk as a Transitional Concern


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


### 12.9 Threat Model Summary


The table below summarises each threat category, the current severity in
SMTP, the treatment under DMCN, and the net outcome for each:

  -----------------------------------------------------------------------
  **Threat          **SMTP Severity** **DMCN            **Net Outcome**
  Category**                          Treatment**       
  ----------------- ----------------- ----------------- -----------------
  Spam / Bulk       Critical —      Protocol-level    **Significantly
  Messaging         protocol endemic  identity cost     Reduced**
                                      eliminates        
                                      economic          
                                      viability         

  Phishing /        Critical —      Cryptographic     **Significantly
  Spoofing          trivially         signing makes     Reduced**
                    executed          spoofing          
                                      mathematically    
                                      infeasible        

  Infrastructure    High —          Distributed       **Partially
  DoS               centralised       architecture      Mitigated**
                    targets           reduces           
                                      single-point risk 

  Relay Node        High —          End-to-end        **Partially
  Misbehaviour      plaintext in      encryption limits Mitigated**
                    transit           relay visibility  
                                      to metadata       

  Sybil Attacks     N/A — no        Non-zero identity **Partially
                    identity system   cost; permanent   Mitigated**
                                      reputation loss   

  State             Critical —      End-to-end        **Significantly
  Surveillance      provider access   encryption; no    Reduced**
                                      centralised       
                                      interception      
                                      point             

  Key Compromise    High —          Hardware keys     **Partially
                    passwords weak    raise bar;        Mitigated**
                                      recovery          
                                      introduces new    
                                      surface           

  Bridge Attacks    N/A —           Bounded to legacy **Partially
                    DMCN-specific     traffic;          Mitigated**
                                      diminishes with   
                                      adoption          
  -----------------------------------------------------------------------


> **Overall Assessment**
> *The DMCN architecture represents a meaningful and substantial
> improvement over SMTP across every threat category where comparison
> is possible. The threats that remain partially mitigated rather than
> eliminated --- Sybil resistance, relay metadata leakage, key recovery
> attacks, and bridge node risks --- are well-understood engineering
> challenges with known mitigation approaches, rather than fundamental
> architectural flaws. None of these residual risks represents a
> regression relative to the current SMTP-based email ecosystem.*

---

## 13. Open Challenges and Research Questions


This whitepaper represents a preliminary investigation into the design
space of a Decentralized Mesh Communication Network. Several significant
challenges remain open and will require further research, prototyping,
and community input.


### 13.1 Scale and Performance


The distributed identity registry and peer-to-peer routing architecture
must be demonstrated to perform adequately at the scale of global email
--- billions of users, hundreds of billions of messages per day. The
performance characteristics of the proposed architecture under realistic
load conditions must be validated through simulation and prototype
deployment.


### 13.2 Key Recovery Without Central Authority


The social recovery model proposed in Section 5 is promising, but its UX
and security properties require careful design and user research. The
threshold for recovery must balance security against the practical
reality that trusted contacts may be unavailable, may themselves lose
access to their accounts, or may be compromised. Alternative recovery
mechanisms should be investigated and compared.


### 13.3 Regulatory and Legal Compliance


End-to-end encrypted communication creates compliance challenges for
regulated industries — financial services, healthcare, law — that
are required to maintain records of communications and provide them to
regulators on demand. The DMCN architecture must provide mechanisms that
allow regulated entities to meet their compliance obligations without
compromising the security properties of the system for other users.


### 13.4 Governance


A truly decentralized network requires a governance model that allows
the protocol to evolve — to address security vulnerabilities,
incorporate technical improvements, and respond to regulatory changes
--- without any central authority having unilateral control. The
governance model for the DMCN is a critical design question with
significant implications for the network's long-term resilience and
trustworthiness.


### 13.5 Sybil Resistance


While the DMCN's identity model prevents spam from unregistered
senders, it must also resist Sybil attacks — scenarios in which a
malicious actor creates a large number of registered identities to
overwhelm spam defenses. The account creation process must impose
sufficient cost or friction to make large-scale Sybil attacks
uneconomical, without imposing unacceptable burden on legitimate users.

---

## 14. Conclusion


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
in Section 13 are real and significant. The competitive landscape
documented in Section 3 demonstrates that the market is beginning to
recognize the problem space, even if existing solutions have not yet
solved it effectively.

The opportunity exists for a solution that combines the cryptographic
soundness of PGP, the decentralized architecture of peer-to-peer
networks, the user experience standards of modern consumer applications,
and a credible migration path from the existing email ecosystem. The
DMCN is a proposal toward that solution.


> **Next Steps**
> *This whitepaper represents Version 0.1 of an ongoing investigation.
> Subsequent versions will incorporate technical prototyping results,
> user research findings, and engagement with the broader cryptographic
> and communications research community. Feedback, critique, and
> collaboration are actively solicited.*

---

*End of Document*

---

## Glossary

Terms are listed alphabetically. Where a term has a common abbreviation used in this document, the abbreviation is shown in parentheses.

---

**Blacklist**
A user-maintained or community-maintained registry of cryptographic identities that are explicitly blocked from delivering messages. In the DMCN, blacklist entries are bound to public keys rather than surface addresses, making them impossible to circumvent by simply creating a new address string. See also: *Greylist*, *Whitelist*, *Shared Reputation Feed*.

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

**Greylist**
In the DMCN trust model, the greylist is the holding area for messages from senders who are cryptographically verified (their identity is registered and their message signature is valid) but not yet explicitly trusted by the recipient. Greylisted messages appear in a pending queue for user review rather than the primary inbox. This differs from the legacy email concept of greylisting, which involves temporary rejection of messages from unknown sending servers.

---

**Key Pair**
A matched pair of cryptographic keys — a public key and a private key — generated together using a mathematical relationship such that data encrypted with the public key can only be decrypted with the corresponding private key, and data signed with the private key can be verified with the public key. In the DMCN, every user identity is represented by a key pair; the public key is published in the identity registry, while the private key never leaves the user's device.

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
An attack on a trust-based network in which a malicious actor creates a large number of fake identities to gain disproportionate influence or to overwhelm defences. In the DMCN context, the primary Sybil attack scenario involves creating many registered identities to conduct spam campaigns before they are blacklisted. The DMCN mitigates this through account creation friction and permanent reputation consequences, but full Sybil resistance is an open design challenge.

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

**Whitelist**
In the DMCN trust model, the whitelist is the user's registry of confirmed trusted contacts. Unlike a simple address book, whitelist entries in the DMCN are cryptographically bound to specific public keys and carry a record of how trust was established (in-person verification, fingerprint check, network vouching, etc.). Messages from whitelisted contacts are delivered directly to the primary inbox without passing through the greylist queue.



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

## 15. Privacy Analysis

This section addresses a question distinct from the threat model in Section 12: not whether the DMCN can be attacked, but what the system *inherently reveals* during normal, correct operation. A communication network can be cryptographically secure against active attackers while still exposing significant information about its users through the ordinary mechanics of message routing, identity discovery, and protocol operation.

The privacy analysis is structured around four areas: metadata exposure at the network layer, the identity registry as a surveillance surface, bridge node privacy, and regulatory compliance in a decentralised architecture. Each area is assessed against a baseline of what the current SMTP email ecosystem reveals, so that the comparison is grounded rather than abstract.

> **Scope**
> *This analysis addresses privacy in the technical sense — what information is exposed to which parties as a structural consequence of the protocol — rather than the policy sense of what operators choose to do with data. Operator conduct is a governance and regulatory matter addressed in Section 13.3 and Section 15.4.*

---

### 15.1 Baseline: What SMTP Reveals

Before assessing the DMCN, it is worth being precise about the privacy properties of the system it proposes to replace. SMTP email, as deployed by major providers, exposes the following to varying parties:

**To the email provider (Gmail, Outlook, etc.):** Full message content, subject lines, sender and recipient addresses, timestamps, device metadata, IP addresses, and — through scanning for features like Smart Reply and spam classification — inferred behavioural patterns and social graphs. Major providers operate under privacy policies that permit substantial use of this data for advertising and product improvement, subject to jurisdiction-specific regulatory constraints.

**To relay infrastructure in transit:** Historically, SMTP transmitted message content in plaintext. Opportunistic TLS between mail transfer agents, now widely deployed, encrypts content in transit between servers — but each server in the relay chain can read content, and TLS is not universally enforced. Message headers, including sender, recipient, routing path, and timestamps, are structurally visible to all relay infrastructure.

**To passive network observers:** On links where TLS is not enforced, full message content is visible. Even with TLS, connection metadata — which servers are communicating, at what times, with what data volumes — is observable at the network layer.

**To recipients:** Full message content, the sender's email address, and whatever metadata the sending client and relay chain have appended to message headers.

This is the baseline against which the DMCN's privacy properties should be measured. The bar is not high.

---

### 15.2 Metadata Exposure at the Network Layer

End-to-end encryption protects message *content* from relay nodes — no node in the DMCN transport layer can read what Alice is saying to Bob. What encryption does not protect is *metadata*: the fact that Alice is communicating with Bob, the frequency of their exchanges, the timing of messages, and approximate message sizes. This metadata can be as revealing as content in many threat models.

#### 15.2.1 What Relay Nodes Can Observe

A DMCN relay node handling a message in transit can observe the following:

- The sender's public key (as the message's cryptographic identifier)
- The recipient's public key (as the routing address)
- The approximate size of the encrypted payload
- The timestamp of receipt and forwarding
- The IP address of the upstream node that delivered the message, and the IP address of the downstream node to which it is forwarded

A relay node cannot observe the message content, subject, or any human-readable metadata. It also cannot — in a correctly implemented onion routing scheme — observe both the originating sender's IP address and the ultimate recipient's IP address simultaneously. Each relay node sees only the previous and next hop in the delivery chain.

This is a material improvement over SMTP, where relay nodes can read full message content and headers. However, it is not equivalent to anonymity. A relay node that handles a high volume of traffic for a small number of users can build a detailed picture of communication patterns between pseudonymous identities (public keys) even without reading content.

#### 15.2.2 What a Global Passive Observer Can Infer

The most sophisticated metadata threat is a global passive adversary — an entity capable of observing a significant fraction of all network traffic simultaneously. This is the same threat that onion routing networks such as Tor are known to be vulnerable to through traffic correlation attacks.

By observing that a message-sized packet left Alice's IP address at time T, and that an equivalently-sized packet arrived at Bob's relay node shortly after T, a global passive observer can probabilistically correlate the two events and infer that Alice sent a message to Bob — even without reading either packet's content.

The DMCN's onion routing layer partially mitigates this by introducing multiple hops and variable timing, increasing the difficulty of correlation. It does not eliminate it. For users whose threat model includes nation-state-level traffic analysis — journalists communicating with sources, activists in authoritarian jurisdictions, legal counsel in sensitive matters — the DMCN should be understood as providing strong content privacy with meaningful but imperfect metadata privacy. Users with these threat models should be directed to Tor-over-DMCN configurations or equivalent overlay networks for the transport layer.

For the vast majority of users whose threat model does not include a global passive adversary, the DMCN's metadata privacy properties represent a substantial improvement over SMTP.

#### 15.2.3 Message Size as a Side Channel

Encrypted message sizes are observable by relay nodes and passive network observers even when content is not. In some contexts, message size is itself informative — a 50KB encrypted message is more likely to contain a document attachment than a brief reply. The DMCN transport layer should implement padding to normalise message sizes into a small number of size classes, reducing the inferential value of size observation. This is a standard technique in privacy-preserving transport protocols and is a recommended design requirement, though its implementation detail is deferred to the protocol specification.

---

### 15.3 The Identity Registry as a Surveillance Surface

The DMCN's distributed identity registry is public by architectural necessity. For the system to function — for any sender to be able to look up a recipient's public key and send them an encrypted message — the registry must be queryable by anyone. This openness is a deliberate design choice, but it creates a surveillance surface that requires explicit attention.

#### 15.3.1 Account Existence Confirmation

Anyone with access to the identity registry can query it to determine whether a given email address has been registered as a DMCN identity. This means:

- An advertiser, stalker, data broker, or intelligence agency can compile lists of DMCN users by querying the registry against known email addresses.
- The registry reveals not just public keys but the fact of account existence, which may itself be sensitive in some contexts.
- Bulk enumeration of the registry — attempting to discover all registered identities — is a privacy concern if not rate-limited.

**Mitigation:** The registry should implement rate limiting on lookups per source IP or per authenticated identity, and should not support bulk enumeration queries. Lookups should return only the specific queried identity, not adjacencies or related entries. The registry design should also consider whether to support unlisted identities — accounts that are reachable by existing contacts but do not appear in registry searches initiated by unknown parties.

#### 15.3.2 Social Graph Inference from the Registry

Because registry entries include the identity's public key and any cross-signatures from the web of trust, a determined observer who accumulates registry data over time can begin to map social relationships: if Alice's key is cross-signed by Bob's and Carol's keys, it is inferable that Alice, Bob, and Carol have a trust relationship. At scale, the web-of-trust attestation data constitutes a partial social graph that is structurally visible without reading any message content.

**Mitigation:** Web-of-trust attestations should be opt-in for public visibility. Users should be able to maintain private attestations — stored locally or exchanged out-of-band — that inform their own trust decisions without being published to the global registry. The registry specification should distinguish between public attestations (visible to all) and private attestations (held locally or shared only with specific contacts).

#### 15.3.3 Timing and Correlation via Registry Lookups

When Alice's client performs a registry lookup for Bob's public key, that lookup is itself observable as a network event. A network observer monitoring Alice's traffic can infer that Alice is about to send a message to Bob — before the message is even sent — simply by observing the registry query.

**Mitigation:** The client should implement a registry prefetching strategy, maintaining a local cache of public keys for recent and likely contacts and refreshing it on a schedule rather than on demand. This decouples the timing of registry lookups from the timing of message composition, reducing the inferential value of lookup timing.

---

### 15.4 Bridge Node Privacy

The SMTP-DMCN bridge architecture, addressed in Section 9 from a security perspective, has distinct privacy implications that require separate treatment.

#### 15.4.1 Outbound Path: What the Bridge Operator Sees

When a DMCN user sends a message to a legacy email address, the message must be decrypted at the bridge node to be re-encoded as SMTP. This is an unavoidable consequence of protocol translation, disclosed in Section 9.2.2. The privacy implication is explicit: the bridge operator has technical access to:

- The full content of every outbound message sent to legacy email recipients
- The sender's DMCN identity
- The recipient's legacy email address
- Timestamps and message sizes

This is structurally equivalent to the trust placed in a conventional email service provider such as Gmail or Outlook, which has identical access to message content. The difference is that users choosing the DMCN are typically doing so with an expectation of enhanced privacy — and the bridge's necessary content access may conflict with that expectation if not clearly disclosed at onboarding.

**Disclosure requirement:** The DMCN client must present a clear, non-technical disclosure at the point where a user first sends a message to a legacy email recipient, explaining that the bridge operator can read the content of messages sent to non-DMCN addresses. This disclosure should be persistent — not a one-time consent flow that users will click through without reading — and the privacy policy of the chosen bridge operator should be linked and surfaced in the client UI.

**Mitigation through operator choice:** Because the bridge architecture is federated (Section 9.5), users can choose bridge operators with strong privacy commitments, including operators that commit to zero message logging and are subject to independent audit. Organisations with strong confidentiality requirements can operate their own bridge nodes, eliminating third-party access entirely.

#### 15.4.2 Inbound Path: Legacy Sender Metadata

When a legacy email sender sends a message to a DMCN user's bridge address, the bridge operator observes the full SMTP headers of the inbound message: sender address, sending server IP, timestamps, and routing path. This metadata is used to perform the authentication classification described in Section 9.3.2 and is necessarily retained for that purpose.

The DMCN specification should define minimum and maximum retention periods for bridge-held metadata, consistent with applicable data protection law, and should require bridge operators to publish their metadata retention policies.

#### 15.4.3 Bridge Operator as Data Controller

Under the EU General Data Protection Regulation and equivalent frameworks, any entity that determines the purposes and means of processing personal data is a data controller with obligations to data subjects. A bridge operator processing message content and metadata on behalf of DMCN users is a data controller for that processing.

This has practical implications: bridge operators must have a lawful basis for processing, must respond to data subject access requests, must implement appropriate technical and organisational security measures, and must be located in or have adequate data transfer arrangements with the jurisdictions in which their users reside. The DMCN Bridge Operator Protocol (BOP) should incorporate these requirements as conditions of registry participation, so that users can trust that any registered bridge operator meets a minimum compliance baseline.

---

### 15.5 Regulatory Privacy Compliance in a Decentralised Architecture

Decentralisation creates genuine tension with data protection frameworks that were designed around the assumption of an identifiable data controller. The DMCN's privacy architecture must honestly engage with this tension rather than treating decentralisation as a compliance shield.

#### 15.5.1 The Data Controller Problem

In the DMCN's native peer-to-peer layer, there is no central operator. Messages are stored encrypted on relay nodes, routed through the mesh, and held until the recipient retrieves them. No single entity controls the processing of any given user's messages. This makes it difficult to identify a data controller in the GDPR sense — and without a data controller, data subjects' rights (access, rectification, erasure, portability) become difficult to exercise.

**Practical positions:**

- For the core protocol layer, the user themselves may be considered the data controller for their own encrypted data, since only they hold the decryption key. This is analogous to the position taken by some self-hosted encrypted services.
- For relay nodes storing encrypted messages, the relay node operator may be considered a data processor acting on behalf of the user-controller, with a data processing agreement required.
- For bridge nodes, as discussed in Section 15.4.3, the operator is a data controller in their own right for the content they can access.
- For the identity registry, the distributed architecture means there is no single controller; each node operator is a processor of the subset of registry data they hold.

These positions are not fully settled in law and will require engagement with data protection authorities in relevant jurisdictions as the DMCN matures. The governance framework (Section 13.4) should include a dedicated working group on regulatory compliance.

#### 15.5.2 The Right to Erasure

GDPR Article 17 grants data subjects the right to erasure of their personal data. In the DMCN, this creates a specific challenge: encrypted messages stored on relay nodes cannot be deleted by the user on demand, because the user has no direct administrative relationship with relay node operators.

**Partial mitigation:** Messages stored on relay nodes are encrypted with the recipient's public key. If the recipient deletes their private key — or if the relay node's retention policy expires the message — the encrypted data becomes permanently inaccessible even if the bytes persist on disk. This achieves functional erasure (the data is unrecoverable) even without literal deletion of the stored bytes.

The DMCN specification should define a maximum message retention period for relay nodes, after which stored messages are deleted regardless of whether they have been retrieved. A default of 30 days with user-configurable extension is a reasonable starting point, consistent with practices in existing encrypted messaging systems.

#### 15.5.3 Data Portability

GDPR Article 20 grants data subjects the right to receive their personal data in a structured, machine-readable format and to transmit it to another controller. In practice, for a messaging system, this means the user's message history, contact list, and trust relationships.

The DMCN client should implement a full data export function that produces a portable, encrypted archive of the user's message history, whitelist, greylist, blacklist, and trust attestations in a documented, open format. This export serves both the regulatory compliance function and the practical function of enabling migration between DMCN client applications without loss of data.

---

### 15.6 Privacy Properties Summary

The table below summarises the DMCN's privacy properties across key dimensions, compared to the SMTP baseline.

| Privacy Dimension | SMTP (Gmail/Outlook) | DMCN Native | DMCN via Bridge |
|---|---|---|---|
| Message content visible to provider | Yes — always | No — E2EE | Yes — to bridge operator |
| Message content visible to relay infrastructure | Yes — historically; TLS in transit only | No — E2EE throughout | Partial — decrypted at bridge |
| Sender/recipient visible to relay nodes | Yes — full headers | Pseudonymous public keys only | Yes — to bridge operator |
| Metadata visible to passive network observer | Yes — sender/recipient, timing, size | Timing and size only (onion routing limits more) | Timing and size |
| Social graph inferable from infrastructure | Yes — from provider data | Partially — from registry attestations | Partially |
| Account existence discoverable | Yes — MX lookup | Yes — registry query | Yes |
| Data subject rights (GDPR etc.) | Provider is data controller; rights exercisable | Distributed; complex controller picture | Bridge operator is data controller |
| Message retention | Provider-controlled; typically indefinite | Relay node retention policy; finite | Bridge operator retention policy |

---

### 15.7 Design Recommendations

The privacy analysis above yields the following concrete design recommendations for the DMCN specification, in priority order:

**Message size normalisation** — implement payload padding in the transport layer to reduce the inferential value of size observation. This is a low-cost, high-value privacy improvement.

**Registry rate limiting and unlisted identity support** — prevent bulk enumeration of the identity registry and allow users to opt out of public discoverability while remaining reachable by existing contacts.

**Private web-of-trust attestations** — allow trust attestations to be held locally rather than published to the global registry, preserving the utility of the web of trust without exposing social graph data.

**Registry lookup prefetching** — decouple registry lookups from message composition timing to reduce the inferential value of lookup timing to network observers.

**Bridge operator disclosure** — require persistent, prominent disclosure in the DMCN client when messages will be processed by a bridge operator, including a link to the operator's privacy policy.

**Relay node retention limits** — specify a maximum message retention period in the protocol, ensuring that unread encrypted messages do not persist indefinitely on relay infrastructure.

**Data export function** — implement a full, portable, encrypted data export in the DMCN client to satisfy data portability obligations and enable client migration.

**Regulatory working group** — establish a dedicated working group within the DMCN governance structure to engage with data protection authorities and develop jurisdiction-specific compliance guidance.

