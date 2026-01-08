package cache

import (
	"context"
	"net/url"
	"strings"
	"time"

	acContextex "github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/core"
	"github.com/spf13/cast"
)

type LimitOption func(option *LimitOptions)

func NewLimitOptions(opts ...LimitOption) *LimitOptions {
	option := &LimitOptions{
		prefix: KEY_LIMITER,
		busi:   KEY_DEFAULT,
		dur:    DEFAULT_DUR,
		max:    DEFAULT_MAX,
		step:   DEFAULT_STEP,
		newCtxFunc: func(src context.Context) context.Context {
			newCtx := context.Background()
			newCtx = context.WithValue(newCtx, core.SessionID, core.SessionID.Get(src))
			newCtx = context.WithValue(newCtx, core.TraceID, core.TraceID.Get(src))
			newCtx = acContextex.WithBizId(newCtx, acContextex.GetBizId(src))
			return newCtx
		},
		idFmt: func(s string) string {
			return url.QueryEscape(s)
		},
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func WithLimitCache(v dependency.ILuaCache) LimitOption {
	return func(o *LimitOptions) {
		o.cache = v
	}
}

func WithLimitPrefix(v string) LimitOption {
	return func(o *LimitOptions) {
		o.prefix = v
	}
}

func WithLimitBusi(v string) LimitOption {
	return func(o *LimitOptions) {
		o.busi = v
	}
}

func WithLimitDur(v time.Duration) LimitOption {
	return func(o *LimitOptions) {
		o.dur = v
	}
}

func WithLimitMax(v int64) LimitOption {
	return func(o *LimitOptions) {
		o.max = v
	}
}

func WithLimitStep(v int64) LimitOption {
	return func(o *LimitOptions) {
		o.step = v
	}
}

func WithLimitSkipFunc(v func(context.Context) bool) LimitOption {
	return func(o *LimitOptions) {
		o.skipFunc = v
	}
}

func WithLimitNewCtxFunc(v func(context.Context) context.Context) LimitOption {
	return func(o *LimitOptions) {
		o.newCtxFunc = v
	}
}

func WithLimitIdFmt(v func(string) string) LimitOption {
	return func(o *LimitOptions) {
		o.idFmt = v
	}
}

type LimitOptions struct {
	cache      dependency.ILuaCache // Cache instance
	prefix     string
	busi       string
	dur        time.Duration // Cache expiration duration
	max        int64
	step       int64
	skipFunc   func(context.Context) bool            // Whether to skip caching
	newCtxFunc func(context.Context) context.Context // ctx
	idFmt      func(string) string
}

func (o LimitOptions) BuildKey(bizId int64, id string) string {
	fmtId := o.idFmt(id)
	keys := []string{
		o.prefix,
		o.busi,
		cast.ToString(bizId),
		fmtId,
	}
	return strings.Join(keys, ":")
}

func (o LimitOptions) Check(ctx context.Context) (bool, error) {
	if o.cache == nil {
		return false, ErrCacheNil
	}
	if o.skipFunc == nil {
		return false, nil
	}
	return o.skipFunc(ctx), nil
}

func LimitClearIncr(ctx context.Context, bizId int64, id string, opts ...LimitOption) (int64, error) {
	option := NewLimitOptions(opts...)
	_, err := option.Check(ctx)
	if err != nil {
		return 0, err
	}
	// 构造Redis键名：前缀+业务+操作+ID 的四段式结构
	key := option.BuildKey(bizId, id)
	affect, err := option.cache.Delete(key)
	if err != nil {
		return affect, err
	}
	return affect, err
}

func LimitGetCursor(ctx context.Context, bizId int64, id string, opts ...LimitOption) (int64, error) {
	option := NewLimitOptions(opts...)
	_, err := option.Check(ctx)
	if err != nil {
		return 0, err
	}
	// 构造Redis键名：前缀+业务+操作+ID 的四段式结构
	key := option.BuildKey(bizId, id)
	res, err := option.cache.Get(key)
	if err != nil {
		return cast.ToInt64(res), err
	}
	return cast.ToInt64(res), err
}
func LimitIncr(ctx context.Context, bizId int64, id string, opts ...LimitOption) (int64, error) {
	option := NewLimitOptions(opts...)
	skip, err := option.Check(ctx)
	if err != nil {
		return 0, err
	}
	if skip {
		return 0, err
	}
	key := option.BuildKey(bizId, id)
	// 执行原子化LUA脚本操作：
	// 1. 检查当前计数是否超过max
	// 2. 未超过时增加step值
	// 3. 设置/刷新过期时间
	res, err := option.cache.EvalContext(ctx, LUA_ST_INC, []string{key},
		option.max,
		option.step,
		int64(option.dur.Seconds()))
	if err != nil {
		return 0, err
	}
	// 处理LUA脚本返回值
	resInt := cast.ToInt(res)
	if resInt == -1 {
		// 达到限流阈值时返回特定错误
		return int64(resInt), ErrLimit
	}
	return int64(resInt), nil
}
