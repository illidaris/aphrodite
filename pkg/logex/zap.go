package logex

import (
	"context"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ = dependency.ILog(Logger{})

// Logger defines a structure for the logger.
type Logger struct {
	core *zap.Logger
}

// Debug logs a message with context at the debug level.
// ctx: The context carrying additional log fields.
// msg: The main log message.
// args: Parameters used to format the msg.
func (l Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.write(ctx, zapcore.DebugLevel, msg, args...)
}

// Info logs a message with context at the info level.
// Similar to Debug but records logs at the info level.
func (l Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.write(ctx, zapcore.InfoLevel, msg, args...)
}

// Warn logs a message with context at the warn level.
// Similar to Debug but records logs at the warn level.
func (l Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.write(ctx, zapcore.WarnLevel, msg, args...)
}

// Error logs a message with context at the error level.
// Similar to Debug but records logs at the error level.
func (l Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.write(ctx, zapcore.ErrorLevel, msg, args...)
}

// write performs the actual logging operation.
// This method of the Logger struct writes the log message with the given parameters.
// ctx: The context from which extra log fields are extracted.
// lvl: The log level.
// msg: The log message to be recorded.
// args: Parameters used to format the msg.
func (l Logger) write(ctx context.Context, lvl zapcore.Level, msg string, args ...interface{}) {
	l.core.Log(lvl, fmt.Sprintf(msg, args...), FieldsFromCtx(ctx)...)
}
