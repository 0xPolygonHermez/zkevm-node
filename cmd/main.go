package main

import (
	"fmt"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/urfave/cli/v2"
)

const (
	// App name
	appName = "zkevm-node"
	// version represents the program based on the git tag
	version = "v0.1.0"
	// commit represents the program based on the git commit
	commit = "dev"
	// date represents the date of application was built
	date = ""
)

const (
	// AGGREGATOR is the aggregator component identifier.
	AGGREGATOR = "aggregator"
	// SEQUENCER is the sequencer component identifier.
	SEQUENCER = "sequencer"
	// RPC is the RPC component identifier.
	RPC = "rpc"
	// SYNCHRONIZER is the synchronizer component identifier.
	SYNCHRONIZER = "synchronizer"
	// BROADCAST is the broadcast component identifier.
	BROADCAST = "broadcast-trusted-state"
)

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: false,
	}
	networkFlag = cli.StringFlag{
		Name:     config.FlagNetwork,
		Aliases:  []string{"n"},
		Usage:    "Network: mainnet, testnet, internaltestnet, local, custom, merge. By default it uses mainnet",
		Required: false,
	}
	customNetworkFlag = cli.StringFlag{
		Name:    config.FlagNetworkCfg,
		Aliases: []string{"nc"},
		Usage:   "Custom network configuration `FILE` when using --network custom parameter",
	}
	baseNetworkFlag = cli.StringFlag{
		Name:     config.FlagNetworkBase,
		Aliases:  []string{"nb"},
		Usage:    "Base existing network configuration to be merged with the custom configuration passed with --network-cfg, by default it uses internaltestnet",
		Value:    "internaltestnet",
		Required: false,
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
		Usage:    fmt.Sprintf("List of JSON RPC apis to be exposed by the server: --http.api=%v,%v,%v,%v,%v,%v", jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIDebug, jsonrpc.APIHez, jsonrpc.APITxPool, jsonrpc.APIWeb3),
		Required: false,
		Value:    cli.NewStringSlice(jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIHez, jsonrpc.APITxPool, jsonrpc.APIWeb3),
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version
	flags := []cli.Flag{
		&configFileFlag,
		&networkFlag,
		&customNetworkFlag,
		&baseNetworkFlag,
		&yesFlag,
		&componentsFlag,
		&httpAPIFlag,
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
			Usage:   "Run the zkevm-node",
			Action:  start,
			Flags:   flags,
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
