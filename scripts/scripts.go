package scripts

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/log"
)

// WaitTxToBeMined waits the tx receipt
func WaitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	log.Infof("waiting tx to be mined: %v", hash.Hex())
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return nil, errors.New("timeout exceed")
		}

		r, err := client.TransactionReceipt(context.Background(), hash)
		if errors.Is(err, ethereum.NotFound) {
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Errorf("Failed to get tx receipt: %v", err)
			return nil, err
		}

		if r.Status == types.ReceiptStatusFailed {
			log.Errorf("tx mined[FAILED]: %v", hash.Hex())
		} else {
			log.Infof("tx mined[SUCCESS]: %v", hash.Hex())
		}

		return r, nil
	}
}
