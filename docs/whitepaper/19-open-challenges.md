## 19. Open Challenges and Research Questions


This whitepaper represents a preliminary investigation into the design
space of a Decentralized Mesh Communication Network. Several significant
challenges remain open and will require further research, prototyping,
and community input.


### 19.1 Scale and Performance


The distributed identity registry and peer-to-peer routing architecture
must be demonstrated to perform adequately at the scale of global email
--- billions of users, hundreds of billions of messages per day. The
performance characteristics of the proposed architecture under realistic
load conditions must be validated through simulation and prototype
deployment.


### 19.2 Key Recovery Without Central Authority


The social recovery model proposed in Section 7 is promising, but its UX
and security properties require careful design and user research. The
threshold for recovery must balance security against the practical
reality that trusted contacts may be unavailable, may themselves lose
access to their accounts, or may be compromised. Alternative recovery
mechanisms should be investigated and compared.


### 19.3 Regulatory and Legal Compliance


End-to-end encrypted communication creates compliance challenges for
regulated industries — financial services, healthcare, law — that
are required to maintain records of communications and provide them to
regulators on demand. The DMCN architecture must provide mechanisms that
allow regulated entities to meet their compliance obligations without
compromising the security properties of the system for other users.

The `ARCHIVE_REQUIRED` policy flag in the Domain Authority Record model (Section 13.5) provides the structural mechanism for compliance archiving under managed domains: outbound messages are BCC-encrypted to an approved archive bridge at send time, giving the domain authority access to a decryptable record without exposing message content to relay infrastructure. Edge cases around legal holds, cross-jurisdiction archiving, and eDiscovery production workflows remain open and require further engagement with compliance counsel and regulators in target jurisdictions.


### 19.4 Governance


A truly decentralized network requires a governance model that allows
the protocol to evolve — to address security vulnerabilities,
incorporate technical improvements, and respond to regulatory changes
--- without any central authority having unilateral control. The
governance model for the DMCN is a critical design question with
significant implications for the network's long-term resilience and
trustworthiness.


### 19.5 Sybil Resistance


While the DMCN's identity model prevents spam from unregistered
senders, it must also resist Sybil attacks — scenarios in which a
malicious actor creates a large number of registered identities to
overwhelm spam defenses. The account creation process must impose
sufficient cost or friction to make large-scale Sybil attacks
uneconomical, without imposing unacceptable burden on legitimate users.


### 19.6 Route Selection Algorithm

The onion routing transport layer specified in Section 15.4.4 requires each client to select three relay nodes for every outbound message path. The hard constraints on this selection — operator diversity, ASN diversity, subnet diversity — are well-defined. The weighted scoring algorithm that balances latency, geographic distribution, reputation, and node capacity within those constraints is not yet empirically validated.

Several open questions require prototype deployment and measurement to resolve:

**Guard node policy.** Section 15.4.4 recommends a stable set of preferred entry nodes rotated every 30 days, following Tor's guard node model. The security benefit of guard nodes — preventing an adversary from gaining statistical entry position through repeated random selection — must be weighed against the privacy cost: a stable entry node that is compromised or surveilled provides persistent observation of the sender's traffic patterns over the rotation period. The optimal rotation interval and guard set size for the DMCN's threat model require empirical analysis.

**Latency versus diversity trade-off.** Geographic diversity in route selection increases the number of distinct network vantage points an adversary requires for traffic correlation, but also increases cumulative path latency. For users in geographically peripheral locations — where the nearest relay nodes are clustered in a small number of regions — enforcing strict continental diversity may impose unacceptable latency penalties. The protocol needs a principled approach to relaxing diversity constraints when strict application would degrade the user experience below acceptable thresholds.

**Relay node directory freshness.** The client's local relay node directory is refreshed on a schedule from the identity registry. A stale directory may include nodes that have gone offline, been delisted, or had their reputation updated since the last refresh. The interaction between directory staleness, route construction failures, and fallback selection behaviour requires careful specification to avoid creating exploitable patterns — for example, an adversary who can predict which nodes a client will fall back to when its preferred nodes are unavailable.

**Adversarial node concentration.** An adversary who operates a large number of relay nodes can increase their probability of appearing in a given path even with operator diversity constraints, by registering nodes under many different operator identities. The extent to which ASN and subnet diversity constraints mitigate this — and whether additional heuristics such as maximum per-operator traffic share are warranted — requires analysis of realistic adversary node deployment strategies.

These questions are research problems as much as engineering problems, and their resolution should draw on the academic literature on anonymous communication network design. Community engagement with the anonymity research community is a recommended early step before the route selection algorithm is finalised.

---

### 19.7 Comparison with Internet Standards Processes and the Path to Standardisation

The DMCN is, at its core, a proposal to replace a foundational internet protocol. That places it in a category with a very small number of historical precedents — and those precedents are instructive about what the standardisation process actually involves, how long it takes, and what the realistic path forward looks like.

#### 19.7.1 How SMTP and Its Successors Were Standardised

SMTP was not the product of a standards body working from first principles. RFC 821, published by Jon Postel in 1982, codified a protocol that was already in use on the ARPANET — standardisation followed deployment, not the reverse. The pattern repeated with every significant addition to the email standards stack: SPF emerged from operator practice and was documented as RFC 4408 in 2006 after years of informal deployment; DKIM was assembled from two competing pre-existing proposals (DomainKeys and Identified Internet Mail) and published as RFC 4871 in 2007; DMARC circulated as an industry specification for three years before reaching the IETF as RFC 7489 in 2015.

