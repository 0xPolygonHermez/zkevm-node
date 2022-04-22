package e2e

import (
	"context"
	"flag"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/ERC20"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// TestStateTransition tests state transitions using the vector
func TestUniswap(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()

	ctx := context.Background()

	// load the same config that will be used by the test
	app := cli.NewApp()
	var n string
	flag.StringVar(&n, "network", "local", "")
	cfg, err := config.Load(cli.NewContext(app, flag.CommandLine, nil))
	require.NoError(t, err)

	opsCfg := &operations.Config{
		Arity: cfg.NetworkConfig.Arity,
		State: &state.Config{
			DefaultChainID:                cfg.NetworkConfig.L2DefaultChainID,
			MaxCumulativeGasUsed:          cfg.NetworkConfig.MaxCumulativeGasUsed,
			GlobalExitRootStoragePosition: cfg.NetworkConfig.GlobalExitRootStoragePosition,
			LocalExitRootStoragePosition:  cfg.NetworkConfig.LocalExitRootStoragePosition,
			L2GlobalExitRootManagerAddr:   cfg.NetworkConfig.L2GlobalExitRootManagerAddr,
		},
	}
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)

	require.NoError(t, opsman.StartNetwork())
	require.NoError(t, opsman.StartProver())

	g := errgroup.Group{}
	g.Go(func() error {
		return opsman.StartCore()
	})
	g.Go(func() error {
		return opsman.InitNetwork()
	})
	err = g.Wait()
	require.NoError(t, err)

	client, err := ethclient.Dial("http://localhost:8123")
	require.NoError(t, err)
	accountAddr := common.HexToAddress("0xC949254d682D8c9ad5682521675b8F43b102aec4")

	balance, err := client.BalanceAt(ctx, accountAddr, nil)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), "10000000000000000000", "invalid ETH Balance")

	require.NoError(t, opsman.DeployUniswap())

	aCoinAddr := common.HexToAddress("0x3A07588DefB088956a2e6dD15C33d63F2E0A2c55")
	bCoinAddr := common.HexToAddress("0x0ef3B0bC8D6313aB7dc03CF7225c872071bE1E6d")
	cCoinAddr := common.HexToAddress("0xd59D09BBEE914015562D95e84a78f1CD4FC347E9")

	aCoin, err := ERC20.NewERC20(aCoinAddr, client)
	require.NoError(t, err)
	balance, err = aCoin.BalanceOf(nil, accountAddr)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), "989000000000000000000", "invalid A Coin Balance")

	bCoin, err := ERC20.NewERC20(bCoinAddr, client)
	require.NoError(t, err)
	balance, err = bCoin.BalanceOf(nil, accountAddr)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), "979906610893880149131", "invalid B Coin Balance")

	cCoin, err := ERC20.NewERC20(cCoinAddr, client)
	require.NoError(t, err)
	balance, err = cCoin.BalanceOf(nil, accountAddr)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), "990906610893880149131", "invalid C Coin Balance")

	balance, err = client.BalanceAt(ctx, accountAddr, nil)
	require.NoError(t, err)
	assert.Equal(t, balance.String(), "10000000000000000000", "invalid ETH Balance after deployments and swaps")

	require.NoError(t, operations.Teardown())
}
