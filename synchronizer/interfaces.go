package synchronizer

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
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
	GetLastBlock(ctx context.Context) (*state.Block, error)
	SetGenesis(ctx context.Context, genesis state.Genesis) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64) (*types.Header, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*state.Block, error)

	BeginStateTransaction(ctx context.Context) (string, error)
	RollbackState(ctx context.Context, txBundleID string) error
	CommitState(ctx context.Context, txBundleID string) error

	AddBlockDBTx(ctx context.Context, txBundleID string, block *state.Block) error
	ConsolidateBatchDBTx(ctx context.Context, txBundleID string, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address) error
	NewBatchProcessorDBTx(ctx context.Context, txBundleID string, sequencerAddress common.Address, stateRoot []byte) (*state.BatchProcessor, error)
	AddSequencerDBTx(ctx context.Context, txBundleID string, seq state.Sequencer) error
	ResetDBTx(ctx context.Context, txBundleID string, block *state.Block) error
}
