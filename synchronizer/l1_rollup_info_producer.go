// package synchronizer
// Implements the logic to retrieve data from L1 and send it to the synchronizer
//   - multiples etherman to do it in parallel
//   - generate blocks to be retrieved
//   - retrieve blocks (parallel)
//   - check last block on L1:
//     The idea is to run all the time and renew the las block from L1
//     but for best fitting current implementation of synchronizer and match the
//     previous behaviour of syncBlocks we update the last block at beginning of the process
//     and finish the process when we reach the last block.
//     To control that:
//   - cte: ttlOfLastBlockDefault
//   - when creating object param renewLastBlockOnL1
//
// TODO:
//   - All the stuff related to update last block on L1 could be moved to another class
//   - Check context usage:
//     It need a context to cancel it self and create another context to cancel workers?
//   - Emit metrics
//   - if nothing to update reduce de code to be executed
//   - Improve the unittest of this object
//   - Check all log.fatals to remove it or add status before the panic
//   - Old syncBlocks method try to ask for blocks over last L1 block, I suppose that is for keep
//     synchronizing even a long the synchronization have new blocks. This is not implemented here
//     This is the behaviour of ethman in that situation:
//   - GetRollupInfoByBlockRange returns no error, zero blocks...
//   - EthBlockByNumber returns error:  "not found"

package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	ttlOfLastBlockDefault                      = time.Second * 60
	timeOutMainLoop                            = time.Minute * 5
	timeForShowUpStatisticsLog                 = time.Second * 60
	conversionFactorPercentage                 = 100
	maxRetriesForRequestnitialValueOfLastBlock = 10
	timeRequestInitialValueOfLastBlock         = time.Second * 5
)

type filter interface {
	toStringBrief() string
	filter(data l1SyncMessage) []l1SyncMessage
}

type syncStatusInterface interface {
	verify(allowModify bool) error
	toStringBrief() string
	getNextRange() *blockRange
	isNodeFullySynchronizedWithL1() bool
	needToRenewLastBlockOnL1() bool

	onStartedNewWorker(br blockRange)
	onFinishWorker(br blockRange, successful bool)
	onNewLastBlockOnL1(lastBlock uint64) onNewLastBlockResponse
}

type workersInterface interface {
	// verify test params, if allowModify = true allow to change things or make connections
	verify(allowModify bool) error
	// initialize object
	initialize() error
	// finalize object
	finalize() error
	// waits until all workers have finish the current task
	waitFinishAllWorkers()

	// asyncRetrieveLastBlock start a async request to retrieve the last block
	asyncRequestLastBlock(ctx context.Context) (chan genericResponse[retrieveL1LastBlockResult], error)
	// asyncGetRollupInfoByBlockRange start a async request to retrieve the rollup info
	asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[responseRollupInfoByBlockRange], error)

	// requestRollupInfoByBlockRange start a sync request to retrieve the last block
	requestLastBlock(ctx context.Context, timeout time.Duration) genericResponse[retrieveL1LastBlockResult]

	getResponseChannelForLastBlock() chan genericResponse[retrieveL1LastBlockResult]
	getResponseChannelForRollupInfo() chan genericResponse[responseRollupInfoByBlockRange]
}

type l1RollupInfoProducer struct {
	mutex           sync.Mutex
	ctx             context.Context
	cancelCtx       context.CancelFunc
	workers         workersInterface
	syncStatus      syncStatusInterface
	outgoingChannel chan l1SyncMessage
	// filter is an object that sort l1DataMessage to be send ordered by block number
	filterToSendOrdererResultsToConsumer filter
	statistics                           l1RollupInfoProducerStatistics
}

