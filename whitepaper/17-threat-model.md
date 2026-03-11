## 17. Threat Model


This section provides a structured analysis of the threat landscape
facing the Decentralized Mesh Communication Network. For each identified
threat category, we document the nature of the attack, how it manifests
in the existing SMTP email ecosystem, how the DMCN architecture changes
the threat surface, and an honest assessment of whether the DMCN
improves, maintains, or introduces new risk relative to the status quo.


> **Methodology**
> *This threat model follows a structured adversarial analysis
> approach: for each threat, we identify the adversary's goal, the
> attack vector available under SMTP, the attack vector (if any)
> available under DMCN, and the net change in risk. Threats are grouped
> by adversary type: mass senders, targeted attackers, infrastructure
> attackers, and state-level actors.*


### 17.1 Threat Category 1: Spam and Bulk Unsolicited Messaging


#### 17.1.1 Nature of the Threat


Spam is the dominant form of abuse in the current email ecosystem. A
spam operator's goal is to deliver promotional, fraudulent, or
malicious messages to as many recipients as possible at the lowest
possible cost per delivery. The economics are straightforward: even at a
conversion rate of a fraction of a percent, the near-zero marginal cost
of sending email makes spam campaigns profitable.


#### 17.1.2 How SMTP Enables This Threat


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


#### 17.1.3 How DMCN Changes the Threat Surface


The DMCN eliminates the conditions that make spam economically viable at
the protocol level. Every sender must possess a registered cryptographic
identity, and every message must bear a valid signature from that
identity. Relay nodes reject unsigned or unverifiably signed messages
before they enter the network. This imposes a non-trivial, though low,
cost on account creation — and critically, each identity's reputation
is permanent and non-transferable.

A spam operator who wishes to send at scale must create a large number
of registered identities. Each identity that is reported and blocklisted
is permanently lost — there is no equivalent of rotating to a new IP
address. The mathematical relationship between spam volume and identity
cost shifts the economics of spam from profitable to uneconomical at
scale.


#### 17.1.4 Residual Risk and Honest Limitations


The DMCN does not make spam creation infinitely expensive — it makes
it non-zero in cost and permanently cumulative in consequence. A
determined, well-resourced spam operation could potentially automate the
account creation process (a Sybil attack), creating large numbers of
identities before they are reported. Section 17.5 addresses Sybil
resistance specifically. The consent-based inbox model (Section 8.2) provides a secondary layer: even a registered
identity cannot reach a user's primary inbox without meeting one of the
allowlisting criteria.


> **Net Assessment**
> *Spam: Significantly mitigated. The protocol-level economics of spam
> are fundamentally changed. Residual risk exists through Sybil attacks
> and requires robust account creation friction to fully close.*


### 17.2 Threat Category 2: Phishing and Identity Spoofing


#### 17.2.1 Nature of the Threat


Phishing attacks exploit the inability of email recipients to reliably
verify sender identity. An attacker impersonates a trusted entity — a
bank, an employer, a government agency — to induce the recipient to
reveal credentials, transfer funds, or install malware. Business Email
Compromise (BEC) is a sophisticated variant in which attackers
impersonate executives or financial officers within an organisation to
authorise fraudulent wire transfers. BEC alone causes billions of
dollars in losses annually.


#### 17.2.2 How SMTP Enables This Threat


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


#### 17.2.3 How DMCN Changes the Threat Surface


In the DMCN, every message carries a cryptographic signature tied to the
sender's registered key pair. It is mathematically impossible to forge
a message that appears to originate from a registered identity without
possessing that identity's private key. The recipient's client
verifies this signature automatically, and surfaces a clear,
non-technical trust indicator based on the sender's position in the
recipient's trust graph.

A phishing attempt from an unregistered identity cannot enter the
network. A phishing attempt from a registered identity that is not in
the recipient's allowlist lands in the pending queue, visibly
distinguished from trusted messages. The only viable phishing vector in
the DMCN is a fully registered identity that the recipient has
explicitly trusted — which requires the attacker to have established a
prior trust relationship.


