package sequencer

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
)

// PendingTxsQueueConfig config for pending tx queue data structure
type PendingTxsQueueConfig struct {
	TxPendingInQueueCheckingFrequency Duration `mapstructure:"TxPendingInQueueCheckingFrequency"`
	TxPoppedCheckingFrequency         Duration `mapstructure:"TxPoppedCheckingFrequency"`
	GetPendingTxsFrequency            Duration `mapstructure:"GetPendingTxsFrequency"`
}

// PendingTxsQueue keeps pending tx queue and gives tx with the highest gas price by request
type PendingTxsQueue struct {
	cfg                 PendingTxsQueueConfig
	poppedTxsHashesChan chan common.Hash
	poppedTxsHashesMap  map[common.Hash]bool
	pendingTxs          []pool.Transaction
	pendingTxsMap       map[common.Hash]bool
	pendingTxsMutex     *sync.RWMutex
	txPool              txPool
}

// NewPendingTxsQueue inits new pending tx queue
func NewPendingTxsQueue(cfg PendingTxsQueueConfig, pool txPool) *PendingTxsQueue {
	poppedTxsChan := make(chan common.Hash, amountOfPendingTxsRequested)
	poppedTxsHashesMap := make(map[common.Hash]bool)
	pendingTxMap := make(map[common.Hash]bool)
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
	txHash := tx.Hash()
	q.poppedTxsHashesMap[txHash] = true
	delete(q.pendingTxsMap, txHash)
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
	tickerPoppedTxs := time.NewTicker(q.cfg.TxPoppedCheckingFrequency.Duration)
	defer tickerPoppedTxs.Stop()
	tickerIsTxPending := time.NewTicker(q.cfg.TxPendingInQueueCheckingFrequency.Duration)
	defer tickerIsTxPending.Stop()
	var err error
	for {
		hash := <-q.poppedTxsHashesChan
		isPending := true
		for isPending {
			isPending, err = q.txPool.IsTxPending(ctx, hash)
			if err != nil {
				log.Warnf("failed to check if tx is still pending, txHash: %s, err: %v", hash.Hex(), err)
			}
			if !isPending {
				delete(q.poppedTxsHashesMap, hash)
			}
			select {
			case <-tickerIsTxPending.C:
				// nothing
			case <-ctx.Done():
				return
			}
		}

		select {
		case <-tickerPoppedTxs.C:
			// nothing
		case <-ctx.Done():
			return
		}
	}
}

// KeepPendingTxsQueue keeps pending txs queue full
func (q *PendingTxsQueue) KeepPendingTxsQueue(ctx context.Context) {
	q.pendingTxsMutex.Lock()
	for len(q.pendingTxs) == 0 {
		txs, err := q.txPool.GetPendingTxs(ctx, false, amountOfPendingTxsRequested)
		q.pendingTxs = txs
		if err != nil {
			log.Errorf("failed to get pending txs, err: %v", err)
		}
		time.Sleep(q.cfg.GetPendingTxsFrequency.Duration)
	}
	for _, tx := range q.pendingTxs {
		q.pendingTxsMap[tx.Hash()] = true
	}
	q.pendingTxsMutex.Unlock()

	for {
		time.Sleep(q.cfg.GetPendingTxsFrequency.Duration)
		if len(q.pendingTxs) < amountOfPendingTxsRequested {
			pendTx, err := q.txPool.GetPendingTxs(ctx, false, 1)
			if err != nil {
				log.Errorf("failed to get pending tx, err: %v", err)
				continue
			}
			if len(pendTx) != 0 && (len(q.pendingTxs) == 0 ||
				(!q.poppedTxsHashesMap[pendTx[0].Hash()] && !q.pendingTxsMap[pendTx[0].Hash()])) {
				q.InsertPendingTx(pendTx[0])
				q.pendingTxsMap[pendTx[0].Hash()] = true
			}
		}
	}
}
