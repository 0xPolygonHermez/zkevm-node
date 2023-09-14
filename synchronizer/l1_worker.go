package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

type ethermanStatusEnum int8

const (
	ethermanIdle    ethermanStatusEnum = 0
	ethermanWorking ethermanStatusEnum = 1
	ethermanError   ethermanStatusEnum = 2
)

func (s ethermanStatusEnum) String() string {
	return [...]string{"idle", "working", "error"}[s]
}

type typeOfRequest int8

const (
	typeRequestNone       typeOfRequest = 0
	typeRequestRollupInfo typeOfRequest = 1
	typeRequestLastBlock  typeOfRequest = 2
	typeRequestEOF        typeOfRequest = 3
)

func (s typeOfRequest) String() string {
	return [...]string{"none", "rollup", "lastBlock", "EOF"}[s]
}

const (
	errWorkerBusy = "worker is busy"
)

// genericResponse struct containts all common data for any kind of transaction
type genericResponse struct {
	err           error
	duration      time.Duration
	typeOfRequest typeOfRequest
}

func (r *genericResponse) toStringBrief() string {
	return fmt.Sprintf("typeOfRequest: [%v] duration: [%v] err: [%v]  ",
		r.typeOfRequest.String(), r.duration, r.err)
}

type responseRollupInfoByBlockRange struct {
	generic genericResponse
	result  *rollupInfoByBlockRangeResult
}

func (r *responseRollupInfoByBlockRange) toStringBrief() string {
	result := fmt.Sprintf(" generic:[%s] ",
		r.generic.toStringBrief())
	if r.result != nil {
		result += fmt.Sprintf(" result:[%s]", r.result.toStringBrief())
	} else {
		result += " result:[nil]"
	}
	return result
}

type rollupInfoByBlockRangeResult struct {
	blockRange blockRange
	blocks     []etherman.Block
	order      map[common.Hash][]etherman.Order
	// If there are no blocks in this range get get the last one
	// so it could be nil if there are blocks.
	lastBlockOfRange *types.Block
}

func (r *rollupInfoByBlockRangeResult) toStringBrief() string {
	isLastBlockOfRangeSet := r.lastBlockOfRange != nil
	return fmt.Sprintf(" blockRange: %s len_blocks: [%d] len_order:[%d] lastBlockOfRangeSet [%t]",
		r.blockRange.toString(),
		len(r.blocks), len(r.order), isLastBlockOfRangeSet)
}

type blockRange struct {
	fromBlock uint64
	toBlock   uint64
}

func (b *blockRange) toString() string {
	return fmt.Sprintf("[%v, %v]", b.fromBlock, b.toBlock)
}

func (b *blockRange) len() uint64 {
	return b.toBlock - b.fromBlock + 1
}

type responseL1LastBlock struct {
	generic genericResponse
	result  *retrieveL1LastBlockResult
}

type retrieveL1LastBlockResult struct {
	block uint64
}

type worker struct {
	mutex                sync.Mutex
	etherman             EthermanInterface
	status               ethermanStatusEnum
	typeOfCurrentRequest typeOfRequest
}

func (w *worker) toString() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return fmt.Sprintf("status:%s req:%s", w.status.String(), w.typeOfCurrentRequest.String())
}

func newWorker(etherman EthermanInterface) *worker {
	return &worker{etherman: etherman, status: ethermanIdle}
}

func (w *worker) asyncRequestRollupInfoByBlockRange(ctx context.Context, ch chan responseRollupInfoByBlockRange, wg *sync.WaitGroup, blockRange blockRange) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w._isBusy() {
		return errors.New(errWorkerBusy)
	}
	w.status = ethermanWorking
	w.typeOfCurrentRequest = typeRequestRollupInfo
	launch := func() {
		if wg != nil {
			defer wg.Done()
		}
		now := time.Now()
		fromBlock := blockRange.fromBlock
		toBlock := blockRange.toBlock
		blocks, order, err := w.etherman.GetRollupInfoByBlockRange(ctx, fromBlock, &toBlock)
		var lastBlock *types.Block = nil
		if err == nil && len(blocks) == 0 {
			log.Debugf("worker: calling EthBlockByNumber(%v)", toBlock)
			lastBlock, err = w.etherman.EthBlockByNumber(ctx, toBlock)
		}
		duration := time.Since(now)
		result := newResponseRollupInfo(err, duration, typeRequestRollupInfo, &rollupInfoByBlockRangeResult{blockRange, blocks, order, lastBlock})
		w.setStatus(ethermanIdle)
		ch <- result
	}
	go launch()
	return nil
}
func (w *worker) requestLastBlock(ctx context.Context) responseL1LastBlock {
	w.mutex.Lock()
	if w._isBusy() {
		w.mutex.Unlock()
		return newResponseL1LastBlock(errors.New(errWorkerBusy), time.Duration(0), typeRequestLastBlock, nil)
	}
	w.status = ethermanWorking
	w.typeOfCurrentRequest = typeRequestLastBlock
	w.mutex.Unlock()
	now := time.Now()
	header, err := w.etherman.HeaderByNumber(ctx, nil)
	duration := time.Since(now)
	var result responseL1LastBlock
	if err == nil {
		result = newResponseL1LastBlock(err, duration, typeRequestLastBlock, &retrieveL1LastBlockResult{header.Number.Uint64()})
	} else {
		result = newResponseL1LastBlock(err, duration, typeRequestLastBlock, nil)
	}
	w.setStatus(ethermanIdle)
	return result
}

func (w *worker) setStatus(status ethermanStatusEnum) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = status
	w.typeOfCurrentRequest = typeRequestNone
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

func newResponseRollupInfo(err error, duration time.Duration, typeOfRequest typeOfRequest, result *rollupInfoByBlockRangeResult) responseRollupInfoByBlockRange {
	return responseRollupInfoByBlockRange{genericResponse{err, duration, typeOfRequest}, result}
}

func newResponseL1LastBlock(err error, duration time.Duration, typeOfRequest typeOfRequest, result *retrieveL1LastBlockResult) responseL1LastBlock {
	return responseL1LastBlock{genericResponse{err, duration, typeOfRequest}, result}
}
