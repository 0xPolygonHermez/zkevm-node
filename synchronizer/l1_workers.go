package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	noLimitLiveRequests                   = 0
	errRequiredEtherman                   = "required etherman"
	errAllWorkersBusy                     = "all workers are busy"
	errReachMaximumLiveRequestsOfThisType = "reach maximum live requests of this type"
)

type workers interface {
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
	asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[getRollupInfoByBlockRangeResult], error)

	getResponseChannelForLastBlock() chan genericResponse[retrieveL1LastBlockResult]
	getResponseChannelForRollupInfo() chan genericResponse[getRollupInfoByBlockRangeResult]
}

type workersImpl struct {
	mutex   sync.Mutex
	workers []*worker
	// Aggregated channels
	chAggregatedRollupInfo chan genericResponse[getRollupInfoByBlockRangeResult]
	chAggregatedLastBlock  chan genericResponse[retrieveL1LastBlockResult]

	waitGroups        [typeRequestEOF]sync.WaitGroup
	limitLiveRequests [typeRequestEOF]uint64
}

func newWorkers(ethermans []EthermanInterface) *workersImpl {
	result := workersImpl{}
	result.limitLiveRequests[typeRequestRollupInfo] = noLimitLiveRequests
	result.limitLiveRequests[typeRequestLastBlock] = 1

	result.workers = make([]*worker, len(ethermans))
	for i, etherman := range ethermans {
		result.workers[i] = newWorker(etherman)
	}
	result.chAggregatedRollupInfo = make(chan genericResponse[getRollupInfoByBlockRangeResult], len(ethermans)+1)
	result.chAggregatedLastBlock = make(chan genericResponse[retrieveL1LastBlockResult], len(ethermans)+1)
	aggregateChannelsGeneric[genericResponse[getRollupInfoByBlockRangeResult]](&result.chAggregatedRollupInfo, result._getAllChannelsRollupInfo(), &result)
	aggregateChannelsGeneric[genericResponse[retrieveL1LastBlockResult]](&result.chAggregatedLastBlock, result._getAllChannelsLastBlock(), &result)
	return &result
}
func (w *workersImpl) verify(allowModify bool) error {
	if len(w.workers) == 0 {
		return errors.New(errRequiredEtherman)
	}
	// TODO: checks that all ethermans have the same chainID
	//verifyChainIDOfEthermans()
	return nil
}

func (w *workersImpl) initialize() error {
	return nil
}

func (w *workersImpl) finalize() error {
	return nil
}

func (w *workersImpl) getResponseChannelForRollupInfo() chan genericResponse[getRollupInfoByBlockRangeResult] {
	return w.chAggregatedRollupInfo
}

func (w *workersImpl) getResponseChannelForLastBlock() chan genericResponse[retrieveL1LastBlockResult] {
	return w.chAggregatedLastBlock
}

func (w *workersImpl) asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[getRollupInfoByBlockRangeResult], error) {
	requestStrForDebug := fmt.Sprintf("GetRollupInfoByBlockRange(%s)", blockRange.toString())
	f := func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error {
		_, res := worker.asyncRequestRollupInfoByBlockRange(ctx, wg, blockRange)
		return res
	}
	res := w.asyncGenericRequest(ctx, typeRequestRollupInfo, requestStrForDebug, f)
	return w.chAggregatedRollupInfo, res
}

// asyncRequestLastBlock launches a request to retrieve the last block of the L1 chain.
func (w *workersImpl) asyncRequestLastBlock(ctx context.Context) (chan genericResponse[retrieveL1LastBlockResult], error) {
	requestStrForDebug := "RetrieveLastBlock()"
	f := func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error {
		_, res := worker.asyncRequestLastBlock(ctx, wg)
		return res
	}
	res := w.asyncGenericRequest(ctx, typeRequestRollupInfo, requestStrForDebug, f)
	return w.chAggregatedLastBlock, res
}

