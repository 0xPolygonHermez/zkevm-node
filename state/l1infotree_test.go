package state_test

import (
	"context"
	"math"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestFirstLeafOfL1InfoTreeIsIndex0(t *testing.T) {
	stateDBCfg := dbutils.NewStateConfigFromEnv()
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}

	stateDb, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	forkID := uint64(state.FORKID_ETROG)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
	ctx := context.Background()
	storage := pgstatestorage.NewPostgresStorage(stateCfg, stateDb)
	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	testState := state.NewState(stateCfg, storage, nil, nil, nil, mt)
	dbTx, err := testState.BeginStateTransaction(ctx)
	defer func() {
		_ = dbTx.Rollback(ctx)
	}()
	require.NoError(t, err)
	block := state.Block{BlockNumber: 123}
	err = testState.AddBlock(ctx, &block, dbTx)
	require.NoError(t, err)

	leaf := state.L1InfoTreeLeaf{
		GlobalExitRoot: state.GlobalExitRoot{
			GlobalExitRoot: common.Hash{},
			BlockNumber:    123,
		},
		PreviousBlockHash: common.Hash{},
	}
	insertedLeaf, err := testState.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
	require.NoError(t, err)
	require.Equal(t, insertedLeaf.L1InfoTreeIndex, uint32(0))
}
