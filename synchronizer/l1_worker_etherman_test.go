package synchronizer

import (
	context "context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExploratoryWorker(t *testing.T) {
	t.Skip("no real test, just exploratory")
	cfg := etherman.Config{
		URL: "http://localhost:8545",
	}

	l1Config := etherman.L1Config{
		L1ChainID:                 1337,
		ZkEVMAddr:                 common.HexToAddress("0x610178dA211FEF7D417bC0e6FeD39F05609AD788"),
		MaticAddr:                 common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		GlobalExitRootManagerAddr: common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"),
	}

	ethermanClient, err := etherman.NewClient(cfg, l1Config)
	require.NoError(t, err)
	worker := newWorker(ethermanClient)
	ch := make(chan responseRollupInfoByBlockRange)
	blockRange := blockRange{
		fromBlock: 9847396,
		toBlock:   9847396,
	}
	err = worker.asyncRequestRollupInfoByBlockRange(newContextWithNone(context.Background()), ch, nil, newRequestNoSleep(blockRange))
	require.NoError(t, err)
	result := <-ch
	log.Info(result.toStringBrief())
	for i := range result.result.blocks {
		for _, element := range result.result.order[result.result.blocks[i].BlockHash] {
			switch element.Name {
			case etherman.SequenceBatchesOrder:
				for i := range result.result.blocks[i].SequencedBatches {
					log.Infof("SequenceBatchesOrder %v %v %v", element.Pos, result.result.blocks[i].SequencedBatches[element.Pos][i].BatchNumber,
						result.result.blocks[i].BlockNumber)
				}
			default:
				log.Info("unknown order", element.Name)
			}
		}
	}
	require.Equal(t, result.generic.err.Error(), "not found")
}

func TestIfRollupRequestReturnsErrorDontRequestEthBlockByNumber(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	blockRange := blockRange{
		fromBlock: 100,
		toBlock:   20000,
	}
	ctx := newContextWithTimeout(context.Background(), time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	expectedCallsForEmptyRollupInfo(mockEtherman, blockRange, errors.New("error"), nil)
	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, &wg, newRequestNoSleep(blockRange))
	require.NoError(t, err)
	wg.Wait()
}

func TestIfWorkerIsBusyReturnsAnErrorUpdateWaitGroupAndCancelContext(t *testing.T) {
	sut, _, ch := setupWorkerEthermanTest(t)
	blockRange := blockRange{
		fromBlock: 100,
		toBlock:   20000,
	}
	ctx := newContextWithTimeout(context.Background(), time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	sut.setStatus(ethermanWorking)
	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, &wg, newRequestNoSleep(blockRange))
	require.Error(t, err)
	wg.Wait()
	select {
	case <-ctx.Done():
	default:
		require.Fail(t, "The context should be cancelled")
	}
}

// Given: a request to get the rollup info by block range that is OK
// When: the request is finished
// Then: the context is canceled
func TestGivenOkRequestWhenFinishThenCancelTheContext(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	blockRange := blockRange{
		fromBlock: 100,
		toBlock:   20000,
	}
	ctx := newContextWithTimeout(context.Background(), time.Second)
	expectedCallsForEmptyRollupInfo(mockEtherman, blockRange, nil, nil)
	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, nil, newRequestNoSleep(blockRange))
	require.NoError(t, err)
	result := <-ch
	require.NoError(t, result.generic.err)
	select {
	case <-ctx.Done():
	default:
		require.Fail(t, "The context should be cancelled")
	}
}

func TestGivenOkRequestWithSleepWhenFinishThenMustExuctedTheSleep(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	blockRange := blockRange{
		fromBlock: 100,
		toBlock:   20000,
	}
	ctx := newContextWithTimeout(context.Background(), time.Second)
	expectedCallsForEmptyRollupInfo(mockEtherman, blockRange, nil, nil)
	startTime := time.Now()
	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, nil, newRequestSleep(blockRange, time.Millisecond*500))
	require.NoError(t, err)
	result := <-ch
	require.NoError(t, result.generic.err)
	require.GreaterOrEqual(t, time.Since(startTime).Milliseconds(), int64(500))
}

func TestCheckIsIdleFunction(t *testing.T) {
	tcs := []struct {
		status         ethermanStatusEnum
		expectedIsIdle bool
	}{
		{status: ethermanIdle, expectedIsIdle: true},
		{status: ethermanWorking, expectedIsIdle: false},
		{status: ethermanError, expectedIsIdle: false},
	}
	for _, tc := range tcs {
		t.Run(tc.status.String(), func(t *testing.T) {
			sut, _, _ := setupWorkerEthermanTest(t)
			sut.setStatus(tc.status)
			require.Equal(t, tc.expectedIsIdle, sut.isIdle())
		})
	}
}

