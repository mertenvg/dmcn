# dmcn-node

Combined DMCN development node that runs a DHT identity registry and message relay service in a single process.

## Use Cases

### 1. Run a standalone node

Start a node that other peers can connect to for identity registration, lookup, and message relay.

```bash
dmcn-node start --listen /ip4/0.0.0.0/tcp/7400
```

The node prints its multiaddr (including peer ID) on startup. Other nodes and bridges use this address to connect.

### 2. Join an existing network

Connect to one or more bootstrap peers to join an existing DHT network. Comma-separate multiple peers.

```bash
dmcn-node start \
    --listen /ip4/0.0.0.0/tcp/7400 \
    --bootstrap /ip4/192.168.1.10/tcp/7400/p2p/12D3KooW...
```

### 3. Generate an identity

Create a new Ed25519 + X25519 key pair and store it in an encrypted keystore file.

```bash
dmcn-node identity generate \
    --address alice@example.com \
    --keystore keys.json \
    --passphrase "my secret passphrase"
```

### 4. Register an identity on the DHT

Publish an identity record to the distributed hash table so other users can look up your public keys by address.

```bash
dmcn-node identity register \
    --address alice@example.com \
    --node /ip4/127.0.0.1/tcp/7400/p2p/12D3KooW... \
    --keystore keys.json \
    --passphrase "my secret passphrase"
```

### 5. Look up an identity

Retrieve a registered identity record by address and verify its self-signature.

```bash
dmcn-node identity lookup \
    --address bob@example.com \
    --node /ip4/127.0.0.1/tcp/7400/p2p/12D3KooW...
```

### 6. Send a message

Compose, sign, encrypt, and deliver a message to a registered recipient via the relay.

```bash
dmcn-node message send \
    --from alice@example.com \
    --to bob@example.com \
    --body "Hello from DMCN" \
    --node /ip4/127.0.0.1/tcp/7400/p2p/12D3KooW... \
    --keystore keys.json \
    --passphrase "my secret passphrase"
```

### 7. Fetch pending messages

Retrieve, decrypt, and display messages waiting on the relay for your address. Uses challenge-response authentication to prove key ownership.

```bash
dmcn-node message fetch \
    --address alice@example.com \
    --node /ip4/127.0.0.1/tcp/7400/p2p/12D3KooW... \
    --keystore keys.json \
    --passphrase "my secret passphrase"
```

## Full End-to-End Example

```bash
# Terminal 1: Start a node
dmcn-node start --listen /ip4/127.0.0.1/tcp/9000

# Terminal 2: Generate and register two identities, then send a message
NODE=/ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID>

dmcn-node identity generate --address alice@localhost --keystore alice.json --passphrase secret
dmcn-node identity register --address alice@localhost --node $NODE --keystore alice.json --passphrase secret

dmcn-node identity generate --address bob@localhost --keystore bob.json --passphrase secret
dmcn-node identity register --address bob@localhost --node $NODE --keystore bob.json --passphrase secret

dmcn-node message send \
    --from alice@localhost --to bob@localhost --body "Hello Bob!" \
    --node $NODE --keystore alice.json --passphrase secret

dmcn-node message fetch \
    --address bob@localhost --node $NODE --keystore bob.json --passphrase secret
```

## Configuration Reference

### `start` command

| Flag | Default | Description |
|---|---|---|
| `--listen` | `/ip4/0.0.0.0/tcp/7400` | libp2p listen multiaddr |
| `--keystore` | `keystore.enc` | Path to encrypted keystore file |
| `--passphrase` | `dmcn-dev` | Keystore encryption passphrase |
| `--bootstrap` | *(none)* | Comma-separated bootstrap peer multiaddrs |

### `identity generate` command

| Flag | Default | Description |
|---|---|---|
| `--address` | *(required)* | DMCN address (local@domain) |
| `--keystore` | `keystore.enc` | Path to encrypted keystore file |
| `--passphrase` | `dmcn-dev` | Keystore encryption passphrase |

### `identity register` command

| Flag | Default | Description |
|---|---|---|
| `--address` | *(required)* | DMCN address to register |
| `--node` | *(required)* | Multiaddr of a running DMCN node |
| `--keystore` | `keystore.enc` | Path to encrypted keystore file |
| `--passphrase` | `dmcn-dev` | Keystore encryption passphrase |

### `identity lookup` command

| Flag | Default | Description |
|---|---|---|
| `--address` | *(required)* | DMCN address to look up |
| `--node` | *(required)* | Multiaddr of a running DMCN node |

### `message send` command

| Flag | Default | Description |
|---|---|---|
| `--from` | *(required)* | Sender's DMCN address |
| `--to` | *(required)* | Recipient's DMCN address |
| `--body` | *(required)* | Message body text |
| `--node` | *(required)* | Multiaddr of a running DMCN node |
| `--keystore` | `keystore.enc` | Path to encrypted keystore file |
| `--passphrase` | `dmcn-dev` | Keystore encryption passphrase |

### `message fetch` command

| Flag | Default | Description |
|---|---|---|
| `--address` | *(required)* | DMCN address to fetch messages for |
| `--node` | *(required)* | Multiaddr of a running DMCN node |
| `--keystore` | `keystore.enc` | Path to encrypted keystore file |
| `--passphrase` | `dmcn-dev` | Keystore encryption passphrase |
