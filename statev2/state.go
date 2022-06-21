package statev2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation/tracers"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor/pb"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
	// Size of the memory in bytes reserved by the zkEVM
	zkEVMReservedMemorySize int = 128
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
	cfg Config
	*PostgresStorage
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage) *State {
	return &State{
		cfg:             cfg,
		PostgresStorage: storage,
	}
}

// BeginDBTransaction starts a state transaction
func (s *State) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Commit commits a state transaction
func (s *State) CommitStateTransaction(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	return err
}

// Rollback rollbacks a state transaction
func (s *State) RollbackStateTransaction(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	return err
}

// ResetDB resets the state to block for the given DB tx .
func (s *State) ResetDB(ctx context.Context, block *Block, tx pgx.Tx) error {
	return s.PostgresStorage.Reset(ctx, block, tx)
}

func (s *State) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, tx pgx.Tx) error {
	return s.PostgresStorage.AddGlobalExitRoot(ctx, exitRoot, tx)
}

func (s *State) GetLatestGlobalExitRoot(ctx context.Context, tx pgx.Tx) (*GlobalExitRoot, error) {
	return s.PostgresStorage.GetLatestGlobalExitRoot(ctx, tx)
}

func (s *State) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, tx pgx.Tx) error {
	return s.PostgresStorage.AddForcedBatch(ctx, forcedBatch, tx)
}

func (s *State) GetForcedBatch(ctx context.Context, tx pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	return s.PostgresStorage.GetForcedBatch(ctx, tx, forcedBatchNumber)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, tx pgx.Tx) error {
	return s.PostgresStorage.AddBlock(ctx, block, tx)
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, blockNumber uint64, txBundleID string) (*big.Int, error) {
	return nil, nil
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, blockNumber uint64, txBundleID string) ([]byte, error) {
	return nil, nil
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error) {
	return 0, nil
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, txBundleID string) (uint64, error) {
	return 0, nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, txBundleID string) (*big.Int, error) {
	return new(big.Int), nil
}

// StoreBatchHeader is used by the Trusted Sequencer to create a new batch
func (s *State) StoreBatchHeader(ctx context.Context, batch Batch) error {
	return nil
}

// ProcessBatch is used by the Trusted Sequencer to add transactions to the batch
func (s *State) ProcessBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction) (*ProcessBatchResponse, error) {
	// check batchNumber is the latest in db
	return nil, nil
}

// StoreTransactions is used by the Trusted Sequencer to add processed transactions into the data base
func (s *State) StoreTransactions(batchNum uint64, processedTxs []*ProcessTransactionResponse) error {
	return nil
}

// ProcessAndStoreWIPBatch is used by the Synchronizer to add a work-in-progress batch into the data base
func (s *State) ProcessAndStoreWIPBatch(ctx context.Context, batch Batch) error {
	return nil
}

// ProcessAndStoreClosedBatch is used by the Synchronizer to a add closed batch into the data base
func (s *State) ProcessAndStoreClosedBatch(ctx context.Context, batch Batch) error {
	processBatchRequest := &pb.ProcessBatchRequest{
		BatchNum: uint32(batch.BatchNum),
		Coinbase: batch.Coinbase.Hex(),
		// BatchL2Data:
		// OldStateRoot:
		GlobalExitRoot:       batch.GlobalExitRootNum.Bytes(),
		OldLocalExitRoot:     batch.OldLocalExitRoot.Bytes(),
		EthTimestamp:         uint64(batch.EthTimestamp.Unix()),
		UpdateMerkleTree:     true,
		GenerateExecuteTrace: false,
		GenerateCallTrace:    false,
	}

	// Create client
	executorClient, _ := executor.NewExecutorClient(s.cfg.ExecutorServerConfig)

	_, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		return err
	}

	// Store Batch into data base

	return nil
}

// SetGenesis
// CondolidateBatch
/*
GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
GetLastBatchNumberSeenOnEthereum(ctx context.Context, txBundleID string) (uint64, error)
GetLastBatchByStateRoot(ctx context.Context, stateRoot []byte, txBundleID string) (*state.Batch, error)
SetGenesis(ctx context.Context, genesis state.Genesis, txBundleID string) error
SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error
SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error
SetInitSyncBatch(ctx context.Context, batchNumber uint64, txBundleID string) error
AddBlock(ctx context.Context, block *state.Block, txBundleID string) error
CreateBatch(ctx context.Context, batch *statev2.Batch) error
ProcessBatch(ctx context.Context, txs []types.Transaction) ProcessBatchResponse
AddTransactionsToBatch(ctx context.Context, batchNumber uint64, txs []ProcessTransactionResponse) error
*/

