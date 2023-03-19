package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

type closingSignalsManager struct {
	ctx                    context.Context
	dbManager              dbManagerInterface
	closingSignalCh        ClosingSignalCh
	cfg                    FinalizerCfg
	lastForcedBatchNumSent uint64
	etherman               etherman
}

func newClosingSignalsManager(ctx context.Context, dbManager dbManagerInterface, closingSignalCh ClosingSignalCh, cfg FinalizerCfg, etherman etherman) *closingSignalsManager {
	return &closingSignalsManager{ctx: ctx, dbManager: dbManager, closingSignalCh: closingSignalCh, cfg: cfg, etherman: etherman}
}

func (c *closingSignalsManager) Start() {
	go c.checkForcedBatches()
	go c.checkGERUpdate()
	go c.checkSendToL1Timeout()
}

func (c *closingSignalsManager) checkSendToL1Timeout() {
	for {
		timestamp, err := c.dbManager.GetLatestVirtualBatchTimestamp(c.ctx, nil)
		if err != nil {
			log.Errorf("error checking latest virtual batch timestamp: %v", err)
			time.Sleep(c.cfg.ClosingSignalsManagerWaitForCheckingL1Timeout.Duration)
		} else {
			limit := time.Now().Unix() - int64(c.cfg.ClosingSignalsManagerWaitForCheckingL1Timeout.Duration.Seconds())

			if timestamp.Unix() < limit {
				log.Debugf("sending to L1 timeout signal (timestamp: %v, limit: %v)", timestamp.Unix(), limit)
				c.closingSignalCh.SendingToL1TimeoutCh <- true
				time.Sleep(c.cfg.ClosingSignalsManagerWaitForCheckingL1Timeout.Duration)
			} else {
				time.Sleep(time.Duration(limit-timestamp.Unix()) * time.Second)
			}
		}
	}
}

func (c *closingSignalsManager) checkGERUpdate() {
	lastBatch, err := c.dbManager.GetLastBatch(c.ctx)
	for err != nil {
		log.Errorf("error getting last batch: %v", err)
		time.Sleep(time.Second)
		lastBatch, err = c.dbManager.GetLastBatch(c.ctx)
	}
	lastGERSent := lastBatch.GlobalExitRoot
	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForCheckingGER.Duration)

		lastL1BlockNumber, err := c.etherman.GetLatestBlockNumber(c.ctx)
		if err != nil {
			log.Errorf("error getting latest L1 block number: %v", err)
			continue
		}

		maxBlockNumber := uint64(0)
		if c.cfg.GERFinalityNumberOfBlocks <= lastL1BlockNumber {
			maxBlockNumber = lastL1BlockNumber - c.cfg.GERFinalityNumberOfBlocks
		}

		ger, _, err := c.dbManager.GetLatestGer(c.ctx, maxBlockNumber)
		if err != nil {
			log.Errorf("error checking GER update: %v", err)
			continue
		}

		if ger.GlobalExitRoot != lastGERSent {
			log.Debugf("sending GER update signal (GER: %v)", ger.GlobalExitRoot)
			c.closingSignalCh.GERCh <- ger.GlobalExitRoot
			lastGERSent = ger.GlobalExitRoot
		}
	}
}

func (c *closingSignalsManager) checkForcedBatches() {
	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForCheckingForcedBatches.Duration)

		if c.lastForcedBatchNumSent == 0 { // TODO: reset c.lastForcedBatchNumSent = 0 on L2 Reorg
			lastTrustedForcedBatchNum, err := c.dbManager.GetLastTrustedForcedBatchNumber(c.ctx, nil)
			if err != nil {
				log.Errorf("error getting last trusted forced batch number: %v", err)
				continue
			}
			if lastTrustedForcedBatchNum > 0 {
				c.lastForcedBatchNumSent = lastTrustedForcedBatchNum
			}
		}
		// Take into account L1 finality
		lastBlock, err := c.dbManager.GetLastBlock(c.ctx, nil)
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

		forcedBatches, err := c.dbManager.GetForcedBatchesSince(c.ctx, c.lastForcedBatchNumSent, maxBlockNumber, nil)
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
