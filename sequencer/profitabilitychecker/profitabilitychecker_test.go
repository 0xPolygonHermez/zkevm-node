package profitabilitychecker_test

import (
	"context"
	"math/big"
	"testing"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func Test_IsSequenceProfitable(t *testing.T) {
	ethman := new(mocks.EthermanMock)
	ethman.On("GetSendSequenceFee").Return(big.NewInt(0), nil)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type:         "default",
		DefaultPrice: pricegetter.TokenPrice{Float: big.NewFloat(2000)},
	})
	require.NoError(t, err)

	pc := profitabilitychecker.New(profitabilitychecker.Config{SendBatchesEvenWhenNotProfitable: false}, ethman, pg)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})

	sequence := ethmanTypes.Sequence{
		Txs: []types.Transaction{*tx1, *tx2, *tx3},
	}
	ctx := context.Background()
	isProfitable, err := pc.IsSequenceProfitable(ctx, sequence)
	require.NoError(t, err)
	require.True(t, isProfitable)
}

func Test_IsSequenceProfitableFalse(t *testing.T) {
	ethman := new(mocks.EthermanMock)
	ethman.On("GetSendSequenceFee").Return(big.NewInt(10000000), nil)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type:         "default",
		DefaultPrice: pricegetter.TokenPrice{Float: big.NewFloat(2000)},
	})
	require.NoError(t, err)

	pc := profitabilitychecker.New(profitabilitychecker.Config{SendBatchesEvenWhenNotProfitable: false}, ethman, pg)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})

	sequence := ethmanTypes.Sequence{
		Txs: []types.Transaction{*tx1, *tx2, *tx3},
	}
	ctx := context.Background()
	isProfitable, err := pc.IsSequenceProfitable(ctx, sequence)
	require.NoError(t, err)
	require.False(t, isProfitable)
}

func Test_IsSendSequencesProfitable(t *testing.T) {
	ethman := new(mocks.EthermanMock)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type:         "default",
		DefaultPrice: pricegetter.TokenPrice{Float: big.NewFloat(2000)},
	})
	require.NoError(t, err)

	pc := profitabilitychecker.New(profitabilitychecker.Config{SendBatchesEvenWhenNotProfitable: false}, ethman, pg)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(1000), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(1000), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(1000), []byte{})

	sequence := ethmanTypes.Sequence{
		Txs: []types.Transaction{*tx1, *tx2, *tx3},
	}

	estGas := big.NewInt(100)
	isProfitable := pc.IsSendSequencesProfitable(estGas, []ethmanTypes.Sequence{sequence})

	require.True(t, isProfitable)
}

func Test_IsSendSequencesFalse(t *testing.T) {
	ethman := new(mocks.EthermanMock)

	pg, err := pricegetter.NewClient(pricegetter.Config{
		Type:         "default",
		DefaultPrice: pricegetter.TokenPrice{Float: big.NewFloat(2000)},
	})
	require.NoError(t, err)

	pc := profitabilitychecker.New(profitabilitychecker.Config{SendBatchesEvenWhenNotProfitable: false}, ethman, pg)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})

	sequence := ethmanTypes.Sequence{
		Txs: []types.Transaction{*tx1, *tx2, *tx3},
	}

	estGas := big.NewInt(100)
	isProfitable := pc.IsSendSequencesProfitable(estGas, []ethmanTypes.Sequence{sequence})

	require.False(t, isProfitable)
}
