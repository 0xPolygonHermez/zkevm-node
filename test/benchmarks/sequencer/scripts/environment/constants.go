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
	// L2NetworkRPCURL is the RPC URL of the L2 network
	L2NetworkRPCURL = testutils.GetEnv("RPC_URL", operations.DefaultL2NetworkURL)

	// PoolDB Credentials
	poolDbName = testutils.GetEnv("POOLDB_DBNAME", "pool_db")
	poolDbUser = testutils.GetEnv("POOLDB_USER", "pool_user")
	poolDbPass = testutils.GetEnv("POOLDB_PASS", "pool_password")
	poolDbHost = testutils.GetEnv("POOLDB_HOST", "localhost")
	poolDbPort = testutils.GetEnv("POOLDB_PORT", "5433")
)
