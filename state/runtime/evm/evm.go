package evm

import "github.com/hermeznetwork/hermez-core/state/runtime"

// EVM is the Ethereum Virtual Machine
type EVM struct {
}

// NewEVM creates a new EVM
func NewEVM() *EVM {
	return &EVM{}
}

func (e *EVM) Run(c *runtime.Contract) *runtime.ExecutionResult {
	contract := acquireState()
	contract.resetReturnData()

	contract.msg = c
	contract.code = c.Code
	contract.evm = e
	contract.gas = c.Gas
	// contract.host = host
	// contract.config = config

	contract.bitmap.setCode(c.Code)

	ret, err := contract.Run()

	// We are probably doing this append magic to make sure that the slice doesn't have more capacity than it needs
	var returnValue []byte
	returnValue = append(returnValue[:0], ret...)

	gasLeft := contract.gas

	releaseState(contract)

	if err != nil && err != errRevert {
		gasLeft = 0
	}

	return &runtime.ExecutionResult{
		ReturnValue: returnValue,
		GasLeft:     gasLeft,
		Err:         err,
	}
}
