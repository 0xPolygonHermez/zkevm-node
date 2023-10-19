package synchronizer

import (
	"context"
	"errors"
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
	ctrl, errTagFound := w.processingRanges.getTagByBlockRange(request.blockRange)
	if errTagFound == nil {
		lastCallElapsedTime := time.Since(ctrl.time)
		if lastCallElapsedTime < w.minTimeBetweenCalls {
			sleepTime := w.minTimeBetweenCalls - lastCallElapsedTime
			log.Infof("workerDecoratorLimitRetriesByTime: br:%s retries:%d last call elapsed time %s < %s, sleeping %s", request.blockRange.String(), ctrl.retries, lastCallElapsedTime, w.minTimeBetweenCalls, sleepTime)
			request.sleepBefore = sleepTime - request.sleepBefore
		}

	}

	res, err := w.workersInterface.asyncRequestRollupInfoByBlockRange(ctx, request)

	if !errors.Is(err, errAllWorkersBusy) {
		// update the tag
		if errTagFound == nil {
			errSetTag := w.processingRanges.setTagByBlockRange(request.blockRange, controlWorkerFlux{time: time.Now(), retries: ctrl.retries + 1})
			if errSetTag != nil {
				log.Warnf("workerDecoratorLimitRetriesByTime: error setting tag %s for blockRange %s err:%s", ctrl, request.blockRange.String(), errSetTag.Error())
			}
		} else {
			ctrl = controlWorkerFlux{time: time.Now(), retries: 0}
			errAddRange := w.processingRanges.addBlockRangeWithTag(request.blockRange, ctrl)
			if errAddRange != nil {
				log.Warnf("workerDecoratorLimitRetriesByTime: error adding blockRange %s err:%s", request.blockRange.String(), errAddRange.Error())
			}
		}

	}
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
