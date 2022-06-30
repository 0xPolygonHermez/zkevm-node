package statev2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor/pb"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/fakevm"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/instrumentation"
)

const ether155V = 27

func EncodeTransactions(txs []types.Transaction) ([]byte, error) {
	var batchL2Data []byte

	// TODO: Check how to encode unsigned transactions

	for _, tx := range txs {
		v, r, s := tx.RawSignatureValues()
		sign := 1 - (v.Uint64() & 1)

		txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
			tx.Nonce(),
			tx.GasPrice(),
			tx.Gas(),
			tx.To(),
			tx.Value(),
			tx.Data(),
			tx.ChainId(), uint(0), uint(0),
		})

		if err != nil {
			return nil, err
		}

		newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
		newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
		newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
		newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
		txData, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
		if err != nil {
			return nil, err
		}

		batchL2Data = append(batchL2Data, txData...)
	}

	return batchL2Data, nil
}

func convertToProcessBatchResponse(response *pb.ProcessBatchResponse) *ProcessBatchResponse {
	return &ProcessBatchResponse{
		CumulativeGasUsed:   response.CumulativeGasUsed,
		Responses:           convertToProcessTransactionResponse(response.Responses),
		NewStateRoot:        common.BytesToHash(response.NewStateRoot),
		NewLocalExitRoot:    common.BytesToHash(response.NewLocalExitRoot),
		CntKeccakHashes:     response.CntKeccakHashes,
		CntPoseidonHashes:   response.CntPoseidonHashes,
		CntPoseidonPaddings: response.CntPoseidonPaddings,
		CntMemAligns:        response.CntMemAligns,
		CntArithmetics:      response.CntArithmetics,
		CntBinaries:         response.CntBinaries,
		CntSteps:            response.CntSteps,
	}
}

func convertToProcessTransactionResponse(responses []*pb.ProcessTransactionResponse) []*ProcessTransactionResponse {
	results := make([]*ProcessTransactionResponse, 0, len(responses))

	for _, response := range responses {
		result := new(ProcessTransactionResponse)
		result.TxHash = common.BytesToHash(response.TxHash)
		result.Type = response.Type
		result.ReturnValue = response.ReturnValue
		result.GasLeft = response.GasLeft
		result.GasUsed = response.GasUsed
		result.GasRefunded = response.GasRefunded
		result.Error = response.Error
		result.CreateAddress = common.HexToAddress(response.CreateAddress)
		result.StateRoot = common.BytesToHash(response.StateRoot)
		result.Logs = convertToLog(response.Logs)
		result.UnprocessedTransaction = response.UnprocessedTransaction
		result.ExecutionTrace = convertToStrucLogArray(response.ExecutionTrace)
		result.CallTrace = convertToExecutorTrace(response.CallTrace)
		results = append(results, result)
	}

	return results
}

func convertToLog(responses []*pb.Log) []types.Log {
	results := make([]types.Log, 0, len(responses))

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
		results = append(results, *result)
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

func convertToStrucLogArray(responses []*pb.ExecutionTraceStep) []instrumentation.StructLog {
	results := make([]instrumentation.StructLog, 0, len(responses))
	for _, response := range responses {
		result := new(instrumentation.StructLog)
		result.Pc = response.Pc
		result.Op = response.Op
		result.Gas = response.RemainingGas
		result.GasCost = response.GasCost
		result.Memory = response.Memory
		result.MemorySize = int(response.MemorySize)
		result.Stack = convertToBigIntArray(response.Stack)
		result.ReturnData = response.ReturnData
		result.Storage = convertToProperMap(response.Storage)
		result.Depth = int(response.Depth)
		result.RefundCounter = response.GasRefund
		result.Err = fmt.Errorf(response.Error)

		results = append(results, *result)
	}
	return results
}

func convertToBigIntArray(responses []uint64) []*big.Int {
	results := make([]*big.Int, 0, len(responses))

	for _, response := range responses {
		results = append(results, new(big.Int).SetUint64(response))
	}
	return results
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
	trace.Context = converToContext(callTrace.Context)
	trace.Steps = convertToInstrumentationSteps(callTrace.Steps)

	return *trace
}

func converToContext(context *pb.TransactionContext) instrumentation.Context {
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
		step.Error = response.Error
		step.Contract = convertToInstrumentationContract(response.Contract)
		step.GasCost = fmt.Sprint(response.GasCost)
		step.Stack = convertUint64ArrayToStringArray(response.Stack)
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

func convertUint64ArrayToStringArray(responses []uint64) []string {
	results := make([]string, 0, len(responses))
	for _, response := range responses {
		results = append(results, fmt.Sprint(response))
	}
	return results
}

func convertByteArrayToStringArray(responses []byte) []string {
	results := make([]string, 0, len(responses))
	for _, response := range responses {
		results = append(results, string(response))
	}
	return results
}
