package etherman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"
	"github.com/0xPolygonHermez/zkevm-node/etherman/ethgasstation"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/globalexitrootmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

var (
	updateGlobalExitRootSignatureHash      = crypto.Keccak256Hash([]byte("UpdateGlobalExitRoot(bytes32,bytes32)"))
	forcedBatchSignatureHash               = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)"))
	sequencedBatchesEventSignatureHash     = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	forceSequencedBatchesSignatureHash     = crypto.Keccak256Hash([]byte("SequenceForceBatches(uint64)"))
	verifyBatchesSignatureHash             = crypto.Keccak256Hash([]byte("VerifyBatches(uint64,bytes32,address)"))
	trustedVerifyBatchesSignatureHash      = crypto.Keccak256Hash([]byte("TrustedVerifyBatches(uint64,bytes32,address)"))
	setTrustedSequencerURLSignatureHash    = crypto.Keccak256Hash([]byte("SetTrustedSequencerURL(string)"))
	setForceBatchAllowedSignatureHash      = crypto.Keccak256Hash([]byte("SetForceBatchAllowed(bool)"))
	setTrustedSequencerSignatureHash       = crypto.Keccak256Hash([]byte("SetTrustedSequencer(address)"))
	transferOwnershipSignatureHash         = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))
	setSecurityCouncilSignatureHash        = crypto.Keccak256Hash([]byte("SetSecurityCouncil(address)"))
	proofDifferentStateSignatureHash       = crypto.Keccak256Hash([]byte("ProofDifferentState(bytes32,bytes32)"))
	emergencyStateActivatedSignatureHash   = crypto.Keccak256Hash([]byte("EmergencyStateActivated()"))
	emergencyStateDeactivatedSignatureHash = crypto.Keccak256Hash([]byte("EmergencyStateDeactivated()"))

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
)

// SequencedBatchesSigHash returns the hash for the `SequenceBatches` event.
func SequencedBatchesSigHash() common.Hash { return sequencedBatchesEventSignatureHash }

// TrustedVerifyBatchesSigHash returns the hash for the `TrustedVerifyBatches` event.
func TrustedVerifyBatchesSigHash() common.Hash { return trustedVerifyBatchesSignatureHash }

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

type externalGasProviders struct {
	MultiGasProvider bool
	Providers        []ethereum.GasPricer
}

// Client is a simple implementation of EtherMan.
type Client struct {
	EthClient             ethereumClient
	PoE                   *proofofefficiency.Proofofefficiency
	GlobalExitRootManager *globalexitrootmanager.Globalexitrootmanager
	Matic                 *matic.Matic
	SCAddresses           []common.Address

	GasProviders externalGasProviders

	auth *bind.TransactOpts // nil in case of read-only client
}

// NewClient creates a new etherman.
func NewClient(cfg Config, auth *bind.TransactOpts) (*Client, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	// Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoEAddr, ethClient)
	if err != nil {
		return nil, err
	}
	globalExitRoot, err := globalexitrootmanager.NewGlobalexitrootmanager(cfg.GlobalExitRootManagerAddr, ethClient)
	if err != nil {
		return nil, err
	}
	matic, err := matic.NewMatic(cfg.MaticAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, cfg.PoEAddr, cfg.GlobalExitRootManagerAddr)

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
		PoE:                   poe,
		Matic:                 matic,
		GlobalExitRootManager: globalExitRoot,
		SCAddresses:           scAddresses,
		GasProviders: externalGasProviders{
			MultiGasProvider: cfg.MultiGasProvider,
			Providers:        gProviders,
		},
		auth: auth,
	}, nil
}

