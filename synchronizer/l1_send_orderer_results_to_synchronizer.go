package synchronizer

import (
	"fmt"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

type sendOrdererResultsToSynchronizer struct {
	mutex                   sync.Mutex
	channel                 chan getRollupInfoByBlockRangeResult
	lastBlockOnSynchronizer uint64
	pendingResults          []getRollupInfoByBlockRangeResult
}

func (s *sendOrdererResultsToSynchronizer) toStringBrief() string {
	return fmt.Sprintf("lastBlockSenedToSync[%v] len(pending_results)[%d]",
		s.lastBlockOnSynchronizer, len(s.pendingResults))
}

func newSendResultsToSynchronizer(ch chan getRollupInfoByBlockRangeResult, lastBlockOnSynchronizer uint64) *sendOrdererResultsToSynchronizer {
	return &sendOrdererResultsToSynchronizer{channel: ch, lastBlockOnSynchronizer: lastBlockOnSynchronizer}
}

func (s *sendOrdererResultsToSynchronizer) addResultAndSendToConsumer(result *getRollupInfoByBlockRangeResult) {
	if result == nil {
		log.Fatal("Nil results, make no sense!")
		return
	}

	log.Debugf("Received: %s", result.toStringBrief())
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if result.blockRange.fromBlock < s.lastBlockOnSynchronizer {
		log.Fatalf("It's not possible to receive a old block [%s] range that have been already send to synchronizer. Ignoring it.  status:[%s]",
			result.blockRange.toString(), s.toStringBrief())
		return
	}

	if !s._matchNextBlock(result) {
		log.Debugf("The range %s is not the next block to be send, 	adding to pending results status:%s",
			result.blockRange.toString(), s.toStringBrief())
	}
	s._addPendingResult(result)
	s._sendResultIfPossible()
}

// _sendResultIfPossible returns true is have send any result
func (s *sendOrdererResultsToSynchronizer) _sendResultIfPossible() bool {
	brToRemove := []blockRange{}
	send := false
	for i := range s.pendingResults {
		result := s.pendingResults[i]
		if s._matchNextBlock(&result) {
			send = true
			log.Infof("Sending results to synchronizer:%s: It could block channel [%d/%d]", result.toStringBrief(), len(s.channel), cap(s.channel))
			s.channel <- result
			s._setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(result.blockRange)
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

func (s *sendOrdererResultsToSynchronizer) _removeBlockRange(toRemove blockRange) {
	for i, result := range s.pendingResults {
		if result.blockRange == toRemove {
			s.pendingResults = append(s.pendingResults[:i], s.pendingResults[i+1:]...)
			break
		}
	}
}

func (s *sendOrdererResultsToSynchronizer) _setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(lastBlock blockRange) {
	newVaule := lastBlock.toBlock
	log.Debug("Moving lastBlockSend from ", s.lastBlockOnSynchronizer, " to ", newVaule)
	s.lastBlockOnSynchronizer = newVaule
}

func (s *sendOrdererResultsToSynchronizer) _matchNextBlock(results *getRollupInfoByBlockRangeResult) bool {
	return results.blockRange.fromBlock == s.lastBlockOnSynchronizer+1
}

func (s *sendOrdererResultsToSynchronizer) _addPendingResult(results *getRollupInfoByBlockRangeResult) {
	s.pendingResults = append(s.pendingResults, *results)
}
