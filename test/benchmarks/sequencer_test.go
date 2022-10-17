package benchmarks

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	defaultInterval = 10 * time.Second
	defaultDeadline = 6000 * time.Second

	gasLimit = 21000
)

var dbConfig = dbutils.NewStateConfigFromEnv()

var (
	ctx                 = context.Background()
	sequencerPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	chainID             = uint64(1000)
	opsCfg              = &operations.Config{
		State: &state.Config{
			MaxCumulativeGasUsed: 800000,
		},

		Sequencer: &operations.SequencerConfig{
			Address:    "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			PrivateKey: sequencerPrivateKey,
		},
	}

	genAccBalance1, _ = new(big.Int).SetString("100000000000000000000", 10)
	genAccBalance2, _ = new(big.Int).SetString("200000000000000000000", 10)
	genAccAddr1       = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	genAccAddr2       = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"

	genesisAccounts = map[string]big.Int{
		genAccAddr1: *genAccBalance1,
		genAccAddr2: *genAccBalance2,
	}

	ethAmount, _  = big.NewInt(0).SetString("100000000000", encoding.Base10)
	privateKey, _ = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	auth, _       = bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))

	table = []struct {
		input int
	}{
		{input: 100},
		{input: 1000},
		{input: 10000},
		{input: 100000},
	}
)

func BenchmarkSequencer(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	for _, v := range table {
		st, pl, gasPrice, l2Client := setUpEnv(b)
		b.Run(fmt.Sprintf("amount_of_txs_%d", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runTxSender(b, l2Client, pl, gasPrice, v.input)
			}
		})
		tearDownEnv(b, st)
	}
}

func setUpEnv(b *testing.B) (*state.State, *pool.Pool, *big.Int, *ethclient.Client) {
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(b, err)

	st := opsman.State()
	s, err := pgpoolstorage.NewPostgresPoolStorage(dbConfig)
	require.NoError(b, err)
	pl := pool.NewPool(s, st, common.Address{}, chainID)
	// store current batch number to check later when the state is updated
	require.NoError(b, opsman.SetGenesis(genesisAccounts))
	require.NoError(b, opsman.Setup())

	// Eth client
	client, err := ethclient.Dial(l1NetworkURL)
	require.NoError(b, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)

	l2Client, err := ethclient.Dial(l2NetworkURL)
	require.NoError(b, err)

	return st, pl, gasPrice, l2Client
}

func tearDownEnv(b *testing.B, st stateInterface) {
	lastBatchNumber, err := st.GetLastL2BlockNumber(ctx, nil)
	require.NoError(b, err)
	fmt.Printf("lastBatchNumber: %v\n", lastBatchNumber)
	require.NoError(b, operations.Teardown())
}

func runTxSender(b *testing.B, l2Client *ethclient.Client, pl *pool.Pool, gasPrice *big.Int, txsAmount int) {
	var err error
	for i := 0; i < txsAmount; i++ {
		tx := types.NewTransaction(uint64(i), common.HexToAddress(genAccAddr2), ethAmount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(b, err)
		err = l2Client.SendTransaction(ctx, signedTx)
		require.NoError(b, err)
	}

	// Wait for sequencer to select txs from pool and propose a new batch
	// Wait for the synchronizer to update state
	err = operations.Poll(defaultInterval, defaultDeadline, func() (bool, error) {
		// using a closure here to capture st and currentBatchNumber
		count, err := pl.CountPendingTransactions(ctx)
		if err != nil {
			return false, err
		}

		fmt.Printf("amount of pending txs: %v\n", count)
		done := count == 0
		return done, nil
	})
	require.NoError(b, err)
}
