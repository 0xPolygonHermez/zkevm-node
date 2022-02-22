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
	"github.com/hermeznetwork/hermez-core/encoding"
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
			txs, raw, err := decodeTxs(data)
			require.NoError(t, err)
			assert.Equal(t, callDataTestCase.BatchL2Data, "0x"+hex.EncodeToString(raw))
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
func newTestingEnv() (ethman *Client, commit func(), maticAddr common.Address) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, commit, maticAddr, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}

	return ethman, commit, maticAddr
}

func TestSCEvents(t *testing.T) {
	// Set up testing environment
	etherman, commit, _ := newTestingEnv()

	// Read currentBlock
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

	matic, err = etherman.GetCurrentSequencerCollateral()
	require.NoError(t, err)
	matic.Add(matic, big.NewInt(1000000000000000000))
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	//Read latestProposedBatch in the smc
	batchNumber, err := etherman.GetLatestProposedBatchNumber()
	require.NoError(t, err)
	assert.Equal(t, uint64(2), batchNumber)

	// Get sequencer collateral
	collateral, err := etherman.GetSequencerCollateralByBatchNumber(2)
	require.NoError(t, err)

	// Check collateral value
	assert.NotEqual(t, big.NewInt(0), collateral)

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(callDataTestCases[1].Txs[k].To))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))

	block, _, err = etherman.GetRollupInfoByBlock(ctx, block[0].BlockNumber, &block[0].BlockHash)
	require.NoError(t, err)
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(callDataTestCases[1].Txs[k].To))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))

	// VerifyBatch event. Consolidate batch event
	var (
		newLocalExitRoot = [32]byte{byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1)}
		newStateRoot     = [32]byte{byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1), byte(1)}
		proofA           = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofC           = [2]*big.Int{big.NewInt(1), big.NewInt(1)}
		proofB           = [2][2]*big.Int{proofC, proofC}
	)
	_, err = etherman.PoE.VerifyBatch(etherman.auth, newLocalExitRoot, newStateRoot, block[0].Batches[0].Number().Uint64(), proofA, proofB, proofC)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	initBlock = finalBlock
	finalBlock, err = etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber = finalBlock.NumberU64()
	block, _, err = etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, block[1].Batches[0].ConsolidatedTxHash)
	assert.Equal(t, 2, len(block[0].Batches))
	assert.Equal(t, 1, len(block[1].Batches))
	log.Debugf("Batch consolidated in txHash: %+v \n", block[1].Batches[0].ConsolidatedTxHash)

	block, _, err = etherman.GetRollupInfoByBlock(ctx, finalBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, block[0].Batches[0].ConsolidatedTxHash)
	log.Debugf("Batch consolidated in txHash: %+v \n", block[0].Batches[0].ConsolidatedTxHash)
}

func TestRegisterSequencerAndEvent(t *testing.T) {
	// Set up testing environment
	etherman, commit, _ := newTestingEnv()
	ctx := context.Background()

	// Read currentBlock
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Send propose batch l1 tx
	_, err = etherman.RegisterSequencer("http://localhost")
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	finalBlockNumber := finalBlock.NumberU64()
	block, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, etherman.auth.From, block[0].NewSequencers[0].Address)
	assert.Equal(t, "http://localhost", block[0].NewSequencers[0].URL)
	assert.Equal(t, big.NewInt(1001), block[0].NewSequencers[0].ChainID)
	log.Debug("Sequencer synced: ", block[0].NewSequencers[0].Address, ", url: ", block[0].NewSequencers[0].URL, ", and chainId: ", block[0].NewSequencers[0].ChainID)
}

