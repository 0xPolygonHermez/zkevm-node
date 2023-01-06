package sequencer

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Pool Loader and DB Updater
type dbManager struct {
	txPool txPool
	state  stateInterface
	worker workerInterface
}

func newDBManager(txPool txPool, state stateInterface, worker *Worker) *dbManager {
	return &dbManager{txPool: txPool, state: state, worker: worker}
}

func (d *dbManager) Start() {
	go d.loadFromPool()
}

func (d *dbManager) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	// TODO: Fetch last BatchNumber from database
	return 0, errors.New("")
}

func (d *dbManager) OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error {
	//TODO: Use state interface to OpenBatch in the DB
	panic("implement me")
}

func (d *dbManager) CreateFirstBatch(ctx context.Context, sequencerAddress common.Address) state.ProcessingContext {
	processingCtx := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: state.ZeroHash,
	}
	dbTx, err := d.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Fatalf("failed to begin state transaction for opening a batch, err: %v", err)
	}
	err = d.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Fatalf(
				"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Fatalf("failed to open a batch, err: %v", err)
	}
	if err := dbTx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit dbTx when opening batch, err: %v", err)
	}
	return processingCtx
}

func (d *dbManager) loadFromPool() {
	// TODO: Endless loop that keeps loading tx from the DB into the worker
}

func (d *dbManager) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := d.state.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (d *dbManager) StoreProcessedTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, dbTx pgx.Tx) error {
	// TODO: Implement store of transaction and adding it to the batch
	return errors.New("")
}

func (d *dbManager) DeleteTxFromPool(ctx context.Context, txHash common.Hash, dbTx pgx.Tx) error {
	// TODO: Delete transaction from Pool DB
	return errors.New("")
}

func (d *dbManager) StoreProcessedTxAndDeleteFromPool(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse) {
	for { // TODO: Finish the retry mechanism
		dbTx, err := d.BeginStateTransaction(ctx)
		if err != nil {
			// TODO: handle
		}
		err = d.StoreProcessedTransaction(ctx, batchNumber, processedTx, dbTx)
		if err != nil {
			err = dbTx.Rollback(ctx)
			if err != nil {
				// TODO: handle
			}
		}
		err = d.DeleteTxFromPool(ctx, processedTx.TxHash, dbTx)
		if err != nil {
			err = dbTx.Rollback(ctx)
			if err != nil {
				// TODO: handle
			}
		}
	}
}

func (d *dbManager) GetWIPBatch(ctx context.Context) (WipBatch, error) {
	// TODO: Make this method to return ready WIP batch it has following cases:
	// if lastBatch IS OPEN - load data from it but set WipBatch.initialStateRoot to Last Closed Batch
	// if lastBatch IS CLOSED - open new batch in the database and load all data from the closed one without the txs and increase batch number
	return WipBatch{}, errors.New("")
}

func (d *dbManager) GetLastClosedBatch(ctx context.Context) (state.Batch, error) {
	// TODO: Returns last closed batch
	return state.Batch{}, errors.New("")
}

func (d *dbManager) GetLastBatch(ctx context.Context) (state.Batch, error) {
	// TODO: Returns last batch
	return state.Batch{}, errors.New("")

}

func (d *dbManager) IsBatchClosed(ctx context.Context, batchNum uint64) (bool, error) {
	// TODO: Returns if the batch with passed batchNum is closed
	return false, errors.New("")
}

func (d *dbManager) GetLastNBatches(ctx context.Context, numBatches uint) ([]*state.Batch, error) {
	// TODO: Returns last N batches
	return []*state.Batch{}, errors.New("")

}
func (d *dbManager) GetLatestGer(ctx context.Context) (state.GlobalExitRoot, time.Time, error) {
	// TODO: Get implementation from old sequencer's batchbuilder
	return state.GlobalExitRoot{}, time.Now(), nil
}

// ClosingBatchParameters contains the necessary parameters to close a batch
type ClosingBatchParameters struct {
	BatchNumber   uint64
	StateRoot     common.Hash
	LocalExitRoot common.Hash
	AccInputHash  common.Hash
	Txs           []TxTracker
}

func (d *dbManager) CloseBatch(ctx context.Context, params ClosingBatchParameters, dbTx pgx.Tx) {
	// TODO: Close current open batch
}

func (d *dbManager) MarkReorgedTxsAsPending(ctx context.Context) error {
	// TODO: call pool.MarkReorgedTxsAsPending and return result
	return errors.New("")
}
