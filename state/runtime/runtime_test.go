package runtime_test

/*
import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/stretchr/testify/assert"
)

var (
	zeroAddr common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
	value                   = new(big.Int)
	gas      uint64         = 5000
	code                    = []byte{
		evm.PUSH1, 0x01, evm.PUSH1, 0x02, evm.ADD,
		evm.PUSH1, 0x00, evm.MSTORE8,
		evm.PUSH1, 0x01, evm.PUSH1, 0x00, evm.RETURN,
	}
)

func TestRuntime(t *testing.T) {
	testEvm := evm.NewEVM()
	contract := runtime.NewContract(1, zeroAddr, zeroAddr, zeroAddr, value, gas, code)
	config := &runtime.ForksInTime{
		EIP158: true,
	}
	host := &runtime.MockHost{}
	res := testEvm.Run(contract, host, config)
	assert.Equal(t, "5000", res.GasLeft)
	assert.NoError(t, res.Err)
}
*/
