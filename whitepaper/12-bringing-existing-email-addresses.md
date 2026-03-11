## 12. Bringing Existing Email Addresses to the DMCN


One of the most significant friction points in any transition away from
legacy email is the email address itself. For most people and
organizations, an email address is not merely a routing string — it is
a persistent identity, published on business cards, embedded in
contracts, known to colleagues and clients accumulated over years or
decades. Any system that requires users to abandon their existing
address in order to participate faces an adoption barrier that is
largely insurmountable. The DMCN's address portability feature directly
addresses this by allowing users to bring their existing email addresses
into the DMCN identity layer without abandoning them.


### 12.1 The Principle of Address Portability


Address portability in the DMCN means that an existing email address can
be registered as a DMCN identity anchor — a verified link between a
known email address and a cryptographic public key. Once this link is
established and published in the DMCN identity registry, the address
functions simultaneously as a legacy SMTP address and as a native DMCN
identity. Senders who are on the DMCN will automatically discover the
cryptographic key and send natively; senders on legacy email will
continue to reach the address through conventional SMTP delivery.

This dual-mode operation is the bridge between the old world and the
new. It requires no change of address, no notification to existing
contacts, and no interruption of legacy email delivery. The upgrade is
invisible to legacy senders and automatic for DMCN senders.


### 12.2 Verification Mechanisms


The strength of the link between an email address and a DMCN identity
depends on the level of control the user has over the address and its
underlying domain. The DMCN supports three verification tiers, each with
distinct trust properties:


#### 12.2.1 Provider-Hosted Address Verification (e.g., Gmail, Outlook)


For addresses hosted by third-party providers — the most common case
for individual users — verification proceeds through an email
confirmation flow analogous to standard account verification practices:

- The user claims ownership of their existing address (e.g., alice@gmail.com) within the DMCN client.

- The DMCN identity service sends a time-limited, single-use verification code to that address via the bridge's outbound SMTP path.

- The user retrieves the code from their legacy inbox and enters it in the DMCN client, confirming they control the address.

- The DMCN identity registry publishes a signed binding record linking the address to the user's public key.

This tier provides practical ownership verification — proof that the
user can receive mail at the address — but does not constitute
cryptographic domain ownership. The binding is valid as long as the user
retains control of the legacy account.


#### 12.2.2 Custom Domain Address Verification


For users and organizations who control their own domain (e.g.,
alice@mycompany.com), a stronger verification path is available that
mirrors the DNS-based ownership proof used by DKIM and SPF:

- The user requests a DMCN verification token for their domain from the DMCN identity service.

- The user publishes the token as a DNS TXT record at a standardized subdomain (e.g., \_dmcn.mycompany.com), alongside their public key fingerprint.

- The DMCN identity service performs a DNS lookup to confirm the record is present and correctly formatted.

- The binding is published in the identity registry with a Domain-Verified status, indicating that the link is backed by DNS control rather than merely inbox access.

Domain-verified bindings are significantly more robust. They can be
updated by the domain owner at any time through DNS changes, they do not
depend on the policies of any email provider, and they are resistant to
account suspension or provider-side interference. For organizations,
this is the recommended verification path.

Domain ownership verification establishes *that* an organisation controls a domain. It does not address *how* the organisation governs which addresses are authorised under that domain, or how it provisions and deprovisions staff identities. Those operational requirements are addressed by the Domain Authority Record model in Section 13.


#### 12.2.3 DANE-Style Cryptographic Domain Binding


For domains that have enabled DNSSEC — the cryptographic extension to
DNS that provides tamper-evident records — a third verification tier
is available that provides the highest level of assurance. In this
model, the domain owner publishes the DMCN public key directly in a
DNSSEC-signed record, creating a chain of cryptographic trust from the
DNS root through the domain to the individual identity. This approach is
analogous to DANE (DNS-based Authentication of Named Entities), which is
already used in some contexts to bind TLS certificates to domain names
without relying on certificate authorities.


### 12.3 Trust Implications of Each Tier


  ------------------- ---------------- ---------------- -----------------
  **Verification      **Proof of       **Resistant to   **Recommended
  Tier**              Control**        Provider         For**
                                       Action**         

  Provider-Hosted     Inbox access     No — provider  Individual users
  (Gmail, Outlook)    only             can suspend      during transition

  Custom Domain DNS   DNS record       Yes — domain   Professionals,
  Verification        control          owner controls   small businesses
                                       DNS              

  DNSSEC / DANE       Cryptographic    Yes — highest  Enterprises,
  Cryptographic       DNS chain        assurance        regulated
  Binding                                               industries
  ------------------- ---------------- ---------------- -----------------


### 12.4 The Honest Limitation: Ownership vs. Control


Address portability is a powerful adoption mechanism, but it requires an
honest disclosure that the whitepaper must not obscure: bringing a
provider-hosted address to the DMCN does not give the user cryptographic
ownership of that address at the domain level. Google still controls
\@gmail.com. Microsoft still controls \@outlook.com. If a user's Google
account is suspended, terminated, or if Google elects to block
DMCN-related traffic, the legacy delivery path for that address breaks
--- though the user's DMCN identity and their cryptographic key persist
independently.

This distinction has practical consequences that users should understand
at onboarding. Provider-hosted address linking is a convenience feature
for the transition period. For users who want long-term,
provider-independent identity, the DMCN should actively encourage
migration to a custom domain. The client application can surface this
recommendation appropriately — not as a barrier, but as a path to
greater identity sovereignty over time.


> **Identity Sovereignty Principle**
> *A provider-hosted address gives you a key to a house you rent. A
> custom domain address gives you a key to a house you own. The DMCN
> supports both — but only one provides true long-term identity
> independence.*


### 12.5 Address Portability and the Spam Problem


Address portability introduces one additional consideration for the spam
model: a user who verifies ownership of an existing email address brings
with them the reputation — positive or negative — associated with
that address in legacy spam databases. The DMCN identity layer should
initialize the reputation of a newly verified address using available
legacy reputation signals as a starting point, rather than treating
every ported address as a clean slate.

Conversely, address portability is a meaningful barrier to spam identity
laundering. A spammer who wishes to port a known-good email address to a
DMCN identity in order to inherit its reputation must actually control
that address — they cannot simply claim it. This is a meaningfully
higher bar than the trivially low cost of sending SMTP mail from an
arbitrary claimed address.


### 12.6 Precedents


The address portability model draws on several well-established
precedents in both identity verification and email infrastructure:

- Keybase — demonstrated the viability of linking a cryptographic identity to multiple existing identities (email, Twitter, GitHub, domain) through a system of cryptographic proofs. The DMCN's verification model is conceptually similar but integrated at the protocol level rather than as a third-party overlay.

- Google Workspace and Microsoft 365 custom domain onboarding — both services use DNS TXT record verification to prove domain ownership before allowing custom domain email. The DMCN's Domain Verification tier follows this exact pattern, which is already familiar to IT administrators globally.

- DKIM public key DNS records — the practice of publishing cryptographic keys in DNS records is already standard email infrastructure. The DMCN's DANE-style binding extends this established pattern.

- Number portability in mobile telephony — the telecommunications industry solved an analogous problem when it allowed consumers to bring their phone numbers between carriers. The lesson from that transition is directly applicable: portability dramatically lowers switching costs and accelerates adoption of superior infrastructure.

Address portability answers the question of how an existing address enters the DMCN identity layer. For organisations deploying DMCN across a workforce, a further layer is required: ongoing governance of who is authorised to operate under the domain, how new addresses are provisioned, and how access is revoked when people leave. This organisational layer is specified in Section 13.



---

