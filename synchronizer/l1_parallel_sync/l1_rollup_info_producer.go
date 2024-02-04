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

package l1_parallel_sync

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
)

const (
	minTTLOfLastBlock                             = time.Second
	minTimeoutForRequestLastBlockOnL1             = time.Second * 1
	minNumOfAllowedRetriesForRequestLastBlockOnL1 = 1
	minTimeOutMainLoop                            = time.Minute * 5
	timeForShowUpStatisticsLog                    = time.Second * 60
	conversionFactorPercentage                    = 100
	lenCommandsChannels                           = 5
	maximumBlockDistanceFromLatestToFinalized     = 96 // https://www.alchemy.com/overviews/ethereum-commitment-levels
)

type filter interface {
	ToStringBrief() string
	Filter(data L1SyncMessage) []L1SyncMessage
	Reset(lastBlockOnSynchronizer uint64)
	numItemBlockedInQueue() int
}

type syncStatusInterface interface {
	// Verify that configuration and lastBlock are right
	Verify() error
	// Reset synchronization to a new starting point
	Reset(lastBlockStoreOnStateDB uint64)
	// String returns a string representation of the object
	String() string
	// GetNextRange returns the next Block to be retrieved
	GetNextRange() *blockRange
	// GetNextRangeOnlyRetries returns the fist Block pending to retry
	GetNextRangeOnlyRetries() *blockRange
	// IsNodeFullySynchronizedWithL1 returns true there nothing pending to retrieved and have finished all workers
	//    so all the rollupinfo are in the channel to be processed by consumer
	IsNodeFullySynchronizedWithL1() bool
	// HaveRequiredAllBlocksToBeSynchronized returns true if have been requested all rollupinfo
	//     but maybe there are some pending retries or still working in some BlockRange
	HaveRequiredAllBlocksToBeSynchronized() bool
	// DoesItHaveAllTheNeedDataToWork returns true if have all the data to start working
	DoesItHaveAllTheNeedDataToWork() bool
	// GetLastBlockOnL1 returns the last block on L1 or InvalidBlock if not set
	GetLastBlockOnL1() uint64

	// OnStartedNewWorker a new worker has been started
	OnStartedNewWorker(br blockRange)
	// OnFinishWorker a worker has finished, returns true if the data have to be processed
	OnFinishWorker(br blockRange, successful bool, highestBlockNumberInResponse uint64) bool
	// OnNewLastBlockOnL1 a new last block on L1 has been received
	OnNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse
	// BlockNumberIsInsideUnsafeArea returns if this block is beyond Finalized (so it could be reorg)
	// If blockNumber == invalidBlockNumber then it uses the highestBlockRequested (the last block requested)
	BlockNumberIsInsideUnsafeArea(blockNumber uint64) bool
	// GetHighestBlockReceived returns the highest block requested
	GetHighestBlockReceived() uint64
}

type workersInterface interface {
	// initialize object
	initialize() error
	// finalize object
	stop()
	// waits until all workers have finish the current task
	waitFinishAllWorkers()
	asyncRequestRollupInfoByBlockRange(ctx context.Context, request requestRollupInfoByBlockRange) (chan responseRollupInfoByBlockRange, error)
	requestLastBlockWithRetries(ctx context.Context, timeout time.Duration, maxPermittedRetries int) responseL1LastBlock
	getResponseChannelForRollupInfo() chan responseRollupInfoByBlockRange
	String() string
	ToStringBrief() string
	howManyRunningWorkers() int
}

type producerStatusEnum int32

const (
	producerIdle         producerStatusEnum = 0
	producerWorking      producerStatusEnum = 1
	producerSynchronized producerStatusEnum = 2
	producerNoRunning    producerStatusEnum = 3
	// producerReseting: is in a reset process, so is going to reject all rollup info
	producerReseting producerStatusEnum = 4
)

func (s producerStatusEnum) String() string {
	return [...]string{"idle", "working", "synchronized", "no_running", "reseting"}[s]
}

