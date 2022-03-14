package state

import "github.com/ethereum/go-ethereum/common"

// Config is state config
type Config struct {
	// DefaultChainID is the common ChainID to all the sequencers
	DefaultChainID uint64
	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64
	// Address of the exit root manager SC
	GlobalExitRootManagerAddr common.Address
	// Position inside SC's storage to read the new local state root
	GlobalExitRootManagerPosition uint64
}
