## 3. Why This Hasn't Been Done Before

The core technical proposition of the DMCN — cryptographic identity at the protocol level, peer-to-peer routing, spam rejection at the network boundary rather than the inbox — has been understood as the correct class of solution since at least the early 1990s. Phil Zimmermann understood it in 1991. The question a technically sophisticated reader will immediately ask is: if this is the right answer, why hasn't it been built?

The honest response is not that nobody thought of it. It is that five compounding obstacles have, until recently, made a successful execution implausible. Understanding those obstacles is important context for evaluating this proposal — and for understanding why the conditions have changed enough to make a serious attempt viable now.

### 3.1 The Network Effect Is Brutal

Email's value derives almost entirely from the fact that everyone is already on it. A replacement network starts at zero users, which means zero value, which means no rational user switches, which means it stays at zero users. This is not a technical problem — it is a cold-start problem of the kind that has ended every serious attempt at email replacement.

The only communication platforms that have broken through this barrier in recent years — Signal, WhatsApp, Telegram — did so by targeting a different use case (mobile messaging) rather than attempting a direct email replacement, and by benefiting from specific external catalysts: the Snowden revelations, Facebook's acquisition of WhatsApp, and government attempts to ban encrypted messaging. No equivalent catalyst has yet destabilised email incumbency at scale.

Any credible proposal for an email replacement must have a specific, credible answer to the cold-start problem. The DMCN's answer — the SMTP bridge architecture and address portability — is the design element that makes or breaks the proposal's real-world viability. The technical architecture is the easier part.

### 3.2 PGP Poisoned the Well

PGP did the cryptographically correct thing and failed so visibly, for so long, that it created a durable belief in the technical community that secure email is a structural impossibility. That belief suppressed serious engineering investment and venture capital attention for three decades.

The diagnosis was wrong. PGP failed on user experience and key discovery — not on cryptography. The underlying model was sound. But "if PGP couldn't do it, nothing can" became a default assumption that was rarely examined. The result is that the problem space has been underfunded and understaffed relative to its importance, and most serious cryptographic engineering talent has been directed toward problems perceived as more tractable.

### 3.3 Incumbent Platforms Have No Incentive to Fix It

Google and Microsoft collectively handle the majority of global email. Their business models are built on the email ecosystem as it currently exists. Gmail's advertising revenue depends on the ability to analyse message content. A cryptographically private email system with end-to-end encryption would directly undermine that model.

Both companies have invested substantially in spam filtering — which keeps users satisfied enough to prevent churn — but have no commercial incentive to solve the identity problem at the protocol level. Solving it would also solve their surveillance capability. The result is a rational corporate decision to treat spam as a manageable nuisance rather than a solvable structural problem, and to direct engineering resources accordingly.

This dynamic means that a solution to the email identity problem is unlikely to come from the incumbent platforms. It must come from outside.

### 3.4 The Transition Cost Is Asymmetric

Every prior attempt at secure or decentralised email forced users to make a binary choice: adopt the new system and abandon your existing address, or stay on SMTP. For most users and organisations, abandoning an established email address is not a realistic option — it is embedded in contracts, business cards, institutional systems, and years of correspondence. The switching cost is effectively prohibitive.

This asymmetry has killed technically sound proposals repeatedly. A system can be cryptographically superior in every measurable way and still achieve zero adoption if the migration path requires users to start over. The DMCN's address portability and bridge architecture are specifically designed to eliminate this barrier — but it is worth being clear that this design challenge is harder than the cryptographic design challenge, and that no previous attempt has solved it convincingly.

### 3.5 The UX Precondition Was Not Yet Met

The DMCN's central UX claim — that public/private key cryptography can be made entirely invisible to mainstream users — has only become demonstrably true in the last two to three years. The mass deployment of passkeys by Apple, Google, and Microsoft from 2022 onwards established, at a scale of hundreds of millions of users, that people can use elliptic curve cryptography daily without any awareness that they are doing so.

Before that precedent existed, the honest answer to "can you hide cryptographic key management from non-technical users at scale" was "we believe so, but we have not proven it." That uncertainty was a legitimate objection to the entire approach. It is no longer a legitimate objection. The passkey deployment provides a real-world existence proof at the exact scale the DMCN requires.

This is the most important change in the conditions since PGP's failure, and it is recent. It is the primary reason why the timing of this proposal differs from prior attempts.

### 3.6 Blockchain Absorbed the Decentralisation Impulse

From approximately 2017 onwards, the engineering energy and venture capital that might otherwise have been directed at decentralised communication infrastructure was largely absorbed by the blockchain and Web3 ecosystem. Projects that correctly identified the decentralised identity problem attempted to solve it using blockchain infrastructure — which introduced transaction latency, economic friction, and a requirement for cryptocurrency wallet ownership that made mainstream adoption structurally impossible.

The result was a decade in which the correct problem was identified by a large number of well-funded teams, and the wrong architectural choice was consistently made. The DMCN's explicit design decision to achieve decentralisation through peer-to-peer networking rather than blockchain infrastructure is a direct response to this pattern.

### 3.7 The Conditions Have Changed

The five obstacles above have not disappeared. The network effect problem remains the hardest unsolved challenge, and this proposal does not claim to have a guaranteed solution to it. What has changed is:

- The UX precondition is now met, demonstrably, at scale
- The cryptographic primitives required (Curve25519, Ed25519, AES-256-GCM) are mature, hardware-accelerated, and universally available
- The blockchain detour has produced useful lessons about what decentralised identity infrastructure should not be
- Regulatory pressure on email security and data privacy (GDPR, NIS2, increasing BEC enforcement) is creating institutional demand for a more trustworthy email substrate
- The distributed systems infrastructure required to run a global peer-to-peer network — cloud compute, global CDN, open-source DHT implementations — is commoditised in a way it was not in 2005 or even 2015

This is not an argument that the DMCN will succeed where others have failed. It is an argument that the conditions under which a well-designed attempt could succeed are better now than they have been at any prior point. The proposal stands or falls on the quality of its execution — particularly the migration strategy — not on the novelty of its core insight.

> **Historical Assessment**
> *The right technical answer to the email identity problem has been known for thirty years. The barriers to implementation have been predominantly economic, social, and experiential rather than cryptographic. The conditions that sustained those barriers are weaker now than at any previous point. That is the case for attempting this now, stated plainly.*




---
