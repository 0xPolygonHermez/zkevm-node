package common

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type stateInterfaceReorgPool interface {
	GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*ethTypes.Transaction, error)
}

type ethermanInterfaceReorgPool interface {
	GetLatestBatchNumber() (uint64, error)
}

type poolInterfaceReorgPool interface {
	DeleteReorgedTransactions(ctx context.Context, txs []*ethTypes.Transaction) error
	StoreTx(ctx context.Context, tx ethTypes.Transaction, ip string, isWIP bool) error
}

type ReorgPool struct {
	state    stateInterfaceReorgPool
	etherMan ethermanInterfaceReorgPool
	pool     poolInterfaceReorgPool
}

func NewReorgPool(state stateInterfaceReorgPool, etherMan ethermanInterfaceReorgPool, pool poolInterfaceReorgPool) *ReorgPool {
	return &ReorgPool{
		state:    state,
		etherMan: etherMan,
		pool:     pool,
	}
}

func (p *ReorgPool) ReorgPool(ctx context.Context, dbTx pgx.Tx) error {
	latestBatchNum, err := p.etherMan.GetLatestBatchNumber()
	if err != nil {
		log.Error("error getting the latestBatchNumber virtualized in the smc. Error: ", err)
		return err
	}
	batchNumber := latestBatchNum + 1
	// Get transactions that have to be included in the pool again
	txs, err := p.state.GetReorgedTransactions(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting txs from trusted state. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Reorged transactions: ", txs)

	// Remove txs from the pool
	err = p.pool.DeleteReorgedTransactions(ctx, txs)
	if err != nil {
		log.Errorf("error deleting txs from the pool. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Delete reorged transactions")

	// Add txs to the pool
	for _, tx := range txs {
		// Insert tx in WIP status to avoid the sequencer to grab them before it gets restarted
		// When the sequencer restarts, it will update the status to pending non-wip
		err = p.pool.StoreTx(ctx, *tx, "", true)
		if err != nil {
			log.Errorf("error storing tx into the pool again. TxHash: %s. BatchNumber: %d, error: %v", tx.Hash().String(), batchNumber, err)
			return err
		}
		log.Debug("Reorged transactions inserted in the pool: ", tx.Hash())
	}
	return nil
}
