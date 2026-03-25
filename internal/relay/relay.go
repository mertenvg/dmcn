package relay

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/mertenvg/logr/v2"
	"google.golang.org/protobuf/proto"

	"github.com/mertenvg/dmcn/internal/core/crypto"
	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/proto/dmcnpb"
	"github.com/mertenvg/dmcn/internal/registry"
)

const (
	// ProtocolID is the libp2p protocol identifier for the relay service.
	ProtocolID = protocol.ID("/dmcn/relay/1.0.0")
	// OrgPeersProtocolID is the libp2p protocol identifier for org peer discovery.
	OrgPeersProtocolID = protocol.ID("/dmcn/org/1.0.0")
	// maxMessageSize is the maximum size of a single protocol message (4 MB).
	maxMessageSize = 4 * 1024 * 1024
	// defaultRateLimit is the PoC rate limit (100 STORE ops/hr/identity).
	defaultRateLimit = 100
)

var (
	// ErrUnregisteredSender is returned when a STORE sender is not in the registry.
	ErrUnregisteredSender = errors.New("relay: sender identity not registered")
	// ErrRateLimited is returned when a sender exceeds the rate limit.
	ErrRateLimited = errors.New("relay: rate limit exceeded")
	// ErrAuthFailed is returned when FETCH authentication fails.
	ErrAuthFailed = errors.New("relay: authentication failed")
)

// LookupFunc looks up an identity in the registry by address.
// This abstraction allows testing without a full DHT.
type LookupFunc func(ctx context.Context, address string) (*identity.IdentityRecord, error)

// Relay implements the DMCN relay protocol over libp2p streams.
type Relay struct {
	host      host.Host
	lookup    LookupFunc
	store     *MessageStore
	limiter   *RateLimiter
	log       logr.Logger
	startTime time.Time
	version   string
	orgPeers  []string

	mu      sync.Mutex
	started bool
}

// New creates a new Relay service.
func New(h host.Host, lookup LookupFunc, opts ...Option) *Relay {
	cfg := &relayOptions{
		rateLimit: defaultRateLimit,
		version:   "dmcn-node/0.1.0",
	}
	for _, o := range opts {
		o(cfg)
	}

	return &Relay{
		host:     h,
		lookup:   lookup,
		store:    NewMessageStore(),
		limiter:  NewRateLimiter(cfg.rateLimit),
		log:      logr.With(logr.M("component", "relay")),
		version:  cfg.version,
		orgPeers: cfg.orgPeers,
	}
}

type relayOptions struct {
	rateLimit int
	version   string
	orgPeers  []string
}

// Option configures a Relay.
type Option func(*relayOptions)

// WithRateLimit sets the maximum STORE operations per hour per sender.
func WithRateLimit(maxPerHour int) Option {
	return func(o *relayOptions) {
		o.rateLimit = maxPerHour
	}
}

// WithOrgPeers sets the organizational peers for this relay node.
func WithOrgPeers(peers []string) Option {
	return func(o *relayOptions) {
		o.orgPeers = peers
	}
}

// Start registers the stream handler and begins serving relay operations.
func (r *Relay) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started {
		return
	}
	r.started = true
	r.startTime = time.Now()
	r.host.SetStreamHandler(ProtocolID, r.handleStream)
	r.host.SetStreamHandler(OrgPeersProtocolID, r.handleOrgPeers)
	r.log.Info("relay started")
}

// Stop removes the stream handler.
func (r *Relay) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.started {
		return
	}
	r.started = false
	r.host.RemoveStreamHandler(ProtocolID)
	r.host.RemoveStreamHandler(OrgPeersProtocolID)
	r.log.Info("relay stopped")
}

// Store returns the underlying message store for direct access in tests.
func (r *Relay) Store() *MessageStore {
	return r.store
}

