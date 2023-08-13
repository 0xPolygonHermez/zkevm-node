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
	outgoingChannel         chan l1PackageData
	lastBlockOnSynchronizer uint64
	// pendingResults is a queue of results that are waiting to be sent to the consumer
	pendingResults []l1PackageData
}

func (s *filterToSendOrdererResultsToConsumer) toStringBrief() string {
	return fmt.Sprintf("lastBlockSenedToSync[%v] len(pending_results)[%d]",
		s.lastBlockOnSynchronizer, len(s.pendingResults))
}

func newFilterToSendOrdererResultsToConsumer(ch chan l1PackageData, lastBlockOnSynchronizer uint64) *filterToSendOrdererResultsToConsumer {
	return &filterToSendOrdererResultsToConsumer{outgoingChannel: ch, lastBlockOnSynchronizer: lastBlockOnSynchronizer}
}

func (s *filterToSendOrdererResultsToConsumer) addResultAndSendToConsumer(result *l1PackageData) {
	if result == nil {
		log.Error("call addResultAndSendToConsumer(result=nil), this never must happen!")
		return
	}

	log.Debugf("Received: %s", result.toStringBrief())
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if result.dataIsValid {
		if result.data.blockRange.fromBlock < s.lastBlockOnSynchronizer {
			log.Fatalf("It's not possible to receive a old block [%s] range that have been already send to synchronizer. Ignoring it.  status:[%s]",
				result.data.blockRange.toString(), s.toStringBrief())
			return
		}

		if !s._matchNextBlock(&result.data) {
			log.Debugf("The range %s is not the next block to be send, 	adding to pending results status:%s",
				result.data.blockRange.toString(), s.toStringBrief())
		}
	}
	s._addPendingResult(result)
	s._sendResultIfPossible()
}

// _sendResultIfPossible returns true is have send any result
func (s *filterToSendOrdererResultsToConsumer) _sendResultIfPossible() bool {
	indexToRemove := []int{}
	send := false
	for i := range s.pendingResults {
		result := s.pendingResults[i]
		if result.dataIsValid {
			if s._matchNextBlock(&result.data) {
				send = true
				log.Infof("Sending results [data] to consumer:%s: It could block channel [%d/%d]", result.toStringBrief(), len(s.outgoingChannel), cap(s.outgoingChannel))
				s.outgoingChannel <- result
				s._setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(result.data.blockRange)
				indexToRemove = append(indexToRemove, i)
				break
			}
		} else {
			// If it's a ctrl package only the first one could be send because it means that the previous one have been send
			if i == 0 {
				log.Infof("Sending results [no data] to consumer:%s: It could block channel [%d/%d]", result.toStringBrief(), len(s.outgoingChannel), cap(s.outgoingChannel))
				s.outgoingChannel <- result
				indexToRemove = append(indexToRemove, i)
				break
			}
		}
	}
	newPendingResults := []l1PackageData{}
	for j := range s.pendingResults {
		if slices.Contains(indexToRemove, j) {
			continue
		}
		newPendingResults = append(newPendingResults, s.pendingResults[j])
	}
	s.pendingResults = newPendingResults

	if send {
		// Try to send more results
		s._sendResultIfPossible()
	}
	return send
}

func (s *filterToSendOrdererResultsToConsumer) _setLastBlockOnSynchronizerCorrespondingLatBlockRangeSend(lastBlock blockRange) {
	newVaule := lastBlock.toBlock
	log.Debug("Moving lastBlockSend from ", s.lastBlockOnSynchronizer, " to ", newVaule)
	s.lastBlockOnSynchronizer = newVaule
}

func (s *filterToSendOrdererResultsToConsumer) _matchNextBlock(results *getRollupInfoByBlockRangeResult) bool {
	return results.blockRange.fromBlock == s.lastBlockOnSynchronizer+1
}

func (s *filterToSendOrdererResultsToConsumer) _addPendingResult(results *l1PackageData) {
	s.pendingResults = append(s.pendingResults, *results)
}
