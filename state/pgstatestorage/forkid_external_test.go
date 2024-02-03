package pgstatestorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddForkIDInterval(t *testing.T) {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
	pgStateStorage = pgstatestorage.NewPostgresStorage(state.Config{}, stateDb)
	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), executorClient, stateTree, nil, nil)

	for i := 1; i <= 6; i++ {
		err = testState.AddForkID(ctx, state.ForkIDInterval{FromBatchNumber: uint64(i * 10), ToBatchNumber: uint64(i*10) + 9, ForkId: uint64(i)}, nil)
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		forkIDToAdd   state.ForkIDInterval
		expectedError error
	}{
		{
			name:          "fails to add because forkID already exists",
			forkIDToAdd:   state.ForkIDInterval{ForkId: 3},
			expectedError: fmt.Errorf("error checking forkID sequence. Last ForkID stored: 6. New ForkID received: 3"),
		},
		{
			name:          "fails to add because forkID is smaller than the latest forkID",
			forkIDToAdd:   state.ForkIDInterval{ForkId: 5},
			expectedError: fmt.Errorf("error checking forkID sequence. Last ForkID stored: 6. New ForkID received: 5"),
		},
		{
			name:          "fails to add because forkID is equal to the latest forkID",
			forkIDToAdd:   state.ForkIDInterval{ForkId: 6},
			expectedError: fmt.Errorf("error checking forkID sequence. Last ForkID stored: 6. New ForkID received: 6"),
		},
		{
			name:          "adds successfully",
			forkIDToAdd:   state.ForkIDInterval{ForkId: 7},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			dbTx, err := testState.BeginStateTransaction(ctx)
			require.NoError(t, err)

			err = testState.AddForkIDInterval(ctx, tc.forkIDToAdd, dbTx)

			if tc.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			}

			require.NoError(t, dbTx.Commit(ctx))
		})
	}
}
