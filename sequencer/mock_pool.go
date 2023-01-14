// Code generated by mockery v2.16.0. DO NOT EDIT.

package sequencer

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	pool "github.com/0xPolygonHermez/zkevm-node/pool"

	state "github.com/0xPolygonHermez/zkevm-node/state"
)

// PoolMock is an autogenerated mock type for the txPool type
type PoolMock struct {
	mock.Mock
}

// DeleteTransactionByHash provides a mock function with given fields: ctx, hash
func (_m *PoolMock) DeleteTransactionByHash(ctx context.Context, hash common.Hash) error {
	ret := _m.Called(ctx, hash)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) error); ok {
		r0 = rf(ctx, hash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTxsByHashes provides a mock function with given fields: ctx, hashes
func (_m *PoolMock) DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error {
	ret := _m.Called(ctx, hashes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Hash) error); ok {
		r0 = rf(ctx, hashes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPendingTxs provides a mock function with given fields: ctx, isClaims, limit
func (_m *PoolMock) GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error) {
	ret := _m.Called(ctx, isClaims, limit)

	var r0 []pool.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, bool, uint64) []pool.Transaction); ok {
		r0 = rf(ctx, isClaims, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pool.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, bool, uint64) error); ok {
		r1 = rf(ctx, isClaims, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxZkCountersByHash provides a mock function with given fields: ctx, hash
func (_m *PoolMock) GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, error) {
	ret := _m.Called(ctx, hash)

	var r0 *state.ZKCounters
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *state.ZKCounters); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.ZKCounters)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkReorgedTxsAsPending provides a mock function with given fields: ctx
func (_m *PoolMock) MarkReorgedTxsAsPending(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTxStatus provides a mock function with given fields: ctx, hash, newStatus
func (_m *PoolMock) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus) error {
	ret := _m.Called(ctx, hash, newStatus)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pool.TxStatus) error); ok {
		r0 = rf(ctx, hash, newStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPoolMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewPoolMock creates a new instance of PoolMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPoolMock(t mockConstructorTestingTNewPoolMock) *PoolMock {
	mock := &PoolMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
