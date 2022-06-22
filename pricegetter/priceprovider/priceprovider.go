package priceprovider

import (
	"context"
	"math/big"
)

// PriceProvider get price from different data sources
type PriceProvider interface {
	// GetEthToMaticPrice getting price from the specified provider
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}
