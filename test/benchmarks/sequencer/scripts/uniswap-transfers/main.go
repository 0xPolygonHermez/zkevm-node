package main

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	uniswaptransfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/e2e/uniswap-transfers"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/environment"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/results"
	uniswap "github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
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
	deployments := uniswap.DeployContractsAndAddLiquidity(l2Client, auth)
	deploymentTxsCount := uniswap.GetExecutedTransactionsCount()
	elapsedTimeForDeployments := time.Since(start)

	// Send Txs
	allTxs, err := transactions.SendAndWait(
		auth,
		l2Client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		nil,
		&deployments,
		uniswaptransfers.TxSender,
	)

	if err != nil {
		panic(err)
	}

	// Wait for Txs to be selected
	err = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, initialCount, params.NumberOfOperations+deploymentTxsCount)
	if err != nil {
		panic(err)
	}

	lastL2BlockTimestamp, err := state.GetLastL2BlockCreatedAt(params.Ctx, nil)
	if err != nil {
		panic(err)
	}
	elapsed := lastL2BlockTimestamp.Sub(start)
	results.PrintUniswapDeployments(elapsedTimeForDeployments, deploymentTxsCount)
	results.Print(l2Client, elapsed, allTxs)

	totalTxsCount := uniswap.GetExecutedTransactionsCount()
	log.Info("##############################")
	log.Info("# Uniswap Total Transactions #")
	log.Info("##############################")
	log.Infof("Number of total txs processed: %d", totalTxsCount)
}