// handleStream processes an incoming relay protocol stream.
func (r *Relay) handleStream(s network.Stream) {
	defer s.Close()

	req, err := readRequest(s)
	if err != nil {
		return
	}

	var resp *dmcnpb.RelayResponse

	switch {
	case req.GetStore() != nil:
		resp = r.handleStore(s.Conn().RemotePeer(), req.GetStore())
	case req.GetFetchInit() != nil:
		r.handleFetch(s, req.GetFetchInit())
		return // handleFetch writes its own responses
	case req.GetAck() != nil:
		resp = r.handleAck(req.GetAck())
	case req.GetPing() != nil:
		resp = r.handlePing()
	default:
		resp = errorResponse("INVALID_REQUEST", "unknown request type")
	}

	writeResponse(s, resp)
}

// handleStore processes a STORE request.
func (r *Relay) handleStore(_ peer.ID, req *dmcnpb.StoreRequest) *dmcnpb.RelayResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Rate limit check
	if !r.limiter.Allow(req.SenderAddress) {
		r.log.Warnf("STORE rate limited for sender %s", req.SenderAddress)
		return errorResponse("RATE_LIMITED", ErrRateLimited.Error())
	}

	// 2. Verify sender is registered
	senderRec, err := r.lookup(ctx, req.SenderAddress)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			r.log.Warnf("STORE rejected: unregistered sender %s", req.SenderAddress)
			return errorResponse("UNREGISTERED_SENDER", ErrUnregisteredSender.Error())
		}
		return errorResponse("LOOKUP_FAILED", fmt.Sprintf("sender lookup: %v", err))
	}

	// 3. Deserialize envelope
	env, err := message.EncryptedEnvelopeFromProto(req.Envelope)
	if err != nil {
		return errorResponse("INVALID_ENVELOPE", fmt.Sprintf("invalid envelope: %v", err))
	}

	// 4. Compute envelope hash and verify sender signature
	envBytes, err := proto.Marshal(req.Envelope)
	if err != nil {
		return errorResponse("INTERNAL_ERROR", "failed to marshal envelope")
	}
	envHash := crypto.SHA256Hash(envBytes)

	if err := crypto.Verify(senderRec.Ed25519Public, envHash[:], req.SenderSignature); err != nil {
		r.log.Warnf("STORE rejected: invalid signature from %s", req.SenderAddress)
		return errorResponse("INVALID_SIGNATURE", "sender signature verification failed")
	}

	// 5. Determine recipient address from envelope
	// The recipient address is not directly in the envelope (it's encrypted).
	// For the PoC, we extract it from the PlaintextMessage sender_address field
	// of the StoreRequest. We need a way to know who to store it for.
	// The relay stores by recipient, but the envelope doesn't expose the recipient.
	// We'll look up who the recipients are by their X25519 public keys in the
	// recipient records. For PoC simplicity, we store for all recipient public
	// keys found in the envelope.
	for _, rec := range env.Recipients {
		// Store indexed by hex-encoded recipient X25519 public key
		addr := fmt.Sprintf("%x", rec.RecipientXPub[:])
		r.store.Store(addr, env, envHash)
	}

	r.log.Debugf("STORE accepted from %s, hash %x, %d recipients", req.SenderAddress, envHash[:8], len(env.Recipients))

	return &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_Store{
			Store: &dmcnpb.StoreResponse{
				EnvelopeHash: envHash[:],
			},
		},
	}
}

// handleFetch processes a FETCH request with challenge-response auth.
func (r *Relay) handleFetch(s network.Stream, init *dmcnpb.FetchInit) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Lookup recipient in registry
	rec, err := r.lookup(ctx, init.Address)
	if err != nil {
		writeResponse(s, errorResponse("LOOKUP_FAILED", fmt.Sprintf("identity not found: %v", err)))
		return
	}

	// 2. Generate challenge nonce
	nonce, err := crypto.RandomBytes(32)
	if err != nil {
		writeResponse(s, errorResponse("INTERNAL_ERROR", "failed to generate challenge"))
		return
	}

	// 3. Send challenge
	challenge := &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_FetchChallenge{
			FetchChallenge: &dmcnpb.FetchChallenge{
				Nonce: nonce,
			},
		},
	}
	if err := writeResponse(s, challenge); err != nil {
		return
	}

	// 4. Read proof
	proofReq, err := readRequest(s)
	if err != nil {
		return
	}
	proof := proofReq.GetFetchProof()
	if proof == nil {
		writeResponse(s, errorResponse("INVALID_REQUEST", "expected fetch proof"))
		return
	}

	// 5. Verify signature
	if err := crypto.Verify(rec.Ed25519Public, nonce, proof.Signature); err != nil {
		r.log.Warnf("FETCH auth failed for %s", init.Address)
		writeResponse(s, errorResponse("AUTH_FAILED", ErrAuthFailed.Error()))
		return
	}

	// 6. Return pending envelopes
	// Look up by X25519 public key (hex-encoded)
	addr := fmt.Sprintf("%x", rec.X25519Public[:])
	envs, hashes := r.store.Fetch(addr)

	pbEnvs := make([]*dmcnpb.EncryptedEnvelope, len(envs))
	pbHashes := make([][]byte, len(hashes))
	for i, env := range envs {
		pbEnvs[i] = env.ToProto()
		hash := hashes[i]
		pbHashes[i] = hash[:]
	}

	r.log.Debugf("FETCH returning %d envelope(s) for %s", len(envs), init.Address)

	writeResponse(s, &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_Fetch{
			Fetch: &dmcnpb.FetchResponse{
				Envelopes:      pbEnvs,
				EnvelopeHashes: pbHashes,
			},
		},
	})
}

