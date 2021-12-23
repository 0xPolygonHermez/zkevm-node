package etherman

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}
func TestDecodeOneTxData(t *testing.T) {
	callDataTestCases := readTests()
	for _, callDataTestCase := range callDataTestCases {
		t.Run("Test id "+strconv.FormatUint(uint64(callDataTestCase.ID), 10), func(t *testing.T) {
			dHex := strings.Replace(callDataTestCase.FullCallData, "0x", "", -1)
			data, err := hex.DecodeString(dHex)
			require.NoError(t, err)
			var auxTxs []vectors.Tx
			for _, tx := range callDataTestCase.Txs {
				if tx.RawTx != "" {
					auxTxs = append(auxTxs, tx)
				}
			}
			txs, err := decodeTxs(data)
			require.NoError(t, err)
			for j := 0; j < len(txs); j++ {
				var addr common.Address
				err = addr.UnmarshalText([]byte(auxTxs[j].To))
				require.NoError(t, err)
				assert.Equal(t, &addr, txs[j].To())
			}
		})
	}
}

//This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (ethman *ClientEtherMan, commit func()) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, commit, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}

	return ethman, commit
}

func TestSCEvents(t *testing.T) {
	// Set up testing environment
	etherman, commit := newTestingEnv()

	//read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	callDataTestCases := readTests()

	//prepare txs
	dHex := strings.Replace(callDataTestCases[1].BatchL2Data, "0x", "", -1)
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	//send propose batch l1 tx
	matic := new(big.Int)
	matic, ok := matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	//prepare txs
	dHex = strings.Replace(callDataTestCases[0].BatchL2Data, "0x", "", -1)
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)

	//send propose batch l1 tx
	matic, ok = matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	//mine the tx in a block
	commit()

	//Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, err := etherman.GetBatchesByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(callDataTestCases[1].Txs[k].To))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))

	block, err = etherman.GetBatchesByBlock(ctx, block[0].BlockNumber, &block[0].BlockHash)
	require.NoError(t, err)
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(callDataTestCases[1].Txs[k].To))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))

	//VerifyBatch event. Consolidate batch event
	var (
		newLocalExitRoot = [32]byte{byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1)}
		newStateRoot     = [32]byte{byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1)}
		proofA           = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofC           = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofB           = [2][2]*big.Int{proofC, proofC}
	)
	_, err = etherman.PoE.VerifyBatch(etherman.auth, newLocalExitRoot, newStateRoot, uint32(block[0].Batches[0].BatchNumber), proofA, proofB, proofC)
	require.NoError(t, err)

	//mine the tx in a block
	commit()

	initBlock = finalBlock
	finalBlock, err = etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber = finalBlock.NumberU64()
	block, err = etherman.GetBatchesByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, block[1].Batches[0].ConsolidatedTxHash)
	assert.Equal(t, 2, len(block[0].Batches))
	assert.Equal(t, 1, len(block[1].Batches))
	log.Debugf("Batch consolidated in txHash: %+v \n", block[1].Batches[0].ConsolidatedTxHash)

	block, err = etherman.GetBatchesByBlock(ctx, finalBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, block[0].Batches[0].ConsolidatedTxHash)
	log.Debugf("Batch consolidated in txHash: %+v \n", block[0].Batches[0].ConsolidatedTxHash)
}

func TestRegisterSequencerAndEvent(t *testing.T) {
	// Set up testing environment
	etherman, commit := newTestingEnv()
	ctx := context.Background()

	//read currentBlock
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	//send propose batch l1 tx
	_, err = etherman.RegisterSequencer("http://localhost")
	require.NoError(t, err)

	//mine the tx in a block
	commit()

	//Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	finalBlockNumber := finalBlock.NumberU64()
	block, err := etherman.GetBatchesByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, etherman.auth.From, block[0].NewSequencers[0].Address)
	assert.Equal(t, "http://localhost", block[0].NewSequencers[0].URL)
	assert.Equal(t, big.NewInt(10001), block[0].NewSequencers[0].ChainID)
	log.Debug("Sequencer synced: ", block[0].NewSequencers[0].Address, ", url: ", block[0].NewSequencers[0].URL, ", and chainId: ", block[0].NewSequencers[0].ChainID)
}

func TestSCSendBatchAndVerify(t *testing.T) {
	callDataTestCases := readTests()
	dHex := strings.Replace(callDataTestCases[1].FullCallData, "0x", "", -1)
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data)
	require.NoError(t, err)

	// Set up testing environment
	etherman, commit := newTestingEnv()
	ctx := context.Background()
	matic := new(big.Int)
	matic, ok := matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	tx, err := etherman.SendBatch(ctx, txs, matic)
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())

	//mine the tx in a block
	commit()

	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, err := etherman.GetBatchesByBlockRange(ctx, finalBlockNumber, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(block[0].Batches))
	assert.Equal(t, 5, len(block[0].Batches[0].Transactions))
	assert.Equal(t, data, block[0].Batches[0].RawTxsData)

	proofSlc := []string{"0", "0"}
	proofBelem := proverclient.ProofX{Proof: proofSlc}
	var proofB []*proverclient.ProofX
	proofB = append(proofB, &proofBelem, &proofBelem)
	newStateRoot, ok := new(big.Int).SetString("1212121212121212121212121212121212121212121212121212121212121212", 16)
	assert.True(t, ok)
	newLocalExitRoot, ok := new(big.Int).SetString("1234123412341234123412341234123412341234123412341234123412341234", 16)
	assert.True(t, ok)
	proof := proverclient.Proof{
		ProofA: proofSlc,
		ProofB: proofB,
		ProofC: proofSlc,
		PublicInputs: &proverclient.PublicInputs{
			NewStateRoot:     newStateRoot.Bytes(),
			NewLocalExitRoot: newLocalExitRoot.Bytes(),
		},
	}

	tx, err = etherman.ConsolidateBatch(big.NewInt(1), &proof)
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())

	commit()

	finalBlock, err = etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber = finalBlock.NumberU64()
	block, err = etherman.GetBatchesByBlockRange(ctx, finalBlockNumber, nil)
	require.NoError(t, err)

	assert.Equal(t, 1, len(block[0].Batches))
}

func TestDefaultChainID(t *testing.T) {
	// Set up testing environment
	etherman, _ := newTestingEnv()

	//get chainID
	defaultChainID, err := etherman.GetDefaultChainID()
	require.NoError(t, err)

	//Check value
	assert.Equal(t, big.NewInt(10000), defaultChainID)
}

func readTests() []vectors.TxEventsSendBatchTestCase {
	// Load test vectors
	txEventsSendBatchTestCases, err := vectors.LoadTxEventsSendBatchTestCases("../test/vectors/smc-txevents-sendbatch-test-vector.json")
	if err != nil {
		log.Fatal(err)
	}
	return txEventsSendBatchTestCases
}
