package eth_transfers

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/stretchr/testify/require"
)

const (
	profilingEnabled = false
)

func BenchmarkSequencerEthTransfersPoolProcess(b *testing.B) {
	start := time.Now()
	//defer func() { require.NoError(b, operations.Teardown()) }()
	opsman, client, pl, auth := setup.Environment(params.Ctx, b)
	initialCount, err := pl.CountTransactionsByStatus(params.Ctx, pool.TxStatusSelected)
	require.NoError(b, err)
	timeForSetup := time.Since(start)
	setup.BootstrapSequencer(b, opsman)
	allTxs, err := transactions.SendAndWait(
		auth,
		client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		nil,
		nil,
		TxSender,
	)
	require.NoError(b, err)

	var (
		elapsed time.Duration
	)
	err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfOperations)
	require.NoError(b, err)
	elapsed = time.Since(start)
	fmt.Printf("Total elapsed time: %s\n", elapsed)

	startMetrics := time.Now()
	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	metrics.CalculateAndPrint(
		"eth",
		uint64(len(allTxs)),
		client,
		profilingResult,
		elapsed,
		0,
		0,
		allTxs,
	)
	fmt.Printf("%s\n", profilingResult)
	timeForFetchAndPrintMetrics := time.Since(startMetrics)
	fmt.Printf("Time for setup: %s\n", timeForSetup)
	fmt.Printf("Time for fetching metrics: %s\n", timeForFetchAndPrintMetrics)
}
