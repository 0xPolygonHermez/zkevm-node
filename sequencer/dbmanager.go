package sequencer

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	wait time.Duration = 5
)

// Pool Loader and DB Updater
type dbManager struct {
	cfg              DBManagerCfg
	txPool           txPool
	state            dbManagerStateInterface
	worker           workerInterface
	txsStore         TxsStore
	l2ReorgCh        chan L2ReorgEvent
	ctx              context.Context
	batchConstraints batchConstraints
	numberOfReorgs   uint64
}

func (d *dbManager) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	return d.state.GetBatchByNumber(ctx, batchNumber, dbTx)
}

// ClosingBatchParameters contains the necessary parameters to close a batch
type ClosingBatchParameters struct {
	BatchNumber   uint64
	StateRoot     common.Hash
	LocalExitRoot common.Hash
	AccInputHash  common.Hash
	Txs           []types.Transaction
}

func newDBManager(ctx context.Context, config DBManagerCfg, txPool txPool, state dbManagerStateInterface, worker *Worker, closingSignalCh ClosingSignalCh, txsStore TxsStore, batchConstraints batchConstraints) *dbManager {
	numberOfReorgs, err := state.CountReorgs(ctx, nil)
	if err != nil {
		log.Error("failed to get number of reorgs: %v", err)
	}

	return &dbManager{ctx: ctx, cfg: config, txPool: txPool, state: state, worker: worker, txsStore: txsStore, l2ReorgCh: closingSignalCh.L2ReorgCh, batchConstraints: batchConstraints, numberOfReorgs: numberOfReorgs}
}

// Start stars the dbManager routines
func (d *dbManager) Start() {
	go d.loadFromPool()
	go func() {
		for {
			// TODO: Move this to a config parameter
			time.Sleep(wait * time.Second)
			d.checkIfReorg()
		}
	}()
	go d.storeProcessedTxAndDeleteFromPool()
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

// checkIfReorg checks if a reorg has happened
func (d *dbManager) checkIfReorg() {
	numberOfReorgs, err := d.state.CountReorgs(d.ctx, nil)
	if err != nil {
		log.Error("failed to get number of reorgs: %v", err)
	}

	if numberOfReorgs != d.numberOfReorgs {
		log.Warnf("New L2 reorg detected")
		d.l2ReorgCh <- L2ReorgEvent{}
	}
}

// loadFromPool keeps loading transactions from the pool
func (d *dbManager) loadFromPool() {
	for {
		time.Sleep(d.cfg.PoolRetrievalInterval.Duration)

		poolTransactions, err := d.txPool.GetNonWIPPendingTxs(d.ctx, false, 0)
		if err != nil && err != pool.ErrNotFound {
			log.Errorf("load tx from pool: %v", err)
		}

		for _, tx := range poolTransactions {
			err := d.addTxToWorker(tx, false)
			if err != nil {
				log.Errorf("error adding transaction to worker: %v", err)
			}
		}

		poolClaims, err := d.txPool.GetNonWIPPendingTxs(d.ctx, true, 0)
		if err != nil && err != pool.ErrNotFound {
			log.Errorf("load claims from pool: %v", err)
		}

		for _, tx := range poolClaims {
			err := d.addTxToWorker(tx, true)
			if err != nil {
				log.Errorf("error adding claim to worker: %v", err)
			}
		}
	}
}

func (d *dbManager) addTxToWorker(tx pool.Transaction, isClaim bool) error {
	txTracker, err := d.worker.NewTxTracker(tx.Transaction, isClaim, tx.ZKCounters, tx.IP)
	if err != nil {
		return err
	}
	d.worker.AddTxTracker(d.ctx, txTracker)
	return d.txPool.UpdateTxWIPStatus(d.ctx, tx.Hash(), true)
}

// BeginStateTransaction starts a db transaction in the state
func (d *dbManager) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	return d.state.BeginStateTransaction(ctx)
}

// StoreProcessedTransaction stores a transaction in the state
func (d *dbManager) StoreProcessedTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error {
	return d.state.StoreTransaction(ctx, batchNumber, processedTx, coinbase, timestamp, dbTx)
}

// DeleteTransactionFromPool deletes a transaction from the pool
func (d *dbManager) DeleteTransactionFromPool(ctx context.Context, txHash common.Hash) error {
	return d.txPool.DeleteTransactionByHash(ctx, txHash)
}

