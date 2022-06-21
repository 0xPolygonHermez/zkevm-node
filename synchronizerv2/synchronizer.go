package synchronizerv2

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	etherman "github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/log"
	state "github.com/hermeznetwork/hermez-core/statev2"
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
			lastEthBlockSynced = &state.Block{
				BlockNumber: s.genBlockNumber,
			}
			// TODO Set Genesis if needed
		} else {
			log.Fatal("unexpected error getting the latest ethereum block. Setting genesis block. Error: ", err)
		}
	} else if lastEthBlockSynced.BlockNumber == 0 {
		lastEthBlockSynced = &state.Block{
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
func (s *ClientSynchronizer) syncBlocks(lastEthBlockSynced *state.Block) (*state.Block, error) {
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

		if lastKnownBlock.Cmp(new(big.Int).SetUint64(fromBlock)) < 1 {
			waitDuration = s.cfg.SyncInterval.Duration
			break
		}
		if len(blocks) == 0 { // If there is no events in the checked blocks range and lastKnownBlock > fromBlock.
			// Store the latest block of the block range. Get block info and process the block
			fb, err := s.etherMan.EthBlockByNumber(s.ctx, toBlock)
			if err != nil {
				return lastEthBlockSynced, err
			}
			b := state.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			s.processBlockRange([]state.Block{b}, order)
			lastEthBlockSynced = &b
			log.Debug("Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
		}
	}

	return lastEthBlockSynced, nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []state.Block, order map[common.Hash][]etherman.Order) {
	// New info has to be included into the db using the state
	for i := range blocks {
		ctx := context.Background()
		// Begin db transaction
		txDB, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to store block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		// Add block information
		err = s.state.AddBlock(ctx, &blocks[i], txDB)
		if err != nil {
			rollbackErr := s.state.RollbackState(ctx, txDB)
			if rollbackErr != nil {
				log.Fatal(fmt.Sprintf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %v", blocks[i].BlockNumber, rollbackErr, err))
			}
			log.Fatalf("error storing block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
		}
		for _, element := range order[blocks[i].BlockHash] {
			// TODO Implement the store methods for each event
			log.Debug(element)
		}
		err = s.state.CommitState(ctx, txDB)
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
func (s *ClientSynchronizer) checkReorg(latestBlock *state.Block) (*state.Block, error) {
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
				return &state.Block{}, nil
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
