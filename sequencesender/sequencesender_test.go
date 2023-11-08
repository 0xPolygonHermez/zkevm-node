package sequencesender

import (
	"context"
	"fmt"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	state_interface    = new(stateMock)
	etherman_interface = new(ethermanMock)
	ethtxman_interface = new(ethTxManagerMock)
	ctx                context.Context
	cfg                = Config{
		WaitPeriodSendSequence: cfgTypes.Duration{
			Duration: 5,
		},
		LastBatchVirtualizationTimeMaxWaitPeriod: cfgTypes.Duration{
			Duration: 5,
		},
		MaxTxSizeForL1: 10,
		PrivateKey: cfgTypes.KeystoreFileConfig{
			Path:     "../test/sequencer.keystore",
			Password: "testonly",
		},
	}

	addr1 = common.Address{0x1}
	addr2 = common.Address{0x2}
)

func TestSequenceSender_getSequencesToSend(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(t, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	sequenceSender, err := New(cfg, state_interface, etherman_interface, ethtxman_interface, eventLog)
	require.NoError(t, err)
	ctx = context.Background()

	state_interface.On("GetTimeForLatestBatchVirtualization", ctx, nil).Return(func(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
		return time.Now().Add(-cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration), nil
	})
	timeSpeical, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 PST 2015")
	require.NoError(t, err)
	testCases := []struct {
		name           string
		times          int
		prepareBatches func(ctx context.Context, number uint64, dbTx pgx.Tx) (*state.Batch, error)
		isBatchClosed  func(ctx context.Context, number uint64, dbTx pgx.Tx) (bool, error)
		expectFunc     func(ctx context.Context, t *testing.T)
	}{
		{
			name:  "only one batch push into sequence",
			times: 1,
			prepareBatches: func(ctx context.Context, number uint64, dbTx pgx.Tx) (*state.Batch, error) {
				if number > 1 {
					return nil, fmt.Errorf("The batch %d is not exist", number)
				}

				return &state.Batch{
					BatchNumber: number,
					Timestamp:   timeSpeical,
					Coinbase:    addr1,
				}, nil
			},
			isBatchClosed: func(ctx context.Context, number uint64, dbTx pgx.Tx) (bool, error) {
				if number > 1 {
					return false, nil
				}
				return true, nil
			},
			expectFunc: func(ctx context.Context, t *testing.T) {
				sequence, coinbase, err := sequenceSender.getSequencesToSend(ctx)
				require.Equal(t, addr1, coinbase)
				require.Equal(t, []types.Sequence{
					{BatchNumber: 1, Timestamp: timeSpeical.Unix()},
				}, sequence)
				require.NoError(t, err)
			},
		},
		{
			name:  "2 batch which has same batch push into sequence",
			times: 2,
			prepareBatches: func(ctx context.Context, number uint64, dbTx pgx.Tx) (*state.Batch, error) {
				if number > 2 {
					return nil, fmt.Errorf("The batch %d is not exist", number)
				}

				return &state.Batch{
					BatchNumber: number,
					Timestamp:   timeSpeical,
					Coinbase:    addr1,
				}, nil
			},
			isBatchClosed: func(ctx context.Context, number uint64, dbTx pgx.Tx) (bool, error) {
				if number > 2 {
					return false, nil
				}
				return true, nil
			},
			expectFunc: func(ctx context.Context, t *testing.T) {
				sequence, coinbase, err := sequenceSender.getSequencesToSend(ctx)
				require.Equal(t, addr1, coinbase)
				require.Equal(t, []types.Sequence{
					{BatchNumber: 1, Timestamp: timeSpeical.Unix()},
					{BatchNumber: 2, Timestamp: timeSpeical.Unix()},
				}, sequence)
				require.NoError(t, err)
			},
		},
		{
			name:  "2 batch which has same coinbase and another 2 batch which has different from first batches  push into sequence",
			times: 5,
			prepareBatches: func(ctx context.Context, number uint64, dbTx pgx.Tx) (*state.Batch, error) {
				if number > 4 {
					return nil, fmt.Errorf("The batch %d is not exist", number)
				}
				if number <= 2 {
					return &state.Batch{
						BatchNumber: number,
						Timestamp:   timeSpeical,
						Coinbase:    addr1,
					}, nil
				} else {
					return &state.Batch{
						BatchNumber: number,
						Timestamp:   timeSpeical,
						Coinbase:    addr2,
					}, nil
				}
			},
			isBatchClosed: func(ctx context.Context, number uint64, dbTx pgx.Tx) (bool, error) {
				if number > 4 {
					return false, nil
				}
				return true, nil
			},
			expectFunc: func(ctx context.Context, t *testing.T) {
				sequence, coinbase, err := sequenceSender.getSequencesToSend(ctx)
				require.Equal(t, addr1, coinbase)
				require.Equal(t, []types.Sequence{
					{BatchNumber: 1, Timestamp: timeSpeical.Unix()},
					{BatchNumber: 2, Timestamp: timeSpeical.Unix()},
				}, sequence)
				require.NoError(t, err)
				state_interface.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(2), nil).Once()

				sequence, coinbase, err = sequenceSender.getSequencesToSend(ctx)
				require.Equal(t, addr2, coinbase)
				require.Equal(t, []types.Sequence{
					{BatchNumber: 3, Timestamp: timeSpeical.Unix()},
					{BatchNumber: 4, Timestamp: timeSpeical.Unix()},
				}, sequence)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state_interface.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(0), nil).Once()
			state_interface.On("GetBatchByNumber", ctx, mock.AnythingOfType("uint64"), nil).Return(tc.prepareBatches).Times(tc.times)
			state_interface.On("IsBatchClosed", ctx, mock.AnythingOfType("uint64"), nil).Return(tc.isBatchClosed).Times(tc.times + 1)
			etherman_interface.On("EstimateGasSequenceBatches", mock.Anything, mock.Anything, mock.Anything).Return(types2.NewTx(&types2.LegacyTx{}), nil)

			tc.expectFunc(ctx, t)

			state_interface.AssertExpectations(t)
			etherman_interface.AssertExpectations(t)
		})
	}
}
