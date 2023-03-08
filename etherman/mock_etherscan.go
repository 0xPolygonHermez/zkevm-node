// Code generated by mockery v2.15.0. DO NOT EDIT.

package etherman

import (
	context "context"
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// etherscanMock is an autogenerated mock type for the GasPricer type
type etherscanMock struct {
	mock.Mock
}

// SuggestGasPrice provides a mock function with given fields: ctx
func (_m *etherscanMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
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

type mockConstructorTestingTnewEtherscanMock interface {
	mock.TestingT
	Cleanup(func())
}

// newEtherscanMock creates a new instance of etherscanMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newEtherscanMock(t mockConstructorTestingTnewEtherscanMock) *etherscanMock {
	mock := &etherscanMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