func TestIfRollupInfoFailGettingLastBlockContainBlockRange(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	var wg sync.WaitGroup
	wg.Add(1)
	ctx := newContextWithTimeout(context.Background(), time.Second)
	blockRange := blockRange{fromBlock: 100, toBlock: 20000}
	request := newRequestNoSleep(blockRange)
	request.requestPreviousBlock = true
	request.requestLastBlockIfNoBlocksInAnswer = requestLastBlockModeAlways

	mockEtherman.
		On("EthBlockByNumber", mock.Anything, blockRange.toBlock).
		Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(blockRange.toBlock))}), fmt.Errorf("error")).
		Once()
	mockEtherman.
		On("GetRollupInfoByBlockRange", mock.Anything, blockRange.fromBlock, mock.Anything).
		Return([]etherman.Block{}, map[common.Hash][]etherman.Order{}, nil).
		Maybe()

	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, &wg, request)
	require.NoError(t, err)
	result := <-ch
	require.Error(t, result.generic.err)
	require.True(t, result.result != nil)
	require.Equal(t, result.result.blockRange, blockRange)
}

func TestIfRollupInfoFailGettingRollupContainBlockRange(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	var wg sync.WaitGroup
	wg.Add(1)
	ctx := newContextWithTimeout(context.Background(), time.Second)
	blockRange := blockRange{fromBlock: 100, toBlock: 20000}
	request := newRequestNoSleep(blockRange)
	request.requestPreviousBlock = true
	request.requestLastBlockIfNoBlocksInAnswer = requestLastBlockModeAlways

	mockEtherman.
		On("EthBlockByNumber", mock.Anything, blockRange.toBlock).
		Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(blockRange.toBlock))}), nil).
		Maybe()
	mockEtherman.
		On("GetRollupInfoByBlockRange", mock.Anything, blockRange.fromBlock, mock.Anything).
		Return([]etherman.Block{}, map[common.Hash][]etherman.Order{}, fmt.Errorf("error")).
		Once()

	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, &wg, request)
	require.NoError(t, err)
	result := <-ch
	require.Error(t, result.generic.err)
	require.True(t, result.result != nil)
	require.Equal(t, result.result.blockRange, blockRange)
}

func TestIfRollupInfoFailPreviousBlockContainBlockRange(t *testing.T) {
	sut, mockEtherman, ch := setupWorkerEthermanTest(t)
	var wg sync.WaitGroup
	wg.Add(1)
	ctx := newContextWithTimeout(context.Background(), time.Second)
	blockRange := blockRange{fromBlock: 100, toBlock: 20000}
	request := newRequestNoSleep(blockRange)
	request.requestPreviousBlock = true
	request.requestLastBlockIfNoBlocksInAnswer = requestLastBlockModeAlways

	mockEtherman.
		On("EthBlockByNumber", mock.Anything, blockRange.toBlock).
		Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(blockRange.toBlock))}), nil).
		Maybe()
	mockEtherman.
		On("GetRollupInfoByBlockRange", mock.Anything, blockRange.fromBlock, mock.Anything).
		Return([]etherman.Block{}, map[common.Hash][]etherman.Order{}, nil).
		Maybe()
	mockEtherman.
		On("EthBlockByNumber", mock.Anything, blockRange.fromBlock-1).
		Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(blockRange.fromBlock - 1))}), fmt.Errorf("error")).
		Once()

	err := sut.asyncRequestRollupInfoByBlockRange(ctx, ch, &wg, request)
	require.NoError(t, err)
	result := <-ch
	require.Error(t, result.generic.err)
	require.True(t, result.result != nil)
	require.Equal(t, result.result.blockRange, blockRange)
}

func expectedCallsForEmptyRollupInfo(mockEtherman *ethermanMock, blockRange blockRange, getRollupError error, ethBlockError error) {
	mockEtherman.
		On("GetRollupInfoByBlockRange", mock.Anything, blockRange.fromBlock, mock.Anything).
		Return([]etherman.Block{}, map[common.Hash][]etherman.Order{}, getRollupError).
		Once()

	if getRollupError == nil {
		mockEtherman.
			On("EthBlockByNumber", mock.Anything, blockRange.toBlock).
			Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{Number: big.NewInt(int64(blockRange.toBlock))}), ethBlockError).
			Once()
	}
}

func setupWorkerEthermanTest(t *testing.T) (*workerEtherman, *ethermanMock, chan responseRollupInfoByBlockRange) {
	mockEtherman := newEthermanMock(t)
	worker := newWorker(mockEtherman)
	ch := make(chan responseRollupInfoByBlockRange, 2)
	return worker, mockEtherman, ch
}

func newRequestNoSleep(blockRange blockRange) requestRollupInfoByBlockRange {
	return requestRollupInfoByBlockRange{
		blockRange:                         blockRange,
		sleepBefore:                        noSleepTime,
		requestLastBlockIfNoBlocksInAnswer: requestLastBlockModeIfNoBlocksInAnswer,
		requestPreviousBlock:               false,
	}
}

func newRequestSleep(blockRange blockRange, sleep time.Duration) requestRollupInfoByBlockRange {
	return requestRollupInfoByBlockRange{
		blockRange:                         blockRange,
		sleepBefore:                        sleep,
		requestLastBlockIfNoBlocksInAnswer: requestLastBlockModeIfNoBlocksInAnswer,
	}
}