func TestSCSendBatchAndVerify(t *testing.T) {
	callDataTestCases := readTests()
	dHex := strings.Replace(callDataTestCases[1].FullCallData, "0x", "", -1)
	txRaw, err := hex.DecodeString(strings.Replace(callDataTestCases[1].BatchL2Data, "0x", "", -1))
	require.NoError(t, err)
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, _, err := decodeTxs(data)
	require.NoError(t, err)

	// Set up testing environment
	etherman, commit, _ := newTestingEnv()
	ctx := context.Background()
	matic := new(big.Int)
	matic, ok := matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	tx, err := etherman.SendBatch(ctx, txs, matic)
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())

	// Mine the tx in a block
	commit()

	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, _, err := etherman.GetRollupInfoByBlockRange(ctx, finalBlockNumber, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(block[0].Batches))
	assert.Equal(t, 5, len(block[0].Batches[0].Transactions))
	assert.Equal(t, txRaw, block[0].Batches[0].RawTxsData)

	proofSlc := []string{"21565562996339338849564494593510043939919477777537254861096838135755987477204", "15747187504567161371139699773055768026064826386340198020294048673214119551381", "1"}
	proofBelem := proverclient.ProofX{Proof: proofSlc}
	var proofB []*proverclient.ProofX
	proofB = append(proofB, &proofBelem, &proofBelem, &proofBelem)
	newStateRoot, ok := new(big.Int).SetString("1212121212121212121212121212121212121212121212121212121212121212", 16)
	assert.True(t, ok)
	newLocalExitRoot, ok := new(big.Int).SetString("1234123412341234123412341234123412341234123412341234123412341234", 16)
	assert.True(t, ok)
	publicInputs := &proverclient.PublicInputs{
		NewStateRoot:     newStateRoot.String(),
		NewLocalExitRoot: newLocalExitRoot.String(),
	}
	proof := proverclient.Proof{
		ProofA: proofSlc,
		ProofB: proofB,
		ProofC: proofSlc,
		PublicInputsExtended: &proverclient.PublicInputsExtended{
			PublicInputs: publicInputs,
		},
	}

	tx, err = etherman.ConsolidateBatch(big.NewInt(1), &proof)
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())

	commit()

	finalBlock, err = etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber = finalBlock.NumberU64()
	block, _, err = etherman.GetRollupInfoByBlockRange(ctx, finalBlockNumber, nil)
	require.NoError(t, err)

	assert.Equal(t, 1, len(block[0].Batches))
}

func TestDefaultChainID(t *testing.T) {
	// Set up testing environment
	etherman, _, _ := newTestingEnv()

	// Get chainID
	defaultChainID, err := etherman.GetDefaultChainID()
	require.NoError(t, err)

	// Check value
	assert.Equal(t, big.NewInt(1000), defaultChainID)
}

func TestCustomChainID(t *testing.T) {
	// Set up testing environment
	etherman, commit, _ := newTestingEnv()

	// Register sequencer
	_, err := etherman.RegisterSequencer("http://localhost")
	require.NoError(t, err)

	commit()

	// Get chainID
	customChainID, err := etherman.GetCustomChainID()
	require.NoError(t, err)

	// Check value
	assert.Equal(t, big.NewInt(1001), customChainID)
}

func readTests() []vectors.TxEventsSendBatchTestCase {
	// Load test vectors
	txEventsSendBatchTestCases, err := vectors.LoadTxEventsSendBatchTestCases("../test/vectors/src/tools/calldata-test-vectors/calldata-test-vector.json")
	if err != nil {
		log.Fatal(err)
	}
	return txEventsSendBatchTestCases
}

