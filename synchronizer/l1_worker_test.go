package synchronizer

import (
	context "context"
	"math/big"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Given_Kk(t *testing.T) {
	Etherman := newEthermanMock(t)
	ctx := context.Background()
	header := new(ethTypes.Header)
	header.Number = big.NewInt(1)
	Etherman.
		On("HeaderByNumber", ctx, mock.Anything).
		Return(header, nil).
		Once()
	worker := newWorker(Etherman)
	ch, err := worker.asyncRequestLastBlock(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, ch)
	result := <-ch
	require.NoError(t, result.err)
	require.Equal(t, result.result.block, uint64(1))
}
