package eth_transfers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/stretchr/testify/require"
)

const (
	profilingEnabled = false
)

func BenchmarkSequencerEthTransfersPoolProcess(b *testing.B) {
	ctx := context.Background()
	//defer func() { require.NoError(b, operations.Teardown()) }()
	opsman, client, pl, senderNonce, gasPrice := setup.Environment(ctx, b)
	shared.Auth.GasPrice = gasPrice
	err := transactions.SendAndWait(shared.Ctx, shared.Auth, senderNonce, client, pl.CountTransactionsByStatus, shared.NumberOfTxs, TxSender)
	require.NoError(b, err)
	setup.BootstrapSequencer(b, opsman)

	var (
		elapsed            time.Duration
		prometheusResponse *http.Response
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", shared.NumberOfTxs), func(b *testing.B) {
		err, _ := transactions.WaitStatusSelected(pl.CountTransactionsByStatus, shared.NumberOfTxs)
		require.NoError(b, err)
		prometheusResponse, err = metrics.FetchPrometheus()

		require.NoError(b, err)
	})

	var profilingResult string
	if profilingEnabled {
		profilingResult, err = metrics.FetchProfiling()
		require.NoError(b, err)
	}

	//err = operations.Teardown()
	if err != nil {
		log.Errorf("failed to teardown: %s", err)
	}

	metrics.CalculateAndPrint(prometheusResponse, profilingResult, elapsed, 0, 0, shared.NumberOfTxs)
	fmt.Printf("%s\n", profilingResult)
}
