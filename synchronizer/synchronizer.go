package synchronizer

import (
	"context"
	"errors"
	"fmt"
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
	etherMan       localEtherman
	state          stateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	genesis        state.Genesis
	cfg            Config
	gpe            gasPriceEstimator
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	ethMan localEtherman,
	st stateInterface,
	genBlockNumber uint64,
	genesis state.Genesis,
	cfg Config,
	gpe gasPriceEstimator) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientSynchronizer{
		state:          st,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: genBlockNumber,
		genesis:        genesis,
		cfg:            cfg,
		gpe:            gpe,
	}, nil
}

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
		waitDuration := time.Duration(0)
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(waitDuration):
				// Check latest Proposed Batch number in the smc
				latestProposedBatchNumber, err := s.etherMan.GetLatestProposedBatchNumber()
				if err != nil {
					log.Warn("error getting latest proposed batch in the rollup. Error: ", err)
					continue
				}
				err = s.state.SetLastBatchNumberSeenOnEthereum(s.ctx, latestProposedBatchNumber, "")
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
				err = s.state.SetLastBatchNumberConsolidatedOnEthereum(s.ctx, latestConsolidatedBatchNumber, "")
				if err != nil {
					log.Warn("error setting latest consolidated batch into db. Error: ", err)
					continue
				}
				if lastEthBlockSynced, err = s.syncBlocks(lastEthBlockSynced); err != nil {
					log.Warn("error syncing blocks: ", err)
					if s.ctx.Err() != nil {
						continue
					}
				}
				if waitDuration != s.cfg.SyncInterval.Duration {
					// Check latest Synced Batch
					latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx, "")
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
		err = s.resetState(block)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Err: %v, Retrying...", err)
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
		return lastEthBlockSynced, err
	}
	lastKnownBlock := header.Number

	for {
		toBlock := fromBlock + s.cfg.SyncChunkSize

		log.Debugf("Getting rollup info from block %d to block %d", fromBlock, toBlock)
		// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
		// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
		// Name can be defferent in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
		// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
		// array index where this value is.
		blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, fromBlock, &toBlock)
		if err != nil {
			return lastEthBlockSynced, err
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

	return lastEthBlockSynced, nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []state.Block, order map[common.Hash][]etherman.Order) {
	// New info has to be included into the db using the state
	for i := range blocks {
		ctx := context.Background()
		// Begin db transaction
		txBundleID, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to store block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		// Add block information
		err = s.state.AddBlock(ctx, &blocks[i], txBundleID)
		if err != nil {
			rollbackErr := s.state.RollbackState(ctx, txBundleID)
			if rollbackErr != nil {
				log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
			}
			log.Fatalf("error storing block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		for _, element := range order[blocks[i].BlockHash] {
			if element.Name == etherman.BatchesOrder {
				batch := &blocks[i].Batches[element.Pos]
				emptyHash := common.Hash{}
				log.Debug("consolidatedTxHash received: ", batch.ConsolidatedTxHash)
				if batch.ConsolidatedTxHash.String() != emptyHash.String() {
					// consolidate batch locally
					err = s.state.ConsolidateBatch(ctx, batch.Number().Uint64(), batch.ConsolidatedTxHash, *batch.ConsolidatedAt, batch.Aggregator, txBundleID)
					if err != nil {
						rollbackErr := s.state.RollbackState(ctx, txBundleID)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
						}
						log.Fatal("failed to consolidate batch locally, batch number: %d, err: %v", batch.Number().Uint64(), err)
					}
				} else {
					// Get latest synced batch number
					latestBatchNumber, err := s.state.GetLastBatchNumber(ctx, txBundleID)
					if err != nil {
						rollbackErr := s.state.RollbackState(ctx, txBundleID)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
						}
						log.Fatal("error getting latest batch. Error: ", err)
					}

					// Get batch header
					latestBatchHeader, err := s.state.GetBatchHeader(ctx, latestBatchNumber, txBundleID)
					if err != nil {
						rollbackErr := s.state.RollbackState(ctx, txBundleID)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
						}
						log.Fatal("error getting latest batch header. Error: ", err)
					}

					sequencerAddress := batch.Sequencer
					batchProcessor, err := s.state.NewBatchProcessor(ctx, sequencerAddress, latestBatchHeader.Root[:], txBundleID)
					if err != nil {
						rollbackErr := s.state.RollbackState(ctx, txBundleID)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
						}
						log.Fatal("error creating new batch processor. Error: ", err)
					}
					// Add batches
					err = batchProcessor.ProcessBatch(ctx, batch)
					if err != nil {
						rollbackErr := s.state.RollbackState(ctx, txBundleID)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
						}
						log.Fatal("error processing batch. BatchNumber: ", batch.Number().Uint64(), ". Error: ", err)
					}
					s.gpe.UpdateGasPriceAvg(new(big.Int).SetUint64(batch.Header.GasUsed))
				}
			} else if element.Name == etherman.NewSequencersOrder {
				// Add new sequencers
				err := s.state.AddSequencer(ctx, blocks[i].NewSequencers[element.Pos], txBundleID)
				if err != nil {
					rollbackErr := s.state.RollbackState(ctx, txBundleID)
					if rollbackErr != nil {
						log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
					}
					log.Fatal("error storing new sequencer in Block: ", blocks[i].BlockNumber, " Sequencer: ", blocks[i].NewSequencers[element.Pos], " err: ", err)
				}
			} else {
				rollbackErr := s.state.RollbackState(ctx, txBundleID)
				if rollbackErr != nil {
					log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
				}
				log.Fatal("error: invalid order element")
			}
		}
		err = s.state.CommitState(ctx, txBundleID)
		if err != nil {
			log.Fatalf("error committing state to store block. BlockNumber: %v, err: %v", blocks[i].BlockNumber, err)
		}
	}
}

