package state

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestConvertToProcessBatchResponse for test purposes
func TestConvertToProcessBatchResponse(txs []types.Transaction, response *pb.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	return convertToProcessBatchResponse(txs, response)
}

func convertToProcessBatchResponse(txs []types.Transaction, response *pb.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	responses, err := convertToProcessTransactionResponse(txs, response.Responses)
	if err != nil {
		return nil, err
	}

	isBatchProcessed := true
	if len(response.Responses) > 0 {
		// Check out of counters
		isBatchProcessed = !(response.Responses[len(response.Responses)-1].Error == pb.Error_ERROR_OUT_OF_COUNTERS)
	}

	return &ProcessBatchResponse{
		CumulativeGasUsed:   response.CumulativeGasUsed,
		IsBatchProcessed:    isBatchProcessed,
		Responses:           responses,
		NewStateRoot:        common.BytesToHash(response.NewStateRoot),
		NewLocalExitRoot:    common.BytesToHash(response.NewLocalExitRoot),
		CntKeccakHashes:     response.CntKeccakHashes,
		CntPoseidonHashes:   response.CntPoseidonHashes,
		CntPoseidonPaddings: response.CntPoseidonPaddings,
		CntMemAligns:        response.CntMemAligns,
		CntArithmetics:      response.CntArithmetics,
		CntBinaries:         response.CntBinaries,
		CntSteps:            response.CntSteps,
	}, nil
}

func isProcessed(err pb.Error) bool {
	return err != pb.Error_ERROR_INTRINSIC_INVALID_TX && err != pb.Error_ERROR_OUT_OF_COUNTERS
}

func isOOC(err pb.Error) bool {
	return err == pb.Error_ERROR_OUT_OF_COUNTERS
}

func convertToProcessTransactionResponse(txs []types.Transaction, responses []*pb.ProcessTransactionResponse) ([]*ProcessTransactionResponse, error) {
	results := make([]*ProcessTransactionResponse, 0, len(responses))
	for i, response := range responses {
		trace, err := convertToStructLogArray(response.ExecutionTrace)
		if err != nil {
			return nil, err
		}

		result := new(ProcessTransactionResponse)
		result.TxHash = common.BytesToHash(response.TxHash)
		result.Type = response.Type
		result.ReturnValue = response.ReturnValue
		result.GasLeft = response.GasLeft
		result.GasUsed = response.GasUsed
		result.GasRefunded = response.GasRefunded
		result.Error = executor.Err(response.Error)
		result.CreateAddress = common.HexToAddress(response.CreateAddress)
		result.StateRoot = common.BytesToHash(response.StateRoot)
		result.Logs = convertToLog(response.Logs)
		result.IsProcessed = isProcessed(response.Error)
		result.ExecutionTrace = *trace
		result.CallTrace = convertToExecutorTrace(response.CallTrace)
		result.Tx = txs[i]
		results = append(results, result)

		log.Debugf("ProcessTransactionResponse[TxHash]: %v", txs[i].Hash().String())
		log.Debugf("ProcessTransactionResponse[Nonce]: %v", txs[i].Nonce())
		log.Debugf("ProcessTransactionResponse[StateRoot]: %v", result.StateRoot.String())
		log.Debugf("ProcessTransactionResponse[Error]: %v", result.Error)
		log.Debugf("ProcessTransactionResponse[GasUsed]: %v", result.GasUsed)
		log.Debugf("ProcessTransactionResponse[GasLeft]: %v", result.GasLeft)
		log.Debugf("ProcessTransactionResponse[GasRefunded]: %v", result.GasRefunded)
		log.Debugf("ProcessTransactionResponse[IsProcessed]: %v", result.IsProcessed)
	}

	return results, nil
}

func convertToLog(protoLogs []*pb.Log) []*types.Log {
	logs := make([]*types.Log, 0, len(protoLogs))

	for _, protoLog := range protoLogs {
		log := new(types.Log)
		log.Address = common.HexToAddress(protoLog.Address)
		log.Topics = convertToTopics(protoLog.Topics)
		log.Data = protoLog.Data
		log.BlockNumber = protoLog.BatchNumber
		log.TxHash = common.BytesToHash(protoLog.TxHash)
		log.TxIndex = uint(protoLog.TxIndex)
		log.BlockHash = common.BytesToHash(protoLog.BatchHash)
		log.Index = uint(protoLog.Index)
		logs = append(logs, log)
	}

	return logs
}

