package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/js"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers/native"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

const (
	two uint = 2
)

// GetSender gets the sender from the transaction's signature
func GetSender(tx types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// RlpFieldsToLegacyTx parses the rlp fields slice into a type.LegacyTx
// in this specific order:
//
// required fields:
// [0] Nonce    uint64
// [1] GasPrice *big.Int
// [2] Gas      uint64
// [3] To       *common.Address
// [4] Value    *big.Int
// [5] Data     []byte
//
// optional fields:
// [6] V        *big.Int
// [7] R        *big.Int
// [8] S        *big.Int
func RlpFieldsToLegacyTx(fields [][]byte, v, r, s []byte) (tx *types.LegacyTx, err error) {
	const (
		fieldsSizeWithoutChainID = 6
		fieldsSizeWithChainID    = 7
	)

	if len(fields) < fieldsSizeWithoutChainID {
		return nil, types.ErrTxTypeNotSupported
	}

	nonce := big.NewInt(0).SetBytes(fields[0]).Uint64()
	gasPrice := big.NewInt(0).SetBytes(fields[1])
	gas := big.NewInt(0).SetBytes(fields[2]).Uint64()
	var to *common.Address

	if fields[3] != nil && len(fields[3]) != 0 {
		tmp := common.BytesToAddress(fields[3])
		to = &tmp
	}
	value := big.NewInt(0).SetBytes(fields[4])
	data := fields[5]

	txV := big.NewInt(0).SetBytes(v)
	if len(fields) >= fieldsSizeWithChainID {
		chainID := big.NewInt(0).SetBytes(fields[6])

		// a = chainId * 2
		// b = v - 27
		// c = a + 35
		// v = b + c
		//
		// same as:
		// v = v-27+chainId*2+35
		a := new(big.Int).Mul(chainID, big.NewInt(double))
		b := new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V))
		c := new(big.Int).Add(a, big.NewInt(etherPre155V))
		txV = new(big.Int).Add(b, c)
	}

	txR := big.NewInt(0).SetBytes(r)
	txS := big.NewInt(0).SetBytes(s)

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
		V:        txV,
		R:        txR,
		S:        txS,
	}, nil
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
		// if the transaction has an intrinsic invalid tx error it means
		// the transaction has not changed the state, so we don't store it
		// and just move to the next
		if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
			continue
		}

		lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return err
		}

		header := &types.Header{
			Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
			ParentHash: lastL2Block.Hash(),
			Coinbase:   processingContext.Coinbase,
			Root:       processedTx.StateRoot,
			GasUsed:    processedTx.GasUsed,
			GasLimit:   s.cfg.MaxCumulativeGasUsed,
			Time:       uint64(processingContext.Timestamp.Unix()),
		}
		transactions := []*types.Transaction{&processedTx.Tx}

		receipt := generateReceipt(header.Number, processedTx)
		receipts := []*types.Receipt{receipt}

		// Create block to be able to calculate its hash
		block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
		block.ReceivedAt = processingContext.Timestamp

		receipt.BlockHash = block.Hash()

		// Store L2 block and its transaction
		if err := s.AddL2Block(ctx, batchNumber, block, receipts, dbTx); err != nil {
			return err
		}
	}
	return nil
}

