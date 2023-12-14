package l2_shared

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// BatchProcessMode is the mode for process a batch (full, incremental, reprocess, nothing)
type BatchProcessMode string

const (
	// FullProcessMode This batch is not on database, so is the first time we process it
	FullProcessMode BatchProcessMode = "full"
	// IncrementalProcessMode We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx
	IncrementalProcessMode BatchProcessMode = "incremental"
	// ReprocessProcessMode We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess it
	ReprocessProcessMode BatchProcessMode = "reprocess"
	// NothingProcessMode The batch is already synchronized, so we don't need to process it
	NothingProcessMode BatchProcessMode = "nothing"
)

// ProcessData contains the data required to process a batch
type ProcessData struct {
	BatchNumber       uint64
	Mode              BatchProcessMode
	OldStateRoot      common.Hash
	OldAccInputHash   common.Hash
	BatchMustBeClosed bool
	// The batch in trusted node, it NEVER will be nil
	TrustedBatch *types.Batch
	// Current batch in state DB, it could be nil
	StateBatch  *state.Batch
	Now         time.Time
	Description string
}

// SyncTrustedStateBatchExecutorSteps is the interface that known how to process a batch
type SyncTrustedStateBatchExecutorSteps interface {
	// FullProcess process a batch that is not on database, so is the first time we process it
	FullProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
	IncrementalProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
	ReProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	// CloseBatch close a batch
	CloseBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) error
}

// SyncTrustedStateBatchExecutorTemplate is a template to sync trusted state. It decide the process mode and call the steps
//
//	the real implementation of the steps is in the SyncTrustedStateBatchExecutorSteps interface that known how to process a batch
type SyncTrustedStateBatchExecutorTemplate struct {
	Steps SyncTrustedStateBatchExecutorSteps
	// CheckBatchTimestampGreaterInsteadOfEqual if true, we consider equal two batches if the timestamp of trusted <= timestamp of state
	// this is because in the permissionless the timestamp of a batch is equal to the timestamp of the l1block where is reported
	// but trusted doesn't known this block and use now() instead. But for sure now() musbe <= l1block.tstamp
	CheckBatchTimestampGreaterInsteadOfEqual bool
	timeProvider                             syncCommon.TimeProvider
}

// NewSyncTrustedStateBatchExecutorTemplate creates a new SyncTrustedStateBatchExecutorTemplate
func NewSyncTrustedStateBatchExecutorTemplate(steps SyncTrustedStateBatchExecutorSteps,
	checkBatchTimestampGreaterInsteadOfEqual bool,
	timeProvider syncCommon.TimeProvider) *SyncTrustedStateBatchExecutorTemplate {
	return &SyncTrustedStateBatchExecutorTemplate{
		Steps:                                    steps,
		CheckBatchTimestampGreaterInsteadOfEqual: checkBatchTimestampGreaterInsteadOfEqual,
		timeProvider:                             timeProvider,
	}
}

// ProcessTrustedBatch processes a trusted batch
func (s *SyncTrustedStateBatchExecutorTemplate) ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, status TrustedState, dbTx pgx.Tx) (*TrustedState, error) {
	log.Debugf("Processing trusted batch: %v", trustedBatch.Number)
	stateCurrentBatch := status.LastTrustedBatches[0]
	statePreviousBatch := status.LastTrustedBatches[1]
	processMode, err := s.getModeForProcessBatch(trustedBatch, stateCurrentBatch, statePreviousBatch, status.LastStateRoot)
	if err != nil {
		log.Error("error getting processMode. Error: ", trustedBatch.Number, err)
		return nil, err
	}
	log.Infof("Batch %v: Processing trusted batch: mode=%s", trustedBatch.Number, processMode.Mode)
	var processBatchResp *state.ProcessBatchResponse = nil
	switch processMode.Mode {
	case NothingProcessMode:
		log.Infof("Batch %v: is already synchronized", trustedBatch.Number)
		err = nil
	case FullProcessMode:
		log.Infof("Batch %v: is not on database, so is the first time we process it", trustedBatch.Number)
		processBatchResp, err = s.Steps.FullProcess(ctx, &processMode, dbTx)
	case IncrementalProcessMode:
		log.Infof("Batch %v: is partially synchronized", trustedBatch.Number)
		processBatchResp, err = s.Steps.IncrementalProcess(ctx, &processMode, dbTx)
	case ReprocessProcessMode:
		log.Infof("Batch %v: is partially synchronized but we don't have intermediate stateRoot so need to be fully reprocessed", trustedBatch.Number)
		processBatchResp, err = s.Steps.ReProcess(ctx, &processMode, dbTx)
	}
	if err != nil {
		log.Errorf("Batch %v: error processing trusted batch. Error: %s", trustedBatch.Number, err)
		return nil, err
	}

	if processMode.BatchMustBeClosed {
		log.Infof("Batch %v: Closing batch", trustedBatch.Number)
		err = checkProcessBatchResultMatchExpected(&processMode, processBatchResp)
		if err != nil {
			log.Error("error closing batch. Error: ", err)
			return nil, err
		}
		err = s.Steps.CloseBatch(ctx, trustedBatch, dbTx)
		if err != nil {
			log.Error("error closing batch. Error: ", err)
			return nil, err
		}
		status.LastStateRoot = nil
	} else {
		if processBatchResp != nil {
			status.LastStateRoot = &StateRootEntry{
				batchNumber: uint64(trustedBatch.Number),
				StateRoot:   processBatchResp.NewStateRoot,
			}
		}
	}

	log.Infof("Batch %v synchronized", trustedBatch.Number)
	return &status, nil
}

