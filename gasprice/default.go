package gasprice

import (
	"context"
	"fmt"
	"math/big"
)

// DefaultGasPricer gas price from config is set.
type DefaultGasPricer struct {
	cfg        Config
	pool       poolInterface
	ctx        context.Context
	l1GasPrice uint64
}

// newDefaultGasPriceSuggester init default gas price suggester.
func newDefaultGasPriceSuggester(ctx context.Context, cfg Config, pool poolInterface) *DefaultGasPricer {
	// Apply factor to calculate l1 gasPrice
	factorAsPercentage := int64(cfg.Factor * 100) // nolint:gomnd
	factor := big.NewInt(factorAsPercentage)
	defaultGasPriceDivByFactor := new(big.Int).Div(new(big.Int).SetUint64(cfg.DefaultGasPriceWei), factor)

	gpe := &DefaultGasPricer{
		ctx:        ctx,
		cfg:        cfg,
		pool:       pool,
		l1GasPrice: new(big.Int).Mul(defaultGasPriceDivByFactor, big.NewInt(100)).Uint64(), // nolint:gomnd
	}
	gpe.setDefaultGasPrice()
	return gpe
}

// UpdateGasPriceAvg not needed for default strategy.
func (d *DefaultGasPricer) UpdateGasPriceAvg() {
	err := d.pool.SetGasPrices(d.ctx, d.cfg.DefaultGasPriceWei, d.l1GasPrice)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}

func (d *DefaultGasPricer) setDefaultGasPrice() {
	err := d.pool.SetGasPrices(d.ctx, d.cfg.DefaultGasPriceWei, d.l1GasPrice)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}
