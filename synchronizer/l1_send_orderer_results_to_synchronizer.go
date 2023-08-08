package synchronizer

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

type SendOrdererResultsToSynchronizer struct {
	mutex                   sync.Mutex
	channel                 chan getRollupInfoByBlockRangeResult
	lastBlockOnSynchronizer uint64
	pendingResults          []getRollupInfoByBlockRangeResult
}

func NewSendResultsToSynchronizer(ch chan getRollupInfoByBlockRangeResult, lastBlockOnSynchronizer uint64) *SendOrdererResultsToSynchronizer {
	return &SendOrdererResultsToSynchronizer{channel: ch, lastBlockOnSynchronizer: lastBlockOnSynchronizer}
}

func (s *SendOrdererResultsToSynchronizer) addResult(results *getRollupInfoByBlockRangeResult) {
	if results == nil {
		log.Fatal("Nil results, make no sense!")
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s._addPendingResult(results)
	s._sendResultIfPossible()
}

func (s *SendOrdererResultsToSynchronizer) sendResultIfPossible() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s._sendResultIfPossible()
}

// _sendResultIfPossible returns true is have send any result
func (s *SendOrdererResultsToSynchronizer) _sendResultIfPossible() bool {
	brToRemove := []blockRange{}
	send := false
	for _, result := range s.pendingResults {
		if s._matchNextBlock(&result) {
			send = true
			log.Infof("Sending results to synchronizer:", result)
			s.channel <- result
			s.lastBlockOnSynchronizer = result.blockRange.toBlock
			brToRemove = append(brToRemove, result.blockRange)
			break
		}
	}
	for _, br := range brToRemove {
		s._removeBlockRange(br)
	}
	if send {
		// Try to send more results
		s._sendResultIfPossible()
	}
	return send
}

func (s *SendOrdererResultsToSynchronizer) _removeBlockRange(toRemove blockRange) {
	for i, result := range s.pendingResults {
		if result.blockRange == toRemove {
			s.pendingResults = append(s.pendingResults[:i], s.pendingResults[i+1:]...)
			break
		}
	}
}

func (s *SendOrdererResultsToSynchronizer) _matchNextBlock(results *getRollupInfoByBlockRangeResult) bool {
	return results.blockRange.fromBlock == s.lastBlockOnSynchronizer+1
}

func (s *SendOrdererResultsToSynchronizer) _addPendingResult(results *getRollupInfoByBlockRangeResult) {
	s.pendingResults = append(s.pendingResults, *results)
}