#### 17.2.4 Account Compromise and the Key Binding Problem


The DMCN does not eliminate the threat of account compromise — it
changes its character. In SMTP, a compromised email account allows an
attacker to send messages that are indistinguishable from legitimate
messages from that sender. In the DMCN, a compromised account requires
the attacker to have stolen the private key itself, not merely the login
credentials. Private keys stored in hardware-backed secure enclaves (as
specified in Section 7.1) cannot be extracted even if the device's
operating system is compromised.

This represents a meaningful improvement over SMTP account compromise,
but it introduces a new concern: if a private key is stolen (e.g., from
a device without hardware security support), the attacker gains the full
trust relationships of that identity with no visible indicator to
contacts. The allowlist key-change notification system (Section 14.1.2)
partially mitigates this: if the attacker uses a new device, contacts
will be alerted that the key has changed and prompted to re-verify.


> **Net Assessment**
> *Phishing: Substantially mitigated for protocol-level spoofing.
> Account compromise risk shifts from credential theft to key theft ---
> a meaningfully higher bar, but not zero. Hardware key storage is
> essential for this property to hold.*


### 17.3 Threat Category 3: Infrastructure Attacks


#### 17.3.1 Denial of Service Against the Network


Any distributed network is a potential target for denial-of-service
attacks. In the DMCN, the primary infrastructure targets are relay
nodes, the distributed identity registry, and bridge nodes. The goal of
such an attack is to prevent message delivery, disrupt identity lookups,
or degrade the network to the point of unusability.


#### 17.3.2 Comparison with SMTP


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


#### 17.3.3 New Infrastructure Risks Introduced by DMCN


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
users. The federated bridge architecture (Section 11.5) distributes this
risk, but organisations using a single bridge provider remain exposed to
single-point-of-failure risk.


> **Net Assessment**
> *Infrastructure DoS: Comparable to SMTP for distributed attacks;
> improved for centralised attack scenarios due to peer-to-peer
> architecture. The identity registry is a novel attack surface
> requiring careful design. Bridge nodes reintroduce some
> centralisation risk during the transition period.*


### 17.4 Threat Category 4: Relay Node Misbehaviour


#### 17.4.1 Nature of the Threat


Unlike centralised email providers, DMCN relay nodes can be operated by
any party. This raises the question of what happens when a relay node
operator acts maliciously or negligently. Potential misbehaviours
include selectively dropping messages, logging message metadata (even
though content is encrypted), injecting false routing information, or
colluding with other nodes to deanonymise communication patterns.


#### 17.4.2 Comparison with SMTP


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


#### 17.4.3 Metadata Privacy and the Onion Routing Layer


The proposed onion-routing-inspired transport protocol (Section 6.2.2)
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


### 17.5 Threat Category 5: Sybil Attacks


#### 17.5.1 Nature of the Threat


A Sybil attack occurs when a malicious actor creates a large number of
fake identities to subvert a trust-based system. In the context of the
DMCN, the primary Sybil attack scenarios are: creating large numbers of
identities to conduct spam campaigns before they are blocklisted;
creating fake identities to inflate web-of-trust vouching for a
malicious identity; and creating fake identities to manipulate shared
reputation feeds.


#### 17.5.2 Comparison with SMTP


SMTP is essentially infinitely susceptible to Sybil attacks — there is
no meaningful identity system to attack, and the cost of registering a
new sending domain is a few dollars. The DMCN's identity model is
inherently more resistant because it requires account creation friction,
and because blocklisted identities cannot be recovered. However, the
DMCN is not immune, and Sybil resistance is one of the most significant
open design challenges.


#### 17.5.3 Proposed Mitigations


Several mechanisms can be combined to raise the cost of Sybil attacks to
uneconomical levels:

- Account creation friction: requiring a verified phone number, email address, or proof-of-work computation during account creation raises the per-identity cost above zero

- Rate limiting on new identity behaviour: newly created identities are subject to stricter pending queue treatment and lower throughput limits until they have established a reputation history

- Web-of-trust bootstrapping requirements: the consent-based inbox model means that a new identity must earn its way into recipients' allowlists; a Sybil farm of identities with no trust relationships has no delivery capability

