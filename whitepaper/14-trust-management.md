## 14. Trust Management: Whitelists, Greylists, and Blocklists


Cryptographic identity verification is the foundation of the DMCN's
trust model — it answers the question of whether a message genuinely
came from a claimed sender. But verification alone does not answer a
second, equally important question: whether the user actually wants to
hear from that sender. Trust management is the user-facing system that
sits on top of cryptographic verification and allows each participant to
define, on their own terms, who they trust, who they are uncertain
about, and who they actively reject.

The DMCN's trust management system operates at three tiers ---
whitelist, greylist, and blocklist — each with distinct delivery
semantics, key storage implications, and sharing properties. Together
they form a layered defence that is more powerful than anything
available in legacy email, precisely because the identities being
managed are cryptographic and persistent rather than superficial and
easily spoofed.


### 14.1 The Whitelist: Confirmed Trusted Senders


The whitelist is the user's registry of confirmed trusted contacts. It
is not merely an address book — it is a cryptographically anchored
record that binds a human identity to a specific public key, with a
record of how and when that binding was established. Every entry in the
whitelist carries a trust provenance — the mechanism by which the user
confirmed the contact's identity — which is surfaced in the client UI
to help users understand the strength of each trust relationship.


#### 14.1.1 Trust Establishment Mechanisms


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


#### 14.1.2 Key Binding and Update Handling


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

The same notification mechanism fires when a contact's address is deprovisioned by a domain authority — for example, when an employee leaves an organisation and their `@company.com` identity is revoked. Contacts who had that address whitelisted are alerted that the identity is no longer active and are prompted to re-verify before sending further messages. The domain authority revocation model that triggers this behaviour is specified in Section 13.3.


> **Key Change Alert**
> *When a whitelisted contact's public key changes, the DMCN client
> suspends delivery from that identity and alerts the user. No message
> is delivered under an unconfirmed new key until the user explicitly
> re-verifies. This is a deliberate friction point — it is the
> correct response to a high-assurance security event.*


#### 14.1.3 Whitelist Portability and Backup


The whitelist is an asset of significant personal value — it
represents years of accumulated trust relationships. It is therefore
backed up as part of the user's encrypted key material and can be
exported in a standardised, encrypted format for migration between
clients or for archival. The export format includes not only the public
keys and addresses but also the full trust provenance record for each
entry, so that the history of how trust was established is preserved
across migrations.


### 14.2 The Greylist: Unknown but Unblocked Senders


The greylist occupies the space between explicit trust and explicit
rejection. It is the default destination for messages from DMCN-verified
senders who are not yet on the user's whitelist — verified in the
cryptographic sense, meaning their signature is valid and their identity
is registered, but not yet confirmed as trusted by the user.


#### 14.2.1 Greylist Delivery Semantics


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
message without any notification to the sender), or Reject and blocklist
(discarding the message and adding the sender to the blocklist to
prevent future delivery attempts).


#### 14.2.2 Greylist Auto-Resolution Rules


To reduce the burden of manual greylist management, the client supports
configurable auto-resolution rules that can automatically promote or
demote senders based on network signals:

- Auto-promote if vouched by N or more whitelist contacts — configurable threshold, default of 3.

- Auto-promote if sender holds a verified organisational identity matching a domain the user has previously whitelisted.

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


### 14.4 Trust Tier Interaction Summary


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
  Blocklist     sender               dropped at      blocked        
                                     relay                          

  Shared        Community-reported   Dropped per     Yes —        Yes ---
  Reputation    bad actor            feed policy     persistent     community opt-in
  Feed                                               listing        
  ------------- -------------------- --------------- -------------- ----------------

---


---

