package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"storytelling-blockchain/internal/api"
	"storytelling-blockchain/internal/app"
	"storytelling-blockchain/internal/blockchain"
	configpkg "storytelling-blockchain/internal/config"
	"storytelling-blockchain/internal/network"
	"storytelling-blockchain/internal/observer"
	"storytelling-blockchain/internal/storage"
	"storytelling-blockchain/internal/supabase"
	"storytelling-blockchain/internal/wallet"
)

type devTokenVerifier struct{}

func (devTokenVerifier) VerifyToken(_ context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("dev: token required")
	}
	return token, nil
}

type config struct {
	NodeID         string
	ClusterNodes   []string
	HTTPAddr       string
	AllowedOrigins []string
	Passphrase     string
	FaultTolerance int
	SeedUsers      []string
	DataDir        string
	Supabase       supabaseSettings
	IPFSEndpoint   string
}

type supabaseSettings struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	PollInterval   time.Duration
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("booting dev node",
		"node", cfg.NodeID,
		"http", cfg.HTTPAddr,
		"cluster", cfg.ClusterNodes,
		"faultTolerance", cfg.FaultTolerance,
		"storage", cfg.DataDir,
		"ipfs", cfg.IPFSEndpoint,
	)

	stateStore, err := storage.NewBadgerStorage(storage.BadgerConfig{Path: cfg.DataDir})
	if err != nil {
		fail(logger, "state store init failed", err)
	}
	defer func() {
		if err := stateStore.Close(); err != nil {
			logger.Warn("state store close failed", "error", err)
		}
	}()

	chain, err := blockchain.LoadBlockchain(stateStore)
	if err != nil {
		fail(logger, "blockchain load failed", err)
	}

	var ipfsClient storage.IPFSClient
	if cfg.IPFSEndpoint != "" {
		client, ipfsErr := storage.NewIPFSShell(cfg.IPFSEndpoint)
		if ipfsErr != nil {
			logger.Warn("ipfs shell init failed", "endpoint", cfg.IPFSEndpoint, "error", ipfsErr)
		} else {
			ipfsClient = client
		}
	}

	if ipfsClient == nil {
		ipfsClient = storage.NewMemoryIPFS()
		logger.Warn("ipfs endpoint not configured, using in-memory store")
	}

	bus := observer.NewBus()
	chain.SetObserver(bus)

	generator, err := wallet.NewGenerator(cfg.Passphrase)
	if err != nil {
		fail(logger, "generator init failed", err)
	}

	manager, err := wallet.NewManager(chain, cfg.Passphrase)
	if err != nil {
		fail(logger, "manager init failed", err)
	}

	walletStorage, err := wallet.NewStorage(chain)
	if err != nil {
		fail(logger, "storage init failed", err)
	}

	var (
		tokenVerifier  supabase.TokenVerifier = devTokenVerifier{}
		supabaseClient *supabase.Client
	)

	if cfg.Supabase.URL != "" && cfg.Supabase.AnonKey != "" {
		client, clientErr := supabase.NewClient(supabase.Config{
			URL:            cfg.Supabase.URL,
			AnonKey:        cfg.Supabase.AnonKey,
			ServiceRoleKey: cfg.Supabase.ServiceRoleKey,
		})
		if clientErr != nil {
			fail(logger, "supabase client init failed", clientErr)
		}
		supabaseClient = client
		tokenVerifier = client
		logger.Info("supabase auth enabled", "url", cfg.Supabase.URL)
	} else {
		logger.Warn("supabase credentials missing, falling back to dev token verifier")
	}

	authMiddleware, err := supabase.NewAuthMiddleware(tokenVerifier)
	if err != nil {
		fail(logger, "auth middleware init failed", err)
	}

	middleware := api.NewMiddleware(api.MiddlewareConfig{
		Auth:           authMiddleware,
		AllowedOrigins: cfg.AllowedOrigins,
	})

	apiServer, err := api.New(api.Config{
		Chain:          chain,
		WalletManager:  manager,
		Middleware:     middleware,
		Observer:       bus,
		ConsensusNode:  cfg.NodeID,
		ConsensusNodes: cfg.ClusterNodes,
		IPFS:           ipfsClient,
	})
	if err != nil {
		fail(logger, "api init failed", err)
	}

	transports := buildClusterTransport(cfg.ClusterNodes)
	peers := append([]string{}, cfg.ClusterNodes...)

	service, err := app.BootstrapConsensus(app.ConsensusBootstrapConfig{
		Context:        ctx,
		Chain:          chain,
		Observer:       bus,
		Storage:        walletStorage,
		API:            apiServer,
		NodeID:         cfg.NodeID,
		Transports:     transports,
		Peers:          peers,
		FaultTolerance: cfg.FaultTolerance,
		ConsensusNodes: cfg.ClusterNodes,
	})
	if err != nil {
		fail(logger, "bootstrap consensus failed", err)
	}
	defer service.Stop()

	if supabaseClient != nil {
		if cfg.Supabase.ServiceRoleKey != "" {
			pollInterval := cfg.Supabase.PollInterval
			if pollInterval <= 0 {
				pollInterval = 30 * time.Second
			}

			p, pollerErr := supabase.NewPoller(supabaseClient, generator, walletStorage, manager, pollInterval)
			if pollerErr != nil {
				fail(logger, "supabase poller init failed", pollerErr)
			}
			go runSupabasePoller(ctx, logger, p)
			logger.Info("supabase poller started", "interval", pollInterval)
		} else {
			logger.Info("supabase poller disabled", "reason", "service role key missing")
		}
	}

	seedWallets(ctx, logger, chain, generator, walletStorage, supabaseClient, cfg.SeedUsers)

	go logObserver(ctx, logger, bus)

	server := &http.Server{Addr: cfg.HTTPAddr, Handler: apiServer.Router()}

	go func() {
		logger.Info("http server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fail(logger, "http server error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Warn("http shutdown error", "error", err)
	}

	logger.Info("shutdown complete")
}

