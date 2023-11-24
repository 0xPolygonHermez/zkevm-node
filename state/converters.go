package state

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ConvertToCounters extracts ZKCounters from a ProcessBatchResponse
func ConvertToCounters(resp *executor.ProcessBatchResponse) ZKCounters {
	return ZKCounters{
		CumulativeGasUsed:    resp.CumulativeGasUsed,
		UsedKeccakHashes:     resp.CntKeccakHashes,
		UsedPoseidonHashes:   resp.CntPoseidonHashes,
		UsedPoseidonPaddings: resp.CntPoseidonPaddings,
		UsedMemAligns:        resp.CntMemAligns,
		UsedArithmetics:      resp.CntArithmetics,
		UsedBinaries:         resp.CntBinaries,
		UsedSteps:            resp.CntSteps,
	}
}

// TestConvertToProcessBatchResponse for test purposes
func (s *State) TestConvertToProcessBatchResponse(response *executor.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	return s.convertToProcessBatchResponse(response)
}

func (s *State) convertToProcessBatchResponse(response *executor.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	responses, err := s.convertToProcessTransactionResponse(response.Responses)
	if err != nil {
		return nil, err
	}

	readWriteAddresses, err := convertToReadWriteAddresses(response.ReadWriteAddresses)
	if err != nil {
		return nil, err
	}

	isExecutorLevelError := (response.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR)
	isRomLevelError := false
	isRomOOCError := false

	if response.Responses != nil {
		for _, resp := range response.Responses {
			if resp.Error != executor.RomError_ROM_ERROR_NO_ERROR {
				isRomLevelError = true
				break
			}
		}

		if len(response.Responses) > 0 {
			// Check out of counters
			errorToCheck := response.Responses[len(response.Responses)-1].Error
			isRomOOCError = executor.IsROMOutOfCountersError(errorToCheck)
		}
	}

	return &ProcessBatchResponse{
		NewStateRoot:         common.BytesToHash(response.NewStateRoot),
		NewAccInputHash:      common.BytesToHash(response.NewAccInputHash),
		NewLocalExitRoot:     common.BytesToHash(response.NewLocalExitRoot),
		NewBatchNumber:       response.NewBatchNum,
		UsedZkCounters:       convertToCounters(response),
		Responses:            responses,
		ExecutorError:        executor.ExecutorErr(response.Error),
		ReadWriteAddresses:   readWriteAddresses,
		FlushID:              response.FlushId,
		StoredFlushID:        response.StoredFlushId,
		ProverID:             response.ProverId,
		IsExecutorLevelError: isExecutorLevelError,
		IsRomLevelError:      isRomLevelError,
		IsRomOOCError:        isRomOOCError,
	}, nil
}

// IsStateRootChanged returns true if the transaction changes the state root
func IsStateRootChanged(err executor.RomError) bool {
	return !executor.IsIntrinsicError(err) && !executor.IsROMOutOfCountersError(err) && err != executor.RomError_ROM_ERROR_INVALID_RLP
}

func convertToReadWriteAddresses(addresses map[string]*executor.InfoReadWrite) (map[common.Address]*InfoReadWrite, error) {
	results := make(map[common.Address]*InfoReadWrite, len(addresses))

	for addr, addrInfo := range addresses {
		var nonce *uint64 = nil
		var balance *big.Int = nil
		var ok bool

		address := common.HexToAddress(addr)

		if addrInfo.Nonce != "" {
			bigNonce, ok := new(big.Int).SetString(addrInfo.Nonce, encoding.Base10)
			if !ok {
				log.Debugf("received nonce as string: %v", addrInfo.Nonce)
				return nil, fmt.Errorf("error while parsing address nonce")
			}
			nonceNp := bigNonce.Uint64()
			nonce = &nonceNp
		}

		if addrInfo.Balance != "" {
			balance, ok = new(big.Int).SetString(addrInfo.Balance, encoding.Base10)
			if !ok {
				log.Debugf("received balance as string: %v", addrInfo.Balance)
				return nil, fmt.Errorf("error while parsing address balance")
			}
		}

		results[address] = &InfoReadWrite{Address: address, Nonce: nonce, Balance: balance}
	}

	return results, nil
}

