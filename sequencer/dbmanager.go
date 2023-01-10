package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	waitSeconds time.Duration = 5
)

// Pool Loader and DB Updater
type dbManager struct {
	txPool    txPool
	state     dbManagerStateInterface
	worker    workerInterface
	txsStore  TxsStore
	l2ReorgCh chan L2ReorgEvent
	ctx       context.Context
}

// ClosingBatchParameters contains the necessary parameters to close a batch
type ClosingBatchParameters struct {
	BatchNumber   uint64
	StateRoot     common.Hash
	LocalExitRoot common.Hash
	AccInputHash  common.Hash
	Txs           []TxTracker
}

func newDBManager(ctx context.Context, txPool txPool, state dbManagerStateInterface, worker *Worker, closingSignalCh ClosingSignalCh, txsStore TxsStore) *dbManager {
	return &dbManager{ctx: ctx, txPool: txPool, state: state, worker: worker, txsStore: txsStore, l2ReorgCh: closingSignalCh.L2ReorgCh}
}

// Start stars the dbManager routines
func (d *dbManager) Start() {
	go d.loadFromPool()
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

// loadFromPool keeps loading transaction from the pool
func (d *dbManager) loadFromPool() {
	for {
		// TODO: Move this to a config parameter
		time.Sleep(waitSeconds * time.Second)

		poolTransactions, err := d.txPool.GetPendingTxs(d.ctx, false, 0)

		if err != nil && err != pgpoolstorage.ErrNotFound {
			log.Errorf("loadFromPool: %v", err)
			continue
		}

		poolClaims, err := d.txPool.GetPendingTxs(d.ctx, true, 0)

		if err != nil && err != pgpoolstorage.ErrNotFound {
			log.Errorf("loadFromPool: %v", err)
			continue
		}

		poolTransactions = append(poolTransactions, poolClaims...)

		for _, tx := range poolTransactions {
			if err != nil {
				log.Errorf("loadFromPool error getting tx sender: %v", err)
				continue
			}

			txTracker := d.worker.NewTxTracker(tx.Transaction, tx.ZKCounters)
			d.worker.AddTx(*txTracker)
			d.txPool.UpdateTxStatus(d.ctx, tx.Hash(), pool.TxStatusWIP)
		}
	}
}

// BeginStateTransaction starts a db transaction in the state
func (d *dbManager) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	return d.state.BeginStateTransaction(ctx)
}

// StoreProcessedTransaction stores a transaction in the state
func (d *dbManager) StoreProcessedTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, dbTx pgx.Tx) error {
	return d.state.StoreTransaction(ctx, batchNumber, processedTx, dbTx)
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

		dbTx, err := d.BeginStateTransaction(d.ctx)
		if err != nil {
			log.Errorf("StoreProcessedTxAndDeleteFromPool :%v", err)
		}
		err = d.StoreProcessedTransaction(d.ctx, txToStore.batchNumber, txToStore.txResponse, dbTx)
		if err != nil {
			err = dbTx.Rollback(d.ctx)
			if err != nil {
				log.Errorf("StoreProcessedTxAndDeleteFromPool :%v", err)
			}
		}

		// Check if the Tx is still valid in the state to detect reorgs
		latestL2Block, err := d.state.GetLastL2Block(d.ctx, dbTx)
		if err != nil {
			log.Errorf("StoreProcessedTxAndDeleteFromPool :%v", err)
		}
		if latestL2Block.Root() != txToStore.previousL2BlockStateRoot {
			log.Info("L2 reorg detected. Old state root: %v New state root: %v", latestL2Block.Root(), txToStore.previousL2BlockStateRoot)
			d.l2ReorgCh <- L2ReorgEvent{}
			continue
		}

		// Change Tx status to selected
		d.txPool.UpdateTxStatus(d.ctx, txToStore.txResponse.TxHash, pool.TxStatusSelected)

		err = dbTx.Commit(d.ctx)
		if err != nil {
			log.Errorf("StoreProcessedTxAndDeleteFromPool error committing: %v", err)
		}

		d.txsStore.Wg.Done()
	}
}

