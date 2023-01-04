// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	coretypes "github.com/ethereum/go-ethereum/core/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/0xPolygonHermez/zkevm-node/etherman/types"
)

// EthermanMock is an autogenerated mock type for the etherman type
type EthermanMock struct {
	mock.Mock
}

// EstimateGasSequenceBatches provides a mock function with given fields: sequences
func (_m *EthermanMock) EstimateGasSequenceBatches(sequences []types.Sequence) (*coretypes.Transaction, error) {
	ret := _m.Called(sequences)

	var r0 *coretypes.Transaction
	if rf, ok := ret.Get(0).(func([]types.Sequence) *coretypes.Transaction); ok {
		r0 = rf(sequences)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]types.Sequence) error); ok {
		r1 = rf(sequences)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBatchTimestamp provides a mock function with given fields:
func (_m *EthermanMock) GetLastBatchTimestamp() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBatchNumber provides a mock function with given fields:
func (_m *EthermanMock) GetLatestBatchNumber() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlockNumber provides a mock function with given fields: ctx
func (_m *EthermanMock) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
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

// GetLatestBlockTimestamp provides a mock function with given fields: ctx
func (_m *EthermanMock) GetLatestBlockTimestamp(ctx context.Context) (uint64, error) {
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

// GetSendSequenceFee provides a mock function with given fields: numBatches
func (_m *EthermanMock) GetSendSequenceFee(numBatches uint64) (*big.Int, error) {
	ret := _m.Called(numBatches)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(uint64) *big.Int); ok {
		r0 = rf(numBatches)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(numBatches)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TrustedSequencer provides a mock function with given fields:
func (_m *EthermanMock) TrustedSequencer() (common.Address, error) {
	ret := _m.Called()

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEthermanMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewEthermanMock creates a new instance of EthermanMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEthermanMock(t mockConstructorTestingTNewEthermanMock) *EthermanMock {
	mock := &EthermanMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
