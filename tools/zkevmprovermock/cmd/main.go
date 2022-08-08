package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	defaultStateDBPort    = 50061
	defaultExecutorPort   = 50071
	defaultTestVectorPath = "../../test/vectors/src/merkle-tree/"
)

func main() {
	app := cli.NewApp()
	serverFlags := []cli.Flag{
		&cli.UintFlag{
			Name:     "statedb-port",
			Usage:    "StateDB server port",
			Required: false,
			Value:    defaultStateDBPort,
		},
		&cli.UintFlag{
			Name:     "executor-port",
			Usage:    "Executor server port",
			Required: false,
			Value:    defaultExecutorPort,
		},
		&cli.StringFlag{
			Name:     "test-vector-path",
			Usage:    "Test vector path",
			Required: false,
			Value:    defaultTestVectorPath,
		},
		&cli.StringFlag{
			Name:     "host",
			Usage:    "Server host",
			Required: false,
			Value:    "0.0.0.0",
		},
	}
	clientFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "state-db-serveruri",
			Usage:    "StateDB server URI",
			Required: false,
			Value:    fmt.Sprintf("127.0.0.1:%d", defaultStateDBPort),
		},
		&cli.StringFlag{
			Name:     "executor-serveruri",
			Usage:    "Executor server URI",
			Required: false,
			Value:    fmt.Sprintf("127.0.0.1:%d", defaultExecutorPort),
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Usage:  "Run zkEVM Prover mock server",
			Action: runServer,
			Flags:  serverFlags,
		},
		{
			Name:   "client",
			Usage:  "Run zkEVM Prover mock client",
			Action: runClient,
			Flags:  clientFlags,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
