package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
)

var cfg = config.Config{
	Log: log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	},
	Database: db.Config{
		Database: "polygon-hermez",
		User:     "hermez",
		Password: "polygon",
		Host:     "localhost",
		Port:     "5432",
	},
	RPC: jsonrpc.Config{
		Host: "",
		Port: 8123,

		ChainID: 2576980377, // 0x99999999,
	},
	Synchronizer: synchronizer.Config{
		Etherman: etherman.Config{},
	},
	Sequencer: sequencer.Config{
		IntervalToProposeBatch: 15 * time.Second,
		Etherman:               etherman.Config{},
	},
	Aggregator: aggregator.Config{
		Etherman: etherman.Config{},
	},
}

func TestStateTransition(t *testing.T) {
	// load vector
	vector, err := loadVector()
	if err != nil {
		t.Error(err)
		return
	}

	// init log
	log.Init(cfg.Log)

	for _, testCase := range vector.StateTests {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			// init database instance
			err = dbutils.InitOrReset(cfg.Database)
			if err != nil {
				t.Error(err)
				return
			}

			// connect to db
			sqlDB, err := db.NewSQLDB(cfg.Database)
			if err != nil {
				t.Error(err)
				return
			}

			// prepare merkle tree
			tr, err := tree.NewReadWriter(sqlDB)
			if err != nil {
				t.Error(err)
				return
			}

			// create pool
			pl, err := pool.NewPostgresPool(cfg.Database)
			if err != nil {
				t.Error(err)
				return
			}

			// create state
			st := state.NewState(sqlDB, tr)
			genesis := state.Genesis{}
			for _, gacc := range testCase.GenesisAccounts {
				genesis.Balances[gacc.Address.Address()] = &gacc.Balance.Int
			}
			st.SetGenesis(genesis)

			// check root
			root, err := st.GetStateRoot(ctx, true)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, testCase.ExpectedOldRoot, root, "Invalid old root")

			// start sequencer
			ethManSeq, err := etherman.NewEtherman(cfg.Sequencer.Etherman)
			if err != nil {
				t.Error(err)
				return
			}
			seq, err := sequencer.NewSequencer(cfg.Sequencer, pl, st, ethManSeq)
			if err != nil {
				t.Error(err)
				return
			}
			go seq.Start()

			// start synchronizer
			ethManSync, err := etherman.NewEtherman(cfg.Synchronizer.Etherman)
			if err != nil {
				t.Error(err)
				return
			}
			sy, err := synchronizer.NewSynchronizer(ethManSync, st, cfg.Synchronizer)
			if err != nil {
				t.Error(err)
				return
			}
			go sy.Sync()

			// start rpc server
			rpcServer := jsonrpc.NewServer(cfg.RPC, pl, st)
			go rpcServer.Start()

			// apply transactions
			for _, tx := range testCase.Txs {
				err := sendRawTransaction(tx)
				if err != nil {
					t.Error(err)
					return
				}
			}

			// check state against the expected state
			root, err = st.GetStateRoot(context.Background(), false)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, testCase.ExpectedNewRoot, root, "invalid new root")
		})
	}
}

func loadVector() (StateTransitionVector, error) {
	var vector StateTransitionVector

	jsonFile, err := os.Open("state-transition.json")
	if err != nil {
		return vector, err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return vector, err
	}

	err = json.Unmarshal(bytes, &vector)
	if err != nil {
		return vector, err
	}

	return vector, nil
}

func sendRawTransaction(tx Tx) error {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", cfg.RPC.Port))
	if err != nil {
		return err
	}

	transaction := types.NewTransaction(tx.Nonce, tx.To.Address(), &tx.Value.Int, tx.GasLimit, &tx.GasPrice.Int, []byte{})
	b, err := transaction.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	args := []interface{}{encoded}

	var result map[string]interface{}

	err = client.Call("eth_sendRawTransaction", args, &result)
	if err != nil {
		return err
	}

	return nil
}
