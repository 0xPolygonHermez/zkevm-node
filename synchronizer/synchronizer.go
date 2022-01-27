package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	etherMan       etherman.EtherMan
	state          state.State
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	genBalances    state.Genesis
	cfg            Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(ethMan etherman.EtherMan, st state.State, genBlockNumber uint64, genBalances state.Genesis, cfg Config) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientSynchronizer{
		state:          st,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: genBlockNumber,
		genBalances:    genBalances,
		cfg:            cfg,
	}, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	go func() {
		// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
		// Get the latest synced block. If there is no block on db, use genesis block
		log.Info("Sync started")
		lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx)
		if err != nil {
			if err == state.ErrStateNotSynchronized {
				log.Warn("error getting the latest ethereum block. No data stored. Setting genesis block. Error: ", err)
				lastEthBlockSynced = &state.Block{
					BlockNumber: s.genBlockNumber,
				}
				// Set genesis
				err := s.state.SetGenesis(s.ctx, s.genBalances)
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
					log.Warn("error setting latest proposed batch into db. Error: ", err)
					continue
				}

				// Check latest Consolidated Batch number in the scm
				latestConsolidatedBatchNumber, err := s.etherMan.GetLatestConsolidatedBatchNumber()
				if err != nil {
					log.Warn("error getting latest consolidated batch in the rollup. Error: ", err)
					continue
				}
				err = s.state.SetLastBatchNumberConsolidatedOnEthereum(s.ctx, latestConsolidatedBatchNumber)
				if err != nil {
					log.Warn("error setting latest consolidated batch into db. Error: ", err)
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

	header, err := s.etherMan.HeaderByNumber(s.ctx, nil)
	if err != nil {
		return nil, err
	}
	lastKnownBlock := header.Number

	for {
		toBlock := fromBlock + s.cfg.SyncChunkSize

		log.Debugf("Getting rollback info from block %d to block %d", fromBlock, toBlock)
		// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
		// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
		// Name can be defferent in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
		// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
		// array index where this value is.
		blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, fromBlock, &toBlock)
		if err != nil {
			return nil, err
		}
		s.processBlockRange(blocks, order)
		if len(blocks) > 0 {
			lastEthBlockSynced = &blocks[len(blocks)-1]
		}
		fromBlock = toBlock + 1

		if lastKnownBlock.Cmp(new(big.Int).SetUint64(fromBlock)) < 1 {
			break
		}
	}

	// in order to prevent repeating querying and checking blocks we return the
	// latest block checked minus some safety number to avoid issues with reorgs.
	// safetyBlocks is the default number of blocks to check to always take into
	// account reorgs.
	const safetyBlocks = 50
	if lastKnownBlock.Cmp(new(big.Int).SetUint64(lastEthBlockSynced.BlockNumber+uint64(safetyBlocks))) == 1 {
		blockHeight := math.Max(0, float64(lastKnownBlock.Uint64()-uint64(safetyBlocks)))
		lastEthBlockSynced = state.NewBlock(uint64(blockHeight))
	}

	return lastEthBlockSynced, nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []state.Block, order map[common.Hash][]etherman.Order) {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Add block information
		err := s.state.AddBlock(context.Background(), &blocks[i])
		if err != nil {
			log.Fatal("error storing block. BlockNumber: ", blocks[i].BlockNumber)
		}
		for _, element := range order[blocks[i].BlockHash] {
			if element.Name == etherman.BatchesOrder {
				batch := &blocks[i].Batches[element.Pos]
				emptyHash := common.Hash{}
				log.Debug("consolidatedTxHash received: ", batch.ConsolidatedTxHash)
				if batch.ConsolidatedTxHash.String() != emptyHash.String() {
					// consolidate batch locally
					err = s.state.ConsolidateBatch(s.ctx, batch.BatchNumber, batch.ConsolidatedTxHash, *batch.ConsolidatedAt)
					if err != nil {
						log.Warnf("failed to consolidate batch locally, batch number: %d, err: %v",
							batch.BatchNumber, err)
						continue
					}
				} else {
					// Get lastest synced batch number
					latestBatchNumber, err := s.state.GetLastBatchNumber(s.ctx)
					if err != nil {
						log.Fatal("error getting latest batch. Error: ", err)
					}

					sequencerAddress := batch.Sequencer
					batchProcessor, err := s.state.NewBatchProcessor(sequencerAddress, latestBatchNumber)
					if err != nil {
						log.Error("error creating new batch processor. Error: ", err)
					}
					// Add batches
					err = batchProcessor.ProcessBatch(batch)
					if err != nil {
						log.Fatal("error processing batch. BatchNumber: ", batch.BatchNumber, ". Error: ", err)
					}
				}
			} else if element.Name == etherman.NewSequencersOrder {
				// Add new sequencers
				err := s.state.AddSequencer(context.Background(), blocks[i].NewSequencers[element.Pos])
				if err != nil {
					log.Fatal("error storing new sequencer in Block: ", blocks[i].BlockNumber, " Sequencer: ", blocks[i].NewSequencers[element.Pos])
				}
			} else if element.Name == etherman.DepositsOrder {
				//TODO Store info into db
				log.Warn("Deposit functionality is not implemented in synchronizer yet")
			} else if element.Name == etherman.GlobalExitRootsOrder {
				//TODO Store info into db
				log.Warn("Consolidate globalExitRoot functionality is not implemented in synchronizer yet")
			} else if element.Name == etherman.ClaimsOrder {
				//TODO Store info into db
				log.Warn("Claim functionality is not implemented in synchronizer yet")
			} else {
				log.Fatal("error: invalid order element")
			}
		}
	}
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
			if errors.Is(err, state.ErrNotFound) {
				return nil, nil
			} else if err != nil {
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
