package operations

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
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress            = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	bridgeAddress         = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	maticTokenAddress     = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec
	globalExitRootAddress = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	defaultInterval        = 2 * time.Second
	defaultDeadline        = 25 * time.Second
	defaultTxMinedDeadline = 5 * time.Second

	makeCmd = "make"
	cmdDir  = "../.."
)

var dbConfig = dbutils.NewConfigFromEnv()

// SequencerConfig is the configuration for the sequencer operations.
type SequencerConfig struct {
	Address, PrivateKey string
	ChainID             uint64
}

// Config is the main Manager configuration.
type Config struct {
	Arity     uint8
	State     *state.Config
	Sequencer *SequencerConfig
}

// Manager controls operations and has knowledge about how to set up and tear
// down a functional environment.
type Manager struct {
	cfg *Config
	ctx context.Context

	st state.State
}

// NewManager returns a manager ready to be used and a potential error caused
// during its creation (which can come from the setup of the db connection).
func NewManager(ctx context.Context, cfg *Config) (*Manager, error) {
	// Init database instance
	err := dbutils.InitOrReset(dbConfig)
	if err != nil {
		return nil, err
	}

	opsman := &Manager{
		cfg: cfg,
		ctx: ctx,
	}
	st, err := initState(cfg.Arity, cfg.State.DefaultChainID, cfg.State.MaxCumulativeGasUsed)
	if err != nil {
		return nil, err
	}
	opsman.st = st

	return opsman, nil
}

// State is a getter for the st field.
func (m *Manager) State() state.State {
	return m.st
}

// CheckVirtualRoot verifies if the given root is the current root of the
// merkletree for virtual state.
func (m *Manager) CheckVirtualRoot(expectedRoot string) error {
	root, err := m.st.GetStateRoot(m.ctx, true)
	if err != nil {
		return err
	}
	return m.checkRoot(root, expectedRoot)
}

// CheckConsolidatedRoot verifies if the given root is the current root of the
// merkletree for consolidated state.
func (m *Manager) CheckConsolidatedRoot(expectedRoot string) error {
	root, err := m.st.GetStateRoot(m.ctx, false)
	if err != nil {
		return err
	}
	return m.checkRoot(root, expectedRoot)
}

// SetGenesis creates the genesis block in the state.
func (m *Manager) SetGenesis(genesisAccounts map[string]big.Int) error {
	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}
	for address, balanceValue := range genesisAccounts {
		// prevent taking the address of a loop variable
		balance := balanceValue
		genesis.Balances[common.HexToAddress(address)] = &balance
	}

	return m.st.SetGenesis(m.ctx, genesis)
}

// ApplyTxs sends the given L2 txs, waits for them to be consolidated and checks
// the final state.
func (m *Manager) ApplyTxs(txs []vectors.Tx, initialRoot, finalRoot string) error {
	// Apply transactions
	l2Client, err := ethclient.Dial(l2NetworkURL)
	if err != nil {
		return err
	}

	// store current batch number to check later when the state is updated
	currentBatchNumber, err := m.st.GetLastBatchNumberSeenOnEthereum(m.ctx)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		if string(tx.RawTx) != "" && tx.Overwrite.S == "" {
			l2tx := new(types.Transaction)

			b, err := hex.DecodeHex(tx.RawTx)
			if err != nil {
				return err
			}

			err = l2tx.UnmarshalBinary(b)
			if err != nil {
				return err
			}

			log.Infof("sending tx: %v - %v, %s", tx.ID, l2tx.Hash(), tx.From)
			err = l2Client.SendTransaction(m.ctx, l2tx)
			if err != nil {
				return err
			}
		}
	}

	// Wait for sequencer to select txs from pool and propose a new batch
	// Wait for the synchronizer to update state
	err = waitPoll(defaultInterval, defaultDeadline, func() (bool, error) {
		// using a closure here to capture st and currentBatchNumber
		latestBatchNumber, err := m.st.GetLastBatchNumberConsolidatedOnEthereum(m.ctx)
		if err != nil {
			return false, err
		}
		done := latestBatchNumber > currentBatchNumber
		return done, nil
	})
	// if the state is not expected to change waitPoll can timeout
	if initialRoot != "" && finalRoot != "" && initialRoot != finalRoot && err != nil {
		return err
	}
	return nil
}

// GetAuth configures and returns an auth object.
func GetAuth(privateKeyStr string, chainID *big.Int) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactorWithChainID(privateKey, chainID)
}

// WaitGRPCHealthy waits for a gRPC endpoint to be responding according to the
// health standard in package grpc.health.v1
func WaitGRPCHealthy(address string) error {
	return waitPoll(defaultInterval, defaultDeadline, func() (bool, error) {
		return grpcHealthyCondition(address)
	})
}

// Setup creates all the required components and initializes them according to
// the manager config.
func (m *Manager) Setup() error {
	// Run network container
	err := startNetwork()
	if err != nil {
		return err
	}

	// Start prover container
	err = startProver()
	if err != nil {
		return err
	}

	err = m.setUpSequencer()
	if err != nil {
		return err
	}

	// Run core container
	err = startCore()
	if err != nil {
		return err
	}

	return m.setSequencerChainID()
}

// Teardown stops all the components.
func Teardown() error {
	err := stopCore()
	if err != nil {
		return err
	}

	err = stopProver()
	if err != nil {
		return err
	}

	err = stopNetwork()
	if err != nil {
		return err
	}

	return nil
}

