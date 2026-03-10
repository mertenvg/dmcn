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

Each DMCN identity has exactly one **primary key pair** at any point in time. The primary key is the canonical representation of the identity: it is what appears in the identity registry, what the trust graph is anchored to, what whitelists are bound to, and what contacts see when they look up an address. The primary key's public half is the stable long-term identifier for the address.

The primary key is generated on the user's first device at account creation and is backed up through the encrypted key backup infrastructure described in Section 7.3. It is never generated or held by any server or relay node.

#### 7.5.2 Device Sub-Keys

Each device on which the user activates their DMCN account generates its own **device sub-key pair** locally. The device sub-key is signed by the primary key, creating a cryptographically verifiable delegation: any party who trusts the primary key can verify that the sub-key is a legitimately authorised device for that identity.

Device sub-keys serve two functions. First, they are the keys actually used for per-message encryption and signing on that device — the primary key is not used for routine message operations, reducing its exposure. Second, they allow per-device revocation: if a device is lost, stolen, or decommissioned, only that device's sub-key need be revoked. The primary key, the identity's trust relationships, and all other active devices are unaffected.

Sub-keys are registered in the identity registry as children of the primary key record and are returned alongside the primary key in registry lookups. When sending to a multi-device user, the message payload is encrypted exactly once using a single randomly generated symmetric content key. That content key is then individually wrapped (encrypted) for each active device sub-key. Any device that holds the corresponding private sub-key can unwrap the content key and decrypt the payload. This approach — a Key Encapsulation Mechanism (KEM) pattern — ensures that message and attachment payload bytes appear on the wire only once regardless of how many devices the recipient has enrolled, eliminating the per-recipient payload duplication that would otherwise result. The protocol structure for this is defined in Section 15.3.3.

#### 7.5.3 Sub-Key Lifecycle

Sub-keys are created when a new device is activated and revoked when a device is decommissioned. Revocation is published to the identity registry and propagated through the same mechanism as primary key revocation. A revoked sub-key cannot be re-activated; a replacement device generates a new sub-key and re-registers it under the same primary key.

Sub-keys carry an optional `expires_at` field, enabling organisations to enforce periodic rotation of device credentials — relevant in managed domain contexts where the domain authority (Section 13) may require regular key refresh as a compliance control.

#### 7.5.4 Key Rotation

Periodic rotation of the primary key — whether on a schedule, following a suspected compromise, or as a policy requirement — is handled by publishing a new primary key to the registry via the `UPDATE` operation (Section 15.2.4), signed by both the old and new primary keys to prove continuity of control. This dual-signature rotation triggers key-change notifications to all whitelisted contacts (Section 14.1.2), prompting them to re-verify before the whitelist binding is updated.

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