func (s *State) convertToProcessTransactionResponse(responses []*executor.ProcessTransactionResponse) ([]*ProcessTransactionResponse, error) {
	results := make([]*ProcessTransactionResponse, 0, len(responses))
	for _, response := range responses {
		result := new(ProcessTransactionResponse)
		result.TxHash = common.BytesToHash(response.TxHash)
		result.Type = response.Type
		result.ReturnValue = response.ReturnValue
		result.GasLeft = response.GasLeft
		result.GasUsed = response.GasUsed
		result.GasRefunded = response.GasRefunded
		result.RomError = executor.RomErr(response.Error)
		result.CreateAddress = common.HexToAddress(response.CreateAddress)
		result.StateRoot = common.BytesToHash(response.StateRoot)
		result.Logs = convertToLog(response.Logs)
		result.ChangesStateRoot = IsStateRootChanged(response.Error)
		fullTrace, err := convertToFullTrace(response.FullTrace)
		if err != nil {
			return nil, err
		}
		result.FullTrace = *fullTrace
		result.EffectiveGasPrice = response.EffectiveGasPrice
		result.EffectivePercentage = response.EffectivePercentage
		result.HasGaspriceOpcode = (response.HasGaspriceOpcode == 1)
		result.HasBalanceOpcode = (response.HasBalanceOpcode == 1)

		tx := new(types.Transaction)

		if response.Error != executor.RomError_ROM_ERROR_INVALID_RLP && len(response.GetRlpTx()) > 0 {
			tx, err = DecodeTx(common.Bytes2Hex(response.GetRlpTx()))
			if err != nil {
				timestamp := time.Now()
				log.Errorf("error decoding rlp returned by executor %v at %v", err, timestamp)

				event := &event.Event{
					ReceivedAt: timestamp,
					Source:     event.Source_Node,
					Level:      event.Level_Error,
					EventID:    event.EventID_ExecutorRLPError,
					Json:       string(response.GetRlpTx()),
				}

				eventErr := s.eventLog.LogEvent(context.Background(), event)
				if eventErr != nil {
					log.Errorf("error storing payload: %v", err)
				}

				return nil, err
			}
		} else {
			log.Warnf("ROM_ERROR_INVALID_RLP returned by the executor")
		}

		if tx != nil {
			result.Tx = *tx
			log.Debugf("ProcessTransactionResponse[TxHash]: %v", result.TxHash)
			if response.Error == executor.RomError_ROM_ERROR_NO_ERROR {
				log.Debugf("ProcessTransactionResponse[Nonce]: %v", result.Tx.Nonce())
			}
			log.Debugf("ProcessTransactionResponse[StateRoot]: %v", result.StateRoot.String())
			log.Debugf("ProcessTransactionResponse[Error]: %v", result.RomError)
			log.Debugf("ProcessTransactionResponse[GasUsed]: %v", result.GasUsed)
			log.Debugf("ProcessTransactionResponse[GasLeft]: %v", result.GasLeft)
			log.Debugf("ProcessTransactionResponse[GasRefunded]: %v", result.GasRefunded)
			log.Debugf("ProcessTransactionResponse[ChangesStateRoot]: %v", result.ChangesStateRoot)
			log.Debugf("ProcessTransactionResponse[EffectiveGasPrice]: %v", result.EffectiveGasPrice)
			log.Debugf("ProcessTransactionResponse[EffectivePercentage]: %v", result.EffectivePercentage)
		}

		results = append(results, result)
	}

	return results, nil
}

