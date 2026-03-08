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
