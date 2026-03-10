## 6. Proposed Architecture: Decentralized Mesh Communication Network


### 6.1 Design Principles


The Decentralized Mesh Communication Network (DMCN) is designed around a
set of foundational principles derived from the failure analysis of
prior approaches. These principles are architectural constraints that
shape every design decision.

- Identity is cryptographic and self-sovereign. Every participant in the network has a unique identity derived from a public/private key pair. Identity is not assigned by any central authority and cannot be revoked by any third party.

- Sender verification is mandatory and protocol-enforced. A message without a valid cryptographic signature from a registered identity cannot enter the network. Verification is not optional, not opt-in, and not a filter applied after the fact — it is a gate at the point of transmission.

- The network is peer-to-peer with no central routing authority. Messages are relayed through a distributed mesh of nodes. No single entity controls routing, storage, or delivery.

- Complexity is hidden from end users. The cryptographic machinery that makes the network trustworthy operates entirely below the user experience layer. Users interact with human-readable identities and familiar communication patterns.

- Legacy email interoperability is a first-class requirement. The network must be capable of sending to and receiving from legacy SMTP addresses during a transition period.


### 6.2 Network Architecture


The DMCN consists of three logical layers, each with distinct
responsibilities:


#### 6.2.1 Identity Layer


The Identity Layer is responsible for the creation, storage, and
discovery of cryptographic identities. Each user identity is represented
by an elliptic curve key pair. The public key serves as the user's
network identifier and is registered in a distributed identity registry
--- a peer-to-peer data structure analogous to a distributed hash table
(DHT) — that allows any node to resolve a human-readable address to
its corresponding public key.

Human-readable addresses follow a format similar to traditional email
(user@domain) but resolve not to a mail server IP address but to a
public key and a set of authorized relay nodes. This means users can be
addressed in a familiar way while the underlying identity is
cryptographic and decentralized.


#### 6.2.2 Transport Layer


The Transport Layer is responsible for routing messages through the mesh
network. Messages are addressed to the recipient's public key,
encrypted with that public key, signed with the sender's private key,
and relayed through a network of nodes using an onion-routing-inspired
protocol that provides metadata privacy in addition to content privacy.

Relay nodes can verify message signatures against the identity layer,
which is the mechanism by which spam is rejected at the network level
--- a node that receives a message signed by an identity not registered
in the identity layer drops the message without delivery.


#### 6.2.3 Storage and Delivery Layer


Unlike real-time messaging systems, email is inherently asynchronous.
The Storage and Delivery Layer provides distributed, encrypted message
storage that holds messages until the recipient's client connects to
retrieve them. Messages are stored encrypted with the recipient's
public key; relay nodes providing storage cannot read message content.


### 6.3 Message Lifecycle


A message in the DMCN follows this lifecycle:

- The sender's client composes a message and signs it with the sender's private key.

- The client encrypts the signed message with the recipient's public key, retrieved from the Identity Layer.

- The encrypted, signed message is submitted to the transport layer with the recipient's public key as the address.

- Relay nodes verify the sender's signature against the Identity Layer. Messages with invalid or absent signatures are dropped.

- The message is routed through the mesh to relay nodes serving the recipient's address, where it is held in encrypted storage.

- When the recipient's client connects, it retrieves and decrypts messages using the recipient's private key.

- The recipient's client verifies the sender's signature, confirming the message genuinely originated from the stated sender.



---
