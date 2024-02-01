package l1_parallel_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
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

type requestLastBlockMode int32

const (
	requestLastBlockModeNone               requestLastBlockMode = 0
	requestLastBlockModeIfNoBlocksInAnswer requestLastBlockMode = 1
	requestLastBlockModeAlways             requestLastBlockMode = 2
)

func (s requestLastBlockMode) String() string {
	return [...]string{"none", "ifNoBlocksInAnswer", "always"}[s]
}

type requestRollupInfoByBlockRange struct {
	blockRange                         blockRange
	sleepBefore                        time.Duration
	requestLastBlockIfNoBlocksInAnswer requestLastBlockMode
	requestPreviousBlock               bool
}

func (r *requestRollupInfoByBlockRange) String() string {
	return fmt.Sprintf("blockRange: %s sleepBefore: %s lastBlock: %s prevBlock:%t",
		r.blockRange.String(), r.sleepBefore, r.requestLastBlockIfNoBlocksInAnswer.String(), r.requestPreviousBlock)
}

func (r *responseRollupInfoByBlockRange) getHighestBlockNumberInResponse() uint64 {
	if r.result == nil {
		return invalidBlockNumber
	}
	return r.result.getHighestBlockNumberInResponse()
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
	lastBlockOfRange     *types.Block
	previousBlockOfRange *types.Block
}

func (r *rollupInfoByBlockRangeResult) toStringBrief() string {
	isLastBlockOfRangeSet := r.lastBlockOfRange != nil
	ispreviousBlockOfRange := r.previousBlockOfRange != nil
	return fmt.Sprintf(" blockRange: %s len_blocks: [%d] len_order:[%d] lastBlockOfRangeSet [%t] previousBlockSet [%t]",
		r.blockRange.String(),
		len(r.blocks), len(r.order), isLastBlockOfRangeSet, ispreviousBlockOfRange)
}

func (r *rollupInfoByBlockRangeResult) getRealHighestBlockNumberInResponse() uint64 {
	highest := invalidBlockNumber
	for _, block := range r.blocks {
		if block.BlockNumber > highest {
			highest = block.BlockNumber
		}
	}
	if r.lastBlockOfRange != nil && r.lastBlockOfRange.Number().Uint64() > highest {
		highest = r.lastBlockOfRange.Number().Uint64()
	}
	return highest
}

// getHighestBlockNumberInResponse returns the highest block number in the response if toBlock or the real one if latestBlockNumber
func (r *rollupInfoByBlockRangeResult) getHighestBlockNumberInResponse() uint64 {
	if r.blockRange.toBlock != latestBlockNumber {
		return r.blockRange.toBlock
	} else {
		highestBlock := r.getRealHighestBlockNumberInResponse()
		if highestBlock == invalidBlockNumber {
			return r.blockRange.fromBlock - 1
		}
		return highestBlock
	}
}

func (r *rollupInfoByBlockRangeResult) getHighestBlockReceived() *state.Block {
	var highest *state.Block = nil
	if r.lastBlockOfRange != nil {
		stateBlock := convertL1BlockToStateBlock(r.lastBlockOfRange)
		return &stateBlock
	}
	for _, block := range r.blocks {
		if highest == nil || block.BlockNumber > highest.BlockNumber {
			blockCopy := block
			stateBlock := convertEthmanBlockToStateBlock(&blockCopy)
			highest = &stateBlock
		}
	}
	return highest
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
	etherman             L1ParallelEthermanInterface
	status               ethermanStatusEnum
	typeOfCurrentRequest typeOfRequest
	request              requestRollupInfoByBlockRange
	startTime            time.Time
}

func (w *workerEtherman) String() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	timeSince := time.Since(w.startTime)
	if w.isBusyUnsafe() {
		return fmt.Sprintf("status:%s br:%s time:%s", w.status.String(), w.request.String(), timeSince.Round(time.Second).String())
	}
	return fmt.Sprintf("status:%s", w.status.String())
}

func newWorker(etherman L1ParallelEthermanInterface) *workerEtherman {
	return &workerEtherman{etherman: etherman, status: ethermanIdle}
}

// sleep returns false if must stop execution
func (w *workerEtherman) sleep(ctx contextWithCancel, ch chan responseRollupInfoByBlockRange, request requestRollupInfoByBlockRange) bool {
	if request.sleepBefore > 0 {
		log.Debugf("worker: RollUpInfo(%s) sleeping %s before executing...", request.blockRange.String(), request.sleepBefore)
		select {
		case <-ctx.ctx.Done():
			log.Debugf("worker: RollUpInfo(%s) cancelled in sleep", request.blockRange.String())
			w.setStatus(ethermanIdle)
			ch <- newResponseRollupInfo(context.Canceled, 0, typeRequestRollupInfo, &rollupInfoByBlockRangeResult{blockRange: request.blockRange})
			return false
		case <-time.After(request.sleepBefore):
		}
	}
	return true
}

func mustRequestLastBlock(mode requestLastBlockMode, lenBlocks int, lastBlockRequest uint64) bool {
	switch mode {
	case requestLastBlockModeNone:
		return false
	case requestLastBlockModeIfNoBlocksInAnswer:
		return lenBlocks == 0 && lastBlockRequest != latestBlockNumber
	case requestLastBlockModeAlways:
		return lastBlockRequest != latestBlockNumber
	default:
		return lastBlockRequest != latestBlockNumber
	}
}

