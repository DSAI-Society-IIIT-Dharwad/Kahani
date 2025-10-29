package blockchain

import (
	"testing"

	"storytelling-blockchain/internal/types"
)

func TestNewBlock(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 1234567890 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	txs := []types.Transaction{{TxID: "tx-1", Type: "test"}}
	block := NewBlock(1, "prev", txs)

	if block.Index != 1 {
		t.Fatalf("expected index 1, got %d", block.Index)
	}

	if block.Timestamp != 1234567890 {
		t.Fatalf("expected timestamp override to apply")
	}

	if block.PrevHash != "prev" {
		t.Fatalf("expected prev hash to match input")
	}

	if block.Hash == "" {
		t.Fatalf("expected hash to be set")
	}

	expected := CalculateHash(block)
	if expected != block.Hash {
		t.Fatalf("hash mismatch: expected %s got %s", expected, block.Hash)
	}
}