// l1DataRetrieverStatistics : create an instance of l1RollupInfoProducer
func newL1DataRetriever(ctx context.Context, ethermans []EthermanInterface,
	startingBlockNumber uint64, SyncChunkSize uint64,
	outgoingChannel chan l1SyncMessage, renewLastBlockOnL1 bool) *l1RollupInfoProducer {
	if cap(outgoingChannel) < len(ethermans) {
		log.Warnf("l1RollupInfoProducer: outgoingChannel must have a capacity (%d) of at least equal to number of ether clients (%d)", cap(outgoingChannel), len(ethermans))
	}
	ctx, cancel := context.WithCancel(ctx)
	ttlOfLastBlock := ttlOfLastBlockDefault
	if !renewLastBlockOnL1 {
		ttlOfLastBlock = ttlOfLastBlockInfinity
	}
	result := l1RollupInfoProducer{
		ctx:                                  ctx,
		cancelCtx:                            cancel,
		syncStatus:                           newSyncStatus(startingBlockNumber, SyncChunkSize, ttlOfLastBlock),
		workers:                              newWorkers(ctx, ethermans),
		filterToSendOrdererResultsToConsumer: newFilterToSendOrdererResultsToConsumer(startingBlockNumber),
		outgoingChannel:                      outgoingChannel,
		statistics: l1RollupInfoProducerStatistics{
			initialBlockNumber: startingBlockNumber,
			startTime:          time.Now(),
		},
	}
	err := result.verify(false)
	if err != nil {
		log.Fatal(err)
	}
	return &result
}

// This object keep track of the statistics of the process, to be able to estimate the ETA
type l1RollupInfoProducerStatistics struct {
	initialBlockNumber  uint64
	lastBlockNumber     uint64
	numRollupInfoOk     uint64
	numRollupInfoErrors uint64
	numRetrievedBlocks  uint64
	startTime           time.Time
	lastShowUpTime      time.Time
}

func (l *l1RollupInfoProducerStatistics) getETA() string {
	numTotalOfBlocks := l.lastBlockNumber - l.initialBlockNumber
	if l.numRetrievedBlocks == 0 {
		return "N/A"
	}
	elapsedTime := time.Since(l.startTime)
	eta := time.Duration(float64(elapsedTime) / float64(l.numRetrievedBlocks) * float64(numTotalOfBlocks-l.numRetrievedBlocks))
	percent := float64(l.numRetrievedBlocks) / float64(numTotalOfBlocks) * conversionFactorPercentage
	blocks_per_seconds := float64(l.numRetrievedBlocks) / float64(elapsedTime.Seconds())
	return fmt.Sprintf("ETA: %s percent:%2.2f  blocks_per_seconds:%2.2f pending_block:%v/%v num_errors:%v",
		eta, percent, blocks_per_seconds, l.numRetrievedBlocks, numTotalOfBlocks, l.numRollupInfoErrors)
}

// TDOO: There is no min function in golang??
func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// verify: test params and status without if not allowModify avoid doing connection or modification of objects
func (l *l1RollupInfoProducer) verify(allowModify bool) error {
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

func (l *l1RollupInfoProducer) initialize() error {
	// TODO: check that all ethermans have the same chainID and get last block in L1
	err := l.verify(true)
	if err != nil {
		log.Fatal(err)
	}
	err = l.workers.initialize()
	if err != nil {
		log.Fatal(err)
	}
	if l.syncStatus.needToRenewLastBlockOnL1() {
		log.Infof("producer: Need a initial value for Last Block On L1, doing the request (maxRetries:%v, timeRequest:%v)",
			maxRetriesForRequestnitialValueOfLastBlock, timeRequestInitialValueOfLastBlock)
		result := l.retrieveInitialValueOfLastBlock(maxRetriesForRequestnitialValueOfLastBlock, timeRequestInitialValueOfLastBlock)
		if result.err != nil {
			log.Error(result.err)
			return result.err
		}
		l.onNewLastBlock(result.result.block, false)
	}
	// We don't want to start request to L1 until calling Start()

	return nil
}

// retrieveInitialValueOfLastBlock: get initial value of Last Block On L1
func (l *l1RollupInfoProducer) retrieveInitialValueOfLastBlock(maxPermittedRetries int, timeout time.Duration) genericResponse[retrieveL1LastBlockResult] {
	for {
		log.Infof("producer: Retrieving last block on L1 (remaining tries=%v, timeout=%v)", maxPermittedRetries, timeout)
		result := l.workers.requestLastBlock(l.ctx, timeout)
		if result.err == nil {
			return result
		}
		log.Info("producer: can't start request because: ", result.err)
		maxPermittedRetries--
		if maxPermittedRetries == 0 {
			log.Error("producer: exhausted retries, returning error: ", result.err)
			return result
		}
		time.Sleep(time.Second)
	}
}

// OnNewLastBlock is called when a new last block on L1 is received
func (l *l1RollupInfoProducer) onNewLastBlock(lastBlock uint64, launchWork bool) {
	resp := l.syncStatus.onNewLastBlockOnL1(lastBlock)
	l.statistics.lastBlockNumber = lastBlock
	log.Infof("producer: New last block on L1: %v -> %s", lastBlock, resp.toString())
	if launchWork {
		l.launchWork()
	}
}

// launchWork: launch new workers if possible and returns new channels created
// returns the number of workers launched
func (l *l1RollupInfoProducer) launchWork() int {
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
		log.Infof("producer: No workers launched because: %s", accDebugStr)
	}
	return launchedWorker
}

