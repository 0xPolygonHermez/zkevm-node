package state

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// ProcessingContext is the necessary data that a batch needs to provide to the runtime,
// without the historical state data (processing receipt from previous batch)
type ProcessingContextV2 struct {
	BatchNumber       uint64
	Coinbase          common.Address
	Timestamp         time.Time    // Batch timeStamp and also TimestampLimit
	L1InfoRoot        *common.Hash // If null is used the current L1InfoRoot
	ForcedBatchNum    *uint64
	BatchL2Data       *[]byte
	ForcedBlockHashL1 *common.Hash
}

// ProcessSequencerBatchV2 is used by the sequencers to process transactions into an open batch for forkID >= ETROG
func (s *State) ProcessSequencerBatchV2(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller metrics.CallerLabel, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessSequencerBatchV2 start")

	processBatchResponse, err := s.processBatchV2(ctx, batchNumber, batchL2Data, nil, nil, nil, caller, dbTx)
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

	if request.SkipFirstChangeL2Block_V2 {
		processBatchRequest.SkipFirstChangeL2Block = cTrue
	}

	if request.SkipWriteBlockInfoRoot_V2 {
		processBatchRequest.SkipWriteBlockInfoRoot = cTrue
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

func (s *State) processBatchV2(ctx context.Context, batchNumber uint64, batchL2Data []byte, timestampLimit *time.Time, l1InfoRoot, forcedBlockHashL1 *common.Hash, caller metrics.CallerLabel, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}

	lastBatches, err := s.GetLastNBatches(ctx, 2, dbTx) // nolint:gomnd
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

	var timestampLimitUnix uint64
	if timestampLimit != nil {
		timestampLimitUnix = uint64(timestampLimit.Unix())
	} else {
		timestampLimitUnix = uint64(time.Now().Unix())
	}
	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:  lastBatch.BatchNumber - 1,
		Coinbase:     lastBatch.Coinbase.String(),
		BatchL2Data:  batchL2Data,
		OldStateRoot: previousBatch.StateRoot.Bytes(),
		// L1InfoRoot:       l1InfoRoot.Bytes() => This can be nil and so is set later
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		TimestampLimit:   timestampLimitUnix,
		UpdateMerkleTree: cTrue,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}
	if forcedBlockHashL1 != nil {
		log.Debug("Setting ForcedBlockhashL1: ", forcedBlockHashL1)
		processBatchRequest.ForcedBlockhashL1 = forcedBlockHashL1.Bytes()
	}

	if l1InfoRoot != nil {
		processBatchRequest.L1InfoRoot = l1InfoRoot.Bytes()
	} else {
		currentl1InfoRoot := s.GetCurrentL1InfoRoot()
		processBatchRequest.L1InfoRoot = currentl1InfoRoot.Bytes()
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
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.SkipFirstChangeL2Block]: %v", processBatchRequest.SkipFirstChangeL2Block)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.SkipWriteBlockInfoRoot]: %v", processBatchRequest.SkipWriteBlockInfoRoot)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
		log.Debugf("sendBatchRequestToExecutorV2[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)
	}
	now := time.Now()
	res, err := s.executorClient.ProcessBatchV2(ctx, processBatchRequest)
	log.Debug(ProcessBatchResponseToString(res, ""))
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
	log.Infof("batch %d took %v to be processed by the executor ", processBatchRequest.OldBatchNum+1, elapsed)

	return res, err
}

func ProcessBatchResponseToString(r *executor.ProcessBatchResponseV2, prefix string) string {
	res := prefix + "ProcessBatchResponseV2: \n"
	res += prefix + fmt.Sprintf("NewStateRoot: 		%v\n", hex.EncodeToHex(r.NewStateRoot))
	res += prefix + fmt.Sprintf("NewAccInputHash: 	%v\n", hex.EncodeToHex(r.NewAccInputHash))
	res += prefix + fmt.Sprintf("NewLocalExitRoot: 	%v\n", hex.EncodeToHex(r.NewLocalExitRoot))
	res += prefix + fmt.Sprintf("NewBatchNumber: 	%v\n", r.NewBatchNum)
	res += prefix + fmt.Sprintf("Error: 			%v\n", r.Error)
	res += prefix + fmt.Sprintf("FlushId: 			%v\n", r.FlushId)
	res += prefix + fmt.Sprintf("StoredFlushId: 	%v\n", r.StoredFlushId)
	res += prefix + fmt.Sprintf("ProverId: 			%v\n", r.ProverId)
	res += prefix + fmt.Sprintf("GasUsed: 			%v\n", r.GasUsed)
	res += prefix + fmt.Sprintf("ForkId: 			%v\n", r.ForkId)
	for blockIndex, block := range r.BlockResponses {
		newPrefix := prefix + "  " + fmt.Sprintf("BlockResponse[%v]: ", blockIndex)
		res += BlockResponseToString(block, newPrefix)
	}
	return res
}
func BlockResponseToString(r *executor.ProcessBlockResponseV2, prefix string) string {
	res := prefix + "ProcessBlockResponseV2:----------------------------- \n"
	res += prefix + fmt.Sprintf("ParentHash:   %v\n", hex.EncodeToHex(r.ParentHash))
	res += prefix + fmt.Sprintf("Coinbase:     %v\n", r.Coinbase)
	res += prefix + fmt.Sprintf("GasLimit:     %v\n", r.GasLimit)
	res += prefix + fmt.Sprintf("BlockNumber:  %v\n", r.BlockNumber)
	res += prefix + fmt.Sprintf("Timestamp:    %v\n", r.Timestamp)
	res += prefix + fmt.Sprintf("GlobalExitRoot: %v\n", hex.EncodeToHex(r.Ger))
	res += prefix + fmt.Sprintf("BlockHashL1:  %v\n", hex.EncodeToHex(r.BlockHashL1))
	res += prefix + fmt.Sprintf("GasUsed:      %v\n", r.GasUsed)
	res += prefix + fmt.Sprintf("BlockInfoRoot:%v\n", hex.EncodeToHex(r.BlockInfoRoot))
	res += prefix + fmt.Sprintf("BlockHash:    %v\n", hex.EncodeToHex(r.BlockHash))
	for txIndex, tx := range r.Responses {
		newPrefix := prefix + "  " + fmt.Sprintf("TransactionResponse[%v]: ", txIndex)
		res += TransactionResponseToString(tx, newPrefix)
	}
	res += prefix + "----------------------------------------------------------------- [Block]\n"

	return res
}

