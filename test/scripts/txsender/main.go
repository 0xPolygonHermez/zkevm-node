package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	flagNetworkName       = "network"
	flagWaitName          = "wait"
	flagVerboseName       = "verbose"
	flagNetworkLocalL1Key = "l1"
	flagNetworkLocalL2Key = "l2"
	defaultNetwork        = flagNetworkLocalL2Key

	receiverAddr = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
)

type networkLayer int

const (
	networkLayer1 networkLayer = iota
	networkLayer2
)

var networks = map[string]struct {
	Name       string
	URL        string
	ChainID    uint64
	PrivateKey string
	networkLayer
}{
	flagNetworkLocalL1Key: {Name: "Local L1", URL: operations.DefaultL1NetworkURL, ChainID: operations.DefaultL1ChainID, PrivateKey: operations.DefaultSequencerPrivateKey, networkLayer: networkLayer1},
	flagNetworkLocalL2Key: {Name: "Local L2", URL: operations.DefaultL2NetworkURL, ChainID: operations.DefaultL2ChainID, PrivateKey: operations.DefaultSequencerPrivateKey, networkLayer: networkLayer2},
}

var (
	flagNetwork = cli.StringSliceFlag{
		Name:     flagNetworkName,
		Aliases:  []string{"n"},
		Usage:    "list of networks on which to send transactions",
		Required: false,
	}
	flagWait = cli.BoolFlag{
		Name:     flagWaitName,
		Aliases:  []string{"w"},
		Usage:    "wait transactions to be confirmed",
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
	txsender := cli.NewApp()
	txsender.Name = "txsender"
	txsender.Flags = []cli.Flag{&flagNetwork, &flagWait, &flagVerbose}
	txsender.Usage = "send transactions"
	txsender.Description = `This tool allows to send a specified number of transactions.
Optionally it can wait for the transactions to be validated.`
	txsender.DefaultCommand = "send"
	txsender.Commands = []*cli.Command{
		{
			Name:    "send",
			Before:  setLogLevel,
			Aliases: []string{},
			Usage:   "Sends the specified number of transactions",
			Description: `This command sends the specified number of transactions.
If --wait flag is used, it waits for the corresponding validation transaction.`,
			ArgsUsage: "number of transactions to be sent (default: 1)",
			Action:    sendTxs,
		},
	}

	err := txsender.Run(os.Args)
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

	log.Init(log.Config{
		Level:   logLevel,
		Outputs: []string{"stderr"},
	})
	return nil
}

func sendTxs(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	nTxs := 1 // send 1 tx by default
	if cliCtx.NArg() > 0 {
		nTxsArgStr := cliCtx.Args().Get(0)
		nTxsArg, err := strconv.Atoi(nTxsArgStr)
		if err == nil {
			nTxs = nTxsArg
		}
	}

	selNetworks := cliCtx.StringSlice(flagNetworkName)

	// if no network selected, pick the default one
	if selNetworks == nil {
		selNetworks = []string{defaultNetwork}
	}

	for _, selNetwork := range selNetworks {
		network, ok := networks[selNetwork]
		if !ok {
			netKeys := make([]string, 0, len(networks))
			for net := range networks {
				netKeys = append(netKeys, net)
			}
			return fmt.Errorf("please specify one or more of these networks: %v", netKeys)
		}
		log.Infof("connecting to %v: %v", network.Name, network.URL)
		client, err := ethclient.Dial(network.URL)
		if err != nil {
			return err
		}
		log.Infof("connected")

		auth, err := operations.GetAuth(network.PrivateKey, network.ChainID)
		if err != nil {
			return err
		}

		log.Infof("Sender Addr: %v", auth.From.String())

		senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
		if err != nil {
			return err
		}
		log.Debugf("ETH Balance for %v: %v", auth.From, senderBalance)

		amount := big.NewInt(10) //nolint:gomnd
		log.Debugf("Transfer Amount: %v", amount)

		senderNonce, err := client.PendingNonceAt(ctx, auth.From)
		if err != nil {
			return err
		}
		log.Debugf("Sender Nonce: %v", senderNonce)

		to := common.HexToAddress(receiverAddr)
		log.Infof("Receiver Addr: %v", to.String())

		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &to, Value: amount})
		if err != nil {
			return err
		}

		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return err
		}

		txs := make([]*types.Transaction, 0, nTxs)
		for i := 0; i < nTxs; i++ {
			tx := types.NewTransaction(senderNonce+uint64(i), to, amount, gasLimit, gasPrice, nil)
			txs = append(txs, tx)
		}

		if cliCtx.Bool(flagWaitName) {
			var err error
			if network.networkLayer == networkLayer1 {
				err = operations.ApplyL1Txs(ctx, txs, auth, client)
			} else if network.networkLayer == networkLayer2 {
				_, err = operations.ApplyL2Txs(ctx, txs, auth, client, operations.VerifiedConfirmationLevel)
			}
			if err != nil {
				return err
			}
		} else {
			for i := 0; i < nTxs; i++ {
				signedTx, err := auth.Signer(auth.From, txs[i])
				if err != nil {
					return err
				}
				log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
				err = client.SendTransaction(ctx, signedTx)
				if err != nil {
					return err
				}
			}
		}

		if nTxs > 1 {
			log.Infof("%d transactions successfully sent", nTxs)
		} else {
			log.Info("transaction successfully sent")
		}
	}

	return nil
}
