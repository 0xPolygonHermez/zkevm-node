package synchronizer

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// - multiples etherman to do it in parallel
// - generate blocks to be retrieved
// - retrieve blocks (parallel)
// - check last block on L1?

type blockStatusEnum int8

const (
	blockIsPending  blockStatusEnum = 0
	blockIsRunning  blockStatusEnum = 1
	blockIsInError  blockStatusEnum = 2
	blockIsFinished blockStatusEnum = 3
)

type blockRangeAlive struct {
	blockRange blockRange
	status     blockStatusEnum
}

type L1DataRetriever struct {
	mutex      sync.Mutex
	ctx        context.Context
	cancelCtx  context.CancelFunc
	workers    workers
	syncStatus syncStatus
	// If there are a running request to get the last block on L1, channelGettingLastBlockOnL1!=nil
	channelGettingLastBlockOnL1    *chan genericResponse[retrieveL1LastBlockResult]
	channelAggregatorForRollupInfo *chan genericResponse[getRollupInfoByBlockRangeResult]
	// Send the block info to this channel ordered by block number
	channel chan getRollupInfoByBlockRangeResult
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// verifyDry: test params and status without any connection of modification of objects
func (this *L1DataRetriever) verifyDry() error {
	err := this.workers.verifyDry()
	if err != nil {
		return err
	}
	err = this.syncStatus.verifyDry()
	if err != nil {
		return err
	}
	return nil
}

// verify: verify params and status with connection to the blockchain
func (this *L1DataRetriever) verify() error {
	err := this.workers.verify()
	if err != nil {
		return err
	}
	return nil
}

func NewL1Sync(ctx context.Context, ethermans []ethermanInterface,
	startingBlockNumber uint64, SyncChunkSize uint64,
	resultChannel chan getRollupInfoByBlockRangeResult) *L1DataRetriever {
	ctx, cancel := context.WithCancel(ctx)
	result := L1DataRetriever{
		ctx:        ctx,
		cancelCtx:  cancel,
		syncStatus: newSyncStatus(startingBlockNumber, SyncChunkSize),
		workers:    newWorkers(ethermans),
		channel:    resultChannel,
	}
	err := result.verifyDry()
	if err != nil {
		log.Fatal(err)
	}
	return &result

}

func (l *L1DataRetriever) Initialize() error {
	// TODO: check that all ethermans have the same chainID and get last block in L1
	err := l.verify()
	if err != nil {
		log.Fatal(err)
	}
	err = l.workers.initialize()
	if err != nil {
		log.Fatal(err)
	}
	ch, err := l.workers.asyncRetrieveLastBlock(l.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := <-ch
	if result.err != nil {
		log.Fatal(result.err)
	}
	l.onNewLastBlock(result.result.block)
	return nil
}

// OnNewLastBlock is called when a new last block is responsed by L1
func (l *L1DataRetriever) onNewLastBlock(lastBlock uint64) {
	resp := l.syncStatus.onNewLastBlockOnL1(lastBlock)
	if resp.extendedRange != nil {
		// TODO: Try to launch new workers
	}
}

// launchWork: launch new workers if possible and returns new channels created
func (l *L1DataRetriever) launchWork() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for {
		br := l.syncStatus.getNextRange()
		if br == nil {
			// No more work to do
			return
		}
		_, err := l.workers.asyncGetRollupInfoByBlockRange(l.ctx, *br)
		if err != nil {
			log.Info(err)
			return
		}
		l.syncStatus.onStartedNewWorker(*br)
	}
}

func (l *L1DataRetriever) Start() error {
	var waitDuration = time.Duration(0)
	l.launchWork()
	for l.step(&waitDuration) {
	}
	return nil
}

func (l *L1DataRetriever) renewLastBlockOnL1IfNeeded() {
	if !l.syncStatus.needToRenewLastBlockOnL1() {
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.channelGettingLastBlockOnL1 != nil {
		// There is already a request to get the last block on L1
		return
	}
	ch, err := l.workers.asyncRetrieveLastBlock(l.ctx)
	if err != nil {
		log.Warnf("Error while trying to get last block on L1: %v", err)
		l.channelGettingLastBlockOnL1 = nil
	}
	l.channelGettingLastBlockOnL1 = &ch
}

func (l *L1DataRetriever) onResponseRetrieveLastBlockOnL1(result genericResponse[retrieveL1LastBlockResult]) {
	l.channelGettingLastBlockOnL1 = nil
	if result.err != nil {
		log.Warnf("Error while trying to get last block on L1: %v", result.err)
		return
	}
	l.onNewLastBlock(result.result.block)
}

func (l *L1DataRetriever) step(waitDuration *time.Duration) bool {
	select {
	case <-l.ctx.Done():
		return false
	case <-time.After(*waitDuration):
		*waitDuration = time.Duration(time.Second)
		l.renewLastBlockOnL1IfNeeded()
	case resultRollupInfo := <-l.workers.chAggregatedRollupInfo:
		l.onResponseRollupInfo(resultRollupInfo)
		l.launchWork()
	case resultLastBlock := <-l.workers.chAggregatedLastBlock:
		l.channelGettingLastBlockOnL1 = nil
		l.onResponseRetrieveLastBlockOnL1(resultLastBlock)
	}
	return true
}

func (l *L1DataRetriever) onResponseRollupInfo(result genericResponse[getRollupInfoByBlockRangeResult]) {
	isOk := (result.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
}

func (l *L1DataRetriever) Stop() {
	l.cancelCtx()
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
