package config_test

import (
	"flag"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state/tree"
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
			expectedValue: sequencer.NewDuration(3 * time.Second),
		},
		{
			path:          "Sequencer.PendingTxsQueue.GetPendingTxsFrequency",
			expectedValue: sequencer.NewDuration(5 * time.Second),
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
			path:          "ExecutorServer.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "ExecutorServer.Port",
			expectedValue: 00000,
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
