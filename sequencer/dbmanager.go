package sequencer

import (
	"context"
	"errors"

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

func (d *dbManager) StoreProcessedTransaction(ctx context.Context, dbTx pgx.Tx, batchNumber uint64, processedTx *state.ProcessTransactionResponse) error {
	// TODO: Implement store of transaction and adding it to the batch
	return errors.New("")
}

func (d *dbManager) DeleteTxFromPool(ctx context.Context, dbTx pgx.Tx, txHash common.Hash) error {
	// TODO: Delete transaction from Pool DB
	return errors.New("")
}

func (d *dbManager) StoreProcessedTxAndDeleteFromPool(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse) {
	for { // TODO: Finish the retry mechanism
		dbTx, err := d.BeginStateTransaction(ctx)
		if err != nil {
			// TODO: handle
		}
		err = d.StoreProcessedTransaction(ctx, dbTx, batchNumber, processedTx)
		if err != nil {
			err = dbTx.Rollback(ctx)
			if err != nil {
				// TODO: handle
			}
		}
		err = d.DeleteTxFromPool(ctx, dbTx, processedTx.TxHash)
		if err != nil {
			err = dbTx.Rollback(ctx)
			if err != nil {
				// TODO: handle
			}
		}
	}
}

func (d *dbManager) CloseBatch(ctx context.Context, receipt state.ProcessingReceipt) {
	// TODO: Close current open batch
}

func (d *dbManager) GetLastBatch(ctx context.Context) (*state.Batch, error) {
	// TODO: Get last batch
	return nil, errors.New("")
}
