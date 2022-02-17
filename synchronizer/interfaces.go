package synchronizer

import "math/big"

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	UpdateGasPriceAvg(newValue *big.Int)
}
