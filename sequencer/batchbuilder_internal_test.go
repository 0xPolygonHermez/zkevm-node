package sequencer

import (
	"context"
	"math/big"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	ethManTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	sequencerMocks "github.com/0xPolygonHermez/zkevm-node/sequencer/mocks"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestIsSynced(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	s := Sequencer{state: st, etherman: eth}
	ctx := context.Background()
	st.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(1), nil)
	eth.On("GetLatestBatchNumber").Return(uint64(1), nil)
	isSynced := s.isSynced(ctx)
	require.Equal(t, true, isSynced)
	st.AssertExpectations(t)
	eth.AssertExpectations(t)
}

func TestIsNotSynced(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	s := Sequencer{state: st, etherman: eth}
	ctx := context.Background()
	st.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(1), nil)
	eth.On("GetLatestBatchNumber").Return(uint64(2), nil)
	isSynced := s.isSynced(ctx)
	require.Equal(t, false, isSynced)
	st.AssertExpectations(t)
	eth.AssertExpectations(t)
}

func TestShouldCloseDueToNewDepositsUpdateGER(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	dbTx := new(sequencerMocks.DbTxMock)
	s := Sequencer{cfg: Config{WaitBlocksToUpdateGER: 10}, state: st, etherman: eth}
	ctx := context.Background()
	mainnetExitRoot := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a53cf2d7d9f1")
	lastGer := &state.GlobalExitRoot{
		BlockNumber:       1,
		GlobalExitRootNum: big.NewInt(2),
		MainnetExitRoot:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		RollupExitRoot:    common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:    common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	s.sequenceInProgress.GlobalExitRoot = lastGer.GlobalExitRoot
	st.On("GetBlockNumAndMainnetExitRootByGER", ctx, s.sequenceInProgress.GlobalExitRoot, nil).Return(lastGer.BlockNumber, mainnetExitRoot, nil)
	st.On("GetLatestGlobalExitRoot", ctx, nil).Return(lastGer, nil)
	eth.On("GetLatestBlockNumber", ctx).Return(uint64(12), nil)
	st.On("BeginStateTransaction", ctx).Return(dbTx, nil)
	dbTx.On("Commit", ctx).Return(nil)
	st.On("UpdateGERInOpenBatch", ctx, lastGer.GlobalExitRoot, dbTx).Return(nil).Once()
	isShouldCloseDueToNewDeposits, err := s.shouldCloseDueToNewDeposits(ctx)
	require.NoError(t, err)
	require.Equal(t, false, isShouldCloseDueToNewDeposits)
	st.AssertExpectations(t)
	eth.AssertExpectations(t)
	dbTx.AssertExpectations(t)
}

func TestShouldCloseDueToNewDepositsSequenceProfitable(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	profitabilityChecker := profitabilitychecker.New(profitabilitychecker.Config{SendBatchesEvenWhenNotProfitable: true}, nil, nil)
	s := Sequencer{cfg: Config{WaitBlocksToUpdateGER: 10}, state: st, etherman: eth, checker: profitabilityChecker}
	ctx := context.Background()
	mainnetExitRoot := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a53cf2d7d9f1")
	lastGer := &state.GlobalExitRoot{
		BlockNumber:       1,
		GlobalExitRootNum: big.NewInt(2),
		MainnetExitRoot:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		RollupExitRoot:    common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:    common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	s.sequenceInProgress.GlobalExitRoot = lastGer.GlobalExitRoot
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	s.sequenceInProgress.Txs = []types.Transaction{*tx}
	st.On("GetBlockNumAndMainnetExitRootByGER", ctx, s.sequenceInProgress.GlobalExitRoot, nil).Return(lastGer.BlockNumber, mainnetExitRoot, nil)
	st.On("GetLatestGlobalExitRoot", ctx, nil).Return(lastGer, nil)
	eth.On("GetLatestBlockNumber", ctx).Return(uint64(12), nil)

	isShouldCloseDueToNewDeposits, err := s.shouldCloseDueToNewDeposits(ctx)
	require.NoError(t, err)
	require.Equal(t, true, isShouldCloseDueToNewDeposits)
	st.AssertExpectations(t)
	eth.AssertExpectations(t)
}

func TestShouldCloseTooLongSinceLastVirtualized(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	s := Sequencer{cfg: Config{MaxTimeForBatchToBeOpen: cfgTypes.NewDuration(1 * time.Second)}, state: st}
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	s.sequenceInProgress.Txs = []types.Transaction{*tx}
	s.sequenceInProgress.Timestamp = time.Now().Add(-s.cfg.MaxTimeForBatchToBeOpen.Duration).Unix()
	ctx := context.Background()
	lastBatchNumber := uint64(10)
	st.On("GetLastBatchNumber", ctx, nil).Return(lastBatchNumber, nil)
	st.On("IsBatchVirtualized", ctx, lastBatchNumber-1, nil).Return(true, nil)
	isShouldCloseTooLongSinceLastVirtualized, err := s.shouldCloseTooLongSinceLastVirtualized(ctx)
	require.NoError(t, err)
	require.True(t, isShouldCloseTooLongSinceLastVirtualized)
	st.AssertExpectations(t)
}

func TestCleanTxsIfTxsDataIsBiggerThanExpected(t *testing.T) {
	s := &Sequencer{}
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	for i := 0; i < 300; i++ {
		s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, *tx)
	}
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Second)
	err := s.cleanTxsIfTxsDataIsBiggerThanExpected(ctx, ticker)
	require.NoError(t, err)
	require.True(t, s.sequenceInProgress.IsSequenceTooBig)
	// 1 transfer equals ~103 bytes, so 291 txs ~= 30000 bytes, what is maximum
	require.Equal(t, 291, len(s.sequenceInProgress.Txs))
}

