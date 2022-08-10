package e2e

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultArity                = 4
	defaultSequencerAddress     = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	defaultSequencerPrivateKey  = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	defaultSequencerBalance     = 400000
	defaultMaxCumulativeGasUsed = 800000

	defaultL1URL            = "http://localhost:8545"
	defaultL1ChainID uint64 = 1337
	defaultL2URL            = "http://localhost:8123"
	defaultL2ChainID uint64 = 1000

	defaultTimeoutTxToBeMined = 1 * time.Minute
)

func getDefaultOperationsConfig() *operations.Config {
	return &operations.Config{
		Arity: defaultArity, State: &state.Config{MaxCumulativeGasUsed: defaultMaxCumulativeGasUsed},
		Sequencer: &operations.SequencerConfig{Address: defaultSequencerAddress, PrivateKey: defaultSequencerPrivateKey},
	}
}

func getClients() (*ethclient.Client, *ethclient.Client, error) {
	l1Client, err := ethclient.Dial(defaultL1URL)
	if err != nil {
		return nil, nil, err
	}

	l2Client, err := ethclient.Dial(defaultL2URL)
	if err != nil {
		return nil, nil, err
	}

	return l1Client, l2Client, nil
}

func getAuth() (*bind.TransactOpts, *bind.TransactOpts, error) {
	chainID := big.NewInt(0).SetUint64(defaultL1ChainID)
	l1Auth, err := operations.GetAuth(defaultSequencerPrivateKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	chainID = big.NewInt(0).SetUint64(defaultL2ChainID)
	l2Auth, err := operations.GetAuth(defaultSequencerPrivateKey, chainID)
	if err != nil {
		return nil, nil, err
	}

	return l1Auth, l2Auth, nil
}
