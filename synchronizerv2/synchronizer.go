package synchronizerv2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	etherman "github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/log"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	etherMan       ethermanInterface
	state          stateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	ethMan ethermanInterface,
	st stateInterface,
	genBlockNumber uint64,
	cfg Config) (Synchronizer, error) {
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

var waitDuration = time.Duration(0)

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Info("Sync started")
	lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx)
	if err != nil {
		if err == state.ErrStateNotSynchronized {
			log.Warn("error getting the latest ethereum block. No data stored. Setting genesis block. Error: ", err)
			lastEthBlockSynced = &etherman.Block{
				BlockNumber: s.genBlockNumber,
			}
			// TODO Set Genesis if needed
		} else {
			log.Fatal("unexpected error getting the latest ethereum block. Setting genesis block. Error: ", err)
		}
	} else if lastEthBlockSynced.BlockNumber == 0 {
		lastEthBlockSynced = &etherman.Block{
			BlockNumber: s.genBlockNumber,
		}
	}
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			latestsequencedBatchNumber, err := s.etherMan.GetLatestBatchNumber()
			if err != nil {
				log.Warn("error getting latest sequenced batch in the rollup. Error: ", err)
				continue
			}
			//Sync L1Blocks
			if lastEthBlockSynced, err = s.syncBlocks(lastEthBlockSynced); err != nil {
				log.Warn("error syncing blocks: ", err)
				if s.ctx.Err() != nil {
					continue
				}
			}
			if waitDuration != s.cfg.SyncInterval.Duration {
				// Check latest Synced Batch
				latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx)
				if err != nil {
					log.Warn("error getting latest batch synced. Error: ", err)
					continue
				}
				if latestSyncedBatch == latestsequencedBatchNumber {
					waitDuration = s.cfg.SyncInterval.Duration
				}
				if latestSyncedBatch > latestsequencedBatchNumber {
					log.Fatal("error: latest Synced BatchNumber is higher than the latest Proposed BatchNumber in the rollup")
				}
			}
			// Sync L2Blocks
			// TODO
		}
	}
}

// This function syncs the node from a specific block to the latest
func (s *ClientSynchronizer) syncBlocks(lastEthBlockSynced *etherman.Block) (*etherman.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %s", err.Error())
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Err: %s, Retrying...", err.Error())
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}

	// Call the blockchain to retrieve data
	header, err := s.etherMan.HeaderByNumber(s.ctx, nil)
	if err != nil {
		return lastEthBlockSynced, err
	}
	lastKnownBlock := header.Number

	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber + 1
	}

	for {
		toBlock := fromBlock + s.cfg.SyncChunkSize

		log.Infof("Getting rollup info from block %d to block %d", fromBlock, toBlock)
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
			for i := range blocks {
				log.Debug("Position: ", i, ". BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
			}
		}
		fromBlock = toBlock + 1

		if lastKnownBlock.Cmp(new(big.Int).SetUint64(toBlock)) < 1 {
			waitDuration = s.cfg.SyncInterval.Duration
			break
		}
		if len(blocks) == 0 { // If there is no events in the checked blocks range and lastKnownBlock > fromBlock.
			// Store the latest block of the block range. Get block info and process the block
			fb, err := s.etherMan.EthBlockByNumber(s.ctx, toBlock)
			if err != nil {
				return lastEthBlockSynced, err
			}
			b := etherman.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			s.processBlockRange([]etherman.Block{b}, order)
			lastEthBlockSynced = &b
			log.Debug("Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
		}
	}

	return lastEthBlockSynced, nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		txDB, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to store block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		// Add block information
		err = s.state.AddBlock(s.ctx, &blocks[i], txDB)
		if err != nil {
			rollbackErr := s.state.RollbackState(s.ctx, txDB)
			if rollbackErr != nil {
				log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
			}
			log.Fatalf("error storing block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		for _, element := range order[blocks[i].BlockHash] {
			// TODO Implement the store methods for each event
			if element.Name == etherman.SequenceBatchesOrder {
				s.processSequenceBatches(blocks[i].SequencedBatches[element.Pos], blocks[i].BlockNumber, txDB)
			}
		}
		err = s.state.CommitState(s.ctx, txDB)
		if err != nil {
			log.Fatalf("error committing state to store block. BlockNumber: %v, err: %v", blocks[i].BlockNumber, err)
		}
	}
}

