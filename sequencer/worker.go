package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Worker represents the worker component of the sequencer
type Worker struct {
	pool         map[string]*addrQueue
	txSortedList *txSortedList
	workerMutex  sync.Mutex
	state        stateInterface
}

// NewWorker creates an init a worker
func NewWorker(state stateInterface) *Worker {
	w := Worker{
		pool:         make(map[string]*addrQueue),
		txSortedList: newTxSortedList(),
		state:        state,
	}

	return &w
}

// NewTxTracker creates and inits a TxTracker
func (w *Worker) NewTxTracker(tx types.Transaction, counters state.ZKCounters, ip string) (*TxTracker, error) {
	return newTxTracker(tx, counters, ip)
}

// AddTxTracker adds a new Tx to the Worker
func (w *Worker) AddTxTracker(ctx context.Context, tx *TxTracker) (replacedTx *TxTracker, dropReason error) {
	w.workerMutex.Lock()

	addr, found := w.pool[tx.FromStr]

	if !found {
		// Unlock the worker to let execute other worker functions while creating the new AddrQueue
		w.workerMutex.Unlock()

		root, err := w.state.GetLastStateRoot(ctx, nil)
		if err != nil {
			dropReason = fmt.Errorf("AddTx GetLastStateRoot error: %v", err)
			log.Error(dropReason)
			return nil, dropReason
		}
		nonce, err := w.state.GetNonceByStateRoot(ctx, tx.From, root)
		if err != nil {
			dropReason = fmt.Errorf("AddTx GetNonceByStateRoot error: %v", err)
			log.Error(dropReason)
			return nil, dropReason
		}
		balance, err := w.state.GetBalanceByStateRoot(ctx, tx.From, root)
		if err != nil {
			dropReason = fmt.Errorf("AddTx GetBalanceByStateRoot error: %v", err)
			log.Error(dropReason)
			return nil, dropReason
		}

		addr = newAddrQueue(tx.From, nonce.Uint64(), balance)

		// Lock again the worker
		w.workerMutex.Lock()

		w.pool[tx.FromStr] = addr
		log.Infof("AddTx new addrQueue created for addr(%s) nonce(%d) balance(%s)", tx.FromStr, nonce.Uint64(), balance.String())
	}

	// Add the txTracker to Addr and get the newReadyTx and prevReadyTx
	log.Infof("AddTx new tx(%s) nonce(%d) cost(%s) to addrQueue(%s) nonce(%d) balance(%d)", tx.HashStr, tx.Nonce, tx.Cost.String(), addr.fromStr, addr.currentNonce, addr.currentBalance)
	var newReadyTx, prevReadyTx, repTx *TxTracker
	newReadyTx, prevReadyTx, repTx, dropReason = addr.addTx(tx)
	if dropReason != nil {
		log.Infof("AddTx tx(%s) dropped from addrQueue(%s), reason: %s", tx.HashStr, tx.FromStr, dropReason.Error())
		w.workerMutex.Unlock()
		return repTx, dropReason
	}

	// Update the txSortedList (if needed)
	if prevReadyTx != nil {
		log.Infof("AddTx prevReadyTx(%s) nonce(%d) cost(%s) deleted from TxSortedList", prevReadyTx.HashStr, prevReadyTx.Nonce, prevReadyTx.Cost.String())
		w.txSortedList.delete(prevReadyTx)
	}
	if newReadyTx != nil {
		log.Infof("AddTx newReadyTx(%s) nonce(%d) cost(%s) added to TxSortedList", newReadyTx.HashStr, newReadyTx.Nonce, newReadyTx.Cost.String())
		w.txSortedList.add(newReadyTx)
	}

	if repTx != nil {
		log.Infof("AddTx replacedTx(%s) nonce(%d) cost(%s) has been replaced", repTx.HashStr, repTx.Nonce, repTx.Cost.String())
	}

	w.workerMutex.Unlock()
	return repTx, nil
}

