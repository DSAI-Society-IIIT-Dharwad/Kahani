package blockchain

import (
	"encoding/json"
	"testing"

	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

func TestLoadBlockchainReconstructsState(t *testing.T) {
	store, err := storage.NewBadgerStorage(storage.BadgerConfig{InMemory: true})
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 5000 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	bc := NewBlockchain()
	if err := bc.WithStorage(store); err != nil {
		t.Fatalf("attach storage failed: %v", err)
	}

	userWallet := types.Wallet{
		Address:             "0xabc123",
		SupabaseUserID:      "user-1",
		PublicKey:           "public-key",
		PrivateKeyEncrypted: "ciphertext",
		CreatedAt:           5500,
		BlockIndex:          -1,
	}

	payload := struct {
		Wallet    types.Wallet `json:"wallet"`
		Timestamp int64        `json:"timestamp"`
	}{Wallet: userWallet, Timestamp: 6000}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	tx := types.Transaction{
		TxID:      utils.ComputeSHA256(payloadBytes),
		Type:      "create_wallet",
		Data:      userWallet,
		Timestamp: payload.Timestamp,
	}

	prev := bc.LatestBlock()
	block := NewBlock(prev.Index+1, prev.Hash, []types.Transaction{tx})

	if err := bc.AddBlock(block); err != nil {
		t.Fatalf("add block failed: %v", err)
	}

	restored, err := LoadBlockchain(store)
	if err != nil {
		t.Fatalf("load blockchain failed: %v", err)
	}

	if len(restored.Blocks()) != 2 {
		t.Fatalf("expected two blocks after restore, got %d", len(restored.Blocks()))
	}

	stored, ok := restored.GetWalletBySupabaseID("user-1")
	if !ok {
		t.Fatalf("wallet not reconstructed")
	}

	if stored.Address != userWallet.Address {
		t.Fatalf("unexpected wallet data: %+v", stored)
	}

	if stored.BlockIndex != 1 {
		t.Fatalf("expected wallet block index 1, got %d", stored.BlockIndex)
	}
}
