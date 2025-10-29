package consensus

import (
	"errors"
	"sync"
	"testing"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/types"
)

type mockNetwork struct {
	mu    sync.Mutex
	nodes map[string]*PBFTNode
}

func newMockNetwork() *mockNetwork {
	return &mockNetwork{nodes: make(map[string]*PBFTNode)}
}

func (m *mockNetwork) register(node *PBFTNode) {
	m.mu.Lock()
	m.nodes[node.id] = node
	m.mu.Unlock()
}

func (m *mockNetwork) Broadcast(sender string, msg Message) error {
	m.mu.Lock()
	nodes := make([]*PBFTNode, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	m.mu.Unlock()

	for _, node := range nodes {
		_ = node.HandleMessage(msg)
	}

	return nil
}

func (m *mockNetwork) Send(sender, recipient string, msg Message) error {
	m.mu.Lock()
	target, ok := m.nodes[recipient]
	m.mu.Unlock()
	if !ok {
		return errors.New("recipient not found")
	}
	return target.HandleMessage(msg)
}

type mockBuilder struct{}

type builderResponse struct {
	sync.Mutex
	counter int
}

func (mockBuilder) BuildBlock(txs []types.Transaction) (types.Block, error) {
	prev := "prev"
	block := blockchain.NewBlock(len(txs), prev, txs)
	return block, nil
}

type mockFinalizer struct {
	mu   sync.Mutex
	hits map[string]int
	wg   *sync.WaitGroup
}

func newMockFinalizer(count int) *mockFinalizer {
	return &mockFinalizer{hits: make(map[string]int), wg: &sync.WaitGroup{}}
}

func (f *mockFinalizer) setExpected(expected int) {
	for i := 0; i < expected; i++ {
		f.wg.Add(1)
	}
}

func (f *mockFinalizer) Finalize(block types.Block, nodeID string) {
	f.mu.Lock()
	f.hits[nodeID]++
	f.mu.Unlock()
	f.wg.Done()
}

type finalizerAdapter struct {
	nodeID   string
	finalize *mockFinalizer
}

func (f finalizerAdapter) call(block types.Block) {
	f.finalize.Finalize(block, f.nodeID)
}

type mockSigner struct{}

func (mockSigner) Sign(data []byte) (string, error) {
	return string(data), nil
}

func (mockSigner) Verify(sender string, data []byte, signature string) bool {
	return string(data) == signature
}

type rejectingSigner struct{}

func (rejectingSigner) Sign(data []byte) (string, error) {
	return string(data), nil
}

func (rejectingSigner) Verify(string, []byte, string) bool {
	return false
}

func TestPBFTConsensusFlow(t *testing.T) {
	network := newMockNetwork()
	finalize := newMockFinalizer(0)

	nodeIDs := []string{"node-1", "node-2", "node-3", "node-4"}
	nodes := make([]*PBFTNode, 0, len(nodeIDs))

	for _, id := range nodeIDs {
		node, err := NewPBFTNode(Config{
			ID:             id,
			Peers:          nodeIDs,
			FaultTolerance: 1,
			Network:        network,
			Signer:         mockSigner{},
			Builder:        mockBuilder{},
			Finalize: func(block types.Block) {
				finalize.Finalize(block, id)
			},
		})
		if err != nil {
			t.Fatalf("failed to create node: %v", err)
		}
		network.register(node)
		nodes = append(nodes, node)
	}

	finalize.setExpected(len(nodes))

	txs := []types.Transaction{{TxID: "tx-1", Type: "test"}}
	if err := nodes[0].ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	finalize.wg.Wait()

	finalize.mu.Lock()
	defer finalize.mu.Unlock()

	if len(finalize.hits) != len(nodes) {
		t.Fatalf("expected finalize called on all nodes, got %d", len(finalize.hits))
	}
}

func TestPBFTRejectsDuplicateMessages(t *testing.T) {
	network := newMockNetwork()
	finalize := newMockFinalizer(0)

	node, err := NewPBFTNode(Config{
		ID:             "node-1",
		Peers:          []string{"node-1"},
		FaultTolerance: 0,
		Network:        network,
		Signer:         mockSigner{},
		Builder:        mockBuilder{},
		Finalize: func(block types.Block) {
			finalize.Finalize(block, "node-1")
		},
	})
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}
	network.register(node)
	finalize.setExpected(1)

	txs := []types.Transaction{{TxID: "1"}}
	if err := node.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	inst := node.instances[node.sequence]
	if inst == nil {
		t.Fatalf("expected instance to exist")
	}

	msg := Message{
		Type:     MessageCommit,
		Sequence: node.sequence,
		Block:    inst.block,
		SenderID: "node-1",
	}

	_ = node.HandleMessage(msg)
	_ = node.HandleMessage(msg)

	finalize.mu.Lock()
	count := finalize.hits["node-1"]
	finalize.mu.Unlock()

	if count != 1 {
		t.Fatalf("expected finalize once, got %d", count)
	}
}

func TestPBFTIgnoresInvalidSignatures(t *testing.T) {
	network := newMockNetwork()

	var mu sync.Mutex
	called := false

	node, err := NewPBFTNode(Config{
		ID:             "node-1",
		Peers:          []string{"node-2"},
		FaultTolerance: 0,
		Network:        network,
		Signer:         rejectingSigner{},
		Builder:        mockBuilder{},
		Finalize: func(types.Block) {
			mu.Lock()
			called = true
			mu.Unlock()
		},
	})
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	network.register(node)

	txs := []types.Transaction{{TxID: "1"}}
	if err := node.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if called {
		t.Fatalf("finalize should not run for invalid signatures")
	}
}
