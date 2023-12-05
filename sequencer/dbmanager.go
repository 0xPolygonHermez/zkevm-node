package sequencer

import (
	"context"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	datastreamChannelMultiplier = 2
)

// Pool Loader and DB Updater
type dbManager struct {
	cfg                          DBManagerCfg
	txPool                       txPool
	state                        stateInterface
	worker                       workerInterface
	l2ReorgCh                    chan L2ReorgEvent
	ctx                          context.Context
	batchConstraints             state.BatchConstraintsCfg
	numberOfStateInconsistencies uint64
	streamServer                 *datastreamer.StreamServer
	dataToStream                 chan state.DSL2FullBlock
}

func (d *dbManager) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	return d.state.GetBatchByNumber(ctx, batchNumber, dbTx)
}

// ClosingBatchParameters contains the necessary parameters to close a batch
type ClosingBatchParameters struct {
	BatchNumber          uint64
	StateRoot            common.Hash
	LocalExitRoot        common.Hash
	AccInputHash         common.Hash
	Txs                  []types.Transaction
	BatchResources       state.BatchResources
	ClosingReason        state.ClosingReason
	EffectivePercentages []uint8
}

func newDBManager(ctx context.Context, config DBManagerCfg, txPool txPool, stateInterface stateInterface, worker *Worker, closingSignalCh ClosingSignalCh, batchConstraints state.BatchConstraintsCfg) *dbManager {
	numberOfReorgs, err := stateInterface.CountReorgs(ctx, nil)
	if err != nil {
		log.Error("failed to get number of reorgs: %v", err)
	}

	return &dbManager{ctx: ctx, cfg: config, txPool: txPool,
		state: stateInterface, worker: worker, l2ReorgCh: closingSignalCh.L2ReorgCh,
		batchConstraints: batchConstraints, numberOfStateInconsistencies: numberOfReorgs,
		dataToStream: make(chan state.DSL2FullBlock, batchConstraints.MaxTxsPerBatch*datastreamChannelMultiplier)}
}

// Start stars the dbManager routines
func (d *dbManager) Start() {
	go d.loadFromPool()
	go func() {
		for {
			time.Sleep(d.cfg.L2ReorgRetrievalInterval.Duration)
			d.checkStateInconsistency()
		}
	}()
	if d.streamServer != nil {
		go d.sendDataToStreamer()
	}
}

// GetLastBatchNumber get the latest batch number from state
func (d *dbManager) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	return d.state.GetLastBatchNumber(ctx, nil)
}

// OpenBatch opens a new batch to star processing transactions
func (d *dbManager) OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error {
	return d.state.OpenBatch(ctx, processingContext, dbTx)
}

// CreateFirstBatch is using during genesis
func (d *dbManager) CreateFirstBatch(ctx context.Context, sequencerAddress common.Address) state.ProcessingContext {
	processingCtx := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: state.ZeroHash,
	}
	dbTx, err := d.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for opening a batch, err: %v", err)
		return processingCtx
	}
	err = d.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to open a batch, err: %v", err)
		return processingCtx
	}
	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when opening batch, err: %v", err)
		return processingCtx
	}
	return processingCtx
}

// checkStateInconsistency checks if state inconsistency happened
func (d *dbManager) checkStateInconsistency() {
	stateInconsistenciesDetected, err := d.state.CountReorgs(d.ctx, nil)
	if err != nil {
		log.Error("failed to get number of reorgs: %v", err)
		return
	}

	if stateInconsistenciesDetected != d.numberOfStateInconsistencies {
		log.Warnf("New State Inconsistency detected")
		d.l2ReorgCh <- L2ReorgEvent{}
	}
}

// loadFromPool keeps loading transactions from the pool
func (d *dbManager) loadFromPool() {
	for {
		time.Sleep(d.cfg.PoolRetrievalInterval.Duration)

		poolTransactions, err := d.txPool.GetNonWIPPendingTxs(d.ctx)
		if err != nil && err != pool.ErrNotFound {
			log.Errorf("load tx from pool: %v", err)
		}

		for _, tx := range poolTransactions {
			err := d.addTxToWorker(tx)
			if err != nil {
				log.Errorf("error adding transaction to worker: %v", err)
			}
		}
	}
}

