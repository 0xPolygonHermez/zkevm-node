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
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	minTTLOfLastBlock                             = time.Second
	minTimeoutForRequestLastBlockOnL1             = time.Second * 1
	minNumOfAllowedRetriesForRequestLastBlockOnL1 = 1
	minTimeOutMainLoop                            = time.Minute * 5
	timeForShowUpStatisticsLog                    = time.Second * 60
	conversionFactorPercentage                    = 100
)

type filter interface {
	ToStringBrief() string
	Filter(data l1SyncMessage) []l1SyncMessage
	Reset(lastBlockOnSynchronizer uint64)
	numItemBlockedInQueue() int
}

type syncStatusInterface interface {
	verify() error
	reset(lastBlockStoreOnStateDB uint64)
	toStringBrief() string
	getNextRange() *blockRange
	getNextRangeOnlyRetries() *blockRange
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
	String() string
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

	timeoutForRequestLastBlockOnL1             time.Duration
	numOfAllowedRetriesForRequestLastBlockOnL1 int

	//timeout for main loop if no is synchronized yet, this time is a safeguard because is not needed
	timeOutMainLoop time.Duration
	//how ofter we show a log with statistics, 0 means disabled
	timeForShowUpStatisticsLog time.Duration
}

func (cfg *configProducer) String() string {
	return fmt.Sprintf("syncChunkSize:%d ttlOfLastBlockOnL1:%s timeoutForRequestLastBlockOnL1:%s numOfAllowedRetriesForRequestLastBlockOnL1:%d timeOutMainLoop:%s timeForShowUpStatisticsLog:%s",
		cfg.syncChunkSize, cfg.ttlOfLastBlockOnL1, cfg.timeoutForRequestLastBlockOnL1, cfg.numOfAllowedRetriesForRequestLastBlockOnL1, cfg.timeOutMainLoop, cfg.timeForShowUpStatisticsLog)
}

func (cfg *configProducer) normalize() {
	if cfg.syncChunkSize == 0 {
		log.Fatalf("producer:config: SyncChunkSize must be greater than 0")
	}
	if cfg.ttlOfLastBlockOnL1 < minTTLOfLastBlock {
		log.Warnf("producer:config: ttlOfLastBlockOnL1 is too low (%s) minimum recomender value %s", cfg.ttlOfLastBlockOnL1, minTTLOfLastBlock)
	}
	if cfg.timeoutForRequestLastBlockOnL1 < minTimeoutForRequestLastBlockOnL1 {
		log.Warnf("producer:config: timeRequestInitialValueOfLastBlock is too low (%s) minimum recomender value%s", cfg.timeoutForRequestLastBlockOnL1, minTimeoutForRequestLastBlockOnL1)
	}
	if cfg.numOfAllowedRetriesForRequestLastBlockOnL1 < minNumOfAllowedRetriesForRequestLastBlockOnL1 {
		log.Warnf("producer:config: retriesForRequestnitialValueOfLastBlock is too low (%d) minimum recomender value %d", cfg.numOfAllowedRetriesForRequestLastBlockOnL1, minNumOfAllowedRetriesForRequestLastBlockOnL1)
	}
	if cfg.timeOutMainLoop < minTimeOutMainLoop {
		log.Warnf("producer:config: timeOutMainLoop is too low (%s) minimum recomender value %s", cfg.timeOutMainLoop, minTimeOutMainLoop)
	}
}

