package synchronizer

import (
	"context"
	"fmt"
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
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
	etherMan  etherman.EtherMan
	state     state.State
	ctx       context.Context
	cancelCtx context.CancelFunc
	config    Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(ethMan etherman.EtherMan, st state.State, cfg Config) (Synchronizer, error) {
	//TODO
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientSynchronizer{
		state:     st,
		etherMan:  ethMan,
		ctx:       ctx,
		cancelCtx: cancel,
		config:    cfg,
	}, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	go func() {
		//If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
		//Get the latest synced block. If there is no block on db, use genesis block
		lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx)
		if err != nil || lastEthBlockSynced.BlockNumber == 0 {
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.config.GenesisBlock,
			}
		}
		waitDuration := time.Duration(0)
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(waitDuration):
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
func (s *ClientSynchronizer) syncBlocks(lastEthBlockSynced *state.Block) (*state.Block, error) {
	//This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(*lastEthBlockSynced)
	if err != nil {
		log.Error("error checking reorgs")
		return nil, fmt.Errorf("error checking reorgs")
	} else if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Error("error resetting the state to a previous block")
			return nil, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}

	//Call the blockchain to retrieve data
	blocks, err := s.etherMan.GetBatchesByBlockRange(s.ctx, lastEthBlockSynced.BlockNumber+1, nil)
	if err != nil {
		return nil, err
	}

	// New info has to be included into the db using the state
	for i := range blocks {
		//get lastest synced batch
		latestBatch, err := s.state.GetLastBatch(s.ctx, false)
		if err != nil {
			log.Error("error getting latest batch. Error: ", err)
			latestBatch = &state.Batch{
				BatchNumber: 0,
			}
		} else if latestBatch == nil {
			latestBatch = &state.Batch{
				BatchNumber: 0,
			}
		}

		batchProcessor := s.state.NewBatchProcessor(latestBatch.BatchNumber, false)

		//Add block information
		err = s.state.AddBlock(&blocks[i])
		if err != nil {
			log.Fatal("error storing block. BlockNumber: ", blocks[i].BlockNumber)
		}
		for _, seq := range blocks[i].NewSequencers {
			//Add new sequencers
			err := s.state.AddNewSequencer(seq)
			if err != nil {
				log.Fatal("error storing new sequencer in Block: ", blocks[i].BlockNumber, " Sequencer: ", seq)
			}
		}
		for j := range blocks[i].Batches {
			//Add batches
			err := batchProcessor.ProcessBatch(&blocks[i].Batches[j])
			if err != nil {
				log.Fatal("error processing batch. BatchNumber: ", blocks[i].Batches[j].BatchNumber, ". Error: ", err)
			}
		}
	}
	if len(blocks) != 0 {
		return &blocks[len(blocks)-1], nil
	}
	return lastEthBlockSynced, nil
}

// This function allows reset the state until an specific ethereum block
func (s *ClientSynchronizer) resetState(ethBlockNum uint64) error {
	err := s.state.Reset(ethBlockNum)
	if err != nil {
		return err
	}

	return nil
}

// This function will check if there is a reorg
func (s *ClientSynchronizer) checkReorg(currentBlock state.Block) (*state.Block, error) {
	//TODO this function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	//getLastEtherblockfromdb and check the hash and parent hash. Using the ethBlockNum, get this info from the blockchain and compare.
	//if the values doesn't match get the previous ethereum block from db (last-1) and get the info for that ethereum block number
	//from the blockchain. Compare the values. If they don't match do this step again. If matches, we have found the good ethereum block.
	// Now, return the ethereum block number
	return nil, nil
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}
