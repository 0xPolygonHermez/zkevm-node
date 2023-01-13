package ethtxmanager

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var defaultEthTxmanagerConfigForTests = Config{
	FrequencyToMonitorTxs: types.NewDuration(time.Millisecond),
	WaitTxToBeMined:       types.NewDuration(time.Second),
}

func TestTxGetMined(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	st := newStateMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(defaultEthTxmanagerConfigForTests, etherman, storage, st)

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	ctx := context.Background()

	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx, from).
		Return(currentNonce, nil).
		Once()

	estimatedGas := uint64(1)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(estimatedGas, nil).
		Once()

	suggestedGasPrice := big.NewInt(1)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(suggestedGasPrice, nil).
		Once()

	signedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      estimatedGas,
		GasPrice: suggestedGasPrice,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(signedTx, nil).
		Once()

	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(signedTx, false, nil).
		Once()

	etherman.
		On("SendTx", ctx, signedTx).
		Return(nil).
		Once()

	etherman.
		On("WaitTxToBeMined", ctx, signedTx, mock.IsType(time.Second)).
		Return(true, nil).
		Once()

	blockNumber := big.NewInt(1)

	receipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
		Status:      ethTypes.ReceiptStatusSuccessful,
	}
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(receipt, nil).
		Once()

	etherman.
		On("GetRevertMessage", ctx, signedTx).
		Return("", nil).
		Once()

	block := &state.Block{
		BlockNumber: blockNumber.Uint64(),
	}
	st.
		On("GetLastBlock", ctx, nil).
		Return(block, nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.Start()

	time.Sleep(time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
	require.Equal(t, 1, len(result.Txs))
	require.Equal(t, signedTx, result.Txs[signedTx.Hash()].Tx)
	require.Equal(t, receipt, result.Txs[signedTx.Hash()].Receipt)
	require.Equal(t, "", result.Txs[signedTx.Hash()].RevertMessage)
}

func TestTxGetMinedAfterReviewed(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	st := newStateMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(defaultEthTxmanagerConfigForTests, etherman, storage, st)

	ctx := context.Background()

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	// Add
	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx, from).
		Return(currentNonce, nil).
		Once()

	firstGasEstimation := uint64(1)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(firstGasEstimation, nil).
		Once()

	firstGasPriceSuggestion := big.NewInt(1)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(firstGasPriceSuggestion, nil).
		Once()

	// Monitoring Cycle 1
	firstSignedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      firstGasEstimation,
		GasPrice: firstGasPriceSuggestion,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(firstSignedTx, nil).
		Once()
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("SendTx", ctx, firstSignedTx).
		Return(nil).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, firstSignedTx, mock.IsType(time.Second)).
		Return(false, errors.New("tx not mined yet")).
		Once()

	// Monitoring Cycle 2
	etherman.
		On("CheckTxWasMined", ctx, firstSignedTx.Hash()).
		Return(false, nil, nil).
		Once()

	secondGasEstimation := uint64(2)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(secondGasEstimation, nil).
		Once()
	secondGasPriceSuggestion := big.NewInt(2)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(secondGasPriceSuggestion, nil).
		Once()

	secondSignedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      secondGasEstimation,
		GasPrice: secondGasPriceSuggestion,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(secondSignedTx, nil).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("SendTx", ctx, secondSignedTx).
		Return(nil).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, secondSignedTx, mock.IsType(time.Second)).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(true, nil).
		Once()

	blockNumber := big.NewInt(1)

	receipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
		Status:      ethTypes.ReceiptStatusSuccessful,
	}
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()

	block := &state.Block{
		BlockNumber: blockNumber.Uint64(),
	}
	st.
		On("GetLastBlock", ctx, nil).
		Return(block, nil).
		Once()

	// Build result
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(firstSignedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, firstSignedTx.Hash()).
		Return(nil, ethereum.NotFound).
		Once()
	etherman.
		On("GetRevertMessage", ctx, firstSignedTx).
		Return("", nil).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(secondSignedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetRevertMessage", ctx, secondSignedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.Start()

	time.Sleep(time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
}

func TestTxGetMinedAfterConfirmedAndReorged(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	st := newStateMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(defaultEthTxmanagerConfigForTests, etherman, storage, st)

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	ctx := context.Background()

	// Add
	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx, from).
		Return(currentNonce, nil).
		Once()

	estimatedGas := uint64(1)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(estimatedGas, nil).
		Once()

	suggestedGasPrice := big.NewInt(1)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(suggestedGasPrice, nil).
		Once()

	// Monitoring Cycle 1
	signedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      estimatedGas,
		GasPrice: suggestedGasPrice,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(signedTx, nil).
		Once()
	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("SendTx", ctx, signedTx).
		Return(nil).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, signedTx, mock.IsType(time.Second)).
		Return(true, nil).
		Once()

	blockNumber := big.NewInt(10)

	receipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
		Status:      ethTypes.ReceiptStatusSuccessful,
	}
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(receipt, nil).
		Once()

	block := &state.Block{
		BlockNumber: blockNumber.Uint64(),
	}
	st.
		On("GetLastBlock", ctx, nil).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(block, nil).
		Once()

	// Build Result 1
	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(signedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetRevertMessage", ctx, signedTx).
		Return("", nil).
		Once()

	// Build Result 2
	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(nil, ethereum.NotFound).
		Once()
	var nilTxPtr *ethTypes.Transaction
	etherman.
		On("GetRevertMessage", ctx, nilTxPtr).
		Return("", nil).
		Once()

	// Monitoring Cycle 2
	etherman.
		On("CheckTxWasMined", ctx, signedTx.Hash()).
		Return(false, nil, nil).
		Once()

	// Monitoring Cycle 3
	etherman.
		On("CheckTxWasMined", ctx, signedTx.Hash()).
		Return(true, receipt, nil).
		Once()
	st.
		On("GetLastBlock", ctx, nil).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(block, nil).
		Once()

	// Build Result 3
	etherman.
		On("GetTx", ctx, signedTx.Hash()).
		Return(signedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetRevertMessage", ctx, signedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.Start()

	time.Sleep(time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
	require.Equal(t, 1, len(result.Txs))
	require.Equal(t, signedTx, result.Txs[signedTx.Hash()].Tx)
	require.Equal(t, receipt, result.Txs[signedTx.Hash()].Receipt)
	require.Equal(t, "", result.Txs[signedTx.Hash()].RevertMessage)

	err = ethTxManagerClient.Reorg(ctx, blockNumber.Uint64()-3, nil)
	require.NoError(t, err)

	result, err = ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, MonitoredTxStatusReorged, result.Status)
	require.Equal(t, 1, len(result.Txs))
	require.Nil(t, result.Txs[signedTx.Hash()].Tx)
	require.Nil(t, result.Txs[signedTx.Hash()].Receipt)
	require.Equal(t, "", result.Txs[signedTx.Hash()].RevertMessage)

	// creates a new instance of client to avoid a race condition in the test code
	ethTxManagerClient = New(defaultEthTxmanagerConfigForTests, etherman, storage, st)

	go ethTxManagerClient.Start()

	time.Sleep(time.Second)
	result, err = ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
	require.Equal(t, 1, len(result.Txs))
	require.Equal(t, signedTx, result.Txs[signedTx.Hash()].Tx)
	require.Equal(t, receipt, result.Txs[signedTx.Hash()].Receipt)
	require.Equal(t, "", result.Txs[signedTx.Hash()].RevertMessage)
}

func TestExecutionReverted(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	st := newStateMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(defaultEthTxmanagerConfigForTests, etherman, storage, st)

	ctx := context.Background()

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	// Add
	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx, from).
		Return(currentNonce, nil).
		Once()

	firstGasEstimation := uint64(1)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(firstGasEstimation, nil).
		Once()

	firstGasPriceSuggestion := big.NewInt(1)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(firstGasPriceSuggestion, nil).
		Once()

	// Monitoring Cycle 1
	firstSignedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      firstGasEstimation,
		GasPrice: firstGasPriceSuggestion,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(firstSignedTx, nil).
		Once()
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("SendTx", ctx, firstSignedTx).
		Return(nil).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, firstSignedTx, mock.IsType(time.Second)).
		Return(true, nil).
		Once()

	blockNumber := big.NewInt(1)
	failedReceipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
		Status:      ethTypes.ReceiptStatusFailed,
		TxHash:      firstSignedTx.Hash(),
	}

	etherman.
		On("GetTxReceipt", ctx, firstSignedTx.Hash()).
		Return(failedReceipt, nil).
		Once()
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(firstSignedTx, false, nil).
		Once()
	etherman.
		On("GetRevertMessage", ctx, firstSignedTx).
		Return("", ErrExecutionReverted).
		Once()

	// Monitoring Cycle 2
	etherman.
		On("CheckTxWasMined", ctx, firstSignedTx.Hash()).
		Return(true, failedReceipt, nil).
		Once()

	currentNonce = uint64(2)
	etherman.
		On("CurrentNonce", ctx, from).
		Return(currentNonce, nil).
		Once()
	secondGasEstimation := uint64(2)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(secondGasEstimation, nil).
		Once()
	secondGasPriceSuggestion := big.NewInt(2)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(secondGasPriceSuggestion, nil).
		Once()

	secondSignedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      secondGasEstimation,
		GasPrice: secondGasPriceSuggestion,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, from, mock.IsType(&ethTypes.Transaction{})).
		Return(secondSignedTx, nil).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("SendTx", ctx, secondSignedTx).
		Return(nil).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, secondSignedTx, mock.IsType(time.Second)).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(true, nil).
		Once()

	blockNumber = big.NewInt(2)
	receipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
		Status:      ethTypes.ReceiptStatusSuccessful,
	}
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()

	block := &state.Block{
		BlockNumber: blockNumber.Uint64(),
	}
	st.
		On("GetLastBlock", ctx, nil).
		Return(block, nil).
		Once()

	// Build result
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(firstSignedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, firstSignedTx.Hash()).
		Return(nil, ethereum.NotFound).
		Once()
	etherman.
		On("GetRevertMessage", ctx, firstSignedTx).
		Return("", nil).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(secondSignedTx, false, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetRevertMessage", ctx, secondSignedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.Start()

	time.Sleep(time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
}
