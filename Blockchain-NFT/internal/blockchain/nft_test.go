package blockchain

import (
	"encoding/json"
	"testing"

	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/types"
)

func TestMintNFTSuccess(t *testing.T) {
	ipfs := storage.NewMemoryIPFS()

	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 12345 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	story := types.Story{
		ID:      "story-1",
		Title:   "Legends of the Chain",
		Summary: "A collaborative sci-fi saga.",
		Contributions: []types.Contribution{
			{ContributorID: "user-a", WalletAddress: "0xaaa", StoryID: "story-1", StoryLine: "Line 1"},
			{ContributorID: "user-a", WalletAddress: "0xaaa", StoryID: "story-1", StoryLine: "Line 2"},
			{ContributorID: "user-b", WalletAddress: "0xbbb", StoryID: "story-1", StoryLine: "Line 3"},
			{ContributorID: "user-a", WalletAddress: "0xaaa", StoryID: "story-1", StoryLine: "Line 4"},
		},
	}

	nft, err := MintNFT(story, ipfs)
	if err != nil {
		t.Fatalf("mint nft failed: %v", err)
	}

	if nft.TokenID == "" {
		t.Fatalf("expected token id to be set")
	}

	if nft.MainAuthor.SupabaseUserID != "user-a" {
		t.Fatalf("expected user-a to be main author")
	}

	if len(nft.CoAuthors) != 1 || nft.CoAuthors[0].SupabaseUserID != "user-b" {
		t.Fatalf("expected user-b to be co-author")
	}

	if nft.MintedAt != 12345 {
		t.Fatalf("expected minted timestamp override, got %d", nft.MintedAt)
	}

	if nft.Summary != "A collaborative sci-fi saga." {
		t.Fatalf("expected summary to propagate into nft")
	}

	// Validate metadata stored in IPFS.
	metadataBytes, err := ipfs.Fetch(nft.MetadataIPFSCID)
	if err != nil {
		t.Fatalf("failed to fetch metadata: %v", err)
	}

	var metadata struct {
		StoryID  string         `json:"story_id"`
		Summary  string         `json:"summary"`
		ImageCID string         `json:"image_cid"`
		Authors  []types.Author `json:"authors"`
		MintedAt int64          `json:"minted_at"`
	}

	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		t.Fatalf("failed to decode metadata: %v", err)
	}

	if metadata.StoryID != "story-1" {
		t.Fatalf("expected metadata story id")
	}

	if metadata.Summary != "A collaborative sci-fi saga." {
		t.Fatalf("expected metadata summary to match story")
	}

	if metadata.ImageCID != nft.ImageIPFSCID {
		t.Fatalf("metadata image cid mismatch")
	}

	if len(metadata.Authors) != 2 {
		t.Fatalf("expected two authors in metadata")
	}

	if metadata.MintedAt != 12345 {
		t.Fatalf("expected metadata minted at to match timestamp")
	}

	if _, err := ipfs.Fetch(nft.ImageIPFSCID); err != nil {
		t.Fatalf("expected image payload to exist: %v", err)
	}
}

func TestMintNFTValidation(t *testing.T) {
	ipfs := storage.NewMemoryIPFS()

	if _, err := MintNFT(types.Story{}, nil); err == nil {
		t.Fatalf("expected error when ipfs nil")
	}

	if _, err := MintNFT(types.Story{ID: "story"}, ipfs); err == nil {
		t.Fatalf("expected error when title missing")
	}

	if _, err := MintNFT(types.Story{ID: "story", Title: "title"}, ipfs); err == nil {
		t.Fatalf("expected error when contributions missing")
	}
}

func TestGenerateNFTImage(t *testing.T) {
	ipfs := storage.NewMemoryIPFS()

	story := types.Story{
		ID:    "story-42",
		Title: "Mystery",
		Contributions: []types.Contribution{
			{ContributorID: "user-a", WalletAddress: "0xaaa", StoryLine: "line"},
		},
	}

	cid, err := GenerateNFTImage(story, ipfs)
	if err != nil {
		t.Fatalf("generate nft image failed: %v", err)
	}

	payload, err := ipfs.Fetch(cid)
	if err != nil {
		t.Fatalf("failed to fetch image payload: %v", err)
	}

	var decoded struct {
		StoryID string         `json:"story_id"`
		Authors []types.Author `json:"authors"`
	}

	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode image payload: %v", err)
	}

	if decoded.StoryID != "story-42" {
		t.Fatalf("expected story id to propagate")
	}

	if len(decoded.Authors) != 1 {
		t.Fatalf("expected author aggregation in image payload")
	}
}
