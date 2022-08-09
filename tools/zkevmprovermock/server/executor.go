package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
)

const (
	bits64 = 64
	bits32 = 32
)

// ExecutorMock represents and Executor mock server
type ExecutorMock struct {
	// address is the address on which the gRPC server will listen, eg. 0.0.0.0:50071
	address string

	tvContainer *testvector.Container

	// srv is an insance of the gRPC server.
	srv *grpc.Server

	// embedding an instance of pb.UnimplementedExecutorServiceServer will allow us
	// to implement all the required method interfaces.
	pb.UnimplementedExecutorServiceServer
}

// NewExecutorMock is the ExecutorMock constructor.
func NewExecutorMock(address string, tvContainer *testvector.Container) *ExecutorMock {
	return &ExecutorMock{
		address:     address,
		tvContainer: tvContainer,
	}
}

// Start sets up the stateDB server to process requests.
func (server *ExecutorMock) Start() {
	lis, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server.srv = grpc.NewServer()
	pb.RegisterExecutorServiceServer(server.srv, server)

	log.Infof("Executor mock server: listening at %s", server.address)
	if err := server.srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (server *ExecutorMock) Stop() {
	log.Info("Executor mock server: stopping...")
	server.srv.Stop()
}

// ProcessBatch implements the ProcessBatch gRPC method.
func (server *ExecutorMock) ProcessBatch(ctx context.Context, request *pb.ProcessBatchRequest) (*pb.ProcessBatchResponse, error) {
	processBatchResponse, err := server.tvContainer.FindProcessBatchResponse(hex.EncodeToString(request.BatchL2Data))
	if err != nil {
		return nil, err
	}

	cumulativeGasUSed, err := strconv.ParseUint(processBatchResponse.CumulativeGasUsed, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", processBatchResponse.CumulativeGasUsed, err)
	}

	responses := []*pb.ProcessTransactionResponse{}
	for _, response := range processBatchResponse.Responses {
		newResponse, err := translateResponse(response)
		if err != nil {
			return nil, err
		}
		responses = append(responses, newResponse)
	}

	return &pb.ProcessBatchResponse{
		CumulativeGasUsed:   cumulativeGasUSed,
		Responses:           responses,
		NewStateRoot:        common.Hex2Bytes(processBatchResponse.NewStateRoot[2:]),
		NewLocalExitRoot:    common.Hex2Bytes(processBatchResponse.NewLocalExitRoot[2:]),
		CntKeccakHashes:     processBatchResponse.CntKeccakHashes,
		CntPoseidonHashes:   processBatchResponse.CntPoseidonHashes,
		CntPoseidonPaddings: processBatchResponse.CntPoseidonPaddings,
		CntMemAligns:        processBatchResponse.CntMemAligns,
		CntArithmetics:      processBatchResponse.CntArithmetics,
		CntBinaries:         processBatchResponse.CntBinaries,
		CntSteps:            processBatchResponse.CntSteps,
	}, nil
}

func translateResponse(response *testvector.ProcessTransactionResponse) (*pb.ProcessTransactionResponse, error) {
	gasLeft, err := strconv.ParseUint(response.GasLeft, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasLeft, err)
	}
	gasUsed, err := strconv.ParseUint(response.GasUsed, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasUsed, err)
	}
	gasRefunded, err := strconv.ParseUint(response.GasRefunded, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasRefunded, err)
	}
	logs, err := translateLogs(response.Logs)
	if err != nil {
		return nil, err
	}

	callTrace, err := translateCallTrace(response.CallTrace)
	if err != nil {
		return nil, err
	}

	return &pb.ProcessTransactionResponse{
		TxHash:      common.Hex2Bytes(response.TxHash[2:]),
		Type:        response.Type,
		GasLeft:     gasLeft,
		GasUsed:     gasUsed,
		GasRefunded: gasRefunded,
		StateRoot:   common.Hex2Bytes(response.StateRoot[2:]),
		Logs:        logs,
		CallTrace:   callTrace,
	}, nil
}

