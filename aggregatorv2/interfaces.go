package aggregator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	// TODO: this interface needs to be reviewed
	// ConsolidateBatch(batchNum *big.Int, proof *pb.GetProofResponse) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	// TODO: State methods needs to be reviewed
}
