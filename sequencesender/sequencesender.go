package sequencesender

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
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
	cfg                 Config
	state               stateInterface
	ethTxManager        ethTxManager
	etherman            etherman
	eventLog            *event.EventLog
	latestVirtualBatch  uint64                     // Latest virtualized batch obtained from L1
	latestSentToL1Batch uint64                     // Latest batch sent to L1
	wipBatch            uint64                     // Work in progress batch
	sequenceList        []uint64                   // Sequence of batch number to be send to L1
	sequenceData        map[uint64]*sequenceData   // All the batch data indexed by batch number
	mutexSequence       sync.Mutex                 // Mutex to update sequence data
	ethTransactions     map[common.Hash]*ethTxData // All the eth tx sent to L1 indexed by hash
	validStream         bool                       // Not valid while receiving data before the desired batch
	fromStreamBatch     uint64                     // Initial batch to connect to the streaming
	latestStreamBatch   uint64                     // Latest batch received by the streaming
	streamClient        *datastreamer.StreamClient
}

type sequenceData struct {
	batchClosed bool
	batch       *types.Sequence
	batchRaw    *state.BatchRawV2
}

type ethTxData struct {
	status    ethtxmanager.MonitoredTxStatus
	fromBatch uint64
	toBatch   uint64
}

// New inits sequence sender
func New(cfg Config, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog) (*SequenceSender, error) {
	s := SequenceSender{
		cfg:               cfg,
		state:             state,
		etherman:          etherman,
		ethTxManager:      manager,
		eventLog:          eventLog,
		ethTransactions:   make(map[common.Hash]*ethTxData),
		sequenceData:      make(map[uint64]*sequenceData),
		validStream:       false,
		latestStreamBatch: 0,
	}

	return &s, nil
}

// Start starts the sequence sender
func (s *SequenceSender) Start(ctx context.Context) {
	var err error

	// Create datastream client
	s.streamClient, err = datastreamer.NewClient(s.cfg.StreamClient.Server, state.StreamTypeSequencer)
	if err != nil {
		log.Fatalf("[SeqSender] failed to create stream client, error: %v", err)
	} else {
		log.Debugf("[SeqSender] new stream client")
	}

	// Set func to handle the streaming
	s.streamClient.SetProcessEntryFunc(s.handleReceivedDataStream)

	// Start datastream client
	err = s.streamClient.Start()
	if err != nil {
		log.Fatalf("[SeqSender] failed to start stream client, error: %v", err)
	}

	// Get latest virtual state batch from L1
	err = s.updateLatestVirtualBatch()
	if err != nil {
		log.Fatalf("[SeqSender] error getting latest sequenced batch, error: %v", err)
	}

	// Set starting point of the streaming
	s.fromStreamBatch = s.latestVirtualBatch
	bookmark := []byte{state.BookMarkTypeBatch}
	bookmark = binary.LittleEndian.AppendUint64(bookmark, s.fromStreamBatch)
	s.streamClient.FromBookmark = bookmark
	log.Debugf("[SeqSender] stream client from bookmark %v", bookmark)

	// Current batch to sequence
	s.wipBatch = s.latestVirtualBatch + 1
	s.latestSentToL1Batch = s.latestVirtualBatch

	// Start receiving the streaming
	err = s.streamClient.ExecCommand(datastreamer.CmdStart)
	if err != nil {
		log.Fatalf("[SeqSender] failed to connect to the streaming")
	}

	// Sequence
	ticker := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	for {
		s.tryToSendSequence(ctx, ticker)
	}
}

