package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Transaction contains metadata of a Tx
type Transaction struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *common.Address
	Value    *big.Int
	Input    []byte
	Hash     common.Hash
	From     common.Address
}

// IsContractCreation check if a Transaction is a contract creation
func (t *Transaction) IsContractCreation() bool {
	return t.To == nil
}

// Copy creates a copy of a Transaction
func (t *Transaction) Copy() *Transaction {
	tt := new(Transaction)
	*tt = *t

	tt.GasPrice = new(big.Int)
	tt.GasPrice.Set(t.GasPrice)

	tt.Value = new(big.Int)
	tt.Value.Set(t.Value)

	tt.Input = make([]byte, len(t.Input))
	copy(tt.Input[:], t.Input[:])
	return tt
}
