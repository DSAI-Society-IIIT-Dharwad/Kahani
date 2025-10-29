package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/consensus"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/supabase"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/internal/wallet"
	"storytelling-blockchain/pkg/utils"
)

type tokenVerifierStub struct {
	accepted string
}

func (t tokenVerifierStub) VerifyToken(_ context.Context, token string) (string, error) {
	if token != t.accepted {
		return "", errors.New("invalid token")
	}
	return "user-123", nil
}

type proposerStub struct {
	nodeIDs []string
	txs     [][]types.Transaction
	err     error
}

func (p *proposerStub) Propose(nodeID string, txs []types.Transaction) error {
	p.nodeIDs = append(p.nodeIDs, nodeID)
	copied := make([]types.Transaction, len(txs))
	copy(copied, txs)
	p.txs = append(p.txs, copied)
	return p.err
}

type consensusSignerStub struct{}

func (consensusSignerStub) Sign(_ []byte) (string, error)      { return "sig", nil }
func (consensusSignerStub) Verify(string, []byte, string) bool { return true }

func setupAPI(t *testing.T, opts ...func(*Config)) (*API, *blockchain.Blockchain, *wallet.Manager, *observer.Bus) {
	t.Helper()

	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 777 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	chain := blockchain.NewBlockchain()
	bus := observer.NewBus()
	chain.SetObserver(bus)
	generator, err := wallet.NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	walletObj, err := generator.GenerateWalletForUser("user-123")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	chain.RegisterWallet(walletObj)

	manager, err := wallet.NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	auth, err := supabase.NewAuthMiddleware(tokenVerifierStub{accepted: "valid-token"})
	if err != nil {
		t.Fatalf("failed to create auth middleware: %v", err)
	}

	middleware := NewMiddleware(MiddlewareConfig{
		Auth:           auth,
		AllowedOrigins: []string{"http://example.com"},
	})

	cfg := Config{Chain: chain, WalletManager: manager, Middleware: middleware, Observer: bus, IPFS: storage.NewMemoryIPFS()}
	for _, opt := range opts {
		opt(&cfg)
	}

	api, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create api: %v", err)
	}

	return api, chain, manager, bus
}

func TestHealthEndpoint(t *testing.T) {
	api, _, _, _ := setupAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Fatalf("expected CORS header to be applied")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse health response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}

	metrics, ok := body["metrics"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected metrics map in health response")
	}

	if _, ok := metrics["contributions"]; !ok {
		t.Fatalf("expected contributions metric present")
	}
}

func TestLivenessEndpoint(t *testing.T) {
	api, _, _, _ := setupAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health/live", nil)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode liveness body: %v", err)
	}

	if body["status"] != "alive" {
		t.Fatalf("expected alive status, got %s", body["status"])
	}
}

func TestReadinessEndpointDegraded(t *testing.T) {
	api, _, _, _ := setupAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health/ready", nil)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when consensus not attached, got %d", w.Code)
	}
}

func TestReadinessEndpointReady(t *testing.T) {
	stub := &proposerStub{}
	api, _, _, _ := setupAPI(t)
	api.WithConsensus("node-primary", stub, "node-backup")

	req := httptest.NewRequest(http.MethodGet, "/api/health/ready", nil)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 when consensus attached, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode readiness body: %v", err)
	}

	if body["status"] != "ready" {
		t.Fatalf("expected ready status, got %v", body["status"])
	}
}

func TestGetWalletEndpoints(t *testing.T) {
	api, _, _, _ := setupAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/api/wallet/user-123", nil)
	w := httptest.NewRecorder()
	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for existing wallet, got %d", w.Code)
	}

	reqMissing := httptest.NewRequest(http.MethodGet, "/api/wallet/unknown", nil)
	wMissing := httptest.NewRecorder()
	api.Router().ServeHTTP(wMissing, reqMissing)
	if wMissing.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing wallet, got %d", wMissing.Code)
	}
}