func (w *Worker) applyAddressUpdate(from common.Address, fromNonce *uint64, fromBalance *big.Int) (*TxTracker, *TxTracker, []*TxTracker) {
	addrQueue, found := w.pool[from.String()]

	if found {
		newReadyTx, prevReadyTx, txsToDelete := addrQueue.updateCurrentNonceBalance(fromNonce, fromBalance)

		// Update the TxSortedList (if needed)
		if prevReadyTx != nil {
			log.Infof("applyAddressUpdate prevReadyTx(%s) nonce(%d) cost(%s) deleted from TxSortedList", prevReadyTx.Hash.String(), prevReadyTx.Nonce, prevReadyTx.Cost.String())
			w.txSortedList.delete(prevReadyTx)
		}
		if newReadyTx != nil {
			log.Infof("applyAddressUpdate newReadyTx(%s) nonce(%d) cost(%s) added to TxSortedList", newReadyTx.Hash.String(), newReadyTx.Nonce, newReadyTx.Cost.String())
			w.txSortedList.add(newReadyTx)
		}

		return newReadyTx, prevReadyTx, txsToDelete
	}

	return nil, nil, nil
}

// UpdateAfterSingleSuccessfulTxExecution updates the touched addresses after execute on Executor a successfully tx
func (w *Worker) UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.InfoReadWrite) []*TxTracker {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()
	if len(touchedAddresses) == 0 {
		log.Warnf("UpdateAfterSingleSuccessfulTxExecution touchedAddresses is nil or empty")
	}
	txsToDelete := make([]*TxTracker, 0)
	touchedFrom, found := touchedAddresses[from]
	if found {
		fromNonce, fromBalance := touchedFrom.Nonce, touchedFrom.Balance
		_, _, txsToDelete = w.applyAddressUpdate(from, fromNonce, fromBalance)
	} else {
		log.Warnf("UpdateAfterSingleSuccessfulTxExecution from(%s) not found in touchedAddresses", from.String())
	}

	for addr, addressInfo := range touchedAddresses {
		if addr != from {
			_, _, txsToDeleteTemp := w.applyAddressUpdate(addr, nil, addressInfo.Balance)
			txsToDelete = append(txsToDelete, txsToDeleteTemp...)
		}
	}
	return txsToDelete
}

// MoveTxToNotReady move a tx to not ready after it fails to execute
func (w *Worker) MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int) []*TxTracker {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()
	log.Infof("MoveTxToNotReady tx(%s) from(%s) actualNonce(%d) actualBalance(%s)", txHash.String(), from.String(), actualNonce, actualBalance.String())

	addrQueue, found := w.pool[from.String()]
	if found {
		// Sanity check. The txHash must be the readyTx
		if addrQueue.readyTx == nil || txHash.String() != addrQueue.readyTx.HashStr {
			readyHashStr := ""
			if addrQueue.readyTx != nil {
				readyHashStr = addrQueue.readyTx.HashStr
			}
			log.Warnf("MoveTxToNotReady txHash(%s) is not the readyTx(%s)", txHash.String(), readyHashStr)
		}
	}
	_, _, txsToDelete := w.applyAddressUpdate(from, actualNonce, actualBalance)

	return txsToDelete
}

// DeleteTx delete the tx after it fails to execute
func (w *Worker) DeleteTx(txHash common.Hash, addr common.Address) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[addr.String()]
	if found {
		deletedReadyTx := addrQueue.deleteTx(txHash)
		if deletedReadyTx != nil {
			log.Infof("DeleteTx tx(%s) deleted from TxSortedList", deletedReadyTx.Hash.String())
			w.txSortedList.delete(deletedReadyTx)
		}
	} else {
		log.Warnf("DeleteTx addrQueue(%s) not found", addr.String())
	}
}