// IsReadOnly returns whether the EtherMan client is in read-only mode.
// Call this before trying to access the `auth` field.
func (c *Client) IsReadOnly() bool { return c.auth == nil }

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
	case trustedVerifyBatchesSignatureHash:
		return etherMan.trustedVerifyBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case verifyBatchesSignatureHash:
		log.Warn("VerifyBatches event not implemented yet")
		return nil
	case forceSequencedBatchesSignatureHash:
		return etherMan.forceSequencedBatchesEvent(ctx, vLog, blocks, blocksOrder)
	case setTrustedSequencerURLSignatureHash:
		log.Debug("SetTrustedSequencerURL event detected")
		return nil
	case setForceBatchAllowedSignatureHash:
		log.Debug("SetForceBatchAllowed event detected")
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
	case setSecurityCouncilSignatureHash:
		log.Debug("SetSecurityCouncil event detected")
		return nil
	case proofDifferentStateSignatureHash:
		log.Debug("ProofDifferentState event detected")
		return nil
	case emergencyStateActivatedSignatureHash:
		log.Debug("EmergencyStateActivated event detected")
		return nil
	case emergencyStateDeactivatedSignatureHash:
		log.Debug("EmergencyStateDeactivated event detected")
		return nil
	}
	log.Warn("Event not registered: ", vLog)
	return nil
}

func (etherMan *Client) updateGlobalExitRootEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("UpdateGlobalExitRoot event detected")
	globalExitRoot, err := etherMan.GlobalExitRootManager.ParseUpdateGlobalExitRoot(vLog)
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	var gExitRoot GlobalExitRoot
	gExitRoot.MainnetExitRoot = common.BytesToHash(globalExitRoot.MainnetExitRoot[:])
	gExitRoot.RollupExitRoot = common.BytesToHash(globalExitRoot.RollupExitRoot[:])
	gExitRoot.BlockNumber = vLog.BlockNumber
	gExitRoot.GlobalExitRoot = hash(globalExitRoot.MainnetExitRoot, globalExitRoot.RollupExitRoot)

	if len(*blocks) == 0 || ((*blocks)[len(*blocks)-1].BlockHash != vLog.BlockHash || (*blocks)[len(*blocks)-1].BlockNumber != vLog.BlockNumber) {
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
func (etherMan *Client) WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error {
	return operations.WaitTxToBeMined(ctx, etherMan.EthClient, tx, timeout)
}

// EstimateGasSequenceBatches estimates gas for sending batches
func (etherMan *Client) EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	if etherMan.IsReadOnly() {
		return nil, ErrIsReadOnlyMode
	}
	noSendOpts := *etherMan.auth
	noSendOpts.NoSend = true
	tx, err := etherMan.sequenceBatches(&noSendOpts, sequences)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SequenceBatches send sequences of batches to the ethereum
func (etherMan *Client) SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error) {
	if etherMan.IsReadOnly() {
		return nil, ErrIsReadOnlyMode
	}
	sendSequencesOpts := *etherMan.auth
	sendSequencesOpts.GasLimit = gasLimit
	if gasPrice != nil {
		sendSequencesOpts.GasPrice = gasPrice
	} else if etherMan.GasProviders.MultiGasProvider {
		sendSequencesOpts.GasPrice = etherMan.getGasPrice(ctx)
	}
	if nonce != nil {
		sendSequencesOpts.Nonce = nonce
	}
	sendSequencesOpts.NoSend = noSend
	return etherMan.sequenceBatches(&sendSequencesOpts, sequences)
}

func (etherMan *Client) sequenceBatches(opts *bind.TransactOpts, sequences []ethmanTypes.Sequence) (*types.Transaction, error) {
	var batches []proofofefficiency.ProofOfEfficiencyBatchData
	for _, seq := range sequences {
		batchL2Data, err := state.EncodeTransactions(seq.Txs)
		if err != nil {
			return nil, fmt.Errorf("failed to encode transactions, err: %v", err)
		}
		batch := proofofefficiency.ProofOfEfficiencyBatchData{
			Transactions:       batchL2Data,
			GlobalExitRoot:     seq.GlobalExitRoot,
			Timestamp:          uint64(seq.Timestamp),
			MinForcedTimestamp: 0, // TODO If this batch is forced, this value must be different to zero. If it is a non forced sequence, then the valio will be valid
		}

		batches = append(batches, batch)
	}

	transaction, err := etherMan.PoE.SequenceBatches(opts, batches)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
	}

	return transaction, err
}

