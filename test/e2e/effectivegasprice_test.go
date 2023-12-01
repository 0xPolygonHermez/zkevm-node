package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestEffectiveGasPrice(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)

	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(t, err)

	// Send tx
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

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	txs := make([]*types.Transaction, 0, 1)
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	txs = append(txs, tx)

	_, err = operations.ApplyL2Txs(ctx, txs, auth, client, operations.TrustedConfirmationLevel)
	require.NoError(t, err)
}
