package synchronizer

import (
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

type Synchronizer struct {
    etherMan *etherman.EtherMan
    state *state.State
    
}
func NewSynchronizer(ethMan *etherman.EtherMan, st *state.State) (*Synchronizer, error) {
	//TODO
	return &Synchronizer{state: st, etherMan: ethMan}, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates. If it is already synced,
// It will keep waiting for a new event
func (s *Synchronizer) Sync() error {
	//TODO
	return nil
}

// This function allows reset the state until an specific ethereum block
func (s *Synchronizer) resetState(ethBlockNum uint64) error {
	//TODO
	return nil
}

// This function will check if there is a reorg
func (s *Synchronizer) reorg() (uint64, error) {
	//TODO
	return 0, nil
}
