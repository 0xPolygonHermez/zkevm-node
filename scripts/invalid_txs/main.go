package main

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ConditionalLoop"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// sending intrinsic invalid transactions
	const numberOfIntrinsicInvalidTransactions = 100
	const transferGas = uint64(22000)
	transferValue := big.NewInt(1)

	url := operations.DefaultL2NetworkURL
	chainID := operations.DefaultL2ChainID
	wrongChainID := operations.DefaultL2ChainID + 1
	privateKey := operations.DefaultSequencerPrivateKey

	ctx := context.Background()

	log.Infof("connecting to %v", url)
	client, err := ethclient.Dial(url)
	chkErr(err)
	log.Infof("connected")

	auth := operations.MustGetAuth(privateKey, chainID)
	wrongAuth := operations.MustGetAuth(privateKey, wrongChainID)

	balance, err := client.BalanceAt(ctx, auth.From, nil)
	chkErr(err)
	log.Debugf("balance of %v: %v", auth.From, balance.String())

	_, scTx, sc, err := ConditionalLoop.DeployConditionalLoop(auth, client)
	log.Debugf("deploying SC tx: %v", scTx.Hash().String())
	chkErr(err)

	err = operations.WaitTxToBeMined(client, scTx.Hash(), operations.DefaultTimeoutTxToBeMined)
	chkErr(err)

	for i := 0; i < numberOfIntrinsicInvalidTransactions; i++ {
		nonce, err := client.PendingNonceAt(ctx, wrongAuth.From)
		chkErr(err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		chkErr(err)

		intrinsicInvalidTx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &wrongAuth.From,
			GasPrice: gasPrice,
			Gas:      transferGas,
			Value:    transferValue,
		})

		intrinsicInvalidSignedTx, err := wrongAuth.Signer(wrongAuth.From, intrinsicInvalidTx)
		chkErr(err)
		log.Debugf("sending intrinsic invalid tx: %v", intrinsicInvalidSignedTx.Hash().String())
		err = client.SendTransaction(ctx, intrinsicInvalidSignedTx)
		chkErr(err)
	}

	// force out of counters error
	times, _ := big.NewInt(0).SetString("1", encoding.Base10)
	callOps := &bind.CallOpts{Pending: false}
	log.Debugf("executing loop %v times", times)
	_, err = sc.ExecuteLoop(callOps, times)
	chkErr(err)
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
