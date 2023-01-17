package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type closingSignalsManager struct {
	ctx            context.Context
	finalizer      *finalizer
	closingSignals state.ClosingSignals
}

func newClosingSignalsManager(ctx context.Context, finalizer *finalizer) *closingSignalsManager {

	closingSignals, err := finalizer.dbManager.GetClosingSignals(ctx, nil)
	if err != nil && err != pgpoolstorage.ErrNotFound {
		log.Errorf("error getting closing signals: %v", err)
		return nil
	}

	if err == pgpoolstorage.ErrNotFound {
		closingSignals = &state.ClosingSignals{
			SentForcedBatchTimestamp: time.Now(),
			SentToL1Timestamp:        time.Now(),
			LastGER:                  common.Hash{},
		}
	}

	return &closingSignalsManager{ctx: ctx, finalizer: finalizer, closingSignals: *closingSignals}
}

func (c *closingSignalsManager) Start() {
	go c.checkForcedBatches()
	go c.checkGERUpdate()
	go c.checkSendToL1Timeout()
}

func (c *closingSignalsManager) checkGERUpdate() {
	for {
		time.Sleep(c.finalizer.cfg.GERDeadlineTimeoutInSec.Duration * time.Second)

		lastL2BlockHeader, err := c.finalizer.dbManager.GetLastL2BlockHeader(c.ctx, nil)
		if err != nil {
			log.Errorf("error getting last L2 block: %v", err)
			continue
		}

		maxBlockNumber := uint64(0)
		if c.finalizer.cfg.WaitBlocksToUpdateGER <= maxBlockNumber {
			maxBlockNumber = lastL2BlockHeader.Number.Uint64() - c.finalizer.cfg.WaitBlocksToUpdateGER
		}

		ger, _, err := c.finalizer.dbManager.GetLatestGer(c.ctx, maxBlockNumber)
		if err != nil {
			log.Errorf("error checking GER update: %v", err)
			continue
		}

		if ger.GlobalExitRoot != c.closingSignals.LastGER {
			c.finalizer.closingSignalCh.GERCh <- ger.GlobalExitRoot
			c.closingSignals.LastGER = ger.GlobalExitRoot
			err = c.finalizer.dbManager.UpdateClosingSignals(c.ctx, c.closingSignals, nil)
			if err != nil {
				log.Errorf("error updating closing signals: %v", err)
			}
		}
	}
}

func (c *closingSignalsManager) checkForcedBatches() {
	for {
		time.Sleep(c.finalizer.cfg.ForcedBatchDeadlineTimeoutInSec.Duration * time.Second)

		backupTimestamp := time.Now()
		forcedBatches, err := c.finalizer.dbManager.GetForcedBatchesSince(c.ctx, c.closingSignals.SentForcedBatchTimestamp, nil)
		if err != nil {
			log.Errorf("error checking forced batches: %v", err)
			continue
		}

		for _, forcedBatch := range forcedBatches {
			c.finalizer.closingSignalCh.ForcedBatchCh <- *forcedBatch
		}

		c.closingSignals.SentForcedBatchTimestamp = backupTimestamp

		err = c.finalizer.dbManager.UpdateClosingSignals(c.ctx, c.closingSignals, nil)
		if err != nil {
			log.Errorf("error updating closing signals: %v", err)
		}
	}
}

func (c *closingSignalsManager) checkSendToL1Timeout() {
	for {
		if time.Now().Sub(c.closingSignals.SentToL1Timestamp) > c.finalizer.cfg.SendingToL1DeadlineTimeoutInSec.Duration {
			c.finalizer.closingSignalCh.SendingToL1TimeoutCh <- true

			c.closingSignals.SentToL1Timestamp = time.Now()

			err := c.finalizer.dbManager.UpdateClosingSignals(c.ctx, c.closingSignals, nil)
			if err != nil {
				log.Errorf("error updating closing signals: %v", err)
			}
		}
	}
}
