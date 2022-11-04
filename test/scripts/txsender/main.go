package main

import (
	"context"
	"math/big"
	"os"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx := context.Background()

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	if err != nil {
		log.Fatal(err)
	}

	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	if err != nil {
		log.Fatal(err)
	}

	// Send 1 tx by default or read the number of txs from args
	nTxs := 1
	if len(os.Args) > 1 {
		nTxs, _ = strconv.Atoi(os.Args[1])
	}

	// Send txs
	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		log.Fatal(err)
	}
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Receiver Addr: %v", toAddress.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		log.Fatal(err)
	}

	txs := make([]*types.Transaction, 0, nTxs)
	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(nonce+uint64(i), toAddress, amount, gasLimit, gasPrice, nil)
		txs = append(txs, tx)
	}

	err = operations.ApplyTxs(ctx, txs)
	if err != nil {
		log.Fatal(err)
	}
}
