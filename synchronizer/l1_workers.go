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
	// Channel to send to outside the responses from worker
	chOutgoingRollupInfo chan responseRollupInfoByBlockRange

	// Channel that receive the responses from worker
	chIncommingRollupInfo chan responseRollupInfoByBlockRange

	// It need a goroutine that listen in chIncomming and send to chOutgoing
	launchedGoRoutineToRouteResponses bool

	waitGroups        [typeRequestEOF]sync.WaitGroup
	limitLiveRequests [typeRequestEOF]uint64
}

func (w *workers) toString() string {
	result := fmt.Sprintf("workers: limitLiveRequests: %v  ch_out:%d ch_in_worker:%d ", w.limitLiveRequests, len(w.chOutgoingRollupInfo), len(w.chIncommingRollupInfo))
	for i := range w.workers {
		result += fmt.Sprintf(" worker[%d]: %s", i, w.workers[i].toString())
	}
	return result
}

func newWorkers(ctx context.Context, ethermans []EthermanInterface) *workers {
	result := workers{ctx: ctx,
		chIncommingRollupInfo:             make(chan responseRollupInfoByBlockRange, len(ethermans)+1),
		launchedGoRoutineToRouteResponses: false,
	}

	result.limitLiveRequests[typeRequestRollupInfo] = noLimitLiveRequests
	result.limitLiveRequests[typeRequestLastBlock] = 1

	result.workers = make([]*worker, len(ethermans))
	for i, etherman := range ethermans {
		result.workers[i] = newWorker(etherman)
	}
	result.chOutgoingRollupInfo = make(chan responseRollupInfoByBlockRange, len(ethermans)+1)
	return &result
}

func (w *workers) initialize() error {
	if len(w.workers) == 0 {
		return errors.New(errRequiredEtherman)
	}
	return nil
}

func (w *workers) stop() {
	// TODO: ctx cancel
	for i := 0; i < len(w.waitGroups); i++ {
		w.waitGroups[i].Wait()
	}
}

func (w *workers) getResponseChannelForRollupInfo() chan responseRollupInfoByBlockRange {
	return w.chOutgoingRollupInfo
}

func (w *workers) asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan responseRollupInfoByBlockRange, error) {
	requestStrForDebug := fmt.Sprintf("GetRollupInfoByBlockRange(%s)", blockRange.toString())
	f := func(worker *worker, ctx context.Context, wg *sync.WaitGroup) error {
		res := worker.asyncRequestRollupInfoByBlockRange(ctx, w.getResponseChannelForRollupInfo(), wg, blockRange)
		return res
	}
	res := w.asyncGenericRequest(ctx, typeRequestRollupInfo, requestStrForDebug, f)
	return w.chOutgoingRollupInfo, res
}

func (w *workers) requestLastBlockWithRetries(ctx context.Context, timeout time.Duration, maxPermittedRetries int) responseL1LastBlock {
	for {
		log.Debugf("workers: Retrieving last block on L1 (remaining tries=%v, timeout=%v)", maxPermittedRetries, timeout)
		result := w.requestLastBlock(ctx, timeout)
		if result.generic.err == nil {
			return result
		}
		maxPermittedRetries--
		log.Debugf("workers: fail request pending retries:%d : err:%s ", maxPermittedRetries, result.generic.err)
		if maxPermittedRetries == 0 {
			log.Error("workers: exhausted retries for last block on L1, returning error: ", result.generic.err)
			return result
		}
		time.Sleep(time.Second)
	}
}

func (w *workers) requestLastBlock(ctx context.Context, timeout time.Duration) responseL1LastBlock {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	w.mutex.Lock()
	defer w.mutex.Unlock()
	worker := w._getIdleWorker()
	if worker == nil {
		log.Debugf("workers: call:[%s] failed err:%s", "requestLastBlock", errAllWorkersBusy)
		return newResponseL1LastBlock(errors.New(errAllWorkersBusy), time.Duration(0), typeRequestLastBlock, nil)
	}
	result := worker.requestLastBlock(ctxTimeout)
	return result
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
		log.Debugf("workers: worker started call:[%s]", requestStrForDebug)
	} else {
		log.Debugf("workers: worker started failed! call:[%s] failed err:[%s]", requestStrForDebug, err.Error())
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
			}
		}
	}()

	w.launchedGoRoutineToRouteResponses = true
}

func (w *workers) onResponseRollupInfo(v responseRollupInfoByBlockRange) {
	msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
	log.Infof(msg)
	w.chOutgoingRollupInfo <- v
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
		log.Debugf("workers: reached limit live request of type [%d] currentWorkes:%d >= maxPermitted:%d", typeOfRequest, numberOfWorkers, maximumLiveRequests)
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
