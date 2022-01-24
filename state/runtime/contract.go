package runtime

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Contract is the instance being called
type Contract struct {
	Code        []byte
	Type        CallType
	CodeAddress common.Address
	Address     common.Address
	Origin      common.Address
	Caller      common.Address
	Depth       int
	Value       *big.Int
	Input       []byte
	Gas         uint64
	Static      bool
}

// NewContract is the contract default constructor
func NewContract(depth int, origin common.Address, from common.Address, to common.Address, value *big.Int, gas uint64, code []byte) *Contract {
	contract := &Contract{
		Caller:      from,
		Origin:      origin,
		CodeAddress: to,
		Address:     to,
		Gas:         gas,
		Value:       value,
		Code:        code,
		Depth:       depth,
	}
	return contract
}

// NewContractCreation is used for contracts creation
func NewContractCreation(depth int, origin common.Address, from common.Address, to common.Address, value *big.Int, gas uint64, code []byte) *Contract {
	c := NewContract(depth, origin, from, to, value, gas, code)
	return c
}

// NewContractCall is used to call a contract
func NewContractCall(depth int, origin common.Address, from common.Address, to common.Address, value *big.Int, gas uint64, code []byte, input []byte) *Contract {
	c := NewContract(depth, origin, from, to, value, gas, code)
	c.Input = input
	return c
}
