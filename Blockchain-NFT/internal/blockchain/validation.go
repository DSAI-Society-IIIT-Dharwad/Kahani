package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"

	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

var (
	errEmptyTransactionType = errors.New("blockchain: transaction type required")
	errMissingWalletID      = errors.New("blockchain: wallet supabase user id required")
	errMissingWalletAddress = errors.New("blockchain: wallet address required")
	errMissingWalletKeys    = errors.New("blockchain: wallet keys required")
	errDuplicateWallet      = errors.New("blockchain: wallet already exists with different details")
	errMissingSignature     = errors.New("blockchain: signature required")
	errMissingWallet        = errors.New("blockchain: wallet not registered")
	errInvalidSignature     = errors.New("blockchain: signature verification failed")
	errDuplicateToken       = errors.New("blockchain: nft token already exists")
)

type contributionPayload struct {
	Contribution types.Contribution `json:"contribution"`
	Timestamp    int64              `json:"timestamp"`
}

// ValidateBlock ensures the block links to its predecessor and that all
// transactions are valid with respect to the provided chain state. A mutated
// copy of the resulting state is returned for application by the caller.
func ValidateBlock(block types.Block, prev types.Block, state types.State) (types.State, error) {
	if block.Index != prev.Index+1 {
		return state, fmt.Errorf("blockchain: expected block index %d, got %d", prev.Index+1, block.Index)
	}

	if block.PrevHash != prev.Hash {
		return state, errors.New("blockchain: previous hash mismatch")
	}

	if CalculateHash(block) != block.Hash {
		return state, errors.New("blockchain: block hash mismatch")
	}

	if len(block.Transactions) == 0 {
		return state, errors.New("blockchain: block must contain transactions")
	}

	nextState := cloneState(state)

	for _, tx := range block.Transactions {
		if err := applyTransaction(&nextState, tx, block.Index); err != nil {
			return state, fmt.Errorf("blockchain: transaction %s invalid: %w", tx.TxID, err)
		}
	}

	return nextState, nil
}

func applyTransaction(state *types.State, tx types.Transaction, blockIndex int) error {
	if tx.Type == "" {
		return errEmptyTransactionType
	}

	switch tx.Type {
	case "create_wallet":
		var wallet types.Wallet
		if err := decodePayload(tx.Data, &wallet); err != nil {
			return err
		}

		if tx.Timestamp <= 0 {
			return errors.New("blockchain: transaction timestamp required")
		}

		if wallet.SupabaseUserID == "" {
			return errMissingWalletID
		}
		if wallet.Address == "" {
			return errMissingWalletAddress
		}
		if wallet.PublicKey == "" || wallet.PrivateKeyEncrypted == "" {
			return errMissingWalletKeys
		}

		hashPayload := struct {
			Wallet    types.Wallet `json:"wallet"`
			Timestamp int64        `json:"timestamp"`
		}{Wallet: wallet, Timestamp: tx.Timestamp}

		if err := verifyTxID(tx.TxID, hashPayload); err != nil {
			return err
		}

		existing, exists := state.WalletRegistry[wallet.SupabaseUserID]
		if exists {
			if existing.Address != wallet.Address || existing.PublicKey != wallet.PublicKey {
				return errDuplicateWallet
			}
		}

		wallet.BlockIndex = blockIndex
		state.WalletRegistry[wallet.SupabaseUserID] = wallet

	case "contribution":
		if tx.Signature == "" {
			return errMissingSignature
		}

		var payload contributionPayload
		if err := decodePayload(tx.Data, &payload); err != nil {
			return err
		}

		if tx.Timestamp <= 0 {
			return errors.New("blockchain: transaction timestamp required")
		}

		if payload.Contribution.ContributorID == "" {
			return errMissingWalletID
		}

		wallet, ok := state.WalletRegistry[payload.Contribution.ContributorID]
		if !ok {
			return errMissingWallet
		}

		if wallet.Address != "" && wallet.Address != payload.Contribution.WalletAddress {
			return errors.New("blockchain: contribution wallet mismatch")
		}

		if payload.Timestamp != tx.Timestamp {
			return errors.New("blockchain: contribution timestamp mismatch")
		}

		if err := verifyTxID(tx.TxID, payload); err != nil {
			return err
		}

		signedBytes, err := json.Marshal(payload.Contribution)
		if err != nil {
			return err
		}

		okSig, err := utils.VerifyEd25519(wallet.PublicKey, signedBytes, tx.Signature)
		if err != nil {
			return err
		}
		if !okSig {
			return errInvalidSignature
		}

	case "mint_nft":
		var nft types.NFT
		if err := decodePayload(tx.Data, &nft); err != nil {
			return err
		}

		if tx.Timestamp <= 0 {
			return errors.New("blockchain: transaction timestamp required")
		}

		if nft.TokenID == "" {
			return errors.New("blockchain: nft token id required")
		}

		if err := verifyTxID(tx.TxID, nft); err != nil {
			return err
		}

		if _, exists := state.NFTRegistry[nft.TokenID]; exists {
			return errDuplicateToken
		}

		nft.BlockIndex = blockIndex
		state.NFTRegistry[nft.TokenID] = nft

	default:
		// Unknown transaction types are accepted without additional validation.
	}

	return nil
}

func decodePayload(src interface{}, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dst)
}

func verifyTxID(txID string, payload interface{}) error {
	if txID == "" {
		return errors.New("blockchain: transaction id required")
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	expected := utils.ComputeSHA256(bytes)
	if expected != txID {
		return errors.New("blockchain: transaction id mismatch")
	}

	return nil
}

func cloneState(state types.State) types.State {
	cloned := types.State{
		WalletRegistry: make(map[string]types.Wallet, len(state.WalletRegistry)),
		NFTRegistry:    make(map[string]types.NFT, len(state.NFTRegistry)),
	}

	for k, v := range state.WalletRegistry {
		cloned.WalletRegistry[k] = v
	}

	for k, v := range state.NFTRegistry {
		cloned.NFTRegistry[k] = v
	}

	return cloned
}
