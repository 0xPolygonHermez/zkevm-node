package etherman

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/bridge"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/globalexitrootmanager"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/mockverifier"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/proofofefficiency"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func NewSimulatedEtherman(cfg Config, auth *bind.TransactOpts) (etherman *Client, ethBackend *backends.SimulatedBackend, maticAddr common.Address, br *bridge.Bridge, err error) {
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
	const maticDecimalPlaces = 18
	totalSupply, _ := new(big.Int).SetString("10000000000000000000000000000", 10) //nolint:gomnd
	maticAddr, _, maticContract, err := matic.DeployMatic(auth, client, "Matic Token", "MATIC", maticDecimalPlaces, totalSupply)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	rollupVerifierAddr, _, _, err := mockverifier.DeployMockverifier(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	nonce, err := client.PendingNonceAt(context.TODO(), auth.From)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	const posBridge = 2
	calculatedBridgeAddr := crypto.CreateAddress(auth.From, nonce+posBridge)
	const posPoE = 4
	calculatedPoEAddr := crypto.CreateAddress(auth.From, nonce+posPoE)
	var genesis [32]byte
	exitManagerAddr, _, globalExitRoot, err := globalexitrootmanager.DeployGlobalexitrootmanager(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	_, err = globalExitRoot.Initialize(auth, calculatedPoEAddr, calculatedBridgeAddr)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	bridgeAddr, _, br, err := bridge.DeployBridge(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	_, err = br.Initialize(auth, 0, exitManagerAddr)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	poeAddr, _, poe, err := proofofefficiency.DeployProofofefficiency(auth, client)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	_, err = poe.Initialize(auth, exitManagerAddr, maticAddr, rollupVerifierAddr, genesis, auth.From, true, "http://localhost", 1000, "L2") //nolint:gomnd
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}

	if calculatedBridgeAddr != bridgeAddr {
		return nil, nil, common.Address{}, nil, fmt.Errorf("bridgeAddr (%s) is different from the expected contract address (%s)",
			bridgeAddr.String(), calculatedBridgeAddr.String())
	}
	if calculatedPoEAddr != poeAddr {
		return nil, nil, common.Address{}, nil, fmt.Errorf("poeAddr (%s) is different from the expected contract address (%s)",
			poeAddr.String(), calculatedPoEAddr.String())
	}

	// Approve the bridge and poe to spend 10000 matic tokens.
	approvedAmount, _ := new(big.Int).SetString("10000000000000000000000", 10) //nolint:gomnd
	_, err = maticContract.Approve(auth, bridgeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}
	_, err = maticContract.Approve(auth, poeAddr, approvedAmount)
	if err != nil {
		return nil, nil, common.Address{}, nil, err
	}

	client.Commit()
	return &Client{EtherClient: client, PoE: poe, Matic: maticContract, GlobalExitRootManager: globalExitRoot, SCAddresses: []common.Address{poeAddr, exitManagerAddr}, auth: auth}, client, maticAddr, br, nil
}
