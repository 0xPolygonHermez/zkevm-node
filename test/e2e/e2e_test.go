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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

var (
	dbConfig = db.Config{
		User:     "test_user",
		Password: "test_password",
		Name:     "test_db",
		Host:     "localhost",
		Port:     "5432",
	}
)

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/state-transition.json")
	require.NoError(t, err)

	buildCore()

	defer stopNodeContainer()
	defer stopNetworkContainer()

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			// init database instance
			err = dbutils.InitOrReset(dbConfig)
			require.NoError(t, err)

			//connect to db
			sqlDB, err := db.NewSQLDB(dbConfig)
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

			// eth client
			client, err := ethclient.Dial(l1NetworkURL)
			require.NoError(t, err)

			// get network chain id
			chainID, err := client.NetworkID(context.Background())
			require.NoError(t, err)

			// send some Ether from l1Acc to sequencer acc
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
			require.NoError(t, err)
			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)
			auth.GasLimit = 999999999
			fromAddress := common.HexToAddress(l1AccHexAddress)
			nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
			require.NoError(t, err)
			gasLimit := uint64(21000)
			gasPrice, err := client.SuggestGasPrice(context.Background())
			require.NoError(t, err)
			toAddress := common.HexToAddress(testCase.SequencerAddress)
			tx := types.NewTransaction(nonce, toAddress, big.NewInt(1000000000000000000), gasLimit, gasPrice, nil)
			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)
			err = client.SendTransaction(context.Background(), signedTx)
			require.NoError(t, err)

			// wait transfer to be mined
			time.Sleep(3 * time.Second)

			// create sequencer auth
			privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(testCase.SequencerPrivateKey, "0x"))
			require.NoError(t, err)

			auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)
			auth.GasLimit = 999999999
			auth.GasPrice = gasPrice

			// register the sequencer
			// ethermanConfig := etherman.Config{
			// 	URL: l1NetworkURL,
			// }
			// etherman, err := etherman.NewEtherman(ethermanConfig, auth, common.HexToAddress(poeAddress))
			// require.NoError(t, err)
			// _, err = etherman.RegisterSequencer(l2NetworkURL)
			// require.NoError(t, err)

			// wait sequencer to be registered
			time.Sleep(3 * time.Second)

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

func buildCore() error {
	cmd := exec.Command(makeCmd, "build-docker")
	cmd.Dir = cmdDir
	return cmd.Run()
}

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
	cmd.Env = []string{"HERMEZCORE_NETWORK=e2e-test"}
	cmd.Dir = cmdDir
	return cmd.Run()
}

func stopNodeContainer() error {
	cmd := exec.Command(makeCmd, "stop-core")
	cmd.Dir = cmdDir
	return cmd.Run()
}

func sendRawTransaction(rawTx string) error {
	contentType := "application/json"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_sendRawTransaction",
		"params":  []string{rawTx},
	}

	jsonStr, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", l2NetworkURL, bytes.NewBuffer(jsonStr))
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
