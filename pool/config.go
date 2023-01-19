package pool

import 	"github.com/0xPolygonHermez/zkevm-node/db"

type Config struct {
	// FreeClaimGasLimit is the max gas allowed use to do a free claim
	FreeClaimGasLimit uint64    `mapstructure:"FreeClaimGasLimit"`
	DB                db.Config `mapstructure:"DB"`
}