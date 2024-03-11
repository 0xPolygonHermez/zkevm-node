package actions_test

import (
	"context"
	"math/big"
	"testing"

	rpctypes "github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type CheckL2BlocksTestData struct {
	sut         *actions.CheckL2BlockHash
	mockState   *mock_syncinterfaces.StateFullInterface
	zKEVMClient *mock_syncinterfaces.ZKEVMClientInterface
}

func TestCheckL2BlockHash_GetMinimumL2BlockToCheck(t *testing.T) {
	// Create an instance of CheckL2BlockHash
	values := []struct {
		initial  uint64
		modulus  uint64
		expected uint64
	}{
		{0, 10, 10},
		{1, 10, 10},
		{9, 10, 10},
		{10, 10, 20},
		{0, 0, 1},
		{1, 0, 2},
	}
	for _, data := range values {
		// Call the GetNextL2BlockToCheck method
		checkL2Block := actions.NewCheckL2BlockHash(nil, nil, data.initial, data.modulus)
		nextL2Block := checkL2Block.GetMinimumL2BlockToCheck()

		// Assert the expected result
		assert.Equal(t, data.expected, nextL2Block)
	}
}

func TestCheckL2BlockHashNotEnoughBlocksToCheck(t *testing.T) {
	data := newCheckL2BlocksTestData(t, 0, 10)
	// Call the CheckL2Block method
	data.mockState.EXPECT().GetLastL2BlockNumber(mock.Anything, mock.Anything).Return(uint64(0), nil)
	err := data.sut.CheckL2Block(context.Background(), nil)
	require.NoError(t, err)
}

func newCheckL2BlocksTestData(t *testing.T, initialL2Block, modulus uint64) CheckL2BlocksTestData {
	res := CheckL2BlocksTestData{
		mockState:   mock_syncinterfaces.NewStateFullInterface(t),
		zKEVMClient: mock_syncinterfaces.NewZKEVMClientInterface(t),
	}
	res.sut = actions.NewCheckL2BlockHash(res.mockState, res.zKEVMClient, initialL2Block, modulus)
	return res
}
func TestCheckL2BlockHash_GetNextL2BlockToCheck(t *testing.T) {
	values := []struct {
		lastLocalL2BlockNumber    uint64
		minL2BlockNumberToCheck   uint64
		expectedShouldCheck       bool
		expectedNextL2BlockNumber uint64
	}{
		{0, 10, false, 0},
		{10, 10, true, 10},
		{9, 10, false, 0},
		{10, 10, true, 10},
		{0, 0, true, 0},
		{1, 0, true, 1},
	}

	for _, data := range values {
		checkL2Block := actions.NewCheckL2BlockHash(nil, nil, 0, 0)
		shouldCheck, nextL2Block := checkL2Block.GetNextL2BlockToCheck(data.lastLocalL2BlockNumber, data.minL2BlockNumberToCheck)

		assert.Equal(t, data.expectedShouldCheck, shouldCheck, data)
		assert.Equal(t, data.expectedNextL2BlockNumber, nextL2Block, data)
	}
}

func TestCheckL2BlockHashMatch(t *testing.T) {
	data := newCheckL2BlocksTestData(t, 1, 10)
	lastL2Block := uint64(14)
	lastL2BlockBigInt := big.NewInt(int64(lastL2Block))
	gethHeader := types.Header{
		Number: big.NewInt(int64(lastL2Block)),
	}
	stateBlock := state.NewL2Block(state.NewL2Header(&gethHeader), nil, nil, nil, nil)

	data.mockState.EXPECT().GetLastL2BlockNumber(mock.Anything, mock.Anything).Return(lastL2Block, nil)
	data.mockState.EXPECT().GetL2BlockByNumber(mock.Anything, lastL2Block, mock.Anything).Return(stateBlock, nil)
	l2blockHash := stateBlock.Hash()
	rpcL2Block := rpctypes.Block{
		Hash:   &l2blockHash,
		Number: rpctypes.ArgUint64(lastL2Block),
	}

	data.zKEVMClient.EXPECT().BlockByNumber(mock.Anything, lastL2BlockBigInt).Return(&rpcL2Block, nil)
	err := data.sut.CheckL2Block(context.Background(), nil)
	require.NoError(t, err)
}

func TestCheckL2BlockHashMissmatch(t *testing.T) {
	data := newCheckL2BlocksTestData(t, 1, 10)
	lastL2Block := uint64(14)
	lastL2BlockBigInt := big.NewInt(int64(lastL2Block))
	gethHeader := types.Header{
		Number: big.NewInt(int64(lastL2Block)),
	}
	stateBlock := state.NewL2Block(state.NewL2Header(&gethHeader), nil, nil, nil, nil)

	data.mockState.EXPECT().GetLastL2BlockNumber(mock.Anything, mock.Anything).Return(lastL2Block, nil)
	data.mockState.EXPECT().GetL2BlockByNumber(mock.Anything, lastL2Block, mock.Anything).Return(stateBlock, nil)
	l2blockHash := common.HexToHash("0x1234")
	rpcL2Block := rpctypes.Block{
		Hash:   &l2blockHash,
		Number: rpctypes.ArgUint64(lastL2Block),
	}

	data.zKEVMClient.EXPECT().BlockByNumber(mock.Anything, lastL2BlockBigInt).Return(&rpcL2Block, nil)
	err := data.sut.CheckL2Block(context.Background(), nil)
	require.Error(t, err)
}
