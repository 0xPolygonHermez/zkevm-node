package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Genesis contains the information to populate State on creation
type Genesis struct {
	Balances       map[common.Address]*big.Int
	SmartContracts map[common.Address][]byte
}
