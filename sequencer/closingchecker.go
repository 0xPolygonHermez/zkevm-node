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
	if s.isSequenceTooBig {
		log.Infof("current sequence should be closed because it has reached the maximum data size")
		s.isSequenceTooBig = false
		return true
	}
	if len(s.sequenceInProgress.Txs) == int(maxTxsPerBatch) {
		log.Infof("current sequence should be closed because it has reached the maximum capacity (%d txs)", maxTxsPerBatch)
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

// shouldCloseDueToNewDeposits return true if there has been new deposits on L1 for more than WaitBlocksToUpdateGER
// and the sequence is profitable (if profitability check is enabled)
func (s *Sequencer) shouldCloseDueToNewDeposits(ctx context.Context) (bool, error) {
	blockNum, mainnetExitRoot, err := s.state.GetBlockNumAndMainnetExitRootByGER(ctx, s.sequenceInProgress.GlobalExitRoot, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get mainnetExitRoot and blockNum by ger, err: %v", err)
		return false, err
	}

	lastGer, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get latest global exit root, err: %v", err)
		return false, err
	}

	if lastGer != nil && lastGer.MainnetExitRoot != mainnetExitRoot {
		latestBlockNumber, err := s.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("failed to get latest batch number from ethereum, err: %v", err)
			return false, err
		}
		if latestBlockNumber-blockNum > s.cfg.WaitBlocksToUpdateGER {
			if len(s.sequenceInProgress.Txs) == 0 {
				err := s.updateGerInBatch(ctx, lastGer)
				if err != nil {
					return false, err
				}
			} else {
				isProfitable := s.isSequenceProfitable(ctx)
				if isProfitable {
					log.Infof("current sequence should be closed because blocks have been mined since last GER and tx is profitable")
					return true, nil
				}
			}
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
