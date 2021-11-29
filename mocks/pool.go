//nolint
package mocks

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
)

type PoolMock struct{}

func NewPool() pool.Pool {
	return &PoolMock{}
}

func (p *PoolMock) AddTx(ctx context.Context, tx types.Transaction) error {
	return nil
}

func (p *PoolMock) GetPendingTxs(ctx context.Context) ([]pool.Transaction, error) {
	return []pool.Transaction{{Transaction: *tx}}, nil
}

func (p *PoolMock) UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error {
	return nil
}

func (p *PoolMock) CleanUpInvalidAndNonSelectedTxs(ctx context.Context) error {
	return nil
}

func (p *PoolMock) SetGasPrice(ctx context.Context, gasPrice uint64) error {
	return nil
}

func (p *PoolMock) GetGasPrice(ctx context.Context) (uint64, error) {
	return gasPrice, nil
}
