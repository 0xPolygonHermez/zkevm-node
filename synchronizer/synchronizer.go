package synchronizer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	isTrustedSequencer bool
	etherMan           ethermanInterface
	state              stateInterface
	ctx                context.Context
	cancelCtx          context.CancelFunc
	genesis            state.Genesis
	cfg                Config
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	isTrustedSequencer bool,
	ethMan ethermanInterface,
	st stateInterface,
	genesis state.Genesis,
	cfg Config) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &ClientSynchronizer{
		isTrustedSequencer: isTrustedSequencer,
		state:              st,
		etherMan:           ethMan,
		ctx:                ctx,
		cancelCtx:          cancel,
		genesis:            genesis,
		cfg:                cfg,
	}, nil
}

var waitDuration = time.Duration(0)

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Info("Sync started")
	dbTx, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Fatalf("error creating db transaction to get latest block")
	}
	lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx, dbTx)
	if err != nil {
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Info("State is empty, setting genesis block")
			header, err := s.etherMan.HeaderByNumber(s.ctx, big.NewInt(0).SetUint64(s.cfg.GenBlockNumber))
			if err != nil {
				log.Fatal("error getting l1 block header for block ", s.cfg.GenBlockNumber, " : ", err)
			}
			lastEthBlockSynced = &state.Block{
				BlockNumber: header.Number.Uint64(),
				BlockHash:   header.Hash(),
				ParentHash:  header.ParentHash,
				ReceivedAt:  time.Unix(int64(header.Time), 0),
			}
			newRoot, err := s.state.SetGenesis(s.ctx, *lastEthBlockSynced, s.genesis, dbTx)
			if err != nil {
				log.Fatal("error setting genesis: ", err)
			}
			var root common.Hash
			root.SetBytes(newRoot)
			if root != s.genesis.Root {
				log.Fatal("Calculated newRoot should be ", s.genesis.Root, " instead of ", root)
			}
			log.Debug("Genesis root matches!")
		} else {
			log.Fatal("unexpected error getting the latest ethereum block. Error: ", err)
		}
	}
	if err := dbTx.Commit(s.ctx); err != nil {
		log.Errorf("error committing dbTx, err: %w", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. RollbackErr: %s, err: %w",
				rollbackErr.Error(), err)
		}
		log.Fatalf("error committing dbTx, err: %w", err)
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			//Sync L1Blocks
			if lastEthBlockSynced, err = s.syncBlocks(lastEthBlockSynced); err != nil {
				log.Warn("error syncing blocks: ", err)
				if s.ctx.Err() != nil {
					continue
				}
			}
			latestSequencedBatchNumber, err := s.etherMan.GetLatestBatchNumber()
			if err != nil {
				log.Warn("error getting latest sequenced batch in the rollup. Error: ", err)
				continue
			}
			latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx, nil)
			if err != nil {
				log.Warn("error getting latest batch synced. Error: ", err)
				continue
			}
			if latestSyncedBatch >= latestSequencedBatchNumber {
				log.Info("L1 state fully synchronized")
				err = s.syncTrustedState(latestSyncedBatch)
				if err != nil {
					log.Warn("error syncing trusted state. Error: ", err)
					continue
				}
				log.Info("Trusted state fully synchronized")
				waitDuration = s.cfg.SyncInterval.Duration
			}
		}
	}
}

// This function syncs the node from a specific block to the latest
func (s *ClientSynchronizer) syncBlocks(lastEthBlockSynced *state.Block) (*state.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %w", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %w", err)
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
			lastEthBlockSynced = &state.Block{
				BlockNumber: blocks[len(blocks)-1].BlockNumber,
				BlockHash:   blocks[len(blocks)-1].BlockHash,
				ParentHash:  blocks[len(blocks)-1].ParentHash,
				ReceivedAt:  blocks[len(blocks)-1].ReceivedAt,
			}
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
			block := state.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			lastEthBlockSynced = &block
			log.Debug("Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
		}
	}

	return lastEthBlockSynced, nil
}

