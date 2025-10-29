package network

import (
	"encoding/json"
	"errors"
)

// GossipMessage represents the payload shared via the gossip protocol.
type GossipMessage struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// GossipHandler processes inbound gossip messages.
type GossipHandler interface {
	HandleGossip(msg GossipMessage)
}

// BroadcastToNetwork sends the gossip message to all connected peers.
func BroadcastToNetwork(node *Node, message GossipMessage) error {
	if node == nil {
		return errors.New("network: node is nil")
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, peer := range node.Peers() {
		if err := node.SendMessage(peer, payload); err != nil {
			return err
		}
	}

	return nil
}

// HandleIncomingMessage decodes an inbound payload and delegates to the handler.
func HandleIncomingMessage(handler GossipHandler, msg Message) error {
	if handler == nil {
		return errors.New("network: handler is nil")
	}

	var gossip GossipMessage
	if err := json.Unmarshal(msg.Payload, &gossip); err != nil {
		return err
	}

	handler.HandleGossip(gossip)
	return nil
}
