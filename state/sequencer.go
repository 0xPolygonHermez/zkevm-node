package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Sequencer struct
type Sequencer struct {
	Sequencer common.Address
	URL       string
	ChainID   *big.Int
}
