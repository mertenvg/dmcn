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
