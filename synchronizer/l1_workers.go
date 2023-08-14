package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	noLimitLiveRequests                   = 0
	errRequiredEtherman                   = "required etherman"
	errAllWorkersBusy                     = "all workers are busy"
	errReachMaximumLiveRequestsOfThisType = "reach maximum live requests of this type"
)

type workers struct {
	ctx     context.Context
	mutex   sync.Mutex
	workers []*worker
	// Channels to send to outside the responses from worker
	chOutgoingRollupInfo chan genericResponse[responseRollupInfoByBlockRange]
	chOutgoingLastBlock  chan genericResponse[retrieveL1LastBlockResult]

	// Channels that receive the responses from worker
	chIncommingRollupInfo chan genericResponse[responseRollupInfoByBlockRange]
	chIncomingLastBlock   chan genericResponse[retrieveL1LastBlockResult]
	// It need a goroutine that listen in chIncomming and send to chOutgoing
	launchedGoRoutineToRouteResponses bool

	waitGroups        [typeRequestEOF]sync.WaitGroup
	limitLiveRequests [typeRequestEOF]uint64
}

func newWorkers(ctx context.Context, ethermans []EthermanInterface) *workers {
	result := workers{ctx: ctx,
		chIncommingRollupInfo:             make(chan genericResponse[responseRollupInfoByBlockRange], len(ethermans)+1),
		chIncomingLastBlock:               make(chan genericResponse[retrieveL1LastBlockResult], len(ethermans)+1),
		launchedGoRoutineToRouteResponses: false,
	}

	result.limitLiveRequests[typeRequestRollupInfo] = noLimitLiveRequests
	result.limitLiveRequests[typeRequestLastBlock] = 1

	result.workers = make([]*worker, len(ethermans))
	for i, etherman := range ethermans {
		result.workers[i] = newWorker(etherman)
	}
	result.chOutgoingRollupInfo = make(chan genericResponse[responseRollupInfoByBlockRange], len(ethermans)+1)
	result.chOutgoingLastBlock = make(chan genericResponse[retrieveL1LastBlockResult], len(ethermans)+1)
	return &result
}
func (w *workers) verify(allowModify bool) error {
	if len(w.workers) == 0 {
		return errors.New(errRequiredEtherman)
	}
	// TODO: checks that all ethermans have the same chainID
	//verifyChainIDOfEthermans()
	return nil
}

func (w *workers) initialize() error {
	return nil
}

func (w *workers) finalize() error {
	return nil
}

func (w *workers) getResponseChannelForRollupInfo() chan genericResponse[responseRollupInfoByBlockRange] {
	return w.chOutgoingRollupInfo
}

func (w *workers) getResponseChannelForLastBlock() chan genericResponse[retrieveL1LastBlockResult] {
	return w.chOutgoingLastBlock
}

func (w *workers) asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[responseRollupInfoByBlockRange], error) {
	requestStrForDebug := fmt.Sprintf("GetRollupInfoByBlockRange(%s)", blockRange.toString())
	f := func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error {
		res := worker.asyncRequestRollupInfoByBlockRange(ctx, w.getResponseChannelForRollupInfo(), wg, blockRange)
		return res
	}
	res := w.asyncGenericRequest(ctx, typeRequestRollupInfo, requestStrForDebug, f)
	return w.chOutgoingRollupInfo, res
}

// asyncRequestLastBlock launches a request to retrieve the last block of the L1 chain.
func (w *workers) asyncRequestLastBlock(ctx context.Context) (chan genericResponse[retrieveL1LastBlockResult], error) {
	requestStrForDebug := "RetrieveLastBlock()"
	f := func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error {
		res := worker.asyncRequestLastBlock(ctx, w.getResponseChannelForLastBlock(), wg)
		return res
	}
	res := w.asyncGenericRequest(ctx, typeRequestRollupInfo, requestStrForDebug, f)
	return w.chOutgoingLastBlock, res
}

