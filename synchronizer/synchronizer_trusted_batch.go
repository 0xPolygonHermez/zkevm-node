package synchronizer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type batchProcessMode string

const (
	// This batch is not on database, so is the first time we process it
	fullProcessMode batchProcessMode = "full"
	// We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx
	incrementalProcessMode batchProcessMode = "incremental"
	// We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess it
	reprocessProcessMode batchProcessMode = "reprocess"
	// The batch is already synchronized, so we don't need to process it
	nothingProcessMode batchProcessMode = "nothing"
)

type processData struct {
	batchNumber       uint64
	mode              batchProcessMode
	oldStateRoot      common.Hash
	oldAccInputHash   common.Hash
	batchMustBeClosed bool
	// The batch in trusted node, it NEVER will be nil
	trustedBatch *types.Batch
	// Current batch in state DB, it could be nil
	stateBatch *state.Batch

	description string
}

// TrustedBatchSynchronizer is an interface for managing the synchronization of trusted batches
type TrustedBatchSynchronizer interface {
	ProcessTrustedBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error
}

type intermediateStateRootEntry struct {
	// Last batch processed
	batchNumber uint64
	// State root for lastBatchNumber.
	// - If not closed is the intermediate state root
	IntermediateStateRoot common.Hash
}

type storeIntermediateStateRoot struct {
	entry *intermediateStateRootEntry
}

// ClientTrustedBatchSynchronizer implements TrustedBatchSynchronizer
type ClientTrustedBatchSynchronizer struct {
	state                      stateInterface
	ctx                        context.Context
	flushIDController          FlushIDController
	storeIntermediateStateRoot storeIntermediateStateRoot
}

func (s *storeIntermediateStateRoot) getIntermediateStateRoot(BatchNumber uint64) (common.Hash, error) {
	if s.entry != nil && s.entry.batchNumber == BatchNumber {
		return s.entry.IntermediateStateRoot, nil
	}
	return common.Hash{}, fmt.Errorf("there is no intermediate state root for batch %v", BatchNumber)
}

func (s *storeIntermediateStateRoot) setIntermediateStateRoot(BatchNumber uint64, stateRoot common.Hash) {
	s.entry = &intermediateStateRootEntry{
		batchNumber:           BatchNumber,
		IntermediateStateRoot: stateRoot,
	}
}

func (s *storeIntermediateStateRoot) clean() {
	s.entry = nil
}

func (s *ClientTrustedBatchSynchronizer) getModeForProcessBatch(trustedNodeBatch *types.Batch, stateBatch *state.Batch, statePreviousBatch *state.Batch) (processData, error) {
	// Check parameters
	if trustedNodeBatch == nil || statePreviousBatch == nil {
		return processData{}, fmt.Errorf("trustedNodeBatch and statePreviousBatch can't be nil")
	}

	var result processData
	if stateBatch == nil {
		result = processData{
			mode:         fullProcessMode,
			oldStateRoot: statePreviousBatch.StateRoot,
			description:  "Batch is not on database, so is the first time we process it",
		}
	} else {
		if checkIfSynced(stateBatch, trustedNodeBatch) {
			result = processData{
				mode:         nothingProcessMode,
				oldStateRoot: common.Hash{},
				description:  "The batch from Node, and the one in database are the same, already synchronized",
			}
		} else {
			// We have a previous batch, but in node something change
			stateRoot, err := s.storeIntermediateStateRoot.getIntermediateStateRoot(uint64(stateBatch.BatchNumber))
			if err == nil {
				result = processData{
					mode:         incrementalProcessMode,
					oldStateRoot: stateRoot,
					description:  "We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx",
				}
			} else {
				result = processData{
					mode:         reprocessProcessMode,
					oldStateRoot: statePreviousBatch.StateRoot,
					description:  "We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess all txs",
				}
			}
		}
	}
	if result.mode == "" {
		return result, fmt.Errorf("failed to get mode for process batch %v", trustedNodeBatch.Number)
	}
	result.batchNumber = uint64(trustedNodeBatch.Number)
	result.batchMustBeClosed = result.mode != nothingProcessMode && isTrustedBatchClosed(trustedNodeBatch)
	result.stateBatch = stateBatch
	result.trustedBatch = trustedNodeBatch
	result.oldAccInputHash = statePreviousBatch.AccInputHash
	return result, nil
}

