package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Claim struct
type Claim struct {
	Index              uint64
	OriginalNetwork    uint
	Token              common.Address
	Amount             *big.Int
	DestinationAddress common.Address
	BlockNumber        uint64
}
