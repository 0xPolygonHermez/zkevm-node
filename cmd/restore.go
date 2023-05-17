package main

import (
	"errors"
	"strings"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	pg "github.com/habx/pg-commands"
	"github.com/urfave/cli/v2"
)

const (
	restoreFlag      = "file"
)

var restoreFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     restoreFlag,
		Aliases:  []string{"i"},
		Usage:    "Input file",
		Required: true,
	},
	&configFileFlag,
}

func restore(ctx *cli.Context) error {
	// Load config
	c, err := config.Load(ctx, false)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	inputFile := ctx.String(restoreFlag)
	if !strings.Contains(inputFile, ".sql.tar.gz") {
		return errors.New("input file must end in .sql.tar.gz")
	}

	// Run migrations to create schemas and tables
	runStateMigrations(c.StateDB)

	port, err := strconv.Atoi(c.StateDB.Port)
    if err != nil {
        log.Error("error converting port to int. Error: ", err)
		return err
    }
	restore, err := pg.NewRestore(&pg.Postgres{
		Host:     c.StateDB.Host,
		Port:     port,
		DB:       c.StateDB.Name,
		Username: c.StateDB.User,
		Password: c.StateDB.Password,
	})
	if err != nil {
		log.Error("error: ", err)
		return err
	}
	restore.Role = c.StateDB.User
	restore.Schemas = append(restore.Schemas, "state")

	restoreExec := restore.Exec(inputFile, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		log.Error("error restoring snapshot. Error: ", restoreExec.Error.Err)
		log.Debug("restoreExec.Output: ", restoreExec.Output)
		return err
	}

	log.Info("Restore snapshot success")
	return nil
}
