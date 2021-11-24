package etherman

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
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
func TestDecodeOneTxData(t *testing.T) {
	dHex := "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a0758b52ed7380ef07d97a26904f6f2340e9437d3f44d4a950db48de846d18d6e5a0562ead2f0619ae253d65196b22026728829dd785c6a489cd2a546c6066d32c2a"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	tx, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	// fmt.Println(tx)
	var addr common.Address
	err = addr.UnmarshalText([]byte("0x1111111111111111111111111111111111111111"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f86b028504a817c8008252089412121212121212121212121212121212121212128506fc23ac0080820226a0a2402d3351e8ec9b0a221d7ff48aca682c646528b1558e74a0b558a943c0f2a3a05c7fff6db65e560833d7f85f07ed092fb63aaa2307b60cdf03b63a2effc6a340"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(257))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1212121212121212121212121212121212121212"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f86e5a8502540be4008259d8941234123412341234123412341234123412341234880214e8348c4f000082123442a0c6ff1e0034458c8dbf64966f49031e44c6509f85545b49d4df2a953e9f4d1324a07403e62dda1922fb1e226632e21e7382c345377ff46e6a43b79f169570e5a725"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(15))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1234123412341234123412341234123412341234"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f8701c8502540be400823a9894987698769876987698769876987698769876987688011c37937e0800008256788202e0a0ae3c16aaf6a780e085f5f919b0d1e5f07a1c014ed4d700c3ab189b4f98677f38a0316608701f846807f96f496b9994e7a821fc9322383e54f84bd235b93fb774b0"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(350))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x9876987698769876987698769876987698769876"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f873528504a817c800824a389480808080808080808080808080808080808080808802c68af0bb140000851234567890820b13a06790f151dab1b65b577532b479938f45abf871d6978e98ee0db53331ad548709a021b5038ec4cf5f7edc3f2ebac9ec5fdaa602f48b53986111a9bddaf5ad5b771c"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(1400))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x8080808080808080808080808080808080808080"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())

	dHex = "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f8702f85055ae826008261a8941111111111222222222233333333334444444444880cc47f20295c0000841122334427a01a3cf5ea05180ac59514dc64d91a1d615857cad466dacfbde4dd06f1988a1074a05ed158fd5c7d54bd827e94683c7f74dbd64be9bb69cbb3e092d317cc4758a146"
	data, err = hex.DecodeString(dHex)
	require.NoError(t, err)
	tx, err = decodeTxs(data, big.NewInt(2))
	require.NoError(t, err)
	err = addr.UnmarshalText([]byte("0x1111111111222222222233333333334444444444"))
	require.NoError(t, err)
	assert.Equal(t, &addr, tx[0].To())
}

func TestDecodeMultipleTxData(t *testing.T) {
	dHex := "06d6490f000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000147f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a08975cf0fe106a0396649d37c5292274f66193346ce5b07bcbaa8dc248f7f5496a0684963f42a662640b6d27e3552151e3f42744c9b4d72dd70b262ca94b9473c94f869028504a817c8008252089412121212121212121212121212121212121212128506fc23ac008025a0528b1dd150ccae6e83fcc44bff11928ca635f0fc6819836a14d526af1ecf0519a02a96710022671e44c81a6f19b88f605567d32dd97508aa84c830ec9d4a4aa0d2"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	txs, err := decodeTxs(data, big.NewInt(1))
	require.NoError(t, err)
	res := []string{"0x3535353535353535353535353535353535353535", "0x1111111111111111111111111111111111111111", "0x1212121212121212121212121212121212121212"}
	for k, tx := range txs {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
}

type testingEnv struct {
	transactOpts *bind.TransactOpts
	blockchain   *backends.SimulatedBackend
	poeAddr      common.Address
	poe          *proofofefficiency.Proofofefficiency
	client       *backends.SimulatedBackend
}

//This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (testingEnv, error) {
	balance := big.NewInt(0)
	balance.SetString("10000000000000000000000000", 10) // 10 ETH in wei
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return testingEnv{}, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		return testingEnv{}, err
	}

	auth.GasLimit = 99999999999
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(999999999999999999)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Deploy contracts
	EmptyAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	poeAddr, _, poe, err := proofofefficiency.DeployProofofefficiency(auth, client, EmptyAddr, EmptyAddr, EmptyAddr)
	if err != nil {
		return testingEnv{}, err
	}

	client.Commit()
	return testingEnv{
		transactOpts: auth,
		blockchain:   client,
		poeAddr:      poeAddr,
		poe:          poe,
		client:       client,
	}, nil
}

func TestSCEvents(t *testing.T) {
	// Set up testing environment
	testEnv, err := newTestingEnv()
	require.NoError(t, err)

	//read currentBlock
	initBlock := testEnv.client.Blockchain().CurrentBlock()

	//prepare txs
	dHex := "f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83f86c088504a817c8008252089411111111111111111111111111111111111111118802c68af0bb1400008026a08975cf0fe106a0396649d37c5292274f66193346ce5b07bcbaa8dc248f7f5496a0684963f42a662640b6d27e3552151e3f42744c9b4d72dd70b262ca94b9473c94f869028504a817c8008252089412121212121212121212121212121212121212128506fc23ac008025a0528b1dd150ccae6e83fcc44bff11928ca635f0fc6819836a14d526af1ecf0519a02a96710022671e44c81a6f19b88f605567d32dd97508aa84c830ec9d4a4aa0d2"
	data, err := hex.DecodeString(dHex)
	require.NoError(t, err)

	//send propose batch l1 tx
	_, err = testEnv.poe.SendBatch(testEnv.transactOpts, data, big.NewInt(2))
	require.NoError(t, err)

	//mine the tx in a block
	testEnv.client.Commit()

	//Now read the event
	conf := Config{
		PoeAddress: testEnv.poeAddr,
	}
	etherman, err := NewTestEtherman(conf, testEnv.client, testEnv.poe)
	require.NoError(t, err)
	finalBlock := testEnv.client.Blockchain().CurrentBlock()
	ctx := context.Background()
	block, err := etherman.GetBatchesFromBlockTo(ctx, initBlock.NumberU64(), finalBlock.NumberU64()+1)
	require.NoError(t, err)
	res := []string{"0x3535353535353535353535353535353535353535", "0x1111111111111111111111111111111111111111", "0x1212121212121212121212121212121212121212"}
	for k, tx := range block[0].Batches[0].Transactions {
		var addr common.Address
		err = addr.UnmarshalText([]byte(res[k]))
		require.NoError(t, err)
		assert.Equal(t, &addr, tx.To())
	}
	fmt.Printf("Block Received with %d txs\n", len(block[0].Batches[0].Transactions))
}
