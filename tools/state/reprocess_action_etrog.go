package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	cTrue             = 1
	cFalse            = 0
	noFlushID  uint64 = 0
	noProverID string = ""
)

type reprocessActionEtrog struct {
	firstBatchNumber uint64
	lastBatchNumber  uint64
	l2ChainId        uint64
	// If true, when execute a batch write the MT in hashDB
	updateHasbDB             bool
	stopOnError              bool
	preferExecutionStateRoot bool

	state       *state.State
	ctx         context.Context
	output      reprocessingOutputer
	flushIdCtrl flushIDController
}

func (r *reprocessActionEtrog) start() error {
	lastBatch := r.lastBatchNumber
	firstBatchNumber := r.firstBatchNumber

	var i uint64
	for i = firstBatchNumber; i < lastBatch; i++ {
		r.output.startProcessingBatch(i)
		dbTx, err := r.state.BeginStateTransaction(r.ctx)
		if err != nil {
			log.Errorf("error starting state transaction: %v", err)
			return err
		}
		//r.executeBatchBlockByBlock(r.ctx, i, dbTx)
		r.executeBatch(r.ctx, i, dbTx)
		rollbackErr := dbTx.Rollback(r.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d,  rollbackErr: %s, error : %v", i, rollbackErr.Error(), err)
			return rollbackErr
		}
	}
	return nil
}

