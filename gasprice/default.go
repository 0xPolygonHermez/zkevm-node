package gasprice

import (
	"context"
	"fmt"
)

// Default gas price from config is set.
type Default struct {
	cfg  Config
	pool pool
	ctx  context.Context
}

// newDefaultSuggestor init default gas price suggestor.
func newDefaultSuggestor(ctx context.Context, cfg Config, pool pool) *Default {
	gpe := &Default{ctx: ctx, cfg: cfg, pool: pool}
	gpe.setDefaultGasPrice()
	return gpe
}

// UpdateGasPriceAvg not needed for default strategy.
func (d *Default) UpdateGasPriceAvg() {}

func (d *Default) setDefaultGasPrice() {
	err := d.pool.SetGasPrice(d.ctx, d.cfg.DefaultGasPriceWei)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}
