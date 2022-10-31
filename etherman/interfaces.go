package etherman

import (
	"context"
	"math/big"
)

type gasPricer interface {
	// SuggestGasPrice returns the gas price
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}
