package e2e

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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
	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/src/state-transition/no-data/general.json")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			opsCfg := &operations.Config{
				State: &state.Config{
					MaxCumulativeGasUsed: 800000,
				},
				Sequencer: &operations.SequencerConfig{
					Address:    testCase.SequencerAddress,
					PrivateKey: testCase.SequencerPrivateKey,
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

			require.NoError(t, opsman.ApplyTxs(testCase.Txs, testCase.ExpectedOldRoot, testCase.ExpectedNewRoot, testCase.GlobalExitRoot))

			st := opsman.State()

			// Check leafs
			l2BlockNumber, err := st.GetLastL2BlockNumber(ctx, nil)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(ctx, addr, l2BlockNumber, nil)
				require.NoError(t, err)
				require.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("addr: %s expected: %s found: %s", addr.Hex(), leaf.Balance.Text(encoding.Base10), actualBalance.Text(encoding.Base10)))

				actualNonce, err := st.GetNonce(ctx, addr, l2BlockNumber, nil)
				require.NoError(t, err)
				require.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, encoding.Base10), fmt.Sprintf("addr: %s expected: %s found: %d", addr.Hex(), leaf.Nonce, actualNonce))
			}

			// Check virtual root against the expected state
			require.NoError(t, opsman.CheckVirtualRoot(testCase.ExpectedNewRoot))

			// Check that last virtual and consolidated l2 block are the same
			lastConsolidatedL2BlockNumber, err := st.GetLastConsolidatedL2BlockNumber(ctx, nil)
			require.NoError(t, err)
			lastVirtualL2BlockNumber, err := st.GetLastL2BlockNumber(ctx, nil)
			require.NoError(t, err)
			t.Logf("lastConsolidatedL2BlockNumber: %d lastVirtualL2BlockNumber: %d", lastConsolidatedL2BlockNumber, lastVirtualL2BlockNumber)
			require.Equal(t, lastConsolidatedL2BlockNumber, lastVirtualL2BlockNumber)

			require.NoError(t, operations.Teardown())
		})
	}
}
