package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Double"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Storage"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	invalidParamsErrorCode = -32602
	defaultErrorCode       = -32000
)

func Setup() {
	var err error
	ctx := context.Background()
	err = operations.Teardown()
	if err != nil {
		panic(err)
	}

	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	if err != nil {
		panic(err)
	}
	err = opsMan.Setup()
	if err != nil {
		panic(err)
	}
}

func Teardown() {
	err := operations.Teardown()
	if err != nil {
		panic(err)
	}
}

// TestJSONRPC tests JSON RPC methods on a running environment.
func TestJSONRPC(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		sc, err := deployContracts(network.URL, operations.DefaultSequencerPrivateKey, network.ChainID)
		require.NoError(t, err)

		callOpts := &bind.CallOpts{Pending: false}

		payload := big.NewInt(5)
		number, err := sc.Double(callOpts, payload)
		require.NoError(t, err)
		expected := big.NewInt(0).Mul(payload, big.NewInt(2))
		require.Equal(t, expected, number)
	}
}

func deployContracts(url, privateKey string, chainId uint64) (*Double.Double, error) {
	client := operations.MustGetClient(url)
	auth := operations.MustGetAuth(privateKey, chainId)

	_, scTx, sc, err := Double.DeployDouble(auth, client)
	if err != nil {
		return nil, err
	}
	err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func createTX(ethdeployment string, chainId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	client, err := ethclient.Dial(ethdeployment)
	if err != nil {
		return nil, err
	}
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, chainId)
	if err != nil {
		return nil, err
	}
	nonce, err := client.NonceAt(context.Background(), auth.From, nil)
	if err != nil {
		return nil, err
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &to, Value: amount})
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	log.Infof("\nTX details:\n\tNonce:    %d\n\tGasLimit: %d\n\tGasPrice: %d", nonce, gasLimit, gasPrice)
	if gasLimit != uint64(21000) {
		return nil, fmt.Errorf("gasLimit %d != 21000", gasLimit)
	}
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}
	log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func Test_Filters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		response, err := jsonrpc.JSONRPCCall(network.URL, "eth_newBlockFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var filterId string
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newPendingTransactionFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		filterId = ""
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"BlockHash": common.HexToHash("0x1"),
			"FromBlock": "0x1",
			"ToBlock":   "0x2",
		})
		require.NoError(t, err)

		require.NotNil(t, response.Error)
		require.Equal(t, invalidParamsErrorCode, response.Error.Code)
		require.Equal(t, "invalid argument 0: cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other", response.Error.Message)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"BlockHash": common.HexToHash("0x1"),
			"Addresses": []common.Address{
				common.HexToAddress("0x2"),
			},
			"Topics": [][]common.Hash{
				{common.HexToHash("0x3")},
			},
		})
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		filterId = ""
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"FromBlock": "0x1",
			"ToBlock":   "0x2",
			"Addresses": []common.Address{
				common.HexToAddress("0x2"),
			},
			"Topics": [][]common.Hash{
				{common.HexToHash("0x3")},
			},
		})
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		filterId = ""
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var uninstalled bool
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.True(t, uninstalled)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		uninstalled = true
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.False(t, uninstalled)

		// logs

		// generate some logs first...
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)
		_, scTx, sc, err := EmitLog2.DeployEmitLog2(auth, client)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scCallTx, err := sc.EmitLogs(auth)
		require.NoError(t, err)

		err = operations.WaitTxToBeMined(ctx, client, scCallTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_getLogs", map[string]interface{}{
			"Addresses": []common.Address{common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff")},
			"Topics":    nil,
			"FromBlock": "0x0",
			"ToBlock":   "latest",
		})
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var logs []types.Log
		err = json.Unmarshal(response.Result, &logs)
		require.NoError(t, err)
		log.Infof("\nHow many logs: %d\n", len(logs))
		if len(logs) >= 2 {
			require.GreaterOrEqual(t, logs[1].BlockNumber, logs[0].BlockNumber)
		}

		// GetFilterChanges - passing an ID from newBlockFilter

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newBlockFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var blockFilterId string
		err = json.Unmarshal(response.Result, &blockFilterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		// Create TX: new block on L2 (l1 block generated every 1s)

		tx, err := createTX(network.URL, network.ChainID, common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"), big.NewInt(1000))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_getFilterChanges", blockFilterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var changes []common.Hash
		err = json.Unmarshal(response.Result, &changes)
		require.NoError(t, err)
		log.Infof("\nHow many changes: %d\n", len(changes))
		require.NotEmpty(t, changes)

		if len(changes) >= 1 {
			for _, change := range changes {
				log.Infof("\n> %s => %s", change, receipt.BlockHash)
			}
		}
		require.Contains(t, changes, receipt.BlockHash)

		// Wrong [any]FilterID

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_getFilterChanges", common.HexToHash("0x42"))
		require.NoError(t, err)
		require.NotNil(t, response.Error)
		require.Equal(t, defaultErrorCode, response.Error.Code)
		require.Equal(t, "filter not found", response.Error.Message)
		require.Nil(t, response.Result)
	}
}

