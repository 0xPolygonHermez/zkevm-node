/*
This file contains some function to check batches
*/

package l2_shared

import (
	"encoding/hex"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// CompareBatchFlags is a flag to ignore some fields when comparing batches
type CompareBatchFlags int

const (
	CMP_BATCH_NONE          CompareBatchFlags = 0x0 // CMP_BATCH_NONE No flag
	CMP_BATCH_IGNORE_WIP    CompareBatchFlags = 0x1 // CMP_BATCH_IGNORE_WIP Ignore WIP field
	CMP_BATCH_IGNORE_TSTAMP CompareBatchFlags = 0x2 // CMP_BATCH_IGNORE_TSTAMP Ignore Timestamp field
)

// IsSet check if a flag is set.
// example of usage:  v.IsSet(CMP_BATCH_IGNORE_WIP)
func (c CompareBatchFlags) IsSet(f CompareBatchFlags) bool {
	return c&f != 0
}

// AreEqualStateBatchAndTrustedBatch check is are equal, it response true|false and a debug string
// you could pass some flags to ignore some fields
func AreEqualStateBatchAndTrustedBatch(stateBatch *state.Batch, trustedBatch *types.Batch, flags CompareBatchFlags) (bool, string) {
	if stateBatch == nil || trustedBatch == nil {
		log.Infof("checkIfSynced stateBatch or trustedBatch is nil, so is not synced")
		return false, "nil pointers"
	}
	matchNumber := stateBatch.BatchNumber == uint64(trustedBatch.Number)
	matchGER := stateBatch.GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := stateBatch.LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := stateBatch.StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := stateBatch.Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := true
	if !flags.IsSet(CMP_BATCH_IGNORE_TSTAMP) {
		matchTimestamp = uint64(trustedBatch.Timestamp) == uint64(stateBatch.Timestamp.Unix())
	}
	matchWIP := true
	if !flags.IsSet(CMP_BATCH_IGNORE_WIP) {
		matchWIP = stateBatch.WIP == !trustedBatch.Closed
	}

	matchL2Data := hex.EncodeToString(stateBatch.BatchL2Data) == hex.EncodeToString(trustedBatch.BatchL2Data)

	if matchNumber && matchGER && matchLER && matchSR &&
		matchCoinbase && matchTimestamp && matchL2Data && matchWIP {
		return true, fmt.Sprintf("Equal batch: %v", stateBatch.BatchNumber)
	}

	debugStrResult := ""
	values := []bool{matchNumber, matchGER, matchLER, matchSR, matchCoinbase, matchL2Data}
	names := []string{"matchNumber", "matchGER", "matchLER", "matchSR", "matchCoinbase", "matchL2Data"}
	if !flags.IsSet(CMP_BATCH_IGNORE_TSTAMP) {
		values = append(values, matchTimestamp)
		names = append(names, "matchTimestamp")
	}
	if !flags.IsSet(CMP_BATCH_IGNORE_WIP) {
		values = append(values, matchWIP)
		names = append(names, "matchWIP")
	}
	for i, v := range values {
		log.Debugf("%s: %v", names[i], v)
		if !v {
			debugStrResult += fmt.Sprintf("%s: %v, ", names[i], v)
		}
	}
	return false, debugStrResult
}
