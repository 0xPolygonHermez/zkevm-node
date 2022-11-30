package main

import (
	"bytes"
	"errors"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	flagWaitName    = "wait"
	flagVerboseName = "verbose"
)

var (
	sequencedBatchesEventSignatureHash = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	verifiedBatchSignatureHash         = crypto.Keccak256Hash([]byte("VerifyBatch(uint64,address)"))

	flagWait = cli.BoolFlag{
		Name:     flagWaitName,
		Aliases:  []string{"w"},
		Usage:    "wait batch transaction to be confirmed",
		Required: false,
	}
	flagVerbose = cli.BoolFlag{
		Name:     flagVerboseName,
		Aliases:  []string{"v"},
		Usage:    "output verbose logs",
		Required: false,
	}
)

func main() {
	batchsender := cli.NewApp()
	batchsender.Name = "batchsender"
	batchsender.Usage = "send batch transactions to L1"
	batchsender.Description = `This tool allows to send a specified number of batch transactions to L1. 
Optionally it can wait for the batches to be validated.`
	batchsender.DefaultCommand = "send"
	batchsender.Flags = []cli.Flag{&flagWait, &flagVerbose}
	batchsender.Commands = []*cli.Command{
		{
			Before:  setLogLevel,
			Name:    "send",
			Aliases: []string{},
			Usage:   "Sends the specified number of batch transactions to L1",
			Description: `This command sends the specified number of transactions to L1.
If --wait flag is used, it waits for the corresponding validation transaction.`,
			ArgsUsage: "number of batches to be sent (default: 1)",
			Action:    sendBatches,
		},
	}

	err := batchsender.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func setLogLevel(ctx *cli.Context) error {
	logLevel := "info"
	if ctx.Bool(flagVerboseName) {
		logLevel = "debug"
	}

	log.Init(log.Config{Level: logLevel, Outputs: []string{"stdout"}})
	return nil
}

func sendBatches(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	nBatches := 1 // send 1 batch by default
	if cliCtx.NArg() > 0 {
		nBatchesArgStr := cliCtx.Args().Get(0)
		nBatchesArg, err := strconv.Atoi(nBatchesArgStr)
		if err == nil {
			nBatches = nBatchesArg
		}
	}

	// retrieve default configuration
	var cfg config.Config
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(config.DefaultValues)))
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return err
	}

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}

	ethMan, err := etherman.NewClient(cfg.Etherman, auth)
	if err != nil {
		return err
	}

	seqAddr, err := ethMan.GetPublicAddress()
	if err != nil {
		return err
	}
	log.Info("Using address: ", seqAddr)

	wait := cliCtx.Bool(flagWaitName)

	var sentTxs []*ethtypes.Transaction
	sentTxsMap := make(map[common.Hash]struct{})

	var duration time.Duration
	if nBatches > 1 {
		duration = 500
	}

	for i := 0; i < nBatches; i++ {
		currentBlock, err := ethMan.EtherClient.BlockByNumber(ctx, nil)
		if err != nil {
			return err
		}
		log.Debug("currentBlock.Time(): ", currentBlock.Time())

		seqs := []ethmanTypes.Sequence{{
			GlobalExitRoot: common.HexToHash("0x"),
			Txs:            []ethtypes.Transaction{},
			Timestamp:      int64(currentBlock.Time() - 1), // fit in latest-sequence < > current-block rage
		}}

		// send empty rollup to L1
		tx, err := ethMan.SequenceBatches(ctx, seqs, 0, nil, nil)
		if err != nil {
			return err
		}

		log.Info("TxHash: ", tx.Hash())
		sentTxs = append(sentTxs, tx)
		sentTxsMap[tx.Hash()] = struct{}{}

		time.Sleep(duration * time.Millisecond)
	}

	sentBatches := len(sentTxs)

	if wait { // wait proofs
		log.Info("Waiting for txs to be confirmed...")
		time.Sleep(time.Second)

		virtualBatches := make(map[uint64]common.Hash)
		verifiedBatches := make(map[uint64]struct{})
		loggedBatches := make(map[uint64]struct{})

		miningTimeout := 180 * time.Second                           //nolint:gomnd
		waitTimeout := time.Duration(180*len(sentTxs)) * time.Second //nolint:gomnd
		done := make(chan struct{})

		for _, tx := range sentTxs {
			err := operations.WaitTxToBeMined(ctx, ethMan.EtherClient, tx, miningTimeout)
			if err != nil {
				return err
			}
		}

		for {
			select {
			case <-time.After(waitTimeout):
				return errors.New("Deadline exceeded")
			case <-done:
				success(sentBatches)
				return nil
			default:
			txLoop:
				for _, tx := range sentTxs {
					// get rollup tx block number
					receipt, err := ethMan.EtherClient.TransactionReceipt(ctx, tx.Hash())
					if err != nil {
						return err
					}

					fromBlock := receipt.BlockNumber
					toBlock := new(big.Int).Add(fromBlock, new(big.Int).SetUint64(cfg.Synchronizer.SyncChunkSize))
					query := ethereum.FilterQuery{
						FromBlock: fromBlock,
						ToBlock:   toBlock,
						Addresses: ethMan.SCAddresses,
					}
					logs, err := ethMan.EtherClient.FilterLogs(ctx, query)
					if err != nil {
						return err
					}
					for _, vLog := range logs {
						switch vLog.Topics[0] {
						case sequencedBatchesEventSignatureHash:
							if vLog.TxHash == tx.Hash() { // ignore other txs happening on L1
								sb, err := ethMan.PoE.ParseSequenceBatches(vLog)
								if err != nil {
									return err
								}

								virtualBatches[sb.NumBatch] = vLog.TxHash

								if _, logged := loggedBatches[sb.NumBatch]; !logged {
									log.Infof("Batch [%d] virtualized in TxHash [%v]", sb.NumBatch, vLog.TxHash)
									loggedBatches[sb.NumBatch] = struct{}{}
								}
							}
						case verifiedBatchSignatureHash:
							vb, err := ethMan.PoE.ParseVerifyBatch(vLog)
							if err != nil {
								return err
							}

							if _, verified := verifiedBatches[vb.NumBatch]; !verified {
								log.Infof("Batch [%d] verified in TxHash [%v]", vb.NumBatch, vLog.TxHash)
								verifiedBatches[vb.NumBatch] = struct{}{}
							}

							// batch is verified, remove it from the txs set
							delete(sentTxsMap, virtualBatches[vb.NumBatch])
							if len(sentTxsMap) == 0 {
								close(done)
								break txLoop
							}
						}
					}

					// wait for verifications
					time.Sleep(time.Second) //nolint:gomnd
				}
			}
		}
	}

	success(sentBatches)
	return nil
}

func success(nBatches int) {
	if nBatches > 1 {
		log.Infof("Successfully sent %d batches", nBatches)
	} else {
		log.Info("Successfully sent 1 batch")
	}
}
