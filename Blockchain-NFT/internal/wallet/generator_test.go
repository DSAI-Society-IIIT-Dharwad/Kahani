package wallet

import (
	"testing"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/types"
)

func TestGenerateWalletForUser(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 999 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	generator, err := NewGenerator("supasecret")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	wallet, err := generator.GenerateWalletForUser("user-123")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	if wallet.SupabaseUserID != "user-123" {
		t.Fatalf("expected user id to propagate")
	}

	if wallet.CreatedAt != 999 {
		t.Fatalf("expected created timestamp override, got %d", wallet.CreatedAt)
	}

	expectedAddress := deriveWalletAddress("user-123")
	if wallet.Address != expectedAddress {
		t.Fatalf("expected deterministic address %s got %s", expectedAddress, wallet.Address)
	}

	chain := blockchain.NewBlockchain()
	manager, err := NewManager(chain, "supasecret")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	plain, err := manager.decryptPrivateKey(wallet.PrivateKeyEncrypted)
	if err != nil {
		t.Fatalf("failed to decrypt private key: %v", err)
	}

	if len(plain) == 0 {
		t.Fatalf("expected decrypted private key to contain data")
	}
}
