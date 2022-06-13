package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/google/uuid"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/runtime/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation/js"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation/tracers"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/holiman/uint256"
	"github.com/umbracle/ethgo/abi"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
	half                       uint64 = 2
	zkEVMReservedMemorySize    int    = 128
)

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
	// ErrStateNotSynchronized indicates the state database may be empty
	ErrStateNotSynchronized = errors.New("state not synchronized")
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")
	// ErrNilDBTransaction indicates the db transaction has not been properly initialized
	ErrNilDBTransaction = errors.New("database transaction not properly initialized")
	// ErrAlreadyInitializedDBTransaction indicates the db transaction was already initialized
	ErrAlreadyInitializedDBTransaction = errors.New("database transaction already initialized")
	// ErrNotEnoughIntrinsicGas indicates the gas is not enough to cover the intrinsic gas cost
	ErrNotEnoughIntrinsicGas = fmt.Errorf("not enough gas supplied for intrinsic gas costs")
	// ErrParsingExecutorTrace indicates an error occurred while parsing the executor trace
	ErrParsingExecutorTrace = fmt.Errorf("error while parsing executor trace")
)

// State is a implementation of the state
type State struct {
	cfg  Config
	tree statetree
	*PostgresStorage

	mu    *sync.Mutex
	dbTxs map[string]bool
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, tree statetree) *State {
	return &State{
		cfg:             cfg,
		tree:            tree,
		PostgresStorage: storage,

		mu:    new(sync.Mutex),
		dbTxs: make(map[string]bool),
	}
}

// GetTree returns State inner tree
func (s *State) GetTree() statetree {
	return s.tree
}

// BeginStateTransaction starts a transaction block
func (s *State) BeginStateTransaction(ctx context.Context) (string, error) {
	const maxAttempts = 3
	var (
		txBundleID string
		found      bool
	)

	s.mu.Lock()
	for i := 0; i < maxAttempts; i++ {
		txBundleID = uuid.NewString()
		_, idExists := s.dbTxs[txBundleID]
		if !idExists {
			found = true
			break
		}
	}
	s.mu.Unlock()

	if !found {
		return "", fmt.Errorf("could not find unused uuid for db tx bundle")
	}

	if err := s.PostgresStorage.BeginDBTransaction(ctx, txBundleID); err != nil {
		return "", err
	}
	if err := s.tree.BeginDBTransaction(ctx, txBundleID); err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.dbTxs[txBundleID] = true

	return txBundleID, nil
}

// CommitState commits a state into db
func (s *State) CommitState(ctx context.Context, txBundleID string) error {
	if err := s.tree.Commit(ctx, txBundleID); err != nil {
		return err
	}

	if err := s.PostgresStorage.Commit(ctx, txBundleID); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dbTxs, txBundleID)
	return nil
}

// RollbackState rollbacks a db state transaction
func (s *State) RollbackState(ctx context.Context, txBundleID string) error {
	if err := s.tree.Rollback(ctx, txBundleID); err != nil {
		return err
	}

	if err := s.PostgresStorage.Rollback(ctx, txBundleID); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dbTxs, txBundleID)
	return nil
}

// NewBatchProcessor creates a new batch processor
func (s *State) NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (*BatchProcessor, error) {
	// Get Sequencer's Chain ID
	chainID := s.cfg.DefaultChainID
	sq, err := s.GetSequencer(ctx, sequencerAddress, txBundleID)
	if err == nil {
		chainID = sq.ChainID.Uint64()
	}

	logs := make(map[common.Hash][]*types.Log)
	host := Host{State: s, stateRoot: stateRoot, txBundleID: txBundleID, logs: logs}
	host.setRuntime(evm.NewEVM())
	blockNumber, err := s.GetLastBlockNumber(ctx, txBundleID)
	if err != nil {
		return nil, err
	}
	host.forks = runtime.AllForksEnabled.At(blockNumber)

	batchProcessor := &BatchProcessor{SequencerAddress: sequencerAddress, SequencerChainID: chainID, MaxCumulativeGasUsed: s.cfg.MaxCumulativeGasUsed, Host: host}
	batchProcessor.Host.setRuntime(evm.NewEVM())

	return batchProcessor, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *State) NewGenesisBatchProcessor(genesisStateRoot []byte, txBundleID string) (*BatchProcessor, error) {
	host := Host{State: s, stateRoot: genesisStateRoot, txBundleID: txBundleID}
	host.setRuntime(evm.NewEVM())
	host.forks = runtime.AllForksEnabled.At(0)
	return &BatchProcessor{Host: host}, nil
}

