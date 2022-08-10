package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
	// Size of the memory in bytes reserved by the zkEVM
	zkEVMReservedMemorySize int  = 128
	two                     uint = 2
	cTrue                        = 1
	cFalse                       = 0
)

var (
	// ZeroHash is the hash 0x0000000000000000000000000000000000000000000000000000000000000000
	ZeroHash = common.Hash{}
	// ZeroAddress is the address 0x0000000000000000000000000000000000000000
	ZeroAddress = common.Address{}
)

// State is a implementation of the state
type State struct {
	cfg Config
	*PostgresStorage
	executorClient pb.ExecutorServiceClient
	tree           *merkletree.StateTree
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, executorClient pb.ExecutorServiceClient, stateTree *merkletree.StateTree) *State {
	return &State{
		cfg:             cfg,
		PostgresStorage: storage,
		executorClient:  executorClient,
		tree:            stateTree,
	}
}

// BeginStateTransaction starts a state transaction
func (s *State) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(ctx, address, l2Block.Root().Bytes())
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetCode(ctx, address, l2Block.Root().Bytes())
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return 0, err
	}

	nonce, err := s.tree.GetNonce(ctx, address, l2Block.Root().Bytes())
	if err != nil {
		return 0, err
	}
	return nonce.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetStorageAt(ctx, address, position, l2Block.Root().Bytes())
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error) {
	var lowEnd uint64
	var highEnd uint64
	ctx := context.Background()

	lastBatch, err := s.GetLastBatch(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return 0, err
	}

	stateRoot := lastBatch.StateRoot

	if s.isContractCreation(transaction) {
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
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, stateRoot.Bytes())
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
		batchL2Data, err := EncodeTransactions([]types.Transaction{*transaction})
		if err != nil {
			return false, err
		}

		// Create a batch to be sent to the executor
		processBatchRequest := &pb.ProcessBatchRequest{
			BatchNum:         lastBatch.BatchNumber + 1,
			Coinbase:         senderAddress.String(),
			BatchL2Data:      batchL2Data,
			OldStateRoot:     stateRoot.Bytes(),
			UpdateMerkleTree: cFalse,
		}

		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		if err != nil {
			return false, err
		}

		// Check if an out of gas error happened during EVM execution
		if processBatchResponse.Responses[0].Error != pb.Error(executor.NO_ERROR) {
			error := fmt.Errorf(executor.ExecutorError(processBatchResponse.Responses[0].Error).Error())

			if (isGasEVMError(error) || isGasApplyError(error)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, nil
			}

			if isEVMRevertError(error) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				return true, constructErrorFromRevert(error, processBatchResponse.Responses[0].ReturnValue)
			}

			return true, error
		}

		return false, nil
	}

	// Start the binary search for the lowest possible gas price
	for lowEnd < highEnd {
		mid := (lowEnd + highEnd) / uint64(two)

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

// OpenBatch adds a new batch into the state, with the necessary data to start processing transactions within it.
// It's meant to be used by sequencers, since they don't necessarely know what transactions are going to be added
// in this batch yet. In other words it's the creation of a WIP batch.
// Note that this will add a batch with batch number N + 1, where N it's the greates batch number on the state.
func (s *State) OpenBatch(ctx context.Context, processingContext ProcessingContext, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}
	// Check if the batch that is being opened has batch num + 1 compared to the latest batch
	lastBatchNum, err := s.PostgresStorage.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum+1 != processingContext.BatchNumber {
		return fmt.Errorf("unexpected batch number %v, should be %v", processingContext.BatchNumber, lastBatchNum+1)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if !isLastBatchClosed {
		return ErrLastBatchShouldBeClosed
	}
	// Check that timestamp is equal or greater compared to previous batch
	prevTimestamp, err := s.GetLastBatchTime(ctx, dbTx)
	if err != nil {
		return err
	}
	if prevTimestamp.Unix() > processingContext.Timestamp.Unix() {
		return ErrTimestampGE
	}
	return s.PostgresStorage.openBatch(ctx, processingContext, dbTx)
}

// ProcessSequencerBatch is used by the sequencers to proceess transactions into an open batch
func (s *State) ProcessSequencerBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	batchL2Data, err := EncodeTransactions(txs)
	if err != nil {
		return nil, err
	}
	processBatchResponse, err := s.processBatch(ctx, batchNumber, batchL2Data, dbTx)
	if err != nil {
		return nil, err
	}
	result, err := convertToProcessBatchResponse(txs, processBatchResponse)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *State) processBatch(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) (*pb.ProcessBatchResponse, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}
	lastBatches, err := s.PostgresStorage.GetLastNBatches(ctx, two, dbTx)
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get GER and Timestamp
	lastBatch := lastBatches[0]
	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[1]
	isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if isBatchClosed {
		return nil, ErrBatchAlreadyClosed
	}

	// Check provided batch number is the latest in db
	if lastBatch.BatchNumber != batchNumber {
		return nil, ErrInvalidBatchNumber
	}
	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		BatchNum:         lastBatch.BatchNumber,
		Coinbase:         lastBatch.Coinbase.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldLocalExitRoot: previousBatch.LocalExitRoot.Bytes(),
		EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
		UpdateMerkleTree: cTrue,
	}

	// Send Batch to the Executor
	return s.executorClient.ProcessBatch(ctx, processBatchRequest)
}