- Reputation feed correlation: feed operators can flag clusters of identities with similar creation timing, device fingerprints, or behaviour patterns as likely Sybil farms


> **Net Assessment**
> *Sybil attacks: Improved relative to SMTP, but not fully solved. The
> per-identity cost is non-zero and cumulative in consequence, raising
> the economics of Sybil attacks. Full resistance requires careful
> design of account creation friction — a balance between security
> and accessibility that requires user research and iteration.*


### 17.6 Threat Category 6: State-Level Surveillance and Censorship


#### 17.6.1 Nature of the Threat


Nation-state actors represent the most sophisticated and well-resourced
adversaries in the threat landscape. Their objectives may include mass
surveillance of communications content, targeted surveillance of
specific individuals, censorship of communications between specific
parties, or disruption of communication infrastructure for geopolitical
purposes.


#### 17.6.2 How SMTP Enables State Surveillance


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


#### 17.6.3 How DMCN Changes the Threat Surface


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


#### 17.6.4 Censorship and Network Disruption


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


### 17.7 Threat Category 7: Key Compromise and Recovery Attacks


#### 17.7.1 Nature of the Threat


The security of the entire DMCN identity model rests on the secrecy of
each user's private key. If a private key is compromised, the attacker
gains the ability to read all future messages sent to that identity, to
send messages that appear to come from that identity, and to modify or
revoke that identity's trust relationships. Key compromise is therefore
the most serious category of attack specific to a cryptographic identity
system.


#### 17.7.2 Comparison with SMTP


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


#### 17.7.3 The Social Recovery Attack Surface


The social recovery mechanism (Section 7.3) — in which trusted
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


#### 17.7.4 Key Revocation and Forward Secrecy


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


### 17.8 Threat Category 8: Bridge Node Attacks


#### 17.8.1 Nature of the Threat


The SMTP-DMCN bridge architecture (Section 10) is a necessary component
of any viable transition strategy, but it reintroduces several trust and
security challenges that the native DMCN architecture otherwise
eliminates. Bridge nodes represent the interface between the trustless
SMTP world and the cryptographically verified DMCN, and they must make
trust determinations about SMTP senders that have no direct equivalent
in the native protocol.


#### 17.8.2 Bridge-Specific Attack Vectors


The following attacks are specific to the bridge architecture and have
no equivalent in the native DMCN:

- Content interception on the outbound path: bridge nodes must decrypt outbound messages to re-encode them as SMTP. A malicious or compromised bridge operator gains access to message content in transit. This is disclosed in Section 11.2.2 and is an unavoidable consequence of protocol translation.

- False trust classification: a malicious bridge could misrepresent the trust tier of an inbound SMTP message — for example, classifying a spam message as 'Verified Legacy Sender' to bypass the recipient's filters. The bridge's classification is signed with the bridge's own DMCN key, creating accountability, but this only helps if users verify which bridge they are trusting.

- Legacy spam injection: spammers may target the bridge's SMTP listener as a pathway into DMCN inboxes, attempting to exploit weaknesses in the bridge's SMTP authentication checks to inject messages that would be rejected if sent natively.

- Bridge impersonation: an attacker could operate a bridge that presents itself as trustworthy in the identity registry but maliciously handles messages. Mitigated by requiring bridge operators to publish their security practices and undergo periodic audits.


#### 17.8.3 Bridge Risk as a Transitional Concern


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


### 17.9 Threat Model Summary


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
> eliminated — Sybil resistance, relay metadata leakage, key recovery
> attacks, and bridge node risks — are well-understood engineering
> challenges with known mitigation approaches, rather than fundamental
> architectural flaws. None of these residual risks represents a
> regression relative to the current SMTP-based email ecosystem.*
>
> *Organisations deploying the Domain Authority Record model (Section 13) introduce an additional threat surface: compromise of the domain authority key. This threat is analogous to enterprise CA root key compromise and should be treated with equivalent operational rigour. Mitigations are specified in Section 13.7.*

---


---