// sendDataToStreamer sends data to the data stream server
func (d *dbManager) sendDataToStreamer() {
	var err error
	for {
		// Read error from previous iteration
		if err != nil {
			err = d.streamServer.RollbackAtomicOp()
			if err != nil {
				log.Errorf("failed to rollback atomic op: %v", err)
			}
			d.streamServer = nil
		}

		// Read data from channel
		fullL2Block := <-d.dataToStream

		l2Block := fullL2Block
		l2Transactions := fullL2Block.Txs

		if d.streamServer != nil {
			err = d.streamServer.StartAtomicOp()
			if err != nil {
				log.Errorf("failed to start atomic op for l2block %v: %v ", l2Block.L2BlockNumber, err)
				continue
			}

			bookMark := state.DSBookMark{
				Type:          state.BookMarkTypeL2Block,
				L2BlockNumber: l2Block.L2BlockNumber,
			}

			_, err = d.streamServer.AddStreamBookmark(bookMark.Encode())
			if err != nil {
				log.Errorf("failed to add stream bookmark for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			blockStart := state.DSL2BlockStart{
				BatchNumber:    l2Block.BatchNumber,
				L2BlockNumber:  l2Block.L2BlockNumber,
				Timestamp:      l2Block.Timestamp,
				GlobalExitRoot: l2Block.GlobalExitRoot,
				Coinbase:       l2Block.Coinbase,
				ForkID:         l2Block.ForkID,
			}

			_, err = d.streamServer.AddStreamEntry(state.EntryTypeL2BlockStart, blockStart.Encode())
			if err != nil {
				log.Errorf("failed to add stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			for _, l2Transaction := range l2Transactions {
				_, err = d.streamServer.AddStreamEntry(state.EntryTypeL2Tx, l2Transaction.Encode())
				if err != nil {
					log.Errorf("failed to add l2tx stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
					continue
				}
			}

			blockEnd := state.DSL2BlockEnd{
				L2BlockNumber: l2Block.L2BlockNumber,
				BlockHash:     l2Block.BlockHash,
				StateRoot:     l2Block.StateRoot,
			}

			_, err = d.streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, blockEnd.Encode())
			if err != nil {
				log.Errorf("failed to add stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			err = d.streamServer.CommitAtomicOp()
			if err != nil {
				log.Errorf("failed to commit atomic op for l2block %v: %v ", l2Block.L2BlockNumber, err)
				continue
			}
		}
	}
}

func (d *dbManager) addTxToWorker(tx pool.Transaction) error {
	txTracker, err := d.worker.NewTxTracker(tx.Transaction, tx.ZKCounters, tx.IP)
	if err != nil {
		return err
	}
	replacedTx, dropReason := d.worker.AddTxTracker(d.ctx, txTracker)
	if dropReason != nil {
		failedReason := dropReason.Error()
		return d.txPool.UpdateTxStatus(d.ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
	} else {
		if replacedTx != nil {
			failedReason := ErrReplacedTransaction.Error()
			error := d.txPool.UpdateTxStatus(d.ctx, replacedTx.Hash, pool.TxStatusFailed, false, &failedReason)
			if error != nil {
				log.Warnf("error when setting as failed replacedTx(%s)", replacedTx.HashStr)
			}
		}
		return d.txPool.UpdateTxWIPStatus(d.ctx, tx.Hash(), true)
	}
}

// BeginStateTransaction starts a db transaction in the state
func (d *dbManager) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	return d.state.BeginStateTransaction(ctx)
}

// DeleteTransactionFromPool deletes a transaction from the pool
func (d *dbManager) DeleteTransactionFromPool(ctx context.Context, txHash common.Hash) error {
	return d.txPool.DeleteTransactionByHash(ctx, txHash)
}

func (d *dbManager) BuildChangeL2Block(deltaTimestamp uint32, l1InfoTreeIndex uint32) []byte {
	changeL2Block := []byte{}
	// changeL2Block transaction mark
	changeL2Block = append(changeL2Block, changeL2BlockMark...)
	// changeL2Block deltaTimeStamp
	deltaTimestampBytes := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(deltaTimestampBytes, deltaTimestamp)
	changeL2Block = append(changeL2Block, deltaTimestampBytes...)
	// changeL2Block l1InfoTreeIndexBytes
	l1InfoTreeIndexBytes := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(l1InfoTreeIndexBytes, uint32(l1InfoTreeIndex))
	changeL2Block = append(changeL2Block, l1InfoTreeIndexBytes...)

	return changeL2Block
}

// StoreL2Block stores a l2 block into the state, update the batch data and changes the status of the block txs in the pool
func (d *dbManager) StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) error {
	return d.state.StoreL2Block(ctx, batchNumber, l2Block, txsEGPLog, dbTx)
}

// GetWIPBatch returns ready WIP batch
func (d *dbManager) GetWIPBatch(ctx context.Context) (*Batch, error) {
	const two = 2
	var lastBatch, previousLastBatch *state.Batch
	dbTx, err := d.BeginStateTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("failed to commit GetWIPBatch: %v", err)
		}
	}()

	lastBatches, err := d.state.GetLastNBatches(ctx, two, dbTx)
	if err != nil {
		return nil, err
	}

	lastBatch = lastBatches[0]
	if len(lastBatches) > 1 {
		previousLastBatch = lastBatches[1]
	}

	forkID := d.state.GetForkIDByBatchNumber(lastBatch.BatchNumber)
	lastBatchTxs, _, _, err := state.DecodeTxs(lastBatch.BatchL2Data, forkID)
	if err != nil {
		return nil, err
	}
	lastBatch.Transactions = lastBatchTxs

	var lastStateRoot common.Hash
	// If the last batch have no txs, the stateRoot can not be retrieved from the l2block because there is no tx.
	// In this case, the stateRoot must be gotten from the previousLastBatch
	if len(lastBatchTxs) == 0 && previousLastBatch != nil {
		lastStateRoot = previousLastBatch.StateRoot
	} else {
		lastStateRoot, err = d.state.GetLastStateRoot(ctx, dbTx)
		if err != nil {
			return nil, err
		}
	}

	wipBatch := &Batch{
		batchNumber:    lastBatch.BatchNumber,
		coinbase:       lastBatch.Coinbase,
		localExitRoot:  lastBatch.LocalExitRoot,
		timestamp:      lastBatch.Timestamp,
		globalExitRoot: lastBatch.GlobalExitRoot,
		countOfTxs:     len(lastBatch.Transactions),
	}

	// Init counters to MAX values
	var totalBytes uint64 = d.batchConstraints.MaxBatchBytesSize
	var batchZkCounters = state.ZKCounters{
		GasUsed:              d.batchConstraints.MaxCumulativeGasUsed,
		UsedKeccakHashes:     d.batchConstraints.MaxKeccakHashes,
		UsedPoseidonHashes:   d.batchConstraints.MaxPoseidonHashes,
		UsedPoseidonPaddings: d.batchConstraints.MaxPoseidonPaddings,
		UsedMemAligns:        d.batchConstraints.MaxMemAligns,
		UsedArithmetics:      d.batchConstraints.MaxArithmetics,
		UsedBinaries:         d.batchConstraints.MaxBinaries,
		UsedSteps:            d.batchConstraints.MaxSteps,
	}

	isClosed, err := d.IsBatchClosed(ctx, lastBatch.BatchNumber)
	if err != nil {
		return nil, err
	}

	if isClosed {
		wipBatch.batchNumber = lastBatch.BatchNumber + 1
		wipBatch.stateRoot = lastBatch.StateRoot
		wipBatch.initialStateRoot = lastBatch.StateRoot

		processingContext := &state.ProcessingContext{
			BatchNumber:    wipBatch.batchNumber,
			Coinbase:       wipBatch.coinbase,
			Timestamp:      wipBatch.timestamp,
			GlobalExitRoot: wipBatch.globalExitRoot,
		}
		err = d.state.OpenBatch(ctx, *processingContext, dbTx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				log.Errorf(
					"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
					rollbackErr, err,
				)
			}
			log.Errorf("failed to open a batch, err: %v", err)
			return nil, err
		}
		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("failed to commit dbTx when opening batch, err: %v", err)
			return nil, err
		}
	} else {
		wipBatch.stateRoot = lastStateRoot
		wipBatch.initialStateRoot = previousLastBatch.StateRoot
		batchL2DataLen := len(lastBatch.BatchL2Data)

		if batchL2DataLen > 0 {
			wipBatch.countOfTxs = len(lastBatch.Transactions)
			batchToExecute := *lastBatch
			batchToExecute.BatchNumber = wipBatch.batchNumber
			batchResponse, err := d.state.ExecuteBatch(ctx, batchToExecute, false, dbTx)
			if err != nil {
				return nil, err
			}

			zkCounters := &state.ZKCounters{
				GasUsed:              batchResponse.GetCumulativeGasUsed(),
				UsedKeccakHashes:     batchResponse.CntKeccakHashes,
				UsedPoseidonHashes:   batchResponse.CntPoseidonHashes,
				UsedPoseidonPaddings: batchResponse.CntPoseidonPaddings,
				UsedMemAligns:        batchResponse.CntMemAligns,
				UsedArithmetics:      batchResponse.CntArithmetics,
				UsedBinaries:         batchResponse.CntBinaries,
				UsedSteps:            batchResponse.CntSteps,
			}

			err = batchZkCounters.Sub(*zkCounters)
			if err != nil {
				return nil, err
			}

			totalBytes -= uint64(batchL2DataLen)
		} else {
			wipBatch.countOfTxs = 0
		}
	}

	wipBatch.remainingResources = state.BatchResources{ZKCounters: batchZkCounters, Bytes: totalBytes}
	return wipBatch, nil
}

