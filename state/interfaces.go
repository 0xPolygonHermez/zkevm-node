package state

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type storage interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	StoreGenesisBatch(ctx context.Context, batch Batch, dbTx pgx.Tx) error
	Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error
	ResetForkID(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddBlock(ctx context.Context, block *Block, dbTx pgx.Tx) error
	GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*Block, error)
	AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, dbTx pgx.Tx) error
	GetLatestGlobalExitRoot(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (GlobalExitRoot, time.Time, error)
	GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetBlockNumAndMainnetExitRootByGER(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (uint64, common.Hash, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, tx pgx.Tx) error
	GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*ForcedBatch, error)
	GetForcedBatchesSince(ctx context.Context, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*ForcedBatch, error)
	AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, dbTx pgx.Tx) error
	GetVerifiedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*VerifiedBatch, error)
	GetLastNBatches(ctx context.Context, numBatches uint, dbTx pgx.Tx) ([]*Batch, error)
	GetLastNBatchesByL2BlockNumber(ctx context.Context, l2BlockNumber *uint64, numBatches uint, dbTx pgx.Tx) ([]*Batch, common.Hash, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastBatchTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLatestVirtualBatchTimestamp(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	SetLastBatchInfoSeenOnEthereum(ctx context.Context, lastBatchNumberSeen, lastBatchNumberVerified uint64, dbTx pgx.Tx) error
	SetInitSyncBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error)
	GetBatchByTxHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*Batch, error)
	GetBatchByL2BlockNumber(ctx context.Context, l2BlockNumber uint64, dbTx pgx.Tx) (*Batch, error)
	GetVirtualBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error)
	IsBatchVirtualized(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error)
	IsBatchConsolidated(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error)
	IsSequencingTXSynced(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (bool, error)
	GetProcessingContext(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*ProcessingContext, error)
	GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encodedTxs []string, effectivePercentages []uint8, err error)
	GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, effectivePercentages []uint8, err error)
	GetTxsHashesByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []common.Hash, err error)
	AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx pgx.Tx) error
	GetVirtualBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*VirtualBatch, error)
	UpdateGERInOpenBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) error
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]ForcedBatch, error)
	GetBatchNumberOfL2Block(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	BatchNumberByL2BlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Block, error)
	GetL2BlocksByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]types.Block, error)
	GetLastL2BlockCreatedAt(ctx context.Context, dbTx pgx.Tx) (*time.Time, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error)
	GetTransactionByL2BlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionByL2BlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetL2BlockTransactionCountByHash(ctx context.Context, blockHash common.Hash, dbTx pgx.Tx) (uint64, error)
	GetL2BlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetTransactionEGPLogByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*EffectiveGasPriceLog, error)
	AddL2Block(ctx context.Context, batchNumber uint64, l2Block *types.Block, receipts []*types.Receipt, txsEGPData []StoreTxEGPData, dbTx pgx.Tx) error
	GetLastVirtualizedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastConsolidatedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastVerifiedL2BlockNumberUntilL1Block(ctx context.Context, l1FinalizedBlockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetLastVerifiedBatchNumberUntilL1Block(ctx context.Context, l1BlockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*types.Header, error)
	GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*types.Block, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*VerifiedBatch, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)
	GetLocalExitRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)
	GetBlockNumVirtualBatchByBatchNum(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (uint64, error)
	GetL2BlockByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Block, error)
	GetTxsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
	GetTxsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
	GetL2BlockHeaderByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Header, error)
	GetL2BlockHeaderByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error)
	GetL2BlockHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error)
	IsL2BlockConsolidated(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (bool, error)
	IsL2BlockVirtualized(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (bool, error)
	GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error)
	GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (SyncingInfo, error)
	AddReceipt(ctx context.Context, receipt *types.Receipt, dbTx pgx.Tx) error
	AddLog(ctx context.Context, l *types.Log, dbTx pgx.Tx) error
	GetExitRootByGlobalExitRoot(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (*GlobalExitRoot, error)
	AddSequence(ctx context.Context, sequence Sequence, dbTx pgx.Tx) error
	GetSequences(ctx context.Context, lastVerifiedBatchNumber uint64, dbTx pgx.Tx) ([]Sequence, error)
	GetVirtualBatchToProve(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*Batch, error)
	CheckProofContainsCompleteSequences(ctx context.Context, proof *Proof, dbTx pgx.Tx) (bool, error)
	GetProofReadyToVerify(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*Proof, error)
	GetProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*Proof, *Proof, error)
	AddGeneratedProof(ctx context.Context, proof *Proof, dbTx pgx.Tx) error
	UpdateGeneratedProof(ctx context.Context, proof *Proof, dbTx pgx.Tx) error
	DeleteGeneratedProofs(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	CleanupGeneratedProofs(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	CleanupLockedProofs(ctx context.Context, duration string, dbTx pgx.Tx) (int64, error)
	DeleteUngeneratedProofs(ctx context.Context, dbTx pgx.Tx) error
	GetLastClosedBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error)
	GetLastClosedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
	AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	GetLastTrustedForcedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	AddTrustedReorg(ctx context.Context, reorg *TrustedReorg, dbTx pgx.Tx) error
	CountReorgs(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
	GetLatestGer(ctx context.Context, maxBlockNumber uint64) (GlobalExitRoot, time.Time, error)
	GetBatchByForcedBatchNum(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*Batch, error)
	AddForkID(ctx context.Context, forkID ForkIDInterval, dbTx pgx.Tx) error
	GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID ForkIDInterval, dbTx pgx.Tx) error
	GetNativeBlockHashesInRange(ctx context.Context, fromBlock, toBlock uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*DSL2Block, error)
	GetDSBatches(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*DSBatch, error)
	GetDSL2Blocks(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*DSL2Block, error)
	GetDSL2Transactions(ctx context.Context, firstL2Block, lastL2Block uint64, dbTx pgx.Tx) ([]*DSL2Transaction, error)
	OpenBatchInStorage(ctx context.Context, batchContext ProcessingContext, dbTx pgx.Tx) error
	CloseBatchInStorage(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error
}