func Test_Gas(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	var Address1 = common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff")
	var Values = []*big.Int{
		big.NewInt(1000),
		big.NewInt(10000000),
		big.NewInt(100000000000),
		big.NewInt(1000000000000000),
	}

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		for _, value := range Values {
			client, err := ethclient.Dial(network.URL)
			require.NoError(t, err)
			msg := ethereum.CallMsg{From: common.HexToAddress(operations.DefaultSequencerAddress),
				To:    &Address1,
				Value: value}

			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(operations.DefaultSequencerAddress), nil)
			require.NoError(t, err)

			log.Infof("Balance: %d", balance)
			require.GreaterOrEqual(t, balance.Cmp(big.NewInt(1)), 1)

			response, err := client.EstimateGas(context.Background(), msg)
			require.NoError(t, err)
			require.NotNil(t, response)
			log.Infof("Estimated gas: %d", response)
			require.GreaterOrEqual(t, response, uint64(21000))
		}
	}
}

func Test_Block(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	type rpcTx struct {
		BlockHash        string `json:"blockHash"`
		BlockNumber      string `json:"blockNumber"`
		ChainID          string `json:"chainId"`
		From             string `json:"from"`
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`
		Hash             string `json:"hash"`
		Input            string `json:"input"`
		Nonce            string `json:"nonce"`
		PublicKey        string `json:"publicKey"`
		R                string `json:"r"`
		Raw              string `json:"raw"`
		S                string `json:"s"`
		To               string `json:"to"`
		TransactionIndex string `json:"transactionIndex"`
		V                string `json:"v"`
		Value            string `json:"value"`
	}

	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		client, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		tx, err := createTX(network.URL, network.ChainID, common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"), big.NewInt(1000))
		require.NoError(t, err)
		// no block number yet... will wait
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)
		require.Equal(t, receipt.TxHash, tx.Hash())
		require.Equal(t, receipt.Type, tx.Type())
		require.Equal(t, uint(0), receipt.TransactionIndex)

		blockNumber, err := client.BlockNumber(ctx)
		require.NoError(t, err)
		log.Infof("\nBlock num %d", blockNumber)
		require.GreaterOrEqual(t, blockNumber, receipt.BlockNumber.Uint64())

		block, err := client.BlockByNumber(ctx, receipt.BlockNumber)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, receipt.BlockNumber.Uint64(), block.Number().Uint64())
		require.Equal(t, receipt.BlockHash.String(), block.Hash().String())

		block, err = client.BlockByHash(ctx, receipt.BlockHash)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, receipt.BlockNumber.Uint64(), block.Number().Uint64())
		require.Equal(t, receipt.BlockHash.String(), block.Hash().String())

		nonExistentBlockNumber := big.NewInt(0).SetUint64(blockNumber + uint64(1))
		block, err = client.BlockByNumber(ctx, nonExistentBlockNumber)
		require.Error(t, err)
		require.Nil(t, block)

		nonExistentBlockHash := common.HexToHash("0xFFFFFF")
		block, err = client.BlockByHash(ctx, nonExistentBlockHash)
		require.Error(t, err)
		require.Nil(t, block)
		// its pending

		response, err := jsonrpc.JSONRPCCall(network.URL, "eth_getBlockTransactionCountByNumber", hexutil.EncodeBig(receipt.BlockNumber))
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		txCount := ""
		err = json.Unmarshal(response.Result, &txCount)
		require.NoError(t, err)
		require.Equal(t, "0x1", txCount)

		// check if block number is correct
		count, err := client.TransactionCount(ctx, receipt.BlockHash)
		require.NoError(t, err)
		require.Equal(t, uint(1), count)

		tx = nil
		tx, err = client.TransactionInBlock(ctx, receipt.BlockHash, receipt.TransactionIndex)
		require.NoError(t, err)
		require.Equal(t, receipt.TxHash, tx.Hash())

		raw, err := jsonrpc.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", hexutil.EncodeBig(receipt.BlockNumber), "0x0")
		require.NoError(t, err)
		require.Nil(t, raw.Error)
		require.NotNil(t, raw.Result)

		var newTx rpcTx
		err = json.Unmarshal(raw.Result, &newTx)
		require.NoError(t, err)

		raw, err = jsonrpc.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", "0x123", "0x8659")
		require.NoError(t, err)
		require.Nil(t, raw.Error)
		require.NotNil(t, raw.Result)

		var empty rpcTx
		err = json.Unmarshal(raw.Result, &empty)
		require.NoError(t, err)

		// Checks for empty, when the lookup fail we get an empty struct and no errors...
		v := reflect.ValueOf(empty)

		for i := 0; i < v.NumField(); i++ {
			require.Empty(t, v.Field(i).Interface())
		}

		// checks for successful query

		require.Equal(t, hexutil.EncodeBig(receipt.BlockNumber), newTx.BlockNumber)
		require.Equal(t, receipt.BlockHash.String(), newTx.BlockHash)
		require.Equal(t, hexutil.EncodeUint64(tx.Nonce()), newTx.Nonce)
		require.Equal(t, hexutil.EncodeBig(tx.ChainId()), newTx.ChainID)
	}
}
func Test_Transactions(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		client, err := ethclient.Dial(network.URL)
		require.NoError(t, err)
		destination := common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff")

		// Test Case: Successful transfer

		tx, err := createTX(network.URL, network.ChainID, destination, big.NewInt(100000))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		// Setup for test cases

		auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, network.ChainID)
		require.NoError(t, err)

		nonce, err := client.NonceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &destination, Value: big.NewInt(10000)})
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(context.Background())
		require.NoError(t, err)

		// Test Case: TX with invalid nonce

		tx = types.NewTransaction(nonce-1, // Nonce will be lower than the current getNonceAt()
			destination, big.NewInt(100), gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)

		log.Infof("Sending Tx %v Nonce (invalid) %v", signedTx.Hash(), signedTx.Nonce())
		err = client.SendTransaction(context.Background(), signedTx)
		require.ErrorContains(t, err, "nonce too low")

		// End Test Case

		// Test Case: TX with no signature (which would fail the EIP-155)

		invalidTx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			Value:    big.NewInt(10000),
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     nil,
		})
		err = client.SendTransaction(context.Background(), invalidTx)
		require.Error(t, err)
		// End Test Case

		// Test Case: TX with amount being higher than balance

		balance, err := client.BalanceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		nonce, err = client.NonceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		log.Infof("Balance: %d", balance)

		tx = types.NewTransaction(nonce, destination, big.NewInt(0).Add(balance, big.NewInt(10)), gasLimit, gasPrice, nil)
		signedTx, err = auth.Signer(auth.From, tx)
		require.NoError(t, err)

		log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
		err = client.SendTransaction(context.Background(), signedTx)
		require.ErrorContains(t, err, pool.ErrInsufficientFunds.Error())
	}
}

func Test_Misc(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	Setup()
	defer Teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		client, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		// ChainId()
		chainId, err := client.ChainID(ctx)
		require.NoError(t, err)
		require.Equal(t, network.ChainID, chainId.Uint64())

		// Syncing()
		progress, err := client.SyncProgress(ctx)
		require.NoError(t, err)
		if progress != nil {
			log.Info("Its syncing")
			blockNumber, err := client.BlockNumber(ctx)
			require.NoError(t, err)
			// if it's actually syncing
			require.Equal(t, uint64(0x0), progress.StartingBlock)
			require.Equal(t, blockNumber, progress.CurrentBlock)
		}

		// GetStorageAt()

		// first deploy sample smart contract
		sc_payload := int64(42)
		sc_retrieve := common.HexToHash("0x2a")
		auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, network.ChainID)
		require.NoError(t, err)
		contractAddress, tx, storageSC, err := Storage.DeployStorage(auth, client)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)
		tx, err = storageSC.Store(auth, big.NewInt(sc_payload))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, client, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		storage, err := client.StorageAt(ctx, contractAddress, common.HexToHash("0x0"), nil)
		require.NoError(t, err)
		// in bytes but has to be hash 0x0...42
		require.Equal(t, sc_retrieve, common.BytesToHash(storage))

		// eth_getCode

		scBytecode, err := client.CodeAt(ctx, contractAddress, nil)
		require.NoError(t, err)
		require.Contains(t, Storage.StorageMetaData.Bin, common.Bytes2Hex(scBytecode))

		emptyBytecode, err := client.CodeAt(ctx, common.HexToAddress("0xdeadbeef"), nil)
		require.NoError(t, err)
		require.Empty(t, emptyBytecode)
	}
}
