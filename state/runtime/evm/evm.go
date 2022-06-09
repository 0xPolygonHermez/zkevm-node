package evm

import (
	"context"

	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// EVM is the Ethereum Virtual Machine
type EVM struct {
	instrumented bool
}

// NewEVM creates a new EVM
func NewEVM() *EVM {
	return &EVM{}
}

// EnableInstrumentation enables EVM instrumentation
func (e *EVM) EnableInstrumentation() {
	e.instrumented = true
}

// DisableInstrumentation enables EVM instrumentation
func (e *EVM) DisableInstrumentation() {
	e.instrumented = false
}

// CanRun implements the runtime interface
func (e *EVM) CanRun(*runtime.Contract, runtime.Host, *runtime.ForksInTime) bool {
	return true
}

// Name implements the runtime interface
func (e *EVM) Name() string {
	return "hermez-evm"
}

// Run implements the runtime interface
func (e *EVM) Run(ctx context.Context, c *runtime.Contract, host runtime.Host, config *runtime.ForksInTime) *runtime.ExecutionResult {
	contract := acquireState()
	contract.resetReturnData()

	contract.msg = c
	contract.code = c.Code
	contract.evm = e
	contract.gas = c.Gas
	contract.host = host
	contract.config = config
	contract.instrumented = e.instrumented

	contract.bitmap.setCode(c.Code)

	ret, vmTrace, structLogs, executorTrace, err := contract.Run(ctx)

	var returnValue []byte
	returnValue = append(returnValue[:0], ret...)

	gasLeft := contract.gas

	releaseState(contract)

	if err != nil && err != errRevert {
		gasLeft = 0
	}

	return &runtime.ExecutionResult{
		ReturnValue:   returnValue,
		GasLeft:       gasLeft,
		Err:           err,
		VMTrace:       vmTrace,
		StructLogs:    structLogs,
		ExecutorTrace: executorTrace,
	}
}
