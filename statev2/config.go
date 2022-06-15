package state

import "github.com/ethereum/go-ethereum/common"

// Config is state config
type Config struct {
	// DefaultChainID is the common ChainID to all the sequencers
	DefaultChainID uint64
	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64
	// L2GlobalExitRootManagerAddr is the L2 address of the exit root manager SC
	L2GlobalExitRootManagerAddr common.Address
	// GlobalExitRootStoragePosition is the position inside SC's storage to write the new global state root
	GlobalExitRootStoragePosition uint64
	// LocalExitRootStoragePosition is the position inside SC's storage to read the new local state root
	LocalExitRootStoragePosition uint64
	// SystemSCAddr is the System Smart Contract Address
	SystemSCAddr common.Address
	// OldStateRootPosition is the position in SystemSCAddr to store the old root
	OldStateRootPosition uint64
}
