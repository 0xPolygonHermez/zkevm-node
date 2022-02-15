package gasprice

import (
	"math/big"

	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Estimator interface for gas price estimator
type Estimator interface {
	GetAvgGasPrice() (*big.Int, error)
	UpdateGasPriceAvg(newValue *big.Int)
}

// NewEstimator init gas price estimator based on type in config
func NewEstimator(cfg Config, state state.State, pool *pool.PostgresPool) Estimator {
	switch cfg.Type {
	case AllBatchesType:
		return NewGasPriceEstimatorAllBatches()
	case LastNBatchesType:
		return NewGasPriceEstimatorLastNBatches(cfg, state)
	case DefaultType:
		return NewDefaultGasPriceEstimator(cfg, pool)
	}
	return nil
}
