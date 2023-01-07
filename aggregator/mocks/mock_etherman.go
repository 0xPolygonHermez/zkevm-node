// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/0xPolygonHermez/zkevm-node/etherman/types"
)

// Etherman is an autogenerated mock type for the etherman type
type Etherman struct {
	mock.Mock
}

// EstimateGasForTrustedVerifyBatches provides a mock function with given fields: lastVerifiedBatch, newVerifiedBatch, inputs
func (_m *Etherman) EstimateGasForTrustedVerifyBatches(lastVerifiedBatch uint64, newVerifiedBatch uint64, inputs *types.FinalProofInputs) (*coretypes.Transaction, error) {
	ret := _m.Called(lastVerifiedBatch, newVerifiedBatch, inputs)

	var r0 *coretypes.Transaction
	if rf, ok := ret.Get(0).(func(uint64, uint64, *types.FinalProofInputs) *coretypes.Transaction); ok {
		r0 = rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64, *types.FinalProofInputs) error); ok {
		r1 = rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestVerifiedBatchNum provides a mock function with given fields:
func (_m *Etherman) GetLatestVerifiedBatchNum() (uint64, error) {
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

// GetPublicAddress provides a mock function with given fields:
func (_m *Etherman) GetPublicAddress() (common.Address, error) {
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

type mockConstructorTestingTNewEtherman interface {
	mock.TestingT
	Cleanup(func())
}

// NewEtherman creates a new instance of Etherman. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEtherman(t mockConstructorTestingTNewEtherman) *Etherman {
	mock := &Etherman{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
