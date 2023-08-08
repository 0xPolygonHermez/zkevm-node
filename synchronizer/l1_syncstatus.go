package synchronizer

import (
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	noLastBlock    = 0
	ttlOfLastBlock = time.Duration(10 * time.Second)
)

type syncLastBlock struct {
	lastBlock uint64
	TTL       time.Time
}

type syncStatus struct {
	mutex                     sync.Mutex
	lastBlockStoreOnStateDB   uint64
	lastBlockOnL1             syncLastBlock
	amountOfBlocksInEachRange uint64
	// This ranges are being processed
	processedRanges LiveBlockRanges
	// This ranges need to be retried because the last execution was an error
	errorRanges LiveBlockRanges
}

func newSyncStatus(lastBlockStoreOnStateDB uint64, amountOfBlocksInEachRange uint64) syncStatus {
	return syncStatus{
		lastBlockStoreOnStateDB:   lastBlockStoreOnStateDB,
		amountOfBlocksInEachRange: amountOfBlocksInEachRange,
		lastBlockOnL1:             newSyncLastBlock(noLastBlock),
		processedRanges:           NewLiveBlockRanges(),
	}
}

func (s *syncStatus) getNextRange() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	blockRangeToRetry, err := s.errorRanges.getFirstBlockRange()
	if err == nil {
		return &blockRangeToRetry
	}

	brs := s.processedRanges.GetSuperBlockRange()
	if brs == nil {
		brs = &blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: s.lastBlockStoreOnStateDB}
	}
	lastBlock, err := s.lastBlockOnL1.getLastBlock()
	if err != nil {
		log.Debugf(fmt.Sprintf("Last block is no valid: %s", err))
		return nil
	}
	if lastBlock <= s.lastBlockStoreOnStateDB {
		log.Debug("No blocks to ask, we are synchronized with L1!")
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
	err := s.processedRanges.addBlockRange(br)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *syncStatus) onFinishWorker(br blockRange, sucessful bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// The work have been done, remove the range from pending list
	// also move the s.lastBlockStoreOnStateDB to the end of the range if needed
	err := s.processedRanges.removeBlockRange(br)
	if err != nil {
		log.Fatal(err)
	}

	if sucessful {
		if s.lastBlockStoreOnStateDB <= br.fromBlock {
			log.Infof("Moving s.lastBlockStoreOnStateDB from %d to %d", s.lastBlockStoreOnStateDB, br.fromBlock)
			s.lastBlockStoreOnStateDB = br.fromBlock
		}
	} else {
		log.Infof("Range %v was not sucessful, adding to errorRanges to be retried", br)
		s.errorRanges.addBlockRange(br)
	}
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

func (s *syncStatus) onNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	return response
}

func (s *syncStatus) needToRenewLastBlockOnL1() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1.isOutdated()
}

func (s *syncStatus) verifyDry() error {
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