// NewClientTrustedBatchSynchronizer creates a new ClientTrustedBatchSynchronizer
func NewClientTrustedBatchSynchronizer(state stateInterface, ctx context.Context, flushIDController FlushIDController) *ClientTrustedBatchSynchronizer {
	return &ClientTrustedBatchSynchronizer{
		state:                      state,
		ctx:                        ctx,
		flushIDController:          flushIDController,
		storeIntermediateStateRoot: storeIntermediateStateRoot{},
	}
}
func (s *ClientTrustedBatchSynchronizer) getCurrentBatches(trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
	batch, err := s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number), dbTx)
	if err != nil && err != state.ErrStateNotSynchronized {
		log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
		return nil, err
	}
	if batch != nil {
		if batch.BatchNumber != uint64(trustedBatch.Number) {
			panic(fmt.Sprintf("batch.BatchNumber %v != uint64(trustedBatch.Number) %v", batch.BatchNumber, trustedBatch.Number))
		}
	}
	prevBatch, err := s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number-1), dbTx)
	if err != nil && err != state.ErrStateNotSynchronized {
		log.Warnf("failed to get prevBatch %v from local trusted state. Error: %v", trustedBatch.Number-1, err)
		return nil, err
	}
	if prevBatch != nil {
		if prevBatch.BatchNumber != uint64(trustedBatch.Number-1) {
			panic(fmt.Sprintf("prevBatch.BatchNumber %v != uint64(trustedBatch.Number-1) %v", prevBatch.BatchNumber, trustedBatch.Number-1))
		}
	}
	batches := []*state.Batch{batch, prevBatch}
	return batches, nil
}

func isTrustedBatchClosed(batch *types.Batch) bool {
	if batch == nil {
		return true
	}
	return batch.StateRoot.String() != state.ZeroHash.String()
}

func getProcessRequest(data *processData) state.ProcessRequest {
	request := state.ProcessRequest{
		BatchNumber:     uint64(data.trustedBatch.Number),
		OldStateRoot:    data.oldStateRoot,
		OldAccInputHash: data.oldAccInputHash,
		Coinbase:        common.HexToAddress(data.trustedBatch.Coinbase.String()),
		Timestamp:       time.Unix(int64(data.trustedBatch.Timestamp), 0),

		GlobalExitRoot: data.trustedBatch.GlobalExitRoot,
		Transactions:   data.trustedBatch.BatchL2Data,
	}
	return request
}

func (s *ClientTrustedBatchSynchronizer) fullProcessTrustedBatch(data *processData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	request := getProcessRequest(data)
	err := s.openBatch(data.trustedBatch, dbTx)
	if err != nil {
		log.Error("error opening batch. Error: ", err)
		return nil, err
	}
	// Update batchL2Data
	// err = s.state.UpdateBatchL2Data(s.ctx, data.batchNumber, data.trustedBatch.BatchL2Data, dbTx)
	// if err != nil {
	// 	log.Errorf("Batch %v: error UpdateBatchL2Data batch", data.batchNumber)
	// 	return nil, err
	// }

	log.Infof("Processing sequencer for batch %v old_state_root %s", data.trustedBatch.Number, request.OldStateRoot)
	processBatchResp, err := s.processAndStoreTxs(data.trustedBatch, request, dbTx)
	if err != nil {
		log.Error("error processingAndStoringTxs. Error: ", err)
		return nil, err
	}

	return processBatchResp, nil
}