// EstimateGasForTrustedVerifyBatches estimates gas for trusted verify batches smart contract call.
func (etherMan *Client) EstimateGasForTrustedVerifyBatches(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (uint64, error) {
	if etherMan.IsReadOnly() {
		return 0, ErrIsReadOnlyMode
	}
	verifyBatchOpts := *etherMan.auth
	verifyBatchOpts.NoSend = true
	tx, err := etherMan.trustedVerifyBatches(&verifyBatchOpts, lastVerifiedBatch, newVerifiedBatch, inputs)
	if err != nil {
		return 0, err
	}
	return tx.Gas(), nil
}

// TrustedVerifyBatches function allows the aggregator send the final proof to L1.
func (etherMan *Client) TrustedVerifyBatches(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error) {
	if etherMan.IsReadOnly() {
		return nil, ErrIsReadOnlyMode
	}
	verifyBatchOpts := *etherMan.auth
	verifyBatchOpts.GasLimit = gasLimit
	if gasPrice != nil {
		verifyBatchOpts.GasPrice = gasPrice
	} else if etherMan.GasProviders.MultiGasProvider {
		verifyBatchOpts.GasPrice = etherMan.getGasPrice(ctx)
	}
	if nonce != nil {
		verifyBatchOpts.Nonce = nonce
	}
	verifyBatchOpts.NoSend = noSend
	return etherMan.trustedVerifyBatches(&verifyBatchOpts, lastVerifiedBatch, newVerifiedBatch, inputs)
}

func (etherMan *Client) trustedVerifyBatches(opts *bind.TransactOpts, lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (*types.Transaction, error) {
	var newLocalExitRoot [32]byte
	copy(newLocalExitRoot[:], inputs.NewLocalExitRoot)

	var newStateRoot [32]byte
	copy(newStateRoot[:], inputs.NewStateRoot)

	proofA, err := strSliceToBigIntArray(inputs.FinalProof.Proof.ProofA)
	if err != nil {
		return nil, err
	}
	proofB, err := proofSlcToIntArray(inputs.FinalProof.Proof.ProofB)
	if err != nil {
		return nil, err
	}
	proofC, err := strSliceToBigIntArray(inputs.FinalProof.Proof.ProofC)
	if err != nil {
		return nil, err
	}

	const pendStateNum = 0 // TODO hardcoded for now until we implement the pending state feature

	tx, err := etherMan.PoE.TrustedVerifyBatches(
		opts,
		pendStateNum,
		lastVerifiedBatch,
		newVerifiedBatch,
		newLocalExitRoot,
		newStateRoot,
		proofA,
		proofB,
		proofC,
	)
	if err != nil {
		if parsedErr, ok := tryParseError(err); ok {
			err = parsedErr
		}
		return nil, err
	}

	return tx, nil
}

// GetSendSequenceFee get super/trusted sequencer fee
func (etherMan *Client) GetSendSequenceFee(numBatches uint64) (*big.Int, error) {
	f, err := etherMan.PoE.GetCurrentBatchFee(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, err
	}
	fee := new(big.Int).Mul(f, new(big.Int).SetUint64(numBatches))
	return fee, nil
}

// TrustedSequencer gets trusted sequencer address
func (etherMan *Client) TrustedSequencer() (common.Address, error) {
	return etherMan.PoE.TrustedSequencer(&bind.CallOpts{Pending: false})
}

func (etherMan *Client) forcedBatchEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("ForceBatch event detected")
	fb, err := etherMan.PoE.ParseForceBatch(vLog)
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
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	if fb.Sequencer == msg.From() {
		txData := tx.Data()
		// Extract coded txs.
		// Load contract ABI
		abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
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
	sb, err := etherMan.PoE.ParseSequenceBatches(vLog)
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
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	sequences, err := decodeSequences(tx.Data(), sb.NumBatch, msg.From(), vLog.TxHash, msg.Nonce())
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
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
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
	var sequences []proofofefficiency.ProofOfEfficiencyBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}

	sequencedBatches := make([]SequencedBatch, len(sequences))
	for i, seq := range sequences {
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:                bn,
			Coinbase:                   sequencer,
			TxHash:                     txHash,
			Nonce:                      nonce,
			ProofOfEfficiencyBatchData: seq,
		}
	}

	return sequencedBatches, nil
}

