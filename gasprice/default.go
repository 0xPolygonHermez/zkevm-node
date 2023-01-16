package gasprice

import (
	"context"
	"fmt"
)

// Default gas price from config is set.
type Default struct {
	cfg  Config
	pool pool
}

// newDefaultEstimator init default gas price estimator.
func newDefaultEstimator(cfg Config, pool pool) *Default {
	gpe := &Default{cfg: cfg, pool: pool}
	gpe.setDefaultGasPrice()
	return gpe
}

// UpdateGasPriceAvg not needed for default strategy.
func (d *Default) UpdateGasPriceAvg() {}

func (d *Default) setDefaultGasPrice() {
	ctx := context.Background()
	err := d.pool.SetGasPrice(ctx, d.cfg.DefaultGasPriceWei)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}
