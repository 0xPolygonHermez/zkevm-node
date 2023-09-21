package synchronizer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	cleanUpOlderThan = time.Hour
)

type controlWorkerFlux struct {
	time time.Time
}

type workerDecoratorLimitRetriesByTime struct {
	*workers
	processingRanges    liveBlockRangesGeneric[controlWorkerFlux]
	minTimeBetweenCalls time.Duration
}

func newWorkerDecoratorLimitRetriesByTime(workers *workers, minTimeBetweenCalls time.Duration) *workerDecoratorLimitRetriesByTime {
	return &workerDecoratorLimitRetriesByTime{workers: workers, minTimeBetweenCalls: minTimeBetweenCalls}
}

func (w *workerDecoratorLimitRetriesByTime) String() string {
	return fmt.Sprintf("[FILTERED_LRBT Active/%s]", w.minTimeBetweenCalls) + w.workers.String()
}

func (w *workerDecoratorLimitRetriesByTime) asyncRequestRollupInfoByBlockRange(ctx context.Context, blockRange blockRange) (chan responseRollupInfoByBlockRange, error) {
	ctrl, err := w.processingRanges.getTagByBlockRange(blockRange)
	if err == nil {
		lastCallElapsedTime := time.Since(ctrl.time)
		if lastCallElapsedTime < w.minTimeBetweenCalls {
			sleepTime := w.minTimeBetweenCalls - lastCallElapsedTime
			log.Debugf("workerDecoratorLimitRetriesByTime: last call elapsed time %s < %s, sleeping %s", lastCallElapsedTime, w.minTimeBetweenCalls, sleepTime)
			time.Sleep(sleepTime)
		}
	} else {
		ctrl = controlWorkerFlux{time: time.Now()}
		err = w.processingRanges.addBlockRangeWithTag(blockRange, ctrl)
		if err != nil {
			log.Warnf("workerDecoratorLimitRetriesByTime: error adding blockRange %s with tag %s", blockRange, ctrl)
		}
	}
	res, err := w.workers.asyncRequestRollupInfoByBlockRange(ctx, blockRange)
	w.cleanUpOlderThan(cleanUpOlderThan)
	return res, err
}

func (w *workerDecoratorLimitRetriesByTime) cleanUpOlderThan(timeout time.Duration) {
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
