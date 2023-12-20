package l2_sync_etrog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrNotImplemented is returned when a method is not implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrBatchDataIsNotIncremental is returned when the new batch has different data than the one in node and is not possible to sync
	ErrBatchDataIsNotIncremental = errors.New("the new batch has different data than the one in node")
	// ErrFailExecuteBatch is returned when the batch is not executed correctly
	ErrFailExecuteBatch = errors.New("fail execute batch")
	// ErrNotExpectedBathResult is returned when the batch result is not the expected (must match Trusted)
	ErrNotExpectedBathResult = errors.New("not expected batch result (differ from Trusted Batch)")
)

// StateInterface contains the methods required to interact with the state.
type StateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	UpdateWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	GetCurrentL1InfoRoot() common.Hash
	StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) error
	GetL1InfoRootLeafByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
}

// SyncTrustedBatchExecutorForEtrog is the implementation of the SyncTrustedStateBatchExecutorSteps that
// have the functions to sync a fullBatch, incrementalBatch and reprocessBatch
type SyncTrustedBatchExecutorForEtrog struct {
	state StateInterface
	sync  syncinterfaces.SynchronizerFlushIDManager
}

// NewSyncTrustedBatchExecutorForEtrog creates a new prcessor for sync with L2 batches
func NewSyncTrustedBatchExecutorForEtrog(zkEVMClient syncinterfaces.ZKEVMClientTrustedBatchesGetter,
	state l2_shared.StateInterface, stateBatchExecutor StateInterface,
	sync syncinterfaces.SynchronizerFlushIDManager, timeProvider syncCommon.TimeProvider) *l2_shared.TrustedBatchesRetrieve {
	executorSteps := &SyncTrustedBatchExecutorForEtrog{
		state: stateBatchExecutor,
		sync:  sync,
	}

	executor := l2_shared.NewProcessorTrustedBatchSync(executorSteps, true, timeProvider)
	a := l2_shared.NewTrustedBatchesRetrieve(executor, zkEVMClient, state, sync, *l2_shared.NewTrustedStateManager(timeProvider, time.Hour))
	return a
}

// NothingProcess process a batch that is already on database and updated, so it is not going to be processed again. Maybe it needs to be close
func (b *SyncTrustedBatchExecutorForEtrog) NothingProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err := b.closeBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Error("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}
	data.StateBatch.WIP = !data.BatchMustBeClosed
	return &l2_shared.ProcessResponse{
		ProcessBatchResponse:                nil,
		ClearCache:                          false,
		UpdateBatchWithProcessBatchResponse: false,
		UpdateBatch:                         data.StateBatch,
	}, nil
}

