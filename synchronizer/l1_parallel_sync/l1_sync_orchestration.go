package l1_parallel_sync

import (
	"context"
	"errors"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

/*
This object is used to coordinate the producer and the consumer process.
*/
type l1RollupProducerInterface interface {
	// Start launch a new process to retrieve data from L1
	Start(ctx context.Context) error
	// Stop cancel current process
	Stop()
	// Abort execution
	Abort()
	// Reset set a new starting point and cancel current process if any
	Reset(startingBlockNumber uint64)
}

type l1RollupConsumerInterface interface {
	Start(ctx context.Context, lastEthBlockSynced *state.Block) error
	StopAfterProcessChannelQueue()
	GetLastEthBlockSynced() (state.Block, bool)
	// Reset set a new starting point
	Reset(startingBlockNumber uint64)
}

// L1SyncOrchestration is the object that coordinates the producer and the consumer process.
type L1SyncOrchestration struct {
	mutex    sync.Mutex
	producer l1RollupProducerInterface
	consumer l1RollupConsumerInterface
	// Producer is running?
	producerRunning bool
	consumerRunning bool
	// The orchestrator is running?
	isRunning     bool
	wg            sync.WaitGroup
	chProducer    chan error
	chConsumer    chan error
	ctxParent     context.Context
	ctxWithCancel contextWithCancel
}

const (
	errMissingLastEthBlockSynced = "orchestration: missing last eth block synced"
)

// NewL1SyncOrchestration create a new L1 sync orchestration object
func NewL1SyncOrchestration(ctx context.Context, producer l1RollupProducerInterface, consumer l1RollupConsumerInterface) *L1SyncOrchestration {
	res := L1SyncOrchestration{
		producer:        producer,
		consumer:        consumer,
		producerRunning: false,
		consumerRunning: false,
		chProducer:      make(chan error, 1),
		chConsumer:      make(chan error, 1),
		ctxParent:       ctx,
	}
	res.ctxWithCancel.createWithCancel(ctx)
	return &res
}

// Reset set a new starting point and cancel current process if any
func (l *L1SyncOrchestration) Reset(startingBlockNumber uint64) {
	log.Warnf("orchestration: Reset L1 sync process to blockNumber %d", startingBlockNumber)
	if l.isRunning {
		log.Infof("orchestration: reset(%d) is going to reset producer", startingBlockNumber)
	}
	l.consumer.Reset(startingBlockNumber)
	l.producer.Reset(startingBlockNumber)
	// If orchestrator is running then producer is going to be started by orchestrate() select  function when detects that producer has finished
}

// Start launch a new process to retrieve and execute data from L1
func (l *L1SyncOrchestration) Start(lastEthBlockSynced *state.Block) (*state.Block, error) {
	l.isRunning = true
	l.launchProducer(l.ctxWithCancel.ctx, lastEthBlockSynced, l.chProducer, &l.wg)
	l.launchConsumer(l.ctxWithCancel.ctx, lastEthBlockSynced, l.chConsumer, &l.wg)
	return l.orchestrate(l.ctxParent, &l.wg, l.chProducer, l.chConsumer)
}

// Abort stop inmediatly the current process
func (l *L1SyncOrchestration) Abort() {
	l.producer.Abort()
	l.ctxWithCancel.cancel()
	l.wg.Wait()
	l.ctxWithCancel.createWithCancel(l.ctxParent)
}

// IsProducerRunning return true if producer is running
func (l *L1SyncOrchestration) IsProducerRunning() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.producerRunning
}

func (l *L1SyncOrchestration) launchProducer(ctx context.Context, lastEthBlockSynced *state.Block, chProducer chan error, wg *sync.WaitGroup) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if !l.producerRunning {
		if wg != nil {
			wg.Add(1)
		}
		log.Infof("orchestration: producer is not running. Resetting the state to start from  block %v (last on DB)", lastEthBlockSynced.BlockNumber)
		l.producer.Reset(lastEthBlockSynced.BlockNumber)
		// Start producer: L1DataRetriever from L1
		l.producerRunning = true

		go func() {
			if wg != nil {
				defer wg.Done()
			}
			log.Infof("orchestration: starting producer")
			err := l.producer.Start(ctx)
			if err != nil {
				log.Warnf("orchestration: producer error . Error: %s", err)
			}
			l.mutex.Lock()
			l.producerRunning = false
			l.mutex.Unlock()
			log.Infof("orchestration: producer finished")
			chProducer <- err
		}()
	}
}

func (l *L1SyncOrchestration) launchConsumer(ctx context.Context, lastEthBlockSynced *state.Block, chConsumer chan error, wg *sync.WaitGroup) {
	l.mutex.Lock()
	if l.consumerRunning {
		l.mutex.Unlock()
		return
	}
	l.consumerRunning = true
	l.mutex.Unlock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("orchestration: starting consumer")
		err := l.consumer.Start(ctx, lastEthBlockSynced)
		l.mutex.Lock()
		l.consumerRunning = false
		l.mutex.Unlock()
		if err != nil {
			log.Warnf("orchestration: consumer error. Error: %s", err)
		}
		log.Infof("orchestration: consumer finished")
		chConsumer <- err
	}()
}

func (l *L1SyncOrchestration) orchestrate(ctx context.Context, wg *sync.WaitGroup, chProducer chan error, chConsumer chan error) (*state.Block, error) {
	// Wait a cond_var for known if consumer have finish
	var err error
	done := false
	for !done {
		select {
		case <-ctx.Done():
			log.Warnf("orchestration: context cancelled")
			done = true
		case err = <-chProducer:
			// Producer has finished
			log.Infof("orchestration: producer has finished. Error: %s, stopping consumer", err)
			l.consumer.StopAfterProcessChannelQueue()
		case err = <-chConsumer:
			if err != nil && err != errAllWorkersBusy {
				log.Warnf("orchestration: consumer have finished with Error: %s", err)
			} else {
				log.Info("orchestration: consumer has finished. No error")
			}
			done = true
		}
	}
	l.isRunning = false
	retBlock, ok := l.consumer.GetLastEthBlockSynced()

	if err == nil {
		if ok {
			log.Infof("orchestration: finished L1 sync orchestration With LastBlock. Last block synced: %d err:nil", retBlock.BlockNumber)
			return &retBlock, nil
		} else {
			err := errors.New(errMissingLastEthBlockSynced)
			log.Warnf("orchestration: finished L1 sync orchestration No LastBlock. Last block synced: %s err:%s", "<no previous block>", err)
			return nil, err
		}
	} else {
		log.Warnf("orchestration: finished L1 sync orchestration With Error. Last block synced: %s err:%s", "IGNORED (nil)", err)
		return nil, err
	}
}
