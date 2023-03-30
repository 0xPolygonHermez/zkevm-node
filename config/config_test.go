package config_test

import (
	"flag"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func Test_Defaults(t *testing.T) {
	tcs := []struct {
		path          string
		expectedValue interface{}
	}{
		{
			path:          "Log.Environment",
			expectedValue: log.LogEnvironment("development"),
		},
		{
			path:          "Log.Level",
			expectedValue: "debug",
		},
		{
			path:          "Log.Outputs",
			expectedValue: []string{"stderr"},
		},
		{
			path:          "Synchronizer.SyncChunkSize",
			expectedValue: uint64(100),
		},
		{
			path:          "PriceGetter.Type",
			expectedValue: pricegetter.DefaultType,
		},
		{
			path:          "PriceGetter.DefaultPrice",
			expectedValue: pricegetter.TokenPrice{Float: new(big.Float).SetInt64(2000)},
		},
		{
			path:          "Sequencer.WaitPeriodPoolIsEmpty",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "Sequencer.LastBatchVirtualizationTimeMaxWaitPeriod",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Sequencer.MaxTxsPerBatch",
			expectedValue: uint64(150),
		},
		{
			path:          "Sequencer.MaxBatchBytesSize",
			expectedValue: uint64(129848),
		},
		{
			path:          "Sequencer.BlocksAmountForTxsToBeDeleted",
			expectedValue: uint64(100),
		},
		{
			path:          "Sequencer.FrequencyToCheckTxsForDelete",
			expectedValue: types.NewDuration(12 * time.Hour),
		},
		{
			path:          "Sequencer.MaxCumulativeGasUsed",
			expectedValue: uint64(30000000),
		},
		{
			path:          "Sequencer.MaxKeccakHashes",
			expectedValue: uint32(468),
		},
		{
			path:          "Sequencer.MaxPoseidonHashes",
			expectedValue: uint32(279620),
		},
		{
			path:          "Sequencer.MaxPoseidonPaddings",
			expectedValue: uint32(149796),
		},
		{
			path:          "Sequencer.MaxMemAligns",
			expectedValue: uint32(262144),
		},
		{
			path:          "Sequencer.MaxArithmetics",
			expectedValue: uint32(262144),
		},
		{
			path:          "Sequencer.MaxBinaries",
			expectedValue: uint32(262144),
		},
		{
			path:          "Sequencer.MaxSteps",
			expectedValue: uint32(8388608),
		},
		{
			path:          "Sequencer.TxLifetimeCheckTimeout",
			expectedValue: types.NewDuration(10 * time.Minute),
		},
		{
			path:          "Sequencer.MaxTxLifetime",
			expectedValue: types.NewDuration(3 * time.Hour),
		},
		{
			path:          "Sequencer.MaxTxSizeForL1",
			expectedValue: uint64(131072),
		},
		{
			path:          "Sequencer.Finalizer.GERDeadlineTimeoutInSec",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.ForcedBatchDeadlineTimeoutInSec",
			expectedValue: types.NewDuration(60 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.SendingToL1DeadlineTimeoutInSec",
			expectedValue: types.NewDuration(20 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.SleepDurationInMs",
			expectedValue: types.NewDuration(100 * time.Millisecond),
		},
		{
			path:          "Sequencer.Finalizer.ResourcePercentageToCloseBatch",
			expectedValue: uint32(10),
		},
		{
			path:          "Sequencer.Finalizer.GERFinalityNumberOfBlocks",
			expectedValue: uint64(64),
		},
		{
			path:          "Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingL1Timeout",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingGER",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.ClosingSignalsManagerWaitForCheckingForcedBatches",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.ForcedBatchesFinalityNumberOfBlocks",
			expectedValue: uint64(64),
		},
		{
			path:          "Sequencer.DBManager.PoolRetrievalInterval",
			expectedValue: types.NewDuration(500 * time.Millisecond),
		},
		{
			path:          "Etherman.URL",
			expectedValue: "http://localhost:8545",
		},
		{
			path:          "Etherman.L1ChainID",
			expectedValue: uint64(1337),
		},
		{
			path:          "Etherman.PoEAddr",
			expectedValue: common.HexToAddress("0x610178dA211FEF7D417bC0e6FeD39F05609AD788"),
		},
		{
			path:          "Etherman.MaticAddr",
			expectedValue: common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		},
		{
			path:          "Etherman.GlobalExitRootManagerAddr",
			expectedValue: common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"),
		},
		{
			path:          "Etherman.MultiGasProvider",
			expectedValue: true,
		},
		{
			path:          "EthTxManager.FrequencyToMonitorTxs",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "EthTxManager.WaitTxToBeMined",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "EthTxManager.WaitTxToBeMined",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "EthTxManager.ForcedGas",
			expectedValue: uint64(0),
		},
		{
			path:          "PriceGetter.Type",
			expectedValue: pricegetter.DefaultType,
		},
		{
			path:          "PriceGetter.DefaultPrice",
			expectedValue: pricegetter.TokenPrice{Float: new(big.Float).SetInt64(2000)},
		},
		{
			path:          "L2GasPriceSuggester.DefaultGasPriceWei",
			expectedValue: uint64(1000000000),
		},
		{
			path:          "MTClient.URI",
			expectedValue: "127.0.0.1:50061",
		},
		{
			path:          "StateDB.User",
			expectedValue: "state_user",
		},
		{
			path:          "StateDB.Password",
			expectedValue: "state_password",
		},
		{
			path:          "StateDB.Name",
			expectedValue: "state_db",
		},
		{
			path:          "StateDB.Host",
			expectedValue: "localhost",
		},
		{
			path:          "StateDB.Port",
			expectedValue: "5432",
		},
		{
			path:          "StateDB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "StateDB.MaxConns",
			expectedValue: 200,
		},
		{
			path:          "Pool.FreeClaimGasLimit",
			expectedValue: uint64(150000),
		},
		{
			path:          "Pool.MaxTxBytesSize",
			expectedValue: uint64(30132),
		},
		{
			path:          "Pool.MaxTxDataBytesSize",
			expectedValue: 30000,
		},

		{
			path:          "Pool.DefaultMinGasPriceAllowed",
			expectedValue: uint64(1000000000),
		},
		{
			path:          "Pool.MinAllowedGasPriceInterval",
			expectedValue: types.NewDuration(5 * time.Minute),
		},
		{
			path:          "Pool.PollMinAllowedGasPriceInterval",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "Pool.DB.User",
			expectedValue: "pool_user",
		},
		{
			path:          "Pool.DB.Password",
			expectedValue: "pool_password",
		},
		{
			path:          "Pool.DB.Name",
			expectedValue: "pool_db",
		},
		{
			path:          "Pool.DB.Host",
			expectedValue: "localhost",
		},
		{
			path:          "Pool.DB.Port",
			expectedValue: "5432",
		},
		{
			path:          "Pool.DB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "Pool.DB.MaxConns",
			expectedValue: 200,
		},
		{
			path:          "RPC.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "RPC.Port",
			expectedValue: int(8123),
		},
		{
			path:          "RPC.ReadTimeoutInSec",
			expectedValue: time.Duration(60),
		},
		{
			path:          "RPC.WriteTimeoutInSec",
			expectedValue: time.Duration(60),
		},
		{
			path:          "RPC.SequencerNodeURI",
			expectedValue: "",
		},
		{
			path:          "RPC.MaxRequestsPerIPAndSecond",
			expectedValue: float64(50),
		},
		{
			path:          "RPC.BroadcastURI",
			expectedValue: "127.0.0.1:61090",
		},
		{
			path:          "RPC.DefaultSenderAddress",
			expectedValue: "0x1111111111111111111111111111111111111111",
		},
		{
			path:          "RPC.WebSockets.Enabled",
			expectedValue: false,
		},

		{
			path:          "RPC.EnableL2SuggestedGasPricePolling",
			expectedValue: true,
		},
		{
			path:          "RPC.WebSockets.Port",
			expectedValue: 8133,
		},
		{
			path:          "Executor.URI",
			expectedValue: "127.0.0.1:50071",
		},
		{
			path:          "BroadcastServer.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "BroadcastServer.Port",
			expectedValue: 61090,
		},
		{
			path:          "Metrics.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "Metrics.Port",
			expectedValue: 9091,
		},
		{
			path:          "Metrics.Enabled",
			expectedValue: false,
		},
		{
			path:          "Aggregator.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "Aggregator.Port",
			expectedValue: 50081,
		},
		{
			path:          "Aggregator.RetryTime",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Aggregator.VerifyProofInterval",
			expectedValue: types.NewDuration(90 * time.Second),
		},
		{
			path:          "Aggregator.TxProfitabilityCheckerType",
			expectedValue: aggregator.TxProfitabilityCheckerType(aggregator.ProfitabilityAcceptAll),
		},
		{
			path:          "Aggregator.TxProfitabilityMinReward",
			expectedValue: aggregator.TokenAmountWithDecimals{Int: big.NewInt(1100000000000000000)},
		},
		{
			path:          "Aggregator.ProofStatePollingInterval",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Aggregator.CleanupLockedProofsInterval",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "Aggregator.GeneratingProofCleanupThreshold",
			expectedValue: "10m",
		},
	}
	file, err := os.CreateTemp("", "genesisConfig")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(file.Name()))
	}()
	require.NoError(t, os.WriteFile(file.Name(), []byte("{}"), 0600))

	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	flagSet.String(config.FlagGenesisFile, file.Name(), "")
	ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
	cfg, err := config.Load(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading default config: %v", err)
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			actual := getValueFromStruct(tc.path, cfg)
			require.Equal(t, tc.expectedValue, actual)
		})
	}
}

func getValueFromStruct(path string, object interface{}) interface{} {
	keySlice := strings.Split(path, ".")
	v := reflect.ValueOf(object)

	for _, key := range keySlice {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		v = v.FieldByName(key)
	}
	return v.Interface()
}

func TestEnvVarArrayDecoding(t *testing.T) {
	file, err := os.CreateTemp("", "genesisConfig")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(file.Name()))
	}()
	require.NoError(t, os.WriteFile(file.Name(), []byte("{}"), 0600))
	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	flagSet.String(config.FlagGenesisFile, file.Name(), "")
	ctx := cli.NewContext(cli.NewApp(), flagSet, nil)

	os.Setenv("ZKEVM_NODE_LOG_OUTPUTS", "a,b,c")
	defer func() {
		os.Unsetenv("ZKEVM_NODE_LOG_OUTPUTS")
	}()

	cfg, err := config.Load(ctx)
	require.NoError(t, err)

	assert.Equal(t, 3, len(cfg.Log.Outputs))
	assert.Equal(t, "a", cfg.Log.Outputs[0])
	assert.Equal(t, "b", cfg.Log.Outputs[1])
	assert.Equal(t, "c", cfg.Log.Outputs[2])
}
