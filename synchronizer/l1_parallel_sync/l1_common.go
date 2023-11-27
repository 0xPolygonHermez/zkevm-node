package l1_parallel_sync

import (
	"context"
	"time"

	"golang.org/x/exp/constraints"
)

// TDOO: There is no min/max function in golang??
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

type contextWithCancel struct {
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (c *contextWithCancel) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *contextWithCancel) isInvalid() bool {
	return c.ctx == nil || c.cancelCtx == nil || (c.ctx != nil && c.ctx.Err() != nil)
}

func (c *contextWithCancel) createWithCancel(ctxParent context.Context) {
	c.ctx, c.cancelCtx = context.WithCancel(ctxParent)
}

func (c *contextWithCancel) createWithTimeout(ctxParent context.Context, timeout time.Duration) {
	c.ctx, c.cancelCtx = context.WithTimeout(ctxParent, timeout)
}

func (c *contextWithCancel) cancel() {
	if c.cancelCtx != nil {
		c.cancelCtx()
	}
}

func newContextWithTimeout(ctxParent context.Context, timeout time.Duration) contextWithCancel {
	ctx := contextWithCancel{}
	ctx.createWithTimeout(ctxParent, timeout)
	return ctx
}

func newContextWithNone(ctxParent context.Context) contextWithCancel {
	ctx := contextWithCancel{ctx: ctxParent}
	return ctx
}
