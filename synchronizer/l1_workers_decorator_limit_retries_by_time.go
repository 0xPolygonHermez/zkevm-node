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

func (c *controlWorkerFlux) String() string {
	return fmt.Sprintf("time:%s retries:%d", c.time, c.retries)
}

// TODO: Change processingRanges by a cache that take full requests in consideration (no sleep time!)
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

func (w *workerDecoratorLimitRetriesByTime) stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.processingRanges = newLiveBlockRangesWithTag[controlWorkerFlux]()
}

func (w *workerDecoratorLimitRetriesByTime) asyncRequestRollupInfoByBlockRange(ctx context.Context, request requestRollupInfoByBlockRange) (chan responseRollupInfoByBlockRange, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	ctrl, err := w.processingRanges.getTagByBlockRange(request.blockRange)
	if err == nil {
		lastCallElapsedTime := time.Since(ctrl.time)
		if lastCallElapsedTime < w.minTimeBetweenCalls {
			sleepTime := w.minTimeBetweenCalls - lastCallElapsedTime
			log.Infof("workerDecoratorLimitRetriesByTime: br:%s retries:%d last call elapsed time %s < %s, sleeping %s", request.blockRange.String(), ctrl.retries, lastCallElapsedTime, w.minTimeBetweenCalls, sleepTime)
			request.sleepBefore = sleepTime - request.sleepBefore
		}
		err = w.processingRanges.setTagByBlockRange(request.blockRange, controlWorkerFlux{time: time.Now(), retries: ctrl.retries + 1})
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error setting tag %s for blockRange %s", ctrl, request.blockRange.String())
		}
	} else {
		ctrl = controlWorkerFlux{time: time.Now(), retries: 0}
		err = w.processingRanges.addBlockRangeWithTag(request.blockRange, ctrl)
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error adding blockRange %s err:%s", request.blockRange.String(), err.Error())
		}
	}

	res, err := w.workersInterface.asyncRequestRollupInfoByBlockRange(ctx, request)
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
