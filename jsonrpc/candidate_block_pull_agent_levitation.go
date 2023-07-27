// LEVITATION_BEGIN

package jsonrpc

import (
	"time"
)

// CandidateBlockPullAgent is pull candidate blocks from decentralized sequencer
type CandidateBlockPullAgent struct {
	config        Config
	chainID       uint64
	exitRequested bool
}

// NewCandidateBlockPullAgent returns the CandidateBlockPullAgent
func NewCandidateBlockPullAgent(
	cfg Config,
	chainID uint64,
) *CandidateBlockPullAgent {
	agent := &CandidateBlockPullAgent{
		config:        cfg,
		chainID:       chainID,
		exitRequested: false,
	}

	return agent
}

// Start initializes the JSON RPC server to listen for request
func (agent *CandidateBlockPullAgent) Start() error {
	go agent.doWork()
	return nil
}

// Stop shutdown the rpc server
func (agent *CandidateBlockPullAgent) Stop() error {
	agent.exitRequested = true
	return nil
}

func (agent *CandidateBlockPullAgent) doWork() {
	for !agent.exitRequested {
		time.Sleep(1 * time.Second)
	}
}
