package e2e

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/state-transition.json")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			opsCfg := &operations.Config{
				Arity:          testCase.Arity,
				DefaultChainID: testCase.DefaultChainID,
				Sequencer: &operations.SequencerConfig{
					Address:    testCase.SequencerAddress,
					PrivateKey: testCase.SequencerPrivateKey,
					ChainID:    testCase.ChainIDSequencer,
				},
			}
			opsman, err := operations.NewManager(ctx, opsCfg)
			require.NoError(t, err)

			err = opsman.SetGenesis(testCase.GenesisAccounts)
			require.NoError(t, err)

			// Check initial root
			err = opsman.CheckVirtualRoot(testCase.ExpectedOldRoot)
			require.NoError(t, err)

			err = opsman.Setup()
			require.NoError(t, err)

			err = opsman.ApplyTxs(testCase.Txs, testCase.ExpectedOldRoot, testCase.ExpectedNewRoot)
			require.NoError(t, err)

			st := opsman.State()

			// Check leafs
			batchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("addr: %s expected: %s found: %s", addr.Hex(), leaf.Balance.Text(encoding.Base10), actualBalance.Text(encoding.Base10)))

				actualNonce, err := st.GetNonce(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, encoding.Base10), fmt.Sprintf("addr: %s expected: %s found: %d", addr.Hex(), leaf.Nonce, actualNonce))
			}

			// Check state against the expected state
			root, err := st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot := new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot, "Invalid new root")

			// Check consolidated state against the expected state
			err = opsman.CheckVirtualRoot(testCase.ExpectedNewRoot)
			require.NoError(t, err)

			// Check that last virtual and consolidated batch are the same
			lastConsolidatedBatchNumber, err := st.GetLastConsolidatedBatchNumber(ctx)
			require.NoError(t, err)
			lastVirtualBatchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			t.Logf("lastConsolidatedBatchNumber: %d lastVirtualBatchNumber: %d", lastConsolidatedBatchNumber, lastVirtualBatchNumber)
			assert.Equal(t, lastConsolidatedBatchNumber, lastVirtualBatchNumber)

			require.NoError(t, operations.Teardown())
		})
	}
}
