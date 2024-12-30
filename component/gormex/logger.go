package gormex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/illidaris/aphrodite/pkg/logex"
	iLog "github.com/illidaris/logger"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

var _ = logger.Interface(&GormLogger{})

// 定义GormLogger相关的配置和选项，提供日志记录功能。

const (
	SqlLogFormat = "elapsed:%dms, affect:%d, err:%v, sql:%s" // SQL日志格式
)

// NewLogger 创建一个新的GormLogger实例。
// opts ...GormLoggerOption: 一个或多个GormLogger选项函数，用于配置GormLogger实例。
// 返回值 logger.Interface: 返回实现了logger.Interface接口的日志实例。
func NewLogger(opts ...GormLoggerOption) logger.Interface {
	l := &GormLogger{
		LogLevel:       logger.Info,            // 默认日志级别为Info
		SlowThreshold:  300 * time.Millisecond, // 默认慢查询阈值为300ms
		IgnoreNoAffect: true,                   // 默认忽略未影响行的日志
	}
	for _, opt := range opts {
		opt(l) // 应用传入的配置选项
	}
	l.core = zap.L().WithOptions(zap.AddCallerSkip(3)) // 核心日志记录器配置
	return l
}

// GormLoggerOption 定义了一个函数类型，用于设置GormLogger的配置。
type GormLoggerOption func(*GormLogger)

// WithCore 通过提供一个*zap.Logger实例来设置GormLogger的核心日志记录器。
func WithCore(core *zap.Logger) GormLoggerOption {
	return func(l *GormLogger) {
		l.core = core
	}
}

// WithIgnoreNoAffect 设置是否忽略未影响行的日志记录。
func WithIgnoreNoAffect(v bool) GormLoggerOption {
	return func(l *GormLogger) {
		l.IgnoreNoAffect = v
	}
}

// WithLogLevel 设置日志级别。
func WithLogLevel(level logger.LogLevel) GormLoggerOption {
	return func(l *GormLogger) {
		l.LogLevel = level
	}
}

// WithIgnoreRecordNotFoundError 设置是否忽略记录未找到的错误。
func WithIgnoreRecordNotFoundError(v bool) GormLoggerOption {
	return func(l *GormLogger) {
		l.IgnoreRecordNotFoundError = v
	}
}

// WithSlowThreshold 设置慢查询的阈值。
func WithSlowThreshold(v time.Duration) GormLoggerOption {
	return func(l *GormLogger) {
		l.SlowThreshold = v
	}
}

// GormLogger 结构体定义了Gorm日志记录器的配置和实现。
type GormLogger struct {
	core                      *zap.Logger     // 核心日志记录器
	LogLevel                  logger.LogLevel // 当前日志级别
	IgnoreRecordNotFoundError bool            // 是否忽略RecordNotFoundError错误
	IgnoreNoAffect            bool            // 是否忽略未影响行的日志
	SlowThreshold             time.Duration   // 慢查询阈值
}

// getLevel 根据提供的logger.LogLevel返回对应的iLog.Level。
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

// LogMode 设置日志记录模式。
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info 记录信息级别的日志。
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Info), fmt.Sprintf(msg, args...))
}

// Warn 记录警告级别的日志。
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Warn), fmt.Sprintf(msg, args...))
}

// Error 记录错误级别的日志。
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, getLevel(logger.Error), fmt.Sprintf(msg, args...))
}

// Trace 记录SQL相关的日志，包括执行时间、影响的行数以及可能的错误。
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.Error(ctx, SqlLogFormat, elapsed.Milliseconds(), rows, err, sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.Warn(ctx, "[SLOW]"+SqlLogFormat, elapsed.Milliseconds(), rows, err, sql)
	case l.LogLevel == logger.Info:
		if l.IgnoreNoAffect && err == nil && rows == 0 {
			return
		}
		l.Info(ctx, SqlLogFormat, elapsed.Milliseconds(), rows, err, sql)
	}
}

// Log 记录指定级别的日志消息。
func Log(ctx context.Context, lvl iLog.Level, msg string) {
	iLog.LogFrmCtx(ctx, getLevel(logger.Info), msg, logex.FieldsFromCtx(ctx)...)
}