// syncTrustedState synchronizes information from the trusted sequencer
// related to the trusted state when the node has all the information from
// l1 synchronized
func (s *ClientSynchronizer) syncTrustedState(latestSyncedBatch uint64) error {
	if s.isTrustedSequencer {
		return nil
	}

	log.Debug("Getting broadcast URI")
	broadcastURI, err := s.getBroadcastURI()
	if err != nil {
		log.Errorf("error getting broadcast URI. Error: %v", err)
		return err
	}
	log.Debug("broadcastURI ", broadcastURI)
	broadcastClient, _, _ := broadcast.NewClient(s.ctx, broadcastURI)

	log.Info("Getting trusted state info")
	lastTrustedStateBatch, err := broadcastClient.GetLastBatch(s.ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("error syncing trusted state. Error: ", err)
		return err
	}

	log.Debug("lastTrustedStateBatch.BatchNumber ", lastTrustedStateBatch.BatchNumber)
	log.Debug("latestSyncedBatch ", latestSyncedBatch)
	if lastTrustedStateBatch.BatchNumber < latestSyncedBatch {
		return nil
	}

	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatch.BatchNumber {
		batchToSync, err := broadcastClient.GetBatch(s.ctx, &pb.GetBatchRequest{BatchNumber: batchNumberToSync})
		if err != nil {
			log.Warnf("failed to get batch %v from trusted state via broadcast. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
		}

		if err := s.processTrustedBatch(batchToSync, dbTx); err != nil {
			log.Errorf("error processing trusted batch %v: %v", batchNumberToSync, err)
			err := dbTx.Rollback(s.ctx)
			if err != nil {
				log.Fatalf("error rolling back db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			}
			break
		}

		if err := dbTx.Commit(s.ctx); err != nil {
			log.Fatalf("error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
		}

		batchNumberToSync++
	}

	return nil
}

// gets the broadcast URI from trusted sequencer JSON RPC server
func (s *ClientSynchronizer) getBroadcastURI() (string, error) {
	log.Debug("getting trusted sequencer URL from smc")
	trustedSequencerURL, err := s.etherMan.GetTrustedSequencerURL()
	if err != nil {
		return "", err
	}
	log.Debug("trustedSequencerURL ", trustedSequencerURL)

	log.Debug("getting broadcast URI from Trusted Sequencer JSON RPC Server")
	res, err := jsonrpc.JSONRPCCall(trustedSequencerURL, "zkevm_getBroadcastURI")
	if err != nil {
		return "", err
	}

	if res.Error != nil {
		errMsg := fmt.Sprintf("%v:%v", res.Error.Code, res.Error.Message)
		return "", errors.New(errMsg)
	}

	var url string
	if err := json.Unmarshal(res.Result, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to store block. BlockNumber: %d, error: %w", blocks[i].BlockNumber, err)
		}
		b := state.Block{
			BlockNumber: blocks[i].BlockNumber,
			BlockHash:   blocks[i].BlockHash,
			ParentHash:  blocks[i].ParentHash,
			ReceivedAt:  blocks[i].ReceivedAt,
		}
		// Add block information
		err = s.state.AddBlock(s.ctx, &b, dbTx)
		if err != nil {
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %w", blocks[i].BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error storing block. BlockNumber: %d, error: %w", blocks[i].BlockNumber, err)
		}
		for _, element := range order[blocks[i].BlockHash] {
			switch element.Name {
			case etherman.SequenceBatchesOrder:
				s.processSequenceBatches(blocks[i].SequencedBatches[element.Pos], blocks[i].BlockNumber, dbTx)
			case etherman.ForcedBatchesOrder:
				s.processForcedBatch(blocks[i].ForcedBatches[element.Pos], dbTx)
			case etherman.GlobalExitRootsOrder:
				s.processGlobalExitRoot(blocks[i].GlobalExitRoots[element.Pos], dbTx)
			case etherman.SequenceForceBatchesOrder:
				s.processSequenceForceBatch(blocks[i].SequencedForceBatches[element.Pos], blocks[i], dbTx)
			case etherman.VerifyBatchOrder:
				s.processVerifyBatches(blocks[i].VerifiedBatches[element.Pos], dbTx)
			}
		}
		err = dbTx.Commit(s.ctx)
		if err != nil {
			log.Errorf("error committing state to store block. BlockNumber: %d, err: %w", blocks[i].BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %w", blocks[i].BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error committing state to store block. BlockNumber: %d, err: %w", blocks[i].BlockNumber, err)
		}
	}
}

// This function allows reset the state until an specific ethereum block
func (s *ClientSynchronizer) resetState(blockNumber uint64) error {
	log.Debug("Reverting synchronization to block: ", blockNumber)
	dbTx, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Error("error starting a db transaction to reset the state. Error: ", err)
		return err
	}
	err = s.state.Reset(s.ctx, blockNumber, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %w", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error resetting the state. Error: ", err)
		return err
	}
	err = dbTx.Commit(s.ctx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %w", blockNumber, rollbackErr.Error(), err)
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
			log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %w", latestBlock.BlockNumber, err)
			return nil, err
		}
		if block.NumberU64() != latestBlock.BlockNumber {
			err = fmt.Errorf("Wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
				latestBlock.BlockNumber, block.NumberU64())
			log.Error("error: ", err)
			return nil, err
		}
		// Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > s.cfg.GenBlockNumber {
			log.Debug("[checkReorg function] => latestBlockNumber: ", latestBlock.BlockNumber)
			log.Debug("[checkReorg function] => latestBlockHash: ", latestBlock.BlockHash)
			log.Debug("[checkReorg function] => latestBlockHashParent: ", latestBlock.ParentHash)
			log.Debug("[checkReorg function] => BlockNumber: ", latestBlock.BlockNumber, block.NumberU64())
			log.Debug("[checkReorg function] => BlockHash: ", block.Hash())
			log.Debug("[checkReorg function] => BlockHashParent: ", block.ParentHash())
			depth++
			log.Debug("REORG: Looking for the latest correct ethereum block. Depth: ", depth)
			// Reorg detected. Getting previous block
			dbTx, err := s.state.BeginStateTransaction(s.ctx)
			if err != nil {
				log.Fatalf("error creating db transaction to get prevoius blocks")
			}
			latestBlock, err = s.state.GetPreviousBlock(s.ctx, depth, dbTx)
			errC := dbTx.Commit(s.ctx)
			if errC != nil {
				log.Errorf("error committing dbTx, err: %w", errC)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. RollbackErr: %w, err: %w",
						rollbackErr, errC)
				}
				log.Fatalf("error committing dbTx, err: %w", errC)
			}
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

