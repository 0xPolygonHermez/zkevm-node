package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	etherMan       etherman.EtherMan
	state          state.State
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(ethMan etherman.EtherMan, st state.State, genBlockNumber uint64, cfg Config) (Synchronizer, error) {
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

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	go func() {
		// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
		// Get the latest synced block. If there is no block on db, use genesis block
		lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx)
		if err != nil {
			log.Warn("error getting the latest ethereum block. Setting genesis block. Error: ", err)
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.genBlockNumber,
			}
		} else if lastEthBlockSynced.BlockNumber == 0 {
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.genBlockNumber,
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
				// Check latest Proposed Batch number in the smc
				latestProposedBatchNumber, err := s.etherMan.GetLatestProposedBatchNumber()
				if err != nil {
					log.Warn("error getting latest proposed batch in the rollup. Error: ", err)
					continue
				}
				err = s.state.SetLastBatchNumberSeenOnEthereum(s.ctx, latestProposedBatchNumber)
				if err != nil {
					log.Warn("error settign latest proposed batch into db. Error: ", err)
					continue
				}
				if waitDuration != s.cfg.SyncInterval.Duration {
					// Check latest Synced Batch
					latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx)
					if err != nil {
						log.Warn("error getting latest batch synced. Error: ", err)
						continue
					}
					if latestSyncedBatch == latestProposedBatchNumber {
						waitDuration = s.cfg.SyncInterval.Duration
					}
					if latestSyncedBatch > latestProposedBatchNumber {
						log.Fatal("error: latest Synced BatchNumber is higher than the latest Proposed BatchNumber in the rollup")
					}
				}
			}
		}
	}()
	return nil
}

// This function syncs the node from a specific block to the latest
func (s *ClientSynchronizer) syncBlocks(lastEthBlockSynced *state.Block) (*state.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	} else if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Error("error resetting the state to a previous block. Retrying...")
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}

	// Call the blockchain to retrieve data
	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber + 1
	}
	blocks, err := s.etherMan.GetBatchesByBlockRange(s.ctx, fromBlock, nil)
	if err != nil {
		return nil, err
	}

	// New info has to be included into the db using the state
	for i := range blocks {
		// Get lastest synced batch number
		latestBatchNumber, err := s.state.GetLastBatchNumber(s.ctx)
		if err != nil {
			log.Warn("error getting latest batch. Error: ", err)
		}

		// Add block information
		err = s.state.AddBlock(context.Background(), &blocks[i])
		if err != nil {
			log.Fatal("error storing block. BlockNumber: ", blocks[i].BlockNumber)
		}
		lastEthBlockSynced = &blocks[i]
		for _, seq := range blocks[i].NewSequencers {
			// Add new sequencers
			err := s.state.AddSequencer(context.Background(), seq)
			if err != nil {
				log.Fatal("error storing new sequencer in Block: ", blocks[i].BlockNumber, " Sequencer: ", seq)
			}
		}
		for j := range blocks[i].Batches {
			sequencerAddress := &blocks[i].Batches[j].Sequencer
			batchProcessor, err := s.state.NewBatchProcessor(*sequencerAddress, latestBatchNumber)
			if err != nil {
				log.Error("error creating new batch processor. Error: ", err)
			}

			// Add batches
			err = batchProcessor.ProcessBatch(&blocks[i].Batches[j])
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
	err := s.state.Reset(s.ctx, ethBlockNum)
	if err != nil {
		return err
	}

	return nil
}

/*
This function will check if there is a reorg.
As input param needs the last ethereum block synced. Retrieve the block info from the blockchain
to compare it with the stored info. If hash and hash parent matches, then no reorg is detected and return a nil.
If hash or hash parent don't match, reorg detected and the function will return the block until the sync process
must be reverted. Then, check the previous ethereum block synced, get block info from the blockchain and check
hash and has parent. This operation has to be done until a match is found.
*/
func (s *ClientSynchronizer) checkReorg(latestBlock *state.Block) (*state.Block, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
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
		// Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > s.genBlockNumber {
			// Reorg detected. Getting previous block
			latestBlock, err = s.state.GetBlockByNumber(s.ctx, latestBlock.BlockNumber-1)
			if err != nil {
				if err.Error() == pgx.ErrNoRows.Error() {
					return nil, nil
				}
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
