package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/consensus/sharding"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/supabase"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/internal/wallet"
	"storytelling-blockchain/pkg/utils"
)

// Config bundles the dependencies required to construct the API server.
type Config struct {
	Chain          *blockchain.Blockchain
	WalletManager  *wallet.Manager
	Middleware     *Middleware
	Observer       *observer.Bus
	Proposer       Proposer
	ConsensusNode  string
	ConsensusNodes []string
	IPFS           storage.IPFSClient
}

// Proposer encapsulates the ability to submit transactions into consensus.
type Proposer interface {
	Propose(nodeID string, txs []types.Transaction) error
}

// API exposes the HTTP handlers for the storytelling blockchain service.
type API struct {
	chain          *blockchain.Blockchain
	walletManager  *wallet.Manager
	middleware     *Middleware
	router         *mux.Router
	observer       *observer.Bus
	upgrader       websocket.Upgrader
	proposer       Proposer
	consensusNode  string
	consensusNodes []string
	metrics        apiMetrics
	startedAt      time.Time
	ipfs           storage.IPFSClient
}

type apiMetrics struct {
	contributions uint64
}

func (m *apiMetrics) incContributions() {
	atomic.AddUint64(&m.contributions, 1)
}

func (m *apiMetrics) snapshot() map[string]uint64 {
	return map[string]uint64{
		"contributions": atomic.LoadUint64(&m.contributions),
	}
}

// New creates an API instance configured with routes and middleware.
func New(cfg Config) (*API, error) {
	if cfg.Chain == nil {
		return nil, errors.New("api: blockchain chain cannot be nil")
	}

	if cfg.WalletManager == nil {
		return nil, errors.New("api: wallet manager cannot be nil")
	}

	if cfg.IPFS == nil {
		return nil, errors.New("api: ipfs client cannot be nil")
	}

	router := mux.NewRouter()
	upgrader := websocket.Upgrader{}

	if cfg.Middleware != nil {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			return cfg.Middleware.OriginAllowed(origin)
		}
	} else {
		upgrader.CheckOrigin = func(*http.Request) bool { return true }
	}

	api := &API{
		chain:          cfg.Chain,
		walletManager:  cfg.WalletManager,
		middleware:     cfg.Middleware,
		router:         router,
		observer:       cfg.Observer,
		upgrader:       upgrader,
		proposer:       cfg.Proposer,
		consensusNode:  cfg.ConsensusNode,
		consensusNodes: append([]string{}, cfg.ConsensusNodes...),
		startedAt:      time.Now().UTC(),
		ipfs:           cfg.IPFS,
	}

	api.registerRoutes()
	return api, nil
}

// Router exposes the configured mux router for mounting in an HTTP server.
func (a *API) Router() http.Handler {
	return a.router
}

func (a *API) registerRoutes() {
	base := a.router.PathPrefix("/api").Subrouter()
	if a.middleware != nil {
		base.Use(a.middleware.Wrap)
	}

	base.HandleFunc("/health", a.handleHealth).Methods(http.MethodGet)
	base.HandleFunc("/health/live", a.handleLiveness).Methods(http.MethodGet)
	base.HandleFunc("/health/ready", a.handleReadiness).Methods(http.MethodGet)
	base.HandleFunc("/blockchain", a.handleBlockchainState).Methods(http.MethodGet)
	base.HandleFunc("/wallet/{userID}", a.handleGetWallet).Methods(http.MethodGet)
	base.HandleFunc("/story/{storyID}", a.handleGetStory).Methods(http.MethodGet)
	base.HandleFunc("/nft/{tokenID}", a.handleGetNFT).Methods(http.MethodGet)
	base.HandleFunc("/nft/{tokenID}/authors", a.handleGetNFTAuthors).Methods(http.MethodGet)
	base.HandleFunc("/events", a.handleEvents).Methods(http.MethodGet)

	authSub := base.PathPrefix("").Subrouter()
	if a.middleware != nil && a.middleware.AuthMiddleware() != nil {
		authSub.Use(a.middleware.AuthMiddleware().Wrap)
	}

	authSub.HandleFunc("/story/contribute", a.handleContributeStory).Methods(http.MethodPost)
	authSub.HandleFunc("/story/{storyID}/mint", a.handleMintStory).Methods(http.MethodPost)
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":               "ok",
		"blocks":               len(a.chain.Blocks()),
		"pending_transactions": len(a.chain.PendingTransactions()),
		"consensus": map[string]interface{}{
			"attached":     a.proposer != nil,
			"primary_node": a.consensusNode,
			"nodes":        append([]string{}, a.consensusNodes...),
		},
		"metrics":        a.metrics.snapshot(),
		"uptime_seconds": int(time.Since(a.startedAt).Round(time.Second) / time.Second),
	}

	writeJSON(w, http.StatusOK, status)
}

