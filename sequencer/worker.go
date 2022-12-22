package sequencer

import (
	"math/big"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Worker struct {
	Pool           map[common.Address]AddrQueue // This should have (almost) all txs from the pool
	efficiencyList efficiencyList
}

func newWorker() *Worker {
	// TODO: Initialize memory structs
	return &Worker{}
}

func (w *Worker) AddTx(tx TxTracker) {
	// 1. Check if addr exists on Pool
	// // If not: create and get nonce and balance from MT
	// 2. Add tx to the AddrQueue if there is a tx with the same nonce and the existing tx has better gas price we keep the existing tx and discard the other
	// 3. Check if the new tx is ready, if so:
	// // A) There wasnt a tx ready => add the tx to the efficiencyList
	// // B) There was a tx ready (and it's worst than the new one) => delete from pool and efficiency list, add new one
}

func (w *Worker) UpdateAfterSingleSuccessfulTxExecution(from common.Address, fromNonce uint64, fromBalance *big.Int, balances map[common.Address]*big.Int) {
	newReadyTx, prevReadyTx := w.Pool[from].UpdateCurrentNonceBalance(&fromNonce, fromBalance)
	if prevReadyTx != nil {
		w.efficiencyList.Delete(prevReadyTx.Hash)
	}
	if newReadyTx != nil {
		w.efficiencyList.Add(*newReadyTx)
	}
	for addr, balance := range balances {
		newReadyTx, prevReadyTx = w.Pool[addr].UpdateCurrentNonceBalance(nil, balance)
		if prevReadyTx != nil {
			w.efficiencyList.Delete(prevReadyTx.Hash)
		}
		if newReadyTx != nil {
			w.efficiencyList.Add(*newReadyTx)
		}
	}
}

// Assume that finalizer detected that a tx was not ready and decides to move to not ready after it fails to execute AND DOESN'T MODIFY THE STATE
func (w *Worker) MoveTxToNotReady(from common.Address, txHash common.Hash, actualNonce *uint64, actualBalance *big.Int) {
	w.UpdateAfterSingleSuccessfulTxExecution(from, *actualNonce, actualBalance, nil)
}

// Assume that finalizer decides to delete the tx after it fails to execute AND DOESN'T MODIFY THE STATE
func (w *Worker) DeleteTx(txHash common.Hash, from common.Address, actualFromNonce *uint64, actualFromBalance *big.Int) {
	/*
		1. Delete from w.Pool and w.efficiencyList
		2. Update w.Pool with nonce/balance if they are not nil
		3. Potentially delete more txs if nonce has been updated
		4. Potentially select new ReadyTx and add it into efficiecnyList
	*/
}

func (w *Worker) UpdateTx(txHash common.Hash, from common.Address, ZKCounters pool.ZkCounters) {
	// 1. Get tx from Pool
	// 2. Set ZKCounters
	// 3. Calculate new wfficiency
	// 4. Resort tx from efficiency
}

// TODO: separate UpdateAfterSingleSuccessfulTxExecution in different functions so it can be nicely reused by MoveTxToNotReady

type AddrQueue struct {
	CurrentNonce   uint64
	CurrentBalance *big.Int
	ReadyTx        *TxTracker
	NotReadyTxs    []TxTracker
}

func (a AddrQueue) UpdateCurrentNonceBalance(nonce *uint64, balance *big.Int) (newReadyTx, prevReadyTx *TxTracker) {
	// 1. Set nonce, balance
	// 2. If ReadyTx != nil, and not longer ready:
	// // prevReadyTx = ReadyTx
	// // tmpTx = ReadyTx
	// // ReadyTx = nil
	// // conisder moving tmpTx to NotReadyTx
	// 3. If ReadyTx == nil, check if NotReadyTxs[0] can be moved to ReadyTx and if so newReadyTx = ReadyTx
	return nil, nil
}

type efficiencyList map[common.Hash]TxTracker // only ready txs. Replace map for sorted map. TODO: find good library

func (e *efficiencyList) Add(tx TxTracker)   {}
func (e *efficiencyList) Delete(common.Hash) {}
func (e *efficiencyList) GetMostEfficientByIndex(i int) *TxTracker {
	return &TxTracker{}
}

type RemainingResources struct {
	remainingZKCounters pool.ZkCounters
	remainingBytes      uint64
	remainingGas        uint64
}

func (r *RemainingResources) sub(tx TxTracker) error {
	// Substract resources
	// error if underflow (restore in this case)
	return nil
}

type TxTracker struct {
	Hash       common.Hash
	addrQueue  *AddrQueue
	Nonce      uint64
	Benefit    *big.Int        // GasLimit * GasPrice
	ZKCounters pool.ZkCounters // To check if it fits into a batch
	Size       uint64          // To check if it fits into a batch
	Gas        uint64          // To check if it fits into a batch
	GasPrice   int64
	Efficiency float64 // To sort. TODO: calculate Benefit / Cost. Cost = some formula taking into account ZKC and Byte Size
	RawTx      []byte
}

func NewTxTracker(tx types.Transaction, counters pool.ZkCounters) *TxTracker {
	txTracker := &TxTracker{
		// Set values
	}
	txTracker.CalculateEfficiency()
	return txTracker
}

func (tx *TxTracker) CalculateEfficiency() {
	// TODO: define efficiency formula, this is just a draft
	const (
		UsedArithmeticsWeight = 0.1
	)
	cost := float64(tx.ZKCounters.UsedArithmetics) * UsedArithmeticsWeight // do for all counters ... AND size
	benefit := tx.ZKCounters.CumulativeGasUsed * tx.GasPrice
	tx.Efficiency = float64(benefit) / cost
}

func (p *Worker) getMostEfficientTx() (TxTracker, error) {
	return TxTracker{}, nil
}

func (p *Worker) len() int {
	return 0
}

func (p *Worker) GetBestFittingTx(resources RemainingResources) *TxTracker {
	var tx *TxTracker
	nGoRoutines := 4 // nCores - K // TODO: Think about this

	// Each go routine looks for a fitting tx
	foundAt := -1 // TODO: add mutex
	wg := sync.WaitGroup{}
	wg.Add(nGoRoutines)
	for i := 0; i < nGoRoutines; i++ {
		go func(n int) {
			defer wg.Done()
			for i := n; i < len(p.efficiencyList); i += nGoRoutines {
				if i > foundAt {
					return
				}
				txCandidate := p.efficiencyList.GetMostEfficientByIndex(i)
				err := resources.sub(*txCandidate)
				if err != nil {
					// We don't add this Tx
					continue
				}

				if foundAt == -1 || foundAt > i {
					foundAt = i
					tx = txCandidate
				}
				return
			}
		}(i)
	}
	wg.Wait()
	return tx
}