// ConfigProducer : configuration for producer
type ConfigProducer struct {
	// SyncChunkSize is the number of blocks to be retrieved in each request
	SyncChunkSize uint64
	// TtlOfLastBlockOnL1 is the time to wait before ask for a new last block on L1
	TtlOfLastBlockOnL1 time.Duration
	// TimeoutForRequestLastBlockOnL1 is the timeout for request a new last block on L1
	TimeoutForRequestLastBlockOnL1 time.Duration
	// NumOfAllowedRetriesForRequestLastBlockOnL1 is the number of retries for request a new last block on L1
	NumOfAllowedRetriesForRequestLastBlockOnL1 int

	//TimeOutMainLoop timeout for main loop if no is synchronized yet, this time is a safeguard because is not needed
	TimeOutMainLoop time.Duration
	//TimeForShowUpStatisticsLog how ofter we show a log with statistics, 0 means disabled
	TimeForShowUpStatisticsLog time.Duration
	// MinTimeBetweenRetriesForRollupInfo is the minimum time between retries for rollup info
	MinTimeBetweenRetriesForRollupInfo time.Duration
}

func (cfg *ConfigProducer) String() string {
	return fmt.Sprintf("syncChunkSize:%d ttlOfLastBlockOnL1:%s timeoutForRequestLastBlockOnL1:%s numOfAllowedRetriesForRequestLastBlockOnL1:%d timeOutMainLoop:%s timeForShowUpStatisticsLog:%s",
		cfg.SyncChunkSize, cfg.TtlOfLastBlockOnL1, cfg.TimeoutForRequestLastBlockOnL1, cfg.NumOfAllowedRetriesForRequestLastBlockOnL1, cfg.TimeOutMainLoop, cfg.TimeForShowUpStatisticsLog)
}

func (cfg *ConfigProducer) normalize() {
	if cfg.SyncChunkSize == 0 {
		log.Fatalf("producer:config: SyncChunkSize must be greater than 0")
	}
	if cfg.TtlOfLastBlockOnL1 < minTTLOfLastBlock {
		log.Warnf("producer:config: ttlOfLastBlockOnL1 is too low (%s) minimum recomender value %s", cfg.TtlOfLastBlockOnL1, minTTLOfLastBlock)
	}
	if cfg.TimeoutForRequestLastBlockOnL1 < minTimeoutForRequestLastBlockOnL1 {
		log.Warnf("producer:config: timeRequestInitialValueOfLastBlock is too low (%s) minimum recomender value%s", cfg.TimeoutForRequestLastBlockOnL1, minTimeoutForRequestLastBlockOnL1)
	}
	if cfg.NumOfAllowedRetriesForRequestLastBlockOnL1 < minNumOfAllowedRetriesForRequestLastBlockOnL1 {
		log.Warnf("producer:config: retriesForRequestnitialValueOfLastBlock is too low (%d) minimum recomender value %d", cfg.NumOfAllowedRetriesForRequestLastBlockOnL1, minNumOfAllowedRetriesForRequestLastBlockOnL1)
	}
	if cfg.TimeOutMainLoop < minTimeOutMainLoop {
		log.Warnf("producer:config: timeOutMainLoop is too low (%s) minimum recomender value %s", cfg.TimeOutMainLoop, minTimeOutMainLoop)
	}
	if cfg.MinTimeBetweenRetriesForRollupInfo <= 0 {
		log.Warnf("producer:config: minTimeBetweenRetriesForRollup is too low (%s)", cfg.MinTimeBetweenRetriesForRollupInfo)
	}
}

type producerCmdEnum int32

const (
	producerNop   producerCmdEnum = 0
	producerStop  producerCmdEnum = 1
	producerReset producerCmdEnum = 2
)

func (s producerCmdEnum) String() string {
	return [...]string{"nop", "stop", "reset"}[s]
}

type producerCmd struct {
	cmd    producerCmdEnum
	param1 uint64
}

