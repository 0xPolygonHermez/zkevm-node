package etrog

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	syncMocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	hashExamplesValues = []string{"0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92",
		"0x9c8fa7ce2e197f9f1b3c30de9f93de3c1cb290e6c118a18446f47a9e1364c3ab",
		"0x896cfc0684057d0560e950dee352189528167f4663609678d19c7a506a03fe4e",
		"0xde6d2dac4b6e0cb39ed1924db533558a23e5c56ab60fadac8c7d21e7eceb121a",
		"0x9883711e78d02992ac1bd6f19de3bf7bb3f926742d4601632da23525e33f8555"}

	addrExampleValues = []string{"0x8dAF17A20c9DBA35f005b6324F493785D239719d",
		"0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e",
		"0x5FbDB2315678afecb367f032d93F642f64180aa3",
		"0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"}
)

type mocksEtrogProcessorL1 struct {
	Etherman             *mock_syncinterfaces.EthermanFullInterface
	State                *mock_syncinterfaces.StateFullInterface
	Pool                 *mock_syncinterfaces.PoolInterface
	Synchronizer         *mock_syncinterfaces.SynchronizerFullInterface
	DbTx                 *syncMocks.DbTxMock
	TimeProvider         *syncCommon.MockTimerProvider
	CriticalErrorHandler *mock_syncinterfaces.CriticalErrorHandler
}

func createMocks(t *testing.T) *mocksEtrogProcessorL1 {
	mocks := &mocksEtrogProcessorL1{
		Etherman:     mock_syncinterfaces.NewEthermanFullInterface(t),
		State:        mock_syncinterfaces.NewStateFullInterface(t),
		Pool:         mock_syncinterfaces.NewPoolInterface(t),
		Synchronizer: mock_syncinterfaces.NewSynchronizerFullInterface(t),
		DbTx:         syncMocks.NewDbTxMock(t),
		//ZKEVMClient:          mock_syncinterfaces.NewZKEVMClientInterface(t),
		TimeProvider:         &syncCommon.MockTimerProvider{},
		CriticalErrorHandler: mock_syncinterfaces.NewCriticalErrorHandler(t),
		//EventLog:     &eventLogMock{},
	}
	return mocks
}

func createSUT(mocks *mocksEtrogProcessorL1) *ProcessorL1SequenceBatchesEtrog {
	return NewProcessorL1SequenceBatches(mocks.State, mocks.Etherman, mocks.Pool, mocks.Synchronizer,
		mocks.TimeProvider, mocks.CriticalErrorHandler)
}

func TestL1SequenceBatchesNoData(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	err := sut.Process(ctx, etherman.Order{}, nil, mocks.DbTx)
	require.ErrorIs(t, err, actions.ErrInvalidParams)
}

func TestL1SequenceBatchesWrongOrder(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	l1Block := etherman.Block{
		SequencedBatches: [][]etherman.SequencedBatch{},
	}
	err := sut.Process(ctx, etherman.Order{Pos: 1}, &l1Block, mocks.DbTx)
	require.Error(t, err)
}

func TestL1SequenceBatchesPermissionlessNewBatchSequenced(t *testing.T) {
	mocks := createMocks(t)
	sut := createSUT(mocks)
	ctx := context.Background()
	l1Block := etherman.Block{
		BlockNumber:      123,
		ReceivedAt:       mocks.TimeProvider.Now(),
		SequencedBatches: [][]etherman.SequencedBatch{},
	}
	l1InfoRoot := common.HexToHash(hashExamplesValues[0])
	l1Block.SequencedBatches = append(l1Block.SequencedBatches, []etherman.SequencedBatch{})
	l1Block.SequencedBatches = append(l1Block.SequencedBatches, []etherman.SequencedBatch{
		{
			BatchNumber:   3,
			L1InfoRoot:    &l1InfoRoot,
			TxHash:        common.HexToHash(hashExamplesValues[1]),
			Coinbase:      common.HexToAddress(addrExampleValues[0]),
			SequencerAddr: common.HexToAddress(addrExampleValues[1]),
			PolygonRollupBaseEtrogBatchData: &polygonzkevm.PolygonRollupBaseEtrogBatchData{
				Transactions: []byte{},
			},
		},
	})
	mocks.State.EXPECT().GetL1InfoTreeDataFromBatchL2Data(ctx, mock.Anything, mocks.DbTx).Return(map[uint32]state.L1DataV2{}, state.ZeroHash, state.ZeroHash, nil)
	mocks.State.EXPECT().GetBatchByNumber(ctx, uint64(3), mocks.DbTx).Return(nil, state.ErrNotFound)
	mocks.Synchronizer.EXPECT().PendingFlushID(mock.Anything, mock.Anything)
	mocks.State.EXPECT().AddVirtualBatch(ctx, mock.Anything, mocks.DbTx).Return(nil)
	mocks.State.EXPECT().AddSequence(ctx, mock.Anything, mocks.DbTx).Return(nil)
	newStateRoot := common.HexToHash(hashExamplesValues[2])
	flushID := uint64(1234)
	proverID := "prover-id"
	mocks.State.EXPECT().ProcessAndStoreClosedBatchV2(ctx, mock.Anything, mocks.DbTx, mock.Anything).Return(newStateRoot, flushID, proverID, nil)
	err := sut.Process(ctx, etherman.Order{Pos: 1}, &l1Block, mocks.DbTx)
	require.NoError(t, err)
}
