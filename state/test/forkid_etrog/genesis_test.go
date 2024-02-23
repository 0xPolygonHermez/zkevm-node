package etrog_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	test "github.com/0xPolygonHermez/zkevm-node/state/test/forkid_common"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenesisTimestamp(t *testing.T) {
	ctx := context.Background()
	genesis := state.Genesis{}

	err := dbutils.InitOrResetState(test.StateDBCfg)
	require.NoError(t, err)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	timeStamp := time.Now()
	block := state.Block{ReceivedAt: timeStamp}

	_, err = testState.SetGenesis(ctx, block, genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)

	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	batchTimeStamp, err := testState.GetBatchTimestamp(ctx, 0, nil, nil)
	require.NoError(t, err)

	assert.Equal(t, 0, timeStamp.Compare(*batchTimeStamp))
}