// DebugTransaction re-executes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, traceConfig TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	// gets the transaction
	tx, err := s.GetTransactionByHash(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// gets the tx receipt
	receipt, err := s.GetTransactionReceipt(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// gets the l2 block including the transaction
	block, err := s.GetL2BlockByNumber(ctx, receipt.BlockNumber.Uint64(), dbTx)
	if err != nil {
		return nil, err
	}

	// get the previous L2 Block
	previousBlockNumber := uint64(0)
	if receipt.BlockNumber.Uint64() > 0 {
		previousBlockNumber = receipt.BlockNumber.Uint64() - 1
	}
	previousBlock, err := s.GetL2BlockByNumber(ctx, previousBlockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	// generate batch l2 data for the transaction
	batchL2Data, err := EncodeTransactions([]types.Transaction{*tx})
	if err != nil {
		return nil, err
	}

	// gets batch that including the l2 block
	batch, err := s.GetBatchByL2BlockNumber(ctx, block.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	// gets batch that including the previous l2 block
	previousBatch, err := s.GetBatchByL2BlockNumber(ctx, previousBlock.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	// Create Batch
	traceConfigRequest := &pb.TraceConfig{
		TxHashToGenerateCallTrace:    transactionHash.Bytes(),
		TxHashToGenerateExecuteTrace: transactionHash.Bytes(),
		// set the defaults to the maximum information we can have.
		// this is needed to process custom tracers later
		DisableStorage:   cFalse,
		DisableStack:     cFalse,
		EnableMemory:     cTrue,
		EnableReturnData: cTrue,
	}

	// if the default tracer is used, then we review the information
	// we want to have in the trace related to the parameters we received.
	if traceConfig.IsDefaultTracer() {
		if traceConfig.DisableStorage {
			traceConfigRequest.DisableStorage = cTrue
		}
		if traceConfig.DisableStack {
			traceConfigRequest.DisableStack = cTrue
		}
		if traceConfig.EnableMemory {
			traceConfigRequest.EnableMemory = cTrue
		}
		if traceConfig.EnableReturnData {
			traceConfigRequest.EnableReturnData = cTrue
		}
	}

	oldStateRoot := previousBlock.Root()
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:     batch.BatchNumber - 1,
		OldStateRoot:    oldStateRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),

		BatchL2Data:      batchL2Data,
		GlobalExitRoot:   batch.GlobalExitRoot.Bytes(),
		EthTimestamp:     uint64(batch.Timestamp.Unix()),
		Coinbase:         batch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		TraceConfig:      traceConfigRequest,
	}

	// Send Batch to the Executor
	startTime := time.Now()
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	endTime := time.Now()
	if err != nil {
		return nil, err
	} else if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
		return nil, err
	}

	// //save process batch response file
	// b, err := json.Marshal(processBatchResponse)
	// if err != nil {
	// 	return nil, err
	// }
	// filePath := "./processBatchResponse.json"
	// err = os.WriteFile(filePath, b, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	txs, _, err := DecodeTxs(batchL2Data)
	if err != nil && !errors.Is(err, ErrInvalidData) {
		return nil, err
	}

	for _, tx := range txs {
		log.Debugf(tx.Hash().String())
	}

	convertedResponse, err := s.convertToProcessBatchResponse(txs, processBatchResponse)
	if err != nil {
		return nil, err
	}

	// Sanity check
	response := convertedResponse.Responses[0]
	log.Debugf(response.TxHash.String())
	if response.TxHash != transactionHash {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	result := &runtime.ExecutionResult{
		CreateAddress: response.CreateAddress,
		GasLeft:       response.GasLeft,
		GasUsed:       response.GasUsed,
		ReturnValue:   response.ReturnValue,
		StateRoot:     response.StateRoot.Bytes(),
		StructLogs:    response.ExecutionTrace,
		ExecutorTrace: response.CallTrace,
	}

	// if is the default trace, return the result
	if traceConfig.IsDefaultTracer() {
		return result, nil
	}

	senderAddress, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	context := instrumentation.Context{
		From:         senderAddress.String(),
		Input:        hex.EncodeToHex(tx.Data()),
		Gas:          strconv.FormatUint(tx.Gas(), encoding.Base10),
		Value:        tx.Value().String(),
		Output:       hex.EncodeToHex(result.ReturnValue),
		GasPrice:     tx.GasPrice().String(),
		OldStateRoot: oldStateRoot.String(),
		Time:         uint64(endTime.Sub(startTime)),
		GasUsed:      strconv.FormatUint(result.GasUsed, encoding.Base10),
	}

	// Fill trace context
	if tx.To() == nil {
		context.Type = "CREATE"
		context.To = result.CreateAddress.Hex()
	} else {
		context.Type = "CALL"
		context.To = tx.To().Hex()
	}

	result.ExecutorTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return nil, fmt.Errorf("failed to parse gasPrice")
	}

	tracerContext := &tracers.Context{
		BlockHash:   receipt.BlockHash,
		BlockNumber: receipt.BlockNumber,
		TxIndex:     int(receipt.TransactionIndex),
		TxHash:      transactionHash,
	}

	var evmTracer tracers.Tracer
	if traceConfig.Is4ByteTracer() {
		evmTracer, err = native.NewFourByteTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create 4byteTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create 4byteTracer, err: %v", err)
		}
	} else if traceConfig.IsCallTracer() {
		evmTracer, err = native.NewCallTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create callTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create callTracer, err: %v", err)
		}
	} else if traceConfig.IsNoopTracer() {
		evmTracer, err = native.NewNoopTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create noopTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create noopTracer, err: %v", err)
		}
	} else if traceConfig.IsPrestateTracer() {
		evmTracer, err = native.NewPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create prestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create prestateTracer, err: %v", err)
		}
	} else if traceConfig.IsJSCustomTracer() {
		evmTracer, err = js.NewJsTracer(*traceConfig.Tracer, tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create jsTracer, err: %v", err)
		}
	} else {
		return nil, fmt.Errorf("invalid tracer: %v, err: %v", traceConfig.Tracer, err)
	}

	fakeDB := &FakeDB{State: s, stateRoot: batch.StateRoot.Bytes()}
	evm := fakevm.NewFakeEVM(fakevm.BlockContext{BlockNumber: big.NewInt(1)}, fakevm.TxContext{GasPrice: gasPrice}, fakeDB, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: evmTracer})

	traceResult, err := s.ParseTheTraceUsingTheTracer(evm, result.ExecutorTrace, evmTracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return nil, fmt.Errorf("failed parse the trace using the tracer: %v", err)
	}

	result.ExecutorTraceResult = traceResult

	return result, nil
}

