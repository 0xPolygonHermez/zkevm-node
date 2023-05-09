package etherman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
	"github.com/0xPolygonHermez/zkevm-node/etherman/ethgasstation"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmglobalexitroot"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

var (
	updateGlobalExitRootSignatureHash              = crypto.Keccak256Hash([]byte("UpdateGlobalExitRoot(bytes32,bytes32)"))
	forcedBatchSignatureHash                       = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)"))
	sequencedBatchesEventSignatureHash             = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	forceSequencedBatchesSignatureHash             = crypto.Keccak256Hash([]byte("SequenceForceBatches(uint64)"))
	verifyBatchesSignatureHash                     = crypto.Keccak256Hash([]byte("VerifyBatches(uint64,bytes32,address)"))
	verifyBatchesTrustedAggregatorSignatureHash    = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	setTrustedSequencerURLSignatureHash            = crypto.Keccak256Hash([]byte("SetTrustedSequencerURL(string)"))
	setTrustedSequencerSignatureHash               = crypto.Keccak256Hash([]byte("SetTrustedSequencer(address)"))
	transferOwnershipSignatureHash                 = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))
	emergencyStateActivatedSignatureHash           = crypto.Keccak256Hash([]byte("EmergencyStateActivated()"))
	emergencyStateDeactivatedSignatureHash         = crypto.Keccak256Hash([]byte("EmergencyStateDeactivated()"))
	updateZkEVMVersionSignatureHash                = crypto.Keccak256Hash([]byte("UpdateZkEVMVersion(uint64,uint64,string)"))
	consolidatePendingStateSignatureHash           = crypto.Keccak256Hash([]byte("ConsolidatePendingState(uint64,bytes32,uint64)"))
	setTrustedAggregatorTimeoutSignatureHash       = crypto.Keccak256Hash([]byte("SetTrustedAggregatorTimeout(uint64)"))
	setTrustedAggregatorSignatureHash              = crypto.Keccak256Hash([]byte("SetTrustedAggregator(address)"))
	setPendingStateTimeoutSignatureHash            = crypto.Keccak256Hash([]byte("SetPendingStateTimeout(uint64)"))
	setMultiplierBatchFeeSignatureHash             = crypto.Keccak256Hash([]byte("SetMultiplierBatchFee(uint16)"))
	setVerifyBatchTimeTargetSignatureHash          = crypto.Keccak256Hash([]byte("SetVerifyBatchTimeTarget(uint64)"))
	setForceBatchTimeoutSignatureHash              = crypto.Keccak256Hash([]byte("SetForceBatchTimeout(uint64)"))
	activateForceBatchesSignatureHash              = crypto.Keccak256Hash([]byte("ActivateForceBatches()"))
	transferAdminRoleSignatureHash                 = crypto.Keccak256Hash([]byte("TransferAdminRole(address)"))
	acceptAdminRoleSignatureHash                   = crypto.Keccak256Hash([]byte("AcceptAdminRole(address)"))
	proveNonDeterministicPendingStateSignatureHash = crypto.Keccak256Hash([]byte("ProveNonDeterministicPendingState(bytes32,bytes32)"))
	overridePendingStateSignatureHash              = crypto.Keccak256Hash([]byte("OverridePendingState(uint64,bytes32,address)"))

	// Proxy events
	initializedSignatureHash    = crypto.Keccak256Hash([]byte("Initialized(uint8)"))
	adminChangedSignatureHash   = crypto.Keccak256Hash([]byte("AdminChanged(address,address)"))
	beaconUpgradedSignatureHash = crypto.Keccak256Hash([]byte("BeaconUpgraded(address)"))
	upgradedSignatureHash       = crypto.Keccak256Hash([]byte("Upgraded(address)"))

	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("not found")
	// ErrIsReadOnlyMode is used when the EtherMan client is in read-only mode.
	ErrIsReadOnlyMode = errors.New("etherman client in read-only mode: no account configured to send transactions to L1. " +
		"please check the [Etherman] PrivateKeyPath and PrivateKeyPassword configuration")
	// ErrPrivateKeyNotFound used when the provided sender does not have a private key registered to be used
	ErrPrivateKeyNotFound = errors.New("can't find sender private key to sign tx")
)

// SequencedBatchesSigHash returns the hash for the `SequenceBatches` event.
func SequencedBatchesSigHash() common.Hash { return sequencedBatchesEventSignatureHash }

// TrustedVerifyBatchesSigHash returns the hash for the `TrustedVerifyBatches` event.
func TrustedVerifyBatchesSigHash() common.Hash { return verifyBatchesTrustedAggregatorSignatureHash }

