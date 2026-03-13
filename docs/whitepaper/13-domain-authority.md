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
| `DOMAIN_APPROVAL_REQUIRED` | 0x01 | Identity records under this domain must carry a valid domain countersignature to be treated as verified |
| `REJECT_UNMANAGED` | 0x02 | Relay nodes must reject messages from addresses under this domain that do not carry a valid domain countersignature |
| `ARCHIVE_REQUIRED` | 0x04 | All messages to/from this domain must be routed through an approved archive bridge (see Section 13.5) |
| `DEVICE_MANAGEMENT_REQUIRED` | 0x08 | Identities under this domain may only be registered from device identifiers pre-approved by the domain authority |
| `SUBDOMAIN_DELEGATION` | 0x10 | Sub-authority delegation is permitted; sub-authority records in this DAR are valid |

---

### 13.2 Identity Provisioning Under a Managed Domain

When `DOMAIN_APPROVAL_REQUIRED` is set, an individual identity record under the domain is only treated as fully verified if it carries a countersignature from the domain authority in addition to the individual's own self-signature.

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

Once a domain revocation record is published to the registry, relay nodes enforcing `DOMAIN_APPROVAL_REQUIRED` or `REJECT_UNMANAGED` will reject messages from the revoked address. The revocation propagates through the DHT within the same timeframe as any other registry update.

The revoked user's underlying DMCN identity — their key pair — is unaffected. They retain their cryptographic identity and can register a new address under a personal domain or provider-hosted address. What they lose is the authorised binding to the organisation's domain.

#### 13.3.2 Key-Change Notification on Revocation

When a domain authority revokes an address, the allowlists of users who had that address as a trusted contact must be updated. The key-change notification mechanism described in Section 14.1.2 fires automatically on revocation: contacts are alerted that the identity previously known as `alice@company.com` is no longer active under that address, and are prompted to re-verify before resuming communication. This prevents a revoked identity from continuing to receive messages intended for the legitimate address holder.

---

### 13.4 Namespace Management and Conflict Resolution

Without a domain authority, nothing prevents two people from both attempting to register as `cfo@company.com`. The first registration wins in the DHT, leaving the domain with no recourse and creating a potential for namespace squatting or impersonation within the domain.

When a DAR is published with `DOMAIN_APPROVAL_REQUIRED`, this problem is structurally eliminated. A registration attempt for an address under a managed domain without the required countersignature is rejected by the registry as unverified. Only addresses that have passed through the domain authority's provisioning flow carry valid countersignatures and are accepted as verified identities under that domain.

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