func (s *SequenceSender) tryToSendSequence(ctx context.Context, ticker *time.Ticker) {
	// process monitored sequences before starting a next cycle
	// retry := false
	// s.ethTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
	// 	if result.Status == ethtxmanager.MonitoredTxStatusFailed {
	// 		retry = true
	// 		mTxResultLogger := ethtxmanager.CreateMonitoredTxResultLogger(ethTxManagerOwner, result)
	// 		mTxResultLogger.Error("failed to send sequence, TODO: review this fatal and define what to do in this case")
	// 	}
	// }, nil)

	// if retry {
	// 	return
	// }

	// Check if should send sequence to L1
	log.Debugf("[SeqSender] getting sequences to send")
	sequences, err := s.getSequencesToSend(ctx)
	if err != nil || len(sequences) == 0 {
		if err != nil {
			log.Errorf("[SeqSender] error getting sequences: %v", err)
		} else {
			log.Debugf("[SeqSender] waiting for sequences to be worth sending to L1")
		}
		waitTick(ctx, ticker)
		return
	}

	// Send sequences to L1
	sequenceCount := len(sequences)
	firstSequence := sequences[0]
	lastSequence := sequences[sequenceCount-1]
	log.Infof("[SeqSender] sending sequences to L1. From batch %d to batch %d", firstSequence.BatchNumber, lastSequence.BatchNumber)
	metrics.SequencesSentToL1(float64(sequenceCount))

	// Add sequence to be monitored
	to, data, err := s.etherman.BuildSequenceBatchesTxData(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
	if err != nil {
		log.Error("[SeqSender] error estimating new sequenceBatches to add to eth tx manager: ", err)
		return
	}
	monitoredTxID := fmt.Sprintf(monitoredIDFormat, firstSequence.BatchNumber, lastSequence.BatchNumber)
	err = s.ethTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to, nil, data, s.cfg.GasOffset, nil)
	if err != nil {
		mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to)
		mTxLogger.Errorf("error to add sequences tx to eth tx manager: ", err)
		return
	}

	// Add new eth tx
	txHash := common.Hash{} // TODO: from response to call ethTxManager.Add
	txData := ethTxData{
		status:    ethtxmanager.MonitoredTxStatusSent,
		fromBatch: firstSequence.BatchNumber,
		toBatch:   lastSequence.BatchNumber,
	}
	s.ethTransactions[txHash] = &txData
	s.latestSentToL1Batch = lastSequence.BatchNumber

	s.printEthTxs()
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *SequenceSender) getSequencesToSend(ctx context.Context) ([]types.Sequence, error) {
	// Update latest virtual batch
	err := s.updateLatestVirtualBatch()
	if err != nil {
		return nil, err
	}

	// Add sequences until too big for a single L1 tx or last batch is reached
	sequences := []types.Sequence{}
	for i := 0; i < len(s.sequenceList); i++ {
		batchNumber := s.sequenceList[i]
		if batchNumber <= s.latestVirtualBatch || batchNumber <= s.latestSentToL1Batch {
			continue
		}

		// Check if the next batch belongs to a new forkid, in this case we need to stop sequencing as we need to
		// wait the upgrade of forkid is completed and s.cfg.NumBatchForkIdUpgrade is disabled (=0) again
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (batchNumber == (s.cfg.ForkUpgradeBatchNumber + 1)) {
			return nil, fmt.Errorf("[SeqSender] aborting sequencing process as we reached the batch %d where a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber+1)
		}

		// Check if batch is closed
		if !s.sequenceData[batchNumber].batchClosed {
			// Reached current wip batch
			break
		}

		// Add new sequence
		batch := *s.sequenceData[batchNumber].batch
		sequences = append(sequences, batch)

		// Check if can be send
		tx, err := s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
		log.Debugf("[SeqSender] estimated gas L1, tx size %d, result %v", tx.Size(), err)

		if err == nil && tx.Size() > s.cfg.MaxTxSizeForL1 {
			metrics.SequencesOvesizedDataError()
			log.Infof("[SeqSender] oversized Data on TX oldHash %s (txSize %d > %d)", tx.Hash(), tx.Size(), s.cfg.MaxTxSizeForL1)
			err = ErrOversizedData
		}

		if err != nil {
			log.Infof("[SeqSender] handling estimate gas send sequence error: %v", err)
			sequences, err = s.handleEstimateGasSendSequenceErr(ctx, sequences, batchNumber, err)
			if sequences != nil {
				// Handling the error gracefully, re-processing the sequence as a sanity check
				_, err = s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, s.cfg.L2Coinbase)
				return sequences, err
			}
			return sequences, err
		}

		// Check if the current batch is the last before a change to a new forkid, in this case we need to close and send the sequence to L1
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (batchNumber == (s.cfg.ForkUpgradeBatchNumber)) {
			log.Infof("[SeqSender] sequence should be sent to L1, as we have reached the batch %d from which a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber)
			return sequences, nil
		}
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("[SeqSender] no batches to be sequenced")
		return nil, nil
	}

	// lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	// if err != nil && !errors.Is(err, state.ErrNotFound) {
	// 	log.Warnf("[SeqSender] failed to get last l1 interaction time, err: %v. Sending sequences as a conservative approach", err)
	// 	return sequences, nil
	// }
	// if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
	// 	log.Info("[SeqSender] sequence should be sent to L1, because too long since didn't send anything to L1")
	// 	return sequences, nil
	// }

	return sequences, nil
	// log.Info("[SeqSender] not enough time has passed since last batch was virtualized, and the sequence could be bigger")
	// return nil, nil
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

// handleReceivedDataStream manages the events received by the streaming
func (s *SequenceSender) handleReceivedDataStream(e *datastreamer.FileEntry, c *datastreamer.StreamClient, ss *datastreamer.StreamServer) error {
	switch e.Type {
	case state.EntryTypeL2BlockStart:
		// Handle stream entry: Start L2 Block
		l2BlockStart := state.DSL2BlockStart{}
		l2BlockStart = l2BlockStart.Decode(e.Data)

		// Already virtualized
		if l2BlockStart.BatchNumber <= s.fromStreamBatch {
			if l2BlockStart.BatchNumber != s.latestStreamBatch {
				log.Infof("[SeqSender] skipped! batch already virtualized, number %d", l2BlockStart.BatchNumber)
			}
		} else {
			s.validStream = true
		}

		// Latest stream batch
		s.latestStreamBatch = l2BlockStart.BatchNumber

		if !s.validStream {
			return nil
		}

		// Manage if it is a new block or new batch
		if l2BlockStart.BatchNumber == s.wipBatch {
			// New block in the current batch
			if s.wipBatch == s.fromStreamBatch+1 {
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

		// Add L2 block
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

	if e.Number%50 == 0 {
		s.printSequences(true)
	}

	return nil
}

// addNewSequenceBatch adds a new batch to the sequence
func (s *SequenceSender) addNewSequenceBatch(l2BlockStart state.DSL2BlockStart) {
	s.mutexSequence.Lock()
	if s.sequenceData[l2BlockStart.BatchNumber] == nil {
		log.Infof("[SeqSender] ...new batch, number %d", l2BlockStart.BatchNumber)

		// Create sequence
		sequence := types.Sequence{
			GlobalExitRoot: l2BlockStart.GlobalExitRoot,
			Timestamp:      l2BlockStart.Timestamp,
			BatchNumber:    l2BlockStart.BatchNumber,
		}

		// Add to the list
		s.sequenceList = append(s.sequenceList, l2BlockStart.BatchNumber)

		// Create initial data
		batchRaw := state.BatchRawV2{}
		data := sequenceData{
			batchClosed: false,
			batch:       &sequence,
			batchRaw:    &batchRaw,
		}
		s.sequenceData[l2BlockStart.BatchNumber] = &data

		// Update wip batch
		s.wipBatch = l2BlockStart.BatchNumber
	}
	s.mutexSequence.Unlock()
}

// addNewBatchL2Block adds a new L2 block to the work in progress batch
func (s *SequenceSender) addNewBatchL2Block(l2BlockStart state.DSL2BlockStart) {
	s.mutexSequence.Lock()
	log.Infof("[SeqSender] .....new L2 block, number %d (batch %d)", l2BlockStart.L2BlockNumber, l2BlockStart.BatchNumber)

	// Current batch
	wipBatchRaw := s.sequenceData[s.wipBatch].batchRaw

	// New L2 block raw
	newBlockRaw := state.L2BlockRaw{}

	// Add L2 block
	wipBatchRaw.Blocks = append(wipBatchRaw.Blocks, newBlockRaw)

	// Get current L2 block
	blockIndex, blockRaw := s.getWipL2Block()
	if blockRaw == nil {
		log.Debugf("[SeqSender] wip block %d not found!")
		return
	}

	// Fill in data
	if blockIndex > 0 {
		blockRaw.DeltaTimestamp = uint32(l2BlockStart.Timestamp) - wipBatchRaw.Blocks[0].DeltaTimestamp
	} else {
		blockRaw.DeltaTimestamp = uint32(l2BlockStart.Timestamp)
	}
	blockRaw.IndexL1InfoTree = 0 //TODO: how to obtain this value
	s.mutexSequence.Unlock()
}

// addNewBlockTx adds a new Tx to the current L2 block
func (s *SequenceSender) addNewBlockTx(l2Tx state.DSL2Transaction) {
	s.mutexSequence.Lock()
	log.Infof("[SeqSender] ........new tx, length %d", l2Tx.EncodedLength)

	// Current L2 block
	_, blockRaw := s.getWipL2Block()

	// New Tx raw
	l2TxRaw := state.L2TxRaw{
		EfficiencyPercentage: l2Tx.EffectiveGasPricePercentage,
	}
	// TODO: how to store data in .Tx

	// Add Tx
	blockRaw.Transactions = append(blockRaw.Transactions, l2TxRaw)
	s.mutexSequence.Unlock()
}

// getWipL2Block returns index of the array and pointer to the current L2 block (helper func)
func (s *SequenceSender) getWipL2Block() (uint64, *state.L2BlockRaw) {
	// Current batch
	var wipBatchRaw *state.BatchRawV2
	if s.sequenceData[s.wipBatch] != nil {
		wipBatchRaw = s.sequenceData[s.wipBatch].batchRaw
	}

	// Current wip block
	if len(wipBatchRaw.Blocks) > 0 {
		blockIndex := uint64(len(wipBatchRaw.Blocks)) - 1
		return blockIndex, &wipBatchRaw.Blocks[blockIndex]
	} else {
		return 0, nil
	}
}

// updateLatestVirtualBatch queries the value in L1 and updates the latest virtual batch field
func (s *SequenceSender) updateLatestVirtualBatch() error {
	// Get latest virtual state batch from L1
	var err error

	s.latestVirtualBatch, err = s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("[SeqSender] error getting latest virtual batch, error: %v", err)
		return errors.New("fail to get latest virtual batch")
	} else {
		log.Debugf("[SeqSender] latest virtual batch %d", s.latestVirtualBatch)
	}
	return nil
}

// printSequences prints the current batches sequence in the memory structure
func (s *SequenceSender) printSequences(showBlock bool) {
	for i := 0; i < len(s.sequenceList); i++ {
		// Batch info
		batchNumber := s.sequenceList[i]
		seq := s.sequenceData[batchNumber]

		var raw *state.BatchRawV2
		if seq != nil {
			raw = seq.batchRaw
		} else {
			raw = &state.BatchRawV2{}
			log.Debugf("[SeqSender] // batch number %d not found in the map!", batchNumber)
		}

		// Total amount of L2 tx in the batch
		totalL2Txs := 0
		for k := 0; k < len(raw.Blocks); k++ {
			totalL2Txs += len(raw.Blocks[k].Transactions)
		}

		log.Debugf("[SeqSender] // seq %d: batch %d (closed? %t, #blocks: %d, #L2txs: %d, GER: %x...)", i, batchNumber, seq.batchClosed, len(raw.Blocks), totalL2Txs, seq.batch.GlobalExitRoot[:8])

		// Blocks info
		if showBlock {
			numBlocks := len(raw.Blocks)
			var firstBlock *state.L2BlockRaw
			var lastBlock *state.L2BlockRaw
			if numBlocks > 0 {
				firstBlock = &raw.Blocks[0]
			}
			if numBlocks > 1 {
				lastBlock = &raw.Blocks[numBlocks-1]
			}
			if firstBlock != nil {
				log.Debugf("[SeqSender] //    block first (delta-timestamp %d, #L2txs: %d)", firstBlock.DeltaTimestamp, len(firstBlock.Transactions))
			}
			if lastBlock != nil {
				log.Debugf("[SeqSender] //    block last (delta-timestamp %d, #L2txs: %d)", lastBlock.DeltaTimestamp, len(lastBlock.Transactions))
			}
		}
	}
}

// printEthTxs prints the current L1 transactions in the memory structure
func (s *SequenceSender) printEthTxs() {
	for hash, data := range s.ethTransactions {
		log.Debugf("[SeqSender] // tx hash %x... (status: %s, from: %d, to: %d) hash %x", hash[:4], data.status, data.fromBatch, data.toBatch, hash)
	}
}
