package storage_test

import (
	"errors"
	"testing"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/types"
)

func TestBadgerStorageBlockRoundTrip(t *testing.T) {
	bs, err := storage.NewBadgerStorage(storage.BadgerConfig{InMemory: true})
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	t.Cleanup(func() {
		_ = bs.Close()
	})

	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 1111 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	block := blockchain.NewBlock(0, "", nil)
	if err := bs.SaveBlock(block); err != nil {
		t.Fatalf("failed to save block: %v", err)
	}

	retrieved, err := bs.GetBlock(0)
	if err != nil {
		t.Fatalf("failed to retrieve block: %v", err)
	}

	if retrieved.Hash != block.Hash {
		t.Fatalf("expected hash %s got %s", block.Hash, retrieved.Hash)
	}
}

func TestBadgerStorageStateRoundTrip(t *testing.T) {
	bs, err := storage.NewBadgerStorage(storage.BadgerConfig{InMemory: true})
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	t.Cleanup(func() {
		_ = bs.Close()
	})

	_, err = bs.GetState()
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("get state failed: %v", err)
	}

	state := types.State{
		WalletRegistry: map[string]types.Wallet{
			"user": {
				Address:        "addr",
				SupabaseUserID: "user",
			},
		},
		NFTRegistry: map[string]types.NFT{
			"token": {
				TokenID: "token",
			},
		},
	}

	if err := bs.SaveState(state); err != nil {
		t.Fatalf("save state failed: %v", err)
	}

	restored, err := bs.GetState()
	if err != nil {
		t.Fatalf("get state failed: %v", err)
	}

	if len(restored.WalletRegistry) != 1 || len(restored.NFTRegistry) != 1 {
		t.Fatalf("unexpected state contents: %+v", restored)
	}
}
