package synchronizer

import (
	"errors"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

/*
This object is to coordinate the producer process and the consumer process.
*/
type l1RollupProducerInterface interface {
	start(startingBlockNumber uint64) error
	stop()
	reset(startingBlockNumber uint64)
}

type l1RollupConsumerInterface interface {
	start() error
	stopAfterProcessChannelQueue()
	getLastEthBlockSynced() (state.Block, bool)
}

type l1SyncOrchestration struct {
	mutex           sync.Mutex
	producer        l1RollupProducerInterface
	consumer        l1RollupConsumerInterface
	producerStarted bool
}

const (
	errMissingLastEthBlockSynced = "orchestration: missing last eth block synced"
)

func newL1SyncOrchestration(producer l1RollupProducerInterface, consumer l1RollupConsumerInterface) *l1SyncOrchestration {
	return &l1SyncOrchestration{
		producer:        producer,
		consumer:        consumer,
		producerStarted: false,
	}
}

func (l *l1SyncOrchestration) reset(startingBlockNumber uint64) {
	log.Warnf("Reset L1 sync process to blockNumber %d", startingBlockNumber)
	l.producer.reset(startingBlockNumber)
}

func (l *l1SyncOrchestration) start(startingBlockNumber uint64) (*state.Block, error) {
	chProducer := make(chan error, 1)
	chConsumer := make(chan error, 1)
	var wg sync.WaitGroup
	l.launch_producer(startingBlockNumber, chProducer, &wg)
	l.launch_consumer(chConsumer, &wg)
	return l.orchestrate(&wg, chProducer, chConsumer)
}

func (l *l1SyncOrchestration) launch_producer(startingBlockNumber uint64, chProducer chan error, wg *sync.WaitGroup) {
	l.mutex.Lock()
	if !l.producerStarted {
		if wg != nil {
			wg.Add(1)
		}

		// Start producer: L1DataRetriever from L1
		l.producerStarted = true
		l.mutex.Unlock()
		go func() {
			if wg != nil {
				defer wg.Done()
			}
			log.Infof("orchestration: starting producer(%v)", startingBlockNumber)
			err := l.producer.start(startingBlockNumber)
			if err != nil {
				log.Warnf("orchestration: producer error . Error: %s", err)
			}
			l.mutex.Lock()
			l.producerStarted = false
			l.mutex.Unlock()
			log.Infof("orchestration: producer finished")
			chProducer <- err
		}()
	} else {
		l.mutex.Unlock()
	}
}

func (l *l1SyncOrchestration) launch_consumer(chConsumer chan error, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("orchestration: starting consumer")
		err := l.consumer.start()
		if err != nil {
			log.Warnf("orchestration: consumer error. Error: %s", err)
		}
		log.Infof("orchestration: consumer finished")
		chConsumer <- err
	}()
}

func (l *l1SyncOrchestration) orchestrate(wg *sync.WaitGroup, hProducer chan error, chConsumer chan error) (*state.Block, error) {
	// Wait a cond_var for known if consumer have finish
	var err error
	done := false
	for !done {
		select {
		case err = <-hProducer:
			// Producer have finish
			log.Warnf("orchestration: consumer have finish! this never have to happen, restarting it. Error:%s", err)
			// to avoid respawn too fast it sleeps a bit
			time.Sleep(time.Second)
			l.launch_producer(invalidBlockNumber, hProducer, wg)
		case err = <-chConsumer:
			if err != nil {
				log.Warnf("orchestration: consumer have finish with Error: %s", err)
			} else {
				log.Info("orchestration: consumer have finish ok")
			}
			done = true
		}
	}
	retBlock, ok := l.consumer.getLastEthBlockSynced()
	if ok {
		log.Infof("orchestration: finished L1 sync orchestration. Last block synced: %d err:%s", retBlock.BlockNumber, err)
		return &retBlock, nil
	} else {
		err := errors.New(errMissingLastEthBlockSynced)
		log.Infof("orchestration: finished L1 sync orchestration. Last block synced: %s err:%s", "???", err)
		return nil, err
	}
}
