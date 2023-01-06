package sequencer

import "github.com/ethereum/go-ethereum/common"

// TBD. Considerations:
// - Should wait for a block to be finalized: https://www.alchemy.com/overviews/ethereum-commitment-levels https://ethereum.github.io/beacon-APIs/#/Beacon/getStateFinalityCheckpoints

type closingSignalsManager struct {
	finalizer *finalizer
}

type L2ReorgEvent struct {
	TxHashes []common.Hash
}

func newClosingSignalsManager(finalizer *finalizer) *closingSignalsManager {
	return &closingSignalsManager{finalizer: finalizer}
}

func (c *closingSignalsManager) Start() {}
