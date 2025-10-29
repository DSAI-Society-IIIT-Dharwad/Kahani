package consensus

import (
	"context"
	"errors"

	"storytelling-blockchain/internal/network"
)

// NodeRuntime represents a running PBFT node attached to the gossip network.
type NodeRuntime struct {
	Node    *PBFTNode
	Handler network.GossipHandler
	cancel  context.CancelFunc
}

// Stop halts the background message pump for the runtime.
func (r *NodeRuntime) Stop() {
	if r == nil || r.cancel == nil {
		return
	}
	r.cancel()
}

// StartNode bootstraps a PBFT node and launches the gossip message pump.
func StartNode(ctx context.Context, opts BootstrapOptions) (*NodeRuntime, error) {
	parent := ctx
	if parent == nil {
		parent = context.Background()
	}

	if opts.Transport == nil {
		return nil, errors.New("consensus: transport node required")
	}

	opts.Transport.DiscoverPeers(opts.Peers)

	sourceID := opts.NodeID
	if sourceID == "" {
		sourceID = opts.Transport.ID()
	}

	for _, peer := range opts.Peers {
		if peer == sourceID {
			continue
		}
		opts.Transport.ConnectToPeer(peer)
	}

	node, handler, err := BootstrapNode(opts)
	if err != nil {
		return nil, err
	}

	pumpCtx, cancel := context.WithCancel(parent)
	go pumpMessages(pumpCtx, opts.Transport, handler)

	return &NodeRuntime{Node: node, Handler: handler, cancel: cancel}, nil
}

// StartCluster bootstraps multiple PBFT nodes and launches their message pumps.
func StartCluster(ctx context.Context, transports map[string]*network.Node, peers []string, builder BlockBuilder, finalizers map[string]Finalizer, signer Signer, faultTolerance int) (map[string]*NodeRuntime, error) {
	if len(transports) == 0 {
		return nil, errors.New("consensus: transports required")
	}
	if builder == nil {
		return nil, errors.New("consensus: block builder required")
	}
	if len(finalizers) != len(transports) {
		return nil, errors.New("consensus: finalizer per node required")
	}

	runtimes := make(map[string]*NodeRuntime, len(transports))
	parent := ctx
	if parent == nil {
		parent = context.Background()
	}

	for id, transport := range transports {
		transport.DiscoverPeers(peers)

		for _, peer := range peers {
			if peer == id {
				continue
			}
			transport.ConnectToPeer(peer)
		}

		finalizer, ok := finalizers[id]
		if !ok {
			stopRuntimes(runtimes)
			return nil, errors.New("consensus: missing finalizer for node " + id)
		}

		runtime, err := StartNode(parent, BootstrapOptions{
			NodeID:         id,
			Peers:          peers,
			FaultTolerance: faultTolerance,
			Transport:      transport,
			Signer:         signer,
			Builder:        builder,
			Finalize:       finalizer,
		})
		if err != nil {
			stopRuntimes(runtimes)
			return nil, err
		}

		runtimes[id] = runtime
	}

	return runtimes, nil
}

func pumpMessages(ctx context.Context, node *network.Node, handler network.GossipHandler) {
	if node == nil || handler == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-node.ReceiveMessages():
			if !ok {
				return
			}

			_ = network.HandleIncomingMessage(handler, msg)
		}
	}
}

func stopRuntimes(runtimes map[string]*NodeRuntime) {
	for _, runtime := range runtimes {
		if runtime != nil {
			runtime.Stop()
		}
	}
}
