package jsonrpc_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlockNumber(t *testing.T) {
	const expectedBlockNumber = uint64(10)

	server, mocks, ethClient := newMockedServer(t)
	defer server.Server.Stop()

	mocks.State.
		On("GetLastBatchNumber", context.Background(), "").
		Return(expectedBlockNumber, nil)

	bn, err := ethClient.BlockNumber(context.Background())
	require.NoError(t, err)

	assert.Equal(t, expectedBlockNumber, bn)
}

func TestCall(t *testing.T) {
	server, mocks, ethClient := newMockedServer(t)
	defer server.Server.Stop()

	testCases := []struct {
		name     string
		from     common.Address
		to       *common.Address
		gas      uint64
		gasPrice *big.Int
		value    *big.Int
		data     []byte

		expectedResult *[]byte
	}{
		{
			name:           "Transaction with all information",
			from:           common.HexToAddress("0x1"),
			to:             ptr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: ptr([]byte("hello world")),
		},
		{
			name:           "Transaction without from and gas",
			to:             ptr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: ptr([]byte("hello world")),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			msg := ethereum.CallMsg{From: testCase.from, To: testCase.to, Gas: testCase.gas, GasPrice: testCase.gasPrice, Value: testCase.value, Data: testCase.data}

			batchNumber := uint64(1)
			txBundleID := ""
			batch := &state.Batch{
				Header: &types.Header{
					Root:     common.Hash{},
					GasLimit: 123456,
				},
			}

			mocks.State.
				On("GetLastBatchNumber", context.Background(), txBundleID).
				Return(batchNumber, nil)

			mocks.State.
				On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).
				Return(batch, nil)

			if testCase.gas == 0 {
				mocks.State.
					On("GetLastBatch", context.Background(), false, txBundleID).
					Return(batch, nil)
			}

			if testCase.expectedResult != nil {
				mocks.State.
					On("NewBatchProcessor", context.Background(), server.SequencerAddress, batch.Header.Root[:], txBundleID).
					Return(mocks.BatchProcessor, nil)

				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					if tx == nil {
						return false
					}

					if testCase.gas == 0 && tx.Gas() != batch.Header.GasLimit {
						return false
					}

					return tx.To().Hex() == testCase.to.Hex() &&
						tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
						tx.Value().Uint64() == testCase.value.Uint64() &&
						hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
				})

				mocks.BatchProcessor.
					On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, server.SequencerAddress).
					Return(&runtime.ExecutionResult{ReturnValue: *testCase.expectedResult})
			}

			result, err := ethClient.CallContract(context.Background(), msg, nil)
			require.NoError(t, err)

			if testCase.expectedResult != nil {
				assert.Equal(t, *testCase.expectedResult, result)
			}
		})
	}

}

func TestChainID(t *testing.T) {
	server, _, ethClient := newMockedServer(t)
	defer server.Server.Stop()

	chainID, err := ethClient.ChainID(context.Background())
	require.NoError(t, err)

	assert.Equal(t, server.ChainID, chainID.Uint64())
}

// ptr returns a pointer to the provided value
func ptr[T interface{}](i T) *T {
	return &i
}

// val return the value from the provided pointer
func val[T interface{}](i *T) T {
	return *i
}
