package synchronizer

import (
	context "context"
	"math/big"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	Etherman     *ethermanMock
	State        *stateMock
	Pool         *poolMock
	EthTxManager *ethTxManagerMock
	DbTx         *dbTxMock
	ZKEVMClient  *zkEVMClientMock
}

// Test commented until we remove the fatal in checkTrustedReorg function
// func TestTrustedStateReorg(t *testing.T) {
// 	type testCase struct {
// 		Name            string
// 		getTrustedBatch func(*mocks, context.Context, etherman.SequencedBatch) *state.Batch
// 		getTrustedReorg func(m *mocks, batchNumber, timestamp uint64) state.TrustedReorg
// 	}

// 	setupMocks := func(m *mocks, tc *testCase) Synchronizer {
// 		genesis := state.Genesis{}
// 		cfg := Config{
// 			SyncInterval:   cfgTypes.Duration{Duration: 1 * time.Second},
// 			SyncChunkSize:  10,
// 			GenBlockNumber: uint64(123456),
// 		}

// 		sync, err := NewSynchronizer(false, m.Etherman, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, genesis, cfg)
// 		require.NoError(t, err)

// 		// state preparation
// 		ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
// 		m.State.
// 			On("BeginStateTransaction", ctxMatchBy).
// 			Run(func(args mock.Arguments) {
// 				ctx := args[0].(context.Context)
// 				parentHash := common.HexToHash("0x111")
// 				ethHeader := &ethTypes.Header{Number: big.NewInt(1), ParentHash: parentHash}
// 				ethBlock := ethTypes.NewBlockWithHeader(ethHeader)
// 				lastBlock := &state.Block{BlockHash: ethBlock.Hash(), BlockNumber: ethBlock.Number().Uint64()}

// 				m.State.
// 					On("GetLastBlock", ctx, m.DbTx).
// 					Return(lastBlock, nil).
// 					Once()

// 				m.DbTx.
// 					On("Commit", ctx).
// 					Return(nil).
// 					Once()

// 				m.Etherman.
// 					On("EthBlockByNumber", ctx, lastBlock.BlockNumber).
// 					Return(ethBlock, nil).
// 					Once()

// 				var n *big.Int
// 				m.Etherman.
// 					On("HeaderByNumber", ctx, n).
// 					Return(ethHeader, nil).
// 					Once()

// 				t := time.Now()
// 				sequencedBatch := etherman.SequencedBatch{
// 					BatchNumber: uint64(1),
// 					Coinbase:    common.HexToAddress("0x222"),
// 					TxHash:      common.HexToHash("0x333"),
// 					PolygonZkEVMBatchData: polygonzkevm.PolygonZkEVMBatchData{
// 						Transactions:       []byte{},
// 						GlobalExitRoot:     [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
// 						Timestamp:          uint64(t.Unix()),
// 						MinForcedTimestamp: 0,
// 					},
// 				}

// 				ethermanBlock := etherman.Block{
// 					BlockHash:        ethBlock.Hash(),
// 					SequencedBatches: [][]etherman.SequencedBatch{{sequencedBatch}},
// 				}
// 				blocks := []etherman.Block{ethermanBlock}
// 				order := map[common.Hash][]etherman.Order{
// 					ethBlock.Hash(): {
// 						{
// 							Name: etherman.SequenceBatchesOrder,
// 							Pos:  0,
// 						},
// 					},
// 				}

// 				fromBlock := ethBlock.NumberU64() + 1
// 				toBlock := fromBlock + cfg.SyncChunkSize

// 				m.Etherman.
// 					On("GetRollupInfoByBlockRange", ctx, fromBlock, &toBlock).
// 					Return(blocks, order, nil).
// 					Once()

// 				m.ZKEVMClient.
// 					On("BatchNumber", ctx).
// 					Return(uint64(1), nil).
// 					Once()

// 				m.State.
// 					On("BeginStateTransaction", ctx).
// 					Return(m.DbTx, nil).
// 					Once()

