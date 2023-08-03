package uniswap_transfers

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	uniswap "github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit  = 21000
	sleepTime = 5 * time.Second
	countTxs  = 0
	txTimeout = 60 * time.Second
)

// TxSender sends eth transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, nonce uint64, auth *bind.TransactOpts, erc20SC *ERC20.ERC20, uniswapDeployments *uniswap.Deployments) ([]*types.Transaction, error) {
	msg := fmt.Sprintf("# swap cycle number: %d #", countTxs)
	delimiter := strings.Repeat("#", len(msg))
	log.Infof("%s\n%s\n%s", delimiter, msg, delimiter)
	var err error

	transactions := uniswap.SwapTokens(l2Client, auth, *uniswapDeployments)
	if errors.Is(err, state.ErrStateNotSynchronized) || errors.Is(err, state.ErrInsufficientFunds) {
		for errors.Is(err, state.ErrStateNotSynchronized) || errors.Is(err, state.ErrInsufficientFunds) {
			time.Sleep(sleepTime)
			transactions = uniswap.SwapTokens(l2Client, auth, *uniswapDeployments)
		}
	}

	if err == nil {
		countTxs += 1
	}

	return transactions, err
}
