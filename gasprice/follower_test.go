package gasprice

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

func TestUpdateGasPriceFollower(t *testing.T) {
	ctx := context.Background()
	var d time.Duration = 1000000000

	cfg := Config{
		Type:               FollowerType,
		DefaultGasPriceWei: 1000000000,
		UpdatePeriod:       types.NewDuration(d),
		Factor:             0.5,
	}
	l1GasPrice := big.NewInt(10000000000)
	l2GasPrice := uint64(5000000000)
	poolM := new(poolMock)
	ethM := new(ethermanMock)
	ethM.On("GetL1GasPrice", ctx).Return(l1GasPrice).Once()
	poolM.On("SetGasPrices", ctx, l2GasPrice, l1GasPrice.Uint64()).Return(nil).Once()
	f := newFollowerGasPriceSuggester(ctx, cfg, poolM, ethM)

	ethM.On("GetL1GasPrice", ctx).Return(l1GasPrice, l1GasPrice).Once()
	poolM.On("SetGasPrices", ctx, l2GasPrice, l1GasPrice.Uint64()).Return(nil).Once()
	f.UpdateGasPriceAvg()
}

func TestLimitMasGasPrice(t *testing.T) {
	ctx := context.Background()
	var d time.Duration = 1000000000

	cfg := Config{
		Type:               FollowerType,
		DefaultGasPriceWei: 100000000,
		MaxGasPriceWei:     50000000,
		UpdatePeriod:       types.NewDuration(d),
		Factor:             0.5,
	}
	l1GasPrice := big.NewInt(1000000000)
	poolM := new(poolMock)
	ethM := new(ethermanMock)
	ethM.On("GetL1GasPrice", ctx).Return(l1GasPrice)
	// Ensure SetGasPrices is called with the MaxGasPriceWei
	poolM.On("SetGasPrices", ctx, cfg.MaxGasPriceWei, l1GasPrice.Uint64()).Return(nil)
	f := newFollowerGasPriceSuggester(ctx, cfg, poolM, ethM)
	f.UpdateGasPriceAvg()
}
