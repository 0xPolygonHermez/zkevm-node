package setup

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const sleepDuration = 5 * time.Second

// Environment sets up the environment for the benchmark
func Environment(ctx context.Context, b *testing.B) (*operations.Manager, *ethclient.Client, *pool.Pool, uint64, *big.Int) {
	if testing.Short() {
		b.Skip()
	}

	err := operations.Teardown()
	require.NoError(b, err)

	shared.OpsCfg.State.MaxCumulativeGasUsed = shared.MaxCumulativeGasUsed
	opsman, err := operations.NewManager(ctx, shared.OpsCfg)
	require.NoError(b, err)

	err = Components(opsman)
	require.NoError(b, err)
	time.Sleep(sleepDuration)

	// Load account with balance on local genesis
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(b, err)

	// Load common client
	client, err := ethclient.Dial(operations.DefaultL2NetworkURL)
	require.NoError(b, err)

	st := opsman.State()
	s, err := pgpoolstorage.NewPostgresPoolStorage(shared.PoolDbConfig)
	require.NoError(b, err)
	pl := pool.NewPool(s, st, common.Address{}, shared.ChainID)

	// Print Info before send
	senderBalance, err := client.BalanceAt(ctx, auth.From, nil)
	require.NoError(b, err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	require.NoError(b, err)

	// Print Initial Stats
	log.Infof("Receiver Addr: %v", shared.To.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)

	return opsman, client, pl, senderNonce, gasPrice
}

// Components runs the network container, starts synchronizer and JSON-RPC components, and approves matic
func Components(opsman *operations.Manager) error {
	// Run network container
	err := opsman.StartNetwork()
	if err != nil {
		return err
	}

	// Approve matic
	err = operations.ApproveMatic()
	if err != nil {
		return err
	}

	err = operations.StartComponent("sync")
	if err != nil {
		return err
	}

	err = operations.StartComponent("json-rpc")
	if err != nil {
		return err
	}
	time.Sleep(sleepDuration)

	return nil
}

// BootstrapSequencer starts the sequencer and waits for it to be ready
func BootstrapSequencer(b *testing.B, opsman *operations.Manager) {
	log.Debug("Starting sequencer ....")
	err := operations.StartComponent("seq")
	require.NoError(b, err)
	log.Debug("Sequencer Started!")
	log.Debug("Setup sequencer ....")
	require.NoError(b, opsman.SetUpSequencer())
	log.Debug("Sequencer Setup ready!")
}
