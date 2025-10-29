package consensus

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/types"
)

type countingFinalizer struct {
	hits map[string]int
	mu   sync.Mutex
	wg   sync.WaitGroup
}

type noopNetwork struct{}

func (noopNetwork) Broadcast(string, Message) error    { return nil }
func (noopNetwork) Send(string, string, Message) error { return nil }

func newCountingFinalizer(expected int) *countingFinalizer {
	f := &countingFinalizer{hits: make(map[string]int)}
	f.wg.Add(expected)
	return f
}

func (f *countingFinalizer) handler(nodeID string) Finalizer {
	return func(block types.Block) {
		f.mu.Lock()
		f.hits[nodeID]++
		f.mu.Unlock()
		f.wg.Done()
	}
}

func (f *countingFinalizer) waitWithTimeout(t *testing.T, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatalf("timed out waiting for finalize: %v", ctx.Err())
	}
}

func forwardMessages(ctx context.Context, node *network.Node, handler network.GossipHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-node.ReceiveMessages():
			_ = network.HandleIncomingMessage(handler, msg)
		}
	}
}

func TestGossipNetworkIntegration(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	nodeA.ConnectToPeer("node-2")
	nodeB.ConnectToPeer("node-1")

	transport.Register(nodeA)
	transport.Register(nodeB)

	gossipNetA, err := NewGossipNetwork(nodeA)
	if err != nil {
		t.Fatalf("failed to create gossip network A: %v", err)
	}

	gossipNetB, err := NewGossipNetwork(nodeB)
	if err != nil {
		t.Fatalf("failed to create gossip network B: %v", err)
	}

	finalizer := newCountingFinalizer(2)

	nodeIDs := []string{"node-1", "node-2"}
	pbftA, err := NewPBFTNode(Config{
		ID:             "node-1",
		Peers:          nodeIDs,
		FaultTolerance: 0,
		Network:        gossipNetA,
		Signer:         mockSigner{},
		Builder:        mockBuilder{},
		Finalize:       finalizer.handler("node-1"),
	})
	if err != nil {
		t.Fatalf("failed to create pbft node A: %v", err)
	}

	pbftB, err := NewPBFTNode(Config{
		ID:             "node-2",
		Peers:          nodeIDs,
		FaultTolerance: 0,
		Network:        gossipNetB,
		Signer:         mockSigner{},
		Builder:        mockBuilder{},
		Finalize:       finalizer.handler("node-2"),
	})
	if err != nil {
		t.Fatalf("failed to create pbft node B: %v", err)
	}

	handlerA, err := NewPBFTGossipHandler(pbftA)
	if err != nil {
		t.Fatalf("failed to create handler A: %v", err)
	}

	handlerB, err := NewPBFTGossipHandler(pbftB)
	if err != nil {
		t.Fatalf("failed to create handler B: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go forwardMessages(ctx, nodeA, handlerA)
	go forwardMessages(ctx, nodeB, handlerB)

	txs := []types.Transaction{{TxID: "tx-1", Type: "story"}}
	if err := pbftA.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	finalizer.waitWithTimeout(t, 2*time.Second)

	finalizer.mu.Lock()
	defer finalizer.mu.Unlock()

	if finalizer.hits["node-1"] != 1 {
		t.Fatalf("expected node-1 to finalize once, got %d", finalizer.hits["node-1"])
	}

	if finalizer.hits["node-2"] != 1 {
		t.Fatalf("expected node-2 to finalize once, got %d", finalizer.hits["node-2"])
	}
}

func TestPBFTGossipHandlerIgnoresOtherTopics(t *testing.T) {
	pbft, err := NewPBFTNode(Config{
		ID:             "node-1",
		Peers:          []string{"node-1"},
		FaultTolerance: 0,
		Network:        noopNetwork{},
		Signer:         mockSigner{},
		Builder:        mockBuilder{},
		Finalize: func(types.Block) {
			t.Fatalf("finalize should not be called")
		},
	})
	if err != nil {
		t.Fatalf("failed to create pbft node: %v", err)
	}

	handler, err := NewPBFTGossipHandler(pbft)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	msg := network.GossipMessage{Topic: "other/topic", Payload: []byte(`{"x":1}`)}
	handler.HandleGossip(msg)
}

func TestPBFTGossipHandlerInvalidPayload(t *testing.T) {
	var finalized bool
	pbft, err := NewPBFTNode(Config{
		ID:             "node-1",
		Peers:          []string{"node-1"},
		FaultTolerance: 0,
		Network:        noopNetwork{},
		Signer:         mockSigner{},
		Builder:        mockBuilder{},
		Finalize: func(types.Block) {
			finalized = true
		},
	})
	if err != nil {
		t.Fatalf("failed to create pbft node: %v", err)
	}

	handler, err := NewPBFTGossipHandler(pbft)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	msg := network.GossipMessage{Topic: gossipTopicPBFT, Payload: []byte("not-json")}
	handler.HandleGossip(msg)

	if finalized {
		t.Fatalf("finalize should not run for invalid payload")
	}
}

func TestNewGossipNetworkErrors(t *testing.T) {
	if _, err := NewGossipNetwork(nil); err == nil {
		t.Fatal("expected error when node is nil")
	}
}

func TestGossipNetworkBroadcast(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	nodeA.ConnectToPeer("node-2")
	nodeB.ConnectToPeer("node-1")

	transport.Register(nodeA)
	transport.Register(nodeB)

	gossipNet, err := NewGossipNetwork(nodeA)
	if err != nil {
		t.Fatalf("new gossip network failed: %v", err)
	}

	msg := Message{Type: MessageCommit, Sequence: 1, Block: types.Block{Hash: "hash"}, SenderID: "node-1"}
	if err := gossipNet.Broadcast("node-1", msg); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	select {
	case received := <-nodeB.ReceiveMessages():
		if received.From != "node-1" {
			t.Fatalf("unexpected sender: %s", received.From)
		}

		var gossip network.GossipMessage
		if err := json.Unmarshal(received.Payload, &gossip); err != nil {
			t.Fatalf("decode gossip failed: %v", err)
		}

		if gossip.Topic != gossipTopicPBFT {
			t.Fatalf("unexpected topic: %s", gossip.Topic)
		}

		var decoded Message
		if err := json.Unmarshal(gossip.Payload, &decoded); err != nil {
			t.Fatalf("decode message failed: %v", err)
		}

		if decoded.Sequence != 1 || decoded.Block.Hash != "hash" {
			t.Fatalf("unexpected decoded message: %+v", decoded)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestGossipNetworkSend(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	nodeA.ConnectToPeer("node-2")
	nodeB.ConnectToPeer("node-1")

	transport.Register(nodeA)
	transport.Register(nodeB)

	gossipNet, err := NewGossipNetwork(nodeA)
	if err != nil {
		t.Fatalf("new gossip network failed: %v", err)
	}

	msg := Message{Type: MessagePrepare, Sequence: 2, Block: types.Block{Hash: "hash-2"}, SenderID: "node-1"}
	if err := gossipNet.Send("node-1", "node-2", msg); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	select {
	case received := <-nodeB.ReceiveMessages():
		if received.From != "node-1" {
			t.Fatalf("unexpected sender: %s", received.From)
		}

		var decoded Message
		if err := json.Unmarshal(received.Payload, &decoded); err != nil {
			t.Fatalf("decode message failed: %v", err)
		}

		if decoded.Type != MessagePrepare || decoded.Sequence != 2 {
			t.Fatalf("unexpected message: %+v", decoded)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for send")
	}
}

func TestGossipNetworkBroadcastPeerError(t *testing.T) {
	nodeA := network.NewNode("node-1", network.NewInMemoryTransport())
	nodeA.ConnectToPeer("missing")

	gossipNet, err := NewGossipNetwork(nodeA)
	if err != nil {
		t.Fatalf("new gossip network failed: %v", err)
	}

	msg := Message{Type: MessageCommit, Sequence: 3, Block: types.Block{Hash: "hash"}}
	if err := gossipNet.Broadcast("node-1", msg); err == nil {
		t.Fatal("expected broadcast to fail when peer missing")
	}
}

func TestGossipNetworkSendPeerError(t *testing.T) {
	nodeA := network.NewNode("node-1", network.NewInMemoryTransport())

	gossipNet, err := NewGossipNetwork(nodeA)
	if err != nil {
		t.Fatalf("new gossip network failed: %v", err)
	}

	msg := Message{Type: MessagePrepare, Sequence: 4, Block: types.Block{Hash: "hash"}}
	if err := gossipNet.Send("node-1", "not-connected", msg); err == nil {
		t.Fatal("expected send to fail for unknown peer")
	}
}
