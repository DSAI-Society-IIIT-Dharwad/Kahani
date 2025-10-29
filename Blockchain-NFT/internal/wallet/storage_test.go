package wallet

import (
	"context"
	"errors"
	"testing"
	"time"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/consensus"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/types"
)

func TestStoreWalletOnChain(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 1111 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	chain := blockchain.NewBlockchain()
	storage, err := NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	generator, err := NewGenerator("secret")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	wallet, err := generator.GenerateWalletForUser("user-abc")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	tx, err := storage.StoreWalletOnChain(wallet)
	if err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}

	if tx.Type != "create_wallet" {
		t.Fatalf("expected create_wallet transaction, got %s", tx.Type)
	}

	if tx.Timestamp != 1111 {
		t.Fatalf("expected timestamp from override, got %d", tx.Timestamp)
	}

	pending := chain.PendingTransactions()
	if len(pending) != 1 {
		t.Fatalf("expected pending transaction to be enqueued")
	}

	state := chain.State()
	if _, ok := state.WalletRegistry["user-abc"]; !ok {
		t.Fatalf("wallet not registered in state")
	}
}

type proposerStub struct {
	calls  int
	nodeID string
	err    error
	nodes  []string
}

func (p *proposerStub) Propose(nodeID string, _ []types.Transaction) error {
	p.calls++
	p.nodeID = nodeID
	p.nodes = append(p.nodes, nodeID)
	return p.err
}

type consensusSignerStub struct{}

func (consensusSignerStub) Sign(_ []byte) (string, error)      { return "sig", nil }
func (consensusSignerStub) Verify(string, []byte, string) bool { return true }

func TestStoreWalletOnChainTriggersConsensus(t *testing.T) {
	chain := blockchain.NewBlockchain()
	storage, err := NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	stub := &proposerStub{}
	storage.WithConsensus("node-1", stub)

	generator, err := NewGenerator("secret")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	wallet, err := generator.GenerateWalletForUser("user-1")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	t.Logf("wallet address: %s", wallet.Address)
	if _, err := storage.StoreWalletOnChain(wallet); err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}

	if stub.calls != 1 || stub.nodeID != "node-1" {
		t.Fatalf("expected proposer to be called once for node-1, got %d calls", stub.calls)
	}
}

func TestStoreWalletOnChainPropagatesConsensusError(t *testing.T) {
	chain := blockchain.NewBlockchain()
	storage, err := NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	stub := &proposerStub{err: errors.New("propose failed")}
	storage.WithConsensus("node-1", stub)

	wallet := types.Wallet{SupabaseUserID: "user-1"}

	if _, err := storage.StoreWalletOnChain(wallet); err == nil {
		t.Fatal("expected error when proposer fails")
	}
}

func TestStoreWalletOnChainConsensusIntegration(t *testing.T) {
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)
	t.Cleanup(bus.Close)

	service, err := consensus.StartService(context.Background(), chain, bus, map[string]*network.Node{"node-1": node}, []string{"node-1"}, consensusSignerStub{}, 0)
	if err != nil {
		t.Fatalf("failed to start consensus service: %v", err)
	}
	t.Cleanup(service.Stop)

	storage, err := NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	storage.WithConsensus("node-1", service)

	subID, events := bus.Subscribe(16)
	defer bus.Unsubscribe(subID)

	generator, err := NewGenerator("secret")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	wallet, err := generator.GenerateWalletForUser("user-1")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	if _, err := storage.StoreWalletOnChain(wallet); err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}

	deadline := time.After(2 * time.Second)
	seen := map[observer.EventType]bool{}

	for !(seen[observer.EventBlockCommitted] && seen[observer.EventTransactionCommitted]) {
		select {
		case ev, ok := <-events:
			if !ok {
				t.Fatal("event channel closed unexpectedly")
			}
			if ev.Type == observer.EventError {
				if payload, ok := ev.Data.(map[string]string); ok {
					t.Fatalf("consensus error event: %s", payload["message"])
				}
				t.Fatalf("consensus error event encountered: %#v", ev.Data)
			}
			seen[ev.Type] = true
		case <-deadline:
			t.Fatalf("timed out waiting for consensus events, saw %#v", seen)
		}
	}

	if len(chain.Blocks()) != 2 {
		t.Fatalf("expected consensus to commit block, got %d blocks", len(chain.Blocks()))
	}

	if len(chain.PendingTransactions()) != 0 {
		t.Fatalf("expected pending transactions to be cleared")
	}
}

func TestStoreWalletOnChainShardedConsensus(t *testing.T) {
	chain := blockchain.NewBlockchain()
	storage, err := NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	stub := &proposerStub{}
	storage.WithConsensus("node-1", stub, "node-2", "node-3")

	users := []string{"user-alpha", "user-beta", "user-gamma", "user-delta"}

	generator, err := NewGenerator("secret")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	for _, user := range users {
		wallet, err := generator.GenerateWalletForUser(user)
		if err != nil {
			t.Fatalf("failed to generate wallet for %s: %v", user, err)
		}
		if _, err := storage.StoreWalletOnChain(wallet); err != nil {
			t.Fatalf("store wallet failed for %s: %v", user, err)
		}
	}

	if len(stub.nodes) != len(users) {
		t.Fatalf("expected proposer to be invoked %d times, got %d", len(users), len(stub.nodes))
	}

	used := map[string]struct{}{}
	for _, node := range stub.nodes {
		used[node] = struct{}{}
	}

	if len(used) < 2 {
		t.Fatalf("expected sharding to target multiple nodes, saw %v", used)
	}
}
