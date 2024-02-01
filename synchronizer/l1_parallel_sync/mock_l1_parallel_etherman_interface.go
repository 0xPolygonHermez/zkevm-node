// Code generated by mockery. DO NOT EDIT.

package l1_parallel_sync

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	etherman "github.com/0xPolygonHermez/zkevm-node/etherman"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// L1ParallelEthermanInterfaceMock is an autogenerated mock type for the L1ParallelEthermanInterface type
type L1ParallelEthermanInterfaceMock struct {
	mock.Mock
}

type L1ParallelEthermanInterfaceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *L1ParallelEthermanInterfaceMock) EXPECT() *L1ParallelEthermanInterfaceMock_Expecter {
	return &L1ParallelEthermanInterfaceMock_Expecter{mock: &_m.Mock}
}

// EthBlockByNumber provides a mock function with given fields: ctx, blockNumber
func (_m *L1ParallelEthermanInterfaceMock) EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	ret := _m.Called(ctx, blockNumber)

	var r0 *types.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*types.Block, error)); ok {
		return rf(ctx, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *types.Block); ok {
		r0 = rf(ctx, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EthBlockByNumber'
type L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call struct {
	*mock.Call
}

// EthBlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - blockNumber uint64
func (_e *L1ParallelEthermanInterfaceMock_Expecter) EthBlockByNumber(ctx interface{}, blockNumber interface{}) *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call {
	return &L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call{Call: _e.mock.On("EthBlockByNumber", ctx, blockNumber)}
}

func (_c *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call) Run(run func(ctx context.Context, blockNumber uint64)) *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call) Return(_a0 *types.Block, _a1 error) *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call) RunAndReturn(run func(context.Context, uint64) (*types.Block, error)) *L1ParallelEthermanInterfaceMock_EthBlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestBatchNumber provides a mock function with given fields:
func (_m *L1ParallelEthermanInterfaceMock) GetLatestBatchNumber() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestBatchNumber'
type L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call struct {
	*mock.Call
}

// GetLatestBatchNumber is a helper method to define mock.On call
func (_e *L1ParallelEthermanInterfaceMock_Expecter) GetLatestBatchNumber() *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call {
	return &L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call{Call: _e.mock.On("GetLatestBatchNumber")}
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call) Run(run func()) *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call) Return(_a0 uint64, _a1 error) *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call) RunAndReturn(run func() (uint64, error)) *L1ParallelEthermanInterfaceMock_GetLatestBatchNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestVerifiedBatchNum provides a mock function with given fields:
func (_m *L1ParallelEthermanInterfaceMock) GetLatestVerifiedBatchNum() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestVerifiedBatchNum'
type L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call struct {
	*mock.Call
}

// GetLatestVerifiedBatchNum is a helper method to define mock.On call
func (_e *L1ParallelEthermanInterfaceMock_Expecter) GetLatestVerifiedBatchNum() *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call {
	return &L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call{Call: _e.mock.On("GetLatestVerifiedBatchNum")}
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call) Run(run func()) *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call) Return(_a0 uint64, _a1 error) *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call) RunAndReturn(run func() (uint64, error)) *L1ParallelEthermanInterfaceMock_GetLatestVerifiedBatchNum_Call {
	_c.Call.Return(run)
	return _c
}

