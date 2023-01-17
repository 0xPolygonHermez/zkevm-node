package sequencer

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

// addrQueue is a struct that stores the ready and notReady txs for a specific from address
type addrQueue struct {
	from           common.Address
	fromStr        string
	currentNonce   uint64
	currentBalance *big.Int
	readyTx        *TxTracker
	notReadyTxs    map[uint64]*TxTracker
}

// newAddrQueue creates and init a addrQueue
func newAddrQueue(addr common.Address, nonce uint64, balance *big.Int) *addrQueue {
	return &addrQueue{
		from:           addr,
		fromStr:        addr.String(),
		currentNonce:   nonce,
		currentBalance: balance,
		readyTx:        nil,
		notReadyTxs:    make(map[uint64]*TxTracker),
	}
}

// addTx adds a tx to the addrQueue and updates the ready a notReady Txs
func (a *addrQueue) addTx(tx *TxTracker) (newReadyTx, prevReadyTx *TxTracker) {
	if a.currentNonce == tx.Nonce { // Is a possible readyTx
		// We set the tx as readyTx if we do not have one assigned or if the gasPrice is better or equal than the current readyTx
		if a.readyTx == nil || ((a.readyTx != nil) && (tx.GasPrice.Cmp(a.readyTx.GasPrice) >= 0)) {
			prevTx := a.readyTx
			if a.currentBalance.Cmp(tx.Cost) >= 0 { //
				a.readyTx = tx
				return tx, prevTx
			} else { // If there is not enought balance we set the new tx as notReadyTxs
				a.readyTx = nil
				a.notReadyTxs[tx.Nonce] = tx
				return nil, prevTx
			}
		}
	} else {
		// TODO: how to manage if the tx.Nonce < currentNonce, we add it to notReady? if not, we need to update in DB?
		// We add the tx to the notReadyTxs list if it does not exists or if it already exists but has a better gasPrice
		nrTx, found := a.notReadyTxs[tx.Nonce]
		if !found || ((found) && (tx.GasPrice.Cmp(nrTx.GasPrice) >= 0)) {
			prevTx := nrTx
			a.notReadyTxs[tx.Nonce] = tx
			return tx, prevTx
		}
	}

	return nil, nil
}

// deleteTx deletes the tx from the addrQueue
func (a *addrQueue) deleteTx(txHash common.Hash) (deletedReadyTx *TxTracker) {
	txHashStr := txHash.String()

	if (a.readyTx != nil) && (a.readyTx.HashStr == txHashStr) {
		prevReadyTx := a.readyTx
		a.readyTx = nil
		return prevReadyTx
	} else {
		for _, txTracker := range a.notReadyTxs {
			if txTracker.HashStr == txHashStr {
				delete(a.notReadyTxs, txTracker.Nonce)
				break
			}
		}
		return nil
	}
}

// updateCurrentNonceBalance updates the nonce and balance of the addrQueue and updates the ready and notReady txs
func (a *addrQueue) updateCurrentNonceBalance(nonce *uint64, balance *big.Int) (newReadyTx, prevReadyTx *TxTracker) {
	var oldReadyTx *TxTracker = nil

	a.currentBalance = balance

	if nonce != nil {
		if a.currentNonce != *nonce {
			a.currentNonce = *nonce

			//TODO: we need to update in the DB the deleted txs?
			txToDelete := []uint64{}
			for _, txTracker := range a.notReadyTxs {
				if txTracker.Nonce < a.currentNonce {
					txToDelete = append(txToDelete, txTracker.Nonce)
				}
			}
			for _, delTxNonce := range txToDelete {
				delete(a.notReadyTxs, delTxNonce)
			}
		}
	}

	if a.readyTx != nil {
		// If readyTX.nonce is not the currentNonce or currentBalance is less that the readyTx.Cost
		// set readyTx=nil. Later we will move the tx to notReadyTxs
		if (a.readyTx.Nonce != a.currentNonce) || (a.currentBalance.Cmp(a.readyTx.Cost) < 0) {
			oldReadyTx = a.readyTx
			a.readyTx = nil
		}
	}

	// We check if we have a new readyTx from the notReadyTxs (at this point, to optmize the code,
	// we are not including the oldReadyTx in notReadyTxs, as it can match again if the nonce has not changed)
	if a.readyTx == nil {
		nrTx, found := a.notReadyTxs[a.currentNonce]
		if found {
			if a.currentBalance.Cmp(nrTx.Cost) >= 0 {
				a.readyTx = nrTx
				delete(a.notReadyTxs, a.currentNonce)
			}
		}
	}

	// We add the oldReadyTx to notReadyTxs (if it has a valid nonce) at this point to avoid check it again in the previous if statement
	if oldReadyTx != nil && oldReadyTx.Nonce > a.currentNonce {
		a.notReadyTxs[oldReadyTx.Nonce] = oldReadyTx
	}

	return a.readyTx, oldReadyTx
}

// UpdateTxZKCounters updates the ZKCounters for the given tx (txHash)
// If the updated tx is the readyTx it returns a copy of the previous readyTx, nil otherwise
func (a *addrQueue) UpdateTxZKCounters(txHash common.Hash, counters state.ZKCounters, constraints batchConstraints, weights batchResourceWeights) (newReadyTx, prevReadyTx *TxTracker) {
	txHashStr := txHash.String()

	if (a.readyTx != nil) && (a.readyTx.HashStr == txHashStr) {
		// We need to assign the new readyTx as a new TxTracker copy of the previous one with the updated efficiency
		// We need to do in this way because the efficiency value is changed and we use this value as key field to
		// add/delete TxTrackers in the efficiencyList
		prevReadyTx := a.readyTx
		newReadyTx := *a.readyTx
		newReadyTx.updateZKCounters(counters, constraints, weights)
		a.readyTx = &newReadyTx
		return a.readyTx, prevReadyTx
	} else { // TODO: This makes sense or we need only to check the readyTx
		txHashStr := txHash.String()
		for _, txTracker := range a.notReadyTxs {
			if txTracker.HashStr == txHashStr {
				txTracker.updateZKCounters(counters, constraints, weights)
				break
			}
		}
		return nil, nil
	}
}
