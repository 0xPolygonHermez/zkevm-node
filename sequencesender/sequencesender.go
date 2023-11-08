package sequencesender

import (
	"bytes"
	"context"
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
	"github.com/ethereum/go-ethereum/common"
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
}

// New inits sequence sender
func New(cfg Config, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog) (*SequenceSender, error) {
	return &SequenceSender{
		cfg:          cfg,
		state:        state,
		etherman:     etherman,
		ethTxManager: manager,
		eventLog:     eventLog,
	}, nil
}

// Start starts the sequence sender
func (s *SequenceSender) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	for {
		s.tryToSendSequence(ctx, ticker)
	}
}

func (s *SequenceSender) tryToSendSequence(ctx context.Context, ticker *time.Ticker) {
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
	sequences, l2coinbase, err := s.getSequencesToSend(ctx)
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
	to, data, err := s.etherman.BuildSequenceBatchesTxData(s.cfg.SenderAddress, sequences, l2coinbase)
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
// return sequences, l2coibase, error
func (s *SequenceSender) getSequencesToSend(ctx context.Context) ([]types.Sequence, common.Address, error) {
	var l2coinbase common.Address
	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to get last virtual batch num, err: %w", err)
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
			return nil, common.Address{}, fmt.Errorf("aborting sequencing process as we reached the batch %d where a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber+1)
		}

		// Check if batch is closed
		isClosed, err := s.state.IsBatchClosed(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, common.Address{}, err
		}
		if !isClosed {
			// Reached current (WIP) batch
			break
		}
		// Add new sequence
		batch, err := s.state.GetBatchByNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, common.Address{}, err
		}

		seq := types.Sequence{
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      batch.Timestamp.Unix(),
			BatchL2Data:    batch.BatchL2Data,
			BatchNumber:    batch.BatchNumber,
		}

		if batch.ForcedBatchNum != nil {
			forcedBatch, err := s.state.GetForcedBatch(ctx, *batch.ForcedBatchNum, nil)
			if err != nil {
				return nil, common.Address{}, err
			}
			seq.ForcedBatchTimestamp = forcedBatch.ForcedAt.Unix()
		}
		//  All coinbase of sequences must be same
		if len(sequences) == 0 {
			l2coinbase.SetBytes(batch.Coinbase.Bytes())
		}
		if !bytes.Equal(l2coinbase.Bytes(), batch.Coinbase.Bytes()) {
			break
		}
		sequences = append(sequences, seq)
		// Check if can be send
		tx, err = s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, l2coinbase)
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
				_, err = s.etherman.EstimateGasSequenceBatches(s.cfg.SenderAddress, sequences, l2coinbase)
				return sequences, l2coinbase, err
			}
			return sequences, l2coinbase, err
		}
		// estimatedGas = tx.Gas()

		//Check if the current batch is the last before a change to a new forkid, in this case we need to close and send the sequence to L1
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber)) {
			log.Infof("sequence should be sent to L1, as we have reached the batch %d from which a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber)
			return sequences, l2coinbase, nil
		}

		// Increase batch num for next iteration
		currentBatchNumToSequence++
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("no batches to be sequenced")
		return nil, common.Address{}, nil
	}

	lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Warnf("failed to get last l1 interaction time, err: %v. Sending sequences as a conservative approach", err)
		return sequences, l2coinbase, nil
	}
	if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
		// TODO: implement check profitability
		// if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
		log.Info("sequence should be sent to L1, because too long since didn't send anything to L1")
		return sequences, l2coinbase, nil
		//}
	}

	log.Info("not enough time has passed since last batch was virtualized, and the sequence could be bigger")
	return nil, common.Address{}, nil
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
	if errors.Is(err, ethman.ErrTimestampMustBeInsideRange) {
		// query the sc about the value of its lastTimestamp variable
		lastTimestamp, err := s.etherman.GetLastBatchTimestamp()
		if err != nil {
			return nil, err
		}
		// check POE SC lastTimestamp against sequences' one
		for _, seq := range sequences {
			if seq.Timestamp < int64(lastTimestamp) {
				// TODO: gracefully handle this situation by creating an L2 reorg
				log.Fatalf("sequence timestamp %d is < POE SC lastTimestamp %d", seq.Timestamp, lastTimestamp)
			}
			lastTimestamp = uint64(seq.Timestamp)
		}
		blockTimestamp, err := s.etherman.GetLatestBlockTimestamp(ctx)
		if err != nil {
			log.Error("error getting block timestamp: ", err)
		}
		log.Debugf("block.timestamp: %d is smaller than seq.Timestamp: %d. A new block must be mined in L1 before the gas can be estimated.", blockTimestamp, sequences[0].Timestamp)
		return nil, nil
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
