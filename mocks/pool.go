package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
)

type PoolMock struct{}

func NewPool() pool.Pool {
	return &PoolMock{}
}

func (p *PoolMock) AddTx(tx types.Transaction) error {
	return nil
}

func (p *PoolMock) GetPendingTxs() ([]pool.Transaction, error) {
	return []pool.Transaction{{LegacyTx: *tx}}, nil
}

func (p *PoolMock) UpdateTxState(hash common.Hash, newState pool.TxState) error {
	return nil
}

func (p *PoolMock) CleanUpInvalidAndNonSelectedTxs() error {
	return nil
}

func (p *PoolMock) GetGasPrice() (uint64, error) {
	return gasPrice, nil
}