func translateLogs(inputLogs []*testvector.Log) ([]*pb.Log, error) {
	logs := []*pb.Log{}

	for _, log := range inputLogs {
		topics := [][]byte{}

		for _, topic := range log.Topics {
			newTopic, err := hex.DecodeString(topic)
			if err != nil {
				return nil, err
			}
			topics = append(topics, newTopic)
		}

		data, err := hex.DecodeString(log.Data[0])
		if err != nil {
			return nil, err
		}

		newLog := &pb.Log{
			Address:     log.Address,
			Topics:      topics,
			Data:        data,
			BatchNumber: log.BatchNumber,
			TxHash:      common.Hex2Bytes(log.TxHash[2:]),
			TxIndex:     log.TxIndex,
			BatchHash:   common.Hex2Bytes(log.BatchHash[2:]),
			Index:       log.Index,
		}

		logs = append(logs, newLog)
	}
	return logs, nil
}

func translateCallTrace(callTrace *testvector.CallTrace) (*pb.CallTrace, error) {
	ctx, err := translateTransactionContext(callTrace.Context)
	if err != nil {
		return nil, err
	}

	steps, err := translateTransactionSteps(callTrace.Steps)
	if err != nil {
		return nil, err
	}

	return &pb.CallTrace{
		Context: ctx,
		Steps:   steps,
	}, nil
}

func translateTransactionContext(ctx *testvector.TransactionContext) (*pb.TransactionContext, error) {
	gas, err := strconv.ParseUint(ctx.Gas, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.Gas, err)
	}
	value, err := strconv.ParseUint(ctx.Value, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.Value, err)
	}
	gasUsed, err := strconv.ParseUint(ctx.GasUsed, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.GasUsed, err)
	}
	gasPrice, err := strconv.ParseUint(ctx.GasPrice, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.GasPrice, err)
	}
	executionTime, err := strconv.ParseUint(ctx.ExecutionTime, encoding.Base10, bits32)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint32: %v", ctx.ExecutionTime, err)
	}

	return &pb.TransactionContext{
		Type:          ctx.Type,
		From:          ctx.From,
		To:            ctx.To,
		Data:          common.Hex2Bytes(ctx.Data[2:]),
		Gas:           gas,
		Value:         value,
		GasUsed:       gasUsed,
		Batch:         common.Hex2Bytes(ctx.Batch[2:]),
		Output:        common.Hex2Bytes(ctx.Output[2:]),
		GasPrice:      gasPrice,
		ExecutionTime: uint32(executionTime),
		OldStateRoot:  common.Hex2Bytes(ctx.OldStateRoot[2:]),
	}, nil
}

func translateTransactionSteps(inputSteps []*testvector.TransactionStep) ([]*pb.TransactionStep, error) {
	steps := []*pb.TransactionStep{}

	for _, inputStep := range inputSteps {
		contract, err := translateContract(inputStep.Contract)
		if err != nil {
			return nil, err
		}

		gas, err := strconv.ParseUint(inputStep.RemainingGas, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.RemainingGas, err)
		}
		gasCost, err := strconv.ParseUint(inputStep.GasCost, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.GasCost, err)
		}
		gasRefund, err := strconv.ParseUint(inputStep.GasRefund, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.GasRefund, err)
		}
		op, err := strconv.ParseUint(inputStep.Op, encoding.Base10, bits32)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint32: %v", inputStep.Op, err)
		}
		pbErr, err := strconv.ParseUint(inputStep.Error, encoding.Base10, bits32)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint32: %v", inputStep.Error, err)
		}

		newStep := &pb.TransactionStep{
			StateRoot:  common.Hex2Bytes(inputStep.StateRoot[2:]),
			Depth:      inputStep.Depth,
			Pc:         inputStep.Pc,
			Gas:        gas,
			GasCost:    gasCost,
			GasRefund:  gasRefund,
			Op:         uint32(op),
			Stack:      inputStep.Stack,
			Memory:     common.Hex2Bytes(inputStep.Memory[0]),
			ReturnData: common.Hex2Bytes(inputStep.ReturnData[0]),
			Contract:   contract,
			Error:      pb.Error(pbErr),
		}
		steps = append(steps, newStep)
	}

	return steps, nil
}

func translateContract(contract *testvector.Contract) (*pb.Contract, error) {
	value, err := strconv.ParseUint(contract.Value, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", contract.Value, err)
	}
	gas, err := strconv.ParseUint(contract.Gas, encoding.Base10, bits64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert %q to uint64: %v", contract.Gas, err)
	}

	return &pb.Contract{
		Address: contract.Address,
		Caller:  contract.Caller,
		Data:    common.Hex2Bytes(contract.Data[2:]),
		Value:   value,
		Gas:     gas,
	}, nil
}
