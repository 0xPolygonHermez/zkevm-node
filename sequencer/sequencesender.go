package sequencer

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core"
)

func (s *Sequencer) tryToSendSequence(ctx context.Context, ticker *time.Ticker) {
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

	// Send sequences to L1
	log.Infof(
		"sending sequences to L1. From batch %d to batch %d",
		s.lastBatchNumSentToL1+1, s.lastBatchNumSentToL1+uint64(len(sequences)),
	)
	s.txManager.SequenceBatches(sequences)
	s.lastBatchNumSentToL1 += uint64(len(sequences))
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *Sequencer) getSequencesToSend(ctx context.Context) ([]types.Sequence, error) {
	currentBatchNumToSequence := s.lastBatchNumSentToL1 + 1
	sequences := []types.Sequence{}
	var estimatedGas uint64

	// Add sequences until too big for a single L1 tx or last batch is reached
	for {
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
		txs, err := s.state.GetTransactionsByBatchNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			return nil, err
		}
		sequences = append(sequences, types.Sequence{
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      batch.Timestamp.Unix(),
			// ForceBatchesNum: TODO,
			Txs: txs,
		})

		// Check if can be send
		estimatedGas, err = s.etherman.EstimateGasSequenceBatches(sequences)
		if err != nil {
			sequences, err = s.handleEstimateGasSendSequenceErr(ctx, sequences, currentBatchNumToSequence, err)
			if sequences != nil {
				// Handling the error gracefully, re-processing the sequence as a sanity check
				_, err = s.etherman.EstimateGasSequenceBatches(sequences)
				return sequences, err
			}
			return sequences, err
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
		// check profitability
		if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
			log.Info("sequence should be sent to L1, because too long since didn't send anything to L1")
			return sequences, nil
		}
	}

	log.Info("not enougth time has passed since last batch was virtualized, and the sequence could be bigger")
	return nil, nil
}

// handleEstimateGasSendSequenceErr handles an error on the estimate gas. It will return:
// nil, error: impossible to handle gracefully
// sequence, nil: handled gracefully. Potentially manipulating the sequences
// nil, nil: a situation that requires waiting
func (s *Sequencer) handleEstimateGasSendSequenceErr(
	ctx context.Context,
	sequences []types.Sequence,
	currentBatchNumToSequence uint64,
	err error,
) ([]types.Sequence, error) {
	// Insufficient allowance
	if strings.Contains(err.Error(), errInsufficientAllowance) {
		return nil, err
	}

	// Data to big for a single ethereum transfer
	if isDataForEthTxTooBig(err) {
		if len(sequences) == 1 {
			// TODO: gracefully handle this situation by crating an L2 reorg
			log.Fatalf(
				"BatchNum %d is too big to be sent to L1, even when it's the only item in the sequence: %v",
				currentBatchNumToSequence, err,
			)
		}
		// Remove latest item and send the sequences
		log.Infof(
			"Done building sequences, selected batches from %d to %d. Batch %d caused the L1 tx to be too big",
			s.lastBatchNumSentToL1+1, currentBatchNumToSequence, currentBatchNumToSequence+1,
		)
		sequences = sequences[:len(sequences)-1]
		return sequences, nil
	}

	// while estimating gas a new block is not created and the POE SC may return
	// an error regarding timestamp verification, this must be handled
	if strings.Contains(err.Error(), errTimestampMustBeInsideRange) {
		// query the sc about the value of its lastTimestamp variable
		lastTimestamp, err := s.etherman.GetLastBatchTimestamp()
		if err != nil {
			return nil, err
		}
		// check POE SC lastTimestamp against sequences' one
		for _, seq := range sequences {
			if seq.Timestamp < int64(lastTimestamp) {
				// TODO: gracefully handle this situation by crating an L2 reorg
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
		// TODO: gracefully handle this situation by crating an L2 reorg
		log.Fatalf(
			"Error when estimating gas for BatchNum %d (alone in the sequences): %v",
			currentBatchNumToSequence, err,
		)
	}
	// Remove latest item and send the sequences
	log.Infof(
		"Done building sequences, selected batches from %d to %d. Batch %d excluded due to unknown error: %v",
		s.lastBatchNumSentToL1+1, currentBatchNumToSequence, currentBatchNumToSequence+1, err,
	)
	sequences = sequences[:len(sequences)-1]
	return sequences, nil
}

func isDataForEthTxTooBig(err error) bool {
	return strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		errors.Is(err, core.ErrOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge)
}