func TestContributeStory(t *testing.T) {
	api, chain, _, _ := setupAPI(t)

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Once upon a blockchain",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 when contribution succeeds, got %d", w.Code)
	}

	var resp struct {
		Transaction types.Transaction `json:"transaction"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Transaction.Type != "contribution" {
		t.Fatalf("unexpected transaction type: %s", resp.Transaction.Type)
	}

	if resp.Transaction.Signature == "" {
		t.Fatalf("expected signature to be populated")
	}

	if len(chain.PendingTransactions()) != 1 {
		t.Fatalf("expected pending transaction to be enqueued")
	}

	// Unauthorized request should fail.
	unauthReq := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	unauthReq.Header.Set("Content-Type", "application/json")
	unauthRes := httptest.NewRecorder()
	api.Router().ServeHTTP(unauthRes, unauthReq)
	if unauthRes.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when auth header missing, got %d", unauthRes.Code)
	}
}

func TestContributeStoryTriggersConsensus(t *testing.T) {
	stub := &proposerStub{}
	api, _, _, _ := setupAPI(t)
	api.WithConsensus("node-1", stub)

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Narrative",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 when consensus succeeds, got %d", w.Code)
	}

	var resp struct {
		Transaction types.Transaction `json:"transaction"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(stub.nodeIDs) != 1 || stub.nodeIDs[0] != "node-1" {
		t.Fatalf("expected proposer to be called once for node-1, got %v", stub.nodeIDs)
	}

	if target := api.selectConsensusNode("user-123"); target != stub.nodeIDs[0] {
		t.Fatalf("expected consensus selection %s to match proposer call %s", target, stub.nodeIDs[0])
	}

	if len(stub.txs) != 1 || len(stub.txs[0]) != 1 {
		t.Fatalf("expected proposer to receive a single transaction payload")
	}

	if stub.txs[0][0].TxID != resp.Transaction.TxID {
		t.Fatalf("expected proposer to receive transaction %s, got %s", resp.Transaction.TxID, stub.txs[0][0].TxID)
	}
}

func TestSelectConsensusNodeShards(t *testing.T) {
	stub := &proposerStub{}
	api, _, _, _ := setupAPI(t)
	api.WithConsensus("node-1", stub, "node-2", "node-3")

	keys := []string{"user-1", "user-2", "user-3", "user-4"}
	used := map[string]struct{}{}

	for _, key := range keys {
		node := api.selectConsensusNode(key)
		if node == "" {
			t.Fatalf("expected node for key %s", key)
		}
		used[node] = struct{}{}
	}

	if len(used) < 2 {
		t.Fatalf("expected sharding across nodes, saw %v", used)
	}
}

func TestContributeStoryConsensusError(t *testing.T) {
	stub := &proposerStub{err: errors.New("consensus failed")}
	api, _, _, _ := setupAPI(t)
	api.WithConsensus("node-1", stub)

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Narrative",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when consensus fails, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp["error"] != "failed to propose transaction" {
		t.Fatalf("unexpected error message: %s", resp["error"])
	}
}

