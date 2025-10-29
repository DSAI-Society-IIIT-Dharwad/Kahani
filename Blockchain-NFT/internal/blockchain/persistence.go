package blockchain

import (
	"errors"

	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/types"
)

// BlockStateStore abstracts the persistence layer used to store blocks and state.
type BlockStateStore interface {
	SaveBlock(block types.Block) error
	GetBlock(index int) (types.Block, error)
	SaveState(state types.State) error
	GetState() (types.State, error)
}

// WithStorage attaches the provided persistence layer to the blockchain and
// synchronises all existing blocks and state to disk.
func (bc *Blockchain) WithStorage(store BlockStateStore) error {
	if store == nil {
		return errors.New("blockchain: storage is nil")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.store = store

	for _, block := range bc.blocks {
		if err := store.SaveBlock(block); err != nil {
			bc.store = nil
			return err
		}
	}

	if err := store.SaveState(cloneStateMaps(bc.walletRegistry, bc.nftRegistry)); err != nil {
		bc.store = nil
		return err
	}

	return nil
}

// LoadBlockchain reconstructs a blockchain instance from the supplied storage
// backend. When no blocks are present, a genesis block is created and persisted.
func LoadBlockchain(store BlockStateStore) (*Blockchain, error) {
	if store == nil {
		return nil, errors.New("blockchain: storage is nil")
	}

	bc := &Blockchain{
		blocks:              make([]types.Block, 0),
		walletRegistry:      make(map[string]types.Wallet),
		nftRegistry:         make(map[string]types.NFT),
		pendingTransactions: make([]types.Transaction, 0),
		store:               store,
	}

	for idx := 0; ; idx++ {
		block, err := store.GetBlock(idx)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				break
			}
			return nil, err
		}
		bc.blocks = append(bc.blocks, block)
	}

	if len(bc.blocks) == 0 {
		genesis := NewBlock(0, "", nil)
		bc.blocks = append(bc.blocks, genesis)
		if err := store.SaveBlock(genesis); err != nil {
			return nil, err
		}
		if err := store.SaveState(cloneStateMaps(bc.walletRegistry, bc.nftRegistry)); err != nil {
			return nil, err
		}
		return bc, nil
	}

	genesis := bc.blocks[0]
	if genesis.Index != 0 || genesis.PrevHash != "" || genesis.Hash != CalculateHash(genesis) {
		return nil, errors.New("blockchain: invalid genesis block in storage")
	}

	state := types.State{
		WalletRegistry: make(map[string]types.Wallet),
		NFTRegistry:    make(map[string]types.NFT),
	}

	for i := 1; i < len(bc.blocks); i++ {
		updated, err := ValidateBlock(bc.blocks[i], bc.blocks[i-1], state)
		if err != nil {
			return nil, err
		}
		state = updated
	}

	if state.WalletRegistry == nil {
		state.WalletRegistry = make(map[string]types.Wallet)
	}
	if state.NFTRegistry == nil {
		state.NFTRegistry = make(map[string]types.NFT)
	}

	bc.walletRegistry = state.WalletRegistry
	bc.nftRegistry = state.NFTRegistry

	if err := store.SaveState(cloneStateMaps(bc.walletRegistry, bc.nftRegistry)); err != nil {
		return nil, err
	}

	return bc, nil
}

func (bc *Blockchain) persistLocked(block types.Block) error {
	if bc.store == nil {
		return nil
	}

	if err := bc.store.SaveBlock(block); err != nil {
		return err
	}

	return bc.store.SaveState(cloneStateMaps(bc.walletRegistry, bc.nftRegistry))
}

func cloneStateMaps(wallets map[string]types.Wallet, nfts map[string]types.NFT) types.State {
	copiedWallets := make(map[string]types.Wallet, len(wallets))
	for k, v := range wallets {
		copiedWallets[k] = v
	}

	copiedNFTs := make(map[string]types.NFT, len(nfts))
	for k, v := range nfts {
		copiedNFTs[k] = v
	}

	return types.State{WalletRegistry: copiedWallets, NFTRegistry: copiedNFTs}
}
