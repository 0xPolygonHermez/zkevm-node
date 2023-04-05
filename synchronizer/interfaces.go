package synchronizer

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
	HeaderByNumber(ctx *context.RequestContext, number *big.Int) (*ethTypes.Header, error)
	GetRollupInfoByBlockRange(ctx *context.RequestContext, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx *context.RequestContext, blockNumber uint64) (*ethTypes.Block, error)
	GetLatestBatchNumber() (uint64, error)
	GetTrustedSequencerURL() (string, error)
	VerifyGenBlockNumber(ctx *context.RequestContext, genBlockNumber uint64) (bool, error)
	GetForks(ctx *context.RequestContext) ([]state.ForkIDInterval, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Block, error)
	AddGlobalExitRoot(ctx *context.RequestContext, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
	AddForcedBatch(ctx *context.RequestContext, forcedBatch *state.ForcedBatch, dbTx pgx.Tx) error
	AddBlock(ctx *context.RequestContext, block *state.Block, dbTx pgx.Tx) error
	Reset(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) error
	GetPreviousBlock(ctx *context.RequestContext, offset uint64, dbTx pgx.Tx) (*state.Block, error)
	GetLastBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ResetTrustedState(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) error
	AddVirtualBatch(ctx *context.RequestContext, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	GetNextForcedBatches(ctx *context.RequestContext, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	AddVerifiedBatch(ctx *context.RequestContext, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error
	ProcessAndStoreClosedBatch(ctx *context.RequestContext, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller state.CallerLabel) (common.Hash, error)
	SetGenesis(ctx *context.RequestContext, block state.Block, genesis state.Genesis, dbTx pgx.Tx) ([]byte, error)
	OpenBatch(ctx *context.RequestContext, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx *context.RequestContext, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessSequencerBatch(ctx *context.RequestContext, batchNumber uint64, batchL2Data []byte, caller state.CallerLabel, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	StoreTransactions(ctx *context.RequestContext, batchNum uint64, processedTxs []*state.ProcessTransactionResponse, dbTx pgx.Tx) error
	GetStateRootByBatchNumber(ctx *context.RequestContext, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)
	ExecuteBatch(ctx *context.RequestContext, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*pb.ProcessBatchResponse, error)
	GetLastVerifiedBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetLastVirtualBatchNum(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	AddSequence(ctx *context.RequestContext, sequence state.Sequence, dbTx pgx.Tx) error
	AddAccumulatedInputHash(ctx *context.RequestContext, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	AddTrustedReorg(ctx *context.RequestContext, trustedReorg *state.TrustedReorg, dbTx pgx.Tx) error
	GetReorgedTransactions(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) ([]*ethTypes.Transaction, error)
	ResetForkID(ctx *context.RequestContext, batchNumber, forkID uint64, version string, dbTx pgx.Tx) error
	GetForkIDTrustedReorgCount(ctx *context.RequestContext, forkID uint64, version string, dbTx pgx.Tx) (uint64, error)
	UpdateForkIDIntervals(ctx *context.RequestContext, intervals []state.ForkIDInterval)

	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
}

type ethTxManager interface {
	Reorg(ctx *context.RequestContext, fromBlockNumber uint64, dbTx pgx.Tx) error
}

type poolInterface interface {
	DeleteReorgedTransactions(ctx *context.RequestContext, txs []*ethTypes.Transaction) error
	StoreTx(ctx *context.RequestContext, tx ethTypes.Transaction, ip string, isWIP bool) error
}

type zkEVMClientInterface interface {
	BatchNumber(ctx *context.RequestContext) (uint64, error)
	BatchByNumber(ctx *context.RequestContext, number *big.Int) (*types.Batch, error)
}
