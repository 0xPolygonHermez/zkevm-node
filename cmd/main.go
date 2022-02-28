package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	flagYes        = "yes"
	flagCfg        = "cfg"
	flagNetwork    = "network"
	flagNetworkCfg = "network-cfg"
	flagAmount     = "amount"
	flagRemoteMT   = "remote-merkletree"
)

const (
	// App name
	appName = "hermez-node"
	// version represents the program based on the git tag
	version = "v0.1.0"
	// commit represents the program based on the git commit
	commit = "dev"
	// date represents the date of application was built
	date = ""
)

func main() {
	app := cli.NewApp()
	app.Name = appName
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
		&cli.StringFlag{
			Name:    flagNetworkCfg,
			Aliases: []string{"nc"},
			Usage:   "Custom network configuration `FILE` when using --network custom parameter",
		},
		&cli.BoolFlag{
			Name:     flagYes,
			Aliases:  []string{"y"},
			Usage:    "Automatically accepts any confirmation to execute the command",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     flagRemoteMT,
			Aliases:  []string{"mt"},
			Usage:    "Connect to merkletree service instead of use local libraries",
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
		{
			Name:    "approve",
			Aliases: []string{"ap"},
			Usage:   "Approve tokens to be spent by the smart contract",
			Action:  approveTokens,
			Flags: append(flags, &cli.StringFlag{
				Name:     flagAmount,
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
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}
