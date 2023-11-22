package state

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// ProcessSequencerBatchV2 is used by the sequencers to process transactions into an open batch for forkID >= ETROG
func (s *State) ProcessSequencerBatchV2(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*ProcessBatchResponseV2, error) {
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
func (s *State) ProcessBatchV2(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponseV2, error) {
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
		L1InfoRoot:       request.SignificantRoot.Bytes(),
		OldAccInputHash:  request.OldAccInputHash.Bytes(),
		TimestampLimit:   request.SignificantTimestamp,
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}
	res, err := s.sendBatchRequestToExecutorV2(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, err
	}

	var result *ProcessBatchResponseV2
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
func (s *State) ExecuteBatchV2(ctx context.Context, batch Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error) {
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
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:     batch.BatchNumber - 1,
		Coinbase:        batch.Coinbase.String(),
		BatchL2Data:     batch.BatchL2Data,
		OldStateRoot:    previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:  batch.GlobalExitRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),
		EthTimestamp:    uint64(batch.Timestamp.Unix()),
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
	log.Debugf("ExecuteBatchV2[processBatchRequest.GlobalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("ExecuteBatchV2[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("ExecuteBatchV2[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("ExecuteBatchV2[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("ExecuteBatchV2[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("ExecuteBatchV2[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
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

// ProcessAndStoreClosedBatchV2 is used by the Synchronizer to add a closed batch into the data base. Values returned are the new stateRoot,
// the flushID (incremental value returned by executor),
// the ProverID (executor running ID) the result of closing the batch.
func (s *State) ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error) {
	BatchL2Data := processingCtx.BatchL2Data
	if BatchL2Data == nil {
		log.Warnf("Batch %v: ProcessAndStoreClosedBatchV2: processingCtx.BatchL2Data is nil, assuming is empty", processingCtx.BatchNumber)
		var BatchL2DataEmpty []byte
		BatchL2Data = &BatchL2DataEmpty
	}
	// Decode transactions
	forkID := s.GetForkIDByBatchNumber(processingCtx.BatchNumber)
	decodedTransactions, _, _, err := DecodeTxs(*BatchL2Data, forkID)
	if err != nil && !errors.Is(err, ErrInvalidData) {
		log.Debugf("error decoding transactions: %v", err)
		return common.Hash{}, noFlushID, noProverID, err
	}

	// Open the batch and process the txs
	if dbTx == nil {
		return common.Hash{}, noFlushID, noProverID, ErrDBTxNil
	}
	// Avoid writing twice to the DB the BatchL2Data that is going to be written also in the call closeBatch
	processingCtx.BatchL2Data = nil
	if err := s.OpenBatch(ctx, processingCtx, dbTx); err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}
	processed, err := s.processBatch(ctx, processingCtx.BatchNumber, *BatchL2Data, caller, dbTx)
	if err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}

	// Sanity check
	if len(decodedTransactions) != len(processed.Responses) {
		log.Errorf("number of decoded (%d) and processed (%d) transactions do not match", len(decodedTransactions), len(processed.Responses))
	}

	// Filter unprocessed txs and decode txs to store metadata
	// note that if the batch is not well encoded it will result in an empty batch (with no txs)
	for i := 0; i < len(processed.Responses); i++ {
		if !IsStateRootChanged(processed.Responses[i].Error) {
			if executor.IsROMOutOfCountersError(processed.Responses[i].Error) {
				processed.Responses = []*executor.ProcessTransactionResponse{}
				break
			}

			// Remove unprocessed tx
			if i == len(processed.Responses)-1 {
				processed.Responses = processed.Responses[:i]
			} else {
				processed.Responses = append(processed.Responses[:i], processed.Responses[i+1:]...)
			}
			i--
		}
	}

	processedBatch, err := s.convertToProcessBatchResponse(processed)
	if err != nil {
		return common.Hash{}, noFlushID, noProverID, err
	}

	if len(processedBatch.TransactionResponses) > 0 {
		// Store processed txs into the batch
		err = s.StoreTransactions(ctx, processingCtx.BatchNumber, processedBatch.TransactionResponses, nil, dbTx)
		if err != nil {
			return common.Hash{}, noFlushID, noProverID, err
		}
	}

	// Close batch
	return common.BytesToHash(processed.NewStateRoot), processed.FlushId, processed.ProverId, s.CloseBatchInStorage(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
		AccInputHash:  processedBatch.NewAccInputHash,
		BatchL2Data:   *BatchL2Data,
	}, dbTx)
}
