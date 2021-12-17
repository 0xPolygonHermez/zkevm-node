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
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/aggregator"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sequencerURL = "http://localhost:8123"
)

//nolint:gomnd
var cfg = config.Config{
	Log: log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	},
	Database: db.Config{
		User:     "test_user",
		Password: "test_password",
		Name:     "test_db",
		Host:     "localhost",
		Port:     "5432",
	},
	Etherman: etherman.Config{
		URL:                "http://localhost:8545",
		PrivateKeyPath:     "../test.keystore",
		PrivateKeyPassword: "testonly",
	},
	RPC: jsonrpc.Config{
		Host: "",
		Port: 8123,
	},
	Sequencer: sequencer.Config{
		IntervalToProposeBatch: sequencer.Duration{Duration: 1 * time.Second},
	},
	Aggregator: aggregator.Config{},
	NetworkConfig: config.NetworkConfig{
		Arity:            4,
		GenBlockNumber:   1,
		PoEAddr:          common.HexToAddress("0x41D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"),
		L1ChainID:        1337,
		L2DefaultChainID: 50000,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
			common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
		},
	},
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

			// set genesis
			mt := tree.NewMerkleTree(sqlDB, testCase.Arity, poseidon.Hash)
			tr := tree.NewStateTree(mt, []byte{})
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

			// Run network container
			err = startNetworkContainer()
			require.NoError(t, err)

			// wait network to be ready
			time.Sleep(5 * time.Second)

			// create auth
			pkHex := strings.TrimPrefix(testCase.SequencerPrivateKey, "0x")
			privateKey, err := crypto.HexToECDSA(pkHex)
			if err != nil {
				log.Fatal(err)
			}
			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
			if err != nil {
				log.Fatal(err)
			}
			auth.GasLimit = 100000000

			// register the sequencer
			etherman, err := etherman.NewEtherman(cfg.Etherman, auth, cfg.NetworkConfig.PoEAddr)
			require.NoError(t, err)
			_, err = etherman.PoE.RegisterSequencer(auth, sequencerURL)
			require.NoError(t, err)

			// Run node container
			err = startNodeContainer()
			require.NoError(t, err)

			// wait node to be ready
			time.Sleep(3 * time.Second)

			// apply transactions
			for _, tx := range testCase.Txs {
				rawTx := tx.RawTx
				err := sendRawTransaction(rawTx)
				require.NoError(t, err)
			}

			// wait for sequencer to select txs from pool and propose a new batch
			// wait for the synchronizer to update state
			time.Sleep(10 * time.Second)

			// stop node
			err = stopNodeContainer()
			require.NoError(t, err)

			// stop network
			err = stopNetworkContainer()
			require.NoError(t, err)

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
				assert.Equal(t, 0, leaf.Balance.Cmp(actualBalance), fmt.Sprintf("addr: %s expected: %s found: %s", addr.Hex(), leaf.Balance.Text(encoding.Base10), actualBalance.Text(encoding.Base10)))

				actualNonce, err := st.GetNonce(addr, batchNumber)
				require.NoError(t, err)
				assert.Equal(t, leaf.Nonce, strconv.FormatUint(actualNonce, encoding.Base10), fmt.Sprintf("addr: %s expected: %s found: %d", addr.Hex(), leaf.Nonce, actualNonce))
			}
		})
	}
}

const (
	makeCmd = "make"
	cmdDir  = "../.."
)

func startNetworkContainer() error {
	if err := stopNetworkContainer(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-network")
	cmd.Dir = cmdDir
	return cmd.Run()
}

func stopNetworkContainer() error {
	cmd := exec.Command(makeCmd, "stop-network")
	cmd.Dir = cmdDir
	return cmd.Run()
}

func startNodeContainer() error {
	if err := stopNodeContainer(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-core")
	cmd.Dir = cmdDir
	return cmd.Run()
}

func stopNodeContainer() error {
	cmd := exec.Command(makeCmd, "stop-core")
	cmd.Dir = cmdDir
	return cmd.Run()
}

func sendRawTransaction(rawTx string) error {
	endpoint := fmt.Sprintf("http://localhost:%d", cfg.RPC.Port)
	contentType := "application/json"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_sendRawTransaction",
		"params":  []string{rawTx},
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
