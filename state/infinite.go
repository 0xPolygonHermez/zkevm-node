package state

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// InfiniteSafeRun executes a function and in case it fails,
// runs the function again infinitely
func InfiniteSafeRun(fn func(), errorMessage string, restartInterval time.Duration) {
	for {
		SafeRun(fn, errorMessage)
		time.Sleep(restartInterval)
	}
}

// SafeRun executes a function with a deferred recover
// to avoid to panic.
func SafeRun(fn func(), errorMessage string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(errorMessage, r)
		}
	}()
	fn()
}
