# dmcn-bridge

SMTP-DMCN bridge node that enables bidirectional message exchange between legacy email clients and the DMCN network.

## How It Works

The bridge operates as both an SMTP server (receiving legacy email) and a DMCN node (connected to the mesh network). It translates between the two protocols:

- **Inbound (SMTP to DMCN):** Receives email via SMTP, verifies sender authentication (SPF/DKIM/DMARC), classifies a trust tier, wraps the message in a signed and encrypted DMCN envelope, and stores it on the relay for the DMCN recipient.
- **Outbound (DMCN to SMTP):** Polls the relay for messages addressed to the bridge, decrypts them, delivers via SMTP to the legacy recipient, and returns a signed delivery receipt to the DMCN sender.

## Use Cases

### 1. Development bridge with stubs

Run a bridge using stub authentication and stub SMTP delivery. Messages are accepted on the SMTP port and wrapped into DMCN envelopes. Outbound messages are logged but not actually delivered via SMTP. This is the default mode and is suitable for local development and testing.

```bash
# Terminal 1: Start a DMCN node
dmcn-node start --listen /ip4/127.0.0.1/tcp/9000

# Terminal 2: Start the bridge (note: use full multiaddr with peer ID from node output)
dmcn-bridge start \
    --node /ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID> \
    --smtp-listen 127.0.0.1:2525
```

The stub auth verifier returns configurable default results (pass by default). The stub SMTP deliverer captures messages in memory for test inspection.

### 2. Receive email from legacy senders

Configure your MX records to point to the bridge's SMTP port. When a legacy email arrives:

1. The bridge verifies SPF/DKIM/DMARC (via the pluggable `AuthVerifier` interface)
2. Assigns a trust tier based on authentication results:
   - **Verified Legacy** — DKIM pass + DMARC pass
   - **Unverified Legacy** — no explicit failures but missing checks
   - **Suspicious** — any DKIM or DMARC failure
3. Constructs a `BridgeClassificationRecord` (signed by the bridge) and attaches it to the message
4. Encrypts the message to the DMCN recipient's public key and stores it on the relay

The DMCN recipient sees the trust tier in the classification attachment and can make informed decisions about the message's authenticity.

```bash
# Send a test email to the bridge via SMTP
echo "Subject: Test\n\nHello from legacy email" | \
    sendmail -S 127.0.0.1:2525 -f sender@gmail.com alice@bridge.localhost
```

The bridge maps `alice@bridge.localhost` to `alice@dmcn.localhost` and delivers the encrypted message to Alice's relay mailbox.

### 3. Send email to legacy recipients

When a DMCN user sends a message to a legacy email address, they encrypt it to the bridge's public key. The bridge:

1. Decrypts the envelope
2. Extracts the legacy recipient address
3. Delivers via the outbound SMTP relay
4. Returns a signed `BridgeDeliveryReceipt` to the sender confirming delivery success or failure

```bash
# Alice sends a message to a legacy address through the bridge
dmcn-node message send \
    --from alice@dmcn.localhost \
    --to external@gmail.com \
    --body "Hello from DMCN" \
    --node /ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID> \
    --keystore alice.json --passphrase secret
```

The bridge picks up the message on its next poll cycle, delivers it via SMTP, and stores a delivery receipt for Alice to fetch.

### 4. Full round-trip demonstration

This demonstrates the complete PRD Section 6.4 scenario:

```bash
# Terminal 1: Start a DMCN node
dmcn-node start --listen /ip4/127.0.0.1/tcp/9000

# Terminal 2: Set up identities and start the bridge
NODE=/ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID>

# Generate and register Alice as a DMCN user
dmcn-node identity generate --address alice@dmcn.localhost --keystore alice.json --passphrase secret
dmcn-node identity register --address alice@dmcn.localhost --node $NODE --keystore alice.json --passphrase secret

# Start the bridge (it generates its own identity and registers with BridgeCapability)
dmcn-bridge start \
    --node $NODE \
    --smtp-listen 127.0.0.1:2525 \
    --bridge-domain bridge.localhost \
    --dmcn-domain dmcn.localhost

# Terminal 3: Simulate inbound legacy email
echo "Hello Alice from legacy land" | \
    sendmail -S 127.0.0.1:2525 -f external@gmail.com alice@bridge.localhost

# Terminal 4: Alice fetches her message (includes BridgeClassificationRecord attachment)
dmcn-node message fetch \
    --address alice@dmcn.localhost --node $NODE \
    --keystore alice.json --passphrase secret
```

## Configuration Reference

| Flag | Default | Description |
|---|---|---|
| `--node` | *(required)* | Multiaddr of a running `dmcn-node` (must include `/p2p/<peer-id>`) |
| `--smtp-listen` | `0.0.0.0:2525` | Address for the SMTP server to listen on |
| `--libp2p-listen` | `/ip4/127.0.0.1/tcp/0` | libp2p listen multiaddr for the bridge's own node |
| `--bridge-domain` | `bridge.localhost` | Email domain the bridge accepts mail for |
| `--dmcn-domain` | `dmcn.localhost` | DMCN address domain to map bridge addresses to |
| `--bridge-address` | `bridge@bridge.localhost` | Bridge's own DMCN identity address |
| `--keystore` | `bridge-keystore.enc` | Path to encrypted keystore file for bridge keys |
| `--passphrase` | `dmcn-bridge-dev` | Keystore encryption passphrase |
| `--poll-interval` | `5s` | How often to poll the relay for outbound messages |
| `--smtp-relay` | *(none, stub mode)* | Outbound SMTP relay `host:port` for real delivery |

## Address Mapping

The bridge translates between SMTP and DMCN address spaces using the configured domains:

| Direction | From | To |
|---|---|---|
| Inbound | `alice@bridge.localhost` | `alice@dmcn.localhost` |
| Outbound (From) | `alice@dmcn.localhost` | `alice@bridge.localhost` |
| Legacy passthrough | `external@gmail.com` | `external@gmail.com` (unchanged) |

## Trust Classification

Inbound legacy messages are classified based on email authentication results:

| Trust Tier | Criteria | Meaning |
|---|---|---|
| Verified Legacy | DKIM pass + DMARC pass | Sender domain authentication succeeded |
| Unverified Legacy | No explicit failures | Authentication checks incomplete or neutral |
| Suspicious | Any DKIM or DMARC failure | Sender authentication failed — potential forgery |

The classification is cryptographically signed by the bridge and attached to the DMCN message as a `BridgeClassificationRecord`. Recipients can verify the bridge's signature and use the trust tier to inform their own trust decisions.

## PoC Limitations

- **Auth verification is stubbed.** The `AuthVerifier` interface is pluggable but only a stub implementation is provided. Production use requires implementing SPF/DKIM/DMARC checks against real DNS.
- **SMTP delivery is stubbed.** The `SMTPDeliverer` interface is pluggable but defaults to a no-op stub. Pass `--smtp-relay` to configure a real outbound relay (not yet implemented in the CLI).
- **In-memory relay storage.** Messages are lost on restart. Persistent storage is post-PoC.
- **No TLS on the SMTP port.** The SMTP server does not support STARTTLS. Do not expose it to the public internet.
