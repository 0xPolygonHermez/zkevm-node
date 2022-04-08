package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type globalExitRoot struct {
	BlockID           uint64
	BlockNumber       uint64
	GlobalExitRootNum *big.Int
	ExitRoots         []common.Hash
}
