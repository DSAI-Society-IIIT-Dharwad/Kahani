package consensus

import (
	"encoding/json"

	"storytelling-blockchain/internal/types"
)

// MessageType represents the stage of the PBFT protocol.
type MessageType string

const (
	MessagePrePrepare MessageType = "PRE_PREPARE"
	MessagePrepare    MessageType = "PREPARE"
	MessageCommit     MessageType = "COMMIT"
	MessageViewChange MessageType = "VIEW_CHANGE"
)

// Message encapsulates the payload exchanged between validators during consensus.
type Message struct {
	Type      MessageType `json:"type"`
	View      int         `json:"view"`
	Sequence  int         `json:"sequence"`
	Block     types.Block `json:"block"`
	SenderID  string      `json:"sender_id"`
	Signature string      `json:"signature"`
}

// Digest returns a deterministic hash of the message contents for signing.
func (m Message) Digest() ([]byte, error) {
	clone := m
	clone.Signature = ""
	payload, err := json.Marshal(clone)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
