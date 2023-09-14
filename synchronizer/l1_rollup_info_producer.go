// package synchronizer
// Implements the logic to retrieve data from L1 and send it to the synchronizer
//   - multiples etherman to do it in parallel
//   - generate blocks to be retrieved
//   - retrieve blocks (parallel)
//   - when reach the update state:
// 		- send a update to channel and  keep retrieving last block to ask for new rollup info
//
//
// TODO:
//   - Check all log.fatals to remove it or add status before the panic

package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"golang.org/x/exp/constraints"
)

const (
	minTTLOfLastBlock                          = time.Second
	timeOutMainLoop                            = time.Minute * 5
	timeForShowUpStatisticsLog                 = time.Second * 60
	conversionFactorPercentage                 = 100
	maxRetriesForRequestnitialValueOfLastBlock = 10
	timeRequestInitialValueOfLastBlock         = time.Second * 5
)

type filter interface {
	ToStringBrief() string
	Filter(data l1SyncMessage) []l1SyncMessage
	Reset(lastBlockOnSynchronizer uint64)
}

type syncStatusInterface interface {
	verify() error
	reset(lastBlockStoreOnStateDB uint64)
	toStringBrief() string
	getNextRange() *blockRange
	isNodeFullySynchronizedWithL1() bool
	haveRequiredAllBlocksToBeSynchronized() bool
	isSetLastBlockOnL1Value() bool
	getLastBlockOnL1() uint64

	onStartedNewWorker(br blockRange)
	onFinishWorker(br blockRange, successful bool)
	onNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse
}

type workersInterface interface {
	// initialize object
	initialize() error
	// finalize object
	stop()
	// waits until all workers have finish the current task
	waitFinishAllWorkers()
	asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan responseRollupInfoByBlockRange, error)
	requestLastBlockWithRetries(ctx context.Context, timeout time.Duration, maxPermittedRetries int) responseL1LastBlock
	getResponseChannelForRollupInfo() chan responseRollupInfoByBlockRange
	toString() string
}

type producerStatusEnum int8

const (
	producerIdle         producerStatusEnum = 0
	producerWorking      producerStatusEnum = 1
	producerSynchronized producerStatusEnum = 2
)

func (s producerStatusEnum) String() string {
	return [...]string{"idle", "working", "synchronized"}[s]
}

type configProducer struct {
	syncChunkSize      uint64
	ttlOfLastBlockOnL1 time.Duration
}

func (cfg *configProducer) normalize() {
	if cfg.syncChunkSize == 0 {
		log.Fatalf("l1RollupInfoProducer: SyncChunkSize must be greater than 0")
	}
	if cfg.ttlOfLastBlockOnL1 < minTTLOfLastBlock {
		log.Warnf("l1RollupInfoProducer: ttlOfLastBlockOnL1 is too low (%s) so setting to %s", cfg.ttlOfLastBlockOnL1, minTTLOfLastBlock)
		cfg.ttlOfLastBlockOnL1 = minTTLOfLastBlock
	}
}

type l1RollupInfoProducer struct {
	mutex             sync.Mutex
	ctxParent         context.Context
	ctx               context.Context
	cancelCtx         context.CancelFunc
	workers           workersInterface
	syncStatus        syncStatusInterface
	outgoingChannel   chan l1SyncMessage
	timeLastBLockOnL1 time.Time
	status            producerStatusEnum
	// filter is an object that sort l1DataMessage to be send ordered by block number
	filterToSendOrdererResultsToConsumer filter
	statistics                           l1RollupInfoProducerStatistics
	cfg                                  configProducer
}

// l1DataRetrieverStatistics : create an instance of l1RollupInfoProducer
func newL1DataRetriever(ctx context.Context, cfg configProducer, ethermans []EthermanInterface,
	startingBlockNumber uint64, outgoingChannel chan l1SyncMessage) *l1RollupInfoProducer {
	if cap(outgoingChannel) < len(ethermans) {
		log.Warnf("l1RollupInfoProducer: outgoingChannel must have a capacity (%d) of at least equal to number of ether clients (%d)", cap(outgoingChannel), len(ethermans))
	}
	newCtx, cancel := context.WithCancel(ctx)
	cfg.normalize()

	result := l1RollupInfoProducer{
		ctxParent:                            ctx,
		ctx:                                  newCtx,
		cancelCtx:                            cancel,
		syncStatus:                           newSyncStatus(startingBlockNumber, cfg.syncChunkSize),
		workers:                              newWorkers(ctx, ethermans),
		filterToSendOrdererResultsToConsumer: newFilterToSendOrdererResultsToConsumer(startingBlockNumber),
		outgoingChannel:                      outgoingChannel,
		statistics:                           newRollupInfoProducerStatistics(startingBlockNumber),
		status:                               producerIdle,
		cfg:                                  cfg,
	}
	return &result
}

