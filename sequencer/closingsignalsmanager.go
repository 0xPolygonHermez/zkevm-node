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
}

func (c *closingSignalsManager) checkGERUpdate() {
	var lastGERSent common.Hash

	for {
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration * time.Second)

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
		time.Sleep(c.cfg.ClosingSignalsManagerWaitForL1OperationsInSec.Duration * time.Second)

		latestSentForcedBatchNumber, err := c.dbManager.GetLastTrustedForcedBatchNumber(c.ctx, nil)

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
