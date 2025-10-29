package wallet

import (
	"encoding/json"
	"errors"
	"fmt"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

// Manager exposes wallet retrieval and signing utilities against the blockchain state.
type Manager struct {
	chain  *blockchain.Blockchain
	aesKey []byte
}

// NewManager creates a wallet manager sharing the encryption key with the generator.
func NewManager(chain *blockchain.Blockchain, passphrase string) (*Manager, error) {
	if chain == nil {
		return nil, errors.New("wallet: blockchain reference is nil")
	}

	aesKey := deriveAESKey(passphrase)
	if len(aesKey) != 32 {
		return nil, errors.New("wallet: invalid AES key length")
	}

	return &Manager{chain: chain, aesKey: aesKey}, nil
}

// GetWalletBySupabaseID returns the wallet stored for the specified Supabase user.
func (m *Manager) GetWalletBySupabaseID(userID string) (types.Wallet, bool) {
	return m.chain.GetWalletBySupabaseID(userID)
}

// SignContribution signs the contribution payload using the wallet's private key.
func (m *Manager) SignContribution(wallet types.Wallet, contribution types.Contribution) (string, error) {
	if wallet.PrivateKeyEncrypted == "" {
		return "", errors.New("wallet: encrypted private key missing")
	}

	plainPrivKey, err := decryptString(m.aesKey, wallet.PrivateKeyEncrypted)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(contribution)
	if err != nil {
		return "", fmt.Errorf("wallet: contribution marshal failed: %w", err)
	}

	signature, err := utils.SignEd25519(plainPrivKey, payload)
	if err != nil {
		return "", fmt.Errorf("wallet: sign contribution failed: %w", err)
	}

	return signature, nil
}

// decryptPrivateKey is exposed for testing to ensure encryption symmetry.
func (m *Manager) decryptPrivateKey(encrypted string) (string, error) {
	return decryptString(m.aesKey, encrypted)
}