func (r *reprocessActionEtrog) executeBatchBlockByBlock(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	previouslBatch, err := r.state.GetBatchByNumber(ctx, batchNumber-1, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	originalBatch, err := r.state.GetBatchByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	vbatch, err := r.state.GetVirtualBatchDataByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetVirtualBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	accBatch, err := r.state.GetBatchByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	accBatch.StateRoot = previouslBatch.StateRoot
	accBatch.AccInputHash = previouslBatch.AccInputHash
	rawPreviousBatch, err := state.DecodeBatchV2(originalBatch.BatchL2Data)
	//timeLimit := l1block.ReceivedAt
	for changeL2blockIndex, changeL2block := range rawPreviousBatch.Blocks {
		partialBatchRaw2 := state.BatchRawV2{
			Blocks: []state.L2BlockRaw{changeL2block},
		}
		partialBatchL2Data, err := state.EncodeBatchV2(&partialBatchRaw2)
		if err != nil {
			log.Errorf("error encoding partialBatchRaw2: %v", err)
			return err
		}
		leaves := make(map[uint32]state.L1DataV2)

		l1block, err := r.state.GetBlockByNumber(ctx, vbatch.BlockNumber, dbTx)
		var maxGER common.Hash
		leaves, l1InfoRoot, maxGER, err := r.state.GetL1InfoTreeDataFromBatchL2Data(ctx, accBatch.BatchL2Data, dbTx)
		if err != nil {
			log.Errorf("error getting L1InfoRootLeafByL1InfoRoot batch: %v", batchNumber)
			return err
		}
		if *vbatch.L1InfoRoot != l1InfoRoot {
			log.Debugf("error no matching L1InfoRoot batch: %v %s != %s", batchNumber, *vbatch.L1InfoRoot, l1InfoRoot)
		}
		accBatch.BatchL2Data = partialBatchL2Data

		l1InfoRoot = *vbatch.L1InfoRoot
		accBatch.GlobalExitRoot = maxGER
		processCtx := state.ProcessingContextV2{
			BatchNumber:          accBatch.BatchNumber,
			Coinbase:             accBatch.Coinbase,
			Timestamp:            &l1block.ReceivedAt,
			L1InfoRoot:           l1InfoRoot,
			L1InfoTreeData:       leaves,
			ForcedBatchNum:       accBatch.ForcedBatchNum,
			BatchL2Data:          &accBatch.BatchL2Data,
			SkipVerifyL1InfoRoot: 1,
			GlobalExitRoot:       accBatch.GlobalExitRoot,
		}
		if accBatch.GlobalExitRoot == (common.Hash{}) {
			if len(leaves) > 0 {
				globalExitRoot := leaves[uint32(len(leaves)-1)].GlobalExitRoot
				log.Debugf("Empty GER detected for batch: %d usign GER of last leaf (%d):%s",
					accBatch.BatchNumber,
					uint32(len(leaves)-1),
					globalExitRoot)

				processCtx.GlobalExitRoot = globalExitRoot
				accBatch.GlobalExitRoot = globalExitRoot
			} else {
				log.Debugf("Empty leaves array detected for batch: %d usign GER:%s", accBatch.BatchNumber, processCtx.GlobalExitRoot.String())
			}
		}
		// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
		batchRespose, err := ExecuteBatchV2(ctx, r.state, *accBatch, accBatch.StateRoot, accBatch.AccInputHash, processCtx.L1InfoRoot, leaves, *processCtx.Timestamp, false, processCtx.SkipVerifyL1InfoRoot, processCtx.ForcedBlockHashL1, dbTx)
		if err != nil {
			log.Errorf("error executing L1 batch: %+v, error: %v", accBatch, err)
			return err
		}
		log.Info("Change L2 block index: ", changeL2blockIndex)
		log.Info(batchResponseToString(batchRespose))
		accBatch.StateRoot = common.Hash(batchRespose.NewStateRoot)
		accBatch.AccInputHash = common.Hash(batchRespose.NewAccInputHash)
		accBatch.LocalExitRoot = common.Hash(batchRespose.NewLocalExitRoot)

	}

	if accBatch.AccInputHash != originalBatch.AccInputHash {
		log.Errorf("Accumulated input hash mismatch for batch %d. Expected: %s, got: %s", accBatch.BatchNumber, accBatch.AccInputHash.String(), originalBatch.AccInputHash.String())
		return fmt.Errorf("Accumulated input hash mismatch")
	}
	status := r.checkTrustedState(ctx, *accBatch, originalBatch, accBatch.StateRoot, dbTx)
	if status {
		log.Errorf("Trusted reorg detected %d", batchNumber)
		return fmt.Errorf("Trusted reorg detected")
	}
	return nil
}

func (r *reprocessActionEtrog) executeBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {

	leaves := make(map[uint32]state.L1DataV2)
	batch, err := r.state.GetBatchByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	originalBatch, err := r.state.GetBatchByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	vbatch, err := r.state.GetVirtualBatchDataByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetVirtualBatchByNumber batch: %v %w", batchNumber, err)
		return err
	}
	l1block, err := r.state.GetBlockByNumber(ctx, vbatch.BlockNumber, dbTx)
	var maxGER common.Hash
	leaves, l1InfoRoot, maxGER, err := r.state.GetL1InfoTreeDataFromBatchL2Data(ctx, batch.BatchL2Data, dbTx)
	if err != nil {
		log.Errorf("error getting L1InfoRootLeafByL1InfoRoot batch: %v", batchNumber)
		return err
	}
	if *vbatch.L1InfoRoot != l1InfoRoot {
		log.Debugf("error no matching L1InfoRoot batch: %v %s != %s", batchNumber, *vbatch.L1InfoRoot, l1InfoRoot)
	}
	l1InfoRoot = *vbatch.L1InfoRoot
	batch.GlobalExitRoot = maxGER
	processCtx := state.ProcessingContextV2{
		BatchNumber:          batch.BatchNumber,
		Coinbase:             batch.Coinbase,
		Timestamp:            &l1block.ReceivedAt,
		L1InfoRoot:           l1InfoRoot,
		L1InfoTreeData:       leaves,
		ForcedBatchNum:       batch.ForcedBatchNum,
		BatchL2Data:          &batch.BatchL2Data,
		SkipVerifyL1InfoRoot: 1,
		GlobalExitRoot:       batch.GlobalExitRoot,
	}
	if batch.GlobalExitRoot == (common.Hash{}) {
		if len(leaves) > 0 {
			globalExitRoot := leaves[uint32(len(leaves)-1)].GlobalExitRoot
			log.Debugf("Empty GER detected for batch: %d usign GER of last leaf (%d):%s",
				batch.BatchNumber,
				uint32(len(leaves)-1),
				globalExitRoot)

			processCtx.GlobalExitRoot = globalExitRoot
			batch.GlobalExitRoot = globalExitRoot
		} else {
			log.Debugf("Empty leaves array detected for batch: %d usign GER:%s", batch.BatchNumber, processCtx.GlobalExitRoot.String())
		}
	}
	// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
	batchRespose, err := r.state.ExecuteBatchV2(ctx, *batch, processCtx.L1InfoRoot, leaves, *processCtx.Timestamp, false, processCtx.SkipVerifyL1InfoRoot, processCtx.ForcedBlockHashL1, dbTx)
	if err != nil {
		log.Errorf("error executing L1 batch: %+v, error: %v", batch, err)
		return err
	}
	if batchRespose == nil {
		log.Errorf("error executing L1 batch batchRespose == nil : %+v, error: %v", batch, err)
		return fmt.Errorf("error executing L1 batch batchRespose == nil")
	}
	log.Info(batchResponseToString(batchRespose))
	newRoot := common.BytesToHash(batchRespose.NewStateRoot)
	accumulatedInputHash := common.BytesToHash(batchRespose.NewAccInputHash)
	if batch.AccInputHash != accumulatedInputHash {
		log.Errorf("Accumulated input hash mismatch for batch %d. Expected: %s, got: %s", batch.BatchNumber, batch.AccInputHash.String(), accumulatedInputHash.String())
		return fmt.Errorf("accumulated input hash mismatch")
	}
	//func (s *State) StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *ProcessBlockResponse, txsEGPLog []*EffectiveGasPriceLog, dbTx pgx.Tx) error {
	kk, err := r.state.ConvertToProcessBatchResponseV2(batchRespose)
	if err != nil {
		log.Errorf("error converting to ProcessBatchResponseV2: %+v, error: %v", batchRespose, err)
		return err
	}
	for _, block := range kk.BlockResponses {
		log.Info("BlockResponse: ", block.BlockNumber)
		log.Info("BlockResponse: tx", len(block.TransactionResponses))
		for _, tx := range block.TransactionResponses {
			//log.Info("TransactionResponse: ", tx.Error)
			log.Info("TransactionResponse: logs: ", len(tx.Logs))

		}

		// err := r.state.StoreL2Block(ctx, batch.BatchNumber, block, nil, dbTx)
		// if err != nil {
		// 	log.Errorf("error storing L2 block: %+v, error: %v", block, err)
		// }
	}

	status := r.checkTrustedState(ctx, *batch, originalBatch, newRoot, dbTx)
	if status {
		log.Errorf("Trusted reorg detected %d", batchNumber)
		return fmt.Errorf("Trusted reorg detected")
	}
	return nil
}

