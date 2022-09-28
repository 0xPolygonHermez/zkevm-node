package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

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

func createTX(ethdeployment string, to common.Address, amount *big.Int) error {
	client, err := ethclient.Dial(ethdeployment)
	if err != nil {
		return err
	}
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	if err != nil {
		return err
	}
	ctx := context.Background()
	nonce, err := client.NonceAt(ctx, auth.From, nil)
	if err != nil {
		return err
	}
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &to, Value: amount})
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce+uint64(rand.Intn(1000))/2, to, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}
	log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	return nil
}
func Test_Filters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	receiverAddr := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	for _, network := range networks {
		log.Debugf(network.Name)

		const (
			CREATE    = 0
			UNINSTALL = 1
			GET       = 2
		)
		ids := map[string]string{
			"newBlockFilter":              "",
			"newPendingTransactionFilter": "",
			"newFilter":                   "",
		}
		testcases := []struct {
			name             string
			action           int
			description      string
			payload          struct{}
			expectedResponse bool
			id               string
			shouldErr        bool
			err              error
		}{
			{
				name:        "eth_newBlockFilter",
				action:      CREATE,
				shouldErr:   false,
				description: "create new block filter successfully",
			},
			{
				name:        "eth_newPendingTransactionFilter",
				action:      CREATE,
				shouldErr:   false,
				description: "create new pending TX filter successfully",
			},
			{

				name:             "eth_newFilter",
				action:           CREATE,
				shouldErr:        false,
				description:      "create a filter successfully",
				id:               "newFilter",
				expectedResponse: true,
			},
			{
				name:             "eth_uninstallFilter",
				action:           UNINSTALL,
				shouldErr:        false,
				description:      "uninstall a filter successfully",
				id:               "newFilter",
				expectedResponse: true,
			},
			{
				name:             "eth_uninstallFilter",
				action:           UNINSTALL,
				shouldErr:        false,
				description:      "uninstall a filter with error because it was already deleted",
				id:               "newFilter",
				expectedResponse: false,
			},
			{

				name:             "eth_newFilter",
				action:           CREATE,
				shouldErr:        false,
				description:      "create a filter successfully again",
				id:               "newFilter",
				expectedResponse: true,
			},
			{
				name:             "eth_getFilterLogs",
				action:           GET,
				shouldErr:        false,
				description:      "get the data from the filter",
				id:               "newFilter",
				expectedResponse: false,
			},
		}

		for _, tc := range testcases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				var response jsonrpc.Response
				var err error
				if tc.action == UNINSTALL || tc.action == GET {
					response, err = jsonrpc.JSONRPCCall(network.URL, tc.name, ids[strings.Replace(tc.id, "eth_", "", -1)])
					require.NoError(t, err)
				} else if tc.action == CREATE && tc.name == "eth_newFilter" {
					response, err = jsonrpc.JSONRPCCall(network.URL, tc.name, &jsonrpc.LogFilter{
						Addresses: []common.Address{
							receiverAddr,
						},
					})
					require.NoError(t, err)

					// create a TX so the filter has some data
					err = createTX(
						network.URL,
						receiverAddr,
						big.NewInt(1),
					)
					require.NoError(t, err)
				} else {
					response, err = jsonrpc.JSONRPCCall(network.URL, tc.name)
					require.NoError(t, err)
				}

				if tc.shouldErr {
					require.NotNil(t, response.Error)
					require.Equal(t, tc.err.Error(), response.Error.Message)

				} else {
					require.Nil(t, response.Error)
					var result interface{}
					err = json.Unmarshal(response.Result, &result)
					require.NoError(t, err)

					require.NotNil(t, result)

					if tc.action == CREATE {
						fmt.Println("Result " + result.(string))
						ids[strings.Replace(tc.id, "eth_", "", -1)] = result.(string)
					}
					if tc.action == UNINSTALL {
						var result bool
						err = json.Unmarshal(response.Result, &result)
						require.NoError(t, err)
						require.NotNil(t, result)
						require.Equal(t, tc.expectedResponse, result)
					}
					if tc.action == GET {
						type filterResponse struct {
							address, data, blockNumber, transactionHash, transactionIndex, blockHash, logIndex string
							topics                                                                             []string
						}
						var result []filterResponse
						err = json.Unmarshal(response.Result, &result)
						require.NoError(t, err)
						require.NotNil(t, result)
						fmt.Printf("\n%+v", result)
					}

				}

				require.Equal(t, float64(1), response.ID)
				require.Equal(t, "2.0", response.JSONRPC)
			})
		}
	}
}

func Test_Filters2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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

	for _, network := range networks {
		response, err := jsonrpc.JSONRPCCall(network.URL, "eth_newBlockFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var filterId string
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		////////////////

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newPendingTransactionFilter")
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		filterId = ""
		err = json.Unmarshal(response.Result, &filterId)
		require.NoError(t, err)
		require.NotEmpty(t, filterId)

		////////////////

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_newFilter", map[string]interface{}{
			"BlockHash": common.HexToHash("0x1"),
			"FromBlock": "0x1",
			"ToBlock":   "0x2",
		})
		require.NoError(t, err)
		require.NotNil(t, response.Error)
		require.Equal(t, -32602, response.Error.Code)
		require.Equal(t, "invalid argument 0: cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other", response.Error.Message)

		////////////////

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

		////////////////

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

		////////////////

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		var uninstalled bool
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.True(t, uninstalled)

		////////////////

		response, err = jsonrpc.JSONRPCCall(network.URL, "eth_uninstallFilter", filterId)
		require.NoError(t, err)
		require.Nil(t, response.Error)
		require.NotNil(t, response.Result)

		uninstalled = true
		err = json.Unmarshal(response.Result, &uninstalled)
		require.NoError(t, err)
		require.False(t, uninstalled)
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

	/*
		func (e *Eth) BlockNumber() (interface{}, rpcError) {

		func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, rpcError) {
		func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, rpcError) {

		func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, rpcError) {
		func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, rpcError) {

		func (e *Eth) GetBlockTransactionCountByHash(hash common.Hash) (interface{}, rpcError) {
		func (e *Eth) GetBlockTransactionCountByNumber(number *BlockNumber) (interface{}, rpcError) {

		func (e *Eth) GetUncleByBlockHashAndIndex() (interface{}, rpcError) {
		func (e *Eth) GetUncleByBlockNumberAndIndex() (interface{}, rpcError) {
		func (e *Eth) GetUncleCountByBlockHash() (interface{}, rpcError) {
		func (e *Eth) GetUncleCountByBlockNumber() (interface{}, rpcError) {
	*/
}

func Test_Transactions(t *testing.T) {

}

func Test_Misc(t *testing.T) {

}
