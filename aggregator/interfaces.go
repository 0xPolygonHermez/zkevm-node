package aggregator

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with
// ethereum.
type etherman interface {
	ConsolidateBatch(batchNum *big.Int, proof *pb.GetProofResponse) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interract with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
	GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, txBundleID string) (uint64, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (*state.Batch, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64, txBundleID string) ([]byte, error)
	GetSequencer(ctx context.Context, address common.Address, txBundleID string) (*state.Sequencer, error)
}
