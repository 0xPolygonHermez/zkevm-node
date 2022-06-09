package fakevm

import "github.com/ethereum/go-ethereum/common"

type Account struct {
	address common.Address
}

func NewAccount(address common.Address) *Account {
	return &Account{address: address}
}

func (a *Account) Address() common.Address { return a.address }
