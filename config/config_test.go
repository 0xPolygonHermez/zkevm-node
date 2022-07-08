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
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/state/tree"
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
			path:          "Sequencer.AllowNonRegistered",
			expectedValue: false,
		},
		{
			path:          "Sequencer.InitBatchProcessorIfDiffType",
			expectedValue: sequencer.InitBatchProcessorIfDiffTypeSynced,
		},
		{
			path:          "Sequencer.PriceGetter.Type",
			expectedValue: pricegetter.DefaultType,
		},
		{
			path:          "Sequencer.PriceGetter.DefaultPrice",
			expectedValue: pricegetter.TokenPrice{Float: new(big.Float).SetInt64(2000)},
		},
		{
			path:          "Sequencer.MaxSendBatchTxRetries",
			expectedValue: uint32(5),
		},
		{
			path:          "Sequencer.FrequencyForResendingFailedSendBatchesInMilliseconds",
			expectedValue: int64(1000),
		},
		{
			path:          "Sequencer.PendingTxsQueue.TxPendingInQueueCheckingFrequency",
			expectedValue: types.NewDuration(3 * time.Second),
		},
		{
			path:          "Sequencer.PendingTxsQueue.GetPendingTxsFrequency",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "SequencerV2.WaitPeriodPoolIsEmpty",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "SequencerV2.LastBatchVirtualizationTimeMaxWaitPeriod",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "SequencerV2.WaitBlocksToUpdateGER",
			expectedValue: uint64(10),
		},
		{
			path:          "SequencerV2.LastTimeBatchMaxWaitPeriod",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "SequencerV2.BlocksAmountForTxsToBeDeleted",
			expectedValue: uint64(100),
		},
		{
			path:          "SequencerV2.FrequencyToCheckTxsForDelete",
			expectedValue: types.NewDuration(12 * time.Hour),
		},
		{
			path:          "SequencerV2.ProfitabilityChecker.SendBatchesEvenWhenNotProfitable",
			expectedValue: true,
		},
		{
			path:          "SequencerV2.MaxGasUsed",
			expectedValue: int64(100000),
		},
		{
			path:          "SequencerV2.MaxKeccakHashes",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxPoseidonHashes",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxPoseidonPaddings",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxMemAligns",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxArithmetics",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxBinaries",
			expectedValue: int32(100),
		},
		{
			path:          "SequencerV2.MaxSteps",
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
			path:          "MTServer.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "MTServer.Port",
			expectedValue: 50060,
		},
		{
			path:          "MTServer.StoreBackend",
			expectedValue: tree.PgMTStoreBackend,
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
			path:          "RPC.SequencerAddress",
			expectedValue: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
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
