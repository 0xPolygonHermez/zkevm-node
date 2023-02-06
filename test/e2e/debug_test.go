package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestDebugTraceTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	const l2NetworkURL = "http://localhost:8124"
	const l2ExplorerRPCComponentName = "l2-explorer-json-rpc"

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
		require.NoError(t, operations.StopComponent(l2ExplorerRPCComponentName))
	}()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	err = operations.StartComponent(l2ExplorerRPCComponentName, func() (bool, error) { return operations.NodeUpCondition(l2NetworkURL) })
	require.NoError(t, err)

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       "Local L1",
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       "Local L2",
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	results := map[string]json.RawMessage{}

	type testCase struct {
		name                   string
		networkPreparationFunc func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) ([]interface{}, error)
		txCreateFunc           func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData []interface{}) (*types.Transaction, error)
	}
	testCases := []testCase{
		{name: "eth transfer", txCreateFunc: createEthTransferTx},
		// {name: "sc deployment", networkPreparationFunc: , txCreateFunc: },
		// {name: "sc call", networkPreparationFunc: ,txCreateFunc: },
		// {name: "erc20 transfer", networkPreparationFunc: , txCreateFunc: },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, network := range networks {
				log.Debugf(network.Name)
				client := operations.MustGetClient(network.URL)
				auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

				var customData []interface{}
				if tc.networkPreparationFunc != nil {
					customData, err = tc.networkPreparationFunc(t, ctx, auth, client)
					require.NoError(t, err)
				}

				signedTx, err := tc.txCreateFunc(t, ctx, auth, client, customData)
				require.NoError(t, err)

				err = client.SendTransaction(ctx, signedTx)
				require.NoError(t, err)

				err = operations.WaitTxToBeMined(ctx, client, signedTx, operations.DefaultTimeoutTxToBeMined)
				require.NoError(t, err)

				log.Debug("***************************************", signedTx.Hash().String())

				response, err := jsonrpc.JSONRPCCall(network.URL, "debug_traceTransaction", signedTx.Hash().String())
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)

				results[network.Name] = response.Result
			}

			referenceValueMap := map[string]interface{}{}
			err = json.Unmarshal(results[networks[0].Name], &referenceValueMap)
			require.NoError(t, err)

			for networkName, result := range results {
				resultMap := map[string]interface{}{}
				err = json.Unmarshal(result, &resultMap)
				require.NoError(t, err)
				diff := deep.Equal(referenceValueMap, resultMap)
				require.Nil(t, diff, fmt.Sprintf("invalid trace for network %s: %v", networkName, diff))
			}
		})
	}
}

func createEthTransferTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData []interface{}) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	to := common.HexToAddress("0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98")

	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: auth.From,
		To:   &to,
	})

	require.NoError(t, err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		GasPrice: gasPrice,
		Gas:      gas,
	})

	return auth.Signer(auth.From, tx)
}