func (s *ClientTrustedBatchSynchronizer) incrementalProcessTrustedBatch(data *processData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	log.Infof("Batch %v: needs to be updated", data.batchNumber)
	request := getProcessRequest(data)
	// Only new txs need to be processed
	forkID := s.state.GetForkIDByBatchNumber(data.batchNumber)
	storedTxs, syncedTxs, _, syncedEfficiencyPercentages, err := decodeTxs(forkID, data.trustedBatch.BatchL2Data, data.stateBatch)
	if err != nil {
		return nil, err
	}
	if len(storedTxs) < len(syncedTxs) {
		txsToBeAdded := syncedTxs[len(storedTxs):]
		if forkID >= forkID5 {
			syncedEfficiencyPercentages = syncedEfficiencyPercentages[len(storedTxs):]
		}

		request.Transactions, err = state.EncodeTransactions(txsToBeAdded, syncedEfficiencyPercentages, forkID)
		if err != nil {
			log.Error("Batch %v: error encoding txs (%d) to be added to the state. Error: %v", data.trustedBatch.Number, len(txsToBeAdded), err)
			return nil, err
		}

		// Update batchL2Data
		err := s.state.UpdateBatchL2Data(s.ctx, data.batchNumber, data.trustedBatch.BatchL2Data, dbTx)
		if err != nil {
			log.Errorf("Batch %v: error UpdateBatchL2Data batch", data.batchNumber)
			return nil, err
		}

		log.Debug("request.Transactions: ", common.Bytes2Hex(request.Transactions))
		processBatchResp, err := s.processAndStoreTxs(data.trustedBatch, request, dbTx)
		if err != nil {
			log.Error("error procesingAndStoringTxs. Error: ", err)
			return nil, err
		}

		return processBatchResp, nil
	} else {
		log.Info("Nothing to sync. Node updated. Checking if it is closed")
	}
	return nil, nil
}

func (s *ClientTrustedBatchSynchronizer) reprocessProcessTrustedBatch(data *processData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	// Delete txs that were stored before restart. We need to reprocess all txs because the intermediary stateRoot is only store in memory
	log.Warnf("Batch %v: needs to be reprocessed! deleting batches from this batch, because it was partially processed but the intermediary stateRoot is lost", data.trustedBatch.Number)
	err := s.state.ResetTrustedState(s.ctx, uint64(data.trustedBatch.Number)-1, dbTx)
	if err != nil {
		log.Error("error resetting trusted state. Error: ", err)
		return nil, err
	}
	// From this point is like a new trusted batch
	return s.fullProcessTrustedBatch(data, dbTx)
}

func (s *ClientTrustedBatchSynchronizer) closeBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   uint64(trustedBatch.Number),
		StateRoot:     trustedBatch.StateRoot,
		LocalExitRoot: trustedBatch.LocalExitRoot,
		BatchL2Data:   trustedBatch.BatchL2Data,
		AccInputHash:  trustedBatch.AccInputHash,
	}
	log.Infof("closing batch %v", trustedBatch.Number)
	if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
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

func checkStateRootAndLER(batchNumber uint64, expectedStateRoot common.Hash, expectedLER common.Hash, calculatedStateRoot common.Hash, calculatedLER common.Hash) error {
	if calculatedStateRoot != expectedStateRoot {
		return fmt.Errorf("Batch %v: stareRoot calculated [%s] is different from the one in the batch [%s]", batchNumber, calculatedStateRoot, expectedStateRoot)
	}
	if calculatedLER != expectedLER {
		return fmt.Errorf("Batch %v: LocalExitRoot calculated [%s] is different from the one in the batch [%s]", batchNumber, calculatedLER, expectedLER)
	}
	return nil
}

func checkProcessBatchResultMatchExpected(data *processData, processBatchResp *state.ProcessBatchResponse) error {
	var err error = nil
	var trustedBatch = data.trustedBatch
	if trustedBatch == nil {
		panic("trustedBatch is nil")
	}
	if processBatchResp == nil {
		log.Warnf("Batch %v: Can't check  processBatchResp because is nil, then check store batch in DB", trustedBatch.Number)
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, data.stateBatch.StateRoot, data.stateBatch.LocalExitRoot)
	} else {
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, processBatchResp.NewStateRoot, processBatchResp.NewLocalExitRoot)
	}
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

