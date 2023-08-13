package synchronizer

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

const (
	errMissingLastBlock                                             = "consumer:the received rollupinfo have no blocks and need to fill last block"
	errCanceled                                                     = "consumer:context canceled"
	numIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfoData = 20
	acceptableTimeWaitingForNewRollupInfoData                       = 1 * time.Second
)

type executionModeEnum int8

const (
	executionModeNormal               executionModeEnum = 0
	executionModeFinishOnEmptyChannel executionModeEnum = 1
)

type l1DataProcessorstatistics struct {
	numProcessedRollupInfo         uint64
	numProcessedBlocks             uint64
	startTime                      time.Time
	timePreviousProcessingDuration time.Duration
}

type synchronizerProcessBlockRangeInterface interface {
	processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error
}
type l1DataProcessor struct {
	mutex                 sync.Mutex
	synchronizer          synchronizerProcessBlockRangeInterface
	chIncommingRollupInfo chan l1PackageData
	ctx                   context.Context
	statistics            l1DataProcessorstatistics
	lastEthBlockSynced    *state.Block
	executionMode         executionModeEnum
}

func newL1DataProcessor(synchronizer synchronizerProcessBlockRangeInterface,
	ctx context.Context, ch chan l1PackageData) *l1DataProcessor {
	return &l1DataProcessor{
		synchronizer:          synchronizer,
		ctx:                   ctx,
		chIncommingRollupInfo: ch,
		executionMode:         executionModeNormal,
		statistics: l1DataProcessorstatistics{
			startTime: time.Now(),
		},
	}
}
func (l *l1DataProcessor) initialize() error {
	return nil
}

func (l *l1DataProcessor) start() error {
	err := l.step()
	for ; err == nil; err = l.step() {
	}
	return err
}
func (l *l1DataProcessor) step() error {
	timeWaitingStart := time.Now()
	var err error
	select {
	case <-l.ctx.Done():
		return errors.New(errCanceled)
	case rollupInfo := <-l.chIncommingRollupInfo:
		if rollupInfo.dataIsValid {
			err = l.processIncommingRollupInfoData(rollupInfo.data, timeWaitingStart)
		}
		if rollupInfo.ctrlIsValid {
			err = l.processIncommingRollupControlData(rollupInfo.ctrl, timeWaitingStart)
		}
	}
	return err
}
func (l *l1DataProcessor) processIncommingRollupControlData(control l1ConsumerControl, timeWaitingStart time.Time) error {
	log.Infof("consumer: processing controlPackage: %s", control.toString())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l._mustStopExecution()
}

func (l *l1DataProcessor) processIncommingRollupInfoData(rollupInfo getRollupInfoByBlockRangeResult, timeWaitingStart time.Time) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var err error
	timeWaitingEnd := time.Now()
	waitingTimeForData := timeWaitingEnd.Sub(timeWaitingStart)
	blocksPerSecond := float64(l.statistics.numProcessedBlocks) / time.Since(l.statistics.startTime).Seconds()
	if l.statistics.numProcessedRollupInfo > numIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfoData && waitingTimeForData > acceptableTimeWaitingForNewRollupInfoData {
		msg := fmt.Sprintf("wasted waiting for new rollupInfo from L1: %s last_process: %s new range: %s block_per_second: %f",
			waitingTimeForData, l.statistics.timePreviousProcessingDuration, rollupInfo.blockRange.toString(), blocksPerSecond)
		log.Warnf("consumer:: Too much wasted time:%s", msg)
	}
	// Process
	l.statistics.numProcessedRollupInfo++
	log.Infof("consumer: processing rollupInfo #%000d: range:%s num_blocks [%d] wasted_time_waiting_for_data [%s] last_process_time [%s] block_per_second [%f]", l.statistics.numProcessedRollupInfo, rollupInfo.blockRange.toString(), len(rollupInfo.blocks),
		waitingTimeForData, l.statistics.timePreviousProcessingDuration, blocksPerSecond)
	timeProcessingStart := time.Now()
	l.lastEthBlockSynced, err = l._Process(rollupInfo)
	l.statistics.timePreviousProcessingDuration = time.Since(timeProcessingStart)
	if err != nil {
		log.Error("consumer: error processing rollupInfo. Error: ", err)
		return err
	}
	l.statistics.numProcessedBlocks += uint64(len(rollupInfo.blocks))
	return l._mustStopExecution()
}