func (l *l1RollupInfoProducer) start() error {
	var waitDuration = time.Duration(0)
	for l.step(&waitDuration) {
	}
	l.workers.waitFinishAllWorkers()
	return nil
}

func (l *l1RollupInfoProducer) step(waitDuration *time.Duration) bool {
	select {
	case <-l.ctx.Done():
		return false
	// That timeout is not need, but just in case that stop launching request
	case <-time.After(*waitDuration):
		log.Infof("producer: Periodic timeout each [%s]: just in case, launching work if need", timeOutMainLoop)
		*waitDuration = timeOutMainLoop
		//log.Debugf(" ** Periodic status: %s", l.toStringBrief())
	case resultRollupInfo := <-l.workers.getResponseChannelForRollupInfo():
		l.onResponseRollupInfo(resultRollupInfo)

	case resultLastBlock := <-l.workers.getResponseChannelForLastBlock():
		l.onResponseRetrieveLastBlockOnL1(resultLastBlock)
	}
	// We check if we have finish the work
	if l.syncStatus.isNodeFullySynchronizedWithL1() {
		log.Infof("producer: we have retieve all rollupInfo from  L1")
		return false
	}
	// Try to nenew last block on L1 if needed
	l.renewLastBlockOnL1IfNeeded()
	// Try to launch retrieve more rollupInfo from L1
	l.launchWork()
	if time.Since(l.statistics.lastShowUpTime) > timeForShowUpStatisticsLog {
		log.Infof("producer: Statistics:%s", l.statistics.getETA())
		l.statistics.lastShowUpTime = time.Now()
	}
	return true
}

func (l *l1RollupInfoProducer) renewLastBlockOnL1IfNeeded() {
	if !l.syncStatus.needToRenewLastBlockOnL1() {
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err := l.workers.asyncRequestLastBlock(l.ctx)
	if err != nil {
		if err.Error() == errReachMaximumLiveRequestsOfThisType {
			log.Debugf("producer: There are a request to get last block on L1 already running")
		} else {
			log.Warnf("producer: Error while trying to get last block on L1: %v", err)
		}
	}
}

func (l *l1RollupInfoProducer) onResponseRetrieveLastBlockOnL1(result genericResponse[retrieveL1LastBlockResult]) {
	if result.err != nil {
		log.Warnf("producer: Error while trying to get last block on L1: %v", result.err)
		return
	}
	l.onNewLastBlock(result.result.block, true)
}

func (l *l1RollupInfoProducer) onResponseRollupInfo(result genericResponse[responseRollupInfoByBlockRange]) {
	isOk := (result.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
	if isOk {
		l.statistics.numRollupInfoOk++
		l.statistics.numRetrievedBlocks += result.result.blockRange.len()
		outgoingPackages := l.filterToSendOrdererResultsToConsumer.filter(*newL1SyncMessageData(result.result))
		l.sendPackages(outgoingPackages)
	} else {
		l.statistics.numRollupInfoErrors++
		log.Warnf("producer: Error while trying to get rollup info by block range: %v", result.err)
	}
}

func (l *l1RollupInfoProducer) stop() {
	l.cancelCtx()
}

func (l *l1RollupInfoProducer) sendPackages(outgoingPackages []l1SyncMessage) {
	for _, pkg := range outgoingPackages {
		log.Infof("producer: Sending results [data] to consumer:%s: It could block channel [%d/%d]", pkg.toStringBrief(), len(l.outgoingChannel), cap(l.outgoingChannel))
		l.outgoingChannel <- pkg
	}
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