func TestOrderReadEvent(t *testing.T) {
	log.Debug("Testing sync order...")

	// Set up testing environment
	etherman, commit, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	callDataTestCases := readTests()

	// Register sequencer
	_, err = etherman.RegisterSequencer("http://localhost")
	require.NoError(t, err)

	// Register sequencer
	_, err = etherman.RegisterSequencer("http://localhost0")
	require.NoError(t, err)

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

	// Register sequencer
	_, err = etherman.RegisterSequencer("http://localhost1")
	require.NoError(t, err)

	//prepare txs
	dHex = strings.Replace(callDataTestCases[0].BatchL2Data, "0x", "", -1)
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)

	matic, err = etherman.GetCurrentSequencerCollateral()
	require.NoError(t, err)
	matic.Add(matic, big.NewInt(1000000000000000000))
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	// Register sequencer
	_, err = etherman.RegisterSequencer("http://localhost2")
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	block, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.Equal(t, NewSequencersOrder, order[block[0].BlockHash][0].Name)
	assert.Equal(t, NewSequencersOrder, order[block[0].BlockHash][1].Name)
	assert.Equal(t, BatchesOrder, order[block[0].BlockHash][2].Name)
	assert.Equal(t, NewSequencersOrder, order[block[0].BlockHash][3].Name)
	assert.Equal(t, BatchesOrder, order[block[0].BlockHash][4].Name)
	assert.Equal(t, NewSequencersOrder, order[block[0].BlockHash][5].Name)
}

func TestConverter(t *testing.T) {
	str := "0x53793751b374bcde3992cb106847589b56c34765e728d55e0ca0afe991b6c16f"
	res, err := stringToFixedByteArray(str)
	if err != nil {
		log.Error(err)
		require.NoError(t, err)
	}
	finalString := hex.EncodeToString(res[:])
	assert.Equal(t, str, "0x"+finalString)
}

func TestDecode(t *testing.T) {
	dHex := strings.Replace("0x06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000001c9f90185808502540be400832dc6c08080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220a92848e725904714428852952fc1c595af6ecdf7f8cfcc7c961cba7c7208044d64736f6c634300080700338203e980803f34a8855378618502c371823a1a6a5d244d9681d6aaf9b35338925ac875c64a4a613f018d4d5842364afdb4a9c1e123983448d82c6d857aba49dd0495bcb2bd1b0000000000000000000000000000000000000000000000", "0x", "", -1)
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)
	txs, _, err := decodeTxs(data)
	require.NoError(t, err)
	v, r, s := txs[0].RawSignatureValues()
	assert.Equal(t, "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220a92848e725904714428852952fc1c595af6ecdf7f8cfcc7c961cba7c7208044d64736f6c63430008070033", hex.EncodeToString(txs[0].Data()))
	vs, _ := new(big.Int).SetString("2037", encoding.Base10)
	rs, _ := new(big.Int).SetString("28588748595963868051531601512753704808665372676636271983291863491367365101130", encoding.Base10)
	ss, _ := new(big.Int).SetString("33642969812305861406427680625077470339025571471898670597741777160820112536253", encoding.Base10)
	assert.Equal(t, vs, v)
	assert.Equal(t, rs, r)
	assert.Equal(t, ss, s)
}

func TestDecodeSmcInteraction(t *testing.T) {
	dHex := strings.Replace("0x06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000001bc16d674ec80000000000000000000000000000000000000000000000000000000000000000008cf84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c0000000000000000000000000000000000000000", "0x", "", -1)
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)
	txs, _, err := decodeTxs(data)
	require.NoError(t, err)
	v, r, s := txs[0].RawSignatureValues()
	assert.Equal(t, "6057361d0000000000000000000000000000000000000000000000000000000000000004", hex.EncodeToString(txs[0].Data()))
	vs, _ := new(big.Int).SetString("2038", encoding.Base10)
	rs, _ := new(big.Int).SetString("52439813673236985820307219467235173403931944419246794350953690522214729109540", encoding.Base10)
	ss, _ := new(big.Int).SetString("21636541153156884113230878396581781324666928638699850950124295339798334610786", encoding.Base10)
	assert.Equal(t, vs, v)
	assert.Equal(t, rs, r)
	assert.Equal(t, ss, s)
	assert.Equal(t, uint64(1), txs[0].Nonce())
	assert.Equal(t, uint64(1001), txs[0].ChainId().Uint64())
}
