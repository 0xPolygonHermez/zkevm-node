package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
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
	pool               poolInterface
	ethTxManager       ethTxManager
	zkEVMClient        zkEVMClientInterface
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
	pool poolInterface,
	ethTxManager ethTxManager,
	zkEVMClient zkEVMClientInterface,
	genesis state.Genesis,
	cfg Config) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &ClientSynchronizer{
		isTrustedSequencer: isTrustedSequencer,
		state:              st,
		etherMan:           ethMan,
		pool:               pool,
		ctx:                ctx,
		cancelCtx:          cancel,
		ethTxManager:       ethTxManager,
		zkEVMClient:        zkEVMClient,
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
			log.Info("State is empty, verifying genesis block")
			valid, err := s.etherMan.VerifyGenBlockNumber(s.ctx, s.cfg.GenBlockNumber)
			if err != nil {
				log.Error("error checking genesis block number. Error: ", err)
				return err
			} else if !valid {
				log.Error("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
				return fmt.Errorf("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
			}
			log.Info("Setting genesis block")
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
		log.Errorf("error committing dbTx, err: %v", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Fatalf("error rolling back state. RollbackErr: %s, err: %v",
				rollbackErr.Error(), err)
		}
		log.Fatalf("error committing dbTx, err: %v", err)
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			//Sync L1Blocks
			if lastEthBlockSynced, err = s.syncBlocks(lastEthBlockSynced); err != nil {
				log.Warn("error syncing blocks: ", err)
				lastEthBlockSynced, err = s.state.GetLastBlock(s.ctx, nil)
				if err != nil {
					log.Fatal("error getting lastEthBlockSynced to resume the synchronization... Error: ", err)
				}
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
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
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
		log.Infof("Syncing block %d of %d", fromBlock, lastKnownBlock.Uint64())
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
		err = s.processBlockRange(blocks, order)
		if err != nil {
			return lastEthBlockSynced, err
		}
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
			err = s.processBlockRange([]etherman.Block{b}, order)
			if err != nil {
				return lastEthBlockSynced, err
			}
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

	log.Info("Getting trusted state info")
	lastTrustedStateBatchNumber, err := s.zkEVMClient.BatchNumber(s.ctx)
	if err != nil {
		log.Warn("error syncing trusted state. Error: ", err)
		return err
	}

	log.Debug("lastTrustedStateBatchNumber ", lastTrustedStateBatchNumber)
	log.Debug("latestSyncedBatch ", latestSyncedBatch)
	if lastTrustedStateBatchNumber < latestSyncedBatch {
		return nil
	}

	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		batchToSync, err := s.zkEVMClient.BatchByNumber(s.ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		if err != nil {
			log.Warnf("failed to get batch %v from trusted state. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Errorf("error creating db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}

		if err := s.processTrustedBatch(batchToSync, dbTx); err != nil {
			log.Errorf("error processing trusted batch %v: %v", batchNumberToSync, err)
			err := dbTx.Rollback(s.ctx)
			if err != nil {
				log.Errorf("error rolling back db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
				return err
			}
			break
		}

		if err := dbTx.Commit(s.ctx); err != nil {
			log.Errorf("error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}

		batchNumberToSync++
	}

	log.Info("Trusted state fully synchronized")
	return nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Errorf("error creating db transaction to store block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
			return err
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
			log.Errorf("error storing block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[i].BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
		for _, element := range order[blocks[i].BlockHash] {
			switch element.Name {
			case etherman.SequenceBatchesOrder:
				err = s.processSequenceBatches(blocks[i].SequencedBatches[element.Pos], blocks[i].BlockNumber, dbTx)
				if err != nil {
					return err
				}
			case etherman.ForcedBatchesOrder:
				err = s.processForcedBatch(blocks[i].ForcedBatches[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.GlobalExitRootsOrder:
				err = s.processGlobalExitRoot(blocks[i].GlobalExitRoots[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.SequenceForceBatchesOrder:
				err = s.processSequenceForceBatch(blocks[i].SequencedForceBatches[element.Pos], blocks[i], dbTx)
				if err != nil {
					return err
				}
			case etherman.TrustedVerifyBatchOrder:
				err = s.processTrustedVerifyBatches(blocks[i].VerifiedBatches[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.ForkIDsOrder:
				err = s.processForkID(blocks[i].ForkIDs[element.Pos], blocks[i].BlockNumber, dbTx)
				if err != nil {
					return err
				}
			}
		}
		err = dbTx.Commit(s.ctx)
		if err != nil {
			log.Errorf("error committing state to store block. BlockNumber: %d, err: %v", blocks[i].BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[i].BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
	}
	return nil
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
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error resetting the state. Error: ", err)
		return err
	}
	err = s.ethTxManager.Reorg(s.ctx, blockNumber+1, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back eth tx manager when reorg detected. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error processing reorg on eth tx manager. Error: ", err)
		return err
	}
	err = dbTx.Commit(s.ctx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
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
			log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %v", latestBlock.BlockNumber, err)
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
				log.Errorf("error creating db transaction to get prevoius blocks")
				return nil, err
			}
			latestBlock, err = s.state.GetPreviousBlock(s.ctx, depth, dbTx)
			errC := dbTx.Commit(s.ctx)
			if errC != nil {
				log.Errorf("error committing dbTx, err: %v", errC)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return nil, rollbackErr
				}
				log.Errorf("error committing dbTx, err: %v", errC)
				return nil, errC
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

func (s *ClientSynchronizer) checkTrustedState(batch state.Batch, tBatch *state.Batch, newRoot common.Hash, dbTx pgx.Tx) bool {
	//Compare virtual state with trusted state
	var reorgReasons strings.Builder
	if newRoot != tBatch.StateRoot {
		log.Warnf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String()))
	}
	if hex.EncodeToString(batch.BatchL2Data) != hex.EncodeToString(tBatch.BatchL2Data) {
		log.Warnf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data))
		reorgReasons.WriteString(fmt.Sprintf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data)))
	}
	if batch.GlobalExitRoot.String() != tBatch.GlobalExitRoot.String() {
		log.Warnf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String()))
	}
	if batch.Timestamp.Unix() != tBatch.Timestamp.Unix() {
		log.Warnf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix())
		reorgReasons.WriteString(fmt.Sprintf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix()))
	}
	if batch.Coinbase.String() != tBatch.Coinbase.String() {
		log.Warnf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String()))
	}

	if reorgReasons.Len() > 0 {
		reason := reorgReasons.String()
		log.Warnf("Trusted Reorg detected for Batch Number: %d. Reasons: %s", tBatch.BatchNumber, reason)
		if s.isTrustedSequencer {
			for {
				log.Error("TRUSTED REORG DETECTED! Batch: ", batch.BatchNumber)
				time.Sleep(5 * time.Second) //nolint:gomnd
			}
		}
		// Store trusted reorg register
		tr := state.TrustedReorg{
			BatchNumber: tBatch.BatchNumber,
			Reason:      reason,
		}
		err := s.state.AddTrustedReorg(s.ctx, &tr, dbTx)
		if err != nil {
			log.Error("error storing tursted reorg register into the db. Error: ", err)
		}
		return true
	}
	return false
}

func (s *ClientSynchronizer) processForkID(forkID etherman.ForkID, blockNumber uint64, dbTx pgx.Tx) error {
	//If the forkID.batchnumber is a future batch
	latestBatchNumber, err := s.state.GetLastBatchNumber(s.ctx, dbTx)
	if err != nil {
		log.Error("error getting last batch number. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if latestBatchNumber < forkID.BatchNumber { //If the forkID will start in a future batch
		// Read Fork ID FROM POE SC
		forkIDIntervals, err := s.etherMan.GetForks(s.ctx)
		if err != nil || len(forkIDIntervals) == 0 {
			log.Error("error getting all forkIDs: ", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
		// Update forkID intervals in the state
		s.state.UpdateForkIDIntervals(forkIDIntervals)
		return nil
	}

	// If forkID affects to a batch from the past. State must be reseted.
	log.Debugf("ForkID: %d, Reverting synchronization to batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	count, err := s.state.GetForkIDTrustedReorgCount(s.ctx, forkID.ForkID, forkID.Version, dbTx)
	if err != nil {
		log.Error("error getting ForkIDTrustedReorg. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state get forkID trusted state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if count > 0 { // If the forkID reset was already done
		return nil
	}

	// Read Fork ID FROM POE SC
	forkIDIntervals, err := s.etherMan.GetForks(s.ctx)
	if err != nil || len(forkIDIntervals) == 0 {
		log.Error("error getting all forkIDs: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	//Reset DB
	err = s.state.ResetForkID(s.ctx, forkID.BatchNumber+1, forkID.ForkID, forkID.Version, dbTx)
	if err != nil {
		log.Error("error resetting the state. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	err = dbTx.Commit(s.ctx)
	if err != nil {
		log.Error("error committing the resetted state. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	// Update forkID intervals in the state
	s.state.UpdateForkIDIntervals(forkIDIntervals)

	return fmt.Errorf("new ForkID detected, reseting synchronizarion")
}

func (s *ClientSynchronizer) processSequenceBatches(sequencedBatches []etherman.SequencedBatch, blockNumber uint64, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	for _, sbatch := range sequencedBatches {
		virtualBatch := state.VirtualBatch{
			BatchNumber:   sbatch.BatchNumber,
			TxHash:        sbatch.TxHash,
			Coinbase:      sbatch.Coinbase,
			BlockNumber:   blockNumber,
			SequencerAddr: sbatch.SequencerAddr,
		}
		batch := state.Batch{
			BatchNumber:    sbatch.BatchNumber,
			GlobalExitRoot: sbatch.GlobalExitRoot,
			Timestamp:      time.Unix(int64(sbatch.Timestamp), 0),
			Coinbase:       sbatch.Coinbase,
			BatchL2Data:    sbatch.Transactions,
		}
		// ForcedBatch must be processed
		if sbatch.MinForcedTimestamp > 0 { // If this is true means that the batch is forced
			log.Debug("FORCED BATCH SEQUENCED!")
			// Read forcedBatches from db
			forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, 1, dbTx)
			if err != nil {
				log.Errorf("error getting forcedBatches. BatchNumber: %d", virtualBatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			if len(forcedBatches) == 0 {
				log.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", sbatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
			}
			if uint64(forcedBatches[0].ForcedAt.Unix()) != sbatch.MinForcedTimestamp ||
				forcedBatches[0].GlobalExitRoot != sbatch.GlobalExitRoot ||
				common.Bytes2Hex(forcedBatches[0].RawTxsData) != common.Bytes2Hex(sbatch.Transactions) {
				log.Warnf("ForcedBatch stored: %+v", forcedBatches)
				log.Warnf("ForcedBatch sequenced received: %+v", sbatch)
				log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", virtualBatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
			}
			batch.ForcedBatchNum = &forcedBatches[0].ForcedBatchNumber
		}

		// Now we need to check the batch. ForcedBatches should be already stored in the batch table because this is done by the sequencer
		processCtx := state.ProcessingContext{
			BatchNumber:    batch.BatchNumber,
			Coinbase:       batch.Coinbase,
			Timestamp:      batch.Timestamp,
			GlobalExitRoot: batch.GlobalExitRoot,
			ForcedBatchNum: batch.ForcedBatchNum,
		}

		var newRoot common.Hash

		// First get trusted batch from db
		tBatch, err := s.state.GetBatchByNumber(s.ctx, batch.BatchNumber, dbTx)
		if err != nil {
			if errors.Is(err, state.ErrNotFound) || errors.Is(err, state.ErrStateNotSynchronized) {
				log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
				// If it is not found, store batch
				newStateRoot, err := s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx, state.SynchronizerCallerLabel)
				if err != nil {
					log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
					rollbackErr := dbTx.Rollback(s.ctx)
					if rollbackErr != nil {
						log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
						return rollbackErr
					}
					log.Errorf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
					return err
				}
				newRoot = newStateRoot
				tBatch = &batch
				tBatch.StateRoot = newRoot
			} else {
				log.Error("error checking trusted state: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		} else {
			// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
			p, err := s.state.ExecuteBatch(s.ctx, batch, false, dbTx)
			if err != nil {
				log.Errorf("error executing L1 batch: %+v, error: %v", batch, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			newRoot = common.BytesToHash(p.NewStateRoot)
			accumulatedInputHash := common.BytesToHash(p.NewAccInputHash)

			//AddAccumulatedInputHash
			err = s.state.AddAccumulatedInputHash(s.ctx, batch.BatchNumber, accumulatedInputHash, dbTx)
			if err != nil {
				log.Errorf("error adding accumulatedInputHash for batch: %d. Error; %v", batch.BatchNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		}

		// Call the check trusted state method to compare trusted and virtual state
		status := s.checkTrustedState(batch, tBatch, newRoot, dbTx)
		if status {
			// Reorg Pool
			err := s.reorgPool(dbTx)
			if err != nil {
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", tBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error: %v. BatchNumber: %d, BlockNumber: %d", err, tBatch.BatchNumber, blockNumber)
				return err
			}

			// Reset trusted state
			previousBatchNumber := batch.BatchNumber - 1
			log.Warnf("Trusted reorg detected, discarding batches until batchNum %d", previousBatchNumber)
			err = s.state.ResetTrustedState(s.ctx, previousBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
			if err != nil {
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				return err
			}
			_, err = s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx, state.SynchronizerCallerLabel)
			if err != nil {
				log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				return err
			}
		}

		// Store virtualBatch
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
			return err
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
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", blockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processSequenceForceBatch(sequenceForceBatch []etherman.SequencedForceBatch, block etherman.Block, dbTx pgx.Tx) error {
	if len(sequenceForceBatch) == 0 {
		log.Warn("Empty sequenceForceBatch array detected, ignoring...")
		return nil
	}
	// First, get last virtual batch number
	lastVirtualizedBatchNumber, err := s.state.GetLastVirtualBatchNum(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	// Second, reset trusted state
	err = s.state.ResetTrustedState(s.ctx, lastVirtualizedBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
	if err != nil {
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		return err
	}
	// Read forcedBatches from db
	forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, len(sequenceForceBatch), dbTx)
	if err != nil {
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d", block.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	if len(sequenceForceBatch) != len(forcedBatches) {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %v", block.BlockNumber, rollbackErr)
			return rollbackErr
		}
		log.Error("error number of forced batches doesn't match")
		return fmt.Errorf("error number of forced batches doesn't match")
	}
	for i, fbatch := range sequenceForceBatch {
		if uint64(forcedBatches[i].ForcedAt.Unix()) != fbatch.MinForcedTimestamp ||
			forcedBatches[i].GlobalExitRoot != fbatch.GlobalExitRoot ||
			common.Bytes2Hex(forcedBatches[i].RawTxsData) != common.Bytes2Hex(fbatch.Transactions) {
			log.Warnf("ForcedBatch stored: %+v", forcedBatches)
			log.Warnf("ForcedBatch sequenced received: %+v", fbatch)
			log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", fbatch.BatchNumber, block.BlockNumber, rollbackErr)
				return rollbackErr
			}
			return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
		}
		virtualBatch := state.VirtualBatch{
			BatchNumber:   fbatch.BatchNumber,
			TxHash:        fbatch.TxHash,
			Coinbase:      fbatch.Coinbase,
			SequencerAddr: fbatch.Coinbase,
			BlockNumber:   block.BlockNumber,
		}
		batch := state.ProcessingContext{
			BatchNumber:    fbatch.BatchNumber,
			GlobalExitRoot: fbatch.GlobalExitRoot,
			Timestamp:      block.ReceivedAt,
			Coinbase:       fbatch.Coinbase,
			ForcedBatchNum: &forcedBatches[i].ForcedBatchNumber,
		}
		// Process batch
		_, err := s.state.ProcessAndStoreClosedBatch(s.ctx, batch, forcedBatches[i].RawTxsData, dbTx, state.SynchronizerCallerLabel)
		if err != nil {
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			return err
		}
		// Store virtualBatch
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			return err
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
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processForcedBatch(forcedBatch etherman.ForcedBatch, dbTx pgx.Tx) error {
	// Store forced batch into the db
	forcedB := state.ForcedBatch{
		BlockNumber:       forcedBatch.BlockNumber,
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
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", forcedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d, error: %v", forcedBatch.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processGlobalExitRoot(globalExitRoot etherman.GlobalExitRoot, dbTx pgx.Tx) error {
	// Store GlobalExitRoot
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := s.state.AddGlobalExitRoot(s.ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the globalExitRoot in processGlobalExitRoot. BlockNumber: %d", globalExitRoot.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", globalExitRoot.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %v", globalExitRoot.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processTrustedVerifyBatches(lastVerifiedBatch etherman.VerifiedBatch, dbTx pgx.Tx) error {
	lastVBatch, err := s.state.GetLastVerifiedBatch(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d", lastVerifiedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing synced blockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
		return err
	}
	nbatches := lastVerifiedBatch.BatchNumber - lastVBatch.BatchNumber
	batch, err := s.state.GetBatchByNumber(s.ctx, lastVerifiedBatch.BatchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing blockNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing blockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BatchNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing blockNumber: %d, error: %v", lastVerifiedBatch.BatchNumber, err)
		return err
	}

	// Checks that calculated state root matches with the verified state root in the smc
	if batch.StateRoot != lastVerifiedBatch.StateRoot {
		log.Warn("nbatches: ", nbatches)
		log.Warnf("Batch from db: %+v", batch)
		log.Warnf("Verified Batch: %+v", lastVerifiedBatch)
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing blockNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing blockNumber: %d, rollbackErr: %v", lastVerifiedBatch.BatchNumber, rollbackErr)
			return rollbackErr
		}
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing blockNumber: %d", lastVerifiedBatch.BatchNumber)
		return fmt.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing blockNumber: %d", lastVerifiedBatch.BatchNumber)
	}
	var i uint64
	for i = 1; i <= nbatches; i++ {
		verifiedB := state.VerifiedBatch{
			BlockNumber: lastVerifiedBatch.BlockNumber,
			BatchNumber: lastVBatch.BatchNumber + i,
			Aggregator:  lastVerifiedBatch.Aggregator,
			StateRoot:   lastVerifiedBatch.StateRoot,
			TxHash:      lastVerifiedBatch.TxHash,
			IsTrusted:   true,
		}
		err = s.state.AddVerifiedBatch(s.ctx, &verifiedB, dbTx)
		if err != nil {
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. verifiedBatch: %+v, lastVerifiedBatch: %+v", verifiedB, lastVerifiedBatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. BlockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
			return err
		}
	}
	return nil
}

func (s *ClientSynchronizer) processTrustedBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("processing trusted batch: %v", trustedBatch.Number)
	txs := []ethTypes.Transaction{}
	for _, transaction := range trustedBatch.Transactions {
		tx := transaction.Tx.CoreTx()
		txs = append(txs, *tx)
	}
	trustedBatchL2Data, err := state.EncodeTransactions(txs)
	if err != nil {
		return err
	}

	batch, err := s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number), nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
		return err
	}

	// check if batch needs to be synchronized
	if batch != nil {
		matchNumber := batch.BatchNumber == uint64(trustedBatch.Number)
		matchGER := batch.GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
		matchLER := batch.LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
		matchSR := batch.StateRoot.String() == trustedBatch.StateRoot.String()
		matchCoinbase := batch.Coinbase.String() == trustedBatch.Coinbase.String()
		matchTimestamp := uint64(batch.Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
		matchL2Data := hex.EncodeToString(batch.BatchL2Data) == hex.EncodeToString(trustedBatchL2Data)

		if matchNumber && matchGER && matchLER && matchSR &&
			matchCoinbase && matchTimestamp && matchL2Data {
			log.Debugf("batch %v already synchronized", trustedBatch.Number)
			return nil
		}
		log.Infof("batch %v needs to be updated", trustedBatch.Number)
	} else {
		log.Infof("batch %v needs to be synchronized", trustedBatch.Number)
	}

	log.Debugf("resetting trusted state from batch %v", trustedBatch.Number)
	previousBatchNumber := trustedBatch.Number - 1
	if err := s.state.ResetTrustedState(s.ctx, uint64(previousBatchNumber), dbTx); err != nil {
		log.Errorf("failed to reset trusted state", trustedBatch.Number)
		return err
	}

	log.Debugf("opening batch %v", trustedBatch.Number)
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
	}
	if err := s.state.OpenBatch(s.ctx, processCtx, dbTx); err != nil {
		log.Errorf("error opening batch %d", trustedBatch.Number)
		return err
	}

	log.Debugf("processing sequencer for batch %v", trustedBatch.Number)

	processBatchResp, err := s.state.ProcessSequencerBatch(s.ctx, uint64(trustedBatch.Number), trustedBatchL2Data, state.SynchronizerCallerLabel, dbTx)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %d", trustedBatch.Number)
		return err
	}

	log.Debugf("storing transactions for batch %v", trustedBatch.Number)
	if err = s.state.StoreTransactions(s.ctx, uint64(trustedBatch.Number), processBatchResp.Responses, dbTx); err != nil {
		log.Errorf("failed to store transactions for batch: %d", trustedBatch.Number)
		return err
	}

	log.Debug("trustedBatch.StateRoot ", trustedBatch.StateRoot)
	isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
	if isBatchClosed {
		receipt := state.ProcessingReceipt{
			BatchNumber:   uint64(trustedBatch.Number),
			StateRoot:     processBatchResp.NewStateRoot,
			LocalExitRoot: processBatchResp.NewLocalExitRoot,
			BatchL2Data:   trustedBatchL2Data,
			AccInputHash:  trustedBatch.AccInputHash,
		}
		log.Debugf("closing batch %v", trustedBatch.Number)
		if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
			log.Errorf("error closing batch %d", trustedBatch.Number)
			return err
		}
	}

	log.Infof("batch %v synchronized", trustedBatch.Number)
	return nil
}

func (s *ClientSynchronizer) reorgPool(dbTx pgx.Tx) error {
	latestBatchNum, err := s.etherMan.GetLatestBatchNumber()
	if err != nil {
		log.Error("error getting the latestBatchNumber virtualized in the smc. Error: ", err)
		return err
	}
	batchNumber := latestBatchNum + 1
	// Get transactions that have to be included in the pool again
	txs, err := s.state.GetReorgedTransactions(s.ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting txs from trusted state. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Reorged transactions: ", txs)

	// Remove txs from the pool
	err = s.pool.DeleteReorgedTransactions(s.ctx, txs)
	if err != nil {
		log.Errorf("error deleting txs from the pool. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Delete reorged transactions")

	// Add txs to the pool
	for _, tx := range txs {
		// Insert tx in WIP status to avoid the sequencer to grab them before it gets restarted
		// When the sequencer restarts, it will update the status to pending non-wip
		err = s.pool.StoreTx(s.ctx, *tx, "", true)
		if err != nil {
			log.Errorf("error storing tx into the pool again. TxHash: %s. BatchNumber: %d, error: %v", tx.Hash().String(), batchNumber, err)
			return err
		}
		log.Debug("Reorged transactions inserted in the pool: ", tx.Hash())
	}
	return nil
}