// asyncGenericRequest launches a generic request to the workers
func (w *workersImpl) asyncGenericRequest(ctx context.Context, requestType typeOfRequest, requestStrForDebug string,
	funcRequest func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._checkReachedLimitLiveRequest(requestType) {
		log.Debugf("workers: call:[%s] failed err:%s", requestStrForDebug, errReachMaximumLiveRequestsOfThisType)
		return errors.New(errReachMaximumLiveRequestsOfThisType)
	}
	worker := w._getIdleWorker()
	if worker == nil {
		log.Debugf("workers: call:[%s] failed err:%s", requestStrForDebug, errAllWorkersBusy)
		return errors.New(errAllWorkersBusy)
	}
	wg := &w.waitGroups[requestType]
	wg.Add(1)

	err := funcRequest(worker, ctx, wg)
	if err == nil {
		log.Infof("workers: worker started call:[%s]", requestStrForDebug)
	} else {
		log.Warnf("workers: worker started failed! call:[%s] failed err:[%s]", requestStrForDebug, err.Error())
	}
	return err
}

func (w *workersImpl) waitFinishAllWorkers() {
	for i := range w.waitGroups {
		wg := &w.waitGroups[i]
		wg.Wait()
	}
}

func (w *workersImpl) _getIdleWorker() *worker {
	for _, worker := range w.workers {
		if worker.isIdle() {
			return worker
		}
	}
	return nil
}

func (w *workersImpl) _checkReachedLimitLiveRequest(typeOfRequest typeOfRequest) bool {
	if w.limitLiveRequests[typeOfRequest] == noLimitLiveRequests {
		return false
	}
	numberOfWorkers := w._countLiveRequestOfType(typeOfRequest)
	maximumLiveRequests := w.limitLiveRequests[typeOfRequest]
	reachedLimit := numberOfWorkers >= maximumLiveRequests
	if reachedLimit {
		log.Infof("workers: reached limit live request of type [%d] currentWorkes:%d >= maxPermitted:%d", typeOfRequest, numberOfWorkers, maximumLiveRequests)
	}
	return reachedLimit
}

func (w *workersImpl) _countLiveRequestOfType(typeOfRequest typeOfRequest) uint64 {
	var n uint64 = 0
	for _, worker := range w.workers {
		if worker.typeOfRequest == typeOfRequest {
			n++
		}
	}
	return n
}

func (w *workersImpl) _onIncommmingResponse(msg interface{}) {
	switch v := msg.(type) {
	case genericResponse[getRollupInfoByBlockRangeResult]:
		msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
		if v.err == nil {
			msg += fmt.Sprintf(" block_range:%s", v.result.blockRange.toString())
		}
		log.Infof(msg)
	case genericResponse[retrieveL1LastBlockResult]:
		msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
		if v.err == nil {
			msg += fmt.Sprintf(" block_number:%v", v.result.block)
		}
		log.Infof(msg)
	}
}

func (w *workersImpl) _getAllChannelsLastBlock() []chan genericResponse[retrieveL1LastBlockResult] {
	result := make([]chan genericResponse[retrieveL1LastBlockResult], len(w.workers))
	for i, worker := range w.workers {
		result[i] = worker.chLastBlock
	}
	return result
}

func (w *workersImpl) _getAllChannelsRollupInfo() []chan genericResponse[getRollupInfoByBlockRangeResult] {
	result := make([]chan genericResponse[getRollupInfoByBlockRangeResult], len(w.workers))
	for i, worker := range w.workers {
		result[i] = worker.chRollupInfo
	}
	return result
}

// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
func aggregateChannelsGeneric[T any](aggChannel *chan T, newChannels []chan T, w *workersImpl) {
	for _, ch := range newChannels {
		go func(ch chan T) {
			for msg := range ch {
				if w != nil {
					w._onIncommmingResponse(msg)
				}
				*aggChannel <- msg
			}
		}(ch)
	}
}
