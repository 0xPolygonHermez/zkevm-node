package txprofitabilitychecker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with
// ethereum.
type etherman interface {
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
}
