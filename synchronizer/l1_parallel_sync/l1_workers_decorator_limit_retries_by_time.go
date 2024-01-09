package l1_parallel_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
)

const (
	timeOfLiveOfEntries = time.Hour
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
	processingRanges    common.Cache[blockRange, controlWorkerFlux]
	minTimeBetweenCalls time.Duration
}

func newWorkerDecoratorLimitRetriesByTime(workers workersInterface, minTimeBetweenCalls time.Duration) *workerDecoratorLimitRetriesByTime {
	return &workerDecoratorLimitRetriesByTime{
		workersInterface:    workers,
		minTimeBetweenCalls: minTimeBetweenCalls,
		processingRanges:    *common.NewCache[blockRange, controlWorkerFlux](common.DefaultTimeProvider{}, timeOfLiveOfEntries),
	}
}

func (w *workerDecoratorLimitRetriesByTime) String() string {
	return fmt.Sprintf("[FILTERED_LRBT Active/%s]", w.minTimeBetweenCalls) + w.workersInterface.String()
}

func (w *workerDecoratorLimitRetriesByTime) stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.processingRanges.Clear()
}

func (w *workerDecoratorLimitRetriesByTime) asyncRequestRollupInfoByBlockRange(ctx context.Context, request requestRollupInfoByBlockRange) (chan responseRollupInfoByBlockRange, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	//ctrl, found := w.processingRanges.getTagByBlockRange(request.blockRange)
	ctrl, found := w.processingRanges.Get(request.blockRange)
	if found {
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
		w.processingRanges.Set(request.blockRange, controlWorkerFlux{time: time.Now(), retries: ctrl.retries + 1})
	}
	w.processingRanges.DeleteOutdated()
	return res, err
}
