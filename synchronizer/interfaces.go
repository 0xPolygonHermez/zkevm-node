package synchronizer

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Consumer interfaces required by the package.

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	UpdateGasPriceAvg(newValue *big.Int)
}

// localEtherman contains the methods required to interact with ethereum.
type localEtherman interface {
	GetLatestProposedBatchNumber() (uint64, error)
	GetLatestConsolidatedBatchNumber() (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNum uint64) (*types.Block, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx context.Context, txBundleID string) (*state.Block, error)
	SetGenesis(ctx context.Context, genesis state.Genesis, txBundleID string) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error
	SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error
	SetInitSyncBatch(ctx context.Context, batchNumber uint64, txBundleID string) error
	GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64, txBundleID string) (*types.Header, error)
	GetPreviousBlock(ctx context.Context, offset uint64, txBundleID string) (*state.Block, error)

	BeginStateTransaction(ctx context.Context) (string, error)
	RollbackState(ctx context.Context, txBundleID string) error
	CommitState(ctx context.Context, txBundleID string) error

	AddBlock(ctx context.Context, block *state.Block, txBundleID string) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address, txBundleID string) error
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (*state.BatchProcessor, error)
	AddSequencer(ctx context.Context, seq state.Sequencer, txBundleID string) error
	Reset(ctx context.Context, block *state.Block, txBundleID string) error
}
