package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// TBD. Considerations:
// - Should wait for a block to be finalized: https://www.alchemy.com/overviews/ethereum-commitment-levels https://ethereum.github.io/beacon-APIs/#/Beacon/getStateFinalityCheckpoints

type closingSignalsManager struct {
	finalizer *finalizer
	timestamp time.Time
}

func newClosingSignalsManager(finalizer *finalizer) *closingSignalsManager {
	return &closingSignalsManager{finalizer: finalizer}
}

func (c *closingSignalsManager) Start() {

	for {

		// Check L2 Reorg
		// ==============
		// Whats the condition to detect a L2 Reorg?

		// Check GER Update
		// Get latest GER from stateDB
		// If latest GER != previousGER -> send new Ger using channel

		// Check Forced Batches

		// Read new forces batches from stateDB
		// Send them using channel
		// Mark them as sended

		// Check Sending to L1 Timeout
		// How do we know when we have sent to L1 to reset the counter and don't do a timeout?

	}
}

func (c *closingSignalsManager) checkGERUpdate() {
}

func (c *closingSignalsManager) checkForcedBatches(ctx context.Context) {
	backupTimestamp := time.Now()
	forcedBatches, err := c.finalizer.dbManager.GetForcedBatchesSince(ctx, c.timestamp, nil)
	if err != nil {
		log.Errorf("error checking forced batches: %v", err)
		return
	}

	for _, forcedBatch := range forcedBatches {
		c.finalizer.closingSignalCh.ForcedBatchCh <- *forcedBatch
	}

	c.timestamp = backupTimestamp
}

func (c *closingSignalsManager) checkSendToL1Timeout() {
}