// EventOrder is the the type used to identify the events order
type EventOrder string

const (
	// GlobalExitRootsOrder identifies a GlobalExitRoot event
	GlobalExitRootsOrder EventOrder = "GlobalExitRoots"
	// SequenceBatchesOrder identifies a VerifyBatch event
	SequenceBatchesOrder EventOrder = "SequenceBatches"
	// ForcedBatchesOrder identifies a ForcedBatches event
	ForcedBatchesOrder EventOrder = "ForcedBatches"
	// TrustedVerifyBatchOrder identifies a TrustedVerifyBatch event
	TrustedVerifyBatchOrder EventOrder = "TrustedVerifyBatch"
	// SequenceForceBatchesOrder identifies a SequenceForceBatches event
	SequenceForceBatchesOrder EventOrder = "SequenceForceBatches"
	// ForkIDsOrder identifies an updateZkevmVersion event
	ForkIDsOrder EventOrder = "forkIDs"
)

type ethereumClient interface {
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.LogFilterer
	ethereum.TransactionReader
	ethereum.TransactionSender

	bind.DeployBackend
}

// L1Config represents the configuration of the network used in L1
type L1Config struct {
	L1ChainID                 uint64         `json:"chainId"`
	ZkEVMAddr                 common.Address `json:"polygonZkEVMAddress"`
	MaticAddr                 common.Address `json:"maticTokenAddress"`
	GlobalExitRootManagerAddr common.Address `json:"polygonZkEVMGlobalExitRootAddress"`
}

type externalGasProviders struct {
	MultiGasProvider bool
	Providers        []ethereum.GasPricer
}

// Client is a simple implementation of EtherMan.
type Client struct {
	EthClient             ethereumClient
	ZkEVM                 *polygonzkevm.Polygonzkevm
	GlobalExitRootManager *polygonzkevmglobalexitroot.Polygonzkevmglobalexitroot
	Matic                 *matic.Matic
	SCAddresses           []common.Address

	GasProviders externalGasProviders

	l1Cfg L1Config
	auth  map[common.Address]bind.TransactOpts // empty in case of read-only client
}

// NewClient creates a new etherman.
func NewClient(cfg Config, l1Config L1Config) (*Client, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	// Create smc clients
	poe, err := polygonzkevm.NewPolygonzkevm(l1Config.ZkEVMAddr, ethClient)
	if err != nil {
		return nil, err
	}
	globalExitRoot, err := polygonzkevmglobalexitroot.NewPolygonzkevmglobalexitroot(l1Config.GlobalExitRootManagerAddr, ethClient)
	if err != nil {
		return nil, err
	}
	matic, err := matic.NewMatic(l1Config.MaticAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, l1Config.ZkEVMAddr, l1Config.GlobalExitRootManagerAddr)

	gProviders := []ethereum.GasPricer{ethClient}
	if cfg.MultiGasProvider {
		if cfg.Etherscan.ApiKey == "" {
			log.Info("No ApiKey provided for etherscan. Ignoring provider...")
		} else {
			log.Info("ApiKey detected for etherscan")
			gProviders = append(gProviders, etherscan.NewEtherscanService(cfg.Etherscan.ApiKey))
		}
		gProviders = append(gProviders, ethgasstation.NewEthGasStationService())
	}

	return &Client{
		EthClient:             ethClient,
		ZkEVM:                 poe,
		Matic:                 matic,
		GlobalExitRootManager: globalExitRoot,
		SCAddresses:           scAddresses,
		GasProviders: externalGasProviders{
			MultiGasProvider: cfg.MultiGasProvider,
			Providers:        gProviders,
		},
		l1Cfg: l1Config,
		auth:  map[common.Address]bind.TransactOpts{},
	}, nil
}

// VerifyGenBlockNumber verifies if the genesis Block Number is valid
func (etherMan *Client) VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error) {
	log.Info("Verifying genesis blockNumber: ", genBlockNumber)
	// Filter query
	genBlock := new(big.Int).SetUint64(genBlockNumber)
	query := ethereum.FilterQuery{
		FromBlock: genBlock,
		ToBlock:   genBlock,
		Addresses: etherMan.SCAddresses,
		Topics:    [][]common.Hash{{updateZkEVMVersionSignatureHash}},
	}
	logs, err := etherMan.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return false, err
	}
	if len(logs) == 0 {
		return false, fmt.Errorf("the specified genBlockNumber in config file does not contain any forkID event. Please use the proper blockNumber.")
	}
	zkevmVersion, err := etherMan.ZkEVM.ParseUpdateZkEVMVersion(logs[0])
	if err != nil {
		log.Error("error parsing the forkID event")
		return false, err
	}
	if zkevmVersion.NumBatch != 0 {
		return false, fmt.Errorf("the specified genBlockNumber in config file does not contain the initial forkID event (BatchNum: %d). Please use the proper blockNumber.", zkevmVersion.NumBatch)
	}
	return true, nil
}

