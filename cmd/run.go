package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	two = 2
	ten = 10
)

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		return err
	}
	setupLog(c.Log)

	if c.Log.Environment == log.EnvironmentDevelopment {
		zkevm.PrintVersion(os.Stdout)
		log.Info("Starting application")
	} else if c.Log.Environment == log.EnvironmentProduction {
		logVersion()
	}

	if c.Metrics.Enabled {
		metrics.Init()
	}
	components := cliCtx.StringSlice(config.FlagComponents)

	// Only runs migration if the component is the synchronizer and if the flag is deactivated
	if !cliCtx.Bool(config.FlagMigrations) {
		for _, comp := range components {
			if comp == SYNCHRONIZER {
				runStateMigrations(c.StateDB)
			}
		}
	}
	checkStateMigrations(c.StateDB)

	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	var (
		cancelFuncs []context.CancelFunc
		etherman    *etherman.Client
	)

	etherman, err = newEtherman(*c)
	if err != nil {
		log.Fatal(err)
	}

	// READ CHAIN ID FROM POE SC
	l2ChainID, err := etherman.GetL2ChainID()
	if err != nil {
		log.Fatal(err)
	}
	// Read Fork ID FROM POE SC
	forkIDIntervals, err := etherman.GetForks(cliCtx.Context)
	if err != nil || len(forkIDIntervals) == 0 {
		log.Fatal("error getting forks: ", err)
	}
	currentForkID := forkIDIntervals[len(forkIDIntervals)-1].ForkId

	c.Aggregator.ChainID = l2ChainID
	c.Aggregator.ForkId = currentForkID
	c.RPC.ChainID = l2ChainID
	log.Infof("Chain ID read from POE SC = %v", l2ChainID)

	ctx := context.Background()
	st := newState(ctx, c, l2ChainID, forkIDIntervals, stateSqlDB)

	ethTxManagerStorage, err := ethtxmanager.NewPostgresStorage(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	etm := ethtxmanager.New(c.EthTxManager, etherman, ethTxManagerStorage, st)

	for _, component := range components {
		switch component {
		case AGGREGATOR:
			log.Info("Running aggregator")
			go runAggregator(ctx, c.Aggregator, etherman, etm, st)
		case SEQUENCER:
			log.Info("Running sequencer")
			poolInstance := createPool(c.Pool, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			seq := createSequencer(*c, poolInstance, ethTxManagerStorage, st)
			go seq.Start(ctx)
		case RPC:
			log.Info("Running JSON-RPC server")
			poolInstance := createPool(c.Pool, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			if c.RPC.EnableL2SuggestedGasPricePolling {
				// Needed for rejecting transactions with too low gas price
				poolInstance.StartPollingMinSuggestedGasPrice(ctx)
			}
			apis := map[string]bool{}
			for _, a := range cliCtx.StringSlice(config.FlagHTTPAPI) {
				apis[a] = true
			}
			go runJSONRPCServer(*c, poolInstance, st, apis)
		case SYNCHRONIZER:
			log.Info("Running synchronizer")
			poolInstance := createPool(c.Pool, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			go runSynchronizer(*c, etherman, etm, st, poolInstance)
		case BROADCAST:
			log.Info("Running broadcast service")
			go runBroadcastServer(c.BroadcastServer, st)
		case ETHTXMANAGER:
			log.Info("Running eth tx manager service")
			etm := createEthTxManager(*c, ethTxManagerStorage, st)
			go etm.Start()
		case L2GASPRICER:
			log.Info("Running L2 gasPricer")
			poolInstance := createPool(c.Pool, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			go runL2GasPriceSuggester(c.L2GasPriceSuggester, st, poolInstance, etherman)
		}
	}

	if c.Metrics.Enabled {
		go startMetricsHttpServer(c.Metrics)
	}

	if c.Metrics.ProfilingEnabled {
		go startProfilingHttpServer(c.Metrics)
	}
	waitSignal(cancelFuncs)

	return nil
}

func setupLog(c log.Config) {
	log.Init(c)
}

func runStateMigrations(c db.Config) {
	runMigrations(c, db.StateMigrationName)
}

func checkStateMigrations(c db.Config) {
	err := db.CheckMigrations(c, db.StateMigrationName)
	if err != nil {
		log.Fatal(err)
	}
}

func runPoolMigrations(c db.Config) {
	runMigrations(c, db.PoolMigrationName)
}

func runMigrations(c db.Config, name string) {
	log.Infof("running migrations for %v", name)
	err := db.RunMigrationsUp(c, name)
	if err != nil {
		log.Fatal(err)
	}
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	etherman, err := etherman.NewClient(c.Etherman)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func runSynchronizer(cfg config.Config, etherman *etherman.Client, ethTxManager *ethtxmanager.Client, st *state.State, pool *pool.Pool) {
	var trustedSequencerURL string
	var err error
	if !cfg.IsTrustedSequencer {
		if cfg.Synchronizer.TrustedSequencerURL != "" {
			trustedSequencerURL = cfg.Synchronizer.TrustedSequencerURL
		} else {
			log.Debug("getting trusted sequencer URL from smc")
			trustedSequencerURL, err = etherman.GetTrustedSequencerURL()
			if err != nil {
				log.Fatal("error getting trusted sequencer URI. Error: %v", err)
			}
		}
		log.Debug("trustedSequencerURL ", trustedSequencerURL)
	}
	zkEVMClient := client.NewClient(trustedSequencerURL)

	sy, err := synchronizer.NewSynchronizer(cfg.IsTrustedSequencer, etherman, st, pool, ethTxManager, zkEVMClient, cfg.NetworkConfig.Genesis, cfg.Synchronizer)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRPCServer(c config.Config, pool *pool.Pool, st *state.State, apis map[string]bool) {
	storage := jsonrpc.NewStorage()
	c.RPC.MaxCumulativeGasUsed = c.Sequencer.MaxCumulativeGasUsed

	if err := jsonrpc.NewServer(c.RPC, pool, st, storage, apis).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(cfg config.Config, pool *pool.Pool, etmStorage *ethtxmanager.PostgresStorage, st *state.State) *sequencer.Sequencer {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, privateKey := range cfg.Sequencer.Finalizer.PrivateKeys {
		_, err := etherman.LoadAuthFromKeyStore(privateKey.Path, privateKey.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	ethTxManager := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)

	seq, err := sequencer.New(cfg.Sequencer, pool, st, etherman, ethTxManager)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(ctx context.Context, c aggregator.Config, etherman *etherman.Client, ethTxManager *ethtxmanager.Client, st *state.State) {
	agg, err := aggregator.New(c, st, ethTxManager, etherman)
	if err != nil {
		log.Fatal(err)
	}
	err = agg.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func runBroadcastServer(c broadcast.ServerConfig, st *state.State) {
	s := grpc.NewServer()

	broadcastSrv := broadcast.NewServer(&c, st)
	pb.RegisterBroadcastServiceServer(s, broadcastSrv)

	broadcastSrv.Start()
}

// runL2GasPriceSuggester init gas price gasPriceEstimator based on type in config.
func runL2GasPriceSuggester(cfg gasprice.Config, state *state.State, pool *pool.Pool, etherman *etherman.Client) {
	ctx := context.Background()
	gasprice.NewL2GasPriceSuggester(ctx, cfg, pool, etherman, state)
}

func waitSignal(cancelFuncs []context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")

			exitStatus := 0
			for _, cancel := range cancelFuncs {
				cancel()
			}
			os.Exit(exitStatus)
		}
	}
}

func newState(ctx context.Context, c *config.Config, l2ChainID uint64, forkIDIntervals []state.ForkIDInterval, sqlDB *pgxpool.Pool) *state.State {
	stateDb := state.NewPostgresStorage(sqlDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, c.Executor)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, c.MTClient)
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := state.Config{
		MaxCumulativeGasUsed: c.Sequencer.MaxCumulativeGasUsed,
		ChainID:              l2ChainID,
		ForkIDIntervals:      forkIDIntervals,
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree)
	return st
}

func createPool(cfgPool pool.Config, l2BridgeAddr common.Address, l2ChainID uint64, st *state.State) *pool.Pool {
	runPoolMigrations(cfgPool.DB)
	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(cfgPool.DB)
	if err != nil {
		log.Fatal(err)
	}
	poolInstance := pool.NewPool(cfgPool, poolStorage, st, l2BridgeAddr, l2ChainID)
	return poolInstance
}

func createEthTxManager(cfg config.Config, etmStorage *ethtxmanager.PostgresStorage, st *state.State) *ethtxmanager.Client {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, privateKey := range cfg.EthTxManager.PrivateKeys {
		_, err := etherman.LoadAuthFromKeyStore(privateKey.Path, privateKey.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	etm := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)
	return etm
}

func startProfilingHttpServer(c metrics.Config) {
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", c.ProfilingHost, c.ProfilingPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener for profiling: %v", err)
		return
	}
	mux.HandleFunc(metrics.ProfilingIndexEndpoint, pprof.Index)
	mux.HandleFunc(metrics.ProfileEndpoint, pprof.Profile)
	mux.HandleFunc(metrics.ProfilingCmdEndpoint, pprof.Cmdline)
	mux.HandleFunc(metrics.ProfilingSymbolEndpoint, pprof.Symbol)
	mux.HandleFunc(metrics.ProfilingTraceEndpoint, pprof.Trace)
	profilingServer := &http.Server{
		Handler:     mux,
		ReadTimeout: two * time.Minute,
	}
	log.Infof("profiling server listening on port %d", c.ProfilingPort)
	if err := profilingServer.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Warnf("http server for profiling stopped")
			return
		}
		log.Errorf("closed http connection for profiling server: %v", err)
		return
	}
}

func startMetricsHttpServer(c metrics.Config) {
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener for metrics: %v", err)
		return
	}
	mux.Handle(metrics.Endpoint, promhttp.Handler())

	metricsServer := &http.Server{
		Handler:     mux,
		ReadTimeout: ten * time.Second,
	}
	log.Infof("metrics server listening on port %d", c.Port)
	if err := metricsServer.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Warnf("http server for metrics stopped")
			return
		}
		log.Errorf("closed http connection for metrics server: %v", err)
		return
	}
}

func logVersion() {
	log.Infow("Starting application",
		// node version is already logged by default
		"gitRevision", zkevm.GitRev,
		"gitBranch", zkevm.GitBranch,
		"goVersion", runtime.Version(),
		"built", zkevm.BuildDate,
		"os/arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	)
}
