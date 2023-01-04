package aggregator

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/mocks"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendFinalProof(t *testing.T) {
	require := require.New(t)
	stateMock := mocks.NewStateMock(t)
	ethTxManager := mocks.NewEthTxManager(t)
	etherman := mocks.NewEtherman(t)
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	currentNonce := uint64(1)
	estimatedGas := uint64(1)
	suggestedGasPrice := big.NewInt(1)
	from := common.BytesToAddress([]byte("aggregator"))
	var to *common.Address
	value := big.NewInt(0)
	var data []byte = nil
	finalBatch := state.Batch{
		LocalExitRoot: common.BytesToHash([]byte("localExitRoot")),
		StateRoot:     common.BytesToHash([]byte("stateRoot")),
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      estimatedGas,
		GasPrice: suggestedGasPrice,
		Data:     data,
	})
	recursiveProof := &state.Proof{
		BatchNumber:      batchNum,
		BatchNumberFinal: batchNumFinal,
	}
	finalProof := &pb.FinalProof{}
	cfg := Config{}
	a, err := New(cfg, stateMock, ethTxManager, etherman)
	a.ctx, a.exit = context.WithCancel(context.Background())
	require.NoError(err)
	// set up mocks
	stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Return(&finalBatch, nil).Once()
	expectedInputs := ethmanTypes.FinalProofInputs{
		FinalProof:       finalProof,
		NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
		NewStateRoot:     finalBatch.StateRoot.Bytes(),
	}
	etherman.On("EstimateGasForTrustedVerifyBatches", batchNum-1, batchNumFinal, &expectedInputs).Return(tx, nil).Once()
	etherman.On("GetPublicAddress").Return(from, nil).Once()
	txID := fmt.Sprintf("%d-%d", batchNum, batchNumFinal)
	ethTxManager.On("Add", mock.Anything, txManagerOwner, txID, from, to, value, data, nil).Return(nil).Once()
	res := ethtxmanager.MonitoredTxResult{
		Status: ethtxmanager.MonitoredTxStatusConfirmed,
	}
	ethTxManager.On("ManageTxs").Once()
	ethTxManager.On("Result", mock.Anything, txManagerOwner, txID, nil).Return(res, nil).Once()
	ethTxManager.On("SetStatusDone", mock.Anything, txManagerOwner, txID, nil).Return(nil).Once()
	verifiedBatch := state.VerifiedBatch{
		BatchNumber: batchNumFinal,
	}
	stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&verifiedBatch, nil).Once()
	etherman.On("GetLatestVerifiedBatchNum").Return(batchNumFinal, nil).Once()
	stateMock.On("DeleteGeneratedProofs", mock.Anything, batchNum, batchNumFinal, nil).Run(func(args mock.Arguments) {
		// test is done, stop the sendFinalProof method
		a.exit()
	}).Return(nil).Once()
	go func() {
		finalMsg := finalProofMsg{
			proverID:       "proverID",
			recursiveProof: recursiveProof,
			finalProof:     finalProof,
		}
		a.finalProof <- finalMsg
	}()

	a.sendFinalProof()
}
