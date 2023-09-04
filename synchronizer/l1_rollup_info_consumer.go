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

	errConsumerStopped = "consumer:stopped by request"
)

// synchronizerProcessBlockRangeInterface is the interface with synchronizer
// to execute blocks. This interface is used to mock the synchronizer in the tests
type synchronizerProcessBlockRangeInterface interface {
	processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error
}

// l1RollupInfoConsumer is the object that process the rollup info data incomming from channel chIncommingRollupInfo
type l1RollupInfoConsumer struct {
	mutex                 sync.Mutex
	synchronizer          synchronizerProcessBlockRangeInterface
	chIncommingRollupInfo chan l1SyncMessage
	ctx                   context.Context
	statistics            ll1RollupInfoConsumerStatistics
	lastEthBlockSynced    *state.Block
}

type ll1RollupInfoConsumerStatistics struct {
	numProcessedRollupInfo         uint64
	numProcessedBlocks             uint64
	startTime                      time.Time
	timePreviousProcessingDuration time.Duration
}

func newL1RollupInfoConsumer(synchronizer synchronizerProcessBlockRangeInterface,
	ctx context.Context, ch chan l1SyncMessage) *l1RollupInfoConsumer {
	return &l1RollupInfoConsumer{
		synchronizer:          synchronizer,
		ctx:                   ctx,
		chIncommingRollupInfo: ch,
		statistics: ll1RollupInfoConsumerStatistics{
			startTime: time.Now(),
		},
	}
}

func (l *l1RollupInfoConsumer) start() error {
	err := l.step()
	for ; err == nil; err = l.step() {
	}
	if err.Error() != errConsumerStopped {
		return err
	}
	// The errConsumerStopped is not an error, so we return nil meaning that the process finished in a normal way
	return nil
}
func (l *l1RollupInfoConsumer) step() error {
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
func (l *l1RollupInfoConsumer) processIncommingRollupControlData(control l1ConsumerControl, timeWaitingStart time.Time) error {
	log.Infof("consumer: processing controlPackage: %s", control.toString())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if control.action == actionStop {
		return errors.New(errConsumerStopped)
	}
	return nil
}

func (l *l1RollupInfoConsumer) processIncommingRollupInfoData(rollupInfo responseRollupInfoByBlockRange, timeWaitingStart time.Time) error {
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
	return nil
}

// getLastEthBlockSynced returns the last block synced, if true is returned, otherwise it returns false
func (l *l1RollupInfoConsumer) getLastEthBlockSynced() (state.Block, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.lastEthBlockSynced == nil {
		return state.Block{}, false
	}
	return *l.lastEthBlockSynced, true
}

func (l *l1RollupInfoConsumer) stopAfterProcessChannelQueue() {
	log.Infof("consumer: Sending stop package: it will stop consumer (current channel len=%d)", len(l.chIncommingRollupInfo))
	l.sendStopPackage()
}

func (l *l1RollupInfoConsumer) sendStopPackage() {
	// Send a dummy result to wake up select
	l.chIncommingRollupInfo <- *newL1SyncMessageControl(actionStop)
}

func (l *l1RollupInfoConsumer) _Process(rollupInfo responseRollupInfoByBlockRange) (*state.Block, error) {
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
