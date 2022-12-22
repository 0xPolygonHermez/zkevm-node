package sequencer

// TBD. Considerations:
// - Should wait for a block to be finalized: https://www.alchemy.com/overviews/ethereum-commitment-levels https://ethereum.github.io/beacon-APIs/#/Beacon/getStateFinalityCheckpoints

type closingSignalsManager struct {
	finalizer *finalizer
}

func newClosingSignalsManager(finalizer *finalizer) *closingSignalsManager {
	return &closingSignalsManager{finalizer: finalizer}
}

func (c *closingSignalsManager) Start() {}
