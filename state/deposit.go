package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Deposit struct
type Deposit struct {
	TokenAddres        common.Address
	Amount             *big.Int
	DestinationNetwork uint
	DestinationAddress common.Address
	BlockNumber        uint64
}
