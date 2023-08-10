package synchronizer

import (
	"fmt"
	"time"
)

const (
	noLastBlock            = 0
	ttlOfLastBlockInfinity = time.Duration(0)
)

type syncLastBlock struct {
	lastBlock    uint64
	TTL          time.Time
	itLastForver bool // If true, TTL is ignored
}

func newSyncLastBlockEmpty() syncLastBlock {
	return syncLastBlock{
		lastBlock:    noLastBlock,
		TTL:          time.Time{},
		itLastForver: false,
	}
}

func newSyncLastBlock(lastBlock uint64, ttlDuration time.Duration) syncLastBlock {
	return syncLastBlock{
		lastBlock:    lastBlock,
		TTL:          time.Now().Add(ttlDuration),
		itLastForver: ttlDuration == ttlOfLastBlockInfinity,
	}
}
func (s *syncLastBlock) toString() string {
	remaining := time.Until(s.TTL)
	return fmt.Sprintf("[lastBlock: %v, TTL remaining: %s]", s.lastBlock, remaining.String())
}

func (s *syncLastBlock) isValid() bool {
	if s.lastBlock == noLastBlock {
		return false
	}
	if s.TTL.IsZero() {
		return false
	}

	return true
}

func (s *syncLastBlock) getLastBlock() (uint64, error) {
	if !s.isValid() {
		return 0, fmt.Errorf("last block is not valid")
	}
	return s.lastBlock, nil
}

func (s *syncLastBlock) isOutdated() bool {
	if !s.isValid() {
		return true
	}
	if s.itLastForver {
		return false
	}
	now := time.Now()
	return now.After(s.TTL)
}
