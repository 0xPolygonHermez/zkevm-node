// Code generated by mockery v2.12.3. DO NOT EDIT.

package profitabilitychecker_test

import (
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// ethermanMock is an autogenerated mock type for the etherman type
type ethermanMock struct {
	mock.Mock
}

// GetSendSequenceFee provides a mock function with given fields:
func (_m *ethermanMock) GetSendSequenceFee() (*big.Int, error) {
	ret := _m.Called()

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
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

type newEthermanMockT interface {
	mock.TestingT
	Cleanup(func())
}

// newEthermanMock creates a new instance of ethermanMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newEthermanMock(t newEthermanMockT) *ethermanMock {
	mock := &ethermanMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
