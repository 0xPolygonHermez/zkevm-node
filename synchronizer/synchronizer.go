package synchronizer

// All this function must be deleted becaue have been move to a l1_executor:
// - pendingFlushID
// - halt
// - reorgPool
// - processSequenceForceBatch
// - processSequenceBatches
// - processForkID
// - checkTrustedState

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_parallel_sync"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1event_orders"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	forkID5 = 5
	// ParallelMode is the value for L1SynchronizationMode to run in parallel mode
	ParallelMode = "parallel"
	// SequentialMode is the value for L1SynchronizationMode to run in sequential mode
	SequentialMode = "sequential"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	isTrustedSequencer bool
	etherMan           EthermanInterface
	latestFlushID      uint64
	// If true the lastFlushID is stored in DB and we don't need to check again
	latestFlushIDIsFulfilled bool
	etherManForL1            []EthermanInterface
	state                    stateInterface
	pool                     poolInterface
	ethTxManager             ethTxManager
	zkEVMClient              zkEVMClientInterface
	eventLog                 *event.EventLog
	ctx                      context.Context
	cancelCtx                context.CancelFunc
	genesis                  state.Genesis
	cfg                      Config
	trustedState             struct {
		lastTrustedBatches []*state.Batch
		lastStateRoot      *common.Hash
	}
	// Id of the 'process' of the executor. Each time that it starts this value changes
	// This value is obtained from the call state.GetStoredFlushID
	// It starts as an empty string and it is filled in the first call
	// later the value is checked to be the same (in function checkFlushID)
	proverID string
	// Previous value returned by state.GetStoredFlushID, is used for decide if write a log or not
	previousExecutorFlushID uint64
	l1SyncOrchestration     *l1_parallel_sync.L1SyncOrchestration
	l1EventProcessors       *processor_manager.L1EventProcessors
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	isTrustedSequencer bool,
	ethMan EthermanInterface,
	etherManForL1 []EthermanInterface,
	st stateInterface,
	pool poolInterface,
	ethTxManager ethTxManager,
	zkEVMClient zkEVMClientInterface,
	eventLog *event.EventLog,
	genesis state.Genesis,
	cfg Config,
	runInDevelopmentMode bool) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	metrics.Register()

	res := &ClientSynchronizer{
		isTrustedSequencer:      isTrustedSequencer,
		state:                   st,
		etherMan:                ethMan,
		etherManForL1:           etherManForL1,
		pool:                    pool,
		ctx:                     ctx,
		cancelCtx:               cancel,
		ethTxManager:            ethTxManager,
		zkEVMClient:             zkEVMClient,
		eventLog:                eventLog,
		genesis:                 genesis,
		cfg:                     cfg,
		proverID:                "",
		previousExecutorFlushID: 0,
		l1SyncOrchestration:     nil,
		l1EventProcessors:       nil,
	}
	res.l1EventProcessors = defaultsL1EventProcessors(res)
	switch cfg.L1SynchronizationMode {
	case ParallelMode:
		log.Info("L1SynchronizationMode is parallel")
		var err error
		res.l1SyncOrchestration, err = newL1SyncParallel(ctx, cfg, etherManForL1, res, runInDevelopmentMode)
		if err != nil {
			log.Fatalf("Can't initialize L1SyncParallel. Error: %s", err)
		}
	case SequentialMode:
		log.Info("L1SynchronizationMode is sequential")
	default:
		log.Fatalf("L1SynchronizationMode is not valid. Valid values are: %s, %s", ParallelMode, SequentialMode)
	}

	return res, nil
}

var waitDuration = time.Duration(0)

func newL1SyncParallel(ctx context.Context, cfg Config, etherManForL1 []EthermanInterface, sync *ClientSynchronizer, runExternalControl bool) (*l1_parallel_sync.L1SyncOrchestration, error) {
	chIncommingRollupInfo := make(chan l1_parallel_sync.L1SyncMessage, cfg.L1ParallelSynchronization.MaxPendingNoProcessedBlocks)
	cfgConsumer := l1_parallel_sync.ConfigConsumer{
		ApplyAfterNumRollupReceived: cfg.L1ParallelSynchronization.PerformanceWarning.ApplyAfterNumRollupReceived,
		AceptableInacctivityTime:    cfg.L1ParallelSynchronization.PerformanceWarning.AceptableInacctivityTime.Duration,
	}
	L1DataProcessor := l1_parallel_sync.NewL1RollupInfoConsumer(cfgConsumer, sync, chIncommingRollupInfo)

	cfgProducer := l1_parallel_sync.ConfigProducer{
		SyncChunkSize:                              cfg.SyncChunkSize,
		TtlOfLastBlockOnL1:                         cfg.L1ParallelSynchronization.RequestLastBlockPeriod.Duration,
		TimeoutForRequestLastBlockOnL1:             cfg.L1ParallelSynchronization.RequestLastBlockTimeout.Duration,
		NumOfAllowedRetriesForRequestLastBlockOnL1: cfg.L1ParallelSynchronization.RequestLastBlockMaxRetries,
		TimeForShowUpStatisticsLog:                 cfg.L1ParallelSynchronization.StatisticsPeriod.Duration,
		TimeOutMainLoop:                            cfg.L1ParallelSynchronization.TimeOutMainLoop.Duration,
		MinTimeBetweenRetriesForRollupInfo:         cfg.L1ParallelSynchronization.RollupInfoRetriesSpacing.Duration,
	}
	// Convert EthermanInterface to l1_sync_parallel.EthermanInterface
	etherManForL1Converted := make([]l1_parallel_sync.L1ParallelEthermanInterface, len(etherManForL1))
	for i, etherMan := range etherManForL1 {
		etherManForL1Converted[i] = etherMan
	}
	l1DataRetriever := l1_parallel_sync.NewL1DataRetriever(cfgProducer, etherManForL1Converted, chIncommingRollupInfo)
	l1SyncOrchestration := l1_parallel_sync.NewL1SyncOrchestration(ctx, l1DataRetriever, L1DataProcessor)
	if runExternalControl {
		log.Infof("Starting external control")
		externalControl := newExternalControl(l1DataRetriever, l1SyncOrchestration)
		externalControl.start()
	}
	return l1SyncOrchestration, nil
}

