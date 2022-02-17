package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/stretchr/testify/require"
)

const (
	defaultArity                = 4
	defaultChainID              = 1000
	defaultSequencerAddress     = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	defaultSequencerPrivateKey  = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	defaultSequencerChainID     = 400
	defaultSequencerBalance     = 400000
	defaultMaxCumulativeGasUsed = 800000
)

// TestJSONRPC tests JSON RPC methods on a running environment.
func TestJSONRPC(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	opsCfg := &operations.Config{
		Arity: defaultArity,
		State: &state.Config{
			DefaultChainID:       defaultChainID,
			MaxCumulativeGasUsed: defaultMaxCumulativeGasUsed,
		},
		Sequencer: &operations.SequencerConfig{
			Address:    defaultSequencerAddress,
			PrivateKey: defaultSequencerPrivateKey,
			ChainID:    defaultSequencerChainID,
		},
	}
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	sequencerBalance := new(big.Int).SetInt64(int64(defaultSequencerBalance))

	genesisAccounts := make(map[string]big.Int)
	genesisAccounts[defaultSequencerAddress] = *sequencerBalance
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
			input:          `{"jsonrpc":"2.0", "method":"eth_call", "params":[{"from": "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D", "to": "0xB08EA26d3D53EC62fD4BD76B5E41844c7041eB6B", "data": "0x6ffa1caa0000000000000000000000000000000000000000000000000000000000000005"}, "latest"], "id":1}`,
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
	var txs []*types.Transaction

	/* bytecode for the contract:
	   // SPDX-License-Identifier: MIT
	   pragma solidity ^0.8.0;
	   contract Test {
	     function double(int a) public pure returns(int) {
	       return 2*a;
	     }
	   }
	*/
	tx0 := types.NewTransaction(0, state.ZeroAddress, new(big.Int), uint64(defaultSequencerBalance), new(big.Int).SetUint64(1), common.Hex2Bytes("608060405234801561001057600080fd5b50610284806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80636ffa1caa14610030575b600080fd5b61004a600480360381019061004591906100b1565b610060565b60405161005791906100ed565b60405180910390f35b600081600261006f9190610137565b9050919050565b600080fd5b6000819050919050565b61008e8161007b565b811461009957600080fd5b50565b6000813590506100ab81610085565b92915050565b6000602082840312156100c7576100c6610076565b5b60006100d58482850161009c565b91505092915050565b6100e78161007b565b82525050565b600060208201905061010260008301846100de565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006101428261007b565b915061014d8361007b565b9250827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211600084136000841316161561018c5761018b610108565b5b817f800000000000000000000000000000000000000000000000000000000000000005831260008412600084131616156101c9576101c8610108565b5b827f8000000000000000000000000000000000000000000000000000000000000000058212600084136000841216161561020657610205610108565b5b827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff058212600084126000841216161561024357610242610108565b5b82820290509291505056fea264697066735822122056fe5e7440cf5d75140a0cf9846a77c37d5214ff1c91d7e8f5ed843f2d65df6a64736f6c63430008090033"))

	auth, err := operations.GetAuth(
		defaultSequencerPrivateKey,
		new(big.Int).SetInt64(defaultSequencerChainID))
	if err != nil {
		return err
	}
	signedTx0, err := auth.Signer(auth.From, tx0)
	if err != nil {
		return err
	}
	txs = append(txs, signedTx0)

	// Create Batch
	sequencerAddress := common.HexToAddress(defaultSequencerAddress)
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
	}

	st := opsman.State()
	bp, err := st.NewBatchProcessor(sequencerAddress, 0)
	if err != nil {
		return err
	}

	return bp.ProcessBatch(batch)
}
