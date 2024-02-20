package etrog_test

import (
	"context"
	"errors"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/ci/vectors"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	test "github.com/0xPolygonHermez/zkevm-node/state/test/forkid_common"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

const (
	testsFolder = "../../../test/vectors/src/etrog/"
)

var (
	testState *state.State
	forkID    = uint64(state.FORKID_ETROG)
	stateCfg  = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
)

func TestMain(m *testing.M) {
	testState = test.InitTestState(stateCfg)
	defer test.CloseTestState()
	result := m.Run()
	os.Exit(result)
}

// TestStateTransition tests state using test vectors
func TestStateTransition(t *testing.T) {
	ctx := context.Background()

	// Get all tests vector in the etrog folder
	files, err := os.ReadDir(testsFolder)
	require.NoError(t, err)

	for _, file := range files {
		// Load test vectors
		testCases, err := vectors.LoadStateTransitionTestCasesEtrog(testsFolder + file.Name())
		require.NoError(t, err)

		// Run test cases
		for i, testCase := range testCases {
			block := state.Block{
				BlockNumber: uint64(i + 1),
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}

			genesisActions := vectors.GenerateGenesisActionsEtrog(testCase.Genesis)

			dbTx, err := testState.BeginStateTransaction(ctx)
			require.NoError(t, err)

			stateRoot, err := testState.SetGenesis(ctx, block, state.Genesis{Actions: genesisActions}, metrics.SynchronizerCallerLabel, dbTx)
			require.NoError(t, err)
			require.Equal(t, testCase.ExpectedOldStateRoot, stateRoot.String())
			err = dbTx.Rollback(ctx)
			require.NoError(t, err)

			// convert vector txs
			txs := make([]state.L2TxRaw, 0, len(testCase.Txs))
			for i := 0; i < len(testCase.Txs); i++ {
				vecTx := testCase.Txs[i]
				if vecTx.Type != 0x0b {
					tx, err := state.DecodeTx(vecTx.RawTx)
					require.NoError(t, err)
					l2Tx := state.L2TxRaw{
						Tx:                   *tx,
						EfficiencyPercentage: 255,
					}
					txs = append(txs, l2Tx)
				}
			}

			timestampLimit, ok := big.NewInt(0).SetString(testCase.TimestampLimit, 10)
			require.True(t, ok)

			if len(txs) > 0 {
				// Generate batchdata from the txs in the test and compared with the vector
				l2block := state.L2BlockRaw{
					DeltaTimestamp:  uint32(timestampLimit.Uint64()),
					IndexL1InfoTree: testCase.Txs[0].IndexL1InfoTree,
					Transactions:    txs,
				}

				batch := state.BatchRawV2{
					Blocks: []state.L2BlockRaw{l2block},
				}

				batchData, err := state.EncodeBatchV2(&batch)
				require.NoError(t, err)

				require.Equal(t, common.FromHex(testCase.BatchL2Data), batchData)
			}

			processRequest := state.ProcessRequest{
				BatchNumber:             uint64(i + 1),
				L1InfoRoot_V2:           common.HexToHash(testCase.L1InfoRoot),
				OldStateRoot:            stateRoot,
				OldAccInputHash:         common.HexToHash(testCase.OldAccInputHash),
				Transactions:            common.FromHex(testCase.BatchL2Data),
				TimestampLimit_V2:       timestampLimit.Uint64(),
				Coinbase:                common.HexToAddress(testCase.SequencerAddress),
				ForkID:                  testCase.ForkID,
				SkipVerifyL1InfoRoot_V2: testCase.L1InfoTree.SkipVerifyL1InfoRoot,
			}

			processResponse, _ := testState.ProcessBatchV2(ctx, processRequest, true)
			require.Nil(t, processResponse.ExecutorError)
			require.Equal(t, testCase.ExpectedNewStateRoot, processResponse.NewStateRoot.String())
		}
	}
}

func TestUnsupportedPrecompile(t *testing.T) {
	ctx := context.Background()
	var chainIDSequencer = new(big.Int).SetUint64(stateCfg.ChainID)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var sequencerBalance = 4000000
	scByteCode, err := testutils.ReadBytecode("customModExp/customModExp.bin")
	require.NoError(t, err)

	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
		{
			Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000000",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	nonce := uint64(0)

	// Deploy contract
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Call SC method
	nonce++
	tx1 := types.NewTransaction(nonce, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("d5665d6f000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"))
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	l2block := state.L2BlockRaw{
		DeltaTimestamp:  3,
		IndexL1InfoTree: 0,
		Transactions:    []state.L2TxRaw{{Tx: *signedTx0, EfficiencyPercentage: 255}},
	}

	batch := state.BatchRawV2{
		Blocks: []state.L2BlockRaw{l2block},
	}

	batchData, err := state.EncodeBatchV2(&batch)
	require.NoError(t, err)

	processRequest := state.ProcessRequest{
		BatchNumber:             1,
		L1InfoRoot_V2:           common.Hash{},
		OldStateRoot:            stateRoot,
		OldAccInputHash:         common.Hash{},
		Transactions:            batchData,
		TimestampLimit_V2:       3,
		Coinbase:                sequencerAddress,
		ForkID:                  forkID,
		SkipVerifyL1InfoRoot_V2: true,
	}

	processResponse, _ := testState.ProcessBatchV2(ctx, processRequest, true)
	require.Nil(t, processResponse.ExecutorError)
	require.NoError(t, err)

	// Call SC method
	l2block = state.L2BlockRaw{
		DeltaTimestamp:  3,
		IndexL1InfoTree: 0,
		Transactions:    []state.L2TxRaw{{Tx: *signedTx1, EfficiencyPercentage: 255}},
	}

	batch = state.BatchRawV2{
		Blocks: []state.L2BlockRaw{l2block},
	}

	batchData, err = state.EncodeBatchV2(&batch)
	require.NoError(t, err)

	processRequest = state.ProcessRequest{
		BatchNumber:             2,
		L1InfoRoot_V2:           common.Hash{},
		OldStateRoot:            processResponse.NewStateRoot,
		OldAccInputHash:         common.Hash{},
		Transactions:            batchData,
		TimestampLimit_V2:       6,
		Coinbase:                sequencerAddress,
		ForkID:                  forkID,
		SkipVerifyL1InfoRoot_V2: true,
	}

	processResponse, err = testState.ProcessBatchV2(ctx, processRequest, true)
	require.Error(t, err)
	require.Nil(t, processResponse)
	require.True(t, errors.Is(err, runtime.ErrExecutorErrorUnsupportedPrecompile))
}