func (s *ClientSynchronizer) checkTrustedState(batch state.Batch, dbTx pgx.Tx) (bool, error) {
	// First get trusted batch from db
	tBatch, err := s.state.GetBatchByNumber(s.ctx, batch.BatchNumber, dbTx)
	if err != nil {
		return false, err
	}

	// Reprocess batch and compare the stateRoot with tBatch.StateRoot
	p, err := s.state.ExecuteBatch(s.ctx, batch.BatchNumber, batch.BatchL2Data, dbTx)
	if err != nil {
		log.Errorf("error executing L1 batch: %+v, error: %w", batch, err)
		return false, err
	}
	newRoot := common.BytesToHash(p.NewStateRoot)

	//Compare virtual state with trusted state
	if hex.EncodeToString(batch.BatchL2Data) == hex.EncodeToString(tBatch.BatchL2Data) &&
		batch.GlobalExitRoot.String() == tBatch.GlobalExitRoot.String() &&
		batch.Timestamp.Unix() == tBatch.Timestamp.Unix() &&
		batch.Coinbase.String() == tBatch.Coinbase.String() &&
		newRoot == tBatch.StateRoot {
		return true, nil
	}
	log.Warn("Trusted Reorg detected")
	log.Debug("batch.BatchL2Data: ", hex.EncodeToString(batch.BatchL2Data))
	log.Debug("batch.GlobalExitRoot: ", batch.GlobalExitRoot)
	log.Debug("batch.Timestamp: ", batch.Timestamp)
	log.Debug("batch.Coinbase: ", batch.Coinbase)
	log.Debug("newRoot: ", newRoot)
	log.Debug("tBatch.BatchL2Data: ", hex.EncodeToString(tBatch.BatchL2Data))
	log.Debug("tBatch.GlobalExitRoot: ", tBatch.GlobalExitRoot)
	log.Debug("tBatch.Timestamp: ", tBatch.Timestamp)
	log.Debug("tBatch.Coinbase: ", tBatch.Coinbase)
	log.Debug("tBatch.StateRoot: ", tBatch.StateRoot)
	return false, nil
}

