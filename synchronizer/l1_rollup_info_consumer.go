package synchronizer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

const (
	minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData = 5
	minAcceptableTimeWaitingForNewRollupInfoData                       = 1 * time.Second
)

var (
	errMissingLastBlock                     = errors.New("consumer:the received rollupinfo have no blocks and need to fill last block")
	errContextCanceled                      = errors.New("consumer:context canceled")
	errConsumerStopped                      = errors.New("consumer:stopped by request")
	errConsumerStoppedBecauseIsSynchronized = errors.New("consumer:stopped because is synchronized")
)

type configConsumer struct {
	numIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData int
	acceptableTimeWaitingForNewRollupInfoData                       time.Duration
}

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
	statistics            l1RollupInfoConsumerStatistics
	lastEthBlockSynced    *state.Block
}

func newL1RollupInfoConsumer(cfg configConsumer,
	synchronizer synchronizerProcessBlockRangeInterface, ch chan l1SyncMessage) *l1RollupInfoConsumer {
	if cfg.acceptableTimeWaitingForNewRollupInfoData < minAcceptableTimeWaitingForNewRollupInfoData {
		log.Warnf("consumer: the acceptableTimeWaitingForNewRollupInfoData is too low (%s) minimum recommended %s", cfg.acceptableTimeWaitingForNewRollupInfoData, minAcceptableTimeWaitingForNewRollupInfoData)
	}
	if cfg.numIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData < minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData {
		log.Warnf("consumer: the numIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfoData is too low (%d) minimum recommended  %d", cfg.numIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData, minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData)
	}

	return &l1RollupInfoConsumer{
		synchronizer:          synchronizer,
		chIncommingRollupInfo: ch,
		statistics: l1RollupInfoConsumerStatistics{
			startTime: time.Now(),
			cfg:       cfg,
		},
	}
}

func (l *l1RollupInfoConsumer) Start(ctx context.Context) error {
	l.ctx = ctx
	l.statistics.onStart()
	err := l.step()
	for ; err == nil; err = l.step() {
	}
	if err != errConsumerStopped && err != errConsumerStoppedBecauseIsSynchronized {
		return err
	}
	// The errConsumerStopped is not an error, so we return nil meaning that the process finished in a normal way
	return nil
}
func (l *l1RollupInfoConsumer) step() error {
	l.statistics.onStartStep()
	var err error
	select {
	case <-l.ctx.Done():
		return errContextCanceled
	case rollupInfo := <-l.chIncommingRollupInfo:
		if rollupInfo.dataIsValid {
			err = l.processIncommingRollupInfoData(rollupInfo.data)
			if err != nil {
				log.Error("consumer: error processing package.RollupInfoData. Error: ", err)
			}
		}
		if rollupInfo.ctrlIsValid {
			err = l.processIncommingRollupControlData(rollupInfo.ctrl)
			if err != nil && !errors.Is(err, errConsumerStoppedBecauseIsSynchronized) && !errors.Is(err, errConsumerStopped) {
				log.Error("consumer: error processing package.ControlData. Error: ", err)
			}
			log.Infof("consumer: processed ControlData[%s]. Result: %s", rollupInfo.ctrl.String(), err)
		}
	}
	return err
}
func (l *l1RollupInfoConsumer) processIncommingRollupControlData(control l1ConsumerControl) error {
	log.Debugf("consumer: processing controlPackage: %s", control.String())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if control.event == eventStop {
		log.Infof("consumer: received a stop, so it stops processing. ignoring rest of items on channel len=%d", len(l.chIncommingRollupInfo))
		return errConsumerStopped
	}
	if control.event == eventProducerIsFullySynced {
		itemsInChannel := len(l.chIncommingRollupInfo)
		if itemsInChannel == 0 {
			log.Infof("consumer: received a fullSync and nothing pending in channel to process, so stopping consumer")
			return errConsumerStoppedBecauseIsSynchronized
		} else {
			log.Infof("consumer: received a fullSync but still have %d items in channel to process, so not stopping consumer", itemsInChannel)
		}
	}
	return nil
}

func (l *l1RollupInfoConsumer) processIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var err error
	statisticsMsg := l.statistics.onStartProcessIncommingRollupInfoData(rollupInfo)
	log.Infof("consumer: processing rollupInfo #%000d: range:%s num_blocks [%d] statistics:%s", l.statistics.numProcessedRollupInfo, rollupInfo.blockRange.String(), len(rollupInfo.blocks), statisticsMsg)
	timeProcessingStart := time.Now()
	l.lastEthBlockSynced, err = l.processUnsafe(rollupInfo)
	l.statistics.onFinishProcessIncommingRollupInfoData(rollupInfo, time.Since(timeProcessingStart), err)
	if err != nil {
		log.Info("consumer: error processing rollupInfo. Error: ", err)
		return err
	}
	l.statistics.numProcessedBlocks += uint64(len(rollupInfo.blocks))
	return nil
}

// GetLastEthBlockSynced returns the last block synced, if true is returned, otherwise it returns false
func (l *l1RollupInfoConsumer) GetLastEthBlockSynced() (state.Block, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.lastEthBlockSynced == nil {
		return state.Block{}, false
	}
	return *l.lastEthBlockSynced, true
}

func (l *l1RollupInfoConsumer) StopAfterProcessChannelQueue() {
	log.Infof("consumer: Sending stop package: it will stop consumer (current channel len=%d)", len(l.chIncommingRollupInfo))
	l.sendStopPackage()
}

func (l *l1RollupInfoConsumer) sendStopPackage() {
	// Send a stop to the channel to stop the consumer when reach this point
	l.chIncommingRollupInfo <- *newL1SyncMessageControl(eventStop)
}

func (l *l1RollupInfoConsumer) processUnsafe(rollupInfo rollupInfoByBlockRangeResult) (*state.Block, error) {
	blocks := rollupInfo.blocks
	order := rollupInfo.order
	var lastEthBlockSynced *state.Block
	err := l.synchronizer.processBlockRange(blocks, order)
	if err != nil {
		log.Info("consumer: Error processing block range: ", rollupInfo.blockRange, " err:", err)
		return nil, err
	}
	if len(blocks) > 0 {
		tmpStateBlock := convertEthmanBlockToStateBlock(&blocks[len(blocks)-1])
		lastEthBlockSynced = &tmpStateBlock
		logBlocks(blocks)
	}
	if len(blocks) == 0 {
		lb := rollupInfo.lastBlockOfRange
		if lb == nil {
			log.Warn("consumer: Error processing block range: ", rollupInfo.blockRange, " err: need the last block of range and got a nil")
			return nil, errMissingLastBlock
		}
		b := convertL1BlockToEthBlock(lb)
		err = l.synchronizer.processBlockRange([]etherman.Block{b}, order)
		if err != nil {
			log.Error("consumer: Error processing last block of range: ", rollupInfo.blockRange, " err:", err)
			return nil, err
		}
		block := convertL1BlockToStateBlock(lb)
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
