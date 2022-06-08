package synchronizerv2

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	etherMan       localEtherman
	state          stateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	genesis        state.Genesis
	cfg            Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	ethMan localEtherman,
	st stateInterface,
	genBlockNumber uint64,
	genesis state.Genesis,
	cfg Config) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientSynchronizer{
		state:          st,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: genBlockNumber,
		genesis:        genesis,
		cfg:            cfg,
	}, nil
}

var waitDuration = time.Duration(0)

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	go func() {
		// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
		// Get the latest synced block. If there is no block on db, use genesis block
		log.Info("Sync started")
		lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx, "")
		if err != nil {
			if err == state.ErrStateNotSynchronized {
				log.Warn("error getting the latest ethereum block. No data stored. Setting genesis block. Error: ", err)
				lastEthBlockSynced = &state.Block{
					BlockNumber: s.genBlockNumber,
				}
				// Set genesis
				err := s.state.SetGenesis(s.ctx, s.genesis, "")
				if err != nil {
					log.Fatal("error setting genesis: ", err)
				}
			} else {
				log.Fatal("unexpected error getting the latest ethereum block. Setting genesis block. Error: ", err)
			}
		} else if lastEthBlockSynced.BlockNumber == 0 {
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.genBlockNumber,
			}
		}
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(waitDuration):
				//Sync L1Blocks and L2Blocks
			}
		}
	}()
	return nil
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}
