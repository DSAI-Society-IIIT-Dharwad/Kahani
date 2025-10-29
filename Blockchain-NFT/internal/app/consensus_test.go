package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"storytelling-blockchain/internal/api"
	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	storagepkg "storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/supabase"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/internal/wallet"
)

type tokenVerifierStub struct {
	accepted string
}

func (t tokenVerifierStub) VerifyToken(_ context.Context, token string) (string, error) {
	if token != t.accepted {
		return "", supabase.ErrMissingAuthorization
	}
	return "user-123", nil
}

type proposerRecorder struct {
	records []struct {
		nodeID string
		txLen  int
	}
}

func (p *proposerRecorder) Propose(nodeID string, txs []types.Transaction) error {
	record := struct {
		nodeID string
		txLen  int
	}{nodeID: nodeID, txLen: len(txs)}
	p.records = append(p.records, record)
	return nil
}

type consensusSignerStub struct{}

func (consensusSignerStub) Sign(_ []byte) (string, error)      { return "sig", nil }
func (consensusSignerStub) Verify(string, []byte, string) bool { return true }

func TestAttachConsensusWiresAPIAndStorage(t *testing.T) {
	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)
	t.Cleanup(bus.Close)

	storage, err := wallet.NewStorage(chain)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	generator, err := wallet.NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("generator init failed: %v", err)
	}

	userWallet, err := generator.GenerateWalletForUser("user-123")
	if err != nil {
		t.Fatalf("wallet generation failed: %v", err)
	}
	chain.RegisterWallet(userWallet)

	manager, err := wallet.NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("manager init failed: %v", err)
	}

	auth, err := supabase.NewAuthMiddleware(tokenVerifierStub{accepted: "valid-token"})
	if err != nil {
		t.Fatalf("auth middleware init failed: %v", err)
	}

	middleware := api.NewMiddleware(api.MiddlewareConfig{Auth: auth})

	apiServer, err := api.New(api.Config{
		Chain:         chain,
		WalletManager: manager,
		Middleware:    middleware,
		Observer:      bus,
		IPFS:          storagepkg.NewMemoryIPFS(),
	})
	if err != nil {
		t.Fatalf("api init failed: %v", err)
	}

	proposer := &proposerRecorder{}
	AttachConsensus(apiServer, storage, "node-1", proposer)

	newWallet, err := generator.GenerateWalletForUser("user-999")
	if err != nil {
		t.Fatalf("wallet generation failed: %v", err)
	}

	if _, err := storage.StoreWalletOnChain(newWallet); err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}

	if len(proposer.records) != 1 || proposer.records[0].nodeID != "node-1" {
		t.Fatalf("expected proposer to be invoked for storage path")
	}

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Attached consensus",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	apiServer.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected contribution to succeed, got %d", w.Code)
	}

	if len(proposer.records) != 2 {
		t.Fatalf("expected proposer to be invoked twice, got %d", len(proposer.records))
	}

	if proposer.records[1].nodeID != "node-1" || proposer.records[1].txLen != 1 {
		t.Fatalf("expected api path to propose single transaction, got %#v", proposer.records[1])
	}
}

func TestAttachConsensusHandlesNilInputs(t *testing.T) {
	AttachConsensus(nil, nil, "", nil)
	// No panic is success; nothing else to assert here.
}

func TestBootstrapConsensus(t *testing.T) {
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)
	t.Cleanup(bus.Close)

	storage, err := wallet.NewStorage(chain)
	if err != nil {
		t.Fatalf("storage init failed: %v", err)
	}

	generator, err := wallet.NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("generator init failed: %v", err)
	}

	baseWallet, err := generator.GenerateWalletForUser("user-123")
	if err != nil {
		t.Fatalf("wallet generation failed: %v", err)
	}
	chain.RegisterWallet(baseWallet)

	manager, err := wallet.NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("manager init failed: %v", err)
	}

	auth, err := supabase.NewAuthMiddleware(tokenVerifierStub{accepted: "valid-token"})
	if err != nil {
		t.Fatalf("auth middleware init failed: %v", err)
	}

	middleware := api.NewMiddleware(api.MiddlewareConfig{Auth: auth})

	apiServer, err := api.New(api.Config{
		Chain:         chain,
		WalletManager: manager,
		Middleware:    middleware,
		Observer:      bus,
		IPFS:          storagepkg.NewMemoryIPFS(),
	})
	if err != nil {
		t.Fatalf("api init failed: %v", err)
	}

	service, err := BootstrapConsensus(ConsensusBootstrapConfig{
		Chain:          chain,
		Observer:       bus,
		Storage:        storage,
		API:            apiServer,
		NodeID:         "node-1",
		Transports:     map[string]*network.Node{"node-1": node},
		Peers:          []string{"node-1"},
		Signer:         consensusSignerStub{},
		FaultTolerance: 0,
	})
	if err != nil {
		t.Fatalf("bootstrap consensus failed: %v", err)
	}
	t.Cleanup(service.Stop)

	subID, events := bus.Subscribe(16)
	defer bus.Unsubscribe(subID)

	newWallet, err := generator.GenerateWalletForUser("user-999")
	if err != nil {
		t.Fatalf("wallet generation failed: %v", err)
	}

	if _, err := storage.StoreWalletOnChain(newWallet); err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}

	waitForEvents(t, events, []observer.EventType{observer.EventBlockCommitted, observer.EventTransactionCommitted})

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Consistent consensus",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	apiServer.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected contribution to succeed, got %d", w.Code)
	}

	waitForEvents(t, events, []observer.EventType{observer.EventBlockCommitted, observer.EventTransactionCommitted})

	if len(chain.PendingTransactions()) != 0 {
		t.Fatalf("expected pending transactions to be cleared")
	}
}

