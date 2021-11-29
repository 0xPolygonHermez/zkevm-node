package state

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	dbutils "github.com/hermeznetwork/hermez-core/test/db"
	"testing"
)

func TestMainTest(t *testing.T) {
	// init db
	db, err := dbutils.ConnectToTestSQLDB()
	if err != nil {
		panic(err)
	}
	state := NewState(db)
	ctx := context.Background()
	addr := common.Hex2Bytes("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	fmt.Println(addr)
	block, err := state.GetBatchByNumber(ctx, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
