package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

type closingSignalsManager struct {
	ctx             context.Context
	dbManager       dbManagerInterface
	closingSignalCh ClosingSignalCh
	cfg             FinalizerCfg
}

func newClosingSignalsManager(ctx context.Context, dbManager dbManagerInterface, closingSignalCh ClosingSignalCh, cfg FinalizerCfg) *closingSignalsManager {
	return &closingSignalsManager{ctx: ctx, dbManager: dbManager, closingSignalCh: closingSignalCh, cfg: cfg}
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
			time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration)
		} else {
			limit := time.Now().Unix() - int64(c.cfg.SendingToL1DeadlineTimeoutInSec.Duration.Seconds())

			if timestamp.Unix() < limit {
				c.closingSignalCh.SendingToL1TimeoutCh <- true
				time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration)
			} else {
				time.Sleep(time.Duration(limit-timestamp.Unix()) * time.Second)
			}
		}
	}
}

func (c *closingSignalsManager) checkGERUpdate() {
	var lastGERSent common.Hash

	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration)

		ger, _, err := c.dbManager.GetLatestGer(c.ctx, c.cfg.GERFinalityNumberOfBlocks)
		if err != nil {
			log.Errorf("error checking GER update: %v", err)
			continue
		}

		if ger.GlobalExitRoot != lastGERSent {
			c.closingSignalCh.GERCh <- ger.GlobalExitRoot
			lastGERSent = ger.GlobalExitRoot
		}
	}
}

func (c *closingSignalsManager) checkForcedBatches() {
	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration)
		log.Debug(time.Now())

		latestSentForcedBatchNumber, err := c.dbManager.GetLastTrustedForcedBatchNumber(c.ctx, nil)
		if err != nil {
			log.Errorf("error getting last trusted forced batch number: %v", err)
			continue
		}

		forcedBatches, err := c.dbManager.GetForcedBatchesSince(c.ctx, latestSentForcedBatchNumber, nil)
		if err != nil {
			log.Errorf("error checking forced batches: %v", err)
			continue
		}

		for _, forcedBatch := range forcedBatches {
			c.closingSignalCh.ForcedBatchCh <- *forcedBatch
		}
	}
}