// This function allows reset the state until an specific ethereum block
func (s *ClientSynchronizer) resetState(block *state.Block) error {
	log.Debug("Reverting synchronization to block: ", block.BlockNumber)
	txBundleID, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Error("error starting a db transaction to reset the state. Error: ", err)
		return err
	}
	err = s.state.Reset(s.ctx, block, txBundleID)
	if err != nil {
		rollbackErr := s.state.RollbackState(s.ctx, txBundleID)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", block.BlockNumber, rollbackErr, err)
			return rollbackErr
		}
		log.Error("error resetting the state. Error: ", err)
		return err
	}
	err = s.state.CommitState(s.ctx, txBundleID)
	if err != nil {
		rollbackErr := s.state.RollbackState(s.ctx, txBundleID)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", block.BlockNumber, rollbackErr, err)
			return rollbackErr
		}
		log.Error("error committing the resetted state. Error: ", err)
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
	var depth uint64
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
			log.Debug("[checkReorg function] => latestBlockNumber: ", latestBlock.BlockNumber)
			log.Debug("[checkReorg function] => latestBlockHash: ", latestBlock.BlockHash)
			log.Debug("[checkReorg function] => latestBlockHashParent: ", latestBlock.ParentHash)
			log.Debug("[checkReorg function] => BlockNumber: ", latestBlock.BlockNumber, block.NumberU64())
			log.Debug("[checkReorg function] => BlockHash: ", block.Hash())
			log.Debug("[checkReorg function] => BlockHashParent: ", block.ParentHash())
			depth++
			log.Debug("REORG: Looking for the latest correct ethereum block. Depth: ", depth)
			// Reorg detected. Getting previous block
			latestBlock, err = s.state.GetPreviousBlock(s.ctx, depth, "")
			if errors.Is(err, state.ErrNotFound) {
				log.Warn("error checking reorg: previous block not found in db: ", err)
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
