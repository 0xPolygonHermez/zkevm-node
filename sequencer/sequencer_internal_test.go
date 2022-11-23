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

func TestShouldCloseSequenceTooBig(t *testing.T) {
	s := Sequencer{}
	s.sequenceInProgress = ethManTypes.Sequence{IsSequenceTooBig: true}
	ctx := context.Background()
	shouldClose := s.shouldCloseSequenceInProgress(ctx)
	require.False(t, s.sequenceInProgress.IsSequenceTooBig)
	require.True(t, shouldClose)
}

func TestShouldCloseSequenceReachedMaxAmountOfTxs(t *testing.T) {
	s := Sequencer{cfg: Config{MaxTxsPerBatch: 150}}
	txs := make([]types.Transaction, 0, s.cfg.MaxTxsPerBatch)
	for i := uint64(0); i < s.cfg.MaxTxsPerBatch; i++ {
		tx := types.NewTransaction(i, common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
		txs = append(txs, *tx)
	}
	s.sequenceInProgress = ethManTypes.Sequence{Txs: txs}
	ctx := context.Background()
	shouldClose := s.shouldCloseSequenceInProgress(ctx)
	require.True(t, shouldClose)
}

func TestShouldCloseDueToNewDeposits(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	dbTx := new(sequencerMocks.DbTxMock)
	s := Sequencer{cfg: Config{WaitBlocksToUpdateGER: 10, WaitBlocksToConsiderGerFinal: 6}, state: st, etherman: eth}
	ctx := context.Background()
	mainnetExitRoot := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a53cf2d7d9f1")
	lastGer := state.GlobalExitRoot{
		BlockNumber:     1,
		Timestamp:       time.Now(),
		MainnetExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		RollupExitRoot:  common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:  common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	s.sequenceInProgress.GlobalExitRoot = lastGer.GlobalExitRoot
	st.On("GetBlockNumAndMainnetExitRootByGER", ctx, s.sequenceInProgress.GlobalExitRoot, nil).Return(lastGer.BlockNumber, mainnetExitRoot, nil)
	st.On("GetLatestGlobalExitRoot", ctx, uint64(6), nil).Return(lastGer, time.Now(), nil)
	eth.On("GetLatestBlockNumber", ctx).Return(uint64(12), nil)
	isShouldCloseDueToNewDeposits, err := s.shouldCloseDueToNewDeposits(ctx)
	require.NoError(t, err)
	require.Equal(t, true, isShouldCloseDueToNewDeposits)
	st.AssertExpectations(t)
	eth.AssertExpectations(t)
	dbTx.AssertExpectations(t)
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
	s := &Sequencer{cfg: Config{MaxTxsPerBatch: 150, MaxBatchBytesSize: 30000}}
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

func TestAppendPendingTxs(t *testing.T) {
	pl := new(sequencerMocks.PoolMock)
	ctx := context.Background()
	s := &Sequencer{pool: pl}
	minGasPrice := big.NewInt(1)
	ticker := time.NewTicker(1 * time.Second)

	poolTx := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	var poolTxs []*pool.Transaction
	poolTxs = append(poolTxs, &pool.Transaction{Transaction: *poolTx})

	pl.On("GetTxs", ctx, pool.TxStatusPending, false, minGasPrice.Uint64(), uint64(150)).Return(poolTxs, nil)
	pendTxsAmount := s.appendPendingTxs(ctx, false, minGasPrice.Uint64(), 150, ticker)
	require.Equal(t, uint64(1), pendTxsAmount)
	require.Equal(t, 1, len(s.sequenceInProgress.Txs))
	pl.AssertExpectations(t)
}

func TestAppendPendingTxsFailedCounter(t *testing.T) {
	pl := new(sequencerMocks.PoolMock)
	ctx := context.Background()
	s := &Sequencer{cfg: Config{MaxAllowedFailedCounter: 5}, pool: pl}
	minGasPrice := big.NewInt(1)
	ticker := time.NewTicker(1 * time.Second)

	poolTx := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	var poolTxs []*pool.Transaction
	poolTxs = append(poolTxs, &pool.Transaction{Transaction: *poolTx, FailedCounter: 55})

	pl.On("GetTxs", ctx, pool.TxStatusPending, false, minGasPrice.Uint64(), uint64(150)).Return(poolTxs, nil)
	pl.On("UpdateTxsStatus", ctx, []string{poolTxs[0].Hash().String()}, pool.TxStatusInvalid).Return(nil)
	pendTxsAmount := s.appendPendingTxs(ctx, false, minGasPrice.Uint64(), 150, ticker)
	require.Equal(t, uint64(0), pendTxsAmount)
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
	txs := []types.Transaction{tx1, tx2}
	s.sequenceInProgress.Txs = txs

	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash:      txs[0].Hash(),
			Tx:          txs[0],
			IsProcessed: true,
		},
		{
			TxHash:      txs[1].Hash(),
			Tx:          txs[1],
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

	processedTxs, processedTxsHashes, unprocessedTxs, unprocessedTxsHashes := state.DetermineProcessedTransactions(processBatchResponse.Responses)

	txsResponse := processTxResponse{
		processedTxs:         processedTxs,
		processedTxsHashes:   processedTxsHashes,
		unprocessedTxs:       unprocessedTxs,
		unprocessedTxsHashes: unprocessedTxsHashes,
		isBatchProcessed:     false,
	}
	txsResponseToReturn := txsResponse
	txsResponseToReturn.isBatchProcessed = true
	lastBatchNumber := uint64(10)
	st.On("GetLastBatchNumber", ctx, dbTx).Return(lastBatchNumber, nil)
	st.On("ProcessSequencerBatch", ctx, lastBatchNumber, txs, dbTx).Return(processBatchResponse, nil)

	unprocessedTxsAfterReprocess, err := s.reprocessBatch(ctx, txsResponse, ethManTypes.Sequence{})
	require.NoError(t, err)
	require.Equal(t, 0, len(unprocessedTxsAfterReprocess))
	require.Equal(t, 2, len(s.sequenceInProgress.Txs))
	st.AssertExpectations(t)
}

func TestUpdateTxsInPool(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	pl := new(sequencerMocks.PoolMock)
	s := &Sequencer{state: st, pool: pl}
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Second)

	tx1 := *types.NewTransaction(0, common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := *types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	tx3 := *types.NewTransaction(3, common.HexToAddress("0x2"), big.NewInt(1), 0, big.NewInt(1), []byte("ddd"))

	txs := []types.Transaction{tx1, tx2, tx3}
	s.sequenceInProgress.Txs = txs

	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash:      txs[0].Hash(),
			Tx:          txs[0],
			IsProcessed: true,
		},
		{
			TxHash:      txs[1].Hash(),
			Tx:          txs[1],
			IsProcessed: true,
		},
		{
			TxHash:      txs[2].Hash(),
			Tx:          txs[2],
			IsProcessed: false,
		},
	}

	processBatchResponse := &state.ProcessBatchResponse{
		CumulativeGasUsed: 100000,
		IsBatchProcessed:  true,
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

	pl.On("UpdateTxsStatus", ctx, processedTxsHashes, pool.TxStatusSelected).Return(nil)

	fromAddress := common.HexToAddress("0x123")
	pl.On("GetTxFromAddressFromByHash", ctx, txs[2].Hash()).Return(fromAddress, txs[2].Nonce(), nil)
	l2BlockNumber := uint64(3)
	st.On("GetLastL2BlockNumber", ctx, nil).Return(l2BlockNumber, nil)
	accNonce := uint64(2)
	st.On("GetNonce", ctx, fromAddress, l2BlockNumber, nil).Return(accNonce, nil)
	pl.On("UpdateTxsStatus", ctx, unprocessedTxsHashes, pool.TxStatusFailed).Return(nil)
	pl.On("IncrementFailedCounter", ctx, unprocessedTxsHashes).Return(nil)
	s.updateTxsInPool(ctx, ticker, txsResponse, unprocessedTxs)
	st.AssertExpectations(t)
	pl.AssertExpectations(t)
}

func TestTryToProcessTxs(t *testing.T) {
	st := new(sequencerMocks.StateMock)
	eth := new(sequencerMocks.EthermanMock)
	dbTx := new(sequencerMocks.DbTxMock)
	pl := new(sequencerMocks.PoolMock)

	gpe := new(sequencerMocks.GasPriceEstimatorMock)
	s := Sequencer{cfg: Config{
		MaxTimeForBatchToBeOpen: cfgTypes.NewDuration(5 * time.Second),
		MaxBatchBytesSize:       30000,
		MaxTxsPerBatch:          150,
	}, state: st, etherman: eth, gpe: gpe, pool: pl}
	ctx := context.Background()
	// Check if synchronizer is up to date
	st.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(1), nil)
	eth.On("GetLatestBatchNumber").Return(uint64(1), nil)

	isSynced := s.isSynced(ctx)
	require.Equal(t, true, isSynced)

	mainnetExitRoot := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a53cf2d7d9f1")
	lastGer := state.GlobalExitRoot{
		BlockNumber:     1,
		Timestamp:       time.Now(),
		MainnetExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a53cf2d7d9f1"),
		RollupExitRoot:  common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:  common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	s.sequenceInProgress.GlobalExitRoot = lastGer.GlobalExitRoot
	eth.On("GetLatestBlockNumber", ctx).Return(uint64(1), nil)
	st.On("GetBlockNumAndMainnetExitRootByGER", ctx, s.sequenceInProgress.GlobalExitRoot, nil).Return(lastGer.BlockNumber, mainnetExitRoot, nil)
	st.On("GetLatestGlobalExitRoot", ctx, uint64(1), nil).Return(lastGer, time.Now(), nil)
	st.On("BeginStateTransaction", ctx).Return(dbTx, nil)
	dbTx.On("Commit", ctx).Return(nil)

	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	s.sequenceInProgress.Txs = []types.Transaction{*tx}
	s.sequenceInProgress.Timestamp = time.Now().Unix()

	lastBatchNumber := uint64(10)
	st.On("GetLastBatchNumber", ctx, nil).Return(lastBatchNumber, nil)
	st.On("GetLastBatchNumber", ctx, dbTx).Return(lastBatchNumber, nil)

	st.On("IsBatchVirtualized", ctx, lastBatchNumber-1, nil).Return(true, nil)

	minGasPrice := big.NewInt(1)
	gpe.On("GetAvgGasPrice", ctx).Return(minGasPrice, nil)

	ticker := time.NewTicker(1 * time.Second)
	poolTx := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	var poolTxs []*pool.Transaction
	poolTxs = append(poolTxs, &pool.Transaction{Transaction: *poolTx})
	pl.On("GetTxs", ctx, pool.TxStatusPending, true, uint64(0), uint64(149)).Return([]*pool.Transaction{}, nil)
	pl.On("GetTxs", ctx, pool.TxStatusFailed, true, uint64(0), uint64(149)).Return([]*pool.Transaction{}, nil)

	pl.On("GetTxs", ctx, pool.TxStatusPending, false, minGasPrice.Uint64(), uint64(149)).Return(poolTxs, nil)

	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash:      s.sequenceInProgress.Txs[0].Hash(),
			Tx:          s.sequenceInProgress.Txs[0],
			IsProcessed: true,
		},
		{
			TxHash:      poolTxs[0].Transaction.Hash(),
			Tx:          poolTxs[0].Transaction,
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
	var txs = s.sequenceInProgress.Txs
	txs = append(txs, poolTxs[0].Transaction)
	st.On("ProcessSequencerBatch", ctx, lastBatchNumber, txs, dbTx).Return(processBatchResponse, nil)

	processedTxs, processedTxsHashes, _, _ := state.DetermineProcessedTransactions(processBatchResponse.Responses)

	st.On("StoreTransactions", ctx, lastBatchNumber, processedTxs, dbTx).Return(nil)
	pl.On("UpdateTxsStatus", ctx, processedTxsHashes, pool.TxStatusSelected).Return(nil)

	s.tryToProcessTx(ctx, ticker)

	st.AssertExpectations(t)
	pl.AssertExpectations(t)
	eth.AssertExpectations(t)
}