// 				stateBlock := &state.Block{
// 					BlockNumber: ethermanBlock.BlockNumber,
// 					BlockHash:   ethermanBlock.BlockHash,
// 					ParentHash:  ethermanBlock.ParentHash,
// 					ReceivedAt:  ethermanBlock.ReceivedAt,
// 				}

// 				m.State.
// 					On("AddBlock", ctx, stateBlock, m.DbTx).
// 					Return(nil).
// 					Once()

// 				trustedBatch := tc.getTrustedBatch(m, ctx, sequencedBatch)

// 				m.State.
// 					On("GetBatchByNumber", ctx, sequencedBatch.BatchNumber, m.DbTx).
// 					Return(trustedBatch, nil).
// 					Once()

// 				sbatch := state.Batch{
// 					BatchNumber:    sequencedBatch.BatchNumber,
// 					Coinbase:       common.HexToAddress("0x222"),
// 					BatchL2Data:    []byte{},
// 					Timestamp:      time.Unix(int64(t.Unix()), 0),
// 					Transactions:   nil,
// 					GlobalExitRoot: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
// 					ForcedBatchNum: nil,
// 				}
// 				m.State.
// 					On("ExecuteBatch", ctx, sbatch, false, m.DbTx).
// 					Return(&pb.ProcessBatchResponse{NewStateRoot: trustedBatch.StateRoot.Bytes()}, nil).
// 					Once()

// 				seq := state.Sequence{
// 					FromBatchNumber: 1,
// 					ToBatchNumber:   1,
// 				}
// 				m.State.
// 					On("AddSequence", ctx, seq, m.DbTx).
// 					Return(nil).
// 					Once()

// 				m.State.
// 					On("AddAccumulatedInputHash", ctx, sequencedBatch.BatchNumber, common.Hash{}, m.DbTx).
// 					Return(nil).
// 					Once()

// 				tr := tc.getTrustedReorg(m, sbatch.BatchNumber, uint64(t.Unix()))
// 				m.State.
// 					On("AddTrustedReorg", ctx, &tr, m.DbTx).
// 					Return(nil).
// 					Once()

// 				m.Etherman.
// 					On("GetLatestBatchNumber").
// 					Return(tr.BatchNumber-1, nil).
// 					Once()

// 				txs := []*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})}
// 				m.State.
// 					On("GetReorgedTransactions", ctx, tr.BatchNumber, m.DbTx).
// 					Return(txs, nil).
// 					Once()

// 				m.Pool.
// 					On("DeleteReorgedTransactions", ctx, txs).
// 					Return(nil).
// 					Once()

// 				m.Pool.
// 					On("StoreTx", ctx, *txs[0], "", true).
// 					Return(nil).
// 					Once()

// 				m.State.
// 					On("ResetTrustedState", ctx, sequencedBatch.BatchNumber-1, m.DbTx).
// 					Return(nil).
// 					Once()

// 				processingContext := state.ProcessingContext{
// 					BatchNumber:    sequencedBatch.BatchNumber,
// 					Coinbase:       sequencedBatch.Coinbase,
// 					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
// 					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
// 				}

// 				m.State.
// 					On("ProcessAndStoreClosedBatch", ctx, processingContext, sequencedBatch.Transactions, m.DbTx, state.SynchronizerCallerLabel).
// 					Return(trustedBatch.StateRoot, nil).
// 					Once()

// 				virtualBatch := &state.VirtualBatch{
// 					BatchNumber: sequencedBatch.BatchNumber,
// 					TxHash:      sequencedBatch.TxHash,
// 					Coinbase:    sequencedBatch.Coinbase,
// 					BlockNumber: ethermanBlock.BlockNumber,
// 				}

// 				m.State.
// 					On("AddVirtualBatch", ctx, virtualBatch, m.DbTx).
// 					Return(nil).
// 					Once()

// 				m.DbTx.
// 					On("Commit", ctx).
// 					Run(func(args mock.Arguments) { sync.Stop() }).
// 					Return(nil).
// 					Once()

// 				m.Etherman.
// 					On("GetLatestBatchNumber").
// 					Return(uint64(10), nil).
// 					Once()

