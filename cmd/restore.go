package main

import (
	"errors"
	"strings"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/db"
	pg "github.com/habx/pg-commands"
	"github.com/urfave/cli/v2"
)

const (
	restorestateDbFlag      = "inputfilestate"
	restoreHashDbFlag      = "inputfileHash"
)

var restoreFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     restorestateDbFlag,
		Aliases:  []string{"is"},
		Usage:    "Input file stateDB",
		Required: true,
	},
	&cli.StringFlag{
		Name:     restoreHashDbFlag,
		Aliases:  []string{"ih"},
		Usage:    "Input file hashDB",
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
	inputFileStateDB := ctx.String(restorestateDbFlag)
	if !strings.Contains(inputFileStateDB, ".sql.tar.gz") {
		return errors.New("stateDB input file must end in .sql.tar.gz")
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
	log.Info("Restore stateDB snapshot started, please wait...")
	restoreExec := restore.Exec(inputFileStateDB, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		log.Error("error restoring snapshot. Error: ", restoreExec.Error.Err)
		log.Debug("restoreExec.Output: ", restoreExec.Output)
		return err
	}
	log.Info("Restore stateDB snapshot success")

	inputFileHashDB := ctx.String(restoreHashDbFlag)
	if !strings.Contains(inputFileHashDB, ".sql.tar.gz") {
		return errors.New("hashDb input file must end in .sql.tar.gz")
	}
	port, err = strconv.Atoi(c.HashDB.Port)
    if err != nil {
        log.Error("error converting port to int. Error: ", err)
		return err
    }
	d, err := db.NewSQLDB(c.HashDB)
	if err != nil {
		log.Error("error conecting to hashdb. Error: ", err)
		return err
	}
	_, err = d.Exec(ctx.Context, "DROP SCHEMA IF EXISTS state; CREATE SCHEMA IF NOT EXISTS state;")
	if err != nil {
		log.Error("error dropping and creating state schema. Error: ", err)
		return err
	}
	restore, err = pg.NewRestore(&pg.Postgres{
		Host:     c.HashDB.Host,
		Port:     port,
		DB:       c.HashDB.Name,
		Username: c.HashDB.User,
		Password: c.HashDB.Password,
	})
	if err != nil {
		log.Error("error: ", err)
		return err
	}
	restore.Role = c.HashDB.User
	restore.Schemas = append(restore.Schemas, "state")
	restore.Options = []string{"--no-owner", "--no-acl"}
	log.Info("Restore HashDB snapshot started, please wait...")
	restoreExec = restore.Exec(inputFileHashDB, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		log.Error("error restoring snapshot. Error: ", restoreExec.Error.Err)
		log.Debug("restoreExec.Output: ", restoreExec.Output)
		return err
	}

	log.Info("Restore HashDB snapshot success")
	return nil
}