// UpdateTxZKCounters updates the ZKCounter of a tx
func (w *Worker) UpdateTxZKCounters(txHash common.Hash, addr common.Address, counters state.ZKCounters) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	log.Infof("UpdateTxZKCounters tx(%s) addr(%s)", txHash.String(), addr.String())
	log.Debugf("UpdateTxZKCounters counters.CumulativeGasUsed: %d", counters.CumulativeGasUsed)
	log.Debugf("UpdateTxZKCounters counters.UsedKeccakHashes: %d", counters.UsedKeccakHashes)
	log.Debugf("UpdateTxZKCounters counters.UsedPoseidonHashes: %d", counters.UsedPoseidonHashes)
	log.Debugf("UpdateTxZKCounters counters.UsedPoseidonPaddings: %d", counters.UsedPoseidonPaddings)
	log.Debugf("UpdateTxZKCounters counters.UsedMemAligns: %d", counters.UsedMemAligns)
	log.Debugf("UpdateTxZKCounters counters.UsedArithmetics: %d", counters.UsedArithmetics)
	log.Debugf("UpdateTxZKCounters counters.UsedBinaries: %d", counters.UsedBinaries)
	log.Debugf("UpdateTxZKCounters counters.UsedSteps: %d", counters.UsedSteps)

	addrQueue, found := w.pool[addr.String()]

	if found {
		addrQueue.UpdateTxZKCounters(txHash, counters)
	} else {
		log.Warnf("UpdateTxZKCounters addrQueue(%s) not found", addr.String())
	}
}

// AddPendingTxToStore adds a tx to the addrQueue list of pending txs to store in the DB (trusted state)
func (w *Worker) AddPendingTxToStore(txHash common.Hash, addr common.Address) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[addr.String()]

	if found {
		addrQueue.addPendingTxToStore(txHash)
	} else {
		log.Warnf("AddPendingTxToStore addrQueue(%s) not found", addr.String())
	}
}

// DeletePendingTxToStore delete a tx from the addrQueue list of pending txs to store in the DB (trusted state)
func (w *Worker) DeletePendingTxToStore(txHash common.Hash, addr common.Address) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[addr.String()]

	if found {
		addrQueue.deletePendingTxToStore(txHash)
	} else {
		log.Warnf("DeletePendingTxToStore addrQueue(%s) not found", addr.String())
	}
}

// GetBestFittingTx gets the most efficient tx that fits in the available batch resources
func (w *Worker) GetBestFittingTx(resources state.BatchResources) *TxTracker {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	var (
		tx         *TxTracker
		foundMutex sync.RWMutex
	)

	nGoRoutines := runtime.NumCPU()
	foundAt := -1

	wg := sync.WaitGroup{}
	wg.Add(nGoRoutines)

	// Each go routine looks for a fitting tx
	for i := 0; i < nGoRoutines; i++ {
		go func(n int, bresources state.BatchResources) {
			defer wg.Done()
			for i := n; i < w.txSortedList.len(); i += nGoRoutines {
				foundMutex.RLock()
				if foundAt != -1 && i > foundAt {
					foundMutex.RUnlock()
					return
				}
				foundMutex.RUnlock()

				txCandidate := w.txSortedList.getByIndex(i)
				err := bresources.Sub(txCandidate.BatchResources)
				if err != nil {
					// We don't add this Tx
					continue
				}

				foundMutex.Lock()
				if foundAt == -1 || foundAt > i {
					foundAt = i
					tx = txCandidate
				}
				foundMutex.Unlock()

				return
			}
		}(i, resources)
	}
	wg.Wait()

	if foundAt != -1 {
		log.Infof("GetBestFittingTx found tx(%s) at index(%d) with gasPrice(%d)", tx.Hash.String(), foundAt, tx.GasPrice)
	} else {
		log.Debugf("GetBestFittingTx no tx found")
	}

	return tx
}

// ExpireTransactions deletes old txs
func (w *Worker) ExpireTransactions(maxTime time.Duration) []*TxTracker {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	var txs []*TxTracker

	log.Info("ExpireTransactions start. addrQueue len: ", len(w.pool))
	for _, addrQueue := range w.pool {
		subTxs, prevReadyTx := addrQueue.ExpireTransactions(maxTime)
		txs = append(txs, subTxs...)

		if prevReadyTx != nil {
			w.txSortedList.delete(prevReadyTx)
		}

		if addrQueue.IsEmpty() {
			delete(w.pool, addrQueue.fromStr)
		}
	}
	log.Info("ExpireTransactions end. addrQueue len: ", len(w.pool), " deleteCount: ", len(txs))

	return txs
}

// HandleL2Reorg handles the L2 reorg signal
func (w *Worker) HandleL2Reorg(txHashes []common.Hash) {
	log.Fatal("L2 Reorg detected. Restarting to sync with the new L2 state...")
}
