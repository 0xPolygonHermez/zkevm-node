package etherman

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
)
type EtherMan struct {
    Poe *proofofefficiency.Proofofefficiency
}

func NewEtherman() (*EtherMan, error) {
	var address common.Address
	var backend bind.ContractBackend
	poe := proofofefficiency.NewProofofefficiency(address, backend)

	return &EtherMan{Poe: poe}, nil
}

//This function retrieves the ethereum block information by ethereum block number
func (etherMan *EtherMan) EthBlockByNumber(blockNum int64) (types.Block, error) {

}

//This function retrieves the batches information that are included in a specific ethereum block
func (etherMan *EtherMan) GetBatchesByBlock(blockNum int64) ([]batch, error) {

}

//This function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *EtherMan) GetBatchesFromBlockTo(fromBlock uint, toBlock uint) ([]batch, error) {

}

//This function allows the sequencer send a new batch proposal to the rollup
func (etherMan *EtherMan) sendBatch(batch) (txId, error) {

}

//This function allows the agregator send the proof for a batch and consolidate it
func (etherMan *EtherMan) consolidateBatch(batch, proof) (txId, error) {

}