func (w *workers) requestLastBlock(ctx context.Context, timeout time.Duration) genericResponse[retrieveL1LastBlockResult] {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	w.mutex.Lock()
	defer w.mutex.Unlock()
	worker := w._getIdleWorker()
	if worker == nil {
		log.Debugf("workers: call:[%s] failed err:%s", "requestLastBlock", errAllWorkersBusy)
		return genericResponse[retrieveL1LastBlockResult]{err: errors.New(errAllWorkersBusy), typeOfRequest: typeRequestLastBlock}
	}
	var err error

	ch := make(chan genericResponse[retrieveL1LastBlockResult], 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	err = worker.asyncRequestLastBlock(ctxTimeout, ch, &wg)
	if err == nil {
		wg.Wait()
		result := <-ch
		return result
	} else {
		return genericResponse[retrieveL1LastBlockResult]{err: err, typeOfRequest: typeRequestLastBlock}
	}
}

// asyncGenericRequest launches a generic request to the workers
func (w *workers) asyncGenericRequest(ctx context.Context, requestType typeOfRequest, requestStrForDebug string,
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
	w._launchGoroutineForRoutingResponsesIfNeed()
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

func (w *workers) _launchGoroutineForRoutingResponsesIfNeed() {
	if w.launchedGoRoutineToRouteResponses {
		return
	}
	log.Infof("workers: launching goroutine to route responses")
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case resultRollupInfo := <-w.chIncommingRollupInfo:
				w.onResponseRollupInfo(resultRollupInfo)
			case resultLastBlock := <-w.chIncomingLastBlock:
				w.onResponseLastBlock(resultLastBlock)
			}
		}
	}()

	w.launchedGoRoutineToRouteResponses = true
}

func (w *workers) onResponseRollupInfo(v genericResponse[responseRollupInfoByBlockRange]) {
	msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
	if v.err == nil {
		msg += fmt.Sprintf(" block_range:%s", v.result.blockRange.toString())
	}
	log.Infof(msg)
	w.chOutgoingRollupInfo <- v
}
func (w *workers) onResponseLastBlock(v genericResponse[retrieveL1LastBlockResult]) {
	msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
	if v.err == nil {
		msg += fmt.Sprintf(" block_number:%v", v.result.block)
	}
	log.Infof(msg)
	w.chOutgoingLastBlock <- v
}

func (w *workers) waitFinishAllWorkers() {
	for i := range w.waitGroups {
		wg := &w.waitGroups[i]
		wg.Wait()
	}
}

func (w *workers) _getIdleWorker() *worker {
	for _, worker := range w.workers {
		if worker.isIdle() {
			return worker
		}
	}
	return nil
}

func (w *workers) _checkReachedLimitLiveRequest(typeOfRequest typeOfRequest) bool {
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

func (w *workers) _countLiveRequestOfType(typeOfRequest typeOfRequest) uint64 {
	var n uint64 = 0
	for _, worker := range w.workers {
		if worker.typeOfCurrentRequest == typeOfRequest {
			n++
		}
	}
	return n
}

// func (w *workersImpl) _onIncommmingResponse(msg interface{}) {
// 	switch v := msg.(type) {
// 	case genericResponse[getRollupInfoByBlockRangeResult]:
// 		msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
// 		if v.err == nil {
// 			msg += fmt.Sprintf(" block_range:%s", v.result.blockRange.toString())
// 		}
// 		log.Infof(msg)
// 	case genericResponse[retrieveL1LastBlockResult]:
// 		msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
// 		if v.err == nil {
// 			msg += fmt.Sprintf(" block_number:%v", v.result.block)
// 		}
// 		log.Infof(msg)
// 	}
// }

// func (w *workersImpl) _getAllChannelsLastBlock() []chan genericResponse[retrieveL1LastBlockResult] {
// 	result := make([]chan genericResponse[retrieveL1LastBlockResult], len(w.workers))
// 	for i, worker := range w.workers {
// 		result[i] = worker.chLastBlock
// 	}
// 	return result
// }

// func (w *workersImpl) _getAllChannelsRollupInfo() []chan genericResponse[getRollupInfoByBlockRangeResult] {
// 	result := make([]chan genericResponse[getRollupInfoByBlockRangeResult], len(w.workers))
// 	for i, worker := range w.workers {
// 		result[i] = worker.chRollupInfo
// 	}
// 	return result
// }

// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
// func aggregateChannelsGeneric[T any](aggChannel *chan T, newChannels []chan T, w *workersImpl) {
// 	for _, ch := range newChannels {
// 		go func(ch chan T) {
// 			for msg := range ch {
// 				if w != nil {
// 					w._onIncommmingResponse(msg)
// 				}
// 				*aggChannel <- msg
// 			}
// 		}(ch)
// 	}
// }
