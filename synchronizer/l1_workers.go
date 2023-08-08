package synchronizer

import (
	"context"
	"errors"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	errAllWorkersBusy = "all workers are busy"
)

type workers struct {
	mutex   sync.Mutex
	workers []*worker
	// Aggregated channels
	chAggregatedRollupInfo chan genericResponse[getRollupInfoByBlockRangeResult]
	chAggregatedLastBlock  chan genericResponse[retrieveL1LastBlockResult]
}

func newWorkers(ethermans []ethermanInterface) workers {
	result := workers{}
	result.workers = make([]*worker, len(ethermans))
	for i, etherman := range ethermans {
		worker := newWorker(etherman)
		result.workers[i] = worker
	}
	result.chAggregatedRollupInfo = make(chan genericResponse[getRollupInfoByBlockRangeResult], len(ethermans))
	result.chAggregatedLastBlock = make(chan genericResponse[retrieveL1LastBlockResult], len(ethermans))
	aggregateChannelsGeneric[genericResponse[getRollupInfoByBlockRangeResult]](&result.chAggregatedRollupInfo, result.getAllChannelsRollupInfo())
	aggregateChannelsGeneric[genericResponse[retrieveL1LastBlockResult]](&result.chAggregatedLastBlock, result.getAllChannelsLastBlock())
	return result
}

func (w *workers) getAllChannelsLastBlock() []chan genericResponse[retrieveL1LastBlockResult] {
	result := make([]chan genericResponse[retrieveL1LastBlockResult], len(w.workers))
	for i, worker := range w.workers {
		result[i] = worker.chLastBlock
	}
	return result
}

func (w *workers) getAllChannelsRollupInfo() []chan genericResponse[getRollupInfoByBlockRangeResult] {
	result := make([]chan genericResponse[getRollupInfoByBlockRangeResult], len(w.workers))
	for i, worker := range w.workers {
		result[i] = worker.chRollupInfo
	}
	return result
}

func (w *workers) verifyDry() error {
	if len(w.workers) == 0 {
		return errors.New("no etherman provided")
	}
	return nil
}

func (this *workers) verify() error {
	// TODO: checks that all ethermans have the same chainID
	//verifyChainIDOfEthermans()
	return nil
}

func (this *workers) initialize() error {
	return nil
}

func (w *workers) asyncRetrieveLastBlock(ctx context.Context) (chan genericResponse[retrieveL1LastBlockResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	worker := w._getIdleWorker()
	if worker == nil {
		return nil, errors.New(errAllWorkersBusy)
	}
	return worker.asyncRetrieveLastBlock(ctx)
}

func (w *workers) asyncGetRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[getRollupInfoByBlockRangeResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	worker := w._getIdleWorker()
	if worker == nil {
		return nil, errors.New(errAllWorkersBusy)
	}
	ch, err := worker.asyncGetRollupInfoByBlockRange(ctx, blockRange)
	if err == nil {
		log.Infof("worker GetRollupInfoByBlockRange launcher for block range %v", blockRange)
	}
	return ch, err
}

func (w *workers) _getIdleWorker() *worker {
	for _, worker := range w.workers {
		if worker.isIdle() {
			return worker
		}
	}
	return nil
}

// https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
func aggregateChannelsGeneric[T any](aggChannel *chan T, newChannels []chan T) {
	for _, ch := range newChannels {
		go func(ch chan T) {
			for msg := range ch {
				*aggChannel <- msg
			}
		}(ch)
	}
}
