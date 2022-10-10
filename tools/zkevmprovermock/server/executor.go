package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"

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

	if strings.HasPrefix(processBatchResponse.NewStateRoot, "0x") { // nolint
		processBatchResponse.NewStateRoot = processBatchResponse.NewStateRoot[2:]
	}
	if strings.HasPrefix(processBatchResponse.NewLocalExitRoot, "0x") { // nolint
		processBatchResponse.NewLocalExitRoot = processBatchResponse.NewLocalExitRoot[2:]
	}

	return &pb.ProcessBatchResponse{
		CumulativeGasUsed:   cumulativeGasUSed,
		Responses:           responses,
		NewStateRoot:        common.Hex2Bytes(processBatchResponse.NewStateRoot),
		NewLocalExitRoot:    common.Hex2Bytes(processBatchResponse.NewLocalExitRoot),
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
	var err error

	var gasLeft uint64
	if response.GasLeft != "" {
		gasLeft, err = strconv.ParseUint(response.GasLeft, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasLeft, err)
		}
	}

	var gasUsed uint64
	if response.GasUsed != "" {
		gasUsed, err = strconv.ParseUint(response.GasUsed, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasUsed, err)
		}
	}

	var gasRefunded uint64
	if response.GasRefunded != "" {
		gasRefunded, err = strconv.ParseUint(response.GasRefunded, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", response.GasRefunded, err)
		}
	}

	logs, err := translateLogs(response.Logs)
	if err != nil {
		return nil, err
	}

	callTrace, err := translateCallTrace(response.CallTrace)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(response.TxHash, "0x") { // nolint
		response.TxHash = response.TxHash[2:]
	}
	if strings.HasPrefix(response.StateRoot, "0x") { // nolint
		response.StateRoot = response.StateRoot[2:]
	}

	return &pb.ProcessTransactionResponse{
		TxHash:      common.Hex2Bytes(response.TxHash),
		Type:        response.Type,
		GasLeft:     gasLeft,
		GasUsed:     gasUsed,
		GasRefunded: gasRefunded,
		StateRoot:   common.Hex2Bytes(response.StateRoot),
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

		data := []byte{}
		var err error
		if len(log.Data) > 0 {
			data, err = hex.DecodeString(log.Data[0])
			if err != nil {
				return nil, err
			}
		}

		if strings.HasPrefix(log.TxHash, "0x") { // nolint
			log.TxHash = log.TxHash[2:]
		}
		if strings.HasPrefix(log.BatchHash, "0x") { // nolint
			log.BatchHash = log.BatchHash[2:]
		}

		newLog := &pb.Log{
			Address:     log.Address,
			Topics:      topics,
			Data:        data,
			BatchNumber: log.BatchNumber,
			TxHash:      common.Hex2Bytes(log.TxHash),
			TxIndex:     log.TxIndex,
			BatchHash:   common.Hex2Bytes(log.BatchHash),
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
	var err error
	var gas uint64
	if ctx.Gas != "" {
		gas, err = strconv.ParseUint(ctx.Gas, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.Gas, err)
		}
	}
	var value uint64
	if ctx.Value != "" {
		value, err = strconv.ParseUint(ctx.Value, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.Value, err)
		}
	}

	var gasUsed uint64
	if ctx.GasUsed != "" {
		gasUsed, err = strconv.ParseUint(ctx.GasUsed, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.GasUsed, err)
		}
	}

	var gasPrice uint64
	if ctx.GasPrice != "" {
		gasPrice, err = strconv.ParseUint(ctx.GasPrice, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", ctx.GasPrice, err)
		}
	}

	var executionTime uint64
	if ctx.ExecutionTime != "" {
		executionTime, err = strconv.ParseUint(ctx.ExecutionTime, encoding.Base10, bits32)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint32: %v", ctx.ExecutionTime, err)
		}
	}
	if strings.HasPrefix(ctx.Data, "0x") { // nolint
		ctx.Data = ctx.Data[2:]
	}
	if strings.HasPrefix(ctx.Batch, "0x") { // nolint
		ctx.Batch = ctx.Batch[2:]
	}
	if strings.HasPrefix(ctx.Output, "0x") { // nolint
		ctx.Output = ctx.Output[2:]
	}
	if strings.HasPrefix(ctx.OldStateRoot, "0x") { // nolint
		ctx.OldStateRoot = ctx.OldStateRoot[2:]
	}

	return &pb.TransactionContext{
		Type:          ctx.Type,
		From:          ctx.From,
		To:            ctx.To,
		Data:          common.Hex2Bytes(ctx.Data),
		Gas:           gas,
		Value:         new(big.Int).SetUint64(value).String(),
		GasUsed:       gasUsed,
		Batch:         common.Hex2Bytes(ctx.Batch),
		Output:        common.Hex2Bytes(ctx.Output),
		GasPrice:      new(big.Int).SetUint64(gasPrice).String(),
		ExecutionTime: uint32(executionTime),
		OldStateRoot:  common.Hex2Bytes(ctx.OldStateRoot),
	}, nil
}

