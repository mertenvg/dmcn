// Package ws implements WebSocket connection management and protocol handling
// for the DMCN web client.
package ws

import "encoding/json"

// WebSocket message type constants.
const (
	TypeFetchRequest  = "fetch_request"
	TypeFetchChallenge = "fetch_challenge"
	TypeFetchProof    = "fetch_proof"
	TypeNewEnvelopes  = "new_envelopes"
	TypeError         = "error"
)

// WSMessage is the top-level JSON envelope for all WebSocket messages.
type WSMessage struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// FetchChallengeData carries the relay challenge nonce to the client.
type FetchChallengeData struct {
	Nonce string `json:"nonce"` // base64
}

// FetchProofData carries the client's signed proof back to the server.
type FetchProofData struct {
	Nonce     string `json:"nonce"`     // base64
	Signature string `json:"signature"` // base64
}

// EnvelopeData represents a single encrypted envelope in transit.
type EnvelopeData struct {
	Hash string `json:"hash"` // hex
	Data string `json:"data"` // base64 protobuf
}

// NewEnvelopesData carries one or more fetched envelopes to the client.
type NewEnvelopesData struct {
	Envelopes []EnvelopeData `json:"envelopes"`
}

// ErrorData carries an error message to the client.
type ErrorData struct {
	Message string `json:"message"`
}
