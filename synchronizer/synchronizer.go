package synchronizer

import (
	"context"
	"errors"
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
		if err != nil {
			log.Warn("error getting the latest ethereum block. Setting genesis block. Error: ", err)
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.config.GenesisBlock,
			}
		} else if lastEthBlockSynced.BlockNumber == 0 {
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
	block, err := s.checkReorg(lastEthBlockSynced)
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
	var fromBlock uint64 = 0
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber + 1
	}
	blocks, err := s.etherMan.GetBatchesByBlockRange(s.ctx, fromBlock, nil)
	if err != nil {
		return nil, err
	}

	// New info has to be included into the db using the state
	for i := range blocks {
		//get lastest synced batch number
		latestBatchNumber, err := s.state.GetLastBatchNumber(s.ctx)
		if err != nil {
			log.Error("error getting latest batch. Error: ", err)
		}

		batchProcessor, err := s.state.NewBatchProcessor(latestBatchNumber, false)
		if err != nil {
			log.Error("error creating new batch processor. Error: ", err)
		}

		//Add block information
		err = s.state.AddBlock(context.Background(), &blocks[i])
		if err != nil {
			log.Fatal("error storing block. BlockNumber: ", blocks[i].BlockNumber)
		}
		for _, seq := range blocks[i].NewSequencers {
			//Add new sequencers
			err := s.state.AddSequencer(context.Background(), seq)
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
	log.Debug("Reverting synchronization to block: ", ethBlockNum)
	err := s.state.Reset(ethBlockNum)
	if err != nil {
		return err
	}

	return nil
}

// This function will check if there is a reorg
func (s *ClientSynchronizer) checkReorg(latestBlock *state.Block) (*state.Block, error) {
	//This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestEthBlockSynced := *latestBlock
	for {
		block, err := s.etherMan.EthBlockByNumber(s.ctx, latestBlock.BlockNumber)
		if err != nil {
			if errors.Is(err, etherman.ErrNotFound) {
				return nil, nil
			}
			return nil, err
		}
		if block.NumberU64() != latestBlock.BlockNumber {
			log.Error("Wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: ",
				latestBlock.BlockNumber, ". BlockNumber retrieved: ", block.NumberU64())
			return nil, fmt.Errorf("Wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
				latestBlock.BlockNumber, block.NumberU64())
		}
		//Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > 0 {
			//Reorg detected. Getting previous block
			latestBlock, err = s.state.GetBlockByNumber(s.ctx, latestBlock.BlockNumber-1)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	if latestEthBlockSynced.BlockHash != latestBlock.BlockHash {
		log.Debug("Reorg detected in block: ", latestEthBlockSynced.BlockNumber)
		return latestBlock, nil
	}
	return nil, nil
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}
