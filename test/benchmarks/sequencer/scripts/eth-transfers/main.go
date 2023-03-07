package main

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	ethtransfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/eth-transfers"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/environment"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/results"
)

func main() {
	var (
		err error
	)
	ctx, pl, state, l2Client, auth := environment.Init()
	initialCount, err := pl.CountTransactionsByStatus(params.Ctx, pool.TxStatusSelected)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	// Send Txs
	err = transactions.SendAndWait(
		ctx,
		auth,
		l2Client,
		pl.CountTransactionsByStatus,
		params.NumberOfTxs,
		nil,
		ethtransfers.TxSender,
	)
	if err != nil {
		panic(err)
	}

	// Wait for Txs to be selected
	err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfTxs)
	if err != nil {
		panic(err)
	}

	lastL2BlockTimestamp, err := state.GetLastL2BlockCreatedAt(params.Ctx, nil)
	if err != nil {
		panic(err)
	}
	elapsed := lastL2BlockTimestamp.Sub(start)
	results.Print(elapsed)
}
