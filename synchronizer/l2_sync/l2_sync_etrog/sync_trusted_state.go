package l2_sync_etrog

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrNotImplemented is returned when a method is not implemented
	ErrNotImplemented = errors.New("not implemented")
)

// BatchStepsExecutorEtrogStateInterface contains the methods required to interact with the state.
type BatchStepsExecutorEtrogStateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	//StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, dbTx pgx.Tx) (*ethTypes.Header, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error

	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error

	//ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	//ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	//ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx state.ProcessingContextV2, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	GetCurrentL1InfoRoot() common.Hash
	StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) error
}

// BatchStepsExecutorEtrogSynchronizerInterface contains the methods required to interact with the synchronizer main class.
type BatchStepsExecutorEtrogSynchronizerInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}

// BatchStepsExecutorEtrog is the implementation of the SyncTrustedStateBatchExecutorSteps that
// have the functions to sync a fullBatch, incrementalBatch and reprocessBatch
type BatchStepsExecutorEtrog struct {
	state BatchStepsExecutorEtrogStateInterface
	sync  BatchStepsExecutorEtrogSynchronizerInterface
}

// NewSyncTrustedStateEtrogExecutor creates a new prcessor for sync with L2 batches
func NewSyncTrustedStateEtrogExecutor(zkEVMClient l2_shared.ZkEVMClientInterface,
	state l2_shared.StateInterface, stateBatchExecutor BatchStepsExecutorEtrogStateInterface,
	sync l2_shared.SyncInterface) *l2_shared.SyncTrustedStateTemplate {
	executorSteps := &BatchStepsExecutorEtrog{
		state: stateBatchExecutor,
		sync:  sync}
	//executor := &l2_shared.SyncTrustedStateBatchExecutorTemplate{
	//	Steps: executorSteps,
	//}
	executor := l2_shared.NewSyncTrustedStateBatchExecutorTemplate(executorSteps, true)
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

	processBatchResp, err := b.processAndStoreTxs(ctx, data.TrustedBatch, getProcessRequest(data, b.state.GetCurrentL1InfoRoot()), dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, err
	}

	return processBatchResp, nil
}

// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
func (b *BatchStepsExecutorEtrog) IncrementalProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	var err error
	if err := checkThatL2DataIsIncremental(data); err != nil {
		log.Error("error checkThatL2DataIsIncremental. Error: ", err)
		return nil, err
	}
	batchNumber := uint64(data.TrustedBatch.Number)
	newBatchL2Data := data.TrustedBatch.BatchL2Data[len(data.StateBatch.BatchL2Data):]
	err = b.state.UpdateBatchL2Data(ctx, batchNumber, data.TrustedBatch.BatchL2Data, dbTx)
	if err != nil {
		log.Errorf("error UpdateBatchL2Data batch %d", batchNumber)
		return nil, err
	}
	data.TrustedBatch.BatchL2Data = newBatchL2Data
	processBatchResp, err := b.processAndStoreTxs(ctx, data.TrustedBatch, getProcessRequest(data, b.state.GetCurrentL1InfoRoot()), dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, err
	}
	return processBatchResp, nil
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
		BatchNumber: uint64(trustedBatch.Number),
		Coinbase:    common.HexToAddress(trustedBatch.Coinbase.String()),
		// Instead of using trustedBatch.Timestamp use now, because the prevBatch could have a newer timestamp because
		// use the tstamp of the L1Block where is the virtualization event
		Timestamp:      time.Now(),
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
		log.Infof("Storing trusted tx %+v", block.BlockNumber)
		if err = b.state.StoreL2Block(ctx, uint64(trustedBatch.Number), block, nil, dbTx); err != nil {
			log.Errorf("failed to store block for batch: %v. BlockNumber: %s err:%v", trustedBatch.Number, block.BlockNumber, err)
			return nil, err
		}
	}
	return processBatchResp, nil
}

func getProcessRequest(data *l2_shared.ProcessData, l1InfoRoot common.Hash) state.ProcessRequest {
	request := state.ProcessRequest{
		BatchNumber:       uint64(data.TrustedBatch.Number),
		OldStateRoot:      data.OldStateRoot,
		OldAccInputHash:   data.OldAccInputHash,
		Coinbase:          common.HexToAddress(data.TrustedBatch.Coinbase.String()),
		L1InfoRoot_V2:     l1InfoRoot,
		TimestampLimit_V2: uint64(data.TrustedBatch.Timestamp),
		Transactions:      data.TrustedBatch.BatchL2Data,
	}
	return request
}

func checkThatL2DataIsIncremental(data *l2_shared.ProcessData) error {
	incommingData := data.TrustedBatch.BatchL2Data
	previousData := data.StateBatch.BatchL2Data
	if len(incommingData) < len(previousData) {
		return errors.New("the new batch has less data than the one in node!")
	}
	if len(incommingData) == len(previousData) {
		return errors.New("the new batch has the same data than the one in node!")
	}
	if hash(incommingData) != hash(previousData) {
		return errors.New("the new batch has different data than the one in node!")
	}
	return nil
}

func hash(data []byte) common.Hash {
	sha := sha3.NewLegacyKeccak256()
	sha.Write(data)
	return common.BytesToHash(sha.Sum(nil))
}
