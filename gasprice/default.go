package gasprice

import (
	"context"
	"fmt"
)

// DefaultGasPricer gas price from config is set.
type DefaultGasPricer struct {
	cfg  Config
	pool pool
	ctx  context.Context
}

// newDefaultGasPriceSuggester init default gas price suggester.
func newDefaultGasPriceSuggester(ctx context.Context, cfg Config, pool pool) *DefaultGasPricer {
	gpe := &DefaultGasPricer{ctx: ctx, cfg: cfg, pool: pool}
	gpe.setDefaultGasPrice()
	return gpe
}

// UpdateGasPriceAvg not needed for default strategy.
func (d *DefaultGasPricer) UpdateGasPriceAvg() {}

func (d *DefaultGasPricer) setDefaultGasPrice() {
	err := d.pool.SetGasPrice(d.ctx, d.cfg.DefaultGasPriceWei)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}
