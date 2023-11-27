package l1_parallel_sync

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExploratoryL1Get(t *testing.T) {
	t.Skip("Exploratory test")
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	header := new(ethTypes.Header)
	header.Number = big.NewInt(150)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Once()

	err := sut.initialize(context.Background())
	require.NoError(t, err)
	_, err = sut.launchWork()
	require.NoError(t, err)
}

func TestGivenNeedSyncWhenStartThenAskForRollupInfo(t *testing.T) {
	sut, ethermans, _ := setup(t)
	expectedForGettingL1LastBlock(t, ethermans[0], 150)
	expectedRollupInfoCalls(t, ethermans[1], 1)
	err := sut.initialize(context.Background())
	require.NoError(t, err)
	_, err = sut.launchWork()
	require.NoError(t, err)
	var waitDuration = time.Duration(0)

	sut.step(&waitDuration)
	sut.step(&waitDuration)
	sut.workers.waitFinishAllWorkers()
}

func TestGivenNoNeedSyncWhenStartsSendAndEventOfSynchronized(t *testing.T) {
	sut, ethermans, ch := setup(t)
	etherman := ethermans[0]
	// Our last block is 100 in DB and it returns 100 as last block on L1
	// so is synchronized
	expectedForGettingL1LastBlock(t, etherman, 100)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	err := sut.Start(ctx)
	require.NoError(t, err)
	// read everything in channel ch
	for len(ch) > 0 {
		data := <-ch
		if data.ctrlIsValid == true && data.ctrl.event == eventProducerIsFullySynced {
			return // ok
		}
	}
	require.Fail(t, "should not have send a eventProducerIsFullySynced in channel")
}

// Given: Need to synchronize
// When:  Ask for last block
// Then:  Ask for rollupinfo
func TestGivenNeedSyncWhenReachLastBlockThenSendAndEventOfSynchronized(t *testing.T) {
	sut, ethermans, ch := setup(t)
	// Our last block is 100 in DB and it returns 101 as last block on L1
	// so it need to retrieve 1 rollupinfo
	expectedForGettingL1LastBlock(t, ethermans[0], 101)
	expectedRollupInfoCalls(t, ethermans[1], 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	res := sut.Start(ctx)
	require.NoError(t, res)

	// read everything in channel ch
	for len(ch) > 0 {
		data := <-ch
		if data.ctrlIsValid == true && data.ctrl.event == eventProducerIsFullySynced {
			return // ok
		}
	}
	require.Fail(t, "should not have send a eventProducerIsFullySynced in channel")
}

func TestGivenNoSetFirstBlockWhenCallStartThenDontReturnError(t *testing.T) {
	sut, ethermans, _ := setupNoResetCall(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()
	expectedForGettingL1LastBlock(t, ethermans[0], 101)
	err := sut.Start(ctx)
	require.NoError(t, err)
}

func setup(t *testing.T) (*L1RollupInfoProducer, []*L1ParallelEthermanInterfaceMock, chan L1SyncMessage) {
	sut, ethermansMock, resultChannel := setupNoResetCall(t)
	sut.Reset(100)
	return sut, ethermansMock, resultChannel
}

func setupNoResetCall(t *testing.T) (*L1RollupInfoProducer, []*L1ParallelEthermanInterfaceMock, chan L1SyncMessage) {
	ethermansMock := []*L1ParallelEthermanInterfaceMock{NewL1ParallelEthermanInterfaceMock(t), NewL1ParallelEthermanInterfaceMock(t)}
	ethermans := []L1ParallelEthermanInterface{ethermansMock[0], ethermansMock[1]}
	resultChannel := make(chan L1SyncMessage, 100)
	cfg := ConfigProducer{
		SyncChunkSize:      100,
		TtlOfLastBlockOnL1: time.Second,
		TimeOutMainLoop:    time.Second,
	}

	sut := NewL1DataRetriever(cfg, ethermans, resultChannel)
	return sut, ethermansMock, resultChannel
}

func expectedForGettingL1LastBlock(t *testing.T, etherman *L1ParallelEthermanInterfaceMock, blockNumber int64) {
	header := new(ethTypes.Header)
	header.Number = big.NewInt(blockNumber)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Maybe()
}

func expectedRollupInfoCalls(t *testing.T, etherman *L1ParallelEthermanInterfaceMock, calls int) {
	etherman.
		On("GetRollupInfoByBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, nil).
		Times(calls)

	etherman.
		On("EthBlockByNumber", mock.Anything, mock.Anything).
		Return(nil, nil).
		Maybe()
}