// handleAck processes an ACK request.
func (r *Relay) handleAck(req *dmcnpb.AckRequest) *dmcnpb.RelayResponse {
	var hash [32]byte
	copy(hash[:], req.EnvelopeHash)

	if err := r.store.Ack(hash); err != nil {
		return errorResponse("NOT_FOUND", err.Error())
	}

	return &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_Ack{
			Ack: &dmcnpb.AckResponse{Success: true},
		},
	}
}

// handlePing processes a PING request.
func (r *Relay) handlePing() *dmcnpb.RelayResponse {
	uptime := time.Since(r.startTime)
	return &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_Ping{
			Ping: &dmcnpb.PingResponse{
				Version:         r.version,
				UptimeSeconds:   int64(uptime.Seconds()),
				StoredEnvelopes: r.store.Count(),
			},
		},
	}
}

// --- Client methods ---

// ClientStore sends a STORE request to a remote relay node.
func (r *Relay) ClientStore(ctx context.Context, peerID peer.ID, senderKP *identity.IdentityKeyPair, env *message.EncryptedEnvelope) ([32]byte, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: open stream: %w", err)
	}
	defer s.Close()

	// Serialize envelope for hashing
	envProto := env.ToProto()
	envBytes, err := proto.Marshal(envProto)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: marshal envelope: %w", err)
	}
	envHash := crypto.SHA256Hash(envBytes)

	// Sign envelope hash
	sig, err := crypto.Sign(senderKP.Ed25519Private, envHash[:])
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: sign: %w", err)
	}

	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_Store{
			Store: &dmcnpb.StoreRequest{
				SenderAddress:   "", // filled by caller context
				SenderSignature: sig,
				Envelope:        envProto,
			},
		},
	}

	if err := writeRequest(s, req); err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: write: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: read response: %w", err)
	}

	if errResp := resp.GetError(); errResp != nil {
		switch errResp.Code {
		case "UNREGISTERED_SENDER":
			return [32]byte{}, ErrUnregisteredSender
		case "RATE_LIMITED":
			return [32]byte{}, ErrRateLimited
		default:
			return [32]byte{}, fmt.Errorf("relay: store: %s: %s", errResp.Code, errResp.Message)
		}
	}

	storeResp := resp.GetStore()
	if storeResp == nil {
		return [32]byte{}, errors.New("relay: store: unexpected response type")
	}

	copy(envHash[:], storeResp.EnvelopeHash)
	return envHash, nil
}

