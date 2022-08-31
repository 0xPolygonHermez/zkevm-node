package config_test

import (
	"flag"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
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
			path:          "Database.MaxConns",
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

	ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("", flag.PanicOnError), nil)
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
