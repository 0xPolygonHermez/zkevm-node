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

// genericResponse struct contains all common data for any kind of transaction
type genericResponse struct {
	err           error
	duration      time.Duration
	typeOfRequest typeOfRequest
}

func (r *genericResponse) String() string {
	return fmt.Sprintf("typeOfRequest: [%v] duration: [%v] err: [%v]  ",
		r.typeOfRequest.String(), r.duration, r.err)
}

type responseRollupInfoByBlockRange struct {
	generic genericResponse
	result  *rollupInfoByBlockRangeResult
}

func (r *responseRollupInfoByBlockRange) toStringBrief() string {
	result := fmt.Sprintf(" generic:[%s] ",
		r.generic.String())
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
	// If there are no blocks in this range, it gets the last one
	// so it could be nil if there are no blocks.
	lastBlockOfRange *types.Block
}

func (r *rollupInfoByBlockRangeResult) toStringBrief() string {
	isLastBlockOfRangeSet := r.lastBlockOfRange != nil
	return fmt.Sprintf(" blockRange: %s len_blocks: [%d] len_order:[%d] lastBlockOfRangeSet [%t]",
		r.blockRange.String(),
		len(r.blocks), len(r.order), isLastBlockOfRangeSet)
}

type blockRange struct {
	fromBlock uint64
	toBlock   uint64
}

func (b *blockRange) String() string {
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

type workerEtherman struct {
	mutex                sync.Mutex
	etherman             EthermanInterface
	status               ethermanStatusEnum
	typeOfCurrentRequest typeOfRequest
	blockRange           blockRange
	startTime            time.Time
}

func (w *workerEtherman) String() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	timeSince := time.Since(w.startTime)
	if w.isBusyUnsafe() {
		return fmt.Sprintf("status:%s br:%s time:%s", w.status.String(), w.blockRange.String(), timeSince.Round(time.Second).String())
	}
	return fmt.Sprintf("status:%s", w.status.String())
}

func newWorker(etherman EthermanInterface) *workerEtherman {
	return &workerEtherman{etherman: etherman, status: ethermanIdle}
}

func (w *workerEtherman) asyncRequestRollupInfoByBlockRange(ctx contextWithCancel, ch chan responseRollupInfoByBlockRange, wg *sync.WaitGroup, blockRange blockRange, sleepBefore time.Duration) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.isBusyUnsafe() {
		ctx.cancel()
		if wg != nil {
			wg.Done()
		}
		return errors.New(errWorkerBusy)
	}
	w.status = ethermanWorking
	w.typeOfCurrentRequest = typeRequestRollupInfo
	w.blockRange = blockRange
	w.startTime = time.Now()
	launch := func() {
		defer ctx.cancel()
		if wg != nil {
			defer wg.Done()
		}
		if sleepBefore > 0 {
			log.Debugf("worker: RollUpInfo(%s) sleeping %s before executing...", blockRange.String(), sleepBefore)
			select {
			case <-ctx.ctx.Done():
				log.Debugf("worker: RollUpInfo(%s) cancelled in sleep", blockRange.String())
				w.setStatus(ethermanIdle)
				ch <- newResponseRollupInfo(context.Canceled, 0, typeRequestRollupInfo, &rollupInfoByBlockRangeResult{blockRange: blockRange})
				return
			case <-time.After(sleepBefore):
			}
		}
		now := time.Now()
		fromBlock := blockRange.fromBlock
		toBlock := blockRange.toBlock
		blocks, order, err := w.etherman.GetRollupInfoByBlockRange(ctx.ctx, fromBlock, &toBlock)
		var lastBlock *types.Block = nil
		if err == nil && len(blocks) == 0 {
			log.Debugf("worker: RollUpInfo(%s) calling EthBlockByNumber(%d)", blockRange.String(), toBlock)
			lastBlock, err = w.etherman.EthBlockByNumber(ctx.ctx, toBlock)
		}
		duration := time.Since(now)
		result := newResponseRollupInfo(err, duration, typeRequestRollupInfo, &rollupInfoByBlockRangeResult{blockRange, blocks, order, lastBlock})
		w.setStatus(ethermanIdle)
		if !errors.Is(err, context.Canceled) {
			log.Debugf("worker: RollUpInfo(%s) cancelled result err=%s", blockRange.String(), err)
		}
		ch <- result
	}
	go launch()
	return nil
}
func (w *workerEtherman) requestLastBlock(ctx context.Context) responseL1LastBlock {
	w.mutex.Lock()
	if w.isBusyUnsafe() {
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

func (w *workerEtherman) setStatus(status ethermanStatusEnum) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = status
	w.typeOfCurrentRequest = typeRequestNone
}

func (w *workerEtherman) isIdle() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.status == ethermanIdle
}

func (w *workerEtherman) isBusyUnsafe() bool {
	return w.status != ethermanIdle
}

func newResponseRollupInfo(err error, duration time.Duration, typeOfRequest typeOfRequest, result *rollupInfoByBlockRangeResult) responseRollupInfoByBlockRange {
	return responseRollupInfoByBlockRange{genericResponse{err, duration, typeOfRequest}, result}
}

func newResponseL1LastBlock(err error, duration time.Duration, typeOfRequest typeOfRequest, result *retrieveL1LastBlockResult) responseL1LastBlock {
	return responseL1LastBlock{genericResponse{err, duration, typeOfRequest}, result}
}
