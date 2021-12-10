package main

import (
	"os"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/urfave/cli/v2"
)

const (
	flagCfg = "cfg"
)

func main() {
	app := cli.NewApp()
	app.Name = "hezcore"
	app.Version = version

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     flagCfg,
			Usage:    "Configuration `FILE`",
			Required: false,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Show the application version and build",
			Action:  cmdVersion,
		},
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Runs the application",
			Action:  cmdRun,
			Flags:   flags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
