// Code generated by mockery v2.39.0. DO NOT EDIT.

package l1_parallel_sync

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// l1RollupProducerInterfaceMock is an autogenerated mock type for the l1RollupProducerInterface type
type l1RollupProducerInterfaceMock struct {
	mock.Mock
}

// Abort provides a mock function with given fields:
func (_m *l1RollupProducerInterfaceMock) Abort() {
	_m.Called()
}

// Reset provides a mock function with given fields: startingBlockNumber
func (_m *l1RollupProducerInterfaceMock) Reset(startingBlockNumber uint64) {
	_m.Called(startingBlockNumber)
}

// Start provides a mock function with given fields: ctx
func (_m *l1RollupProducerInterfaceMock) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *l1RollupProducerInterfaceMock) Stop() {
	_m.Called()
}

// newL1RollupProducerInterfaceMock creates a new instance of l1RollupProducerInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newL1RollupProducerInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *l1RollupProducerInterfaceMock {
	mock := &l1RollupProducerInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
