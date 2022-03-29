package e2e

import (
	"context"
	"fmt"
	"io/fs"
	"math/big"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
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

	// Load test vectors
	var testCases []vectors.StateTransitionTestCase
	root := filepath.Clean("./../vectors/src/state-transition/no-data")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		tmpStateTransitionTestCases, err := vectors.LoadStateTransitionTestCases(path)
		if err != nil {
			return err
		}

		testCases = append(testCases, tmpStateTransitionTestCases...)
		return nil
	})
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			opsCfg := &operations.Config{
				Arity: testCase.Arity,
				State: &state.Config{
					DefaultChainID:       testCase.DefaultChainID,
					MaxCumulativeGasUsed: 800000,
				},

				Sequencer: &operations.SequencerConfig{
					Address:    testCase.SequencerAddress,
					PrivateKey: testCase.SequencerPrivateKey,
					ChainID:    testCase.ChainIDSequencer,
				},
			}
			opsman, err := operations.NewManager(ctx, opsCfg)
			require.NoError(t, err)

			genesisAccounts := make(map[string]big.Int)
			for _, gacc := range testCase.GenesisAccounts {
				genesisAccounts[gacc.Address] = gacc.Balance.Int
			}
			require.NoError(t, opsman.SetGenesis(genesisAccounts))

			// Check initial root
			require.NoError(t, opsman.CheckVirtualRoot(testCase.ExpectedOldRoot))

			require.NoError(t, opsman.Setup())

			require.NoError(t, opsman.ApplyTxs(testCase.Txs, testCase.ExpectedOldRoot, testCase.ExpectedNewRoot))

			st := opsman.State()

			// Check leafs
			batchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(ctx, addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("addr: %s expected: %s found: %s", addr.Hex(), leaf.Balance.Text(encoding.Base10), actualBalance.Text(encoding.Base10)))

				actualNonce, err := st.GetNonce(ctx, addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, encoding.Base10), fmt.Sprintf("addr: %s expected: %s found: %d", addr.Hex(), leaf.Nonce, actualNonce))
			}

			// Check state against the expected state
			root, err := st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot := hex.EncodeToHex(root)
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot, "Invalid new root")

			// Check consolidated state against the expected state
			require.NoError(t, opsman.CheckVirtualRoot(testCase.ExpectedNewRoot))

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
