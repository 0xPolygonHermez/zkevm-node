package synchronizerv2

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	etherman "github.com/hermeznetwork/hermez-core/ethermanv2"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
	GetLatestBatchNumber() (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx context.Context) (*etherman.Block, error)
	AddGlobalExitRoot(ctx context.Context, exitRoot *etherman.GlobalExitRoot, tx pgx.Tx) error
	AddForcedBatch(ctx context.Context, forcedBatch *etherman.ForcedBatch, tx pgx.Tx) error
	AddBlock(ctx context.Context, block *etherman.Block, tx pgx.Tx) error
	Reset(ctx context.Context, blockNumber uint64, tx pgx.Tx) error
	GetPreviousBlock(ctx context.Context, offset uint64) (*etherman.Block, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetTrustedBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*state.TrustedBatch, error)
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch state.VirtualBatch, tx pgx.Tx) error
	AddTrustedBatch(ctx context.Context, trustedBatch state.TrustedBatch, tx pgx.Tx) error
	// GetNextForcedBatches returns the next forcedBatches in FIFO order
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, tx pgx.Tx) (*[]state.ForcedBatch, error)

	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	RollbackState(ctx context.Context, tx pgx.Tx) error
	CommitState(ctx context.Context, tx pgx.Tx) error
}