// GetLastClosedBatch gets the latest closed batch from state
func (d *dbManager) GetLastClosedBatch(ctx context.Context) (*state.Batch, error) {
	return d.state.GetLastClosedBatch(ctx, nil)
}

// GetLastBatch gets the latest batch from state
func (d *dbManager) GetLastBatch(ctx context.Context) (*state.Batch, error) {
	batch, err := d.state.GetLastBatch(d.ctx, nil)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

// IsBatchClosed checks if a batch is closed
func (d *dbManager) IsBatchClosed(ctx context.Context, batchNum uint64) (bool, error) {
	return d.state.IsBatchClosed(ctx, batchNum, nil)
}

// GetLastNBatches gets the latest N batches from state
func (d *dbManager) GetLastNBatches(ctx context.Context, numBatches uint) ([]*state.Batch, error) {
	return d.state.GetLastNBatches(ctx, numBatches, nil)
}

// GetLatestGer gets the latest global exit root
func (d *dbManager) GetLatestGer(ctx context.Context, gerFinalityNumberOfBlocks uint64) (state.GlobalExitRoot, time.Time, error) {
	return d.state.GetLatestGer(ctx, gerFinalityNumberOfBlocks)
}

// GetLatestL1InfoRoot gets the latest L1InfoRoot
func (d *dbManager) GetLatestL1InfoRoot(ctx context.Context, maxBlockNumber uint64) (state.L1InfoTreeExitRootStorageEntry, error) {
	return d.state.GetLatestL1InfoRoot(ctx, maxBlockNumber)
}

// CloseBatch closes a batch in the state
func (d *dbManager) CloseBatch(ctx context.Context, params ClosingBatchParameters) error {
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:    params.BatchNumber,
		StateRoot:      params.StateRoot,
		LocalExitRoot:  params.LocalExitRoot,
		AccInputHash:   params.AccInputHash,
		BatchResources: params.BatchResources,
		ClosingReason:  params.ClosingReason,
	}

	forkID := d.state.GetForkIDByBatchNumber(params.BatchNumber)
	batchL2Data, err := state.EncodeTransactions(params.Txs, params.EffectivePercentages, forkID)
	if err != nil {
		return err
	}

	processingReceipt.BatchL2Data = batchL2Data

	dbTx, err := d.BeginStateTransaction(ctx)
	if err != nil {
		return err
	}

	err = d.state.CloseBatch(ctx, processingReceipt, dbTx)
	if err != nil {
		err2 := dbTx.Rollback(ctx)
		if err2 != nil {
			log.Errorf("CloseBatch error rolling back: %v", err2)
		}
		return err
	} else {
		err := dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("CloseBatch error committing: %v", err)
			return err
		}
	}

	return nil
}

