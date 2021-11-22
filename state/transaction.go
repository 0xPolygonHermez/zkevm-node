package state

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *common.Address
	Value    *big.Int
	Input    []byte
	V        []byte
	R        []byte
	S        []byte
	hash     atomic.Value
	From     common.Address
}

func (t *Transaction) IsContractCreation() bool {
	return t.To == nil
}

// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	h := rlpHash(tx)

	tx.hash.Store(h)
	return h
}

func (t *Transaction) Copy() *Transaction {
	tt := new(Transaction)
	*tt = *t

	tt.GasPrice = new(big.Int)
	tt.GasPrice.Set(t.GasPrice)

	tt.Value = new(big.Int)
	tt.Value.Set(t.Value)

	tt.R = make([]byte, len(t.R))
	copy(tt.R[:], t.R[:])
	tt.S = make([]byte, len(t.S))
	copy(tt.S[:], t.S[:])

	tt.Input = make([]byte, len(t.Input))
	copy(tt.Input[:], t.Input[:])
	return tt
}

// Cost returns gas * gasPrice + value
func (t *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(t.GasPrice, new(big.Int).SetUint64(t.Gas))
	total.Add(total, t.Value)
	return total
}
