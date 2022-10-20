package operations

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	poeAddress         = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	maticTokenAddress  = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec
	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

// Public constants
const (
	DefaultSequencerAddress     = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	DefaultSequencerPrivateKey  = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DefaultSequencerBalance     = 400000
	DefaultMaxCumulativeGasUsed = 800000

	DefaultL1NetworkURL        = "http://localhost:8545"
	DefaultL1ChainID    uint64 = 1337
	DefaultL2NetworkURL        = "http://localhost:8123"
	DefaultL2ChainID    uint64 = 1001

	DefaultTimeoutTxToBeMined = 1 * time.Minute
)

var (
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	poolDBCfg  = dbutils.NewPoolConfigFromEnv()
	rpcDBCfg   = dbutils.NewRPCConfigFromEnv()

	executorURI      = testutils.GetEnv(constants.ENV_ZKPROVER_URI, "127.0.0.1:50071")
	merkleTreeURI    = testutils.GetEnv(constants.ENV_MERKLETREE_URI, "127.0.0.1:50061")
	executorConfig   = executor.Config{URI: executorURI}
	merkleTreeConfig = merkletree.Config{URI: merkleTreeURI}
)

// SequencerConfig is the configuration for the sequencer operations.
type SequencerConfig struct {
	Address, PrivateKey string
}

// Config is the main Manager configuration.
type Config struct {
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
	initOrResetDB()

	opsman := &Manager{
		cfg:  cfg,
		ctx:  ctx,
		wait: NewWait(),
	}
	st, err := initState(cfg.State.MaxCumulativeGasUsed)
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
	panic("not implemented yet")
	// root, err := m.st.Getroot(m.ctx, true, "")
	// if err != nil {
	// 	return err
	// }
	// return m.checkRoot(root, expectedRoot)
}

// CheckConsolidatedRoot verifies if the given root is the current root of the
// merkletree for consolidated state.
func (m *Manager) CheckConsolidatedRoot(expectedRoot string) error {
	panic("not implemented yet")
	// root, err := m.st.GetStateRoot(m.ctx, false, "")
	// if err != nil {
	// 	return err
	// }
	// return m.checkRoot(root, expectedRoot)
}

// SetGenesis creates the genesis block in the state.
func (m *Manager) SetGenesis(genesisAccounts map[string]big.Int) error {
	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{},
	}
	for address, balanceValue := range genesisAccounts {
		action := &state.GenesisAction{
			Address: address,
			Type:    int(merkletree.LeafTypeBalance),
			Value:   balanceValue.String(),
		}
		genesis.Actions = append(genesis.Actions, action)
	}

	dbTx, err := m.st.BeginStateTransaction(m.ctx)
	if err != nil {
		return err
	}

	_, err = m.st.SetGenesis(m.ctx, genesisBlock, genesis, dbTx)

	return err
}

// ApplyTxs sends the given L2 txs, waits for them to be consolidated and checks
// the final state.
func (m *Manager) ApplyTxs(vectorTxs []vectors.Tx, initialRoot, finalRoot, globalExitRoot string) error {
	panic("not implemented yet")
	// // store current batch number to check later when the state is updated
	// currentBatchNumber, err := m.st.GetLastBatchNumberSeenOnEthereum(m.ctx, nil)
	// if err != nil {
	// 	return err
	// }

	// var txs []*types.Transaction
	// for _, vectorTx := range vectorTxs {
	// 	if string(vectorTx.RawTx) != "" && vectorTx.Overwrite.S == "" {
	// 		var tx types.LegacyTx
	// 		bytes, err := hex.DecodeHex(vectorTx.RawTx)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		err = rlp.DecodeBytes(bytes, &tx)
	// 		if err == nil {
	// 			txs = append(txs, types.NewTx(&tx))
	// 		}
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	// // Create Batch
	// batch := &state.Batch{
	// 	BlockNumber:        uint64(0),
	// 	Sequencer:          common.HexToAddress(m.cfg.Sequencer.Address),
	// 	Aggregator:         common.HexToAddress(aggregatorAddress),
	// 	ConsolidatedTxHash: common.Hash{},
	// 	Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
	// 	Uncles:             nil,
	// 	Transactions:       txs,
	// 	RawTxsData:         nil,
	// 	MaticCollateral:    big.NewInt(1),
	// 	ChainID:            big.NewInt(int64(m.cfg.NetworkConfig.ChainID)),
	// 	GlobalExitRoot:     common.HexToHash(globalExitRoot),
	// }

	// // Create Batch Processor
	// bp, err := m.st.NewBatchProcessor(m.ctx, common.HexToAddress(m.cfg.Sequencer.Address), common.Hex2Bytes(strings.TrimPrefix(initialRoot, "0x")), "")
	// if err != nil {
	// 	return err
	// }

	// err = bp.ProcessBatch(m.ctx, batch)
	// if err != nil {
	// 	return err
	// }

	// // Wait for sequencer to select txs from pool and propose a new batch
	// // Wait for the synchronizer to update state
	// err = Poll(DefaultInterval, DefaultDeadline, func() (bool, error) {
	// 	// using a closure here to capture st and currentBatchNumber
	// 	latestBatchNumber, err := m.st.GetLastBatchNumberConsolidatedOnEthereum(m.ctx, "")
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	done := latestBatchNumber > currentBatchNumber
	// 	return done, nil
	// })
	// // if the state is not expected to change waitPoll can timeout
	// if initialRoot != "" && finalRoot != "" && initialRoot != finalRoot && err != nil {
	// 	return err
	// }
	// return nil
}