// ParseTheTraceUsingTheTracer parses the given trace with the given tracer.
func (s *State) ParseTheTraceUsingTheTracer(evm *fakevm.FakeEVM, trace instrumentation.ExecutorTrace, tracer tracers.Tracer) (json.RawMessage, error) {
	var previousDepth int
	var previousOp, previousGas *big.Int
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

	tracer.CaptureTxStart(contextGas.Uint64())
	decodedInput, err := hex.DecodeHex(trace.Context.Input)
	if err != nil {
		log.Debugf("error while decoding context input from hex to bytes")
		return nil, fmt.Errorf("%w %w", err, ErrParsingExecutorTrace)
	}
	tracer.CaptureStart(evm, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", decodedInput, contextGas.Uint64(), value)

	bigStateRoot, ok := new(big.Int).SetString(trace.Context.OldStateRoot, 0)
	if !ok {
		log.Debugf("error while parsing context oldStateRoot")
		return nil, ErrParsingExecutorTrace
	}
	stateRoot = bigStateRoot.Bytes()
	evm.StateDB.SetStateRoot(stateRoot)

	output := common.FromHex(trace.Context.Output)

	var stepError error
	for i, step := range trace.Steps {
		stepErrorMsg := strings.TrimSpace(step.Error)
		if stepErrorMsg != "" {
			stepError = fmt.Errorf(stepErrorMsg)
		}

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

		op, ok := new(big.Int).SetString(step.Op, 0)
		if !ok {
			log.Debugf("error while parsing step op")
			return nil, ErrParsingExecutorTrace
		}

		// set Stack
		stack := fakevm.NewStack()
		for _, stackContent := range step.Stack {
			valueBigInt, ok := new(big.Int).SetString(stackContent, hex.Base)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			value, _ := uint256.FromBig(valueBigInt)
			stack.Push(value)
		}

		// set Memory
		memory := fakevm.NewMemory()
		if len(step.Memory) > 0 {
			memory.Resize(uint64(len(step.Memory)))
			memory.Set(0, uint64(len(step.Memory)), step.Memory)
		} else {
			memory = fakevm.NewMemory()
		}

		value := hex.DecodeBig(step.Contract.Value)
		scope := &fakevm.ScopeContext{
			Contract: fakevm.NewContract(fakevm.NewAccount(common.HexToAddress(step.Contract.Caller)), fakevm.NewAccount(common.HexToAddress(step.Contract.Address)), value, gas.Uint64()),
			Memory:   memory,
			Stack:    stack,
		}

		codeAddr := common.HexToAddress(step.Contract.Address)
		scope.Contract.CodeAddr = &codeAddr

		// when a revert is detected, we stop the execution
		if step.OpCode == "REVERT" {
			stepError = fakevm.ErrExecutionReverted
			break
		}

		if previousOpcode == "CALL" && step.Pc != 0 {
			tracer.CaptureExit(step.ReturnData, gasCost.Uint64(), stepError)
		}

		if step.OpCode != "CALL" || trace.Steps[i+1].Pc == 0 {
			if stepError != nil {
				tracer.CaptureFault(step.Pc, fakevm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.Depth, stepError)
			} else {
				tracer.CaptureState(step.Pc, fakevm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.ReturnData, step.Depth, nil)
			}
		}

		previousOpCodeCanBeSubCall := previousOpcode == "CREATE" ||
			previousOpcode == "CREATE2" ||
			previousOpcode == "DELEGATECALL" ||
			previousOpcode == "CALL" ||
			previousOpcode == "STATICCALL" ||
			// deprecated ones
			previousOpcode == "CALLCODE" ||
			previousOpcode == "SELFDESTRUCT"

		// when a sub call or create is detected, the next step contains the contract updated
		if previousOpCodeCanBeSubCall {
			// shadowing "value" to override its value without compromising the external code
			value := value
			// value is not carried over when the capture enter handles STATIC CALL
			if previousOpcode == "STATICCALL" {
				value = nil
			}
			tracer.CaptureEnter(fakevm.OpCode(previousOp.Uint64()), common.HexToAddress(step.Contract.Caller), common.HexToAddress(step.Contract.Address), []byte(step.Contract.Input), previousGas.Uint64(), value)
		}

		// returning from a call or create
		if previousDepth > step.Depth {
			tracer.CaptureExit(step.ReturnData, gasCost.Uint64(), stepError)
		}

		// set StateRoot
		stateRoot = []byte(step.StateRoot)
		evm.StateDB.SetStateRoot(stateRoot)

		// set previous step values
		previousDepth = step.Depth
		previousOp = op
		previousGas = gas
		previousOpcode = step.OpCode
	}

	gasUsed, ok := new(big.Int).SetString(trace.Context.GasUsed, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing gasUsed")
		return nil, ErrParsingExecutorTrace
	}

	restGas := contextGas.Uint64() - gasUsed.Uint64()
	tracer.CaptureTxEnd(restGas)
	tracer.CaptureEnd(output, gasUsed.Uint64(), stepError)

	return tracer.GetResult()
}

// PreProcessTransaction processes the transaction in order to calculate its zkCounters before adding it to the pool
func (s *State) PreProcessTransaction(ctx context.Context, tx *types.Transaction, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	sender, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	response, err := s.internalProcessUnsignedTransaction(ctx, tx, sender, nil, false, dbTx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	result := new(runtime.ExecutionResult)
	response, err := s.internalProcessUnsignedTransaction(ctx, tx, senderAddress, l2BlockNumber, noZKEVMCounters, dbTx)
	if err != nil {
		return nil, err
	}

	r := response.Responses[0]
	result.ReturnValue = r.ReturnValue
	result.GasLeft = r.GasLeft
	result.GasUsed = r.GasUsed
	result.CreateAddress = r.CreateAddress
	result.StateRoot = r.StateRoot.Bytes()

	if errors.Is(r.RomError, runtime.ErrExecutionReverted) {
		result.Err = constructErrorFromRevert(r.RomError, r.ReturnValue)
	} else {
		result.Err = r.RomError
	}

	return result, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) internalProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, two, dbTx)
	if err != nil {
		return nil, err
	}

	stateRoot := l2BlockStateRoot
	if l2BlockNumber != nil {
		l2Block, err := s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
		if err != nil {
			return nil, err
		}
		stateRoot = l2Block.Root()
	}

	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
	if err != nil {
		return nil, err
	}
	nonce := loadedNonce.Uint64()

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	timestamp := uint64(lastBatch.Timestamp.Unix())

	if l2BlockNumber != nil {
		latestL2BlockNumber, err := s.PostgresStorage.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		if *l2BlockNumber == latestL2BlockNumber {
			timestamp = uint64(time.Now().Unix())
		}
	}

	batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, &nonce)
	if err != nil {
		log.Errorf("error encoding unsigned transaction ", err)
		return nil, err
	}

	forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, lastBatch.BatchNumber)
	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber,
		BatchL2Data:      batchL2Data,
		From:             senderAddress.String(),
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     timestamp,
		Coinbase:         lastBatch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
	}

	if noZKEVMCounters {
		processBatchRequest.NoCounters = cTrue
	}

	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)

	// Send Batch to the Executor
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		// Log this error as an executor unspecified error
		s.eventLog.LogExecutorError(ctx, pb.ExecutorError_EXECUTOR_ERROR_UNSPECIFIED, processBatchRequest)
		log.Errorf("error processing unsigned transaction ", err)
		return nil, err
	} else if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
		return nil, err
	}

	response, err := s.convertToProcessBatchResponse([]types.Transaction{*tx}, processBatchResponse)
	if err != nil {
		return nil, err
	}

	if processBatchResponse.Responses[0].Error != pb.RomError(executor.ROM_ERROR_NO_ERROR) {
		err := executor.RomErr(processBatchResponse.Responses[0].Error)
		if !isEVMRevertError(err) {
			return response, err
		}
	}

	return response, nil
}

