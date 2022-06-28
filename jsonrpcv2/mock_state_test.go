// Code generated by mockery v2.12.3. DO NOT EDIT.

package jsonrpcv2

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"

	runtime "github.com/hermeznetwork/hermez-core/statev2/runtime"

	statev2 "github.com/hermeznetwork/hermez-core/statev2"

	time "time"

	types "github.com/ethereum/go-ethereum/core/types"
)

// stateMock is an autogenerated mock type for the stateInterface type
type stateMock struct {
	mock.Mock
}

// BeginStateTransaction provides a mock function with given fields: ctx
func (_m *stateMock) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	ret := _m.Called(ctx)

	var r0 pgx.Tx
	if rf, ok := ret.Get(0).(func(context.Context) pgx.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Tx)
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

// CommitState provides a mock function with given fields: ctx, tx
func (_m *stateMock) CommitState(ctx context.Context, tx pgx.Tx) error {
	ret := _m.Called(ctx, tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DebugTransaction provides a mock function with given fields: ctx, transactionHash, tracer
func (_m *stateMock) DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error) {
	ret := _m.Called(ctx, transactionHash, tracer)

	var r0 *runtime.ExecutionResult
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, string) *runtime.ExecutionResult); ok {
		r0 = rf(ctx, transactionHash, tracer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.ExecutionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, string) error); ok {
		r1 = rf(ctx, transactionHash, tracer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EstimateGas provides a mock function with given fields: transaction, senderAddress
func (_m *stateMock) EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error) {
	ret := _m.Called(transaction, senderAddress)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(*types.Transaction, common.Address) uint64); ok {
		r0 = rf(transaction, senderAddress)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*types.Transaction, common.Address) error); ok {
		r1 = rf(transaction, senderAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalance provides a mock function with given fields: ctx, address, blockNumber, dbTx
func (_m *stateMock) GetBalance(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	ret := _m.Called(ctx, address, blockNumber, dbTx)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint64, pgx.Tx) *big.Int); ok {
		r0 = rf(ctx, address, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByHash provides a mock function with given fields: ctx, hash, dbTx
func (_m *stateMock) GetBlockByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*statev2.L2Block, error) {
	ret := _m.Called(ctx, hash, dbTx)

	var r0 *statev2.L2Block
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) *statev2.L2Block); ok {
		r0 = rf(ctx, hash, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*statev2.L2Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r1 = rf(ctx, hash, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockByNumber provides a mock function with given fields: ctx, blockNumber, dbTx
func (_m *stateMock) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*statev2.L2Block, error) {
	ret := _m.Called(ctx, blockNumber, dbTx)

	var r0 *statev2.L2Block
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) *statev2.L2Block); ok {
		r0 = rf(ctx, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*statev2.L2Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHashesSince provides a mock function with given fields: ctx, since, dbTx
func (_m *stateMock) GetBlockHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error) {
	ret := _m.Called(ctx, since, dbTx)

	var r0 []common.Hash
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, pgx.Tx) []common.Hash); ok {
		r0 = rf(ctx, since, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, pgx.Tx) error); ok {
		r1 = rf(ctx, since, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeader provides a mock function with given fields: ctx, blockNumber, dbTx
func (_m *stateMock) GetBlockHeader(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error) {
	ret := _m.Called(ctx, blockNumber, dbTx)

	var r0 *types.Header
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) *types.Header); ok {
		r0 = rf(ctx, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockTransactionCountByHash provides a mock function with given fields: ctx, hash, dbTx
func (_m *stateMock) GetBlockTransactionCountByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, hash, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) uint64); ok {
		r0 = rf(ctx, hash, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r1 = rf(ctx, hash, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockTransactionCountByNumber provides a mock function with given fields: ctx, blockNumber, dbTx
func (_m *stateMock) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, blockNumber, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) uint64); ok {
		r0 = rf(ctx, blockNumber, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCode provides a mock function with given fields: ctx, address, blockNumber, dbTx
func (_m *stateMock) GetCode(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	ret := _m.Called(ctx, address, blockNumber, dbTx)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint64, pgx.Tx) []byte); ok {
		r0 = rf(ctx, address, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBlock provides a mock function with given fields: ctx, dbTx
func (_m *stateMock) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*statev2.L2Block, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 *statev2.L2Block
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) *statev2.L2Block); ok {
		r0 = rf(ctx, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*statev2.L2Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBlockNumber provides a mock function with given fields: ctx, dbTx
func (_m *stateMock) GetLastBlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) uint64); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastConsolidatedBlockNumber provides a mock function with given fields: ctx, dbTx
func (_m *stateMock) GetLastConsolidatedBlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) uint64); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogs provides a mock function with given fields: ctx, fromBlock, toBlock, addresses, topics, blockHash, since, dbTx
func (_m *stateMock) GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error) {
	ret := _m.Called(ctx, fromBlock, toBlock, addresses, topics, blockHash, since, dbTx)

	var r0 []*types.Log
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, []common.Address, [][]common.Hash, *common.Hash, *time.Time, pgx.Tx) []*types.Log); ok {
		r0 = rf(ctx, fromBlock, toBlock, addresses, topics, blockHash, since, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Log)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64, []common.Address, [][]common.Hash, *common.Hash, *time.Time, pgx.Tx) error); ok {
		r1 = rf(ctx, fromBlock, toBlock, addresses, topics, blockHash, since, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonce provides a mock function with given fields: ctx, address, blockNumber, dbTx
func (_m *stateMock) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, address, blockNumber, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint64, pgx.Tx) uint64); ok {
		r0 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStorageAt provides a mock function with given fields: ctx, address, position, blockNumber, dbTx
func (_m *stateMock) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	ret := _m.Called(ctx, address, position, blockNumber, dbTx)

	var r0 *big.Int
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, uint64, pgx.Tx) *big.Int); ok {
		r0 = rf(ctx, address, position, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, address, position, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSyncingInfo provides a mock function with given fields: ctx, dbTx
func (_m *stateMock) GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (statev2.SyncingInfo, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 statev2.SyncingInfo
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) statev2.SyncingInfo); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(statev2.SyncingInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByBlockHashAndIndex provides a mock function with given fields: ctx, blockHash, index, dbTx
func (_m *stateMock) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	ret := _m.Called(ctx, blockHash, index, dbTx)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, uint64, pgx.Tx) *types.Transaction); ok {
		r0 = rf(ctx, blockHash, index, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, blockHash, index, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByBlockNumberAndIndex provides a mock function with given fields: ctx, blockNumber, index, dbTx
func (_m *stateMock) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	ret := _m.Called(ctx, blockNumber, index, dbTx)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, pgx.Tx) *types.Transaction); ok {
		r0 = rf(ctx, blockNumber, index, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, blockNumber, index, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByHash provides a mock function with given fields: ctx, transactionHash, dbTx
func (_m *stateMock) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	ret := _m.Called(ctx, transactionHash, dbTx)

	var r0 *types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) *types.Transaction); ok {
		r0 = rf(ctx, transactionHash, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r1 = rf(ctx, transactionHash, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionReceipt provides a mock function with given fields: ctx, transactionHash, dbTx
func (_m *stateMock) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error) {
	ret := _m.Called(ctx, transactionHash, dbTx)

	var r0 *types.Receipt
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) *types.Receipt); ok {
		r0 = rf(ctx, transactionHash, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Receipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r1 = rf(ctx, transactionHash, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessUnsignedTransaction provides a mock function with given fields: ctx, tx, senderAddress, sequencerAddress, blockNumber, dbTx
func (_m *stateMock) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, sequencerAddress common.Address, blockNumber uint64, dbTx pgx.Tx) *runtime.ExecutionResult {
	ret := _m.Called(ctx, tx, senderAddress, sequencerAddress, blockNumber, dbTx)

	var r0 *runtime.ExecutionResult
	if rf, ok := ret.Get(0).(func(context.Context, *types.Transaction, common.Address, common.Address, uint64, pgx.Tx) *runtime.ExecutionResult); ok {
		r0 = rf(ctx, tx, senderAddress, sequencerAddress, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.ExecutionResult)
		}
	}

	return r0
}

// RollbackState provides a mock function with given fields: ctx, tx
func (_m *stateMock) RollbackState(ctx context.Context, tx pgx.Tx) error {
	ret := _m.Called(ctx, tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type newStateMockT interface {
	mock.TestingT
	Cleanup(func())
}

// newStateMock creates a new instance of stateMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newStateMock(t newStateMockT) *stateMock {
	mock := &stateMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
