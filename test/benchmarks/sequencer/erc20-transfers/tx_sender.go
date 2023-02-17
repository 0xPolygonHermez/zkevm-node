package erc20_transfers

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	mintAmount, _  = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	transferAmount = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Div(mintAmount, big.NewInt(shared.NumberOfTxs)), big.NewInt(90)), big.NewInt(100))
	erc20SC        *ERC20.ERC20
)

func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64) error {
	log.Debugf("sending nonce: %d", nonce)
	var actualTransferAmount *big.Int
	if nonce%2 == 0 {
		actualTransferAmount = big.NewInt(0).Sub(transferAmount, big.NewInt(int64(nonce)))
	} else {
		actualTransferAmount = big.NewInt(0).Add(transferAmount, big.NewInt(int64(nonce)))
	}
	_, err := erc20SC.Transfer(shared.Auth, shared.To, actualTransferAmount)
	return err
}