// ProcessForcedBatch process a forced batch
func (d *dbManager) ProcessForcedBatch(ForcedBatchNumber uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error) {
	// Open Batch
	processingCtx := state.ProcessingContext{
		BatchNumber:    request.BatchNumber,
		Coinbase:       request.Coinbase,
		Timestamp:      request.Timestamp_V1,
		GlobalExitRoot: request.GlobalExitRoot_V1,
		ForcedBatchNum: &ForcedBatchNumber,
	}
	dbTx, err := d.state.BeginStateTransaction(d.ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for opening a forced batch, err: %v", err)
		return nil, err
	}

	err = d.state.OpenBatch(d.ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(d.ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when opening a forced batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to open a batch, err: %v", err)
		return nil, err
	}

	// Fetch Forced Batch
	forcedBatch, err := d.state.GetForcedBatch(d.ctx, ForcedBatchNumber, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(d.ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when getting forced batch err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to get a forced batch, err: %v", err)
		return nil, err
	}

	// Process Batch
	processBatchResponse, err := d.state.ProcessSequencerBatch(d.ctx, request.BatchNumber, forcedBatch.RawTxsData, request.Caller, dbTx)
	if err != nil {
		log.Errorf("failed to process a forced batch, err: %v", err)
		return nil, err
	}

	// Close Batch
	txsBytes := uint64(0)
	for _, blockResp := range processBatchResponse.BlockResponses {
		for _, resp := range blockResp.TransactionResponses {
			if !resp.ChangesStateRoot {
				continue
			}
			txsBytes += resp.Tx.Size()
		}
	}
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   request.BatchNumber,
		StateRoot:     processBatchResponse.NewStateRoot,
		LocalExitRoot: processBatchResponse.NewLocalExitRoot,
		AccInputHash:  processBatchResponse.NewAccInputHash,
		BatchL2Data:   forcedBatch.RawTxsData,
		BatchResources: state.BatchResources{
			ZKCounters: processBatchResponse.UsedZkCounters,
			Bytes:      txsBytes,
		},
		ClosingReason: state.ForcedBatchClosingReason,
	}

	isClosed := false
	tryToCloseAndCommit := true
	for tryToCloseAndCommit {
		if !isClosed {
			closingErr := d.state.CloseBatch(d.ctx, processingReceipt, dbTx)
			tryToCloseAndCommit = closingErr != nil
			if tryToCloseAndCommit {
				continue
			}
			isClosed = true
		}

		if err := dbTx.Commit(d.ctx); err != nil {
			log.Errorf("failed to commit dbTx when processing a forced batch, err: %v", err)
		}
		tryToCloseAndCommit = err != nil
	}

	return processBatchResponse, nil
}

// GetForcedBatchesSince gets L1 forced batches since timestamp
func (d *dbManager) GetForcedBatchesSince(ctx context.Context, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error) {
	return d.state.GetForcedBatchesSince(ctx, forcedBatchNumber, maxBlockNumber, dbTx)
}

// GetLastL2BlockHeader gets the last l2 block number
func (d *dbManager) GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*state.L2Header, error) {
	return d.state.GetLastL2BlockHeader(ctx, dbTx)
}

