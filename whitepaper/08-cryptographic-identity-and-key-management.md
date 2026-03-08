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

- Multi-device access is handled through secure key synchronization — a flow identical to the device migration process used by Signal.

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

