package l2_sync_etrog

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type BatchStepsExecutorEtrogStateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, dbTx pgx.Tx) (*ethTypes.Header, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error

	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error

	//ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	//ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
}

type BatchStepsExecutorEtrogSynchronizerInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}

type BatchStepsExecutorEtrog struct {
	state BatchStepsExecutorEtrogStateInterface
	sync  BatchStepsExecutorEtrogSynchronizerInterface
}

// NewSyncTrustedStateEtrogExecutor creates a new prcessor for sync with L2 batches
func NewSyncTrustedStateEtrogExecutor(zkEVMClient l2_shared.ZkEVMClientInterface,
	state l2_shared.StateInterface, stateBatchExecutor BatchStepsExecutorEtrogStateInterface,
	sync l2_shared.SyncInterface) *l2_shared.SyncTrustedStateTemplate {
	executorSteps := &BatchStepsExecutorEtrog{state: stateBatchExecutor}
	executor := &l2_shared.SyncTrustedStateBatchExecutorTemplate{
		Steps: executorSteps,
	}
	a := l2_shared.NewSyncTrustedStateTemplate(executor, zkEVMClient, state, sync)
	return a
}

// FullProcess process a batch that is not on database, so is the first time we process it
func (b *BatchStepsExecutorEtrog) FullProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	log.Infof("Batch %d needs to be synchronized", uint64(data.TrustedBatch.Number))
	err := b.openBatch(ctx, data.TrustedBatch, dbTx)
	if err != nil {
		log.Error("error openning batch. Error: ", err)
		return nil, err
	}
	processBatchResp, err := b.processAndStoreTxs(ctx, trustedBatch, getProcessRequest(data), dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, nil, err
	}
	return nil, ErrNotImplemented
}

// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
func (b *BatchStepsExecutorEtrog) IncrementalProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	return nil, ErrNotImplemented
}

// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
func (b *BatchStepsExecutorEtrog) ReProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	log.Warnf("Batch %v: needs to be reprocessed! deleting batches from this batch, because it was partially processed but the intermediary stateRoot is lost", data.TrustedBatch.Number)
	err := b.state.ResetTrustedState(ctx, uint64(data.TrustedBatch.Number)-1, dbTx)
	if err != nil {
		log.Warnf("Batch %v: error deleting batches from this batch: %v", data.TrustedBatch.Number, err)
		return nil, err
	}
	// From this point is like a new trusted batch
	return b.FullProcess(ctx, data, dbTx)
}

// CloseBatch close a batch
func (b *BatchStepsExecutorEtrog) CloseBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   uint64(trustedBatch.Number),
		StateRoot:     trustedBatch.StateRoot,
		LocalExitRoot: trustedBatch.LocalExitRoot,
		BatchL2Data:   trustedBatch.BatchL2Data,
		AccInputHash:  trustedBatch.AccInputHash,
	}
	log.Infof("closing batch %v", trustedBatch.Number)
	if err := b.state.CloseBatch(ctx, receipt, dbTx); err != nil {
		// This is a workaround to avoid closing a batch that was already closed
		if err.Error() != state.ErrBatchAlreadyClosed.Error() {
			log.Errorf("error closing batch %d", trustedBatch.Number)
			return err
		} else {
			log.Warnf("CASE 02: the batch [%d] looks like were not close but in STATE was closed", trustedBatch.Number)
		}
	}
	return nil
}

func (b *BatchStepsExecutorEtrog) openBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("Opening batch %d", trustedBatch.Number)
	var batchL2Data []byte = trustedBatch.BatchL2Data
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	if trustedBatch.ForcedBatchNumber != nil {
		fb := uint64(*trustedBatch.ForcedBatchNumber)
		processCtx.ForcedBatchNum = &fb
	}
	err := b.state.OpenBatch(ctx, processCtx, dbTx)
	if err != nil {
		log.Error("error opening batch: ", trustedBatch.Number)
		return err
	}
	return nil
}

func (b *BatchStepsExecutorEtrog) processAndStoreTxs(ctx context.Context, trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {

	// Now we need to check the batch. ForcedBatches should be already stored in the batch table because this is done by the sequencer
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       trustedBatch.Coinbase,
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		ForcedBatchNum: nil,
		BatchL2Data:    (*[]byte)(&trustedBatch.BatchL2Data),
	}
	newStateRoot, flushID, proverID, err := b.state.ProcessAndStoreClosedBatch(ctx, processCtx, batch.BatchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)

	processBatchResp, err := b.state.ProcessBatchV2(ctx, request, true)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %v", trustedBatch.Number)
		return nil, err
	}
	b.sync.PendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("Storing %d blocks for batch %v", len(processBatchResp.BlockResponses), trustedBatch.Number)
	if processBatchResp.IsExecutorLevelError {
		log.Warn("executorLevelError detected. Avoid store txs...")
		return processBatchResp, nil
	} else if processBatchResp.IsRomOOCError {
		log.Warn("romOOCError detected. Avoid store txs...")
		return processBatchResp, nil
	}
	for _, block := range processBatchResp.BlockResponses {
		for _, tx := range block.TransactionResponses {
			if state.IsStateRootChanged(executor.RomErrorCode(tx.RomError)) {
				log.Infof("TrustedBatch info: %+v", processBatchResp)
				log.Infof("Storing trusted tx %+v", tx)
				if _, err = b.state.StoreTransaction(ctx, uint64(trustedBatch.Number), tx, trustedBatch.Coinbase, uint64(trustedBatch.Timestamp), nil, dbTx); err != nil {
					log.Errorf("failed to store transactions for batch: %v. Tx: %s", trustedBatch.Number, tx.TxHash.String())
					return nil, err
				}
			}
		}
	}
	return processBatchResp, nil
}

/*
	type ProcessRequest struct {
		BatchNumber               uint64
		GlobalExitRoot_V1         common.Hash
		L1InfoRoot_V2             common.Hash
		OldStateRoot              common.Hash
		OldAccInputHash           common.Hash
		Transactions              []byte
		Coinbase                  common.Address
		Timestamp_V1              time.Time
		TimestampLimit_V2         uint64
		Caller                    metrics.CallerLabel
		SkipFirstChangeL2Block_V2 bool
		ForkID                    uint64
	}
*/
func getProcessRequest(data *l2_shared.ProcessData) state.ProcessRequest {
	request := state.ProcessRequest{
		BatchNumber:     uint64(data.TrustedBatch.Number),
		OldStateRoot:    data.OldStateRoot,
		OldAccInputHash: data.OldAccInputHash,
		Coinbase:        common.HexToAddress(data.TrustedBatch.Coinbase.String()),
		Timestamp:       time.Unix(int64(data.TrustedBatch.Timestamp), 0),

		GlobalExitRoot: data.trustedBatch.GlobalExitRoot,
		Transactions:   data.trustedBatch.BatchL2Data,
	}
	return request
}