func (l *l1DataProcessor) getLastEthBlockSynced() *state.Block {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.lastEthBlockSynced
}

func (l *l1DataProcessor) finishExecutionWhenChannelIsEmpty() {
	log.Infof("consumer: Setting executionMode to executionModeFinishOnEmptyChannel (current channel len=%d)", len(l.chIncommingRollupInfo))
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.executionMode = executionModeFinishOnEmptyChannel
	log.Infof("consumer: sending a dummy result to wake up select and evaluate if must consumer finish execution")
	l.sendIgnoreResultToWakeUpSelect()
}

func (l *l1DataProcessor) sendIgnoreResultToWakeUpSelect() {
	// Send a dummy result to wake up select
	l.chIncommingRollupInfo <- *newL1PackageDataControl(actionStop)

}

func (l *l1DataProcessor) _mustStopExecution() error {
	if l.executionMode == executionModeFinishOnEmptyChannel && len(l.chIncommingRollupInfo) == 0 {
		log.Infof("consumer:  executionModeFinishOnEmptyChannel so Finishing execution because the channel is empty")
		return errors.New("executionModeFinishOnEmptyChannel")
	}
	return nil
}

func (l *l1DataProcessor) _Process(rollupInfo getRollupInfoByBlockRangeResult) (*state.Block, error) {
	blocks := rollupInfo.blocks
	order := rollupInfo.order
	err := l.synchronizer.processBlockRange(blocks, order)
	var lastEthBlockSynced *state.Block
	if err != nil {
		log.Error("consumer: Error processing block range: ", rollupInfo.blockRange, " err:", err)
		return nil, err
	}
	if len(blocks) > 0 {
		tmpStateBlock := convertEthmanBlockToStateBlock(&blocks[len(blocks)-1])
		lastEthBlockSynced = &tmpStateBlock
		logBlocks(blocks)
	}
	if len(blocks) == 0 {
		fb := rollupInfo.lastBlockOfRange
		if fb == nil {
			log.Warn("consumer: Error processing block range: ", rollupInfo.blockRange, " err: need the last block of range and got a nil")
			return nil, errors.New(errMissingLastBlock)
		}
		b := convertL1BlockToEthBlock(fb)
		err = l.synchronizer.processBlockRange([]etherman.Block{b}, order)
		if err != nil {
			log.Error("consumer: Error processing last block of range: ", rollupInfo.blockRange, " err:", err)
			return nil, err
		}
		block := convertL1BlockToStateBlock(fb)
		lastEthBlockSynced = &block
		log.Debug("consumer: Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
	}
	return lastEthBlockSynced, nil
}

func logBlocks(blocks []etherman.Block) {
	for i := range blocks {
		log.Debug("consumer: Position: [", i, "/", len(blocks), "] . BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
	}
}

func convertL1BlockToEthBlock(fb *types.Block) etherman.Block {
	return etherman.Block{
		BlockNumber: fb.NumberU64(),
		BlockHash:   fb.Hash(),
		ParentHash:  fb.ParentHash(),
		ReceivedAt:  time.Unix(int64(fb.Time()), 0),
	}
}

func convertL1BlockToStateBlock(fb *types.Block) state.Block {
	return state.Block{
		BlockNumber: fb.NumberU64(),
		BlockHash:   fb.Hash(),
		ParentHash:  fb.ParentHash(),
		ReceivedAt:  time.Unix(int64(fb.Time()), 0),
	}
}

func convertEthmanBlockToStateBlock(fb *etherman.Block) state.Block {
	return state.Block{
		BlockNumber: fb.BlockNumber,
		BlockHash:   fb.BlockHash,
		ParentHash:  fb.ParentHash,
		ReceivedAt:  fb.ReceivedAt,
	}
}
