package state_test

/* WIP: This test is not working yet
// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	ctx := context.Background()
	// Load test vectors
	testCases, err := vectors.LoadStateTransitionTestCasesEtrog("./test/vectors/src/state-transition/etrog/balances.json")
	require.NoError(t, err)

	// Run test cases
	for i, testCase := range testCases {
		block := state.Block{
			BlockNumber: 1,
			BlockHash:   state.ZeroHash,
			ParentHash:  state.ZeroHash,
			ReceivedAt:  time.Now(),
		}

		genesisActions := vectors.GenerateGenesisActionsEtrog(testCase.Genesis)

		dbTx, err := testState.BeginStateTransaction(ctx)
		require.NoError(t, err)
		defer dbTx.Commit(ctx)

		stateRoot, err := testState.SetGenesis(ctx, block, state.Genesis{Actions: genesisActions}, metrics.SynchronizerCallerLabel, dbTx)
		require.NoError(t, err)
		require.Equal(t, testCase.ExpectedOldStateRoot, stateRoot.String())

		/*
			// convert vector txs
			txs := make([]*types.Transaction, 0, len(testCase.Txs))
			for i := 0; i < len(testCase.Txs); i++ {
				vecTx := testCase.Txs[i]
				tx, err := state.DecodeTx(vecTx.RawTx)
				require.NoError(t, err)
				txs = append(txs, tx)
			}
*/

/*
		timestampLimit, ok := big.NewInt(0).SetString(testCase.TimestampLimit, 10)
		require.True(t, ok)

		processRequest := state.ProcessRequest{
			BatchNumber:       uint64(i + 1),
			L1InfoRoot_V2:     common.HexToHash(testCase.L1InfoRoot),
			OldStateRoot:      stateRoot,
			OldAccInputHash:   common.HexToHash(testCase.OldAccInputHash),
			Transactions:      common.Hex2Bytes(strings.TrimLeft(testCase.BatchL2Data, "0x")),
			TimestampLimit_V2: timestampLimit.Uint64(),
			ForkID:            testCase.ForkID,
		}

		fmt.Printf("processRequest: %+v\n", processRequest)

		processResponse, err := testState.ProcessBatchV2(ctx, processRequest, false)
		require.NoError(t, err)

		fmt.Printf("processResponse: %+v\n", processResponse)

		require.Nil(t, processResponse.ExecutorError)
		require.Equal(t, testCase.ExpectedNewStateRoot, processResponse.NewStateRoot.String())

		break
	}
}
*/
