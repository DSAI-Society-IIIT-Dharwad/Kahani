package consensus

import (
	"context"
	"errors"
	"fmt"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

// Service orchestrates PBFT runtimes paired with the shared blockchain state.
type Service struct {
	chain    *blockchain.Blockchain
	bus      *observer.Bus
	runtimes map[string]*NodeRuntime
}

// StartService boots a PBFT cluster wired to the provided blockchain and bus.
func StartService(ctx context.Context, chain *blockchain.Blockchain, bus *observer.Bus, transports map[string]*network.Node, peers []string, signer Signer, faultTolerance int) (*Service, error) {
	if chain == nil {
		return nil, errors.New("consensus: blockchain required")
	}
	if len(transports) == 0 {
		return nil, errors.New("consensus: transports required")
	}

	builder := &ChainBlockBuilder{Chain: chain}
	finalizer := NewChainFinalizer(chain, bus)

	finals := make(map[string]Finalizer, len(transports))
	for id := range transports {
		finals[id] = finalizer
	}

	runtimes, err := StartCluster(ctx, transports, peers, builder, finals, signer, faultTolerance)
	if err != nil {
		return nil, err
	}

	return &Service{chain: chain, bus: bus, runtimes: runtimes}, nil
}

// Stop terminates all running PBFT runtimes.
func (s *Service) Stop() {
	if s == nil {
		return
	}
	stopRuntimes(s.runtimes)
}

// Runtimes returns a copy of the managed runtimes keyed by node identifier.
func (s *Service) Runtimes() map[string]*NodeRuntime {
	if s == nil {
		return nil
	}

	out := make(map[string]*NodeRuntime, len(s.runtimes))
	for id, rt := range s.runtimes {
		out[id] = rt
	}
	return out
}

// Propose submits transactions into consensus using the specified validator.
// When txs is empty, the service falls back to pending chain transactions.
func (s *Service) Propose(nodeID string, txs []types.Transaction) error {
	if s == nil {
		return errors.New("consensus: service not initialized")
	}

	runtime, ok := s.runtimes[nodeID]
	if !ok {
		return fmt.Errorf("consensus: node %s not registered", nodeID)
	}

	if len(txs) == 0 {
		txs = s.chain.PendingTransactions()
	}

	if len(txs) == 0 {
		return errors.New("consensus: no transactions to propose")
	}

	return runtime.Node.ProposeBlock(txs)
}
