package synchronizer

import (
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state/db"
)

type Synchronizer struct {
    etherMan *etherman.EtherMan
    state *state.State
    
}
func NewSynchronizer() (*Synchronizer, error) {
	//TODO
	var db db.KeyValuer
	st := state.NewState(db)
	ethMan, err := etherman.NewEtherman()
	if err != nil {
		return nil, err
	}
	return &Synchronizer{state: st, etherMan: ethMan}, nil
}

// This function will read the last state synced and will continue from that point
func (s *Synchronizer) Sync() error {
	//TODO
	return nil
}

// This function allows reset the state until an specific ethereum block
func (s *Synchronizer) resetState(ethBlockNum uint64) error {
	//TODO
	return nil
}

// This function will check if there is a reorg and if so, fix the sync state
func (s *Synchronizer) reorg() (uint64, error) {
	//TODO
	return 0, nil
}