// GetStateRoot returns the root of the state tree
func (s *State) GetStateRoot(ctx context.Context, virtual bool, txBundleID string) ([]byte, error) {
	batch, err := s.GetLastBatch(ctx, virtual, txBundleID)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetStateRootByBatchNumber returns state root by batch number from the MT
func (s *State) GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64, txBundleID string) ([]byte, error) {
	batch, err := s.GetBatchByNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(ctx, address, root, txBundleID)
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) ([]byte, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetCode(ctx, address, root, txBundleID)
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error) {
	var lowEnd uint64
	var highEnd uint64
	ctx := context.Background()
	sequencerAddress := common.Address{}

	lastBatch, err := s.GetLastBatch(ctx, true, "")
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return 0, err
	}

	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, lastBatch.Header.Root[:], "")
	if err != nil {
		log.Errorf("failed to get create a new batch processor, err: %v", err)
		return 0, err
	}
	bp.SetSimulationMode(true)

	if bp.isContractCreation(transaction) {
		lowEnd = TxSmartContractCreationGas
	} else {
		lowEnd = TxTransferGas
	}

	if transaction.Gas() != 0 && transaction.Gas() > lowEnd {
		highEnd = transaction.Gas()
	} else {
		highEnd = s.cfg.MaxCumulativeGasUsed
	}

	var availableBalance *big.Int

	if senderAddress != ZeroAddress {
		senderBalance, err := bp.Host.State.tree.GetBalance(ctx, senderAddress, bp.Host.stateRoot, "")
		if err != nil {
			if err == ErrNotFound {
				senderBalance = big.NewInt(0)
			} else {
				return 0, err
			}
		}

		availableBalance = new(big.Int).Set(senderBalance)

		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, ErrInsufficientFunds
			}

			availableBalance.Sub(availableBalance, transaction.Value())
		}
	}

	if transaction.GasPrice().BitLen() != 0 && // Gas price has been set
		availableBalance != nil && // Available balance is found
		availableBalance.Cmp(big.NewInt(0)) > 0 { // Available balance > 0
		gasAllowance := new(big.Int).Div(availableBalance, transaction.GasPrice())

		// Check the gas allowance for this account, make sure high end is capped to it
		if gasAllowance.IsUint64() && highEnd > gasAllowance.Uint64() {
			log.Debugf("Gas estimation high-end capped by allowance [%d]", gasAllowance.Uint64())
			highEnd = gasAllowance.Uint64()
		}
	}

	// Checks if executor level valid gas errors occurred
	isGasApplyError := func(err error) bool {
		return errors.As(err, &ErrNotEnoughIntrinsicGas)
	}

	// Checks if EVM level valid gas errors occurred
	isGasEVMError := func(err error) bool {
		return errors.Is(err, runtime.ErrOutOfGas) ||
			errors.Is(err, runtime.ErrCodeStoreOutOfGas)
	}

	// Checks if the EVM reverted during execution
	isEVMRevertError := func(err error) bool {
		return errors.Is(err, runtime.ErrExecutionReverted)
	}

	// Run the transaction with the specified gas value.
	// Returns a status indicating if the transaction failed and the accompanying error
	testTransaction := func(gas uint64, shouldOmitErr bool) (bool, error) {
		var testResult *runtime.ExecutionResult
		receiverAddress := transaction.To()

		txBundleID, err := s.BeginStateTransaction(ctx)
		if err != nil {
			log.Errorf("estimate gas: failed to begin db transaction, err: %v", err)
			return false, err
		}

		testBp, err := s.NewBatchProcessor(ctx, sequencerAddress, lastBatch.Header.Root[:], txBundleID)
		if err != nil {
			log.Errorf("failed to get create a new batch processor, err: %v", err)
			return false, err
		}
		testBp.SetSimulationMode(true)

		testBp.Host.transactionContext.currentTransaction = transaction
		testBp.Host.transactionContext.currentOrigin = senderAddress
		testBp.Host.transactionContext.coinBase = sequencerAddress

		if testBp.isContractCreation(transaction) {
			testResult = testBp.create(ctx, transaction, senderAddress, sequencerAddress, gas)
		} else if testBp.isSmartContractExecution(ctx, transaction) {
			testResult = testBp.execute(ctx, transaction, senderAddress, *receiverAddress, sequencerAddress, gas, transaction.ChainId())
		} else {
			testResult = testBp.transfer(ctx, transaction, senderAddress, *receiverAddress, sequencerAddress, gas)
		}

		err = s.RollbackState(ctx, txBundleID)
		if err != nil {
			log.Errorf("estimate gas: failed to rollback transaction, err: %v", err)
			return false, err
		}

		// Check if an out of gas error happened during EVM execution
		if testResult.Failed() {
			if (isGasEVMError(testResult.Err) || isGasApplyError(testResult.Err)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, nil
			}

			if isEVMRevertError(testResult.Err) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				return true, constructErrorFromRevert(testResult)
			}

			return true, testResult.Err
		}

		return false, nil
	}

	// Start the binary search for the lowest possible gas price
	for lowEnd < highEnd {
		mid := (lowEnd + highEnd) / half

		failed, testErr := testTransaction(mid, true)
		if testErr != nil &&
			!isEVMRevertError(testErr) {
			// Reverts are ignored in the binary search, but are checked later on
			// during the execution for the optimal gas limit found
			return 0, testErr
		}

		if failed {
			// If the transaction failed => increase the gas
			lowEnd = mid + 1
		} else {
			// If the transaction didn't fail => make this ok value the high end
			highEnd = mid
		}
	}

	// Check if the highEnd is a good value to make the transaction pass
	failed, err := testTransaction(highEnd, false)
	if failed {
		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, fmt.Errorf(
			"unable to apply transaction even for the highest gas limit %d: %w",
			highEnd,
			err,
		)
	}
	return highEnd, nil
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, genesis Genesis, txBundleID string) error {
	// Generate Genesis Block
	block := &Block{
		BlockNumber: genesis.Block.NumberU64(),
		BlockHash:   genesis.Block.Hash(),
		ParentHash:  genesis.Block.ParentHash(),
		ReceivedAt:  genesis.Block.ReceivedAt,
	}

	// Add Block
	err := s.PostgresStorage.AddBlock(ctx, block, txBundleID)
	if err != nil {
		return err
	}

	var (
		root    common.Hash
		newRoot []byte
	)

	if genesis.Balances != nil {
		for address, balance := range genesis.Balances {
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
	}

	if genesis.SmartContracts != nil {
		for address, sc := range genesis.SmartContracts {
			newRoot, _, err = s.tree.SetCode(ctx, address, sc, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
	}

	if len(genesis.Storage) > 0 {
		for address, storage := range genesis.Storage {
			for key, value := range storage {
				newRoot, _, err = s.tree.SetStorageAt(ctx, address, key, value, newRoot, txBundleID)
				if err != nil {
					return err
				}
			}
		}
		root.SetBytes(newRoot)
	}

	if genesis.Nonces != nil {
		for address, nonce := range genesis.Nonces {
			newRoot, _, err = s.tree.SetNonce(ctx, address, nonce, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
	}

	// Generate Genesis Batch
	receivedAt := genesis.Block.ReceivedAt
	batch := &Batch{
		Header: &types.Header{
			Number: big.NewInt(0),
		},
		BlockNumber:        genesis.Block.NumberU64(),
		ConsolidatedTxHash: common.HexToHash("0x1"),
		ConsolidatedAt:     &receivedAt,
		MaticCollateral:    big.NewInt(0),
		ReceivedAt:         time.Now(),
		ChainID:            new(big.Int).SetUint64(genesis.L2ChainID),
		GlobalExitRoot:     common.Hash{},
	}

	// Store batch into db
	bp, err := s.NewGenesisBatchProcessor(root[:], txBundleID)
	if err != nil {
		return err
	}
	err = bp.ProcessBatch(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

// GetNonce returns the nonce of the given account at the given batch number
func (s *State) GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (uint64, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return 0, err
	}

	n, err := s.tree.GetNonce(ctx, address, root, txBundleID)
	if errors.Is(err, tree.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, txBundleID string) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetStorageAt(ctx, address, position, root, txBundleID)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.AddBlock(ctx, block, txBundleID)
}

// ConsolidateBatch consolidates a batch for the given DB tx bundle.
func (s *State) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address, txBundleID string) error {
	return s.PostgresStorage.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash, consolidatedAt, aggregator, txBundleID)
}

// ResetDB resets the state to block for the given DB tx bundle.
func (s *State) ResetDB(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.Reset(ctx, block, txBundleID)
}

func constructErrorFromRevert(result *runtime.ExecutionResult) error {
	revertErrMsg, unpackErr := abi.UnpackRevertError(result.ReturnValue)
	if unpackErr != nil {
		return result.Err
	}

	return fmt.Errorf("%w: %s", result.Err, revertErrMsg)
}

func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error) {
	txBundleID, err := s.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("debug transaction: failed to begin db transaction, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	tx, err := s.GetTransactionByHash(ctx, transactionHash, txBundleID)
	if err != nil {
		log.Errorf("debug transaction: failed to get transaction by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	receipt, err := s.GetTransactionReceipt(ctx, transactionHash, txBundleID)
	if err != nil {
		log.Errorf("debug transaction: failed to get receipt by tx hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	batch, err := s.GetBatchByHash(ctx, receipt.BlockHash, txBundleID)
	if err != nil {
		log.Errorf("debug transaction: failed to get batch by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	var stateRoot []byte

	if receipt.TransactionIndex > 0 {
		previousTX, err := s.GetTransactionByBatchHashAndIndex(ctx, receipt.BlockHash, uint64(receipt.TransactionIndex-1), txBundleID)
		if err != nil {
			log.Errorf("debug transaction: failed to get previous tx, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return nil, err
		}

		previousReceipt, err := s.GetTransactionReceipt(ctx, previousTX.Hash(), txBundleID)
		if err != nil {
			log.Errorf("debug transaction: failed to get receipt by previous tx hash, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return nil, err
		}

		stateRoot = previousReceipt.PostState
	} else {
		previousBatch, err := s.GetBatchByHash(ctx, batch.Header.ParentHash, txBundleID)
		if err == ErrNotFound {
			previousBatch, err = s.GetLastBatch(ctx, true, txBundleID)
			if err != nil {
				log.Errorf("debug transaction: failed to get last batch, err: %v", err)
				rbErr := s.RollbackState(ctx, txBundleID)
				if rbErr != nil {
					log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
				}
				return nil, err
			}
		} else if err != nil {
			log.Errorf("debug transaction: failed to get batch by hash, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return nil, err
		}

		stateRoot = previousBatch.Header.Root.Bytes()
	}

	sequencerAddress := batch.Header.Coinbase

	log.Debugf("debug root: %v", common.Bytes2Hex(stateRoot))

	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, stateRoot, txBundleID)
	if err != nil {
		log.Errorf("debug transaction: failed to create a new batch processor, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("debug transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	// Activate EVM Instrumentation
	bp.Host.runtimes = []runtime.Runtime{}
	evmRT := evm.NewEVM()
	evmRT.EnableInstrumentation()
	bp.Host.setRuntime(evmRT)
	bp.SetSimulationMode(true)

	startTime := time.Now()
	result := bp.processTransaction(ctx, tx, receipt.From, sequencerAddress, tx.ChainId())
	endTime := time.Now()

	// Rollback
	err = s.RollbackState(ctx, txBundleID)
	if err != nil {
		log.Errorf("debug transaction: failed to rollback transaction, err: %v", err)
		return nil, err
	}

	if tracer == "" {
		return result, nil
	}

	// Parse the executor-like trace using the FakeEVM
	jsTracer, err := js.NewJsTracer(tracer, new(tracers.Context))
	if err != nil {
		log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
		return nil, err
	}

	context := instrumentation.Context{}

	// Fill trace context
	if tx.To() == nil {
		context.Type = "CREATE"
	} else {
		context.Type = "CALL"
		context.To = tx.To().Hex()
	}
	context.From = receipt.From.Hex()
	context.Input = "0x" + hex.EncodeToString(tx.Data())
	context.Gas = strconv.FormatUint(tx.Gas(), encoding.Base10)
	context.Value = tx.Value().String()
	context.Output = "0x" + hex.EncodeToString(result.ReturnValue)
	context.GasPrice = tx.GasPrice().String()
	context.ChainID = tx.ChainId().Uint64()
	context.OldStateRoot = "0x" + hex.EncodeToString(stateRoot)
	context.Time = uint64(endTime.Sub(startTime))
	context.GasUsed = strconv.FormatUint(result.GasUsed, encoding.Base10)

	result.ExecutorTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return result, nil
	}

	env := fakevm.NewFakeEVM(vm.BlockContext{BlockNumber: big.NewInt(1)}, vm.TxContext{GasPrice: gasPrice}, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: jsTracer})
	fakeDB := &FakeDB{State: s, stateRoot: common.Hex2Bytes(context.OldStateRoot)}
	env.SetStateDB(fakeDB)

	traceResult, err := s.ParseTheTraceUsingTheTracer(env, result.ExecutorTrace, jsTracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return result, nil
	}

	result.ExecutorTraceResult = traceResult

	return result, nil
}

func (s *State) ParseTheTraceUsingTheTracer(env *fakevm.FakeEVM, trace instrumentation.ExecutorTrace, jsTracer tracers.Tracer) (json.RawMessage, error) {
	var previousDepth int
	var stateRoot []byte

	contextGas, ok := new(big.Int).SetString(trace.Context.Gas, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing contextGas")
		return nil, ErrParsingExecutorTrace
	}
	value, ok := new(big.Int).SetString(trace.Context.Value, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing value")
		return nil, ErrParsingExecutorTrace
	}

	jsTracer.CaptureTxStart(contextGas.Uint64())
	jsTracer.CaptureStart(env, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", common.Hex2Bytes(strings.TrimLeft(trace.Context.Input, "0x")), contextGas.Uint64(), value)

	stack := fakevm.Newstack()
	memory := fakevm.NewMemory()

	stateRoot = common.Hex2Bytes(trace.Context.OldStateRoot)

	for _, step := range trace.Steps {
		gas, ok := new(big.Int).SetString(step.Gas, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step gas")
			return nil, ErrParsingExecutorTrace
		}

		gasCost, ok := new(big.Int).SetString(step.GasCost, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step gasCost")
			return nil, ErrParsingExecutorTrace
		}

		value, ok := new(big.Int).SetString(step.Contract.Value, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step value")
			return nil, ErrParsingExecutorTrace
		}

		op, ok := new(big.Int).SetString(step.Op, 0)
		if !ok {
			log.Debugf("error while parsing step op")
			return nil, ErrParsingExecutorTrace
		}

		scope := &fakevm.ScopeContext{
			Contract: vm.NewContract(fakevm.NewAccount(common.HexToAddress(step.Contract.Caller)), fakevm.NewAccount(common.HexToAddress(step.Contract.Address)), value, gas.Uint64()),
			Memory:   memory,
			Stack:    stack,
		}

		opcode := vm.OpCode(op.Uint64()).String()

		if opcode == "CREATE" || opcode == "CREATE2" || opcode == "CALL" || opcode == "CALLCODE" || opcode == "DELEGATECALL" || opcode == "STATICCALL" || opcode == "SELFDESTRUCT" {
			jsTracer.CaptureEnter(vm.OpCode(op.Uint64()), common.HexToAddress(step.Contract.Caller), common.HexToAddress(step.Contract.Address), common.Hex2Bytes(strings.TrimLeft(step.Contract.Input, "0x")), gas.Uint64(), value)
		}

		if step.Error != "" {
			err := fmt.Errorf(step.Error)
			jsTracer.CaptureFault(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.Depth, err)
		} else {
			jsTracer.CaptureState(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, common.Hex2Bytes(strings.TrimLeft(step.ReturnData, "0x")), step.Depth, nil)
		}

		// Set Memory
		if len(step.Memory) > 0 {
			memory.Resize(uint64(fakevm.MemoryItemSize*len(step.Memory) + zkEVMReservedMemorySize))
			for offset, memoryContent := range step.Memory {
				memory.Set(uint64((offset*fakevm.MemoryItemSize)+zkEVMReservedMemorySize), uint64(fakevm.MemoryItemSize), common.Hex2Bytes(memoryContent))
			}
		}

		// Set Stack
		stack = fakevm.Newstack()
		for _, stackContent := range step.Stack {
			// log.Debugf(stackContent)
			valueBigInt, ok := new(big.Int).SetString(stackContent, 0)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			value, _ := uint256.FromBig(valueBigInt)
			stack.Push(value)
		}

		// Set Storage
		for storageKey, storageValue := range step.Storage {
			key, ok := new(big.Int).SetString("0x"+storageKey, 0)
			if !ok {
				log.Debugf("error while parsing storage key")
				return nil, ErrParsingExecutorTrace
			}
			value, ok := new(big.Int).SetString("0x"+storageValue, 0)
			if !ok {
				log.Debugf("error while parsing storage value")
				return nil, ErrParsingExecutorTrace
			}
			_, _, err := s.tree.SetStorageAt(context.Background(), common.HexToAddress(step.Contract.Address), key, value, stateRoot, "")
			if err != nil {
				return nil, err
			}
		}

		// Returning from a call or create
		if previousDepth > step.Depth {
			jsTracer.CaptureExit([]byte{}, gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		// Set StateRoot
		stateRoot = common.Hex2Bytes(step.StateRoot)
		env.StateDB.SetStateRoot(stateRoot)
		previousDepth = step.Depth
	}

	gasUsed, ok := new(big.Int).SetString(trace.Context.GasUsed, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing gasUsed")
		return nil, ErrParsingExecutorTrace
	}

	jsTracer.CaptureTxEnd(gasUsed.Uint64())
	jsTracer.CaptureEnd([]byte(trace.Context.Output), gasUsed.Uint64(), time.Duration(trace.Context.Time), nil)

	return jsTracer.GetResult()
}