func (s *SyncTrustedStateBatchExecutorTemplate) getModeForProcessBatch(trustedNodeBatch *types.Batch, stateBatch *state.Batch, statePreviousBatch *state.Batch, lastStateRoot *StateRootEntry) (ProcessData, error) {
	// Check parameters
	if trustedNodeBatch == nil || statePreviousBatch == nil {
		return ProcessData{}, fmt.Errorf("trustedNodeBatch and statePreviousBatch can't be nil")
	}

	var result ProcessData
	if stateBatch == nil {
		result = ProcessData{
			Mode:         FullProcessMode,
			OldStateRoot: statePreviousBatch.StateRoot,
			Description:  "Batch is not on database, so is the first time we process it",
		}
	} else {
		if checkIfSynced(stateBatch, trustedNodeBatch, s.CheckBatchTimestampGreaterInsteadOfEqual) {
			result = ProcessData{
				Mode:         NothingProcessMode,
				OldStateRoot: common.Hash{},
				Description:  "The batch from Node, and the one in database are the same, already synchronized",
			}
		} else {
			// We have a previous batch, but in node something change
			if lastStateRoot != nil && lastStateRoot.batchNumber == stateBatch.BatchNumber {
				result = ProcessData{
					Mode:         IncrementalProcessMode,
					OldStateRoot: lastStateRoot.StateRoot,
					Description:  "We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx",
				}
			} else {
				result = ProcessData{
					Mode:         ReprocessProcessMode,
					OldStateRoot: statePreviousBatch.StateRoot,
					Description:  "We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess all txs",
				}
			}
		}
	}
	if result.Mode == "" {
		return result, fmt.Errorf("failed to get mode for process batch %v", trustedNodeBatch.Number)
	}
	result.BatchNumber = uint64(trustedNodeBatch.Number)
	result.BatchMustBeClosed = result.Mode != NothingProcessMode && isTrustedBatchClosed(trustedNodeBatch)
	result.StateBatch = stateBatch
	result.TrustedBatch = trustedNodeBatch
	result.OldAccInputHash = statePreviousBatch.AccInputHash
	result.Now = s.timeProvider.Now()
	return result, nil
}

func isTrustedBatchClosed(batch *types.Batch) bool {
	if batch == nil {
		return true
	}
	return batch.StateRoot.String() != state.ZeroHash.String()
}

func checkIfSynced(stateBatch *state.Batch, trustedBatch *types.Batch, checkTimestampGreater bool) bool {
	if stateBatch == nil || trustedBatch == nil {
		log.Infof("checkIfSynced stateBatch or trustedBatch is nil, so is not synced")
		return false
	}
	matchNumber := stateBatch.BatchNumber == uint64(trustedBatch.Number)
	matchGER := stateBatch.GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := stateBatch.LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := stateBatch.StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := stateBatch.Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := false
	if checkTimestampGreater {
		matchTimestamp = uint64(stateBatch.Timestamp.Unix()) >= uint64(trustedBatch.Timestamp)
	} else {
		matchTimestamp = uint64(stateBatch.Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
	}
	matchL2Data := hex.EncodeToString(stateBatch.BatchL2Data) == hex.EncodeToString(trustedBatch.BatchL2Data)

	if matchNumber && matchGER && matchLER && matchSR &&
		matchCoinbase && matchTimestamp && matchL2Data {
		return true
	}
	log.Info("matchNumber", matchNumber)
	log.Info("matchGER", matchGER)
	log.Info("matchLER", matchLER)
	log.Info("matchSR", matchSR)
	log.Info("matchCoinbase", matchCoinbase)
	log.Info("matchTimestamp", matchTimestamp)
	log.Info("matchL2Data", matchL2Data)
	return false
}

func checkStateRootAndLER(batchNumber uint64, expectedStateRoot common.Hash, expectedLER common.Hash, calculatedStateRoot common.Hash, calculatedLER common.Hash) error {
	if calculatedStateRoot != expectedStateRoot {
		return fmt.Errorf("Batch %v: stareRoot calculated [%s] is different from the one in the batch [%s]", batchNumber, calculatedStateRoot, expectedStateRoot)
	}
	if calculatedLER != expectedLER {
		return fmt.Errorf("Batch %v: LocalExitRoot calculated [%s] is different from the one in the batch [%s]", batchNumber, calculatedLER, expectedLER)
	}
	return nil
}

func checkProcessBatchResultMatchExpected(data *ProcessData, processBatchResp *state.ProcessBatchResponse) error {
	var err error = nil
	var trustedBatch = data.TrustedBatch
	if trustedBatch == nil {
		panic("trustedBatch is nil")
	}
	if processBatchResp == nil {
		log.Warnf("Batch %v: Can't check  processBatchResp because is nil, then check store batch in DB", trustedBatch.Number)
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, data.StateBatch.StateRoot, data.StateBatch.LocalExitRoot)
	} else {
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, processBatchResp.NewStateRoot, processBatchResp.NewLocalExitRoot)
	}
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
