// package synchronizer
// This file contains common struct definitions and functions used by L1 sync.
// l1DataMessage : struct to hold L1 rollup info data package send from producer to consumer
//
//	This packages could contain data or control information.
//	 - data is a real rollup info
//	 - control send actions to consumer
//
// Constructors:
// - newL1PackageDataControl: create a l1PackageData with only control information
// - newL1PackageData: create a l1PackageData with data and control information
package synchronizer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// l1SyncMessage : struct to hold L1 rollup info data package
// It could contain data or control information, or both.
// A control package is used to send actions to consumer or to notify that producer is fully synced.
type l1SyncMessage struct {
	// dataIsValid : true if data field is valid
	dataIsValid bool
	// data: is the rollup info data
	data rollupInfoByBlockRangeResult
	// ctrlIsValid : true if ctrl field is valid
	ctrlIsValid bool
	// ctrl : control package, it send actions to consumer
	ctrl l1ConsumerControl
}

type l1ConsumerControl struct {
	event eventEnum
}

type eventEnum int8

const (
	eventNone                  eventEnum = 0
	eventStop                  eventEnum = 1
	eventProducerIsFullySynced eventEnum = 2
)

func newL1SyncMessageControl(event eventEnum) *l1SyncMessage {
	return &l1SyncMessage{
		dataIsValid: false,
		ctrlIsValid: true,
		ctrl: l1ConsumerControl{
			event: event,
		},
	}
}

func newL1SyncMessageData(result *rollupInfoByBlockRangeResult) *l1SyncMessage {
	if result == nil {
		log.Fatal("newL1PackageDataFromResult: result is nil, the idea of this func is create packages with data")
	}
	return &l1SyncMessage{
		dataIsValid: true,
		data:        *result,
		ctrlIsValid: false,
	}
}

func (a eventEnum) String() string {
	switch a {
	case eventNone:
		return "actionNone"
	case eventStop:
		return "actionStop"
	case eventProducerIsFullySynced:
		return "eventIsFullySynced"
	default:
		return "actionUnknown"
	}
}

func (l *l1ConsumerControl) String() string {
	return fmt.Sprintf("action:%s", l.event.String())
}

func (l *l1SyncMessage) toStringBrief() string {
	res := ""
	if l.dataIsValid {
		res += fmt.Sprintf("data:%v ", l.data.toStringBrief())
	} else {
		res += " NO_DATA "
	}
	if l.ctrlIsValid {
		res += fmt.Sprintf("ctrl:%v ", l.ctrl.String())
	} else {
		res += " NO_CTRL "
	}

	return res
}