// 				var nilDbTx pgx.Tx
// 				m.State.
// 					On("GetLastBatchNumber", ctx, nilDbTx).
// 					Return(uint64(10), nil).
// 					Once()
// 			}).
// 			Return(m.DbTx, nil).
// 			Once()

// 		return sync
// 	}

// 	testCases := []testCase{
// 		{
// 			Name: "Transactions are different",
// 			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
// 				return &state.Batch{
// 					BatchNumber:    1,
// 					BatchL2Data:    []byte{1},
// 					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
// 					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
// 					Coinbase:       sequencedBatch.Coinbase,
// 				}
// 			},
// 			getTrustedReorg: func(m *mocks, batchNumber, timestamp uint64) state.TrustedReorg {
// 				return state.TrustedReorg{
// 					BatchNumber: batchNumber,
// 					Reason:      "Different field BatchL2Data. Virtual: , Trusted: 01\n",
// 				}
// 			},
// 		},
// 		{
// 			Name: "Global Exit Root is different",
// 			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
// 				return &state.Batch{
// 					BatchNumber:    1,
// 					BatchL2Data:    sequencedBatch.Transactions,
// 					GlobalExitRoot: common.HexToHash("0x999888777"),
// 					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
// 					Coinbase:       sequencedBatch.Coinbase,
// 				}
// 			},
// 			getTrustedReorg: func(m *mocks, batchNumber, timestamp uint64) state.TrustedReorg {
// 				return state.TrustedReorg{
// 					BatchNumber: batchNumber,
// 					Reason:      "Different field GlobalExitRoot. Virtual: 0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20, Trusted: 0x0000000000000000000000000000000000000000000000000000000999888777\n",
// 				}
// 			},
// 		},
// 		{
// 			Name: "Timestamp is different",
// 			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
// 				return &state.Batch{
// 					BatchNumber:    1,
// 					BatchL2Data:    sequencedBatch.Transactions,
// 					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
// 					Timestamp:      time.Unix(int64(0), 0),
// 					Coinbase:       sequencedBatch.Coinbase,
// 				}
// 			},
// 			getTrustedReorg: func(m *mocks, batchNumber, timestamp uint64) state.TrustedReorg {
// 				return state.TrustedReorg{
// 					BatchNumber: batchNumber,
// 					Reason:      "Different field Timestamp. Virtual: " + strconv.FormatUint(timestamp, 10) + ", Trusted: 0\n",
// 				}
// 			},
// 		},
// 		{
// 			Name: "Coinbase is different",
// 			getTrustedBatch: func(m *mocks, ctx context.Context, sequencedBatch etherman.SequencedBatch) *state.Batch {
// 				return &state.Batch{
// 					BatchNumber:    1,
// 					BatchL2Data:    sequencedBatch.Transactions,
// 					GlobalExitRoot: sequencedBatch.GlobalExitRoot,
// 					Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
// 					Coinbase:       common.HexToAddress("0x999888777"),
// 				}
// 			},
// 			getTrustedReorg: func(m *mocks, batchNumber, timestamp uint64) state.TrustedReorg {
// 				return state.TrustedReorg{
// 					BatchNumber: batchNumber,
// 					Reason:      "Different field Coinbase. Virtual: 0x0000000000000000000000000000000000000222, Trusted: 0x0000000000000000000000000000000999888777\n",
// 				}
// 			},
// 		},
// 	}

// 	m := mocks{
// 		Etherman:     newEthermanMock(t),
// 		State:        newStateMock(t),
// 		Pool:         newPoolMock(t),
// 		EthTxManager: newEthTxManagerMock(t),
// 		DbTx:         newDbTxMock(t),
// 		ZKEVMClient:  newZkEVMClientMock(t),
// 	}

// 	// start synchronizing
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			testCase := tc
// 			sync := setupMocks(&m, &testCase)
// 			err := sync.Sync()
// 			require.NoError(t, err)
// 		})
// 	}
// }