func TestContributeStoryEndToEndConsensus(t *testing.T) {
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

	generator, err := wallet.NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	walletObj, err := generator.GenerateWalletForUser("user-123")
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	chain.RegisterWallet(walletObj)

	manager, err := wallet.NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	auth, err := supabase.NewAuthMiddleware(tokenVerifierStub{accepted: "valid-token"})
	if err != nil {
		t.Fatalf("failed to create auth middleware: %v", err)
	}

	middleware := NewMiddleware(MiddlewareConfig{Auth: auth})

	ipfs := storage.NewMemoryIPFS()

	api, err := New(Config{
		Chain:         chain,
		WalletManager: manager,
		Middleware:    middleware,
		Observer:      bus,
		IPFS:          ipfs,
	})
	if err != nil {
		t.Fatalf("failed to create api: %v", err)
	}

	api.WithConsensus("node-1", service)

	subID, events := bus.Subscribe(16)
	defer bus.Unsubscribe(subID)

	body, _ := json.Marshal(map[string]string{
		"story_id":   "story-1",
		"story_line": "Once upon a consensus",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/story/contribute", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp struct {
		Transaction types.Transaction `json:"transaction"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Transaction.TxID == "" {
		t.Fatalf("expected transaction id to be set")
	}

	deadline := time.After(2 * time.Second)
	seen := map[observer.EventType]bool{}

	for !(seen[observer.EventBlockCommitted] && seen[observer.EventTransactionCommitted]) {
		select {
		case ev, ok := <-events:
			if !ok {
				t.Fatal("event channel closed unexpectedly")
			}
			seen[ev.Type] = true
		case <-deadline:
			t.Fatalf("timed out waiting for consensus events, saw %#v", seen)
		}
	}

	if len(chain.Blocks()) != 2 {
		t.Fatalf("expected committed block, got %d blocks", len(chain.Blocks()))
	}

	if len(chain.PendingTransactions()) != 0 {
		t.Fatalf("expected pending transactions to be cleared")
	}
}

func TestStoryAndNFTEndpointsDefaultResponses(t *testing.T) {
	api, _, _, _ := setupAPI(t)

	tests := []struct {
		Path        string
		Method      string
		requireAuth bool
		Body        []byte
		Want        int
	}{
		{"/api/story/story-1", http.MethodGet, false, nil, http.StatusNotFound},
		{"/api/story/story-1/mint", http.MethodPost, true, []byte(`{"title":"Test","summary":"Story"}`), http.StatusNotFound},
		{"/api/nft/token-1", http.MethodGet, false, nil, http.StatusNotFound},
		{"/api/nft/token-1/authors", http.MethodGet, false, nil, http.StatusNotFound},
	}

	for _, tc := range tests {
		var bodyReader io.Reader
		if tc.Body != nil {
			bodyReader = bytes.NewReader(tc.Body)
		}

		req := httptest.NewRequest(tc.Method, tc.Path, bodyReader)
		if tc.requireAuth {
			req.Header.Set("Authorization", "Bearer valid-token")
		}
		if tc.Body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		api.Router().ServeHTTP(w, req)
		if w.Code != tc.Want {
			t.Fatalf("expected %d for %s %s, got %d", tc.Want, tc.Method, tc.Path, w.Code)
		}
	}
}

func TestGetStoryIncludesSummary(t *testing.T) {
	api, chain, manager, _ := setupAPI(t)

	walletObj, ok := manager.GetWalletBySupabaseID("user-123")
	if !ok {
		t.Fatalf("expected wallet to exist")
	}

	contribution := types.Contribution{
		ContributorID: "user-123",
		WalletAddress: walletObj.Address,
		StoryID:       "story-42",
		StoryLine:     "Opening line",
		Timestamp:     555,
	}

	signature, err := manager.SignContribution(walletObj, contribution)
	if err != nil {
		t.Fatalf("failed to sign contribution: %v", err)
	}

	contribPayload := struct {
		Contribution types.Contribution `json:"contribution"`
		Timestamp    int64              `json:"timestamp"`
	}{Contribution: contribution, Timestamp: contribution.Timestamp}

	contribBytes, err := json.Marshal(contribPayload)
	if err != nil {
		t.Fatalf("failed to marshal contribution payload: %v", err)
	}

	contribTx := types.Transaction{
		Type:      "contribution",
		Data:      contribPayload,
		Timestamp: contribution.Timestamp,
		Signature: signature,
	}
	contribTx.TxID = utils.ComputeSHA256(contribBytes)

	prev := chain.LatestBlock()
	block := blockchain.NewBlock(prev.Index+1, prev.Hash, []types.Transaction{contribTx})
	if err := chain.AddBlock(block); err != nil {
		t.Fatalf("failed to add contribution block: %v", err)
	}

	story := types.Story{
		ID:            "story-42",
		Title:         "Adventure",
		Summary:       "A quest across chains.",
		Contributions: []types.Contribution{contribution},
	}

	nft, err := blockchain.MintNFT(story, storage.NewMemoryIPFS())
	if err != nil {
		t.Fatalf("failed to mint nft: %v", err)
	}

	mintPayload, err := json.Marshal(nft)
	if err != nil {
		t.Fatalf("failed to marshal nft payload: %v", err)
	}

	mintTx := types.Transaction{
		Type:      "mint_nft",
		Data:      nft,
		Timestamp: nft.MintedAt,
	}
	mintTx.TxID = utils.ComputeSHA256(mintPayload)

	prev = chain.LatestBlock()
	mintBlock := blockchain.NewBlock(prev.Index+1, prev.Hash, []types.Transaction{mintTx})
	if err := chain.AddBlock(mintBlock); err != nil {
		t.Fatalf("failed to add mint block: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/story/story-42", nil)
	resp := httptest.NewRecorder()
	api.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var body struct {
		StoryID       string               `json:"story_id"`
		Title         string               `json:"title"`
		Summary       string               `json:"summary"`
		Contributions []types.Contribution `json:"contributions"`
		NFTs          []types.NFT          `json:"nfts"`
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Title != "Adventure" {
		t.Fatalf("expected title to match minted story, got %q", body.Title)
	}

	if body.Summary != "A quest across chains." {
		t.Fatalf("expected summary to match minted story, got %q", body.Summary)
	}

	if len(body.NFTs) != 1 {
		t.Fatalf("expected single nft entry, got %d", len(body.NFTs))
	}

	if body.NFTs[0].Summary != "A quest across chains." {
		t.Fatalf("expected nft summary to propagate, got %q", body.NFTs[0].Summary)
	}
}

func TestEventsWebsocket(t *testing.T) {
	api, _, _, bus := setupAPI(t)

	srv := httptest.NewServer(api.Router())
	defer srv.Close()

	wsURL, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("failed to parse server url: %v", err)
	}
	wsURL.Scheme = "ws"
	wsURL.Path = "/api/events"

	header := http.Header{}
	header.Set("Origin", "http://example.com")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), header)
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}
	defer conn.Close()

	// read connection ready message
	var ready map[string]string
	if err := conn.ReadJSON(&ready); err != nil {
		t.Fatalf("failed to read ready message: %v", err)
	}

	event := observer.Event{Type: observer.EventTransactionQueued, Timestamp: time.Now().UTC(), Data: map[string]string{"tx_id": "tx-1"}}
	bus.Publish(event)

	var received observer.Event
	if err := conn.ReadJSON(&received); err != nil {
		t.Fatalf("failed to read event: %v", err)
	}

	if received.Type != observer.EventTransactionQueued {
		t.Fatalf("expected transaction queued event, got %s", received.Type)
	}
}