// StoreTransactions is used by the sequencer to add processed transactions into
// an open batch. If the batch already has txs, the processedTxs must be a super
// set of the existing ones, preserving order.
func (s *State) StoreTransactions(ctx context.Context, batchNumber uint64, processedTxs []*ProcessTransactionResponse, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// check existing txs vs parameter txs
	existingTxs, err := s.GetTxsHashesByBatchNumber(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if err := CheckSupersetBatchTransactions(existingTxs, processedTxs); err != nil {
		return err
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if isBatchClosed {
		return ErrBatchAlreadyClosed
	}

	processingContext, err := s.GetProcessingContext(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}

	firstTxToInsert := len(existingTxs)

	for i := firstTxToInsert; i < len(processedTxs); i++ {
		processedTx := processedTxs[i]

		lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return err
		}

		header := &types.Header{
			Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
			ParentHash: lastL2Block.Hash(),
			Coinbase:   processingContext.Coinbase,
			Root:       processedTx.StateRoot,
		}
		transactions := []*types.Transaction{&processedTx.Tx}

		// Create block to be able to calculate its hash
		block := types.NewBlock(header, transactions, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
		block.ReceivedAt = processingContext.Timestamp

		receipts := []*types.Receipt{generateReceipt(block, processedTx)}

		// Store L2 block and its transaction
		if err := s.PostgresStorage.AddL2Block(ctx, batchNumber, block, receipts, dbTx); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) isBatchClosable(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	// Check if the batch that is being closed is the last batch
	lastBatchNum, err := s.PostgresStorage.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum != receipt.BatchNumber {
		return fmt.Errorf("unexpected batch number %v, should be %v", receipt.BatchNumber, lastBatchNum)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if isLastBatchClosed {
		return ErrBatchAlreadyClosed
	}

	return nil
}

// closeSynchronizedBatch is used by Synchronizer to close the current batch.
func (s *State) closeSynchronizedBatch(ctx context.Context, receipt ProcessingReceipt, txs []types.Transaction, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	err := s.isBatchClosable(ctx, receipt, dbTx)
	if err != nil {
		return err
	}

	if len(txs) == 0 {
		return ErrClosingBatchWithoutTxs
	}

	batchL2Data, err := EncodeTransactions(txs)
	if err != nil {
		return err
	}

	return s.PostgresStorage.closeBatch(ctx, receipt, batchL2Data, dbTx)
}

// CloseBatch is used by sequencer to close the current batch. It will set the processing receipt and
// the raw txs data based on the txs included on that batch that are already in the state
func (s *State) CloseBatch(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	// TODO: differentiate the case where sequencer / sync calls the function so it's possible
	// to use L2BatchData from L1 rather than from stored txs
	if dbTx == nil {
		return ErrDBTxNil
	}

	err := s.isBatchClosable(ctx, receipt, dbTx)
	if err != nil {
		return err
	}

	// Generate raw txs data
	encodedTxsArray, err := s.GetEncodedTransactionsByBatchNumber(ctx, receipt.BatchNumber, dbTx)
	if err != nil {
		return err
	}
	if len(encodedTxsArray) == 0 {
		return ErrClosingBatchWithoutTxs
	}
	txs := []types.Transaction{}
	for i := 0; i < len(encodedTxsArray); i++ {
		tx, err := DecodeTx(encodedTxsArray[i])
		if err != nil {
			return err
		}
		txs = append(txs, *tx)
	}
	batchL2Data, err := EncodeTransactions(txs)
	if err != nil {
		return err
	}

	return s.PostgresStorage.closeBatch(ctx, receipt, batchL2Data, dbTx)
}

// ProcessAndStoreClosedBatch is used by the Synchronizer to add a closed batch into the data base
func (s *State) ProcessAndStoreClosedBatch(ctx context.Context, processingCtx ProcessingContext, encodedTxs []byte, dbTx pgx.Tx) error {
	// Decode transactions
	decodedTransactions, _, err := DecodeTxs(encodedTxs)
	if err != nil {
		log.Debugf("error decoding transactions: %w", err)
		return err
	}

	// Open the batch and process the txs
	if dbTx == nil {
		return ErrDBTxNil
	}
	if err := s.OpenBatch(ctx, processingCtx, dbTx); err != nil {
		return err
	}
	processed, err := s.processBatch(ctx, processingCtx.BatchNumber, encodedTxs, dbTx)
	if err != nil {
		return err
	}

	// Sanity check
	if len(decodedTransactions) != len(processed.Responses) {
		return fmt.Errorf("number of decoded (%d) and processed (%d) transactions do not match", len(decodedTransactions), len(processed.Responses))
	}

	// Filter unprocessed txs and decode txs to store metadata
	// note that if the batch is not well encoded it will result in an empty batch (with no txs)
	for i := 0; i < len(processed.Responses); i++ {
		//TODO: Also check this
		if !isProcessed(processed.Responses[i].Error) {
			// Remove unprocessed tx
			if i == len(processed.Responses)-1 {
				processed.Responses = processed.Responses[:i]
				decodedTransactions = decodedTransactions[:i]
			} else {
				processed.Responses = append(processed.Responses[:i], processed.Responses[i+1:]...)
				decodedTransactions = append(decodedTransactions[:i], decodedTransactions[i+1:]...)
			}
			i--
		}
	}

	processedBatch, err := convertToProcessBatchResponse(decodedTransactions, processed)
	if err != nil {
		return err
	}
	// Store processed txs into the batch
	err = s.StoreTransactions(ctx, processingCtx.BatchNumber, processedBatch.Responses, dbTx)
	if err != nil {
		return err
	}

	// Close batch
	return s.closeSynchronizedBatch(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
	}, decodedTransactions, dbTx)
}

// GetLastBatch gets latest batch (closed or not) on the data base
func (s *State) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error) {
	batches, err := s.PostgresStorage.GetLastNBatches(ctx, 1, dbTx)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, ErrNotFound
	}
	return batches[0], nil
}

// DebugTransaction re-executes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error) {
	// TODO: Implement
	return new(runtime.ExecutionResult), nil
}

