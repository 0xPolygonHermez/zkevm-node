package priceprovider

import (
	"context"
	"math/big"
)

// PriceProvider get price from different data sources
type PriceProvider interface {
	// GetPrice getting price from the specified provider
	GetPrice(ctx context.Context) (*big.Float, error)
}
