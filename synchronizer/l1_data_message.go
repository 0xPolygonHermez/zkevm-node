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
type l1SyncMessage struct {
	// dataIsValid : true if data is valid
	dataIsValid bool
	data        responseRollupInfoByBlockRange
	// ctrlIsValid : true if ctrl is valid
	ctrlIsValid bool
	// ctrl : control package, it send actions to consumer
	ctrl l1ConsumerControl
}

type l1ConsumerControl struct {
	action actionsEnum
}

type actionsEnum int8

const (
	actionNone actionsEnum = 0
	actionStop actionsEnum = 1
)

func newL1SyncMessageControl(action actionsEnum) *l1SyncMessage {
	return &l1SyncMessage{
		dataIsValid: false,
		ctrlIsValid: true,
		ctrl: l1ConsumerControl{
			action: action,
		},
	}
}

func newL1SyncMessageData(result *responseRollupInfoByBlockRange) *l1SyncMessage {
	if result == nil {
		log.Fatal("newL1PackageDataFromResult: result is nil, the idea of this func is create packages with data")
	}
	return &l1SyncMessage{
		dataIsValid: true,
		data:        *result,
		ctrlIsValid: false,
	}
}

func (a actionsEnum) toString() string {
	switch a {
	case actionNone:
		return "actionNone"
	case actionStop:
		return "actionStop"
	default:
		return "actionUnknown"
	}
}

func (l *l1ConsumerControl) toString() string {
	return fmt.Sprintf("action:%s", l.action.toString())
}

func (l *l1SyncMessage) toStringBrief() string {
	res := ""
	if l.dataIsValid {
		res += fmt.Sprintf("data:%v ", l.data.toStringBrief())
	} else {
		res += " NO_DATA "
	}
	if l.ctrlIsValid {
		res += fmt.Sprintf("ctrl:%v ", l.ctrl.toString())
	} else {
		res += " NO_CTRL "
	}

	return res
}