func (s *ClientSynchronizer) processSequenceBatches(sequencedBatches []etherman.SequencedBatch, blockNumber uint64, dbTx pgx.Tx) {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return
	}
	for _, sbatch := range sequencedBatches {
		virtualBatch := state.VirtualBatch{
			BatchNumber: sbatch.BatchNumber,
			TxHash:      sbatch.TxHash,
			Coinbase:    sbatch.Coinbase,
			BlockNumber: blockNumber,
		}
		batch := state.Batch{
			BatchNumber:    sbatch.BatchNumber,
			GlobalExitRoot: sbatch.GlobalExitRoot,
			Timestamp:      time.Unix(int64(sbatch.Timestamp), 0),
			Coinbase:       sbatch.Coinbase,
			BatchL2Data:    sbatch.Transactions,
		}
		// ForcedBatch must be processed
		if sbatch.MinForcedTimestamp > 0 {
			// Read forcedBatches from db
			forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, 1, dbTx)
			if err != nil {
				log.Errorf("error getting forcedBatches. BatchNumber: %d", virtualBatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				}
				log.Fatalf("error getting forcedBatches. BatchNumber: %d, BlockNumber: %d, error: %w", virtualBatch.BatchNumber, blockNumber, err)
			}
			if len(forcedBatches) == 0 {
				log.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", sbatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				}
				log.Fatal("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
			}
			if uint64(forcedBatches[0].ForcedAt.Unix()) != sbatch.MinForcedTimestamp ||
				forcedBatches[0].GlobalExitRoot != sbatch.GlobalExitRoot ||
				common.Bytes2Hex(forcedBatches[0].RawTxsData) != common.Bytes2Hex(sbatch.Transactions) ||
				forcedBatches[0].Sequencer != sbatch.Coinbase {
				log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %w", virtualBatch.BatchNumber, blockNumber, rollbackErr)
				}
				log.Fatalf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
			}
			// Store batchNumber in forced_batch table
			err = s.state.AddBatchNumberInForcedBatch(s.ctx, forcedBatches[0].ForcedBatchNumber, sbatch.BatchNumber, dbTx)
			if err != nil {
				log.Errorf("error adding the batchNumber to forcedBatch in processSequenceBatches. BlockNumber: %d", blockNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", blockNumber, rollbackErr.Error(), err)
				}
				log.Fatalf("error adding the batchNumber to forcedBatch in processSequenceBatches. BlockNumber: %d, error: %w", blockNumber, err)
			}
		}

		// Now we need to check the batch. ForcedBatches should be already stored in the batch table because this is done by the sequencer
		processCtx := state.ProcessingContext{
			BatchNumber:    batch.BatchNumber,
			Coinbase:       batch.Coinbase,
			Timestamp:      batch.Timestamp,
			GlobalExitRoot: batch.GlobalExitRoot,
		}
		// Call the check trusted state method to compare trusted and virtual state
		status, err := s.checkTrustedState(batch, dbTx)
		if err != nil {
			if errors.Is(err, state.ErrNotFound) || errors.Is(err, state.ErrStateNotSynchronized) {
				log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
				// If it is not found, store batch
				err = s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx)
				if err != nil {
					log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
					rollbackErr := dbTx.Rollback(s.ctx)
					if rollbackErr != nil {
						log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					}
					log.Fatalf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
				}
				status = true
			} else {
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				}
				log.Fatal("error checking trusted state: ", err)
			}
		}
		if !status {
			// Reset trusted state
			previousBatchNumber := batch.BatchNumber - 1
			log.Warnf("Trusted reorg detected, discarding batches until batchNum %d", previousBatchNumber)
			err := s.state.ResetTrustedState(s.ctx, previousBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
			if err != nil {
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				}
				log.Fatalf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
			}
			err = s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx)
			if err != nil {
				log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				}
				log.Fatalf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, blockNumber, err)
			}
		}
		// Store virtualBatch
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %w", virtualBatch.BatchNumber, blockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %w", virtualBatch.BatchNumber, blockNumber, err)
		}
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
	}
	err := s.state.AddSequence(s.ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", blockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error getting adding sequence. BlockNumber: %d, error: %w", blockNumber, err)
	}
}

