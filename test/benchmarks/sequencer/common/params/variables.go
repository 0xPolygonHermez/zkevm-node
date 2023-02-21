package params

import (
	"context"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Ctx is the context
	Ctx = context.Background()
	// PoolDbConfig is the pool db config
	PoolDbConfig = dbutils.NewPoolConfigFromEnv()
	// SequencerPrivateKey is the private key of the sequencer
	SequencerPrivateKey = operations.DefaultSequencerPrivateKey
	// ChainID is the chain id
	ChainID = operations.DefaultL2ChainID
	// OpsCfg is the operations config
	OpsCfg = operations.GetDefaultOperationsConfig()
	// ToAddress is the address to send the txs
	ToAddress = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	// To is the address to send the txs
	To = common.HexToAddress(ToAddress)
	// PrivateKey is the private key of the sender
	PrivateKey, _ = crypto.HexToECDSA(strings.TrimPrefix(SequencerPrivateKey, "0x"))
)
