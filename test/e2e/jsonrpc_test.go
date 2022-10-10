package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	invalidParamsErrorCode = -32602
)

func TestMain(t *testing.T) {
	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)
}

// TestJSONRPC tests JSON RPC methods on a running environment.
func TestJSONRPC(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	opsCfg := &operations.Config{
		State:     &state.Config{MaxCumulativeGasUsed: operations.DefaultMaxCumulativeGasUsed},
		Sequencer: &operations.SequencerConfig{Address: operations.DefaultSequencerAddress, PrivateKey: operations.DefaultSequencerPrivateKey},
	}
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	sequencerBalance := new(big.Int).SetInt64(int64(operations.DefaultSequencerBalance))

	genesisAccounts := make(map[string]big.Int)
	genesisAccounts[operations.DefaultSequencerAddress] = *sequencerBalance
	require.NoError(t, opsman.SetGenesis(genesisAccounts))

	require.NoError(t, opsman.Setup())

	require.NoError(t, deployContracts(opsman))

	tcs := []struct {
		description, input, expectedOutput string
		expectedErr                        bool
		expectedErrMsg                     string
	}{
		{
			description:    "eth_call, calling double(int256) with data 5",
			input:          `{"jsonrpc":"2.0", "method":"eth_call", "params":[{"from": "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D", "to": "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", "data": "0x6ffa1caa0000000000000000000000000000000000000000000000000000000000000005"}, "latest"], "id":1}`,
			expectedOutput: `{"jsonrpc":"2.0","id":1,"result":"0x000000000000000000000000000000000000000000000000000000000000000a"}`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			actualOutput, err := httpQuery(tc.input)
			if err := checkError(err, tc.expectedErr, tc.expectedErrMsg); err != nil {
				t.Fatalf(err.Error())
			}

			if actualOutput != tc.expectedOutput {
				t.Fatalf("Query return value did not match expectation, got %q, want %q", actualOutput, tc.expectedOutput)
			}
		})
	}
}

func httpQuery(payload string) (string, error) {
	const target = "http://localhost:8123"

	var jsonStr = []byte(payload)
	req, err := http.NewRequest(
		"POST", target,
		bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer func() {
			err = res.Body.Close()
		}()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func checkError(err error, expected bool, msg string) error {
	if !expected && err != nil {
		return fmt.Errorf("Unexpected error %v", err)
	}
	if !expected {
		return nil
	}
	if err == nil {
		return fmt.Errorf("Expected error didn't happen")
	}
	if msg == "" {
		return fmt.Errorf("Expected error message not defined")
	}
	if !strings.HasPrefix(err.Error(), msg) {
		return fmt.Errorf("Wrong error, expected %q, got %q", msg, err.Error())
	}
	return nil
}

func deployContracts(opsman *operations.Manager) error {
	panic("not implemented yet")
	// var txs []*types.Transaction

	// bytecode, err := testutils.ReadBytecode("Double/Double.bin")
	// if err != nil {
	// 	return err
	// }
	// tx0 := types.NewTx(&types.LegacyTx{
	// 	Nonce:    0,
	// 	To:       nil,
	// 	Value:    new(big.Int),
	// 	Gas:      uint64(defaultSequencerBalance),
	// 	GasPrice: new(big.Int).SetUint64(1),
	// 	Data:     common.Hex2Bytes(bytecode),
	// })

	// auth, err := operations.GetAuth(
	// 	defaultSequencerPrivateKey,
	// 	new(big.Int).SetInt64(defaultSequencerChainID))
	// if err != nil {
	// 	return err
	// }
	// signedTx0, err := auth.Signer(auth.From, tx0)
	// if err != nil {
	// 	return err
	// }
	// txs = append(txs, signedTx0)

	// // Create Batch
	// sequencerAddress := common.HexToAddress(defaultSequencerAddress)
	// batch := &state.Batch{
	// 	BlockNumber:        uint64(0),
	// 	Sequencer:          sequencerAddress,
	// 	Aggregator:         sequencerAddress,
	// 	ConsolidatedTxHash: common.Hash{},
	// 	Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
	// 	Uncles:             nil,
	// 	Transactions:       txs,
	// 	RawTxsData:         nil,
	// 	MaticCollateral:    big.NewInt(1),
	// 	ReceivedAt:         time.Now(),
	// 	ChainID:            big.NewInt(defaultSequencerChainID),
	// 	GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	// }

	// st := opsman.State()
	// ctx := context.Background()

	// lastVirtualBatch, err := st.GetLastBatch(ctx, true, "")
	// if err != nil {
	// 	return err
	// }

	// bp, err := st.NewBatchProcessor(ctx, sequencerAddress, lastVirtualBatch.Header.Root[:], "")
	// if err != nil {
	// 	return err
	// }

	// return bp.ProcessBatch(ctx, batch)
}

func createTX(ethdeployment string, to common.Address, amount *big.Int) (*types.Transaction, error) {
	client, err := ethclient.Dial(ethdeployment)
	if err != nil {
		return nil, err
	}
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
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
		TestEmitLog2(t)

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
		fmt.Printf("\nHow many logs : %d\n", len(logs))
		if len(logs) >= 2 {
			require.GreaterOrEqual(t, logs[1].BlockNumber, logs[0].BlockNumber)
		}
	}
}

func Test_Gas(t *testing.T) {
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
	ctx := context.Background()
	for _, network := range networks {
		log.Infof("Network %s", network.Name)
		client, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		tx, err := createTX(network.URL, common.HexToAddress("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"), big.NewInt(1000))
		require.NoError(t, err)
		// no block number yet... will wait
		err = operations.WaitTxToBeMined(client, tx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		require.NoError(t, err)

		blockNumber, err := client.BlockNumber(ctx)
		require.NoError(t, err)
		log.Infof("\nBlock num %d", blockNumber)
		require.GreaterOrEqual(t, blockNumber, receipt.BlockNumber.Uint64())

		blockHash, err := client.BlockByNumber(ctx, big.NewInt(0))
		require.NotNil(t, blockHash)
		require.NoError(t, err)

		blockHash, err = client.BlockByHash(ctx, common.HexToHash("0x0"))
		require.Nil(t, blockHash)
		require.Error(t, err)

		blockHash, err = client.BlockByHash(ctx, common.HexToHash("0x2"))
		require.Nil(t, blockHash)
		require.Error(t, err)

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
		require.Equal(t, uint(0x1), count)

		tx, err = client.TransactionInBlock(ctx, receipt.BlockHash, receipt.TransactionIndex)
		require.NoError(t, err)
		require.Equal(t, tx.Hash(), receipt.TxHash)

		raw, err := jsonrpc.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", hexutil.EncodeBig(receipt.BlockNumber), "0x0")
		require.NoError(t, err)
		require.Nil(t, raw.Error)
		require.NotNil(t, raw.Result)

		var newTx rpcTx
		err = json.Unmarshal(raw.Result, &newTx)
		require.NoError(t, err)

		raw, err = jsonrpc.JSONRPCCall(network.URL, "eth_getTransactionByBlockNumberAndIndex", "0x123", "0x865")
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

}

func Test_Misc(t *testing.T) {

}
