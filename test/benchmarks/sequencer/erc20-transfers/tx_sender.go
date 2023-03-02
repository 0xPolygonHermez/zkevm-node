package erc20_transfers

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	mintAmount     = 1000000000000000000
	transferAmount = 1
)

var (
	mintAmountBig     = big.NewInt(mintAmount)
	transferAmountBig = big.NewInt(transferAmount)
	countTxs          = 0
)

// TxSender sends ERC20 transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts, erc20SC *ERC20.ERC20) error {
	log.Debugf("sending tx num: %d nonce: %d", countTxs, nonce)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	var actualTransferAmount *big.Int
	if nonce%2 == 0 {
		actualTransferAmount = big.NewInt(0).Sub(transferAmountBig, auth.Nonce)
	} else {
		actualTransferAmount = big.NewInt(0).Add(transferAmountBig, auth.Nonce)
	}
	_, err := erc20SC.Transfer(auth, params.To, actualTransferAmount)
	countTxs += 1
	return err
}
