package gasprice

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// GasPricerI interface for gas price estimator.
type GasPricerI interface {
	UpdateGasPriceAvg()
}

// GasPricer interface for gas price estimator.
type GasPricer struct {
	cfg Config
	ctx context.Context
	gpI GasPricerI
}

// NewGasPricer init.
func NewGasPricer(ctx context.Context, cfg Config, pool pool, ethMan *etherman.Client, state *state.State) *GasPricer {
	var gpricer GasPricerI
	switch cfg.Type {
	case LastNBatchesType:
		gpricer = newEstimatorLastNL2Blocks(ctx, cfg, state, pool)
	case FollowerType:
		gpricer = newFollowerGasEstimator(ctx, cfg, pool, ethMan)
	case DefaultType:
		gpricer = newDefaultEstimator(ctx, cfg, pool)
	}
	gpe := &GasPricer{
		cfg: cfg,
		ctx: ctx,
		gpI: gpricer,
	}
	return gpe
}

// Start function runs the gasPricer
func (g GasPricer) Start() error {
	for {
		select {
		case <-g.ctx.Done():
			return nil
		case <-time.After(g.cfg.UpdatePeriod.Duration): // TODO meter esto en config files
			g.gpI.UpdateGasPriceAvg()
		}
	}
}
