package synchronizer

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

const (
	errMissingLastBlock = "the received rollupinfo have no blocks and need to fill last block"
	errCanceled         = "context canceled"
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
	synchronizer          synchronizerProcessBlockRangeInterface
	chIncommingRollupInfo chan getRollupInfoByBlockRangeResult
	ctx                   context.Context
	statistics            l1DataProcessorstatistics
}

func newL1DataProcessor(synchronizer synchronizerProcessBlockRangeInterface,
	ctx context.Context, ch chan getRollupInfoByBlockRangeResult) *l1DataProcessor {
	return &l1DataProcessor{
		synchronizer:          synchronizer,
		ctx:                   ctx,
		chIncommingRollupInfo: ch,
		statistics: l1DataProcessorstatistics{
			startTime: time.Now(),
		},
	}
}
func (l *l1DataProcessor) initialize() error {
	return nil
}

func (l *l1DataProcessor) Start() error {
	err := l.step()
	for ; err == nil; err = l.step() {
	}
	return err
}
func (l *l1DataProcessor) step() error {
	timeWaitingStart := time.Now()
	select {
	case <-l.ctx.Done():
		return errors.New(errCanceled)
	case rollupInfo := <-l.chIncommingRollupInfo:
		timeWaitingEnd := time.Now()
		log.Debugf("Time wasted waiting for new rollupInfo from L1: %s last_process: %s new range: %s block_per_second: %f",
			timeWaitingEnd.Sub(timeWaitingStart), l.statistics.timePreviousProcessingDuration, rollupInfo.blockRange.toString(),
			float64(l.statistics.numProcessedBlocks)/time.Since(l.statistics.startTime).Seconds())
		// Process
		l.statistics.numProcessedRollupInfo++
		log.Infof("Processing rollupInfo [%000d]: range:%s num_blocks [%d]", l.statistics.numProcessedRollupInfo, rollupInfo.blockRange.toString(), len(rollupInfo.blocks))
		timeProcessingStart := time.Now()
		_, err := l.Process(rollupInfo)
		l.statistics.timePreviousProcessingDuration = time.Since(timeProcessingStart)
		if err != nil {
			log.Error("error processing rollupInfo. Error: ", err)
			return err
		}
		l.statistics.numProcessedBlocks += uint64(len(rollupInfo.blocks))
	}
	return nil
}

func (l *l1DataProcessor) Process(rollupInfo getRollupInfoByBlockRangeResult) (*state.Block, error) {
	blocks := rollupInfo.blocks
	order := rollupInfo.order
	err := l.synchronizer.processBlockRange(blocks, order)
	var lastEthBlockSynced *state.Block
	if err != nil {
		log.Error("Error processing block range: ", rollupInfo.blockRange, " err:", err)
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
			log.Warn("Error processing block range: ", rollupInfo.blockRange, " err: need the last block of range and got a nil")
			return nil, errors.New(errMissingLastBlock)
		}
		b := convertL1BlockToEthBlock(fb)
		err = l.synchronizer.processBlockRange([]etherman.Block{b}, order)
		if err != nil {
			log.Error("Error processing last block of range: ", rollupInfo.blockRange, " err:", err)
			return nil, err
		}
		block := convertL1BlockToStateBlock(fb)
		lastEthBlockSynced = &block
		log.Debug("Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
	}
	return lastEthBlockSynced, nil
}

func logBlocks(blocks []etherman.Block) {
	for i := range blocks {
		log.Debug("Position: [", i, "/", len(blocks), "] . BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
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
