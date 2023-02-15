package sequencer

import (
	"context"
	"math/big"
	"runtime"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Worker represents the worker component of the sequencer
type Worker struct {
	pool           map[string]*addrQueue
	efficiencyList *efficiencyList
	workerMutex    sync.Mutex
	// dbManager            dbManagerInterface
	state                stateInterface
	batchConstraints     batchConstraints
	batchResourceWeights batchResourceWeights
}

// NewWorker creates an init a worker
func NewWorker(state stateInterface, constraints batchConstraints, weights batchResourceWeights) *Worker {
	w := Worker{
		pool:                 make(map[string]*addrQueue),
		efficiencyList:       newEfficiencyList(),
		state:                state,
		batchConstraints:     constraints,
		batchResourceWeights: weights,
	}

	return &w
}

// NewTxTracker creates and inits a TxTracker
func (w *Worker) NewTxTracker(tx types.Transaction, isClaim bool, counters state.ZKCounters) (*TxTracker, error) {
	return newTxTracker(tx, isClaim, counters, w.batchConstraints, w.batchResourceWeights)
}

// AddTx adds a new Tx to the Worker
// TODO: Rename to AddTxTracker?
func (w *Worker) AddTx(ctx context.Context, tx *TxTracker) {
	// TODO: Review if additional mutex is needed to lock GetBestFittingTx
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addr, found := w.pool[tx.FromStr]

	if !found {
		// Unlock the worker to let execute other worker functions while creating the new AddrQueue
		w.workerMutex.Unlock()

		root, err := w.state.GetLastStateRoot(ctx, nil)
		if err != nil {
			// TODO: How to manage this
			return
		}
		nonce, err := w.state.GetNonceByStateRoot(ctx, tx.From, root)
		if err != nil {
			// TODO: How to manage this
			return
		}
		balance, err := w.state.GetBalanceByStateRoot(ctx, tx.From, root)
		if err != nil {
			// TODO: How to manage this
			return
		}

		addr = newAddrQueue(tx.From, nonce.Uint64(), balance)

		// Lock again the worker
		w.workerMutex.Lock()

		w.pool[tx.FromStr] = addr
	}

	// Add the txTracker to Addr and get the newReadyTx and prevReadyTx
	newReadyTx, prevReadyTx := addr.addTx(tx)

	// Update the EfficiencyList (if needed)
	if prevReadyTx != nil {
		w.efficiencyList.delete(prevReadyTx)
	}
	if newReadyTx != nil {
		w.efficiencyList.add(newReadyTx)
	}
}

func (w *Worker) applyAddressUpdate(from common.Address, fromNonce *uint64, fromBalance *big.Int) (*TxTracker, *TxTracker) {
	addrQueue, found := w.pool[from.String()]

	// TODO: What happens if addr no found. Could it be possible if addrQueue has not been yet created for this from addr (touchedAddresses)
	if found {
		newReadyTx, prevReadyTx := addrQueue.updateCurrentNonceBalance(fromNonce, fromBalance)

		// Update the EfficiencyList (if needed)
		if prevReadyTx != nil {
			w.efficiencyList.delete(prevReadyTx)
		}
		if newReadyTx != nil {
			w.efficiencyList.add(newReadyTx)
		}

		return newReadyTx, prevReadyTx
	}

	return nil, nil
}

// UpdateAfterSingleSuccessfulTxExecution updates the touched addresses after execute on Executor a successfully tx
func (w *Worker) UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.InfoReadWrite) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()
	if touchedAddresses == nil || len(touchedAddresses) == 0 {
		log.Errorf("UpdateAfterSingleSuccessfulTxExecution touchedAddresses is nil or empty")
	}

	touchedFrom, found := touchedAddresses[from]
	if found {
		fromNonce, fromBalance := touchedFrom.Nonce, touchedFrom.Balance
		w.applyAddressUpdate(from, fromNonce, fromBalance)
	} else {
		log.Errorf("UpdateAfterSingleSuccessfulTxExecution from(%s) not found in touchedAddresses", from.String())
	}

	for addr, addressInfo := range touchedAddresses {
		if addr != from {
			w.applyAddressUpdate(addr, nil, addressInfo.Balance)
		}
	}
}

// MoveTxToNotReady move a tx to not ready after it fails to execute
func (w *Worker) MoveTxToNotReady(txHash common.Hash, from common.Address, actualNonce *uint64, actualBalance *big.Int) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[from.String()]

	if found {
		// Sanity check. The txHash must be the readyTx
		if addrQueue.readyTx == nil || txHash.String() != addrQueue.readyTx.HashStr {
			readyHashStr := ""
			if addrQueue.readyTx != nil {
				readyHashStr = addrQueue.readyTx.HashStr
			}
			log.Errorf("MoveTxToNotReady txHash(s) is not the readyTx(%s)", txHash.String(), readyHashStr)
			// TODO: how to manage this?
		}
	}

	w.applyAddressUpdate(from, actualNonce, actualBalance)
}

// DeleteTx delete the tx after it fails to execute
func (w *Worker) DeleteTx(txHash common.Hash, addr common.Address) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[addr.String()]
	if found {
		deletedReadyTx := addrQueue.deleteTx(txHash)
		if deletedReadyTx != nil {
			w.efficiencyList.delete(deletedReadyTx)
		}
	} else {
		log.Errorf("DeleteTx addrQueue(%s) not found", addr.String())
	}
}

// UpdateTx updates the ZKCounter of a tx and resort the tx in the efficiency list if needed
func (w *Worker) UpdateTx(txHash common.Hash, addr common.Address, counters state.ZKCounters) {
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	addrQueue, found := w.pool[addr.String()]

	if found {
		newReadyTx, prevReadyTx := addrQueue.UpdateTxZKCounters(txHash, counters, w.batchConstraints, w.batchResourceWeights)

		// Resort the newReadyTx in efficiencyList
		if prevReadyTx != nil {
			w.efficiencyList.delete(prevReadyTx)
		}
		if newReadyTx != nil {
			w.efficiencyList.add(newReadyTx)
		}
	} else {
		log.Errorf("UpdateTx addrQueue(%s) not found", addr.String())
	}
}

// GetBestFittingTx gets the most efficient tx that fits in the available batch resources
func (w *Worker) GetBestFittingTx(resources batchResources) *TxTracker {
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
		go func(n int, bresources batchResources) {
			defer wg.Done()
			for i := n; i < w.efficiencyList.len(); i += nGoRoutines {
				foundMutex.RLock()
				if foundAt != -1 && i > foundAt {
					foundMutex.RUnlock()
					return
				}
				foundMutex.RUnlock()

				txCandidate := w.efficiencyList.getByIndex(i)
				err := bresources.sub(txCandidate.BatchResources)
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

	return tx
}

// HandleL2Reorg handles the L2 reorg signal
func (w *Worker) HandleL2Reorg(txHashes []common.Hash) {
	// 1. Delete related txs from w.efficiencyList
	// 2. Mark the affected addresses as "reorged" in w.Pool
	// 3. Update these addresses (go to MT, update nonce and balance into w.Pool)
}
