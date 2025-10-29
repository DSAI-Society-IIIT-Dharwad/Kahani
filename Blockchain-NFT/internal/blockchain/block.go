package blockchain

import (
	"encoding/json"

	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

// NewBlock builds a block ready to append to the chain.
func NewBlock(index int, prevHash string, transactions []types.Transaction) types.Block {
	block := types.Block{
		Index:               index,
		Timestamp:           types.NowUnix(),
		Transactions:        transactions,
		PrevHash:            prevHash,
		ValidatorSignatures: make(map[string]string),
	}
	block.Hash = CalculateHash(block)
	return block
}

// CalculateHash deterministically hashes the block fields.
func CalculateHash(block types.Block) string {
	clone := block
	clone.Hash = ""

	for i := range clone.Transactions {
		if clone.Transactions[i].Data == nil {
			continue
		}

		canonical, err := canonicalizeJSON(clone.Transactions[i].Data)
		if err != nil {
			return ""
		}

		clone.Transactions[i].Data = canonical
	}

	payload, err := json.Marshal(clone)
	if err != nil {
		return ""
	}

	return utils.ComputeSHA256(payload)
}

// canonicalizeJSON normalizes JSON-compatible values to avoid hash drift after reloading from storage.
func canonicalizeJSON(value interface{}) (interface{}, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var decoded interface{}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}
