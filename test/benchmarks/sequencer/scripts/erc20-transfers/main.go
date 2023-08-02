package main

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	erc20transfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/e2e/erc20-transfers"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/environment"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common/results"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/ethereum/go-ethereum/common"
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
	erc20SC, err := ERC20.NewERC20(common.HexToAddress(environment.Erc20TokenAddress), l2Client)
	if err != nil {
		panic(err)
	}
	// Send Txs
	err = transactions.SendAndWait(
		auth,
		l2Client,
		pl.GetTxsByStatus,
		params.NumberOfOperations,
		erc20SC,
		nil,
		erc20transfers.TxSender,
	)
	if err != nil {
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
	results.Print(elapsed)
}
