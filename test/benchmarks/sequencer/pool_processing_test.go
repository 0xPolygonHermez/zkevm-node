package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	nTxs                 = 300
	gasLimit             = 21000
	prometheusPort       = 9092
	defaultDeadline      = 6000 * time.Second
	maxCumulativeGasUsed = 80000000000
)

var (
	ctx                 = context.Background()
	poolDbConfig        = dbutils.NewPoolConfigFromEnv()
	sequencerPrivateKey = operations.DefaultSequencerPrivateKey
	chainID             = operations.DefaultL2ChainID
	opsCfg              = operations.GetDefaultOperationsConfig()

	toAddress     = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	to            = common.HexToAddress(toAddress)
	ethAmount, _  = big.NewInt(0).SetString("100000000000", encoding.Base10)
	privateKey, _ = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	auth, _       = bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
)

func BenchmarkSequencerPoolProcess(b *testing.B) {
	ctx := context.Background()
	defer func() { require.NoError(b, operations.Teardown()) }()
	opsman, client, pl, senderNonce, gasPrice := setup(ctx, b)
	sendAndWaitTxs(b, senderNonce, client, gasPrice, pl, ctx)
	startAndSetupSequencer(b, opsman)

	var (
		elapsed  time.Duration
		response *http.Response
		err      error
	)

	b.Run(fmt.Sprintf("sequencer_selecting_%d_txs", nTxs), func(b *testing.B) {
		// Wait all txs to be selected by the sequencer
		start := time.Now()
		log.Debug("Wait for sequencer to select all txs from the pool")
		err := operations.Poll(1*time.Second, defaultDeadline, func() (bool, error) {
			selectedCount, err := pl.CountTransactionsByStatus(ctx, pool.TxStatusSelected)
			if err != nil {
				return false, err
			}

			log.Debugf("amount of selected txs: %d", selectedCount)
			done := selectedCount == nTxs
			return done, nil
		})
		require.NoError(b, err)
		elapsed = time.Since(start)
		response, err = http.Get(fmt.Sprintf("http://localhost:%d%s", prometheusPort, metrics.Endpoint))
		if err != nil {
			log.Errorf("failed to get metrics data: %s", err)
		}
	})

	err = operations.Teardown()
	if err != nil {
		log.Errorf("failed to teardown: %s", err)
	}

	printResults(response, elapsed)
}

func printResults(metricsResponse *http.Response, elapsed time.Duration) {
	mf, err := testutils.ParseMetricFamilies(metricsResponse.Body)
	if err != nil {
		return
	}
	sequencerTotalProcessingTimeHisto := mf["sequencer_processing_time"].Metric[0].Histogram
	sequencerTotalProcessingTime := sequencerTotalProcessingTimeHisto.GetSampleSum()

	executorTotalProcessingTimeHisto := mf["state_executor_processing_time"].Metric[0].Histogram
	executorTotalProcessingTime := executorTotalProcessingTimeHisto.GetSampleSum()

	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time took for the sequencer to select all txs from the pool: %v", elapsed)
	log.Info("######################")
	log.Info("# Prometheus Metrics #")
	log.Info("######################")
	log.Infof("[sequencer_processing_time]: %v s", sequencerTotalProcessingTime)
	log.Infof("[state_executor_processing_time (sequencer)]: %v s", executorTotalProcessingTime)
	log.Infof("[sequencer_processing_time_without_executor]: %v s", sequencerTotalProcessingTime-executorTotalProcessingTime)
}

func startAndSetupSequencer(b *testing.B, opsman *operations.Manager) {
	log.Debug("Starting sequencer ....")
	err := operations.StartComponent("seq")
	require.NoError(b, err)
	log.Debug("Sequencer Started!")
	log.Debug("Setup sequencer ....")
	require.NoError(b, opsman.SetUpSequencer())
	log.Debug("Sequencer setup ready!")
}

func sendAndWaitTxs(b *testing.B, senderNonce uint64, client *ethclient.Client, gasPrice *big.Int, pl *pool.Pool, ctx context.Context) {
	log.Debugf("Sending %d txs ...", nTxs)
	maxNonce := uint64(nTxs) + senderNonce

	for nonce := senderNonce; nonce < maxNonce; nonce++ {
		runTxSender(b, client, gasPrice, nonce)
	}
	log.Debug("All txs were sent!")

	log.Debug("Waiting pending transactions to be added in the pool ...")
	err := operations.Poll(1*time.Second, defaultDeadline, func() (bool, error) {
		// using a closure here to capture st and currentBatchNumber
		count, err := pl.CountPendingTransactions(ctx)
		if err != nil {
			return false, err
		}

		log.Debugf("amount of pending txs: %d\n", count)
		done := count == uint64(nTxs)
		return done, nil
	})
	require.NoError(b, err)
	log.Debug("All pending txs are added in the pool!")
}

func setup(ctx context.Context, b *testing.B) (*operations.Manager, *ethclient.Client, *pool.Pool, uint64, *big.Int) {
	if testing.Short() {
		b.Skip()
	}

	err := operations.Teardown()
	require.NoError(b, err)

	opsCfg.State.MaxCumulativeGasUsed = maxCumulativeGasUsed
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(b, err)

	err = setupComponents(opsman)
	require.NoError(b, err)
	time.Sleep(5 * time.Second)

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(b, err)

	// Load eth client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(b, err)

	st := opsman.State()
	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDbConfig)
	require.NoError(b, err)
	pl := pool.NewPool(s, st, common.Address{}, chainID)

	// Print Info before send
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(b, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(b, err)

	// Print Initial Stats
	log.Infof("Receiver Addr: %v", to.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)

	return opsman, client, pl, senderNonce, gasPrice
}

func setupComponents(opsman *operations.Manager) error {
	// Run network container
	err := opsman.StartNetwork()
	if err != nil {
		return err
	}

	// Approve matic
	err = operations.ApproveMatic()
	if err != nil {
		return err
	}

	err = operations.StartComponent("sync")
	if err != nil {
		return err
	}

	err = operations.StartComponent("json-rpc")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	return nil
}

func runTxSender(b *testing.B, l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) {
	log.Debugf("sending nonce: %d", nonce)
	tx := types.NewTransaction(nonce, to, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(b, err)
	err = l2Client.SendTransaction(ctx, signedTx)
	if errors.Is(err, state.ErrStateNotSynchronized) {
		for errors.Is(err, state.ErrStateNotSynchronized) {
			time.Sleep(5 * time.Second)
			err = l2Client.SendTransaction(ctx, signedTx)
		}
	}
	require.NoError(b, err)
}
