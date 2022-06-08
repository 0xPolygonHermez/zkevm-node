package fakevm

import "github.com/ethereum/go-ethereum/common"

type account struct {
	address common.Address
}

func NewAccount(address common.Address) *account {
	return &account{address: address}
}

func (a *account) Address() common.Address { return a.address }
