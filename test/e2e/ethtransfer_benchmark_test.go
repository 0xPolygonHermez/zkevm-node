package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var table = []struct {
	input int
}{
	{input: 100},
	{input: 500},
	{input: 1000},
}

func BenchmarkEthTransfer(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	for _, v := range table {
		ctx := context.Background()

		defer func() { require.NoError(b, operations.Teardown()) }()

		err := operations.Teardown()
		require.NoError(b, err)
		opsCfg := operations.GetDefaultOperationsConfig()
		opsCfg.State.MaxCumulativeGasUsed = 80000000000
		opsman, err := operations.NewManager(ctx, opsCfg)
		require.NoError(b, err)
		err = opsman.Setup()
		require.NoError(b, err)
		time.Sleep(5 * time.Second)
		// Load account with balance on local genesis
		auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
		require.NoError(b, err)
		// Load eth client
		client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
		require.NoError(b, err)
		// Send txs
		amount := big.NewInt(10000)
		toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
		time.Sleep(10 * time.Second)
		senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
		require.NoError(b, err)
		senderNonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(b, err)

		log.Infof("Receiver Addr: %v", toAddress.String())
		log.Infof("Sender Addr: %v", auth.From.String())
		log.Infof("Sender Balance: %v", senderBalance.String())
		log.Infof("Sender Nonce: %v", senderNonce)

		log.Infof("Sending %d transactions...", v.input)
		var lastTxHash common.Hash

		var sentTxs []*types.Transaction

		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
		require.NoError(b, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(b, err)

		nonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(b, err)
		b.Run(fmt.Sprintf("amount of txs: %d", v.input), func(b *testing.B) {
			for i := 0; i < v.input; i++ {
				tx := types.NewTransaction(nonce+uint64(i), toAddress, amount, gasLimit, gasPrice, nil)
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(b, err)
				log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
				err = client.SendTransaction(context.Background(), signedTx)
				require.NoError(b, err)
				lastTxHash = signedTx.Hash()

				sentTxs = append(sentTxs, signedTx)
			}
			// wait for TX to be mined
			timeout := 180 * time.Second
			for _, tx := range sentTxs {
				log.Infof("Waiting Tx %s to be mined", tx.Hash())
				err = operations.WaitTxToBeMined(ctx, client, tx, timeout)
				require.NoError(b, err)
				log.Infof("Tx %s mined successfully", tx.Hash())

				// check transaction nonce against transaction reported L2 block number
				receipt, err := client.TransactionReceipt(ctx, tx.Hash())
				require.NoError(b, err)

				// get block L2 number
				blockL2Number := receipt.BlockNumber
				require.Equal(b, tx.Nonce(), blockL2Number.Uint64()-1)
			}
			log.Infof("%d transactions added into the trusted state successfully.", v.input)

			// get block L2 number of the last transaction sent
			receipt, err := client.TransactionReceipt(ctx, lastTxHash)
			require.NoError(b, err)
			l2BlockNumber := receipt.BlockNumber

			// wait for l2 block to be virtualized
			log.Infof("waiting for the block number %v to be virtualized", l2BlockNumber.String())
			err = operations.WaitL2BlockToBeVirtualized(l2BlockNumber, 4*time.Minute)
			require.NoError(b, err)
		})
	}
}
