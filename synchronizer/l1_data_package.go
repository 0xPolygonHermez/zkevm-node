package synchronizer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

type l1PackageData struct {
	dataIsValid bool
	data        getRollupInfoByBlockRangeResult
	ctrlIsValid bool
	ctrl        l1ConsumerControl
}

type l1ConsumerControl struct {
	action actionsEnum
}

type actionsEnum int8

const (
	actionNone actionsEnum = 0
	actionStop actionsEnum = 1
)

func (l *l1ConsumerControl) toString() string {
	return fmt.Sprintf("action:%v", l.action)
}

func (l *l1PackageData) toStringBrief() string {
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

func newL1PackageDataControl(action actionsEnum) *l1PackageData {
	return &l1PackageData{
		dataIsValid: false,
		ctrlIsValid: true,
		ctrl: l1ConsumerControl{
			action: action,
		},
	}
}

func newL1PackageDataFromResult(result *getRollupInfoByBlockRangeResult) *l1PackageData {
	if result == nil {
		log.Fatal("newL1PackageDataFromResult: result is nil, the idea of this func is create packages with data")
	}
	return &l1PackageData{
		dataIsValid: true,
		data:        *result,
		ctrlIsValid: false,
	}
}
