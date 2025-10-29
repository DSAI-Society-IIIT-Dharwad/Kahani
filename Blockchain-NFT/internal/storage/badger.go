package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"storytelling-blockchain/internal/types"
)

// BadgerConfig controls the storage backend creation.
type BadgerConfig struct {
	Path     string
	InMemory bool
}

// BadgerStorage wraps a BadgerDB instance for block persistence.
type BadgerStorage struct {
	db *badger.DB
}

// ErrNotFound indicates the requested record does not exist in storage.
var ErrNotFound = errors.New("storage: not found")

// NewBadgerStorage opens a Badger database with the provided configuration.
func NewBadgerStorage(cfg BadgerConfig) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(cfg.Path)
	if cfg.InMemory {
		opts = opts.WithInMemory(true)
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerStorage{db: db}, nil
}

// Close releases resources held by the underlying database.
func (bs *BadgerStorage) Close() error {
	return bs.db.Close()
}

// SaveBlock persists the block keyed by its index.
func (bs *BadgerStorage) SaveBlock(block types.Block) error {
	payload, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(blockKey(block.Index)), payload)
	})
}

// GetBlock retrieves the block for the provided index.
func (bs *BadgerStorage) GetBlock(index int) (types.Block, error) {
	var block types.Block

	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(blockKey(index)))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &block)
		})
	})

	return block, err
}

// SaveState writes the aggregated chain state.
func (bs *BadgerStorage) SaveState(state types.State) error {
	payload, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(stateKey()), payload)
	})
}

// GetState retrieves the latest persisted chain state.
func (bs *BadgerStorage) GetState() (types.State, error) {
	state := types.State{
		WalletRegistry: make(map[string]types.Wallet),
		NFTRegistry:    make(map[string]types.NFT),
	}

	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(stateKey()))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &state)
		})
	})

	return state, err
}

func blockKey(index int) string {
	return fmt.Sprintf("block:%d", index)
}

func stateKey() string {
	return "state:latest"
}
