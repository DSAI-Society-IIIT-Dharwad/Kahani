package consensus

import (
	"context"
	"testing"
	"time"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

func TestStartServiceLifecycle(t *testing.T) {
	transport := network.NewInMemoryTransport()
	nodeA := network.NewNode("node-1", transport)
	nodeB := network.NewNode("node-2", transport)

	transport.Register(nodeA)
	transport.Register(nodeB)

	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)

	id, ch := bus.Subscribe(8)
	t.Cleanup(func() { bus.Unsubscribe(id) })

	service, err := StartService(context.Background(), chain, bus, map[string]*network.Node{
		"node-1": nodeA,
		"node-2": nodeB,
	}, []string{"node-1", "node-2"}, mockSigner{}, 0)
	if err != nil {
		t.Fatalf("failed to start service: %v", err)
	}
	defer service.Stop()

	tx := types.Transaction{TxID: "tx-1", Type: "story"}
	chain.EnqueueTransaction(tx)

	if err := service.Propose("node-1", nil); err != nil {
		t.Fatalf("propose failed: %v", err)
	}

	timeout := time.After(2 * time.Second)
	received := make(map[observer.EventType]int)

	for len(received) < 3 {
		select {
		case ev := <-ch:
			received[ev.Type]++
		case <-timeout:
			t.Fatalf("timed out waiting for events, got %#v", received)
		}
	}

	if received[observer.EventBlockCommitted] == 0 {
		t.Fatalf("expected block committed event, got %#v", received)
	}

	if received[observer.EventTransactionCommitted] == 0 {
		t.Fatalf("expected transaction committed event, got %#v", received)
	}

	if received[observer.EventTransactionQueued] == 0 {
		t.Fatalf("expected transaction queued event, got %#v", received)
	}

	if len(chain.Blocks()) != 2 {
		t.Fatalf("expected block to be committed, got %d", len(chain.Blocks()))
	}

	if len(chain.PendingTransactions()) != 0 {
		t.Fatalf("expected pending transactions to be cleared")
	}
}

func TestStartServiceErrors(t *testing.T) {
	_, err := StartService(context.Background(), nil, nil, map[string]*network.Node{}, nil, nil, 0)
	if err == nil {
		t.Fatal("expected error when blockchain missing")
	}

	chain := blockchain.NewBlockchain()
	if _, err := StartService(context.Background(), chain, nil, map[string]*network.Node{}, nil, nil, 0); err == nil {
		t.Fatal("expected error when transports missing")
	}
}

func TestServiceProposeValidation(t *testing.T) {
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	chain := blockchain.NewBlockchain()
	service, err := StartService(context.Background(), chain, nil, map[string]*network.Node{"node-1": node}, []string{"node-1"}, mockSigner{}, 0)
	if err != nil {
		t.Fatalf("failed to start service: %v", err)
	}
	defer service.Stop()

	if err := service.Propose("node-2", nil); err == nil {
		t.Fatal("expected error for unknown node")
	}

	if err := service.Propose("node-1", nil); err == nil {
		t.Fatal("expected error when no transactions available")
	}

	tx := types.Transaction{TxID: "tx-1"}
	if err := service.Propose("node-1", []types.Transaction{tx}); err != nil {
		t.Fatalf("propose with transactions should succeed: %v", err)
	}
}
