package etherman

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
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
	dHex := "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000074f872b870f86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820344a0a9683074bcc81dba07fad2ac4015cf2eba4807c1aa1a8d291e77317a45fc2023a03d9ad247102178817ab2714984b4deb48bedd6ec06da0471745a81c60d1ab0b5"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	tx, err := decodeTxs(data)
	require.NoError(t, err)
	var addr common.Address
	err = addr.UnmarshalText([]byte("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000076f874b872f87080843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff8a021e19e0c9bab240000080820344a001b48d462a4d85e850d36a297d79cd78235b6aa98d0e76318a8e2c4dcd39d881a0277093ef3e6d2f2970e6386eefaac403f7575a5a7855ab0b4ccad6bc2dab081d"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data)
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x4d5Cf5032B2a844602278b01199ED191A86c93ff"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000071f86fb86df86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820344a0a9683074bcc81dba07fad2ac4015cf2eba4807c1aa1a8d291e77317a45fc2023a03d9ad247102178817ab2714984b4deb48bedd6ec06da0471745a81c60d"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data)
	require.NoError(t, err)
	assert.Equal(t, 0, len(tx))
}

func TestDecodeMultipleTxData(t *testing.T) {
	dHex := "06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000023Af90237b870f86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820344a0a9683074bcc81dba07fad2ac4015cf2eba4807c1aa1a8d291e77317a45fc2023a03d9ad247102178817ab2714984b4deb48bedd6ec06da0471745a81c60d1ab0b5b871f86f80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478893635c9adc5dea0000080820343a0d0f11c506e8606f49759222a30e767212c6f6f671c6c4073664b49cc1ff3dde2a02dd180dea53b190e6053fe9c467672c5fdb299b3b2fd5e9cf7497b52d1bcd80fb86ef86c80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478880de0b6b3a76400008025a0803fc443307ddda6c4a8e922f0a02c3a00df8de7701ab99d78954c3dc4aa7009a05407070d5186dd0a95a232f4276949e43d54f4ed120732f98776dd42ceb7cadab870f86e80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a764000080820343a0b3672951e7ad1c60799a7c5ea89eee4165b2bedce92517ca410b5efcf5471f2ba012c3894e1d00ef3ad457495e46ce38fdbbc901802ca61db9832a9a024f491e51b86ef86c80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a7640000801ba02c92035bc11227e9e94ba066a6d77a65c43d29dfb4a855c9464e1b60fabd6334a07171c6dc84816ffcf025040cd6193ecef3928a0c4e4964ddba320826b76c725d000000000000"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data)
	require.NoError(t, err)
	res := []string{"0x4d5Cf5032B2a844602278b01199ED191A86c93ff", "0x187Bd40226A7073b49163b1f6c2b73d8F2aa8478", "0x187Bd40226A7073b49163b1f6c2b73d8F2aa8478", "0xabCcEd19d7f290B84608feC510bEe872CC8F5112", "0xabCcEd19d7f290B84608feC510bEe872CC8F5112"}
	for k, tx := range txs {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}

	// Id 3
	dHex = "06d6490f0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000e7f8e5b870f86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820346a06e209c61ca92c2b980d6197e7ac9ccc3f547bf13be6455dfe682aa5dda9655efa016819a7edcc3fefec81ca97c7a6f3d10ec774440e409adbba693ce8b698d41f1b871f86f80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff89056bc75e2d6310000080820344a04d33f299985c01a539ca9f7c04c700c49cd93ab8af32d61df807b7c66bec2af7a01fec7d7cd8ac1587fdfff9da57fca66117dc86cf264cd13bf9d7660d7f7b0ef4"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err = decodeTxs(data)
	require.NoError(t, err)
	assert.Equal(t, 2, len(txs))
	res = []string{"0x4d5Cf5032B2a844602278b01199ED191A86c93ff", "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"}
	for k, tx := range txs {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
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
	auth.GasLimit = 99999999999
	ethman, commit, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}

	return ethman, commit
}

