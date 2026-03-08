## 9. User Experience: Hiding Complexity Without Sacrificing Security


### 9.1 The Fundamental Principle


The history of secure communication tools is, in large part, a history
of UX failures. The DMCN's user experience layer is designed around a
single principle: the security model must be invisible to the user in
normal operation. Users should experience the DMCN as a familiar
messaging application — with inboxes, contacts, compose windows, and
threads — and should never encounter the words 'key', 'signature',
'certificate', or 'encryption' in the normal flow of using the
product.


### 9.2 Familiar Addressing


Users are addressed using a format that mirrors traditional email: a
local identifier and a domain, separated by the @ symbol. Internally,
this address resolves to a public key — but from the user's
perspective, it is simply their address, just as a phone number is
simply a phone number without any awareness of the SS7 routing protocol
underneath.


### 9.3 Onboarding Flow


Account creation in the DMCN is designed to be comparable in friction to
creating a Gmail account. The user provides a chosen identifier,
authenticates with biometric or PIN, and the application generates their
key pair silently in the background. The entire process takes under two
minutes. There is no seed phrase, no key download, no certificate
request, and no cryptographic terminology.


### 9.4 Contact Discovery


Finding another user on the DMCN requires only their address. The
application resolves the address to a public key through the distributed
identity registry, and the contact appears in the user's contact list.
All messages to that contact are automatically encrypted and signed. The
user does not need to take any additional steps to enable security ---
it is on by default and cannot be turned off.


### 9.5 Trust Indicators


The application surfaces trust information in intuitive, non-technical
ways. Verified organizational identities display a checkmark alongside
the sender's domain. Messages from unknown senders appear in a separate
pending section. A simple trust indicator shows whether a contact's
identity has been verified by mutual connections in the user's network.

