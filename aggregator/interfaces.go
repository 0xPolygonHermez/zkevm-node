package aggregator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with
// ethereum.
type etherman interface {
	ConsolidateBatch(batchNum *big.Int, proof *proverclient.Proof) (*types.Transaction, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interract with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error)
	GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context) (uint64, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*state.Batch, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64) ([]byte, error)
	GetSequencer(ctx context.Context, address common.Address) (*state.Sequencer, error)
}
