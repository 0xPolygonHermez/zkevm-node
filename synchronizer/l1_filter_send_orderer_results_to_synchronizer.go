// Impelements

package synchronizer

import (
	"fmt"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"golang.org/x/exp/slices"
)

type filterToSendOrdererResultsToConsumer struct {
	mutex                   sync.Mutex
	lastBlockOnSynchronizer uint64
	// pendingResults is a queue of results that are waiting to be sent to the consumer
	pendingResults []l1SyncMessage
}

func newFilterToSendOrdererResultsToConsumer(lastBlockOnSynchronizer uint64) *filterToSendOrdererResultsToConsumer {
	return &filterToSendOrdererResultsToConsumer{lastBlockOnSynchronizer: lastBlockOnSynchronizer}
}

func (s *filterToSendOrdererResultsToConsumer) toStringBrief() string {
	return fmt.Sprintf("lastBlockSenedToSync[%v] len(pending_results)[%d]",
		s.lastBlockOnSynchronizer, len(s.pendingResults))
}

func (s *filterToSendOrdererResultsToConsumer) reset(lastBlockOnSynchronizer uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBlockOnSynchronizer = lastBlockOnSynchronizer
	s.pendingResults = []l1SyncMessage{}
}

func (s *filterToSendOrdererResultsToConsumer) filter(data l1SyncMessage) []l1SyncMessage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s._checkValidData(&data)
	s._addPendingResult(&data)
	res := []l1SyncMessage{}
	res = s._sendResultIfPossible(res)
	return res
}

func (s *filterToSendOrdererResultsToConsumer) _checkValidData(result *l1SyncMessage) {
	if result.dataIsValid {
		if result.data.blockRange.fromBlock < s.lastBlockOnSynchronizer {
			log.Warnf("It's not possible to receive a old block [%s] range that have been already send to synchronizer. Ignoring it.  status:[%s]",
				result.data.blockRange.toString(), s.toStringBrief())
			return
		}

		if !s._matchNextBlock(&result.data) {
			log.Debugf("The range %s is not the next block to be send, 	adding to pending results status:%s",
				result.data.blockRange.toString(), s.toStringBrief())
		}
	}
}

// _sendResultIfPossible returns true is have send any result
func (s *filterToSendOrdererResultsToConsumer) _sendResultIfPossible(previous []l1SyncMessage) []l1SyncMessage {
	result_list_packages := previous
	indexToRemove := []int{}
	send := false
	for i := range s.pendingResults {
		result := s.pendingResults[i]
		if result.dataIsValid {
			if s._matchNextBlock(&result.data) {
				send = true
				result_list_packages = append(result_list_packages, result)
				s._setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(result.data.blockRange)
				indexToRemove = append(indexToRemove, i)
				break
			}
		} else {
			// If it's a ctrl package only the first one could be send because it means that the previous one have been send
			if i == 0 {
				result_list_packages = append(result_list_packages, result)
				indexToRemove = append(indexToRemove, i)
				break
			}
		}
	}
	s._removeIndexFromPendingResults(indexToRemove)

	if send {
		// Try to send more results
		result_list_packages = s._sendResultIfPossible(result_list_packages)
	}
	return result_list_packages
}

func (s *filterToSendOrdererResultsToConsumer) _removeIndexFromPendingResults(indexToRemove []int) {
	newPendingResults := []l1SyncMessage{}
	for j := range s.pendingResults {
		if slices.Contains(indexToRemove, j) {
			continue
		}
		newPendingResults = append(newPendingResults, s.pendingResults[j])
	}
	s.pendingResults = newPendingResults
}

func (s *filterToSendOrdererResultsToConsumer) _setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(lastBlock blockRange) {
	newVaule := lastBlock.toBlock
	log.Debug("Moving lastBlockSend from ", s.lastBlockOnSynchronizer, " to ", newVaule)
	s.lastBlockOnSynchronizer = newVaule
}

func (s *filterToSendOrdererResultsToConsumer) _matchNextBlock(results *responseRollupInfoByBlockRange) bool {
	return results.blockRange.fromBlock == s.lastBlockOnSynchronizer+1
}

func (s *filterToSendOrdererResultsToConsumer) _addPendingResult(results *l1SyncMessage) {
	s.pendingResults = append(s.pendingResults, *results)
}