// FullProcess process a batch that is not on database, so is the first time we process it
func (b *SyncTrustedBatchExecutorForEtrog) FullProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	log.Debugf("%s FullProcess", data.DebugPrefix, uint64(data.TrustedBatch.Number))

	err := b.openBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
	if err != nil {
		log.Errorf("%s error openning batch. Error: %v", data.DebugPrefix, err)
		return nil, err
	}
	l1InfoRoot := b.state.GetCurrentL1InfoRoot()
	l1InfoTree, err := b.state.GetL1InfoRootLeafByL1InfoRoot(ctx, l1InfoRoot, dbTx)
	if err != nil {
		log.Errorf("%s error getting L1InfoRootLeafByL1InfoRoot: %v.", data.DebugPrefix, l1InfoRoot)
		return nil, err
	}
	debugStr := data.DebugPrefix
	processBatchResp, err := b.processAndStoreTxs(ctx, data.TrustedBatch, b.getProcessRequest(data, l1InfoTree), dbTx, debugStr)
	if err != nil {
		log.Error("%s error procesingAndStoringTxs. Error: ", debugStr, err)
		return nil, err
	}

	err = batchResultSanityCheck(data, processBatchResp, debugStr)
	if err != nil {
		// TODO: Remove this fatal
		log.Fatalf("%s error batchResultSanityCheck. Error: %s", data.DebugPrefix, err.Error())
		return nil, err
	}

	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err = b.closeBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Error("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	} else {
		log.Debugf("%s updateWIPBatch", data.DebugPrefix)
		err = b.updateWIPBatch(ctx, data, processBatchResp, dbTx)
		if err != nil {
			log.Errorf("%s error updateWIPBatch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}

	resultBatch, err := b.state.GetBatchByNumber(ctx, uint64(data.TrustedBatch.Number), dbTx)
	if err != nil {
		log.Error("%s error getting batch. Error: ", data.DebugPrefix, err)
		return nil, err
	}

	res := l2_shared.ProcessResponse{
		ProcessBatchResponse:                processBatchResp,
		ClearCache:                          false,
		UpdateBatch:                         resultBatch,
		UpdateBatchWithProcessBatchResponse: true,
	}
	return &res, nil
}

// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
func (b *SyncTrustedBatchExecutorForEtrog) IncrementalProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	var err error
	if data == nil || data.TrustedBatch == nil || data.StateBatch == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if err := checkThatL2DataIsIncremental(data); err != nil {
		log.Errorf("%s error checkThatL2DataIsIncremental. Error: %v", data.DebugPrefix, err)
		return nil, err
	}
	madeUpBatch := *data.TrustedBatch

	madeUpBatch.BatchL2Data, err = b.composePartialBatch(data.StateBatch, data.TrustedBatch)
	if err != nil {
		log.Errorf("%s error composePartialBatch batch Error:%w", data.DebugPrefix, err)
		return nil, err
	}
	l1InfoRoot := b.state.GetCurrentL1InfoRoot()
	l1InfoTree, err := b.state.GetL1InfoRootLeafByL1InfoRoot(ctx, l1InfoRoot, dbTx)
	if err != nil {
		log.Errorf("%s error getting L1InfoRootLeafByL1InfoRoot: %v. Error:%w", data.DebugPrefix, l1InfoRoot, err)
		return nil, err
	}
	debugStr := fmt.Sprintf("%s: Batch %d:", data.Mode, uint64(data.TrustedBatch.Number))
	processBatchResp, err := b.processAndStoreTxs(ctx, &madeUpBatch, b.getProcessRequest(data, l1InfoTree), dbTx, debugStr)
	if err != nil {
		log.Errorf("%s error procesingAndStoringTxs. Error: ", data.DebugPrefix, err)
		return nil, err
	}

	err = batchResultSanityCheck(data, processBatchResp, debugStr)
	if err != nil {
		// TODO: Remove this fatal
		log.Fatalf("%s error batchResultSanityCheck. Error: %s", data.DebugPrefix, err.Error())
		return nil, err
	}

	if data.BatchMustBeClosed {
		log.Debugf("%s Closing batch", data.DebugPrefix)
		err = b.closeBatch(ctx, data.TrustedBatch, dbTx, data.DebugPrefix)
		if err != nil {
			log.Errorf("%s error closing batch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	} else {
		log.Debugf("%s updateWIPBatch", data.DebugPrefix)
		err = b.updateWIPBatch(ctx, data, processBatchResp, dbTx)
		if err != nil {
			log.Errorf("%s error updateWIPBatch. Error: ", data.DebugPrefix, err)
			return nil, err
		}
	}

	updatedBatch := *data.StateBatch
	updatedBatch.BatchL2Data = data.TrustedBatch.BatchL2Data
	updatedBatch.WIP = !data.BatchMustBeClosed
	res := l2_shared.ProcessResponse{
		ProcessBatchResponse:                processBatchResp,
		ClearCache:                          false,
		UpdateBatchWithProcessBatchResponse: true,
		UpdateBatch:                         &updatedBatch,
	}
	return &res, nil
}

func (b *SyncTrustedBatchExecutorForEtrog) updateWIPBatch(ctx context.Context, data *l2_shared.ProcessData, processBatchResp *state.ProcessBatchResponse, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   data.BatchNumber,
		StateRoot:     processBatchResp.NewStateRoot,
		LocalExitRoot: data.TrustedBatch.RollupExitRoot,
		BatchL2Data:   data.TrustedBatch.BatchL2Data,
		AccInputHash:  data.TrustedBatch.AccInputHash,
		// TODO: Check what to put here
		//BatchResources: processBatchResp.UsedZkCounters,
	}

	err := b.state.UpdateWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		log.Errorf("%s error UpdateWIPBatch. Error: ", data.DebugPrefix, err)
		return err
	}
	return err
}

// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
func (b *SyncTrustedBatchExecutorForEtrog) ReProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*l2_shared.ProcessResponse, error) {
	log.Warnf("%s needs to be reprocessed! deleting batches from this batch, because it was partially processed but the intermediary stateRoot is lost", data.DebugPrefix)
	err := b.state.ResetTrustedState(ctx, uint64(data.TrustedBatch.Number)-1, dbTx)
	if err != nil {
		log.Warnf("%s error deleting batches from this batch: %v", data.DebugPrefix, err)
		return nil, err
	}
	// From this point is like a new trusted batch
	return b.FullProcess(ctx, data, dbTx)
}

