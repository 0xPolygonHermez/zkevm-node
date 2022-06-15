package ethermanv2

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/ethermanv2/smartcontracts/bridge"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

//This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (ethman *Client, commit func(), maticAddr common.Address, br *bridge.Bridge) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, commit, maticAddr, br, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}
	return ethman, commit, maticAddr, br
}

func TestGEREvent(t *testing.T) {
	// Set up testing environment
	etherman, commit, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount := big.NewInt(1000000000000000)
	a := etherman.auth
	a.Value = amount
	_, err = br.Bridge(a, common.Address{}, amount, 1, etherman.auth.From)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)

	assert.Equal(t, uint64(2), blocks[0].GlobalExitRoots[0].BlockNumber)
	assert.Equal(t, big.NewInt(1), blocks[0].GlobalExitRoots[0].GlobalExitRootNum)
	assert.NotEqual(t, common.Hash{}, blocks[0].GlobalExitRoots[0].MainnetExitRoot)
	assert.Equal(t, common.Hash{}, blocks[0].GlobalExitRoots[0].RollupExitRoot)
}

func TestForcedBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, commit, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount, err := etherman.PoE.CalculateForceProverFee(&bind.CallOpts{Pending: false})
	require.NoError(t, err)
	_, err = etherman.PoE.ForceBatch(etherman.auth, []byte{}, amount)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), blocks[0].ForcedBatches[0].BlockNumber)
	assert.NotEqual(t, common.Hash{}, blocks[0].ForcedBatches[0].GlobalExitRoot)
	assert.NotEqual(t, time.Time{}, blocks[0].ForcedBatches[0].ForcedAt)
	assert.Equal(t, uint64(1), blocks[0].ForcedBatches[0].ForcedBatchNumber)
	assert.Equal(t, []byte{}, blocks[0].ForcedBatches[0].RawTxsData)
	assert.Equal(t, etherman.auth.From, blocks[0].ForcedBatches[0].Sequencer)
}