func TestForcedBatch(t *testing.T) {
	genesis := state.Genesis{}
	cfg := Config{
		SyncInterval:   cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:  10,
		GenBlockNumber: uint64(123456),
	}

	m := mocks{
		Etherman:    newEthermanMock(t),
		State:       newStateMock(t),
		Pool:        newPoolMock(t),
		DbTx:        newDbTxMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
	}

	sync, err := NewSynchronizer(false, m.Etherman, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, genesis, cfg)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader := &ethTypes.Header{Number: big.NewInt(1), ParentHash: parentHash}
			ethBlock := ethTypes.NewBlockWithHeader(ethHeader)
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

			t := time.Now()
			sequencedBatch := etherman.SequencedBatch{
				BatchNumber:   uint64(2),
				Coinbase:      common.HexToAddress("0x222"),
				SequencerAddr: common.HexToAddress("0x00"),
				TxHash:        common.HexToHash("0x333"),
				PolygonZkEVMBatchData: polygonzkevm.PolygonZkEVMBatchData{
					Transactions:       []byte{},
					GlobalExitRoot:     [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					Timestamp:          uint64(t.Unix()),
					MinForcedTimestamp: 1000, //ForcedBatch
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.GlobalExitRoot,
				RawTxsData:        sequencedBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedBatch.MinForcedTimestamp), 0),
			}}

			ethermanBlock := etherman.Block{
				BlockHash:        ethBlock.Hash(),
				SequencedBatches: [][]etherman.SequencedBatch{{sequencedBatch}},
				ForcedBatches:    forceb,
			}
			blocks := []etherman.Block{ethermanBlock}
			order := map[common.Hash][]etherman.Order{
				ethBlock.Hash(): {
					{
						Name: etherman.ForcedBatchesOrder,
						Pos:  0,
					},
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

			m.ZKEVMClient.
				On("BatchNumber", ctx).
				Return(uint64(1), nil).
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

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedBatch.Coinbase,
				GlobalExitRoot:    sequencedBatch.GlobalExitRoot,
				RawTxsData:        sequencedBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedBatch.MinForcedTimestamp), 0),
			}}

			m.State.
				On("AddForcedBatch", ctx, &fb[0], m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetNextForcedBatches", ctx, 1, m.DbTx).
				Return(fb, nil).
				Once()

			trustedBatch := &state.Batch{
				BatchL2Data:    sequencedBatch.Transactions,
				GlobalExitRoot: sequencedBatch.GlobalExitRoot,
				Timestamp:      time.Unix(int64(sequencedBatch.Timestamp), 0),
				Coinbase:       sequencedBatch.Coinbase,
			}

			m.State.
				On("GetBatchByNumber", ctx, sequencedBatch.BatchNumber, m.DbTx).
				Return(trustedBatch, nil).
				Once()

			var forced uint64 = 1
			sbatch := state.Batch{
				BatchNumber:    sequencedBatch.BatchNumber,
				Coinbase:       common.HexToAddress("0x222"),
				BatchL2Data:    []byte{},
				Timestamp:      time.Unix(int64(t.Unix()), 0),
				Transactions:   nil,
				GlobalExitRoot: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
				ForcedBatchNum: &forced,
			}
			m.State. //ExecuteBatch(s.ctx, batch.BatchNumber, batch.BatchL2Data, dbTx
					On("ExecuteBatch", ctx, sbatch, false, m.DbTx).
					Return(&pb.ProcessBatchResponse{NewStateRoot: trustedBatch.StateRoot.Bytes()}, nil).
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

			seq := state.Sequence{
				FromBatchNumber: 2,
				ToBatchNumber:   2,
			}
			m.State.
				On("AddSequence", ctx, seq, m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("AddAccumulatedInputHash", ctx, sequencedBatch.BatchNumber, common.Hash{}, m.DbTx).
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

	err = sync.Sync()
	require.NoError(t, err)
}

func TestSequenceForcedBatch(t *testing.T) {
	genesis := state.Genesis{}
	cfg := Config{
		SyncInterval:   cfgTypes.Duration{Duration: 1 * time.Second},
		SyncChunkSize:  10,
		GenBlockNumber: uint64(123456),
	}

	m := mocks{
		Etherman:    newEthermanMock(t),
		State:       newStateMock(t),
		Pool:        newPoolMock(t),
		DbTx:        newDbTxMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
	}

	sync, err := NewSynchronizer(true, m.Etherman, m.State, m.Pool, m.EthTxManager, m.ZKEVMClient, genesis, cfg)
	require.NoError(t, err)

	// state preparation
	ctxMatchBy := mock.MatchedBy(func(ctx context.Context) bool { return ctx != nil })
	m.State.
		On("BeginStateTransaction", ctxMatchBy).
		Run(func(args mock.Arguments) {
			ctx := args[0].(context.Context)
			parentHash := common.HexToHash("0x111")
			ethHeader := &ethTypes.Header{Number: big.NewInt(1), ParentHash: parentHash}
			ethBlock := ethTypes.NewBlockWithHeader(ethHeader)
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

			sequencedForceBatch := etherman.SequencedForceBatch{
				BatchNumber: uint64(2),
				Coinbase:    common.HexToAddress("0x222"),
				TxHash:      common.HexToHash("0x333"),
				PolygonZkEVMForcedBatchData: polygonzkevm.PolygonZkEVMForcedBatchData{
					Transactions:       []byte{},
					GlobalExitRoot:     [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					MinForcedTimestamp: 1000, //ForcedBatch
				},
			}

			forceb := []etherman.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.GlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.MinForcedTimestamp), 0),
			}}

			ethermanBlock := etherman.Block{
				BlockHash:             ethBlock.Hash(),
				SequencedForceBatches: [][]etherman.SequencedForceBatch{{sequencedForceBatch}},
				ForcedBatches:         forceb,
			}
			blocks := []etherman.Block{ethermanBlock}
			order := map[common.Hash][]etherman.Order{
				ethBlock.Hash(): {
					{
						Name: etherman.ForcedBatchesOrder,
						Pos:  0,
					},
					{
						Name: etherman.SequenceForceBatchesOrder,
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

			fb := []state.ForcedBatch{{
				BlockNumber:       lastBlock.BlockNumber,
				ForcedBatchNumber: 1,
				Sequencer:         sequencedForceBatch.Coinbase,
				GlobalExitRoot:    sequencedForceBatch.GlobalExitRoot,
				RawTxsData:        sequencedForceBatch.Transactions,
				ForcedAt:          time.Unix(int64(sequencedForceBatch.MinForcedTimestamp), 0),
			}}

			m.State.
				On("AddForcedBatch", ctx, &fb[0], m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetLastVirtualBatchNum", ctx, m.DbTx).
				Return(uint64(1), nil).
				Once()

			m.State.
				On("ResetTrustedState", ctx, uint64(1), m.DbTx).
				Return(nil).
				Once()

			m.State.
				On("GetNextForcedBatches", ctx, 1, m.DbTx).
				Return(fb, nil).
				Once()

			f := uint64(1)
			processingContext := state.ProcessingContext{
				BatchNumber:    sequencedForceBatch.BatchNumber,
				Coinbase:       sequencedForceBatch.Coinbase,
				Timestamp:      ethBlock.ReceivedAt,
				GlobalExitRoot: sequencedForceBatch.GlobalExitRoot,
				ForcedBatchNum: &f,
			}

			m.State.
				On("ProcessAndStoreClosedBatch", ctx, processingContext, sequencedForceBatch.Transactions, m.DbTx, state.SynchronizerCallerLabel).
				Return(common.Hash{}, nil).
				Once()

			virtualBatch := &state.VirtualBatch{
				BatchNumber:   sequencedForceBatch.BatchNumber,
				TxHash:        sequencedForceBatch.TxHash,
				Coinbase:      sequencedForceBatch.Coinbase,
				SequencerAddr: sequencedForceBatch.Coinbase,
				BlockNumber:   ethermanBlock.BlockNumber,
			}

			m.State.
				On("AddVirtualBatch", ctx, virtualBatch, m.DbTx).
				Return(nil).
				Once()

			seq := state.Sequence{
				FromBatchNumber: 2,
				ToBatchNumber:   2,
			}
			m.State.
				On("AddSequence", ctx, seq, m.DbTx).
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

	err = sync.Sync()
	require.NoError(t, err)
}
