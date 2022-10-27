package etherman

import (
	"context"
	"math/big"
)

type gasPricer interface {
	// GetGasPrice returns the gas price
	GetGasPrice(ctx context.Context) (*big.Int, error)
}
