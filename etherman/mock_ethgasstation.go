// Code generated by mockery v2.16.0. DO NOT EDIT.

package etherman

import (
	context "context"
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// ethGasStationMock is an autogenerated mock type for the GasPricer type
type ethGasStationMock struct {
	mock.Mock
}

// SuggestGasPrice provides a mock function with given fields: ctx
func (_m *ethGasStationMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewEthGasStationMock interface {
	mock.TestingT
	Cleanup(func())
}

// newEthGasStationMock creates a new instance of ethGasStationMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newEthGasStationMock(t mockConstructorTestingTnewEthGasStationMock) *ethGasStationMock {
	mock := &ethGasStationMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
