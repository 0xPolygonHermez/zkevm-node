package synchronizerv2

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/log"
	state "github.com/hermeznetwork/hermez-core/statev2"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	etherMan       ethermanInterface
	state          stateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	ethMan ethermanInterface,
	st stateInterface,
	genBlockNumber uint64,
	cfg Config) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientSynchronizer{
		state:          st,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: genBlockNumber,
		cfg:            cfg,
	}, nil
}

var waitDuration = time.Duration(0)

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Info("Sync started")
	lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx, "")
	if err != nil {
		if err == state.ErrStateNotSynchronized {
			log.Warn("error getting the latest ethereum block. No data stored. Setting genesis block. Error: ", err)
			lastEthBlockSynced = &ethermanv2.Block{
				BlockNumber: s.genBlockNumber,
			}
		} else {
			log.Fatal("unexpected error getting the latest ethereum block. Setting genesis block. Error: ", err)
		}
	} else if lastEthBlockSynced.BlockNumber == 0 {
		lastEthBlockSynced = &ethermanv2.Block{
			BlockNumber: s.genBlockNumber,
		}
	}
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			//Sync L1Blocks and L2Blocks
		}
	}
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}