func batchResultSanityCheck(data *l2_shared.ProcessData, processBatchResp *state.ProcessBatchResponse, debugStr string) error {
	if processBatchResp == nil {
		return nil
	}
	if processBatchResp.NewStateRoot == state.ZeroHash {
		return fmt.Errorf("%s processBatchResp.NewStateRoot is ZeroHash. Err: %w", debugStr, ErrNotExpectedBathResult)
	}
	if processBatchResp.NewStateRoot != data.TrustedBatch.StateRoot {
		return fmt.Errorf("%s processBatchResp.NewStateRoot(%s) != data.TrustedBatch.StateRoot(%s). Err: %w",
			processBatchResp.NewStateRoot.String(), data.TrustedBatch.StateRoot.String(), debugStr, ErrNotExpectedBathResult)
	}
	if processBatchResp.NewLocalExitRoot != data.TrustedBatch.LocalExitRoot {
		return fmt.Errorf("%s processBatchResp.NewLocalExitRoot(%s) != data.StateBatch.LocalExitRoot(%s). Err: %w", debugStr,
			processBatchResp.NewLocalExitRoot.String(), data.TrustedBatch.LocalExitRoot.String(), ErrNotExpectedBathResult)
	}
	// We can't check AccInputHash because we dont have timeLimit neither L1InfoRoot used to create the batch
	// is going to be update from L1
	// if processBatchResp.NewAccInputHash != data.TrustedBatch.AccInputHash {
	// 	return fmt.Errorf("%s processBatchResp.	if processBatchResp.NewAccInputHash(%s) != data.TrustedBatch.AccInputHash(%s). Err:%w", debugStr,
	// 		processBatchResp.NewAccInputHash.String(), data.TrustedBatch.AccInputHash.String(), ErrNotExpectedBathResult)
	// }
	return nil
}

// CloseBatch close a batch
func (b *SyncTrustedBatchExecutorForEtrog) closeBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx, debugStr string) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   uint64(trustedBatch.Number),
		StateRoot:     trustedBatch.StateRoot,
		LocalExitRoot: trustedBatch.LocalExitRoot,
		BatchL2Data:   trustedBatch.BatchL2Data,
		AccInputHash:  trustedBatch.AccInputHash,
	}
	log.Debugf("%s closing batch %v", debugStr, trustedBatch.Number)
	if err := b.state.CloseBatch(ctx, receipt, dbTx); err != nil {
		// This is a workaround to avoid closing a batch that was already closed
		if err.Error() != state.ErrBatchAlreadyClosed.Error() {
			log.Errorf("%s error closing batch %d", debugStr, trustedBatch.Number)
			return err
		} else {
			log.Warnf("%s CASE 02: the batch [%d] looks like were not close but in STATE was closed", debugStr, trustedBatch.Number)
		}
	}
	return nil
}

func (b *SyncTrustedBatchExecutorForEtrog) openBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx, debugStr string) error {
	log.Debugf("%s Opening batch %d", debugStr, trustedBatch.Number)
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
		log.Error("%s error opening batch: ", debugStr, trustedBatch.Number)
		return err
	}
	return nil
}

func (b *SyncTrustedBatchExecutorForEtrog) processAndStoreTxs(ctx context.Context, trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx, debugPrefix string) (*state.ProcessBatchResponse, error) {
	if request.OldStateRoot == state.ZeroHash {
		log.Warnf("%s Processing batch with oldStateRoot == zero....", debugPrefix)
	}
	processBatchResp, err := b.state.ProcessBatchV2(ctx, request, true)
	if err != nil {
		log.Errorf("%s error processing sequencer batch for batch: %v error:%v ", debugPrefix, trustedBatch.Number, err)
		return nil, err
	}
	b.sync.PendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("%s Storing %d blocks for batch %v", debugPrefix, len(processBatchResp.BlockResponses), trustedBatch.Number)
	if processBatchResp.IsExecutorLevelError {
		log.Warnf("%s executorLevelError detected. Avoid store txs...", debugPrefix)
		return nil, fmt.Errorf("%s executorLevelError detected err: %w", debugPrefix, ErrFailExecuteBatch)
	} else if processBatchResp.IsRomOOCError {
		log.Warnf("%s romOOCError detected. Avoid store txs...", debugPrefix)
		return nil, fmt.Errorf("%s romOOCError detected.err: %w", debugPrefix, ErrFailExecuteBatch)
	}
	for _, block := range processBatchResp.BlockResponses {
		log.Debugf("%s Storing trusted tx %+v", block.BlockNumber, debugPrefix)
		if err = b.state.StoreL2Block(ctx, uint64(trustedBatch.Number), block, nil, dbTx); err != nil {
			newErr := fmt.Errorf("%s failed to store l2block: %v  err:%w", debugPrefix, block.BlockNumber, err)
			log.Error(newErr.Error())
			return nil, newErr
		}
	}
	log.Infof("%s Batch %v: batchl2data len:%d processed and stored: %s oldStateRoot: %s -> newStateRoot:%s", debugPrefix, trustedBatch.Number, len(request.Transactions), getResponseInfo(processBatchResp),
		request.OldStateRoot.String(), processBatchResp.NewStateRoot.String())
	return processBatchResp, nil
}

