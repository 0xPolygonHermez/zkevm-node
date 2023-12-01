package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/urfave/cli/v2"
)

const (
	dumpStateFlagDescription = "description"
	dumpStateFlagOutput      = "output"
)

var dumpStateFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     dumpStateFlagDescription,
		Aliases:  []string{"desc"},
		Usage:    "Description of the state being dumped",
		Required: true,
	},
	&cli.StringFlag{
		Name:     dumpStateFlagOutput,
		Aliases:  []string{"o"},
		Usage:    "Output file to save the dump, should end in .json",
		Required: true,
	},
	&configFileFlag,
	&networkFlag,
}

type dumpedState struct {
	Description string
	Genesis     genesis
	Batches     []batchMeta
}

type genesis state.Genesis

func (g genesis) MarshalJSON() ([]byte, error) {
	for _, action := range g.Actions {
		if !strings.HasPrefix(action.Value, "0x") {
			action.Value = fmt.Sprintf("0x%s", action.Value)
		}
		if action.Bytecode != "" && !strings.HasPrefix(action.Bytecode, "0x") {
			action.Bytecode = fmt.Sprintf("0x%s", action.Bytecode)
		}
		if action.StoragePosition != "" && !strings.HasPrefix(action.StoragePosition, "0x") {
			action.StoragePosition = fmt.Sprintf("0x%s", action.StoragePosition)
		}
	}

	// Create JSON
	type Alias genesis
	return json.Marshal(&struct {
		Alias
		Actions []*state.GenesisAction
	}{
		Alias:   (Alias)(g),
		Actions: g.Actions,
	})
}

type batchMeta struct {
	Virtualized  bool
	Consolidated bool
	Batch        batch
}

type batch state.Batch

func (b *batch) MarshalJSON() ([]byte, error) {
	type Alias batch
	var strDataPtr *string
	if b.BatchL2Data != nil {
		strData := "0x" + hex.EncodeToString(b.BatchL2Data)
		strDataPtr = &strData
	}
	return json.Marshal(&struct {
		*Alias
		Timestamp   int64
		BatchL2Data *string
	}{
		Alias:       (*Alias)(b),
		Timestamp:   b.Timestamp.Unix(),
		BatchL2Data: strDataPtr,
	})
}

func dumpState(ctx *cli.Context) error {
	// Load config
	c, err := config.Load(ctx, true)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	description := ctx.String(dumpStateFlagDescription)
	outputFile := ctx.String(dumpStateFlagOutput)
	if !strings.Contains(outputFile, ".json") {
		return errors.New("output file must end in .json")
	}

	// Connect to SQL
	stateSqlDB, err := db.NewSQLDB(c.State.DB)
	if err != nil {
		return err
	}
	stateDB := pgstatestorage.NewPostgresStorage(state.Config{}, stateSqlDB)

	dump := dumpedState{
		Description: description,
		Genesis:     genesis(c.NetworkConfig.Genesis),
	}

	// Load batches
	dbCtx := context.Background()
	lastBatchNum, err := stateDB.GetLastBatchNumber(dbCtx, nil)
	if err != nil {
		return err
	}
	lastVirtualBatchNum, err := stateDB.GetLastVirtualBatchNum(dbCtx, nil)
	if err != nil {
		return err
	}
	lastConsolidatedBatch, err := stateDB.GetLastVerifiedBatch(dbCtx, nil)
	if err != nil {
		return err
	}
	lastConsolidatedBatchNum := lastConsolidatedBatch.BatchNumber
	for i := uint64(0); i <= lastBatchNum; i++ {
		b, err := stateDB.GetBatchByNumber(dbCtx, i, nil)
		if err != nil {
			return err
		}
		isClosed, err := stateDB.IsBatchClosed(dbCtx, i, nil)
		if err != nil {
			return err
		}
		if !isClosed {
			continue
		}
		dump.Batches = append(dump.Batches, batchMeta{
			Virtualized:  i <= lastVirtualBatchNum,
			Consolidated: i <= lastConsolidatedBatchNum,
			Batch:        batch(*b),
		})
	}

	// Dump JSON
	file, err := json.MarshalIndent(dump, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, file, 0600) //nolint:gomnd
}