func convertToLog(protoLogs []*executor.Log) []*types.Log {
	logs := make([]*types.Log, 0, len(protoLogs))

	for _, protoLog := range protoLogs {
		log := new(types.Log)
		log.Address = common.HexToAddress(protoLog.Address)
		log.Topics = convertToTopics(protoLog.Topics)
		log.Data = protoLog.Data
		log.TxHash = common.BytesToHash(protoLog.TxHash)
		log.TxIndex = uint(protoLog.TxIndex)
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

func convertToFullTrace(fullTrace *executor.FullTrace) (*instrumentation.FullTrace, error) {
	trace := new(instrumentation.FullTrace)
	if fullTrace != nil {
		trace.Context = convertToContext(fullTrace.Context)
		steps, err := convertToInstrumentationSteps(fullTrace.Steps)
		if err != nil {
			return nil, err
		}
		trace.Steps = steps
	}

	return trace, nil
}

func convertToContext(context *executor.TransactionContext) instrumentation.Context {
	return instrumentation.Context{
		Type:         context.Type,
		From:         context.From,
		To:           context.To,
		Input:        context.Data,
		Gas:          context.Gas,
		Value:        hex.DecodeBig(context.Value),
		Output:       context.Output,
		GasPrice:     context.GasPrice,
		OldStateRoot: common.BytesToHash(context.OldStateRoot),
		Time:         uint64(context.ExecutionTime),
		GasUsed:      context.GasUsed,
	}
}

func convertToInstrumentationSteps(responses []*executor.TransactionStep) ([]instrumentation.Step, error) {
	results := make([]instrumentation.Step, 0, len(responses))
	for _, response := range responses {
		step := new(instrumentation.Step)
		step.StateRoot = common.BytesToHash(response.StateRoot)
		step.Depth = int(response.Depth)
		step.Pc = response.Pc
		step.Gas = response.Gas
		step.OpCode = fakevm.OpCode(response.Op).String()
		step.Refund = response.GasRefund
		step.Op = uint64(response.Op)
		err := executor.RomErr(response.Error)
		if err != nil {
			step.Error = err
		}
		step.Contract = convertToInstrumentationContract(response.Contract)
		step.GasCost = response.GasCost
		step.Stack = make([]*big.Int, 0, len(response.Stack))
		for _, s := range response.Stack {
			if len(s)%2 != 0 {
				s = "0" + s
			}
			bi, ok := new(big.Int).SetString(s, hex.Base)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			step.Stack = append(step.Stack, bi)
		}
		step.MemorySize = response.MemorySize
		step.MemoryOffset = response.MemoryOffset
		step.Memory = make([]byte, len(response.Memory))
		copy(step.Memory, response.Memory)
		step.ReturnData = make([]byte, len(response.ReturnData))
		copy(step.ReturnData, response.ReturnData)
		step.Storage = make(map[common.Hash]common.Hash, len(response.Storage))
		for k, v := range response.Storage {
			addr := common.BytesToHash(hex.DecodeBig(k).Bytes())
			value := common.BytesToHash(hex.DecodeBig(v).Bytes())
			step.Storage[addr] = value
		}
		results = append(results, *step)
	}
	return results, nil
}

func convertToInstrumentationContract(response *executor.Contract) instrumentation.Contract {
	return instrumentation.Contract{
		Address: common.HexToAddress(response.Address),
		Caller:  common.HexToAddress(response.Caller),
		Value:   hex.DecodeBig(response.Value),
		Input:   response.Data,
		Gas:     response.Gas,
	}
}

func convertToCounters(resp *executor.ProcessBatchResponse) ZKCounters {
	return ZKCounters{
		CumulativeGasUsed:    resp.CumulativeGasUsed,
		UsedKeccakHashes:     resp.CntKeccakHashes,
		UsedPoseidonHashes:   resp.CntPoseidonHashes,
		UsedPoseidonPaddings: resp.CntPoseidonPaddings,
		UsedMemAligns:        resp.CntMemAligns,
		UsedArithmetics:      resp.CntArithmetics,
		UsedBinaries:         resp.CntBinaries,
		UsedSteps:            resp.CntSteps,
	}
}
