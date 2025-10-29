package wallet

import (
	"errors"
	"fmt"

	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

// Generator is responsible for constructing wallets for Supabase users.
type Generator struct {
	aesKey []byte
}

// NewGenerator creates a wallet generator using the provided passphrase.
func NewGenerator(passphrase string) (*Generator, error) {
	key := deriveAESKey(passphrase)
	if len(key) != 32 {
		return nil, errors.New("wallet: invalid AES key length")
	}
	return &Generator{aesKey: key}, nil
}

// GenerateWalletForUser creates a new wallet using Ed25519 keys and AES-GCM encryption.
func (g *Generator) GenerateWalletForUser(supabaseUserID string) (types.Wallet, error) {
	if supabaseUserID == "" {
		return types.Wallet{}, errors.New("wallet: empty supabase user id")
	}

	pubKey, privKey, err := utils.GenerateEd25519Keypair()
	if err != nil {
		return types.Wallet{}, fmt.Errorf("wallet: keypair generation failed: %w", err)
	}

	encryptedPrivKey, err := encryptString(g.aesKey, privKey)
	if err != nil {
		return types.Wallet{}, err
	}

	address := deriveWalletAddress(supabaseUserID)

	wallet := types.Wallet{
		Address:             address,
		SupabaseUserID:      supabaseUserID,
		PublicKey:           pubKey,
		PrivateKeyEncrypted: encryptedPrivKey,
		CreatedAt:           types.NowUnix(),
		BlockIndex:          -1,
	}

	return wallet, nil
}

func deriveWalletAddress(supabaseUserID string) string {
	hash := utils.ComputeSHA256([]byte(supabaseUserID))
	// Keep address readable while deterministic.
	return "0x" + hash[:40]
}
