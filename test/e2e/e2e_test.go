package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress        = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	maticTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

var dbConfig = dbutils.NewConfigFromEnv()

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testCases, err := vectors.LoadStateTransitionTestCases("./../vectors/state-transition.json")
	require.NoError(t, err)

	err = stopCoreContainer()
	require.NoError(t, err)

	err = stopNetworkContainer()
	require.NoError(t, err)

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
			store := tree.NewPostgresStore(sqlDB)
			mt := tree.NewMerkleTree(store, testCase.Arity, poseidon.Hash)
			tr := tree.NewStateTree(mt, []byte{})

			stateCfg := state.Config{
				DefaultChainID: 1000,
			}

			st := state.NewState(stateCfg, sqlDB, tr)
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

			// preparing l1 acc info
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
			require.NoError(t, err)
			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)

			// getting l1 info
			gasPrice, err := client.SuggestGasPrice(context.Background())
			require.NoError(t, err)

			// send some Ether from l1Acc to sequencer acc
			fromAddress := common.HexToAddress(l1AccHexAddress)
			nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
			require.NoError(t, err)
			gasLimit := uint64(21000)
			toAddress := common.HexToAddress(testCase.SequencerAddress)
			tx := types.NewTransaction(nonce, toAddress, big.NewInt(1000000000000000000), gasLimit, gasPrice, nil)
			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)
			err = client.SendTransaction(context.Background(), signedTx)
			require.NoError(t, err)

			// wait eth transfer to be mined
			err = waitTxToBeMined(client, signedTx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// create matic maticTokenSC sc instance
			maticTokenSC, err := NewToken(common.HexToAddress(maticTokenAddress), client)
			require.NoError(t, err)

			// Send matic to sequencer
			maticAmount, ok := big.NewInt(0).SetString("100000000000000000000000", encoding.Base10)
			require.True(t, ok)
			tx, err = maticTokenSC.Transfer(auth, toAddress, maticAmount)
			require.NoError(t, err)

			// wait matic transfer to be mined
			err = waitTxToBeMined(client, tx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// check matic balance
			require.NoError(t, err)
			b, err := maticTokenSC.BalanceOf(&bind.CallOpts{}, toAddress)
			require.NoError(t, err)
			assert.Equal(t, b.Cmp(maticAmount), 0, fmt.Sprintf("expected: %v found %v", maticAmount.Text(encoding.Base10), b.Text(encoding.Base10)))

			// create sequencer auth
			privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(testCase.SequencerPrivateKey, "0x"))
			require.NoError(t, err)
			auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)

			// approve tokens to be used by PoE SC on behalf of the sequencer
			tx, err = maticTokenSC.Approve(auth, common.HexToAddress(poeAddress), maticAmount)
			require.NoError(t, err)
			err = waitTxToBeMined(client, tx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// register the sequencer
			ethermanConfig := etherman.Config{
				URL: l1NetworkURL,
			}
			etherman, err := etherman.NewEtherman(ethermanConfig, auth, common.HexToAddress(poeAddress))
			require.NoError(t, err)
			tx, err = etherman.RegisterSequencer(l2NetworkURL)
			require.NoError(t, err)

			// wait sequencer to be registered
			err = waitTxToBeMined(client, tx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// Run core container
			err = startCoreContainer(testCase.SequencerPrivateKey)
			require.NoError(t, err)

			// wait core to be ready
			time.Sleep(5 * time.Second)

			// update Sequencer ChainID to the one in the test vector
			_, err = sqlDB.Exec(ctx, "UPDATE state.sequencer SET chain_id = $1 WHERE address = $2", testCase.ChainIDSequencer, common.HexToAddress(testCase.SequencerAddress).Bytes())
			require.NoError(t, err)

			// apply transactions
			for _, tx := range testCase.Txs {
				if string(tx.RawTx) != "" && tx.Overwrite.S == "" {
					rawTx := tx.RawTx
					err := sendRawTransaction(rawTx)
					require.NoError(t, err)
				}
			}

			// wait for sequencer to select txs from pool and propose a new batch
			// wait for the synchronizer to update state
			time.Sleep(10 * time.Second)

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

			// check state against the expected state
			root, err = st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot = new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot, "Invalid new root")

			err = stopCoreContainer()
			require.NoError(t, err)

			err = stopNetworkContainer()
			require.NoError(t, err)
		})
	}
}

const (
	makeCmd = "make"
	cmdDir  = "../.."
)

func startNetworkContainer() error {
	cmd := exec.Command(makeCmd, "run-network")
	return runCmd(cmd)
}

func stopNetworkContainer() error {
	cmd := exec.Command(makeCmd, "stop-network")
	return runCmd(cmd)
}

func startCoreContainer(sequencerPrivateKey string) error {
	cmd := exec.Command(makeCmd, "run-core")
	cmd.Env = []string{
		"HERMEZCORE_NETWORK=e2e-test",
		"HERMEZCORE_KEYSTORE_FILEPATH=./test/e2e/e2e.keystore",
	}
	return runCmd(cmd)
}

func stopCoreContainer() error {
	cmd := exec.Command(makeCmd, "stop-core")
	return runCmd(cmd)
}

func runCmd(c *exec.Cmd) error {
	c.Dir = cmdDir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
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

func waitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return errors.New("timeout exceed")
		}

		time.Sleep(1 * time.Second)

		_, isPending, err := client.TransactionByHash(context.Background(), hash)
		if err == ethereum.NotFound {
			continue
		}

		if err != nil {
			return err
		}

		if !isPending {
			r, err := client.TransactionReceipt(context.Background(), hash)
			if err != nil {
				return err
			}

			if r.Status == types.ReceiptStatusFailed {
				return fmt.Errorf("transaction has failed: %s", string(r.PostState))
			}

			return nil
		}
	}
}
