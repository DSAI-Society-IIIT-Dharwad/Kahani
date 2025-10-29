package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

func TestValidateBlockAppliesCreateWallet(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 5000 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	bc := NewBlockchain()
	prev := bc.LatestBlock()

	wallet := types.Wallet{
		Address:             "0xabc",
		SupabaseUserID:      "user-123",
		PublicKey:           "pub",
		PrivateKeyEncrypted: "enc",
		CreatedAt:           4000,
	}

	payload := struct {
		Wallet    types.Wallet `json:"wallet"`
		Timestamp int64        `json:"timestamp"`
	}{Wallet: wallet, Timestamp: 5000}

	txBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	tx := types.Transaction{
		TxID:      utils.ComputeSHA256(txBytes),
		Type:      "create_wallet",
		Data:      wallet,
		Timestamp: payload.Timestamp,
	}

	block := NewBlock(prev.Index+1, prev.Hash, []types.Transaction{tx})

	state := types.State{WalletRegistry: map[string]types.Wallet{}, NFTRegistry: map[string]types.NFT{}}

	updated, err := ValidateBlock(block, prev, state)
	if err != nil {
		t.Fatalf("validate block: %v", err)
	}

	stored, ok := updated.WalletRegistry[wallet.SupabaseUserID]
	if !ok {
		t.Fatalf("wallet not inserted into state")
	}

	if stored.BlockIndex != block.Index {
		t.Fatalf("expected wallet block index %d, got %d", block.Index, stored.BlockIndex)
	}
}

func TestValidateBlockRejectsInvalidSignature(t *testing.T) {
	bc := NewBlockchain()
	prev := bc.LatestBlock()

	pub, priv, err := utils.GenerateEd25519Keypair()
	if err != nil {
		t.Fatalf("generate keypair: %v", err)
	}

	wallet := types.Wallet{
		Address:             "0xabc",
		SupabaseUserID:      "user-123",
		PublicKey:           pub,
		PrivateKeyEncrypted: "enc",
	}

	state := types.State{
		WalletRegistry: map[string]types.Wallet{"user-123": wallet},
		NFTRegistry:    map[string]types.NFT{},
	}

	payload := contributionPayload{
		Contribution: types.Contribution{
			ContributorID: "user-123",
			WalletAddress: "0xabc",
			StoryID:       "story-1",
			StoryLine:     "line",
			Timestamp:     6000,
		},
		Timestamp: 6000,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	validSig, err := utils.SignEd25519(priv, payloadBytes)
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}

	// Deliberately corrupt the signature while preserving base64 structure.
	decoded, err := base64.StdEncoding.DecodeString(validSig)
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}

	decoded[0] ^= 0xFF

	invalidSig := base64.StdEncoding.EncodeToString(decoded)

	tx := types.Transaction{
		TxID:      utils.ComputeSHA256(payloadBytes),
		Type:      "contribution",
		Data:      payload,
		Timestamp: payload.Timestamp,
		Signature: invalidSig,
	}

	block := NewBlock(prev.Index+1, prev.Hash, []types.Transaction{tx})

	if _, err := ValidateBlock(block, prev, state); err == nil {
		t.Fatalf("expected validation to fail for invalid signature")
	}
}
