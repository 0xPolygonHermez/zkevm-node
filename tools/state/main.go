package main

import (
	"log"
	"os"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/urfave/cli/v2"
)

const (
	appName = "zkevm-statedb-tool"
)

const (
	flagChainID = "l2_chain_id"
)

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: false,
	}
	configChainIDFlag = cli.StringFlag{
		Name:     flagChainID,
		Usage:    "forced L2 chain id instead of asking SMC",
		Required: false,
	}
	networkFlag = cli.StringFlag{
		Name:     config.FlagNetwork,
		Aliases:  []string{"net"},
		Usage:    "Load default network configuration. Supported values: [`custom`]",
		Required: false,
	}
	customNetworkFlag = cli.StringFlag{
		Name:     config.FlagCustomNetwork,
		Aliases:  []string{"net-file"},
		Usage:    "Load the network configuration file if --network=custom",
		Required: false,
	}
	firstBatchNumberFlag = cli.StringFlag{
		Name:     "first_batch_number",
		Aliases:  []string{"start"},
		Usage:    "First batch number (default:1)",
		Required: false,
	}
	lastBatchNumberFlag = cli.StringFlag{
		Name:     "last_batch_number",
		Aliases:  []string{"end"},
		Usage:    "Last batch number (default:last one on batch table)",
		Required: false,
	}
	writeOnHashDBFlag = cli.BoolFlag{
		Name:     "write_on_hash_db",
		Usage:    "When process batches say to exectuor to write on the MT in a persistent way (default:false)",
		Required: false,
	}
	dontStopOnErrorFlag = cli.BoolFlag{
		Name:     "dont_stop_on_error",
		Usage:    "Keep processing even if a batch have an error (default:false)",
		Required: false,
	}
	preferExecutionStateRootFlag = cli.BoolFlag{
		Name:     "prefer_execution_state_root",
		Usage:    "Instaed of using the state_root from previous batch use the stateRoot from previous execution (default:false)",
		Required: false,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = zkevm.Version

	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Application version and build",
			Action:  versionCmd,
		},
		{
			Name:    "reprocess",
			Aliases: []string{},
			Usage:   "reprocess batches",
			Action:  reprocessCmd,
			Flags: []cli.Flag{&configFileFlag, &networkFlag, &customNetworkFlag, &configChainIDFlag, &firstBatchNumberFlag,
				&lastBatchNumberFlag, &writeOnHashDBFlag, &dontStopOnErrorFlag, &preferExecutionStateRootFlag},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
