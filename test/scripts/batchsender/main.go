package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func main() {
	// send 1 batch by default or read the number of batches from args
	nBatches := 1
	usage := `Usage: batchsender <nbatches>
 nbatches	Number of batches to be sent (defaults to 1)`
	if len(os.Args) > 1 {
		var err error
		arg := os.Args[1]
		if arg == "--help" || arg == "-h" {
			fmt.Println(usage)
			os.Exit(0)
		}
		nBatches, err = strconv.Atoi(arg)
		if err != nil {
			fmt.Println(usage)
			os.Exit(1)
		}
	}

	// retrieve default configuration
	var cfg config.Config
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(config.DefaultValues)))
	checkErr(err)
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	checkErr(err)

	client, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	checkErr(err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	checkErr(err)

	ethMan, err := etherman.NewClient(cfg.Etherman, auth)
	checkErr(err)

	seqAddr, err := ethMan.GetPublicAddress()
	checkErr(err)
	log.Info("Using address: ", seqAddr)

	ctx := context.Background()

	for i := 0; i < nBatches; i++ {
		currentBlock, err := client.BlockByNumber(ctx, nil)
		checkErr(err)
		log.Debug("currentBlock.Time(): ", currentBlock.Time())

		seqs := []types.Sequence{{
			GlobalExitRoot: common.HexToHash("0x"),
			Txs:            []ethtypes.Transaction{},
			Timestamp:      int64(currentBlock.Time() - 1), // fit in latest-sequence < > current-block rage
		}}
		tx, err := ethMan.SequenceBatches(ctx, seqs, 0, nil, nil)
		checkErr(err)
		log.Info("TxHash: ", tx.Hash())

		var duration time.Duration
		if nBatches > 1 {
			duration = 1
		}
		time.Sleep(duration * time.Second)
	}
	if nBatches > 1 {
		log.Infof("Successfully sent %d batches", nBatches)
	} else {
		log.Info("Successfully sent 1 batch")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
