package state

import (
	"github.com/ethereum/go-ethereum/common"
)

// TokenWrapped struct
type TokenWrapped struct {
	OriginalNetwork      uint
	OriginalTokenAddress common.Address
	WrappedTokenAddress  common.Address
}
