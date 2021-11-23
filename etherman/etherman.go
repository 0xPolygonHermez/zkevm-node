package etherman

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

type EtherMan interface {
	EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error)
	GetBatchesByBlock(blockNum int64) ([]state.Batch, error)
	GetBatchesFromBlockTo(fromBlock uint, toBlock uint) ([]state.Batch, error)
	SendBatch(batch state.Batch) (common.Hash, error)
	ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error)
}

type BasicEtherMan struct {
	EtherClient *ethclient.Client
	PoE         *proofofefficiency.Proofofefficiency

	key *keystore.Key
}

func NewEtherman(cfg Config) (EtherMan, error) {
	//TODO
	//Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.Url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.Url, err)
		return nil, err
	}
	//Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoeAddress, ethClient)
	if err != nil {
		return nil, err
	}

	key, err := decryptKeystore(cfg.PrivateKeyPath, cfg.PrivateKeyPassword)
	if err != nil {
		return nil, err
	}

	return &BasicEtherMan{EtherClient: ethClient, PoE: poe, key: key}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *BasicEtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, nil
	}
	return block, nil //TODO Change types.Block. It only needs hash, hash parent and block number
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *BasicEtherMan) GetBatchesByBlock(blockNum int64) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *BasicEtherMan) GetBatchesFromBlockTo(fromBlock uint, toBlock uint) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup
func (etherMan *BasicEtherMan) SendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

// ConsolidateBatch function allows the agregator send the proof for a batch and consolidate it
func (etherMan *BasicEtherMan) ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}
