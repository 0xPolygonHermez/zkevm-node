package state

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// ProcessSequencerBatchV2 is used by the sequencers to process transactions into an open batch for forkID >= ETROG
func (s *State) ProcessSequencerBatchV2(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessSequencerBatchV2 start")

	processBatchResponse, err := s.processBatchV2(ctx, batchNumber, batchL2Data, caller, dbTx)
	if err != nil {
		return nil, err
	}

	result, err := s.convertToProcessBatchResponseV2(processBatchResponse)
	if err != nil {
		return nil, err
	}
	log.Debugf("ProcessSequencerBatchV2 end")
	log.Debugf("*******************************************")
	return result, nil
}

// ProcessBatchV2 processes a batch for forkID >= ETROG
func (s *State) ProcessBatchV2(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessBatchV2 start")

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	forkID := s.GetForkIDByBatchNumber(request.BatchNumber)

	// Create Batch
	var processBatchRequest = &executor.ProcessBatchRequestV2{
		OldBatchNum:      request.BatchNumber - 1,
		Coinbase:         request.Coinbase.String(),
		BatchL2Data:      request.Transactions,
		OldStateRoot:     request.OldStateRoot.Bytes(),
		L1InfoRoot:       request.GlobalExitRoot_V1.Bytes(),
		OldAccInputHash:  request.OldAccInputHash.Bytes(),
		TimestampLimit:   request.TimestampLimit_V2,
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}
	res, err := s.sendBatchRequestToExecutorV2(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, err
	}

	var result *ProcessBatchResponse
	result, err = s.convertToProcessBatchResponseV2(res)
	if err != nil {
		return nil, err
	}

	log.Debugf("ProcessBatchV2 end")
	log.Debugf("*******************************************")

	return result, nil
}

// ExecuteBatchV2 is used by the synchronizer to reprocess batches to compare generated state root vs stored one
// It is also used by the sequencer in order to calculate used zkCounter of a WIPBatch
func (s *State) ExecuteBatchV2(ctx context.Context, batch Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}

	// Get previous batch to get state root and local exit root
	previousBatch, err := s.GetBatchByNumber(ctx, batch.BatchNumber-1, dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:  batch.BatchNumber - 1,
		Coinbase:     batch.Coinbase.String(),
		BatchL2Data:  batch.BatchL2Data,
		OldStateRoot: previousBatch.StateRoot.Bytes(),
		// TODO: Change this to L1InfoRoot
		L1InfoRoot:      batch.GlobalExitRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),
		// TODO: Change this to TimestampLimit
		TimestampLimit: uint64(batch.Timestamp.Unix()),
		// Changed for new sequencer strategy
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		ContextId:        uuid.NewString(),
	}

	// Send Batch to the Executor
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("ExecuteBatchV2[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
	log.Debugf("ExecuteBatchV2[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("ExecuteBatchV2[processBatchRequest.L1InfoRoot]: %v", hex.EncodeToHex(processBatchRequest.L1InfoRoot))
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("ExecuteBatchV2[processBatchRequest.TimestampLimit]: %v", processBatchRequest.TimestampLimit)
	log.Debugf("ExecuteBatchV2[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("ExecuteBatchV2[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

	processBatchResponse, err := s.executorClient.ProcessBatchV2(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorErrorV2(ctx, processBatchResponse.Error, processBatchRequest)
	}

	return processBatchResponse, err
}

func (s *State) processBatchV2(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}

	lastBatches, err := s.GetLastNBatches(ctx, two, dbTx)
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	isBatchClosed, err := s.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if isBatchClosed {
		return nil, ErrBatchAlreadyClosed
	}

	// Check provided batch number is the latest in db
	if lastBatch.BatchNumber != batchNumber {
		return nil, ErrInvalidBatchNumber
	}
	forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:  lastBatch.BatchNumber - 1,
		Coinbase:     lastBatch.Coinbase.String(),
		BatchL2Data:  batchL2Data,
		OldStateRoot: previousBatch.StateRoot.Bytes(),
		// TODO: Update this to L1InfoRoot
		L1InfoRoot:       lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		TimestampLimit:   uint64(lastBatch.Timestamp.Unix()),
		UpdateMerkleTree: cTrue,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	return s.sendBatchRequestToExecutorV2(ctx, processBatchRequest, caller)
}

func (s *State) sendBatchRequestToExecutorV2(ctx context.Context, processBatchRequest *executor.ProcessBatchRequestV2, caller metrics.CallerLabel) (*executor.ProcessBatchResponseV2, error) {
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}
	// Send Batch to the Executor
	if caller != metrics.DiscardCallerLabel {
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.L1InfoRoot]: %v", hex.EncodeToHex(processBatchRequest.L1InfoRoot))
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.TimestampLimit]: %v", processBatchRequest.TimestampLimit)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)
	}
	now := time.Now()
	res, err := s.executorClient.ProcessBatchV2(ctx, processBatchRequest)
	if err != nil {
		log.Errorf("Error s.executorClient.ProcessBatchV2: %v", err)
		log.Errorf("Error s.executorClient.ProcessBatchV2: %s", err.Error())
		log.Errorf("Error s.executorClient.ProcessBatchV2 response: %v", res)
	} else if res.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(res.Error)
		s.eventLog.LogExecutorErrorV2(ctx, res.Error, processBatchRequest)
	}
	elapsed := time.Since(now)
	if caller != metrics.DiscardCallerLabel {
		metrics.ExecutorProcessingTime(string(caller), elapsed)
	}
	log.Infof("Batch: %d took %v to be processed by the executor ", processBatchRequest.OldBatchNum+1, elapsed)

	return res, err
}
