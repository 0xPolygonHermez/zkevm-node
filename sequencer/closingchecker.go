package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// shouldCloseSequenceInProgress checks if sequence should be closed or not
// in case it's enough blocks since last GER update, long time since last batch and sequence is profitable
func (s *Sequencer) shouldCloseSequenceInProgress(ctx context.Context) bool {
	// Check if sequence is full
	if s.sequenceInProgress.IsSequenceTooBig {
		log.Infof("current sequence should be closed because it has reached the maximum data size")
		s.sequenceInProgress.IsSequenceTooBig = false
		return true
	}
	if len(s.sequenceInProgress.Txs) >= int(s.cfg.MaxTxsPerBatch) {
		log.Infof("current sequence should be closed because it has reached the maximum capacity (%d txs)", s.cfg.MaxTxsPerBatch)
		return true
	}
	// Check if there are any deposits or GER needs to be updated
	if isThereAnyDeposits, err := s.shouldCloseDueToNewDeposits(ctx); err != nil || isThereAnyDeposits {
		return err == nil
	}
	// Check if it has been too long since a previous batch was virtualized
	if isBatchVirtualized, err := s.shouldCloseTooLongSinceLastVirtualized(ctx); err != nil || isBatchVirtualized {
		return err == nil
	}

	return false
}

// shouldCloseDueToNewDeposits return true if there has been new deposits on L1
// for more than WaitBlocksToUpdateGER, enough time has passed since the
// sequence was opened and the sequence is profitable (if profitability check
// is enabled).
func (s *Sequencer) shouldCloseDueToNewDeposits(ctx context.Context) (bool, error) {
	lastGer, gerReceivedAt, err := s.getLatestGer(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get latest global exit root, err: %w", err)
		return false, err
	}

	// get current ger and compare it with the last ger
	blockNum, mainnetExitRoot, err := s.state.GetBlockNumAndMainnetExitRootByGER(ctx, s.sequenceInProgress.GlobalExitRoot, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get mainnetExitRoot and blockNum by ger, err: %w", err)
		return false, err
	}

	if lastGer.MainnetExitRoot != mainnetExitRoot {
		latestBlockNumber, err := s.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("failed to get latest batch number from ethereum, err: %w", err)
			return false, err
		}

		sequenceInProgressTimestamp := time.Unix(s.sequenceInProgress.Timestamp, 0)

		if latestBlockNumber-blockNum > s.cfg.WaitBlocksToUpdateGER &&
			gerReceivedAt.Before(time.Now().Add(-s.cfg.ElapsedTimeToCloseBatchWithoutTxsDueToNewGER.Duration)) &&
			sequenceInProgressTimestamp.Before(time.Now().Add(-s.cfg.MinTimeToCloseBatch.Duration)) {
			log.Info("current sequence should be closed because blocks have been mined since last GER and enough time has passed")
			return true, nil
		}
	}

	return false, nil
}

// shouldCloseTooLongSinceLastVirtualized returns true if last batch virtualization happened
// more than MaxTimeForBatchToBeOpen ago and there are transactions in the current sequence
func (s *Sequencer) shouldCloseTooLongSinceLastVirtualized(ctx context.Context) (bool, error) {
	lastBatchNumber, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last batch number, err: %w", err)
		return false, err
	}
	isPreviousBatchVirtualized, err := s.state.IsBatchVirtualized(ctx, lastBatchNumber-1, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %w", err)
		return false, err
	}
	if time.Unix(s.sequenceInProgress.Timestamp, 0).Add(s.cfg.MaxTimeForBatchToBeOpen.Duration).Before(time.Now()) &&
		isPreviousBatchVirtualized && len(s.sequenceInProgress.Txs) > 0 {
		log.Info("current sequence should be closed because because there are enough time to close a batch, previous batch is virtualized and batch has txs")
		return true, nil
	}
	return false, nil
}
