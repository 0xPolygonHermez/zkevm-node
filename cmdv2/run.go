package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	aggregator "github.com/hermeznetwork/hermez-core/aggregatorv2"
	config "github.com/hermeznetwork/hermez-core/configv2"
	"github.com/hermeznetwork/hermez-core/db"
	etherman "github.com/hermeznetwork/hermez-core/ethermanv2"
	gasprice "github.com/hermeznetwork/hermez-core/gaspricev2"
	jsonrpc "github.com/hermeznetwork/hermez-core/jsonrpcv2"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/merkletree"
	pool "github.com/hermeznetwork/hermez-core/poolv2"
	"github.com/hermeznetwork/hermez-core/poolv2/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/proverclient"
	proverclientpb "github.com/hermeznetwork/hermez-core/proverclient/pb"
	sequencer "github.com/hermeznetwork/hermez-core/sequencerv2"
	"github.com/hermeznetwork/hermez-core/sequencerv2/broadcast"
	"github.com/hermeznetwork/hermez-core/sequencerv2/broadcast/pb"
	"github.com/hermeznetwork/hermez-core/state/tree"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor"
	synchronizer "github.com/hermeznetwork/hermez-core/synchronizerv2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mtStore interface {
	SupportsDBTransactions() bool
	BeginDBTransaction(ctx context.Context, txBundleID string) error
	Commit(ctx context.Context, txBundleID string) error
	Rollback(ctx context.Context, txBundleID string) error
	Get(ctx context.Context, key []byte, txBundleID string) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte, txBundleID string) error
}

// slice contains method
func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func start(ctx *cli.Context) error {
	c, err := config.Load(ctx)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	runMigrations(c.Database)

	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		return err
	}

	stateDb := state.NewPostgresStorage(sqlDB)
	executorClient, _ := executor.NewExecutorClient(c.Executor)
	stateDBClient, _ := merkletree.NewStateDBServiceClient(merkletree.Config(c.MTClient))
	stateTree := merkletree.NewStateTree(stateDBClient)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: c.NetworkConfig.MaxCumulativeGasUsed,
	}
	st := state.NewState(stateCfg, stateDb, executorClient, stateTree)

	var (
		grpcClientConns []*grpc.ClientConn
		cancelFuncs     []context.CancelFunc
	)

	poolDb, err := pgpoolstorage.NewPostgresPoolStorage(c.Database)
	if err != nil {
		log.Fatal(err)
	}

	proverClient, proverConn := newProverClient(c.Prover)

	var npool *pool.Pool
	var gpe gasPriceEstimator
	var etherman *etherman.Client

	if contains(ctx.StringSlice(config.FlagComponents), AGGREGATOR) ||
		contains(ctx.StringSlice(config.FlagComponents), SEQUENCER) ||
		contains(ctx.StringSlice(config.FlagComponents), SYNCHRONIZER) {
		var err error
		etherman, err = newEtherman(*c)
		if err != nil {
			log.Fatal(err)
		}
	}

	npool = pool.New(poolDb, st, c.NetworkConfig.L2GlobalExitRootManagerAddr)
	gpe = createGasPriceEstimator(c.GasPriceEstimator, st, npool)

	for _, item := range ctx.StringSlice(config.FlagComponents) {
		switch item {
		case AGGREGATOR:
			log.Info("Running aggregator")
			go runAggregator(c.Aggregator, etherman, proverClient, st)
		case SEQUENCER:
			log.Info("Running sequencer")
			seq := createSequencer(c.Sequencer, etherman, npool, st)
			go seq.Start(ctx.Context)
		case RPC:
			log.Info("Running JSON-RPC server")
			apis := map[string]bool{}
			for _, a := range ctx.StringSlice(config.FlagHTTPAPI) {
				apis[a] = true
			}
			go runJSONRpcServer(*c, npool, st, c.RPC.ChainID, gpe, apis)
		case SYNCHRONIZER:
			log.Info("Running synchronizer")
			go runSynchronizer(c.NetworkConfig, etherman, st, c.Synchronizer)
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

func runMigrations(c db.Config) {
	err := db.RunMigrationsUp(c)
	if err != nil {
		log.Fatal(err)
	}
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	auth, err := newAuthFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword, c.NetworkConfig.ChainID)
	if err != nil {
		return nil, err
	}
	etherman, err := etherman.NewClient(c.Etherman, auth, c.NetworkConfig.PoEAddr, c.NetworkConfig.MaticAddr, c.NetworkConfig.L2GlobalExitRootManagerAddr)
	if err != nil {
		return nil, err
	}
	return etherman, nil
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

func runSynchronizer(networkConfig config.NetworkConfig, etherman *etherman.Client, st *state.State, cfg synchronizer.Config) {
	sy, err := synchronizer.NewSynchronizer(etherman, st, networkConfig.GenBlockNumber, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(c config.Config, pool *pool.Pool, st *state.State, chainID uint64, gpe gasPriceEstimator, apis map[string]bool) {
	storage, err := jsonrpc.NewPostgresStorage(c.Database)
	if err != nil {
		log.Fatal(err)
	}

	if err := jsonrpc.NewServer(c.RPC, chainID, pool, st, gpe, storage, apis).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(c sequencer.Config, etherman *etherman.Client, pool *pool.Pool, state *state.State) *sequencer.Sequencer {
	// TODO: define priceGetter and txManager
	seq, err := sequencer.New(c, pool, state, etherman, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(c aggregator.Config, etherman *etherman.Client, proverclient proverclientpb.ZKProverServiceClient, state *state.State) {
	agg, err := aggregator.NewAggregator(c, state, etherman, proverclient)
	if err != nil {
		log.Fatal(err)
	}
	agg.Start()
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
		return gasprice.NewEstimatorLastNBatches(cfg, state)
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

func newMTStores(c *config.Config, sqlDB *pgxpool.Pool) (mtStore, mtStore, error) {
	switch c.MTServer.StoreBackend {
	case tree.PgMTStoreBackend:
		store := tree.NewPostgresStore(sqlDB)
		scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)

		return store, scCodeStore, nil
	case tree.PgRistrettoMTStoreBackend:
		cache, err := tree.NewStoreCache()
		if err != nil {
			return nil, nil, err
		}
		store := tree.NewPgRistrettoStore(sqlDB, cache)
		scCodeStore := tree.NewPgRistrettoSCCodeStore(sqlDB, cache)
		return store, scCodeStore, nil
	case tree.BadgerRistrettoMTStoreBackend:
		cache, err := tree.NewStoreCache()
		if err != nil {
			return nil, nil, err
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, err
		}
		dataDir := path.Join(home, ".hermezcore", "db")
		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}
		db, err := tree.NewBadgerDB(dataDir)
		if err != nil {
			return nil, nil, err
		}
		store := tree.NewBadgerRistrettoStore(db, cache)
		return store, store, nil
	}
	return nil, nil, fmt.Errorf("Unknown MT store backend: %q", c.MTServer.StoreBackend)
}
