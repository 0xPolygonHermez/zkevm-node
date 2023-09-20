package e2e

import (
	"context"
	"math/big"
	"os/exec"
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

func TestEthTransferGasless(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// Edit config
	const path = "../../test/config/test.node.config.toml"
	require.NoError(t,
		exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 1000000000/DefaultMinGasPriceAllowed = 0/g", path).Run(),
	)
	require.NoError(t,
		exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = true/EnableL2SuggestedGasPricePolling = false/g", path).Run(),
	)
	// Undo edit config
	defer func() {
		require.NoError(t,
			exec.Command("sed", "-i", "s/DefaultMinGasPriceAllowed = 0/DefaultMinGasPriceAllowed = 1000000000/g", path).Run(),
		)
		require.NoError(t,
			exec.Command("sed", "-i", "s/EnableL2SuggestedGasPricePolling = false/EnableL2SuggestedGasPricePolling = true/g", path).Run(),
		)
	}()

	ctx := context.Background()
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
	nTxs := 10
	amount := big.NewInt(0)
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

	// Force gas price to be 0
	gasPrice := big.NewInt(0)
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	txs := make([]*types.Transaction, 0, nTxs)
	for i := 0; i < nTxs; i++ {
		tx := types.NewTransaction(nonce+uint64(i), toAddress, amount, gasLimit, gasPrice, nil)
		txs = append(txs, tx)
	}

	_, err = operations.ApplyL2Txs(ctx, txs, auth, client, operations.VerifiedConfirmationLevel)
	require.NoError(t, err)
}
