package e2e

import (
	"context"
	"math/big"
	"os"
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

func init() {
	os.Setenv("CONFIG_MODE", "test")
}

func TestEthTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() { require.NoError(t, operations.Teardown()) }()

	err := operations.Teardown()
	require.NoError(t, err)
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsman.Setup()
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)
	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(t, err)
	// Send txs
	nTxs := 50
	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(t, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	log.Infof("Receiver Addr: %v", toAddress.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	log.Infof("Sending %d transactions...", nTxs)
	var lastTxHash common.Hash

	var sentTxs []*types.Transaction

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(nonce+uint64(i), toAddress, amount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
		err = client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)
		lastTxHash = signedTx.Hash()

		sentTxs = append(sentTxs, signedTx)
	}
	// wait for TX to be mined
	timeout := 180 * time.Second
	for _, tx := range sentTxs {
		log.Infof("Waiting Tx %s to be mined", tx.Hash())
		err = operations.WaitTxToBeMined(client, tx.Hash(), timeout)
		require.NoError(t, err)
		log.Infof("Tx %s mined successfully", tx.Hash())

		// check transaction nonce against transaction reported L2 block number
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		// get block L2 number
		blockL2Number := receipt.BlockNumber
		require.Equal(t, tx.Nonce(), blockL2Number.Uint64()-1)
	}
	log.Infof("%d transactions added into the trusted state successfully.", nTxs)

	// get block L2 number of the last transaction sent
	receipt, err := client.TransactionReceipt(ctx, lastTxHash)
	require.NoError(t, err)
	l2BlockNumber := receipt.BlockNumber

	// wait for l2 block to be virtualized
	log.Infof("waiting for the block number %v to be virtualized", l2BlockNumber.String())
	err = operations.WaitL2BlockToBeVirtualized(l2BlockNumber, 4*time.Minute)
	require.NoError(t, err)

	// wait for l2 block number to be consolidated
	log.Infof("waiting for the block number %v to be consolidated", l2BlockNumber.String())
	err = operations.WaitL2BlockToBeConsolidated(l2BlockNumber, 4*time.Minute)
	require.NoError(t, err)
}
