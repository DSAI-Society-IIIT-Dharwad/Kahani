package blockchain

import (
	"testing"
	"time"

	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

func TestEnqueueTransactionEmitsEvent(t *testing.T) {
	chain := NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)

	id, ch := bus.Subscribe(1)
	t.Cleanup(func() { bus.Unsubscribe(id) })

	tx := types.Transaction{TxID: "tx-1", Type: "test"}
	chain.EnqueueTransaction(tx)

	select {
	case ev := <-ch:
		if ev.Type != observer.EventTransactionQueued {
			t.Fatalf("expected transaction queued event, got %s", ev.Type)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for event")
	}
}
