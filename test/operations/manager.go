package operations

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress        = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	maticTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	aggregatorAddress = "0xb94f5374fce5edbc8e2a8697c15331677e6ebf0b"
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

	st   *state.State
	wait *Wait
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
		cfg:  cfg,
		ctx:  ctx,
		wait: NewWait(),
	}
	st, err := initState(cfg.Arity, cfg.State.DefaultChainID, cfg.State.MaxCumulativeGasUsed)
	if err != nil {
		return nil, err
	}
	opsman.st = st

	return opsman, nil
}

// State is a getter for the st field.
func (m *Manager) State() *state.State {
	return m.st
}

// CheckVirtualRoot verifies if the given root is the current root of the
// merkletree for virtual state.
func (m *Manager) CheckVirtualRoot(expectedRoot string) error {
	root, err := m.st.GetStateRoot(m.ctx, true, "")
	if err != nil {
		return err
	}
	return m.checkRoot(root, expectedRoot)
}

// CheckConsolidatedRoot verifies if the given root is the current root of the
// merkletree for consolidated state.
func (m *Manager) CheckConsolidatedRoot(expectedRoot string) error {
	root, err := m.st.GetStateRoot(m.ctx, false, "")
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

	return m.st.SetGenesis(m.ctx, genesis, "")
}

// ApplyTxs sends the given L2 txs, waits for them to be consolidated and checks
// the final state.
func (m *Manager) ApplyTxs(vectorTxs []vectors.Tx, initialRoot, finalRoot, globalExitRoot string) error {
	// store current batch number to check later when the state is updated
	currentBatchNumber, err := m.st.GetLastBatchNumberSeenOnEthereum(m.ctx, "")
	if err != nil {
		return err
	}

	var txs []*types.Transaction
	for _, vectorTx := range vectorTxs {
		if string(vectorTx.RawTx) != "" && vectorTx.Overwrite.S == "" {
			var tx types.LegacyTx
			bytes, err := hex.DecodeHex(vectorTx.RawTx)
			if err != nil {
				return err
			}

			err = rlp.DecodeBytes(bytes, &tx)
			if err == nil {
				txs = append(txs, types.NewTx(&tx))
			}
			if err != nil {
				return err
			}
		}
	}
	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          common.HexToAddress(m.cfg.Sequencer.Address),
		Aggregator:         common.HexToAddress(aggregatorAddress),
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ChainID:            big.NewInt(int64(m.cfg.State.DefaultChainID)),
		GlobalExitRoot:     common.HexToHash(globalExitRoot),
	}

	// Create Batch Processor
	bp, err := m.st.NewBatchProcessor(m.ctx, common.HexToAddress(m.cfg.Sequencer.Address), common.Hex2Bytes(strings.TrimPrefix(initialRoot, "0x")), "")
	if err != nil {
		return err
	}

	err = bp.ProcessBatch(m.ctx, batch)
	if err != nil {
		return err
	}

	// Wait for sequencer to select txs from pool and propose a new batch
	// Wait for the synchronizer to update state
	err = Poll(DefaultInterval, DefaultDeadline, func() (bool, error) {
		// using a closure here to capture st and currentBatchNumber
		latestBatchNumber, err := m.st.GetLastBatchNumberConsolidatedOnEthereum(m.ctx, "")
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

// Setup creates all the required components and initializes them according to
// the manager config.
func (m *Manager) Setup() error {
	// Run network container
	err := m.StartNetwork()
	if err != nil {
		return err
	}

	// Start prover container
	err = m.StartProver()
	if err != nil {
		return err
	}

	err = m.SetUpSequencer()
	if err != nil {
		return err
	}

	// Run core container
	err = m.StartCore()
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

func initState(arity uint8, defaultChainID uint64, maxCumulativeGasUsed uint64) (*state.State, error) {
	sqlDB, err := db.NewSQLDB(dbConfig)
	if err != nil {
		return nil, err
	}

	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, arity)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore)

	stateCfg := state.Config{
		DefaultChainID:              defaultChainID,
		MaxCumulativeGasUsed:        maxCumulativeGasUsed,
		L2GlobalExitRootManagerAddr: common.HexToAddress("0xAE4bB80bE56B819606589DE61d5ec3b522EEB032"),
	}

	stateDB := state.NewPostgresStorage(sqlDB)
	return state.NewState(stateCfg, stateDB, tr), nil
}

func (m *Manager) checkRoot(root []byte, expectedRoot string) error {
	actualRoot := hex.EncodeToHex(root)

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

// SetUpSequencer provide ETH, Matic to and register the sequencer
func (m *Manager) SetUpSequencer() error {
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
	err = WaitTxToBeMined(client, signedTx.Hash(), DefaultTxMinedDeadline)
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
	err = WaitTxToBeMined(client, tx.Hash(), DefaultTxMinedDeadline)
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

	err = WaitTxToBeMined(client, tx.Hash(), DefaultTxMinedDeadline)
	if err != nil {
		return err
	}

	// Register the sequencer
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewClient(ethermanConfig, auth, common.HexToAddress(poeAddress), common.HexToAddress(maticTokenAddress))
	if err != nil {
		return err
	}
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	if err != nil {
		return err
	}

	// Wait sequencer to be registered
	err = WaitTxToBeMined(client, tx.Hash(), DefaultTxMinedDeadline)
	if err != nil {
		return err
	}
	return nil
}

// StartNetwork starts the L1 network container
func (m *Manager) StartNetwork() error {
	return StartComponent("run-network", "stop-network", networkUpCondition)
}

// InitNetwork Initializes the L2 network registering the sequencer and adding funds via the bridge
func (m *Manager) InitNetwork() error {
	if err := RunMakeTarget("init-network"); err != nil {
		return err
	}

	// Wait network to be ready
	return Poll(DefaultInterval, DefaultDeadline, networkUpCondition)
}

// DeployUniswap deploys a uniswap environment and perform swaps
func (m *Manager) DeployUniswap() error {
	if err := RunMakeTarget("deploy-uniswap"); err != nil {
		return err
	}
	// Wait network to be ready
	return Poll(DefaultInterval, DefaultDeadline, networkUpCondition)
}

func stopNetwork() error {
	return RunMakeTarget("stop-network")
}

// StartCore starts the core container
func (m *Manager) StartCore() error {
	return StartComponent("run-core", "stop-core", coreUpCondition)
}

func stopCore() error {
	return RunMakeTarget("stop-core")
}

// StartProver starts the prover container
func (m *Manager) StartProver() error {
	return StartComponent("run-prover", "stop-prover", ProverUpCondition)
}

func stopProver() error {
	return RunMakeTarget("stop-prover")
}

func runCmd(c *exec.Cmd) error {
	c.Dir = ".."
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// StartComponent starts a docker-compose component.
func StartComponent(cmdUp, cmdDown string, conditions ...ConditionFunc) error {
	if err := RunMakeTarget(cmdDown); err != nil {
		return err
	}
	if err := RunMakeTarget(cmdUp); err != nil {
		return err
	}

	// Wait component to be ready
	for _, condition := range conditions {
		if err := Poll(DefaultInterval, DefaultDeadline, condition); err != nil {
			return err
		}
	}
	return nil
}

// RunMakeTarget runs a Makefile target.
func RunMakeTarget(target string) error {
	cmd := exec.Command("make", target)
	return runCmd(cmd)
}