// The order of the request are important:
//
//	 The previous and last block are used to guarantee that the blocks belongs to the same chain.
//	 Check next example:
//	    Request1: LAST(200) Rollup(100-200) PREVIOUS(99)
//		       Request2: LAST(300) Rollup(201-300) PREVIOUS(200)
//	             Request3: LAST(400) Rollup(301-400) PREVIOUS(300)
//
// If there are a reorg in Request2:
//
//	Request2: [P1] LAST(300) [P2] Rollup(201-300) [P3] PREVIOUS(200) [P4]
//
// P1: PREVIOUS(200) are not going to match with the same in Request1 LAST(200)
// P2: PREVIOUS(200) are not going to match with the same in Request1 LAST(200)
// P3: PREVIOUS(200) are not going to match with the same in Request1 LAST(200)
// P4: LAST(300) are not going to match with Request3 PREVIOUS(300)
//
// In case of Rollup(100-latest):
// 	 Request1:  -----  Rollup(100..)[B120]  PREVIOUS(99)
//    													 Request2: ----- Rollup(121..)[B122]  PREVIOUS(120)
// Works in the same way

func (w *workerEtherman) asyncRequestRollupInfoByBlockRange(ctx contextWithCancel, ch chan responseRollupInfoByBlockRange, wg *sync.WaitGroup, request requestRollupInfoByBlockRange) error {
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
	w.request = request
	w.startTime = time.Now()
	launch := func() {
		defer ctx.cancel()
		if wg != nil {
			defer wg.Done()
		}
		if !w.sleep(ctx, ch, request) {
			return
		}

		// Uncomment these lines to respond with a nil result to generate fast responses (just for develop!)
		//w.setStatus(ethermanIdle)
		//ch <- newResponseRollupInfo(nil, time.Second, typeRequestRollupInfo, &rollupInfoByBlockRangeResult{blockRange, nil, nil, nil})

		now := time.Now()
		data, err := w.executeRequestRollupInfoByBlockRange(ctx, ch, request)
		duration := time.Since(now)
		result := newResponseRollupInfo(err, duration, typeRequestRollupInfo, data)
		w.setStatus(ethermanIdle)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Debugf("worker: RollUpInfo(%s) result err=%s", request.blockRange.String(), err.Error())
		}
		ch <- result
	}
	go launch()
	return nil
}

func (w *workerEtherman) executeRequestRollupInfoByBlockRange(ctx contextWithCancel, ch chan responseRollupInfoByBlockRange, request requestRollupInfoByBlockRange) (*rollupInfoByBlockRangeResult, error) {
	resultRollupInfo := rollupInfoByBlockRangeResult{request.blockRange, nil, nil, nil, nil}
	if err := w.fillLastBlock(&resultRollupInfo, ctx, request, false); err != nil {
		return &resultRollupInfo, err
	}
	if err := w.fillRollup(&resultRollupInfo, ctx, request); err != nil {
		return &resultRollupInfo, err
	}
	if err := w.fillLastBlock(&resultRollupInfo, ctx, request, true); err != nil {
		return &resultRollupInfo, err
	}
	if err := w.fillPreviousBlock(&resultRollupInfo, ctx, request); err != nil {
		return &resultRollupInfo, err
	}
	return &resultRollupInfo, nil
}

func (w *workerEtherman) fillPreviousBlock(result *rollupInfoByBlockRangeResult, ctx contextWithCancel, request requestRollupInfoByBlockRange) error {
	if request.requestPreviousBlock && request.blockRange.fromBlock > 2 {
		log.Debugf("worker: RollUpInfo(%s) request previousBlock calling EthBlockByNumber(%d)", request.blockRange.String(), request.blockRange.fromBlock)
		var err error
		result.previousBlockOfRange, err = w.etherman.EthBlockByNumber(ctx.ctx, request.blockRange.fromBlock-1)
		return err
	}
	return nil
}

func (w *workerEtherman) fillRollup(result *rollupInfoByBlockRangeResult, ctx contextWithCancel, request requestRollupInfoByBlockRange) error {
	var toBlock *uint64 = nil
	// If latest we send a nil
	if request.blockRange.toBlock != latestBlockNumber {
		toBlock = &request.blockRange.toBlock
	}
	var err error
	result.blocks, result.order, err = w.etherman.GetRollupInfoByBlockRange(ctx.ctx, request.blockRange.fromBlock, toBlock)
	if err != nil {
		return err
	}
	return nil
}

func (w *workerEtherman) fillLastBlock(result *rollupInfoByBlockRangeResult, ctx contextWithCancel, request requestRollupInfoByBlockRange, haveExecutedRollupInfo bool) error {
	if result.lastBlockOfRange != nil {
		return nil
	}
	lenBlocks := len(result.blocks)
	if !haveExecutedRollupInfo {
		lenBlocks = -1
	}
	if mustRequestLastBlock(request.requestLastBlockIfNoBlocksInAnswer, lenBlocks, request.blockRange.toBlock) {
		log.Debugf("worker: RollUpInfo(%s) request lastBlock calling EthBlockByNumber(%d) (before rollup) ", request.blockRange.String(), request.blockRange.toBlock)
		lastBlock, err := w.etherman.EthBlockByNumber(ctx.ctx, request.blockRange.toBlock)
		if err != nil {
			return err
		}
		result.lastBlockOfRange = lastBlock
	}
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
