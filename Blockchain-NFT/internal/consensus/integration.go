package consensus

import (
	"errors"
	"time"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

// ChainBlockBuilder constructs blocks using the current blockchain head.
type ChainBlockBuilder struct {
	Chain *blockchain.Blockchain
}

// BuildBlock assembles a new block from the supplied transactions.
func (b *ChainBlockBuilder) BuildBlock(transactions []types.Transaction) (types.Block, error) {
	if b == nil || b.Chain == nil {
		return types.Block{}, errors.New("consensus: blockchain is required to build blocks")
	}

	if len(transactions) == 0 {
		return types.Block{}, errors.New("consensus: transactions required to build block")
	}

	prev := b.Chain.LatestBlock()
	block := blockchain.NewBlock(prev.Index+1, prev.Hash, transactions)
	return block, nil
}

// NewChainFinalizer returns a Finalizer that commits blocks to the blockchain and emits events.
func NewChainFinalizer(chain *blockchain.Blockchain, bus *observer.Bus) Finalizer {
	return func(block types.Block) {
		if chain == nil {
			return
		}
		latest := chain.LatestBlock()
		if block.Index <= latest.Index {
			// Skip already committed blocks for shared in-process runtimes.
			if block.Index == latest.Index && block.Hash == latest.Hash {
				return
			}
		}

		if err := chain.AddBlock(block); err != nil {
			if bus != nil {
				bus.Publish(observer.Event{
					Type:      observer.EventError,
					Timestamp: time.Now().UTC(),
					Data: map[string]string{
						"message": err.Error(),
					},
				})
			}
			return
		}

		chain.ClearPendingTransactions()

		if bus != nil {
			bus.Publish(observer.Event{
				Type:      observer.EventBlockCommitted,
				Timestamp: time.Now().UTC(),
				Data:      block,
			})

			for _, tx := range block.Transactions {
				bus.Publish(observer.Event{
					Type:      observer.EventTransactionCommitted,
					Timestamp: time.Now().UTC(),
					Data:      tx,
				})
			}
		}
	}
}
