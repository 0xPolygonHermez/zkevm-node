package synchronizer

import (
	context "context"
	"math/big"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	Etherman *ethermanMock
	State    *stateMock
	DbTx     *dbTxMock
}

func TestTrustedStateReorg(t *testing.T) {
	type testCase struct {
		Name            string
		getTrustedBatch func(*mocks, context.Context, etherman.SequencedBatch) *state.Batch
	}

	setupMocks := func(m *mocks, tc *testCase) Synchronizer {
		genesis := state.Genesis{}
		cfg := Config{
			SyncInterval:   cfgTypes.Duration{Duration: 1 * time.Second},
			SyncChunkSize:  10,
			GenBlockNumber: uint64(123456),
		}
		sync, err := NewSynchronizer(true, m.Etherman, m.State, genesis, cfg)
		require.NoError(t, err)

		// state preparation
		ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
		m.State.
			On("BeginStateTransaction", ctxMatchBy).
			Run(func(args mock.Arguments) {
				ctx := args[0].(context.Context)
				parentHash := common.HexToHash("0x111")
				ethHeader := &types.Header{Number: big.NewInt(1), ParentHash: parentHash}
				ethBlock := types.NewBlockWithHeader(ethHeader)
				lastBlock := &state.Block{BlockHash: ethBlock.Hash(), BlockNumber: ethBlock.Number().Uint64()}

				m.State.
					On("GetLastBlock", ctx, m.DbTx).
					Return(lastBlock, nil).
					Once()

				m.DbTx.
					On("Commit", ctx).
					Return(nil).
					Once()

				m.Etherman.
					On("EthBlockByNumber", ctx, lastBlock.BlockNumber).
					Return(ethBlock, nil).
					Once()

				var n *big.Int
				m.Etherman.
					On("HeaderByNumber", ctx, n).
					Return(ethHeader, nil).
					Once()

				sequencedBatch := etherman.SequencedBatch{
					BatchNumber: uint64(1),
					Coinbase:    common.HexToAddress("0x222"),
					TxHash:      common.HexToHash("0x333"),
					ProofOfEfficiencyBatchData: proofofefficiency.ProofOfEfficiencyBatchData{
						Transactions:          []byte{},
						GlobalExitRoot:        [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
						Timestamp:             uint64(time.Now().Unix()),
						ForceBatchesTimestamp: []uint64{},
					},
				}

				ethermanBlock := etherman.Block{
					BlockHash:        ethBlock.Hash(),
					SequencedBatches: [][]etherman.SequencedBatch{{sequencedBatch}},
				}
				blocks := []etherman.Block{ethermanBlock}
				order := map[common.Hash][]etherman.Order{
					ethBlock.Hash(): {
						{
							Name: etherman.SequenceBatchesOrder,
							Pos:  0,
						},
					},
				}

				fromBlock := ethBlock.NumberU64() + 1
				toBlock := fromBlock + cfg.SyncChunkSize

				m.Etherman.
					On("GetRollupInfoByBlockRange", ctx, fromBlock, &toBlock).
					Return(blocks, order, nil).
					Once()

				m.State.
					On("BeginStateTransaction", ctx).
					Return(m.DbTx, nil).
					Once()

				stateBlock := &state.Block{
					BlockNumber: ethermanBlock.BlockNumber,
					BlockHash:   ethermanBlock.BlockHash,
					ParentHash:  ethermanBlock.ParentHash,
					ReceivedAt:  ethermanBlock.ReceivedAt,
				}

				m.State.
					On("AddBlock", ctx, stateBlock, m.DbTx).
					Return(nil).
					Once()

				trustedBatch := tc.getTrustedBatch(m, ctx, sequencedBatch)

				m.State.
					On("GetBatchByNumber", ctx, sequencedBatch.BatchNumber, m.DbTx).
					Return(trustedBatch, nil).
					Once()

				m.State. //ExecuteBatch(s.ctx, batch.BatchNumber, batch.BatchL2Data, dbTx
						On("ExecuteBatch", ctx, sequencedBatch.BatchNumber, sequencedBatch.Transactions, m.DbTx).
						Return(&pb.ProcessBatchResponse{NewStateRoot: trustedBatch.StateRoot.Bytes()}, nil).
						Once()

				m.State.
					On("ResetTrustedState", ctx, sequencedBatch.BatchNumber-1, m.DbTx).
					Return(nil).
					Once()

				processingContext := state.ProcessingContext{
					BatchNumber:    sequencedBatch.BatchNumber,
					Coinbase:       sequencedBatch.Coinbase,
					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
				}

				m.State.
					On("ProcessAndStoreClosedBatch", ctx, processingContext, sequencedBatch.Transactions, m.DbTx).
					Return(nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					BatchNumber: sequencedBatch.BatchNumber,
					TxHash:      sequencedBatch.TxHash,
					Coinbase:    sequencedBatch.Coinbase,
					BlockNumber: ethermanBlock.BlockNumber,
				}

				m.State.
					On("AddVirtualBatch", ctx, virtualBatch, m.DbTx).
					Return(nil).
					Once()

				m.DbTx.
					On("Commit", ctx).
					Run(func(args mock.Arguments) { sync.Stop() }).
					Return(nil).
					Once()

				m.Etherman.
					On("GetLatestBatchNumber").
					Return(uint64(10), nil).
					Once()

				var nilDbTx pgx.Tx
				m.State.
					On("GetLastBatchNumber", ctx, nilDbTx).
					Return(uint64(10), nil).
					Once()
			}).
			Return(m.DbTx, nil).
			Once()

		return sync
	}

	testCases := []testCase{
		{
			Name: "Transactions are different",
			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
				return &state.Batch{
					BatchL2Data:    []byte{1},
					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
					Coinbase:       sequencedBatch.Coinbase,
				}
			},
		},
		{
			Name: "Global Exit Root is different",
			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
				return &state.Batch{
					BatchL2Data:    sequencedBatch.Transactions,
					GlobalExitRoot: common.HexToHash("0x999888777"),
					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
					Coinbase:       sequencedBatch.Coinbase,
				}
			},
		},
		{
			Name: "Timestamp is different",
			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
				return &state.Batch{
					BatchL2Data:    sequencedBatch.Transactions,
					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
					Timestamp:      time.Unix(int64(0), 0),
					Coinbase:       sequencedBatch.Coinbase,
				}
			},
		},
		{
			Name: "Coinbase is different",
			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
				return &state.Batch{
					BatchL2Data:    sequencedBatch.Transactions,
					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
					Coinbase:       common.HexToAddress("0x999888777"),
				}
			},
		},
	}

	m := mocks{
		Etherman: newEthermanMock(t),
		State:    newStateMock(t),
		DbTx:     newDbTxMock(t),
	}

	// start synchronizing
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testCase := tc
			sync := setupMocks(&m, &testCase)
			err := sync.Sync()
			require.NoError(t, err)
		})
	}
}