func loadConfig() (config, error) {
	const (
		defaultNodeID       = "node-1"
		defaultHTTPAddr     = ":8080"
		defaultPassphrase   = "local-passphrase"
		defaultOrigin       = "http://localhost:3000"
		defaultSeedUser     = "user-123"
		defaultDataDir      = "devnode-data"
		defaultPollInterval = 30 * time.Second
	)

	configFlag := flag.String("config", strings.TrimSpace(os.Getenv("DEVNODE_CONFIG")), "path to YAML config file")
	nodeFlag := flag.String("node", "", "node identifier")
	httpFlag := flag.String("http", "", "HTTP listen address")
	peersFlag := flag.String("peers", "", "comma separated peer identifiers")
	originsFlag := flag.String("origins", "", "comma separated allowed origins")
	passFlag := flag.String("passphrase", "", "wallet passphrase")
	faultFlag := flag.Int("fault", -1, "fault tolerance threshold")
	clusterFlag := flag.Int("cluster-size", -1, "number of nodes to auto-provision when peers omitted")
	seedFlag := flag.String("seed-users", "", "comma separated supabase user IDs to pre-provision")
	dataDirFlag := flag.String("data-dir", "", "path to persistent storage directory")
	ipfsFlag := flag.String("ipfs-api", "", "IPFS API endpoint")
	pollFlag := flag.Duration("supabase-poll-interval", 0, "Supabase poll interval (e.g. 30s)")
	flag.Parse()

	setFlags := map[string]bool{}
	flag.CommandLine.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	configPath := strings.TrimSpace(*configFlag)
	if configPath == "" {
		if _, err := os.Stat("config/config.yaml"); err == nil {
			configPath = "config/config.yaml"
		}
	}

	var fileCfg *configpkg.FileConfig
	if configPath != "" {
		cfg, err := configpkg.LoadFile(configPath)
		if err != nil {
			return config{}, err
		}
		fileCfg = cfg
	}

	envNode := strings.TrimSpace(os.Getenv("DEVNODE_NODE_ID"))
	envHTTP := strings.TrimSpace(os.Getenv("DEVNODE_HTTP_ADDR"))
	envPeers := strings.TrimSpace(os.Getenv("DEVNODE_PEERS"))
	envOrigins := strings.TrimSpace(os.Getenv("DEVNODE_ALLOWED_ORIGINS"))
	envPass := strings.TrimSpace(os.Getenv("DEVNODE_WALLET_PASSPHRASE"))
	envFault := strings.TrimSpace(os.Getenv("DEVNODE_FAULT_TOLERANCE"))
	envCluster := strings.TrimSpace(os.Getenv("DEVNODE_CLUSTER_SIZE"))
	envSeeds := strings.TrimSpace(os.Getenv("DEVNODE_SEED_USERS"))
	envDataDir := strings.TrimSpace(os.Getenv("DEVNODE_DATA_DIR"))
	envIPFS := strings.TrimSpace(os.Getenv("DEVNODE_IPFS_API"))
	envPoll := strings.TrimSpace(os.Getenv("SUPABASE_POLL_INTERVAL"))
	envSupabaseURL := strings.TrimSpace(os.Getenv("SUPABASE_URL"))
	envSupabaseAnon := strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
	envSupabaseService := strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_KEY"))
	if envSupabaseService == "" {
		envSupabaseService = strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	}
	if envIPFS == "" {
		envIPFS = strings.TrimSpace(os.Getenv("IPFS_API"))
	}

	var (
		fileNodeID       string
		fileHTTP         string
		filePeers        []string
		fileOrigins      []string
		fileDataDir      string
		fileIPFS         string
		filePoll         time.Duration
		fileSupabaseURL  string
		fileSupabaseAnon string
		fileSupabaseServ string
	)

	if fileCfg != nil {
		fileNodeID = fileCfg.Node.ID
		if fileCfg.Node.Port > 0 {
			fileHTTP = fmt.Sprintf(":%d", fileCfg.Node.Port)
		}
		filePeers = append([]string{}, fileCfg.Network.BootstrapPeers...)
		fileOrigins = append([]string{}, fileCfg.API.AllowedOrigins...)
		fileDataDir = strings.TrimSpace(fileCfg.Storage.BadgerPath)
		fileIPFS = strings.TrimSpace(fileCfg.Storage.IPFSAPI)
		filePoll = fileCfg.Supabase.PollInterval.Duration
		fileSupabaseURL = fileCfg.Supabase.URL
		fileSupabaseAnon = fileCfg.Supabase.AnonKey
		fileSupabaseServ = fileCfg.Supabase.ServiceRoleKey
	}

	nodeID := pickString(setFlags["node"], *nodeFlag, envNode, fileNodeID, defaultNodeID)
	httpAddr := pickString(setFlags["http"], *httpFlag, envHTTP, fileHTTP, defaultHTTPAddr)
	passphrase := pickString(setFlags["passphrase"], *passFlag, envPass, "", defaultPassphrase)
	dataDir := pickString(setFlags["data-dir"], *dataDirFlag, envDataDir, fileDataDir, defaultDataDir)
	ipfsEndpoint := pickString(setFlags["ipfs-api"], *ipfsFlag, envIPFS, fileIPFS, "")

	peersList := pickStringSlice(setFlags["peers"], *peersFlag, envPeers, filePeers, nil)
	clusterSize := pickInt(setFlags["cluster-size"], *clusterFlag, envCluster, len(filePeers), 1)
	if len(peersList) == 0 {
		for i := 1; i <= max(1, clusterSize); i++ {
			peersList = append(peersList, fmt.Sprintf("node-%d", i))
		}
	}
	clusterNodes := orderWithPrimary(nodeID, peersList)
	if len(clusterNodes) == 0 {
		clusterNodes = []string{nodeID}
	}

	allowedOrigins := pickStringSlice(setFlags["origins"], *originsFlag, envOrigins, fileOrigins, []string{defaultOrigin})
	seedUsers := pickStringSlice(setFlags["seed-users"], *seedFlag, envSeeds, nil, []string{defaultSeedUser})
	faultTolerance := pickInt(setFlags["fault"], *faultFlag, envFault, 0, 0)

	supabaseURL := pickString(false, "", envSupabaseURL, fileSupabaseURL, "")
	supabaseAnon := pickString(false, "", envSupabaseAnon, fileSupabaseAnon, "")
	supabaseService := pickString(false, "", envSupabaseService, fileSupabaseServ, "")
	pollInterval := pickDuration(setFlags["supabase-poll-interval"], *pollFlag, envPoll, filePoll, defaultPollInterval)

	return config{
		NodeID:         nodeID,
		ClusterNodes:   clusterNodes,
		HTTPAddr:       httpAddr,
		AllowedOrigins: allowedOrigins,
		Passphrase:     passphrase,
		FaultTolerance: faultTolerance,
		SeedUsers:      seedUsers,
		DataDir:        dataDir,
		Supabase: supabaseSettings{
			URL:            supabaseURL,
			AnonKey:        supabaseAnon,
			ServiceRoleKey: supabaseService,
			PollInterval:   pollInterval,
		},
		IPFSEndpoint: ipfsEndpoint,
	}, nil
}

