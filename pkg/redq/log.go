package redq

import (
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

var _ = asynq.Logger(&ZapLog{})

func NewZapLog() *ZapLog {
	z := zap.L().WithOptions(zap.AddCallerSkip(2))
	return &ZapLog{
		core: z,
	}
}

type ZapLog struct {
	core *zap.Logger
}

// Debug logs a message at Debug level.
func (l ZapLog) Debug(args ...interface{}) {
	l.core.Debug(fmt.Sprintln(args...))
}

// Info logs a message at Info level.
func (l ZapLog) Info(args ...interface{}) {
	l.core.Info(fmt.Sprintln(args...))
}

// Warn logs a message at Warning level.
func (l ZapLog) Warn(args ...interface{}) {
	l.core.Warn(fmt.Sprintln(args...))
}

// Error logs a message at Error level.
func (l ZapLog) Error(args ...interface{}) {
	l.core.Error(fmt.Sprintln(args...))
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l ZapLog) Fatal(args ...interface{}) {
	l.core.Fatal(fmt.Sprintln(args...))
}
