package synchronizer

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

type Synchronizer struct {
	etherMan  *etherman.EtherMan
	state     *state.State
	ctx       context.Context
	cancelCtx context.CancelFunc

	newBatchProposalHandlers     []NewBatchProposalHandler
	newConsolidatedStateHandlers []NewConsolidatedStateHandler
	stateResetHandlers           []StateResetHandler
}

type NewBatchProposalHandler func()
type NewConsolidatedStateHandler func()
type StateResetHandler func()

func NewSynchronizer(ethMan *etherman.EtherMan, st *state.State, ag chan int, sq chan int) (*Synchronizer, error) {
	//TODO
	ctx, cancel := context.WithCancel(context.Background())
	return &Synchronizer{
		state:     st,
		etherMan:  ethMan,
		ctx:       ctx,
		cancelCtx: cancel,
	}, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates. If it is already synced,
// It will keep waiting for a new event
func (s *Synchronizer) Sync() error {
	go func() {
		var lastEthBlockSynced uint64
		var err error
		waitDuration := time.Duration(0)
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(waitDuration):
				//TODO
				if lastEthBlockSynced, err = s.syncBlocks(lastEthBlockSynced); err != nil {
					if s.ctx.Err() != nil {
						continue
					}
				}
			}
		}
	}()
	return nil
}

func (s *Synchronizer) RegisterNewBatchProposalHandler(handler NewBatchProposalHandler) {
	s.newBatchProposalHandlers = append(s.newBatchProposalHandlers, handler)
}

func (s *Synchronizer) RegisterNewConsolidatedStateHandler(handler NewConsolidatedStateHandler) {
	s.newConsolidatedStateHandlers = append(s.newConsolidatedStateHandlers, handler)
}

func (s *Synchronizer) RegisterStateResetHandler(handler StateResetHandler) {
	s.stateResetHandlers = append(s.stateResetHandlers, handler)
}

// This function syncs the node from a specific block to the latest
func (s *Synchronizer) syncBlocks(lastEthBlockSynced uint64) (uint64, error) {
	//TODO
	//This function will read events fromBlockNum to latestEthBlock. First It has to retrieve the latestEthereumBlock and check reorg to be sure that everything is ok.
	//if there is no lastEthereumBlock means that sync from the begining is necesary. If not, it continues from the retrieved ethereum block
	// New info has to be included into the db using the state

	// When a new batch propostal is synchronized, we notify it
	// go s.notifyNewBathProposal()

	// When a new consolidated state is synchronized, we notify it
	// go s.notifyNewConsolidatedState()

	return 0, nil
}

// This function allows reset the state until an specific ethereum block
func (s *Synchronizer) resetState(ethBlockNum uint64) error {
	err := s.state.Reset(ethBlockNum)
	if err != nil {
		return err
	}

	go s.notifyResetState()

	return nil
}

// This function will check if there is a reorg
func (s *Synchronizer) checkReorg() (uint64, error) {
	//TODO this function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	//getLastEtherblockfromdb and check the hash and parent hash. Using the ethBlockNum, get this info from the blockchain and compare.
	//if the values doesn't match get the previous ethereum block from db (last-1) and get the info for that ethereum block number
	//from the blockchain. Compare the values. If they don't match do this step again. If matches, we have found the good ethereum block.
	// Now, return the ethereum block number
	return 0, nil
}

// Stop function stops the synchronizer
func (s *Synchronizer) Stop() {
	s.cancelCtx()
}

// notifyResetState notifies all registered reset state handlers that the state was reset
func (s *Synchronizer) notifyResetState() {
	for _, handler := range s.stateResetHandlers {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					//log.Error(err)
				}
			}()
			handler()
		}()
	}
}

// notifyNewBathProposal notifies all registered new batch proposal handlers
// that a new batch proposal was synchronized
func (s *Synchronizer) notifyNewBathProposal() {
	for _, handler := range s.newBatchProposalHandlers {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					//log.Error(err)
				}
			}()
			handler()
		}()
	}
}

// notifyNewConsolidatedState notifies all registered new consolidated state handlers
// that a new consolidated state was synchronized
func (s *Synchronizer) notifyNewConsolidatedState() {
	for _, handler := range s.newConsolidatedStateHandlers {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					//log.Error(err)
				}
			}()
			handler()
		}()
	}
}
