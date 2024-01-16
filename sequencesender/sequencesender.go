package sequencesender

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	ethTxManagerOwner = "sequencer"
	monitoredIDFormat = "sequence-from-%v-to-%v"
)

var (
	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")
)

// SequenceSender represents a sequence sender
type SequenceSender struct {
	cfg                Config
	state              stateInterface
	ethTxManager       ethTxManager
	etherman           etherman
	eventLog           *event.EventLog
	fromVirtualBatch   uint64                   // Initial value to connect to the streaming
	latestVirtualBatch uint64                   // Latest virtualized batch
	wipBatch           uint64                   // Work in progress batch
	sequenceList       []*types.Sequence        // Sequence of batches to be send to L1
	sequenceData       map[uint64]*sequenceData // All the batch info with L2 blocks and tx data
	validStream        bool
	streamClient       *datastreamer.StreamClient
}

type sequenceData struct {
	batchClosed bool
	batchRaw    *state.BatchRawV2
}

// New inits sequence sender
func New(cfg Config, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog) (*SequenceSender, error) {
	s := SequenceSender{
		cfg:          cfg,
		state:        state,
		etherman:     etherman,
		ethTxManager: manager,
		eventLog:     eventLog,
		sequenceData: make(map[uint64]*sequenceData),
		validStream:  false,
	}

	return &s, nil
}

// Start starts the sequence sender
func (s *SequenceSender) Start(ctx context.Context) {
	var err error

	// Create datastream client
	s.streamClient, err = datastreamer.NewClient(s.cfg.StreamClient.Server, state.StreamTypeSequencer)
	if err != nil {
		log.Fatalf("failed to create stream client, error: %v", err)
	} else {
		log.Debugf("[SequenceSender] new stream client")
	}

	// Start datastream client
	err = s.streamClient.Start()
	if err != nil {
		log.Fatalf("failed to start stream client, error: %v", err)
	}

	// Handle the streaming
	s.streamClient.SetProcessEntryFunc(s.handleReceivedDataStream)

	// Get latest virtual state batch from L1
	s.latestVirtualBatch, err = s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Fatalf("error getting latest sequenced batch, error: %v", err)
	} else {
		log.Debugf("[SequenceSender] latest batch number %d", s.latestVirtualBatch)
	}

	// Set starting point of the streaming
	s.fromVirtualBatch = s.latestVirtualBatch
	bookmark := []byte{state.BookMarkTypeBatch}
	bookmark = binary.LittleEndian.AppendUint64(bookmark, s.fromVirtualBatch)
	s.streamClient.FromBookmark = bookmark

	// Current batch to sequence
	s.wipBatch = s.latestVirtualBatch + 1

	// Start receiving the streaming
	err = s.streamClient.ExecCommand(datastreamer.CmdStart)
	if err != nil {
		log.Fatalf("failed to connect to the streaming")
	}

	// Sequence
	ticker := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	for {
		s.tryToSendSequence(ctx, ticker)
	}
}

