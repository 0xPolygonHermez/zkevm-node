package sequencer

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/context"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	pb "github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	DeleteTransactionsByHashes(ctx *context.RequestContext, hashes []common.Hash) error
	DeleteTransactionByHash(ctx *context.RequestContext, hash common.Hash) error
	MarkWIPTxsAsPending(ctx *context.RequestContext) error
	GetNonWIPPendingTxs(ctx *context.RequestContext, isClaims bool, limit uint64) ([]pool.Transaction, error)
	UpdateTxStatus(ctx *context.RequestContext, hash common.Hash, newStatus pool.TxStatus, isWIP bool) error
	GetTxZkCountersByHash(ctx *context.RequestContext, hash common.Hash) (*state.ZKCounters, error)
	UpdateTxWIPStatus(ctx *context.RequestContext, hash common.Hash, isWIP bool) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	EstimateGasSequenceBatches(sender common.Address, sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetSendSequenceFee(numBatches uint64) (*big.Int, error)
	TrustedSequencer() (common.Address, error)
	GetLatestBatchNumber() (uint64, error)
	GetLastBatchTimestamp() (uint64, error)
	GetLatestBlockTimestamp(ctx *context.RequestContext) (uint64, error)
	BuildSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence) (to *common.Address, data []byte, err error)
	GetLatestBlockNumber(ctx *context.RequestContext) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetTimeForLatestBatchVirtualization(ctx *context.RequestContext, dbTx pgx.Tx) (time.Time, error)
	GetTxsOlderThanNL1Blocks(ctx *context.RequestContext, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
	GetLastVirtualBatchNum(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	IsBatchClosed(ctx *context.RequestContext, batchNum uint64, dbTx pgx.Tx) (bool, error)
	GetBalanceByStateRoot(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	GetNonceByStateRoot(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	GetLastStateRoot(ctx *context.RequestContext, dbTx pgx.Tx) (common.Hash, error)
	ProcessBatch(ctx *context.RequestContext, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	CloseBatch(ctx *context.RequestContext, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ExecuteBatch(ctx *context.RequestContext, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*pb.ProcessBatchResponse, error)
	GetForcedBatch(ctx *context.RequestContext, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error)
	GetLastBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Batch, error)
	GetLastBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	OpenBatch(ctx *context.RequestContext, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	GetLastNBatches(ctx *context.RequestContext, numBatches uint, dbTx pgx.Tx) ([]*state.Batch, error)
	StoreTransaction(ctx *context.RequestContext, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error
	GetLastClosedBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Batch, error)
	GetLastL2Block(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Block, error)
	GetLastBlock(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Block, error)
	GetLatestGlobalExitRoot(ctx *context.RequestContext, maxBlockNumber uint64, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error)
	GetLastL2BlockHeader(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Header, error)
	UpdateBatchL2Data(ctx *context.RequestContext, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
	ProcessSequencerBatch(ctx *context.RequestContext, batchNumber uint64, batchL2Data []byte, caller state.CallerLabel, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	GetForcedBatchesSince(ctx *context.RequestContext, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error)
	GetLastTrustedForcedBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLatestVirtualBatchTimestamp(ctx *context.RequestContext, dbTx pgx.Tx) (time.Time, error)
	CountReorgs(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLatestGer(ctx *context.RequestContext, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error)
	FlushMerkleTree(ctx *context.RequestContext) error
}

type workerInterface interface {
	GetBestFittingTx(resources batchResources) *TxTracker
	UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.InfoReadWrite) []*TxTracker
	UpdateTx(txHash common.Hash, from common.Address, ZKCounters state.ZKCounters)
	AddTxTracker(ctx *context.RequestContext, txTracker *TxTracker)
	MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int) []*TxTracker
	DeleteTx(txHash common.Hash, from common.Address)
	HandleL2Reorg(txHashes []common.Hash)
	NewTxTracker(tx types.Transaction, isClaim bool, counters state.ZKCounters, ip string) (*TxTracker, error)
}

// The dbManager will need to handle the errors inside the functions which don't return error as they will be used async in the other abstractions.
// Also if dbTx is missing this needs also to be handled in the dbManager
type dbManagerInterface interface {
	OpenBatch(ctx *context.RequestContext, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
	CreateFirstBatch(ctx *context.RequestContext, sequencerAddress common.Address) state.ProcessingContext
	GetLastBatchNumber(ctx *context.RequestContext) (uint64, error)
	StoreProcessedTransaction(ctx *context.RequestContext, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error
	DeleteTransactionFromPool(ctx *context.RequestContext, txHash common.Hash) error
	CloseBatch(ctx *context.RequestContext, params ClosingBatchParameters) error
	GetWIPBatch(ctx *context.RequestContext) (*WipBatch, error)
	GetTransactionsByBatchNumber(ctx *context.RequestContext, batchNumber uint64) (txs []types.Transaction, err error)
	GetLastBatch(ctx *context.RequestContext) (*state.Batch, error)
	GetLastNBatches(ctx *context.RequestContext, numBatches uint) ([]*state.Batch, error)
	GetLastClosedBatch(ctx *context.RequestContext) (*state.Batch, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	IsBatchClosed(ctx *context.RequestContext, batchNum uint64) (bool, error)
	GetLatestGer(ctx *context.RequestContext, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error)
	ProcessForcedBatch(forcedBatchNum uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error)
	GetForcedBatchesSince(ctx *context.RequestContext, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error)
	GetLastL2BlockHeader(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Header, error)
	GetLastBlock(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Block, error)
	GetLastTrustedForcedBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetBalanceByStateRoot(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	UpdateTxStatus(ctx *context.RequestContext, hash common.Hash, newStatus pool.TxStatus, isWIP bool) error
	GetLatestVirtualBatchTimestamp(ctx *context.RequestContext, dbTx pgx.Tx) (time.Time, error)
	CountReorgs(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	FlushMerkleTree(ctx *context.RequestContext) error
}

type dbManagerStateInterface interface {
	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
	OpenBatch(ctx *context.RequestContext, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	GetLastVirtualBatchNum(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLastNBatches(ctx *context.RequestContext, numBatches uint, dbTx pgx.Tx) ([]*state.Batch, error)
	StoreTransaction(ctx *context.RequestContext, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error
	CloseBatch(ctx *context.RequestContext, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	IsBatchClosed(ctx *context.RequestContext, batchNum uint64, dbTx pgx.Tx) (bool, error)
	GetTransactionsByBatchNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	GetLastClosedBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Batch, error)
	GetLastBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLastBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Batch, error)
	GetLatestGlobalExitRoot(ctx *context.RequestContext, maxBlockNumber uint64, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error)
	GetLastStateRoot(ctx *context.RequestContext, dbTx pgx.Tx) (common.Hash, error)
	GetLastL2BlockHeader(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Header, error)
	GetLastBlock(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Block, error)
	ExecuteBatch(ctx *context.RequestContext, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*pb.ProcessBatchResponse, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	UpdateBatchL2Data(ctx *context.RequestContext, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
	GetForcedBatch(ctx *context.RequestContext, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error)
	ProcessSequencerBatch(ctx *context.RequestContext, batchNumber uint64, batchL2Data []byte, caller state.CallerLabel, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	GetForcedBatchesSince(ctx *context.RequestContext, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error)
	GetLastTrustedForcedBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetBalanceByStateRoot(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	GetLatestVirtualBatchTimestamp(ctx *context.RequestContext, dbTx pgx.Tx) (time.Time, error)
	CountReorgs(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLatestGer(ctx *context.RequestContext, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error)
	FlushMerkleTree(ctx *context.RequestContext) error
}

type ethTxManager interface {
	Add(ctx *context.RequestContext, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, dbTx pgx.Tx) error
	Result(ctx *context.RequestContext, owner, id string, dbTx pgx.Tx) (ethtxmanager.MonitoredTxResult, error)
	ResultsByStatus(ctx *context.RequestContext, owner string, statuses []ethtxmanager.MonitoredTxStatus, dbTx pgx.Tx) ([]ethtxmanager.MonitoredTxResult, error)
	ProcessPendingMonitoredTxs(ctx *context.RequestContext, owner string, failedResultHandler ethtxmanager.ResultHandler, dbTx pgx.Tx)
}
