package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	ttlOfLastBlockDefault      = time.Second * 60
	timeOutMainLoop            = time.Minute * 5
	timeForShowUpStatisticsLog = time.Second * 60
	conversionFactorPercentage = 100
)

// - multiples etherman to do it in parallel
// - generate blocks to be retrieved
// - retrieve blocks (parallel)
// - check last block on L1?

type l1DataRetrieverStatistics struct {
	initialBlockNumber  uint64
	lastBlockNumber     uint64
	numRollupInfoOk     uint64
	numRollupInfoErrors uint64
	numRetrievedBlocks  uint64
	startTime           time.Time
	lastShowUpTime      time.Time
}

func (l *l1DataRetrieverStatistics) getETA() string {
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

type l1DataRetriever struct {
	mutex      sync.Mutex
	ctx        context.Context
	cancelCtx  context.CancelFunc
	workers    workers
	syncStatus syncStatusInterface
	// Send the block info to this channel ordered by block number
	//channel chan getRollupInfoByBlockRangeResult
	sender     *filterToSendOrdererResultsToConsumer
	statistics l1DataRetrieverStatistics
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// verify: test params and status without if not allowModify avoid doing connection or modification of objects
func (l *l1DataRetriever) verify(allowModify bool) error {
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

func newL1DataRetriever(ctx context.Context, ethermans []EthermanInterface,
	startingBlockNumber uint64, SyncChunkSize uint64,
	resultChannel chan l1PackageData, renewLastBlockOnL1 bool) *l1DataRetriever {
	if cap(resultChannel) < len(ethermans) {
		log.Warnf("resultChannel must have a capacity (%d) of at least equal to number of ether clients (%d)", cap(resultChannel), len(ethermans))
	}
	ctx, cancel := context.WithCancel(ctx)
	ttlOfLastBlock := ttlOfLastBlockDefault
	if !renewLastBlockOnL1 {
		ttlOfLastBlock = ttlOfLastBlockInfinity
	}
	result := l1DataRetriever{
		ctx:        ctx,
		cancelCtx:  cancel,
		syncStatus: newSyncStatus(startingBlockNumber, SyncChunkSize, ttlOfLastBlock),
		workers:    newWorkers(ethermans),
		sender:     newFilterToSendOrdererResultsToConsumer(resultChannel, startingBlockNumber),
		statistics: l1DataRetrieverStatistics{
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

func (l *l1DataRetriever) initialize() error {
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

// retrieveInitialValueOfLastBlock do a synchronous request to get the last block on L1
// that is needed to start the synchronizer
func (l *l1DataRetriever) retrieveInitialValueOfLastBlock() genericResponse[retrieveL1LastBlockResult] {
	maxPermittedRetries := 10
	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		ch, err := l.workers.asyncRequestLastBlock(ctx)
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
func (l *l1DataRetriever) onNewLastBlock(lastBlock uint64, launchWork bool) {
	resp := l.syncStatus.onNewLastBlockOnL1(lastBlock)
	l.statistics.lastBlockNumber = lastBlock
	log.Infof("New last block on L1: %v -> %s", lastBlock, resp.toString())
	if launchWork {
		l.launchWork()
	}
}

// launchWork: launch new workers if possible and returns new channels created
// returns true if new workers were launched
func (l *l1DataRetriever) launchWork() int {
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

func (l *l1DataRetriever) start() error {
	var waitDuration = time.Duration(0)
	for l.step(&waitDuration) {
	}
	l.workers.waitFinishAllWorkers()
	return nil
}

func (l *l1DataRetriever) step(waitDuration *time.Duration) bool {
	select {
	case <-l.ctx.Done():
		return false
	// That timeout is not need, but just in case that stop launching request
	case <-time.After(*waitDuration):
		log.Infof("Periodic timeout each [%s]: just in case, launching work if need", timeOutMainLoop)
		*waitDuration = timeOutMainLoop
		//log.Debugf(" ** Periodic status: %s", l.toStringBrief())
	case resultRollupInfo := <-l.workers.getResponseChannelForRollupInfo():
		l.onResponseRollupInfo(resultRollupInfo)

	case resultLastBlock := <-l.workers.getResponseChannelForLastBlock():
		l.onResponseRetrieveLastBlockOnL1(resultLastBlock)
	}
	// We check if we have finish the work
	if l.syncStatus.isNodeFullySynchronizedWithL1() {
		log.Infof("we have retieve all rollupInfo from  L1")
		return false
	}
	// Try to nenew last block on L1 if needed
	l.renewLastBlockOnL1IfNeeded()
	// Try to launch retrieve more rollupInfo from L1
	l.launchWork()
	if time.Since(l.statistics.lastShowUpTime) > timeForShowUpStatisticsLog {
		log.Infof("Statistics:%s", l.statistics.getETA())
		l.statistics.lastShowUpTime = time.Now()
	}
	return true
}

func (l *l1DataRetriever) renewLastBlockOnL1IfNeeded() {
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

func (l *l1DataRetriever) onResponseRetrieveLastBlockOnL1(result genericResponse[retrieveL1LastBlockResult]) {
	if result.err != nil {
		log.Warnf("Error while trying to get last block on L1: %v", result.err)
		return
	}
	l.onNewLastBlock(result.result.block, true)
}

func (l *l1DataRetriever) onResponseRollupInfo(result genericResponse[getRollupInfoByBlockRangeResult]) {
	isOk := (result.err == nil)
	l.syncStatus.onFinishWorker(result.result.blockRange, isOk)
	if isOk {
		l.statistics.numRollupInfoOk++
		l.statistics.numRetrievedBlocks += result.result.blockRange.len()
		l.sender.addResultAndSendToConsumer(newL1PackageDataFromResult(result.result))
	} else {
		l.statistics.numRollupInfoErrors++
		log.Warnf("Error while trying to get rollup info by block range: %v", result.err)
	}
}

func (l *l1DataRetriever) stop() {
	l.cancelCtx()
}

// https://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go
