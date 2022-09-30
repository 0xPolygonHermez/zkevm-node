package main

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/proverclient"
	proverclientpb "github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	runStateMigrations(c.StateDB)
	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	st := newState(ctx, c, stateSqlDB)

	var (
		grpcClientConns []*grpc.ClientConn
		cancelFuncs     []context.CancelFunc
		etherman        *etherman.Client
	)

	if strings.Contains(cliCtx.String(config.FlagComponents), AGGREGATOR) ||
		strings.Contains(cliCtx.String(config.FlagComponents), SEQUENCER) ||
		strings.Contains(cliCtx.String(config.FlagComponents), SYNCHRONIZER) ||
		strings.Contains(cliCtx.String(config.FlagComponents), RPC) {
		var err error
		etherman, err = newEtherman(*c)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: READ CHAIN ID FROM POE SC
		// log.Info("VALUE = %v", etherman.GetPublicAddress().String())
		c.NetworkConfig.L2ChainID = 1001
		log.Infof("Chain ID read from POE SC = %v", c.NetworkConfig.L2ChainID)

		if strings.Contains(cliCtx.String(config.FlagComponents), RPC) {
			etherman = nil
		}
	}

	ethTxManager := ethtxmanager.New(c.EthTxManager, etherman)
	proverClient, proverConn := newProverClient(c.Prover)
	for _, item := range cliCtx.StringSlice(config.FlagComponents) {
		switch item {
		case AGGREGATOR:
			log.Info("Running aggregator")
			c.Aggregator.ChainID = c.NetworkConfig.L2ChainID
			go runAggregator(ctx, c.Aggregator, etherman, ethTxManager, proverClient, st)
		case SEQUENCER:
			log.Info("Running sequencer")
			poolInstance := createPool(c.PoolDB, c.NetworkConfig, st)
			seq := createSequencer(*c, poolInstance, st, etherman, ethTxManager)
			go seq.Start(ctx)
		case RPC:
			log.Info("Running JSON-RPC server")
			runRPCMigrations(c.RPC.DB)
			poolInstance := createPool(c.PoolDB, c.NetworkConfig, st)
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

	grpcClientConns = append(grpcClientConns, proverConn)

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
	auth, err := newAuthFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword, c.NetworkConfig.L1ChainID)
	if err != nil {
		return nil, err
	}
	etherman, err := etherman.NewClient(c.Etherman, auth, c.NetworkConfig.PoEAddr, c.NetworkConfig.MaticAddr, c.NetworkConfig.GlobalExitRootManagerAddr)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func runSynchronizer(cfg config.Config, etherman *etherman.Client, st *state.State) {
	sy, err := synchronizer.NewSynchronizer(cfg.IsTrustedSequencer, etherman, st, cfg.NetworkConfig.GenBlockNumber, cfg.NetworkConfig.Genesis, cfg.Synchronizer)
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
	c.RPC.ChainID = c.NetworkConfig.L2ChainID

	if err := jsonrpc.NewServer(c.RPC, pool, st, gpe, storage, apis).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(c config.Config, pool *pool.Pool, state *state.State, etherman *etherman.Client,
	ethTxManager *ethtxmanager.Client) *sequencer.Sequencer {
	pg, err := pricegetter.NewClient(c.PriceGetter)
	if err != nil {
		log.Fatal(err)
	}

	seq, err := sequencer.New(c.Sequencer, pool, state, etherman, pg, ethTxManager)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(ctx context.Context, c aggregator.Config, ethman *etherman.Client, ethTxManager *ethtxmanager.Client,
	proverClient proverclientpb.ZKProverServiceClient, state *state.State) {
	agg, err := aggregator.NewAggregator(c, state, ethTxManager, ethman, proverClient)
	if err != nil {
		log.Fatal(err)
	}
	agg.Start(ctx)
}

func newProverClient(c proverclient.Config) (proverclientpb.ZKProverServiceClient, *grpc.ClientConn) {
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	proverConn, err := grpc.Dial(c.ProverURI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	proverClient := proverclientpb.NewZKProverServiceClient(proverConn)
	return proverClient, proverConn
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
		log.Fatal(err)
	}
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		log.Fatal(err)
	}
	return auth, nil
}

func newState(ctx context.Context, c *config.Config, sqlDB *pgxpool.Pool) *state.State {
	stateDb := state.NewPostgresStorage(sqlDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, c.Executor)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, c.MTClient)
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := state.Config{
		MaxCumulativeGasUsed: c.Sequencer.MaxCumulativeGasUsed,
		ChainID:              c.NetworkConfig.L2ChainID,
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree)
	return st
}

func createPool(poolDBConfig db.Config, networkConfig config.NetworkConfig, st *state.State) *pool.Pool {
	runPoolMigrations(poolDBConfig)
	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(poolDBConfig)
	if err != nil {
		log.Fatal(err)
	}
	poolInstance := pool.NewPool(poolStorage, st, networkConfig.L2GlobalExitRootManagerAddr)
	return poolInstance
}