// ClientStoreWithAddress sends a STORE request with explicit sender address.
func (r *Relay) ClientStoreWithAddress(ctx context.Context, peerID peer.ID, senderAddr string, senderKP *identity.IdentityKeyPair, env *message.EncryptedEnvelope) ([32]byte, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: open stream: %w", err)
	}
	defer s.Close()

	envProto := env.ToProto()
	envBytes, err := proto.Marshal(envProto)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: marshal envelope: %w", err)
	}
	envHash := crypto.SHA256Hash(envBytes)

	sig, err := crypto.Sign(senderKP.Ed25519Private, envHash[:])
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: sign: %w", err)
	}

	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_Store{
			Store: &dmcnpb.StoreRequest{
				SenderAddress:   senderAddr,
				SenderSignature: sig,
				Envelope:        envProto,
			},
		},
	}

	if err := writeRequest(s, req); err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: write: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: read response: %w", err)
	}

	if errResp := resp.GetError(); errResp != nil {
		switch errResp.Code {
		case "UNREGISTERED_SENDER":
			return [32]byte{}, ErrUnregisteredSender
		case "RATE_LIMITED":
			return [32]byte{}, ErrRateLimited
		default:
			return [32]byte{}, fmt.Errorf("relay: store: %s: %s", errResp.Code, errResp.Message)
		}
	}

	storeResp := resp.GetStore()
	if storeResp == nil {
		return [32]byte{}, errors.New("relay: store: unexpected response type")
	}

	copy(envHash[:], storeResp.EnvelopeHash)
	return envHash, nil
}

// ClientFetch authenticates to a remote relay and retrieves pending envelopes.
func (r *Relay) ClientFetch(ctx context.Context, peerID peer.ID, kp *identity.IdentityKeyPair, address string) ([]*message.EncryptedEnvelope, [][32]byte, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: open stream: %w", err)
	}
	defer s.Close()

	// 1. Send FetchInit
	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_FetchInit{
			FetchInit: &dmcnpb.FetchInit{Address: address},
		},
	}
	if err := writeRequest(s, req); err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: write init: %w", err)
	}

	// 2. Read challenge
	resp, err := readResponse(s)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: read challenge: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		return nil, nil, fmt.Errorf("relay: fetch: %s: %s", errResp.Code, errResp.Message)
	}
	challenge := resp.GetFetchChallenge()
	if challenge == nil {
		return nil, nil, errors.New("relay: fetch: expected challenge response")
	}

	// 3. Sign nonce and send proof
	sig, err := crypto.Sign(kp.Ed25519Private, challenge.Nonce)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: sign challenge: %w", err)
	}

	proofReq := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_FetchProof{
			FetchProof: &dmcnpb.FetchProof{
				Address:   address,
				Nonce:     challenge.Nonce,
				Signature: sig,
			},
		},
	}
	if err := writeRequest(s, proofReq); err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: write proof: %w", err)
	}

	// 4. Read envelopes
	resp, err = readResponse(s)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: read envelopes: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		if errResp.Code == "AUTH_FAILED" {
			return nil, nil, ErrAuthFailed
		}
		return nil, nil, fmt.Errorf("relay: fetch: %s: %s", errResp.Code, errResp.Message)
	}

	fetchResp := resp.GetFetch()
	if fetchResp == nil {
		return nil, nil, errors.New("relay: fetch: unexpected response type")
	}

	envs := make([]*message.EncryptedEnvelope, len(fetchResp.Envelopes))
	hashes := make([][32]byte, len(fetchResp.EnvelopeHashes))
	for i, pb := range fetchResp.Envelopes {
		env, err := message.EncryptedEnvelopeFromProto(pb)
		if err != nil {
			return nil, nil, fmt.Errorf("relay: fetch: envelope %d: %w", i, err)
		}
		envs[i] = env
		if i < len(fetchResp.EnvelopeHashes) {
			copy(hashes[i][:], fetchResp.EnvelopeHashes[i])
		}
	}

	return envs, hashes, nil
}

// ClientStorePreSigned sends a STORE request with a pre-computed signature.
// This is used by the web client where the browser signs the envelope hash
// and the server relays it without having access to private keys.
func (r *Relay) ClientStorePreSigned(ctx context.Context, peerID peer.ID, senderAddr string, signature []byte, env *message.EncryptedEnvelope) ([32]byte, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: open stream: %w", err)
	}
	defer s.Close()

	envProto := env.ToProto()

	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_Store{
			Store: &dmcnpb.StoreRequest{
				SenderAddress:   senderAddr,
				SenderSignature: signature,
				Envelope:        envProto,
			},
		},
	}

	if err := writeRequest(s, req); err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: write: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return [32]byte{}, fmt.Errorf("relay: store: read response: %w", err)
	}

	if errResp := resp.GetError(); errResp != nil {
		switch errResp.Code {
		case "UNREGISTERED_SENDER":
			return [32]byte{}, ErrUnregisteredSender
		case "RATE_LIMITED":
			return [32]byte{}, ErrRateLimited
		default:
			return [32]byte{}, fmt.Errorf("relay: store: %s: %s", errResp.Code, errResp.Message)
		}
	}

	storeResp := resp.GetStore()
	if storeResp == nil {
		return [32]byte{}, errors.New("relay: store: unexpected response type")
	}

	var envHash [32]byte
	copy(envHash[:], storeResp.EnvelopeHash)
	return envHash, nil
}

