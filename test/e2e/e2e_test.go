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
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/synchronizer"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	Etherman: etherman.Config{
		PrivateKeyPath:     "../test.keystore",
		PrivateKeyPassword: "testonly",
	},
	RPC: jsonrpc.Config{
		Host: "",
		Port: 8123,
	},
	Synchronizer: synchronizer.Config{},
	Sequencer: sequencer.Config{
		IntervalToProposeBatch: 1 * time.Second,
		URL:                    "https://localhost",
	},
	Aggregator: aggregator.Config{},
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

			//connect to db
			sqlDB, err := db.NewSQLDB(cfg.Database)
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

			// create etherman
			auth, err := newAuthFromKeystore(cfg.Etherman.PrivateKeyPath, cfg.Etherman.PrivateKeyPassword)
			if err != nil {
				t.Error(err)
				return
			}
			etherman, commit, err := etherman.NewSimulatedEtherman(cfg.Etherman, auth)
			if err != nil {
				t.Error(err)
				return
			}

			// create state
			tr, err := tree.NewReadWriter(sqlDB)
			if err != nil {
				t.Error(err)
				return
			}
			st := state.NewState(sqlDB, tr)
			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}
			for _, gacc := range testCase.GenesisAccounts {
				b := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &b
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
			// expectedOldRoot, ok := big.NewInt(0).SetString(testCase.ExpectedOldRoot, 10)
			// if !ok {
			// 	t.Error(fmt.Errorf("Failed to read ExpectedOldRoot"))
			// 	return
			// }
			// assert.Equal(t, expectedOldRoot.Cmp(root), 0, "Invalid old root")

			// start synchronizer
			sy, err := synchronizer.NewSynchronizer(etherman, st, cfg.Synchronizer)
			if err != nil {
				t.Error(err)
				return
			}
			go func(t *testing.T, s synchronizer.Synchronizer) {
				if err := sy.Sync(); err != nil {
					t.Error(err)
					return
				}
			}(t, sy)

			// start sequencer
			_, err = etherman.PoE.RegisterSequencer(auth, cfg.Sequencer.URL)
			if err != nil {
				t.Error(err)
				return
			}
			// mine next block with sequencer registration
			commit()

			// wait sequencer registration to be synchronized
			time.Sleep(3 * time.Second)

			seq, err := sequencer.NewSequencer(cfg.Sequencer, pl, st, etherman)
			if err != nil {
				t.Error(err)
				return
			}
			go seq.Start()

			// start rpc server
			stSeq, err := st.GetSequencer(ctx, cfg.Sequencer.URL)
			if err != nil {
				t.Error(err)
				return
			}
			rpcServer := jsonrpc.NewServer(cfg.RPC, stSeq.ChainID.Uint64(), pl, st)
			go func(t *testing.T, s *jsonrpc.Server) {
				if err := s.Start(); err != nil {
					t.Error(err)
					return
				}
			}(t, rpcServer)

			// wait RPC server to be ready
			time.Sleep(1 * time.Second)

			// apply transactions
			for _, tx := range testCase.Txs {
				err := sendRawTransaction(tx)
				if err != nil {
					t.Error(err)
					return
				}
			}

			// wait for sequencer to select txs from pool and propose a new batch
			time.Sleep(3 * time.Second)

			// mine next block with batch propostal
			commit()

			// wait for the synchronizer to update state
			time.Sleep(3 * time.Second)

			// shutdown rpc server
			if err := rpcServer.Stop(); err != nil {
				t.Error(err)
				return
			}

			// stop synchronizer
			sy.Stop()

			// stop sequencer
			seq.Stop()

			// check state against the expected state
			// root, err = st.GetStateRoot(ctx, true)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }
			// expectedNewRoot, ok := big.NewInt(0).SetString(testCase.ExpectedNewRoot, 10)
			// if !ok {
			// 	t.Error(fmt.Errorf("Failed to read ExpectedNewRoot"))
			// 	return
			// }
			// assert.Equal(t, expectedNewRoot.Cmp(root), 0, "Invalid new root")

			// check leafs
			batchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, 0, leaf.Balance.Cmp(actualBalance))

				actualNonce, err := st.GetNonce(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, leaf.Nonce, actualNonce)
			}
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

func newAuthFromKeystore(path, password string) (*bind.TransactOpts, error) {
	if path == "" && password == "" {
		log.Info("lol")
		return nil, nil
	}
	keystoreEncrypted, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keystoreEncrypted, password)
	if err != nil {
		return nil, err
	}
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(1337)) //nolint:gomnd
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = 99999999999
	return auth, nil
}
