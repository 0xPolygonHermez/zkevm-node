package profitabilitychecker

import (
	"context"
	"math/big"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	GetSendSequenceFee() (*big.Int, error)
}

// priceGetter is for getting eth/matic price, used for the base tx profitability checker
type priceGetter interface {
	Start(ctx context.Context)
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}
