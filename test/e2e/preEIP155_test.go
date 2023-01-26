package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestPreEIP155Tx(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		nonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)

		to := common.HexToAddress("0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98")
		data := hex.DecodeBig("0x64fbb77c").Bytes()

		gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From: auth.From,
			To:   &to,
			Data: data,
		})
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			GasPrice: gasPrice,
			Gas:      gas,
			Data:     data,
		})

		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(network.PrivateKey, "0x"))
		require.NoError(t, err)

		signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
		require.NoError(t, err)

		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, signedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		require.NoError(t, err)

		// wait for l2 block to be virtualized
		if network.ChainID == operations.DefaultL2ChainID {
			log.Infof("waiting for the block number %v to be virtualized", receipt.BlockNumber.String())
			err = operations.WaitL2BlockToBeVirtualized(receipt.BlockNumber, 4*time.Minute) //nolint:gomnd
			require.NoError(t, err)
		}
	}
}
