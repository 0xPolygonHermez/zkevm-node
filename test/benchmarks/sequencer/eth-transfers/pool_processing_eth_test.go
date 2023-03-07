package eth_transfers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
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
	setup.BootstrapSequencer(b, opsman)
	err = transactions.SendAndWait(params.Ctx, auth, client, pl.CountTransactionsByStatus, params.NumberOfTxs, nil, TxSender)
	require.NoError(b, err)

	var (
		elapsed            time.Duration
		prometheusResponse *http.Response
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", params.NumberOfTxs), func(b *testing.B) {
		err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfTxs)
		require.NoError(b, err)
		elapsed = time.Since(start)
		log.Infof("Total elapsed time: %s", elapsed)
		prometheusResponse, err = metrics.FetchPrometheus()
		require.NoError(b, err)
	})

	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	metrics.CalculateAndPrint(prometheusResponse, profilingResult, elapsed, 0, 0, params.NumberOfTxs)
	fmt.Printf("%s\n", profilingResult)
}