// L1RollupInfoProducer is the object that retrieves data from L1
type L1RollupInfoProducer struct {
	mutex             sync.Mutex
	ctxParent         context.Context
	ctxWithCancel     contextWithCancel
	workers           workersInterface
	syncStatus        syncStatusInterface
	outgoingChannel   chan L1SyncMessage
	timeLastBLockOnL1 time.Time
	status            producerStatusEnum
	// filter is an object that sort l1DataMessage to be send ordered by block number
	filterToSendOrdererResultsToConsumer filter
	statistics                           l1RollupInfoProducerStatistics
	cfg                                  ConfigProducer
	channelCmds                          chan producerCmd
}

func (l *L1RollupInfoProducer) toStringBrief() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.toStringBriefUnsafe()
}

func (l *L1RollupInfoProducer) toStringBriefUnsafe() string {
	return fmt.Sprintf("status:%s syncStatus:[%s] workers:[%s] filter:[%s] cfg:[%s]", l.getStatus(), l.syncStatus.String(), l.workers.String(), l.filterToSendOrdererResultsToConsumer.ToStringBrief(), l.cfg.String())
}

// NewL1DataRetriever creates a new object
func NewL1DataRetriever(cfg ConfigProducer, ethermans []L1ParallelEthermanInterface, outgoingChannel chan L1SyncMessage) *L1RollupInfoProducer {
	if cap(outgoingChannel) < len(ethermans) {
		log.Warnf("producer: outgoingChannel must have a capacity (%d) of at least equal to number of ether clients (%d)", cap(outgoingChannel), len(ethermans))
	}
	cfg.normalize()
	// The timeout for clients are set to infinite because the time to process a rollup segment is not known
	// TODO: move this to config file
	workersConfig := workersConfig{timeoutRollupInfo: time.Duration(math.MaxInt64)}

	result := L1RollupInfoProducer{
		syncStatus:                           newSyncStatus(invalidBlockNumber, cfg.SyncChunkSize),
		workers:                              newWorkerDecoratorLimitRetriesByTime(newWorkers(ethermans, workersConfig), cfg.MinTimeBetweenRetriesForRollupInfo),
		filterToSendOrdererResultsToConsumer: newFilterToSendOrdererResultsToConsumer(invalidBlockNumber),
		outgoingChannel:                      outgoingChannel,
		statistics:                           newRollupInfoProducerStatistics(invalidBlockNumber, common.DefaultTimeProvider{}),
		status:                               producerNoRunning,
		cfg:                                  cfg,
		channelCmds:                          make(chan producerCmd, lenCommandsChannels),
	}
	return &result
}

// Reset reset the object and stop the current process. Set first block to be retrieved
func (l *L1RollupInfoProducer) Reset(startingBlockNumber uint64) {
	log.Infof("producer: Reset(%d) queue cmd and discarding all info in channel", startingBlockNumber)
	l.setStatusReseting()
	l.emptyChannel()
	l.channelCmds <- producerCmd{cmd: producerReset, param1: startingBlockNumber}
}

func (l *L1RollupInfoProducer) resetUnsafe(startingBlockNumber uint64) {
	log.Debugf("producer: Reset L1 sync process to blockNumber %d st=%s", startingBlockNumber, l.toStringBrief())
	l.setStatusReseting()
	log.Debugf("producer: Reset(%d): stop previous run (state=%s)", startingBlockNumber, l.getStatus().String())
	log.Debugf("producer: Reset(%d): syncStatus.reset", startingBlockNumber)
	l.syncStatus.Reset(startingBlockNumber)
	l.statistics.reset(startingBlockNumber)
	log.Debugf("producer: Reset(%d): stopping workers", startingBlockNumber)
	l.workers.stop()
	// Empty pending rollupinfos
	log.Debugf("producer: Reset(%d): emptyChannel", startingBlockNumber)
	l.emptyChannel()
	log.Debugf("producer: Reset(%d): reset Filter", startingBlockNumber)
	l.filterToSendOrdererResultsToConsumer.Reset(startingBlockNumber)
	l.setStatus(producerIdle)
	log.Infof("producer: Reset(%d): reset producer done!", startingBlockNumber)
}

