package transactions

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendAndWait sends a number of transactions and waits for them to be marked as pending in the pool
func SendAndWait(
	ctx context.Context,
	auth *bind.TransactOpts,
	senderNonce uint64,
	client *ethclient.Client,
	countByStatusFunc func(ctx context.Context, status pool.TxStatus) (uint64, error),
	nTxs int,
	txSenderFunc func(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) error,
) error {
	auth.GasLimit = 2100000
	log.Debugf("Sending %d txs ...", nTxs)
	maxNonce := uint64(nTxs) + senderNonce

	for nonce := senderNonce; nonce < maxNonce; nonce++ {
		err := txSenderFunc(client, auth.GasPrice, nonce)
		if err != nil {
			return err
		}
	}
	log.Debug("All txs were sent!")

	log.Debug("Waiting pending transactions To be added in the pool ...")
	err := operations.Poll(1*time.Second, shared.DefaultDeadline, func() (bool, error) {
		// using a closure here To capture st and currentBatchNumber
		count, err := countByStatusFunc(ctx, pool.TxStatusPending)
		if err != nil {
			return false, err
		}

		log.Debugf("amount of pending txs: %d\n", count)
		done := count == uint64(nTxs)
		return done, nil
	})
	if err != nil {
		return err
	}

	log.Debug("All pending txs are added in the pool!")

	return nil
}

func WaitStatusSelected(countByStatusFunc func(ctx context.Context, status pool.TxStatus) (uint64, error), nTxs uint64) (error, time.Duration) {
	start := time.Now()
	log.Debug("Wait for sequencer to select all txs from the pool")
	err := operations.Poll(200*time.Millisecond, shared.DefaultDeadline, func() (bool, error) {
		selectedCount, err := countByStatusFunc(shared.Ctx, pool.TxStatusSelected)
		if err != nil {
			return false, err
		}

		log.Debugf("amount of selected txs: %d", selectedCount)
		done := selectedCount >= nTxs
		return done, nil
	})
	elapsed := time.Since(start)
	return err, elapsed
}
