package blockchain

import (
	"errors"
	"sync"
	"time"

	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

var (
	errInvalidGenesisPrev = errors.New("invalid genesis previous hash")
	errIndexOutOfSequence = errors.New("block index out of sequence")
	errPrevHashMismatch   = errors.New("previous hash mismatch")
	errHashMismatch       = errors.New("block hash mismatch")
)

// Blockchain manages the in-memory view of blocks and state.
type Blockchain struct {
	mu                  sync.RWMutex
	blocks              []types.Block
	walletRegistry      map[string]types.Wallet
	nftRegistry         map[string]types.NFT
	pendingTransactions []types.Transaction
	observer            *observer.Bus
	store               BlockStateStore
}

// NewBlockchain bootstraps a chain with a genesis block.
func NewBlockchain() *Blockchain {
	genesis := NewBlock(0, "", nil)
	return &Blockchain{
		blocks:              []types.Block{genesis},
		walletRegistry:      make(map[string]types.Wallet),
		nftRegistry:         make(map[string]types.NFT),
		pendingTransactions: make([]types.Transaction, 0),
	}
}

// SetObserver attaches the provided event bus to the blockchain.
func (bc *Blockchain) SetObserver(bus *observer.Bus) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.observer = bus
}

// LatestBlock returns the most recent block.
func (bc *Blockchain) LatestBlock() types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.blocks[len(bc.blocks)-1]
}

// Blocks returns a shallow copy of the chain.
func (bc *Blockchain) Blocks() []types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	blocks := make([]types.Block, len(bc.blocks))
	copy(blocks, bc.blocks)
	return blocks
}

// AddBlock validates and appends the block.
func (bc *Blockchain) AddBlock(block types.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.blocks) == 0 {
		return errors.New("blockchain not initialized")
	}

	if block.Index == 0 {
		if block.PrevHash != "" {
			return errInvalidGenesisPrev
		}
		if block.Hash != CalculateHash(block) {
			return errHashMismatch
		}

		previous := bc.blocks
		bc.blocks = []types.Block{block}

		if err := bc.persistLocked(block); err != nil {
			bc.blocks = previous
			return err
		}

		return nil
	}

	prev := bc.blocks[len(bc.blocks)-1]

	state := types.State{
		WalletRegistry: make(map[string]types.Wallet, len(bc.walletRegistry)),
		NFTRegistry:    make(map[string]types.NFT, len(bc.nftRegistry)),
	}

	for k, v := range bc.walletRegistry {
		state.WalletRegistry[k] = v
	}

	for k, v := range bc.nftRegistry {
		state.NFTRegistry[k] = v
	}

	updatedState, err := ValidateBlock(block, prev, state)
	if err != nil {
		return err
	}

	prevWallets := bc.walletRegistry
	prevNFTs := bc.nftRegistry

	bc.blocks = append(bc.blocks, block)
	bc.walletRegistry = updatedState.WalletRegistry
	bc.nftRegistry = updatedState.NFTRegistry

	if err := bc.persistLocked(block); err != nil {
		bc.blocks = bc.blocks[:len(bc.blocks)-1]
		bc.walletRegistry = prevWallets
		bc.nftRegistry = prevNFTs
		return err
	}
	return nil
}

// ValidateChain ensures the entire chain is internally consistent.
func (bc *Blockchain) ValidateChain() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.blocks) == 0 {
		return false
	}

	for i := 1; i < len(bc.blocks); i++ {
		curr := bc.blocks[i]
		prev := bc.blocks[i-1]

		if curr.Index != prev.Index+1 {
			return false
		}

		if curr.PrevHash != prev.Hash {
			return false
		}

		if CalculateHash(curr) != curr.Hash {
			return false
		}
	}

	return true
}

// EnqueueTransaction stages a transaction for inclusion in the next block.
func (bc *Blockchain) EnqueueTransaction(tx types.Transaction) {
	bc.mu.Lock()
	bc.pendingTransactions = append(bc.pendingTransactions, tx)
	bc.mu.Unlock()

	bc.emitEvent(observer.Event{
		Type:      observer.EventTransactionQueued,
		Timestamp: time.Now().UTC(),
		Data:      tx,
	})
}

// PendingTransactions returns a copy of currently staged transactions.
func (bc *Blockchain) PendingTransactions() []types.Transaction {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	txs := make([]types.Transaction, len(bc.pendingTransactions))
	copy(txs, bc.pendingTransactions)
	return txs
}

// ClearPendingTransactions drops all staged transactions.
func (bc *Blockchain) ClearPendingTransactions() {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.pendingTransactions = bc.pendingTransactions[:0]
}

func (bc *Blockchain) emitEvent(event observer.Event) {
	if bc == nil {
		return
	}

	if bus := bc.observer; bus != nil {
		bus.Publish(event)
	}
}

// State returns a snapshot of the chain state maps.
func (bc *Blockchain) State() types.State {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return cloneStateMaps(bc.walletRegistry, bc.nftRegistry)
}

// RegisterWallet inserts or updates a wallet in the chain state keyed by Supabase user ID.
func (bc *Blockchain) RegisterWallet(wallet types.Wallet) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.walletRegistry[wallet.SupabaseUserID] = wallet
}

// GetWalletBySupabaseID fetches the wallet for the provided Supabase user ID.
func (bc *Blockchain) GetWalletBySupabaseID(userID string) (types.Wallet, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	wallet, ok := bc.walletRegistry[userID]
	return wallet, ok
}
