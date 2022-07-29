package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
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
	&customNetworkFlag,
	&baseNetworkFlag,
}

type dumpedState struct {
	Description string
	Genesis     genesis
	Batches     []batchMeta
}

type genesis state.Genesis

func (g genesis) MarshalJSON() ([]byte, error) {
	type Alias genesis
	contractsHex := map[common.Address]string{}
	for addr, code := range g.SmartContracts {
		contractsHex[addr] = "0x" + hex.EncodeToString(code)
	}
	storageHex := map[common.Address]map[string]string{}
	for addr, storage := range g.Storage {
		addrStorage := map[string]string{}
		for position, value := range storage {
			addrStorage["0x"+position.Text(16)] = "0x" + value.Text(16)
		}
		storageHex[addr] = addrStorage
	}
	return json.Marshal(&struct {
		Alias
		SmartContracts map[common.Address]string            `json:"smartContracts"`
		Storage        map[common.Address]map[string]string `json:"storage"`
	}{
		Alias:          (Alias)(g),
		SmartContracts: contractsHex,
		Storage:        storageHex,
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
	c, err := config.Load(ctx)
	if err != nil {
		return err
	}
	setupLog(c.Log)
	description := ctx.String(dumpStateFlagDescription)
	outputFile := ctx.String(dumpStateFlagOutput)
	if !strings.Contains(outputFile, ".json") {
		return errors.New("Output file must end in .json")
	}

	// Connect to SQL
	sqlDB, err := db.NewSQLDB(c.Database)
	if err != nil {
		return err
	}
	stateDB := state.NewPostgresStorage(sqlDB)

	// Load g
	g := state.Genesis{
		Balances:       c.NetworkConfig.Genesis.Balances,
		SmartContracts: c.NetworkConfig.Genesis.SmartContracts,
		Storage:        c.NetworkConfig.Genesis.Storage,
		Nonces:         c.NetworkConfig.Genesis.Nonces,
	}
	dump := dumpedState{
		Description: description,
		Genesis:     genesis(g),
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
	return ioutil.WriteFile(outputFile, file, 0600) //nolint:gomnd
}
