package network

import (
	"encoding/json"
	"testing"
)

type mockHandler struct {
	received []GossipMessage
}

func (m *mockHandler) HandleGossip(msg GossipMessage) {
	m.received = append(m.received, msg)
}

func TestBroadcastAndHandleGossip(t *testing.T) {
	transport := NewInMemoryTransport()
	nodeA := NewNode("A", transport)
	nodeB := NewNode("B", transport)
	nodeC := NewNode("C", transport)

	nodeA.ConnectToPeer("B")
	nodeA.ConnectToPeer("C")

	transport.Register(nodeA)
	transport.Register(nodeB)
	transport.Register(nodeC)

	handlerB := &mockHandler{}
	handlerC := &mockHandler{}

	payload, _ := json.Marshal(map[string]string{"value": "update"})
	gossip := GossipMessage{Topic: "story", Payload: payload}

	if err := BroadcastToNetwork(nodeA, gossip); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	msgB := <-nodeB.ReceiveMessages()
	if err := HandleIncomingMessage(handlerB, msgB); err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	msgC := <-nodeC.ReceiveMessages()
	if err := HandleIncomingMessage(handlerC, msgC); err != nil {
		t.Fatalf("handler failed: %v", err)
	}

	if len(handlerB.received) != 1 || handlerB.received[0].Topic != "story" {
		t.Fatalf("expected handler B to receive gossip")
	}

	if len(handlerC.received) != 1 {
		t.Fatalf("expected handler C to receive gossip")
	}
}