// GetRollupInfoByBlockRange provides a mock function with given fields: ctx, fromBlock, toBlock
func (_m *L1ParallelEthermanInterfaceMock) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error) {
	ret := _m.Called(ctx, fromBlock, toBlock)

	var r0 []etherman.Block
	var r1 map[common.Hash][]etherman.Order
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)); ok {
		return rf(ctx, fromBlock, toBlock)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, *uint64) []etherman.Block); ok {
		r0 = rf(ctx, fromBlock, toBlock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]etherman.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, *uint64) map[common.Hash][]etherman.Order); ok {
		r1 = rf(ctx, fromBlock, toBlock)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[common.Hash][]etherman.Order)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, uint64, *uint64) error); ok {
		r2 = rf(ctx, fromBlock, toBlock)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRollupInfoByBlockRange'
type L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call struct {
	*mock.Call
}

// GetRollupInfoByBlockRange is a helper method to define mock.On call
//   - ctx context.Context
//   - fromBlock uint64
//   - toBlock *uint64
func (_e *L1ParallelEthermanInterfaceMock_Expecter) GetRollupInfoByBlockRange(ctx interface{}, fromBlock interface{}, toBlock interface{}) *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call {
	return &L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call{Call: _e.mock.On("GetRollupInfoByBlockRange", ctx, fromBlock, toBlock)}
}

func (_c *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call) Run(run func(ctx context.Context, fromBlock uint64, toBlock *uint64)) *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(*uint64))
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call) Return(_a0 []etherman.Block, _a1 map[common.Hash][]etherman.Order, _a2 error) *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call) RunAndReturn(run func(context.Context, uint64, *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)) *L1ParallelEthermanInterfaceMock_GetRollupInfoByBlockRange_Call {
	_c.Call.Return(run)
	return _c
}

// GetTrustedSequencerURL provides a mock function with given fields:
func (_m *L1ParallelEthermanInterfaceMock) GetTrustedSequencerURL() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTrustedSequencerURL'
type L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call struct {
	*mock.Call
}

// GetTrustedSequencerURL is a helper method to define mock.On call
func (_e *L1ParallelEthermanInterfaceMock_Expecter) GetTrustedSequencerURL() *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call {
	return &L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call{Call: _e.mock.On("GetTrustedSequencerURL")}
}

func (_c *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call) Run(run func()) *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call) Return(_a0 string, _a1 error) *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call) RunAndReturn(run func() (string, error)) *L1ParallelEthermanInterfaceMock_GetTrustedSequencerURL_Call {
	_c.Call.Return(run)
	return _c
}

// HeaderByNumber provides a mock function with given fields: ctx, number
func (_m *L1ParallelEthermanInterfaceMock) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ret := _m.Called(ctx, number)

	var r0 *types.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*types.Header, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Header); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_HeaderByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HeaderByNumber'
type L1ParallelEthermanInterfaceMock_HeaderByNumber_Call struct {
	*mock.Call
}

// HeaderByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - number *big.Int
func (_e *L1ParallelEthermanInterfaceMock_Expecter) HeaderByNumber(ctx interface{}, number interface{}) *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call {
	return &L1ParallelEthermanInterfaceMock_HeaderByNumber_Call{Call: _e.mock.On("HeaderByNumber", ctx, number)}
}

func (_c *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call) Run(run func(ctx context.Context, number *big.Int)) *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*big.Int))
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call) Return(_a0 *types.Header, _a1 error) *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call) RunAndReturn(run func(context.Context, *big.Int) (*types.Header, error)) *L1ParallelEthermanInterfaceMock_HeaderByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyGenBlockNumber provides a mock function with given fields: ctx, genBlockNumber
func (_m *L1ParallelEthermanInterfaceMock) VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error) {
	ret := _m.Called(ctx, genBlockNumber)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (bool, error)); ok {
		return rf(ctx, genBlockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) bool); ok {
		r0 = rf(ctx, genBlockNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, genBlockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyGenBlockNumber'
type L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call struct {
	*mock.Call
}

// VerifyGenBlockNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - genBlockNumber uint64
func (_e *L1ParallelEthermanInterfaceMock_Expecter) VerifyGenBlockNumber(ctx interface{}, genBlockNumber interface{}) *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call {
	return &L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call{Call: _e.mock.On("VerifyGenBlockNumber", ctx, genBlockNumber)}
}

func (_c *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call) Run(run func(ctx context.Context, genBlockNumber uint64)) *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call) Return(_a0 bool, _a1 error) *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call) RunAndReturn(run func(context.Context, uint64) (bool, error)) *L1ParallelEthermanInterfaceMock_VerifyGenBlockNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewL1ParallelEthermanInterfaceMock creates a new instance of L1ParallelEthermanInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewL1ParallelEthermanInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *L1ParallelEthermanInterfaceMock {
	mock := &L1ParallelEthermanInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