// GetForks returns fork information
func (etherMan *Client) GetForks(ctx context.Context, genBlockNumber uint64) ([]state.ForkIDInterval, error) {
	log.Debug("Getting forkIDs from blockNumber: ", genBlockNumber)
	// Filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(genBlockNumber),
		Addresses: etherMan.SCAddresses,
		Topics:    [][]common.Hash{{updateZkEVMVersionSignatureHash}},
	}
	logs, err := etherMan.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.ForkIDInterval{}, err
	}
	var forks []state.ForkIDInterval
	for i, l := range logs {
		zkevmVersion, err := etherMan.ZkEVM.ParseUpdateZkEVMVersion(l)
		if err != nil {
			return []state.ForkIDInterval{}, err
		}
		var fork state.ForkIDInterval
		if i == 0 {
			fork = state.ForkIDInterval{
				FromBatchNumber: zkevmVersion.NumBatch + 1,
				ToBatchNumber:   math.MaxUint64,
				ForkId:          zkevmVersion.ForkID,
				Version:         zkevmVersion.Version,
			}
		} else {
			forks[len(forks)-1].ToBatchNumber = zkevmVersion.NumBatch
			fork = state.ForkIDInterval{
				FromBatchNumber: zkevmVersion.NumBatch + 1,
				ToBatchNumber:   math.MaxUint64,
				ForkId:          zkevmVersion.ForkID,
				Version:         zkevmVersion.Version,
			}
		}
		forks = append(forks, fork)
	}
	log.Debugf("Forks decoded: %+v", forks)
	return forks, nil
}

// GetRollupInfoByBlockRange function retrieves the Rollup information that are included in all this ethereum blocks
// from block x to block y.
func (etherMan *Client) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]Block, map[common.Hash][]Order, error) {
	// Filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: etherMan.SCAddresses,
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, blocksOrder, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return blocks, blocksOrder, nil
}

// Order contains the event order to let the synchronizer store the information following this order.
type Order struct {
	Name EventOrder
	Pos  int
}

func (etherMan *Client) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]Block, map[common.Hash][]Order, error) {
	logs, err := etherMan.EthClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	var blocks []Block
	blocksOrder := make(map[common.Hash][]Order)
	for _, vLog := range logs {
		err := etherMan.processEvent(ctx, vLog, &blocks, &blocksOrder)
		if err != nil {
			log.Warnf("error processing event. Retrying... Error: %s. vLog: %+v", err.Error(), vLog)
			return nil, nil, err
		}
	}
	return blocks, blocksOrder, nil
}

func (etherMan *Client) processEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	switch vLog.Topics[0] {
	case sequencedBatchesEventSignatureHash:
		return etherMan.sequencedBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case updateGlobalExitRootSignatureHash:
		return etherMan.updateGlobalExitRootEvent(ctx, vLog, blocks, blocksOrder)
	case forcedBatchSignatureHash:
		return etherMan.forcedBatchEvent(ctx, vLog, blocks, blocksOrder)
	case verifyBatchesTrustedAggregatorSignatureHash:
		return etherMan.verifyBatchesTrustedAggregatorEvent(ctx, vLog, blocks, blocksOrder)
	case verifyBatchesSignatureHash:
		log.Warn("VerifyBatches event not implemented yet")
		return nil
	case forceSequencedBatchesSignatureHash:
		return etherMan.forceSequencedBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case setTrustedSequencerURLSignatureHash:
		log.Debug("SetTrustedSequencerURL event detected")
		return nil
	case setTrustedSequencerSignatureHash:
		log.Debug("SetTrustedSequencer event detected")
		return nil
	case initializedSignatureHash:
		log.Debug("Initialized event detected")
		return nil
	case adminChangedSignatureHash:
		log.Debug("AdminChanged event detected")
		return nil
	case beaconUpgradedSignatureHash:
		log.Debug("BeaconUpgraded event detected")
		return nil
	case upgradedSignatureHash:
		log.Debug("Upgraded event detected")
		return nil
	case transferOwnershipSignatureHash:
		log.Debug("TransferOwnership event detected")
		return nil
	case emergencyStateActivatedSignatureHash:
		log.Debug("EmergencyStateActivated event detected")
		return nil
	case emergencyStateDeactivatedSignatureHash:
		log.Debug("EmergencyStateDeactivated event detected")
		return nil
	case updateZkEVMVersionSignatureHash:
		return etherMan.updateZkevmVersion(ctx, vLog, blocks, blocksOrder)
	case consolidatePendingStateSignatureHash:
		log.Debug("ConsolidatePendingState event detected")
		return nil
	case setTrustedAggregatorTimeoutSignatureHash:
		log.Debug("SetTrustedAggregatorTimeout event detected")
		return nil
	case setTrustedAggregatorSignatureHash:
		log.Debug("setTrustedAggregator event detected")
		return nil
	case setPendingStateTimeoutSignatureHash:
		log.Debug("SetPendingStateTimeout event detected")
		return nil
	case setMultiplierBatchFeeSignatureHash:
		log.Debug("SetMultiplierBatchFee event detected")
		return nil
	case setVerifyBatchTimeTargetSignatureHash:
		log.Debug("SetVerifyBatchTimeTarget event detected")
		return nil
	case setForceBatchTimeoutSignatureHash:
		log.Debug("SetForceBatchTimeout event detected")
		return nil
	case activateForceBatchesSignatureHash:
		log.Debug("ActivateForceBatches event detected")
		return nil
	case transferAdminRoleSignatureHash:
		log.Debug("TransferAdminRole event detected")
		return nil
	case acceptAdminRoleSignatureHash:
		log.Debug("AcceptAdminRole event detected")
		return nil
	case proveNonDeterministicPendingStateSignatureHash:
		log.Debug("ProveNonDeterministicPendingState event detected")
		return nil
	case overridePendingStateSignatureHash:
		log.Debug("OverridePendingState event detected")
		return nil
	}
	log.Warn("Event not registered: ", vLog)
	return nil
}

