package main

import (
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/metrics"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	erc20transfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/e2e/erc20-transfers"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/environment"
)

func ExecuteERC20Transfers(numOps uint64) uint64 {
	var (
		err error
	)

	pl, l2Client, auth := environment.Init()
	initialCount, err := pl.CountTransactionsByStatus(params.Ctx, pool.TxStatusSelected)
	if err != nil {
		panic(err)
	}

	erc20SC, err := erc20transfers.DeployERC20Contract(l2Client, params.Ctx, auth)
	allTxs, err := transactions.SendAndWait(
		auth,
		l2Client,
		pl.GetTxsByStatus,
		numOps,
		erc20SC,
		nil,
		erc20transfers.TxSender,
	)
	if err != nil {
		panic(err)
	}

	// Wait for Txs to be selected
	err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, numOps)
	if err != nil {
		panic(err)
	}

	totalGas := metrics.GetTotalGasUsedFromTxs(l2Client, allTxs)

	return totalGas
}
