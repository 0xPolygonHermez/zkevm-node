package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Revert2"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const fixedTxGasLimit uint64 = 100000

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
		{name: "eth transfer", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment", createSignedTx: createScDeploySignedTx},
		{name: "sc call", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
		// failed transactions
		{name: "sc deployment reverted", createSignedTx: createScDeployRevertedSignedTx},
		{name: "sc call reverted", prepare: prepareScCallReverted, createSignedTx: createScCallRevertedSignedTx},
		{name: "erc20 transfer reverted", prepare: prepareERC20TransferReverted, createSignedTx: createERC20TransferRevertedSignedTx},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debugf(tc.name)
			for _, network := range networks {
				log.Debugf(network.Name)
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

				debugOptions := map[string]interface{}{
					"disableStorage":   false,
					"disableStack":     false,
					"enableMemory":     true,
					"enableReturnData": true,
				}

				response, err := client.JSONRPCCall(network.URL, "debug_traceTransaction", signedTx.Hash().String(), debugOptions)
				require.NoError(t, err)
				require.Nil(t, response.Error)
				require.NotNil(t, response.Result)

				results[network.Name] = response.Result

				// // save result in a file
				// sanitizedNetworkName := strings.ReplaceAll(network.Name+"_"+tc.name, " ", "_")
				// filePath := fmt.Sprintf("/home/tclemos/github.com/0xPolygonHermez/zkevm-node/dist/%v.json", sanitizedNetworkName)
				// b, _ := signedTx.MarshalBinary()
				// fileContent := struct {
				// 	Tx    *ethTypes.Transaction
				// 	RLP   string
				// 	Trace json.RawMessage
				// }{
				// 	Tx:    signedTx,
				// 	RLP:   hex.EncodeToHex(b),
				// 	Trace: response.Result,
				// }
				// c, err := json.MarshalIndent(fileContent, "", "    ")
				// require.NoError(t, err)
				// err = os.WriteFile(filePath, c, 0644)
				// require.NoError(t, err)
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
		{name: "eth transfer by number", blockNumberOrHash: "number", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by number", blockNumberOrHash: "number", createSignedTx: createScDeploySignedTx},
		{name: "sc call by number", blockNumberOrHash: "number", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by number", blockNumberOrHash: "number", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},

		{name: "eth transfer by hash", blockNumberOrHash: "hash", createSignedTx: createEthTransferSignedTx},
		{name: "sc deployment by hash", blockNumberOrHash: "hash", createSignedTx: createScDeploySignedTx},
		{name: "sc call by hash", blockNumberOrHash: "hash", prepare: prepareScCall, createSignedTx: createScCallSignedTx},
		{name: "erc20 transfer by hash", blockNumberOrHash: "hash", prepare: prepareERC20Transfer, createSignedTx: createERC20TransferSignedTx},
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
			for _, network := range networks {
				log.Debugf(network.Name)
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

func createEthTransferSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
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

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		GasPrice: gasPrice,
		Gas:      gas,
	})

	return auth.Signer(auth.From, tx)
}

func createScDeploySignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	scByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)
	data := common.Hex2Bytes(scByteCode)

	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: auth.From,
		Data: data,
	})
	require.NoError(t, err)

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		Data:     data,
	})

	return auth.Signer(auth.From, tx)
}

func prepareScCall(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error) {
	_, tx, sc, err := EmitLog.DeployEmitLog(auth, client)
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	return map[string]interface{}{
		"sc": sc,
	}, nil
}

func createScCallSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	scInterface := customData["sc"]
	sc := scInterface.(*EmitLog.EmitLog)

	opts := *auth
	opts.NoSend = true

	tx, err := sc.EmitLogs(&opts)
	require.NoError(t, err)

	return tx, nil
}

func prepareERC20Transfer(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error) {
	_, tx, sc, err := ERC20.DeployERC20(auth, client, "MyToken", "MT")
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	tx, err = sc.Mint(auth, big.NewInt(1000000000))
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	return map[string]interface{}{
		"sc": sc,
	}, nil
}

func createERC20TransferSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	scInterface := customData["sc"]
	sc := scInterface.(*ERC20.ERC20)

	opts := *auth
	opts.NoSend = true

	to := common.HexToAddress("0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98")

	tx, err := sc.Transfer(&opts, to, big.NewInt(123456))
	require.NoError(t, err)

	return tx, nil
}

func createScDeployRevertedSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	scByteCode, err := testutils.ReadBytecode("Revert/Revert.bin")
	require.NoError(t, err)
	data := common.Hex2Bytes(scByteCode)

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      fixedTxGasLimit,
		Data:     data,
	})

	return auth.Signer(auth.From, tx)
}

func prepareScCallReverted(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error) {
	_, tx, sc, err := Revert2.DeployRevert2(auth, client)
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	return map[string]interface{}{
		"sc": sc,
	}, nil
}

func createScCallRevertedSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	scInterface := customData["sc"]
	sc := scInterface.(*Revert2.Revert2)

	opts := *auth
	opts.NoSend = true
	opts.GasLimit = fixedTxGasLimit

	tx, err := sc.GenerateError(&opts)
	require.NoError(t, err)

	return tx, nil
}

func prepareERC20TransferReverted(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) (map[string]interface{}, error) {
	_, tx, sc, err := ERC20.DeployERC20(auth, client, "MyToken2", "MT2")
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
	require.NoError(t, err)

	return map[string]interface{}{
		"sc": sc,
	}, nil
}

func createERC20TransferRevertedSignedTx(t *testing.T, ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, customData map[string]interface{}) (*ethTypes.Transaction, error) {
	scInterface := customData["sc"]
	sc := scInterface.(*ERC20.ERC20)

	opts := *auth
	opts.NoSend = true
	opts.GasLimit = fixedTxGasLimit

	to := common.HexToAddress("0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98")

	tx, err := sc.Transfer(&opts, to, big.NewInt(123456))
	require.NoError(t, err)

	return tx, nil
}