func (d *dbManager) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error) {
	return d.state.GetLastBlock(ctx, dbTx)
}

func (d *dbManager) GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*state.L2Block, error) {
	return d.state.GetLastL2Block(ctx, dbTx)
}

func (d *dbManager) GetLastTrustedForcedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	return d.state.GetLastTrustedForcedBatchNumber(ctx, dbTx)
}

func (d *dbManager) GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	return d.state.GetBalanceByStateRoot(ctx, address, root)
}

func (d *dbManager) GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64) (txs []types.Transaction, effectivePercentages []uint8, err error) {
	return d.state.GetTransactionsByBatchNumber(ctx, batchNumber, nil)
}

func (d *dbManager) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool, failedReason *string) error {
	return d.txPool.UpdateTxStatus(ctx, hash, newStatus, isWIP, failedReason)
}

func (d *dbManager) UpdateBatch(ctx context.Context, batchNumber uint64, batchL2Data []byte, localExitRoot common.Hash, dbTx pgx.Tx) error {
	return d.state.UpdateBatch(ctx, batchNumber, batchL2Data, localExitRoot, dbTx)
}

// GetLatestVirtualBatchTimestamp gets last virtual batch timestamp
func (d *dbManager) GetLatestVirtualBatchTimestamp(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	return d.state.GetLatestVirtualBatchTimestamp(ctx, dbTx)
}

