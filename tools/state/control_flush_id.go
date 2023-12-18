package main

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// FlushIDController is an interface to control the flushID and ProverID
type flushIDController interface {
	// UpdateAndCheckProverID check the incomming proverID from executor with the last one, if no match finalize synchronizer
	// if there are no previous one it keep this value as the current one
	UpdateAndCheckProverID(proverID string)
	// BlockUntilLastFlushIDIsWritten blocks until the last flushID is written in DB. It keep in a loop asking to executor
	// the flushid written, also check ProverID
	BlockUntilLastFlushIDIsWritten() error
	// SetPendingFlushIDAndCheckProverID set the pending flushID to be written in DB and check proverID
	SetPendingFlushIDAndCheckProverID(flushID uint64, proverID string, callDescription string)
}

// ClientFlushIDControl is a struct to control the flushID and ProverID, implements FlushIDController interface
type ClientFlushIDControl struct {
	state *state.State
	ctx   context.Context

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

// NewFlushIDController create a new struct ClientFlushIDControl
func NewFlushIDController(state *state.State, ctx context.Context) *ClientFlushIDControl {
	return &ClientFlushIDControl{
		state:                   state,
		ctx:                     ctx,
		proverID:                "",
		previousExecutorFlushID: 0,
	}
}

// SetPendingFlushIDAndCheckProverID set the pending flushID to be written in DB and check proverID
func (s *ClientFlushIDControl) SetPendingFlushIDAndCheckProverID(flushID uint64, proverID string, callDescription string) {
	log.Infof("new executor [%s] pending flushID: %d", callDescription, flushID)
	s.latestFlushID = flushID
	s.latestFlushIDIsFulfilled = false
	s.UpdateAndCheckProverID(proverID)
}

// UpdateAndCheckProverID check the incomming proverID from executor with the last one, if no match finalize synchronizer
// if there are no previous one it keep this value as the current one
func (s *ClientFlushIDControl) UpdateAndCheckProverID(proverID string) {
	if s.proverID == "" {
		log.Infof("Current proverID is %s", proverID)
		s.proverID = proverID
		return
	}
	if s.proverID != proverID {
		log.Fatal("restarting synchronizer because  executor have restarted (old=%s, new=%s)", s.proverID, proverID)
	}
}

// BlockUntilLastFlushIDIsWritten blocks until the last flushID is written in DB. It keep in a loop asking to executor
// the flushid written, also check ProverID
func (s *ClientFlushIDControl) BlockUntilLastFlushIDIsWritten() error {
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
		log.Infof("Synchornizer BLOCKED!: Wating for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d",
			s.latestFlushID,
			storedFlushID)
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
