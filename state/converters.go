package state

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func convertToProcessBatchResponse(txs []types.Transaction, response *pb.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	responses, err := convertToProcessTransactionResponse(txs, response.Responses)
	if err != nil {
		return nil, err
	}
	return &ProcessBatchResponse{
		CumulativeGasUsed:   response.CumulativeGasUsed,
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

func isProcessed(error pb.Error) bool {
	//TODO: Implement this check

	return true
}

func convertToProcessTransactionResponse(txs []types.Transaction, responses []*pb.ProcessTransactionResponse) ([]*ProcessTransactionResponse, error) {
	results := make([]*ProcessTransactionResponse, 0, len(responses))

	for i, response := range responses {
		trace, err := convertToStrucLogArray(response.ExecutionTrace)
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
		result.Error = executor.ExecutorError(response.Error).Error()
		result.CreateAddress = common.HexToAddress(response.CreateAddress)
		result.StateRoot = common.BytesToHash(response.StateRoot)
		result.Logs = convertToLog(response.Logs)
		result.IsProcessed = isProcessed(response.Error)
		result.ExecutionTrace = *trace
		result.CallTrace = convertToExecutorTrace(response.CallTrace)
		result.Tx = txs[i]
		results = append(results, result)
	}

	return results, nil
}

func convertToLog(responses []*pb.Log) []*types.Log {
	results := make([]*types.Log, 0, len(responses))

	for _, response := range responses {
		result := new(types.Log)
		result.Address = common.HexToAddress(response.Address)
		result.Topics = convertToTopics(response.Topics)
		result.Data = response.Data
		result.BlockNumber = response.BatchNumber
		result.TxHash = common.BytesToHash(response.TxHash)
		result.TxIndex = uint(response.TxIndex)
		result.BlockHash = common.BytesToHash(response.BatchHash)
		result.Index = uint(response.Index)
		results = append(results, result)
	}

	return results
}

func convertToTopics(responses [][]byte) []common.Hash {
	results := make([]common.Hash, 0, len(responses))

	for _, response := range responses {
		results = append(results, common.BytesToHash(response))
	}
	return results
}

func convertToStrucLogArray(responses []*pb.ExecutionTraceStep) (*[]instrumentation.StructLog, error) {
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
		result.Err = fmt.Errorf(executor.ExecutorError(response.Error).Error())

		results = append(results, *result)
	}
	return &results, nil
}

func convertToBigIntArray(responses []string) ([]*big.Int, error) {
	results := make([]*big.Int, 0, len(responses))

	for _, response := range responses {
		result, ok := new(big.Int).SetString(response, 16)
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
		Value:        fmt.Sprint(context.Value),
		Output:       string(context.Output),
		GasPrice:     fmt.Sprint(context.GasPrice),
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
		step.Error = executor.ExecutorError(response.Error).Error()
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
		Value:   fmt.Sprint(response.Value),
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
