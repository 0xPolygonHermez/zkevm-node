package synchronizer

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
	// ResetAndStop set a new starting point and cancel current process if any
	ResetAndStop(startingBlockNumber uint64)
}

type l1RollupConsumerInterface interface {
	Start(ctx context.Context) error
	StopAfterProcessChannelQueue()
	GetLastEthBlockSynced() (state.Block, bool)
}

type l1SyncOrchestration struct {
	mutex    sync.Mutex
	producer l1RollupProducerInterface
	consumer l1RollupConsumerInterface
	// Producer is running?
	producerRunning bool
	consumerRunning bool
	// The orchestrator is running?
	isStarted  bool
	wg         sync.WaitGroup
	chProducer chan error
	chConsumer chan error
	ctxParent  context.Context
}

const (
	errMissingLastEthBlockSynced = "orchestration: missing last eth block synced"
)

func newL1SyncOrchestration(ctx context.Context, producer l1RollupProducerInterface, consumer l1RollupConsumerInterface) *l1SyncOrchestration {
	return &l1SyncOrchestration{
		producer:        producer,
		consumer:        consumer,
		producerRunning: false,
		consumerRunning: false,
		chProducer:      make(chan error, 1),
		chConsumer:      make(chan error, 1),
		ctxParent:       ctx,
	}
}

func (l *l1SyncOrchestration) reset(startingBlockNumber uint64) {
	log.Warnf("Reset L1 sync process to blockNumber %d", startingBlockNumber)
	if l.isStarted {
		log.Infof("orchestration: reset(%d) is going to stop producer", startingBlockNumber)
	}
	l.producer.ResetAndStop(startingBlockNumber)
	// If orchestrator is running then producer is going to be started by orchestrate() select  function when detects that producer has finished
}

func (l *l1SyncOrchestration) start() (*state.Block, error) {
	l.isStarted = true
	l.launchProducer(l.ctxParent, l.chProducer, &l.wg)
	l.launchConsumer(l.ctxParent, l.chConsumer, &l.wg)
	return l.orchestrate(l.ctxParent, &l.wg, l.chProducer, l.chConsumer)
}

func (l *l1SyncOrchestration) isProducerRunning() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.producerRunning
}

func (l *l1SyncOrchestration) launchProducer(ctx context.Context, chProducer chan error, wg *sync.WaitGroup) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if !l.producerRunning {
		if wg != nil {
			wg.Add(1)
		}
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

func (l *l1SyncOrchestration) launchConsumer(ctx context.Context, chConsumer chan error, wg *sync.WaitGroup) {
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
		err := l.consumer.Start(ctx)
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

func (l *l1SyncOrchestration) orchestrate(ctx context.Context, wg *sync.WaitGroup, chProducer chan error, chConsumer chan error) (*state.Block, error) {
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
			// if l.isStarted {
			// 	log.Warnf("orchestration: consumer have finished! respawn. Error:%s", err)
			// 	// to avoid respawn too fast it sleeps a bit
			// 	time.Sleep(time.Second)
			// 	l.launchProducer(ctx, chProducer, wg)
			// } else {
			// 	log.Infof("orchestration: consumer has finished. Error: %s", err)
			// }
		case err = <-chConsumer:
			if err != nil && err != errAllWorkersBusy {
				log.Warnf("orchestration: consumer have finished with Error: %s", err)
			} else {
				log.Info("orchestration: consumer has finished. No error")
			}
			done = true
		}
	}
	l.isStarted = false
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
