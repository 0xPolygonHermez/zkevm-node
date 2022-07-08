package runtime_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/evm"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
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
	cfg = dbutils.NewConfigFromEnv()
)

func TestRuntime(t *testing.T) {
	testEvm := evm.NewEVM()
	contract := runtime.NewContract(1, zeroAddr, zeroAddr, zeroAddr, value, gas, code)
	config := &runtime.ForksInTime{
		EIP158: true,
	}

	stateDb, err := db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	res := testEvm.Run(context.Background(), contract, nil, config)
	assert.Equal(t, uint64(4976), res.GasLeft)
	assert.NoError(t, res.Err)
}
