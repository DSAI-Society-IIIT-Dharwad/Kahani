package consensus

import (
	"errors"

	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/types"
)

// BootstrapOptions bundle parameters required to spin up consensus participants.
type BootstrapOptions struct {
	NodeID         string
	Peers          []string
	FaultTolerance int
	Transport      *network.Node
	Signer         Signer
	Builder        BlockBuilder
	Finalize       Finalizer
}

// BootstrapNode creates a PBFT node tied into the gossip system.
func BootstrapNode(opts BootstrapOptions) (*PBFTNode, network.GossipHandler, error) {
	if opts.Transport == nil {
		return nil, nil, errors.New("consensus: transport node required")
	}
	if opts.Builder == nil {
		return nil, nil, errors.New("consensus: block builder required")
	}
	if opts.Finalize == nil {
		return nil, nil, errors.New("consensus: finalize callback required")
	}

	gossipNet, err := NewGossipNetwork(opts.Transport)
	if err != nil {
		return nil, nil, err
	}

	cfg := Config{
		ID:             opts.NodeID,
		Peers:          opts.Peers,
		FaultTolerance: opts.FaultTolerance,
		Network:        gossipNet,
		Signer:         opts.Signer,
		Builder:        opts.Builder,
		Finalize:       opts.Finalize,
	}

	node, err := NewPBFTNode(cfg)
	if err != nil {
		return nil, nil, err
	}

	handler, err := NewPBFTGossipHandler(node)
	if err != nil {
		return nil, nil, err
	}

	return node, handler, nil
}

// BootstrapCluster assists in constructing a set of PBFT nodes over the network transport.
func BootstrapCluster(transports map[string]*network.Node, peers []string, builder BlockBuilder, finalizers map[string]Finalizer, signer Signer, faultTolerance int) (map[string]*PBFTNode, map[string]network.GossipHandler, error) {
	if len(transports) == 0 {
		return nil, nil, errors.New("consensus: transports required")
	}
	if builder == nil {
		return nil, nil, errors.New("consensus: block builder required")
	}
	if len(finalizers) != len(transports) {
		return nil, nil, errors.New("consensus: finalizer per node required")
	}

	nodes := make(map[string]*PBFTNode, len(transports))
	handlers := make(map[string]network.GossipHandler, len(transports))

	for id, transport := range transports {
		transport.DiscoverPeers(peers)

		fin, ok := finalizers[id]
		if !ok {
			return nil, nil, errors.New("consensus: missing finalizer for node " + id)
		}

		node, handler, err := BootstrapNode(BootstrapOptions{
			NodeID:         id,
			Peers:          peers,
			FaultTolerance: faultTolerance,
			Transport:      transport,
			Signer:         signer,
			Builder:        builder,
			Finalize:       fin,
		})
		if err != nil {
			return nil, nil, err
		}

		nodes[id] = node
		handlers[id] = handler
	}

	return nodes, handlers, nil
}

// ApplyBlocks finalizes a fetched chain by replaying the finalize callback.
func ApplyBlocks(blocks []types.Block, finalize Finalizer) error {
	if finalize == nil {
		return errors.New("consensus: finalize required")
	}

	for _, block := range blocks {
		finalize(block)
	}

	return nil
}