// isContractCreation checks if the tx is a contract creation
func (s *State) isContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// DetermineProcessedTransactions splits the given tx process responses
// returning a slice with only processed and a map unprocessed txs
// respectively.
func DetermineProcessedTransactions(responses []*ProcessTransactionResponse) (
	[]*ProcessTransactionResponse, []string, map[string]*ProcessTransactionResponse, []string) {
	processedTxResponses := []*ProcessTransactionResponse{}
	processedTxsHashes := []string{}
	unprocessedTxResponses := map[string]*ProcessTransactionResponse{}
	unprocessedTxsHashes := []string{}
	for _, response := range responses {
		if response.IsProcessed {
			processedTxResponses = append(processedTxResponses, response)
			processedTxsHashes = append(processedTxsHashes, response.TxHash.String())
		} else {
			log.Infof("Tx %s has not been processed", response.TxHash)
			unprocessedTxResponses[response.TxHash.String()] = response
			unprocessedTxsHashes = append(unprocessedTxsHashes, response.TxHash.String())
		}
	}
	return processedTxResponses, processedTxsHashes, unprocessedTxResponses, unprocessedTxsHashes
}

// StoreTransaction is used by the sequencer to add process a transaction
func (s *State) StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	/*
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
	*/
	// if the transaction has an intrinsic invalid tx error it means
	// the transaction has not changed the state, so we don't store it
	if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
		return nil
	}

	lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
	if err != nil {
		return err
	}

	header := &types.Header{
		Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
		ParentHash: lastL2Block.Hash(),
		Coinbase:   coinbase,
		Root:       processedTx.StateRoot,
		GasUsed:    processedTx.GasUsed,
		GasLimit:   s.cfg.MaxCumulativeGasUsed,
		Time:       timestamp,
	}
	transactions := []*types.Transaction{&processedTx.Tx}

	receipt := generateReceipt(header.Number, processedTx)
	receipts := []*types.Receipt{receipt}

	// Create block to be able to calculate its hash
	block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
	block.ReceivedAt = time.Unix(int64(timestamp), 0)

	receipt.BlockHash = block.Hash()

	// Store L2 block and its transaction
	if err := s.AddL2Block(ctx, batchNumber, block, receipts, dbTx); err != nil {
		return err
	}

	return nil
}

