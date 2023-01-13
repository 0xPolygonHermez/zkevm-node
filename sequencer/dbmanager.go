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
	l2ReorgCh    chan L2ReorgEvent
	ctx          context.Context
}

func (d *dbManager) ProcessForcedBatch(forcedBatchNum uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error) {
	//TODO implement me
	panic("implement me")
}

func newDBManager(ctx context.Context, txPool txPool, state dbManagerStateInterface, worker *Worker, closingSignalCh ClosingSignalCh, txsStore TxsStore) *dbManager {
	return &dbManager{ctx: ctx, txPool: txPool, state: state, worker: worker, txsToStoreCh: txsStore.Ch, wgTxsToStore: txsStore.Wg, l2ReorgCh: closingSignalCh.L2ReorgCh}
}

func (d *dbManager) Start() {
	go d.loadFromPool()
	go d.StoreProcessedTxAndDeleteFromPool()
}

func (d *dbManager) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	return d.state.GetLastBatchNumber(ctx, nil)
}

func (d *dbManager) OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error {
	//TODO: Use state interface to OpenBatch in the DB
	panic("implement me")
}

func (d *dbManager) CreateFirstBatch(ctx context.Context, sequencerAddress common.Address) state.ProcessingContext {
	// TODO: Retry in case of error
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

func (d *dbManager) loadFromPool() {

	ctx := context.Background()

	for {
		// TODO: Define how to do this
		time.Sleep(5 * time.Second)

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

			txTracker, err := d.worker.NewTx(tx.Transaction, tx.ZKCounters)
			if err != nil {
				// TODO: How to handle this error?
			}
			d.worker.AddTx(ctx, txTracker)
			d.txPool.UpdateTxStatus(ctx, tx.Hash(), pool.TxStatusWIP)
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
		latestL2Block, err := d.state.GetLastL2Block(d.ctx, dbTx)
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

		d.wgTxsToStore.Done()
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
		batchNumber:  lastBatch.BatchNumber,
		coinbase:     lastBatch.Coinbase,
		accInputHash: lastBatch.AccInputHash,
		// initialStateRoot: lastBatch.StateRoot,
		stateRoot:      lastBatch.StateRoot,
		timestamp:      uint64(lastBatch.Timestamp.Unix()),
		globalExitRoot: lastBatch.GlobalExitRoot,

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

		wipBatch.stateRoot = lastClosedBatch.StateRoot
	}

	return wipBatch, nil
}

func (d *dbManager) GetLastClosedBatch(ctx context.Context) (*state.Batch, error) {
	return d.state.GetLastClosedBatch(ctx, nil)
}

func (d *dbManager) GetLastBatch(ctx context.Context) (*state.Batch, error) {
	batch, err := d.state.GetLastBatch(ctx)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (d *dbManager) IsBatchClosed(ctx context.Context, batchNum uint64) (bool, error) {
	return d.state.IsBatchClosed(ctx, batchNum, nil)
}

func (d *dbManager) GetLastNBatches(ctx context.Context, numBatches uint) ([]*state.Batch, error) {
	return d.state.GetLastNBatches(ctx, numBatches, nil)
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

func (d *dbManager) MarkReorgedTxsAsPending(ctx context.Context) {
	err := d.txPool.MarkReorgedTxsAsPending(ctx)
	if err != nil {
		log.Errorf("error marking reorged txs as pending: %v", err)
	}
}
