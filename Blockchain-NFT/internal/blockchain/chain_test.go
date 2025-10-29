package blockchain

import (
	"testing"

	"storytelling-blockchain/internal/types"
)

func TestNewBlockchainCreatesGenesis(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 1000 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	bc := NewBlockchain()

	blocks := bc.Blocks()
	if len(blocks) != 1 {
		t.Fatalf("expected genesis block, got %d blocks", len(blocks))
	}

	genesis := blocks[0]
	if genesis.Index != 0 {
		t.Fatalf("expected genesis index 0, got %d", genesis.Index)
	}

	if genesis.PrevHash != "" {
		t.Fatalf("genesis previous hash should be empty")
	}

	if !bc.ValidateChain() {
		t.Fatalf("genesis chain should validate")
	}
}

func TestAddBlockAndValidate(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 2000 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	bc := NewBlockchain()
	genesis := bc.LatestBlock()

	tx := types.Transaction{TxID: "tx-1", Type: "test", Timestamp: types.NowUnix()}

	block := NewBlock(genesis.Index+1, genesis.Hash, []types.Transaction{tx})

	if err := bc.AddBlock(block); err != nil {
		t.Fatalf("unexpected error adding block: %v", err)
	}

	if !bc.ValidateChain() {
		t.Fatalf("expected valid chain after adding block")
	}

	latest := bc.LatestBlock()
	if latest.Index != 1 {
		t.Fatalf("expected latest block index 1, got %d", latest.Index)
	}

	pending := bc.PendingTransactions()
	if len(pending) != 0 {
		t.Fatalf("pending transactions should be unaffected in this phase")
	}
}

func TestPendingTransactionQueue(t *testing.T) {
	bc := NewBlockchain()

	tx := types.Transaction{TxID: "tx-queue", Type: "pending"}
	bc.EnqueueTransaction(tx)

	pending := bc.PendingTransactions()
	if len(pending) != 1 || pending[0].TxID != "tx-queue" {
		t.Fatalf("pending queue not updated as expected")
	}

	bc.ClearPendingTransactions()
	if len(bc.PendingTransactions()) != 0 {
		t.Fatalf("pending queue should be empty after clear")
	}
}

func TestRegisterWallet(t *testing.T) {
	bc := NewBlockchain()

	wallet := types.Wallet{SupabaseUserID: "user-1", Address: "0xabc"}
	bc.RegisterWallet(wallet)

	stored, ok := bc.GetWalletBySupabaseID("user-1")
	if !ok {
		t.Fatalf("expected wallet to be stored")
	}

	if stored.Address != "0xabc" {
		t.Fatalf("unexpected wallet data")
	}
}
