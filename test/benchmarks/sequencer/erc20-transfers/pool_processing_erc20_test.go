package erc20_transfers

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
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
	nTxs             = 10000
	txTimeout        = 60 * time.Second
	profilingEnabled = true
)

var (
	mintAmount, _  = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	transferAmount = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Div(mintAmount, big.NewInt(nTxs)), big.NewInt(90)), big.NewInt(100))
	erc20SC        *ERC20.ERC20
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

	transactions.SendAndWait(b, senderNonce, client, gasPrice, pl, shared.Ctx, nTxs, runERC20TxSender)
	require.NoError(b, err)

	var (
		elapsed  time.Duration
		response *http.Response
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", nTxs), func(b *testing.B) {
		// Wait all txs to be selected by the sequencer
		start := time.Now()
		log.Debug("Wait for sequencer to select all txs from the pool")
		err := operations.Poll(1*time.Second, shared.DefaultDeadline, func() (bool, error) {
			selectedCount, err := pl.CountTransactionsByStatus(shared.Ctx, pool.TxStatusSelected)
			if err != nil {
				return false, err
			}

			log.Debugf("amount of selected txs: %d", selectedCount)
			done := selectedCount >= nTxs
			return done, nil
		})
		require.NoError(b, err)
		elapsed = time.Since(start)
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

	metrics.CalculateAndPrint(response, profilingResult, elapsed-deploySCElapsed, deploySCSequencerTime, deploySCExecutorOnlyTime, nTxs)
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

func runERC20TxSender(b *testing.B, l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) {
	log.Debugf("sending nonce: %d", nonce)
	var actualTransferAmount *big.Int
	if nonce%2 == 0 {
		actualTransferAmount = big.NewInt(0).Sub(transferAmount, big.NewInt(int64(nonce)))
	} else {
		actualTransferAmount = big.NewInt(0).Add(transferAmount, big.NewInt(int64(nonce)))
	}
	_, err := erc20SC.Transfer(shared.Auth, shared.To, actualTransferAmount)
	require.NoError(b, err)
}
