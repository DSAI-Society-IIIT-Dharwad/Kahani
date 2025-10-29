package consensus

import (
	"encoding/json"
	"errors"

	"storytelling-blockchain/internal/network"
)

const gossipTopicPBFT = "consensus/pbft"

// GossipNetwork adapts a network.Node to the consensus.Network interface.
type GossipNetwork struct {
	node *network.Node
}

// NewGossipNetwork wraps the provided node so it can be used by a PBFT node.
func NewGossipNetwork(node *network.Node) (*GossipNetwork, error) {
	if node == nil {
		return nil, errors.New("consensus: node is nil")
	}

	return &GossipNetwork{node: node}, nil
}

// Broadcast encodes the PBFT message and gossips it across the network.
func (g *GossipNetwork) Broadcast(sender string, msg Message) error {
	if g == nil {
		return errors.New("consensus: gossip network is nil")
	}
	if g.node == nil {
		return errors.New("consensus: node is nil")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	gossip := network.GossipMessage{Topic: gossipTopicPBFT, Payload: payload}
	return network.BroadcastToNetwork(g.node, gossip)
}

// Send delivers the message to a specific peer via direct messaging.
func (g *GossipNetwork) Send(sender, recipient string, msg Message) error {
	if g == nil {
		return errors.New("consensus: gossip network is nil")
	}
	if g.node == nil {
		return errors.New("consensus: node is nil")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return g.node.SendMessage(recipient, payload)
}

// PBFTGossipHandler forwards PBFT gossip messages into the consensus node.
type PBFTGossipHandler struct {
	node *PBFTNode
}

// NewPBFTGossipHandler creates a gossip handler that feeds messages to node.
func NewPBFTGossipHandler(node *PBFTNode) (*PBFTGossipHandler, error) {
	if node == nil {
		return nil, errors.New("consensus: pbft node is nil")
	}
	return &PBFTGossipHandler{node: node}, nil
}

// HandleGossip decodes consensus messages and hands them to the PBFT node.
func (h *PBFTGossipHandler) HandleGossip(msg network.GossipMessage) {
	if h == nil || h.node == nil {
		return
	}
	if msg.Topic != gossipTopicPBFT {
		return
	}

	var consensusMsg Message
	if err := json.Unmarshal(msg.Payload, &consensusMsg); err != nil {
		return
	}

	_ = h.node.HandleMessage(consensusMsg)
}