func batchResponseToString(resp *executor.ProcessBatchResponseV2) string {
	if resp == nil {
		return "nil"
	}
	res := "\n"
	res += fmt.Sprintf("NewStateRoot: %s\n", common.Bytes2Hex(resp.NewStateRoot))
	res += fmt.Sprintf("NewAccInputHash: %s\n", common.Bytes2Hex(resp.NewAccInputHash))
	res += fmt.Sprintf("NewLocalExitRoot: %s\n", common.Bytes2Hex(resp.NewLocalExitRoot))
	res += fmt.Sprintf("NewBatchNum: %d\n", resp.NewBatchNum)
	if resp.BlockResponses != nil {
		for i, block := range resp.BlockResponses {
			res += fmt.Sprintf("BlockResponse[%d]:BlockNumber %d\n", i, block.BlockNumber)
			res += fmt.Sprintf("BlockResponse[%d]:Timestamp %d\n", i, block.Timestamp)
			res += fmt.Sprintf("BlockResponse[%d]:Error %d\n", i, block.Error)
			res += fmt.Sprintf("BlockResponse[%d]:BlockHash %s\n", i, common.Bytes2Hex(block.BlockHash))
			if block.Logs != nil {
				for j, logs := range block.Logs {
					res += fmt.Sprintf("BlockResponse[%d]:Logs[%d]:%s\n", i, j, logs.String())
				}
			} else {
				res += fmt.Sprintf("BlockResponse[%d]:Logs: nil\n", i)

			}
			if block.Responses != nil {
				for k, tx := range block.Responses {
					res += fmt.Sprintf("BlockResponse[%d]:Responses[%d]:hash:%d\n", i, k, common.Bytes2Hex(tx.TxHash))
					res += fmt.Sprintf("BlockResponse[%d]:Responses[%d]:hashL2:%d\n", i, k, common.Bytes2Hex(tx.TxHashL2))
					res += fmt.Sprintf("BlockResponse[%d]:Responses[%d]:Error:%d\n", i, k, tx.Error)
					res += fmt.Sprintf("BlockResponse[%d]:Responses[%d]:nLogs:%d\n", i, k, len(tx.Logs))
					for l, txLog := range tx.Logs {
						res += fmt.Sprintf("BlockResponse[%d]:Responses[%d]:logs[%d]:%s\n", i, k, l, txLog)
					}

				}
			}
		}
	}

	return res

}

