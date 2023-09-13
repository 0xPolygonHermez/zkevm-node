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
	ctx      context.Context
}

func Test_Given_Consumer_When_ReceiveAFullSyncAndChannelIsEmpty_Then_StopOk(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t, &ctxTimeout)
	defer cancel()
	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	err := data.sut.start()
	require.NoError(t, err)
}
func Test_Given_Consumer_When_ReceiveAFullSyncAndChannelIsNotEmpty_Then_DontStop(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t, &ctxTimeout)
	defer cancel()

	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	data.ch <- *newL1SyncMessageControl(eventNone)
	err := data.sut.start()
	require.Error(t, err)
	require.Equal(t, errContextCanceled, err.Error())
}

func Test_Given_Consumer_When_FailsToProcessRollup_Then_IDontKnown(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t, &ctxTimeout)
	defer cancel()
	responseRollupInfoByBlockRange := responseRollupInfoByBlockRange{
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
	err := data.sut.start()
	require.Error(t, err)
	_, ok := data.sut.getLastEthBlockSynced()
	require.False(t, ok)
}

func setupConsumerTest(t *testing.T, ctx *context.Context) consumerTestData {
	syncMock := newSynchronizerProcessBlockRangeMock(t)
	ch := make(chan l1SyncMessage, 10)
	if ctx == nil {
		rctx := context.Background()
		ctx = &rctx
	}
	sut := newL1RollupInfoConsumer(syncMock, *ctx, ch)
	return consumerTestData{sut, syncMock, ch, *ctx}
}
