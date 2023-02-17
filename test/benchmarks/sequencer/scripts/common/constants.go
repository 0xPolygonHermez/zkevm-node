package common

import (
	"context"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/test/operations"

	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

var (
	Ctx             = context.Background()
	PoolDbName      = testutils.GetEnv("POOL_DB_NAME", "pool_db")
	PoolDbUser      = testutils.GetEnv("POOL_DB_USER", "pool_user")
	PoolDbPass      = testutils.GetEnv("POOL_DB_PASS", "pool_password")
	PoolDbHost      = testutils.GetEnv("POOL_DB_HOST", "localhost")
	PoolDbPort      = testutils.GetEnv("POOL_DB_PORT", "5432")
	L2NetworkRPCURL = testutils.GetEnv("L2_NETWORK_RPC_URL", "http://localhost:8545")
	PrivateKey      = testutils.GetEnv("PRIVATE_KEY", operations.DefaultSequencerPrivateKey)
	ChainId         = testutils.GetEnv("CHAIN_ID", strconv.FormatUint(operations.DefaultL2ChainID, 10))
)
