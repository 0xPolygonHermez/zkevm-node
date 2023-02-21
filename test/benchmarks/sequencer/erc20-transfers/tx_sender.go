package erc20_transfers

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	mintAmount, _               = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	maxPercentage               = 100
	percentageOfMintAmountToUse = 10
	transferAmount              = big.NewInt(0).Div(big.NewInt(0).Mul(
		big.NewInt(0).Div(mintAmount, big.NewInt(params.NumberOfTxs)), // 1/10 of the minted amount
		big.NewInt(int64(percentageOfMintAmountToUse)),                // 10% of the minted amount
	), big.NewInt(int64(maxPercentage))) // Divide by 100 to get the percentage
	erc20SC *ERC20.ERC20
)

// TxSender sends ERC20 transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts) error {
	log.Debugf("sending nonce: %d", nonce)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	var actualTransferAmount *big.Int
	if nonce%2 == 0 {
		actualTransferAmount = big.NewInt(0).Sub(transferAmount, auth.Nonce)
	} else {
		actualTransferAmount = big.NewInt(0).Add(transferAmount, auth.Nonce)
	}
	_, err := erc20SC.Transfer(auth, params.To, actualTransferAmount)
	return err
}
