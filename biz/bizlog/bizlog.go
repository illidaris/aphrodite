package bizlog

import (
	"context"
	"fmt"
	"sync"

	"github.com/illidaris/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once
var l *Logger

type Entry struct {
	Uid      int64    `json:"ibzl_uid,omitempty"`      // 【WHO】主角Id
	Pid      string   `josn:"ibzl_pid,omitempty"`      // 【WHO】主角PID
	Name     string   `json:"ibzl_name,omitempty"`     // 【WHO】主角名
	BizId    int64    `json:"ibzl_bizid,omitempty"`    // 【WHERE】业务Id
	System   string   `json:"ibzl_system,omitempty"`   // 【WHERE】所属系统
	Category int32    `json:"ibzl_cgy,omitempty"`      // 大类
	Action   int64    `json:"ibzl_action,omitempty"`   // 行为
	Message  string   `json:"ibzl_msg,omitempty"`      // 【WHAT】什么事
	Tags     []string `json:"ibzl_tags,omitempty"`     // 标签
	OpAt     int64    `json:"ibzl_opat,omitempty"`     // 【WHEN】发生时间点
	CreateAt int64    `json:"ibzl_createat,omitempty"` // 入库时间点
	Dst      int32    `json:"log_dst,omitempty"`       // 日志用途 兼容之前的设计
	TraceId  string   `json:"traceId,omitempty"`       // 追踪ID 兼容之前的设计
}

type Logger struct {
	Opts []Option
	log  *zap.Logger
}

func NewInstance(opts ...Option) *Logger {
	once.Do(func() {
		l = &Logger{
			Opts: opts,
			log:  zap.L().WithOptions(zap.AddCallerSkip(1)),
		}
	})
	return l
}

func (l *Logger) DebugFunc(opts ...Option) func(ctx context.Context, msg string, args ...any) {
	return func(ctx context.Context, msg string, args ...any) {
		l.Log(ctx, opts, l.log.Debug, msg, args...)
	}
}

func (l *Logger) InfoFunc(opts ...Option) func(ctx context.Context, msg string, args ...any) {
	return func(ctx context.Context, msg string, args ...any) {
		l.Log(ctx, opts, l.log.Debug, msg, args...)
	}
}

func (l *Logger) WarnFunc(opts ...Option) func(ctx context.Context, msg string, args ...any) {
	return func(ctx context.Context, msg string, args ...any) {
		l.Log(ctx, opts, l.log.Debug, msg, args...)
	}
}

func (l *Logger) ErrorFunc(opts ...Option) func(ctx context.Context, msg string, args ...any) {
	return func(ctx context.Context, msg string, args ...any) {
		l.Log(ctx, opts, l.log.Debug, msg, args...)
	}
}

func (l *Logger) Log(ctx context.Context,
	opts []Option,
	f func(msg string, fields ...zapcore.Field),
	msg string, args ...any) {

	fields := WithTrace(ctx)

	fields = append(fields,
		NewOptions().
			WithOption(l.Opts...).
			WithOption(opts...).
			Fields()...)

	f(fmt.Sprintf(msg, args...), fields...)
}

func WithTrace(ctx context.Context) []zapcore.Field {
	traceID := core.TraceID.GetString(ctx)
	sessionID := core.SessionID.GetString(ctx)
	return []zapcore.Field{
		zap.String(core.TraceID.String(), traceID),
		zap.String(core.SessionID.String(), sessionID),
	}
}
