package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

type ethermanStatusEnum int8

const (
	ethermanIdle    ethermanStatusEnum = 0
	ethermanWorking ethermanStatusEnum = 1
	ethermanError   ethermanStatusEnum = 2
)

type typeOfRequest int8

const (
	typeRequestNone       typeOfRequest = 0
	typeRequestRollupInfo typeOfRequest = 1
	typeRequestLastBlock  typeOfRequest = 2
	typeRequestEOF        typeOfRequest = 3
)

func (t typeOfRequest) toString() string {
	switch t {
	case typeRequestNone:
		return "typeRequestNone"
	case typeRequestRollupInfo:
		return "typeRequestRollupInfo"
	case typeRequestLastBlock:
		return "typeRequestLastBlock"
	case typeRequestEOF:
		return "typeRequestEOF"
	default:
		return "unknown"
	}
}

const (
	errWorkerBusy = "worker is busy"
)

type genericResponse[T any] struct {
	err           error
	duration      time.Duration
	typeOfRequest typeOfRequest
	result        *T
}

func (r *genericResponse[T]) toStringBrief() string {
	return fmt.Sprintf("typeOfRequest: [%v] duration: [%v] err: [%v]  ",
		r.typeOfRequest.toString(), r.duration, r.err)
}

type blockRange struct {
	fromBlock uint64
	toBlock   uint64
}

type getRollupInfoByBlockRangeResult struct {
	blockRange blockRange
	blocks     []etherman.Block
	order      map[common.Hash][]etherman.Order
	// If there are no blocks in this range get get the last one
	// so it could be nil if there are blocks.
	lastBlockOfRange *types.Block
	// If true the consumer will ignore this package
	ignoreThisResult bool
}

func (r *getRollupInfoByBlockRangeResult) toStringBrief() string {
	isLastBlockOfRangeSet := r.lastBlockOfRange != nil
	return fmt.Sprintf(" blockRange: %s len_blocks: [%d] len_order:[%d] lastBlockOfRangeSet [%t]",
		r.blockRange.toString(),
		len(r.blocks), len(r.order), isLastBlockOfRangeSet)
}

func (b *blockRange) toString() string {
	return fmt.Sprintf("[%v, %v]", b.fromBlock, b.toBlock)
}

func (b *blockRange) len() uint64 {
	return b.toBlock - b.fromBlock + 1
}

type retrieveL1LastBlockResult struct {
	block uint64
}

type worker struct {
	mutex         sync.Mutex
	etherman      EthermanInterface
	status        ethermanStatusEnum
	typeOfRequest typeOfRequest
	// channels
	chRollupInfo chan genericResponse[getRollupInfoByBlockRangeResult]
	chLastBlock  chan genericResponse[retrieveL1LastBlockResult]
}

func newWorker(etherman EthermanInterface) *worker {
	return &worker{etherman: etherman, status: ethermanIdle,
		chRollupInfo: make(chan genericResponse[getRollupInfoByBlockRangeResult], 1),
		chLastBlock:  make(chan genericResponse[retrieveL1LastBlockResult], 1)}
}

func (w *worker) asyncRequestRollupInfoByBlockRange(ctx context.Context, wg *sync.WaitGroup, blockRange blockRange) (chan genericResponse[getRollupInfoByBlockRangeResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._isBusy() {
		return nil, errors.New(errWorkerBusy)
	}
	ch := w.chRollupInfo
	w.status = ethermanWorking
	w.typeOfRequest = typeRequestRollupInfo
	launch := func() {
		if wg != nil {
			defer wg.Done()
		}
		now := time.Now()
		blocks, order, err := w.etherman.GetRollupInfoByBlockRange(ctx, blockRange.fromBlock, &blockRange.toBlock)
		var lastBlock *types.Block = nil
		if err == nil && len(blocks) == 0 {
			lastBlock, err = w.etherman.EthBlockByNumber(ctx, blockRange.toBlock)
		}
		duration := time.Since(now)
		result := newGenericAnswer(err, duration, typeRequestRollupInfo, &getRollupInfoByBlockRangeResult{blockRange, blocks, order, lastBlock, false})
		w.setStatus(ethermanIdle)
		ch <- result
	}
	go launch()
	return ch, nil
}

func (w *worker) asyncRequestLastBlock(ctx context.Context, wg *sync.WaitGroup) (chan genericResponse[retrieveL1LastBlockResult], error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._isBusy() {
		return nil, errors.New(errWorkerBusy)
	}
	ch := w.chLastBlock
	w.status = ethermanWorking
	w.typeOfRequest = typeRequestLastBlock
	launch := func() {
		if wg != nil {
			defer wg.Done()
		}
		now := time.Now()
		header, err := w.etherman.HeaderByNumber(ctx, nil)
		duration := time.Since(now)
		var result genericResponse[retrieveL1LastBlockResult]
		if err == nil {
			result = newGenericAnswer(err, duration, typeRequestLastBlock, &retrieveL1LastBlockResult{header.Number.Uint64()})
		} else {
			result = newGenericAnswer[retrieveL1LastBlockResult](err, duration, typeRequestLastBlock, nil)
		}
		w.setStatus(ethermanIdle)
		ch <- result
	}
	go launch()
	return ch, nil
}

func (w *worker) setStatus(status ethermanStatusEnum) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = status
	w.typeOfRequest = typeRequestNone
}

func (w *worker) isIdle() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w._isIdle()
}

func (w *worker) _isBusy() bool {
	return w.status != ethermanIdle
}

func (w *worker) _isIdle() bool {
	return w.status == ethermanIdle
}

func newGenericAnswer[T any](err error, duration time.Duration, typeOfRequest typeOfRequest, result *T) genericResponse[T] {
	return genericResponse[T]{err, duration, typeOfRequest, result}
}
