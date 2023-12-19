package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// TODO: closingSignalsManager now is only used to notifiy new forced batches to process, maybe it's better to remove this struct and add the
// check for new forced batches as a go func of the finalizer/sequencer
type closingSignalsManager struct {
	ctx                    context.Context
	state                  stateInterface
	closingSignalCh        ClosingSignalCh
	cfg                    FinalizerCfg
	lastForcedBatchNumSent uint64
	etherman               etherman
}

func newClosingSignalsManager(ctx context.Context, state stateInterface, closingSignalCh ClosingSignalCh, cfg FinalizerCfg, etherman etherman) *closingSignalsManager {
	return &closingSignalsManager{ctx: ctx, state: state, closingSignalCh: closingSignalCh, cfg: cfg, etherman: etherman}
}

func (c *closingSignalsManager) Start() {
	go c.checkForcedBatches()
}

func (c *closingSignalsManager) checkForcedBatches() {
	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForCheckingForcedBatches.Duration)

		if c.lastForcedBatchNumSent == 0 {
			lastTrustedForcedBatchNum, err := c.state.GetLastTrustedForcedBatchNumber(c.ctx, nil)
			if err != nil {
				log.Errorf("error getting last trusted forced batch number: %v", err)
				continue
			}
			if lastTrustedForcedBatchNum > 0 {
				c.lastForcedBatchNumSent = lastTrustedForcedBatchNum
			}
		}
		// Take into account L1 finality
		lastBlock, err := c.state.GetLastBlock(c.ctx, nil)
		if err != nil {
			log.Errorf("failed to get latest eth block number, err: %v", err)
			continue
		}

		blockNumber := lastBlock.BlockNumber

		maxBlockNumber := uint64(0)
		finalityNumberOfBlocks := c.cfg.ForcedBatchesFinalityNumberOfBlocks

		if finalityNumberOfBlocks <= blockNumber {
			maxBlockNumber = blockNumber - finalityNumberOfBlocks
		}

		forcedBatches, err := c.state.GetForcedBatchesSince(c.ctx, c.lastForcedBatchNumSent, maxBlockNumber, nil)
		if err != nil {
			log.Errorf("error checking forced batches: %v", err)
			continue
		}

		for _, forcedBatch := range forcedBatches {
			log.Debugf("sending forced batch signal (forced batch number: %v)", forcedBatch.ForcedBatchNumber)
			c.closingSignalCh.ForcedBatchCh <- *forcedBatch
			c.lastForcedBatchNumSent = forcedBatch.ForcedBatchNumber
		}
	}
}