// storeProcessedTxAndDeleteFromPool stores a tx into the state and changes it status in the pool
func (d *dbManager) storeProcessedTxAndDeleteFromPool() {
	// TODO: Finish the retry mechanism and error handling
	for {
		txToStore := <-d.txsStore.Ch
		d.checkIfReorg()

		// Flush the state db
		err := d.state.FlushMerkleTree(d.ctx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool. Error flushing state db: %v", err)
		}

		log.Debugf("Storing tx %v", txToStore.txResponse.TxHash)
		dbTx, err := d.BeginStateTransaction(d.ctx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}

		err = d.StoreProcessedTransaction(d.ctx, txToStore.batchNumber, txToStore.txResponse, txToStore.coinbase, txToStore.timestamp, dbTx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}

		// Update batch l2 data
		batch, err := d.state.GetBatchByNumber(d.ctx, txToStore.batchNumber, dbTx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}

		txData, err := state.EncodeTransaction(txToStore.txResponse.Tx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}
		batch.BatchL2Data = append(batch.BatchL2Data, txData...)

		err = d.state.UpdateBatchL2Data(d.ctx, txToStore.batchNumber, batch.BatchL2Data, dbTx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}

		err = dbTx.Commit(d.ctx)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool error committing : %v", err)
		}

		// Change Tx status to selected
		err = d.txPool.UpdateTxStatus(d.ctx, txToStore.txResponse.TxHash, pool.TxStatusSelected, false)
		if err != nil {
			log.Fatalf("StoreProcessedTxAndDeleteFromPool: %v", err)
		}

		log.Infof("StoreProcessedTxAndDeleteFromPool: successfully stored tx: %v for batch: %v", txToStore.txResponse.TxHash.String(), txToStore.batchNumber)
		d.txsStore.Wg.Done()
	}
}

// GetWIPBatch returns ready WIP batch
func (d *dbManager) GetWIPBatch(ctx context.Context) (*WipBatch, error) {
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

	lastBatchTxs, _, err := state.DecodeTxs(lastBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}

	lastBatch.Transactions = lastBatchTxs

	lastStateRoot, err := d.state.GetLastStateRoot(ctx, dbTx)
	if err != nil {
		return nil, err
	}

	wipBatch := &WipBatch{
		batchNumber:    lastBatch.BatchNumber,
		coinbase:       lastBatch.Coinbase,
		localExitRoot:  lastBatch.LocalExitRoot,
		timestamp:      uint64(lastBatch.Timestamp.Unix()),
		globalExitRoot: lastBatch.GlobalExitRoot,
		countOfTxs:     len(lastBatch.Transactions),
	}

	// Init counters to MAX values
	var totalBytes uint64 = d.batchConstraints.MaxBatchBytesSize
	var batchZkCounters state.ZKCounters = state.ZKCounters{
		CumulativeGasUsed:    d.batchConstraints.MaxCumulativeGasUsed,
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
			Timestamp:      time.Unix(int64(wipBatch.timestamp), 0),
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
				CumulativeGasUsed:    batchResponse.GetCumulativeGasUsed(),
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

	wipBatch.remainingResources = batchResources{zKCounters: batchZkCounters, bytes: totalBytes}
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

// CloseBatch closes a batch in the state
func (d *dbManager) CloseBatch(ctx context.Context, params ClosingBatchParameters) error {
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   params.BatchNumber,
		StateRoot:     params.StateRoot,
		LocalExitRoot: params.LocalExitRoot,
		AccInputHash:  params.AccInputHash,
	}

	batchL2Data, err := state.EncodeTransactions(params.Txs)
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
func (d *dbManager) ProcessForcedBatch(forcedBatchNum uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error) {
	// Open Batch
	processingCtx := state.ProcessingContext{
		BatchNumber:    request.BatchNumber,
		Coinbase:       request.Coinbase,
		Timestamp:      time.Unix(int64(request.Timestamp), 0),
		GlobalExitRoot: request.GlobalExitRoot,
		ForcedBatchNum: &forcedBatchNum,
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
	forcedBatch, err := d.state.GetForcedBatch(d.ctx, forcedBatchNum, dbTx)
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
	}

	// Close Batch
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   request.BatchNumber,
		StateRoot:     processBatchResponse.NewStateRoot,
		LocalExitRoot: processBatchResponse.NewLocalExitRoot,
		AccInputHash:  processBatchResponse.NewAccInputHash,
		BatchL2Data:   forcedBatch.RawTxsData,
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
func (d *dbManager) GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*types.Header, error) {
	return d.state.GetLastL2BlockHeader(ctx, dbTx)
}

func (d *dbManager) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error) {
	return d.state.GetLastBlock(ctx, dbTx)
}

func (d *dbManager) GetLastTrustedForcedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	return d.state.GetLastTrustedForcedBatchNumber(ctx, dbTx)
}

func (d *dbManager) GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	return d.state.GetBalanceByStateRoot(ctx, address, root)
}

func (d *dbManager) GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64) (txs []types.Transaction, err error) {
	return d.state.GetTransactionsByBatchNumber(ctx, batchNumber, nil)
}

func (d *dbManager) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool) error {
	return d.txPool.UpdateTxStatus(ctx, hash, newStatus, isWIP)
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

// AddDebugInfo is used to store debug info useful during runtime
func (d *dbManager) AddDebugInfo(ctx context.Context, info *state.DebugInfo, dbTx pgx.Tx) error {
	return d.state.AddDebugInfo(ctx, info, dbTx)
}

// AddEvent is used to store and event in the database
func (d *dbManager) AddEvent(ctx context.Context, event *state.Event, dbTx pgx.Tx) error {
	return d.state.AddEvent(ctx, event, dbTx)
}
