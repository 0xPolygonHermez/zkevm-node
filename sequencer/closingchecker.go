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
	// Check if there are any deposits or GER needs to be updated
	if isThereAnyDeposits, err := s.isThereAnyDeposits(ctx); err != nil || isThereAnyDeposits {
		if err != nil {
			return false
		}
		return true
	}
	// Check if it has been too long since a previous batch was virtualized
	if isBatchVirtualized, err := s.isBatchVirtualized(ctx); err != nil || isBatchVirtualized {
		if err != nil {
			return false
		}
		return true
	}
	// Check ZK counters
	zkCounters := s.calculateZkCounters()
	if zkCounters.IsZkCountersBelowZero() && len(s.sequenceInProgress.Txs) != 0 {
		log.Info("closing sequence because at least some ZK counter is bellow 0")
		return true
	}

	return false
}

func (s *Sequencer) isThereAnyDeposits(ctx context.Context) (bool, error) {
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

func (s *Sequencer) isBatchVirtualized(ctx context.Context) (bool, error) {
	isPreviousBatchVirtualized, err := s.state.IsBatchVirtualized(ctx, s.lastBatchNum-1, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return false, err
	}
	if time.Unix(s.sequenceInProgress.Timestamp, 0).Add(s.cfg.MaxTimeForBatchToBeOpen.Duration).Before(time.Now()) &&
		isPreviousBatchVirtualized && len(s.sequenceInProgress.Txs) > 0 {
		log.Info("closing batch because there are enough time to close a batch, previous batch is virtualized and batch has txs")
		return true, nil
	}
	return false, nil
}