func waitForEvents(t *testing.T, events <-chan observer.Event, expected []observer.EventType) {
	t.Helper()
	deadline := time.After(2 * time.Second)
	seen := make(map[observer.EventType]int, len(expected))

	needed := make(map[observer.EventType]int, len(expected))
	for _, et := range expected {
		needed[et]++
	}

	for len(seen) < len(needed) {
		select {
		case ev, ok := <-events:
			if !ok {
				t.Fatal("event channel closed unexpectedly")
			}

			if ev.Type == observer.EventError {
				if msg, ok := ev.Data.(map[string]string); ok {
					t.Fatalf("consensus error event: %s", msg["message"])
				}
				t.Fatalf("consensus error event encountered: %#v", ev.Data)
			}
			if _, interested := needed[ev.Type]; interested {
				seen[ev.Type]++
				if seen[ev.Type] >= needed[ev.Type] {
					// nothing to do, map entry already updated
				}
			}
		case <-deadline:
			t.Fatalf("timed out waiting for events, saw %#v", seen)
		}
	}
}

func TestBootstrapConsensusValidation(t *testing.T) {
	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	storage, _ := wallet.NewStorage(chain)

	testCases := []struct {
		name string
		cfg  ConsensusBootstrapConfig
	}{
		{
			name: "missing chain",
			cfg:  ConsensusBootstrapConfig{NodeID: "node-1", Storage: storage, Transports: map[string]*network.Node{"node-1": network.NewNode("node-1", network.NewInMemoryTransport())}},
		},
		{
			name: "missing targets",
			cfg:  ConsensusBootstrapConfig{Chain: chain, NodeID: "node-1", Transports: map[string]*network.Node{"node-1": network.NewNode("node-1", network.NewInMemoryTransport())}},
		},
		{
			name: "missing node id",
			cfg:  ConsensusBootstrapConfig{Chain: chain, Storage: storage, Transports: map[string]*network.Node{"node-1": network.NewNode("node-1", network.NewInMemoryTransport())}},
		},
		{
			name: "missing transports",
			cfg:  ConsensusBootstrapConfig{Chain: chain, Storage: storage, NodeID: "node-1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := BootstrapConsensus(tc.cfg); err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
		})
	}

	// ensure success path still works when requirements satisfied.
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	service, err := BootstrapConsensus(ConsensusBootstrapConfig{
		Chain:      chain,
		Observer:   bus,
		Storage:    storage,
		NodeID:     "node-1",
		Transports: map[string]*network.Node{"node-1": node},
		Signer:     consensusSignerStub{},
	})
	if err != nil {
		t.Fatalf("expected success when config valid: %v", err)
	}
	service.Stop()
}

func TestBootstrapConsensusDetectsClusterFailures(t *testing.T) {
	transport := network.NewInMemoryTransport()
	node := network.NewNode("node-1", transport)
	transport.Register(node)

	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	storage, _ := wallet.NewStorage(chain)

	_, err := BootstrapConsensus(ConsensusBootstrapConfig{
		Chain:          chain,
		Observer:       bus,
		Storage:        storage,
		NodeID:         "node-1",
		Transports:     map[string]*network.Node{"node-1": node},
		FaultTolerance: -1,
	})
	if err == nil {
		t.Fatalf("expected fault tolerance validation to fail")
	}
}

func TestDeriveConsensusNodesFromConfig(t *testing.T) {
	cfg := ConsensusBootstrapConfig{
		NodeID:         "node-primary",
		ConsensusNodes: []string{"node-primary", "node-b", "node-c", ""},
	}

	nodes := deriveConsensusNodes(cfg)
	if len(nodes) != 2 {
		t.Fatalf("expected two derived nodes, got %v", nodes)
	}

	for _, id := range nodes {
		if id == "node-primary" || id == "" {
			t.Fatalf("unexpected node in derived list: %s", id)
		}
	}
}
