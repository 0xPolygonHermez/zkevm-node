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
	poolM := new(poolMock)
	ethM := new(ethermanMock)
	ethM.On("GetL1GasPrice", ctx).Return(big.NewInt(10000000000)).Once()
	poolM.On("SetGasPrice", ctx, uint64(5000000000)).Return(nil).Once()
	f := newFollowerGasPriceSuggester(ctx, cfg, poolM, ethM)

	ethM.On("GetL1GasPrice", ctx).Return(big.NewInt(10000000000)).Once()
	poolM.On("SetGasPrice", ctx, uint64(5000000000)).Return(nil).Once()
	f.UpdateGasPriceAvg()
}
