package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"gotest.tools/assert"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress        = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	bridgeAddress     = "0xffffffffffffffffffffffffffffffffffffffff"
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

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()

			// Init database instance
			err = dbutils.InitOrReset(dbConfig)
			require.NoError(t, err)

			// Connect to db
			sqlDB, err := db.NewSQLDB(dbConfig)
			require.NoError(t, err)

			// Set genesis
			store := tree.NewPostgresStore(sqlDB)
			mt := tree.NewMerkleTree(store, testCase.Arity, poseidon.Hash)
			tr := tree.NewStateTree(mt, []byte{})

			stateCfg := state.Config{
				DefaultChainID: 1000,
			}

			stateDB := pgstatestorage.NewPostgresStorage(sqlDB)
			st := state.NewState(stateCfg, stateDB, tr)
			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}
			for _, gacc := range testCase.GenesisAccounts {
				b := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &b
			}
			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			// Check initial root
			root, err := st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot := new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedOldRoot, strRoot, "Invalid old root")

			// Run network container
			err = startNetworkContainer()
			require.NoError(t, err)

			// Wait network to be ready
			err = waitPoll(1*time.Second, 15*time.Second, networkUpCondition)
			require.NoError(t, err)

			// Start prover container
			err = startProverContainer()
			require.NoError(t, err)

			// Wait prover to be ready
			err = waitPoll(1*time.Second, 5*time.Second, proverUpCondition)
			require.NoError(t, err)

			// Eth client
			client, err := ethclient.Dial(l1NetworkURL)
			require.NoError(t, err)

			// Get network chain id
			chainID, err := client.NetworkID(context.Background())
			require.NoError(t, err)

			// Preparing l1 acc info
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
			require.NoError(t, err)
			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)

			// Getting l1 info
			gasPrice, err := client.SuggestGasPrice(context.Background())
			require.NoError(t, err)

			// Send some Ether from l1Acc to sequencer acc
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

			// Wait eth transfer to be mined
			err = waitTxToBeMined(client, signedTx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// Create matic maticTokenSC sc instance
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

			// Check matic balance
			require.NoError(t, err)
			b, err := maticTokenSC.BalanceOf(&bind.CallOpts{}, toAddress)
			require.NoError(t, err)
			assert.Equal(t, b.Cmp(maticAmount), 0, fmt.Sprintf("expected: %v found %v", maticAmount.Text(encoding.Base10), b.Text(encoding.Base10)))

			// Create sequencer auth
			privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(testCase.SequencerPrivateKey, "0x"))
			require.NoError(t, err)
			auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			require.NoError(t, err)

			// approve tokens to be used by PoE SC on behalf of the sequencer
			tx, err = maticTokenSC.Approve(auth, common.HexToAddress(poeAddress), maticAmount)
			require.NoError(t, err)
			err = waitTxToBeMined(client, tx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// Register the sequencer
			ethermanConfig := etherman.Config{
				URL: l1NetworkURL,
			}
			etherman, err := etherman.NewEtherman(ethermanConfig, auth, common.HexToAddress(poeAddress), common.HexToAddress(bridgeAddress), common.HexToAddress(maticTokenAddress))
			require.NoError(t, err)
			tx, err = etherman.RegisterSequencer(l2NetworkURL)
			require.NoError(t, err)

			// Wait sequencer to be registered
			err = waitTxToBeMined(client, tx.Hash(), 5*time.Second)
			require.NoError(t, err)

			// Run core container
			err = startCoreContainer()
			require.NoError(t, err)

			// Wait core to be ready
			err = waitPoll(1*time.Second, 10*time.Second, coreUpCondition)
			require.NoError(t, err)

			// Update Sequencer ChainID to the one in the test vector
			_, err = sqlDB.Exec(ctx, "UPDATE state.sequencer SET chain_id = $1 WHERE address = $2", testCase.ChainIDSequencer, common.HexToAddress(testCase.SequencerAddress).Bytes())
			require.NoError(t, err)

			// Apply transactions
			l2Client, err := ethclient.Dial(l2NetworkURL)
			require.NoError(t, err)

			// store current batch number to check later when the state is updated
			currentBatchNumber, err := st.GetLastBatchNumberSeenOnEthereum(ctx)
			require.NoError(t, err)

			for _, tx := range testCase.Txs {
				if string(tx.RawTx) != "" && tx.Overwrite.S == "" {
					l2tx := new(types.Transaction)

					b, err := hex.DecodeHex(tx.RawTx)
					require.NoError(t, err)

					err = l2tx.UnmarshalBinary(b)
					require.NoError(t, err)

					t.Logf("sending tx: %v - %v, %s", tx.ID, l2tx.Hash(), tx.From)
					err = l2Client.SendTransaction(ctx, l2tx)
					require.NoError(t, err)
				}
			}

			// Wait for sequencer to select txs from pool and propose a new batch
			// Wait for the synchronizer to update state
			err = waitPoll(1*time.Second, 15*time.Second, func() (bool, error) {
				// using a closure here to capture st and currentBatchNumber
				latestBatchNumber, err := st.GetLastBatchNumberConsolidatedOnEthereum(ctx)
				if err != nil {
					return false, err
				}
				done := latestBatchNumber > currentBatchNumber
				return done, nil
			})
			// if the state is not expected to change waitPoll can timeout
			if testCase.ExpectedNewRoot != testCase.ExpectedOldRoot {
				require.NoError(t, err)
			}

			// Check leafs
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

			// Check state against the expected state
			root, err = st.GetStateRoot(ctx, true)
			require.NoError(t, err)
			strRoot = new(big.Int).SetBytes(root).String()
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot, "Invalid new root")

			// Check consolidated state against the expected state
			consolidatedRoot, err := st.GetStateRoot(ctx, true)

			require.NoError(t, err)
			strRoot = new(big.Int).SetBytes(consolidatedRoot).String()
			assert.Equal(t, testCase.ExpectedNewRoot, strRoot)

			// Check that last virtual and consolidated batch are the same
			lastConsolidatedBatchNumber, err := st.GetLastConsolidatedBatchNumber(ctx)
			require.NoError(t, err)
			lastVirtualBatchNumber, err := st.GetLastBatchNumber(ctx)
			require.NoError(t, err)
			t.Logf("lastConsolidatedBatchNumber: %d lastVirtualBatchNumber: %d", lastConsolidatedBatchNumber, lastVirtualBatchNumber)
			assert.Equal(t, lastConsolidatedBatchNumber, lastVirtualBatchNumber)

			err = stopCoreContainer()
			require.NoError(t, err)

			err = stopProverContainer()
			require.NoError(t, err)

			err = stopNetworkContainer()
			require.NoError(t, err)
		})
	}

	err = stopCoreContainer()
	require.NoError(t, err)

	err = stopProverContainer()
	require.NoError(t, err)

	err = stopNetworkContainer()
	require.NoError(t, err)
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
	return runCmd(cmd)
}