func translateTransactionSteps(inputSteps []*testvector.TransactionStep) ([]*pb.TransactionStep, error) {
	steps := []*pb.TransactionStep{}

	for _, inputStep := range inputSteps {
		contract, err := translateContract(inputStep.Contract)
		if err != nil {
			return nil, err
		}

		var gas uint64
		if inputStep.RemainingGas != "" {
			gas, err = strconv.ParseUint(inputStep.RemainingGas, encoding.Base10, bits64)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.RemainingGas, err)
			}
		}
		var gasCost uint64
		if inputStep.GasCost != "" {
			gasCost, err = strconv.ParseUint(inputStep.GasCost, encoding.Base10, bits64)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.GasCost, err)
			}
		}
		var gasRefund uint64
		if inputStep.GasRefund != "" {
			gasRefund, err = strconv.ParseUint(inputStep.GasRefund, encoding.Base10, bits64)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %q to uint64: %v", inputStep.GasRefund, err)
			}
		}
		var op uint64
		if inputStep.Op != "" {
			if strings.HasPrefix(inputStep.Op, "0x") { // nolint
				inputStep.Op = inputStep.Op[2:]
			}
			opBI, ok := new(big.Int).SetString(inputStep.Op, encoding.Base16)
			if !ok {
				return nil, fmt.Errorf("Could not convert base16 %q to big int", inputStep.Op)
			}
			op = opBI.Uint64()
		}
		var pbErr uint64
		if inputStep.Error != "" {
			pbErr, err = strconv.ParseUint(inputStep.Error, encoding.Base10, bits32)
			if err != nil {
				return nil, fmt.Errorf("Could not convert %q to uint32: %v", inputStep.Error, err)
			}
		}

		if strings.HasPrefix(inputStep.StateRoot, "0x") { // nolint
			inputStep.StateRoot = inputStep.StateRoot[2:]
		}

		memory := []byte{}
		if len(inputStep.Memory) > 0 {
			memory = common.Hex2Bytes(inputStep.Memory[0])
		}
		returnData := []byte{}
		if len(inputStep.ReturnData) > 0 {
			returnData = common.Hex2Bytes(inputStep.ReturnData[0])
		}

		newStep := &pb.TransactionStep{
			StateRoot:  common.Hex2Bytes(inputStep.StateRoot),
			Depth:      inputStep.Depth,
			Pc:         inputStep.Pc,
			Gas:        gas,
			GasCost:    gasCost,
			GasRefund:  gasRefund,
			Op:         uint32(op),
			Stack:      inputStep.Stack,
			Memory:     memory,
			ReturnData: returnData,
			Contract:   contract,
			Error:      pb.Error(pbErr),
		}
		steps = append(steps, newStep)
	}

	return steps, nil
}

func translateContract(contract *testvector.Contract) (*pb.Contract, error) {
	var err error
	var value uint64
	if contract.Value != "" {
		value, err = strconv.ParseUint(contract.Value, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", contract.Value, err)
		}
	}
	var gas uint64
	if contract.Gas != "" {
		gas, err = strconv.ParseUint(contract.Gas, encoding.Base10, bits64)
		if err != nil {
			return nil, fmt.Errorf("Could not convert %q to uint64: %v", contract.Gas, err)
		}
	}

	if strings.HasPrefix(contract.Data, "0x") { // nolint
		contract.Data = contract.Data[2:]
	}

	return &pb.Contract{
		Address: contract.Address,
		Caller:  contract.Caller,
		Data:    common.Hex2Bytes(contract.Data),
		Value:   new(big.Int).SetUint64(value).String(),
		Gas:     gas,
	}, nil
}