func TransactionResponseToString(r *executor.ProcessTransactionResponseV2, prefix string) string {
	res := prefix + "ProcessTransactionResponseV2:----------------------------------- \n"
	res += prefix + fmt.Sprintf("TxHash: 	%v\n", hex.EncodeToHex(r.TxHash))
	res += prefix + fmt.Sprintf("TxHashL2: 	%v\n", hex.EncodeToHex(r.TxHashL2))
	res += prefix + fmt.Sprintf("Type: 		%v\n", r.Type)
	res += prefix + fmt.Sprintf("Error: 	%v\n", r.Error)
	res += prefix + fmt.Sprintf("GasUsed: 	%v\n", r.GasUsed)
	res += prefix + fmt.Sprintf("GasLeft: 	%v\n", r.GasLeft)
	res += prefix + fmt.Sprintf("GasRefund: %v\n", r.GasRefunded)
	res += prefix + fmt.Sprintf("StateRoot: %v\n", hex.EncodeToHex(r.StateRoot))
	res += prefix + "----------------------------------------------------------------- [Transaction]\n"

	return res
}

// ProcessAndStoreClosedBatchV2 is used by the Synchronizer to add a closed batch into the data base. Values returned are the new stateRoot,
// the flushID (incremental value returned by executor),
// the ProverID (executor running ID) the result of closing the batch.
func (s *State) ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx ProcessingContextV2, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error) {
	debugPrefix := fmt.Sprint("Batch ", processingCtx.BatchNumber, ": ProcessAndStoreClosedBatchV2: ")

	BatchL2Data := processingCtx.BatchL2Data
	if BatchL2Data == nil {
		log.Warnf("%s processingCtx.BatchL2Data is nil, assuming is empty", debugPrefix, processingCtx.BatchNumber)
		var BatchL2DataEmpty []byte
		BatchL2Data = &BatchL2DataEmpty
	}

	if dbTx == nil {
		return common.Hash{}, noFlushID, noProverID, ErrDBTxNil
	}
	// Avoid writing twice to the DB the BatchL2Data that is going to be written also in the call closeBatch
	// TODO: check if is need this

	convertedProcessingContextV1, err := convertProcessingContext(&processingCtx)
	if err != nil {
		log.Errorf("%s error convertProcessingContext: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}
	convertedProcessingContextV1.BatchL2Data = nil
	if err := s.OpenBatch(ctx, *convertedProcessingContextV1, dbTx); err != nil {
		log.Errorf("%s error OpenBatch: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}
	processed, err := s.processBatchV2(ctx, processingCtx.BatchNumber, *BatchL2Data,
		&processingCtx.Timestamp, processingCtx.L1InfoRoot, processingCtx.ForcedBlockHashL1, caller, dbTx)
	if err != nil {
		log.Errorf("%s error processBatchV2: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}

	processedBatch, err := s.convertToProcessBatchResponseV2(processed)
	if err != nil {
		log.Errorf("%s error convertToProcessBatchResponseV2: %v", debugPrefix, err)
		return common.Hash{}, noFlushID, noProverID, err
	}

	if len(processedBatch.BlockResponses) > 0 {
		for _, blockResponse := range processedBatch.BlockResponses {
			err = s.StoreL2Block(ctx, processingCtx.BatchNumber, blockResponse, nil, dbTx)
			if err != nil {
				log.Errorf("%s error StoreL2Block: %v", debugPrefix, err)
				return common.Hash{}, noFlushID, noProverID, err
			}
		}
	}
	return common.BytesToHash(processed.NewStateRoot), processed.FlushId, processed.ProverId, s.CloseBatchInStorage(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
		AccInputHash:  processedBatch.NewAccInputHash,
		BatchL2Data:   *BatchL2Data,
	}, dbTx)
}