// TDOO: There is no min/max function in golang??
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func (l *l1RollupInfoProducer) reset(startingBlockNumber uint64) {
	log.Infof("producer: Reset L1 sync process to blockNumber %d", startingBlockNumber)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	log.Debugf("producer: Reset(%d): context cancel", startingBlockNumber)
	l.cancelCtx()
	l.ctx, l.cancelCtx = context.WithCancel(l.ctxParent)
	log.Debugf("producer: Reset(%d): syncStatus.reset", startingBlockNumber)
	l.syncStatus.reset(startingBlockNumber)
	l.statistics.reset(startingBlockNumber)
	log.Debugf("producer: Reset(%d): stop workers (%s)", startingBlockNumber, l.workers.toString())
	l.workers.stop()
	// Empty pending rollupinfos
	log.Debugf("producer: Reset(%d): emptyChannel", startingBlockNumber)
	l.emptyChannel()
	log.Debugf("producer: Reset(%d): reset Filter", startingBlockNumber)
	l.filterToSendOrdererResultsToConsumer.Reset(startingBlockNumber)
	l.status = producerIdle
	log.Debugf("producer: Reset(%d): reset done!", startingBlockNumber)
}

func (l *l1RollupInfoProducer) emptyChannel() {
	for len(l.outgoingChannel) > 0 {
		<-l.outgoingChannel
	}
}

// verify: test params and status without if not allowModify avoid doing connection or modification of objects
func (l *l1RollupInfoProducer) verify() error {
	return l.syncStatus.verify()
}

func (l *l1RollupInfoProducer) initialize() error {
	err := l.verify()
	if err != nil {
		log.Fatal(err)
	}
	err = l.workers.initialize()
	if err != nil {
		log.Fatal(err)
	}
	if l.syncStatus.isSetLastBlockOnL1Value() {
		log.Infof("producer: Need a initial value for Last Block On L1, doing the request (maxRetries:%v, timeRequest:%v)",
			maxRetriesForRequestnitialValueOfLastBlock, timeRequestInitialValueOfLastBlock)
		//result := l.retrieveInitialValueOfLastBlock(maxRetriesForRequestnitialValueOfLastBlock, timeRequestInitialValueOfLastBlock)
		result := l.workers.requestLastBlockWithRetries(l.ctx, timeRequestInitialValueOfLastBlock, maxRetriesForRequestnitialValueOfLastBlock)
		if result.generic.err != nil {
			log.Error(result.generic.err)
			return result.generic.err
		}
		l.onNewLastBlock(result.result.block, false)
	}

	return nil
}

// If startingBlockNumber is invalidBlockNumber it will use the last block on stateDB
func (l *l1RollupInfoProducer) start(startingBlockNumber uint64) error {
	if startingBlockNumber != invalidBlockNumber {
		log.Infof("producer: starting L1 sync from %v, previous status:%s", startingBlockNumber, l.syncStatus.toStringBrief())
		l.reset(startingBlockNumber)
	} else {
		log.Infof("producer: starting L1 sync with no changed status:%s (startingBlock:%v)", l.syncStatus.toStringBrief(), l.syncStatus.getLastBlockOnL1())
	}
	err := l.initialize()
	if err != nil {
		return err
	}
	var waitDuration = time.Duration(0)
	for l.step(&waitDuration) {
	}
	l.workers.waitFinishAllWorkers()
	return nil
}

func (l *l1RollupInfoProducer) step(waitDuration *time.Duration) bool {
	previousStatus := l.status
	res := l.stepInner(waitDuration)
	newStatus := l.status
	if previousStatus != newStatus {
		log.Infof("producer: Status changed from [%s] to [%s]", previousStatus.String(), newStatus.String())
		if newStatus == producerSynchronized {
			log.Infof("producer: send a message to consumer to indicate that we are synchronized")
			l.sendPackages([]l1SyncMessage{*newL1SyncMessageControl(eventProducerIsFullySynced)})
		}
	}
	return res
}

func (l *l1RollupInfoProducer) stepInner(waitDuration *time.Duration) bool {
	select {
	case <-l.ctx.Done():
		log.Debugf("producer: context canceled")
		return false
	// That timeout is not need, but just in case that stop launching request
	case <-time.After(*waitDuration):
		log.Debugf("producer: step reach periodic timeout of [%s]", *waitDuration)
	case resultRollupInfo := <-l.workers.getResponseChannelForRollupInfo():
		l.onResponseRollupInfo(resultRollupInfo)
	}
	if l.syncStatus.haveRequiredAllBlocksToBeSynchronized() {
		// Try to nenew last block on L1 if needed
		log.Debugf("producer: we have required (maybe not responsed yet) all blocks, so  getting last block on L1")
		l.renewLastBlockOnL1IfNeeded(false)
	}
	// Try to launch retrieve more rollupInfo from L1
	l.launchWork()
	if time.Since(l.statistics.lastShowUpTime) > timeForShowUpStatisticsLog {
		log.Infof("producer: Statistics:%s", l.statistics.getETA())
		l.statistics.lastShowUpTime = time.Now()
	}
	if l.syncStatus.isNodeFullySynchronizedWithL1() {
		l.status = producerSynchronized
	} else {
		l.status = producerWorking
	}
	*waitDuration = l.getNextTimeout()
	log.Debugf("producer: Next timeout: %s status:%s sync_status: %s", *waitDuration, l.status, l.syncStatus.toStringBrief())
	return true
}

