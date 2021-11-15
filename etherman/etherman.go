package etherman

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/state"
)
type EtherMan struct {
    Poe *proofofefficiency.Proofofefficiency
}

func NewEtherman() (*EtherMan, error) {
	//TODO
	var address common.Address
	var backend bind.ContractBackend
	poe, err := proofofefficiency.NewProofofefficiency(address, backend)
	if err != nil {
		return nil, err
	}

	return &EtherMan{Poe: poe}, nil
}

//This function retrieves the ethereum block information by ethereum block number
func (etherMan *EtherMan) EthBlockByNumber(blockNum int64) (types.Block, error) {
	//TODO
	return types.Block{}, nil
}

//This function retrieves the batches information that are included in a specific ethereum block
func (etherMan *EtherMan) GetBatchesByBlock(blockNum int64) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

//This function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *EtherMan) GetBatchesFromBlockTo(fromBlock uint, toBlock uint) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

//This function allows the sequencer send a new batch proposal to the rollup
func (etherMan *EtherMan) sendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

//This function allows the agregator send the proof for a batch and consolidate it
func (etherMan *EtherMan) consolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}