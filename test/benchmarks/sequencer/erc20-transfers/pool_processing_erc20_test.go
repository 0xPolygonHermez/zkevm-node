package erc20_transfers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	txTimeout        = 60 * time.Second
	profilingEnabled = false
)

func BenchmarkSequencerERC20TransfersPoolProcess(b *testing.B) {
	//defer func() { require.NoError(b, operations.Teardown()) }()
	opsman, client, pl, senderNonce, gasPrice := setup.Environment(shared.Ctx, b)

	setup.BootstrapSequencer(b, opsman)
	startDeploySCTime := time.Now()
	err := deployERC20Contract(b, client, shared.Ctx)
	require.NoError(b, err)
	deploySCElapsed := time.Since(startDeploySCTime)
	deploySCSequencerTime, deploySCExecutorOnlyTime, _, err := metrics.GetValues(nil)
	if err != nil {
		return
	}
	shared.Auth.GasPrice = gasPrice
	err = transactions.SendAndWait(shared.Ctx, shared.Auth, senderNonce, client, pl.CountTransactionsByStatus, shared.NumberOfTxs, TxSender)
	require.NoError(b, err)

	var (
		elapsed  time.Duration
		response *http.Response
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", shared.NumberOfTxs), func(b *testing.B) {
		// Wait all txs to be selected by the sequencer
		err, _ := transactions.WaitStatusSelected(pl.CountTransactionsByStatus, shared.NumberOfTxs)
		require.NoError(b, err)
		response, err = metrics.FetchPrometheus()
		require.NoError(b, err)
	})

	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	err = operations.Teardown()
	if err != nil {
		log.Errorf("failed to teardown: %s", err)
	}

	metrics.CalculateAndPrint(response, profilingResult, elapsed-deploySCElapsed, deploySCSequencerTime, deploySCExecutorOnlyTime, shared.NumberOfTxs)
	log.Infof("########################################")
	log.Infof("# Deploying ERC20 SC and Mint Tx took: #")
	log.Infof("########################################")
	metrics.PrintPrometheus(deploySCSequencerTime, deploySCExecutorOnlyTime, 0)
}

func deployERC20Contract(b *testing.B, client *ethclient.Client, ctx context.Context) error {
	var (
		tx  *types.Transaction
		err error
	)
	log.Debugf("Sending TX to deploy ERC20 SC")
	_, tx, erc20SC, err = ERC20.DeployERC20(shared.Auth, client, "Test Coin", "TCO")
	require.NoError(b, err)
	err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	require.NoError(b, err)
	log.Debugf("Sending TX to do a ERC20 mint")
	tx, err = erc20SC.Mint(shared.Auth, mintAmount)
	require.NoError(b, err)
	err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	require.NoError(b, err)
	return err
}