// GetWIPBatch returns ready WIP batch
// if lastBatch IS OPEN - load data from it but set wipBatch.initialStateRoot to Last Closed Batch
// if lastBatch IS CLOSED - open new batch in the database and load all data from the closed one without the txs and increase batch number
func (d *dbManager) GetWIPBatch(ctx context.Context) (*WipBatch, error) {
	lastBatch, err := d.GetLastBatch(ctx)
	if err != nil {
		return nil, err
	}

	wipBatch := &WipBatch{
		batchNumber:    lastBatch.BatchNumber,
		coinbase:       lastBatch.Coinbase,
		accInputHash:   lastBatch.AccInputHash,
		stateRoot:      lastBatch.StateRoot,
		localExitRoot:  lastBatch.LocalExitRoot,
		timestamp:      uint64(lastBatch.Timestamp.Unix()),
		globalExitRoot: lastBatch.GlobalExitRoot,
	}

	txs, _, err := state.DecodeTxs(lastBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}

	txTrackerArray := make([]TxTracker, len(txs), len(txs))

	// TODO: Init counters and totals to MAX values
	var totalBytes uint64
	var totalGas uint64
	var batchZkCounters state.ZKCounters

	for _, tx := range txs {
		zkCounters, err := d.txPool.GetTxZkCountersByHash(ctx, tx.Hash())
		if err != nil {
			return nil, err
		}
		txTracker := d.worker.NewTxTracker(tx, *zkCounters)

		err = batchZkCounters.Sub(*zkCounters)
		if err != nil {
			return nil, err
		}

		if uint64(tx.Size()) <= totalBytes {
			totalBytes -= uint64(tx.Size())
		} else {
			return nil, fmt.Errorf("tx does not fit into the batch because of byte size")
		}

		if zkCounters.CumulativeGasUsed <= totalGas {
			totalGas -= zkCounters.CumulativeGasUsed
		} else {
			return nil, fmt.Errorf("tx does not fit into the batch because of gas")
		}

		txTrackerArray = append(txTrackerArray, *txTracker)
	}

	wipBatch.txs = txTrackerArray
	wipBatch.remainingResources = BatchResources{zKCounters: batchZkCounters, bytes: totalBytes, gas: totalGas}

	isClosed, err := d.IsBatchClosed(ctx, lastBatch.BatchNumber)
	if err != nil {
		return nil, err
	}

	if isClosed {
		wipBatch.batchNumber = lastBatch.BatchNumber + 1

		processingContext := &state.ProcessingContext{
			BatchNumber:    wipBatch.batchNumber,
			Coinbase:       wipBatch.coinbase,
			Timestamp:      time.Now(),
			GlobalExitRoot: wipBatch.globalExitRoot,
		}

		dbTx, err := d.BeginStateTransaction(ctx)
		if err != nil {
			return nil, err
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
		lastClosedBatch, err := d.GetLastClosedBatch(ctx)
		if err != nil {
			return nil, err
		}

		wipBatch.stateRoot = lastClosedBatch.StateRoot
	}

	return wipBatch, nil
}

// GetLastClosedBatch gets the latest closed batch from state
func (d *dbManager) GetLastClosedBatch(ctx context.Context) (*state.Batch, error) {
	return d.state.GetLastClosedBatch(ctx, nil)
}

// GetLastBatch gets the latest batch from state
func (d *dbManager) GetLastBatch(ctx context.Context) (*state.Batch, error) {
	batch, err := d.state.GetLastBatch(ctx)
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
func (d *dbManager) GetLatestGer(ctx context.Context, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error) {
	ger, receivedAt, err := d.state.GetLatestGlobalExitRoot(ctx, maxBlockNumber, nil)
	if err != nil && errors.Is(err, state.ErrNotFound) {
		return state.GlobalExitRoot{}, time.Time{}, nil
	} else if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, fmt.Errorf("failed to get latest global exit root, err: %w", err)
	} else {
		return ger, receivedAt, nil
	}
}

// CloseBatch closes a batch in the state
func (d *dbManager) CloseBatch(ctx context.Context, params ClosingBatchParameters) error {
	// TODO: Create new type txManagerArray and refactor CloseBatch method in state

	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   params.BatchNumber,
		StateRoot:     params.StateRoot,
		LocalExitRoot: params.LocalExitRoot,
		AccInputHash:  params.AccInputHash,
	}

	transactions := make([]types.Transaction, len(params.Txs))

	for _, tx := range params.Txs {
		transaction, err := state.DecodeTx(string(tx.RawTx))

		if err != nil {
			return err
		}

		transactions = append(transactions, *transaction)
	}

	processingReceipt.Txs = transactions

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

// MarkReorgedTxsAsPending marks all reorged tx as pending in the pool
func (d *dbManager) MarkReorgedTxsAsPending(ctx context.Context) {
	err := d.txPool.MarkReorgedTxsAsPending(ctx)
	if err != nil {
		log.Errorf("error marking reorged txs as pending: %v", err)
	}
}
