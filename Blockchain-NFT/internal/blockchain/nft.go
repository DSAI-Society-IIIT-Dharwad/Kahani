package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/pkg/utils"
)

var (
	errNoContributions = errors.New("nft: contributions required")
	errMissingStoryID  = errors.New("nft: story id required")
	errMissingTitle    = errors.New("nft: title required")
	errNilIPFSClient   = errors.New("nft: ipfs client is nil")
)

// Exported error aliases for consumers outside the blockchain package.
var (
	ErrNoContributions = errNoContributions
	ErrMissingStoryID  = errMissingStoryID
	ErrMissingTitle    = errMissingTitle
	ErrNilIPFSClient   = errNilIPFSClient
)

// MintNFT aggregates contributions and uploads metadata to IPFS, returning an NFT struct.
func MintNFT(story types.Story, ipfs storage.IPFSClient) (types.NFT, error) {
	if ipfs == nil {
		return types.NFT{}, errNilIPFSClient
	}

	if story.ID == "" {
		return types.NFT{}, errMissingStoryID
	}

	if story.Title == "" {
		return types.NFT{}, errMissingTitle
	}

	contributions := story.Contributions
	if len(contributions) == 0 {
		return types.NFT{}, errNoContributions
	}

	authors := AggregateAuthors(contributions)
	if len(authors) == 0 {
		return types.NFT{}, errNoContributions
	}

	imageCID, err := GenerateNFTImage(story, ipfs)
	if err != nil {
		return types.NFT{}, err
	}

	metadataCID, err := uploadMetadata(story, authors, imageCID, ipfs)
	if err != nil {
		return types.NFT{}, err
	}

	tokenID := fmt.Sprintf("nft_%s_%s", story.ID, utils.ComputeSHA256([]byte(metadataCID))[:12])
	mintedAt := types.NowUnix()

	nft := types.NFT{
		TokenID:         tokenID,
		StoryID:         story.ID,
		Title:           story.Title,
		Summary:         story.Summary,
		MainAuthor:      authors[0],
		CoAuthors:       authors[1:],
		ImageIPFSCID:    imageCID,
		MetadataIPFSCID: metadataCID,
		MintedAt:        mintedAt,
		BlockIndex:      -1,
	}

	return nft, nil
}

// GenerateNFTImage creates a simple visualization payload and pushes it to IPFS.
func GenerateNFTImage(story types.Story, ipfs storage.IPFSClient) (string, error) {
	if ipfs == nil {
		return "", errNilIPFSClient
	}

	if story.ID == "" || story.Title == "" {
		return "", errMissingStoryID
	}

	authors := AggregateAuthors(story.Contributions)
	imagePayload := struct {
		StoryID string               `json:"story_id"`
		Title   string               `json:"title"`
		Summary string               `json:"summary"`
		Authors []types.Author       `json:"authors"`
		Lines   []types.Contribution `json:"lines"`
	}{
		StoryID: story.ID,
		Title:   story.Title,
		Summary: story.Summary,
		Authors: authors,
		Lines:   story.Contributions,
	}

	bytes, err := json.MarshalIndent(imagePayload, "", "  ")
	if err != nil {
		return "", err
	}

	cid, err := ipfs.UploadBytes(bytes)
	if err != nil {
		return "", err
	}

	return cid, nil
}

func uploadMetadata(story types.Story, authors []types.Author, imageCID string, ipfs storage.IPFSClient) (string, error) {
	metadata := struct {
		StoryID       string               `json:"story_id"`
		Title         string               `json:"title"`
		Summary       string               `json:"summary"`
		ImageCID      string               `json:"image_cid"`
		Authors       []types.Author       `json:"authors"`
		Contributions []types.Contribution `json:"contributions"`
		MintedAt      int64                `json:"minted_at"`
		TokenHint     string               `json:"token_hint"`
	}{
		StoryID:       story.ID,
		Title:         story.Title,
		Summary:       story.Summary,
		ImageCID:      imageCID,
		Authors:       authors,
		Contributions: story.Contributions,
		MintedAt:      types.NowUnix(),
		TokenHint:     utils.ComputeSHA256([]byte(story.ID + story.Title))[:16],
	}

	return ipfs.UploadJSON(metadata)
}

func AggregateAuthors(contributions []types.Contribution) []types.Author {
	if len(contributions) == 0 {
		return nil
	}

	total := len(contributions)
	counts := make(map[string]*types.Author)

	for _, c := range contributions {
		author, ok := counts[c.ContributorID]
		if !ok {
			author = &types.Author{
				SupabaseUserID: c.ContributorID,
				WalletAddress:  c.WalletAddress,
			}
			counts[c.ContributorID] = author
		}

		author.ContributionCount++
		author.GithubCommitStyleContributions = append(author.GithubCommitStyleContributions, c)
	}

	authors := make([]types.Author, 0, len(counts))
	for _, author := range counts {
		author.OwnershipPercentage = float64(author.ContributionCount) / float64(total) * 100
		authors = append(authors, *author)
	}

	sort.Slice(authors, func(i, j int) bool {
		if authors[i].ContributionCount == authors[j].ContributionCount {
			return authors[i].SupabaseUserID < authors[j].SupabaseUserID
		}
		return authors[i].ContributionCount > authors[j].ContributionCount
	})

	return authors
}