// ClientFetchChallenge sends a FetchInit and returns the challenge nonce and
// the open stream. The caller must complete the exchange by calling
// ClientFetchComplete with the signed nonce.
func (r *Relay) ClientFetchChallenge(ctx context.Context, peerID peer.ID, address string) ([]byte, network.Stream, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: open stream: %w", err)
	}

	// Send FetchInit
	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_FetchInit{
			FetchInit: &dmcnpb.FetchInit{Address: address},
		},
	}
	if err := writeRequest(s, req); err != nil {
		s.Close()
		return nil, nil, fmt.Errorf("relay: fetch: write init: %w", err)
	}

	// Read challenge
	resp, err := readResponse(s)
	if err != nil {
		s.Close()
		return nil, nil, fmt.Errorf("relay: fetch: read challenge: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		s.Close()
		return nil, nil, fmt.Errorf("relay: fetch: %s: %s", errResp.Code, errResp.Message)
	}
	challenge := resp.GetFetchChallenge()
	if challenge == nil {
		s.Close()
		return nil, nil, errors.New("relay: fetch: expected challenge response")
	}

	return challenge.Nonce, s, nil
}

// ClientFetchComplete sends the signed proof on an open stream (from
// ClientFetchChallenge) and returns the envelopes. The stream is closed
// when this method returns.
func (r *Relay) ClientFetchComplete(s network.Stream, address string, nonce, signature []byte) ([]*message.EncryptedEnvelope, [][32]byte, error) {
	defer s.Close()

	proofReq := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_FetchProof{
			FetchProof: &dmcnpb.FetchProof{
				Address:   address,
				Nonce:     nonce,
				Signature: signature,
			},
		},
	}
	if err := writeRequest(s, proofReq); err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: write proof: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return nil, nil, fmt.Errorf("relay: fetch: read envelopes: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		if errResp.Code == "AUTH_FAILED" {
			return nil, nil, ErrAuthFailed
		}
		return nil, nil, fmt.Errorf("relay: fetch: %s: %s", errResp.Code, errResp.Message)
	}

	fetchResp := resp.GetFetch()
	if fetchResp == nil {
		return nil, nil, errors.New("relay: fetch: unexpected response type")
	}

	envs := make([]*message.EncryptedEnvelope, len(fetchResp.Envelopes))
	hashes := make([][32]byte, len(fetchResp.EnvelopeHashes))
	for i, pb := range fetchResp.Envelopes {
		env, err := message.EncryptedEnvelopeFromProto(pb)
		if err != nil {
			return nil, nil, fmt.Errorf("relay: fetch: envelope %d: %w", i, err)
		}
		envs[i] = env
		if i < len(fetchResp.EnvelopeHashes) {
			copy(hashes[i][:], fetchResp.EnvelopeHashes[i])
		}
	}

	return envs, hashes, nil
}

// ClientAck sends an ACK for a delivered envelope.
func (r *Relay) ClientAck(ctx context.Context, peerID peer.ID, envelopeHash [32]byte) error {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return fmt.Errorf("relay: ack: open stream: %w", err)
	}
	defer s.Close()

	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_Ack{
			Ack: &dmcnpb.AckRequest{
				EnvelopeHash: envelopeHash[:],
			},
		},
	}
	if err := writeRequest(s, req); err != nil {
		return fmt.Errorf("relay: ack: write: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return fmt.Errorf("relay: ack: read: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		return fmt.Errorf("relay: ack: %s: %s", errResp.Code, errResp.Message)
	}

	return nil
}

