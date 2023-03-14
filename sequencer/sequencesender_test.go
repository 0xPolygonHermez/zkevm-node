package sequencer

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	ethermanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

var (
	s *Sequencer
)

func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

// Test_getSequencesToSend_OversizedDataError tests the error when the data is oversized
func Test_getSequencesToSend_OversizedDataError(t *testing.T) {
	stateMock := NewStateMock(t)
	ethermanMock := NewEthermanMock(t)
	s = &Sequencer{
		state:    stateMock,
		etherman: ethermanMock,
	}
	ctx := context.Background()
	toAddress := common.HexToAddress(operations.DefaultSequencerAddress)
	// data is 250 kilobytes
	data := make([]byte, txMaxSize+1)
	tx := types.NewTransaction(1, toAddress, big.NewInt(1), 21000, big.NewInt(30), data)
	batchL2Data, err := state.EncodeTransaction(*tx)
	batch := &state.Batch{
		BatchNumber:    1,
		BatchL2Data:    batchL2Data,
		Timestamp:      time.Now(),
		Transactions:   []types.Transaction{*tx},
		GlobalExitRoot: common.HexToHash("0x1"),
	}
	sequences := []ethermanTypes.Sequence{
		{
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      batch.Timestamp.Unix(),
			BatchL2Data:    batch.BatchL2Data,
			BatchNumber:    batch.BatchNumber,
		},
	}
	stateMock.On("GetLastVirtualBatchNum", ctx, nil).Return(uint64(0), nil)
	stateMock.On("IsBatchClosed", ctx, uint64(1), nil).Return(true, nil).Once()
	stateMock.On("GetBatchByNumber", ctx, uint64(1), nil).Return(batch, nil)
	ethermanMock.On("EstimateGasSequenceBatches", common.HexToAddress("0x0000000000000000000000000000000000000000"), sequences).Return(tx, nil)
	stateMock.On("IsBatchClosed", ctx, uint64(2), nil).Return(false, nil).Once()
	stateMock.On("GetTimeForLatestBatchVirtualization", ctx, nil).Return(batch.Timestamp, nil)
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	var (
		seq []ethermanTypes.Sequence
	)
	logOutput := captureLogOutput(func() {
		seq, err = s.getSequencesToSend(ctx)
	})
	require.Empty(t, seq)
	require.Error(t, err)
	if logOutput == "" {
		t.Error("Expected fatal error did not occur")
	} else if !bytes.Contains([]byte(logOutput), []byte("Expected fatal error occurred")) {
		t.Error("Expected fatal error did not contain the correct message")
	}
}
