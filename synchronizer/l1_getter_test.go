package synchronizer

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

func Test_L1Get(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	header := new(ethTypes.Header)
	header.Number = big.NewInt(150)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Once()

	sut.Initialize()
	sut.launchWork()
}

func Test_L1Get2(t *testing.T) {
	sut, ethermans, _ := setup(t)
	etherman := ethermans[0]
	header := new(ethTypes.Header)
	header.Number = big.NewInt(150)
	etherman.
		On("HeaderByNumber", mock.Anything, mock.Anything).
		Return(header, nil).
		Maybe()
	expectedCalls(t, etherman, 1)
	sut.Initialize()
	sut.launchWork()
	var waitDuration = time.Duration(0)

	sut.step(&waitDuration)
}

func setup(t *testing.T) (*L1DataRetriever, []*ethermanMock, chan getRollupInfoByBlockRangeResult) {
	etherman := newEthermanMock(t)
	ethermansMock := []*ethermanMock{etherman}
	ethermans := []ethermanInterface{etherman}
	resultChannel := make(chan getRollupInfoByBlockRangeResult)
	sut := NewL1Sync(context.Background(), ethermans, 100, 10, resultChannel)
	return sut, ethermansMock, resultChannel
}

func expectedCalls(t *testing.T, etherman *ethermanMock, calls int) {
	etherman.
		On("GetRollupInfoByBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil, nil).
		Times(calls)
}
