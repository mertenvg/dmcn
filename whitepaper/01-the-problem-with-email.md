## 1. The Problem with Email


### 1.1 A Protocol Designed for a Different World


The Simple Mail Transfer Protocol (SMTP), first defined in RFC 821 in
1982, was engineered for a network of a few hundred trusted nodes ---
universities, government research labs, and military institutions. In
that context, identity verification was unnecessary. Everyone on the
network was known. The protocol's defining characteristic — openness
--- was a feature, not a vulnerability.

Today, SMTP underpins communication for over 4 billion email users
globally. The network it was designed for no longer exists. The trust
assumptions it was built on are not merely strained — they are
completely inverted. The openness that made email powerful is precisely
what makes it exploitable.


### 1.2 The Scale of the Spam Problem


Spam is not a minor nuisance. It is the dominant form of email
communication on the planet by volume. Industry estimates consistently
place spam at between 45% and 85% of all global email traffic, with some
peak periods substantially higher. On any given day, hundreds of
billions of spam messages are transmitted across global email
infrastructure.

The consequences extend well beyond annoyance. Email-based phishing is
among the most costly forms of cybercrime. Business Email Compromise
(BEC) — a form of fraud in which attackers impersonate executives or
trusted partners — costs organizations billions of dollars annually.
The FBI's Internet Crime Complaint Center consistently ranks BEC among
the highest-impact cybercrime categories by financial loss.


> **Scale of the Problem**
> *Spam accounts for the majority of all email ever sent. Email fraud
> costs the global economy billions annually. These are not edge cases
> — they are the normal operating conditions of the current system.*


### 1.3 The Structural Root Cause


The spam problem is an identity problem. SMTP provides no mechanism for
a receiving server to verify that the sending server is who it claims to
be, and no mechanism for a recipient to verify that a message actually
came from the stated sender. Sending an email that appears to come from
any address — including your bank, your employer, or a government
agency — requires no special access, no credentials, and no technical
sophistication beyond a basic SMTP client.

This means that spam and phishing are not aberrations that a better
filter can eliminate. They are rational economic behaviors enabled by
the protocol itself. As long as sending a message claiming to be from
any sender costs essentially nothing and carries no verifiable identity,
the conditions for spam will persist regardless of how sophisticated
filtering becomes.


### 1.4 Existing Mitigations and Their Limitations


The email ecosystem has accumulated several layers of mitigation over
the decades, none of which address the root cause:

- Sender Policy Framework (SPF) — allows domain owners to specify which IP addresses are authorized to send email on their behalf. Widely adopted but easily circumvented through compromised authorized servers, and provides no per-message signing.

- DomainKeys Identified Mail (DKIM) — adds a cryptographic signature to outgoing messages, allowing receivers to verify that content was not altered in transit. Addresses integrity but not spam; a spammer controlling their own domain can produce valid DKIM signatures.

- Domain-based Message Authentication, Reporting and Conformance (DMARC) — a policy layer on top of SPF and DKIM. Adoption is inconsistent and enforcement is frequently weak.

- AI-based spam filtering — the approach used by major providers like Google and Microsoft. Highly effective at classifying known spam patterns, but reactive by nature, computationally expensive, and produces significant false positive rates that affect legitimate communication.

Each of these mitigations is a layer of additional complexity applied to
a fundamentally trust-less protocol. They reduce spam volumes in
practice, but they cannot eliminate spam in principle, because they do
not address sender identity at the protocol level.



---
