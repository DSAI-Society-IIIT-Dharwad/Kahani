package consensus

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/types"
)

type waitFinalizer struct {
	mu   sync.Mutex
	hits map[string]int
	wg   sync.WaitGroup
}

type failingTransport struct{}

func (failingTransport) Send(string, string, []byte) error {
	return errors.New("transport unavailable")
}

func newWaitFinalizer(expected int) *waitFinalizer {
	wf := &waitFinalizer{hits: make(map[string]int)}
	wf.wg.Add(expected)
	return wf
}

func (wf *waitFinalizer) finalize(id string) Finalizer {
	return func(block types.Block) {
		wf.mu.Lock()
		wf.hits[id]++
		wf.mu.Unlock()
		wf.wg.Done()
	}
}

func (wf *waitFinalizer) wait(t *testing.T) {
	done := make(chan struct{})
	go func() {
		wf.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for finalizer")
	}
}

func TestBootstrapNodeSuccess(t *testing.T) {
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	finalizer := newWaitFinalizer(1)

	runtime, err := StartNode(context.Background(), BootstrapOptions{
		NodeID:         "node-1",
		Peers:          []string{"node-1"},
		FaultTolerance: 0,
		Transport:      node,
		Builder:        mockBuilder{},
		Finalize:       finalizer.finalize("node-1"),
		Signer:         mockSigner{},
	})
	if err != nil {
		t.Fatalf("bootstrap node failed: %v", err)
	}
	defer runtime.Stop()

	txs := []types.Transaction{{TxID: "tx-1", Type: "story"}}
	if err := runtime.Node.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	finalizer.wait(t)

	finalizer.mu.Lock()
	defer finalizer.mu.Unlock()
	if finalizer.hits["node-1"] != 1 {
		t.Fatalf("expected finalize called once, got %d", finalizer.hits["node-1"])
	}
}

func TestBootstrapNodeErrors(t *testing.T) {
	_, _, err := BootstrapNode(BootstrapOptions{})
	if err == nil {
		t.Fatal("expected error for missing transport")
	}

	transport := network.NewNode("node-1", network.NewInMemoryTransport())
	_, _, err = BootstrapNode(BootstrapOptions{
		NodeID:    "node-1",
		Transport: transport,
		Finalize: func(types.Block) {
		},
	})
	if err == nil {
		t.Fatal("expected error for missing builder")
	}

	_, _, err = BootstrapNode(BootstrapOptions{
		NodeID:    "node-1",
		Transport: transport,
		Builder:   mockBuilder{},
	})
	if err == nil {
		t.Fatal("expected error for missing finalize")
	}
}

func TestBootstrapCluster(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	transport.Register(nodeA)
	transport.Register(nodeB)

	transports := map[string]*network.Node{
		"node-1": nodeA,
		"node-2": nodeB,
	}

	finalizers := map[string]Finalizer{
		"node-1": func(types.Block) {},
		"node-2": func(types.Block) {},
	}

	nodes, handlers, err := BootstrapCluster(transports, []string{"node-1", "node-2"}, mockBuilder{}, finalizers, mockSigner{}, 0)
	if err != nil {
		t.Fatalf("bootstrap cluster failed: %v", err)
	}

	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if len(handlers) != 2 {
		t.Fatalf("expected 2 handlers, got %d", len(handlers))
	}
}

func TestStartClusterIntegration(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	transport.Register(nodeA)
	transport.Register(nodeB)

	transports := map[string]*network.Node{
		"node-1": nodeA,
		"node-2": nodeB,
	}

	finalizer := newWaitFinalizer(2)
	finalizers := map[string]Finalizer{
		"node-1": finalizer.finalize("node-1"),
		"node-2": finalizer.finalize("node-2"),
	}

	runtimes, err := StartCluster(context.Background(), transports, []string{"node-1", "node-2"}, mockBuilder{}, finalizers, mockSigner{}, 0)
	if err != nil {
		t.Fatalf("start cluster failed: %v", err)
	}

	for _, runtime := range runtimes {
		defer runtime.Stop()
	}

	txs := []types.Transaction{{TxID: "tx-1", Type: "story"}}
	if err := runtimes["node-1"].Node.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	finalizer.wait(t)

	finalizer.mu.Lock()
	defer finalizer.mu.Unlock()
	if finalizer.hits["node-1"] != 1 {
		t.Fatalf("expected node-1 finalize once, got %d", finalizer.hits["node-1"])
	}
	if finalizer.hits["node-2"] != 1 {
		t.Fatalf("expected node-2 finalize once, got %d", finalizer.hits["node-2"])
	}
}

