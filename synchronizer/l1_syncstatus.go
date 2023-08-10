package synchronizer

import (
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	noLastBlock    = 0
	ttlOfLastBlock = time.Duration(600 * time.Second)
)

type syncLastBlock struct {
	lastBlock uint64
	TTL       time.Time
}

func (s *syncLastBlock) toString() string {
	remaining := s.TTL.Sub(time.Now())
	return fmt.Sprintf("[lastBlock: %v, TTL remaining: %s]", s.lastBlock, remaining.String())
}

type syncStatus struct {
	mutex                     sync.Mutex
	lastBlockStoreOnStateDB   uint64
	highestBlockRequested     uint64
	lastBlockOnL1             syncLastBlock
	amountOfBlocksInEachRange uint64
	// This ranges are being processed
	processingRanges LiveBlockRanges
	// This ranges need to be retried because the last execution was an error
	errorRanges LiveBlockRanges
}

func (s *syncStatus) toStringBrief() string {
	return fmt.Sprintf("lastBlockStoreOnStateDB: %v, lastBlockOnL1: %s, amountOfBlocksInEachRange: %d, processingRanges: %s, errorRanges: %s",
		s.lastBlockStoreOnStateDB, s.lastBlockOnL1.toString(), s.amountOfBlocksInEachRange, s.processingRanges.toStringBrief(), s.errorRanges.toStringBrief())
}

func newSyncStatus(lastBlockStoreOnStateDB uint64, amountOfBlocksInEachRange uint64) syncStatus {
	return syncStatus{
		lastBlockStoreOnStateDB:   lastBlockStoreOnStateDB,
		highestBlockRequested:     lastBlockStoreOnStateDB,
		amountOfBlocksInEachRange: amountOfBlocksInEachRange,
		lastBlockOnL1:             newSyncLastBlock(noLastBlock),
		processingRanges:          NewLiveBlockRanges(),
	}
}

func (s *syncStatus) isNodeFullySynchronizedWithL1() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	lastBlock, err := s.lastBlockOnL1.getLastBlock()
	if err != nil {
		log.Warnf(fmt.Sprintf("Can't decide if it's fully sync because Last block on L1  is no valid: %s", err))
		return false
	}
	return s._isNodeFullySynchronizedWithL1(lastBlock)
}

func (s *syncStatus) _isNodeFullySynchronizedWithL1(lastBlock uint64) bool {
	if lastBlock <= s.lastBlockStoreOnStateDB {
		log.Debug("No blocks to ask, we are synchronized with L1!")
		return true
	}
	return false
}

func (s *syncStatus) getNextRange() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Check if there are any range that need to be retried
	blockRangeToRetry, err := s.errorRanges.getFirstBlockRange()
	if err == nil {
		return &blockRangeToRetry
	}

	brs := &blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: s.highestBlockRequested} //s.processingRanges.GetSuperBlockRange()
	// if brs == nil {
	// 	brs = &blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: s.lastBlockStoreOnStateDB}
	// }
	lastBlock, err := s.lastBlockOnL1.getLastBlock()
	if err != nil {
		log.Debug("Last block is no valid: ", err)
		return nil
	}
	if s._isNodeFullySynchronizedWithL1(lastBlock) {
		log.Infof("No blocks to ask, we are synchronized with L1! status:%s", s.toStringBrief())
		return nil
	}
	br := _getNextBlockRangeFrom(brs.toBlock, lastBlock, s.amountOfBlocksInEachRange)
	return br
}