func (s *SequenceSender) tryToSendSequence(ctx context.Context, ticker *time.Ticker) {
	return

	retry := false
	// process monitored sequences before starting a next cycle
	s.ethTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
		if result.Status == ethtxmanager.MonitoredTxStatusFailed {
			retry = true
			mTxResultLogger := ethtxmanager.CreateMonitoredTxResultLogger(ethTxManagerOwner, result)
			mTxResultLogger.Error("failed to send sequence, TODO: review this fatal and define what to do in this case")
		}
	}, nil)

	if retry {
		return
	}

	// Check if synchronizer is up to date
	if !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}

	// Check if should send sequence to L1
	log.Infof("getting sequences to send")
	sequences, err := s.getSequencesToSend(ctx)
	if err != nil || len(sequences) == 0 {
		if err != nil {
			log.Errorf("error getting sequences: %v", err)
		} else {
			log.Info("waiting for sequences to be worth sending to L1")
		}
		waitTick(ctx, ticker)
		return
	}

	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return
	}

	// Send sequences to L1
	sequenceCount := len(sequences)
	log.Infof(
		"sending sequences to L1. From batch %d to batch %d",
		lastVirtualBatchNum+1, lastVirtualBatchNum+uint64(sequenceCount),
	)
	metrics.SequencesSentToL1(float64(sequenceCount))

	// add sequence to be monitored
	to, data, err := s.etherman.BuildSequenceBatchesTxData(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
	if err != nil {
		log.Error("error estimating new sequenceBatches to add to eth tx manager: ", err)
		return
	}
	firstSequence := sequences[0]
	lastSequence := sequences[len(sequences)-1]
	monitoredTxID := fmt.Sprintf(monitoredIDFormat, firstSequence.BatchNumber, lastSequence.BatchNumber)
	err = s.ethTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to, nil, data, s.cfg.GasOffset, nil)
	if err != nil {
		mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to)
		mTxLogger.Errorf("error to add sequences tx to eth tx manager: ", err)
		return
	}
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *SequenceSender) getSequencesToSend(ctx context.Context) ([]types.Sequence, error) {
	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get last virtual batch num, err: %w", err)
	}

	currentBatchNumToSequence := lastVirtualBatchNum + 1
	sequences := []types.Sequence{}
	// var estimatedGas uint64

	var tx *ethTypes.Transaction

	// Add sequences until too big for a single L1 tx or last batch is reached
	for {
		//Check if the next batch belongs to a new forkid, in this case we need to stop sequencing as we need to
		//wait the upgrade of forkid is completed and s.cfg.NumBatchForkIdUpgrade is disabled (=0) again
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber + 1)) {
			return nil, fmt.Errorf("aborting sequencing process as we reached the batch %d where a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber+1)
		}

		// Check if batch is closed
		isClosed, err := s.state.IsBatchClosed(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, err
		}
		if !isClosed {
			// Reached current (WIP) batch
			break
		}
		// Add new sequence
		batch, err := s.state.GetBatchByNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, err
		}

		seq := types.Sequence{
			GlobalExitRoot: batch.GlobalExitRoot,   //TODO: set empty for regular batches
			Timestamp:      batch.Timestamp.Unix(), //TODO: set empty for regular batches
			BatchL2Data:    batch.BatchL2Data,
			BatchNumber:    batch.BatchNumber,
		}

		if batch.ForcedBatchNum != nil {
			//TODO: Assign GER, timestamp(forcedAt) and l1block.parentHash to seq
			forcedBatch, err := s.state.GetForcedBatch(ctx, *batch.ForcedBatchNum, nil)
			if err != nil {
				return nil, err
			}

			// Get L1 block for the forced batch
			fbL1Block, err := s.state.GetBlockByNumber(ctx, forcedBatch.BlockNumber, nil)
			if err != nil {
				return nil, err
			}

			seq.GlobalExitRoot = forcedBatch.GlobalExitRoot
			seq.ForcedBatchTimestamp = forcedBatch.ForcedAt.Unix()
			seq.PrevBlockHash = fbL1Block.ParentHash
		}

		sequences = append(sequences, seq)
		// Check if can be send
		tx, err = s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
		if err == nil && tx.Size() > s.cfg.MaxTxSizeForL1 {
			metrics.SequencesOvesizedDataError()
			log.Infof("oversized Data on TX oldHash %s (txSize %d > %d)", tx.Hash(), tx.Size(), s.cfg.MaxTxSizeForL1)
			err = ErrOversizedData
		}
		if err != nil {
			log.Infof("Handling estimage gas send sequence error: %v", err)
			sequences, err = s.handleEstimateGasSendSequenceErr(ctx, sequences, currentBatchNumToSequence, err)
			if sequences != nil {
				// Handling the error gracefully, re-processing the sequence as a sanity check
				_, err = s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
				return sequences, err
			}
			return sequences, err
		}
		// estimatedGas = tx.Gas()

		//Check if the current batch is the last before a change to a new forkid, in this case we need to close and send the sequence to L1
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber)) {
			log.Infof("sequence should be sent to L1, as we have reached the batch %d from which a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber)
			return sequences, nil
		}

		// Increase batch num for next iteration
		currentBatchNumToSequence++
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("no batches to be sequenced")
		return nil, nil
	}

	lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Warnf("failed to get last l1 interaction time, err: %v. Sending sequences as a conservative approach", err)
		return sequences, nil
	}
	if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
		// TODO: implement check profitability
		// if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
		log.Info("sequence should be sent to L1, because too long since didn't send anything to L1")
		return sequences, nil
		//}
	}

	log.Info("not enough time has passed since last batch was virtualized, and the sequence could be bigger")
	return nil, nil
}

