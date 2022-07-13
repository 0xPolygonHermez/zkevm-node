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
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "Sequencer.LastBatchVirtualizationTimeMaxWaitPeriod",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "Sequencer.WaitBlocksToUpdateGER",
			expectedValue: uint64(10),
		},
		{
			path:          "Sequencer.LastTimeBatchMaxWaitPeriod",
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
			path:          "Sequencer.MaxGasUsed",
			expectedValue: int64(100000),
		},
		{
			path:          "Sequencer.MaxKeccakHashes",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxPoseidonHashes",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxPoseidonPaddings",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxMemAligns",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxArithmetics",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxBinaries",
			expectedValue: int32(100),
		},
		{
			path:          "Sequencer.MaxSteps",
			expectedValue: int32(100),
		},
		{
			path:          "EthTxManager.MaxSendBatchTxRetries",
			expectedValue: uint32(10),
		},
		{
			path:          "EthTxManager.FrequencyForResendingFailedSendBatchesInMilliseconds",
			expectedValue: int64(1000),
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
			expectedValue: "127.0.0.1:50060",
		},
		{
			path:          "Database.MaxConns",
			expectedValue: 200,
		},
		{
			path:          "RPC.MaxRequestsPerIPAndSecond",
			expectedValue: float64(50),
		},
		{
			path:          "Executor.URI",
			expectedValue: "51.210.116.237:50071",
		},
		{
			path:          "RPC.ChainID",
			expectedValue: uint64(1001),
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
			path:          "BroadcastClient.URI",
			expectedValue: "127.0.0.1:61090",
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