func TestStartClusterErrors(t *testing.T) {
	_, err := StartCluster(context.Background(), map[string]*network.Node{}, nil, mockBuilder{}, map[string]Finalizer{}, mockSigner{}, 0)
	if err == nil {
		t.Fatal("expected error for missing transports")
	}

	transports := map[string]*network.Node{
		"node-1": network.NewNode("node-1", network.NewInMemoryTransport()),
	}

	_, err = StartCluster(context.Background(), transports, []string{"node-1"}, nil, map[string]Finalizer{"node-1": func(types.Block) {}}, mockSigner{}, 0)
	if err == nil {
		t.Fatal("expected error for missing builder")
	}

	_, err = StartCluster(context.Background(), transports, []string{"node-1"}, mockBuilder{}, map[string]Finalizer{}, mockSigner{}, 0)
	if err == nil {
		t.Fatal("expected error for missing finalizer")
	}
}

func TestApplyBlocks(t *testing.T) {
	blocks := []types.Block{
		{Hash: "hash-1"},
		{Hash: "hash-2"},
	}

	var mu sync.Mutex
	var seen []string
	finalize := func(block types.Block) {
		mu.Lock()
		seen = append(seen, block.Hash)
		mu.Unlock()
	}

	if err := ApplyBlocks(blocks, finalize); err != nil {
		t.Fatalf("apply blocks failed: %v", err)
	}

	if len(seen) != 2 {
		t.Fatalf("expected 2 blocks finalized, got %d", len(seen))
	}

	mu.Lock()
	defer mu.Unlock()
	if seen[0] != "hash-1" || seen[1] != "hash-2" {
		t.Fatalf("blocks applied out of order: %v", seen)
	}
}

func TestStartClusterIgnoresBadGossip(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	transport.Register(nodeA)
	transport.Register(nodeB)

	transports := map[string]*network.Node{
		"node-1": nodeA,
		"node-2": nodeB,
	}

	finalizer := newWaitFinalizer(2)
	finalizers := map[string]Finalizer{
		"node-1": finalizer.finalize("node-1"),
		"node-2": finalizer.finalize("node-2"),
	}

	runtimes, err := StartCluster(context.Background(), transports, []string{"node-1", "node-2"}, mockBuilder{}, finalizers, mockSigner{}, 0)
	if err != nil {
		t.Fatalf("start cluster failed: %v", err)
	}

	for _, runtime := range runtimes {
		defer runtime.Stop()
	}

	// Inject unrelated gossip; handlers should ignore it without affecting consensus.
	runtimes["node-1"].Handler.HandleGossip(network.GossipMessage{Topic: "unrelated", Payload: []byte(`{"ok":true}`)})

	txs := []types.Transaction{{TxID: "tx-2", Type: "story"}}
	if err := runtimes["node-2"].Node.ProposeBlock(txs); err != nil {
		t.Fatalf("propose block failed: %v", err)
	}

	finalizer.wait(t)

	finalizer.mu.Lock()
	defer finalizer.mu.Unlock()
	if finalizer.hits["node-1"] != 1 || finalizer.hits["node-2"] != 1 {
		t.Fatalf("expected both nodes to finalize exactly once, got %+v", finalizer.hits)
	}
}

func TestStartNodeFailsWhenTransportDrops(t *testing.T) {
	node := network.NewNode("node-1", failingTransport{})

	var mu sync.Mutex
	called := false

	runtime, err := StartNode(context.Background(), BootstrapOptions{
		NodeID:         "node-1",
		Peers:          []string{"node-2"},
		FaultTolerance: 0,
		Transport:      node,
		Builder:        mockBuilder{},
		Finalize: func(types.Block) {
			mu.Lock()
			called = true
			mu.Unlock()
		},
	})
	if err != nil {
		t.Fatalf("start node failed: %v", err)
	}
	defer runtime.Stop()

	txs := []types.Transaction{{TxID: "tx-err"}}
	if err := runtime.Node.ProposeBlock(txs); err == nil {
		t.Fatal("expected propose block to fail when transport drops messages")
	}

	mu.Lock()
	defer mu.Unlock()
	if called {
		t.Fatalf("finalize should not run when transport drops messages")
	}
}
