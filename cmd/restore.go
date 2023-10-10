package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	pg "github.com/habx/pg-commands"
	"github.com/urfave/cli/v2"
)

const (
	restorestateDbFlag = "inputfilestate"
	restoreHashDbFlag  = "inputfileHash"
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

	d, err := db.NewSQLDB(c.State.DB)
	if err != nil {
		log.Error("error conecting to stateDB. Error: ", err)
		return err
	}
	_, err = d.Exec(ctx.Context, "DROP SCHEMA IF EXISTS state CASCADE; DROP TABLE IF EXISTS gorp_migrations;")
	if err != nil {
		log.Error("error dropping state schema or migration table. Error: ", err)
		return err
	}
	port, err := strconv.Atoi(c.State.DB.Port)
	if err != nil {
		log.Error("error converting port to int. Error: ", err)
		return err
	}
	restore, err := pg.NewRestore(&pg.Postgres{
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
	params := []string{"--no-owner", "--no-acl", "--format=c"}
	log.Info("Restore stateDB snapshot started, please wait...")
	restoreExec := execCommand(restore, inputFileStateDB, pg.ExecOptions{StreamPrint: false}, params)
	if restoreExec.Error != nil {
		log.Error("error restoring stateDB snapshot. Error: ", restoreExec.Error.Err)
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
	d, err = db.NewSQLDB(c.HashDB)
	if err != nil {
		log.Error("error conecting to hashdb. Error: ", err)
		return err
	}
	_, err = d.Exec(ctx.Context, "DROP SCHEMA IF EXISTS state CASCADE;")
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

	log.Info("Restore HashDB snapshot started, please wait...")
	restoreExec = execCommand(restore, inputFileHashDB, pg.ExecOptions{StreamPrint: false}, params)
	if restoreExec.Error != nil {
		log.Error("error restoring hashDB snapshot. Error: ", restoreExec.Error.Err)
		log.Debug("restoreExec.Output: ", restoreExec.Output)
		return err
	}

	log.Info("Restore HashDB snapshot success")
	return nil
}

func execCommand(x *pg.Restore, filename string, opts pg.ExecOptions, params []string) pg.Result {
	result := pg.Result{}
	options := append(params, x.Postgres.Parse()...)
	options = append(options, fmt.Sprintf("%s%s", x.Path, filename))
	log.Debug("Options: ", options)

	result.FullCommand = strings.Join(options, " ")
	cmd := exec.Command(pg.PGRestoreCmd, options...) //nolint:gosec
	cmd.Env = append(os.Environ(), x.EnvPassword)
	stderrIn, _ := cmd.StderrPipe()
	go func(stderrIn io.ReadCloser, opts pg.ExecOptions, result *pg.Result) {
		output := ""
		reader := bufio.NewReader(stderrIn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					result.Output = output
					break
				}
				result.Error = &pg.ResultError{Err: fmt.Errorf("error reading output: %w", err), CmdOutput: output}
				break
			}

			if opts.StreamPrint {
				_, err = fmt.Fprint(opts.StreamDestination, line)
				if err != nil {
					result.Error = &pg.ResultError{Err: fmt.Errorf("error writing output: %w", err), CmdOutput: output}
					break
				}
			}

			output += line
		}
	}(stderrIn, opts, &result)
	err := cmd.Start()
	if err != nil {
		result.Error = &pg.ResultError{Err: err, CmdOutput: result.Output}
	}
	err = cmd.Wait()
	if exitError, ok := err.(*exec.ExitError); ok {
		result.Error = &pg.ResultError{Err: err, ExitCode: exitError.ExitCode(), CmdOutput: result.Output}
	}

	return result
}
