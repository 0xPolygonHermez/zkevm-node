package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Genesis contains the information to populate State on creation
type Genesis struct {
	Block          *types.Block
	Balances       map[common.Address]*big.Int
	SmartContracts map[common.Address][]byte
	Storage        map[common.Address]map[*big.Int]*big.Int
	L2ChainID      uint64
}