// ProcessTrustedBatch process the txs and store the batch in the DB
func (s *ClientTrustedBatchSynchronizer) ProcessTrustedBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("Processing trusted batch: %v", trustedBatch.Number)
	tempBatches, err := s.getCurrentBatches(trustedBatch, dbTx)
	if err != nil {
		log.Error("error getting currentBatches. Error: ", trustedBatch.Number, err)
		return err
	}
	stateCurrentBatch := tempBatches[0]
	statePreviousBatch := tempBatches[1]
	processMode, err := s.getModeForProcessBatch(trustedBatch, stateCurrentBatch, statePreviousBatch)
	if err != nil {
		log.Error("error getting processMode. Error: ", trustedBatch.Number, err)
		return err
	}
	log.Infof("Batch %v: Processing trusted batch: mode=%s", trustedBatch.Number, processMode.mode)
	var processBatchResp *state.ProcessBatchResponse = nil
	switch processMode.mode {
	case nothingProcessMode:
		log.Infof("Batch %v: is already synchronized", trustedBatch.Number)
		err = nil
	case fullProcessMode:
		log.Infof("Batch %v: is not on database, so is the first time we process it", trustedBatch.Number)
		processBatchResp, err = s.fullProcessTrustedBatch(&processMode, dbTx)
	case incrementalProcessMode:
		log.Infof("Batch %v: is partially synchronized", trustedBatch.Number)
		processBatchResp, err = s.incrementalProcessTrustedBatch(&processMode, dbTx)
	case reprocessProcessMode:
		log.Infof("Batch %v: is partially synchronized but we don't have intermediate stateRoot so need to be fully reprocessed", trustedBatch.Number)
		processBatchResp, err = s.reprocessProcessTrustedBatch(&processMode, dbTx)
	}
	if err != nil {
		log.Errorf("Batch %v: error processing trusted batch. Error: %s", trustedBatch.Number, err)
		return err
	}

	if processMode.batchMustBeClosed {
		log.Infof("Batch %v: Closing batch", trustedBatch.Number)
		err = checkProcessBatchResultMatchExpected(&processMode, processBatchResp)
		if err != nil {
			log.Error("error closing batch. Error: ", err)
			return err
		}
		err = s.closeBatch(trustedBatch, dbTx)
		if err != nil {
			log.Error("error closing batch. Error: ", err)
			return err
		}
		s.storeIntermediateStateRoot.clean()
	} else {
		if processBatchResp != nil {
			s.storeIntermediateStateRoot.setIntermediateStateRoot(uint64(trustedBatch.Number), processBatchResp.NewStateRoot)
		}
	}

	log.Infof("Batch %v synchronized", trustedBatch.Number)
	return nil
}

func checkIfSynced(stateBatch *state.Batch, trustedBatch *types.Batch) bool {
	if stateBatch == nil || trustedBatch == nil {
		log.Infof("checkIfSynced stateBatch or trustedBatch is nil, so is not synced")
		return false
	}
	matchNumber := stateBatch.BatchNumber == uint64(trustedBatch.Number)
	matchGER := stateBatch.GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := stateBatch.LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := stateBatch.StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := stateBatch.Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := uint64(stateBatch.Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
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

func decodeTxs(forkID uint64, trustedBatchL2Data types.ArgBytes, batch *state.Batch) ([]ethTypes.Transaction, []ethTypes.Transaction, []uint8, []uint8, error) {
	syncedTxs, _, syncedEfficiencyPercentages, err := state.DecodeTxs(trustedBatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, trustedBatchL2Data.Hex())
		return nil, nil, nil, nil, err
	}
	storedTxs, _, storedEfficiencyPercentages, err := state.DecodeTxs(batch.BatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding stored txs from trustedstate. Error: %v, batch.BatchL2Data: %s", err, common.Bytes2Hex(batch.BatchL2Data))
		return nil, nil, nil, nil, err
	}
	log.Debug("len(storedTxs): ", len(storedTxs))
	log.Debug("len(syncedTxs): ", len(syncedTxs))
	return storedTxs, syncedTxs, storedEfficiencyPercentages, syncedEfficiencyPercentages, nil
}

func (s *ClientTrustedBatchSynchronizer) openBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
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
	err := s.state.OpenBatch(s.ctx, processCtx, dbTx)
	if err != nil {
		log.Error("error opening batch: ", trustedBatch.Number)
		return err
	}
	return nil
}

func (s *ClientTrustedBatchSynchronizer) processAndStoreTxs(trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	processBatchResp, err := s.state.ProcessBatch(s.ctx, request, true)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %v", trustedBatch.Number)
		return nil, err
	}
	s.flushIDController.SetPendingFlushIDAndCheckProverID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("Storing transactions %d for batch %v", len(processBatchResp.Responses), trustedBatch.Number)
	for _, tx := range processBatchResp.Responses {
		if err = s.state.StoreTransaction(s.ctx, uint64(trustedBatch.Number), tx, trustedBatch.Coinbase, uint64(trustedBatch.Timestamp), dbTx); err != nil {
			log.Errorf("failed to store transactions for batch: %v", trustedBatch.Number)
			return nil, err
		}
	}
	return processBatchResp, nil
}
