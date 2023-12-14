package test

import (
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
)

func InitOrResetDB(cfg db.Config) {
	if err := dbutils.InitOrResetState(cfg); err != nil {
		panic(err)
	}
}
