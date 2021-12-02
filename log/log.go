package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/hermeznetwork/tracerr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger
var logLevel *zap.AtomicLevel

func getDefaultLoggerOrPanic() *zap.SugaredLogger {
	var err error
	if log != nil {
		return log
	}
	// default level: debug
	log, logLevel, err = NewLogger(Config{"debug", []string{"stdout"}})
	if err != nil {
		panic(err)
	}
	return log
}

// Init the logger with defined level. outputs defines the outputs where the
// logs will be sent. By default outputs contains "stdout", which prints the
// logs at the output of the process. To add a log file as output, the path
// should be added at the outputs array. To avoid printing the logs but storing
// them on a file, can use []string{"pathtofile.log"}
func Init(cfg Config) {
	var err error
	log, logLevel, err = NewLogger(cfg)
	if err != nil {
		panic(err)
	}
}

// NewLogger creates the logger with defined level. outputs defines the outputs where the
// logs will be sent. By default, outputs contains "stdout", which prints the
// logs at the output of the process. To add a log file as output, the path
// should be added at the outputs array. To avoid printing the logs but storing
// them on a file, can use []string{"pathtofile.log"}
func NewLogger(cfg Config) (*zap.SugaredLogger, *zap.AtomicLevel, error) {
	var level zap.AtomicLevel
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, nil, fmt.Errorf("error on setting log level: %s", err)
	}

	zapCfg := zap.Config{
		Level:            level,
		Encoding:         "console",
		OutputPaths:      cfg.Outputs,
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey: "timestamp",
			EncodeTime: func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(ts.Local().Format(time.RFC3339))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			// StacktraceKey: "stacktrace",
			StacktraceKey: "",
			LineEnding:    zapcore.DefaultLineEnding,
		},
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, nil, err
	}
	defer logger.Sync()
	withOptions := logger.WithOptions(zap.AddCallerSkip(1))
	return withOptions.Sugar(), &level, nil
}

func sprintStackTrace(st []tracerr.Frame) string {
	builder := strings.Builder{}
	// Skip deepest frame because it belongs to the go runtime and we don't
	// care about it.
	if len(st) > 0 {
		st = st[:len(st)-1]
	}
	for _, f := range st {
		builder.WriteString(fmt.Sprintf("\n%s:%d %s()", f.Path, f.Line, f.Func))
	}
	builder.WriteString("\n")
	return builder.String()
}

// appendStackTraceMaybeArgs will append the stacktrace to the args if one of them
// is a tracerr.Error
func appendStackTraceMaybeArgs(args []interface{}) []interface{} {
	for i := range args {
		if err, ok := args[i].(tracerr.Error); ok {
			st := err.StackTrace()
			return append(args, sprintStackTrace(st))
		}
	}
	return args
}

// Debug calls log.Debug
func Debug(args ...interface{}) {
	getDefaultLoggerOrPanic().Debug(args...)
}

// Info calls log.Info
func Info(args ...interface{}) {
	getDefaultLoggerOrPanic().Info(args...)
}

// Warn calls log.Warn
func Warn(args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLoggerOrPanic().Warn(args...)
}

// Error calls log.Error
func Error(args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLoggerOrPanic().Error(args...)
}

// Fatal calls log.Fatal
func Fatal(args ...interface{}) {
	args = appendStackTraceMaybeArgs(args)
	getDefaultLoggerOrPanic().Fatal(args...)
}

// Debugf calls log.Debugf
func Debugf(template string, args ...interface{}) {
	getDefaultLoggerOrPanic().Debugf(template, args...)
}

// Infof calls log.Infof
func Infof(template string, args ...interface{}) {
	getDefaultLoggerOrPanic().Infof(template, args...)
}

// Warnf calls log.Warnf
func Warnf(template string, args ...interface{}) {
	getDefaultLoggerOrPanic().Warnf(template, args...)
}

// Fatalf calls log.Warnf
func Fatalf(template string, args ...interface{}) {
	getDefaultLoggerOrPanic().Fatalf(template, args...)
}

// Errorf calls log.Errorf and stores the error message into the ErrorFile
func Errorf(template string, args ...interface{}) {
	getDefaultLoggerOrPanic().Errorf(template, args...)
}

// appendStackTraceMaybeKV will append the stacktrace to the KV if one of them
// is a tracerr.Error
func appendStackTraceMaybeKV(msg string, kv []interface{}) string {
	for i := range kv {
		if i%2 == 0 {
			continue
		}
		if err, ok := kv[i].(tracerr.Error); ok {
			st := err.StackTrace()
			return fmt.Sprintf("%v: %v%v\n", msg, err, sprintStackTrace(st))
		}
	}
	return msg
}

// Debugw calls log.Debugw
func Debugw(template string, kv ...interface{}) {
	getDefaultLoggerOrPanic().Debugw(template, kv...)
}

// Infow calls log.Infow
func Infow(template string, kv ...interface{}) {
	getDefaultLoggerOrPanic().Infow(template, kv...)
}

// Warnw calls log.Warnw
func Warnw(template string, kv ...interface{}) {
	template = appendStackTraceMaybeKV(template, kv)
	getDefaultLoggerOrPanic().Warnw(template, kv...)
}

// Errorw calls log.Errorw
func Errorw(template string, kv ...interface{}) {
	template = appendStackTraceMaybeKV(template, kv)
	getDefaultLoggerOrPanic().Errorw(template, kv...)
}

// Fatalw calls log.Fatalw
func Fatalw(template string, kv ...interface{}) {
	template = appendStackTraceMaybeKV(template, kv)
	getDefaultLoggerOrPanic().Fatalw(template, kv...)
}
