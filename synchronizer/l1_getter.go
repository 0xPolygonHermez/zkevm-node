package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// - multiples etherman to do it in parallel
// - generate blocks to be retrieved
// - retrieve blocks (parallel)
// - check last block on L1?

type L1DataRetriever struct {
	mutex      sync.Mutex
	ctx        context.Context
	cancelCtx  context.CancelFunc
	workers    workers
	syncStatus syncStatus
	// Send the block info to this channel ordered by block number
	//channel chan getRollupInfoByBlockRangeResult
	sender *SendOrdererResultsToSynchronizer
}

func (l *L1DataRetriever) toStringBrief() string {
	return fmt.Sprintf("syncStatus:%s sender:%s ", l.syncStatus.toStringBrief(), l.sender.toStringBrief())

}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// verify: test params and status without if not allowModify avoid doing connection or modification of objects
func (l *L1DataRetriever) verify(allowModify bool) error {
	err := l.workers.verify(allowModify)
	if err != nil {
		return err
	}
	err = l.syncStatus.verify(allowModify)
	if err != nil {
		return err
	}
	return nil
}

func NewL1DataRetriever(ctx context.Context, ethermans []EthermanInterface,
	startingBlockNumber uint64, SyncChunkSize uint64,
	resultChannel chan getRollupInfoByBlockRangeResult) *L1DataRetriever {
	ctx, cancel := context.WithCancel(ctx)
	result := L1DataRetriever{
		ctx:        ctx,
		cancelCtx:  cancel,
		syncStatus: newSyncStatus(startingBlockNumber, SyncChunkSize),
		workers:    newWorkers(ethermans),
		sender:     NewSendResultsToSynchronizer(resultChannel, startingBlockNumber),
	}
	err := result.verify(false)
	if err != nil {
		log.Fatal(err)
	}
	return &result

}

func (l *L1DataRetriever) Initialize() error {
	// TODO: check that all ethermans have the same chainID and get last block in L1
	err := l.verify(true)
	if err != nil {
		log.Fatal(err)
	}
	err = l.workers.initialize()
	if err != nil {
		log.Fatal(err)
	}
	result := l.retrieveInitialValueOfLastBlock()
	// We don't want to start request to L1 until calling Start()
	l.onNewLastBlock(result.result.block, false)
	return nil
}

func (l *L1DataRetriever) retrieveInitialValueOfLastBlock() genericResponse[retrieveL1LastBlockResult] {
	maxPermittedRetries := 10
	for {
		ch, err := l.workers.asyncRequestLastBlock(l.ctx)
		if err != nil {
			log.Error(err)
		}
		result := <-ch
		if result.err == nil {
			return result
		}
		log.Error(result.err)
		maxPermittedRetries--
		if maxPermittedRetries == 0 {
			log.Fatal("Cannot get last block on L1")
		}
		time.Sleep(time.Second)

	}
}

// OnNewLastBlock is called when a new last block is responsed by L1
func (l *L1DataRetriever) onNewLastBlock(lastBlock uint64, launchWork bool) {
	resp := l.syncStatus.onNewLastBlockOnL1(lastBlock)
	log.Infof("New last block on L1: %v -> %s", lastBlock, resp.toString())
	if launchWork {
		l.launchWork()
	}
}

// launchWork: launch new workers if possible and returns new channels created
// returns true if new workers were launched
func (l *L1DataRetriever) launchWork() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	launchedWorker := 0
	accDebugStr := ""
	for {
		br := l.syncStatus.getNextRange()
		if br == nil {
			// No more work to do
			accDebugStr += "[NoNextRange] "
			break
		}
		_, err := l.workers.asyncRequestRollupInfoByBlockRange(l.ctx, *br)
		if err != nil {
			accDebugStr += fmt.Sprintf(" [Error:%s] ", err.Error())
			break
		}
		launchedWorker++
		l.syncStatus.onStartedNewWorker(*br)
	}
	if launchedWorker == 0 {
		log.Infof("No workers launched because: %s", accDebugStr)
	}
	return launchedWorker
}

func (l *L1DataRetriever) Start() error {
	var waitDuration = time.Duration(0)
	l.launchWork()
	for l.step(&waitDuration) {
	}
	l.workers.waitFinishAllWorkers()
	return nil
}

func (l *L1DataRetriever) renewLastBlockOnL1IfNeeded() {
	if !l.syncStatus.needToRenewLastBlockOnL1() {
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err := l.workers.asyncRequestLastBlock(l.ctx)
	if err != nil {
		if err.Error() == errReachMaximumLiveRequestsOfThisType {
			log.Debugf("There are a request to get last block on L1 already running")
		} else {
			log.Warnf("Error while trying to get last block on L1: %v", err)
		}

	}
}

func (l *L1DataRetriever) onResponseRetrieveLastBlockOnL1(result genericResponse[retrieveL1LastBlockResult]) {
	if result.err != nil {
		log.Warnf("Error while trying to get last block on L1: %v", result.err)
		return
	}
	l.onNewLastBlock(result.result.block, true)
}

func (l *L1DataRetriever) step(waitDuration *time.Duration) bool {
	select {
	case <-l.ctx.Done():
		return false
	case <-time.After(*waitDuration):
		*waitDuration = time.Duration(time.Second)
		//log.Debugf(" ** Periodic status: %s", l.toStringBrief())
	case resultRollupInfo := <-l.workers.getResponseChannelForRollupInfo():
		l.onResponseRollupInfo(resultRollupInfo)
		l.renewLastBlockOnL1IfNeeded()
		l.launchWork()
	case resultLastBlock := <-l.workers.getResponseChannelForLastBlock():
		l.onResponseRetrieveLastBlockOnL1(resultLastBlock)
		if l.syncStatus.isNodeFullySynchronizedWithL1() {
			log.Infof("Node is fully synchronized with L1")
			return false
		}
		l.launchWork()
	}
	return true
}

func (l *L1DataRetriever) onResponseRollupInfo(result genericResponse[getRollupInfoByBlockRangeResult]) {
	isOk := (result.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
	if isOk {
		l.sender.addResultAndSendToConsumer(result.result)
	} else {
		log.Warnf("Error while trying to get rollup info by block range: %v", result.err)
	}
}

func (l *L1DataRetriever) Stop() {
	l.cancelCtx()
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
