package ethtxmanager

import (
	"context"
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

func TestSend(t *testing.T) {
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

	err = ethTxManagerClient.Add(ctx, id, from, to, value, data, nil)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	status, err := ethTxManagerClient.Status(ctx, id, nil)
	require.NoError(t, err)
	require.Equal(t, MonitoredTxStatusConfirmed, status)
}
