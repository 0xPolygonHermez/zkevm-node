package gasprice

import (
	"context"
	"math/big"
	"sync"
)

// AllBatches struct for all batches avg price strategy.
type AllBatches struct {
	// Average gas price (rolling average)
	averageGasPrice      *big.Int // The average gas price that gets queried
	averageGasPriceCount *big.Int // Param used in the avg. gas price calculation

	agpMux sync.Mutex // Mutex for the averageGasPrice calculation
}

// NewEstimatorAllBatches init gas price estimator for all batches strategy.
func NewEstimatorAllBatches() *AllBatches {
	return &AllBatches{
		averageGasPrice:      big.NewInt(0),
		averageGasPriceCount: big.NewInt(0),
	}
}

// UpdateGasPriceAvg Updates the rolling average value of the gas price.
func (g *AllBatches) UpdateGasPriceAvg(newValue *big.Int) {
	g.agpMux.Lock()

	g.averageGasPriceCount.Add(g.averageGasPriceCount, big.NewInt(1))

	differential := big.NewInt(0)
	differential.Div(newValue.Sub(newValue, g.averageGasPrice), g.averageGasPriceCount)

	g.averageGasPrice.Add(g.averageGasPrice, differential)

	g.agpMux.Unlock()
}

// GetAvgGasPrice get avg gas price from all blocks.
func (g *AllBatches) GetAvgGasPrice(ctx context.Context) (*big.Int, error) {
	return g.averageGasPrice, nil
}
