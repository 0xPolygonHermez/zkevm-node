// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EthTxManager is an autogenerated mock type for the ethTxManager type
type EthTxManager struct {
	mock.Mock
}

type mockConstructorTestingTNewEthTxManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewEthTxManager creates a new instance of EthTxManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEthTxManager(t mockConstructorTestingTNewEthTxManager) *EthTxManager {
	mock := &EthTxManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