// CountReorgs returns the number of reorgs
func (d *dbManager) CountReorgs(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	return d.state.CountReorgs(ctx, dbTx)
}

// FlushMerkleTree persists updates in the Merkle tree
func (d *dbManager) FlushMerkleTree(ctx context.Context) error {
	return d.state.FlushMerkleTree(ctx)
}

// GetGasPrices returns the current L2 Gas Price and L1 Gas Price
func (d *dbManager) GetGasPrices(ctx context.Context) (pool.GasPrices, error) {
	return d.txPool.GetGasPrices(ctx)
}

// GetDefaultMinGasPriceAllowed return the configured DefaultMinGasPriceAllowed value
func (d *dbManager) GetDefaultMinGasPriceAllowed() uint64 {
	return d.txPool.GetDefaultMinGasPriceAllowed()
}

func (d *dbManager) GetL1AndL2GasPrice() (uint64, uint64) {
	return d.txPool.GetL1AndL2GasPrice()
}

// GetStoredFlushID returns the stored flush ID and prover ID
func (d *dbManager) GetStoredFlushID(ctx context.Context) (uint64, string, error) {
	return d.state.GetStoredFlushID(ctx)
}

// GetForcedBatch gets a forced batch by number
func (d *dbManager) GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error) {
	return d.state.GetForcedBatch(ctx, forcedBatchNumber, dbTx)
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func (d *dbManager) GetForkIDByBatchNumber(batchNumber uint64) uint64 {
	return d.state.GetForkIDByBatchNumber(batchNumber)
}

func (d *dbManager) DSSendL2Block(l2Block *L2Block) error {
	blockResponse := l2Block.batchResponse.BlockResponses[0]
	forkID := d.GetForkIDByBatchNumber(l2Block.batchNumber)

	// Send data to streamer
	if d.streamServer != nil {
		l2Block := state.DSL2Block{
			BatchNumber:    l2Block.batchNumber,
			L2BlockNumber:  blockResponse.BlockNumber,
			Timestamp:      int64(blockResponse.Timestamp),
			GlobalExitRoot: blockResponse.BlockInfoRoot, //TODO: is it ok?
			Coinbase:       l2Block.coinbase,
			ForkID:         uint16(forkID),
			BlockHash:      blockResponse.BlockHash,
			StateRoot:      blockResponse.BlockHash, //TODO: in etrog the blockhash is the block root
		}

		l2Transactions := []state.DSL2Transaction{}

		for _, txResponse := range blockResponse.TransactionResponses {
			binaryTxData, err := txResponse.Tx.MarshalBinary()
			if err != nil {
				return err
			}

			l2Transaction := state.DSL2Transaction{
				L2BlockNumber:               blockResponse.BlockNumber,
				EffectiveGasPricePercentage: uint8(txResponse.EffectivePercentage),
				IsValid:                     1,
				EncodedLength:               uint32(len(binaryTxData)),
				Encoded:                     binaryTxData,
			}

			l2Transactions = append(l2Transactions, l2Transaction)
		}

		d.dataToStream <- state.DSL2FullBlock{ //TODO: review channel manages the slice of txs
			DSL2Block: l2Block,
			Txs:       l2Transactions,
		}
	}

	return nil
}
