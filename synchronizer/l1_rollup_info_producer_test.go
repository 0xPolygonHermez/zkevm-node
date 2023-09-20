package synchronizer

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
	sut.launchWork()
}

func TestGivenNeedSyncWhenStartThenAskForRollupInfo(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	expectedForGettingL1LastBlock(t, etherman, 150)
	expectedRollupInfoCalls(t, etherman, 1)
	err := sut.initialize(context.Background())
	require.NoError(t, err)
	sut.launchWork()
	var waitDuration = time.Duration(0)

	sut.stepInner(&waitDuration)
	sut.workers.waitFinishAllWorkers()
}

func TestGivenNoNeedSyncWhenStartsSendAndEventOfSynchronized(t *testing.T) {
	sut, ethermans, ch := setup(t)
	etherman := ethermans[0]
	// Our last block is 100 in DB and it returns 100 as last block on L1
	// so is synchronized
	expectedForGettingL1LastBlock(t, etherman, 100)
	//expectedRollupInfoCalls(t, etherman, 1)
	err := sut.initialize(context.Background())
	require.NoError(t, err)
	sut.launchWork()
	var waitDuration = time.Duration(0)

	sut.step(&waitDuration)

	waitDuration = time.Duration(0)
	res := sut.step(&waitDuration)
	require.True(t, res)
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
	etherman := ethermans[0]
	// Our last block is 100 in DB and it returns 101 as last block on L1
	// so it need to retrieve 1 rollupinfo
	expectedForGettingL1LastBlock(t, etherman, 101)
	expectedRollupInfoCalls(t, etherman, 1)
	err := sut.initialize(context.Background())
	require.NoError(t, err)
	var waitDuration = time.Duration(0)

	// Is going to ask for last block again because it'll launch all request
	expectedForGettingL1LastBlock(t, etherman, 101)
	sut.step(&waitDuration)
	require.Equal(t, sut.status, producerWorking)
	waitDuration = time.Millisecond * 100 // need a bit of time to receive the response to rollupinfo
	res := sut.step(&waitDuration)
	require.True(t, res)
	require.Equal(t, sut.status, producerSynchronized)
	// read everything in channel ch
	for len(ch) > 0 {
		data := <-ch
		if data.ctrlIsValid == true && data.ctrl.event == eventProducerIsFullySynced {
			return // ok
		}
	}
	require.Fail(t, "should not have send a eventProducerIsFullySynced in channel")
}

func TestGivenNoSetFirstBlockWhenCallStartThenReturnError(t *testing.T) {
	sut, _, _ := setupNoResetCall(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := sut.Start(ctx)
	require.Error(t, err)
	require.Equal(t, errStartingBlockNumberMustBeDefined, err)
}

func setup(t *testing.T) (*l1RollupInfoProducer, []*ethermanMock, chan l1SyncMessage) {
	sut, ethermansMock, resultChannel := setupNoResetCall(t)
	sut.ResetAndStop(100)
	return sut, ethermansMock, resultChannel
}

func setupNoResetCall(t *testing.T) (*l1RollupInfoProducer, []*ethermanMock, chan l1SyncMessage) {
	etherman := newEthermanMock(t)
	ethermansMock := []*ethermanMock{etherman}
	ethermans := []EthermanInterface{etherman}
	resultChannel := make(chan l1SyncMessage, 100)
	cfg := configProducer{
		syncChunkSize:      100,
		ttlOfLastBlockOnL1: time.Second,
	}

	sut := newL1DataRetriever(cfg, ethermans, resultChannel)
	return sut, ethermansMock, resultChannel
}

func expectedForGettingL1LastBlock(t *testing.T, etherman *ethermanMock, blockNumber int64) {
	header := new(ethTypes.Header)
	header.Number = big.NewInt(blockNumber)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Maybe()
}

func expectedRollupInfoCalls(t *testing.T, etherman *ethermanMock, calls int) {
	etherman.
		On("GetRollupInfoByBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, nil).
		Times(calls)

	etherman.
		On("EthBlockByNumber", mock.Anything, mock.Anything).
		Return(nil, nil).
		Maybe()
}
