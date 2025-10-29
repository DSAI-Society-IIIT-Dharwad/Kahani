package wallet

import (
	"encoding/json"
	"testing"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

func TestManagerGetAndSign(t *testing.T) {
	chain := blockchain.NewBlockchain()
	generator, err := NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("generator init failed: %v", err)
	}

	wallet, err := generator.GenerateWalletForUser("user-xyz")
	if err != nil {
		t.Fatalf("wallet generation failed: %v", err)
	}

	chain.RegisterWallet(wallet)

	manager, err := NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("manager init failed: %v", err)
	}

	retrieved, ok := manager.GetWalletBySupabaseID("user-xyz")
	if !ok {
		t.Fatalf("expected wallet to be retrievable")
	}

	if retrieved.Address != wallet.Address {
		t.Fatalf("retrieved wallet mismatch")
	}

	contribution := types.Contribution{
		ContributorID: "user-xyz",
		WalletAddress: retrieved.Address,
		StoryID:       "story-1",
		StoryLine:     "Once upon a time",
		Timestamp:     555,
	}

	signature, err := manager.SignContribution(retrieved, contribution)
	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}

	payload, _ := json.Marshal(contribution)
	valid, err := utils.VerifyEd25519(retrieved.PublicKey, payload, signature)
	if err != nil {
		t.Fatalf("verification error: %v", err)
	}

	if !valid {
		t.Fatalf("expected contribution signature to verify")
	}

	if _, err := manager.SignContribution(types.Wallet{}, contribution); err == nil {
		t.Fatalf("expected error when encrypted key missing")
	}
}