func (a *API) handleLiveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

func (a *API) handleReadiness(w http.ResponseWriter, r *http.Request) {
	checks := map[string]bool{
		"blockchain":    a.chain != nil,
		"walletManager": a.walletManager != nil,
		"consensus":     a.proposer != nil && len(a.consensusNodes) > 0,
	}

	for _, ok := range checks {
		if !ok {
			writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status": "degraded",
				"checks": checks,
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ready",
		"checks": checks,
	})
}

func (a *API) handleBlockchainState(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"blocks": a.chain.Blocks(),
		"state":  a.chain.State(),
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *API) handleGetWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user id is required")
		return
	}

	wallet, ok := a.chain.GetWalletBySupabaseID(userID)
	if !ok {
		writeError(w, http.StatusNotFound, "wallet not found")
		return
	}
	writeJSON(w, http.StatusOK, wallet)
}

func (a *API) handleContributeStory(w http.ResponseWriter, r *http.Request) {
	userID, ok := supabase.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	var request struct {
		StoryID   string `json:"story_id"`
		StoryLine string `json:"story_line"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if request.StoryID == "" || request.StoryLine == "" {
		writeError(w, http.StatusBadRequest, "story_id and story_line are required")
		return
	}

	wallet, ok := a.walletManager.GetWalletBySupabaseID(userID)
	if !ok {
		writeError(w, http.StatusNotFound, "wallet not found")
		return
	}

	contribution := types.Contribution{
		ContributorID: userID,
		WalletAddress: wallet.Address,
		StoryID:       request.StoryID,
		StoryLine:     request.StoryLine,
		Timestamp:     types.NowUnix(),
	}

	signature, err := a.walletManager.SignContribution(wallet, contribution)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to sign contribution")
		return
	}

	payload := struct {
		Contribution types.Contribution `json:"contribution"`
		Timestamp    int64              `json:"timestamp"`
	}{Contribution: contribution, Timestamp: contribution.Timestamp}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to marshal contribution")
		return
	}

	tx := types.Transaction{
		Type:      "contribution",
		Data:      payload,
		Timestamp: contribution.Timestamp,
		Signature: signature,
	}
	tx.TxID = utils.ComputeSHA256(payloadBytes)

	a.chain.EnqueueTransaction(tx)

	nodeID := a.selectConsensusNode(userID)

	if a.proposer != nil && nodeID != "" {
		if err := a.proposer.Propose(nodeID, []types.Transaction{tx}); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to propose transaction")
			return
		}
	}

	a.metrics.incContributions()

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"transaction": tx,
	})
}

func (a *API) handleGetStory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storyID := vars["storyID"]
	if storyID == "" {
		writeError(w, http.StatusBadRequest, "story id is required")
		return
	}

	contributions := a.chain.StoryContributions(storyID)
	if len(contributions) == 0 {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}

	authors := blockchain.AggregateAuthors(contributions)
	nfts := a.chain.NFTsByStory(storyID)

	var (
		title   string
		summary string
	)

	if len(nfts) > 0 {
		latest := nfts[0]
		for _, nft := range nfts[1:] {
			if nft.BlockIndex > latest.BlockIndex {
				latest = nft
				continue
			}
			if nft.BlockIndex == latest.BlockIndex && nft.MintedAt > latest.MintedAt {
				latest = nft
			}
		}

		title = latest.Title
		summary = latest.Summary
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"story_id":      storyID,
		"title":         title,
		"summary":       summary,
		"contributions": contributions,
		"authors":       authors,
		"nfts":          nfts,
	})
}

func (a *API) handleGetNFT(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	if tokenID == "" {
		writeError(w, http.StatusBadRequest, "token id is required")
		return
	}

	nft, ok := a.chain.GetNFT(tokenID)
	if !ok {
		writeError(w, http.StatusNotFound, "nft not found")
		return
	}

	writeJSON(w, http.StatusOK, nft)
}

func (a *API) handleGetNFTAuthors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	if tokenID == "" {
		writeError(w, http.StatusBadRequest, "token id is required")
		return
	}

	nft, ok := a.chain.GetNFT(tokenID)
	if !ok {
		writeError(w, http.StatusNotFound, "nft not found")
		return
	}

	authors := append([]types.Author{nft.MainAuthor}, nft.CoAuthors...)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token_id": tokenID,
		"authors":  authors,
	})
}

func (a *API) handleMintStory(w http.ResponseWriter, r *http.Request) {
	userID, ok := supabase.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	vars := mux.Vars(r)
	storyID := vars["storyID"]
	if storyID == "" {
		writeError(w, http.StatusBadRequest, "story id is required")
		return
	}

	var request struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(request.Title) == "" || strings.TrimSpace(request.Summary) == "" {
		writeError(w, http.StatusBadRequest, "title and summary are required")
		return
	}

	contributions := a.chain.StoryContributions(storyID)
	if len(contributions) == 0 {
		writeError(w, http.StatusNotFound, "story has no contributions")
		return
	}

	authors := blockchain.AggregateAuthors(contributions)
	if len(authors) == 0 {
		writeError(w, http.StatusNotFound, "story has no authors")
		return
	}

	if authors[0].SupabaseUserID != userID {
		writeError(w, http.StatusForbidden, "only the main author can mint the story")
		return
	}

	if existing := a.chain.NFTsByStory(storyID); len(existing) > 0 {
		writeError(w, http.StatusConflict, "story already minted")
		return
	}

	if a.ipfs == nil {
		writeError(w, http.StatusServiceUnavailable, "ipfs unavailable")
		return
	}

	story := types.Story{
		ID:            storyID,
		Title:         strings.TrimSpace(request.Title),
		Summary:       strings.TrimSpace(request.Summary),
		Contributions: contributions,
	}

	nft, err := blockchain.MintNFT(story, a.ipfs)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, blockchain.ErrNoContributions),
			errors.Is(err, blockchain.ErrMissingStoryID),
			errors.Is(err, blockchain.ErrMissingTitle):
			status = http.StatusBadRequest
		case errors.Is(err, blockchain.ErrNilIPFSClient):
			status = http.StatusServiceUnavailable
		}

		writeError(w, status, err.Error())
		return
	}

	payload, err := json.Marshal(nft)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encode nft")
		return
	}

	tx := types.Transaction{
		Type:      "mint_nft",
		Data:      nft,
		Timestamp: nft.MintedAt,
	}
	tx.TxID = utils.ComputeSHA256(payload)

	a.chain.EnqueueTransaction(tx)

	nodeID := a.selectConsensusNode(storyID)
	if a.proposer != nil && nodeID != "" {
		if err := a.proposer.Propose(nodeID, []types.Transaction{tx}); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to propose transaction")
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"nft":         nft,
		"transaction": tx,
	})
}

func (a *API) handleNotImplemented(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "endpoint not implemented yet")
}

// WithConsensus wires the consensus proposer into the API at runtime.
func (a *API) WithConsensus(nodeID string, proposer Proposer, additional ...string) {
	if a == nil {
		return
	}

	a.proposer = proposer
	a.consensusNode = nodeID
	a.consensusNodes = append([]string{}, additional...)
	if nodeID != "" {
		a.consensusNodes = append([]string{nodeID}, a.consensusNodes...)
	}
}

func (a *API) selectConsensusNode(key string) string {
	if len(a.consensusNodes) == 0 {
		return ""
	}

	return sharding.SelectNode(a.consensusNodes, key)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