// handleEstimateGasSendSequenceErr handles an error on the estimate gas. It will return:
// nil, error: impossible to handle gracefully
// sequence, nil: handled gracefully. Potentially manipulating the sequences
// nil, nil: a situation that requires waiting
func (s *SequenceSender) handleEstimateGasSendSequenceErr(
	ctx context.Context,
	sequences []types.Sequence,
	currentBatchNumToSequence uint64,
	err error,
) ([]types.Sequence, error) {
	// Insufficient allowance
	if errors.Is(err, ethman.ErrInsufficientAllowance) {
		return nil, err
	}
	if isDataForEthTxTooBig(err) {
		// Remove the latest item and send the sequences
		log.Infof(
			"Done building sequences, selected batches to %d. Batch %d caused the L1 tx to be too big",
			currentBatchNumToSequence-1, currentBatchNumToSequence,
		)
		sequences = sequences[:len(sequences)-1]
		return sequences, nil
	}

	// Unknown error
	if len(sequences) == 1 {
		// TODO: gracefully handle this situation by creating an L2 reorg
		log.Errorf(
			"Error when estimating gas for BatchNum %d (alone in the sequences): %v",
			currentBatchNumToSequence, err,
		)
	}
	// Remove the latest item and send the sequences
	log.Infof(
		"Done building sequences, selected batches to %d. Batch %d excluded due to unknown error: %v",
		currentBatchNumToSequence, currentBatchNumToSequence+1, err,
	)
	sequences = sequences[:len(sequences)-1]

	return sequences, nil
}

func isDataForEthTxTooBig(err error) bool {
	return errors.Is(err, ethman.ErrGasRequiredExceedsAllowance) ||
		errors.Is(err, ErrOversizedData) ||
		errors.Is(err, ethman.ErrContentLengthTooLarge)
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}

func (s *SequenceSender) isSynced(ctx context.Context) bool {
	lastSyncedBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last isSynced batch, err: %v", err)
		return false
	}
	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last batch num, err: %v", err)
		return false
	}
	if lastBatchNum > lastSyncedBatchNum {
		return true
	}
	lastEthBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last eth batch, err: %v", err)
		return false
	}
	if lastSyncedBatchNum < lastEthBatchNum {
		log.Infof("waiting for the state to be isSynced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return false
	}

	return true
}

