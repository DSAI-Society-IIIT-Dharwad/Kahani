package wallet

import (
	"encoding/json"
	"errors"
	"fmt"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/consensus/sharding"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

// Storage handles persisting wallet transactions onto the blockchain queue.
type Storage struct {
	chain          *blockchain.Blockchain
	proposer       consensusProposer
	consensusNodes []string
}

type consensusProposer interface {
	Propose(nodeID string, txs []types.Transaction) error
}

// NewStorage constructs a wallet storage layer backed by a blockchain instance.
func NewStorage(chain *blockchain.Blockchain) (*Storage, error) {
	if chain == nil {
		return nil, errors.New("wallet: blockchain reference is nil")
	}
	return &Storage{chain: chain}, nil
}

// WithConsensus registers a consensus proposer and optional sharded node set.
func (s *Storage) WithConsensus(nodeID string, proposer consensusProposer, additional ...string) {
	s.consensusNodes = append([]string{}, additional...)
	if nodeID != "" {
		s.consensusNodes = append([]string{nodeID}, s.consensusNodes...)
	}
	s.proposer = proposer
}

// StoreWalletOnChain creates a create_wallet transaction and registers the wallet in state.
func (s *Storage) StoreWalletOnChain(wallet types.Wallet) (types.Transaction, error) {
	if wallet.SupabaseUserID == "" {
		return types.Transaction{}, errors.New("wallet: missing supabase user id")
	}

	timestamp := types.NowUnix()
	payload, err := json.Marshal(struct {
		Wallet    types.Wallet `json:"wallet"`
		Timestamp int64        `json:"timestamp"`
	}{Wallet: wallet, Timestamp: timestamp})
	if err != nil {
		return types.Transaction{}, fmt.Errorf("wallet: marshal wallet failed: %w", err)
	}

	tx := types.Transaction{
		Type:      "create_wallet",
		Data:      wallet,
		Timestamp: timestamp,
	}
	tx.TxID = utils.ComputeSHA256(payload)

	s.chain.RegisterWallet(wallet)
	s.chain.EnqueueTransaction(tx)

	targetNode := s.selectConsensusNode(wallet.SupabaseUserID)

	if s.proposer != nil && targetNode != "" {
		if err := s.proposer.Propose(targetNode, nil); err != nil {
			return types.Transaction{}, err
		}
	}

	return tx, nil
}

func (s *Storage) selectConsensusNode(key string) string {
	if len(s.consensusNodes) == 0 {
		return ""
	}

	if key == "" {
		return s.consensusNodes[0]
	}

	return sharding.SelectNode(s.consensusNodes, key)
}