// CheckSupersetBatchTransactions verifies that processedTransactions is a
// superset of existingTxs and that the existing txs have the same order,
// returns a non-nil error if that is not the case.
func CheckSupersetBatchTransactions(existingTxHashes []common.Hash, processedTxs []*ProcessTransactionResponse) error {
	if len(existingTxHashes) > len(processedTxs) {
		return ErrExistingTxGreaterThanProcessedTx
	}
	for i, existingTxHash := range existingTxHashes {
		if existingTxHash != processedTxs[i].TxHash {
			return ErrOutOfOrderProcessedTx
		}
	}
	return nil
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, dbTx pgx.Tx) (uint64, []byte, error) {
	const ethTransferGas = 21000

	var lowEnd uint64
	var highEnd uint64

	ctx := context.Background()

	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, two, dbTx)
	if err != nil {
		return 0, nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	lowEnd, err = core.IntrinsicGas(transaction.Data(), transaction.AccessList(), s.isContractCreation(transaction), true, false, false)
	if err != nil {
		return 0, nil, err
	}

	if lowEnd == ethTransferGas && transaction.To() != nil {
		code, err := s.tree.GetCode(ctx, *transaction.To(), l2BlockStateRoot.Bytes())
		if err != nil {
			log.Warnf("error while getting transaction.to() code %v", err)
		} else if len(code) == 0 {
			return lowEnd, nil, nil
		}
	}

	if transaction.Gas() != 0 && transaction.Gas() > lowEnd {
		highEnd = transaction.Gas()
	} else {
		highEnd = s.cfg.MaxCumulativeGasUsed
	}

	var availableBalance *big.Int

	if senderAddress != ZeroAddress {
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, l2BlockStateRoot.Bytes())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				senderBalance = big.NewInt(0)
			} else {
				return 0, nil, err
			}
		}

		availableBalance = new(big.Int).Set(senderBalance)

		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, nil, ErrInsufficientFunds
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

	// Run the transaction with the specified gas value.
	// Returns a status indicating if the transaction failed, if it was reverted and the accompanying error
	testTransaction := func(gas uint64, shouldOmitErr bool) (failed, reverted bool, gasUsed uint64, returnValue []byte, err error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    transaction.Nonce(),
			To:       transaction.To(),
			Value:    transaction.Value(),
			Gas:      gas,
			GasPrice: transaction.GasPrice(),
			Data:     transaction.Data(),
		})

		batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, nil)
		if err != nil {
			log.Errorf("error encoding unsigned transaction ", err)
			return false, false, gasUsed, nil, err
		}

		forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, lastBatch.BatchNumber)
		// Create a batch to be sent to the executor
		processBatchRequest := &pb.ProcessBatchRequest{
			OldBatchNum:      lastBatch.BatchNumber,
			BatchL2Data:      batchL2Data,
			From:             senderAddress.String(),
			OldStateRoot:     l2BlockStateRoot.Bytes(),
			GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
			OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
			EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
			Coinbase:         lastBatch.Coinbase.String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkID,
		}

		log.Debugf("EstimateGas[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		// log.Debugf("EstimateGas[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("EstimateGas[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("EstimateGas[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("EstimateGas[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("EstimateGas[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("EstimateGas[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("EstimateGas[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("EstimateGas[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("EstimateGas[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("EstimateGas[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)

		txExecutionOnExecutorTime := time.Now()
		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		log.Debugf("executor time: %vms", time.Since(txExecutionOnExecutorTime).Milliseconds())
		if err != nil {
			log.Errorf("error estimating gas: %v", err)
			return false, false, gasUsed, nil, err
		}
		gasUsed = processBatchResponse.Responses[0].GasUsed
		if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponse.Error)
			s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
			return false, false, gasUsed, nil, err
		}

		// Check if an out of gas error happened during EVM execution
		if processBatchResponse.Responses[0].Error != pb.RomError(executor.ROM_ERROR_NO_ERROR) {
			err := executor.RomErr(processBatchResponse.Responses[0].Error)

			if (isGasEVMError(err) || isGasApplyError(err)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, false, gasUsed, nil, nil
			}

			if isEVMRevertError(err) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				returnValue := processBatchResponse.Responses[0].ReturnValue
				return true, true, gasUsed, returnValue, constructErrorFromRevert(err, returnValue)
			}

			return true, false, gasUsed, nil, err
		}

		return false, false, gasUsed, nil, nil
	}

	txExecutions := []time.Duration{}
	var totalExecutionTime time.Duration

	// Check if the highEnd is a good value to make the transaction pass
	failed, reverted, gasUsed, returnValue, err := testTransaction(highEnd, false)
	log.Debugf("Estimate gas. Trying to execute TX with %v gas", highEnd)
	if failed {
		if reverted {
			return 0, returnValue, err
		}

		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, nil, fmt.Errorf(
			"unable to apply transaction even for the highest gas limit %d: %w",
			highEnd,
			err,
		)
	}

	if lowEnd < gasUsed {
		lowEnd = gasUsed
	}

	// Start the binary search for the lowest possible gas price
	for (lowEnd < highEnd) && (highEnd-lowEnd) > 4096 {
		txExecutionStart := time.Now()
		mid := (lowEnd + highEnd) / uint64(two)

		log.Debugf("Estimate gas. Trying to execute TX with %v gas", mid)

		failed, reverted, _, _, testErr := testTransaction(mid, true)
		executionTime := time.Since(txExecutionStart)
		totalExecutionTime += executionTime
		txExecutions = append(txExecutions, executionTime)
		if testErr != nil && !reverted {
			// Reverts are ignored in the binary search, but are checked later on
			// during the execution for the optimal gas limit found
			return 0, nil, testErr
		}

		if failed {
			// If the transaction failed => increase the gas
			lowEnd = mid + 1
		} else {
			// If the transaction didn't fail => make this ok value the high end
			highEnd = mid
		}
	}

	executions := int64(len(txExecutions))
	if executions > 0 {
		log.Debugf("EstimateGas executed the TX %v times", executions)
		averageExecutionTime := totalExecutionTime.Milliseconds() / executions
		log.Debugf("EstimateGas tx execution average time is %v milliseconds", averageExecutionTime)
	} else {
		log.Error("Estimate gas. Tx not executed")
	}
	return highEnd, nil, nil
}

// Checks if executor level valid gas errors occurred
func isGasApplyError(err error) bool {
	return errors.Is(err, ErrNotEnoughIntrinsicGas)
}

// Checks if EVM level valid gas errors occurred
func isGasEVMError(err error) bool {
	return errors.Is(err, runtime.ErrOutOfGas)
}

// Checks if the EVM reverted during execution
func isEVMRevertError(err error) bool {
	return errors.Is(err, runtime.ErrExecutionReverted)
}