func (p *reprocessActionEtrog) checkTrustedState(ctx context.Context, batch state.Batch, tBatch *state.Batch, newRoot common.Hash, dbTx pgx.Tx) bool {
	//Comp/are virtual state with trusted state
	var reorgReasons strings.Builder
	batchNumStr := fmt.Sprintf("Batch: %d.", batch.BatchNumber)
	if newRoot != tBatch.StateRoot {
		errMsg := batchNumStr + fmt.Sprintf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if hex.EncodeToString(batch.BatchL2Data) != hex.EncodeToString(tBatch.BatchL2Data) {
		errMsg := batchNumStr + fmt.Sprintf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data))
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if batch.GlobalExitRoot.String() != tBatch.GlobalExitRoot.String() {
		errMsg := batchNumStr + fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String()))
	}
	if batch.Timestamp.Unix() < tBatch.Timestamp.Unix() { // TODO: this timestamp will be different in permissionless nodes and the trusted node
		errMsg := batchNumStr + fmt.Sprintf("Invalid timestamp. Virtual timestamp limit(%d) must be greater or equal than Trusted timestamp (%d)\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if batch.Coinbase.String() != tBatch.Coinbase.String() {
		errMsg := batchNumStr + fmt.Sprintf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}

	if reorgReasons.Len() > 0 {
		reason := reorgReasons.String()
		log.Warnf("Trusted reorg detected %d: %s", batch.BatchNumber, reason)
		return true
	}
	return false
}

// ExecuteBatchV2 is used by the synchronizer to reprocess batches to compare generated state root vs stored one
func ExecuteBatchV2(ctx context.Context, st *state.State, batch state.Batch, oldStateRoot common.Hash, oldAccInputHash common.Hash, L1InfoTreeRoot common.Hash, l1InfoTreeData map[uint32]state.L1DataV2, timestampLimit time.Time, updateMerkleTree bool, skipVerifyL1InfoRoot uint32, forcedBlockHashL1 *common.Hash, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error) {

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequestV2{
		OldBatchNum:     batch.BatchNumber - 1,
		Coinbase:        batch.Coinbase.String(),
		BatchL2Data:     batch.BatchL2Data,
		OldStateRoot:    oldStateRoot.Bytes(),
		L1InfoRoot:      L1InfoTreeRoot.Bytes(),
		OldAccInputHash: oldAccInputHash.Bytes(),
		TimestampLimit:  uint64(timestampLimit.Unix()),
		// Changed for new sequencer strategy
		UpdateMerkleTree:     updateMT,
		ChainId:              st.GetChainID(),
		ForkId:               uint64(7),
		ContextId:            uuid.NewString(),
		SkipVerifyL1InfoRoot: skipVerifyL1InfoRoot,
	}

	if forcedBlockHashL1 != nil {
		processBatchRequest.ForcedBlockhashL1 = forcedBlockHashL1.Bytes()
	} else {
		l1InfoTree := make(map[uint32]*executor.L1DataV2)
		for i, v := range l1InfoTreeData {
			l1InfoTree[i] = &executor.L1DataV2{
				GlobalExitRoot: v.GlobalExitRoot.Bytes(),
				BlockHashL1:    v.BlockHashL1.Bytes(),
				MinTimestamp:   v.MinTimestamp,
			}
		}
		processBatchRequest.L1InfoTreeData = l1InfoTree
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
	log.Debugf("ExecuteBatchV2[processBatchRequest.SkipVerifyL1InfoRoot]: %v", processBatchRequest.SkipVerifyL1InfoRoot)
	log.Debugf("ExecuteBatchV2[processBatchRequest.L1InfoTreeData]: %+v", l1InfoTreeData)
	executorClient := st.GetExecutorClient()
	processBatchResponse, err := executorClient.ProcessBatchV2(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
	}

	return processBatchResponse, err
}
