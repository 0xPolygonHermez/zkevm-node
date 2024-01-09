package common

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

// HaltSlowDownSynchronizer is a Synchronizer halter, implements syncinterfaces.Halter
// basically it slows down the synchronizer sleeping for a while
type HaltSlowDownSynchronizer struct {
	EventLog  syncinterfaces.EventLogInterface
	SleepTime time.Duration
}

// NewHaltSlowDownSynchronizer creates a new Halter that doesnt not block execution
func NewHaltSlowDownSynchronizer(sleepTime time.Duration) *HaltInfinteLoop {
	return &HaltInfinteLoop{
		SleepTime: sleepTime,
	}
}

// Halt just sleep to avoid fast respawning
func (g *HaltSlowDownSynchronizer) Halt(ctx context.Context, err error) {
	log.Errorf("halting sync (no blocking): fatal error: %s", err)
	log.Error("halting the Synchronizer just for %v", g.SleepTime)
	time.Sleep(g.SleepTime) //nolint:gomnd
}