// GetAuth configures and returns an auth object.
func GetAuth(privateKeyStr string, chainID uint64) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(chainID))
}

// MustGetAuth GetAuth but panics if err
func MustGetAuth(privateKeyStr string, chainID uint64) *bind.TransactOpts {
	auth, err := GetAuth(privateKeyStr, chainID)
	if err != nil {
		panic(err)
	}
	return auth
}

// Setup creates all the required components and initializes them according to
// the manager config.
func (m *Manager) Setup() error {
	// Run network container
	err := m.StartNetwork()
	if err != nil {
		return err
	}

	// Approve matic
	err = approveMatic()
	if err != nil {
		return err
	}

	err = m.SetUpSequencer()
	if err != nil {
		return err
	}

	// Run node container
	return m.StartNode()
}

// Teardown stops all the components.
func Teardown() error {
	err := stopNode()
	if err != nil {
		return err
	}

	err = stopNetwork()
	if err != nil {
		return err
	}

	return nil
}

func initState(maxCumulativeGasUsed uint64) (*state.State, error) {
	sqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	stateDb := state.NewPostgresStorage(sqlDB)
	executorClient, _, _ := executor.NewExecutorClient(ctx, executorConfig)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, merkleTreeConfig)
	stateTree := merkletree.NewStateTree(stateDBClient)

	stateCfg := state.Config{
		MaxCumulativeGasUsed: maxCumulativeGasUsed,
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree)
	return st, nil
}

// func (m *Manager) checkRoot(root []byte, expectedRoot string) error {
// 	actualRoot := hex.EncodeToHex(root)

// 	if expectedRoot != actualRoot {
// 		return fmt.Errorf("Invalid root, want %q, got %q", expectedRoot, actualRoot)
// 	}
// 	return nil
// }

// SetUpSequencer provide ETH, Matic to and register the sequencer
func (m *Manager) SetUpSequencer() error {
	// Eth client
	client, err := ethclient.Dial(DefaultL1NetworkURL)
	if err != nil {
		return err
	}

	// Get network chain id
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	auth, err := GetAuth(l1AccHexPrivateKey, chainID.Uint64())
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

	if b.Cmp(maticAmount) < 0 {
		return fmt.Errorf("Minimum amount is: %v but found: %v", maticAmount.Text(encoding.Base10), b.Text(encoding.Base10))
	}

	// Create sequencer auth
	auth, err = GetAuth(m.cfg.Sequencer.PrivateKey, chainID.Uint64())
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
	return nil
}

// StartNetwork starts the L1 network container
func (m *Manager) StartNetwork() error {
	return StartComponent("network", networkUpCondition)
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
	return StopComponent("network")
}

// StartNode starts the node container
func (m *Manager) StartNode() error {
	return StartComponent("node", nodeUpCondition)
}

func approveMatic() error {
	return StartComponent("approve-matic")
}

func stopNode() error {
	return StopComponent("node")
}

func runCmd(c *exec.Cmd) error {
	c.Dir = "../.."
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// StartComponent starts a docker-compose component.
func StartComponent(component string, conditions ...ConditionFunc) error {
	cmdDown := fmt.Sprintf("stop-%s", component)
	if err := RunMakeTarget(cmdDown); err != nil {
		return err
	}
	cmdUp := fmt.Sprintf("run-%s", component)
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

// StopComponent stops a docker-compose component.
func StopComponent(component string) error {
	cmdDown := fmt.Sprintf("stop-%s", component)
	return RunMakeTarget(cmdDown)
}

// RunMakeTarget runs a Makefile target.
func RunMakeTarget(target string) error {
	cmd := exec.Command("make", target)
	return runCmd(cmd)
}

// GetDefaultOperationsConfig provides a default configuration to run the environment
func GetDefaultOperationsConfig() *Config {
	return &Config{
		State:     &state.Config{MaxCumulativeGasUsed: DefaultMaxCumulativeGasUsed},
		Sequencer: &SequencerConfig{Address: DefaultSequencerAddress, PrivateKey: DefaultSequencerPrivateKey},
	}
}

// GetClient returns an ethereum client to the provided URL
func GetClient(URL string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(URL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// MustGetClient GetClient but panic if err
func MustGetClient(URL string) *ethclient.Client {
	client, err := GetClient(URL)
	if err != nil {
		panic(err)
	}
	return client
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
	if err := dbutils.InitOrResetPool(poolDBCfg); err != nil {
		panic(err)
	}
	if err := dbutils.InitOrResetRPC(rpcDBCfg); err != nil {
		panic(err)
	}
}
