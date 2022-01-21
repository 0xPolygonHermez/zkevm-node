package runtime_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
)

var (
	zeroAddr common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
	value                   = new(big.Int)
	gas      uint64         = 5000
	// code                    = []byte("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220404e37f487a89a932dca5e77faaf6ca2de3b991f93d230604b1b8daaef64766264736f6c63430008070033")
	code = []byte{
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

	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)

	host := runtime.NewMerkleTreeHost(tree.NewStateTree(mt, nil))
	res := testEvm.Run(contract, host, config)
	assert.Equal(t, uint64(4976), res.GasLeft)
	assert.NoError(t, res.Err)
}