func (etherMan *Client) updateZkevmVersion(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateZkEVMVersion event detected")
	zkevmVersion, err := etherMan.ZkEVM.ParseUpdateZkEVMVersion(vLog)
	if err != nil {
		log.Error("error parsing UpdateZkEVMVersion event. Error: ", err)
		return err
	}
	fork := ForkID{
		BatchNumber: zkevmVersion.NumBatch,
		ForkID:      zkevmVersion.ForkID,
		Version:     zkevmVersion.Version,
	}
	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		t := time.Unix(int64(fullBlock.Time()), 0)
		block := prepareBlock(vLog, t, fullBlock)
		block.ForkIDs = append(block.ForkIDs, fork)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].ForkIDs = append((*blocks)[len(*blocks)-1].ForkIDs, fork)
	} else {
		log.Error("Error processing updateZkevmVersion event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing updateZkevmVersion event")
	}
	or := Order{
		Name: ForkIDsOrder,
		Pos:  len((*blocks)[len(*blocks)-1].ForkIDs) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) updateGlobalExitRootEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateGlobalExitRoot event detected")
	globalExitRoot, err := etherMan.GlobalExitRootManager.ParseUpdateGlobalExitRoot(vLog)
	if err != nil {
		return err
	}
	var gExitRoot GlobalExitRoot
	gExitRoot.MainnetExitRoot = common.BytesToHash(globalExitRoot.MainnetExitRoot[:])
	gExitRoot.RollupExitRoot = common.BytesToHash(globalExitRoot.RollupExitRoot[:])
	gExitRoot.BlockNumber = vLog.BlockNumber
	gExitRoot.GlobalExitRoot = hash(globalExitRoot.MainnetExitRoot, globalExitRoot.RollupExitRoot)

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		t := time.Unix(int64(fullBlock.Time()), 0)
		block := prepareBlock(vLog, t, fullBlock)
		block.GlobalExitRoots = append(block.GlobalExitRoots, gExitRoot)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].GlobalExitRoots = append((*blocks)[len(*blocks)-1].GlobalExitRoots, gExitRoot)
	} else {
		log.Error("Error processing UpdateGlobalExitRoot event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing UpdateGlobalExitRoot event")
	}
	or := Order{
		Name: GlobalExitRootsOrder,
		Pos:  len((*blocks)[len(*blocks)-1].GlobalExitRoots) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

// WaitTxToBeMined waits for an L1 tx to be mined. It will return error if the tx is reverted or timeout is exceeded
func (etherMan *Client) WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) (bool, error) {
	err := operations.WaitTxToBeMined(ctx, etherMan.EthClient, tx, timeout)
	if errors.Is(err, context.DeadlineExceeded) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// EstimateGasSequenceBatches estimates gas for sending batches
func (etherMan *Client) EstimateGasSequenceBatches(sender common.Address, sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, ErrPrivateKeyNotFound
	}
	opts.NoSend = true

	tx, err := etherMan.sequenceBatches(opts, sequences)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// BuildSequenceBatchesTxData builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceBatches(opts, sequences)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceBatches(opts bind.TransactOpts, sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	var batches []polygonzkevm.PolygonZkEVMBatchData
	for _, seq := range sequences {
		batch := polygonzkevm.PolygonZkEVMBatchData{
			Transactions:       seq.BatchL2Data,
			GlobalExitRoot:     seq.GlobalExitRoot,
			Timestamp:          uint64(seq.Timestamp),
			MinForcedTimestamp: uint64(seq.ForcedBatchTimestamp),
		}

		batches = append(batches, batch)
	}

	tx, err := etherMan.ZkEVM.SequenceBatches(&opts, batches, opts.From)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
	}

	return tx, err
}