// ParseTheTraceUsingTheTracer parses the given trace with the given tracer.
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

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, blockNumber uint64, dbTx pgx.Tx) *runtime.ExecutionResult {
	result := new(runtime.ExecutionResult)

	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByBlockNumber(ctx, blockNumber, two, dbTx)
	if err != nil {
		result.Err = err
		return result
	}

	// Get latest batch from the database to get GER and Timestamp
	lastBatch := lastBatches[0]
	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[1]

	batchL2Data, err := EncodeUnsignedTransaction(*tx)
	if err != nil {
		log.Errorf("error encoding unsigned transaction ", err)
		result.Err = err
		return result
	}

	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		BatchL2Data:      batchL2Data,
		From:             senderAddress.String(),
		OldStateRoot:     l2BlockStateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldLocalExitRoot: previousBatch.LocalExitRoot.Bytes(),
		EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
		UpdateMerkleTree: cFalse,
	}

	// Send Batch to the Executor
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Errorf("error processing unsigned transaction ", err)
		result.Err = err
		return result
	}
	response, err := convertToProcessBatchResponse([]types.Transaction{*tx}, processBatchResponse)
	if err != nil {
		result.Err = err
		return result
	}
	// Todo populate result
	r := response.Responses[0]
	result.ReturnValue = r.ReturnValue
	result.GasLeft = r.GasLeft
	result.GasUsed = r.GasUsed
	result.Err = fmt.Errorf(r.Error)
	result.CreateAddress = r.CreateAddress
	result.StateRoot = r.StateRoot.Bytes()

	return result
}