func getResponseInfo(response *state.ProcessBatchResponse) string {
	if len(response.BlockResponses) == 0 {
		return "no blocks, no txs"
	}
	minBlock := response.BlockResponses[0].BlockNumber
	maxBlock := response.BlockResponses[len(response.BlockResponses)-1].BlockNumber
	totalTx := 0
	for _, block := range response.BlockResponses {
		totalTx += len(block.TransactionResponses)
	}
	return fmt.Sprintf(" l2block[%v-%v] txs[%v]", minBlock, maxBlock, totalTx)
}

func (b *SyncTrustedBatchExecutorForEtrog) getProcessRequest(data *l2_shared.ProcessData, l1InfoTree state.L1InfoTreeExitRootStorageEntry) state.ProcessRequest {
	request := state.ProcessRequest{
		BatchNumber:     uint64(data.TrustedBatch.Number),
		OldStateRoot:    data.OldStateRoot,
		OldAccInputHash: data.OldAccInputHash,
		Coinbase:        common.HexToAddress(data.TrustedBatch.Coinbase.String()),
		L1InfoRoot_V2:   l1InfoTree.L1InfoTreeRoot,
		//TODO: Fill L1InfoTreeData
		L1InfoTreeData_V2:       map[uint32]state.L1DataV2{0: {GlobalExitRoot: l1InfoTree.L1InfoTreeRoot}},
		TimestampLimit_V2:       uint64(data.TrustedBatch.Timestamp),
		Transactions:            data.TrustedBatch.BatchL2Data,
		ForkID:                  b.state.GetForkIDByBatchNumber(uint64(data.TrustedBatch.Number)),
		SkipVerifyL1InfoRoot_V2: true,
	}
	return request
}

func checkThatL2DataIsIncremental(data *l2_shared.ProcessData) error {
	incommingData := data.TrustedBatch.BatchL2Data
	previousData := data.StateBatch.BatchL2Data
	if len(incommingData) < len(previousData) {
		return fmt.Errorf("L2Data check: the new batch has less data than the one in node err:%w", ErrBatchDataIsNotIncremental)
	}

	if hash(incommingData[:len(previousData)]) != hash(previousData) {
		strDiff := syncCommon.LogComparedBytes("trusted L2BatchData", "state   L2BatchData", incommingData, previousData, 10, 10) //nolint:gomnd
		err := fmt.Errorf("L2Data check: the common part with state dont have same hash (different at: %s) err:%w", strDiff, ErrBatchDataIsNotIncremental)
		log.Error(err.Error())
		return err
	}
	return nil
}

func sumAllL2BlockDeltaTimestamp(rawBatch *state.BatchRawV2) uint32 {
	var sum uint32 = 0
	for _, l2block := range rawBatch.Blocks {
		sum += l2block.DeltaTimestamp
	}
	return sum
}

func (b *SyncTrustedBatchExecutorForEtrog) composePartialBatch(previousBatch *state.Batch, newBatch *types.Batch) (types.ArgBytes, error) {
	debugStr := " composePartialBatch: "
	rawPreviousBatch, err := state.DecodeBatchV2(previousBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}
	debugStr += fmt.Sprintf("previousBatch.blocks: %v (%v) ", len(rawPreviousBatch.Blocks), len(previousBatch.BatchL2Data))
	if len(previousBatch.BatchL2Data) >= len(newBatch.BatchL2Data) {
		return nil, fmt.Errorf("previousBatch.BatchL2Data>=newBatch.BatchL2Data")
	}
	newData := newBatch.BatchL2Data[len(previousBatch.BatchL2Data):]
	rawPartialBatch, err := state.DecodeBatchV2(newData)
	if err != nil {
		return nil, err
	}
	debugStr += fmt.Sprintf(" deltaBatch.blocks: %v (%v) ", len(rawPartialBatch.Blocks), len(newData))

	if len(rawPreviousBatch.Blocks) > 0 {
		// We put in first block the absolute timestamp
		rawPartialBatch.Blocks[0].DeltaTimestamp += sumAllL2BlockDeltaTimestamp(rawPreviousBatch)
		debugStr += fmt.Sprintf(" firstBlock tstamp: %v", rawPartialBatch.Blocks[0].DeltaTimestamp)
	}
	newBatchEncoded, err := state.EncodeBatchV2(rawPartialBatch)
	if err != nil {
		return nil, err
	}
	log.Debug(debugStr)
	return newBatchEncoded, nil
}

func hash(data []byte) common.Hash {
	sha := sha3.NewLegacyKeccak256()
	sha.Write(data)
	return common.BytesToHash(sha.Sum(nil))
}
