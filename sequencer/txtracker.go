package sequencer

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxTracker is a struct that contains all the tx data needed to be managed by the worker
type TxTracker struct {
	Hash                              common.Hash
	HashStr                           string
	From                              common.Address
	FromStr                           string
	Nonce                             uint64
	Gas                               uint64 // To check if it fits into a batch
	GasPrice                          *big.Int
	Cost                              *big.Int             // Cost = Amount + Benefit
	Benefit                           *big.Int             // GasLimit * GasPrice
	BatchResources                    state.BatchResources // To check if it fits into a batch
	RawTx                             []byte
	ReceivedAt                        time.Time // To check if it has been in the txSortedList for too long
	IP                                string    // IP of the tx sender
	FailedReason                      *string   // FailedReason is the reason why the tx failed, if it failed
	BreakEvenGasPrice                 *big.Int
	GasPriceEffectivePercentage       uint8
	EffectiveGasPriceProcessCount     uint8
	IsEffectiveGasPriceFinalExecution bool
	L1GasPrice                        uint64
}

// newTxTracker creates and inti a TxTracker
func newTxTracker(tx types.Transaction, counters state.ZKCounters, ip string) (*TxTracker, error) {
	addr, err := state.GetSender(tx)
	if err != nil {
		return nil, err
	}

	rawTx, err := state.EncodeTransactionWithoutEffectivePercentage(tx)
	if err != nil {
		return nil, err
	}
	txTracker := &TxTracker{
		Hash:     tx.Hash(),
		HashStr:  tx.Hash().String(),
		From:     addr,
		FromStr:  addr.String(),
		Nonce:    tx.Nonce(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Cost:     tx.Cost(),
		Benefit:  new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice()),
		BatchResources: state.BatchResources{
			Bytes:      tx.Size(),
			ZKCounters: counters,
		},
		RawTx:                             rawTx,
		ReceivedAt:                        time.Now(),
		IP:                                ip,
		BreakEvenGasPrice:                 new(big.Int).SetUint64(0),
		EffectiveGasPriceProcessCount:     0,
		IsEffectiveGasPriceFinalExecution: false,
	}

	return txTracker, nil
}

// updateZKCounters updates the counters of the tx
func (tx *TxTracker) updateZKCounters(counters state.ZKCounters) {
	tx.BatchResources.ZKCounters = counters
}
