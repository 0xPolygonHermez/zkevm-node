package etherman

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
)

func NewSimulatedEtherman(cfg Config, auth *bind.TransactOpts) (etherman *ClientEtherMan, commit func(), err error) {
	balance, _ := new(big.Int).SetString("10000000000000000000000000", 10) // 10 ETH in wei
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(999999999999999999)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Deploy contracts
	emptyAddr := common.Address{}
	poeAddr, _, poe, err := proofofefficiency.DeployProofofefficiency(auth, client, emptyAddr, emptyAddr, emptyAddr)
	if err != nil {
		return nil, nil, err
	}

	client.Commit()
	return &ClientEtherMan{EtherClient: client, PoE: poe, SCAddresses: []common.Address{poeAddr}, auth: auth}, client.Commit, nil
}
