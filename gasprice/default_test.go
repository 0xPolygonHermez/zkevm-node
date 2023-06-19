package gasprice

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

func TestUpdateGasPriceDefault(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Type:               DefaultType,
		Factor:             0.5,
		DefaultGasPriceWei: 1000000000,
	}
	factorAsPercentage := int64(cfg.Factor * 100) // nolint:gomnd
	factor := big.NewInt(factorAsPercentage)
	defaultGasPriceDivByFactor := new(big.Int).Div(new(big.Int).SetUint64(cfg.DefaultGasPriceWei), factor)
	l1GasPrice := new(big.Int).Mul(defaultGasPriceDivByFactor, big.NewInt(100)).Uint64() // nolint:gomnd

	poolM := new(poolMock)
	poolM.On("SetGasPrices", ctx, cfg.DefaultGasPriceWei, l1GasPrice).Return(nil).Twice()
	dge := newDefaultGasPriceSuggester(ctx, cfg, poolM)
	dge.UpdateGasPriceAvg()
}
