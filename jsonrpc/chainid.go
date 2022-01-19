package jsonrpc

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state"
)

type chainIDSelector struct {
	defaultChainID   uint64
	sequencerAddress common.Address
	s                state.State
}

func newChainIDSelector(defaultChainID uint64, sequencerAddress common.Address, s state.State) *chainIDSelector {
	return &chainIDSelector{
		defaultChainID:   defaultChainID,
		sequencerAddress: sequencerAddress,
		s:                s,
	}
}

func (s *chainIDSelector) getChainID() (uint64, error) {
	sequencer, err := s.s.GetSequencer(context.Background(), s.sequencerAddress)
	if err == nil {
		return sequencer.ChainID.Uint64(), nil
	}

	if errors.Is(err, state.ErrNotFound) {
		return s.defaultChainID, nil
	}

	return 0, err
}
