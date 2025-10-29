package types

import "time"

// Block represents a single block in the blockchain.
type Block struct {
	Index               int               `json:"index"`
	Timestamp           int64             `json:"timestamp"`
	Transactions        []Transaction     `json:"transactions"`
	PrevHash            string            `json:"prev_hash"`
	Hash                string            `json:"hash"`
	ValidatorSignatures map[string]string `json:"validator_signatures"`
	Nonce               int               `json:"nonce"`
}

// Transaction captures the actions recorded on-chain.
type Transaction struct {
	TxID      string      `json:"tx_id"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	Signature string      `json:"signature"`
}

// Contribution holds information for a single story contribution.
type Contribution struct {
	ContributorID string `json:"contributor_id"`
	WalletAddress string `json:"wallet_address"`
	StoryID       string `json:"story_id"`
	StoryLine     string `json:"story_line"`
	Timestamp     int64  `json:"timestamp"`
}

// NFT represents the minted storytelling NFT metadata.
type NFT struct {
	TokenID         string   `json:"token_id"`
	StoryID         string   `json:"story_id"`
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	MainAuthor      Author   `json:"main_author"`
	CoAuthors       []Author `json:"co_authors"`
	ImageIPFSCID    string   `json:"image_ipfs_cid"`
	MetadataIPFSCID string   `json:"metadata_ipfs_cid"`
	MintedAt        int64    `json:"minted_at"`
	BlockIndex      int      `json:"block_index"`
}

// Author represents a collaborative writer on a story.
type Author struct {
	SupabaseUserID                 string         `json:"supabase_user_id"`
	WalletAddress                  string         `json:"wallet_address"`
	ContributionCount              int            `json:"contribution_count"`
	OwnershipPercentage            float64        `json:"ownership_percentage"`
	GithubCommitStyleContributions []Contribution `json:"github_commit_style_contributions"`
}

// Wallet encapsulates a Supabase user's wallet details.
type Wallet struct {
	Address             string `json:"address"`
	SupabaseUserID      string `json:"supabase_user_id"`
	PublicKey           string `json:"public_key"`
	PrivateKeyEncrypted string `json:"private_key_encrypted"`
	CreatedAt           int64  `json:"created_at"`
	BlockIndex          int    `json:"block_index"`
}

// Story represents the aggregate context used when minting NFTs.
type Story struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Summary       string         `json:"summary"`
	Contributions []Contribution `json:"contributions"`
}

// State aggregates the on-chain registries required for querying.
type State struct {
	WalletRegistry map[string]Wallet `json:"wallet_registry"`
	NFTRegistry    map[string]NFT    `json:"nft_registry"`
}

// NowUnix returns the current unix timestamp to aid testing hooks.
var NowUnix = func() int64 { return time.Now().Unix() }
