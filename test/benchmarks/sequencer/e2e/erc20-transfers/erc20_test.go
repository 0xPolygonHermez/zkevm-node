package erc20_transfers

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/stretchr/testify/require"
)

const (
	profilingEnabled = false
)

var (
	erc20SC *ERC20.ERC20
)

func BenchmarkSequencerERC20TransfersPoolProcess(b *testing.B) {
	var err error
	start := time.Now()
	opsman, client, pl, auth := setup.Environment(params.Ctx, b)
	setup.BootstrapSequencer(b, opsman)
	timeForSetup := time.Since(start)
	startDeploySCTime := time.Now()
	erc20SC, err = DeployERC20Contract(client, params.Ctx, auth)
	require.NoError(b, err)
	deploySCElapsed := time.Since(startDeploySCTime)
	deployMetricsValues, err := metrics.GetValues(nil)
	if err != nil {
		return
	}
	allTxs, err := transactions.SendAndWait(
		auth,
		client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		erc20SC,
		nil,
		TxSender,
	)
	require.NoError(b, err)

	var (
		elapsed time.Duration
	)

	elapsed = time.Since(start)
	fmt.Printf("Total elapsed time: %s\n", elapsed)

	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	startMetrics := time.Now()
	metrics.CalculateAndPrint(
		"erc20",
		uint64(len(allTxs)),
		client,
		profilingResult,
		elapsed,
		deployMetricsValues.SequencerTotalProcessingTime,
		deployMetricsValues.ExecutorTotalProcessingTime,
		allTxs,
	)
	timeForFetchAndPrintMetrics := time.Since(startMetrics)
	fmt.Println("########################################")
	fmt.Println("# Deploying ERC20 SC and Mint Tx took: #")
	fmt.Println("########################################")
	fmt.Printf("%s\n", deploySCElapsed)
	fmt.Printf("Time for setup: %s\n", timeForSetup)
	fmt.Printf("Time for fetching metrics: %s\n", timeForFetchAndPrintMetrics)
}
