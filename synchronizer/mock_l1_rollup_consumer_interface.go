// Code generated by mockery v2.32.0. DO NOT EDIT.

package synchronizer

import (
	context "context"

	state "github.com/0xPolygonHermez/zkevm-node/state"
	mock "github.com/stretchr/testify/mock"
)

// l1RollupConsumerInterfaceMock is an autogenerated mock type for the l1RollupConsumerInterface type
type l1RollupConsumerInterfaceMock struct {
	mock.Mock
}

// GetLastEthBlockSynced provides a mock function with given fields:
func (_m *l1RollupConsumerInterfaceMock) GetLastEthBlockSynced() (state.Block, bool) {
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

// Reset provides a mock function with given fields: startingBlockNumber
func (_m *l1RollupConsumerInterfaceMock) Reset(startingBlockNumber uint64) {
	_m.Called(startingBlockNumber)
}

// Start provides a mock function with given fields: ctx, lastEthBlockSynced
func (_m *l1RollupConsumerInterfaceMock) Start(ctx context.Context, lastEthBlockSynced *state.Block) error {
	ret := _m.Called(ctx, lastEthBlockSynced)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *state.Block) error); ok {
		r0 = rf(ctx, lastEthBlockSynced)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StopAfterProcessChannelQueue provides a mock function with given fields:
func (_m *l1RollupConsumerInterfaceMock) StopAfterProcessChannelQueue() {
	_m.Called()
}

// newL1RollupConsumerInterfaceMock creates a new instance of l1RollupConsumerInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
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