func (l *l1RollupInfoProducer) ttlOfLastBlockOnL1() time.Duration {
	return l.cfg.ttlOfLastBlockOnL1
}

func (l *l1RollupInfoProducer) getNextTimeout() time.Duration {
	switch l.status {
	case producerIdle:
		return timeOutMainLoop
	case producerWorking:
		return timeOutMainLoop
	case producerSynchronized:
		nextRenewLastBlock := time.Since(l.timeLastBLockOnL1) + l.ttlOfLastBlockOnL1()
		return max(nextRenewLastBlock, time.Second)
	default:
		log.Fatalf("producer: Unknown status: %s", l.status)
	}
	return time.Second
}

// OnNewLastBlock is called when a new last block on L1 is received
func (l *l1RollupInfoProducer) onNewLastBlock(lastBlock uint64, launchWork bool) onNewLastBlockResponse {
	resp := l.syncStatus.onNewLastBlockOnL1(lastBlock)
	l.statistics.updateLastBlockNumber(resp.fullRange.toBlock)
	l.timeLastBLockOnL1 = time.Now()
	if resp.extendedRange != nil {
		log.Infof("producer: New last block on L1: %v -> %s", resp.fullRange.toBlock, resp.toString())
	}
	if launchWork {
		l.launchWork()
	}
	return resp
}

// launchWork: launch new workers if possible and returns new channels created
// returns the number of workers launched
func (l *l1RollupInfoProducer) launchWork() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	launchedWorker := 0
	accDebugStr := ""
	thereAreAnError := false
	for {
		br := l.syncStatus.getNextRange()
		if br == nil {
			// No more work to do
			accDebugStr += "[NoNextRange] "
			break
		}
		_, err := l.workers.asyncRequestRollupInfoByBlockRange(l.ctx, *br)
		if err != nil {
			thereAreAnError = true
			accDebugStr += fmt.Sprintf(" segment %s -> [Error:%s] ", br.String(), err.Error())
			break
		}
		launchedWorker++
		log.Infof("producer: Launched worker for segment %s, num_workers_in_this_iteration: %d", br.String(), launchedWorker)
		l.syncStatus.onStartedNewWorker(*br)
	}
	if launchedWorker == 0 {
		log.Debugf("producer: No workers launched because: %s", accDebugStr)
	}
	if thereAreAnError && launchedWorker == 0 {
		log.Warnf("producer: launched workers: %d , but there are an error: %s", launchedWorker, accDebugStr)
	}
	return launchedWorker
}

func (l *l1RollupInfoProducer) renewLastBlockOnL1IfNeeded(forced bool) {
	l.mutex.Lock()
	elapsed := time.Since(l.timeLastBLockOnL1)
	ttl := l.ttlOfLastBlockOnL1()
	oldBlock := l.syncStatus.getLastBlockOnL1()
	l.mutex.Unlock()
	if elapsed > ttl || forced {
		log.Infof("producer: Need a new value for Last Block On L1, doing the request")
		result := l.workers.requestLastBlockWithRetries(l.ctx, timeRequestInitialValueOfLastBlock, maxRetriesForRequestnitialValueOfLastBlock)
		log.Infof("producer: Need a new value for Last Block On L1, doing the request old_block:%v -> new block:%v", oldBlock, result.result.block)
		if result.generic.err != nil {
			log.Error(result.generic.err)
			return
		}
		l.onNewLastBlock(result.result.block, true)
	}
}

func (l *l1RollupInfoProducer) onResponseRollupInfo(result responseRollupInfoByBlockRange) {
	l.statistics.onResponseRollupInfo(result)
	isOk := (result.generic.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
	if isOk {
		outgoingPackages := l.filterToSendOrdererResultsToConsumer.Filter(*newL1SyncMessageData(result.result))
		l.sendPackages(outgoingPackages)
	} else {
		log.Warnf("producer: Error while trying to get rollup info by block range: %v", result.generic.err)
	}
}

func (l *l1RollupInfoProducer) stop() {
	l.cancelCtx()
}

func (l *l1RollupInfoProducer) sendPackages(outgoingPackages []l1SyncMessage) {
	for _, pkg := range outgoingPackages {
		log.Infof("producer: Sending results [data] to consumer:%s:  channel status [%d/%d]", pkg.toStringBrief(), len(l.outgoingChannel), cap(l.outgoingChannel))
		l.outgoingChannel <- pkg
	}
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