func (s *syncStatus) onStartedNewWorker(br blockRange) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Try to remove from error Blocks
	s.errorRanges.removeBlockRange(br)
	err := s.processingRanges.addBlockRange(br)
	if br.toBlock > s.highestBlockRequested {
		s.highestBlockRequested = br.toBlock
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (s *syncStatus) onFinishWorker(br blockRange, sucessful bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Debugf("onFinishWorker initial_status: %s", s.toStringBrief())
	// The work have been done, remove the range from pending list
	// also move the s.lastBlockStoreOnStateDB to the end of the range if needed
	err := s.processingRanges.removeBlockRange(br)
	if err != nil {
		log.Fatal(err)
	}

	if sucessful {
		// If this range is the first in the window, we need to move the s.lastBlockStoreOnStateDB to next range
		// example:
		// 		 lbs  = 99
		// 		 pending = [100, 200], [201, 300], [301, 400]
		// 		 if process the [100,200] -> lbs = 200
		if s.lastBlockStoreOnStateDB+1 == br.fromBlock {
			newValue := br.toBlock
			log.Infof("Moving s.lastBlockStoreOnStateDB from %d to %d (diff %d)", s.lastBlockStoreOnStateDB, newValue, newValue-s.lastBlockStoreOnStateDB)
			s.lastBlockStoreOnStateDB = newValue

		}
	} else {
		log.Infof("Range %s was not sucessful, adding to errorRanges to be retried", br.toString())
		s.errorRanges.addBlockRange(br)
	}
	log.Debugf("onFinishWorker final_status: %s", s.toStringBrief())
}

func _getNextBlockRangeFrom(lastBlockInState uint64, lastBlockInL1 uint64, amountOfBlocksInEachRange uint64) *blockRange {
	fromBlock := lastBlockInState + 1
	toBlock := min(lastBlockInL1, fromBlock+amountOfBlocksInEachRange)
	return &blockRange{fromBlock: fromBlock, toBlock: toBlock}
}

func (s *syncStatus) setLastBlockOnL1(lastBlock uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBlockOnL1 = newSyncLastBlock(lastBlock)
}

func (s *syncStatus) _setLastBlockOnL1(lastBlock uint64) {
	s.lastBlockOnL1 = newSyncLastBlock(lastBlock)
}

func newSyncLastBlock(lastBlock uint64) syncLastBlock {
	return syncLastBlock{
		lastBlock: lastBlock,
		TTL:       time.Now().Add(ttlOfLastBlock),
	}
}

type onNewLastBlockResponse struct {
	// New fullRange of pending blocks
	fullRange blockRange
	// New extendedRange of pending blocks due to new last block
	extendedRange *blockRange
}

func (n *onNewLastBlockResponse) toString() string {
	res := fmt.Sprintf("fullRange: [%s]", n.fullRange.toString())
	if n.extendedRange != nil {
		res += fmt.Sprintf(" extendedRange: [%s]", n.extendedRange.toString())
	} else {
		res += " extendedRange: nil"
	}
	return res
}

func (s *syncStatus) onNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Debugf("onNewLastBlockOnL1(%v) initial_status: %s", lastBlock, s.toStringBrief())
	response := onNewLastBlockResponse{
		fullRange: blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: lastBlock},
	}
	oldLastBlock, err := s.lastBlockOnL1.getLastBlock()
	if err != nil {
		// No previous last block
		response.extendedRange = &blockRange{
			fromBlock: s.lastBlockStoreOnStateDB,
			toBlock:   lastBlock,
		}
		s._setLastBlockOnL1(lastBlock)
		return response

	}
	if lastBlock > oldLastBlock {
		response.extendedRange = &blockRange{
			fromBlock: oldLastBlock + 1,
			toBlock:   lastBlock,
		}
		s._setLastBlockOnL1(lastBlock)
		return response
	}
	if lastBlock == oldLastBlock {
		response.extendedRange = nil
		s._setLastBlockOnL1(lastBlock)
		return response
	}
	if lastBlock < oldLastBlock {
		log.Warnf("new block [%v] is less than old block [%v]!", lastBlock, oldLastBlock)
		lastBlock = oldLastBlock
		response.fullRange = blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: lastBlock}
		return response
	}
	log.Debugf("onNewLastBlockOnL1(%v) final_status: %s", lastBlock, s.toStringBrief())
	return response
}

func (s *syncStatus) needToRenewLastBlockOnL1() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1.isOutdated()
}

func (s *syncStatus) verify(allowModify bool) error {
	if s.amountOfBlocksInEachRange == 0 {
		return fmt.Errorf("SyncChunkSize must be greater than 0")
	}
	if s.lastBlockStoreOnStateDB == 0 {
		return fmt.Errorf("startingBlockNumber must be greater than 0")
	}
	return nil
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
	now := time.Now()
	if now.After(s.TTL) {
		return true
	}
	return false
}
