package main

import (
	"fmt"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/urfave/cli/v2"
)

const (
	// App name
	appName = "zkevm-node"
	// version represents the program based on the git tag
	version = "v0.1.0"
	// date represents the date of application was built
	date = ""
)

const (
	// AGGREGATOR is the aggregator component identifier.
	AGGREGATOR = "aggregator"
	// AGGREGATOR2 is the aggregator v2 component identifier.
	AGGREGATOR2 = "aggregator2"
	// SEQUENCER is the sequencer component identifier.
	SEQUENCER = "sequencer"
	// RPC is the RPC component identifier.
	RPC = "rpc"
	// SYNCHRONIZER is the synchronizer component identifier.
	SYNCHRONIZER = "synchronizer"
	// BROADCAST is the broadcast component identifier.
	BROADCAST = "broadcast-trusted-state"
)

const (
	//envCommitHash environment variable name for COMMIT_HASH
	envCommitHash = "COMMIT_HASH"
)

var (
	// commit represents the program based on the git commit
	commit         = testutils.GetEnv(envCommitHash, "dev")
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: false,
	}
	genesisFlag = cli.StringFlag{
		Name:     config.FlagGenesisFile,
		Aliases:  []string{"gen"},
		Usage:    "Loads the genesis `FILE`",
		Required: true,
	}
	yesFlag = cli.BoolFlag{
		Name:     config.FlagYes,
		Aliases:  []string{"y"},
		Usage:    "Automatically accepts any confirmation to execute the command",
		Required: false,
	}
	componentsFlag = cli.StringSliceFlag{
		Name:     config.FlagComponents,
		Aliases:  []string{"co"},
		Usage:    "List of components to run",
		Required: false,
		Value:    cli.NewStringSlice(AGGREGATOR, SEQUENCER, RPC, SYNCHRONIZER),
	}
	httpAPIFlag = cli.StringSliceFlag{
		Name:     config.FlagHTTPAPI,
		Aliases:  []string{"ha"},
		Usage:    fmt.Sprintf("List of JSON RPC apis to be exposed by the server: --http.api=%v,%v,%v,%v,%v,%v", jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIDebug, jsonrpc.APIZKEVM, jsonrpc.APITxPool, jsonrpc.APIWeb3),
		Required: false,
		Value:    cli.NewStringSlice(jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIZKEVM, jsonrpc.APITxPool, jsonrpc.APIWeb3),
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version
	flags := []cli.Flag{
		&configFileFlag,
		&yesFlag,
		&componentsFlag,
		&httpAPIFlag,
	}
	log.Infof("Starting application [Commit Hash: %s, Version: %s] ...", commit, version)
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
			Usage:   "Run the zkevm-node",
			Action:  start,
			Flags:   append(flags, &genesisFlag),
		},
		{
			Name:    "approve",
			Aliases: []string{"ap"},
			Usage:   "Approve tokens to be spent by the smart contract",
			Action:  approveTokens,
			Flags: append(flags, &cli.StringFlag{
				Name:     config.FlagAmount,
				Aliases:  []string{"am"},
				Usage:    "Amount that is gonna be approved",
				Required: true,
			},
			),
		},
		{
			Name:    "encryptKey",
			Aliases: []string{},
			Usage:   "Encrypts the privatekey with a password and create a keystore file",
			Action:  encryptKey,
			Flags:   encryptKeyFlags,
		},
		{
			Name:    "dumpState",
			Aliases: []string{},
			Usage:   "Dumps the state in a JSON file, for debug purposes",
			Action:  dumpState,
			Flags:   dumpStateFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
