package network

import (
	"testing"
)

func TestNodeMessaging(t *testing.T) {
	transport := NewInMemoryTransport()
	nodeA := NewNode("A", transport)
	nodeB := NewNode("B", transport)

	nodeA.ConnectToPeer("B")
	nodeB.ConnectToPeer("A")

	transport.Register(nodeA)
	transport.Register(nodeB)

	payload := []byte("hello")
	if err := nodeA.SendMessage("B", payload); err != nil {
		t.Fatalf("send message failed: %v", err)
	}

	msg := <-nodeB.ReceiveMessages()
	if msg.From != "A" {
		t.Fatalf("expected message from A")
	}

	if string(msg.Payload) != "hello" {
		t.Fatalf("unexpected payload %s", string(msg.Payload))
	}
}

func TestSendToUnknownPeer(t *testing.T) {
	node := NewNode("A", nil)
	if err := node.SendMessage("B", []byte("test")); err == nil {
		t.Fatalf("expected error when transport missing")
	}
}
