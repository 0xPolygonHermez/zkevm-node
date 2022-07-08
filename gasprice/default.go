package gasprice

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

// Default gas price from config is set.
type Default struct {
	cfg  Config
	pool pool
}

// GetAvgGasPrice get default gas price from the pool.
func (d *Default) GetAvgGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := d.pool.GetGasPrice(ctx)
	if errors.Is(err, state.ErrNotFound) {
		return big.NewInt(0), nil
	} else if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(gasPrice), nil
}

// UpdateGasPriceAvg not needed for default strategy.
func (d *Default) UpdateGasPriceAvg(newValue *big.Int) {}

func (d *Default) setDefaultGasPrice() {
	ctx := context.Background()
	err := d.pool.SetGasPrice(ctx, d.cfg.DefaultGasPriceWei)
	if err != nil {
		panic(fmt.Errorf("failed to set default gas price, err: %v", err))
	}
}

// NewDefaultEstimator init default gas price estimator.
func NewDefaultEstimator(cfg Config, pool pool) *Default {
	gpe := &Default{cfg: cfg, pool: pool}
	gpe.setDefaultGasPrice()
	return gpe
}
