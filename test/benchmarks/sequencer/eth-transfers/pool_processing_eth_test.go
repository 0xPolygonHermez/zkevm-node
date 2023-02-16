package eth_transfers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/setup"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	ethAmount, _ = big.NewInt(0).SetString("100000000000", encoding.Base10)
)

const (
	nTxs             = 10000
	gasLimit         = 21000
	profilingEnabled = true
)

func BenchmarkSequencerEthTransfersPoolProcess(b *testing.B) {
	ctx := context.Background()
	//defer func() { require.NoError(b, operations.Teardown()) }()
	opsman, client, pl, senderNonce, gasPrice := setup.Environment(ctx, b)
	transactions.SendAndWait(b, senderNonce, client, gasPrice, pl, ctx, nTxs, runTxSender)
	setup.BootstrapSequencer(b, opsman)

	var (
		elapsed            time.Duration
		prometheusResponse *http.Response
		err                error
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", nTxs), func(b *testing.B) {
		// Wait all txs to be selected by the sequencer
		start := time.Now()
		log.Debug("Wait for sequencer to select all txs from the pool")
		err := operations.Poll(1*time.Second, shared.DefaultDeadline, func() (bool, error) {
			selectedCount, err := pl.CountTransactionsByStatus(ctx, pool.TxStatusSelected)
			if err != nil {
				return false, err
			}

			log.Debugf("amount of selected txs: %d", selectedCount)
			done := selectedCount >= nTxs
			return done, nil
		})
		require.NoError(b, err)
		elapsed = time.Since(start)
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

	metrics.CalculateAndPrint(prometheusResponse, profilingResult, elapsed, 0, 0, nTxs)
	fmt.Printf("%s\n", profilingResult)
}

func runTxSender(b *testing.B, l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) {
	log.Debugf("sending nonce: %d", nonce)
	tx := types.NewTransaction(nonce, shared.To, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := shared.Auth.Signer(shared.Auth.From, tx)
	require.NoError(b, err)
	err = l2Client.SendTransaction(shared.Ctx, signedTx)
	if errors.Is(err, state.ErrStateNotSynchronized) {
		for errors.Is(err, state.ErrStateNotSynchronized) {
			time.Sleep(5 * time.Second)
			err = l2Client.SendTransaction(shared.Ctx, signedTx)
		}
	}
	require.NoError(b, err)
}
