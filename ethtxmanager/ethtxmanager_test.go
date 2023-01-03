package ethtxmanager

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTxGetMined(t *testing.T) {
	cfg := Config{
		FrequencyForResendingFailedTxs: types.NewDuration(time.Second),
		WaitTxToBeMined:                types.NewDuration(1 * time.Minute),
	}
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(cfg, etherman, storage)

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	ctx := context.Background()

	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx).
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
		On("SignTx", ctx, mock.IsType(&ethTypes.Transaction{})).
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
		Return(nil).
		Once()

	receipt := &ethTypes.Receipt{
		Status: ethTypes.ReceiptStatusSuccessful,
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
		On("GetRevertMessage", ctx, *signedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.ManageTxs()

	time.Sleep(5 * time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
	require.Equal(t, 1, len(result.Txs))
	require.Equal(t, *signedTx, result.Txs[signedTx.Hash()].Tx)
	require.Equal(t, receipt, result.Txs[signedTx.Hash()].Receipt)
	require.Equal(t, "", result.Txs[signedTx.Hash()].RevertMessage)
}

func TestTxGetMinedAfterReviewed(t *testing.T) {
	cfg := Config{
		FrequencyForResendingFailedTxs: types.NewDuration(time.Second),
		WaitTxToBeMined:                types.NewDuration(1 * time.Minute),
	}
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(cfg, etherman, storage)

	ctx := context.Background()

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx).
		Return(currentNonce, nil).
		Once()

	firstGasEstimation := uint64(1)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(firstGasEstimation, nil).
		Once()
	secondGasEstimation := uint64(2)
	etherman.
		On("EstimateGas", ctx, from, to, value, data).
		Return(secondGasEstimation, nil).
		Once()

	firstGasPriceSuggestion := big.NewInt(1)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(firstGasPriceSuggestion, nil).
		Once()
	secondGasPriceSuggestion := big.NewInt(2)
	etherman.
		On("SuggestedGasPrice", ctx).
		Return(secondGasPriceSuggestion, nil).
		Once()

	firstSignedTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      firstGasEstimation,
		GasPrice: firstGasPriceSuggestion,
		Data:     data,
	})
	etherman.
		On("SignTx", ctx, mock.IsType(&ethTypes.Transaction{})).
		Return(firstSignedTx, nil).
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
		On("SignTx", ctx, mock.IsType(&ethTypes.Transaction{})).
		Return(secondSignedTx, nil).
		Once()

	etherman.
		On("CheckTxWasMined", ctx, firstSignedTx.Hash()).
		Return(false, nil, nil).
		Once()

	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(nil, false, ethereum.NotFound).
		Once()
	etherman.
		On("GetTx", ctx, firstSignedTx.Hash()).
		Return(firstSignedTx, false, nil).
		Once()
	etherman.
		On("GetTx", ctx, secondSignedTx.Hash()).
		Return(secondSignedTx, false, nil).
		Once()

	etherman.
		On("SendTx", ctx, firstSignedTx).
		Return(nil).
		Once()
	etherman.
		On("SendTx", ctx, secondSignedTx).
		Return(nil).
		Once()

	etherman.
		On("WaitTxToBeMined", ctx, firstSignedTx, mock.IsType(time.Second)).
		Return(errors.New("tx not mined yet")).
		Once()
	etherman.
		On("WaitTxToBeMined", ctx, secondSignedTx, mock.IsType(time.Second)).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(nil).
		Once()

	receipt := &ethTypes.Receipt{
		Status: ethTypes.ReceiptStatusSuccessful,
	}
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()
	etherman.
		On("GetTxReceipt", ctx, firstSignedTx.Hash()).
		Return(nil, ethereum.NotFound).
		Once()
	etherman.
		On("GetTxReceipt", ctx, secondSignedTx.Hash()).
		Return(receipt, nil).
		Once()

	etherman.
		On("GetRevertMessage", ctx, *secondSignedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.ManageTxs()

	time.Sleep(5 * time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
}

func TestTxGetMinedAfterFailingToWaiting(t *testing.T) {
	cfg := Config{
		FrequencyForResendingFailedTxs: types.NewDuration(time.Second),
		WaitTxToBeMined:                types.NewDuration(1 * time.Minute),
	}
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newEthermanMock(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(cfg, etherman, storage)

	ctx := context.Background()

	owner := "owner"
	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	currentNonce := uint64(1)
	etherman.
		On("CurrentNonce", ctx).
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
		On("SignTx", ctx, mock.IsType(&ethTypes.Transaction{})).
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
		Return(errors.New("tx not mined yet")).
		Once()

	receipt := &ethTypes.Receipt{
		Status: ethTypes.ReceiptStatusSuccessful,
	}

	etherman.
		On("GetTxReceipt", ctx, signedTx.Hash()).
		Return(receipt, nil).
		Once()

	etherman.
		On("CheckTxWasMined", ctx, signedTx.Hash()).
		Run(func(args mock.Arguments) { ethTxManagerClient.Stop() }). // stops the management cycle to avoid problems with mocks
		Return(true, receipt, nil).
		Once()

	etherman.
		On("GetRevertMessage", ctx, *signedTx).
		Return("", nil).
		Once()

	err = ethTxManagerClient.Add(ctx, owner, id, from, to, value, data, nil)
	require.NoError(t, err)

	go ethTxManagerClient.ManageTxs()

	time.Sleep(5 * time.Second)
	result, err := ethTxManagerClient.Result(ctx, owner, id, nil)
	require.NoError(t, err)
	require.Equal(t, MonitoredTxStatusConfirmed, result.Status)
}
