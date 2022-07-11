// Code generated by mockery v2.13.1. DO NOT EDIT.

package jsonrpcv2

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	pool "github.com/0xPolygonHermez/zkevm-node/pool"

	time "time"

	types "github.com/ethereum/go-ethereum/core/types"
)

// poolMock is an autogenerated mock type for the jsonRPCTxPool type
type poolMock struct {
	mock.Mock
}

// AddTx provides a mock function with given fields: ctx, tx
func (_m *poolMock) AddTx(ctx context.Context, tx types.Transaction) error {
	ret := _m.Called(ctx, tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Transaction) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetGasPrice provides a mock function with given fields: ctx
func (_m *poolMock) GetGasPrice(ctx context.Context) (uint64, error) {
	ret := _m.Called(ctx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPendingTxHashesSince provides a mock function with given fields: ctx, since
func (_m *poolMock) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	ret := _m.Called(ctx, since)

	var r0 []common.Hash
	if rf, ok := ret.Get(0).(func(context.Context, time.Time) []common.Hash); ok {
		r0 = rf(ctx, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time) error); ok {
		r1 = rf(ctx, since)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPendingTxs provides a mock function with given fields: ctx, isClaims, limit
func (_m *poolMock) GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error) {
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

type mockConstructorTestingTnewPoolMock interface {
	mock.TestingT
	Cleanup(func())
}

// newPoolMock creates a new instance of poolMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newPoolMock(t mockConstructorTestingTnewPoolMock) *poolMock {
	mock := &poolMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
