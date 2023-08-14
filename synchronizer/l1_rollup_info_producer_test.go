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

func Test_L1Get(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	header := new(ethTypes.Header)
	header.Number = big.NewInt(150)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Once()

	err := sut.initialize()
	require.NoError(t, err)
	sut.launchWork()
}

func Test_Given_NeedSync_When_Start_Then_AskForRollupInfo(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	expectedForGettingL1LastBlock(t, etherman, 150)
	expectedCalls(t, etherman, 1)
	err := sut.initialize()
	require.NoError(t, err)
	sut.launchWork()
	var waitDuration = time.Duration(0)

	sut.step(&waitDuration)
	sut.workers.waitFinishAllWorkers()
}

func Test_Given_NeedSync_When_ReachLastBlock_Then_Finish(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	expectedForGettingL1LastBlock(t, etherman, 101)
	expectedCalls(t, etherman, 1)
	err := sut.initialize()
	require.NoError(t, err)
	sut.launchWork()
	var waitDuration = time.Duration(0)

	sut.step(&waitDuration)
	sut.workers.waitFinishAllWorkers()
	res := sut.step(&waitDuration)
	require.False(t, res)
}

func setup(t *testing.T) (*l1RollupInfoProducer, []*ethermanMock, chan l1SyncMessage) {
	etherman := newEthermanMock(t)
	ethermansMock := []*ethermanMock{etherman}
	ethermans := []EthermanInterface{etherman}
	resultChannel := make(chan l1SyncMessage, 100)
	sut := newL1DataRetriever(context.Background(), ethermans, 100, 10, resultChannel, false)
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

func expectedCalls(t *testing.T, etherman *ethermanMock, calls int) {
	etherman.
		On("GetRollupInfoByBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, nil).
		Times(calls)

	etherman.
		On("EthBlockByNumber", mock.Anything, mock.Anything).
		Return(nil, nil).
		Maybe()
}
