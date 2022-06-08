package ethermanv2

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/matic"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

var ownershipTransferredSignatureHash = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))

type ethClienter interface {
	ethereum.ChainReader
	ethereum.LogFilterer
	ethereum.TransactionReader
}

// Client is a simple implementation of EtherMan.
type Client struct {
	EtherClient ethClienter
	PoE         *proofofefficiency.Proofofefficiency
	Matic       *matic.Matic
	SCAddresses []common.Address

	auth *bind.TransactOpts
}

// NewClient creates a new etherman.
func NewClient(cfg Config, auth *bind.TransactOpts, PoEAddr common.Address, maticAddr common.Address) (*Client, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	// Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(PoEAddr, ethClient)
	if err != nil {
		return nil, err
	}
	matic, err := matic.NewMatic(maticAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, PoEAddr)

	return &Client{EtherClient: ethClient, PoE: poe, Matic: matic, SCAddresses: scAddresses, auth: auth}, nil
}

// GetRollupInfoByBlockRange function retrieves the Rollup information that are included in all this ethereum blocks
// from block x to block y.
func (etherMan *Client) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, error) {
	// First filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: etherMan.SCAddresses,
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (etherMan *Client) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]state.Block, error) {
	logs, err := etherMan.EtherClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.Block{}, err
	}
	var blockArr []state.Block
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warnf("error processing event. Retrying... Error: %s. vLog: %+v", err.Error(), vLog)
			return nil, err
		}
		if block == nil {
			continue
		}
	}
	return blockArr, nil
}

func (etherMan *Client) processEvent(ctx context.Context, vLog types.Log) (*state.Block, error) {
	switch vLog.Topics[0] {
	case ownershipTransferredSignatureHash:
		ownership, err := etherMan.PoE.ParseOwnershipTransferred(vLog)
		if err != nil {
			return nil, err
		}
		emptyAddr := common.Address{}
		if ownership.PreviousOwner == emptyAddr {
			log.Debug("New rollup smc deployment detected. Deployment account: ", ownership.NewOwner)
		} else {
			log.Debug("Rollup smc OwnershipTransferred from account ", ownership.PreviousOwner, " to ", ownership.NewOwner)
		}
		return nil, nil
	}
	log.Warn("Event not registered: ", vLog)
	return nil, nil
}
