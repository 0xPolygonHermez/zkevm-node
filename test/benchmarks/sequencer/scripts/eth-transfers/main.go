package eth_transfers

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/clients"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	ethtransfers "github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/eth-transfers"
)

func main() {
	var (
		err     error
		elapsed time.Duration
	)
	ctx, pl, l2Client, auth, senderNonce := clients.Init()

	// Send Txs
	err = transactions.SendAndWait(
		ctx,
		auth,
		senderNonce,
		l2Client,
		pl.CountTransactionsByStatus,
		shared.NumberOfTxs,
		ethtransfers.TxSender,
	)
	if err != nil {
		panic(err)
	}

	// Wait for Txs to be selected
	err, elapsed = transactions.WaitStatusSelected(pl.CountTransactionsByStatus, shared.NumberOfTxs)
	if err != nil {
		panic(err)
	}

	// Print results
	log.Info("##########")
	log.Info("# Result #")
	log.Info("##########")
	log.Infof("Total time took for the sequencer to select all txs from the pool: %v", elapsed)
}
