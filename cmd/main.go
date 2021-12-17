package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

const (
	flagCfg     = "cfg"
	flagNetwork = "network"
)

var (
	// version represents the program based on the git tag
	version = "v0.1.0"
	// commit represents the program based on the git commit
	commit = "dev"
	// date represents the date of application was built
	date = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "hermez-node"
	app.Version = version
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     flagCfg,
			Aliases:  []string{"c"},
			Usage:    "Configuration `FILE`",
			Required: false,
		},
		&cli.StringFlag{
			Name:     flagNetwork,
			Aliases:  []string{"n"},
			Usage:    "Network: mainnet, testnet, internaltestnet, local. By default it uses mainnet",
			Required: false,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Application version and build",
			Action:  versionCmd,
		},
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Run the hermez core",
			Action:  start,
			Flags:   flags,
		},
		{
			Name:    "register",
			Aliases: []string{"reg"},
			Usage:   "Register sequencer in the smart contract",
			Action:  registerSequencer,
			Flags:   flags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}

func start(ctx *cli.Context) error {
	configFilePath := ctx.String(flagCfg)
	network := ctx.String(flagNetwork)
	c, err := config.Load(configFilePath, network)
	if err != nil {
		return err
	}
	setupLog(c.Log)

	runMigrations(c.Database)

	etherman, err := newSimulatedEtherman(c.Etherman)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		log.Fatal(err)
		return err
	}
	mt := tree.NewMerkleTree(sqlDB, c.NetworkConfig.Arity, poseidon.Hash)
	tr := tree.NewStateTree(mt, []byte{})
	st := state.NewState(sqlDB, tr)

	pool, err := pool.NewPostgresPool(c.Database)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//proverClient, conn := newProverClient(c.Prover)
	go runSynchronizer(c.NetworkConfig.GenBlockNumber, etherman, st)
	go runJSONRpcServer(c.RPC, c.Etherman, c.NetworkConfig, pool, st)
	go runSequencer(c.Sequencer, etherman, pool, st)
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

func runSynchronizer(genBlockNumber uint64, etherman *etherman.ClientEtherMan, state state.State) {
	sy, err := synchronizer.NewSynchronizer(etherman, state, genBlockNumber)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRpcServer(jc jsonrpc.Config, ec etherman.Config, nc config.NetworkConfig, pool pool.Pool, st state.State) {
	var err error
	var seq *state.Sequencer

	key, err := newKeyFromKeystore(ec.PrivateKeyPath, ec.PrivateKeyPassword)
	if err != nil {
		log.Fatal(err)
	}

	seqAddress := key.Address

	const intervalToCheckSequencerRegistrationInSeconds = 10

	for {
		seq, err = st.GetSequencer(context.Background(), seqAddress)
		if err != nil {
			log.Warnf("Make sure the address %s has been registered in the smart contract as a sequencer, err: %v", seqAddress.Hex(), err)
			time.Sleep(intervalToCheckSequencerRegistrationInSeconds * time.Second)
			continue
		}
		break
	}

	if err := jsonrpc.NewServer(jc, nc.L2DefaultChainID, seq.ChainID.Uint64(), pool, st).Start(); err != nil {
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

func newAuthFromKeystore(path, password string) (*bind.TransactOpts, error) {
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(1337)) //nolint:gomnd
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = 99999999999
	return auth, nil
}

func versionCmd(*cli.Context) error {
	fmt.Printf("Version = \"%v\"\n", version)
	fmt.Printf("Build = \"%v\"\n", commit)
	fmt.Printf("Date = \"%v\"\n", date)
	return nil
}

func registerSequencer(ctx *cli.Context) error {
	configFilePath := ctx.String(flagCfg)
	network := ctx.String(flagNetwork)
	url := ctx.Args().First()
	fmt.Print("*WARNING* Are you sure you want to register " +
		"the sequencer in the rollup using the domain <" + url + ">? [y/N]: ")
	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	input = strings.ToLower(input)
	if !(input == "y" || input == "yes") {
		return nil
	}

	c, err := config.Load(configFilePath, network)
	if err != nil {
		return err
	}

	setupLog(c.Log)

	runMigrations(c.Database)

	//Check if it is already registered
	etherman, err := newSimulatedEtherman(c.Etherman)
	if err != nil {
		log.Fatal(err)
		return err
	}
	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		log.Fatal(err)
		return err
	}
	mt := tree.NewMerkleTree(sqlDB, c.NetworkConfig.Arity, poseidon.Hash)
	tr := tree.NewStateTree(mt, []byte{})
	st := state.NewState(sqlDB, tr)
	_, err = st.GetSequencer(ctx.Context, etherman.GetAddress())
	if err == pgx.ErrNoRows { //If It doesn't exist, register the sequencer
		tx, err := etherman.RegisterSequencer(url)
		if err != nil {
			return err
		}
		log.Info("Sequencer registered. Check this tx to see the status: ", tx.Hash())
		return nil
	} else if err != nil {
		return err
	}
	// If Sequencer exists in the db
	fmt.Print("*WARNING* Sequencer is already registered. Do you want to update " +
		"the sequencer url in the rollup usign the domain <" + url + ">? [y/N]: ")
	if _, err := fmt.Scanln(&input); err != nil {
		return err
	}
	input = strings.ToLower(input)
	if !(input == "y" || input == "yes") {
		return nil
	}

	tx, err := etherman.RegisterSequencer(url)
	if err != nil {
		return err
	}
	log.Info("Sequencer updated. Check this tx to see the status: ", tx.Hash())

	return nil
}
