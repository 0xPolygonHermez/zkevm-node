package eth_transfers

import (
	"errors"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit  = 21000
	ethAmount = big.NewInt(0)
	sleepTime = 5 * time.Second
	countTxs  = 0
)

// TxSender sends eth transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts, erc20SC *ERC20.ERC20) error {
	log.Debugf("sending tx num: %d nonce: %d", countTxs, nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	tx := types.NewTransaction(nonce, params.To, ethAmount, uint64(gasLimit), gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	err = l2Client.SendTransaction(params.Ctx, signedTx)
	if errors.Is(err, state.ErrStateNotSynchronized) || errors.Is(err, state.ErrInsufficientFunds) {
		for errors.Is(err, state.ErrStateNotSynchronized) || errors.Is(err, state.ErrInsufficientFunds) {
			time.Sleep(sleepTime)
			err = l2Client.SendTransaction(params.Ctx, signedTx)
		}
	}

	if err == nil {
		countTxs += 1
	}

	return err
}
