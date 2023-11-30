package l1_parallel_sync

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
	pendingResults []L1SyncMessage
}

func newFilterToSendOrdererResultsToConsumer(lastBlockOnSynchronizer uint64) *filterToSendOrdererResultsToConsumer {
	return &filterToSendOrdererResultsToConsumer{lastBlockOnSynchronizer: lastBlockOnSynchronizer}
}

func (s *filterToSendOrdererResultsToConsumer) ToStringBrief() string {
	return fmt.Sprintf("lastBlockSenedToSync[%v] len(pending_results)[%d]",
		s.lastBlockOnSynchronizer, len(s.pendingResults))
}

func (s *filterToSendOrdererResultsToConsumer) numItemBlockedInQueue() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.pendingResults)
}
func (s *filterToSendOrdererResultsToConsumer) Reset(lastBlockOnSynchronizer uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBlockOnSynchronizer = lastBlockOnSynchronizer
	s.pendingResults = []L1SyncMessage{}
}

func (s *filterToSendOrdererResultsToConsumer) Filter(data L1SyncMessage) []L1SyncMessage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.checkValidDataUnsafe(&data)
	s.addPendingResultUnsafe(&data)
	res := []L1SyncMessage{}
	res = s.sendResultIfPossibleUnsafe(res)
	return res
}

func (s *filterToSendOrdererResultsToConsumer) checkValidDataUnsafe(result *L1SyncMessage) {
	if result.dataIsValid {
		if result.data.blockRange.fromBlock < s.lastBlockOnSynchronizer {
			log.Warnf("It's not possible to receive a old block [%s] range that have been already send to synchronizer. Ignoring it.  status:[%s]",
				result.data.blockRange.String(), s.ToStringBrief())
			return
		}

		if !s.matchNextBlockUnsafe(&result.data) {
			log.Debugf("The range %s is not the next block to be send, 	adding to pending results status:%s",
				result.data.blockRange.String(), s.ToStringBrief())
		}
	}
}

// sendResultIfPossibleUnsafe returns true is have send any result
func (s *filterToSendOrdererResultsToConsumer) sendResultIfPossibleUnsafe(previous []L1SyncMessage) []L1SyncMessage {
	resultListPackages := previous
	indexToRemove := []int{}
	send := false
	for i := range s.pendingResults {
		result := s.pendingResults[i]
		if result.dataIsValid {
			if s.matchNextBlockUnsafe(&result.data) {
				send = true
				resultListPackages = append(resultListPackages, result)
				highestBlockNumber := result.data.getHighestBlockNumberInResponse()

				s.setLastBlockOnSynchronizerCorrespondingLatBlockRangeSendUnsafe(highestBlockNumber)
				indexToRemove = append(indexToRemove, i)
				break
			}
		} else {
			// If it's a ctrl package only the first one could be send because it means that the previous one have been send
			if i == 0 {
				resultListPackages = append(resultListPackages, result)
				indexToRemove = append(indexToRemove, i)
				send = true
				break
			}
		}
	}
	s.removeIndexFromPendingResultsUnsafe(indexToRemove)

	if send {
		// Try to send more results
		resultListPackages = s.sendResultIfPossibleUnsafe(resultListPackages)
	}
	return resultListPackages
}

func (s *filterToSendOrdererResultsToConsumer) removeIndexFromPendingResultsUnsafe(indexToRemove []int) {
	newPendingResults := []L1SyncMessage{}
	for j := range s.pendingResults {
		if slices.Contains(indexToRemove, j) {
			continue
		}
		newPendingResults = append(newPendingResults, s.pendingResults[j])
	}
	s.pendingResults = newPendingResults
}

func (s *filterToSendOrdererResultsToConsumer) setLastBlockOnSynchronizerCorrespondingLatBlockRangeSendUnsafe(highestBlockNumber uint64) {
	if highestBlockNumber == invalidBlockNumber {
		return
	}
	log.Debug("Moving lastBlockSend from ", s.lastBlockOnSynchronizer, " to ", highestBlockNumber)
	s.lastBlockOnSynchronizer = highestBlockNumber
}

func (s *filterToSendOrdererResultsToConsumer) matchNextBlockUnsafe(results *rollupInfoByBlockRangeResult) bool {
	return results.blockRange.fromBlock == s.lastBlockOnSynchronizer+1
}

func (s *filterToSendOrdererResultsToConsumer) addPendingResultUnsafe(results *L1SyncMessage) {
	s.pendingResults = append(s.pendingResults, *results)
}