// BuildTrustedVerifyBatchesTxData builds a []bytes to be sent to the PoE SC method TrustedVerifyBatches.
func (etherMan *Client) BuildTrustedVerifyBatchesTxData(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.generateRandomAuth()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build trusted verify batches, err: %w", err)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	var newLocalExitRoot [32]byte
	copy(newLocalExitRoot[:], inputs.NewLocalExitRoot)

	var newStateRoot [32]byte
	copy(newStateRoot[:], inputs.NewStateRoot)

	proof, err := encoding.DecodeBytes(&inputs.FinalProof.Proof)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode proof, err: %w", err)
	}

	const pendStateNum = 0 // TODO hardcoded for now until we implement the pending state feature

	tx, err := etherMan.ZkEVM.VerifyBatchesTrustedAggregator(
		&opts,
		pendStateNum,
		lastVerifiedBatch,
		newVerifiedBatch,
		newLocalExitRoot,
		newStateRoot,
		proof,
	)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

// GetSendSequenceFee get super/trusted sequencer fee
func (etherMan *Client) GetSendSequenceFee(numBatches uint64) (*big.Int, error) {
	f, err := etherMan.ZkEVM.BatchFee(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, err
	}
	fee := new(big.Int).Mul(f, new(big.Int).SetUint64(numBatches))
	return fee, nil
}

// TrustedSequencer gets trusted sequencer address
func (etherMan *Client) TrustedSequencer() (common.Address, error) {
	return etherMan.ZkEVM.TrustedSequencer(&bind.CallOpts{Pending: false})
}

func (etherMan *Client) forcedBatchEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("ForceBatch event detected")
	fb, err := etherMan.ZkEVM.ParseForceBatch(vLog)
	if err != nil {
		return err
	}
	var forcedBatch ForcedBatch
	forcedBatch.BlockNumber = vLog.BlockNumber
	forcedBatch.ForcedBatchNumber = fb.ForceBatchNum
	forcedBatch.GlobalExitRoot = fb.LastGlobalExitRoot
	// Read the tx for this batch.
	tx, isPending, err := etherMan.EthClient.TransactionByHash(ctx, vLog.TxHash)
	if err != nil {
		return err
	} else if isPending {
		return fmt.Errorf("error: tx is still pending. TxHash: %s", tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	if fb.Sequencer == msg.From {
		txData := tx.Data()
		// Extract coded txs.
		// Load contract ABI
		abi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
		if err != nil {
			return err
		}

		// Recover Method from signature and ABI
		method, err := abi.MethodById(txData[:4])
		if err != nil {
			return err
		}

		// Unpack method inputs
		data, err := method.Inputs.Unpack(txData[4:])
		if err != nil {
			return err
		}
		bytedata := data[0].([]byte)
		forcedBatch.RawTxsData = bytedata
	} else {
		forcedBatch.RawTxsData = fb.Transactions
	}
	forcedBatch.Sequencer = fb.Sequencer
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	t := time.Unix(int64(fullBlock.Time()), 0)
	forcedBatch.ForcedAt = t
	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, t, fullBlock)
		block.ForcedBatches = append(block.ForcedBatches, forcedBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].ForcedBatches = append((*blocks)[len(*blocks)-1].ForcedBatches, forcedBatch)
	} else {
		log.Error("Error processing ForceBatch event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing ForceBatch event")
	}
	or := Order{
		Name: ForcedBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].ForcedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) sequencedBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("SequenceBatches event detected")
	sb, err := etherMan.ZkEVM.ParseSequenceBatches(vLog)
	if err != nil {
		return err
	}
	// Read the tx for this event.
	tx, isPending, err := etherMan.EthClient.TransactionByHash(ctx, vLog.TxHash)
	if err != nil {
		return err
	} else if isPending {
		return fmt.Errorf("error tx is still pending. TxHash: %s", tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	sequences, err := decodeSequences(tx.Data(), sb.NumBatch, msg.From, vLog.TxHash, msg.Nonce)
	if err != nil {
		return fmt.Errorf("error decoding the sequences: %v", err)
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedBatches = append(block.SequencedBatches, sequences)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedBatches = append((*blocks)[len(*blocks)-1].SequencedBatches, sequences)
	} else {
		log.Error("Error processing SequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing SequencedBatches event")
	}
	or := Order{
		Name: SequenceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func decodeSequences(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}
	var sequences []polygonzkevm.PolygonZkEVMBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}
	coinbase := (data[1]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))
	for i, seq := range sequences {
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:           bn,
			SequencerAddr:         sequencer,
			TxHash:                txHash,
			Nonce:                 nonce,
			Coinbase:              coinbase,
			PolygonZkEVMBatchData: seq,
		}
	}

	return sequencedBatches, nil
}

