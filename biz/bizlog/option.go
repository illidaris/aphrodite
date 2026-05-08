package bizlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// OptionFunc 用于配置 Options 的函数选项类型
type Option func(*Options)

func NewOptions() *Options {
	return &Options{
		Entry: Entry{
			Tags: []string{},
		},
	}
}

// Options 业务日志可选配置
type Options struct {
	Entry
}

func (o *Options) WithOption(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) Fields() []zapcore.Field {
	return []zapcore.Field{
		zap.Int32("log_dst", o.Dst),
		zap.Int64("ibzl_uid", o.Uid),
		zap.String("ibzl_pid", o.Pid),
		zap.String("ibzl_name", o.Name),
		zap.Int64("ibzl_bizid", o.BizId),
		zap.String("ibzl_system", o.System),
		zap.Int32("ibzl_cgy", o.Category),
		zap.Int64("ibzl_action", o.Action),
		zap.String("ibzl_msg", o.Message),
		zap.Strings("ibzl_tags", o.Tags),
		zap.Int64("ibzl_opat", o.OpAt),
		zap.Int64("ibzl_createat", o.CreateAt),
	}
}

// WithDst 日志分类
func WithDst(v int32) Option {
	return func(o *Options) {
		o.Dst = v
	}
}

// WithUid 设置主角 Id（WHO）
func WithUid(uid int64) Option {
	return func(o *Options) {
		o.Uid = uid
	}
}

// WithPid 设置主角 PID（WHO）
func WithPid(pid string) Option {
	return func(o *Options) {
		o.Pid = pid
	}
}

// WithName 设置主角名（WHO）
func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// WithBizId 设置业务 Id（WHERE）
func WithBizId(bizId int64) Option {
	return func(o *Options) {
		o.BizId = bizId
	}
}

// WithSystem 设置所属系统（WHERE）
func WithSystem(system string) Option {
	return func(o *Options) {
		o.System = system
	}
}

// WithCategory 设置大类
func WithCategory(category int32) Option {
	return func(o *Options) {
		o.Category = category
	}
}

// WithAction 设置行为
func WithAction(action int64) Option {
	return func(o *Options) {
		o.Action = action
	}
}

// WithMessage 设置事件描述（WHAT）
func WithMessage(message string) Option {
	return func(o *Options) {
		o.Message = message
	}
}

// WithTags 追加标签
func WithTags(tags ...string) Option {
	return func(o *Options) {
		o.Tags = append(o.Tags, tags...)
	}
}

// WithOpAt 设置事件发生时间点（WHEN）
func WithOpAt(opAt int64) Option {
	return func(o *Options) {
		o.OpAt = opAt
	}
}

// WithCreateAt 设置入库时间点
func WithCreateAt(createAt int64) Option {
	return func(o *Options) {
		o.CreateAt = createAt
	}
}
