package sequencer

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	cmap "github.com/orcaman/concurrent-map"
)

// PendingTxsQueueConfig config for pending tx queue data structure
type PendingTxsQueueConfig struct {
	TxPendingInQueueCheckingFrequency Duration `mapstructure:"TxPendingInQueueCheckingFrequency"`
	GetPendingTxsFrequency            Duration `mapstructure:"GetPendingTxsFrequency"`
}

// PendingTxsQueue keeps pending tx queue and gives tx with the highest gas price by request
type PendingTxsQueue struct {
	cfg                 PendingTxsQueueConfig
	poppedTxsHashesChan chan common.Hash
	poppedTxsHashesMap  cmap.ConcurrentMap

	pendingTxs      []pool.Transaction
	pendingTxsMutex sync.RWMutex
	pendingTxsMap   cmap.ConcurrentMap
	txPool          txPool
}

// NewPendingTxsQueue inits new pending tx queue
func NewPendingTxsQueue(cfg PendingTxsQueueConfig, pool txPool) *PendingTxsQueue {
	poppedTxsChan := make(chan common.Hash, amountOfPendingTxsRequested)
	poppedTxsHashesMap := cmap.New()
	pendingTxMap := cmap.New()
	return &PendingTxsQueue{
		cfg:                 cfg,
		txPool:              pool,
		pendingTxsMap:       pendingTxMap,
		poppedTxsHashesChan: poppedTxsChan,
		poppedTxsHashesMap:  poppedTxsHashesMap,
	}
}

// PopPendingTx pops top pending tx from the queue
func (q *PendingTxsQueue) PopPendingTx() *pool.Transaction {
	var tx *pool.Transaction
	q.pendingTxsMutex.Lock()
	defer q.pendingTxsMutex.Unlock()
	if len(q.pendingTxs) > 1 {
		tx, q.pendingTxs = &q.pendingTxs[0], q.pendingTxs[1:]
	} else if len(q.pendingTxs) == 1 {
		tx = &q.pendingTxs[0]
		q.pendingTxs = []pool.Transaction{}
	} else {
		return nil
	}
	txHash := tx.Hash().Hex()
	q.poppedTxsHashesMap.Set(txHash, true)
	q.pendingTxsMap.Remove(txHash)
	q.poppedTxsHashesChan <- tx.Hash()

	return tx
}

// findPlaceInSlice finds place where to insert tx to the queue according to gas amount
func (q *PendingTxsQueue) findPlaceInSlice(pendingTx pool.Transaction) int {
	q.pendingTxsMutex.RLock()
	defer q.pendingTxsMutex.RUnlock()
	low := 0
	high := len(q.pendingTxs) - 1
	for low <= high {
		median := low + (high-low)/2 // nolint:gomnd
		if q.pendingTxs[median].Gas() == pendingTx.Gas() {
			return median
		} else if q.pendingTxs[median].Gas() < pendingTx.Gas() {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	return high + 1
}

// InsertPendingTx insert pending tx from the pool to the queue
func (q *PendingTxsQueue) InsertPendingTx(tx pool.Transaction) {
	index := q.findPlaceInSlice(tx)
	q.pendingTxsMutex.Lock()
	defer q.pendingTxsMutex.Unlock()
	if index <= 1 {
		q.pendingTxs = append(q.pendingTxs, tx)
	} else {
		q.pendingTxs = append(q.pendingTxs[:index+1], q.pendingTxs[index:]...)
		q.pendingTxs[index] = tx
	}
}

// CleanPendTxsChan cleans pending tx that is already popped from the queue and selected/rejected
func (q *PendingTxsQueue) CleanPendTxsChan(ctx context.Context) {
	for {
		select {
		case hash := <-q.poppedTxsHashesChan:
			q.waitForTxToBeProcessed(ctx, hash)
		case <-ctx.Done():
			return
		}
	}
}

// waitForTxToBeProcessed for the tx to change it's status from pending to invalid or selected
func (q *PendingTxsQueue) waitForTxToBeProcessed(ctx context.Context, hash common.Hash) {
	var err error
	tickerIsTxPending := time.NewTicker(q.cfg.TxPendingInQueueCheckingFrequency.Duration)
	isPending := true
	for isPending {
		isPending, err = q.txPool.IsTxPending(ctx, hash)
		if err != nil {
			log.Warnf("failed to check if tx is still pending, txHash: %s, err: %v", hash.Hex(), err)
		}

		if !isPending {
			q.poppedTxsHashesMap.Remove(hash.Hex())
			return
		}
		select {
		case <-tickerIsTxPending.C:
			// nothing
		case <-ctx.Done():
			return
		}
	}
}

// KeepPendingTxsQueue keeps pending txs queue full
func (q *PendingTxsQueue) KeepPendingTxsQueue(ctx context.Context) {
	var err error
	q.pendingTxsMutex.Lock()
	for len(q.pendingTxs) == 0 {
		q.pendingTxs, err = q.txPool.GetPendingTxs(ctx, false, amountOfPendingTxsRequested)
		if err != nil {
			log.Errorf("failed to get pending txs, err: %v", err)
		}
		if len(q.pendingTxs) == 0 {
			time.Sleep(q.cfg.GetPendingTxsFrequency.Duration)
		}
	}
	q.pendingTxsMutex.Unlock()

	for _, tx := range q.pendingTxs {
		q.pendingTxsMap.Set(tx.Hash().Hex(), true)
	}

	for {
		time.Sleep(q.cfg.GetPendingTxsFrequency.Duration)
		lenPendingTxs := q.GetPendingTxsQueueLength()
		if lenPendingTxs >= amountOfPendingTxsRequested {
			continue
		}
		pendTx, err := q.txPool.GetPendingTxs(ctx, false, 1)
		if err != nil {
			log.Errorf("failed to get pending tx, err: %v", err)
			continue
		}
		if len(pendTx) != 0 && (lenPendingTxs == 0 ||
			(!q.poppedTxsHashesMap.Has(pendTx[0].Hash().Hex()) && !q.pendingTxsMap.Has(pendTx[0].Hash().Hex()))) {
			q.InsertPendingTx(pendTx[0])
			q.pendingTxsMap.Set(pendTx[0].Hash().Hex(), true)
		}
	}
}

// GetPendingTxsQueueLength get length
func (q *PendingTxsQueue) GetPendingTxsQueueLength() int {
	q.pendingTxsMutex.RLock()
	defer q.pendingTxsMutex.RUnlock()
	return len(q.pendingTxs)
}
