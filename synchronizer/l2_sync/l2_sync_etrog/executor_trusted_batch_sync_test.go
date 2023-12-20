package l2_sync_etrog

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// changeL2Block + deltaTimeStamp + indexL1InfoTree
	codedL2BlockHeader = "0b73e6af6f00000000"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs1 = "ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff"
)

func TestIncrementalProcess(t *testing.T) {
	// Arrange
	stateMock := NewStateInterfaceMock(t)
	syncMock := NewSynchronizerInterfaceMock(t)
	sut := SyncTrustedBatchExecutorForEtrog{
		state: stateMock,
		sync:  syncMock,
	}
	ctx := context.Background()

	stateBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1)
	trustedBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1 + codedL2BlockHeader + codedRLP2Txs1)
	//deltaBatchL2Data := []byte{4}
	batchNumber := uint64(123)
	data := l2_shared.ProcessData{
		BatchNumber:  batchNumber,
		OldStateRoot: common.Hash{},
		TrustedBatch: &types.Batch{
			Number:      123,
			BatchL2Data: trustedBatchL2Data,
		},
		StateBatch: &state.Batch{
			BatchNumber: batchNumber,
			BatchL2Data: stateBatchL2Data,
		},
	}
	stateMock.
		On("UpdateBatchL2Data", mock.Anything, batchNumber, trustedBatchL2Data, mock.Anything).
		Return(nil).
		Once()

	stateMock.
		On("GetCurrentL1InfoRoot").
		Return(state.ZeroHash).
		Once()
	stateMock.
		On("GetL1InfoRootLeafByL1InfoRoot", mock.Anything, mock.Anything, mock.Anything).
		Return(state.L1InfoTreeExitRootStorageEntry{}, nil).
		Once()

	stateMock.
		On("GetForkIDByBatchNumber", batchNumber).
		Return(uint64(7)).
		Once()

	processBatchResp := &state.ProcessBatchResponse{}
	stateMock.
		On("ProcessBatchV2", mock.Anything, mock.Anything, true).
		Return(processBatchResp, nil).
		Once()

	syncMock.
		On("PendingFlushID", mock.Anything, mock.Anything).
		Once()

	syncMock.
		On("CheckFlushID", mock.Anything).
		Return(nil).
		Maybe()
	// Act
	res, err := sut.IncrementalProcess(ctx, &data, nil)
	// Assert
	log.Info(res)
	require.NoError(t, err)
	require.Equal(t, trustedBatchL2Data, res.UpdateBatch.BatchL2Data)
}
