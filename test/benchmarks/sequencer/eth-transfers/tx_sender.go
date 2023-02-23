package eth_transfers

import (
	"errors"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit     = 21000
	ethAmount, _ = big.NewInt(0).SetString("1", encoding.Base10)
	sleepTime    = 5 * time.Second
)

// TxSender sends eth transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts) error {
	log.Debugf("sending nonce: %d", nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	tx := types.NewTransaction(nonce, params.To, ethAmount, uint64(gasLimit), gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	err = l2Client.SendTransaction(params.Ctx, signedTx)
	if errors.Is(err, state.ErrStateNotSynchronized) {
		for errors.Is(err, state.ErrStateNotSynchronized) {
			time.Sleep(sleepTime)
			err = l2Client.SendTransaction(params.Ctx, signedTx)
		}
	}

	return err
}
