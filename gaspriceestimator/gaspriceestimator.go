package gaspriceestimator

import (
	"math/big"
	"sync"
)

type GasPriceEstimator interface {
	GetAvgGasPrice() *big.Int
	UpdateGasPriceAvg(newValue *big.Int)
}

type AllBlocks struct {
	// Average gas price (rolling average)
	averageGasPrice      *big.Int // The average gas price that gets queried
	averageGasPriceCount *big.Int // Param used in the avg. gas price calculation

	agpMux sync.Mutex // Mutex for the averageGasPrice calculation
}

func NewGasPriceEstimatorAllBlocks() *AllBlocks {
	return &AllBlocks{
		averageGasPrice:      big.NewInt(0),
		averageGasPriceCount: big.NewInt(0),
	}
}

// UpdateGasPriceAvg Updates the rolling average value of the gas price
func (g *AllBlocks) UpdateGasPriceAvg(newValue *big.Int) {
	g.agpMux.Lock()

	g.averageGasPriceCount.Add(g.averageGasPriceCount, big.NewInt(1))

	differential := big.NewInt(0)
	differential.Div(newValue.Sub(newValue, g.averageGasPrice), g.averageGasPriceCount)

	g.averageGasPrice.Add(g.averageGasPrice, differential)

	g.agpMux.Unlock()
}

func (g *AllBlocks) GetAvgGasPrice() *big.Int {
	return g.averageGasPrice
}