func (etherMan *Client) verifyBatchesTrustedAggregatorEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("TrustedVerifyBatches event detected")
	vb, err := etherMan.ZkEVM.ParseVerifyBatchesTrustedAggregator(vLog)
	if err != nil {
		return err
	}
	var trustedVerifyBatch VerifiedBatch
	trustedVerifyBatch.BlockNumber = vLog.BlockNumber
	trustedVerifyBatch.BatchNumber = vb.NumBatch
	trustedVerifyBatch.TxHash = vLog.TxHash
	trustedVerifyBatch.StateRoot = vb.StateRoot
	trustedVerifyBatch.Aggregator = vb.Aggregator

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.VerifiedBatches = append(block.VerifiedBatches, trustedVerifyBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].VerifiedBatches = append((*blocks)[len(*blocks)-1].VerifiedBatches, trustedVerifyBatch)
	} else {
		log.Error("Error processing trustedVerifyBatch event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing trustedVerifyBatch event")
	}
	or := Order{
		Name: TrustedVerifyBatchOrder,
		Pos:  len((*blocks)[len(*blocks)-1].VerifiedBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)
	return nil
}

func (etherMan *Client) forceSequencedBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("SequenceForceBatches event detect")
	fsb, err := etherMan.ZkEVM.ParseSequenceForceBatches(vLog)
	if err != nil {
		return err
	}

	// Read the tx for this batch.
	tx, isPending, err := etherMan.EthClient.TransactionByHash(ctx, vLog.TxHash)
	if err != nil {
		return err
	} else if isPending {
		return fmt.Errorf("error: tx is still pending. TxHash: %s", tx.Hash().String())
	}
	msg, err := core.TransactionToMessage(tx, types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	sequencedForceBatch, err := decodeSequencedForceBatches(tx.Data(), fsb.NumBatch, msg.From, vLog.TxHash, fullBlock, msg.Nonce)
	if err != nil {
		return err
	}

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
		block := prepareBlock(vLog, time.Unix(int64(fullBlock.Time()), 0), fullBlock)
		block.SequencedForceBatches = append(block.SequencedForceBatches, sequencedForceBatch)
		*blocks = append(*blocks, block)
	} else if (*blocks)[len(*blocks)-1].BlockHash == vLog.BlockHash && (*blocks)[len(*blocks)-1].BlockNumber == vLog.BlockNumber {
		(*blocks)[len(*blocks)-1].SequencedForceBatches = append((*blocks)[len(*blocks)-1].SequencedForceBatches, sequencedForceBatch)
	} else {
		log.Error("Error processing ForceSequencedBatches event. BlockHash:", vLog.BlockHash, ". BlockNumber: ", vLog.BlockNumber)
		return fmt.Errorf("error processing ForceSequencedBatches event")
	}
	or := Order{
		Name: SequenceForceBatchesOrder,
		Pos:  len((*blocks)[len(*blocks)-1].SequencedForceBatches) - 1,
	}
	(*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash] = append((*blocksOrder)[(*blocks)[len(*blocks)-1].BlockHash], or)

	return nil
}

func decodeSequencedForceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, block *types.Block, nonce uint64) ([]SequencedForceBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		return nil, err
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}

	var forceBatches []polygonzkevm.PolygonZkEVMForcedBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &forceBatches)
	if err != nil {
		return nil, err
	}

	sequencedForcedBatches := make([]SequencedForceBatch, len(forceBatches))
	for i, force := range forceBatches {
		bn := lastBatchNumber - uint64(len(forceBatches)-(i+1))
		sequencedForcedBatches[i] = SequencedForceBatch{
			BatchNumber:                 bn,
			Coinbase:                    sequencer,
			TxHash:                      txHash,
			Timestamp:                   time.Unix(int64(block.Time()), 0),
			Nonce:                       nonce,
			PolygonZkEVMForcedBatchData: force,
		}
	}
	return sequencedForcedBatches, nil
}

