package erc20_transfers

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	uniswap "github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	mintAmount = 1000000000000000000
)

var (
	mintAmountBig = big.NewInt(mintAmount)
	countTxs      = 0
)

// TxSender sends ERC20 transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, auth *bind.TransactOpts, erc20SC *ERC20.ERC20, uniswapDeployments *uniswap.Deployments) ([]*types.Transaction, error) {
	fmt.Printf("sending tx num: %d\n", countTxs+1)
	var actualTransferAmount *big.Int
	if countTxs%2 == 0 {
		actualTransferAmount = big.NewInt(0)
	} else {
		actualTransferAmount = big.NewInt(1)
	}
	tx, err := erc20SC.Transfer(auth, params.To, actualTransferAmount)
	if err == nil {
		countTxs += 1
	}

	return []*types.Transaction{tx}, err
}
