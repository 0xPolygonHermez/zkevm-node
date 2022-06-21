package txprofitabilitychecker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	GetCurrentSequencerCollateral() (*big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
}

// priceGetter is for getting eth/matic price, used for the base tx profitability checker
type priceGetter interface {
	Start(ctx context.Context)
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}