func prepareBlock(vLog types.Log, t time.Time, fullBlock *types.Block) Block {
	var block Block
	block.BlockNumber = vLog.BlockNumber
	block.BlockHash = vLog.BlockHash
	block.ParentHash = fullBlock.ParentHash()
	block.ReceivedAt = t
	return block
}

func hash(data ...[32]byte) [32]byte {
	var res [32]byte
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		hash.Write(d[:]) //nolint:errcheck,gosec
	}
	copy(res[:], hash.Sum(nil))
	return res
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (etherMan *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return etherMan.EthClient.HeaderByNumber(ctx, number)
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number.
func (etherMan *Client) EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := etherMan.EthClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) || err.Error() == "block does not exist in blockchain" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return block, nil
}

// GetLastBatchTimestamp function allows to retrieve the lastTimestamp value in the smc
func (etherMan *Client) GetLastBatchTimestamp() (uint64, error) {
	return etherMan.ZkEVM.LastTimestamp(&bind.CallOpts{Pending: false})
}

// GetLatestBatchNumber function allows to retrieve the latest proposed batch in the smc
func (etherMan *Client) GetLatestBatchNumber() (uint64, error) {
	return etherMan.ZkEVM.LastBatchSequenced(&bind.CallOpts{Pending: false})
}

// GetLatestBlockNumber gets the latest block number from the ethereum
func (etherMan *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	header, err := etherMan.EthClient.HeaderByNumber(ctx, nil)
	if err != nil || header == nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// GetLatestBlockTimestamp gets the latest block timestamp from the ethereum
func (etherMan *Client) GetLatestBlockTimestamp(ctx context.Context) (uint64, error) {
	header, err := etherMan.EthClient.HeaderByNumber(ctx, nil)
	if err != nil || header == nil {
		return 0, err
	}
	return header.Time, nil
}

// GetLatestVerifiedBatchNum gets latest verified batch from ethereum
func (etherMan *Client) GetLatestVerifiedBatchNum() (uint64, error) {
	return etherMan.ZkEVM.LastVerifiedBatch(&bind.CallOpts{Pending: false})
}

// GetTx function get ethereum tx
func (etherMan *Client) GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	return etherMan.EthClient.TransactionByHash(ctx, txHash)
}

// GetTxReceipt function gets ethereum tx receipt
func (etherMan *Client) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return etherMan.EthClient.TransactionReceipt(ctx, txHash)
}

// ApproveMatic function allow to approve tokens in matic smc
func (etherMan *Client) ApproveMatic(ctx context.Context, account common.Address, maticAmount *big.Int, to common.Address) (*types.Transaction, error) {
	opts, err := etherMan.getAuthByAddress(account)
	if err == ErrNotFound {
		return nil, errors.New("can't find account private key to sign tx")
	}
	if etherMan.GasProviders.MultiGasProvider {
		opts.GasPrice = etherMan.GetL1GasPrice(ctx)
	}
	tx, err := etherMan.Matic.Approve(&opts, etherMan.l1Cfg.ZkEVMAddr, maticAmount)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		return nil, fmt.Errorf("error approving balance to send the batch. Error: %w", err)
	}

	return tx, nil
}

// GetTrustedSequencerURL Gets the trusted sequencer url from rollup smc
func (etherMan *Client) GetTrustedSequencerURL() (string, error) {
	return etherMan.ZkEVM.TrustedSequencerURL(&bind.CallOpts{Pending: false})
}

// GetL2ChainID returns L2 Chain ID
func (etherMan *Client) GetL2ChainID() (uint64, error) {
	return etherMan.ZkEVM.ChainID(&bind.CallOpts{Pending: false})
}

// GetL2ForkID returns current L2 Fork ID
func (etherMan *Client) GetL2ForkID() (uint64, error) {
	// TODO: implement this
	return 1, nil
}

// GetL2ForkIDIntervals return L2 Fork ID intervals
func (etherMan *Client) GetL2ForkIDIntervals() ([]state.ForkIDInterval, error) {
	// TODO: implement this
	return []state.ForkIDInterval{{FromBatchNumber: 0, ToBatchNumber: math.MaxUint64, ForkId: 1}}, nil
}

