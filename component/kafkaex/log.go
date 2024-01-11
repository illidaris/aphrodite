package kafkaex

import (
	"context"
	"io"
	"log"

	"github.com/IBM/sarama"
)

var _ = sarama.StdLogger(&defaultLogger{})

var logger ILogger = newdDefaultLogger()

func SetLogger(l ILogger) {
	logger = l
	sarama.Logger = l
}

type ILogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
}

type defaultLogger struct {
	core *log.Logger
}

func newdDefaultLogger() *defaultLogger {
	return &defaultLogger{
		core: log.New(io.Discard, "[KafkaexSarama] ", log.LstdFlags),
	}
}

func (l *defaultLogger) Print(v ...interface{}) {
	l.core.Print(v...)
}

func (l *defaultLogger) Printf(format string, v ...interface{}) {
	l.core.Printf(format, v...)
}

func (l *defaultLogger) Println(v ...interface{}) {
	l.core.Println(v...)
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.core.Printf(msg, args...)
}

func (l *defaultLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.core.Printf(msg, args...)
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.core.Printf(msg, args...)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.core.Printf(msg, args...)
}
