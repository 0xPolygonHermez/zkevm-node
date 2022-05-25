package executor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/js"
	"github.com/stretchr/testify/require"
)

type account struct{}

func (account) SubBalance(amount *big.Int)                          {}
func (account) AddBalance(amount *big.Int)                          {}
func (account) SetAddress(common.Address)                           {}
func (account) Value() *big.Int                                     { return nil }
func (account) SetBalance(*big.Int)                                 {}
func (account) SetNonce(uint64)                                     {}
func (account) Balance() *big.Int                                   { return nil }
func (account) Address() common.Address                             { return common.Address{} }
func (account) SetCode(common.Hash, []byte)                         {}
func (account) ForEachStorage(cb func(key, value common.Hash) bool) {}

func Test_Trace(t *testing.T) {
	var (
		trace  Trace
		tracer Tracer
	)

	traceFile, err := os.Open("demo_trace_2.json")
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

	jsTracer, err := js.NewJsTracer(string(tracer.Code))
	require.NoError(t, err)

	contextGas, ok := new(big.Int).SetString(trace.Context.Gas, 10)
	require.Equal(t, true, ok)

	value, ok := new(big.Int).SetString(trace.Context.Value, 10)
	require.Equal(t, true, ok)

	gasPrice, ok := new(big.Int).SetString(trace.Context.GasPrice, 10)
	require.Equal(t, true, ok)

	env := fakevm.NewFakeEVM(vm.BlockContext{BlockNumber: big.NewInt(1)}, vm.TxContext{GasPrice: gasPrice}, fakevm.FakeDB{}, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: jsTracer})

	jsTracer.CaptureStart(env, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", common.Hex2Bytes(strings.TrimLeft(trace.Context.Input, "0x")), contextGas.Uint64(), value)

	log.Debugf("%v Steps", len(trace.Steps))

	for _, step := range trace.Steps {
		// Prepare Scope
		memory := fakevm.NewMemory()
		for offset, memoryContent := range step.Memory {
			memory.Set(uint64(offset), uint64(len(memoryContent)), common.Hex2Bytes(memoryContent))
		}

		stack := fakevm.Newstack()

		scope := &fakevm.ScopeContext{
			Contract: vm.NewContract(&account{}, &account{}, big.NewInt(0), 0),
			Memory:   memory,
			Stack:    stack,
		}
		gas, ok := new(big.Int).SetString(step.Gas, 10)
		require.Equal(t, true, ok)

		gasCost, ok := new(big.Int).SetString(step.GasCost, 10)
		require.Equal(t, true, ok)

		value, ok := new(big.Int).SetString(step.Contract.Value, 10)
		require.Equal(t, true, ok)

		// log.Debugf("Op=%v", step.Op)

		op, ok := new(big.Int).SetString(step.Op, 0)
		require.Equal(t, true, ok)

		opcode := vm.OpCode(op.Uint64()).String()
		log.Debugf("OpCode=%v", opcode)

		if step.Error != "" {
			err := fmt.Errorf(step.Error)
			jsTracer.CaptureFault(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.Depth, err)
		} else {
			if opcode == "CALL" {
				log.Debugf("THIS IS JUST TO SET A BREAKPOINT")
			}
			jsTracer.CaptureState(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, common.Hex2Bytes(strings.TrimLeft(step.Contract.Input, "0x")), step.Depth, nil)
		}

		if opcode == "CREATE" || opcode == "CREATE2" || opcode == "CALL" || opcode == "CALLCODE" || opcode == "DELEGATECALL" || opcode == "STATICCALL" {
			log.Debugf(opcode)
			jsTracer.CaptureEnter(vm.OpCode(op.Uint64()), common.HexToAddress(step.Contract.Caller), common.HexToAddress(step.Contract.Address), common.Hex2Bytes(strings.TrimLeft(step.Contract.Input, "0x")), gas.Uint64(), value)
			jsTracer.CaptureExit([]byte{}, gasCost.Uint64(), fmt.Errorf(step.Error))
		}
	}

	jsTracer.CaptureEnd([]byte(trace.Context.Output), 0, 1, nil)
	result, err := jsTracer.GetResult()
	require.NoError(t, err)
	log.Debugf("%v", string(result))
}
