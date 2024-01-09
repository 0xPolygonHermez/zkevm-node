package main

import (
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonrollupmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	flagL1URLName   = "url"
	flagZkevmAddrName = "zkevm"
	flagRollupManagerAddrName = "rollupmanager"
	miningTimeout   = 180
)

var (
	flagL1URL = cli.StringFlag{
		Name:     flagL1URLName,
		Aliases:  []string{"u"},
		Usage:    "L1 node url",
		Required: true,
	}
	flagZkevmAddr = cli.StringFlag{
		Name:     flagZkevmAddrName,
		Aliases:  []string{"zk"},
		Usage:    "Zkevm smart contract address",
		Required: true,
	}
	flagRollupManagerAddr = cli.StringFlag{
		Name:     flagRollupManagerAddrName,
		Aliases:  []string{"r"},
		Usage:    "RollupmManager smart contract address",
		Required: true,
	}
)

func main() {
	fbatchsender := cli.NewApp()
	fbatchsender.Name = "forcedBatchsender"
	fbatchsender.Usage = "send forced batch transactions to L1"
	fbatchsender.DefaultCommand = "send"
	flags := []cli.Flag{&flagL1URL, &flagZkevmAddr, &flagRollupManagerAddr}
	fbatchsender.Commands = []*cli.Command{
		{
			Before:  setLogLevel,
			Name:    "send",
			Aliases: []string{},
			Flags:   flags,
			Action:  sendForcedBatches,
		},
	}

	err := fbatchsender.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setLogLevel(ctx *cli.Context) error {
	logLevel := "debug"
	log.Init(log.Config{
		Level:   logLevel,
		Outputs: []string{"stderr"},
	})
	return nil
}

func sendForcedBatches(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	url := cliCtx.String(flagL1URLName)
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", url, err)
		return err
	}
	// Create smc client
	zkevmAddr := common.HexToAddress(cliCtx.String(flagZkevmAddrName))
	zkevm, err := polygonzkevm.NewPolygonzkevm(zkevmAddr, ethClient)
	if err != nil {
		return err
	}

	rollupManagerAddr := common.HexToAddress(cliCtx.String(flagRollupManagerAddrName))
	rollupManager, err := polygonrollupmanager.NewPolygonrollupmanager(rollupManagerAddr, ethClient)
	if err != nil {
		return err
	}

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}

	log.Info("Using address: ", auth.From)

	num, err := zkevm.LastForceBatch(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("error getting lastForBatch number. Error : ", err)
		return err
	}
	log.Info("Number of forceBatches in the smc: ", num)

	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		log.Error("error getting blockByNumber. Error: ", err)
		return err
	}
	log.Debug("currentBlock.Time(): ", currentBlock.Time())

	// Get tip
	tip, err := rollupManager.GetForcedBatchFee(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("error getting tip. Error: ", err)
		return err
	}

	// Allow forced batches in smart contract if disallowed
	disallowed, err := zkevm.IsForcedBatchAllowed(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("error getting IsForcedBatchAllowed. Error: ", err)
		return err
	}
	if disallowed {
		tx, err := zkevm.ActivateForceBatches(auth)
		if err != nil {
			log.Error("error sending activateForceBatches. Error: ", err)
			return err
		}
		err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
		if err != nil {

			log.Error("error waiting tx to be mined. Error: ", err)
			return err
		}
	}

	// Send forceBatch
	tx, err := zkevm.ForceBatch(auth, []byte{}, tip)
	if err != nil {
		log.Error("error sending forceBatch. Error: ", err)
		return err
	}

	log.Info("TxHash: ", tx.Hash())

	time.Sleep(1 * time.Second)

	err = operations.WaitTxToBeMined(ctx, ethClient, tx, miningTimeout*time.Second)
	if err != nil {
		return err
	}

	query := ethereum.FilterQuery{
		FromBlock: currentBlock.Number(),
		Addresses: []common.Address{zkevmAddr},
	}
	logs, err := ethClient.FilterLogs(ctx, query)
	if err != nil {
		return err
	}
	for _, vLog := range logs {
		fb, err := zkevm.ParseForceBatch(vLog)
		if err == nil {
			log.Debugf("log decoded: %+v", fb)
			ger := fb.LastGlobalExitRoot
			log.Info("GlobalExitRoot: ", ger)
			log.Info("Transactions: ", common.Bytes2Hex(fb.Transactions))
			fullBlock, err := ethClient.BlockByHash(ctx, vLog.BlockHash)
			if err != nil {
				log.Errorf("error getting hashParent. BlockNumber: %d. Error: %v", vLog.BlockNumber, err)
				return err
			}
			log.Info("MinForcedTimestamp: ", fullBlock.Time())
		}
	}

	return nil
}