func convertToTopics(responses [][]byte) []common.Hash {
	results := make([]common.Hash, 0, len(responses))

	for _, response := range responses {
		results = append(results, common.BytesToHash(response))
	}
	return results
}

func convertToStructLogArray(responses []*pb.ExecutionTraceStep) (*[]instrumentation.StructLog, error) {
	results := make([]instrumentation.StructLog, 0, len(responses))

	for _, response := range responses {
		convertedStack, err := convertToBigIntArray(response.Stack)
		if err != nil {
			return nil, err
		}
		result := new(instrumentation.StructLog)
		result.Pc = response.Pc
		result.Op = response.Op
		result.Gas = response.RemainingGas
		result.GasCost = response.GasCost
		result.Memory = response.Memory
		result.MemorySize = int(response.MemorySize)
		result.Stack = convertedStack
		result.ReturnData = response.ReturnData
		result.Storage = convertToProperMap(response.Storage)
		result.Depth = int(response.Depth)
		result.RefundCounter = response.GasRefund
		result.Err = executor.Err(response.Error)

		results = append(results, *result)
	}
	return &results, nil
}

func convertToBigIntArray(responses []string) ([]*big.Int, error) {
	results := make([]*big.Int, 0, len(responses))

	for _, response := range responses {
		result, ok := new(big.Int).SetString(response, hex.Base)
		if ok {
			results = append(results, result)
		} else {
			return nil, fmt.Errorf("String %s is not valid", response)
		}
	}
	return results, nil
}

func convertToProperMap(responses map[string]string) map[common.Hash]common.Hash {
	results := make(map[common.Hash]common.Hash, len(responses))
	for key, response := range responses {
		results[common.HexToHash(key)] = common.HexToHash(response)
	}
	return results
}

func convertToExecutorTrace(callTrace *pb.CallTrace) instrumentation.ExecutorTrace {
	trace := new(instrumentation.ExecutorTrace)
	if callTrace != nil {
		trace.Context = convertToContext(callTrace.Context)
		trace.Steps = convertToInstrumentationSteps(callTrace.Steps)
	}

	return *trace
}

func convertToContext(context *pb.TransactionContext) instrumentation.Context {
	return instrumentation.Context{
		Type:         context.Type,
		From:         context.From,
		To:           context.To,
		Input:        string(context.Data),
		Gas:          fmt.Sprint(context.Gas),
		Value:        context.Value,
		Output:       string(context.Output),
		GasPrice:     context.GasPrice,
		OldStateRoot: string(context.OldStateRoot),
		Time:         uint64(context.ExecutionTime),
		GasUsed:      fmt.Sprint(context.GasUsed),
	}
}

func convertToInstrumentationSteps(responses []*pb.TransactionStep) []instrumentation.Step {
	results := make([]instrumentation.Step, 0, len(responses))
	for _, response := range responses {
		step := new(instrumentation.Step)
		step.StateRoot = string(response.StateRoot)
		step.Depth = int(response.Depth)
		step.Pc = response.Pc
		step.Gas = fmt.Sprint(response.Gas)
		step.OpCode = fakevm.OpCode(response.Op).String()
		step.Refund = fmt.Sprint(response.GasRefund)
		step.Op = fmt.Sprint(response.Op)
		err := executor.Err(response.Error)
		if err != nil {
			step.Error = err.Error()
		}
		step.Contract = convertToInstrumentationContract(response.Contract)
		step.GasCost = fmt.Sprint(response.GasCost)
		step.Stack = response.Stack
		step.Memory = convertByteArrayToStringArray(response.Memory)
		step.ReturnData = string(response.ReturnData)

		results = append(results, *step)
	}
	return results
}

func convertToInstrumentationContract(response *pb.Contract) instrumentation.Contract {
	return instrumentation.Contract{
		Address: response.Address,
		Caller:  response.Caller,
		Value:   response.Value,
		Input:   string(response.Data),
		Gas:     fmt.Sprint(response.Gas),
	}
}

func convertByteArrayToStringArray(responses []byte) []string {
	results := make([]string, 0, len(responses))
	for _, response := range responses {
		results = append(results, string(response))
	}
	return results
}