// GetL1GasPrice gets the l1 gas price
func (etherMan *Client) GetL1GasPrice(ctx context.Context) *big.Int {
	// Get gasPrice from providers
	gasPrice := big.NewInt(0)
	for i, prov := range etherMan.GasProviders.Providers {
		gp, err := prov.SuggestGasPrice(ctx)
		if err != nil {
			log.Warnf("error getting gas price from provider %d. Error: %s", i+1, err.Error())
		} else if gasPrice.Cmp(gp) == -1 { // gasPrice < gp
			gasPrice = gp
		}
	}
	log.Debug("gasPrice chose: ", gasPrice)
	return gasPrice
}

// SendTx sends a tx to L1
func (etherMan *Client) SendTx(ctx context.Context, tx *types.Transaction) error {
	return etherMan.EthClient.SendTransaction(ctx, tx)
}

// CurrentNonce returns the current nonce for the provided account
func (etherMan *Client) CurrentNonce(ctx context.Context, account common.Address) (uint64, error) {
	return etherMan.EthClient.NonceAt(ctx, account, nil)
}

// SuggestedGasPrice returns the suggest nonce for the network at the moment
func (etherMan *Client) SuggestedGasPrice(ctx context.Context) (*big.Int, error) {
	suggestedGasPrice := etherMan.GetL1GasPrice(ctx)
	if suggestedGasPrice.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("failed to get the suggested gas price")
	}
	return suggestedGasPrice, nil
}

// EstimateGas returns the estimated gas for the tx
func (etherMan *Client) EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error) {
	return etherMan.EthClient.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	})
}

// CheckTxWasMined check if a tx was already mined
func (etherMan *Client) CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error) {
	receipt, err := etherMan.EthClient.TransactionReceipt(ctx, txHash)
	if errors.Is(err, ethereum.NotFound) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}

	return true, receipt, nil
}

// SignTx tries to sign a transaction accordingly to the provided sender
func (etherMan *Client) SignTx(ctx context.Context, sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	auth, err := etherMan.getAuthByAddress(sender)
	if err == ErrNotFound {
		return nil, ErrPrivateKeyNotFound
	}
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// GetRevertMessage tries to get a revert message of a transaction
func (etherMan *Client) GetRevertMessage(ctx context.Context, tx *types.Transaction) (string, error) {
	if tx == nil {
		return "", nil
	}

	receipt, err := etherMan.GetTxReceipt(ctx, tx.Hash())
	if err != nil {
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		revertMessage, err := operations.RevertReason(ctx, etherMan.EthClient, tx, receipt.BlockNumber)
		if err != nil {
			return "", err
		}
		return revertMessage, nil
	}
	return "", nil
}

// AddOrReplaceAuth adds an authorization or replace an existent one to the same account
func (etherMan *Client) AddOrReplaceAuth(auth bind.TransactOpts) error {
	log.Infof("added or replaced authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return nil
}

// LoadAuthFromKeyStore loads an authorization from a key store file
func (etherMan *Client) LoadAuthFromKeyStore(path, password string) (*bind.TransactOpts, error) {
	auth, err := newAuthFromKeystore(path, password, etherMan.l1Cfg.L1ChainID)
	if err != nil {
		return nil, err
	}

	log.Infof("loaded authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return &auth, nil
}

// newKeyFromKeystore creates an instance of a keystore key from a keystore file
func newKeyFromKeystore(path, password string) (*keystore.Key, error) {
	if path == "" && password == "" {
		return nil, nil
	}
	keystoreEncrypted, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	log.Infof("decrypting key from: %v", path)
	key, err := keystore.DecryptKey(keystoreEncrypted, password)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// newAuthFromKeystore an authorization instance from a keystore file
func newAuthFromKeystore(path, password string, chainID uint64) (bind.TransactOpts, error) {
	log.Infof("reading key from: %v", path)
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		return bind.TransactOpts{}, err
	}
	if key == nil {
		return bind.TransactOpts{}, nil
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		return bind.TransactOpts{}, err
	}
	return *auth, nil
}

// getAuthByAddress tries to get an authorization from the authorizations map
func (etherMan *Client) getAuthByAddress(addr common.Address) (bind.TransactOpts, error) {
	auth, found := etherMan.auth[addr]
	if !found {
		return bind.TransactOpts{}, ErrNotFound
	}
	return auth, nil
}

// generateRandomAuth generates an authorization instance from a
// randomly generated private key to be used to estimate gas for PoE
// operations NOT restricted to the Trusted Sequencer
func (etherMan *Client) generateRandomAuth() (bind.TransactOpts, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a private key to estimate L1 txs")
	}
	chainID := big.NewInt(0).SetUint64(etherMan.l1Cfg.L1ChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return bind.TransactOpts{}, errors.New("failed to generate a fake authorization to estimate L1 txs")
	}

	return *auth, nil
}
