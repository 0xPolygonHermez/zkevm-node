package synchronizer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

type FlushIDController interface {
	UpdateAndCheckProverID(proverID string)
	CheckFlushID(dbTx pgx.Tx) error
	PendingFlushID(flushID uint64, proverID string)
}

type ClientFlushIDControl struct {
	state    stateInterface
	ctx      context.Context
	eventLog *event.EventLog

	// Id of the 'process' of the executor. Each time that it starts this value changes
	// This value is obtained from the call state.GetStoredFlushID
	// It starts as an empty string and it is filled in the first call
	// later the value is checked to be the same (in function checkFlushID)
	proverID string
	// Previous value returned by state.GetStoredFlushID, is used for decide if write a log or not
	previousExecutorFlushID uint64
	latestFlushID           uint64
	// If true the lastFlushID is stored in DB and we don't need to check again
	latestFlushIDIsFulfilled bool
}

func NewFlushIDController(state stateInterface, ctx context.Context, eventLog *event.EventLog) *ClientFlushIDControl {
	return &ClientFlushIDControl{
		state:                   state,
		ctx:                     ctx,
		eventLog:                eventLog,
		proverID:                "",
		previousExecutorFlushID: 0,
	}
}

func (s *ClientFlushIDControl) PendingFlushID(flushID uint64, proverID string) {
	log.Infof("pending flushID: %d", flushID)
	s.latestFlushID = flushID
	s.latestFlushIDIsFulfilled = false
	s.UpdateAndCheckProverID(proverID)
}

func (s *ClientFlushIDControl) UpdateAndCheckProverID(proverID string) {
	if s.proverID == "" {
		log.Infof("Current proverID is %s", proverID)
		s.proverID = proverID
		return
	}
	if s.proverID != proverID {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			Source:      event.Source_Node,
			Component:   event.Component_Synchronizer,
			Level:       event.Level_Critical,
			EventID:     event.EventID_SynchronizerRestart,
			Description: fmt.Sprintf("proverID changed from %s to %s, restarting Synchonizer ", s.proverID, proverID),
		}

		err := s.eventLog.LogEvent(context.Background(), event)
		if err != nil {
			log.Errorf("error storing event payload: %v", err)
		}

		log.Fatal("restarting synchronizer because  executor have restarted (old=%s, new=%s)", s.proverID, proverID)
	}
}

func (s *ClientFlushIDControl) CheckFlushID(dbTx pgx.Tx) error {
	if s.latestFlushIDIsFulfilled {
		log.Debugf("no pending flushID, nothing to do. Last pending fulfilled flushID: %d, last executor flushId received: %d", s.latestFlushID, s.latestFlushID)
		return nil
	}
	storedFlushID, proverID, err := s.state.GetStoredFlushID(s.ctx)
	if err != nil {
		log.Error("error getting stored flushID. Error: ", err)
		return err
	}
	if (s.previousExecutorFlushID != storedFlushID) || (s.proverID != proverID) {
		log.Infof("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	} else {
		log.Debugf("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	}
	s.UpdateAndCheckProverID(proverID)
	log.Debugf("storedFlushID (executor reported): %d, latestFlushID (pending): %d", storedFlushID, s.latestFlushID)
	if storedFlushID < s.latestFlushID {
		log.Infof("Synchornized BLOCKED!: Wating for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d", s.latestFlushID, storedFlushID)
		iteration := 0
		start := time.Now()
		for storedFlushID < s.latestFlushID {
			log.Debugf("Waiting for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d iteration:%d elpased:%s",
				s.latestFlushID, storedFlushID, iteration, time.Since(start))
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
			storedFlushID, _, err = s.state.GetStoredFlushID(s.ctx)
			if err != nil {
				log.Error("error getting stored flushID. Error: ", err)
				return err
			}
			iteration++
		}
		log.Infof("Synchornizer resumed, flushID stored: %d", s.latestFlushID)
	}
	log.Infof("Pending Flushid fullfiled: %d, executor have write %d", s.latestFlushID, storedFlushID)
	s.latestFlushIDIsFulfilled = true
	s.previousExecutorFlushID = storedFlushID
	return nil
}
