package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestDebugTraceTransactionNotFoundTx(t *testing.T) {
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

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	for _, network := range networks {
		log.Debugf(network.Name)
		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce: 10,
		})

		response, err := client.JSONRPCCall(network.URL, "debug_traceTransaction", tx.Hash().String())
		require.NoError(t, err)
		require.Nil(t, response.Result)
		require.NotNil(t, response.Error)

		require.Equal(t, -32000, response.Error.Code)
		require.Equal(t, "transaction not found", response.Error.Message)
		require.Nil(t, response.Error.Data)
	}
}

func TestDebugTraceBlockByNumberNotFoundTx(t *testing.T) {
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

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	for _, network := range networks {
		log.Debugf(network.Name)

		response, err := client.JSONRPCCall(network.URL, "debug_traceBlockByNumber", hex.EncodeBig(big.NewInt(999999999999)))
		require.NoError(t, err)
		require.Nil(t, response.Result)
		require.NotNil(t, response.Error)

		require.Equal(t, -32000, response.Error.Code)
		require.Equal(t, "block #999999999999 not found", response.Error.Message)
		require.Nil(t, response.Error.Data)
	}
}

func TestDebugTraceBlockByHashNotFoundTx(t *testing.T) {
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

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	for _, network := range networks {
		log.Debugf(network.Name)

		response, err := client.JSONRPCCall(network.URL, "debug_traceBlockByHash", common.Hash{}.String())
		require.NoError(t, err)
		require.Nil(t, response.Result)
		require.NotNil(t, response.Error)

		require.Equal(t, -32000, response.Error.Code)
		require.Equal(t, "block 0x0000000000000000000000000000000000000000000000000000000000000000 not found", response.Error.Message)
		require.Nil(t, response.Error.Data)
	}
}

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

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	results := map[string]json.RawMessage{}

	type testCase struct {
		name           string
		prepare        func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error)
		createSignedTx func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error)
	}
	testCases := []testCase{
		// successful transactions
		// {name: "eth transfer", createSignedTx: createEthTransferSignedTx},
		// {name: "sc deployment", createSignedTx: createScDeploySignedTx},
		// {name: "sc call", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		// {name: "erc20 transfer", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		// {name: "create", prepare: prepareCreate, createSignedTx: createCreateSignedTx},
		// {name: "create2", prepare: prepareCreate, createSignedTx: createCreate2SignedTx},
		{name: "delegatecall", prepare: prepareDelegateCall, createSignedTx: createDelegateCallSignedTx},
		// failed transactions
		// {name: "sc deployment reverted", createSignedTx: createScDeployRevertedSignedTx},
		// {name: "sc call reverted", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		// {name: "erc20 transfer reverted", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
	}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	for _, network := range networks {
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(network.ChainID))
		require.NoError(t, err)

		ethereumClient := operations.MustGetClient(network.URL)
		sourceAuth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		nonce, err := ethereumClient.NonceAt(ctx, sourceAuth.From, nil)
		require.NoError(t, err)

		balance, err := ethereumClient.BalanceAt(ctx, sourceAuth.From, nil)
		require.NoError(t, err)

		gasPrice, err := ethereumClient.SuggestGasPrice(ctx)
		require.NoError(t, err)

		value := big.NewInt(0).Quo(balance, big.NewInt(2))

		gas, err := ethereumClient.EstimateGas(ctx, ethereum.CallMsg{
			From:     sourceAuth.From,
			To:       &auth.From,
			GasPrice: gasPrice,
			Value:    value,
		})
		require.NoError(t, err)

		tx := ethTypes.NewTx(&ethTypes.LegacyTx{
			To:       &auth.From,
			Nonce:    nonce,
			GasPrice: gasPrice,
			Value:    value,
			Gas:      gas,
		})

		signedTx, err := sourceAuth.Signer(sourceAuth.From, tx)
		require.NoError(t, err)

		err = ethereumClient.SendTransaction(ctx, signedTx)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")

			for _, network := range networks {
				log.Debug("------------------------ ", network.Name, " ------------------------")
				ethereumClient := operations.MustGetClient(network.URL)
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(network.ChainID))
				require.NoError(t, err)

				var customData map[string]interface{}
				if tc.prepare != nil {
					customData, err = tc.prepare(t, ctx, auth, ethereumClient)
					require.NoError(t, err)
				}

				signedTx, err := tc.createSignedTx(t, ctx, auth, ethereumClient, customData)
				require.NoError(t, err)

				balance, err := ethereumClient.BalanceAt(ctx, auth.From, nil)
				require.NoError(t, err)

				log.Debugf("balance of %v: %v", auth.From, balance.String())

				err = ethereumClient.SendTransaction(ctx, signedTx)
				require.NoError(t, err)

				log.Debugf("tx sent: %v", signedTx.Hash().String())

				err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
				if err != nil && !strings.HasPrefix(err.Error(), "transaction has failed, reason:") {
					require.NoError(t, err)
				}

				debugOptions := map[string]interface{}{
					"disableStorage":   true,
					"disableStack":     true,
					"enableMemory":     false,
					"enableReturnData": false,
				}

				response, err := client.JSONRPCCall(network.URL, "debug_traceTransaction", signedTx.Hash().String(), debugOptions)
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)

				results[network.Name] = response.Result

				// save result in a file
				sanitizedNetworkName := strings.ReplaceAll(network.Name+"_"+tc.name, " ", "_")
				filePath := fmt.Sprintf("/Users/thiago/github.com/0xPolygonHermez/zkevm-node/dist/%v.json", sanitizedNetworkName)
				b, _ := signedTx.MarshalBinary()
				fileContent := struct {
					Tx    *ethTypes.Transaction
					RLP   string
					Trace json.RawMessage
				}{
					Tx:    signedTx,
					RLP:   hex.EncodeToHex(b),
					Trace: response.Result,
				}
				c, err := json.MarshalIndent(fileContent, "", "    ")
				require.NoError(t, err)
				err = os.WriteFile(filePath, c, 0644)
				require.NoError(t, err)
			}

			referenceValueMap := map[string]interface{}{}
			err = json.Unmarshal(results[l1NetworkName], &referenceValueMap)
			require.NoError(t, err)

			referenceStructLogsMap := referenceValueMap["structLogs"].([]interface{})

			for networkName, result := range results {
				if networkName == l1NetworkName {
					continue
				}

				resultMap := map[string]interface{}{}
				err = json.Unmarshal(result, &resultMap)
				require.NoError(t, err)

				require.Equal(t, referenceValueMap["failed"], resultMap["failed"], fmt.Sprintf("invalid `failed` for network %s", networkName))

				resultStructLogsMap := resultMap["structLogs"].([]interface{})
				require.Equal(t, len(referenceStructLogsMap), len(resultStructLogsMap))

				for structLogIndex := range referenceStructLogsMap {
					referenceStructLogMap := referenceStructLogsMap[structLogIndex].(map[string]interface{})
					resultStructLogMap := resultStructLogsMap[structLogIndex].(map[string]interface{})

					require.Equal(t, referenceStructLogMap["pc"], resultStructLogMap["pc"], fmt.Sprintf("invalid struct log pc for network %s", networkName))
					require.Equal(t, referenceStructLogMap["op"], resultStructLogMap["op"], fmt.Sprintf("invalid struct log op for network %s", networkName))
					require.Equal(t, referenceStructLogMap["depth"], resultStructLogMap["depth"], fmt.Sprintf("invalid struct log depth for network %s", networkName))

					pc := referenceStructLogMap["pc"]
					op := referenceStructLogMap["op"]

					referenceStack, found := referenceStructLogMap["stack"].([]interface{})
					if found {
						resultStack := resultStructLogMap["stack"].([]interface{})

						require.Equal(t, len(referenceStack), len(resultStack), fmt.Sprintf("stack size doesn't match for pc %v op %v", pc, op))
						for stackIndex := range referenceStack {
							require.Equal(t, referenceStack[stackIndex], resultStack[stackIndex], fmt.Sprintf("stack index %v doesn't match for pc %v op %v", stackIndex, pc, op))
						}
					}

					referenceMemory, found := referenceStructLogMap["memory"].([]interface{})
					if found {
						resultMemory := resultStructLogMap["memory"].([]interface{})

						require.Equal(t, len(referenceMemory), len(resultMemory), fmt.Sprintf("memory size doesn't match for pc %v op %v", pc, op))
						for memoryIndex := range referenceMemory {
							require.Equal(t, referenceMemory[memoryIndex], resultMemory[memoryIndex], fmt.Sprintf("memory index %v doesn't match for pc %v op %v", memoryIndex, pc, op))
						}
					}

					referenceStorage, found := referenceStructLogMap["storage"].(map[string]interface{})
					if found {
						resultStorage := resultStructLogMap["storage"].(map[string]interface{})

						require.Equal(t, len(referenceStorage), len(resultStorage), fmt.Sprintf("storage size doesn't match for pc %v op %v", pc, op))
						for storageKey, referenceStorageValue := range referenceStorage {
							resultStorageValue, found := resultStorage[storageKey]
							require.True(t, found, "storage address not found")
							require.Equal(t, referenceStorageValue, resultStorageValue, fmt.Sprintf("storage value doesn't match for address %v for pc %v op %v", storageKey, pc, op))
						}
					}
				}
			}
		})
	}
}