func splitAndClean(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func orderWithPrimary(primary string, nodes []string) []string {
	ordered := []string{}
	if primary != "" {
		ordered = append(ordered, primary)
	}

	for _, node := range nodes {
		if node == primary {
			continue
		}
		ordered = append(ordered, node)
	}

	return uniqueStrings(ordered)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func pickString(flagUsed bool, flagVal, envVal, fileVal, fallback string) string {
	if flagUsed {
		if v := strings.TrimSpace(flagVal); v != "" {
			return v
		}
	}

	if v := strings.TrimSpace(envVal); v != "" {
		return v
	}

	if v := strings.TrimSpace(fileVal); v != "" {
		return v
	}

	return fallback
}

func pickStringSlice(flagUsed bool, flagVal, envVal string, fileVal []string, fallback []string) []string {
	if flagUsed {
		if values := splitAndClean(flagVal); len(values) > 0 {
			return uniqueStrings(values)
		}
	}

	if trimmed := strings.TrimSpace(envVal); trimmed != "" {
		if values := splitAndClean(trimmed); len(values) > 0 {
			return uniqueStrings(values)
		}
	}

	if len(fileVal) > 0 {
		values := make([]string, 0, len(fileVal))
		for _, v := range fileVal {
			if trimmed := strings.TrimSpace(v); trimmed != "" {
				values = append(values, trimmed)
			}
		}
		if len(values) > 0 {
			return uniqueStrings(values)
		}
	}

	if len(fallback) == 0 {
		return nil
	}

	values := make([]string, len(fallback))
	copy(values, fallback)
	return values
}

func pickInt(flagUsed bool, flagVal int, envVal string, fileVal int, fallback int) int {
	if flagUsed {
		return flagVal
	}

	if trimmed := strings.TrimSpace(envVal); trimmed != "" {
		if parsed, err := strconv.Atoi(trimmed); err == nil {
			return parsed
		}
	}

	if fileVal != 0 {
		return fileVal
	}

	return fallback
}

func pickDuration(flagUsed bool, flagVal time.Duration, envVal string, fileVal time.Duration, fallback time.Duration) time.Duration {
	if flagUsed {
		return flagVal
	}

	if trimmed := strings.TrimSpace(envVal); trimmed != "" {
		if parsed, err := time.ParseDuration(trimmed); err == nil {
			return parsed
		}
	}

	if fileVal > 0 {
		return fileVal
	}

	return fallback
}

func buildClusterTransport(nodes []string) map[string]*network.Node {
	transport := network.NewInMemoryTransport()
	registry := make(map[string]*network.Node, len(nodes))

	for _, id := range nodes {
		if id == "" {
			continue
		}
		node := network.NewNode(id, transport)
		transport.Register(node)
		registry[id] = node
	}

	return registry
}

func seedWallets(ctx context.Context, logger *slog.Logger, chain *blockchain.Blockchain, generator *wallet.Generator, storage *wallet.Storage, supabaseClient *supabase.Client, users []string) {
	for _, userID := range users {
		if chain != nil {
			if _, exists := chain.GetWalletBySupabaseID(userID); exists {
				logger.Info("wallet already seeded", "user", userID)
				continue
			}
		}
		walletObj, err := generator.GenerateWalletForUser(userID)
		if err != nil {
			logger.Warn("wallet generation failed", "user", userID, "error", err)
			continue
		}

		if supabaseClient != nil {
			if err := supabaseClient.UpsertWallet(ctx, walletObj); err != nil {
				logger.Warn("supabase wallet sync failed", "user", userID, "error", err)
				continue
			}
		}

		if _, err := storage.StoreWalletOnChain(walletObj); err != nil {
			logger.Warn("wallet provisioning failed", "user", userID, "error", err)
			continue
		}

		logger.Info("seeded wallet", "user", userID)
	}

}

func runSupabasePoller(ctx context.Context, logger *slog.Logger, poller *supabase.Poller) {
	if poller == nil {
		return
	}

	ticker := time.NewTicker(poller.Interval())
	defer ticker.Stop()

	if _, err := poller.PollNewUsers(ctx); err != nil {
		logger.Warn("supabase poll failed", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("supabase poller stopping")
			return
		case <-ticker.C:
			if _, err := poller.PollNewUsers(ctx); err != nil {
				logger.Warn("supabase poll failed", "error", err)
			}
		}
	}
}

func logObserver(ctx context.Context, logger *slog.Logger, bus *observer.Bus) {
	if bus == nil {
		return
	}

	id, ch := bus.Subscribe(64)
	if id == "" {
		return
	}
	defer bus.Unsubscribe(id)

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			logger.Info("observer event", "type", string(event.Type))
		}
	}
}

func fail(logger *slog.Logger, message string, err error) {
	logger.Error(message, "error", err)
	os.Exit(1)
}
