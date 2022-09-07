package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"

	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	st "github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func init() {

}

func TestSequenceTooBig(t *testing.T) {

	ctx := context.Background()
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	eth_man, _, _, _, err := ethman.NewSimulatedEtherman(ethman.Config{}, auth)
	require.NoError(t, err)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type: "default",
	})
	require.NoError(t, err)

	dbConfig := db.Config{
		User:      "test_user",
		Password:  "test_password",
		Name:      "test_db",
		Host:      "localhost",
		Port:      "5432",
		EnableLog: false,
		MaxConns:  200,
	}
	err = dbutils.InitOrReset(dbConfig)
	require.NoError(t, err)

	poolDb, err := pgpoolstorage.NewPostgresPoolStorage(dbConfig)
	require.NoError(t, err)

	sqlDB, err := db.NewSQLDB(dbConfig)
	require.NoError(t, err)

	stateDb := st.NewPostgresStorage(sqlDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, executor.Config{
		URI: "localhost:50071",
	})
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, merkletree.Config{
		URI: "localhost:50071",
	})
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := st.Config{
		MaxCumulativeGasUsed: 30000000,
		ChainID:              1000,
	}

	state := st.NewState(stateCfg, stateDb, executorClient, stateTree)

	pool := pool.NewPool(poolDb, state, common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"))
	ethtxmanager := ethtxmanager.New(ethtxmanager.Config{}, eth_man)

	seq, err := New(Config{
		MaxSequenceSize: big.NewInt(1000000),
	}, pool, state, eth_man, pg, ethtxmanager)
	require.NoError(t, err)

	// generate fake data

	dbTx, err := state.BeginStateTransaction(ctx)
	require.NoError(t, err)

	state.PostgresStorage.Pool.Exec(ctx, "INSERT INTO state.batch VALUES(batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data) ")

	err = dbTx.Commit(ctx)
	require.NoError(t, err)
	sequences, err := seq.getSequencesToSend(ctx)
	require.NoError(t, err)
	fmt.Printf("\n%+v", sequences)
}
