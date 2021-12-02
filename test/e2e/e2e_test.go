package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/mocks"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
)

//nolint:gomnd
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

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/state-transition.json")
	if err != nil {
		t.Error(err)
		return
	}

	// init log
	log.Init(cfg.Log)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			// init database instance
			err = dbutils.InitOrReset(cfg.Database)
			if err != nil {
				t.Error(err)
				return
			}

			// connect to db
			// sqlDB, err := db.NewSQLDB(cfg.Database)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			// prepare merkle tree
			// tr, err := tree.NewReadWriter(sqlDB)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			// create pool
			pl, err := pool.NewPostgresPool(cfg.Database)
			if err != nil {
				t.Error(err)
				return
			}

			// create state
			// st := state.NewState(sqlDB, tr)
			st := mocks.NewState()
			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}
			for _, gacc := range testCase.GenesisAccounts {
				genesis.Balances[common.HexToAddress(gacc.Address)] = &gacc.Balance.Int
			}
			err = st.SetGenesis(ctx, genesis)
			if err != nil {
				t.Error(err)
				return
			}

			// check root
			// root, err := st.GetStateRoot(ctx, true)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// assert.Equal(t, testCase.ExpectedOldRoot, root, "Invalid old root")

			// start sequencer
			// ethManSeq, err := etherman.NewEtherman(cfg.Sequencer.Etherman)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// seq, err := sequencer.NewSequencer(cfg.Sequencer, pl, st, ethManSeq)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// go seq.Start()

			// start synchronizer
			// ethManSync, err := etherman.NewEtherman(cfg.Synchronizer.Etherman)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// sy, err := synchronizer.NewSynchronizer(ethManSync, st, cfg.Synchronizer)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// go func(t *testing.T, s synchronizer.Synchronizer) {
			// 	if err := sy.Sync(); err != nil {
			// 		t.Error(err)
			// 		return
			// 	}
			// }(t, sy)

			// start rpc server
			rpcServer := jsonrpc.NewServer(cfg.RPC, pl, st)
			go func(t *testing.T, s *jsonrpc.Server) {
				if err := s.Start(); err != nil {
					t.Error(err)
					return
				}
			}(t, rpcServer)

			time.Sleep(1 * time.Second)

			// apply transactions
			for _, tx := range testCase.Txs {
				err := sendRawTransaction(tx)
				if err != nil {
					t.Error(err)
					return
				}
			}

			// shutdown rpc server
			if err := rpcServer.Stop(); err != nil {
				t.Error(err)
				return
			}

			// check state against the expected state
			// root, err = st.GetStateRoot(context.Background(), false)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			// assert.Equal(t, testCase.ExpectedNewRoot, root, "invalid new root")
		})
	}
}

func sendRawTransaction(tx vectors.Tx) error {
	endpoint := fmt.Sprintf("http://localhost:%d", cfg.RPC.Port)
	contentType := "application/json"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_sendRawTransaction",
		"params":  []string{tx.RawTx},
	}

	jsonStr, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
