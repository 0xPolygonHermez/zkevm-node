package statev2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Genesis contains the information to populate State on creation
type Genesis struct {
	Block          *types.Block                             `json:"-"`
	Balances       map[common.Address]*big.Int              `json:"balances"`
	SmartContracts map[common.Address][]byte                `json:"smartContracts"`
	Storage        map[common.Address]map[*big.Int]*big.Int `json:"storage"`
	Nonces         map[common.Address]*big.Int              `json:"nonces"`
	L2ChainID      uint64                                   `json:"-"`
}
