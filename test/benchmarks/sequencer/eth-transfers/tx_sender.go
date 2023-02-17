package eth_transfers

import (
	"errors"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit     = uint64(21000)
	ethAmount, _ = big.NewInt(0).SetString("100000000000", encoding.Base10)
)

func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) error {
	log.Debugf("sending nonce: %d", nonce)
	tx := types.NewTransaction(nonce, shared.To, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := shared.Auth.Signer(shared.Auth.From, tx)
	if err != nil {
		return err
	}

	err = l2Client.SendTransaction(shared.Ctx, signedTx)
	if errors.Is(err, state.ErrStateNotSynchronized) {
		for errors.Is(err, state.ErrStateNotSynchronized) {
			time.Sleep(5 * time.Second)
			err = l2Client.SendTransaction(shared.Ctx, signedTx)
		}
	}
	return err
}
