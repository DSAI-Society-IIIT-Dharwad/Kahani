package network

import (
	"errors"

	"storytelling-blockchain/internal/types"
)

// BlockchainFetcher retrieves blocks from a peer during synchronization.
type BlockchainFetcher interface {
	Fetch(peerID string) ([]types.Block, error)
}

// DiscoverPeers adds the provided bootstrap nodes to the peer list.
func (n *Node) DiscoverPeers(bootstrap []string) {
	for _, peer := range bootstrap {
		if peer == "" || peer == n.id {
			continue
		}
		n.ConnectToPeer(peer)
	}
}

// SyncBlockchain attempts to fetch blocks from connected peers using the provided fetcher.
func (n *Node) SyncBlockchain(fetcher BlockchainFetcher) ([]types.Block, error) {
	if fetcher == nil {
		return nil, errors.New("network: blockchain fetcher not provided")
	}

	for _, peer := range n.Peers() {
		blocks, err := fetcher.Fetch(peer)
		if err != nil {
			continue
		}
		if len(blocks) > 0 {
			return blocks, nil
		}
	}

	return nil, errors.New("network: no blocks available from peers")
}
