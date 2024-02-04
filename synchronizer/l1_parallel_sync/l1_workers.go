package l1_parallel_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	noSleepTime              = time.Duration(0)
	minimumNumberOfEthermans = 2
)

var (
	errAllWorkersBusy   = errors.New("all workers are busy")
	errRequiredEtherman = errors.New("required etherman")
)

// worker: is the expected functions of a worker
type worker interface {
	String() string
	asyncRequestRollupInfoByBlockRange(ctx contextWithCancel, ch chan responseRollupInfoByBlockRange, wg *sync.WaitGroup, request requestRollupInfoByBlockRange) error
	requestLastBlock(ctx context.Context) responseL1LastBlock
	isIdle() bool
}

type workersConfig struct {
	timeoutRollupInfo time.Duration
}

type workerData struct {
	worker worker
	ctx    contextWithCancel
}

func (w *workerData) String() string {
	return fmt.Sprintf("worker:%s ctx:%v", w.worker.String(), w.ctx)
}

type workers struct {
	mutex sync.Mutex
	// worker for asking lastBlock on L1 (to avoid that all of them are busy)
	workerForLastBlock workerData
	workers            []workerData
	// Channel to send to outside the responses from worker | workers --> client
	chOutgoingRollupInfo chan responseRollupInfoByBlockRange

	// Channel that receive the responses from worker  | worker --> workers
	chIncommingRollupInfo chan responseRollupInfoByBlockRange

	waitGroups [typeRequestEOF]sync.WaitGroup

	cfg workersConfig
}

func (w *workers) String() string {
	result := fmt.Sprintf("num_workers:%d ch[%d,%d] ", len(w.workers), len(w.chOutgoingRollupInfo), len(w.chIncommingRollupInfo))
	for i := range w.workers {
		if !w.workers[i].worker.isIdle() {
			result += fmt.Sprintf(" worker[%d]: %s", i, w.workers[i].worker.String())
		}
	}
	return result
}

func (w *workers) ToStringBrief() string {
	return fmt.Sprintf(" working: %d of %d ", w.howManyRunningWorkers(), len(w.workers))
}

func newWorkers(ethermans []L1ParallelEthermanInterface, cfg workersConfig) *workers {
	result := workers{chIncommingRollupInfo: make(chan responseRollupInfoByBlockRange, len(ethermans)+1),
		cfg: cfg}
	if (len(ethermans)) < minimumNumberOfEthermans {
		log.Fatalf("workers: at least %d ethermans are required, got %d", minimumNumberOfEthermans, len(ethermans))
	}
	workers := make([]workerData, len(ethermans))
	for i, etherman := range ethermans {
		workers[i].worker = newWorker(etherman)
	}
	result.workers = workers[1:]
	result.workerForLastBlock = workers[0]
	result.chOutgoingRollupInfo = make(chan responseRollupInfoByBlockRange, len(ethermans)+1)
	return &result
}

func (w *workers) initialize() error {
	if len(w.workers) == 0 {
		return errRequiredEtherman
	}
	return nil
}

func (w *workers) howManyRunningWorkers() int {
	result := 0
	for _, worker := range w.workers {
		if !worker.worker.isIdle() {
			result++
		}
	}
	return result
}

func (w *workers) stop() {
	log.Infof("workers: stopping workers %s", w.String())
	for i := range w.workers {
		wd := &w.workers[i]
		if !wd.worker.isIdle() {
			log.Debugf("workers: stopping worker[%d] %s", i, wd.String())
		}
		wd.ctx.cancel()
	}
	//w.waitFinishAllWorkers()
}

func (w *workers) getResponseChannelForRollupInfo() chan responseRollupInfoByBlockRange {
	return w.chOutgoingRollupInfo
}

func (w *workers) asyncRequestRollupInfoByBlockRange(ctx context.Context, request requestRollupInfoByBlockRange) (chan responseRollupInfoByBlockRange, error) {
	requestStrForDebug := fmt.Sprintf("GetRollupInfoByBlockRange(%s)", request.String())
	f := func(worker worker, ctx contextWithCancel, wg *sync.WaitGroup) error {
		res := worker.asyncRequestRollupInfoByBlockRange(ctx, w.getResponseChannelForRollupInfo(), wg, request)
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
	ctxTimeout := newContextWithTimeout(ctx, timeout)
	defer ctxTimeout.cancel()
	w.mutex.Lock()
	defer w.mutex.Unlock()
	//workerIndex, worker := w.getIdleWorkerUnsafe()
	worker := &w.workerForLastBlock
	if worker == nil {
		log.Debugf("workers: call:[%s] failed err:%s", "requestLastBlock", errAllWorkersBusy)
		return newResponseL1LastBlock(errAllWorkersBusy, time.Duration(0), typeRequestLastBlock, nil)
	}
	worker.ctx = ctxTimeout

	log.Debugf("workers: worker : launching requestLatBlock (timeout=%s)", timeout.String())
	result := worker.worker.requestLastBlock(ctxTimeout.ctx)
	return result
}

// asyncGenericRequest launches a generic request to the workers
func (w *workers) asyncGenericRequest(ctx context.Context, requestType typeOfRequest, requestStrForDebug string,
	funcRequest func(worker worker, ctx contextWithCancel, wg *sync.WaitGroup) error) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	workerIndex, worker := w.getIdleWorkerUnsafe()
	if worker == nil {
		log.Debugf("workers: call:[%s] failed err:%s", requestStrForDebug, errAllWorkersBusy)
		return errAllWorkersBusy
	}
	ctxWithCancel := newContextWithTimeout(ctx, w.cfg.timeoutRollupInfo)
	w.workers[workerIndex].ctx = ctxWithCancel
	w.launchGoroutineForRoutingResponse(ctxWithCancel.ctx, workerIndex)
	wg := &w.waitGroups[requestType]
	wg.Add(1)

	err := funcRequest(worker, ctxWithCancel, wg)
	if err == nil {
		log.Debugf("workers: worker[%d] started call:[%s]", workerIndex, requestStrForDebug)
	} else {
		log.Debugf("workers: worker[%d] started failed! call:[%s] failed err:[%s]", workerIndex, requestStrForDebug, err.Error())
	}
	return err
}

func (w *workers) launchGoroutineForRoutingResponse(ctx context.Context, workerIndex int) {
	log.Debugf("workers: launching goroutine to route response for worker[%d]", workerIndex)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resultRollupInfo := <-w.chIncommingRollupInfo:
				w.onResponseRollupInfo(resultRollupInfo)
			}
		}
	}()
}

func (w *workers) onResponseRollupInfo(v responseRollupInfoByBlockRange) {
	msg := fmt.Sprintf("workers: worker finished:[ %s ]", v.toStringBrief())
	log.Infof(msg)
	w.chOutgoingRollupInfo <- v
}

func (w *workers) waitFinishAllWorkers() {
	for i := 0; i < len(w.waitGroups); i++ {
		log.Debugf("workers: waiting for waitGroup[%d]", i)
		w.waitGroups[i].Wait()
	}
}

func (w *workers) getIdleWorkerUnsafe() (int, worker) {
	for idx, worker := range w.workers {
		if worker.worker.isIdle() {
			return idx, worker.worker
		}
	}
	return -1, nil
}
