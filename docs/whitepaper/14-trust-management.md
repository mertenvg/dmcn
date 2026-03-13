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

Messages held in the pending queue are presented in a section of the client visually distinct from the primary inbox.

**Protocol-level content visibility rule:** For any message whose sender does not hold an effective allowlist entry with the recipient, the client MUST display the sender's verified DMCN identity and the message subject line, and MUST NOT render the message body or any attachment content until the sender is allowlisted. This is a protocol requirement, not a client implementation guideline. Clients that expose message body content prior to allowlisting are non-conformant.

The rationale for this rule is twofold. First, it prevents senders from using message content as a vector to manipulate the recipient's trust decision — unsolicited content cannot be used to manufacture urgency, distress, or social pressure before the recipient has consented to receive it. Second, it provides a consistent privacy guarantee for senders: message content is not exposed to infrastructure or clients that have not yet established a trust relationship with the sender's identity.

From the pending queue the user has four options for each pending message: Accept and allowlist the sender (promoting all future messages to the primary inbox and unlocking content for the current message), Accept this message only (unlocking content for this message without allowlisting the sender, so future messages from the same sender return to the pending queue), Reject and ignore (discarding the message without any notification to the sender), or Reject and blocklist (discarding the message and adding the sender's cryptographic identity to the personal blocklist).


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

### 14.5 Guardian Trust Controls

The trust management model described in Sections 14.1–14.3 assumes full user autonomy: each identity manages its own allowlist, pending queue, and blocklist independently. This assumption is appropriate for adult users operating their own identities, but it does not accommodate accounts provisioned for minors or other dependants, where a supervising party — a parent, legal guardian, or family domain authority — has a legitimate interest in overseeing trust relationship formation.

The DMCN addresses this through a **guardian trust control** model: a lightweight extension to the identity record that places allowlist entries under countersignature requirement, while leaving the child's cryptographic identity fully intact and their blocklist under their own control. The guardian never holds the child's private key, never has access to message content, and cannot read the child's trust list. Their authority is limited to approving or rejecting the child's requests to form new trust relationships.

> **Design Principle**
> *Guardian trust controls give supervising parties authority over trust relationship formation, not over identity or message content. The child's cryptographic self-sovereignty is preserved. The guardian countersigns allowlist entries; they do not co-hold keys.*

---

#### 14.5.1 The Guardian Policy Record

Guardian trust controls are activated by the presence of a `guardian_policy` record attached to an identity record in the DHT. This record is published by the domain authority at address provisioning time and is signed by the domain authority's Ed25519 key.

```
guardian_policy {
    version:                    uint32
    subject_address:            string       // the address this policy applies to
    subject_ed25519_pubkey:     bytes[32]    // Ed25519 public key of the subject identity
    guardian_ed25519_pubkey:    bytes[32]    // Ed25519 public key of the countersigning authority
    guardian_address:           string       // DMCN address of the guardian identity
    valid_until:                uint64       // Unix timestamp; 0 = indefinite
    policy_flags:               uint32       // bitmask; see Section 14.5.2
    domain_authority_signature: bytes[64]    // Ed25519 signature by the DAR authority over all preceding fields
}
```

The `guardian_ed25519_pubkey` identifies the key whose countersignature is required for allowlist entries to become effective. This is typically the domain authority key (for a family provider scenario) or a parent's own DMCN identity key (for a self-hosted family domain scenario). Both are valid; the protocol makes no distinction between them.

The `domain_authority_signature` anchors the policy to the domain's administrative authority. A `guardian_policy` record without a valid domain authority signature is rejected by conformant clients.

---

#### 14.5.2 Policy Flags

The `policy_flags` bitmask supports the following flags for the guardian policy:

| Flag | Value | Behaviour |
|---|---|---|
| `GUARDIAN_APPROVAL_REQUIRED` | `0x01` | Allowlist entries require guardian countersignature to become effective |
| `GUARDIAN_NOTIFY_PENDING` | `0x02` | Guardian receives a notification when a new sender reaches the subject's pending queue |
| `GUARDIAN_CONTENT_GATE_STRICT` | `0x04` | Message content is withheld from the subject until the guardian countersigns (guardian-first mode); without this flag, the subject sees sender and subject line and initiates the approval request (child-initiated mode) |

`GUARDIAN_APPROVAL_REQUIRED` is the foundational flag. The remaining flags extend it. A `guardian_policy` record with no flags set is valid but has no effect.

---

#### 14.5.3 Allowlist Countersignature Requirement

When `GUARDIAN_APPROVAL_REQUIRED` is set, an allowlist entry added by the subject identity is **pending** until it carries a valid countersignature from the `guardian_ed25519_pubkey`. A pending allowlist entry exists in the subject's allowlist but is inert: the sender named in the entry is treated as unknown (pending queue behaviour) until the countersignature is provided.

A conformant allowlist entry under guardian control carries an additional field:

```
guardian_countersignature {
    guardian_ed25519_pubkey:    bytes[32]    // must match guardian_policy.guardian_ed25519_pubkey
    approved_at:                uint64       // Unix timestamp of approval
    signature:                  bytes[64]    // Ed25519 signature by guardian over:
                                             //   allowlist_entry fields +
                                             //   guardian_ed25519_pubkey +
                                             //   approved_at
}
```

An allowlist entry is effective if and only if:
- The `guardian_countersignature` is present, and
- The `guardian_ed25519_pubkey` in the countersignature matches the `guardian_ed25519_pubkey` in the active `guardian_policy` record, and
- The signature is valid

A client receiving a message from a sender whose allowlist entry is pending (countersignature absent or invalid) applies standard pending queue behaviour for that sender, regardless of the allowlist entry's presence.

**Blocklist entries do not require countersignature.** Either the subject or the guardian may add a sender to the blocklist unilaterally. A sender blocked by the guardian is treated as blocked for the subject even if the subject has not independently blocklisted them. This union behaviour means guardian blocklist entries propagate downward without requiring any action from the subject.

---

#### 14.5.4 Approval Flow

The approval flow under `GUARDIAN_APPROVAL_REQUIRED` operates as follows:

1. An unknown sender's message arrives and is held in the subject's pending queue. The subject's client displays the sender's verified DMCN identity and subject line per the content visibility rule in Section 14.2.1. Message body content is not rendered.

2. The subject may initiate an approval request: a signed notification sent to the guardian's DMCN address identifying the sender and requesting countersignature of an allowlist entry.

3. The guardian's client presents the approval request with the sender's verified identity, trust provenance signals (web-of-trust vouching, shared reputation feed status), and the subject line of the pending message. The guardian approves or rejects.

4. On approval, the guardian's client produces a `guardian_countersignature` and delivers it to the subject's client through the standard DMCN message transport. The subject's client attaches the countersignature to the allowlist entry, which becomes effective immediately.

5. On rejection, the subject's client may optionally notify the subject that the request was declined, at the guardian's discretion.

If `GUARDIAN_CONTENT_GATE_STRICT` is set, step 2 is modified: the approval request is sent to the guardian automatically on message arrival, and the subject is not notified of the pending message until the guardian has made a decision. This implements a guardian-first content gate appropriate for younger children, where the subject should not be aware of unsolicited contact attempts until a guardian has reviewed them.

---

#### 14.5.5 Policy Transition and Lapse

The `valid_until` field in the `guardian_policy` record governs automatic lapse. When the current time exceeds `valid_until`:

- The `guardian_policy` record becomes inert. Clients treat it as absent.
- No action is required from either the subject or the guardian for lapse to take effect.
- Existing allowlist entries that carry valid countersignatures remain effective — the approval already occurred and is not retroactively withdrawn.
- New allowlist entries no longer require countersignature. The subject gains full trust autonomy.

A `valid_until` value of `0` indicates an indefinite policy with no automatic lapse. The domain authority may explicitly revoke or modify the `guardian_policy` at any time before `valid_until` by publishing an updated identity record. Early revocation grants the subject full autonomy ahead of schedule; extension defers it.

> **Protocol Scope Note**
> *The protocol does not mandate any particular `valid_until` value or enforce age-appropriate transition. A domain authority may set `valid_until: 0` for any identity. Determination of appropriate policy duration is a matter for the domain authority, the provider's terms of service, and applicable jurisdictional requirements. The DMCN provides the mechanism; it does not prescribe the policy.*

---

#### 14.5.6 Relationship to Domain Authority Flags

Guardian trust controls use a distinct flag vocabulary from the domain-level provisioning controls defined in Section 13. The two flag sets operate at different layers and must not be conflated:

| Flag | Layer | Governs |
|---|---|---|
| `DOMAIN_APPROVAL_REQUIRED` | Domain Authority Record (DAR) | Whether an identity may be provisioned under the domain namespace |
| `GUARDIAN_APPROVAL_REQUIRED` | Identity record (`guardian_policy`) | Whether a provisioned identity may form trust relationships autonomously |

An account provisioned for a child under a family domain would typically carry both: `DOMAIN_APPROVAL_REQUIRED` ensures the child's address is domain-authority-countersigned at registration (consistent with all managed identities under the domain), and `GUARDIAN_APPROVAL_REQUIRED` imposes ongoing trust relationship supervision specific to that identity. An adult employee account under a corporate domain would carry `DOMAIN_APPROVAL_REQUIRED` but not `GUARDIAN_APPROVAL_REQUIRED`.

The flags are orthogonal. Either may be set independently of the other.

---


---

