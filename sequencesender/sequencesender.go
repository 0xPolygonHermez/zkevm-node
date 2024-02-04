package sequencesender

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

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
	cfg          Config
	state        stateInterface
	ethTxManager ethTxManager
	etherman     etherman
	eventLog     *event.EventLog
	privKey      *ecdsa.PrivateKey
}

// New inits sequence sender
func New(cfg Config, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog, privKey *ecdsa.PrivateKey) (*SequenceSender, error) {
	return &SequenceSender{
		cfg:          cfg,
		state:        state,
		etherman:     etherman,
		ethTxManager: manager,
		eventLog:     eventLog,
		privKey:      privKey,
	}, nil
}

// Start starts the sequence sender
func (s *SequenceSender) Start(ctx context.Context) {
	for {
		s.tryToSendSequenceX1(ctx)
	}
}

// marginTimeElapsed checks if the time between currentTime and l2BlockTimestamp is greater than timeMargin.
// If it's greater returns true, otherwise it returns false and the waitTime needed to achieve this timeMargin
func (s *SequenceSender) marginTimeElapsed(ctx context.Context, l2BlockTimestamp uint64, currentTime uint64, timeMargin int64) (bool, int64) {
	// Check the time difference between L2 block and currentTime
	var timeDiff int64
	if l2BlockTimestamp >= currentTime {
		//L2 block timestamp is above currentTime, negative timeDiff. We do in this way to avoid uint64 overflow
		timeDiff = int64(-(l2BlockTimestamp - currentTime))
	} else {
		timeDiff = int64(currentTime - l2BlockTimestamp)
	}

	// Check if the time difference is less than timeMargin (L1BlockTimestampMargin)
	if timeDiff < timeMargin {
		var waitTime int64
		if timeDiff < 0 { //L2 block timestamp is above currentTime
			waitTime = timeMargin + (-timeDiff)
		} else {
			waitTime = timeMargin - timeDiff
		}
		return false, waitTime
	} else { // timeDiff is greater than timeMargin
		return true, 0
	}
}

