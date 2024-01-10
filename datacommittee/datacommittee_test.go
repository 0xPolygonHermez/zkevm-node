package datacommittee

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/cdkdatacommittee"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (
	dacman *DataCommitteeMan,
	ethBackend *backends.SimulatedBackend,
	auth *bind.TransactOpts,
	da *cdkdatacommittee.Cdkdatacommittee,
) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	dacman, ethBackend, da, err = newSimulatedDacman(auth)
	if err != nil {
		log.Fatal(err)
	}
	return dacman, ethBackend, auth, da
}

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func newSimulatedDacman(auth *bind.TransactOpts) (
	dacman *DataCommitteeMan,
	ethBackend *backends.SimulatedBackend,
	da *cdkdatacommittee.Cdkdatacommittee,
	err error,
) {
	if auth == nil {
		// read only client
		return &DataCommitteeMan{}, nil, nil, nil
	}
	// 10000000 ETH in wei
	balance, _ := new(big.Int).SetString("10000000000000000000000000", 10) //nolint:gomnd
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(999999999999999999) //nolint:gomnd
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// DAC Setup
	_, _, da, err = cdkdatacommittee.DeployCdkdatacommittee(auth, client)
	if err != nil {
		return &DataCommitteeMan{}, nil, nil, nil
	}
	_, err = da.Initialize(auth)
	if err != nil {
		return &DataCommitteeMan{}, nil, nil, nil
	}
	_, err = da.SetupCommittee(auth, big.NewInt(0), []string{}, []byte{})
	if err != nil {
		return &DataCommitteeMan{}, nil, nil, nil
	}
	client.Commit()

	c := &DataCommitteeMan{
		DataCommitteeContract: da,
	}
	return c, client, da, nil
}