func (l *L1RollupInfoProducer) isProducerRunning() bool {
	return l.getStatus() != producerNoRunning
}

func (l *L1RollupInfoProducer) setStatusReseting() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.setStatus(producerReseting)
}

func (l *L1RollupInfoProducer) getStatus() producerStatusEnum {
	return producerStatusEnum(atomic.LoadInt32((*int32)(&l.status)))
}

func (l *L1RollupInfoProducer) setStatus(newStatus producerStatusEnum) {
	previousStatus := l.getStatus()
	atomic.StoreInt32((*int32)(&l.status), int32(newStatus))
	if previousStatus != newStatus {
		log.Infof("producer: Status changed from [%s] to [%s]", previousStatus.String(), newStatus.String())
		if newStatus == producerSynchronized {
			highestBlock := l.syncStatus.GetHighestBlockReceived()
			log.Infof("producer: send a message to consumer to indicate that we are synchronized. highestBlockRequested:%d", highestBlock)
			l.sendPackages([]L1SyncMessage{*newL1SyncMessageControlWProducerIsFullySynced(highestBlock)})
		}
	}
}

// Abort stop inmediatly the current process
func (l *L1RollupInfoProducer) Abort() {
	l.emptyChannel()
	l.ctxWithCancel.cancelCtx()
	l.ctxWithCancel.createWithCancel(l.ctxParent)
}

// Stop stop the current process sending a stop command to the process queue
// so it stops when finish to process all packages in queue
func (l *L1RollupInfoProducer) Stop() {
	log.Infof("producer: Stop() queue cmd")
	l.channelCmds <- producerCmd{cmd: producerStop}
}

func (l *L1RollupInfoProducer) stopUnsafe() {
	log.Infof("producer: stop() called st=%s", l.toStringBrief())

	if l.isProducerRunning() {
		log.Infof("producer:Stop:was running -> stopping producer")
		l.ctxWithCancel.cancel()
	}

	l.setStatus(producerNoRunning)
	log.Debugf("producer:Stop: stop workers and wait for finish (%s)", l.workers.String())
	l.workers.stop()
}

func (l *L1RollupInfoProducer) emptyChannel() {
	for len(l.outgoingChannel) > 0 {
		<-l.outgoingChannel
	}
}

// verify: test params and status without if not allowModify avoid doing connection or modification of objects
func (l *L1RollupInfoProducer) verify() error {
	return l.syncStatus.Verify()
}