The lesson is not that the IETF process is slow — it is that the IETF process is not where internet protocols originate. It is where protocols that have already demonstrated real-world viability go to be documented, reviewed, and given the interoperability guarantees that come with formal standardisation. The IETF's rough consensus model requires that working group participants have implementation experience; a protocol that exists only as a whitepaper cannot progress through the process in any meaningful sense.

This is the honest baseline against which the DMCN's standardisation path should be evaluated. The question is not whether the IETF will eventually standardise the DMCN — it is whether the DMCN can first demonstrate the real-world deployment that makes standardisation meaningful.

#### 19.7.2 What the IETF Process Would Require

If the DMCN were to pursue IETF standardisation, the process would follow a well-established path: an Internet-Draft submitted to a relevant area (most likely the Applications and Real-Time Area, which covers email standards), possible formation of a dedicated working group, a series of revisions through rough consensus among working group participants, IETF Last Call, and publication as one or more RFCs.

This process has several practical implications for the DMCN:

**Multiple documents, not one.** The DMCN as described in this whitepaper would decompose into at minimum five or six separate RFCs: the identity layer, the DHT registry protocol, the message format, the transport and onion routing layer, the bridge protocol, and the trust management framework. Each would be reviewed independently by different sets of experts and would need to stand on its own.

**Running code is required.** IETF culture places significant weight on the principle that standards should describe deployed technology, not speculative design. Working group participants who have implemented competing approaches will challenge design choices that have not been validated through prototyping. The sections of this whitepaper that defer design questions to future research — forward secrecy, Sybil resistance, the route selection algorithm — are precisely the sections that would face the most scrutiny.

**Deep incumbent expertise is present.** The relevant working groups include engineers with decades of operational experience deploying and maintaining the current email infrastructure at scale. This is a double-edged resource: the working group participants who best understand the failure modes of SMTP are also the people best positioned to identify weaknesses in any proposed replacement. The DMCN would need to demonstrate not only cryptographic soundness but also credible answers to the practical interoperability and transition concerns that experienced email operators will raise — concerns that are legitimate engineering questions, not merely institutional resistance to change.

**The timeline is measured in years, not months.** The DKIM standardisation process took approximately four years from the initial Internet-Draft to published RFC. DMARC took three years. A protocol of the DMCN's scope and novelty would realistically take five to ten years to progress from initial Internet-Draft to a set of finalisable RFCs — and that is contingent on active working group engagement and sustained contributor momentum throughout.

#### 19.7.3 Alternative Paths: Pre-Standardisation Deployment

The historical pattern in internet protocol development suggests that the productive near-term path is not IETF standardisation but pre-standardisation deployment — building a working system, validating the design, accumulating operational experience, and producing the implementation evidence that the IETF process requires.

This is precisely the path that Signal followed. The Signal Protocol — which underpins the Double Ratchet algorithm referenced in this whitepaper — has never been standardised by the IETF. It is documented in Signal Foundation's own specifications, has been adopted by WhatsApp, Google Messages, and Meta Messenger, and has been analysed extensively by academic cryptographers. Its legitimacy derives from deployment scale and cryptographic scrutiny, not from RFC status. The IETF's MLS (Messaging Layer Security) protocol, RFC 9420, published in 2023, subsequently provided a standardised group messaging security layer that draws on the same conceptual foundations — but Signal did not wait for it.

The DMCN is in an earlier position than Signal was when it began its deployment phase. The practical near-term standardisation strategy is therefore:

- Publish the protocol specification as a series of versioned technical documents analogous to the Signal Protocol documentation — open, citable, subject to public review, but not dependent on the IETF process for legitimacy
- Engage the academic cryptography and anonymity research communities for peer review of the security model, particularly the threat model and the onion routing design
- Submit Internet-Drafts to relevant IETF working groups as an early signal of intent and to invite technical engagement from the standards community, without committing to the full working group process until deployment experience justifies it
- Build a community of independent implementors — the existence of multiple independent implementations of a protocol is a strong signal of standardisation readiness and a prerequisite for IETF progression

#### 19.7.4 The Governance and Standards Tension

There is a genuine tension between the DMCN's decentralised governance ambitions and the reality of how internet standards are made. The IETF is itself a relatively open body — anyone can submit an Internet-Draft, anyone can join a working group mailing list — but its working practices favour participants who can sustain years of engagement, attend multiple in-person meetings per year, and navigate a culture with its own norms and institutional memory.

A governance model for the DMCN (Section 19.4) that is designed around decentralised community participation will need to interface with this reality. The most successful open protocols have generally had a combination of a small, technically credible core team that drives the standards process, and a broader community that validates the design through deployment and use. The tension between these two groups — and between the community's desire for agility and the IETF's demand for stability and backward compatibility — is a governance challenge that the DMCN will need to resolve explicitly rather than assuming it will resolve itself.

> **Summary Assessment**
> *The IETF standardisation path is open to the DMCN in principle but is not the appropriate immediate priority. The correct sequence is: prototype implementation, peer review of the security model, deployment at meaningful scale, and accumulation of operational experience. Standardisation follows demonstrated viability — it does not precede it. The DMCN should engage the standards community early as a source of technical scrutiny, not as the primary vehicle for legitimacy.*

---
