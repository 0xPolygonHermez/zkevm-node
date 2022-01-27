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
func newTestingEnv() (ethman *ClientEtherMan, commit func(), maticAddr common.Address) {
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

	//send propose batch l1 tx
	matic, ok = matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	//Read latestProposedBatch in the smc
	batchNumber, err := etherman.GetLatestProposedBatchNumber()
	require.NoError(t, err)
	assert.Equal(t, uint64(2), batchNumber)

	// Get sequencer collateral
	collateral, err := etherman.GetSequencerCollateral(2)
	require.NoError(t, err)

	// Check collateral value
	assert.NotEqual(t, big.NewInt(0), collateral)

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	for k, tx := range block[1].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(callDataTestCases[1].Txs[k].To))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[1].Batches[0].Transactions))

	block, _, err = etherman.GetRollupInfoByBlock(ctx, block[1].BlockNumber, &block[1].BlockHash)
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
	_, err = etherman.PoE.VerifyBatch(etherman.auth, newLocalExitRoot, newStateRoot, uint32(block[0].Batches[0].BatchNumber), proofA, proofB, proofC)
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
	assert.Equal(t, etherman.auth.From, block[1].NewSequencers[0].Address)
	assert.Equal(t, "http://localhost", block[1].NewSequencers[0].URL)
	assert.Equal(t, big.NewInt(1001), block[1].NewSequencers[0].ChainID)
	log.Debug("Sequencer synced: ", block[1].NewSequencers[0].Address, ", url: ", block[1].NewSequencers[0].URL, ", and chainId: ", block[1].NewSequencers[0].ChainID)
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
	txEventsSendBatchTestCases, err := vectors.LoadTxEventsSendBatchTestCases("../test/vectors/smc-txevents-sendbatch-test-vector.json")
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

	//send propose batch l1 tx
	matic, ok = matic.SetString(callDataTestCases[1].MaticAmount, 10)
	if !ok {
		log.Fatal("error decoding maticAmount")
	}
	_, err = etherman.PoE.SendBatch(etherman.auth, data, matic)
	require.NoError(t, err)

	// Register sequencer
	_, err = etherman.RegisterSequencer("http://localhost2")
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	block, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.Equal(t, NewSequencersOrder, order[block[1].BlockHash][0].Name)
	assert.Equal(t, NewSequencersOrder, order[block[1].BlockHash][1].Name)
	assert.Equal(t, BatchesOrder, order[block[1].BlockHash][2].Name)
	assert.Equal(t, NewSequencersOrder, order[block[1].BlockHash][3].Name)
	assert.Equal(t, BatchesOrder, order[block[1].BlockHash][4].Name)
	assert.Equal(t, NewSequencersOrder, order[block[1].BlockHash][5].Name)
}

func TestDepositAndGlobalExitRootEvent(t *testing.T) {
	// Set up testing environment
	etherman, commit, maticAddr := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Deposit funds
	amount := big.NewInt(9000000000000000000)
	var destNetwork uint32 = 1 // 0 is reserved to mainnet. This variable is set in the smc
	destinationAddr := common.HexToAddress("0x61A1d716a74fb45d29f148C6C20A2eccabaFD753")
	_, err = etherman.Bridge.Bridge(etherman.auth, maticAddr, amount, destNetwork, destinationAddr)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	block, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.Equal(t, DepositsOrder, order[block[1].BlockHash][0].Name)
	assert.Equal(t, GlobalExitRootsOrder, order[block[0].BlockHash][0].Name)
	assert.Equal(t, GlobalExitRootsOrder, order[block[1].BlockHash][1].Name)
	assert.Equal(t, uint64(2), block[1].BlockNumber)
	assert.Equal(t, big.NewInt(9000000000000000000), block[1].Deposits[0].Amount)
	assert.Equal(t, uint(1), block[1].Deposits[0].DestinationNetwork)
	assert.Equal(t, destinationAddr, block[1].Deposits[0].DestinationAddress)
	assert.Equal(t, 1, len(block[0].GlobalExitRoots))
	assert.Equal(t, 1, len(block[1].GlobalExitRoots))

	//Claim funds
	var (
		network  uint32
		smtProof [][32]byte
		index    uint64
	)
	mainnetExitRoot := block[1].GlobalExitRoots[0].MainnetExitRoot
	rollupExitRoot := block[1].GlobalExitRoots[0].RollupExitRoot

	_, err = etherman.Bridge.Claim(etherman.auth, maticAddr, big.NewInt(1000000000000000000), network,
		network, etherman.auth.From, smtProof, index, big.NewInt(2), mainnetExitRoot, rollupExitRoot)
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	//Read claim event
	initBlock, err = etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	block, order, err = etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), nil)
	require.NoError(t, err)
	assert.Equal(t, ClaimsOrder, order[block[0].BlockHash][0].Name)
	assert.Equal(t, big.NewInt(1000000000000000000), block[0].Claims[0].Amount)
	assert.Equal(t, uint64(3), block[0].BlockNumber)
	assert.NotEqual(t, common.Address{}, block[0].Claims[0].Token)
	assert.Equal(t, etherman.auth.From, block[0].Claims[0].DestinationAddress)
	assert.Equal(t, uint64(0), block[0].Claims[0].Index)
	assert.Equal(t, uint(0), block[0].Claims[0].OriginalNetwork)
	assert.Equal(t, uint64(3), block[0].Claims[0].BlockNumber)
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
