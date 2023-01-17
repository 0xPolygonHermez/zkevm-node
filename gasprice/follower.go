package gasprice

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// Follower struct.
type Follower struct {
	cfg  Config
	pool pool
	ctx  context.Context
	eth  ethermanInterface
}

// newFollowerGasEstimator init default gas price estimator.
func newFollowerGasEstimator(ctx context.Context, cfg Config, pool pool, ethMan ethermanInterface) *Follower {
	gpe := &Follower{
		cfg:  cfg,
		pool: pool,
		ctx:  ctx,
		eth:  ethMan,
	}
	gpe.UpdateGasPriceAvg()
	return gpe
}

// UpdateGasPriceAvg updates the gas price.
func (f *Follower) UpdateGasPriceAvg() {
	ctx := context.Background()
	// Get L1 gasprice
	gp := f.eth.GetL1GasPrice(f.ctx)
	if big.NewInt(0).Cmp(gp) == 0 {
		log.Warn("gas price 0 received. Skipping update...")
		return
	}
	// Apply factor to calculate l2 gasPrice
	factor := big.NewFloat(0).SetFloat64(f.cfg.Factor)
	res := new(big.Float).Mul(factor, big.NewFloat(0).SetInt(gp))

	// Store l2 gasPrice calculated
	result := new(big.Int)
	res.Int(result)
	minGasPrice := big.NewInt(0).SetUint64(f.cfg.DefaultGasPriceWei)
	if minGasPrice.Cmp(result) == 1 { // minGasPrice > result
		log.Warn("setting minGasPrice for L2")
		result = minGasPrice
	}
	log.Debug("Storing L2 gas price: ", result)
	err := f.pool.SetGasPrice(ctx, result.Uint64())
	if err != nil {
		log.Errorf("failed to update gas price in poolDB, err: %v", err)
	}
}
