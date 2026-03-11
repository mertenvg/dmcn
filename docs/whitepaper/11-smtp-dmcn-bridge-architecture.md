## 11. The SMTP-DMCN Bridge Architecture


A Decentralized Mesh Communication Network that cannot communicate with
the existing email ecosystem is, for practical purposes, a closed
system. The SMTP-DMCN Bridge is the infrastructure component that makes
gradual, real-world adoption possible — allowing DMCN users to send
and receive messages with the 4 billion people on legacy email, without
compromising the security properties of native DMCN communication.

This section provides a detailed architectural treatment of the bridge,
covering both directions of message flow, the trust model the bridge
operates under, its internal components, and the mechanisms by which it
preserves as much security as possible within the constraints of
protocol translation.


### 11.1 Architectural Role and Placement


Bridge nodes occupy a distinct position in the DMCN topology. They are
DMCN-native nodes — registered in the identity layer with their own
cryptographic identity — that additionally operate SMTP listener and
sender services on the legacy internet. A bridge node is, in effect,
bilingual: it speaks DMCN natively and SMTP at the boundary.

Bridge nodes can be operated by the DMCN foundation, by commercial
providers, or by organizations that wish to run their own bridge
infrastructure for compliance or privacy reasons. The bridge is not a
single central server — multiple independent bridge operators can
coexist, and DMCN users can choose which bridge operator handles their
legacy email traffic, much as organizations today choose their own mail
exchanger (MX) records.


> **Design Principle**
> *The bridge is a temporary necessity during the transition period,
> not a permanent fixture of the architecture. Every design decision
> prioritizes making native DMCN adoption more attractive, so that
> dependence on the bridge diminishes naturally over time as the
> network grows.*


### 11.2 Outbound Path: DMCN to SMTP


When a DMCN user addresses a message to a legacy email address, the
message follows a modified delivery path through a bridge node. This is
the simpler of the two directions, and the security properties are
well-defined.


#### 11.2.1 Message Flow


- The sender's DMCN client composes and signs the message with the sender's private key, as in a standard DMCN message.

- The client detects that the recipient address resolves to a legacy email address (no DMCN public key found in the identity registry) and routes the message to a bridge node rather than the standard transport layer.

