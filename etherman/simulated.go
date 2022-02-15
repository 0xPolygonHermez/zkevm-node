package etherman

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/bridge"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/matic"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
)

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func NewSimulatedEtherman(cfg Config, auth *bind.TransactOpts) (etherman *Client, commit func(), maticAddr common.Address, err error) {
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

	// Deploy contracts
	totalSupply, _ := new(big.Int).SetString("10000000000000000000000000000", 10) //nolint:gomnd
	maticAddr, _, maticContract, err := matic.DeployMatic(auth, client, "Matic Token", "MATIC", 18, totalSupply)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	rollupVerifierAddr := common.Address{}
	nonce, err := client.PendingNonceAt(context.TODO(), auth.From)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	calculatedMaticAddr := crypto.CreateAddress(auth.From, nonce+1)
	var genesis [32]byte
	poeAddr, _, poe, err := proofofefficiency.DeployProofofefficiency(auth, client, calculatedMaticAddr, maticAddr, rollupVerifierAddr, genesis)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	bridgeAddr, _, bridge, err := bridge.DeployBridge(auth, client, poeAddr)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	if calculatedMaticAddr != bridgeAddr {
		return nil, nil, common.Address{}, fmt.Errorf("bridgeAddr (" + bridgeAddr.String() +
			") is different from the expected contract address (" + calculatedMaticAddr.String() + ")")
	}

	// Approve the bridge and poe to spend 10000 matic tokens
	approvedAmount, _ := new(big.Int).SetString("10000000000000000000000", 10) //nolint:gomnd
	_, err = maticContract.Approve(auth, bridgeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	_, err = maticContract.Approve(auth, poeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, err
	}

	client.Commit()
	return &Client{EtherClient: client, PoE: poe, Bridge: bridge, Matic: maticContract, SCAddresses: []common.Address{poeAddr, bridgeAddr}, auth: auth}, client.Commit, maticAddr, nil
}