// GetTree returns State inner tree
func (s *State) GetTree() *merkletree.StateTree {
	return s.tree
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, block Block, genesis Genesis, dbTx pgx.Tx) ([]byte, error) {
	var (
		root    common.Hash
		newRoot []byte
		err     error
	)
	if dbTx == nil {
		return newRoot, ErrDBTxNil
	}

	for _, action := range genesis.Actions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, ok := new(big.Int).SetString(action.Value, encoding.Base10)
			if !ok {
				return newRoot, fmt.Errorf("Could not set base10 balance %q to big.Int", action.Value)
			}
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, ok := new(big.Int).SetString(action.Value, encoding.Base10)
			if !ok {
				return newRoot, fmt.Errorf("Could not set base10 nonce %q to big.Int", action.Value)
			}
			newRoot, _, err = s.tree.SetNonce(ctx, address, nonce, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return newRoot, fmt.Errorf("Could not decode SC bytecode for address %q: %v", address, err)
			}
			newRoot, _, err = s.tree.SetCode(ctx, address, code, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeStorage):
			if strings.HasPrefix(action.StoragePosition, "0x") { // nolint
				action.StoragePosition = action.StoragePosition[2:]
			}
			positionBI := new(big.Int).SetBytes(common.Hex2Bytes(action.StoragePosition))
			if strings.HasPrefix(action.Value, "0x") { // nolint
				action.StoragePosition = action.Value[2:]
			}
			valueBI := new(big.Int).SetBytes(common.Hex2Bytes(action.Value))

			newRoot, _, err = s.tree.SetStorageAt(ctx, address, positionBI, valueBI, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeSCLength):
			log.Debug("Skipped genesis action of type merkletree.LeafTypeSCLength, these actions will be handled as part of merkletree.LeafTypeCode actions")
		default:
			return newRoot, fmt.Errorf("Unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(newRoot)

	receivedAt := time.Unix(0, 0)

	// store L1 block related to genesis batch
	err = s.AddBlock(ctx, &block, dbTx)
	if err != nil {
		return newRoot, err
	}

	// store genesis batch
	batch := Batch{
		BatchNumber:    0,
		Coinbase:       ZeroAddress,
		BatchL2Data:    nil,
		StateRoot:      root,
		LocalExitRoot:  ZeroHash,
		Timestamp:      receivedAt,
		Transactions:   []types.Transaction{},
		GlobalExitRoot: ZeroHash,
	}

	err = s.storeGenesisBatch(ctx, batch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// mark the genesis batch as virtualized
	virtualBatch := &VirtualBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Coinbase:    ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVirtualBatch(ctx, virtualBatch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// mark the genesis batch as verified/consolidated
	verifiedBatch := &VerifiedBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Aggregator:  ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVerifiedBatch(ctx, verifiedBatch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// store L2 genesis block
	header := &types.Header{
		Number:     big.NewInt(0),
		ParentHash: ZeroHash,
		Coinbase:   ZeroAddress,
		Root:       root,
	}
	rootHex := root.Hex()
	log.Info("Genesis root ", rootHex)
	l2Block := types.NewBlock(header, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	l2Block.ReceivedAt = receivedAt

	return newRoot, s.PostgresStorage.AddL2Block(ctx, batch.BatchNumber, l2Block, []*types.Receipt{}, dbTx)
}

// CheckSupersetBatchTransactions verifies that processedTransactions is a
// superset of existingTxs and that the existing txs have the same order,
// returns a non-nil error if that is not the case.
func CheckSupersetBatchTransactions(existingTxHashes []common.Hash, processedTxs []*ProcessTransactionResponse) error {
	if len(existingTxHashes) > len(processedTxs) {
		return ErrExistingTxGreaterThanProcessedTx
	}
	for i, existingTxHash := range existingTxHashes {
		if existingTxHash != processedTxs[i].Tx.Hash() {
			return ErrOutOfOrderProcessedTx
		}
	}
	return nil
}

// isContractCreation checks if the tx is a contract creation
func (s *State) isContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// DetermineProcessedTransactions splits the given tx process responses
// returning a slice with only processed and a map unprocessed txs
// respectively.
func DetermineProcessedTransactions(responses []*ProcessTransactionResponse) ([]*ProcessTransactionResponse, map[string]*ProcessTransactionResponse) {
	processedTxResponses := []*ProcessTransactionResponse{}
	unprocessedTxResponses := map[string]*ProcessTransactionResponse{}
	for _, response := range responses {
		if response.IsProcessed {
			processedTxResponses = append(processedTxResponses, response)
		} else {
			log.Infof("Tx %s has not been processed", response.TxHash)
			unprocessedTxResponses[response.TxHash.String()] = response
		}
	}
	return processedTxResponses, unprocessedTxResponses
}
