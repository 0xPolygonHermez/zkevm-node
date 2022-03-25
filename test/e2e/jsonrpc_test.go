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
	"github.com/hermeznetwork/hermez-core/test/testutils"
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
	var txs []*types.Transaction

	bytecode, err := testutils.ReadBytecode("Double.bin")
	if err != nil {
		return err
	}
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(defaultSequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(bytecode),
	})

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
		ChainID:            big.NewInt(defaultSequencerChainID),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	st := opsman.State()
	ctx := context.Background()

	lastVirtualBatch, err := st.GetLastBatch(ctx, true)
	if err != nil {
		return err
	}

	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, lastVirtualBatch.Header.Root[:])
	if err != nil {
		return err
	}

	return bp.ProcessBatch(ctx, batch)
}
