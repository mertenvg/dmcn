# DMCN — Questions & Answers

Common questions about the DMCN architecture, transition strategy, and bridge design.

---

## Bridge and Legacy Email

### How does SMTP route email to a bridge node?

SMTP routes email based on the MX (Mail Exchange) DNS records of the recipient's domain. To receive legacy email, the bridge operator must control a domain and configure its MX records to point to the bridge's SMTP server.

For example, if the bridge domain is `bridge.example.com`:

1. The operator sets `MX 10 smtp.bridge.example.com` in DNS
2. A legacy sender composes email to `alice@bridge.example.com`
3. The sender's mail server queries `MX bridge.example.com`, gets the bridge host
4. The mail server connects to the bridge's SMTP port and delivers the message
5. The bridge verifies sender authentication, classifies trust, maps `alice@bridge.example.com` to `alice@dmcn.example.com`, encrypts to Alice's public key, and stores on the relay

From the perspective of legacy SMTP infrastructure, the bridge is just another mail server. The translation to DMCN is invisible to the sending side.

### Can a user with a Gmail (or other provider) address receive DMCN messages through the bridge?

Not via their Gmail address. Google controls the MX records for `gmail.com`, so all mail to `alice@gmail.com` is delivered to Google's servers. No bridge can intercept that traffic.

The bridge handles two directions:

- **Legacy to DMCN:** A legacy sender sends to `alice@bridge.example.com` (a domain the bridge operator controls). The bridge wraps the message into an encrypted DMCN envelope.
- **DMCN to Legacy:** A DMCN user sends to `bob@gmail.com`. The bridge decrypts, delivers via SMTP to Gmail's MX servers. The Gmail user receives it as a normal email.

What does not work is intercepting inbound email for domains you do not control.

### How does a Gmail user participate in DMCN?

A user on Gmail (or any provider-controlled domain) has several options:

1. **Register a DMCN identity.** Create a new DMCN address (e.g. `alice@dmcn.example.com`) and use a DMCN client. Messages between DMCN users stay entirely within DMCN — no bridge, no SMTP, full end-to-end encryption. The Gmail address remains separate for legacy communication.

2. **Get a bridge domain address.** Obtain an address on a bridge-controlled domain (e.g. `alice@bridge.example.com`) and share it with legacy senders who should reach them via DMCN. This is effectively a second email-facing address that routes through the bridge into the DMCN network.

3. **Migrate to a domain they control.** If the user owns a domain (e.g. `alice@example.com`), they can point its MX records at a bridge and transition fully, receiving all legacy email through the DMCN bridge.

The typical transition path:

1. Register a DMCN identity — communicate natively with other DMCN users
2. Optionally get a bridge domain address — let legacy senders reach them via DMCN
3. Over time, share the DMCN address more, the provider address less
4. The provider address becomes a fallback, not the primary

---

## Message Routing

### If both users have DMCN identities, does the message touch SMTP at all?

No. If both sender and recipient have identities registered in the DHT, the message stays entirely within DMCN:

1. Sender looks up the recipient's address in the DHT
2. Gets the recipient's public keys (Ed25519 for verification, X25519 for encryption)
3. Signs the message with their own Ed25519 key
4. Encrypts directly to the recipient's X25519 key
5. Stores the encrypted envelope on the recipient's relay node

No SMTP, no bridge, no DNS. The message is end-to-end encrypted from sender to recipient. The bridge only enters the picture when one side is on legacy SMTP.

### How does a DMCN user send to a legacy email address?

The sender encrypts the message to the bridge's public key (the bridge has a registered DMCN identity with `BridgeCapability` set). The bridge:

1. Decrypts the envelope using its own private key
2. Verifies the sender's signature
3. Extracts the legacy recipient address from the message
4. Delivers via outbound SMTP to the recipient's mail server
5. Constructs and signs a `BridgeDeliveryReceipt` confirming success or failure
6. Stores the receipt on the relay for the sender to fetch

The sender receives cryptographic confirmation that the bridge processed and attempted delivery of their message.

---

## Identity and Trust

### What trust does a DMCN user place in the bridge?

The bridge is a trusted intermediary only for messages crossing the SMTP boundary. Specifically:

- **Inbound:** The bridge decrypts nothing — it receives plaintext SMTP email, wraps it in an encrypted DMCN envelope. The trust question is whether the bridge's `BridgeClassificationRecord` (SPF/DKIM/DMARC results) is honest. The classification is cryptographically signed by the bridge, so the recipient knows which bridge made the assertion and can decide whether to trust that bridge.

- **Outbound:** The bridge must decrypt the DMCN message to deliver it via SMTP. This is an inherent limitation of bridging to a protocol without end-to-end encryption. The bridge logs a warning when this occurs. Users should be aware that outbound-to-legacy messages lose their end-to-end encryption guarantee at the bridge boundary.

Messages between two DMCN users never involve a bridge and require no trust in any intermediary.

### Can there be multiple bridge operators?

Yes. Any party can run a bridge node. Each bridge has its own DMCN identity (with `BridgeCapability` set), its own domain, and its own signing keys. Users choose which bridge to trust based on the bridge's identity and reputation, similar to choosing an email provider today.
