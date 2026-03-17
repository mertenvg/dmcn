# DMCN Web Client

Browser-based client for the DMCN messaging network. Private keys never leave the browser — the server stores only encrypted key material.

**Architecture:** Go backend as network proxy + JavaScript frontend doing all crypto client-side via Web Crypto API.

---

## Prerequisites

- Go 1.25 or later
- Node.js 20+ and npm
- A running `dmcn-node` (see root README)
- TLS certificates (self-signed for localhost, autocert for production)

---

## Local Development

### 1. Start a DMCN node

```bash
# From the repository root
make build
./bin/dmcn-node start --listen /ip4/127.0.0.1/tcp/9000
```

Note the multiaddr with peer ID printed at startup (e.g. `/ip4/127.0.0.1/tcp/9000/p2p/12D3Koo...`).

### 2. Generate localhost TLS certificates

The repository includes pre-generated certificates in `certs/`. If you need to regenerate them:

```bash
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
  -keyout certs/localhost.key -out certs/localhost.crt \
  -days 365 -nodes -subj '/CN=localhost' \
  -addext 'subjectAltName=DNS:localhost,IP:127.0.0.1'
```

### 3. Build and run the web client

```bash
# Build frontend + backend (from repository root)
make build-web

# Or build just the Go binary (uses placeholder frontend)
go build -o bin/dmcn-web ./cmd/dmcn-web

# Run in development mode
DMCN_WEB_DEV=true \
DMCN_WEB_NODE=/ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID> \
./bin/dmcn-web
```

The server starts on `https://localhost:8443` by default. Your browser will warn about the self-signed certificate — accept it to proceed.

### 4. Frontend development with hot reload

For iterating on the frontend without rebuilding the Go binary:

```bash
cd cmd/dmcn-web/web
npm install
npm run dev
```

Vite dev server starts on `http://localhost:5173` and proxies `/api` and `/ws` requests to the Go backend at `https://localhost:8443`.

---

## Production Deployment

### Build

```bash
make build-web
```

This produces a single static binary at `bin/dmcn-web` with the frontend embedded via `go:embed`. No separate file serving is needed.

### Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|---|---|---|
| `DMCN_WEB_LISTEN` | `:8443` | HTTP/TLS listen address |
| `DMCN_WEB_DOMAIN` | `dmcn.me` | Domain for TLS autocert |
| `DMCN_WEB_NODE` | *(required)* | Multiaddr of running dmcn-node |
| `DMCN_WEB_DATA_DIR` | `data` | Directory for user files and envelopes |
| `DMCN_WEB_TLS_CERT` | *(empty)* | Path to TLS certificate (overrides autocert) |
| `DMCN_WEB_TLS_KEY` | *(empty)* | Path to TLS private key |
| `DMCN_WEB_DEV` | `false` | Development mode (uses `certs/localhost.*`, relaxed CORS) |
| `DMCN_WEB_POLL_INTERVAL` | `10s` | Relay poll interval for connected users |

### Run

```bash
# With autocert (production — Let's Encrypt via ACME)
DMCN_WEB_NODE=/ip4/<NODE_IP>/tcp/9000/p2p/<PEER_ID> \
DMCN_WEB_DOMAIN=dmcn.me \
DMCN_WEB_LISTEN=:443 \
DMCN_WEB_DATA_DIR=/var/lib/dmcn-web \
./bin/dmcn-web

# With existing certificates
DMCN_WEB_NODE=/ip4/<NODE_IP>/tcp/9000/p2p/<PEER_ID> \
DMCN_WEB_TLS_CERT=/etc/ssl/dmcn.me.crt \
DMCN_WEB_TLS_KEY=/etc/ssl/dmcn.me.key \
DMCN_WEB_LISTEN=:443 \
DMCN_WEB_DATA_DIR=/var/lib/dmcn-web \
./bin/dmcn-web
```

### Data directory

The data directory (default `data/`) contains:

```
data/
  users/              # One JSON file per registered user
    alice@dmcn.me.json
  envelopes/          # Encrypted envelopes pending delivery
    alice@dmcn.me/
      <hash>.bin
```

User files store only public keys and an encrypted payload (AES-256-GCM). The server never has access to plaintext private keys.

Sessions are in-memory only — a server restart logs out all users.

### systemd example

```ini
[Unit]
Description=DMCN Web Client
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/dmcn-web
Environment=DMCN_WEB_NODE=/ip4/127.0.0.1/tcp/9000/p2p/<PEER_ID>
Environment=DMCN_WEB_DOMAIN=dmcn.me
Environment=DMCN_WEB_LISTEN=:443
Environment=DMCN_WEB_DATA_DIR=/var/lib/dmcn-web
Restart=on-failure
User=dmcn
Group=dmcn

[Install]
WantedBy=multi-user.target
```

---

## Security Model

- Private keys are generated in the browser and encrypted with a user-chosen passphrase before being sent to the server
- Key encryption uses HKDF-SHA256 + AES-256-GCM with info string `"dmcn-webkeys-v1"` (domain-separated from the CLI keystore's `"dmcn-keystore-v1"`)
- Login uses challenge-response: the server sends a random nonce, the browser signs it with the user's Ed25519 key to prove key possession
- Message encryption and signing happen entirely in the browser using the Web Crypto API
- The server acts as a network proxy — it relays pre-signed envelopes to the DMCN relay network without access to plaintext content
- TLS is mandatory in production (autocert or configured certificates)
- CSP headers restrict script sources to prevent XSS
- Session tokens are 32 random bytes, not JWTs

---

## API Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/register` | Register new identity |
| `POST` | `/api/v1/login` | Start login (returns challenge) |
| `POST` | `/api/v1/login/verify` | Complete login (verify signed challenge) |

### Authenticated (Bearer token)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/logout` | Logout (invalidate session) |
| `POST` | `/api/v1/messages/send` | Send pre-signed encrypted envelope |
| `GET` | `/api/v1/messages` | List pending encrypted envelopes |
| `POST` | `/api/v1/messages/ack` | Acknowledge envelope delivery |
| `GET` | `/api/v1/identity/lookup?address=...` | Proxy DHT identity lookup |
| `PUT` | `/api/v1/user/payload` | Update encrypted payload |

### WebSocket

| Path | Description |
|---|---|
| `GET /ws` | WebSocket for real-time envelope delivery via FETCH challenge-response |

The WebSocket does not accept the session token in the URL query string (tokens in URLs leak into server logs, browser history, and Referer headers). Instead, the client must send an `authenticate` message as the first frame after connecting:

```json
{"id": "auth", "type": "authenticate", "data": {"token": "<session_token>"}}
```

The server responds with `{"type": "authenticated"}` on success or closes the connection on failure.
