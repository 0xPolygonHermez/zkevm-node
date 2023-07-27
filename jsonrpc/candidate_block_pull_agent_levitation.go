// LEVITATION_BEGIN

package jsonrpc

// CandidateBlockPullAgent is pull candidate blocks from decentralized sequencer
type CandidateBlockPullAgent struct {
	config  Config
	chainID uint64
	handler *Handler
}

// NewCandidateBlockPullAgent returns the CandidateBlockPullAgent
func NewCandidateBlockPullAgent(
	cfg Config,
	chainID uint64,
) *CandidateBlockPullAgent {
	agent := &CandidateBlockPullAgent{
		config:  cfg,
		chainID: chainID,
	}
	return agent
}

// Start initializes the JSON RPC server to listen for request
func (agent *CandidateBlockPullAgent) Start() error {
	return nil
}

// Stop shutdown the rpc server
func (agent *CandidateBlockPullAgent) Stop() error {
	return nil
}
