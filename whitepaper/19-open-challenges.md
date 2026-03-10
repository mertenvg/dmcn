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


---
