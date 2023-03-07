package environment

import (
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

var (
	// IntBase is the base for the conversion of strings to integers
	IntBase = 10
	// PrivateKey is the private key of the sequencer
	PrivateKey = testutils.GetEnv("PRIVATE_KEY", operations.DefaultSequencerPrivateKey)
	// L2ChainId is the chain id of the L2 network
	L2ChainId = testutils.GetEnv("CHAIN_ID", strconv.FormatUint(operations.DefaultL2ChainID, IntBase))
	//Erc20TokenAddress is the address of the ERC20 token
	Erc20TokenAddress = testutils.GetEnv("ERC20_TOKEN_ADDRESS", "0x729fc461b26f69cf75a31182788eaf722b08c240")

	l2NetworkRPCURL = testutils.GetEnv("L2_NETWORK_RPC_URL", operations.DefaultL2NetworkURL)

	// StateDB Credentials
	stateDbName = testutils.GetEnv("STATE_DB_NAME", "state_db")
	stateDbUser = testutils.GetEnv("STATE_DB_USER", "state_user")
	stateDbPass = testutils.GetEnv("STATE_DB_PASS", "state_password")
	stateDbHost = testutils.GetEnv("STATE_DB_HOST", "localhost")
	stateDbPort = testutils.GetEnv("STATE_DB_PORT", "5432")

	// PoolDB Credentials
	poolDbName = testutils.GetEnv("POOL_DB_NAME", "pool_db")
	poolDbUser = testutils.GetEnv("POOL_DB_USER", "pool_user")
	poolDbPass = testutils.GetEnv("POOL_DB_PASS", "pool_password")
	poolDbHost = testutils.GetEnv("POOL_DB_HOST", "localhost")
	poolDbPort = testutils.GetEnv("POOL_DB_PORT", "5433")
)
