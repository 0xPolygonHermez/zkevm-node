package synchronizer

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

/*
This object is to coordinate the producer process and the consumer process.
*/
type l1DataRetrieverInterface interface {
	start() error
	stop()
}

type l1DataProcessorInterface interface {
	start() error
	finishExecutionWhenChannelIsEmpty()
	getLastEthBlockSynced() *state.Block
}

type l1SyncOrchestration struct {
	producer l1DataRetrieverInterface
	consumer l1DataProcessorInterface
}

func newL1SyncOrchestration(producer l1DataRetrieverInterface, consumer l1DataProcessorInterface) *l1SyncOrchestration {
	return &l1SyncOrchestration{
		producer: producer,
		consumer: consumer,
	}
}

// There are 2 reason for finish:
// 1) The producer process finish (have requested all the data from L1):
//   - Wait until consumer run out of data on channel
//
// 2) The consumer process finish if there are an error in the process of the data:
//   - Abort cosumer
func (l *l1SyncOrchestration) start() (*state.Block, error) {
	log.Info("orchestration: starting L1 sync orchestration")
	chProducer := make(chan error, 1)
	chConsumer := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	// Start producer: L1DataRetriever from L1
	go func() {
		defer wg.Done()
		err := l.producer.start()
		if err != nil {
			log.Warnf("orchestration: error starting L1DataRetriever. Error: %s", err)
		}
		chProducer <- err
	}()
	// Start consumer: L1DataProcessor execute the RollupInfo
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := l.consumer.start()
		if err != nil {
			log.Warnf("orchestration: error starting L1DataProcessor. Error: %s", err)
		}
		chConsumer <- err
	}()

	return l.orchestrate(&wg, chProducer, chConsumer)
}

func (l *l1SyncOrchestration) orchestrate(wg *sync.WaitGroup, hProducer chan error, chConsumer chan error) (*state.Block, error) {
	// Wait a cond_var for known if consumer have finish
	var err error
	done := false
	for !done {
		select {
		case err = <-hProducer:
			// Producer have finish
			if err != nil {
				log.Warnf("orchestration: DataRetriever (producer) have finish with  Error: %s", err)
			} else {
				log.Info("orchestration: DataRetriever (producer) have finish")
			}
			// process all pending RollupInfo and finish
			log.Info("orchestration: forcing to consumer to consume all pending RollupInfo and finish")
			l.consumer.finishExecutionWhenChannelIsEmpty()

		case err = <-chConsumer:
			if err != nil {
				log.Warnf("orchestration: DataProcessor (consumer) have finish with Error: %s", err)
			} else {
				log.Info("orchestration: DataProcessor (consumer) have finish")
			}
			log.Info("orchestration: Stoping producer because we don't need more rollupinfo")
			l.producer.stop()
			// Consumer have finish, return
			done = true
		}
	}
	log.Info("orchestration: waiting to finish producer and consumer")
	wg.Wait()
	retBlock := l.consumer.getLastEthBlockSynced()
	if retBlock != nil {
		log.Infof("orchestration: finished L1 sync orchestration. Last block synced: %d err:%s", retBlock.BlockNumber, err)
	} else {
		log.Infof("orchestration: finished L1 sync orchestration. Last block synced: %s err:%s", "nil", err)
	}
	return retBlock, err
}
