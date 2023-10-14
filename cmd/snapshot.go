package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	pg "github.com/habx/pg-commands"
	"github.com/urfave/cli/v2"
)

var snapshotFlags = []cli.Flag{
	&configFileFlag,
	&outputFileFlag,
}

func snapshot(ctx *cli.Context) error {
	// Load config
	c, err := config.Load(ctx, false)
	if err != nil {
		log.Error("error loading configuration. Error: ", err)
		return err
	}
	setupLog(c.Log)

	port, err := strconv.Atoi(c.State.DB.Port)
	if err != nil {
		log.Error("error converting port to int. Error: ", err)
		return err
	}
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     c.State.DB.Host,
		Port:     port,
		DB:       c.State.DB.Name,
		Username: c.State.DB.User,
		Password: c.State.DB.Password,
	})
	if err != nil {
		log.Error("error: ", err)
		return err
	}
	dump.Options = append(dump.Options, "-Z 9")
	log.Info("StateDB snapshot is being created...")
	dump.Path = ctx.String(config.FlagOutputFile)
	dump.SetFileName(fmt.Sprintf(`%v_%v_%v_%v.sql.tar.gz`, dump.DB, time.Now().Unix(), zkevm.Version, zkevm.GitRev))
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		log.Error("error dumping statedb. Error: ", dumpExec.Error.Err)
		log.Debug("dumpExec.Output: ", dumpExec.Output)
		return err
	}

	log.Info("StateDB snapshot success. Saved in ", dumpExec.File)

	port, err = strconv.Atoi(c.HashDB.Port)
	if err != nil {
		log.Error("error converting port to int. Error: ", err)
		return err
	}
	dump, err = pg.NewDump(&pg.Postgres{
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
	dump.Options = append(dump.Options, "-Z 9")
	log.Info("HashDB snapshot is being created...")
	dump.Path = ctx.String(config.FlagOutputFile)
	dump.SetFileName(fmt.Sprintf(`%v_%v_%v_%v.sql.tar.gz`, dump.DB, time.Now().Unix(), zkevm.Version, zkevm.GitRev))
	dumpExec = dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		log.Error("error dumping hashdb. Error: ", dumpExec.Error.Err)
		log.Debug("dumpExec.Output: ", dumpExec.Output)
		return err
	}

	log.Info("HashDB snapshot success. Saved in ", dumpExec.File)
	return nil
}