func (s *ClientSynchronizer) processSequenceForceBatch(sequenceForceBatch []etherman.SequencedForceBatch, block etherman.Block, dbTx pgx.Tx) {
	if len(sequenceForceBatch) == 0 {
		log.Warn("Empty sequenceForceBatch array detected, ignoring...")
		return
	}
	// First, get last virtual batch number
	lastVirtualizedBatchNumber, err := s.state.GetLastVirtualBatchNum(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %w", block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %w", block.BlockNumber, err)
	}
	// Second, reset trusted state
	err = s.state.ResetTrustedState(s.ctx, lastVirtualizedBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
	if err != nil {
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %w", lastVirtualizedBatchNumber, block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %w", lastVirtualizedBatchNumber, block.BlockNumber, err)
	}
	// Read forcedBatches from db
	forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, len(sequenceForceBatch), dbTx)
	if err != nil {
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d", block.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", block.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d, error: %w", block.BlockNumber, err)
	}
	if len(sequenceForceBatch) != len(forcedBatches) {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", block.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatal("error number of forced batches doesn't match")
	}
	for i, fbatch := range sequenceForceBatch {
		if uint64(forcedBatches[i].ForcedAt.Unix()) != fbatch.MinForcedTimestamp ||
			forcedBatches[i].GlobalExitRoot != fbatch.GlobalExitRoot ||
			common.Bytes2Hex(forcedBatches[i].RawTxsData) != common.Bytes2Hex(fbatch.Transactions) ||
			forcedBatches[i].Sequencer != fbatch.Coinbase {
			log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %w", fbatch.BatchNumber, block.BlockNumber, rollbackErr)
			}
			log.Fatalf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
		}
		virtualBatch := state.VirtualBatch{
			BatchNumber: fbatch.BatchNumber,
			TxHash:      fbatch.TxHash,
			Coinbase:    fbatch.Coinbase,
			BlockNumber: block.BlockNumber,
		}
		batch := state.ProcessingContext{
			BatchNumber:    fbatch.BatchNumber,
			GlobalExitRoot: fbatch.GlobalExitRoot,
			Timestamp:      block.ReceivedAt,
			Coinbase:       fbatch.Coinbase,
		}
		// Process batch
		err := s.state.ProcessAndStoreClosedBatch(s.ctx, batch, forcedBatches[i].RawTxsData, dbTx)
		if err != nil {
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", batch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %w", batch.BatchNumber, block.BlockNumber, err)
		}
		// Store virtualBatch
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %w", virtualBatch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %w", virtualBatch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %w", virtualBatch.BatchNumber, block.BlockNumber, err)
		}
		// Store batchNumber in forced_batch table
		err = s.state.AddBatchNumberInForcedBatch(s.ctx, forcedBatches[i].ForcedBatchNumber, virtualBatch.BatchNumber, dbTx)
		if err != nil {
			log.Errorf("error adding the batchNumber to forcedBatch in processSequenceForceBatch. BlockNumber: %d", block.BlockNumber)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", block.BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error adding the batchNumber to forcedBatch in processSequenceForceBatch. BlockNumber: %d, error: %w", block.BlockNumber, err)
		}
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequenceForceBatch[0].BatchNumber,
		ToBatchNumber:   sequenceForceBatch[len(sequenceForceBatch)-1].BatchNumber,
	}
	err = s.state.AddSequence(s.ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", block.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error getting adding sequence. BlockNumber: %d, error: %w", block.BlockNumber, err)
	}
}

