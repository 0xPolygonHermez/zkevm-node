package common

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

// HaltSynchronizer is a Synchronizer halter, implements syncinterfaces.Halter
// basically it logs an error and keep in a infinite loop to halt the synchronizer
type HaltSynchronizer struct {
	EventLog  syncinterfaces.EventLogInterface
	SleepTime time.Duration
}

// NewHaltSynchronizer creates a new HaltSynchronizer
func NewHaltSynchronizer(eventLog syncinterfaces.EventLogInterface, sleepTime time.Duration) *HaltSynchronizer {
	return &HaltSynchronizer{
		EventLog:  eventLog,
		SleepTime: sleepTime,
	}
}

// Halt halts the Synchronizer and write a eventLog on Database
func (g *HaltSynchronizer) Halt(ctx context.Context, err error) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Synchronizer,
		Level:       event.Level_Critical,
		EventID:     event.EventID_SynchronizerHalt,
		Description: fmt.Sprintf("Synchronizer halted due to error: %s", err),
	}

	eventErr := g.EventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing Synchronizer halt event: %v", eventErr)
	}

	for {
		log.Errorf("halting sync: fatal error: %s", err)
		log.Error("halting the Synchronizer")
		time.Sleep(g.SleepTime) //nolint:gomnd
	}
}
