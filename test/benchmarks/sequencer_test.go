package benchmarks

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
	"time"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress        = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	bridgeAddress     = "0xffffffffffffffffffffffffffffffffffffffff"
	maticTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	defaultInterval        = 5 * time.Second
	defaultDeadline        = 300 * time.Second
	defaultTxMinedDeadline = 5 * time.Second

	makeCmd = "make"
	cmdDir  = "../.."
)

var dbConfig = dbutils.NewConfigFromEnv()

var (
	ctx                 = context.Background()
	sequencerPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	chainID             = uint64(400)
	opsCfg              = &operations.Config{
		Arity: 4,
		State: &state.Config{
			DefaultChainID:       1000,
			MaxCumulativeGasUsed: 800000,
		},

		Sequencer: &operations.SequencerConfig{
			Address:    "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			PrivateKey: sequencerPrivateKey,
			ChainID:    chainID,
		},
	}

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
	for _, v := range table {
		b.Run(fmt.Sprintf("amount_of_txs_%d", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runTxSender(b, v.input)
			}
		})
	}

}

func runTxSender(b *testing.B, txsAmount int) {
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(b, err)

	st := opsman.State()
	pl, err := pool.NewPostgresPool(dbConfig)
	require.NoError(b, err)
	// store current batch number to check later when the state is updated
	currentBatchNumber, err := st.GetLastBatchNumberSeenOnEthereum(ctx)
	require.NoError(b, err)

	genAccBalance1, _ := new(big.Int).SetString("100000000000000000000", 10)
	genAccBalance2, _ := new(big.Int).SetString("200000000000000000000", 10)
	genAccAddr1 := "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	genAccAddr2 := "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
	genesisAccounts := map[string]big.Int{
		genAccAddr1: *genAccBalance1,
		genAccAddr2: *genAccBalance2,
	}
	require.NoError(b, opsman.SetGenesis(genesisAccounts))

	require.NoError(b, opsman.Setup())

	ethAmount, _ := big.NewInt(0).SetString("100000000000", encoding.Base10)

	// Eth client
	fmt.Println("Connecting to l1")
	client, err := ethclient.Dial(l1NetworkURL)
	require.NoError(b, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(b, err)
	const gasLimit = 21000

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	require.NoError(b, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
	l2Client, err := ethclient.Dial(l2NetworkURL)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < txsAmount; i++ {
		tx := types.NewTransaction(uint64(i), common.HexToAddress(genAccAddr2), ethAmount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(b, err)
		err = l2Client.SendTransaction(ctx, signedTx)
		require.NoError(b, err)
	}

	// Wait for sequencer to select txs from pool and propose a new batch
	// Wait for the synchronizer to update state
	err = operations.WaitPoll(defaultInterval, defaultDeadline, func() (bool, error) {
		// using a closure here to capture st and currentBatchNumber
		txs, err := pl.GetPendingTxs(ctx)
		if err != nil {
			return false, err
		}
		fmt.Println("txs")
		fmt.Println(len(txs))

		latestBatchNumber, err := st.GetLastBatchNumber(ctx)
		if err != nil {
			return false, err
		}
		done := len(txs) == 0 && latestBatchNumber > currentBatchNumber
		return done, nil
	})

	fmt.Println(st.GetLastBatchNumber(ctx))
	require.NoError(b, operations.Teardown())
}
