package main

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	ethtransfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/e2e/eth-transfers"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/environment"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/results"
)

func main() {
	var (
		err error
	)
	pl, state, l2Client, auth := environment.Init()
	initialCount, err := pl.CountTransactionsByStatus(params.Ctx, pool.TxStatusSelected)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	// Send Txs
	allTxs, err := transactions.SendAndWait(
		auth,
		l2Client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		nil,
		nil,
		ethtransfers.TxSender,
	)
	if err != nil {
		fmt.Println(auth.Nonce)
		panic(err)
	}

	// Wait for Txs to be selected
	err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfOperations)
	if err != nil {
		panic(err)
	}

	lastL2BlockTimestamp, err := state.GetLastL2BlockCreatedAt(params.Ctx, nil)
	if err != nil {
		panic(err)
	}
	elapsed := lastL2BlockTimestamp.Sub(start)
	results.Print(l2Client, elapsed, allTxs)
}
