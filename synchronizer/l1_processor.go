package synchronizer

import (
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	types "github.com/ethereum/go-ethereum/core/types"
)

const (
	errMissingLastBlock = "the received rollupinfo have no blocks and need to fill last block"
)

type L1DataProcessor struct {
	synchronizer *ClientSynchronizer
}

func NewL1DataProcessor(synchronizer *ClientSynchronizer) *L1DataProcessor {
	return &L1DataProcessor{
		synchronizer: synchronizer,
	}
}

func (l *L1DataProcessor) Process(rollupInfo getRollupInfoByBlockRangeResult) (*state.Block, error) {
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