func initState(arity uint8, defaultChainID uint64, maxCumulativeGasUsed uint64) (state.State, error) {
	sqlDB, err := db.NewSQLDB(dbConfig)
	if err != nil {
		return nil, err
	}

	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, arity, poseidon.Hash)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore, []byte{})

	stateCfg := state.Config{
		DefaultChainID:       defaultChainID,
		MaxCumulativeGasUsed: maxCumulativeGasUsed,
	}

	stateDB := pgstatestorage.NewPostgresStorage(sqlDB)
	return state.NewState(stateCfg, stateDB, tr), nil
}

func (m *Manager) checkRoot(root []byte, expectedRoot string) error {
	actualRoot := new(big.Int).SetBytes(root).String()

	if expectedRoot != actualRoot {
		return fmt.Errorf("Invalid root, want %q, got %q", expectedRoot, actualRoot)
	}
	return nil
}

func (m *Manager) setSequencerChainID() error {
	// Update Sequencer ChainID to the one in the test vector
	sqlDB, err := db.NewSQLDB(dbConfig)
	if err != nil {
		return err
	}

	_, err = sqlDB.Exec(m.ctx, "UPDATE state.sequencer SET chain_id = $1 WHERE address = $2", m.cfg.Sequencer.ChainID, common.HexToAddress(m.cfg.Sequencer.Address).Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) setUpSequencer() error {
	// Eth client
	client, err := ethclient.Dial(l1NetworkURL)
	if err != nil {
		return err
	}

	// Get network chain id
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	auth, err := GetAuth(l1AccHexPrivateKey, chainID)
	if err != nil {
		return err
	}

	// Getting l1 info
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Send some Ether from l1Acc to sequencer acc
	fromAddress := common.HexToAddress(l1AccHexAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	const (
		gasLimit = 21000
		OneEther = 1000000000000000000
	)
	toAddress := common.HexToAddress(m.cfg.Sequencer.Address)
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(OneEther), uint64(gasLimit), gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	// Wait eth transfer to be mined
	err = waitTxToBeMined(client, signedTx.Hash(), defaultTxMinedDeadline)
	if err != nil {
		return err
	}

	// Create matic maticTokenSC sc instance
	maticTokenSC, err := NewToken(common.HexToAddress(maticTokenAddress), client)
	if err != nil {
		return err
	}

	// Send matic to sequencer
	maticAmount, ok := big.NewInt(0).SetString("100000000000000000000000", encoding.Base10)
	if !ok {
		return fmt.Errorf("Error setting matic amount")
	}

	tx, err = maticTokenSC.Transfer(auth, toAddress, maticAmount)
	if err != nil {
		return err
	}

	// wait matic transfer to be mined
	err = waitTxToBeMined(client, tx.Hash(), defaultTxMinedDeadline)
	if err != nil {
		return err
	}

	// Check matic balance
	b, err := maticTokenSC.BalanceOf(&bind.CallOpts{}, toAddress)
	if err != nil {
		return err
	}

	if 0 != b.Cmp(maticAmount) {
		return fmt.Errorf("expected: %v found %v", maticAmount.Text(encoding.Base10), b.Text(encoding.Base10))
	}

	// Create sequencer auth
	auth, err = GetAuth(m.cfg.Sequencer.PrivateKey, chainID)
	if err != nil {
		return err
	}

	// approve tokens to be used by PoE SC on behalf of the sequencer
	tx, err = maticTokenSC.Approve(auth, common.HexToAddress(poeAddress), maticAmount)
	if err != nil {
		return err
	}

	err = waitTxToBeMined(client, tx.Hash(), defaultTxMinedDeadline)
	if err != nil {
		return err
	}

	// Register the sequencer
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewEtherman(ethermanConfig, auth, common.HexToAddress(poeAddress), common.HexToAddress(bridgeAddress), common.HexToAddress(maticTokenAddress), common.HexToAddress(globalExitRootAddress))
	if err != nil {
		return err
	}
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	if err != nil {
		return err
	}

	// Wait sequencer to be registered
	err = waitTxToBeMined(client, tx.Hash(), defaultTxMinedDeadline)
	if err != nil {
		return err
	}
	return nil
}

func startNetwork() error {
	if err := stopNetwork(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-network")
	err := runCmd(cmd)
	if err != nil {
		return err
	}
	// Wait network to be ready
	return waitPoll(defaultInterval, defaultDeadline, networkUpCondition)
}

func stopNetwork() error {
	cmd := exec.Command(makeCmd, "stop-network")
	return runCmd(cmd)
}

func startCore() error {
	if err := stopCore(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-core")
	err := runCmd(cmd)
	if err != nil {
		return err
	}
	// Wait core to be ready
	return waitPoll(defaultInterval, defaultDeadline, coreUpCondition)
}

func stopCore() error {
	cmd := exec.Command(makeCmd, "stop-core")
	return runCmd(cmd)
}

func startProver() error {
	if err := stopProver(); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, "run-prover")
	err := runCmd(cmd)
	if err != nil {
		return err
	}
	// Wait prover to be ready
	return waitPoll(defaultInterval, defaultDeadline, proverUpCondition)
}

func stopProver() error {
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

	r := struct {
		Result bool
	}{}
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:50051", opts...)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}
	defer func() {
		err = conn.Close()
	}()

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

func grpcHealthyCondition(address string) (bool, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}
	defer func() {
		err = conn.Close()
	}()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	state, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		// we allow connection errors to wait for the container up
		return false, nil
	}

	done := state.Status == grpc_health_v1.HealthCheckResponse_SERVING

	return done, nil
}

func waitPoll(interval, deadline time.Duration, condition conditionFunc) error {
	timeout := time.After(deadline)
	tick := time.NewTicker(interval)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("Condition not met after %s", deadline)
		case <-tick.C:
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
