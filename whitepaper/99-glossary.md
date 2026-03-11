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