func (l *L1RollupInfoProducer) initialize(ctx context.Context) error {
	log.Debug("producer: initialize")
	err := l.verify()
	if err != nil {
		log.Debug("producer: initialize, syncstatus not ready: %s", err.Error())
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
	return nil
}

// Start a producer
func (l *L1RollupInfoProducer) Start(ctx context.Context) error {
	log.Infof("producer: starting L1 sync from:%s", l.syncStatus.String())
	err := l.initialize(ctx)
	if err != nil {
		log.Infof("producer:  can't start because: %s", err.Error())
		return err
	}
	l.setStatus(producerIdle)
	log.Debugf("producer:  starting configuration: %s", l.cfg.String())
	var waitDuration = time.Duration(0)
	for l.step(&waitDuration) {
	}
	l.setStatus(producerNoRunning)
	l.workers.waitFinishAllWorkers()
	return nil
}

func (l *L1RollupInfoProducer) step(waitDuration *time.Duration) bool {
	if atomic.CompareAndSwapInt32((*int32)(&l.status), int32(producerNoRunning), int32(producerIdle)) { // l.getStatus() == producerNoRunning
		log.Info("producer: step: status is no running, changing to idle  %s", l.getStatus().String())
	}
	log.Debugf("producer: step: status:%s", l.toStringBrief())
	select {
	case <-l.ctxWithCancel.Done():
		log.Debugf("producer: context canceled")
		return false
	case cmd := <-l.channelCmds:
		log.Debugf("producer: received a command")
		res := l.executeCmd(cmd)
		if !res {
			log.Info("producer: cmd %s stop the process", cmd.cmd.String())
			return false
		}
	// That timeout is not need, but just in case that stop launching request
	case <-time.After(*waitDuration):
		log.Debugf("producer: reach timeout of step loop it was of %s", *waitDuration)
	case resultRollupInfo := <-l.workers.getResponseChannelForRollupInfo():
		l.onResponseRollupInfo(resultRollupInfo)
	}
	switch l.getStatus() {
	case producerIdle:
		// Is ready to start working?
		l.renewLastBlockOnL1IfNeeded("idle")
		if l.syncStatus.DoesItHaveAllTheNeedDataToWork() {
			log.Infof("producer: producerIdle: have all the data to work, moving to working status.  status:%s", l.syncStatus.String())
			l.setStatus(producerWorking)
			// This is for wakeup the step again to launch a new work
			l.channelCmds <- producerCmd{cmd: producerNop}
		} else {
			log.Infof("producer: producerIdle: still dont have all the data to work status:%s", l.syncStatus.String())
		}
	case producerWorking:
		// launch new Work
		_, err := l.launchWork()
		if err != nil {
			log.Errorf("producer: producerWorking: error launching work: %s", err.Error())
			return false
		}
		// If I'm have required all blocks to L1?
		if l.syncStatus.HaveRequiredAllBlocksToBeSynchronized() {
			log.Debugf("producer: producerWorking: haveRequiredAllBlocksToBeSynchronized -> renewLastBlockOnL1IfNeeded")
			l.renewLastBlockOnL1IfNeeded("HaveRequiredAllBlocksToBeSynchronized")
		}
		if l.syncStatus.BlockNumberIsInsideUnsafeArea(invalidBlockNumber) {
			log.Debugf("producer: producerWorking: we are inside unsafe area!, renewLastBlockOnL1IfNeeded")
			l.renewLastBlockOnL1IfNeeded("unsafe block area")
		}
		// If after asking for a new lastBlockOnL1 we are still synchronized then we are synchronized
		if l.syncStatus.IsNodeFullySynchronizedWithL1() {
			l.setStatus(producerSynchronized)
		} else {
			log.Infof("producer: producerWorking: still not synchronized with the new block range launch workers again")
			_, err := l.launchWork()
			if err != nil {
				log.Errorf("producer: producerSynchronized: error launching work: %s", err.Error())
				return false
			}
		}
	case producerSynchronized:
		// renew last block on L1 if needed
		log.Debugf("producer: producerSynchronized")
		l.renewLastBlockOnL1IfNeeded("producerSynchronized")

		numLaunched, err := l.launchWork()
		if err != nil {
			log.Errorf("producer: producerSynchronized: error launching work: %s", err.Error())
			return false
		}
		if numLaunched > 0 {
			l.setStatus(producerWorking)
		}
	case producerReseting:
		log.Infof("producer: producerReseting")
	}

	if l.cfg.TimeForShowUpStatisticsLog != 0 && time.Since(l.statistics.lastShowUpTime) > l.cfg.TimeForShowUpStatisticsLog {
		log.Infof("producer: Statistics:%s", l.statistics.getStatisticsDebugString())
		l.statistics.lastShowUpTime = time.Now()
	}
	*waitDuration = l.getNextTimeout()
	log.Debugf("producer: Next timeout: %s status:%s ", *waitDuration, l.toStringBrief())
	return true
}

// return if the producer must keep running (false -> stop)
func (l *L1RollupInfoProducer) executeCmd(cmd producerCmd) bool {
	switch cmd.cmd {
	case producerStop:
		log.Infof("producer: received a stop, so it stops requesting new rollup info and stop current requests")
		l.stopUnsafe()
		return false
	case producerReset:
		log.Infof("producer: received a reset(%d)", cmd.param1)
		l.resetUnsafe(cmd.param1)
		return true
	}
	return true
}

func (l *L1RollupInfoProducer) ttlOfLastBlockOnL1() time.Duration {
	return l.cfg.TtlOfLastBlockOnL1
}

func (l *L1RollupInfoProducer) getNextTimeout() time.Duration {
	timeOutMainLoop := l.cfg.TimeOutMainLoop
	status := l.getStatus()
	switch status {
	case producerIdle:
		return timeOutMainLoop
	case producerWorking:
		return timeOutMainLoop
	case producerSynchronized:
		nextRenewLastBlock := time.Since(l.timeLastBLockOnL1) + l.ttlOfLastBlockOnL1()
		return max(nextRenewLastBlock, time.Second)
	case producerNoRunning:
		return timeOutMainLoop
	case producerReseting:
		return timeOutMainLoop
	default:
		log.Fatalf("producer: Unknown status: %s", status.String())
	}
	return timeOutMainLoop
}

// OnNewLastBlock is called when a new last block on L1 is received
func (l *L1RollupInfoProducer) onNewLastBlock(lastBlock uint64) onNewLastBlockResponse {
	resp := l.syncStatus.OnNewLastBlockOnL1(lastBlock)
	l.statistics.updateLastBlockNumber(resp.fullRange.toBlock)
	l.timeLastBLockOnL1 = time.Now()
	if resp.extendedRange != nil {
		log.Infof("producer: New last block on L1: %v -> %s", resp.fullRange.toBlock, resp.toString())
	}
	return resp
}

func (l *L1RollupInfoProducer) canISendNewRequestsUnsafe() (bool, string) {
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
func (l *L1RollupInfoProducer) launchWork() (int, error) {
	launchedWorker := 0
	allowNewRequests, allowNewRequestMsg := l.canISendNewRequestsUnsafe()
	accDebugStr := "[" + allowNewRequestMsg + "] "
	for {
		var br *blockRange
		if allowNewRequests {
			br = l.syncStatus.GetNextRange()
		} else {
			br = l.syncStatus.GetNextRangeOnlyRetries()
		}
		if br == nil {
			// No more work to do
			accDebugStr += "[NoNextRange] "
			break
		}
		// The request include previous block only if a latest request, because then it starts from l
		request := requestRollupInfoByBlockRange{blockRange: *br,
			sleepBefore:                        noSleepTime,
			requestLastBlockIfNoBlocksInAnswer: requestLastBlockModeIfNoBlocksInAnswer,
			requestPreviousBlock:               false,
		}
		unsafeAreaMsg := ""
		// GetLastBlockOnL1 is the lastest block on L1, if we are not in safe zone of reorgs we ask for previous and last block
		//   to be able to check that there is no reorgs
		if l.syncStatus.BlockNumberIsInsideUnsafeArea(br.fromBlock) {
			log.Debugf("L1 unsafe zone: asking for previous and last block")
			request.requestLastBlockIfNoBlocksInAnswer = requestLastBlockModeAlways
			request.requestPreviousBlock = true
			unsafeAreaMsg = "/UNSAFE"
		}
		blockRangeMsg := br.String() + unsafeAreaMsg
		_, err := l.workers.asyncRequestRollupInfoByBlockRange(l.ctxWithCancel.ctx, request)
		if err != nil {
			if !errors.Is(err, errAllWorkersBusy) {
				accDebugStr += fmt.Sprintf(" segment %s -> [Error:%s] ", blockRangeMsg, err.Error())
			}
			break
		} else {
			accDebugStr += fmt.Sprintf(" segment %s -> [LAUNCHED] ", blockRangeMsg)
		}
		launchedWorker++
		log.Debugf("producer: launch_worker: Launched worker for segment %s, num_workers_in_this_iteration: %d", blockRangeMsg, launchedWorker)
		l.syncStatus.OnStartedNewWorker(*br)
	}
	if launchedWorker > 0 {
		log.Infof("producer: launch_worker: num of launched workers: %d  (%s) result: %s ", launchedWorker, l.workers.ToStringBrief(), accDebugStr)
	}
	log.Debugf("producer: launch_worker:  num of launched workers: %d  result: %s status_comm:%s", launchedWorker, accDebugStr, l.outgoingPackageStatusDebugString())

	return launchedWorker, nil
}

func (l *L1RollupInfoProducer) outgoingPackageStatusDebugString() string {
	return fmt.Sprintf("outgoint_channel[%d/%d], filter:%s workers:%s", len(l.outgoingChannel), cap(l.outgoingChannel), l.filterToSendOrdererResultsToConsumer.ToStringBrief(), l.workers.String())
}

func (l *L1RollupInfoProducer) renewLastBlockOnL1IfNeeded(reason string) {
	elapsed := time.Since(l.timeLastBLockOnL1)
	ttl := l.ttlOfLastBlockOnL1()
	oldBlock := l.syncStatus.GetLastBlockOnL1()
	if elapsed > ttl {
		log.Debugf("producer: Need a new value for Last Block On L1, doing the request reason:%s", reason)
		result := l.workers.requestLastBlockWithRetries(l.ctxWithCancel.ctx, l.cfg.TimeoutForRequestLastBlockOnL1, l.cfg.NumOfAllowedRetriesForRequestLastBlockOnL1)
		if result.generic.err != nil {
			return
		}
		log.Infof("producer: Need a new value for Last Block On L1, doing the request old_block:%v -> new block:%v", oldBlock, result.result.block)

		l.onNewLastBlock(result.result.block)
	}
}

func (l *L1RollupInfoProducer) onResponseRollupInfo(result responseRollupInfoByBlockRange) {
	log.Infof("producer: Received responseRollupInfoByBlockRange: %s", result.toStringBrief())
	if l.getStatus() == producerReseting {
		log.Infof("producer: Ignoring result because is in reseting status: %s", result.toStringBrief())
		return
	}
	l.statistics.onResponseRollupInfo(result)
	isOk := (result.generic.err == nil)
	var highestBlockNumberInResponse uint64 = invalidBlockNumber
	if isOk {
		highestBlockNumberInResponse = result.getHighestBlockNumberInResponse()
	}
	if !l.syncStatus.OnFinishWorker(result.result.blockRange, isOk, highestBlockNumberInResponse) {
		log.Infof("producer: Ignoring result because the range is not longer valid: %s", result.toStringBrief())
		return
	}
	if isOk {
		outgoingPackages := l.filterToSendOrdererResultsToConsumer.Filter(*newL1SyncMessageData(result.result))
		log.Debugf("producer: filtered Br[%s/%d], outgoing %d filter_status:%s", result.result.blockRange.String(), result.result.getHighestBlockNumberInResponse(), len(outgoingPackages), l.filterToSendOrdererResultsToConsumer.ToStringBrief())
		if len(outgoingPackages) > 0 {
			for idx, msg := range outgoingPackages {
				log.Infof("producer: sendind data to consumer: [%d/%d] -> range:[%s] Sending results [data] to consumer:%s ", idx, len(outgoingPackages), result.result.blockRange.String(), msg.toStringBrief())
			}
		}
		l.sendPackages(outgoingPackages)
	} else {
		if errors.Is(result.generic.err, context.Canceled) {
			log.Infof("producer: Error while trying to get rollup info by block range: %v", result.generic.err)
		} else {
			log.Warnf("producer: Error while trying to get rollup info by block range: %v", result.generic.err)
		}
	}
}

func (l *L1RollupInfoProducer) sendPackages(outgoingPackages []L1SyncMessage) {
	for _, pkg := range outgoingPackages {
		log.Debugf("producer: Sending results [data] to consumer:%s:  status_comm:%s", pkg.toStringBrief(), l.outgoingPackageStatusDebugString())
		l.outgoingChannel <- pkg
	}
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
