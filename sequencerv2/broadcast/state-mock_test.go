// Code generated by mockery v2.12.3. DO NOT EDIT.

package broadcast_test

import (
	context "context"

	broadcast "github.com/hermeznetwork/hermez-core/sequencerv2/broadcast"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"
)

// stateMock is an autogenerated mock type for the stateInterface type
type stateMock struct {
	mock.Mock
}

// GetBatchByNumber provides a mock function with given fields: ctx, batchNumber, tx
func (_m *stateMock) GetBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*broadcast.Batch, error) {
	ret := _m.Called(ctx, batchNumber, tx)

	var r0 *broadcast.Batch
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) *broadcast.Batch); ok {
		r0 = rf(ctx, batchNumber, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*broadcast.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEncodedTransactionsByBatchNumber provides a mock function with given fields: ctx, batchNumber, tx
func (_m *stateMock) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) ([]string, error) {
	ret := _m.Called(ctx, batchNumber, tx)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) []string); ok {
		r0 = rf(ctx, batchNumber, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBatch provides a mock function with given fields: ctx, tx
func (_m *stateMock) GetLastBatch(ctx context.Context, tx pgx.Tx) (*broadcast.Batch, error) {
	ret := _m.Called(ctx, tx)

	var r0 *broadcast.Batch
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) *broadcast.Batch); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*broadcast.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type newStateMockT interface {
	mock.TestingT
	Cleanup(func())
}

// newStateMock creates a new instance of stateMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newStateMock(t newStateMockT) *stateMock {
	mock := &stateMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
