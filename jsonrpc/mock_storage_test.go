// Code generated by mockery v2.14.0. DO NOT EDIT.

package jsonrpc

import mock "github.com/stretchr/testify/mock"

// storageMock is an autogenerated mock type for the storageInterface type
type storageMock struct {
	mock.Mock
}

// GetFilter provides a mock function with given fields: filterID
func (_m *storageMock) GetFilter(filterID uint64) (*Filter, error) {
	ret := _m.Called(filterID)

	var r0 *Filter
	if rf, ok := ret.Get(0).(func(uint64) *Filter); ok {
		r0 = rf(filterID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Filter)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(filterID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBlockFilter provides a mock function with given fields:
func (_m *storageMock) NewBlockFilter() (uint64, error) {
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

// NewLogFilter provides a mock function with given fields: filter
func (_m *storageMock) NewLogFilter(filter LogFilter) (uint64, error) {
	ret := _m.Called(filter)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(LogFilter) uint64); ok {
		r0 = rf(filter)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(LogFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPendingTransactionFilter provides a mock function with given fields:
func (_m *storageMock) NewPendingTransactionFilter() (uint64, error) {
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

// UninstallFilter provides a mock function with given fields: filterID
func (_m *storageMock) UninstallFilter(filterID uint64) (bool, error) {
	ret := _m.Called(filterID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64) bool); ok {
		r0 = rf(filterID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(filterID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFilterLastPoll provides a mock function with given fields: filterID
func (_m *storageMock) UpdateFilterLastPoll(filterID uint64) error {
	ret := _m.Called(filterID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(filterID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewStorageMock interface {
	mock.TestingT
	Cleanup(func())
}

// newStorageMock creates a new instance of storageMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newStorageMock(t mockConstructorTestingTnewStorageMock) *storageMock {
	mock := &storageMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
