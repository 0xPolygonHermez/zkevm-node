package etherman

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

type EtherMan struct {
	EtherClient *ethclient.Client
	PoE         *proofofefficiency.Proofofefficiency
}

func NewEtherman(url string, poeAddr common.Address) (*EtherMan, error) {
	//TODO
	//Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", url, err)
		return nil, err
	}
	//Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(poeAddr, ethClient)
	if err != nil {
		return nil, err
	}

	return &EtherMan{EtherClient: ethClient, PoE: poe}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *EtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, nil
	}
	return block, nil //TODO Change types.Block. It only needs hash, hash parent and block number
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *EtherMan) GetBatchesByBlock(blockNum int64) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *EtherMan) GetBatchesFromBlockTo(fromBlock uint, toBlock uint) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup
func (etherMan *EtherMan) SendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

// ConsolidateBatch function allows the agregator send the proof for a batch and consolidate it
func (etherMan *EtherMan) ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}
