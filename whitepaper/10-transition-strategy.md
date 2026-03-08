## 10. Transition Strategy: Coexistence with Legacy Email


### 10.1 The Migration Problem


No communication platform has ever achieved mainstream adoption by
requiring users to abandon their existing communication infrastructure.
The transition strategy for the DMCN is built on the principle of
graceful degradation — the system provides maximum value and security
to users communicating with each other on the DMCN, while maintaining
the ability to communicate with legacy email users at reduced security
levels during the transition period.


### 10.2 DMCN-to-DMCN Communication


When both sender and recipient are on the DMCN, messages are fully
encrypted, cryptographically signed, peer-to-peer routed, and spam-free
by protocol. This is the target state of the system and the experience
that should be promoted as the default.


### 10.3 DMCN-to-Email Communication


When a DMCN user sends a message to a legacy email address, the message
passes through a gateway node that translates it to SMTP format for
delivery. The message can include a footer inviting the recipient to
join the DMCN. Security properties are reduced in this path — message
content must be decrypted at the gateway for SMTP delivery — but
sender identity remains verifiable at the gateway level.


### 10.4 Email-to-DMCN Communication


Receiving a message from a legacy email sender requires a verified
gateway address system, where legacy emails pass through a gateway that
performs basic spam filtering and sender verification at the SMTP level
before delivering to the DMCN inbox. Users may also maintain a connected
legacy email address displayed in a separate, clearly labeled section of
their DMCN client.


### 10.5 Progressive Migration Incentives


The transition strategy includes mechanisms that actively incentivize
migration to native DMCN communication: visible trust indicators that
distinguish DMCN-verified messages from legacy email; organizational
compliance features requiring DMCN for sensitive internal
correspondence; and developer APIs allowing third-party applications to
integrate DMCN identity as a communication primitive.

