package synchronizer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
)

type ethermanStatusEnum int8

const (
	ethermanIdle    ethermanStatusEnum = 0
	ethermanWorking ethermanStatusEnum = 1
	ethermanError   ethermanStatusEnum = 2
)

const (
	errWorkerBusy = "worker is busy"
)

type genericResponse[T any] struct {
	err      error
	duration time.Duration
	result   *T
}

type blockRange struct {
	fromBlock uint64
	toBlock   uint64
}

type getRollupInfoByBlockRangeResult struct {
	blockRange blockRange
	blocks     []etherman.Block
	order      map[common.Hash][]etherman.Order
}

type retrieveL1LastBlockResult struct {
	block uint64
}

type worker struct {
	mutex    sync.Mutex
	etherman ethermanInterface
	status   ethermanStatusEnum
	// channels
	chRollupInfo chan genericResponse[getRollupInfoByBlockRangeResult]
	chLastBlock  chan genericResponse[retrieveL1LastBlockResult]
}

func newWorker(etherman ethermanInterface) *worker {
	return &worker{etherman: etherman, status: ethermanIdle,
		chRollupInfo: make(chan genericResponse[getRollupInfoByBlockRangeResult], 1),
		chLastBlock:  make(chan genericResponse[retrieveL1LastBlockResult], 1)}
}

func newGenericAnswer[T any](err error, duration time.Duration, result *T) genericResponse[T] {
	return genericResponse[T]{err, duration, result}
}

func (w *worker) asyncGetRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan genericResponse[getRollupInfoByBlockRangeResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._isBusy() {
		return nil, errors.New(errWorkerBusy)
	}
	ch := w.chRollupInfo
	w.status = ethermanWorking
	launch := func() {
		now := time.Now()
		blocks, order, err := w.etherman.GetRollupInfoByBlockRange(ctx, blockRange.fromBlock, &blockRange.toBlock)
		duration := time.Since(now)
		result := newGenericAnswer(err, duration, &getRollupInfoByBlockRangeResult{blockRange, blocks, order})
		ch <- result
		w.setStatus(ethermanIdle)
	}
	go launch()
	return ch, nil
}

func (w *worker) asyncRetrieveLastBlock(ctx context.Context) (chan genericResponse[retrieveL1LastBlockResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._isBusy() {
		return nil, errors.New(errWorkerBusy)
	}
	ch := w.chLastBlock
	w.status = ethermanWorking
	launch := func() {

		now := time.Now()
		header, err := w.etherman.HeaderByNumber(ctx, nil)
		duration := time.Since(now)
		result := newGenericAnswer(err, duration, &retrieveL1LastBlockResult{header.Number.Uint64()})
		ch <- result
		w.setStatus(ethermanIdle)
	}
	go launch()
	return ch, nil
}

func (w *worker) setStatus(status ethermanStatusEnum) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = status
}

func (w *worker) _isBusy() bool {
	return w.status != ethermanIdle
}

func (w *worker) isIdle() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w._isIdle()
}

func (w *worker) _isIdle() bool {
	return w.status == ethermanIdle
}
