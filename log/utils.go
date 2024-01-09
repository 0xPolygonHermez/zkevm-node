package log

import (
	"context"
)

type contextKey string

const (
	ctxKeyLogger contextKey = "logger"
)

func CtxWithLogger(parentCtx context.Context, logger *Logger) context.Context {
	return context.WithValue(parentCtx, ctxKeyLogger, logger)
}

func LoggerFromCtx(ctx context.Context) *Logger {
	res := ctx.Value(ctxKeyLogger)
	if res == nil {
		return getDefaultLog()
	}
	if logger, ok := res.(*Logger); ok {
		return logger
	}
	return getDefaultLog()
}
