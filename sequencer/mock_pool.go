// Code generated by mockery v2.32.0. DO NOT EDIT.

package sequencer

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	pool "github.com/0xPolygonHermez/zkevm-node/pool"

	state "github.com/0xPolygonHermez/zkevm-node/state"

	time "time"
)

// PoolMock is an autogenerated mock type for the txPool type
type PoolMock struct {
	mock.Mock
}

// DeleteFailedTransactionsOlderThan provides a mock function with given fields: ctx, date
func (_m *PoolMock) DeleteFailedTransactionsOlderThan(ctx context.Context, date time.Time) error {
	ret := _m.Called(ctx, date)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) error); ok {
		r0 = rf(ctx, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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

// DeleteTransactionsByHashes provides a mock function with given fields: ctx, hashes
func (_m *PoolMock) DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error {
	ret := _m.Called(ctx, hashes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Hash) error); ok {
		r0 = rf(ctx, hashes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDefaultMinGasPriceAllowed provides a mock function with given fields:
func (_m *PoolMock) GetDefaultMinGasPriceAllowed() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// GetGasPrices provides a mock function with given fields: ctx
func (_m *PoolMock) GetGasPrices(ctx context.Context) (pool.GasPrices, error) {
	ret := _m.Called(ctx)

	var r0 pool.GasPrices
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (pool.GasPrices, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) pool.GasPrices); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(pool.GasPrices)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetL1GasPrice provides a mock function with given fields:
func (_m *PoolMock) GetL1GasPrice() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// GetNonWIPPendingTxs provides a mock function with given fields: ctx
func (_m *PoolMock) GetNonWIPPendingTxs(ctx context.Context) ([]pool.Transaction, error) {
	ret := _m.Called(ctx)

	var r0 []pool.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]pool.Transaction, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []pool.Transaction); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pool.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxZkCountersByHash provides a mock function with given fields: ctx, hash
func (_m *PoolMock) GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, error) {
	ret := _m.Called(ctx, hash)

	var r0 *state.ZKCounters
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*state.ZKCounters, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *state.ZKCounters); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.ZKCounters)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkWIPTxsAsPending provides a mock function with given fields: ctx
func (_m *PoolMock) MarkWIPTxsAsPending(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTxStatus provides a mock function with given fields: ctx, hash, newStatus, isWIP, failedReason
func (_m *PoolMock) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool, failedReason *string) error {
	ret := _m.Called(ctx, hash, newStatus, isWIP, failedReason)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pool.TxStatus, bool, *string) error); ok {
		r0 = rf(ctx, hash, newStatus, isWIP, failedReason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTxWIPStatus provides a mock function with given fields: ctx, hash, isWIP
func (_m *PoolMock) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	ret := _m.Called(ctx, hash, isWIP)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, bool) error); ok {
		r0 = rf(ctx, hash, isWIP)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPoolMock creates a new instance of PoolMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPoolMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PoolMock {
	mock := &PoolMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
