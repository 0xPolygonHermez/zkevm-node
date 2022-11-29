package sequencer

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
	st "github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSequenceTooBig(t *testing.T) {
	// before running:
	// make run-db
	// make run-network
	// make run-zkprover

	const (
		CONFIG_MAX_GAS_PER_SEQUENCE     = 200000
		CONFIG_ENCRYPTION_KEY_FILE_PATH = "./../test/test.keystore"
		CONFIG_ENCRYPTION_KEY_PASSWORD  = "testonly"
		CONFIG_CHAIN_ID                 = 1337
		CONFIG_ETH_URL                  = "http://localhost:8545"

		CONFIG_NAME_POE   = "poe"
		CONFIG_NAME_MATIC = "matic"
		CONFIG_NAME_GER   = "ger"
	)

	var (
		CONFIG_ADDRESSES = map[string]common.Address{
			CONFIG_NAME_POE:   common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"), // <= PoE
			CONFIG_NAME_MATIC: common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), // <= Matic
			CONFIG_NAME_GER:   common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"), // <= GER
		}
		CONFIG_DB_STATE = db.Config{
			User:      "state_user",
			Password:  "state_password",
			Name:      "state_db",
			Host:      "localhost",
			Port:      "5432",
			EnableLog: false,
			MaxConns:  200,
		}
		CONFIG_DB_POOL = db.Config{
			User:      "pool_user",
			Password:  "pool_password",
			Name:      "pool_db",
			Host:      "localhost",
			Port:      "5433",
			EnableLog: false,
			MaxConns:  200,
		}
		CONFIG_EXECUTOR_URL = fmt.Sprintf("%s:50071", testutils.GetEnv("ZKPROVER_URI", "localhost"))
	)
	type TestCase struct {
		Input  []int // slice of batch sizes
		Output int   // split into N sequences

	}

	var testcases = []TestCase{
		{
			Input: []int{
				1000,
				500,
			},
			Output: 2, // two sequences (of 1 batch each) fit inside
		},

		{
			Input: []int{
				1,
			},
			Output: 1, // only one sequence fits
		},
		{
			Input: []int{
				100000000,
				1000000,
				1000,
				100,
				1,
			},
			Output: 2, // only two sequences fit inside
		},
		{
			Input: []int{
				1, 1, 1, 1,
			},
			Output: 2, // all sequences fit inside
		},
	}
	ctx := context.Background()

	keystoreEncrypted, err := ioutil.ReadFile(CONFIG_ENCRYPTION_KEY_FILE_PATH)
	require.NoError(t, err)
	key, err := keystore.DecryptKey(keystoreEncrypted, CONFIG_ENCRYPTION_KEY_PASSWORD)
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(CONFIG_CHAIN_ID))
	require.NoError(t, err)
	//	eth_man, _, _, _, err := ethman.NewSimulatedEtherman(ethman.Config{}, auth)
	eth_man, err := ethman.NewClient(ethman.Config{
		URL:                       CONFIG_ETH_URL,
		L1ChainID:                 CONFIG_CHAIN_ID,
		PoEAddr:                   CONFIG_ADDRESSES[CONFIG_NAME_POE],
		MaticAddr:                 CONFIG_ADDRESSES[CONFIG_NAME_MATIC],
		GlobalExitRootManagerAddr: CONFIG_ADDRESSES[CONFIG_NAME_GER],
	}, auth)

	require.NoError(t, err)

	const decimals = 1000000000000000000
	amount := big.NewFloat(10000000000000000)
	amountInWei := new(big.Float).Mul(amount, big.NewFloat(decimals))
	amountB := new(big.Int)
	amountInWei.Int(amountB)

	_, err = eth_man.ApproveMatic(amountB, CONFIG_ADDRESSES[CONFIG_NAME_POE])
	require.NoError(t, err)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type: "default",
	})
	require.NoError(t, err)

	err = dbutils.InitOrResetState(CONFIG_DB_STATE)
	require.NoError(t, err)

	err = dbutils.InitOrResetPool(CONFIG_DB_POOL)
	require.NoError(t, err)

	poolDb, err := pgpoolstorage.NewPostgresPoolStorage(CONFIG_DB_POOL)
	require.NoError(t, err)

	sqlStateDB, err := db.NewSQLDB(CONFIG_DB_STATE)
	require.NoError(t, err)

	stateDb := st.NewPostgresStorage(sqlStateDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, executor.Config{
		URI: CONFIG_EXECUTOR_URL,
	})
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, merkletree.Config{
		URI: CONFIG_EXECUTOR_URL,
	})
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := st.Config{
		MaxCumulativeGasUsed: 30000000,
		ChainID:              CONFIG_CHAIN_ID,
	}

	state := st.NewState(stateCfg, stateDb, executorClient, stateTree)

	pool := pool.NewPool(poolDb, state, CONFIG_ADDRESSES[CONFIG_NAME_GER], big.NewInt(CONFIG_CHAIN_ID).Uint64())
	ethtxmanager := ethtxmanager.New(ethtxmanager.Config{}, eth_man, state)
	gpe := gasprice.NewDefaultEstimator(gasprice.Config{
		Type:               gasprice.DefaultType,
		DefaultGasPriceWei: 1000000000,
	}, pool)
	seq, err := New(Config{
		MaxSequenceSize:                          MaxSequenceSize{Int: big.NewInt(CONFIG_MAX_GAS_PER_SEQUENCE)},
		LastBatchVirtualizationTimeMaxWaitPeriod: types.NewDuration(1 * time.Second),
		ProfitabilityChecker: profitabilitychecker.Config{
			SendBatchesEvenWhenNotProfitable: true,
		},
	}, pool, state, eth_man, pg, ethtxmanager, gpe)
	require.NoError(t, err)

	// generate fake data

	mainnetExitRoot := common.HexToHash("caffe")
	rollupExitRoot := common.HexToHash("bead")

	if _, err := stateDb.Exec(ctx, "DELETE FROM state.block"); err != nil {
		t.Fail()
	}
	if _, err := stateDb.Exec(ctx, "DELETE FROM state.batch"); err != nil {
		t.Fail()
	}
	if _, err := stateDb.Exec(ctx, "DELETE FROM state.exit_root"); err != nil {
		t.Fail()
	}

	const sqlAddBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	_, err = stateDb.Exec(ctx, sqlAddBlock, 1, time.Now(), "")
	require.NoError(t, err)

	_, err = stateDb.Exec(ctx, sqlAddBlock, 2, time.Now(), "") // for use in lastVirtualized time
	require.NoError(t, err)

	const sqlAddExitRoots = "INSERT INTO state.exit_root (block_num, global_exit_root, mainnet_exit_root, rollup_exit_root, global_exit_root_num) VALUES ($1, $2, $3, $4, $5)"
	_, err = stateDb.Exec(ctx, sqlAddExitRoots, 1, common.Address{}, mainnetExitRoot, rollupExitRoot, 3)
	require.NoError(t, err)

	for _, testCase := range testcases {
		innerDbTx, err := state.BeginStateTransaction(ctx)
		require.NoError(t, err)
		err = dbutils.InitOrResetState(CONFIG_DB_STATE)
		require.NoError(t, err)

		err = dbutils.InitOrResetPool(CONFIG_DB_POOL)
		require.NoError(t, err)

		if _, err := stateDb.Exec(ctx, "DELETE FROM state.block"); err != nil {
			t.Fail()
		}
		if _, err := stateDb.Exec(ctx, "DELETE FROM state.batch"); err != nil {
			t.Fail()
		}

		for i := 0; i < len(testCase.Input); i++ {
			fmt.Printf("\niteration: [%d]: %d\n", testCase.Output, testCase.Input[i])

			payload := make([]byte, testCase.Input[i]) // 10mb

			_, err = stateDb.Exec(ctx, "INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				i+1,
				common.Address{}.String(),
				common.Hash{}.String(),
				common.Hash{}.String(),
				time.Unix(9, 0).UTC(),
				common.HexToAddress("").String(),
				payload,
			)
			require.NoError(t, err)
		}

		//needed for completion: wip batch

		_, err = stateDb.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)",
			len(testCase.Input)+1,
		)

		require.NoError(t, err)

		// make L2 equivalences

		err = innerDbTx.Commit(ctx)
		require.NoError(t, err)

		sequences, err := seq.getSequencesToSend(ctx)
		require.NoError(t, err)

		fmt.Printf("%+v", sequences)

		require.Equal(t, testCase.Output, len(sequences))
	}
}
