package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/cmd/dmcn-web/internal/store"
	"github.com/mertenvg/dmcn/internal/core/message"
)

// PendingFetch holds a relay stream awaiting the client's signature proof.
type PendingFetch struct {
	Nonce     []byte
	Stream    network.Stream
	ExpiresAt time.Time
}

// clientConn wraps a WebSocket connection with its owner address and
// cancellation function for background goroutines.
type clientConn struct {
	conn    *websocket.Conn
	address string
	cancel  context.CancelFunc
}

// SessionValidator validates a session token and returns the associated address.
type SessionValidator interface {
	Validate(token string) (string, error)
}

// RelayProxy abstracts the two-phase relay FETCH protocol. FetchChallenge
// initiates a FETCH request and returns the challenge nonce and the open
// libp2p stream. FetchComplete sends the signed proof over the same stream
// and returns the resulting envelopes with their content hashes.
type RelayProxy interface {
	FetchChallenge(ctx context.Context, address string) (nonce []byte, stream network.Stream, err error)
	FetchComplete(stream network.Stream, address string, nonce, signature []byte) ([]*message.EncryptedEnvelope, [][32]byte, error)
}

// ConnManager upgrades HTTP connections to WebSocket and manages the
// lifecycle of connected clients including periodic message polling.
type ConnManager struct {
	mu             sync.RWMutex
	conns          map[string]*clientConn
	sessions       SessionValidator
	relay          RelayProxy
	envelopes      *store.EnvelopeStore
	pollInterval   time.Duration
	log            logr.Logger
	upgrader       websocket.Upgrader
	pendingFetches sync.Map // correlationID → *PendingFetch
}

// NewConnManager creates a ConnManager that validates sessions via the
// provided SessionValidator and proxies FETCH requests through relayProxy.
// Poll interval controls how often connected clients are polled for new
// messages from the relay.
func NewConnManager(
	sessions SessionValidator,
	relayProxy RelayProxy,
	envelopes *store.EnvelopeStore,
	pollInterval time.Duration,
	log logr.Logger,
) *ConnManager {
	return &ConnManager{
		conns:        make(map[string]*clientConn),
		sessions:     sessions,
		relay:        relayProxy,
		envelopes:    envelopes,
		pollInterval: pollInterval,
		log:          log,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // CORS handled by middleware
			},
		},
	}
}

// HandleUpgrade upgrades an HTTP request to a WebSocket connection. The
// connection is not considered authenticated until the client sends an
// "authenticate" message containing its session token. This avoids placing
// the token in the URL query string where it would leak into server logs,
// browser history, and Referer headers.
func (cm *ConnManager) HandleUpgrade(w http.ResponseWriter, r *http.Request) {
	wsConn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		cm.log.Error("websocket upgrade failed", logr.M("error", err.Error()))
		return
	}

	// Wait for the first message which must be an authenticate message.
	// Apply a generous timeout so slow clients aren't immediately dropped.
	wsConn.SetReadDeadline(time.Now().Add(30 * time.Second))

	_, rawMsg, err := wsConn.ReadMessage()
	if err != nil {
		wsConn.Close()
		return
	}

	var msg WSMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil || msg.Type != TypeAuthenticate {
		cm.writeWSError(wsConn, "", "first message must be authenticate")
		wsConn.Close()
		return
	}

	var authData AuthenticateData
	if err := json.Unmarshal(msg.Data, &authData); err != nil || authData.Token == "" {
		cm.writeWSError(wsConn, msg.ID, "missing token")
		wsConn.Close()
		return
	}

	address, err := cm.sessions.Validate(authData.Token)
	if err != nil {
		cm.writeWSError(wsConn, msg.ID, "invalid session")
		wsConn.Close()
		return
	}

	// Clear the read deadline now that authentication succeeded.
	wsConn.SetReadDeadline(time.Time{})

	// Confirm authentication to the client.
	ack, _ := json.Marshal(WSMessage{ID: msg.ID, Type: TypeAuthenticated})
	wsConn.WriteMessage(websocket.TextMessage, ack)

	ctx, cancel := context.WithCancel(r.Context())
	cc := &clientConn{
		conn:    wsConn,
		address: address,
		cancel:  cancel,
	}

	cm.mu.Lock()
	// Close any existing connection for this address.
	if old, ok := cm.conns[address]; ok {
		old.cancel()
		old.conn.Close()
	}
	cm.conns[address] = cc
	cm.mu.Unlock()

	cm.log.Info("websocket connected", logr.M("address", address))

	go cm.readPump(cc)
	go cm.pollLoop(ctx, cc)
}

// readPump reads JSON messages from the WebSocket and dispatches them by type.
func (cm *ConnManager) readPump(cc *clientConn) {
	defer func() {
		cm.Remove(cc.address)
	}()

	for {
		_, rawMsg, err := cc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				cm.log.Warn("websocket read error", logr.M("address", cc.address), logr.M("error", err.Error()))
			}
			return
		}

		var msg WSMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			cm.log.Warn("invalid websocket message", logr.M("address", cc.address), logr.M("error", err.Error()))
			continue
		}

		switch msg.Type {
		case TypeFetchRequest:
			cm.handleFetchRequest(cc, msg.ID)
		case TypeFetchProof:
			var data FetchProofData
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				cm.sendError(cc, msg.ID, "invalid fetch proof data")
				continue
			}
			cm.handleFetchProof(cc, msg.ID, data)
		default:
			cm.sendError(cc, msg.ID, "unknown message type: "+msg.Type)
		}
	}
}

// pollLoop periodically triggers a FETCH for connected clients.
func (cm *ConnManager) pollLoop(ctx context.Context, cc *clientConn) {
	ticker := time.NewTicker(cm.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cm.triggerFetch(cc)
		}
	}
}

// Send marshals msg as JSON and writes it to the WebSocket for address.
func (cm *ConnManager) Send(address string, msg interface{}) {
	cm.mu.RLock()
	cc, ok := cm.conns[address]
	cm.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		cm.log.Error("failed to marshal ws message", logr.M("error", err.Error()))
		return
	}

	if err := cc.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		cm.log.Warn("failed to send ws message", logr.M("address", address), logr.M("error", err.Error()))
	}
}

// Remove closes and removes the connection for the given address.
func (cm *ConnManager) Remove(address string) {
	cm.mu.Lock()
	cc, ok := cm.conns[address]
	if ok {
		delete(cm.conns, address)
	}
	cm.mu.Unlock()

	if ok {
		cc.cancel()
		cc.conn.Close()
		cm.log.Info("websocket disconnected", logr.M("address", address))
	}
}

// writeWSError writes an error frame directly to a WebSocket connection that
// may not yet be registered in the ConnManager (e.g. during authentication).
func (cm *ConnManager) writeWSError(conn *websocket.Conn, msgID, errMsg string) {
	data, _ := json.Marshal(ErrorData{Message: errMsg})
	resp, _ := json.Marshal(WSMessage{ID: msgID, Type: TypeError, Data: data})
	conn.WriteMessage(websocket.TextMessage, resp)
}

// sendError sends an error message over the WebSocket.
func (cm *ConnManager) sendError(cc *clientConn, msgID, errMsg string) {
	data, _ := json.Marshal(ErrorData{Message: errMsg})
	resp := WSMessage{
		ID:   msgID,
		Type: TypeError,
		Data: data,
	}
	cm.Send(cc.address, resp)
}
