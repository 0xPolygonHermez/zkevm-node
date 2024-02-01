package sequencer

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error
	DeleteFailedTransactionsOlderThan(ctx context.Context, date time.Time) error
	DeleteTransactionByHash(ctx context.Context, hash common.Hash) error
	MarkWIPTxsAsPending(ctx context.Context) error
	GetNonWIPPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool, failedReason *string) error
	GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, error)
	UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error
	GetGasPrices(ctx context.Context) (pool.GasPrices, error)
	GetDefaultMinGasPriceAllowed() uint64
	GetL1AndL2GasPrice() (uint64, uint64)
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	GetSendSequenceFee(numBatches uint64) (*big.Int, error)
	TrustedSequencer() (common.Address, error)
	GetLatestBatchNumber() (uint64, error)
	GetLatestBlockTimestamp(ctx context.Context) (uint64, error)
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, effectivePercentages []uint8, err error)
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error)
	GetNonceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error)
	GetLastStateRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error)
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	CloseWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ExecuteBatch(ctx context.Context, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error)
	ExecuteBatchV2(ctx context.Context, batch state.Batch, L1InfoTreeRoot common.Hash, l1InfoTreeData map[uint32]state.L1DataV2, timestampLimit time.Time, updateMerkleTree bool, skipVerifyL1InfoRoot uint32, forcedBlockHashL1 *common.Hash, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error)
	GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error)
	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	OpenWIPBatch(ctx context.Context, batch state.Batch, dbTx pgx.Tx) error
	GetWIPBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetLastNBatches(ctx context.Context, numBatches uint, dbTx pgx.Tx) ([]*state.Batch, error)
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, globalExitRoot, blockInfoRoot common.Hash, dbTx pgx.Tx) (*state.L2Header, error)
	GetLastClosedBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*state.L2Block, error)
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
	GetLatestGlobalExitRoot(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error)
	GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*state.L2Header, error)
	UpdateWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	GetForcedBatchesSince(ctx context.Context, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error)
	GetLastTrustedForcedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLatestVirtualBatchTimestamp(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	CountReorgs(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLatestGer(ctx context.Context, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error)
	GetLatestL1InfoRoot(ctx context.Context, maxBlockNumber uint64) (state.L1InfoTreeExitRootStorageEntry, error)
	FlushMerkleTree(ctx context.Context, newStateRoot common.Hash) error
	GetStoredFlushID(ctx context.Context) (uint64, string, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	AddL2Block(ctx context.Context, batchNumber uint64, l2Block *state.L2Block, receipts []*types.Receipt, txsEGPData []state.StoreTxEGPData, dbTx pgx.Tx) error
	GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*state.DSL2Block, error)
	GetDSBatches(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, readWIPBatch bool, dbTx pgx.Tx) ([]*state.DSBatch, error)
	GetDSL2Blocks(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*state.DSL2Block, error)
	GetDSL2Transactions(ctx context.Context, firstL2Block, lastL2Block uint64, dbTx pgx.Tx) ([]*state.DSL2Transaction, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root common.Hash) (*big.Int, error)
	StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) error
	BuildChangeL2Block(deltaTimestamp uint32, l1InfoTreeIndex uint32) []byte
	GetL1InfoTreeDataFromBatchL2Data(ctx context.Context, batchL2Data []byte, dbTx pgx.Tx) (map[uint32]state.L1DataV2, common.Hash, common.Hash, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
	GetVirtualBatchParentHash(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetForcedBatchParentHash(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
	GetVirtualBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.VirtualBatch, error)
	GetLatestBatchGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error)
}

type workerInterface interface {
	GetBestFittingTx(resources state.BatchResources) (*TxTracker, error)
	UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.InfoReadWrite) []*TxTracker
	UpdateTxZKCounters(txHash common.Hash, from common.Address, ZKCounters state.ZKCounters)
	AddTxTracker(ctx context.Context, txTracker *TxTracker) (replacedTx *TxTracker, dropReason error)
	MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int) []*TxTracker
	DeleteTx(txHash common.Hash, from common.Address)
	AddPendingTxToStore(txHash common.Hash, addr common.Address)
	DeletePendingTxToStore(txHash common.Hash, addr common.Address)
	NewTxTracker(tx types.Transaction, counters state.ZKCounters, ip string) (*TxTracker, error)
	AddForcedTx(txHash common.Hash, addr common.Address)
	DeleteForcedTx(txHash common.Hash, addr common.Address)
}
