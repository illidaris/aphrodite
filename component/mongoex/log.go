package mongoex

import (
	"context"
	"fmt"
	"sync"

	"github.com/illidaris/aphrodite/pkg/logex"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger
var once sync.Once
var logger *Logger

func NewLoggerOptions() *options.LoggerOptions {
	return &options.LoggerOptions{
		ComponentLevels: map[options.LogComponent]options.LogLevel{
			options.LogComponentAll: options.LogLevelInfo,
		},
		MaxDocumentLength: 1024,
		Sink:              NewLogger(),
	}
}
func NewLogger() *Logger {
	once.Do(func() {
		logger = &Logger{}
		l = zap.L().WithOptions(zap.AddCallerSkip(3))
	})
	return logger
}

var _ = options.LogSink(&Logger{})

type Logger struct{}

func (l *Logger) Info(_ int, message string, _ ...interface{}) {
	Log(context.TODO(), message, zapcore.InfoLevel)
}
func (l *Logger) Error(err error, message string, _ ...interface{}) {
	Log(context.TODO(), fmt.Sprintf("mongo_err %s, err: %v", message, err), zapcore.ErrorLevel)
}

func Log(ctx context.Context, msg string, lvl zapcore.Level) {
	l.Log(lvl, msg, logex.FieldsFromCtx(ctx)...)
}