func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error) {
	return new(runtime.ExecutionResult), nil
}

func (s *State) ParseTheTraceUsingTheTracer(env *fakevm.FakeEVM, trace instrumentation.ExecutorTrace, jsTracer tracers.Tracer) (json.RawMessage, error) {
	var previousDepth int
	var previousOpcode string
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

	bigStateRoot, ok := new(big.Int).SetString(trace.Context.OldStateRoot, 0)
	if !ok {
		log.Debugf("error while parsing context oldStateRoot")
		return nil, ErrParsingExecutorTrace
	}
	stateRoot = bigStateRoot.Bytes()
	env.StateDB.SetStateRoot(stateRoot)

	for i, step := range trace.Steps {
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

		codeAddr := common.HexToAddress(step.Contract.Address)
		scope.Contract.CodeAddr = &codeAddr

		opcode := vm.OpCode(op.Uint64()).String()

		if previousOpcode == "CALL" && step.Pc != 0 {
			jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		if opcode != "CALL" || trace.Steps[i+1].Pc == 0 {
			if step.Error != "" {
				err := fmt.Errorf(step.Error)
				jsTracer.CaptureFault(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.Depth, err)
			} else {
				jsTracer.CaptureState(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, common.Hex2Bytes(strings.TrimLeft(step.ReturnData, "0x")), step.Depth, nil)
			}
		}

		if opcode == "CREATE" || opcode == "CREATE2" || opcode == "CALL" || opcode == "CALLCODE" || opcode == "DELEGATECALL" || opcode == "STATICCALL" || opcode == "SELFDESTRUCT" {
			jsTracer.CaptureEnter(vm.OpCode(op.Uint64()), common.HexToAddress(step.Contract.Caller), common.HexToAddress(step.Contract.Address), common.Hex2Bytes(strings.TrimLeft(step.Contract.Input, "0x")), gas.Uint64(), value)
			if step.OpCode == "SELFDESTRUCT" {
				jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
			}
		}

		// Set Memory
		if len(step.Memory) > 0 {
			memory.Resize(uint64(fakevm.MemoryItemSize*len(step.Memory) + zkEVMReservedMemorySize))
			for offset, memoryContent := range step.Memory {
				memory.Set(uint64((offset*fakevm.MemoryItemSize)+zkEVMReservedMemorySize), uint64(fakevm.MemoryItemSize), common.Hex2Bytes(memoryContent))
			}
		} else {
			memory = fakevm.NewMemory()
		}

		// Set Stack
		stack = fakevm.Newstack()
		for _, stackContent := range step.Stack {
			valueBigInt, ok := new(big.Int).SetString(stackContent, 0)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			value, _ := uint256.FromBig(valueBigInt)
			stack.Push(value)
		}

		// Returning from a call or create
		if previousDepth > step.Depth {
			jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		// Set StateRoot
		bigStateRoot, ok := new(big.Int).SetString(step.StateRoot, 0)
		if !ok {
			log.Debugf("error while parsing step stateRoot")
			return nil, ErrParsingExecutorTrace
		}

		stateRoot = bigStateRoot.Bytes()
		env.StateDB.SetStateRoot(stateRoot)
		previousDepth = step.Depth
		previousOpcode = step.OpCode
	}

	gasUsed, ok := new(big.Int).SetString(trace.Context.GasUsed, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing gasUsed")
		return nil, ErrParsingExecutorTrace
	}

	jsTracer.CaptureTxEnd(gasUsed.Uint64())
	jsTracer.CaptureEnd(common.Hex2Bytes(trace.Context.Output), gasUsed.Uint64(), time.Duration(trace.Context.Time), nil)

	return jsTracer.GetResult()
}

// AddVerifiedBatch adds a new VerifiedBatch to the db
func (s *State) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, tx pgx.Tx) error {
	return s.PostgresStorage.AddVerifiedBatch(ctx, verifiedBatch, tx)
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (s *State) GetVerifiedBatch(ctx context.Context, tx pgx.Tx, batchNumber uint64) (*VerifiedBatch, error) {
	return s.PostgresStorage.GetVerifiedBatch(ctx, tx, batchNumber)
}
