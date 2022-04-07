package config_test

import (
	"flag"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state/tree"
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
			path:          "GasPriceEstimator.DefaultGasPriceWei",
			expectedValue: uint64(1000000000),
		},
		{
			path:          "MTServer.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "MTServer.Port",
			expectedValue: 50052,
		},
		{
			path:          "MTServer.StoreBackend",
			expectedValue: tree.PgMTStoreBackend,
		},
		{
			path:          "MTClient.URI",
			expectedValue: "127.0.0.1:50052",
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

			if actual != tc.expectedValue {
				t.Fatalf("Unexpected default value for path %q, want %d, got %d", tc.path, tc.expectedValue, actual)
			}
		})
	}
}

func Test_CustomNetwork(t *testing.T) {
	var err error

	app := cli.NewApp()
	var n string
	flag.StringVar(&n, "network", "custom", "")
	var nc string
	flag.StringVar(&nc, "network-cfg", "./network-config.example.json", "")
	ctx := cli.NewContext(app, flag.CommandLine, nil)

	cfg, err := config.Load(ctx)
	require.NoError(t, err)

	assert.Equal(t, uint8(4), cfg.NetworkConfig.Arity)
	assert.Equal(t, uint64(1), cfg.NetworkConfig.GenBlockNumber)
	assert.Equal(t, common.HexToAddress("0xCF7ED3ACCA5A467E9E704C703E8D87F634FB0FC9").Hex(), cfg.NetworkConfig.PoEAddr.Hex())
	assert.Equal(t, common.HexToAddress("0x37AFFAF737C3683AB73F6E1B0933B725AB9796AA").Hex(), cfg.NetworkConfig.MaticAddr.Hex())
	assert.Equal(t, uint64(1337), cfg.NetworkConfig.L1ChainID)
	assert.Equal(t, uint64(1000), cfg.NetworkConfig.L2DefaultChainID)
	assert.Equal(t, uint64(123456), cfg.NetworkConfig.MaxCumulativeGasUsed)

	assert.Equal(t, 3, len(cfg.NetworkConfig.Genesis.Balances))

	assertBalance := func(t *testing.T, a, b string) {
		balance, ok := big.NewInt(0).SetString(b, encoding.Base10)
		assert.True(t, ok)

		addr := common.HexToAddress(a)
		balanceFound, found := cfg.NetworkConfig.Genesis.Balances[addr]
		assert.True(t, found)

		if !found {
			return
		}

		assert.Equal(t, 0, balance.Cmp(balanceFound))
	}

	assertBalance(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "1000000000000000000000")
	assertBalance(t, "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "2000000000000000000000")
	assertBalance(t, "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "3000000000000000000000")
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
