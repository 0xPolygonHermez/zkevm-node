package main

import (
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/mocks"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/urfave/cli/v2"
)

func cmdRun(*cli.Context) error {
	log.Debugf("run command executed")

	c := config.Load()
	setupLog(c.Log)
	runMigrations(c.Database)
	etherman, err := newSimulatedEtherman(c.Synchronizer.Etherman)
	if err != nil {
		return err
	}
	state := mocks.NewState()
	pool, err := pool.NewPostgresPool(c.Database)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	//proverClient, conn := newProverClient(c.Prover)
	go runSynchronizer(c.Synchronizer, etherman, state)
	go runJSONRpcServer(c.RPC, pool, state)
	go runSequencer(c.Sequencer, etherman, pool, state)
	//go runAggregator(c.Aggregator, etherman, proverClient, state)
	//waitSignal(conn)

	waitSignal()
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

func newSimulatedEtherman(c etherman.Config) (*etherman.ClientEtherMan, error) {
	auth, err := newAuthFromKeystore(c.PrivateKeyPath, c.PrivateKeyPassword)
	if err != nil {
		return nil, err
	}
	etherman, commit, err := etherman.NewSimulatedEtherman(c, auth)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			time.Sleep(time.Second)
			commit()
		}
	}()
	return etherman, nil
}

//func newProverClient(c proverclient.Config) (proverclient.ZKProverClient, *grpc.ClientConn) {
//	opts := []grpc.DialOption{
//		// TODO: once we have user and password for prover server, change this
//		grpc.WithInsecure(),
//	}
//	conn, err := grpc.Dial(c.ProverURI, opts...)
//	if err != nil {
//		log.Fatalf("fail to dial: %v", err)
//	}
//
//	proverClient := proverclient.NewZKProverClient(conn)
//	return proverClient, conn
//}

func runSynchronizer(c synchronizer.Config, etherman *etherman.ClientEtherMan, state state.State) {
	sy, err := synchronizer.NewSynchronizer(etherman, state, c)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(jc jsonrpc.Config, pool pool.Pool, state state.State) {
	if err := jsonrpc.NewServer(jc, pool, state).Start(); err != nil {
		log.Fatal(err)
	}
}

func runSequencer(c sequencer.Config, etherman *etherman.ClientEtherMan, pool pool.Pool, state state.State) {
	seq, err := sequencer.NewSequencer(c, pool, state, etherman)
	if err != nil {
		log.Fatal(err)
	}
	seq.Start()
}

//func runAggregator(c aggregator.Config, etherman *etherman.ClientEtherMan, proverclient proverclient.ZKProverClient, state state.State) {
//	agg, err := aggregator.NewAggregator(c, state, etherman, proverclient)
//	if err != nil {
//		log.Fatal(err)
//	}
//	agg.Start()
//}

//func waitSignal(conn *grpc.ClientConn) {
func waitSignal() {
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

func newAuthFromKeystore(path, password string) (*bind.TransactOpts, error) {
	if path == "" && password == "" {
		log.Info("lol")
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
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(1337)) //nolint:gomnd
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = 99999999999
	return auth, nil
}
