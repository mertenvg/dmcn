package ws

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mertenvg/logr/v2"
	"google.golang.org/protobuf/proto"
)

// handleFetchRequest initiates a relay FETCH challenge on behalf of the
// connected client. The resulting nonce is sent back as a FetchChallenge
// message and the open stream is stored keyed by the message ID so that
// handleFetchProof can complete the exchange.
func (cm *ConnManager) handleFetchRequest(cc *clientConn, msgID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	nonce, stream, err := cm.relay.FetchChallenge(ctx, cc.address)
	if err != nil {
		cm.log.Error("fetch challenge failed", logr.M("address", cc.address), logr.M("error", err.Error()))
		cm.sendError(cc, msgID, "fetch challenge failed: "+err.Error())
		return
	}

	cm.pendingFetches.Store(msgID, &PendingFetch{
		Nonce:     nonce,
		Stream:    stream,
		ExpiresAt: time.Now().Add(30 * time.Second),
	})

	challengeData, _ := json.Marshal(FetchChallengeData{
		Nonce: base64.StdEncoding.EncodeToString(nonce),
	})

	cm.Send(cc.address, WSMessage{
		ID:   msgID,
		Type: TypeFetchChallenge,
		Data: challengeData,
	})
}

// handleFetchProof completes the relay FETCH exchange using the client's
// signed proof. Received envelopes are persisted to the envelope store and
// forwarded to the client as a NewEnvelopes message.
func (cm *ConnManager) handleFetchProof(cc *clientConn, msgID string, data FetchProofData) {
	val, ok := cm.pendingFetches.LoadAndDelete(msgID)
	if !ok {
		cm.sendError(cc, msgID, "no pending fetch for this message ID")
		return
	}
	pf := val.(*PendingFetch)

	if time.Now().After(pf.ExpiresAt) {
		pf.Stream.Reset()
		cm.sendError(cc, msgID, "fetch challenge expired")
		return
	}

	nonce, err := base64.StdEncoding.DecodeString(data.Nonce)
	if err != nil {
		pf.Stream.Reset()
		cm.sendError(cc, msgID, "invalid nonce encoding")
		return
	}

	signature, err := base64.StdEncoding.DecodeString(data.Signature)
	if err != nil {
		pf.Stream.Reset()
		cm.sendError(cc, msgID, "invalid signature encoding")
		return
	}

	envelopes, hashes, err := cm.relay.FetchComplete(pf.Stream, cc.address, nonce, signature)
	if err != nil {
		cm.log.Error("fetch complete failed", logr.M("address", cc.address), logr.M("error", err.Error()))
		cm.sendError(cc, msgID, "fetch complete failed: "+err.Error())
		return
	}

	var envData []EnvelopeData
	for i, env := range envelopes {
		pbBytes, err := proto.Marshal(env.ToProto())
		if err != nil {
			cm.log.Error("failed to marshal envelope", logr.M("error", err.Error()))
			continue
		}

		hashHex := hex.EncodeToString(hashes[i][:])

		// Persist to envelope store.
		if err := cm.envelopes.Store(cc.address, hashes[i], pbBytes); err != nil {
			cm.log.Error("failed to store envelope", logr.M("error", err.Error()))
		}

		envData = append(envData, EnvelopeData{
			Hash: hashHex,
			Data: base64.StdEncoding.EncodeToString(pbBytes),
		})
	}

	respData, _ := json.Marshal(NewEnvelopesData{
		Envelopes: envData,
	})

	cm.Send(cc.address, WSMessage{
		ID:   msgID,
		Type: TypeNewEnvelopes,
		Data: respData,
	})
}

// triggerFetch initiates a FETCH request on behalf of a connected client
// during periodic polling. It generates a correlation ID and follows the
// same challenge flow as handleFetchRequest.
func (cm *ConnManager) triggerFetch(cc *clientConn) {
	msgID := fmt.Sprintf("poll-%d", time.Now().UnixNano())
	cm.handleFetchRequest(cc, msgID)
}