// This function allows reset the state until an specific ethereum block
func (s *ClientSynchronizer) resetState(blockNumber uint64) error {
	log.Debug("Reverting synchronization to block: ", blockNumber)
	txDB, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Error("error starting a db transaction to reset the state. Error: ", err)
		return err
	}
	err = s.state.Reset(s.ctx, blockNumber, txDB)
	if err != nil {
		rollbackErr := s.state.RollbackState(s.ctx, txDB)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blockNumber, rollbackErr, err)
			return rollbackErr
		}
		log.Error("error resetting the state. Error: ", err)
		return err
	}
	err = s.state.CommitState(s.ctx, txDB)
	if err != nil {
		rollbackErr := s.state.RollbackState(s.ctx, txDB)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blockNumber, rollbackErr, err)
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
func (s *ClientSynchronizer) checkReorg(latestBlock *etherman.Block) (*etherman.Block, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestEthBlockSynced := *latestBlock
	var depth uint64
	for {
		block, err := s.etherMan.EthBlockByNumber(s.ctx, latestBlock.BlockNumber)
		if err != nil {
			log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %s", latestBlock.BlockNumber, err.Error())
			return nil, err
		}
		if block.NumberU64() != latestBlock.BlockNumber {
			err = fmt.Errorf("Wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
				latestBlock.BlockNumber, block.NumberU64())
			log.Error("error: ", err)
			return nil, err
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
			latestBlock, err = s.state.GetPreviousBlock(s.ctx, depth)
			if errors.Is(err, state.ErrNotFound) {
				log.Warn("error checking reorg: previous block not found in db: ", err)
				return &etherman.Block{}, nil
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

func (s *ClientSynchronizer) checkTrustedState(batch state.TrustedBatch, txDB pgx.Tx) (bool, error) {
	// First get trusted batch from db
	tBatch, err := s.state.GetTrustedBatchByNumber(s.ctx, batch.BatchNumber, txDB)
	if err != nil {
		return false, err
	}
	//Compare virtual state with trusted state
	if batch.RawTxs == tBatch.RawTxs &&
		batch.GlobalExitRoot == tBatch.GlobalExitRoot &&
		batch.Timestamp == tBatch.Timestamp &&
		batch.Sequencer == tBatch.Sequencer {
		return true, nil
	}
	return false, nil
}

func (s *ClientSynchronizer) processSequenceBatches(batches []etherman.SequencedBatch, blockNumber uint64, txDB pgx.Tx) {
	for _, batch := range batches {
		vb := state.VirtualBatch{
			BatchNumber: batch.BatchNumber,
			TxHash:      batch.TxHash,
			Sequencer:   batch.Sequencer,
			BlockNumber: blockNumber,
		}
		virtualBatches := []state.VirtualBatch{vb}
		tb := state.TrustedBatch{
			BatchNumber:    batch.BatchNumber,
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      time.Unix(int64(batch.Timestamp), 0),
			Sequencer:      batch.Sequencer,
			RawTxs:         hex.EncodeToString(batch.Transactions),
		}
		trustedBatches := []state.TrustedBatch{tb}
		// ForcedBatchesmust be processed after the trusted batch.
		numForcedBatches := len(batch.ForceBatchesTimestamp)
		if numForcedBatches > 0 {
			// Read forcedBatches from db
			forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, numForcedBatches, txDB)
			if err != nil {
				log.Errorf("error getting forcedBatches. BatchNumber: %d", vb.BatchNumber)
				rollbackErr := s.state.RollbackState(s.ctx, txDB)
				if rollbackErr != nil {
					log.Fatal(fmt.Sprintf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v, error : %v", vb.BatchNumber, blockNumber, rollbackErr, err))
				}
				log.Fatalf("error getting forcedBatches. BatchNumber: %d, BlockNumber: %d, error: %v", vb.BatchNumber, blockNumber, err)
			}
			if numForcedBatches != len(*forcedBatches) {
				log.Fatal("error number of forced batches doesn't match")
			}
			for i, forcedBatch := range *forcedBatches {
				vb := state.VirtualBatch{
					BatchNumber: batch.BatchNumber + uint64(i),
					TxHash:      batch.TxHash,
					Sequencer:   batch.Sequencer,
					BlockNumber: blockNumber,
				}
				virtualBatches = append(virtualBatches, vb)
				tb := state.TrustedBatch{
					BatchNumber:    batch.BatchNumber + uint64(i), // First process the trusted and then the forcedBatches
					GlobalExitRoot: forcedBatch.GlobalExitRoot,
					Timestamp:      time.Unix(int64(batch.ForceBatchesTimestamp[i]), 0), // ForceBatchesTimestamp instead of forcedAt because it is the timestamp selected by the sequencer, not when the forced batch was sent. This forcedAt is the min timestamp allowed.
					Sequencer:      forcedBatch.Sequencer,
					RawTxs:         forcedBatch.RawTxsData,
				}
				trustedBatches = append(trustedBatches, tb)
			}
		}

		if len(virtualBatches) != len(trustedBatches) {
			log.Fatal("error: length of trustedBatches and virtualBatches don't match.\nvirtualBatches: %+v \ntrustedBatches: %+v", virtualBatches, trustedBatches)
		}

		// Now we need to check all the trusted batches. ForcedBatches should be already stored as trusted because this is don by the trusted sequencer
		for i, trustedBatch := range trustedBatches {
			// Call the check trusted state method to compare trusted and virtual state
			status, err := s.checkTrustedState(trustedBatch, txDB)
			if err != nil {
				if errors.Is(err, state.ErrNotFound) {
					log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
					// If it is not found, store trustedBatch
					err = s.state.AddTrustedBatch(s.ctx, trustedBatch, txDB)
					if err != nil {
						rollbackErr := s.state.RollbackState(s.ctx, txDB)
						if rollbackErr != nil {
							log.Fatal(fmt.Sprintf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v, error : %v", trustedBatch.BatchNumber, blockNumber, rollbackErr, err))
						}
						log.Fatalf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", trustedBatch.BatchNumber, blockNumber, err)
					}
					status = true
				} else {
					log.Fatal("error checking trusted state: ", err)
				}
			}
			if !status {
				// Reset trusted state
				err := s.state.ResetTrustedState(s.ctx, trustedBatch.BatchNumber, txDB) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
				if err != nil {
					rollbackErr := s.state.RollbackState(s.ctx, txDB)
					if rollbackErr != nil {
						log.Fatal(fmt.Sprintf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v, error : %v", trustedBatch.BatchNumber, blockNumber, rollbackErr, err))
					}
					log.Fatalf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", trustedBatch.BatchNumber, blockNumber, err)
				}
				err = s.state.AddTrustedBatch(s.ctx, trustedBatch, txDB)
				if err != nil {
					rollbackErr := s.state.RollbackState(s.ctx, txDB)
					if rollbackErr != nil {
						log.Fatal(fmt.Sprintf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v, error : %v", trustedBatch.BatchNumber, blockNumber, rollbackErr, err))
					}
					log.Fatalf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", trustedBatch.BatchNumber, blockNumber, err)
				}
			}
			// Store virtualBatch
			err = s.state.AddVirtualBatch(s.ctx, virtualBatches[i], txDB)
			if err != nil {
				rollbackErr := s.state.RollbackState(s.ctx, txDB)
				if rollbackErr != nil {
					log.Fatal(fmt.Sprintf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v, error : %v", virtualBatches[i].BatchNumber, blockNumber, rollbackErr, err))
				}
				log.Fatalf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatches[i].BatchNumber, blockNumber, err)
			}
		}
	}
}
