package jsonrpc

type chainIDSelector struct {
	chainID uint64
}

func newChainIDSelector(chainID uint64) *chainIDSelector {
	return &chainIDSelector{
		chainID: chainID,
	}
}

func (s *chainIDSelector) getChainID() (uint64, error) {
	return s.chainID, nil
}