func TestDebugTraceBlock(t *testing.T) {
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

	const l1NetworkName, l2NetworkName = "Local L1", "Local L2"

	networks := []struct {
		Name         string
		URL          string
		WebSocketURL string
		ChainID      uint64
		PrivateKey   string
	}{
		{
			Name:       l1NetworkName,
			URL:        operations.DefaultL1NetworkURL,
			ChainID:    operations.DefaultL1ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
		{
			Name:       l2NetworkName,
			URL:        l2NetworkURL,
			ChainID:    operations.DefaultL2ChainID,
			PrivateKey: operations.DefaultSequencerPrivateKey,
		},
	}

	results := map[string]json.RawMessage{}

	type testCase struct {
		name              string
		blockNumberOrHash string
		prepare           func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error)
		createSignedTx    func(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error)
	}
	testCases := []testCase{
		// successful transactions
		// by block number
		{name: "eth transfer by number", blockNumberOrHash: "number", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by number", blockNumberOrHash: "number", createSignedTx: createScDeploySignedTx},
		{name: "sc call by number", blockNumberOrHash: "number", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by number", blockNumberOrHash: "number", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		{name: "create by number", blockNumberOrHash: "number", prepare: prepareCreate, createSignedTx: createCreateSignedTx},
		{name: "create2 by number", blockNumberOrHash: "number", prepare: prepareCreate, createSignedTx: createCreate2SignedTx},
		{name: "delegatecall by number", blockNumberOrHash: "number", prepare: prepareDelegateCall, createSignedTx: createDelegateCallSignedTx},
		// by block hash
		{name: "eth transfer by hash", blockNumberOrHash: "hash", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by hash", blockNumberOrHash: "hash", createSignedTx: createScDeploySignedTx},
		{name: "sc call by hash", blockNumberOrHash: "hash", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by hash", blockNumberOrHash: "hash", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		{name: "create by hash", blockNumberOrHash: "hash", prepare: prepareCreate, createSignedTx: createCreateSignedTx},
		{name: "create2 by hash", blockNumberOrHash: "hash", prepare: prepareCreate, createSignedTx: createCreate2SignedTx},
		{name: "delegatecall by hash", blockNumberOrHash: "hash", prepare: prepareDelegateCall, createSignedTx: createDelegateCallSignedTx},

		// failed transactions
		{name: "sc deployment reverted by number", blockNumberOrHash: "number", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc call reverted by number", blockNumberOrHash: "number", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted by number", blockNumberOrHash: "number", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},

		{name: "sc deployment reverted by hash", blockNumberOrHash: "hash", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc call reverted by hash", blockNumberOrHash: "hash", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted by hash", blockNumberOrHash: "hash", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")

			for _, network := range networks {
				log.Debug("------------------------ ", network.Name, " ------------------------")
				ethereumClient := operations.MustGetClient(network.URL)
				auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

				var customData map[string]interface{}
				if tc.prepare != nil {
					customData, err = tc.prepare(t, ctx, auth, ethereumClient)
					require.NoError(t, err)
				}

				signedTx, err := tc.createSignedTx(t, ctx, auth, ethereumClient, customData)
				require.NoError(t, err)

				err = ethereumClient.SendTransaction(ctx, signedTx)
				require.NoError(t, err)

				log.Debugf("tx sent: %v", signedTx.Hash().String())

				err = operations.WaitTxToBeMined(ctx, ethereumClient, signedTx, operations.DefaultTimeoutTxToBeMined)
				if err != nil && !strings.HasPrefix(err.Error(), "transaction has failed, reason:") {
					require.NoError(t, err)
				}

				receipt, err := ethereumClient.TransactionReceipt(ctx, signedTx.Hash())
				require.NoError(t, err)

				debugOptions := map[string]interface{}{
					"disableStorage":   false,
					"disableStack":     false,
					"enableMemory":     true,
					"enableReturnData": true,
				}

				var response types.Response
				if tc.blockNumberOrHash == "number" {
					response, err = client.JSONRPCCall(network.URL, "debug_traceBlockByNumber", hex.EncodeBig(receipt.BlockNumber), debugOptions)
				} else {
					response, err = client.JSONRPCCall(network.URL, "debug_traceBlockByHash", receipt.BlockHash.String(), debugOptions)
				}
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)

				results[network.Name] = response.Result
			}

			referenceTransactions := []interface{}{}
			err = json.Unmarshal(results[l1NetworkName], &referenceTransactions)
			require.NoError(t, err)

			for networkName, result := range results {
				if networkName == l1NetworkName {
					continue
				}

				resultTransactions := []interface{}{}
				err = json.Unmarshal(result, &resultTransactions)
				require.NoError(t, err)

				for transactionIndex := range referenceTransactions {
					referenceTransactionMap := referenceTransactions[transactionIndex].(map[string]interface{})
					referenceResultMap := referenceTransactionMap["result"].(map[string]interface{})
					referenceStructLogsMap := referenceResultMap["structLogs"].([]interface{})

					resultTransactionMap := resultTransactions[transactionIndex].(map[string]interface{})
					resultResultMap := resultTransactionMap["result"].(map[string]interface{})
					resultStructLogsMap := resultResultMap["structLogs"].([]interface{})

					require.Equal(t, len(referenceStructLogsMap), len(resultStructLogsMap))

					for structLogIndex := range referenceStructLogsMap {
						referenceStructLogMap := referenceStructLogsMap[structLogIndex].(map[string]interface{})
						resultStructLogMap := resultStructLogsMap[structLogIndex].(map[string]interface{})

						require.Equal(t, referenceStructLogMap["pc"], resultStructLogMap["pc"], fmt.Sprintf("invalid struct log pc for network %s", networkName))
						require.Equal(t, referenceStructLogMap["op"], resultStructLogMap["op"], fmt.Sprintf("invalid struct log op for network %s", networkName))
						require.Equal(t, referenceStructLogMap["depth"], resultStructLogMap["depth"], fmt.Sprintf("invalid struct log depth for network %s", networkName))

						pc := referenceStructLogMap["pc"]
						op := referenceStructLogMap["op"]

						referenceStack, found := referenceStructLogMap["stack"].([]interface{})
						if found {
							resultStack := resultStructLogMap["stack"].([]interface{})

							require.Equal(t, len(referenceStack), len(resultStack), fmt.Sprintf("stack size doesn't match for pc %v op %v", pc, op))
							for stackIndex := range referenceStack {
								require.Equal(t, referenceStack[stackIndex], resultStack[stackIndex], fmt.Sprintf("stack index %v doesn't match for pc %v op %v", stackIndex, pc, op))
							}
						}

						referenceMemory, found := referenceStructLogMap["memory"].([]interface{})
						if found {
							resultMemory := resultStructLogMap["memory"].([]interface{})

							require.Equal(t, len(referenceMemory), len(resultMemory), fmt.Sprintf("memory size doesn't match for pc %v op %v", pc, op))
							for memoryIndex := range referenceMemory {
								require.Equal(t, referenceMemory[memoryIndex], resultMemory[memoryIndex], fmt.Sprintf("memory index %v doesn't match for pc %v op %v", memoryIndex, pc, op))
							}
						}

						referenceStorage, found := referenceStructLogMap["storage"].(map[string]interface{})
						if found {
							resultStorage := resultStructLogMap["storage"].(map[string]interface{})

							require.Equal(t, len(referenceStorage), len(resultStorage), fmt.Sprintf("storage size doesn't match for pc %v op %v", pc, op))
							for storageKey, referenceStorageValue := range referenceStorage {
								resultStorageValue, found := resultStorage[storageKey]
								require.True(t, found, "storage address not found")
								require.Equal(t, referenceStorageValue, resultStorageValue, fmt.Sprintf("storage value doesn't match for address %v for pc %v op %v", storageKey, pc, op))
							}
						}
					}
				}
			}
		})
	}
}
