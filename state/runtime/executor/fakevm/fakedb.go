package fakevm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type FakeDB struct {
	StateRoot []byte
	state.StateDB
}

func (f *FakeDB) GetBalance(address common.Address) *big.Int {
	panic("GetBalance NOT IMPLEMENTED")
}

func (f *FakeDB) GetNonce(address common.Address) uint64 {
	panic("GetNonce NOT IMPLEMENTED")
}

func (f *FakeDB) GetCode(address common.Address) []byte {
	panic("GetCode NOT IMPLEMENTED")
}

func (f *FakeDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	panic("GetState NOT IMPLEMENTED")
}

func (f *FakeDB) Exist(addr common.Address) bool {
	panic("GetState NOT IMPLEMENTED")
}