func stopNetworkContainer() error {
	cmd := exec.Command(makeCmd, "stop-network")
	return runCmd(cmd)
}

func startCoreContainer() error {
	if err := stopCoreContainer(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-core")
	return runCmd(cmd)
}

func stopCoreContainer() error {
	cmd := exec.Command(makeCmd, "stop-core")
	return runCmd(cmd)
}

func startProverContainer() error {
	if err := stopProverContainer(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-prover")
	return runCmd(cmd)
}

func stopProverContainer() error {
	cmd := exec.Command(makeCmd, "stop-prover")
	return runCmd(cmd)
}

func runCmd(c *exec.Cmd) error {
	c.Dir = cmdDir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
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

func nodeUpCondition(target string) (bool, error) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}`)
	req, err := http.NewRequest(
		"POST", target,
		bytes.NewBuffer(jsonStr))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	if res.Body != nil {
		defer func() {
			err = res.Body.Close()
		}()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return false, err
	}

	type response struct {
		Result bool
	}

	r := response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return false, err
	}

	done := !r.Result

	return done, nil
}

type conditionFunc func() (done bool, err error)

func networkUpCondition() (bool, error) {
	return nodeUpCondition(l1NetworkURL)
}

func proverUpCondition() (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	proverClient := proverclient.NewZKProverClient(conn)
	state, err := proverClient.GetStatus(context.Background(), &proverclient.NoParams{})
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	done := state.Status == proverclient.State_IDLE

	return done, nil
}

func coreUpCondition() (done bool, err error) {
	return nodeUpCondition(l2NetworkURL)
}

func waitPoll(interval, deadline time.Duration, condition conditionFunc) error {
	timeout := time.After(deadline)
	tick := time.Tick(interval)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("Condition not met after %s", deadline)
		case <-tick:
			ok, err := condition()
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
}
