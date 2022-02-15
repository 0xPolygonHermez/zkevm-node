package main

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"

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
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func start(ctx *cli.Context) error {
	configFilePath := ctx.String(flagCfg)
	network := ctx.String(flagNetwork)
	c, err := config.Load(configFilePath, network)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	runMigrations(c.Database)

	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		log.Fatal(err)
		return err
	}
	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, c.NetworkConfig.Arity, poseidon.Hash)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore, []byte{})

	stateCfg := state.Config{
		DefaultChainID: c.NetworkConfig.L2DefaultChainID,
	}

	stateDb := pgstatestorage.NewPostgresStorage(sqlDB)
	st := state.NewState(stateCfg, stateDb, tr)

	pool, err := pool.NewPostgresPool(c.Database)
	if err != nil {
		log.Fatal(err)
	}

	c.Sequencer.DefaultChainID = c.NetworkConfig.L2DefaultChainID
	seq := createSequencer(c.Sequencer, etherman, pool, st)

	gpe := gasprice.NewEstimator(c.GasPriceEstimator, st, pool)
	go runSynchronizer(c.NetworkConfig, etherman, st, c.Synchronizer, gpe)
	go seq.Start()
	go runJSONRpcServer(*c, pool, st, seq.ChainID, gpe)

	proverClient, conn := newProverClient(c.Prover)
	go runAggregator(c.Aggregator, etherman, proverClient, st)
	waitSignal(conn)
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

func newEtherman(c config.Config) (*etherman.ClientEtherMan, error) {
	auth, err := newAuthFromKeystore(c.Etherman.PrivateKeyPath, c.Etherman.PrivateKeyPassword, c.NetworkConfig.L1ChainID)
	if err != nil {
		return nil, err
	}
	etherman, err := etherman.NewEtherman(c.Etherman, auth, c.NetworkConfig.PoEAddr, c.NetworkConfig.BridgeAddr, c.NetworkConfig.MaticAddr)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func newProverClient(c proverclient.Config) (proverclient.ZKProverClient, *grpc.ClientConn) {
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(c.ProverURI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	proverClient := proverclient.NewZKProverClient(conn)
	return proverClient, conn
}

func runSynchronizer(networkConfig config.NetworkConfig, etherman *etherman.ClientEtherMan, st state.State, cfg synchronizer.Config, gpe gasprice.Estimator) {
	genesisBlock, err := etherman.EtherClient.BlockByNumber(context.Background(), big.NewInt(0).SetUint64(networkConfig.GenBlockNumber))
	if err != nil {
		log.Fatal(err)
	}
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: networkConfig.Balances,
	}
	sy, err := synchronizer.NewSynchronizer(etherman, st, networkConfig.GenBlockNumber, genesis, cfg, gpe)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(c config.Config, pool *pool.PostgresPool, st state.State, chainID uint64, gpe gasprice.Estimator) {
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

func createSequencer(c sequencer.Config, etherman *etherman.ClientEtherMan, pool *pool.PostgresPool, state state.State) sequencer.Sequencer {
	seq, err := sequencer.NewSequencer(c, pool, state, etherman)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func runAggregator(c aggregator.Config, etherman *etherman.ClientEtherMan, proverclient proverclient.ZKProverClient, state state.State) {
	agg, err := aggregator.NewAggregator(c, state, etherman, proverclient)
	if err != nil {
		log.Fatal(err)
	}
	agg.Start()
}

func waitSignal(conn *grpc.ClientConn) {
	//func waitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")
			//conn.Close() //nolint:gosec,errcheck
			os.Exit(0)
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