func TestCleanTxsIfTxsDataIsBiggerThanExpectedTxIsTooBig(t *testing.T) {
	pl := new(sequencerMocks.PoolMock)
	s := &Sequencer{pool: pl}
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	txs := []types.Transaction{}
	for i := 0; i < 300; i++ {
		txs = append(txs, *tx)
	}
	data, err := state.EncodeTransactions(txs)
	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), data)
	s.sequenceInProgress.Txs = []types.Transaction{*tx1}
	require.NoError(t, err)
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Second)
	pl.On("UpdateTxStatus", ctx, s.sequenceInProgress.Txs[0].Hash(), pool.TxStatusInvalid).Return(nil)
	err = s.cleanTxsIfTxsDataIsBiggerThanExpected(ctx, ticker)
	require.NoError(t, err)
	require.False(t, s.sequenceInProgress.IsSequenceTooBig)
	// 1 transfer equals ~103 bytes, so 291 txs ~= 30000 bytes, what is maximum
	require.Equal(t, 0, len(s.sequenceInProgress.Txs))
	pl.AssertExpectations(t)
}

func TestProcessBatch(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	s := &Sequencer{state: st}
	dbTx := new(sequencerMocks.DbTxMock)
	ctx := context.Background()
	st.On("BeginStateTransaction", ctx).Return(dbTx, nil)
	dbTx.On("Commit", ctx).Return(nil)
	lastBatchNumber := uint64(10)
	tx1 := *types.NewTransaction(0, common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))

	s.sequenceInProgress.Txs = []types.Transaction{tx1, tx2}

	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash:      tx1.Hash(),
			Tx:          tx1,
			IsProcessed: true,
		},
		{
			TxHash:      tx2.Hash(),
			Tx:          tx2,
			IsProcessed: true,
		},
	}

	processBatchResponse := &state.ProcessBatchResponse{
		CumulativeGasUsed: 100000,
		IsBatchProcessed:  true,
		Responses:         txsBatch1,
		NewStateRoot:      common.HexToHash("0x123"),
		NewLocalExitRoot:  common.HexToHash("0x123"),
	}
	st.On("GetLastBatchNumber", ctx, dbTx).Return(lastBatchNumber, nil)
	st.On("ProcessSequencerBatch", ctx, lastBatchNumber, s.sequenceInProgress.Txs, dbTx).Return(processBatchResponse, nil)
	procResponse, err := s.processTxs(ctx)
	require.NoError(t, err)
	require.True(t, procResponse.isBatchProcessed)
	require.Equal(t, 2, len(procResponse.processedTxs))
	st.AssertExpectations(t)
}

func TestReprocessBatch(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	s := &Sequencer{state: st}
	dbTx := new(sequencerMocks.DbTxMock)
	ctx := context.Background()
	st.On("BeginStateTransaction", ctx).Return(dbTx, nil)
	dbTx.On("Commit", ctx).Return(nil)
	tx1 := *types.NewTransaction(0, common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	s.sequenceInProgress.Txs = []types.Transaction{tx1, tx2}

	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash:      tx1.Hash(),
			Tx:          tx1,
			IsProcessed: true,
		},
		{
			TxHash:      tx2.Hash(),
			Tx:          tx2,
			IsProcessed: true,
		},
	}

	processBatchResponse := &state.ProcessBatchResponse{
		CumulativeGasUsed: 100000,
		IsBatchProcessed:  false,
		Responses:         txsBatch1,
		NewStateRoot:      common.HexToHash("0x123"),
		NewLocalExitRoot:  common.HexToHash("0x123"),
	}

	processedTxs, processedTxsHashes, unprocessedTxs, unprocessedTxsHashes := state.DetermineProcessedTransactions(processBatchResponse.Responses)

	txsResponse := processTxResponse{
		processedTxs:         processedTxs,
		processedTxsHashes:   processedTxsHashes,
		unprocessedTxs:       unprocessedTxs,
		unprocessedTxsHashes: unprocessedTxsHashes,
		isBatchProcessed:     false,
	}
	lastBatchNumber := uint64(10)
	st.On("GetLastBatchNumber", ctx, dbTx).Return(lastBatchNumber, nil)
	st.On("ProcessSequencerBatch", ctx, lastBatchNumber, s.sequenceInProgress.Txs, dbTx).Return(processBatchResponse, nil)

	unprocessedTxsAfterReprocess, err := s.reprocessBatch(ctx, txsResponse, ethManTypes.Sequence{})
	require.NoError(t, err)
	require.Equal(t, 0, len(unprocessedTxsAfterReprocess))
	require.Equal(t, 2, len(s.sequenceInProgress.Txs))

}