func TestSCEvents(t *testing.T) {
	// Set up testing environment
	etherman, commit := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Prepare txs
	dHex := "f90237b870f86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820344a0a9683074bcc81dba07fad2ac4015cf2eba4807c1aa1a8d291e77317a45fc2023a03d9ad247102178817ab2714984b4deb48bedd6ec06da0471745a81c60d1ab0b5b871f86f80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478893635c9adc5dea0000080820343a0d0f11c506e8606f49759222a30e767212c6f6f671c6c4073664b49cc1ff3dde2a02dd180dea53b190e6053fe9c467672c5fdb299b3b2fd5e9cf7497b52d1bcd80fb86ef86c80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478880de0b6b3a76400008025a0803fc443307ddda6c4a8e922f0a02c3a00df8de7701ab99d78954c3dc4aa7009a05407070d5186dd0a95a232f4276949e43d54f4ed120732f98776dd42ceb7cadab870f86e80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a764000080820343a0b3672951e7ad1c60799a7c5ea89eee4165b2bedce92517ca410b5efcf5471f2ba012c3894e1d00ef3ad457495e46ce38fdbbc901802ca61db9832a9a024f491e51b86ef86c80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a7640000801ba02c92035bc11227e9e94ba066a6d77a65c43d29dfb4a855c9464e1b60fabd6334a07171c6dc84816ffcf025040cd6193ecef3928a0c4e4964ddba320826b76c725d"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	// Send propose batch l1 tx
	_, err = etherman.PoE.SendBatch(etherman.auth, data, big.NewInt(2))
	require.NoError(t, err)

	// Prepare txs
	dHex = "f874b872f87080843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff8a021e19e0c9bab240000080820344a001b48d462a4d85e850d36a297d79cd78235b6aa98d0e76318a8e2c4dcd39d881a0277093ef3e6d2f2970e6386eefaac403f7575a5a7855ab0b4ccad6bc2dab081d"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)

	// Send propose batch l1 tx
	_, err = etherman.PoE.SendBatch(etherman.auth, data, big.NewInt(2))
	require.NoError(t, err)

	// Mine the tx in a block
	commit()

	// Now read the event
	finalBlock, err := etherman.EtherClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	block, err := etherman.GetBatchesByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	res := []string{"0x4d5Cf5032B2a844602278b01199ED191A86c93ff", "0x187Bd40226A7073b49163b1f6c2b73d8F2aa8478", "0x187Bd40226A7073b49163b1f6c2b73d8F2aa8478", "0xabCcEd19d7f290B84608feC510bEe872CC8F5112", "0xabCcEd19d7f290B84608feC510bEe872CC8F5112"}
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	log.Debugf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))

	block, err = etherman.GetBatchesByBlock(ctx, block[0].BlockNumber, &block[0].BlockHash)
	require.NoError(t, err)
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
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
	block, err := etherman.GetBatchesByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, etherman.auth.From, block[0].NewSequencers[0].Address)
	assert.Equal(t, "http://localhost", block[0].NewSequencers[0].URL)
	assert.Equal(t, big.NewInt(10001), block[0].NewSequencers[0].ChainID)
	log.Debug("Sequencer synced: ", block[0].NewSequencers[0].Address, ", url: ", block[0].NewSequencers[0].URL, ", and chainId: ", block[0].NewSequencers[0].ChainID)
}

func TestSCSendBatchAndVerify(t *testing.T) {
	dHex := "06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000023Af90237b870f86e80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a000080820344a0a9683074bcc81dba07fad2ac4015cf2eba4807c1aa1a8d291e77317a45fc2023a03d9ad247102178817ab2714984b4deb48bedd6ec06da0471745a81c60d1ab0b5b871f86f80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478893635c9adc5dea0000080820343a0d0f11c506e8606f49759222a30e767212c6f6f671c6c4073664b49cc1ff3dde2a02dd180dea53b190e6053fe9c467672c5fdb299b3b2fd5e9cf7497b52d1bcd80fb86ef86c80843b9aca00830186a094187bd40226a7073b49163b1f6c2b73d8f2aa8478880de0b6b3a76400008025a0803fc443307ddda6c4a8e922f0a02c3a00df8de7701ab99d78954c3dc4aa7009a05407070d5186dd0a95a232f4276949e43d54f4ed120732f98776dd42ceb7cadab870f86e80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a764000080820343a0b3672951e7ad1c60799a7c5ea89eee4165b2bedce92517ca410b5efcf5471f2ba012c3894e1d00ef3ad457495e46ce38fdbbc901802ca61db9832a9a024f491e51b86ef86c80843b9aca00830186a094abcced19d7f290b84608fec510bee872cc8f5112880de0b6b3a7640000801ba02c92035bc11227e9e94ba066a6d77a65c43d29dfb4a855c9464e1b60fabd6334a07171c6dc84816ffcf025040cd6193ecef3928a0c4e4964ddba320826b76c725d000000000000"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data)
	require.NoError(t, err)

	// Set up testing environment
	etherman, commit := newTestingEnv()
	ctx := context.Background()
	tx, err := etherman.SendBatch(ctx, txs, big.NewInt(2))
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())

	// Mine the tx in a block
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

func TestDefaultChainIDAndCollateral(t *testing.T) {
	// Set up testing environment
	etherman, _ := newTestingEnv()

	// Get chainID
	defaultChainID, err := etherman.GetDefaultChainID()
	require.NoError(t, err)

	// Check value
	assert.Equal(t, big.NewInt(10000), defaultChainID)

	// Get sequencer collateral
	collateral, err := etherman.GetSequencerCollateral()
	require.NoError(t, err)

	// Check collateral value
	assert.NotEqual(t, big.NewInt(0), collateral)
}