func (s *SequenceSender) tryToSendSequence(ctx context.Context) {
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
		log.Info("wait virtual state to be synced...")
		time.Sleep(5 * time.Second) // nolint:gomnd
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
		time.Sleep(s.cfg.WaitPeriodSendSequence.Duration)
		return
	}

	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return
	}

	// Send sequences to L1
	sequenceCount := len(sequences)
	log.Infof("sending sequences to L1. From batch %d to batch %d", lastVirtualBatchNum+1, lastVirtualBatchNum+uint64(sequenceCount))
	metrics.SequencesSentToL1(float64(sequenceCount))

	// Check if we need to wait until last L1 block timestamp is L1BlockTimestampMargin seconds above the timestamp of the last L2 block in the sequence
	// Get last batch in the sequence
	lastBatchNumInSequence := sequences[sequenceCount-1].BatchNumber

	// Get L2 blocks for the last batch
	lastBatchL2Blocks, err := s.state.GetL2BlocksByBatchNumber(ctx, lastBatchNumInSequence, nil)
	if err != nil {
		log.Errorf("failed to get L2 blocks for batch %d, err: %v", lastBatchNumInSequence, err)
		return
	}

	// Check there are L2 blocks for the last batch
	if len(lastBatchL2Blocks) == 0 {
		log.Errorf("no L2 blocks returned from the state for batch %d", lastBatchNumInSequence)
		return
	}

	// Get timestamp of the last L2 block in the sequence
	lastL2Block := lastBatchL2Blocks[len(lastBatchL2Blocks)-1]
	lastL2BlockTimestamp := uint64(lastL2Block.ReceivedAt.Unix())

	timeMargin := int64(s.cfg.L1BlockTimestampMargin.Seconds())

	// Wait until last L1 block timestamp is timeMargin (L1BlockTimestampMargin) seconds above the timestamp of the last L2 block in the sequence
	for {
		// Get header of the last L1 block
		lastL1BlockHeader, err := s.etherman.GetLatestBlockHeader(ctx)
		if err != nil {
			log.Errorf("failed to get last L1 block timestamp, err: %v", err)
			return
		}

		elapsed, waitTime := s.marginTimeElapsed(ctx, lastL2BlockTimestamp, lastL1BlockHeader.Time, timeMargin)

		if !elapsed {
			log.Infof("waiting at least %d seconds to send sequences, time difference between last L1 block %d (ts: %d) and last L2 block %d (ts: %d) in the sequence is lower than %d seconds",
				waitTime, lastL1BlockHeader.Number, lastL1BlockHeader.Time, lastL2Block.Number(), lastL2BlockTimestamp, timeMargin)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			log.Infof("continuing, time difference between last L1 block %d (ts: %d) and last L2 block %d (ts: %d) in the sequence is greater than %d seconds",
				lastL1BlockHeader.Number, lastL1BlockHeader.Time, lastL2Block.Number(), lastL2BlockTimestamp, timeMargin)
			break
		}
	}

	// Sanity check. Wait also until current time (now) is timeMargin (L1BlockTimestampMargin) seconds above the timestamp of the last L2 block in the sequence
	for {
		currentTime := uint64(time.Now().Unix())

		elapsed, waitTime := s.marginTimeElapsed(ctx, lastL2BlockTimestamp, currentTime, timeMargin)

		// Wait if the time difference is less than timeMargin (L1BlockTimestampMargin)
		if !elapsed {
			log.Infof("waiting at least %d seconds to send sequences, time difference between now (ts: %d) and last L2 block %d (ts: %d) in the sequence is lower than %d seconds",
				waitTime, currentTime, lastL2Block.Number(), lastL2BlockTimestamp, timeMargin)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			log.Infof("sending sequences now, time difference between now (ts: %d) and last L2 block %d (ts: %d) in the sequence is also greater than %d seconds",
				currentTime, lastL2Block.Number(), lastL2BlockTimestamp, timeMargin)
			break
		}
	}

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
	log.Debugf("last virtual batch number: %d", lastVirtualBatchNum)

	currentBatchNumToSequence := lastVirtualBatchNum + 1
	log.Debugf("current batch number to sequence: %d", currentBatchNumToSequence)

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

		// Add new sequence
		batch, err := s.state.GetBatchByNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			if err == state.ErrNotFound {
				break
			}
			log.Debugf("failed to get batch by number %d, err: %w", currentBatchNumToSequence, err)
			return nil, err
		}

		// Check if batch is closed
		isClosed, err := s.state.IsBatchClosed(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			log.Debugf("failed to check if batch %d is closed, err: %w", currentBatchNumToSequence, err)
			return nil, err
		}

		if !isClosed {
			// Reached current (WIP) batch
			break
		}

		seq := types.Sequence{
			BatchL2Data: batch.BatchL2Data,
			BatchNumber: batch.BatchNumber,
		}

		if batch.ForcedBatchNum != nil {
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

	// while estimating gas a new block is not created and the POE SC may return
	// an error regarding timestamp verification, this must be handled
	// if errors.Is(err, ethman.ErrTimestampMustBeInsideRange) {
	// 	// query the sc about the value of its lastTimestamp variable
	// 	lastTimestamp, err := s.etherman.GetLastBatchTimestamp()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// check POE SC lastTimestamp against sequences' one
	// 	for _, seq := range sequences {
	// 		if seq.Timestamp < int64(lastTimestamp) {
	// 			// TODO: gracefully handle this situation by creating an L2 reorg
	// 			log.Fatalf("sequence timestamp %d is < POE SC lastTimestamp %d", seq.Timestamp, lastTimestamp)
	// 		}
	// 		lastTimestamp = uint64(seq.Timestamp)
	// 	}
	// 	blockTimestamp, err := s.etherman.GetLatestBlockTimestamp(ctx)
	// 	if err != nil {
	// 		log.Error("error getting block timestamp: ", err)
	// 	}
	// 	log.Debugf("block.timestamp: %d is smaller than seq.Timestamp: %d. A new block must be mined in L1 before the gas can be estimated.", blockTimestamp, sequences[0].Timestamp)
	// 	return nil, nil
	// }

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

func (s *SequenceSender) isSynced(ctx context.Context) bool {
	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last virtual batch number, err: %v", err)
		return false
	}

	lastTrustedBatchClosed, err := s.state.GetLastClosedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last trusted batch closed, err: %v", err)
		return false
	}

	lastSCBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Warnf("failed to get from the SC last sequenced batch number, err: %v", err)
		return false
	}

	if lastVirtualBatchNum < lastSCBatchNum {
		log.Infof("waiting for the state to be synced, last virtual batch: %d, last SC sequenced batch: %d", lastVirtualBatchNum, lastSCBatchNum)
		return false
	} else if lastVirtualBatchNum > lastSCBatchNum { // Sanity check: virtual batch number cannot be greater than last batch sequenced in the SC
		s.halt(ctx, fmt.Errorf("last virtual batch %d is greater than last SC sequenced batch %d", lastVirtualBatchNum, lastSCBatchNum))
		return false
	}

	// At this point lastVirtualBatchNum = lastEthBatchNum. Check trusted batches
	if lastTrustedBatchClosed.BatchNumber >= lastVirtualBatchNum {
		return true
	} else { // Sanity check: virtual batch number cannot be greater than last trusted batch closed
		s.halt(ctx, fmt.Errorf("last virtual batch %d is greater than last trusted batch closed %d", lastVirtualBatchNum, lastTrustedBatchClosed.BatchNumber))
		return false
	}
}

// halt halts the SequenceSender
func (s *SequenceSender) halt(ctx context.Context, err error) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Sequence_Sender,
		Level:       event.Level_Critical,
		EventID:     event.EventID_FinalizerHalt,
		Description: fmt.Sprintf("SequenceSender halted due to error, error: %s", err),
	}

	eventErr := s.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing SequenceSender halt event, error: %v", eventErr)
	}

	log.Errorf("halting SequenceSender, fatal error: %v", err)
	for {
		time.Sleep(300 * time.Second) //nolint:gomnd
	}
}
