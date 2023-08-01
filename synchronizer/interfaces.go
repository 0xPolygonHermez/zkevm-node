package synchronizer

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*ethTypes.Block, error)
	GetLatestBatchNumber() (uint64, error)
	GetTrustedSequencerURL() (string, error)
	VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error)
	GetForks(ctx context.Context, genBlockNumber uint64) ([]state.ForkIDInterval, error)
	GetLatestVerifiedBatchNum() (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
	AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, dbTx pgx.Tx) error
	AddBlock(ctx context.Context, block *state.Block, dbTx pgx.Tx) error
	Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error
	GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*state.Block, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	AddVerifiedBatch(ctx context.Context, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error
	ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	SetGenesis(ctx context.Context, block state.Block, genesis state.Genesis, dbTx pgx.Tx) ([]byte, error)
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error
	GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)
	ExecuteBatch(ctx context.Context, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
	AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	AddTrustedReorg(ctx context.Context, trustedReorg *state.TrustedReorg, dbTx pgx.Tx) error
	GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*ethTypes.Transaction, error)
	ResetForkID(ctx context.Context, batchNumber, forkID uint64, version string, dbTx pgx.Tx) error
	GetForkIDTrustedReorgCount(ctx context.Context, forkID uint64, version string, dbTx pgx.Tx) (uint64, error)
	UpdateForkIDIntervals(intervals []state.ForkIDInterval)
	SetLastBatchInfoSeenOnEthereum(ctx context.Context, lastBatchNumberSeen, lastBatchNumberVerified uint64, dbTx pgx.Tx) error
	SetInitSyncBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	GetStoredFlushID(ctx context.Context) (uint64, string, error)
}

type ethTxManager interface {
	Reorg(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) error
}

type poolInterface interface {
	DeleteReorgedTransactions(ctx context.Context, txs []*ethTypes.Transaction) error
	StoreTx(ctx context.Context, tx ethTypes.Transaction, ip string, isWIP bool) error
}

type zkEVMClientInterface interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}
