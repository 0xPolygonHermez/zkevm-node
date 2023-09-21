package synchronizer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	invalidLastBlock = 0
)

var (
	errSyncChunkSizeMustBeGreaterThanZero = errors.New("SyncChunkSize must be greater than 0")
	errStartingBlockNumberMustBeDefined   = errors.New("startingBlockNumber must be defined, call producer ResetAndStop() to set a new starting point")
)

type syncStatus struct {
	mutex                     sync.Mutex
	lastBlockStoreOnStateDB   uint64
	highestBlockRequested     uint64
	lastBlockOnL1             uint64
	amountOfBlocksInEachRange uint64
	// This ranges are being processed
	processingRanges liveBlockRanges
	// This ranges need to be retried because the last execution was an error
	errorRanges liveBlockRanges
}

func (s *syncStatus) toStringBrief() string {
	return fmt.Sprintf(" lastBlockStoreOnStateDB: %v, highestBlockRequested:%v, lastBlockOnL1: %v, amountOfBlocksInEachRange: %d, processingRanges: %s, errorRanges: %s",
		s.lastBlockStoreOnStateDB, s.highestBlockRequested, s.lastBlockOnL1, s.amountOfBlocksInEachRange, s.processingRanges.toStringBrief(), s.errorRanges.toStringBrief())
}

func (s *syncStatus) toString() string {
	brief := s.toStringBrief()
	brief += fmt.Sprintf(" processingRanges:{ %s }", s.processingRanges.String())
	brief += fmt.Sprintf(" errorRanges:{ %s }", s.errorRanges.String())
	return brief
}

// newSyncStatus create a new syncStatus object
// lastBlockStoreOnStateDB: last block stored on stateDB
// amountOfBlocksInEachRange: amount of blocks to be retrieved in each range
// lastBlockTTLDuration: TTL of the last block on L1 (it could be ttlOfLastBlockInfinity that means that is no renewed)
func newSyncStatus(lastBlockStoreOnStateDB uint64, amountOfBlocksInEachRange uint64) *syncStatus {
	return &syncStatus{
		lastBlockStoreOnStateDB:   lastBlockStoreOnStateDB,
		highestBlockRequested:     lastBlockStoreOnStateDB,
		amountOfBlocksInEachRange: amountOfBlocksInEachRange,
		lastBlockOnL1:             invalidLastBlock,
		processingRanges:          newLiveBlockRanges(),
	}
}
func (s *syncStatus) reset(lastBlockStoreOnStateDB uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBlockStoreOnStateDB = lastBlockStoreOnStateDB
	s.highestBlockRequested = lastBlockStoreOnStateDB
	s.processingRanges = newLiveBlockRanges()
	s.lastBlockOnL1 = invalidLastBlock
}

func (s *syncStatus) getLastBlockOnL1() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1
}

func (s *syncStatus) haveRequiredAllBlocksToBeSynchronized() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1 <= s.highestBlockRequested && s.errorRanges.len() == 0
}

// isNodeFullySynchronizedWithL1 returns true if the node is fully synchronized with L1
// it means that all blocks until the last block on L1 are requested (maybe not finish yet) and there are no pending errors
func (s *syncStatus) isNodeFullySynchronizedWithL1() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.lastBlockOnL1 == invalidLastBlock {
		log.Warnf("Can't decide if it's fully synced because last block on L1  is no valid: %d", s.lastBlockOnL1)
		return false
	}

	if s.lastBlockOnL1 <= s.highestBlockRequested && s.errorRanges.len() == 0 && s.processingRanges.len() == 0 {
		log.Debug("No blocks to ask, we have requested and responsed all blocks from L1!")
		return true
	}
	return false
}

func (s *syncStatus) getNextRangeOnlyRetries() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Check if there are any range that need to be retried
	blockRangeToRetry, err := s.errorRanges.getFirstBlockRange()
	if err == nil {
		return &blockRangeToRetry
	}
	return nil
}

// getNextRange: if there are pending work it returns the next block to ask for
//
//	it could be a retry from a previous error or a new range
func (s *syncStatus) getNextRange() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Check if there are any range that need to be retried
	blockRangeToRetry, err := s.errorRanges.getFirstBlockRange()
	if err == nil {
		return &blockRangeToRetry
	}

	brs := &blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: s.highestBlockRequested} //s.processingRanges.GetSuperBlockRange()

	if s.lastBlockOnL1 == invalidLastBlock {
		log.Debug("Last block is no valid: ", s.lastBlockOnL1)
		return nil
	}
	if s.lastBlockOnL1 <= s.highestBlockRequested {
		log.Debug("No blocks to ask, we have requested all blocks from L1!")
		return nil
	}

	br := getNextBlockRangeFromUnsafe(brs.toBlock, s.lastBlockOnL1, s.amountOfBlocksInEachRange)
	err = br.isValid()
	if err != nil {
		log.Error(s.toString())
		log.Fatal(err)
	}
	return br
}