func (etherMan *Client) trustedVerifyBatchesEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	log.Debug("TrustedVerifyBatches event detected")
	vb, err := etherMan.PoE.ParseTrustedVerifyBatches(vLog)
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
	fsb, err := etherMan.PoE.ParseSequenceForceBatches(vLog)
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
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), big.NewInt(0))
	if err != nil {
		return err
	}
	fullBlock, err := etherMan.EthClient.BlockByHash(ctx, vLog.BlockHash)
	if err != nil {
		return fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
	}
	sequencedForceBatch, err := decodeSequencedForceBatches(tx.Data(), fsb.NumBatch, msg.From(), vLog.TxHash, fullBlock, msg.Nonce())
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
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
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

	var forceBatches []proofofefficiency.ProofOfEfficiencyForcedBatchData
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
			BatchNumber:                      bn,
			Coinbase:                         sequencer,
			TxHash:                           txHash,
			Timestamp:                        time.Unix(int64(block.Time()), 0),
			Nonce:                            nonce,
			ProofOfEfficiencyForcedBatchData: force,
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
	return etherMan.PoE.LastTimestamp(&bind.CallOpts{Pending: false})
}

// GetLatestBatchNumber function allows to retrieve the latest proposed batch in the smc
func (etherMan *Client) GetLatestBatchNumber() (uint64, error) {
	return etherMan.PoE.LastBatchSequenced(&bind.CallOpts{Pending: false})
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
	return etherMan.PoE.LastVerifiedBatch(&bind.CallOpts{Pending: false})
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
func (etherMan *Client) ApproveMatic(ctx context.Context, maticAmount *big.Int, to common.Address) (*types.Transaction, error) {
	if etherMan.IsReadOnly() {
		return nil, ErrIsReadOnlyMode
	}
	opts := *etherMan.auth
	if etherMan.GasProviders.MultiGasProvider {
		opts.GasPrice = etherMan.getGasPrice(ctx)
	}
	tx, err := etherMan.Matic.Approve(&opts, etherMan.SCAddresses[0], maticAmount)
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
	return etherMan.PoE.TrustedSequencerURL(&bind.CallOpts{Pending: false})
}

// GetPublicAddress returns eth client public address
func (etherMan *Client) GetPublicAddress() (common.Address, error) {
	if etherMan.IsReadOnly() {
		return common.Address{}, ErrIsReadOnlyMode
	}
	return etherMan.auth.From, nil
}

// GetL2ChainID returns L2 Chain ID
func (etherMan *Client) GetL2ChainID() (uint64, error) {
	return etherMan.PoE.ChainID(&bind.CallOpts{Pending: false})
}

func (etherMan *Client) getGasPrice(ctx context.Context) *big.Int {
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
	log.Debug("gasPrice choosed: ", gasPrice)
	return gasPrice
}

// SendTx sends a tx to L1
func (etherMan *Client) SendTx(ctx context.Context, tx *types.Transaction) error {
	return etherMan.EthClient.SendTransaction(ctx, tx)
}

// CurrentNonce returns the current nonce for the account signing the L1 txs
func (etherMan *Client) CurrentNonce(ctx context.Context) (uint64, error) {
	return etherMan.EthClient.NonceAt(ctx, etherMan.auth.From, nil)
}

// SuggestedGasPrice returns the suggest nonce for the network at the moment
func (etherMan *Client) SuggestedGasPrice(ctx context.Context) (*big.Int, error) {
	return etherMan.EthClient.SuggestGasPrice(ctx)
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
func (etherMan *Client) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := etherMan.auth.Signer(etherMan.auth.From, tx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// GetRevertMessage tries to get a revert message of a transaction
func (etherMan *Client) GetRevertMessage(ctx context.Context, tx types.Transaction) (string, error) {
	receipt, err := etherMan.GetTxReceipt(ctx, tx.Hash())
	if err != nil {
		return "", err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		revertMessage, err := operations.RevertReason(ctx, etherMan.EthClient, &tx, receipt.BlockNumber)
		if err != nil {
			return "", err
		}
		return revertMessage, nil
	}
	return "", nil
}
