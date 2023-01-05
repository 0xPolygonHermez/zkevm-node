package sequencer

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Pool Loader and DB Updater
type dbManager struct {
	txPool       txPool
	state        dbManagerStateInterface
	worker       workerInterface
	txsToStoreCh chan *txToStore
	wgTxsToStore *sync.WaitGroup
	ctx          context.Context
}

func newDBManager(ctx context.Context, txPool txPool, state dbManagerStateInterface, worker *Worker, txsToStoreCh chan *txToStore, wgTxsToStore *sync.WaitGroup) *dbManager {
	return &dbManager{ctx: ctx, txPool: txPool, state: state, worker: worker, txsToStoreCh: txsToStoreCh, wgTxsToStore: wgTxsToStore}
}

func (d *dbManager) Start() {
	go d.loadFromPool()
	go d.StoreProcessedTxAndDeleteFromPool()
}

func (d *dbManager) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	return d.state.GetLastBatchNumber(ctx, nil)
}

func (d *dbManager) CreateFirstBatch(ctx context.Context, sequencerAddress common.Address) (*state.ProcessingContext, error) {
	processingCtx := &state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: state.ZeroHash,
	}
	dbTx, err := d.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for opening a batch, err: %v", err)
		return nil, err
	}
	err = d.state.OpenBatch(ctx, *processingCtx, dbTx)
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
	return processingCtx, nil
}

func (d *dbManager) loadFromPool() {

	ctx := context.Background()

	for {
		// TODO: Define how to do this
		time.Sleep(5 * time.Second)

		// TODO: Decide about the creation of a new GetPending function
		poolTransactions, err := d.txPool.GetPendingTxs(ctx, false, 0)

		if err != nil && err != pgpoolstorage.ErrNotFound {
			log.Errorf("loadFromPool: %v", err)
			continue
		}

		poolClaims, err := d.txPool.GetPendingTxs(ctx, true, 0)

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

			txTracker := TxTracker{
				Hash: tx.Hash(),
				// TODO: Complete
			}
			d.worker.AddTx(txTracker)
			// TODO: Redefine pool tx statuses
			d.txPool.UpdateTxStatus(ctx, tx.Hash(), pool.TxStatusSelected)
		}

	}
}

func (d *dbManager) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	return d.BeginStateTransaction(ctx)
}

func (d *dbManager) StoreProcessedTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, dbTx pgx.Tx) error {
	return d.state.StoreTransaction(ctx, batchNumber, processedTx, dbTx)
}

func (d *dbManager) DeleteTransactionFromPool(ctx context.Context, txHash common.Hash) error {
	return d.txPool.DeleteTransactionByHash(ctx, txHash)
}

func (d *dbManager) StoreProcessedTxAndDeleteFromPool() {
	// TODO: Finish the retry mechanism and error handling
	for {
		txToStore := <-d.txsToStoreCh

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
		// TODO: GetLatestL2Block from database and compare with txToStore.previousL2BlockStateRoot
		// Send signal to L2ReorgCh            chan struct{}

		// TODO: Change this to update status to selected
		err = d.DeleteTransactionFromPool(d.ctx, txToStore.txResponse.TxHash)
		if err != nil {
			err = dbTx.Rollback(d.ctx)
			if err != nil {
				log.Errorf("StoreProcessedTxAndDeleteFromPool :%v", err)
			}
		}

		err = dbTx.Commit(d.ctx)
		if err != nil {
			log.Errorf("StoreProcessedTxAndDeleteFromPool error committing: %v", err)
		}

		d.wgTxsToStore.Done()
	}
}

// GetWIPBatch returns ready WIP batch
// if lastBatch IS OPEN - load data from it but set wipBatch.initialStateRoot to Last Closed Batch
// if lastBatch IS CLOSED - open new batch in the database and load all data from the closed one without the txs and increase batch number
func (d *dbManager) GetWIPBatch(ctx context.Context) (*wipBatch, error) {
	lastBatch, err := d.GetLastBatch(ctx)
	if err != nil {
		return nil, err
	}

	wipBatch := &wipBatch{
		batchNumber:      lastBatch.BatchNumber,
		coinbase:         lastBatch.Coinbase,
		accInputHash:     lastBatch.AccInputHash,
		initialStateRoot: lastBatch.StateRoot,
		stateRoot:        lastBatch.StateRoot,
		timestamp:        uint64(lastBatch.Timestamp.Unix()),
		globalExitRoot:   lastBatch.GlobalExitRoot,

		// TODO: txs
		// TODO: remainingResources
	}

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

		wipBatch.initialStateRoot = lastClosedBatch.StateRoot
	}

	return wipBatch, nil
}

func (d *dbManager) GetLastClosedBatch(ctx context.Context) (*state.Batch, error) {
	return d.state.GetLastClosedBatch(ctx, nil)
}

func (d *dbManager) GetLastBatch(ctx context.Context) (*state.Batch, error) {
	// TODO: Implement new method in state
	batches, err := d.state.GetLastNBatches(ctx, 1, nil)
	if err != nil {
		return nil, err
	}
	return batches[0], nil
}

func (d *dbManager) IsBatchClosed(ctx context.Context, batchNum uint64) (bool, error) {
	return d.state.IsBatchClosed(ctx, batchNum, nil)
}

func (d *dbManager) GetLastNBatches(ctx context.Context, numBatches uint) ([]*state.Batch, error) {
	return d.state.GetLastNBatches(ctx, numBatches, nil)
}

// ClosingBatchParameters contains the necessary parameters to close a batch
type ClosingBatchParameters struct {
	BatchNumber   uint64
	StateRoot     common.Hash
	LocalExitRoot common.Hash
	AccInputHash  common.Hash
	Txs           []TxTracker
}

func (d *dbManager) CloseBatch(ctx context.Context, params ClosingBatchParameters) error {

	// Create new type txManagerArray and refactor CloseBatch method in state

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

func (d *dbManager) MarkReorgedTxsAsPending(ctx context.Context) {
	// TODO: Handle error
	err := d.txPool.MarkReorgedTxsAsPending(ctx)
}
