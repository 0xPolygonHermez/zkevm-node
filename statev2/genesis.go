package statev2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Genesis contains the information to populate Statev2 on creation
type Genesis struct {
	Balances       map[common.Address]*big.Int              `json:"balances"`
	SmartContracts map[common.Address][]byte                `json:"smartContracts"`
	Storage        map[common.Address]map[*big.Int]*big.Int `json:"storage"`
	Nonces         map[common.Address]*big.Int              `json:"nonces"`
}
