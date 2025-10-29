package app

import (
	"context"
	"errors"

	"storytelling-blockchain/internal/api"
	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/consensus"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/internal/wallet"
)

// proposer is the shared interface implemented by consensus services.
type proposer interface {
	Propose(nodeID string, txs []types.Transaction) error
}

// AttachConsensus wires the provided proposer into the API and wallet storage layers.
func AttachConsensus(apiServer *api.API, storage *wallet.Storage, nodeID string, p proposer, additional ...string) {
	if p == nil {
		return
	}

	if apiServer != nil {
		apiServer.WithConsensus(nodeID, p, additional...)
	}

	if storage != nil {
		storage.WithConsensus(nodeID, p, additional...)
	}
}

// ConsensusBootstrapConfig describes the dependencies required to start and wire consensus.
type ConsensusBootstrapConfig struct {
	Context        context.Context
	Chain          *blockchain.Blockchain
	Observer       *observer.Bus
	Storage        *wallet.Storage
	API            *api.API
	NodeID         string
	Transports     map[string]*network.Node
	Peers          []string
	Signer         consensus.Signer
	FaultTolerance int
	ConsensusNodes []string
}

// BootstrapConsensus starts the consensus service and attaches it to the provided components.
func BootstrapConsensus(cfg ConsensusBootstrapConfig) (*consensus.Service, error) {
	if cfg.Chain == nil {
		return nil, errors.New("app: blockchain chain is required")
	}

	if cfg.Storage == nil && cfg.API == nil {
		return nil, errors.New("app: api or wallet storage must be provided")
	}

	if cfg.NodeID == "" {
		return nil, errors.New("app: consensus node id is required")
	}

	if len(cfg.Transports) == 0 {
		return nil, errors.New("app: consensus transports are required")
	}

	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}

	service, err := consensus.StartService(ctx, cfg.Chain, cfg.Observer, cfg.Transports, cfg.Peers, cfg.Signer, cfg.FaultTolerance)
	if err != nil {
		return nil, err
	}

	nodes := deriveConsensusNodes(cfg)
	AttachConsensus(cfg.API, cfg.Storage, cfg.NodeID, service, nodes...)
	return service, nil
}

func deriveConsensusNodes(cfg ConsensusBootstrapConfig) []string {
	if len(cfg.ConsensusNodes) > 0 {
		nodes := make([]string, 0, len(cfg.ConsensusNodes))
		for _, id := range cfg.ConsensusNodes {
			if id == "" || id == cfg.NodeID {
				continue
			}
			nodes = append(nodes, id)
		}
		return nodes
	}

	nodes := make([]string, 0, len(cfg.Transports))
	for id := range cfg.Transports {
		if id == cfg.NodeID {
			continue
		}
		nodes = append(nodes, id)
	}

	return nodes
}
