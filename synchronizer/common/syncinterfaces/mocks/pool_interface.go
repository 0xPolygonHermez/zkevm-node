// Code generated by mockery. DO NOT EDIT.

package mock_syncinterfaces

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// PoolInterface is an autogenerated mock type for the PoolInterface type
type PoolInterface struct {
	mock.Mock
}

type PoolInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *PoolInterface) EXPECT() *PoolInterface_Expecter {
	return &PoolInterface_Expecter{mock: &_m.Mock}
}

// DeleteReorgedTransactions provides a mock function with given fields: ctx, txs
func (_m *PoolInterface) DeleteReorgedTransactions(ctx context.Context, txs []*types.Transaction) error {
	ret := _m.Called(ctx, txs)

	if len(ret) == 0 {
		panic("no return value specified for DeleteReorgedTransactions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*types.Transaction) error); ok {
		r0 = rf(ctx, txs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PoolInterface_DeleteReorgedTransactions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteReorgedTransactions'
type PoolInterface_DeleteReorgedTransactions_Call struct {
	*mock.Call
}

// DeleteReorgedTransactions is a helper method to define mock.On call
//   - ctx context.Context
//   - txs []*types.Transaction
func (_e *PoolInterface_Expecter) DeleteReorgedTransactions(ctx interface{}, txs interface{}) *PoolInterface_DeleteReorgedTransactions_Call {
	return &PoolInterface_DeleteReorgedTransactions_Call{Call: _e.mock.On("DeleteReorgedTransactions", ctx, txs)}
}

func (_c *PoolInterface_DeleteReorgedTransactions_Call) Run(run func(ctx context.Context, txs []*types.Transaction)) *PoolInterface_DeleteReorgedTransactions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*types.Transaction))
	})
	return _c
}

func (_c *PoolInterface_DeleteReorgedTransactions_Call) Return(_a0 error) *PoolInterface_DeleteReorgedTransactions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PoolInterface_DeleteReorgedTransactions_Call) RunAndReturn(run func(context.Context, []*types.Transaction) error) *PoolInterface_DeleteReorgedTransactions_Call {
	_c.Call.Return(run)
	return _c
}

// StoreTx provides a mock function with given fields: ctx, tx, ip, isWIP
func (_m *PoolInterface) StoreTx(ctx context.Context, tx types.Transaction, ip string, isWIP bool) error {
	ret := _m.Called(ctx, tx, ip, isWIP)

	if len(ret) == 0 {
		panic("no return value specified for StoreTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Transaction, string, bool) error); ok {
		r0 = rf(ctx, tx, ip, isWIP)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PoolInterface_StoreTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreTx'
type PoolInterface_StoreTx_Call struct {
	*mock.Call
}

// StoreTx is a helper method to define mock.On call
//   - ctx context.Context
//   - tx types.Transaction
//   - ip string
//   - isWIP bool
func (_e *PoolInterface_Expecter) StoreTx(ctx interface{}, tx interface{}, ip interface{}, isWIP interface{}) *PoolInterface_StoreTx_Call {
	return &PoolInterface_StoreTx_Call{Call: _e.mock.On("StoreTx", ctx, tx, ip, isWIP)}
}

func (_c *PoolInterface_StoreTx_Call) Run(run func(ctx context.Context, tx types.Transaction, ip string, isWIP bool)) *PoolInterface_StoreTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Transaction), args[2].(string), args[3].(bool))
	})
	return _c
}

func (_c *PoolInterface_StoreTx_Call) Return(_a0 error) *PoolInterface_StoreTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PoolInterface_StoreTx_Call) RunAndReturn(run func(context.Context, types.Transaction, string, bool) error) *PoolInterface_StoreTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewPoolInterface creates a new instance of PoolInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPoolInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *PoolInterface {
	mock := &PoolInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