type l1RollupInfoProducer struct {
	mutex             sync.Mutex
	ctxParent         context.Context
	ctxWithCancel     contextWithCancel
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

func (l *l1RollupInfoProducer) toStringBrief() string {
	return fmt.Sprintf("status:%s syncStatus:[%s] workers:[%s] filter:[%s] cfg:[%s]", l.status, l.syncStatus.toStringBrief(), l.workers.String(), l.filterToSendOrdererResultsToConsumer.ToStringBrief(), l.cfg.String())
}

// l1DataRetrieverStatistics : create an instance of l1RollupInfoProducer
func newL1DataRetriever(cfg configProducer, ethermans []EthermanInterface, outgoingChannel chan l1SyncMessage) *l1RollupInfoProducer {
	if cap(outgoingChannel) < len(ethermans) {
		log.Warnf("producer: outgoingChannel must have a capacity (%d) of at least equal to number of ether clients (%d)", cap(outgoingChannel), len(ethermans))
	}
	cfg.normalize()
	// The timeout for clients are set to infinite because the time to process a rollup segment is not known
	// TODO: move this to config file
	workersConfig := workersConfig{timeoutRollupInfo: time.Duration(math.MaxInt64)}

	result := l1RollupInfoProducer{
		syncStatus:                           newSyncStatus(invalidBlockNumber, cfg.syncChunkSize),
		workers:                              newWorkers(ethermans, workersConfig),
		filterToSendOrdererResultsToConsumer: newFilterToSendOrdererResultsToConsumer(invalidBlockNumber),
		outgoingChannel:                      outgoingChannel,
		statistics:                           newRollupInfoProducerStatistics(invalidBlockNumber),
		status:                               producerIdle,
		cfg:                                  cfg,
	}
	return &result
}

// ResetAndStop: reset the object and stop the current process. Set first block to be retrieved
func (l *l1RollupInfoProducer) ResetAndStop(startingBlockNumber uint64) {
	log.Infof("producer: Reset L1 sync process to blockNumber %d st=%s", startingBlockNumber, l.toStringBrief())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	log.Debugf("producer: Reset(%d): stop previous run (state=%s)", startingBlockNumber, l.status.String())
	l.stopUnsafe()
	log.Debugf("producer: Reset(%d): syncStatus.reset", startingBlockNumber)
	l.syncStatus.reset(startingBlockNumber)
	l.statistics.reset(startingBlockNumber)
	// Empty pending rollupinfos
	log.Debugf("producer: Reset(%d): emptyChannel", startingBlockNumber)
	l.emptyChannel()
	log.Debugf("producer: Reset(%d): reset Filter", startingBlockNumber)
	l.filterToSendOrdererResultsToConsumer.Reset(startingBlockNumber)
	log.Debugf("producer: Reset(%d): reset done!", startingBlockNumber)
}

func (l *l1RollupInfoProducer) Stop() {
	log.Debugf("producer: stop() called st=%s", l.toStringBrief())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.stopUnsafe()
}

// stopUnsafe: stop the object without locking the mutex (need to be locked before call it!)
func (l *l1RollupInfoProducer) stopUnsafe() {
	if l.status != producerIdle {
		log.Infof("producer: stopping producer")
		l.ctxWithCancel.cancel()
		l.status = producerIdle
	}
	log.Debugf("producer: stopUnsafe: stop workers (%s)", l.workers.String())
	l.workers.stop()
	l.workers.waitFinishAllWorkers()
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

func (l *l1RollupInfoProducer) initialize(ctx context.Context) error {
	log.Debug("producer: initialize")
	err := l.verify()
	if err != nil {
		return err
	}
	if l.ctxParent != ctx || l.ctxWithCancel.isInvalid() {
		log.Debug("producer: start called and need to create a new context")
		l.ctxParent = ctx
		l.ctxWithCancel.createWithCancel(l.ctxParent)
	}
	err = l.workers.initialize()
	if err != nil {
		return err
	}
	if l.syncStatus.isSetLastBlockOnL1Value() {
		log.Infof("producer: Need a initial value for Last Block On L1, doing the request (maxRetries:%v, timeRequest:%v)",
			l.cfg.numOfAllowedRetriesForRequestLastBlockOnL1, l.cfg.timeoutForRequestLastBlockOnL1)
		//result := l.retrieveInitialValueOfLastBlock(maxRetriesForRequestnitialValueOfLastBlock, timeRequestInitialValueOfLastBlock)
		result := l.workers.requestLastBlockWithRetries(l.ctxWithCancel.ctx, l.cfg.timeoutForRequestLastBlockOnL1, l.cfg.numOfAllowedRetriesForRequestLastBlockOnL1)
		if result.generic.err != nil {
			log.Error(result.generic.err)
			return result.generic.err
		}
		l.onNewLastBlock(result.result.block, false)
	}

	return nil
}

// Before calling Start you must set lastBlockOnDB calling ResetAndStop
func (l *l1RollupInfoProducer) Start(ctx context.Context) error {
	log.Infof("producer: starting L1 sync from:%s", l.syncStatus.toStringBrief())
	err := l.initialize(ctx)
	if err != nil {
		log.Infof("producer:  can't start because: %s", err.Error())
		return err
	}
	log.Debugf("producer:  starting configuration: %s", l.cfg.String())
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
	case <-l.ctxWithCancel.Done():
		log.Debugf("producer: context canceled")
		return false
	// That timeout is not need, but just in case that stop launching request
	case <-time.After(*waitDuration):
		log.Debugf("producer: reach timeout of step loop it was of %s", *waitDuration)
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
	if l.cfg.timeForShowUpStatisticsLog != 0 && time.Since(l.statistics.lastShowUpTime) > l.cfg.timeForShowUpStatisticsLog {
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
	timeOutMainLoop := l.cfg.timeOutMainLoop
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
	return timeOutMainLoop
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

func (l *l1RollupInfoProducer) canISendNewRequestsUnsafe() (bool, string) {
	queued := l.filterToSendOrdererResultsToConsumer.numItemBlockedInQueue()
	inChannel := len(l.outgoingChannel)
	maximum := cap(l.outgoingChannel)
	msg := fmt.Sprintf("inFilter:%d + inChannel:%d > maximum:%d?", queued, inChannel, maximum)
	if queued+inChannel > maximum {
		msg = msg + " ==> only allow retries"
		return false, msg
	}
	msg = msg + " ==> allow new req"
	return true, msg
}

// launchWork: launch new workers if possible and returns new channels created
// returns the number of workers launched
func (l *l1RollupInfoProducer) launchWork() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	launchedWorker := 0
	allowNewRequests, allowNewRequestMsg := l.canISendNewRequestsUnsafe()
	accDebugStr := "[" + allowNewRequestMsg + "] "
	for {
		var br *blockRange
		if allowNewRequests {
			br = l.syncStatus.getNextRange()
		} else {
			br = l.syncStatus.getNextRangeOnlyRetries()
		}
		if br == nil {
			// No more work to do
			accDebugStr += "[NoNextRange] "
			break
		}
		_, err := l.workers.asyncRequestRollupInfoByBlockRange(l.ctxWithCancel.ctx, *br)
		if err != nil {
			if errors.Is(err, errAllWorkersBusy) {
				accDebugStr += fmt.Sprintf(" segment %s -> [Error:%s] ", br.String(), err.Error())
			}
			break
		} else {
			accDebugStr += fmt.Sprintf(" segment %s -> [LAUNCHED] ", br.String())
		}
		launchedWorker++
		log.Debugf("producer: launch_worker: Launched worker for segment %s, num_workers_in_this_iteration: %d", br.String(), launchedWorker)
		l.syncStatus.onStartedNewWorker(*br)
	}
	log.Infof("producer: launch_worker:  num of launched workers: %d  result: %s status_comm:%s", launchedWorker, accDebugStr, l.outgoingPackageStatusDebugString())

	return launchedWorker
}

func (l *l1RollupInfoProducer) outgoingPackageStatusDebugString() string {
	return fmt.Sprintf("outgoint_channel[%d/%d], filter:%s workers:%s", len(l.outgoingChannel), cap(l.outgoingChannel), l.filterToSendOrdererResultsToConsumer.ToStringBrief(), l.workers.String())
}

func (l *l1RollupInfoProducer) renewLastBlockOnL1IfNeeded(forced bool) {
	l.mutex.Lock()
	elapsed := time.Since(l.timeLastBLockOnL1)
	ttl := l.ttlOfLastBlockOnL1()
	oldBlock := l.syncStatus.getLastBlockOnL1()
	l.mutex.Unlock()
	if elapsed > ttl || forced {
		log.Infof("producer: Need a new value for Last Block On L1, doing the request")
		result := l.workers.requestLastBlockWithRetries(l.ctxWithCancel.ctx, l.cfg.timeoutForRequestLastBlockOnL1, l.cfg.numOfAllowedRetriesForRequestLastBlockOnL1)
		log.Infof("producer: Need a new value for Last Block On L1, doing the request old_block:%v -> new block:%v", oldBlock, result.result.block)
		if result.generic.err != nil {
			log.Error(result.generic.err)
			return
		}
		l.onNewLastBlock(result.result.block, true)
	}
}

func (l *l1RollupInfoProducer) onResponseRollupInfo(result responseRollupInfoByBlockRange) {
	log.Infof("producer: Received responseRollupInfoByBlockRange: %s", result.toStringBrief())
	l.statistics.onResponseRollupInfo(result)
	isOk := (result.generic.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
	if isOk {
		outgoingPackages := l.filterToSendOrdererResultsToConsumer.Filter(*newL1SyncMessageData(result.result))
		l.sendPackages(outgoingPackages)
	} else {
		if errors.Is(result.generic.err, context.Canceled) {
			log.Infof("producer: Error while trying to get rollup info by block range: %v", result.generic.err)
		} else {
			log.Warnf("producer: Error while trying to get rollup info by block range: %v", result.generic.err)
		}
	}
}

func (l *l1RollupInfoProducer) sendPackages(outgoingPackages []l1SyncMessage) {
	for _, pkg := range outgoingPackages {
		log.Infof("producer: Sending results [data] to consumer:%s:  status_comm:%s", pkg.toStringBrief(), l.outgoingPackageStatusDebugString())
		l.outgoingChannel <- pkg
	}
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
