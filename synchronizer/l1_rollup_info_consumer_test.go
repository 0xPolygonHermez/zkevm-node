package synchronizer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type consumerTestData struct {
	sut      *l1RollupInfoConsumer
	syncMock *synchronizerProcessBlockRangeMock
	ch       chan l1SyncMessage
}

func TestGivenConsumerWhenReceiveAFullSyncAndChannelIsEmptyThenStopOk(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
}
func TestGivenConsumerWhenReceiveAFullSyncAndChannelIsNotEmptyThenDontStop(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()

	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	data.ch <- *newL1SyncMessageControl(eventNone)
	err := data.sut.Start(ctxTimeout, nil)
	require.Error(t, err)
	require.Equal(t, errContextCanceled, err)
}

func TestGivenConsumerWhenFailsToProcessRollupThenDontKnownLastEthBlock(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: nil,
	}
	data.syncMock.
		On("processBlockRange", mock.Anything, mock.Anything).
		Return(errors.New("error")).
		Once()
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	err := data.sut.Start(ctxTimeout, nil)
	require.Error(t, err)
	_, ok := data.sut.GetLastEthBlockSynced()
	require.False(t, ok)
}

func setupConsumerTest(t *testing.T) consumerTestData {
	syncMock := newSynchronizerProcessBlockRangeMock(t)
	ch := make(chan l1SyncMessage, 10)

	cfg := configConsumer{
		numIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData: minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData,
		acceptableTimeWaitingForNewRollupInfoData:                       minAcceptableTimeWaitingForNewRollupInfoData,
	}
	sut := newL1RollupInfoConsumer(cfg, syncMock, ch)
	return consumerTestData{sut, syncMock, ch}
}
