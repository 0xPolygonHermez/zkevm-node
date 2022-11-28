package config_test

import (
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func Test_Defaults(t *testing.T) {
	tcs := []struct {
		path          string
		expectedValue interface{}
	}{
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
			expectedValue: types.NewDuration(300 * time.Second),
		},
		{
			path:          "Sequencer.WaitBlocksToUpdateGER",
			expectedValue: uint64(10),
		},
		{
			path:          "Sequencer.MaxTimeForBatchToBeOpen",
			expectedValue: types.NewDuration(15 * time.Second),
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
			path:          "Sequencer.ProfitabilityChecker.SendBatchesEvenWhenNotProfitable",
			expectedValue: true,
		},
		{
			path:          "Sequencer.MaxCumulativeGasUsed",
			expectedValue: uint64(30000000),
		},
		{
			path:          "Sequencer.MaxKeccakHashes",
			expectedValue: int32(468),
		},
		{
			path:          "Sequencer.MaxPoseidonHashes",
			expectedValue: int32(279620),
		},
		{
			path:          "Sequencer.MaxPoseidonPaddings",
			expectedValue: int32(149796),
		},
		{
			path:          "Sequencer.MaxMemAligns",
			expectedValue: int32(262144),
		},
		{
			path:          "Sequencer.MaxArithmetics",
			expectedValue: int32(262144),
		},
		{
			path:          "Sequencer.MaxBinaries",
			expectedValue: int32(262144),
		},
		{
			path:          "Sequencer.MaxSteps",
			expectedValue: int32(8388608),
		},
		{
			path:          "Sequencer.MaxSequenceSize",
			expectedValue: sequencer.MaxSequenceSize{Int: new(big.Int).SetInt64(2000000)},
		},
		{
			path:          "Sequencer.MaxAllowedFailedCounter",
			expectedValue: uint64(50),
		},
		{
			path:          "EthTxManager.MaxSendBatchTxRetries",
			expectedValue: uint32(10),
		},
		{
			path:          "EthTxManager.MaxVerifyBatchTxRetries",
			expectedValue: uint32(10),
		},
		{
			path:          "EthTxManager.FrequencyForResendingFailedSendBatches",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "EthTxManager.FrequencyForResendingFailedVerifyBatch",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "EthTxManager.WaitTxToBeMined",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "EthTxManager.WaitTxToBeSynced",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "EthTxManager.PercentageToIncreaseGasPrice",
			expectedValue: uint64(10),
		},
		{
			path:          "EthTxManager.PercentageToIncreaseGasLimit",
			expectedValue: uint64(10),
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
			path:          "GasPriceEstimator.DefaultGasPriceWei",
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
			path:          "PoolDB.User",
			expectedValue: "pool_user",
		},
		{
			path:          "PoolDB.Password",
			expectedValue: "pool_password",
		},
		{
			path:          "PoolDB.Name",
			expectedValue: "pool_db",
		},
		{
			path:          "PoolDB.Host",
			expectedValue: "localhost",
		},
		{
			path:          "PoolDB.Port",
			expectedValue: "5432",
		},
		{
			path:          "PoolDB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "PoolDB.MaxConns",
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
			path:          "RPC.DB.User",
			expectedValue: "rpc_user",
		},
		{
			path:          "RPC.DB.Password",
			expectedValue: "rpc_password",
		},
		{
			path:          "RPC.DB.Name",
			expectedValue: "rpc_db",
		},
		{
			path:          "RPC.DB.Host",
			expectedValue: "localhost",
		},
		{
			path:          "RPC.DB.Port",
			expectedValue: "5432",
		},
		{
			path:          "RPC.DB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "RPC.DB.MaxConns",
			expectedValue: 200,
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
	}
	file, err := ioutil.TempFile("", "genesisConfig")
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
