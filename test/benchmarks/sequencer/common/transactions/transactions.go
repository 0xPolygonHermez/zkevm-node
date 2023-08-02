package transactions

import (
	"context"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendAndWait sends a number of transactions and waits for them to be marked as pending in the pool
func SendAndWait(auth *bind.TransactOpts, client *ethclient.Client, getTxsByStatus func(ctx context.Context, status pool.TxStatus, limit uint64) ([]pool.Transaction, error), nTxs int, erc20SC *ERC20.ERC20, uniswapDeployments *pkg.Deployments, txSenderFunc func(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts, erc20SC *ERC20.ERC20, uniswapDeployments *pkg.Deployments) error) error {
	auth.GasLimit = 2100000
	log.Debugf("Sending %d txs ...", nTxs)
	startingNonce := uint64(0)
	if auth.Nonce != nil {
		startingNonce = auth.Nonce.Uint64()
	}
	maxNonce := uint64(nTxs) + startingNonce
	IP := getPublicIP()

	for nonce := startingNonce; nonce < maxNonce; nonce++ {
		err := txSenderFunc(client, auth.GasPrice, nonce, auth, erc20SC, uniswapDeployments)
		if err != nil {
			for err != nil && err.Error() == "nonce intrinsic error" {
				log.Warnf("nonce intrinsic error, retrying with nonce %d", nonce)
				err = txSenderFunc(client, auth.GasPrice, nonce, auth, erc20SC, uniswapDeployments)
			}
			if err == nil {
				continue
			}
			return err
		}
	}
	log.Debug("All txs were sent!")
	log.Debug("Waiting pending transactions To be added in the pool ...")
	err := operations.Poll(1*time.Second, params.DefaultDeadline, func() (bool, error) {
		// using a closure here To capture st and currentBatchNumber
		pendingTxs, err := getTxsByStatus(params.Ctx, pool.TxStatusPending, 0)
		if err != nil {
			panic(err)
		}
		pendingTxsCount := 0
		for _, tx := range pendingTxs {
			if tx.IP == IP {
				pendingTxsCount++
			}
		}

		log.Debugf("amount of pending txs: %d\n", pendingTxsCount)
		done := pendingTxsCount == 0
		return done, nil
	})
	if err != nil {
		return err
	}

	log.Debug("All pending txs are added in the pool!")

	return nil
}

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(ip)
}

// WaitStatusSelected waits for a number of transactions to be marked as selected in the pool
func WaitStatusSelected(countByStatusFunc func(ctx context.Context, status ...pool.TxStatus) (uint64, error), initialCount uint64, nTxs uint64) error {
	log.Debug("Wait for sequencer to select all txs from the pool")
	pollingInterval := 1 * time.Second

	prevCount := uint64(0)
	txsPerSecond := 0
	txsPerSecondAsStr := "N/A"
	estimatedTimeToFinish := "N/A"
	err := operations.Poll(pollingInterval, params.DefaultDeadline, func() (bool, error) {
		selectedCount, err := countByStatusFunc(params.Ctx, pool.TxStatusSelected)
		if err != nil {
			return false, err
		}
		currCount := selectedCount - initialCount
		remainingTxs := nTxs - currCount
		if prevCount > 0 {
			txsPerSecond = int(currCount - prevCount)
			if txsPerSecond == 0 {
				estimatedTimeToFinish = "N/A"
			} else {
				estimatedTimeToFinish = (time.Duration(int(remainingTxs)/txsPerSecond) * time.Second).String()
			}
			txsPerSecondAsStr = strconv.Itoa(txsPerSecond)
		}
		log.Debugf("amount of selected txs: %d/%d, estimated txs per second: %s, time to finish: %s", selectedCount-initialCount, nTxs, txsPerSecondAsStr, estimatedTimeToFinish)
		prevCount = currCount

		done := (int64(selectedCount) - int64(initialCount)) >= int64(nTxs)
		return done, nil
	})

	return err
}
