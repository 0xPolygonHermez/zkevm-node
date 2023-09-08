// Code generated by mockery v2.32.0. DO NOT EDIT.

package synchronizer

import (
	state "github.com/0xPolygonHermez/zkevm-node/state"
	mock "github.com/stretchr/testify/mock"
)

// l1RollupConsumerInterfaceMock is an autogenerated mock type for the l1RollupConsumerInterfaceMock type
type l1RollupConsumerInterfaceMock struct {
	mock.Mock
}

// getLastEthBlockSynced provides a mock function with given fields:
func (_m *l1RollupConsumerInterfaceMock) getLastEthBlockSynced() (state.Block, bool) {
	ret := _m.Called()

	var r0 state.Block
	var r1 bool
	if rf, ok := ret.Get(0).(func() (state.Block, bool)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() state.Block); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(state.Block)
	}

	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// start provides a mock function with given fields:
func (_m *l1RollupConsumerInterfaceMock) start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// stopAfterProcessChannelQueue provides a mock function with given fields:
func (_m *l1RollupConsumerInterfaceMock) stopAfterProcessChannelQueue() {
	_m.Called()
}

// newL1RollupConsumerInterfaceMock creates a new instance of l1RollupConsumerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newL1RollupConsumerInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *l1RollupConsumerInterfaceMock {
	mock := &l1RollupConsumerInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
