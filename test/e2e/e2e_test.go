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
	"strconv"
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

const (
	sequencerURL = "http://localhost"
)

//nolint:gomnd
var cfg = config.Config{
	Log: log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	},
	Database: db.Config{
		Name:     "polygon-hermez",
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
		IntervalToProposeBatch: sequencer.Duration{Duration: 1 * time.Second},
	},
	Aggregator: aggregator.Config{},
}

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/state-transition.json")
	require.NoError(t, err)

	// init log
	log.Init(cfg.Log)

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			// init database instance
			err = dbutils.InitOrReset(cfg.Database)
			require.NoError(t, err)

			//connect to db
			sqlDB, err := db.NewSQLDB(cfg.Database)
			require.NoError(t, err)

			// create pool
			pl, err := pool.NewPostgresPool(cfg.Database)
			require.NoError(t, err)

			// create etherman
			auth, err := newAuthFromKeystore(cfg.Etherman.PrivateKeyPath, cfg.Etherman.PrivateKeyPassword)
			require.NoError(t, err)

			etherman, commit, err := etherman.NewSimulatedEtherman(cfg.Etherman, auth)
			require.NoError(t, err)

			// create state
			tr := tree.NewStateTree(sqlDB, []byte{})
			st := state.NewState(sqlDB, tr)
			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}
			for _, gacc := range testCase.GenesisAccounts {
				b := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &b
			}
			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			// check initial root
			root, err := st.GetStateRoot(ctx, true)
			require.NoError(t, err)

			strRoot := new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedOldRoot, strRoot, "Invalid old root")

			// start synchronizer
			sy, err := synchronizer.NewSynchronizer(etherman, st, cfg.Synchronizer)
			require.NoError(t, err)
			go func(t *testing.T, s synchronizer.Synchronizer) {
				err := sy.Sync()
				require.NoError(t, err)
			}(t, sy)

			// start sequencer
			_, err = etherman.PoE.RegisterSequencer(auth, sequencerURL)
			require.NoError(t, err)

			// mine next block with sequencer registration
			commit()

			// wait sequencer registration to be synchronized
			time.Sleep(10 * time.Second)

			// create sequencer
			seq, err := sequencer.NewSequencer(cfg.Sequencer, pl, st, etherman)
			require.NoError(t, err)
			go seq.Start()

			// start rpc server
			key, err := newKeyFromKeystore(cfg.Etherman.PrivateKeyPath, cfg.Etherman.PrivateKeyPassword)
			require.NoError(t, err)
			stSeq, err := st.GetSequencer(ctx, key.Address)
			require.NoError(t, err)

			rpcServer := jsonrpc.NewServer(cfg.RPC, stSeq.ChainID.Uint64(), pl, st)
			go func(t *testing.T, s *jsonrpc.Server) {
				err := s.Start()
				require.NoError(t, err)
			}(t, rpcServer)

			// wait RPC server to be ready
			time.Sleep(10 * time.Second)

			// apply transactions
			for _, tx := range testCase.Txs {
				err := sendRawTransaction(tx)
				require.NoError(t, err)
			}

			// wait for sequencer to select txs from pool and propose a new batch
			time.Sleep(10 * time.Second)

			// mine next block with batch propostal
			commit()

			// wait for the synchronizer to update state
			time.Sleep(10 * time.Second)

			// shutdown rpc server
			err = rpcServer.Stop()
			require.NoError(t, err)

			// stop synchronizer
			sy.Stop()

			// stop sequencer
			seq.Stop()

			// check state against the expected state
			root, err = st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot = new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot, "Invalid new root")

			// check leafs
			batchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			for addrStr, leaf := range testCase.ExpectedNewLeafs {
				addr := common.HexToAddress(addrStr)

				actualBalance, err := st.GetBalance(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("leaf.Balance: %s actualBalance: %s", leaf.Balance.Text(10), actualBalance.Text(10)))

				actualNonce, err := st.GetNonce(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, 10), fmt.Sprintf("leaf.Nonce: %s actualNonce: %d", leaf.Nonce, actualNonce))
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

func newKeyFromKeystore(path, password string) (*keystore.Key, error) {
	if path == "" && password == "" {
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
	return key, nil
}

func newAuthFromKeystore(path, password string) (*bind.TransactOpts, error) {
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("addr: ", key.Address.Hex())
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(1337)) //nolint:gomnd
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = 99999999999
	return auth, nil
}