- The bridge node receives the encrypted, signed DMCN message, verifies the sender's signature against the identity registry, and decrypts the message content using the bridge's own private key (the message is re-encrypted to the bridge's key rather than a non-existent recipient key).

- The bridge constructs a standard SMTP message from the decrypted content, applying DKIM signing using the bridge operator's domain key. The `From` header is set to the sender's human-readable DMCN address (e.g. `alice@mycompany.com`), so that the recipient's email client displays the sender's familiar identity. The `Sender` header is set to the bridge's own address (e.g. `bridge@bridge.dmcn.net`), identifying the mailbox actually responsible for transmission per RFC 5322. The `Reply-To` header is set to the sender's DMCN bridge receive address, so that replies from legacy clients are routed back through the bridge correctly. This combination — `From` showing the human author, `Sender` showing the transmitting agent — is the same pattern used by legitimate bulk email providers (such as Mailchimp and SendGrid) when sending on behalf of their customers. SPF and DKIM pass against the bridge's own domain; the `From` address is not required to be an authorised sender for those checks under standard DMARC evaluation when the `Sender` header is present.

- The bridge delivers the SMTP message to the recipient's mail server using standard MX lookup and SMTP relay.

- A standardized footer is appended to the message body, displaying the sender's verified DMCN identity and an invitation link for the recipient to join the DMCN and receive future messages natively.


#### 11.2.2 Trust Properties


The outbound path involves one unavoidable trust compromise: the bridge
node decrypts message content in order to re-encode it as SMTP. This
means the bridge operator has technical access to message content in
transit. This is an honest limitation that must be clearly disclosed to
users and is analogous to the trust placed in any email service provider
today.

This limitation is mitigated by several factors: bridge operators are
registered DMCN identities with cryptographic accountability; users can
choose their bridge operator and can migrate between operators;
organizations with strong confidentiality requirements can operate their
own bridge nodes; and the limitation affects only messages sent to
legacy email recipients, not native DMCN-to-DMCN communication.


### 11.3 Inbound Path: SMTP to DMCN


Receiving messages from legacy email senders is the more complex
direction, because the sender has no DMCN identity and no cryptographic
signing capability. The bridge must make a trust determination about an
unverified sender and communicate that determination clearly to the
recipient.


#### 11.3.1 Addressing


DMCN users who wish to receive legacy email obtain a bridge address ---
a standard email address managed by the bridge operator (e.g.,
username@bridge.dmcn.net). This address is registered as an MX record
pointing to the bridge node's SMTP listener. Users can publish this
address on websites, business cards, and legacy systems as their contact
point for people who have not yet adopted DMCN.


#### 11.3.2 Message Flow


- The bridge's SMTP listener receives an inbound message addressed to the user's bridge address.

- The bridge performs a full suite of legacy authentication checks: SPF validation, DKIM signature verification, DMARC policy evaluation, and reverse DNS lookup on the sending mail server.

- The bridge queries distributed reputation databases (analogous to existing RBL/DNSBL services) for the sending IP address and domain.

- Messages that fail hard authentication checks (invalid DKIM, SPF failure with DMARC reject policy) are dropped with a delivery failure response to the sending server.

- Messages that pass or partially pass authentication are classified into trust tiers: Verified Legacy Sender (valid DKIM + positive reputation), Unverified Legacy Sender (no DKIM or neutral reputation), and Suspicious Legacy Sender (reputation flags present).

- The bridge wraps the classified message in a DMCN envelope, signed by the bridge's own private key as an attestation of the classification outcome. The DMCN envelope includes the full authentication result metadata, the trust tier assignment, and the original SMTP headers.

- The wrapped message is encrypted to the recipient's DMCN public key and delivered through the standard DMCN transport layer to the recipient's inbox.


#### 11.3.3 Recipient Experience


Inbound legacy messages appear in a clearly distinguished section of the
recipient's DMCN client, visually separated from native DMCN messages.
Each message displays its trust tier — Verified Legacy, Unverified
Legacy, or Suspicious — with a plain-language explanation of what the
classification means. The recipient can promote a legacy sender to their
contact list (which triggers the DMCN client to send the sender an
invitation to join the network natively) or block the sender
permanently.


### 11.4 Bridge Node Security Model


Bridge nodes are high-value infrastructure components and require a
rigorous security model. Several specific threats must be addressed:

- Bridge impersonation — a malicious actor operating a fraudulent bridge that misrepresents message authentication results. Mitigated by requiring bridge operators to register their identity in the DMCN identity registry, publish their security practices, and submit to periodic independent audits.

- Content interception — a bridge operator reading or modifying message content in transit on the outbound path. Mitigated by end-to-end message signing (recipients can verify the sender's signature even after bridge translation), audit logging, and regulatory accountability for commercial operators.

- SMTP relay abuse — spammers attempting to use the inbound bridge path to inject spam into DMCN inboxes. Mitigated by the authentication classification system and by rate limiting on the bridge's SMTP listener per sending domain and IP.

- Bridge node compromise — an attacker gaining control of a bridge node. Mitigated by the decentralized bridge model (no single bridge handles all traffic), key rotation protocols, and the ability for users to revoke trust in a specific bridge operator.


### 11.5 Federated Bridge Architecture


To avoid reintroducing centralization through the bridge layer, the DMCN
specification defines an open Bridge Operator Protocol (BOP) that any
qualified operator can implement. Bridge operators are discoverable
through the DMCN identity registry, and DMCN clients implement automatic
bridge selection based on operator reputation, geographic proximity,
organizational policy, and user preference.

This federated model mirrors the existing email MX record system — the
delivery path is determined by the recipient's published preferences,
not by any central routing authority — while adding the cryptographic
accountability layer that SMTP lacks. An organization can run its own
bridge, use a commercial bridge provider, or designate different bridges
for inbound and outbound traffic.


### 11.6 Precedents and Comparable Systems


The bridge architecture is a well-established pattern in communication
infrastructure. Several analogous systems demonstrate that protocol
translation at scale is an engineering challenge with proven solutions:

- Matrix bridges (Synapse) — the Matrix protocol has operated production bridges to Slack, Discord, WhatsApp, Telegram, SMS, and IRC simultaneously, with hundreds of thousands of active users. The Matrix bridge architecture provides a detailed precedent for federated, bidirectional protocol translation.

- Email-to-SMS gateways — carriers have operated SMTP-to-SMS translation services for decades, handling billions of messages. The protocol mismatch between email and SMS is in many respects more severe than the mismatch between SMTP and DMCN.

- SIP-PSTN gateways in VoIP — Voice over IP systems routinely bridge between SIP-native networks and the legacy public switched telephone network, including cryptographic signaling translation. The architecture for handling identity and trust across this boundary is directly analogous to the SMTP-DMCN bridge problem.


> **Engineering Assessment**
> *The SMTP-DMCN bridge is a significant but tractable engineering
> project. It does not require novel cryptographic research — it
> requires careful implementation of known patterns applied to a new
> protocol boundary. A working prototype bridge is a realistic
> deliverable for an initial proof-of-concept phase.*



---

