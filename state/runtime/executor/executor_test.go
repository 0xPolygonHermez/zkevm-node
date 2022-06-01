package executor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/js"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/tracers"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

type account struct {
	address common.Address
}

func newAccount(address common.Address) *account {
	return &account{address: address}
}

func (a *account) Address() common.Address { return a.address } // { return common.Address{} }

func Test_Trace(t *testing.T) {
	var (
		trace         Trace
		tracer        Tracer
		previousDepth int
	)

	traceFile, err := os.Open("traces/op-call_2__full_trace_0.json")
	require.NoError(t, err)
	defer traceFile.Close()

	tracerFile, err := os.Open("tracer.json")
	require.NoError(t, err)
	defer tracerFile.Close()

	byteValue, err := ioutil.ReadAll(traceFile)
	require.NoError(t, err)

	err = json.Unmarshal(byteValue, &trace)
	require.NoError(t, err)

	byteCode, err := ioutil.ReadAll(tracerFile)
	require.NoError(t, err)

	err = json.Unmarshal(byteCode, &tracer)
	require.NoError(t, err)

	jsTracer, err := js.NewJsTracer(string(tracer.Code), new(tracers.Context))
	require.NoError(t, err)

	contextGas, ok := new(big.Int).SetString(trace.Context.Gas, 10)
	require.Equal(t, true, ok)

	value, ok := new(big.Int).SetString(trace.Context.Value, 10)
	require.Equal(t, true, ok)

	gasPrice, ok := new(big.Int).SetString(trace.Context.GasPrice, 10)
	require.Equal(t, true, ok)

	env := fakevm.NewFakeEVM(vm.BlockContext{BlockNumber: big.NewInt(1)}, vm.TxContext{GasPrice: gasPrice}, fakevm.FakeDB{StateRoot: []byte(trace.Context.OldStateRoot)}, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: jsTracer})

	jsTracer.CaptureTxStart(contextGas.Uint64())
	jsTracer.CaptureStart(env, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", common.Hex2Bytes(strings.TrimLeft(trace.Context.Input, "0x")), contextGas.Uint64(), value)

	log.Debugf("%v Steps", len(trace.Steps))

	stack := fakevm.Newstack()
	memory := fakevm.NewMemory()

	for _, step := range trace.Steps {
		gas, ok := new(big.Int).SetString(step.Gas, 10)
		require.Equal(t, true, ok)

		gasCost, ok := new(big.Int).SetString(step.GasCost, 10)
		require.Equal(t, true, ok)

		value, ok := new(big.Int).SetString(step.Contract.Value, 10)
		require.Equal(t, true, ok)

		op, ok := new(big.Int).SetString(step.Op, 0)
		require.Equal(t, true, ok)

		scope := &fakevm.ScopeContext{
			Contract: vm.NewContract(newAccount(common.HexToAddress(step.Contract.Caller)), newAccount(common.HexToAddress(step.Contract.Address)), value, gas.Uint64()),
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
			// TODO: Add return data when available
			jsTracer.CaptureState(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, common.Hex2Bytes(strings.TrimLeft(step.ReturnData, "0x")), step.Depth, nil)
		}

		// Set Memory
		if len(step.Memory) > 0 {
			memory.Resize(uint64(32*len(step.Memory) + 128))
			for offset, memoryContent := range step.Memory {
				memory.Set(uint64(offset*32)+128, 32, common.Hex2Bytes(memoryContent))
			}
		}

		// Set Stack
		stack = fakevm.Newstack()
		for _, stackContent := range step.Stack {
			// log.Debugf(stackContent)
			valueBigInt, ok := new(big.Int).SetString(stackContent, 0)
			require.Equal(t, true, ok)
			value, _ := uint256.FromBig(valueBigInt)
			stack.Push(value)
		}

		// Returning from a call or create
		if previousDepth < step.Depth {
			jsTracer.CaptureExit([]byte{}, gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		// Set StateRoot
		env.StateDB.StateRoot = []byte(step.StateRoot)
		previousDepth = step.Depth
	}

	gasUsed, ok := new(big.Int).SetString(trace.Context.GasUsed, 10)
	require.Equal(t, true, ok)

	jsTracer.CaptureTxEnd(gasUsed.Uint64())
	jsTracer.CaptureEnd([]byte(trace.Context.Output), gasUsed.Uint64(), time.Duration(trace.Context.Time), nil)
	result, err := jsTracer.GetResult()
	require.NoError(t, err)
	log.Debugf("%v", string(result))
}
