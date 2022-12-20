package docs

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/common"
)

type WorkerPool map[common.Address]AddrQueue // Replace map for sorted map. TODO: find good library
type AddrQueue struct {
	CurrentNonce   uint64
	CurrentBalance *big.Int
	ReadyTxs       []TxTracker
	NotReadyTxs    []TxTracker
}

type RemainingResources struct {
	remainingZKCounters pool.ZkCounters
	remainingBytes      uint64
	remainingGas        uint64
}

func (r *RemainingResources) sub(tx TxTracker) error {
	// Substract resources
	// error if underflow (restore in this case)
}

type TxTracker struct {
	Addr       *AddrQueue
	Nonce      uint64
	Benefit    *big.Int        // GasLimit * GasPrice
	ZKCounters pool.ZkCounters // To check if it fits into a batch
	Size       uint64          // To check if it fits into a batch
	Gas        uint64          // To check if it fits into a batch
	Efficiency float64         // To sort. TODO: calculate Benefit / Cost. Cost = some formula taking into account ZKC and Byte Size
	RawTx      []byte
}

func (p *WorkerPool) getMostEfficientTx() (TxTracker, error) {
	return TxTracker{}, nil
}

func (p *WorkerPool) len() int {
	return 0
}

func (p *WorkerPool) fitTransactions(resources RemainingResources) []TxTracker {
	var txs []TxTracker

	for {
		// Get Most eficient TX in the Pool
		tx, err := p.getMostEfficientTx()
		if err != nil {
			break
		}

		// Check if the tx fits into the batch (AKA also check Gas and bytes)
		err = resources.sub(tx)
		if err != nil {
			// We don't add this Tx
			break
		}

		txs = append(txs, tx)
		/*
		Updates: (convert this to reusable function)
			1 - Remove tx from the efficiency list
			2 - Remove tx from the address list
			3 - if tx.addr.ReadyTxs[0] != nil add tx.addr.ReadyTxs[0] to efficiency list

			* Consider Processing already the next transaction (Niko's Idea)
				
		*/

	nGoRoutines := nCores - K

	for i := 0; i < nGoRoutines; i++ {
		go func(n int) {
			for i := n; i < len(efficiencyList); i += nGoRoutines {
				tx := efficiencyList[i]
				err = resources.sub(tx)
				if err != nil {
					// We don't add this Tx
					continue
				}

				txs = append(txs, tx)

				/*
				Updates: (convert this to reusable function)
					1 - Remove tx from the efficiency list
					2 - Remove tx from the address list
					3 - if tx.addr.ReadyTxs[0] != nil add tx.addr.ReadyTxs[0] to efficiency list

					* Consider Processing already the next transaction (Niko's Idea)
				
				*/

			}
		}(i)
	}

	return txs
}
