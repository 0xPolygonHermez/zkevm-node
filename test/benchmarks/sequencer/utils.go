package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	Ctx                 = context.Background()
	PoolDbConfig        = dbutils.NewPoolConfigFromEnv()
	SequencerPrivateKey = operations.DefaultSequencerPrivateKey
	ChainID             = operations.DefaultL2ChainID
	OpsCfg              = operations.GetDefaultOperationsConfig()

	ToAddress     = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	To            = common.HexToAddress(ToAddress)
	PrivateKey, _ = crypto.HexToECDSA(strings.TrimPrefix(SequencerPrivateKey, "0x"))
	Auth, _       = bind.NewKeyedTransactorWithChainID(PrivateKey, new(big.Int).SetUint64(ChainID))
)

const (
	DefaultDeadline      = 6000 * time.Second
	MaxCumulativeGasUsed = 80000000000
	PrometheusPort       = 9092
)

func CalculateAndPrintResults(response *http.Response, elapsed time.Duration, sequencerTimeSub, executorTimeSub float64) {
	sequencerTotalProcessingTime, executorTotalProcessingTime, err := GetPrometheusMetricValues(response)
	if err != nil {
		log.Fatalf("error getting prometheus metrics: %v", err)
	}

	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time took for the sequencer To select all txs from the pool: %v", elapsed)
	log.Info("######################")
	log.Info("# Prometheus Metrics #")
	log.Info("######################")
	PrintPrometheusMetrics(sequencerTotalProcessingTime-sequencerTimeSub, executorTotalProcessingTime-executorTimeSub)
}

func PrintPrometheusMetrics(sequencerTotalProcessingTime float64, executorTotalProcessingTime float64) {
	log.Infof("[total_processing_time]: %v s", sequencerTotalProcessingTime)
	log.Infof("[executor_processing_time]: %v s", executorTotalProcessingTime)
	log.Infof("[sequencer_processing_time]: %v s", sequencerTotalProcessingTime-executorTotalProcessingTime)
}

func GetPrometheusMetricValues(metricsResponse *http.Response) (float64, float64, error) {
	var err error
	if metricsResponse == nil {
		metricsResponse, err = http.Get(fmt.Sprintf("http://localhost:%d%s", PrometheusPort, metrics.Endpoint))
		if err != nil {
			log.Errorf("failed to get metrics data: %s", err)
		}
	}

	mf, err := testutils.ParseMetricFamilies(metricsResponse.Body)
	if err != nil {
		return 0, 0, err
	}
	sequencerTotalProcessingTimeHisto := mf["sequencer_processing_time"].Metric[0].Histogram
	sequencerTotalProcessingTime := sequencerTotalProcessingTimeHisto.GetSampleSum()

	executorTotalProcessingTimeHisto := mf["state_executor_processing_time"].Metric[0].Histogram
	executorTotalProcessingTime := executorTotalProcessingTimeHisto.GetSampleSum()
	return sequencerTotalProcessingTime, executorTotalProcessingTime, nil
}

func StartAndSetupSequencer(b *testing.B, opsman *operations.Manager) {
	log.Debug("Starting sequencer ....")
	err := operations.StartComponent("seq")
	require.NoError(b, err)
	log.Debug("Sequencer Started!")
	log.Debug("Setup sequencer ....")
	require.NoError(b, opsman.SetUpSequencer())
	log.Debug("Sequencer Setup ready!")
}

func SendAndWaitTxs(
	b *testing.B,
	senderNonce uint64,
	client *ethclient.Client,
	gasPrice *big.Int,
	pl *pool.Pool,
	ctx context.Context,
	nTxs int,
	txSenderFunc func(b *testing.B, l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64),
) {
	Auth.GasPrice = gasPrice
	Auth.GasLimit = 2100000
	log.Debugf("Sending %d txs ...", nTxs)
	maxNonce := uint64(nTxs) + senderNonce

	for nonce := senderNonce; nonce < maxNonce; nonce++ {
		txSenderFunc(b, client, gasPrice, nonce)
	}
	log.Debug("All txs were sent!")

	log.Debug("Waiting pending transactions To be added in the pool ...")
	err := operations.Poll(1*time.Second, DefaultDeadline, func() (bool, error) {
		// using a closure here To capture st and currentBatchNumber
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

func Setup(ctx context.Context, b *testing.B) (*operations.Manager, *ethclient.Client, *pool.Pool, uint64, *big.Int) {
	if testing.Short() {
		b.Skip()
	}

	err := operations.Teardown()
	require.NoError(b, err)

	OpsCfg.State.MaxCumulativeGasUsed = MaxCumulativeGasUsed
	opsman, err := operations.NewManager(ctx, OpsCfg)
	require.NoError(b, err)

	err = SetupComponents(opsman)
	require.NoError(b, err)
	time.Sleep(5 * time.Second)

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(b, err)

	// Load common client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(b, err)

	st := opsman.State()
	s, err := pgpoolstorage.NewPostgresPoolStorage(PoolDbConfig)
	require.NoError(b, err)
	pl := pool.NewPool(s, st, common.Address{}, ChainID)

	// Print Info before send
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(b, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(b, err)

	// Print Initial Stats
	log.Infof("Receiver Addr: %v", To.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)

	return opsman, client, pl, senderNonce, gasPrice
}

func SetupComponents(opsman *operations.Manager) error {
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
