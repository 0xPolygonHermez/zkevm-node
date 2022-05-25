package fakevm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type FakeDB struct {
	state.StateDB
}

func (f *FakeDB) GetBalance(address common.Address) *big.Int {
	return new(big.Int)
}

func (f *FakeDB) GetNonce(address common.Address) uint64 {
	return 0
}

func (f *FakeDB) GetCode(address common.Address) []byte {
	return []byte{}
}

func (f *FakeDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return hash
}

func (f *FakeDB) Exist(addr common.Address) bool {
	return true
}