func (s *ClientSynchronizer) processForcedBatch(forcedBatch etherman.ForcedBatch, dbTx pgx.Tx) {
	// Store forced batch into the db
	forcedB := state.ForcedBatch{
		BlockNumber:       forcedBatch.BlockNumber,
		BatchNumber:       nil,
		ForcedBatchNumber: forcedBatch.ForcedBatchNumber,
		Sequencer:         forcedBatch.Sequencer,
		GlobalExitRoot:    forcedBatch.GlobalExitRoot,
		RawTxsData:        forcedBatch.RawTxsData,
		ForcedAt:          forcedBatch.ForcedAt,
	}
	err := s.state.AddForcedBatch(s.ctx, &forcedB, dbTx)
	if err != nil {
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d", forcedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", forcedBatch.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d, error: %w", forcedBatch.BlockNumber, err)
	}
}

func (s *ClientSynchronizer) processGlobalExitRoot(globalExitRoot etherman.GlobalExitRoot, dbTx pgx.Tx) {
	// Store GlobalExitRoot
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		Timestamp:       globalExitRoot.Timestamp,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := s.state.AddGlobalExitRoot(s.ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d", globalExitRoot.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", globalExitRoot.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %w", globalExitRoot.BlockNumber, err)
	}
}

func (s *ClientSynchronizer) processVerifyBatches(lastVerifiedBatch etherman.VerifiedBatch, dbTx pgx.Tx) {
	lastVBatch, err := s.state.GetLastVerifiedBatch(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVerifiedBatch stored in db in processVerifyBatches. Processing synced blockNumber: %d", lastVerifiedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. Processing synced blockNumber: %d, rollbackErr: %s, error : %w", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
		}
		log.Fatalf("error getting lastVerifiedBatch stored in db in processVerifyBatches. Processing synced blockNumber: %d, error: %w", lastVerifiedBatch.BlockNumber, err)
	}
	nbatches := lastVerifiedBatch.BatchNumber - lastVBatch.BatchNumber
	var i uint64
	for i = 1; i <= nbatches; i++ {
		verifiedB := state.VerifiedBatch{
			BlockNumber: lastVerifiedBatch.BlockNumber,
			BatchNumber: lastVBatch.BatchNumber + i,
			Aggregator:  lastVerifiedBatch.Aggregator,
			StateRoot:   lastVerifiedBatch.StateRoot,
			TxHash:      lastVerifiedBatch.TxHash,
		}
		batch, err := s.state.GetBatchByNumber(s.ctx, verifiedB.BatchNumber, dbTx)
		if err != nil {
			log.Errorf("error getting GetBatchByNumber stored in db in processVerifyBatches. Processing blockNumber: %d", verifiedB.BatchNumber)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. Processing blockNumber: %d, rollbackErr: %s, error : %w", verifiedB.BatchNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error getting GetBatchByNumber stored in db in processVerifyBatches. Processing blockNumber: %d, error: %w", verifiedB.BatchNumber, err)
		}

		// Checks that calculated state root matches with the verified state root in the smc
		if batch.StateRoot != verifiedB.StateRoot {
			log.Errorf("error: stateRoot calculated and state root verified don't match in processVerifyBatches. Processing blockNumber: %d, error: %w", verifiedB.BatchNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. Processing blockNumber: %d, rollbackErr: %s, error : %w", verifiedB.BatchNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error: stateRoot calculated and state root verified don't match in processVerifyBatches. Processing blockNumber: %d, error: %w", verifiedB.BatchNumber, err)
		}

		err = s.state.AddVerifiedBatch(s.ctx, &verifiedB, dbTx)
		if err != nil {
			log.Errorf("error storing the verifiedB in processVerifyBatches. verifiedBatch: %+v, lastVerifiedBatch: %+v", verifiedB, lastVerifiedBatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Fatalf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %w", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
			}
			log.Fatalf("error storing the verifiedB in processVerifyBatches. BlockNumber: %d, error: %w", lastVerifiedBatch.BlockNumber, err)
		}
	}
}

