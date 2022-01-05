package runtime

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	Code        []byte
	Type        int
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
