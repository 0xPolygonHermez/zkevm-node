package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	cleanUpOlderThan = time.Hour
)

type controlWorkerFlux struct {
	time    time.Time
	retries int
}

type workerDecoratorLimitRetriesByTime struct {
	mutex sync.Mutex
	workersInterface
	processingRanges    liveBlockRangesGeneric[controlWorkerFlux]
	minTimeBetweenCalls time.Duration
}

func newWorkerDecoratorLimitRetriesByTime(workers workersInterface, minTimeBetweenCalls time.Duration) *workerDecoratorLimitRetriesByTime {
	return &workerDecoratorLimitRetriesByTime{workersInterface: workers, minTimeBetweenCalls: minTimeBetweenCalls}
}

func (w *workerDecoratorLimitRetriesByTime) String() string {
	return fmt.Sprintf("[FILTERED_LRBT Active/%s]", w.minTimeBetweenCalls) + w.workersInterface.String()
}

func (w *workerDecoratorLimitRetriesByTime) asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange, sleepBefore time.Duration) (chan responseRollupInfoByBlockRange, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	ctrl, err := w.processingRanges.getTagByBlockRange(blockRange)
	if err == nil {
		lastCallElapsedTime := time.Since(ctrl.time)
		if lastCallElapsedTime < w.minTimeBetweenCalls {
			sleepTime := w.minTimeBetweenCalls - lastCallElapsedTime
			log.Infof("workerDecoratorLimitRetriesByTime: br:%s retries:%d last call elapsed time %s < %s, sleeping %s", blockRange.String(), ctrl.retries, lastCallElapsedTime, w.minTimeBetweenCalls, sleepTime)
			sleepBefore = sleepTime - sleepBefore
		}
		err = w.processingRanges.setTagByBlockRange(blockRange, controlWorkerFlux{time: time.Now(), retries: ctrl.retries + 1})
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error setting tag %s for blockRange %s", ctrl, blockRange)
		}
	} else {
		ctrl = controlWorkerFlux{time: time.Now(), retries: 0}
		err = w.processingRanges.addBlockRangeWithTag(blockRange, ctrl)
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error adding blockRange %s with tag %s", blockRange, ctrl)
		}
	}

	res, err := w.workersInterface.asyncRequestRollupInfoByBlockRange(ctx, blockRange, sleepBefore)
	w.cleanUpOlderThanUnsafe(cleanUpOlderThan)
	return res, err
}

func (w *workerDecoratorLimitRetriesByTime) cleanUpOlderThanUnsafe(timeout time.Duration) {
	brs := w.processingRanges.filterBlockRangesByTag(func(br blockRange, tag controlWorkerFlux) bool {
		return tag.time.Add(timeout).Before(time.Now())
	})
	for _, br := range brs {
		err := w.processingRanges.removeBlockRange(br)
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error removing outdated blockRange %s", br)
		}
	}
}
