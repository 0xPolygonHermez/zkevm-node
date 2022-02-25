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
	OriginNetwork      uint
	DestinationAddress common.Address
	DepositCount       uint
	BlockNumber        uint64
}
