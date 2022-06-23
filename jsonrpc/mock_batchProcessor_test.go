// Code generated by mockery v2.12.3. DO NOT EDIT.

package jsonrpc

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	runtime "github.com/hermeznetwork/hermez-core/state/runtime"

	types "github.com/ethereum/go-ethereum/core/types"
)

// batchProcessorMock is an autogenerated mock type for the BatchProcessorInterface type
type batchProcessorMock struct {
	mock.Mock
}

// ProcessUnsignedTransaction provides a mock function with given fields: ctx, tx, senderAddress, sequencerAddress
func (_m *batchProcessorMock) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, sequencerAddress common.Address) *runtime.ExecutionResult {
	ret := _m.Called(ctx, tx, senderAddress, sequencerAddress)

	var r0 *runtime.ExecutionResult
	if rf, ok := ret.Get(0).(func(context.Context, *types.Transaction, common.Address, common.Address) *runtime.ExecutionResult); ok {
		r0 = rf(ctx, tx, senderAddress, sequencerAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.ExecutionResult)
		}
	}

	return r0
}

type newBatchProcessorMockT interface {
	mock.TestingT
	Cleanup(func())
}

// newBatchProcessorMock creates a new instance of batchProcessorMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newBatchProcessorMock(t newBatchProcessorMockT) *batchProcessorMock {
	mock := &batchProcessorMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
