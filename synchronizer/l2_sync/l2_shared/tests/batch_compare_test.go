package test_l2_shared

import (
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func newExampleStateBatch() state.Batch {
	return state.Batch{
		BatchNumber:    1,
		Coinbase:       common.HexToAddress("0x01"),
		StateRoot:      common.HexToHash("0x02"),
		LocalExitRoot:  common.HexToHash("0x03"),
		GlobalExitRoot: common.HexToHash("0x04"),
		Timestamp:      time.Unix(0, 0),
		WIP:            true,
		BatchL2Data:    []byte("0x05"),
	}
}

func newExampleTrustedBatch() types.Batch {
	return types.Batch{
		Number:         1,
		Coinbase:       common.HexToAddress("0x01"),
		StateRoot:      common.HexToHash("0x02"),
		LocalExitRoot:  common.HexToHash("0x03"),
		GlobalExitRoot: common.HexToHash("0x04"),
		Timestamp:      0,
		Closed:         false,
		BatchL2Data:    []byte("0x05"),
	}
}

func TestA(t *testing.T) {
	stateBatch := newExampleStateBatch()
	trustedBatch := newExampleTrustedBatch()
	equal, _ := l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_NONE)
	require.True(t, equal)

	stateBatch = newExampleStateBatch()
	trustedBatch = newExampleTrustedBatch()
	trustedBatch.Number = 2
	equal, _ = l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_NONE)
	require.False(t, equal)

	stateBatch = newExampleStateBatch()
	trustedBatch = newExampleTrustedBatch()
	trustedBatch.Timestamp = 123
	equal, _ = l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_NONE)
	require.False(t, equal)
	equal, _ = l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_IGNORE_TSTAMP)
	require.True(t, equal)

	stateBatch = newExampleStateBatch()
	stateBatch.WIP = true
	trustedBatch = newExampleTrustedBatch()
	trustedBatch.Closed = true
	equal, _ = l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_NONE)
	require.False(t, equal)
	equal, _ = l2_shared.AreEqualStateBatchAndTrustedBatch(&stateBatch, &trustedBatch, l2_shared.CMP_BATCH_IGNORE_WIP)
	require.True(t, equal)
}