func (s *syncStatus) onStartedNewWorker(br blockRange) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Try to remove from error Blocks
	err := s.errorRanges.removeBlockRange(br)
	if err == nil {
		log.Infof("Retrying ranges: %s ", br.String())
	}
	err = s.processingRanges.addBlockRange(br)
	if err != nil {
		log.Error(s.toString())
		log.Fatal(err)
	}

	if br.toBlock > s.highestBlockRequested {
		s.highestBlockRequested = br.toBlock
	}
}

func (s *syncStatus) onFinishWorker(br blockRange, successful bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Debugf("onFinishWorker(br=%s, successful=%v) initial_status: %s", br.String(), successful, s.toStringBrief())
	// The work have been done, remove the range from pending list
	// also move the s.lastBlockStoreOnStateDB to the end of the range if needed
	err := s.processingRanges.removeBlockRange(br)
	if err != nil {
		log.Infof("Unexpected finished block_range %s, ignoring it: %s", br.String(), err)
		return
	}

	if successful {
		// If this range is the first in the window, we need to move the s.lastBlockStoreOnStateDB to next range
		// example:
		// 		 lbs  = 99
		// 		 pending = [100, 200], [201, 300], [301, 400]
		// 		 if process the [100,200] -> lbs = 200
		if s.lastBlockStoreOnStateDB+1 == br.fromBlock {
			newValue := br.toBlock
			log.Debugf("Moving s.lastBlockStoreOnStateDB from %d to %d (diff %d)", s.lastBlockStoreOnStateDB, newValue, newValue-s.lastBlockStoreOnStateDB)
			s.lastBlockStoreOnStateDB = newValue
		}
	} else {
		log.Infof("Range %s was not successful, adding to errorRanges to be retried", br.String())
		err := s.errorRanges.addBlockRange(br)
		if err != nil {
			log.Error(s.toString())
			log.Fatal(err)
		}
	}
	log.Debugf("onFinishWorker final_status: %s", s.toStringBrief())
}

func getNextBlockRangeFromUnsafe(lastBlockInState uint64, lastBlockInL1 uint64, amountOfBlocksInEachRange uint64) *blockRange {
	fromBlock := lastBlockInState + 1
	toBlock := min(lastBlockInL1, fromBlock+amountOfBlocksInEachRange)
	return &blockRange{fromBlock: fromBlock, toBlock: toBlock}
}

func (s *syncStatus) setLastBlockOnL1(lastBlock uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.setLastBlockOnL1Unsafe(lastBlock)
}

func (s *syncStatus) setLastBlockOnL1Unsafe(lastBlock uint64) {
	s.lastBlockOnL1 = lastBlock
}

type onNewLastBlockResponse struct {
	// New fullRange of pending blocks
	fullRange blockRange
	// New extendedRange of pending blocks due to new last block
	extendedRange *blockRange
}

func (n *onNewLastBlockResponse) toString() string {
	res := fmt.Sprintf("fullRange: %s", n.fullRange.String())
	if n.extendedRange != nil {
		res += fmt.Sprintf(" extendedRange: %s", n.extendedRange.String())
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

	if s.lastBlockOnL1 == invalidLastBlock {
		// No previous last block
		response.extendedRange = &blockRange{
			fromBlock: s.lastBlockStoreOnStateDB,
			toBlock:   lastBlock,
		}
		s.setLastBlockOnL1Unsafe(lastBlock)
		return response
	}
	oldLastBlock := s.lastBlockOnL1
	if lastBlock > oldLastBlock {
		response.extendedRange = &blockRange{
			fromBlock: oldLastBlock + 1,
			toBlock:   lastBlock,
		}
		s.setLastBlockOnL1Unsafe(lastBlock)
		return response
	}
	if lastBlock == oldLastBlock {
		response.extendedRange = nil
		s.setLastBlockOnL1Unsafe(lastBlock)
		return response
	}
	if lastBlock < oldLastBlock {
		log.Warnf("new block [%d] is less than old block [%d]!", lastBlock, oldLastBlock)
		lastBlock = oldLastBlock
		response.fullRange = blockRange{fromBlock: s.lastBlockStoreOnStateDB, toBlock: lastBlock}
		return response
	}
	log.Debugf("onNewLastBlockOnL1(%d) final_status: %s", lastBlock, s.toStringBrief())
	return response
}

func (s *syncStatus) isSetLastBlockOnL1Value() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1 == invalidLastBlock
}

func (s *syncStatus) verify() error {
	if s.amountOfBlocksInEachRange == 0 {
		return errSyncChunkSizeMustBeGreaterThanZero
	}
	if s.lastBlockStoreOnStateDB == invalidBlockNumber {
		return errStartingBlockNumberMustBeDefined
	}
	return nil
}
