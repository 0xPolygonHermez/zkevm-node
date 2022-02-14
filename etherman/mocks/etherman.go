// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	etherman "github.com/hermeznetwork/hermez-core/etherman"

	mock "github.com/stretchr/testify/mock"

	proverclient "github.com/hermeznetwork/hermez-core/proverclient"

	state "github.com/hermeznetwork/hermez-core/state"

	types "github.com/ethereum/go-ethereum/core/types"
)

// EtherMan is an autogenerated mock type for the EtherMan type
type EtherMan struct {
	mock.Mock
}

// ConsolidateBatch provides a mock function with given fields: batchNum, proof
func (_m *EtherMan) ConsolidateBatch(batchNum *big.Int, proof *proverclient.Proof) (*types.Transaction, error) {
	ret := _m.Called(batchNum, proof)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(*big.Int, *proverclient.Proof) *types.Transaction); ok {
		r0 = rf(batchNum, proof)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*big.Int, *proverclient.Proof) error); ok {
		r1 = rf(batchNum, proof)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EstimateSendBatchCost provides a mock function with given fields: ctx, txs, maticAmount
func (_m *EtherMan) EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, txs, maticAmount)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(context.Context, []*types.Transaction, *big.Int) *big.Int); ok {
		r0 = rf(ctx, txs, maticAmount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*types.Transaction, *big.Int) error); ok {
		r1 = rf(ctx, txs, maticAmount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EthBlockByNumber provides a mock function with given fields: ctx, blockNum
func (_m *EtherMan) EthBlockByNumber(ctx context.Context, blockNum uint64) (*types.Block, error) {
	ret := _m.Called(ctx, blockNum)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *types.Block); ok {
		r0 = rf(ctx, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAddress provides a mock function with given fields:
func (_m *EtherMan) GetAddress() common.Address {
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

// GetCustomChainID provides a mock function with given fields:
func (_m *EtherMan) GetCustomChainID() (*big.Int, error) {
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

// GetDefaultChainID provides a mock function with given fields:
func (_m *EtherMan) GetDefaultChainID() (*big.Int, error) {
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

// GetLatestConsolidatedBatchNumber provides a mock function with given fields:
func (_m *EtherMan) GetLatestConsolidatedBatchNumber() (uint64, error) {
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

// GetLatestProposedBatchNumber provides a mock function with given fields:
func (_m *EtherMan) GetLatestProposedBatchNumber() (uint64, error) {
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

// GetRollupInfoByBlock provides a mock function with given fields: ctx, blockNum, blockHash
func (_m *EtherMan) GetRollupInfoByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, map[common.Hash][]etherman.Order, error) {
	ret := _m.Called(ctx, blockNum, blockHash)

	var r0 []state.Block
	if rf, ok := ret.Get(0).(func(context.Context, uint64, *common.Hash) []state.Block); ok {
		r0 = rf(ctx, blockNum, blockHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]state.Block)
		}
	}

	var r1 map[common.Hash][]etherman.Order
	if rf, ok := ret.Get(1).(func(context.Context, uint64, *common.Hash) map[common.Hash][]etherman.Order); ok {
		r1 = rf(ctx, blockNum, blockHash)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[common.Hash][]etherman.Order)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, uint64, *common.Hash) error); ok {
		r2 = rf(ctx, blockNum, blockHash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetRollupInfoByBlockRange provides a mock function with given fields: ctx, fromBlock, toBlock
func (_m *EtherMan) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, map[common.Hash][]etherman.Order, error) {
	ret := _m.Called(ctx, fromBlock, toBlock)

	var r0 []state.Block
	if rf, ok := ret.Get(0).(func(context.Context, uint64, *uint64) []state.Block); ok {
		r0 = rf(ctx, fromBlock, toBlock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]state.Block)
		}
	}

	var r1 map[common.Hash][]etherman.Order
	if rf, ok := ret.Get(1).(func(context.Context, uint64, *uint64) map[common.Hash][]etherman.Order); ok {
		r1 = rf(ctx, fromBlock, toBlock)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[common.Hash][]etherman.Order)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, uint64, *uint64) error); ok {
		r2 = rf(ctx, fromBlock, toBlock)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetSequencerCollateral provides a mock function with given fields: batchNumber
func (_m *EtherMan) GetSequencerCollateralByBatchNumber(batchNumber uint64) (*big.Int, error) {
	ret := _m.Called(batchNumber)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(uint64) *big.Int); ok {
		r0 = rf(batchNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(batchNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HeaderByNumber provides a mock function with given fields: ctx, number
func (_m *EtherMan) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ret := _m.Called(ctx, number)

	var r0 *types.Header
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Header); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterSequencer provides a mock function with given fields: url
func (_m *EtherMan) RegisterSequencer(url string) (*types.Transaction, error) {
	ret := _m.Called(url)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(string) *types.Transaction); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendBatch provides a mock function with given fields: ctx, txs, maticAmount
func (_m *EtherMan) SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error) {
	ret := _m.Called(ctx, txs, maticAmount)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, []*types.Transaction, *big.Int) *types.Transaction); ok {
		r0 = rf(ctx, txs, maticAmount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*types.Transaction, *big.Int) error); ok {
		r1 = rf(ctx, txs, maticAmount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentSequencerCollateral provides a mock function to get the current collateral
func (_m *EtherMan) GetCurrentSequencerCollateral() (*big.Int, error) {
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