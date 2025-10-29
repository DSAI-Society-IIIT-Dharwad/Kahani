package consensus

import (
	"testing"
	"time"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

func TestChainBlockBuilder(t *testing.T) {
	chain := blockchain.NewBlockchain()
	builder := &ChainBlockBuilder{Chain: chain}

	txs := []types.Transaction{{TxID: "tx-1", Type: "test"}}
	block, err := builder.BuildBlock(txs)
	if err != nil {
		t.Fatalf("expected block to be built: %v", err)
	}

	if block.Index != 1 {
		t.Fatalf("expected block index 1, got %d", block.Index)
	}

	if len(block.Transactions) != 1 {
		t.Fatalf("expected transactions to be populated")
	}
}

func TestChainBlockBuilderRequiresTransactions(t *testing.T) {
	builder := &ChainBlockBuilder{Chain: blockchain.NewBlockchain()}
	if _, err := builder.BuildBlock(nil); err == nil {
		t.Fatal("expected error when transactions missing")
	}
}

func TestChainFinalizerPublishesEvents(t *testing.T) {
	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()

	id, ch := bus.Subscribe(4)
	t.Cleanup(func() { bus.Unsubscribe(id) })

	builder := &ChainBlockBuilder{Chain: chain}
	txs := []types.Transaction{{TxID: "tx-1", Type: "test"}}
	block, err := builder.BuildBlock(txs)
	if err != nil {
		t.Fatalf("failed to build block: %v", err)
	}

	finalizer := NewChainFinalizer(chain, bus)
	finalizer(block)

	// Ensure block added to chain.
	if len(chain.Blocks()) != 2 {
		t.Fatalf("expected chain to contain new block")
	}

	// Expect block committed and transaction events.
	received := make([]observer.Event, 0)
	timeout := time.After(time.Second)

	for len(received) < 2 {
		select {
		case ev := <-ch:
			received = append(received, ev)
		case <-timeout:
			t.Fatalf("timed out waiting for events")
		}
	}

	if received[0].Type != observer.EventBlockCommitted {
		t.Fatalf("expected first event to be block committed, got %s", received[0].Type)
	}

	if received[1].Type != observer.EventTransactionCommitted {
		t.Fatalf("expected transaction committed event, got %s", received[1].Type)
	}
}
