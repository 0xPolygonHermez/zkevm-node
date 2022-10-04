// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"
)

// Etherman is an autogenerated mock type for the etherman type
type Etherman struct {
	mock.Mock
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
func (_m *Etherman) GetPublicAddress() common.Address {
	ret := _m.Called()

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
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
