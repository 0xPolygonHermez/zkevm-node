package synchronizer

import (
	"fmt"
	"time"
)

const (
	noLastBlock            = 0
	ttlOfLastBlockInfinity = time.Duration(0)
)

type lastBlockOnL1 struct {
	lastBlock    uint64
	TTL          time.Time
	itLastForver bool // If true, TTL is ignored
}

func newSyncLastBlockEmpty() *lastBlockOnL1 {
	return &lastBlockOnL1{
		lastBlock:    noLastBlock,
		TTL:          time.Time{},
		itLastForver: false,
	}
}

func newSyncLastBlock(lastBlock uint64, ttlDuration time.Duration) *lastBlockOnL1 {
	return &lastBlockOnL1{
		lastBlock:    lastBlock,
		TTL:          time.Now().Add(ttlDuration),
		itLastForver: ttlDuration == ttlOfLastBlockInfinity,
	}
}
func (s *lastBlockOnL1) toString() string {
	remaining := time.Until(s.TTL)
	if s.itLastForver {
		return fmt.Sprintf("[lastBlock: %v, TTL remaining:  INFINITE]", s.lastBlock)
	}
	return fmt.Sprintf("[lastBlock: %v, TTL remaining: %s]", s.lastBlock, remaining.String())
}

func (s *lastBlockOnL1) isValid() bool {
	if s.lastBlock == noLastBlock {
		return false
	}
	if s.TTL.IsZero() {
		return false
	}

	return true
}

func (s *lastBlockOnL1) getLastBlock() (uint64, error) {
	if !s.isValid() {
		return 0, fmt.Errorf("last block is not valid")
	}
	return s.lastBlock, nil
}

func (s *lastBlockOnL1) isOutdated() bool {
	if !s.isValid() {
		return true
	}
	if s.itLastForver {
		return false
	}
	now := time.Now()
	return now.After(s.TTL)
}
