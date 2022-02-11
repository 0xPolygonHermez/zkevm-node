package state

import "math/big"

// GlobalExitRoot struct
type GlobalExitRoot struct {
	GlobalExitRootNum *big.Int
	MainnetExitRoot   [32]byte
	RollupExitRoot    [32]byte
}
