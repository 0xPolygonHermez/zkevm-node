package uniswap_transfers

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	uniswap "github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
	"github.com/stretchr/testify/require"
)

const (
	profilingEnabled = false
)

func BenchmarkSequencerUniswapTransfersPoolProcess(b *testing.B) {
	start := time.Now()
	//defer func() { require.NoError(b, operations.Teardown()) }()

	opsman, client, pl, auth := setup.Environment(params.Ctx, b)
	timeForSetup := time.Since(start)
	setup.BootstrapSequencer(b, opsman)
	deployments := uniswap.DeployContractsAndAddLiquidity(client, auth)
	deploymentTxsCount := uniswap.GetExecutedTransactionsCount()
	elapsedTimeForDeployments := time.Since(start)
	allTxs, err := transactions.SendAndWait(
		auth,
		client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		nil,
		&deployments,
		TxSender,
	)
	require.NoError(b, err)

	elapsed := time.Since(start)
	fmt.Printf("Total elapsed time: %s\n", elapsed)

	startMetrics := time.Now()
	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	metrics.CalculateAndPrint(
		"uniswap",
		deploymentTxsCount+uint64(len(allTxs)),
		client,
		profilingResult,
		elapsed,
		0,
		0,
		allTxs,
	)
	fmt.Printf("%s\n", profilingResult)
	timeForFetchAndPrintMetrics := time.Since(startMetrics)
	metrics.PrintUniswapDeployments(elapsedTimeForDeployments, deploymentTxsCount)
	fmt.Printf("Time for setup: %s\n", timeForSetup)
	fmt.Printf("Time for fetching metrics: %s\n", timeForFetchAndPrintMetrics)
}
