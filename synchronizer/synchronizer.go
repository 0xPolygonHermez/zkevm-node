package synchronizer

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// BasicSynchronizer connects L1 and L2
type BasicSynchronizer struct {
	etherMan  etherman.EtherMan
	state     state.State
	ctx       context.Context
	cancelCtx context.CancelFunc
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(ethMan etherman.EtherMan, st state.State) (Synchronizer, error) {
	//TODO
	ctx, cancel := context.WithCancel(context.Background())
	return &BasicSynchronizer{
		state:     st,
		etherMan:  ethMan,
		ctx:       ctx,
		cancelCtx: cancel,
	}, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates. If it is already synced,
// It will keep waiting for a new event
func (s *BasicSynchronizer) Sync() error {
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

// This function syncs the node from a specific block to the latest
func (s *BasicSynchronizer) syncBlocks(lastEthBlockSynced uint64) (uint64, error) {
	//TODO
	//This function will read events fromBlockNum to latestEthBlock. First It has to retrieve the latestEthereumBlock and check reorg to be sure that everything is ok.
	//if there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// New info has to be included into the db using the state

	return 0, nil
}

// // This function allows reset the state until an specific ethereum block
// func (s *BasicSynchronizer) resetState(ethBlockNum uint64) error {
// 	err := s.state.Reset(ethBlockNum)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // This function will check if there is a reorg
// func (s *BasicSynchronizer) checkReorg() (uint64, error) {
// 	//TODO this function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
// 	//getLastEtherblockfromdb and check the hash and parent hash. Using the ethBlockNum, get this info from the blockchain and compare.
// 	//if the values doesn't match get the previous ethereum block from db (last-1) and get the info for that ethereum block number
// 	//from the blockchain. Compare the values. If they don't match do this step again. If matches, we have found the good ethereum block.
// 	// Now, return the ethereum block number
// 	return 0, nil
// }

// Stop function stops the synchronizer
func (s *BasicSynchronizer) Stop() {
	s.cancelCtx()
}
