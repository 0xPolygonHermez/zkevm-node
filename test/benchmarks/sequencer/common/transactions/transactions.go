package transactions

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// SendAndWait sends a number of transactions and waits for them to be marked as pending in the pool
func SendAndWait(
	b *testing.B,
	senderNonce uint64,
	client *ethclient.Client,
	gasPrice *big.Int,
	pl *pool.Pool,
	ctx context.Context,
	nTxs int,
	txSenderFunc func(b *testing.B, l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64),
) {
	shared.Auth.GasPrice = gasPrice
	shared.Auth.GasLimit = 2100000
	log.Debugf("Sending %d txs ...", nTxs)
	maxNonce := uint64(nTxs) + senderNonce

	for nonce := senderNonce; nonce < maxNonce; nonce++ {
		txSenderFunc(b, client, gasPrice, nonce)
	}
	log.Debug("All txs were sent!")

	log.Debug("Waiting pending transactions To be added in the pool ...")
	err := operations.Poll(1*time.Second, shared.DefaultDeadline, func() (bool, error) {
		// using a closure here To capture st and currentBatchNumber
		count, err := pl.CountPendingTransactions(ctx)
		if err != nil {
			return false, err
		}

		log.Debugf("amount of pending txs: %d\n", count)
		done := count == uint64(nTxs)
		return done, nil
	})
	require.NoError(b, err)
	log.Debug("All pending txs are added in the pool!")
}