// CleanTrustedState Clean cache of TrustedBatches and StateRoot
func (s *ClientSynchronizer) CleanTrustedState() {
	s.trustedState.lastTrustedBatches = nil
	s.trustedState.lastStateRoot = nil
}

// IsTrustedSequencer returns true is a running in a trusted sequencer
func (s *ClientSynchronizer) IsTrustedSequencer() bool {
	return s.isTrustedSequencer
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	startInitialization := time.Now()
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Info("Sync started")
	dbTx, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Errorf("error creating db transaction to get latest block. Error: %v", err)
		return err
	}
	lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx, dbTx)
	if err != nil {
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Info("State is empty, verifying genesis block")
			valid, err := s.etherMan.VerifyGenBlockNumber(s.ctx, s.genesis.GenesisBlockNum)
			if err != nil {
				log.Error("error checking genesis block number. Error: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			} else if !valid {
				log.Error("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
			}
			log.Info("Setting genesis block")
			header, err := s.etherMan.HeaderByNumber(s.ctx, big.NewInt(0).SetUint64(s.genesis.GenesisBlockNum))
			if err != nil {
				log.Errorf("error getting l1 block header for block %d. Error: %v", s.genesis.GenesisBlockNum, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			lastEthBlockSynced = &state.Block{
				BlockNumber: header.Number.Uint64(),
				BlockHash:   header.Hash(),
				ParentHash:  header.ParentHash,
				ReceivedAt:  time.Unix(int64(header.Time), 0),
			}
			genesisRoot, err := s.state.SetGenesis(s.ctx, *lastEthBlockSynced, s.genesis, stateMetrics.SynchronizerCallerLabel, dbTx)
			if err != nil {
				log.Error("error setting genesis: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}

			blocks, _, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, lastEthBlockSynced.BlockNumber, &lastEthBlockSynced.BlockNumber)
			if err != nil {
				log.Error("error getting rollupInfoByBlockRange after set the genesis: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			//err = s.processForkID(blocks[0].ForkIDs[0], blocks[0].BlockNumber, dbTx)
			err = s.l1EventProcessors.Process(s.ctx, 1, etherman.Order{Name: etherman.ForkIDsOrder, Pos: 0}, &blocks[0], dbTx)

			if err != nil {
				log.Error("error storing genesis forkID: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			if s.genesis.FirstBatchData != nil {
				log.Info("Initial transaction found in genesis file. Applying...")
				err = s.setInitialBatch(blocks[0].BlockNumber, dbTx)
				if err != nil {
					log.Error("error setting initial tx Batch. BatchNum: ", blocks[0].SequencedBatches[0][0].BatchNumber)
					rollbackErr := dbTx.Rollback(s.ctx)
					if rollbackErr != nil {
						log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[0].BlockNumber, rollbackErr.Error(), err)
						return rollbackErr
					}
					return err
				}
			} else {
				log.Info("No initial transaction found in genesis file")
			}

			if genesisRoot != s.genesis.Root {
				log.Errorf("Calculated newRoot should be %s instead of %s", s.genesis.Root.String(), genesisRoot.String())
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("Calculated newRoot should be %s instead of %s", s.genesis.Root.String(), genesisRoot.String())
			}
			// Waiting for the flushID to be stored
			err = s.checkFlushID(dbTx)
			if err != nil {
				log.Error("error checking genesis flushID: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			log.Debug("Genesis root matches!")
		} else {
			log.Error("unexpected error getting the latest ethereum block. Error: ", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
				return rollbackErr
			}
			return err
		}
	}
	initBatchNumber, err := s.state.GetLastBatchNumber(s.ctx, dbTx)
	if err != nil {
		log.Error("error getting latest batchNumber synced. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	err = s.state.SetInitSyncBatch(s.ctx, initBatchNumber, dbTx)
	if err != nil {
		log.Error("error setting initial batch number. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	if err := dbTx.Commit(s.ctx); err != nil {
		log.Errorf("error committing dbTx, err: %v", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	metrics.InitializationTime(time.Since(startInitialization))

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			start := time.Now()
			latestSequencedBatchNumber, err := s.etherMan.GetLatestBatchNumber()
			if err != nil {
				log.Warn("error getting latest sequenced batch in the rollup. Error: ", err)
				continue
			}
			latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx, nil)
			if err != nil {
				log.Warn("error getting latest batch synced in the db. Error: ", err)
				continue
			}
			// Check the latest verified Batch number in the smc
			lastVerifiedBatchNumber, err := s.etherMan.GetLatestVerifiedBatchNum()
			if err != nil {
				log.Warn("error getting last verified batch in the rollup. Error: ", err)
				continue
			}
			err = s.state.SetLastBatchInfoSeenOnEthereum(s.ctx, latestSequencedBatchNumber, lastVerifiedBatchNumber, nil)
			if err != nil {
				log.Warn("error setting latest batch info into db. Error: ", err)
				continue
			}
			log.Infof("latestSequencedBatchNumber: %d, latestSyncedBatch: %d, lastVerifiedBatchNumber: %d", latestSequencedBatchNumber, latestSyncedBatch, lastVerifiedBatchNumber)
			// Sync trusted state
			if latestSyncedBatch >= latestSequencedBatchNumber {
				startTrusted := time.Now()
				log.Info("Syncing trusted state")
				err = s.syncTrustedState(latestSyncedBatch)
				metrics.FullTrustedSyncTime(time.Since(startTrusted))
				if err != nil {
					log.Warn("error syncing trusted state. Error: ", err)
					s.trustedState.lastTrustedBatches = nil
					s.trustedState.lastStateRoot = nil
					continue
				}
				waitDuration = s.cfg.SyncInterval.Duration
			}
			//Sync L1Blocks
			startL1 := time.Now()
			if s.l1SyncOrchestration != nil && (latestSyncedBatch < latestSequencedBatchNumber || !s.cfg.L1ParallelSynchronization.FallbackToSequentialModeOnSynchronized) {
				log.Infof("Syncing L1 blocks in parallel lastEthBlockSynced=%d", lastEthBlockSynced.BlockNumber)
				lastEthBlockSynced, err = s.syncBlocksParallel(lastEthBlockSynced)
			} else {
				if s.l1SyncOrchestration != nil {
					log.Infof("Switching to sequential mode, stopping parallel sync and deleting object")
					s.l1SyncOrchestration.Abort()
					s.l1SyncOrchestration = nil
				}
				log.Infof("Syncing L1 blocks sequentially lastEthBlockSynced=%d", lastEthBlockSynced.BlockNumber)
				lastEthBlockSynced, err = s.syncBlocksSequential(lastEthBlockSynced)
			}
			metrics.FullL1SyncTime(time.Since(startL1))
			if err != nil {
				log.Warn("error syncing blocks: ", err)
				lastEthBlockSynced, err = s.state.GetLastBlock(s.ctx, nil)
				if err != nil {
					log.Fatal("error getting lastEthBlockSynced to resume the synchronization... Error: ", err)
				}
				if s.l1SyncOrchestration != nil {
					// If have failed execution and get starting point from DB, we must reset parallel sync to this point
					// producer must start requesting this block
					s.l1SyncOrchestration.Reset(lastEthBlockSynced.BlockNumber)
				}
				if s.ctx.Err() != nil {
					continue
				}
			}
			metrics.FullSyncIterationTime(time.Since(start))
			log.Info("L1 state fully synchronized")
		}
	}
}

// This function syncs the node from a specific block to the latest
// lastEthBlockSynced -> last block synced in the db
func (s *ClientSynchronizer) syncBlocksParallel(lastEthBlockSynced *state.Block) (*state.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		log.Infof("reorg detected. Resetting the state from block %v to block %v", lastEthBlockSynced.BlockNumber, block.BlockNumber)
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
			s.l1SyncOrchestration.Reset(lastEthBlockSynced.BlockNumber)
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}
	log.Infof("Starting L1 sync orchestrator in parallel block: %d", lastEthBlockSynced.BlockNumber)
	return s.l1SyncOrchestration.Start(lastEthBlockSynced)
}

// This function syncs the node from a specific block to the latest
func (s *ClientSynchronizer) syncBlocksSequential(lastEthBlockSynced *state.Block) (*state.Block, error) {
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
		// Name can be different in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
		// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
		// array index where this value is.
		start := time.Now()
		blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, fromBlock, &toBlock)
		metrics.ReadL1DataTime(time.Since(start))
		if err != nil {
			return lastEthBlockSynced, err
		}
		start = time.Now()
		err = s.ProcessBlockRange(blocks, order)
		metrics.ProcessL1DataTime(time.Since(start))
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
			err = s.ProcessBlockRange([]etherman.Block{b}, order)
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

	log.Info("syncTrustedState: Getting trusted state info")
	start := time.Now()
	lastTrustedStateBatchNumber, err := s.zkEVMClient.BatchNumber(s.ctx)
	metrics.GetTrustedBatchNumberTime(time.Since(start))
	if err != nil {
		log.Warn("syncTrustedState: error syncing trusted state. Error: ", err)
		return err
	}

	log.Debug("syncTrustedState: lastTrustedStateBatchNumber ", lastTrustedStateBatchNumber)
	log.Debug("syncTrustedState: latestSyncedBatch ", latestSyncedBatch)
	if lastTrustedStateBatchNumber < latestSyncedBatch {
		return nil
	}

	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		if batchNumberToSync == 0 {
			batchNumberToSync++
			continue
		}
		start = time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(s.ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("syncTrustedState: failed to get batch %d from trusted state. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Errorf("syncTrustedState: error creating db transaction to sync trusted batch %d: %v", batchNumberToSync, err)
			return err
		}
		start = time.Now()
		cbatches, lastStateRoot, err := s.processTrustedBatch(batchToSync, dbTx)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("syncTrustedState: error processing trusted batch %d: %v", batchNumberToSync, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back db transaction to sync trusted batch %d: %v", batchNumberToSync, rollbackErr)
				return rollbackErr
			}
			return err
		}
		log.Debug("syncTrustedState: Checking FlushID to commit trustedState data to db")
		err = s.checkFlushID(dbTx)
		if err != nil {
			log.Errorf("syncTrustedState: error checking flushID. Error: %v", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}

		if err := dbTx.Commit(s.ctx); err != nil {
			log.Errorf("syncTrustedState: error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}
		s.trustedState.lastTrustedBatches = cbatches
		s.trustedState.lastStateRoot = lastStateRoot
		batchNumberToSync++
	}

	log.Info("syncTrustedState: Trusted state fully synchronized")
	return nil
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *ClientSynchronizer) ProcessBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error {
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
			batchSequence := l1event_orders.GetSequenceFromL1EventOrder(element.Name, &blocks[i], element.Pos)
			var forkId uint64
			if batchSequence != nil {
				forkId = s.state.GetForkIDByBatchNumber(batchSequence.FromBatchNumber)
				log.Debug("EventOrder:", element.Name, "Batch Sequence: ", batchSequence, "forkId:", forkId)
			} else {
				forkId = s.state.GetForkIDByBlockNumber(blocks[i].BlockNumber)
				log.Debug("EventOrder:", element.Name, "BlockNumber: ", blocks[i].BlockNumber, "forkId:", forkId)
			}
			forkIdTyped := actions.ForkIdType(forkId)
			var err error
			switch element.Name {
			case etherman.SequenceBatchesOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.ForcedBatchesOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.GlobalExitRootsOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.SequenceForceBatchesOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.TrustedVerifyBatchOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.VerifyBatchOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.ForkIDsOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			case etherman.L1InfoTreeOrder:
				err = s.l1EventProcessors.Process(s.ctx, forkIdTyped, element, &blocks[i], dbTx)
			}
			if err != nil {
				log.Error("error: ", err)
				// If any goes wrong we ensure that the state is rollbacked
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
					log.Warnf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[i].BlockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
		}
		log.Debug("Checking FlushID to commit L1 data to db")
		err = s.checkFlushID(dbTx)
		if err != nil {
			log.Errorf("error checking flushID. Error: %v", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
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
	log.Info("Reverting synchronization to block: ", blockNumber)
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
	if s.l1SyncOrchestration != nil {
		s.l1SyncOrchestration.Reset(blockNumber)
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
			err = fmt.Errorf("wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
				latestBlock.BlockNumber, block.NumberU64())
			log.Error("error: ", err)
			return nil, err
		}
		// Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > s.genesis.GenesisBlockNum {
			log.Infof("checkReorg: Bad block %d hashOk %t parentHashOk %t", latestBlock.BlockNumber, block.Hash() == latestBlock.BlockHash, block.ParentHash() == latestBlock.ParentHash)
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
		log.Info("Reorg detected in block: ", latestEthBlockSynced.BlockNumber, " last block OK: ", latestBlock.BlockNumber)
		return latestBlock, nil
	}
	return nil, nil
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}

func (s *ClientSynchronizer) processTrustedBatch(trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, *common.Hash, error) {
	log.Debugf("Processing trusted batch: %d", uint64(trustedBatch.Number))
	trustedBatchL2Data := trustedBatch.BatchL2Data
	batches := s.trustedState.lastTrustedBatches
	log.Debug("len(batches): ", len(batches))
	batches, err := s.getCurrentBatches(batches, trustedBatch, dbTx)
	if err != nil {
		log.Error("error getting currentBatches. Error: ", err)
		return nil, nil, err
	}

	if batches[0] != nil && (((trustedBatch.StateRoot == common.Hash{}) && (batches[0].StateRoot != common.Hash{})) ||
		len(batches[0].BatchL2Data) > len(trustedBatchL2Data)) {
		log.Error("error: inconsistency in data received from trustedNode")
		log.Infof("BatchNumber. stored: %d. synced: %d", batches[0].BatchNumber, uint64(trustedBatch.Number))
		log.Infof("GlobalExitRoot. stored:  %s. synced: %s", batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
		log.Infof("LocalExitRoot. stored:  %s. synced: %s", batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
		log.Infof("StateRoot. stored:  %s. synced: %s", batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
		log.Infof("Coinbase. stored:  %s. synced: %s", batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
		log.Infof("Timestamp. stored:  %d. synced: %d", uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
		log.Infof("BatchL2Data. stored: %s. synced: %s", common.Bytes2Hex(batches[0].BatchL2Data), common.Bytes2Hex(trustedBatchL2Data))
		return nil, nil, fmt.Errorf("error: inconsistency in data received from trustedNode")
	}

	if s.trustedState.lastStateRoot == nil && (batches[0] == nil || (batches[0].StateRoot == common.Hash{})) {
		log.Debug("Setting stateRoot of previous batch. StateRoot: ", batches[1].StateRoot)
		// Previous synchronization incomplete. Needs to reprocess all txs again
		s.trustedState.lastStateRoot = &batches[1].StateRoot
	} else if batches[0] != nil && (batches[0].StateRoot != common.Hash{}) {
		// Previous synchronization completed
		s.trustedState.lastStateRoot = &batches[0].StateRoot
	}

	request := state.ProcessRequest{
		BatchNumber:     uint64(trustedBatch.Number),
		OldStateRoot:    *s.trustedState.lastStateRoot,
		OldAccInputHash: batches[1].AccInputHash,
		Coinbase:        common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp_V1:    time.Unix(int64(trustedBatch.Timestamp), 0),
	}
	// check if batch needs to be synchronized
	if batches[0] != nil {
		if checkIfSynced(batches, trustedBatch) {
			log.Debugf("Batch %d already synchronized", uint64(trustedBatch.Number))
			return batches, s.trustedState.lastStateRoot, nil
		}
		log.Infof("Batch %d needs to be updated", uint64(trustedBatch.Number))

		// Find txs to be processed and included in the trusted state
		if *s.trustedState.lastStateRoot == batches[1].StateRoot {
			prevBatch := uint64(trustedBatch.Number) - 1
			log.Infof("ResetTrustedState: processTrustedBatch: %d Cleaning state until batch:%d  ", trustedBatch.Number, prevBatch)
			// Delete txs that were stored before restart. We need to reprocess all txs because the intermediary stateRoot is only stored in memory
			err := s.state.ResetTrustedState(s.ctx, prevBatch, dbTx)
			if err != nil {
				log.Error("error resetting trusted state. Error: ", err)
				return nil, nil, err
			}
			// All txs need to be processed
			request.Transactions = trustedBatchL2Data
			// Reopen batch
			err = s.openBatch(trustedBatch, dbTx)
			if err != nil {
				log.Error("error openning batch. Error: ", err)
				return nil, nil, err
			}
			request.GlobalExitRoot_V1 = trustedBatch.GlobalExitRoot
			request.Transactions = trustedBatchL2Data
		} else {
			// Only new txs need to be processed
			storedTxs, syncedTxs, _, syncedEfficiencyPercentages, err := s.decodeTxs(trustedBatchL2Data, batches)
			if err != nil {
				return nil, nil, err
			}
			if len(storedTxs) < len(syncedTxs) {
				forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
				txsToBeAdded := syncedTxs[len(storedTxs):]
				if forkID >= forkID5 {
					syncedEfficiencyPercentages = syncedEfficiencyPercentages[len(storedTxs):]
				}
				log.Infof("Processing %d new txs with forkID: %d", len(txsToBeAdded), forkID)

				request.Transactions, err = state.EncodeTransactions(txsToBeAdded, syncedEfficiencyPercentages, forkID)
				if err != nil {
					log.Error("error encoding txs (%d) to be added to the state. Error: %v", len(txsToBeAdded), err)
					return nil, nil, err
				}
				log.Debug("request.Transactions: ", common.Bytes2Hex(request.Transactions))
			} else {
				log.Info("Nothing to sync. Node updated. Checking if it is closed")
				isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
				if isBatchClosed {
					//Sanity check
					if s.trustedState.lastStateRoot != nil && trustedBatch.StateRoot != *s.trustedState.lastStateRoot {
						log.Errorf("batch %d, different batchL2Datas (trustedBatchL2Data: %s, batches[0].BatchL2Data: %s). Decoded txs are len(storedTxs): %d, len(syncedTxs): %d", uint64(trustedBatch.Number), trustedBatchL2Data.Hex(), "0x"+common.Bytes2Hex(batches[0].BatchL2Data), len(storedTxs), len(syncedTxs))
						for _, tx := range storedTxs {
							log.Error("stored txHash : ", tx.Hash())
						}
						for _, tx := range syncedTxs {
							log.Error("synced txHash : ", tx.Hash())
						}
						log.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.trustedState.lastStateRoot, trustedBatch.StateRoot)
						return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.trustedState.lastStateRoot, trustedBatch.StateRoot)
					}
					receipt := state.ProcessingReceipt{
						BatchNumber:   uint64(trustedBatch.Number),
						StateRoot:     trustedBatch.StateRoot,
						LocalExitRoot: trustedBatch.LocalExitRoot,
						BatchL2Data:   trustedBatchL2Data,
						AccInputHash:  trustedBatch.AccInputHash,
					}
					log.Debugf("closing batch %d", uint64(trustedBatch.Number))
					if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
						// This is a workaround to avoid closing a batch that was already closed
						if err.Error() != state.ErrBatchAlreadyClosed.Error() {
							log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
							return nil, nil, err
						} else {
							log.Warnf("CASE 02: the batch [%d] was already closed", uint64(trustedBatch.Number))
							log.Info("batches[0].BatchNumber: ", batches[0].BatchNumber)
							log.Info("batches[0].AccInputHash: ", batches[0].AccInputHash)
							log.Info("batches[0].StateRoot: ", batches[0].StateRoot)
							log.Info("batches[0].LocalExitRoot: ", batches[0].LocalExitRoot)
							log.Info("batches[0].GlobalExitRoot: ", batches[0].GlobalExitRoot)
							log.Info("batches[0].Coinbase: ", batches[0].Coinbase)
							log.Info("batches[0].ForcedBatchNum: ", batches[0].ForcedBatchNum)
							log.Info("####################################")
							log.Info("batches[1].BatchNumber: ", batches[1].BatchNumber)
							log.Info("batches[1].AccInputHash: ", batches[1].AccInputHash)
							log.Info("batches[1].StateRoot: ", batches[1].StateRoot)
							log.Info("batches[1].LocalExitRoot: ", batches[1].LocalExitRoot)
							log.Info("batches[1].GlobalExitRoot: ", batches[1].GlobalExitRoot)
							log.Info("batches[1].Coinbase: ", batches[1].Coinbase)
							log.Info("batches[1].ForcedBatchNum: ", batches[1].ForcedBatchNum)
							log.Info("###############################")
							log.Info("trustedBatch.BatchNumber: ", trustedBatch.Number)
							log.Info("trustedBatch.AccInputHash: ", trustedBatch.AccInputHash)
							log.Info("trustedBatch.StateRoot: ", trustedBatch.StateRoot)
							log.Info("trustedBatch.LocalExitRoot: ", trustedBatch.LocalExitRoot)
							log.Info("trustedBatch.GlobalExitRoot: ", trustedBatch.GlobalExitRoot)
							log.Info("trustedBatch.Coinbase: ", trustedBatch.Coinbase)
							log.Info("trustedBatch.ForcedBatchNum: ", trustedBatch.ForcedBatchNumber)
						}
					}
					batches[0].AccInputHash = trustedBatch.AccInputHash
					batches[0].StateRoot = trustedBatch.StateRoot
					batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
				}
				return batches, &trustedBatch.StateRoot, nil
			}
		}
		// Update batchL2Data
		err := s.state.UpdateBatchL2Data(s.ctx, batches[0].BatchNumber, trustedBatchL2Data, dbTx)
		if err != nil {
			log.Errorf("error opening batch %d", uint64(trustedBatch.Number))
			return nil, nil, err
		}
		batches[0].BatchL2Data = trustedBatchL2Data
		log.Debug("BatchL2Data updated for batch: ", batches[0].BatchNumber)
	} else {
		log.Infof("Batch %d needs to be synchronized", uint64(trustedBatch.Number))
		err := s.openBatch(trustedBatch, dbTx)
		if err != nil {
			log.Error("error openning batch. Error: ", err)
			return nil, nil, err
		}
		request.GlobalExitRoot_V1 = trustedBatch.GlobalExitRoot
		request.Transactions = trustedBatchL2Data
	}

	log.Debugf("Processing sequencer for batch %d", uint64(trustedBatch.Number))

	processBatchResp, err := s.processAndStoreTxs(trustedBatch, request, dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, nil, err
	}

	log.Debug("TrustedBatch.StateRoot ", trustedBatch.StateRoot)
	isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
	if isBatchClosed {
		//Sanity check
		if trustedBatch.StateRoot != processBatchResp.NewStateRoot {
			log.Error("trustedBatchL2Data: ", trustedBatchL2Data)
			log.Error("request.Transactions: ", request.Transactions)
			log.Errorf("batch: %d after processing some txs, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
			return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
		}
		receipt := state.ProcessingReceipt{
			BatchNumber:   uint64(trustedBatch.Number),
			StateRoot:     processBatchResp.NewStateRoot,
			LocalExitRoot: processBatchResp.NewLocalExitRoot,
			BatchL2Data:   trustedBatchL2Data,
			AccInputHash:  trustedBatch.AccInputHash,
		}

		log.Debugf("closing batch %d", uint64(trustedBatch.Number))
		if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
			// This is a workarround to avoid closing a batch that was already closed
			if err.Error() != state.ErrBatchAlreadyClosed.Error() {
				log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
				return nil, nil, err
			} else {
				log.Warnf("CASE 01: batch [%d] was already closed", uint64(trustedBatch.Number))
			}
		}
		log.Info("Batch closed right after processing some tx")
		if batches[0] != nil {
			log.Debug("Updating batches[0] values...")
			batches[0].AccInputHash = trustedBatch.AccInputHash
			batches[0].StateRoot = trustedBatch.StateRoot
			batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
			batches[0].BatchL2Data = trustedBatchL2Data
		}
	}

	log.Infof("Batch %d synchronized", uint64(trustedBatch.Number))
	return batches, &processBatchResp.NewStateRoot, nil
}

func (s *ClientSynchronizer) processAndStoreTxs(trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	processBatchResp, err := s.state.ProcessBatch(s.ctx, request, true)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %v", trustedBatch.Number)
		return nil, err
	}
	s.PendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("Storing %d blocks for batch %v", len(processBatchResp.BlockResponses), trustedBatch.Number)
	if processBatchResp.IsExecutorLevelError {
		log.Warn("executorLevelError detected. Avoid store txs...")
		return processBatchResp, nil
	} else if processBatchResp.IsRomOOCError {
		log.Warn("romOOCError detected. Avoid store txs...")
		return processBatchResp, nil
	}
	for _, block := range processBatchResp.BlockResponses {
		for _, tx := range block.TransactionResponses {
			if state.IsStateRootChanged(executor.RomErrorCode(tx.RomError)) {
				log.Infof("TrustedBatch info: %+v", processBatchResp)
				log.Infof("Storing trusted tx %+v", tx)
				if _, err = s.state.StoreTransaction(s.ctx, uint64(trustedBatch.Number), tx, trustedBatch.Coinbase, uint64(trustedBatch.Timestamp), nil, dbTx); err != nil {
					log.Errorf("failed to store transactions for batch: %v. Tx: %s", trustedBatch.Number, tx.TxHash.String())
					return nil, err
				}
			}
		}
	}
	return processBatchResp, nil
}

func (s *ClientSynchronizer) openBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("Opening batch %d", trustedBatch.Number)
	var batchL2Data []byte = trustedBatch.BatchL2Data
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	if trustedBatch.ForcedBatchNumber != nil {
		fb := uint64(*trustedBatch.ForcedBatchNumber)
		processCtx.ForcedBatchNum = &fb
	}
	err := s.state.OpenBatch(s.ctx, processCtx, dbTx)
	if err != nil {
		log.Error("error opening batch: ", trustedBatch.Number)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) decodeTxs(trustedBatchL2Data types.ArgBytes, batches []*state.Batch) ([]ethTypes.Transaction, []ethTypes.Transaction, []uint8, []uint8, error) {
	forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
	syncedTxs, _, syncedEfficiencyPercentages, err := state.DecodeTxs(trustedBatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, trustedBatchL2Data.Hex())
		return nil, nil, nil, nil, err
	}
	storedTxs, _, storedEfficiencyPercentages, err := state.DecodeTxs(batches[0].BatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding stored txs from trustedstate. Error: %v, batch.BatchL2Data: %s", err, common.Bytes2Hex(batches[0].BatchL2Data))
		return nil, nil, nil, nil, err
	}
	log.Debug("len(storedTxs): ", len(storedTxs))
	log.Debug("len(syncedTxs): ", len(syncedTxs))
	return storedTxs, syncedTxs, storedEfficiencyPercentages, syncedEfficiencyPercentages, nil
}

func checkIfSynced(batches []*state.Batch, trustedBatch *types.Batch) bool {
	matchNumber := batches[0].BatchNumber == uint64(trustedBatch.Number)
	matchGER := batches[0].GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := batches[0].LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := batches[0].StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := batches[0].Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := uint64(batches[0].Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
	matchL2Data := hex.EncodeToString(batches[0].BatchL2Data) == hex.EncodeToString(trustedBatch.BatchL2Data)

	if matchNumber && matchGER && matchLER && matchSR &&
		matchCoinbase && matchTimestamp && matchL2Data {
		return true
	}
	log.Infof("matchNumber %v %d %d", matchNumber, batches[0].BatchNumber, uint64(trustedBatch.Number))
	log.Infof("matchGER %v %s %s", matchGER, batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
	log.Infof("matchLER %v %s %s", matchLER, batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
	log.Infof("matchSR %v %s %s", matchSR, batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
	log.Infof("matchCoinbase %v %s %s", matchCoinbase, batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
	log.Infof("matchTimestamp %v %d %d", matchTimestamp, uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
	log.Infof("matchL2Data %v", matchL2Data)
	return false
}

func (s *ClientSynchronizer) getCurrentBatches(batches []*state.Batch, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
	if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number) != batches[0].BatchNumber) {
		log.Debug("Updating batch[0] value!")
		batch, err := s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number), dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
			return nil, err
		}
		var prevBatch *state.Batch
		if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number-1) != batches[0].BatchNumber) {
			log.Debug("Updating batch[1] value!")
			prevBatch, err = s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number-1), dbTx)
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get prevBatch %v from local trusted state. Error: %v", trustedBatch.Number-1, err)
				return nil, err
			}
		} else {
			prevBatch = batches[0]
		}
		log.Debug("batch: ", batch)
		log.Debug("prevBatch: ", prevBatch)
		batches = []*state.Batch{batch, prevBatch}
	}
	return batches, nil
}

// PendingFlushID is called when a flushID is pending to be stored in the db
func (s *ClientSynchronizer) PendingFlushID(flushID uint64, proverID string) {
	log.Infof("pending flushID: %d", flushID)
	if flushID == 0 {
		log.Fatal("flushID is 0. Please check that prover/executor config parameter dbReadOnly is false")
	}
	s.latestFlushID = flushID
	s.latestFlushIDIsFulfilled = false
	s.updateAndCheckProverID(proverID)
}

// deprecated: use PendingFlushID instead
//
//nolint:unused
func (s *ClientSynchronizer) pendingFlushID(flushID uint64, proverID string) {
	s.PendingFlushID(flushID, proverID)
}

func (s *ClientSynchronizer) updateAndCheckProverID(proverID string) {
	if s.proverID == "" {
		log.Infof("Current proverID is %s", proverID)
		s.proverID = proverID
		return
	}
	if s.proverID != proverID {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			Source:      event.Source_Node,
			Component:   event.Component_Synchronizer,
			Level:       event.Level_Critical,
			EventID:     event.EventID_SynchronizerRestart,
			Description: fmt.Sprintf("proverID changed from %s to %s, restarting Synchronizer ", s.proverID, proverID),
		}

		err := s.eventLog.LogEvent(context.Background(), event)
		if err != nil {
			log.Errorf("error storing event payload: %v", err)
		}

		log.Fatal("restarting synchronizer because executor has been restarted (old=%s, new=%s)", s.proverID, proverID)
	}
}

func (s *ClientSynchronizer) checkFlushID(dbTx pgx.Tx) error {
	if s.latestFlushIDIsFulfilled {
		log.Debugf("no pending flushID, nothing to do. Last pending fulfilled flushID: %d, last executor flushId received: %d", s.latestFlushID, s.latestFlushID)
		return nil
	}
	storedFlushID, proverID, err := s.state.GetStoredFlushID(s.ctx)
	if err != nil {
		log.Error("error getting stored flushID. Error: ", err)
		return err
	}
	if s.previousExecutorFlushID != storedFlushID || s.proverID != proverID {
		log.Infof("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	} else {
		log.Debugf("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	}
	s.updateAndCheckProverID(proverID)
	log.Debugf("storedFlushID (executor reported): %d, latestFlushID (pending): %d", storedFlushID, s.latestFlushID)
	if storedFlushID < s.latestFlushID {
		log.Infof("Synchronized BLOCKED!: Wating for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d", s.latestFlushID, storedFlushID)
		iteration := 0
		start := time.Now()
		for storedFlushID < s.latestFlushID {
			log.Debugf("Waiting for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d iteration:%d elpased:%s",
				s.latestFlushID, storedFlushID, iteration, time.Since(start))
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
			storedFlushID, _, err = s.state.GetStoredFlushID(s.ctx)
			if err != nil {
				log.Error("error getting stored flushID. Error: ", err)
				return err
			}
			iteration++
		}
		log.Infof("Synchronizer resumed, flushID stored: %d", s.latestFlushID)
	}
	log.Infof("Pending Flushid fullfiled: %d, executor have write %d", s.latestFlushID, storedFlushID)
	s.latestFlushIDIsFulfilled = true
	s.previousExecutorFlushID = storedFlushID
	return nil
}

func (s *ClientSynchronizer) setInitialBatch(blockNumber uint64, dbTx pgx.Tx) error {
	log.Debug("Setting initial transaction batch 1")
	// Process FirstTransaction included in batch 1
	batchL2Data := common.Hex2Bytes(s.genesis.FirstBatchData.Transactions[2:])
	processCtx := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       s.genesis.FirstBatchData.Sequencer,
		Timestamp:      time.Unix(int64(s.genesis.FirstBatchData.Timestamp), 0),
		GlobalExitRoot: s.genesis.FirstBatchData.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	_, flushID, proverID, err := s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)
	if err != nil {
		log.Error("error storing batch 1. Error: ", err)
		return err
	}
	s.pendingFlushID(flushID, proverID)

	// Virtualize Batch and add sequence
	virtualBatch1 := state.VirtualBatch{
		BatchNumber:   1,
		TxHash:        state.ZeroHash,
		Coinbase:      s.genesis.FirstBatchData.Sequencer,
		BlockNumber:   blockNumber,
		SequencerAddr: s.genesis.FirstBatchData.Sequencer,
	}
	err = s.state.AddVirtualBatch(s.ctx, &virtualBatch1, dbTx)
	if err != nil {
		log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch1.BatchNumber, s.genesis.GenesisBlockNum, err)
		return err
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: 1,
		ToBatchNumber:   1,
	}
	err = s.state.AddSequence(s.ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		return err
	}
	return nil
}
