package e2e

import (
	"context"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Double"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJSONRPC tests JSON RPC methods on a running environment.
func TestJSONRPC(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()
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
	ctx := context.Background()
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

func Test_Filters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	setup()
	defer teardown()
	for _, network := range networks {
		// test newBlockFilter creation
		log.Infof("Network %s", network.Name)
		response, err := client.JSONRPCCall(network.URL, "eth_newBlockFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var filterId string
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		// test newFilter creation with block range and block hash
		response, err = client.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"BlockHash": common.HexToHash("0x1"),
			"FromBlock": "0x1",
			"ToBlock":   "0x2",
		})
		require.NoError(t, err)
		require.NotNil(t, response.Error)
		require.Equal(t, invalidParamsErrorCode, response.Error.Code)
		require.Equal(t, "invalid argument 0: cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other", response.Error.Message)

		// test newFilter creation with block hash
		response, err = client.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
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

		// test newFilter creation with block range
		response, err = client.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
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

		// test uninstallFilter when filter is installed
		response, err = client.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var uninstalled bool
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.True(t, uninstalled)

		// test uninstallFilter when filter doesn't exist or was already uninstalled
		response, err = client.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		uninstalled = true
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.False(t, uninstalled)

		ethereumClient := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		// test getFilterChanges for a blockFilter ID
		blockBeforeFilter, err := ethereumClient.BlockByNumber(ctx, nil)
		require.NoError(t, err)

		response, err = client.JSONRPCCall(network.URL, "eth_newBlockFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var blockFilterId string
		err = json.Unmarshal(response.Result, &blockFilterId)
		require.NoError(t, err)
		require.NotEmpty(t, blockFilterId)

		// force a block to be generated sending a eth transfer tx
		tx, err := createTX(ethereumClient, auth, toAddress, big.NewInt(1000))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		blockAfterFilter, err := ethereumClient.BlockByNumber(ctx, nil)
		require.NoError(t, err)

		response, err = client.JSONRPCCall(network.URL, "eth_getFilterChanges", blockFilterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var blockFilterChanges []common.Hash
		err = json.Unmarshal(response.Result, &blockFilterChanges)
		require.NoError(t, err)

		assert.NotEqual(t, blockBeforeFilter.Hash().String(), blockFilterChanges[0].String())
		assert.Equal(t, blockAfterFilter.Hash().String(), blockFilterChanges[len(blockFilterChanges)-1].String())

		// test getFilterChanges for a logFilter ID
		// create a SC to emit some logs
		scAddr, scTx, sc, err := EmitLog.DeployEmitLog(auth, ethereumClient)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, scTx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		response, err = client.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"Addresses": []common.Address{scAddr},
		})
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		logFilterId := ""
		err = json.Unmarshal(response.Result, &logFilterId)
		require.NoError(t, err)
		require.NotEmpty(t, logFilterId)

		// emit logs
		tx, err = sc.EmitLogs(auth)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		logs, err := ethereumClient.FilterLogs(ctx, ethereum.FilterQuery{Addresses: []common.Address{scAddr}})
		require.NoError(t, err)

		response, err = client.JSONRPCCall(network.URL, "eth_getFilterChanges", logFilterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var logFilterChanges []ethTypes.Log
		err = json.Unmarshal(response.Result, &logFilterChanges)
		require.NoError(t, err)

		assert.Equal(t, 10, len(logs))
		assert.Equal(t, 10, len(logFilterChanges))
		assert.True(t, reflect.DeepEqual(logs, logFilterChanges))

		// emit more logs
		tx, err = sc.EmitLogs(auth)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		tx, err = sc.EmitLogs(auth)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		logs, err = ethereumClient.FilterLogs(ctx, ethereum.FilterQuery{Addresses: []common.Address{scAddr}})
		require.NoError(t, err)

		response, err = client.JSONRPCCall(network.URL, "eth_getFilterChanges", logFilterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		err = json.Unmarshal(response.Result, &logFilterChanges)
		require.NoError(t, err)

		assert.Equal(t, 30, len(logs))
		assert.Equal(t, 20, len(logFilterChanges))
	}
}

func Test_Gas(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()
	var Address1 = toAddress
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
	ctx := context.Background()
	setup()
	defer teardown()
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
		ethereumClient, err := ethclient.Dial(network.URL)
		require.NoError(t, err)
		auth, err := operations.GetAuth(network.PrivateKey, network.ChainID)
		require.NoError(t, err)

		tx, err := createTX(ethereumClient, auth, toAddress, big.NewInt(1000))
		require.NoError(t, err)
		// no block number yet... will wait
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := ethereumClient.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)
		require.Equal(t, receipt.TxHash, tx.Hash())
		require.Equal(t, receipt.Type, tx.Type())
		require.Equal(t, uint(0), receipt.TransactionIndex)

		blockNumber, err := ethereumClient.BlockNumber(ctx)
		require.NoError(t, err)
		log.Infof("\nBlock num %d", blockNumber)
		require.GreaterOrEqual(t, blockNumber, receipt.BlockNumber.Uint64())

		block, err := ethereumClient.BlockByNumber(ctx, receipt.BlockNumber)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, receipt.BlockNumber.Uint64(), block.Number().Uint64())
		require.Equal(t, receipt.BlockHash.String(), block.Hash().String())

		block, err = ethereumClient.BlockByHash(ctx, receipt.BlockHash)
		require.NoError(t, err)
		require.NotNil(t, block)
		require.Equal(t, receipt.BlockNumber.Uint64(), block.Number().Uint64())
		require.Equal(t, receipt.BlockHash.String(), block.Hash().String())

		nonExistentBlockNumber := big.NewInt(0).SetUint64(blockNumber + uint64(1))
		block, err = ethereumClient.BlockByNumber(ctx, nonExistentBlockNumber)
		require.Error(t, err)
		require.Nil(t, block)

		nonExistentBlockHash := common.HexToHash("0xFFFFFF")
		block, err = ethereumClient.BlockByHash(ctx, nonExistentBlockHash)
		require.Error(t, err)
		require.Nil(t, block)
		// its pending

		response, err := client.JSONRPCCall(network.URL, "eth_getBlockTransactionCountByNumber", hexutil.EncodeBig(receipt.BlockNumber))
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		txCount := ""
		err = json.Unmarshal(response.Result, &txCount)
		require.NoError(t, err)
		require.Equal(t, "0x1", txCount)

		// check if block number is correct
		count, err := ethereumClient.TransactionCount(ctx, receipt.BlockHash)
		require.NoError(t, err)
		require.Equal(t, uint(1), count)

		tx = nil
		tx, err = ethereumClient.TransactionInBlock(ctx, receipt.BlockHash, receipt.TransactionIndex)
		require.NoError(t, err)
		require.Equal(t, receipt.TxHash, tx.Hash())

		raw, err := client.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", hexutil.EncodeBig(receipt.BlockNumber), "0x0")
		require.NoError(t, err)
		require.Nil(t, raw.Error)
		require.NotNil(t, raw.Result)

		var newTx rpcTx
		err = json.Unmarshal(raw.Result, &newTx)
		require.NoError(t, err)

		raw, err = client.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", "0x123", "0x8659")
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
	ctx := context.Background()
	setup()
	defer teardown()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		ethClient, err := ethclient.Dial(network.URL)
		require.NoError(t, err)
		auth, err := operations.GetAuth(network.PrivateKey, network.ChainID)
		require.NoError(t, err)

		// Test Case: Successful transfer
		tx, err := createTX(ethClient, auth, toAddress, big.NewInt(100000))
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		// Test Case: get transaction by block number and index
		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)
		require.NotNil(t, receipt)
		res, err := client.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", hex.EncodeBig(receipt.BlockNumber), hex.EncodeUint64(uint64(receipt.TransactionIndex)))
		require.NoError(t, err)
		require.Nil(t, res.Error)
		require.NotNil(t, res.Result)
		var txByBlockNumberAndIndex *types.Transaction
		err = json.Unmarshal(res.Result, &txByBlockNumberAndIndex)
		require.NoError(t, err)

		require.Equal(t, tx.Hash().String(), txByBlockNumberAndIndex.Hash.String())

		// Test Case: get transaction by block hash and index
		receipt, err = ethClient.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)
		require.NotNil(t, receipt)
		txByBlockHashAndIndex, err := ethClient.TransactionInBlock(ctx, receipt.BlockHash, receipt.TransactionIndex)
		require.NoError(t, err)
		require.Equal(t, tx.Hash().String(), txByBlockHashAndIndex.Hash().String())

		// Setup for test cases
		nonce, err := ethClient.NonceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &toAddress, Value: big.NewInt(10000)})
		require.NoError(t, err)

		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		require.NoError(t, err)

		// Test Case: TX with invalid nonce
		tx = ethTypes.NewTransaction(nonce-1, // Nonce will be lower than the current getNonceAt()
			toAddress, big.NewInt(100), gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)

		log.Infof("Sending Tx %v Nonce (invalid) %v", signedTx.Hash(), signedTx.Nonce())
		err = ethClient.SendTransaction(context.Background(), signedTx)
		require.ErrorContains(t, err, "nonce too low")
		// End Test Case

		// Test Case: TX with no signature (which would fail the EIP-155)
		invalidTx := ethTypes.NewTx(&ethTypes.LegacyTx{
			Nonce:    nonce,
			Value:    big.NewInt(10000),
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     nil,
		})
		err = ethClient.SendTransaction(context.Background(), invalidTx)
		require.Error(t, err)
		// End Test Case

		// Test Case: TX with amount being higher than balance
		balance, err := ethClient.BalanceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		nonce, err = ethClient.NonceAt(context.Background(), auth.From, nil)
		require.NoError(t, err)

		log.Infof("Balance: %d", balance)

		tx = ethTypes.NewTransaction(nonce, toAddress, big.NewInt(0).Add(balance, big.NewInt(10)), gasLimit, gasPrice, nil)
		signedTx, err = auth.Signer(auth.From, tx)
		require.NoError(t, err)

		log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
		err = ethClient.SendTransaction(context.Background(), signedTx)
		require.ErrorContains(t, err, pool.ErrInsufficientFunds.Error())

		// no contract code at given address test
		// deploy contract with not enough gas for storage, just execution
		address := common.HexToAddress("0xDEADBEEF596a836C9063a7EE35dA94DDA3b57B62")
		instance, err := Double.NewDouble(address, ethClient)
		require.NoError(t, err)

		callOpts := &bind.CallOpts{Pending: false}

		payload := big.NewInt(5)
		_, err = instance.Double(callOpts, payload)
		require.ErrorContains(t, err, "no contract code at given address")
	}
}
