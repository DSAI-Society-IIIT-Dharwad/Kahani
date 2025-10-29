package blockchain

import (
	"encoding/json"

	"storytelling-blockchain/internal/types"
)

// StoryContributions returns all contribution transactions matching the story ID.
func (bc *Blockchain) StoryContributions(storyID string) []types.Contribution {
	if storyID == "" {
		return nil
	}

	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var results []types.Contribution

	envelope := struct {
		Contribution types.Contribution `json:"contribution"`
		Timestamp    int64              `json:"timestamp"`
	}{}

	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			if tx.Type != "contribution" {
				continue
			}

			payload, err := json.Marshal(tx.Data)
			if err != nil {
				continue
			}

			if err := json.Unmarshal(payload, &envelope); err != nil {
				continue
			}

			if envelope.Contribution.StoryID != storyID {
				continue
			}

			results = append(results, envelope.Contribution)
		}
	}

	return results
}

// GetNFT retrieves the NFT for the provided token ID from the chain state.
func (bc *Blockchain) GetNFT(tokenID string) (types.NFT, bool) {
	if tokenID == "" {
		return types.NFT{}, false
	}

	bc.mu.RLock()
	defer bc.mu.RUnlock()

	nft, ok := bc.nftRegistry[tokenID]
	return nft, ok
}

// NFTsByStory returns all NFTs minted for the given story ID.
func (bc *Blockchain) NFTsByStory(storyID string) []types.NFT {
	if storyID == "" {
		return nil
	}

	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var results []types.NFT
	for _, nft := range bc.nftRegistry {
		if nft.StoryID == storyID {
			results = append(results, nft)
		}
	}

	return results
}
