package erc20_transfers

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	txTimeout        = 60 * time.Second
	profilingEnabled = false
)

var (
	erc20SC *ERC20.ERC20
)

func BenchmarkSequencerERC20TransfersPoolProcess(b *testing.B) {
	start := time.Now()
	opsman, client, pl, auth := setup.Environment(params.Ctx, b)
	setup.BootstrapSequencer(b, opsman)
	startDeploySCTime := time.Now()
	err := deployERC20Contract(b, client, params.Ctx, auth)
	require.NoError(b, err)
	deploySCElapsed := time.Since(startDeploySCTime)
	deploySCSequencerTime, deploySCExecutorOnlyTime, _, err := metrics.GetValues(nil)
	if err != nil {
		return
	}
	initialCount, err := pl.CountTransactionsByStatus(params.Ctx, pool.TxStatusSelected)
	require.NoError(b, err)
	err = transactions.SendAndWait(params.Ctx, auth, client, pl.CountTransactionsByStatus, params.NumberOfTxs, erc20SC, TxSender)
	require.NoError(b, err)

	var (
		elapsed  time.Duration
		response *http.Response
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", params.NumberOfTxs), func(b *testing.B) {
		// Wait all txs to be selected by the sequencer
		err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfTxs)
		require.NoError(b, err)
		elapsed = time.Since(start)
		log.Infof("Total elapsed time: %s", elapsed)
		response, err = metrics.FetchPrometheus()
		require.NoError(b, err)
	})

	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	metrics.CalculateAndPrint(response, profilingResult, elapsed-deploySCElapsed, deploySCSequencerTime, deploySCExecutorOnlyTime, params.NumberOfTxs)
	log.Infof("########################################")
	log.Infof("# Deploying ERC20 SC and Mint Tx took: #")
	log.Infof("########################################")
	metrics.PrintPrometheus(deploySCSequencerTime, deploySCExecutorOnlyTime, 0)
}

func deployERC20Contract(b *testing.B, client *ethclient.Client, ctx context.Context, auth *bind.TransactOpts) error {
	var (
		tx  *types.Transaction
		err error
	)
	log.Debugf("Sending TX to deploy ERC20 SC")
	_, tx, erc20SC, err = ERC20.DeployERC20(auth, client, "Test Coin", "TCO")
	require.NoError(b, err)
	err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	require.NoError(b, err)
	log.Debugf("Sending TX to do a ERC20 mint")
	auth.Nonce = big.NewInt(1) // for the mint tx
	tx, err = erc20SC.Mint(auth, mintAmountBig)
	auth.Nonce = big.NewInt(2)
	require.NoError(b, err)
	err = operations.WaitTxToBeMined(ctx, client, tx, txTimeout)
	require.NoError(b, err)
	return err
}
