package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Sequencer) tryToCreateSequence(ctx context.Context, ticker *time.Ticker) {
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

	lastBatchNumber := uint64(0)
	lastSequence, err := s.state.GetLastSequence(ctx, nil)
	if errors.Is(err, state.ErrNotFound) {
		lastBatchNumber, err = s.state.GetLastVirtualBatchNum(ctx, nil)
		if err != nil {
			log.Errorf("failed to get last virtual batch num, err: %v", err)
			return
		}
	} else if err != nil {
		log.Errorf("failed to get last sequence, err: %v", err)
		return
	} else {
		lastBatchNumber = lastSequence.BatchNumber
	}

	// Send sequences to L1
	log.Infof(
		"creating sequences from batch %d to batch %d",
		lastBatchNumber+1, lastBatchNumber+uint64(len(sequences)),
	)

	tx, err := s.etherman.SequenceBatches(ctx, sequences, 0, nil, nil)
	if err != nil {
		log.Errorf("failed to send sequence batches: %v", err)
		return
	}

	sequenceGroup := state.SequenceGroup{
		TxHash:       tx.Hash(),
		TxNonce:      tx.Nonce(),
		Status:       state.SequenceGroupStatusPending,
		CreatedAt:    time.Now(),
		FromBatchNum: sequences[0].BatchNumber,
		ToBatchNum:   sequences[len(sequences)-1].BatchNumber,
	}

	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin a tx on state, err: %v", err)
		return
	}

	for _, sequence := range sequences {
		err = s.state.CreateSequence(ctx, sequence.BatchNumber,
			sequence.GlobalExitRoot, sequence.StateRoot, sequence.LocalExitRoot,
			sequence.Timestamp, sequence.Txs, dbTx)
		if err != nil {
			log.Errorf("failed to create a sequence for batch %v, err: %v", sequence.BatchNumber, err)
			if err := dbTx.Rollback(ctx); err != nil {
				log.Errorf("failed to rollback dbTx to create sequences: %v", err)
			}
			return
		}
	}

	err = s.state.AddSequenceGroup(ctx, sequenceGroup, dbTx)
	if err != nil {
		log.Errorf("failed to create sequence group: %v", err)
		return
	}
	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx to create sequences: %v", err)
	}

	log.Infof("sequence group created for batches from %d to %d: %s", sequenceGroup.FromBatchNum, sequenceGroup.ToBatchNum, sequenceGroup.TxHash.String())
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *Sequencer) getSequencesToSend(ctx context.Context) ([]state.Sequence, error) {
	lastBatchNumber := uint64(0)
	lastSequence, err := s.state.GetLastSequence(ctx, nil)
	if errors.Is(err, state.ErrNotFound) {
		lastBatchNumber, err = s.state.GetLastVirtualBatchNum(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get last virtual batch num, err: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get last sequence, err: %v", err)
	} else {
		lastBatchNumber = lastSequence.BatchNumber
	}

	currentBatchNumToSequence := lastBatchNumber + 1
	sequences := []state.Sequence{}
	var estimatedGas uint64

	var tx *types.Transaction

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
		sequences = append(sequences, state.Sequence{
			BatchNumber:    batch.BatchNumber,
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      batch.Timestamp,
			// ForceBatchesNum: TODO,
			Txs: txs,
		})

		// Check if can be send
		tx, err = s.etherman.EstimateGasSequenceBatches(sequences)
		if err == nil && new(big.Int).SetUint64(tx.Gas()).Cmp(s.cfg.MaxSequenceSize.Int) >= 1 {
			log.Infof("oversized Data on TX hash %s (%d > %d)", tx.Hash(), tx.Gas(), s.cfg.MaxSequenceSize)
			err = core.ErrOversizedData
		}

		if err != nil {
			sequences, err = s.handleEstimateGasSendSequenceErr(ctx, sequences, currentBatchNumToSequence, err)
			if sequences != nil {
				// Handling the error gracefully, re-processing the sequence as a sanity check
				_, err = s.etherman.EstimateGasSequenceBatches(sequences)
				return sequences, err
			}
			return sequences, err
		}
		estimatedGas = tx.Gas()

		// Increase batch num for next iteration
		currentBatchNumToSequence++
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("no batches to be sequenced")
		return nil, nil
	}

	// check profitability
	if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
		log.Info("sequence is profitable and should be sent to L1")
		return sequences, nil
	}

	log.Info("sequence not profitable yet")
	return nil, nil
}

// handleEstimateGasSendSequenceErr handles an error on the estimate gas. It will return:
// nil, error: impossible to handle gracefully
// sequence, nil: handled gracefully. Potentially manipulating the sequences
// nil, nil: a situation that requires waiting
func (s *Sequencer) handleEstimateGasSendSequenceErr(
	ctx context.Context,
	sequences []state.Sequence,
	currentBatchNumToSequence uint64,
	err error,
) ([]state.Sequence, error) {
	// Insufficient allowance
	if errors.Is(err, ethman.ErrInsufficientAllowance) {
		return nil, err
	}

	if isDataForEthTxTooBig(err) {
		if len(sequences) == 1 {
			// TODO: gracefully handle this situation by creating an L2 reorg
			log.Fatalf(
				"BatchNum %d is too big to be sent to L1, even when it's the only item in the sequence: %v",
				currentBatchNumToSequence, err,
			)
		}
		// Remove the latest item and send the sequences
		log.Infof(
			"Done building sequences, selected batches to %d. Batch %d caused the L1 tx to be too big",
			currentBatchNumToSequence, currentBatchNumToSequence+1,
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
			if seq.Timestamp.Unix() < int64(lastTimestamp) {
				// TODO: gracefully handle this situation by creating an L2 reorg
				log.Fatalf("sequence timestamp %d is < POE SC lastTimestamp %d", seq.Timestamp, lastTimestamp)
			}
			lastTimestamp = uint64(seq.Timestamp.Unix())
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
		log.Fatalf(
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
		errors.Is(err, core.ErrOversizedData) ||
		errors.Is(err, ethman.ErrContentLengthTooLarge)
}
