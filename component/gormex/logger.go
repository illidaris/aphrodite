package gormex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/illidaris/extensions/pkg/logex"
	iLog "github.com/illidaris/logger"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

var _ = logger.Interface(&GormLogger{})

const (
	SqlLogFormat = "elapsed:%dms,affect:%d,err:%s,sql:%s"
)

func NewLogger() logger.Interface {
	l := &GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 200 * time.Millisecond,
	}
	l.core = zap.L().WithOptions(zap.AddCallerSkip(3))
	return l
}

type GormLogger struct {
	core                      *zap.Logger
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func getLevel(level logger.LogLevel) iLog.Level {
	switch level {
	case logger.Info:
		return iLog.InfoLevel
	case logger.Warn:
		return iLog.WarnLevel
	case logger.Error:
		return iLog.ErrorLevel
	default:
		return iLog.DebugLevel
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Info), fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Warn), fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Error), fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	msg := fmt.Sprintf(SqlLogFormat, elapsed.Milliseconds(), rows, err, sql)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.Error(ctx, msg)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.Warn(ctx, "[SLOW]"+msg)
	case l.LogLevel == logger.Info:
		l.Info(ctx, msg)
	}
}

func Log(ctx context.Context, lvl iLog.Level, msg string) {
	iLog.LogFrmCtx(ctx, getLevel(logger.Info), msg, logex.FieldsFromCtx(ctx)...)
}
