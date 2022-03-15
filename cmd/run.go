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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/gasprice"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/state/tree/pb"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mtStore interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte) error
}

func start(ctx *cli.Context) error {
	c, err := config.Load(ctx)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	runMigrations(c.Database)

	etherman, err := newEtherman(*c)
	if err != nil {
		return err
	}

	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		return err
	}
	store, scCodeStore, err := newMTStores(c, sqlDB)
	if err != nil {
		return err
	}

	mt := tree.NewMerkleTree(store, c.NetworkConfig.Arity, poseidon.Hash)
	tr := tree.NewStateTree(mt, scCodeStore)

	stateCfg := state.Config{
		DefaultChainID:       c.NetworkConfig.L2DefaultChainID,
		MaxCumulativeGasUsed: c.NetworkConfig.MaxCumulativeGasUsed,
	}

	stateDb := pgstatestorage.NewPostgresStorage(sqlDB)

	var (
		st              *state.State
		grpcClientConns []*grpc.ClientConn
		cancelFuncs     []context.CancelFunc
	)
	if ctx.Bool(flagRemoteMT) {
		log.Debugf("running with remote MT")
		srvCfg := &tree.ServerConfig{
			Host: c.MTServer.Host,
			Port: c.MTServer.Port,
		}
		s := grpc.NewServer()
		mtSrv := tree.NewServer(srvCfg, tr)
		go mtSrv.Start()
		pb.RegisterMTServiceServer(s, mtSrv)

		mtClient, mtConn, mtCancel := newMTClient(c.MTClient)
		treeAdapter := tree.NewAdapter(mtClient)

		grpcClientConns = append(grpcClientConns, mtConn)
		cancelFuncs = append(cancelFuncs, mtCancel)

		st = state.NewState(stateCfg, stateDb, treeAdapter)
	} else {
		log.Debugf("running with local MT")
		st = state.NewState(stateCfg, stateDb, tr)
	}

	poolDb, err := pgpoolstorage.NewPostgresPoolStorage(c.Database)
	if err != nil {
		log.Fatal(err)
	}

	pool := pool.NewPool(poolDb, st)
	c.Sequencer.DefaultChainID = c.NetworkConfig.L2DefaultChainID
	seq := createSequencer(c.Sequencer, etherman, pool, st)

	gpe := createGasPriceEstimator(c.GasPriceEstimator, st, pool)
	go runSynchronizer(c.NetworkConfig, etherman, st, c.Synchronizer, gpe)
	go seq.Start()
	go runJSONRpcServer(*c, pool, st, seq.ChainID, gpe)

	proverClient, proverConn := newProverClient(c.Prover)
	go runAggregator(c.Aggregator, etherman, proverClient, st)

	grpcClientConns = append(grpcClientConns, proverConn)

	waitSignal(grpcClientConns, cancelFuncs)

	return nil
}

func setupLog(c log.Config) {
	log.Init(c)
}

func runMigrations(c db.Config) {
	err := db.RunMigrations(c)
	if err != nil {
		log.Fatal(err)
	}
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	auth, err := newAuthFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword, c.NetworkConfig.L1ChainID)
	if err != nil {
		return nil, err
	}
	etherman, err := etherman.NewClient(c.Etherman, auth, c.NetworkConfig.PoEAddr, c.NetworkConfig.MaticAddr)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func newProverClient(c proverclient.Config) (proverclient.ZKProverClient, *grpc.ClientConn) {
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	proverConn, err := grpc.Dial(c.ProverURI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	proverClient := proverclient.NewZKProverClient(proverConn)
	return proverClient, proverConn
}

func newMTClient(c tree.ClientConfig) (pb.MTServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	mtConn, err := grpc.DialContext(ctx, c.URI, opts...)

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	mtClient := pb.NewMTServiceClient(mtConn)

	return mtClient, mtConn, cancel
}

func runSynchronizer(networkConfig config.NetworkConfig, etherman *etherman.Client, st *state.State, cfg synchronizer.Config, gpe gasPriceEstimator) {
	genesisBlock, err := etherman.EtherClient.BlockByNumber(context.Background(), big.NewInt(0).SetUint64(networkConfig.GenBlockNumber))
	if err != nil {
		log.Fatal(err)
	}
	genesis := state.Genesis{
		Block:     genesisBlock,
		Balances:  networkConfig.Balances,
		L2ChainID: networkConfig.L2DefaultChainID,
	}
	sy, err := synchronizer.NewSynchronizer(etherman, st, networkConfig.GenBlockNumber, genesis, cfg, gpe)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(c config.Config, pool *pool.Pool, st *state.State, chainID uint64, gpe gasPriceEstimator) {
	var err error
	key, err := newKeyFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword)
	if err != nil {
		log.Fatal(err)
	}

	seqAddress := key.Address

	if err := jsonrpc.NewServer(c.RPC, c.NetworkConfig.L2DefaultChainID, seqAddress, pool, st, chainID, gpe).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(c sequencer.Config, etherman *etherman.Client, pool *pool.Pool, state *state.State) sequencer.Sequencer {
	seq, err := sequencer.NewSequencer(c, pool, state, etherman)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(c aggregator.Config, etherman *etherman.Client, proverclient proverclient.ZKProverClient, state *state.State) {
	agg, err := aggregator.NewAggregator(c, state, etherman, proverclient)
	if err != nil {
		log.Fatal(err)
	}
	agg.Start()
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
