package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Sequencer struct
type Sequencer struct {
	Address     common.Address
	URL         string
	ChainID     *big.Int
	BlockNumber uint64
}