// ClientPing sends a PING to a remote relay.
func (r *Relay) ClientPing(ctx context.Context, peerID peer.ID) (*dmcnpb.PingResponse, error) {
	s, err := r.host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return nil, fmt.Errorf("relay: ping: open stream: %w", err)
	}
	defer s.Close()

	req := &dmcnpb.RelayRequest{
		Request: &dmcnpb.RelayRequest_Ping{
			Ping: &dmcnpb.PingRequest{},
		},
	}
	if err := writeRequest(s, req); err != nil {
		return nil, fmt.Errorf("relay: ping: write: %w", err)
	}

	resp, err := readResponse(s)
	if err != nil {
		return nil, fmt.Errorf("relay: ping: read: %w", err)
	}
	if errResp := resp.GetError(); errResp != nil {
		return nil, fmt.Errorf("relay: ping: %s: %s", errResp.Code, errResp.Message)
	}

	return resp.GetPing(), nil
}

// --- Org peers protocol ---

// orgPeersResponse is the JSON response for the org peers protocol.
type orgPeersResponse struct {
	Peers []string `json:"peers"`
}

// handleOrgPeers responds to org peer discovery requests with the configured
// org peers list.
func (r *Relay) handleOrgPeers(s network.Stream) {
	defer s.Close()

	resp := orgPeersResponse{Peers: r.orgPeers}
	if resp.Peers == nil {
		resp.Peers = []string{}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))
	s.Write(lenBuf[:])
	s.Write(data)
}

// ClientOrgPeers queries a remote peer for its organizational peers.
func (r *Relay) ClientOrgPeers(ctx context.Context, peerID peer.ID) ([]string, error) {
	s, err := r.host.NewStream(ctx, peerID, OrgPeersProtocolID)
	if err != nil {
		return nil, fmt.Errorf("relay: org-peers: open stream: %w", err)
	}
	defer s.Close()

	// Read length-prefixed JSON response.
	var lenBuf [4]byte
	if _, err := io.ReadFull(s, lenBuf[:]); err != nil {
		return nil, fmt.Errorf("relay: org-peers: read length: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	if length > maxMessageSize {
		return nil, errors.New("relay: org-peers: response too large")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(s, data); err != nil {
		return nil, fmt.Errorf("relay: org-peers: read data: %w", err)
	}

	var resp orgPeersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("relay: org-peers: unmarshal: %w", err)
	}

	return resp.Peers, nil
}

// OrgPeers returns the configured organizational peers.
func (r *Relay) OrgPeers() []string {
	return r.orgPeers
}

// --- Wire protocol helpers ---
// Messages are length-prefixed: [4-byte big-endian length][protobuf data]

func writeRequest(w io.Writer, req *dmcnpb.RelayRequest) error {
	return writeMessage(w, req)
}

func writeResponse(w io.Writer, resp *dmcnpb.RelayResponse) error {
	return writeMessage(w, resp)
}

func readRequest(r io.Reader) (*dmcnpb.RelayRequest, error) {
	req := &dmcnpb.RelayRequest{}
	if err := readMessage(r, req); err != nil {
		return nil, err
	}
	return req, nil
}

func readResponse(r io.Reader) (*dmcnpb.RelayResponse, error) {
	resp := &dmcnpb.RelayResponse{}
	if err := readMessage(r, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func writeMessage(w io.Writer, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if len(data) > maxMessageSize {
		return errors.New("message too large")
	}

	// Write 4-byte big-endian length prefix
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return fmt.Errorf("write length: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write data: %w", err)
	}
	return nil
}

func readMessage(r io.Reader, msg proto.Message) error {
	// Read 4-byte big-endian length prefix
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return fmt.Errorf("read length: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	if length > maxMessageSize {
		return errors.New("message too large")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return fmt.Errorf("read data: %w", err)
	}
	if err := proto.Unmarshal(data, msg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func errorResponse(code, msg string) *dmcnpb.RelayResponse {
	return &dmcnpb.RelayResponse{
		Response: &dmcnpb.RelayResponse_Error{
			Error: &dmcnpb.ErrorResponse{
				Code:    code,
				Message: msg,
			},
		},
	}
}