func (s *SequenceSender) handleReceivedDataStream(e *datastreamer.FileEntry, c *datastreamer.StreamClient, ss *datastreamer.StreamServer) error {
	switch e.Type {
	case state.EntryTypeL2BlockStart:
		// Handle stream entry: Start L2 Block
		l2BlockStart := state.DSL2BlockStart{}
		l2BlockStart = l2BlockStart.Decode(e.Data)

		// Already virtualized
		if l2BlockStart.BatchNumber <= s.fromVirtualBatch {
			return nil
		} else {
			s.validStream = true
		}

		if l2BlockStart.BatchNumber == s.wipBatch {
			// New block in the current batch
			if s.wipBatch == s.fromVirtualBatch+1 {
				// Initial case after startup
				s.addNewSequenceBatch(l2BlockStart)
			}
		} else if l2BlockStart.BatchNumber > s.wipBatch {
			// New batch in the sequence
			// Close current batch
			s.sequenceData[s.wipBatch].batchClosed = true

			// Create new sequential batch
			s.addNewSequenceBatch(l2BlockStart)
		}

		// Add L2 block data
		s.addNewBatchL2Block(l2BlockStart)

	case state.EntryTypeL2Tx:
		// Handle stream entry: L2 Tx
		if !s.validStream {
			return nil
		}

		l2Tx := state.DSL2Transaction{}
		l2Tx = l2Tx.Decode(e.Data)

		// Add tx data
		s.addNewBlockTx(l2Tx)

	case state.EntryTypeL2BlockEnd:
		// Handle stream entry: End L2 Block
		if !s.validStream {
			return nil
		}

		// TODO: Add end block data

	case state.EntryTypeUpdateGER:
		// Handle stream entry: Update GER
		// TODO: What should I do
	}

	return nil
}

func (s *SequenceSender) addNewSequenceBatch(l2BlockStart state.DSL2BlockStart) {
	if s.sequenceData[l2BlockStart.BatchNumber] == nil {
		log.Infof("..new batch, number %d", l2BlockStart.BatchNumber)

		// Create sequence
		sequence := types.Sequence{
			GlobalExitRoot: l2BlockStart.GlobalExitRoot,
			Timestamp:      l2BlockStart.Timestamp,
			BatchNumber:    l2BlockStart.BatchNumber,
		}

		// Add to the list
		s.sequenceList = append(s.sequenceList, &sequence)

		// Create initial data
		batchRaw := state.BatchRawV2{}
		data := sequenceData{
			batchClosed: false,
			batchRaw:    &batchRaw,
		}
		s.sequenceData[l2BlockStart.BatchNumber] = &data

		// Update wip batch
		s.wipBatch = l2BlockStart.BatchNumber

		// Add first block
		s.addNewBatchL2Block(l2BlockStart)
	}
}

func (s *SequenceSender) addNewBatchL2Block(l2BlockStart state.DSL2BlockStart) {
	log.Infof("....new L2 block, number %d (batch %d)", l2BlockStart.L2BlockNumber, l2BlockStart.BatchNumber)

	// Current batch
	wipBatchRaw := s.sequenceData[s.wipBatch].batchRaw

	// New L2 block raw
	newBlockRaw := state.L2BlockRaw{}

	// Add L2 block
	wipBatchRaw.Blocks = append(wipBatchRaw.Blocks, newBlockRaw)

	// Get current L2 block
	blockIndex, blockRaw := s.getWipL2Block()

	// Fill in data
	if blockIndex > 0 {
		blockRaw.DeltaTimestamp = uint32(time.Now().Unix()) - wipBatchRaw.Blocks[0].DeltaTimestamp
	} else {
		blockRaw.DeltaTimestamp = uint32(time.Now().Unix())
	}
}

func (s *SequenceSender) addNewBlockTx(l2Tx state.DSL2Transaction) {
	log.Infof("......new tx, length %d", l2Tx.EncodedLength)

	// Current L2 block
	_, blockRaw := s.getWipL2Block()

	// New Tx raw
	l2TxRaw := state.L2TxRaw{
		EfficiencyPercentage: l2Tx.EffectiveGasPricePercentage,
	}

	// Add Tx
	blockRaw.Transactions = append(blockRaw.Transactions, l2TxRaw)
}

func (s *SequenceSender) getWipL2Block() (uint64, *state.L2BlockRaw) {
	// Current batch
	wipBatchRaw := s.sequenceData[s.wipBatch].batchRaw

	// Current block
	if len(wipBatchRaw.Blocks) > 0 {
		blockIndex := uint64(len(wipBatchRaw.Blocks)) - 1
		blockRaw := wipBatchRaw.Blocks[blockIndex]
		return blockIndex, &blockRaw
	} else {
		return 0, nil
	}
}