func (s *ClientSynchronizer) processTrustedBatch(trustedBatch *pb.GetBatchResponse, dbTx pgx.Tx) error {
	log.Debugf("processing trusted batch: %v", trustedBatch.BatchNumber)
	txs := []types.Transaction{}
	for _, transaction := range trustedBatch.Transactions {
		tx, err := state.DecodeTx(transaction.Encoded)
		if err != nil {
			return err
		}
		txs = append(txs, *tx)
	}
	trustedBatchL2Data, err := state.EncodeTransactions(txs)
	if err != nil {
		return err
	}

	batch, err := s.state.GetBatchByNumber(s.ctx, trustedBatch.BatchNumber, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.BatchNumber, err)
		return err
	}

	// check if batch needs to be synchronized
	if batch != nil {
		matchNumber := batch.BatchNumber == trustedBatch.BatchNumber
		matchGER := batch.GlobalExitRoot.String() == trustedBatch.GlobalExitRoot
		matchLER := batch.LocalExitRoot.String() == trustedBatch.LocalExitRoot
		matchSR := batch.StateRoot.String() == trustedBatch.StateRoot
		matchCoinbase := batch.Coinbase.String() == trustedBatch.Sequencer
		matchTimestamp := uint64(batch.Timestamp.Unix()) == trustedBatch.Timestamp
		matchL2Data := hex.EncodeToString(batch.BatchL2Data) == hex.EncodeToString(trustedBatchL2Data)

		if matchNumber && matchGER && matchLER && matchSR &&
			matchCoinbase && matchTimestamp && matchL2Data {
			log.Debugf("batch %v already synchronized", trustedBatch.BatchNumber)
			return nil
		}
		log.Infof("batch %v needs to be updated", trustedBatch.BatchNumber)
	} else {
		log.Infof("batch %v needs to be synchronized", trustedBatch.BatchNumber)
	}

	log.Debugf("resetting trusted state from batch %v", trustedBatch.BatchNumber)
	previousBatchNumber := trustedBatch.BatchNumber - 1
	if err := s.state.ResetTrustedState(s.ctx, previousBatchNumber, dbTx); err != nil {
		log.Errorf("failed to reset trusted state", trustedBatch.BatchNumber)
		return err
	}

	log.Debugf("opening batch %v", trustedBatch.BatchNumber)
	processCtx := state.ProcessingContext{
		BatchNumber:    trustedBatch.BatchNumber,
		Coinbase:       common.HexToAddress(trustedBatch.Sequencer),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: common.HexToHash(trustedBatch.GlobalExitRoot),
	}
	if err := s.state.OpenBatch(s.ctx, processCtx, dbTx); err != nil {
		log.Errorf("error opening batch %d", trustedBatch.BatchNumber)
		return err
	}

	log.Debugf("processing sequencer for batch %v", trustedBatch.BatchNumber)

	processBatchResp, err := s.state.ProcessSequencerBatch(s.ctx, trustedBatch.BatchNumber, txs, dbTx)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %d", trustedBatch.BatchNumber)
		return err
	}

	log.Debugf("storing transactions for batch %v", trustedBatch.BatchNumber)
	if err = s.state.StoreTransactions(s.ctx, trustedBatch.BatchNumber, processBatchResp.Responses, dbTx); err != nil {
		log.Errorf("failed to store transactions for batch: %d", trustedBatch.BatchNumber)
		return err
	}

	log.Debug("trustedBatch.StateRoot ", trustedBatch.StateRoot)
	isBatchClosed := trustedBatch.StateRoot != state.ZeroHash.String()
	if isBatchClosed {
		receipt := state.ProcessingReceipt{
			BatchNumber:   trustedBatch.BatchNumber,
			StateRoot:     processBatchResp.NewStateRoot,
			LocalExitRoot: processBatchResp.NewLocalExitRoot,
		}
		log.Debugf("closing batch %v", trustedBatch.BatchNumber)
		if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
			log.Errorf("error closing batch %d", trustedBatch.BatchNumber)
			return err
		}
	}

	if trustedBatch.ForcedBatchNumber > 0 {
		log.Debugf("adding batch num %v for forced batch %v", trustedBatch.BatchNumber, trustedBatch.ForcedBatchNumber)
		if err := s.state.AddBatchNumberInForcedBatch(s.ctx, trustedBatch.ForcedBatchNumber, trustedBatch.BatchNumber, dbTx); err != nil {
			log.Errorf("error adding batch %v for forced batch %v", trustedBatch.BatchNumber, trustedBatch.ForcedBatchNumber)
			return err
		}
	}

	log.Infof("batch %v synchronized", trustedBatch.BatchNumber)
	return nil
}
