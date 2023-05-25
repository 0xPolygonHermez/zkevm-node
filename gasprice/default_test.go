package gasprice

import (
	"context"
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
		DefaultGasPriceWei: 1000000000,
	}
	poolM := new(poolMock)
	poolM.On("SetGasPrice", ctx, cfg.DefaultGasPriceWei).Return(nil).Once()
	dge := newDefaultGasPriceSuggester(ctx, cfg, poolM)
	dge.UpdateGasPriceAvg()
}
