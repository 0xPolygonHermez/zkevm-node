package l1_parallel_sync

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

func (s *syncStatus) String() string {
	return fmt.Sprintf(" lastBlockStoreOnStateDB: %s, highestBlockRequested:%s, lastBlockOnL1: %s, amountOfBlocksInEachRange: %d, processingRanges: %s, errorRanges: %s",
		blockNumberToString(s.lastBlockStoreOnStateDB),
		blockNumberToString(s.highestBlockRequested),
		blockNumberToString(s.lastBlockOnL1), s.amountOfBlocksInEachRange, s.processingRanges.toStringBrief(), s.errorRanges.toStringBrief())
}

func (s *syncStatus) toString() string {
	brief := s.String()
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
func (s *syncStatus) Reset(lastBlockStoreOnStateDB uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBlockStoreOnStateDB = lastBlockStoreOnStateDB
	s.highestBlockRequested = lastBlockStoreOnStateDB
	s.processingRanges = newLiveBlockRanges()
	//s.lastBlockOnL1 = invalidLastBlock
}

func (s *syncStatus) GetLastBlockOnL1() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1
}

// All pending blocks have been requested or are currently being requested
func (s *syncStatus) HaveRequiredAllBlocksToBeSynchronized() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1 <= s.highestBlockRequested
}

// IsNodeFullySynchronizedWithL1 returns true if the node is fully synchronized with L1
// it means that all blocks until the last block on L1 are requested (maybe not finish yet) and there are no pending errors
func (s *syncStatus) IsNodeFullySynchronizedWithL1() bool {
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

func (s *syncStatus) GetNextRangeOnlyRetries() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.getNextRangeOnlyRetriesUnsafe()
}

func (s *syncStatus) getNextRangeOnlyRetriesUnsafe() *blockRange {
	// Check if there are any range that need to be retried
	blockRangeToRetry, err := s.errorRanges.getFirstBlockRange()
	if err == nil {
		if blockRangeToRetry.toBlock == latestBlockNumber {
			// If is a latestBlockNumber must be discarded
			log.Debugf("Discarding error block range: %s because it's a latestBlockNumber", blockRangeToRetry.String())
			err := s.errorRanges.removeBlockRange(blockRangeToRetry)
			if err != nil {
				log.Errorf("syncstatus: error removing an error br: %s current_status:%s err:%s", blockRangeToRetry.String(), s.String(), err.Error())
			}
			return nil
		}
		return &blockRangeToRetry
	}
	return nil
}

func (s *syncStatus) getHighestBlockRequestedUnsafe() uint64 {
	res := invalidBlockNumber
	for _, r := range s.processingRanges.ranges {
		if r.blockRange.toBlock > res {
			res = r.blockRange.toBlock
		}
	}

	for _, r := range s.errorRanges.ranges {
		if r.blockRange.toBlock > res {
			res = r.blockRange.toBlock
		}
	}

	return res
}

// GetNextRange: if there are pending work it returns the next block to ask for
//
//	it could be a retry from a previous error or a new range
func (s *syncStatus) GetNextRange() *blockRange {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Check if there are any range that need to be retried
	blockRangeToRetry := s.getNextRangeOnlyRetriesUnsafe()
	if blockRangeToRetry != nil {
		return blockRangeToRetry
	}

	if s.lastBlockOnL1 == invalidLastBlock {
		log.Debug("Last block is no valid: ", s.lastBlockOnL1)
		return nil
	}
	if s.lastBlockOnL1 <= s.highestBlockRequested {
		log.Debug("No blocks to ask, we have requested all blocks from L1!")
		return nil
	}
	highestBlockInProcess := s.getHighestBlockRequestedUnsafe()
	if highestBlockInProcess == latestBlockNumber {
		log.Debug("No blocks to ask, we have requested all blocks from L1!")
		return nil
	}
	br := getNextBlockRangeFromUnsafe(max(s.lastBlockStoreOnStateDB, s.getHighestBlockRequestedUnsafe()), s.lastBlockOnL1, s.amountOfBlocksInEachRange)
	err := br.isValid()
	if err != nil {
		log.Error(s.toString())
		log.Fatal(err)
	}
	return br
}

func (s *syncStatus) OnStartedNewWorker(br blockRange) {
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
	if br.toBlock == latestBlockNumber {
		s.highestBlockRequested = s.lastBlockOnL1
	} else if br.toBlock > s.highestBlockRequested {
		s.highestBlockRequested = br.toBlock
	}
}

// return true is a valid blockRange
func (s *syncStatus) OnFinishWorker(br blockRange, successful bool, highestBlockNumberInResponse uint64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Debugf("onFinishWorker(br=%s, successful=%v) initial_status: %s", br.String(), successful, s.String())
	// The work have been done, remove the range from pending list
	// also move the s.lastBlockStoreOnStateDB to the end of the range if needed
	err := s.processingRanges.removeBlockRange(br)
	if err != nil {
		log.Infof("Unexpected finished block_range %s, ignoring it: %s", br.String(), err)
		return false
	}

	if successful {
		// If this range is the first in the window, we need to move the s.lastBlockStoreOnStateDB to next range
		// example:
		// 		 lbs  = 99
		// 		 pending = [100, 200], [201, 300], [301, 400]
		// 		 if process the [100,200] -> lbs = 200
		if highestBlockNumberInResponse != invalidBlockNumber && highestBlockNumberInResponse > s.lastBlockStoreOnStateDB {
			newValue := highestBlockNumberInResponse
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
	log.Debugf("onFinishWorker final_status: %s", s.String())
	return true
}

func getNextBlockRangeFromUnsafe(lastBlockInState uint64, lastBlockInL1 uint64, amountOfBlocksInEachRange uint64) *blockRange {
	fromBlock := lastBlockInState + 1
	toBlock := min(lastBlockInL1, fromBlock+amountOfBlocksInEachRange)
	if toBlock == lastBlockInL1 {
		toBlock = latestBlockNumber
	}
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

func (s *syncStatus) OnNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Debugf("onNewLastBlockOnL1(%v) initial_status: %s", lastBlock, s.String())
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
	log.Debugf("onNewLastBlockOnL1(%d) final_status: %s", lastBlock, s.String())
	return response
}

func (s *syncStatus) DoesItHaveAllTheNeedDataToWork() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockOnL1 != invalidLastBlock && s.lastBlockStoreOnStateDB != invalidBlockNumber
}

func (s *syncStatus) Verify() error {
	if s.amountOfBlocksInEachRange == 0 {
		return errSyncChunkSizeMustBeGreaterThanZero
	}
	if s.lastBlockStoreOnStateDB == invalidBlockNumber {
		return errStartingBlockNumberMustBeDefined
	}
	return nil
}

// It returns if this block is beyond Finalized (so it could be reorg)
// If blockNumber == invalidBlockNumber then it uses the highestBlockRequested (the last block requested)
func (s *syncStatus) BlockNumberIsInsideUnsafeArea(blockNumber uint64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if blockNumber == invalidBlockNumber {
		blockNumber = s.highestBlockRequested
	}
	distanceInBlockToLatest := s.lastBlockOnL1 - blockNumber
	return distanceInBlockToLatest < maximumBlockDistanceFromLatestToFinalized
}

func (s *syncStatus) GetHighestBlockReceived() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastBlockStoreOnStateDB
}
