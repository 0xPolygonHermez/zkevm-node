package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	if c.Metrics.Enabled {
		metrics.Init()
	}
	runStateMigrations(c.StateDB)
	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	var (
		grpcClientConns []*grpc.ClientConn
		cancelFuncs     []context.CancelFunc
		etherman        *etherman.Client
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
	c.Aggregator.ChainID = l2ChainID
	c.RPC.ChainID = l2ChainID
	log.Infof("Chain ID read from POE SC = %v", l2ChainID)

	ctx := context.Background()
	st := newState(ctx, c, l2ChainID, stateSqlDB)

	ethTxManager := ethtxmanager.New(c.EthTxManager, etherman)

	for _, item := range cliCtx.StringSlice(config.FlagComponents) {
		switch item {
		case AGGREGATOR:
			log.Info("Running aggregator")
			go runAggregator(ctx, c.Aggregator, etherman, ethTxManager, st, grpcClientConns)
		case SEQUENCER:
			log.Info("Running sequencer")
			poolInstance := createPool(c.PoolDB, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			gpe := createGasPriceEstimator(c.GasPriceEstimator, st, poolInstance)
			seq := createSequencer(*c, poolInstance, st, etherman, ethTxManager, gpe)
			go seq.Start(ctx)
		case RPC:
			log.Info("Running JSON-RPC server")
			runRPCMigrations(c.RPC.DB)
			poolInstance := createPool(c.PoolDB, c.NetworkConfig.L2BridgeAddr, l2ChainID, st)
			gpe := createGasPriceEstimator(c.GasPriceEstimator, st, poolInstance)
			apis := map[string]bool{}
			for _, a := range cliCtx.StringSlice(config.FlagHTTPAPI) {
				apis[a] = true
			}
			go runJSONRPCServer(*c, poolInstance, st, gpe, apis)
		case SYNCHRONIZER:
			log.Info("Running synchronizer")
			go runSynchronizer(*c, etherman, st)
		case BROADCAST:
			log.Info("Running broadcast service")
			go runBroadcastServer(c.BroadcastServer, st)
		}
	}

	if c.Metrics.Enabled {
		go startMetricsHttpServer(c)
	}

	waitSignal(grpcClientConns, cancelFuncs)

	return nil
}

func setupLog(c log.Config) {
	log.Init(c)
}

func runStateMigrations(c db.Config) {
	runMigrations(c, db.StateMigrationName)
}

func runPoolMigrations(c db.Config) {
	runMigrations(c, db.PoolMigrationName)
}

func runRPCMigrations(c db.Config) {
	runMigrations(c, db.RPCMigrationName)
}

func runMigrations(c db.Config, name string) {
	err := db.RunMigrationsUp(c, name)
	if err != nil {
		log.Fatal(err)
	}
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	auth, err := newAuthFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword, c.Etherman.L1ChainID)
	if err != nil {
		return nil, err
	}
	etherman, err := etherman.NewClient(c.Etherman, auth)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func runSynchronizer(cfg config.Config, etherman *etherman.Client, st *state.State) {
	sy, err := synchronizer.NewSynchronizer(cfg.IsTrustedSequencer, etherman, st, cfg.NetworkConfig.Genesis, cfg.Synchronizer)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRPCServer(c config.Config, pool *pool.Pool, st *state.State, gpe gasPriceEstimator, apis map[string]bool) {
	storage, err := jsonrpc.NewPostgresStorage(c.RPC.DB)
	if err != nil {
		log.Fatal(err)
	}

	c.RPC.MaxCumulativeGasUsed = c.Sequencer.MaxCumulativeGasUsed

	if err := jsonrpc.NewServer(c.RPC, pool, st, gpe, storage, apis).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(c config.Config, pool *pool.Pool, state *state.State, etherman *etherman.Client,
	ethTxManager *ethtxmanager.Client, gpe gasPriceEstimator) *sequencer.Sequencer {
	pg, err := pricegetter.NewClient(c.PriceGetter)
	if err != nil {
		log.Fatal(err)
	}

	seq, err := sequencer.New(c.Sequencer, pool, state, etherman, pg, ethTxManager, gpe)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(ctx context.Context, c aggregator.Config, ethman *etherman.Client, ethTxManager *ethtxmanager.Client, state *state.State, grpcClientConns []*grpc.ClientConn) {
	agg, err := aggregator.New(c, state, ethTxManager, ethman, grpcClientConns)
	if err != nil {
		log.Fatal(err)
	}
	agg.Start(ctx)
}

func runBroadcastServer(c broadcast.ServerConfig, st *state.State) {
	s := grpc.NewServer()

	broadcastSrv := broadcast.NewServer(&c, st)
	pb.RegisterBroadcastServiceServer(s, broadcastSrv)

	broadcastSrv.Start()
}

// gasPriceEstimator interface for gas price estimator.
type gasPriceEstimator interface {
	GetAvgGasPrice(ctx context.Context) (*big.Int, error)
	UpdateGasPriceAvg(newValue *big.Int)
}

// createGasPriceEstimator init gas price gasPriceEstimator based on type in config.
func createGasPriceEstimator(cfg gasprice.Config, state *state.State, pool *pool.Pool) gasPriceEstimator {
	switch cfg.Type {
	case gasprice.AllBatchesType:
		return gasprice.NewEstimatorAllBatches()
	case gasprice.LastNBatchesType:
		return gasprice.NewEstimatorLastNL2Blocks(cfg, state)
	case gasprice.DefaultType:
		return gasprice.NewDefaultEstimator(cfg, pool)
	}
	return nil
}

func waitSignal(conns []*grpc.ClientConn, cancelFuncs []context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")

			exitStatus := 0
			for _, conn := range conns {
				if err := conn.Close(); err != nil {
					log.Errorf("Could not properly close gRPC connection: %v", err)
					exitStatus = -1
				}
			}
			for _, cancel := range cancelFuncs {
				cancel()
			}
			os.Exit(exitStatus)
		}
	}
}

func newKeyFromKeystore(path, password string) (*keystore.Key, error) {
	if path == "" && password == "" {
		return nil, nil
	}
	keystoreEncrypted, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keystoreEncrypted, password)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func newAuthFromKeystore(path, password string, chainID uint64) (*bind.TransactOpts, error) {
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, nil
	}
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func newState(ctx context.Context, c *config.Config, l2ChainID uint64, sqlDB *pgxpool.Pool) *state.State {
	stateDb := state.NewPostgresStorage(sqlDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, c.Executor)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, c.MTClient)
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := state.Config{
		MaxCumulativeGasUsed: c.Sequencer.MaxCumulativeGasUsed,
		ChainID:              l2ChainID,
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree)
	return st
}

func createPool(poolDBConfig db.Config, l2BridgeAddr common.Address, l2ChainID uint64, st *state.State) *pool.Pool {
	runPoolMigrations(poolDBConfig)
	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(poolDBConfig)
	if err != nil {
		log.Fatal(err)
	}
	poolInstance := pool.NewPool(poolStorage, st, l2BridgeAddr, l2ChainID)
	return poolInstance
}

func startMetricsHttpServer(c *config.Config) {
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", c.Metrics.Host, c.Metrics.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener for metrics: %v", err)
		return
	}
	mux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Handler: mux,
	}
	log.Infof("metrics server listening on port %d", c.Metrics.Port)
	if err := metricsServer.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Warnf("http server for metrics stopped")
			return
		}
		log.Errorf("closed http connection for metrics server: %v", err)
		return
	}
}
