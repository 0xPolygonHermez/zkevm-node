package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func Test1000EthTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsman.Setup()
	require.NoError(t, err)

	// Load account with balance on local genesis
	auth, err := operations.GetAuth("0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", big.NewInt(1001))
	require.NoError(t, err)
	// Load eth client
	client, err := ethclient.Dial("http://localhost:8123")
	require.NoError(t, err)
	// Send txs
	nTxs := 100
	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000)
	log.Infof("Sending %d transactions...", nTxs)
	var lastTxHash common.Hash

	var sentTxs []*types.Transaction

	for i := 0; i < nTxs; i++ {
		nonce := uint64(i + 1)
		tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)
		if i == nTxs-1 {
			lastTxHash = signedTx.Hash()
		}

		sentTxs = append(sentTxs, signedTx)
	}

	for _, tx := range sentTxs {
		// wait for TX to be mined
		timeout := 20 * time.Minute
		err = operations.WaitTxToBeMined(client, tx.Hash(), timeout)
		require.NoError(t, err)

		// check transaction nonce against transaction reported L2 block number
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		// get block L2 number
		blockL2Number := receipt.BlockNumber

		require.Equal(t, tx.Nonce(), blockL2Number.Uint64())
	}
	log.Infof("\n%d transactions sent without error. Waiting for all the transactions to be mined", nTxs)

	receipt, err := client.TransactionReceipt(ctx, lastTxHash)
	require.NoError(t, err)

	// get block L2 number
	blockL2Number := receipt.BlockNumber

	// wait for l2 Block number to be virtualized

	fmt.Printf("\nL2 Block number: %s", blockL2Number)
	fmt.Printf("\nLast TX Hash %s", lastTxHash.String())
	err = operations.WaitL2BlockToBeVirtualized(blockL2Number, 30*time.Second)
	require.NoError(t, err)

	// wait for l2 block number to be consolidated

	err = operations.WaitL2BlockToBeConsolidated(blockL2Number, 30*time.Second)
	require.NoError(t, err)
}
