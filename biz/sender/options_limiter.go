package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/spf13/cast"
)

type LimiterOption func(*Limiter)

func WithLimitCode(code, name string) LimiterOption {
	return func(l *Limiter) {
		l.Code = code
		l.Name = name
		if l.Name == "" {
			l.Name = l.Code
		}
	}
}
func WithLimitTtl(v time.Duration) LimiterOption {
	return func(l *Limiter) {
		l.Ttl = v
	}
}

func WithLimitInterval(v time.Duration) LimiterOption {
	return func(l *Limiter) {
		l.Interval = v
	}
}
func WithLimitMax(v uint64) LimiterOption {
	return func(l *Limiter) {
		l.Max = v
	}
}
func WithLimitStep(v uint64) LimiterOption {
	return func(l *Limiter) {
		l.Step = v
	}
}

var _ = ILimiter(&Limiter{})

func NewLimiter(opts ...LimiterOption) ILimiter {
	l := &Limiter{
		Code:     KEY_SECTION_LIMITERS_BASE,
		Name:     "验证码",
		Ttl:      time.Hour * 24,
		Interval: time.Minute * 1,
		Step:     1,
		Max:      5,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

type ILimiter interface {
	WithStore(store IStore) ILimiter
	Check(ctx context.Context) (int64, exception.Exception)
	Exec(ctx context.Context) (int64, exception.Exception)
	Clear(ctx context.Context) (int64, exception.Exception)
}

type Limiter struct {
	Code     string        // code 唯一编码
	Name     string        // 名称
	Ttl      time.Duration // 发送限制周期
	Interval time.Duration // 发送每次间隔
	Step     uint64        // 步长
	Max      uint64        // 累计限制数

	BuldKeyFunc func(section string) string
	Store       IStore // 持久化
}

func (l *Limiter) WithStore(store IStore) ILimiter {
	l.Store = store
	return l
}

func (l Limiter) LimitError() exception.Exception {
	msg := fmt.Sprintf("%v达到发送上限", l.Name)
	return exception.ERR_VERIFYCODE_SENDLIMIT.New(msg)
}

func (l Limiter) GetCode() string {
	return l.Code
}

func (l Limiter) Check(ctx context.Context) (int64, exception.Exception) {
	if l.Store == nil {
		return 0, ExStoreNil
	}
	// 获取当前值
	key := l.BuldKeyFunc(l.Code)
	v, err := l.Store.GetContext(ctx, key)
	cur := cast.ToInt64(v)
	if err != nil {
		return cur, exception.ERR_VERIFYCODE_SENDLIMIT.Wrap(err)
	}
	// 判断逻辑
	if cur+int64(l.Step) >= int64(l.Max) {
		return cur, l.LimitError()
	}
	return cur, nil
}

func (l Limiter) Exec(ctx context.Context) (int64, exception.Exception) {
	if l.Store == nil {
		return 0, ExStoreNil
	}
	// 获取当前值
	key := l.BuldKeyFunc(l.Code)
	v, err := l.Store.EvalContext(ctx, LUA_ST_INC, []string{key}, l.Max, l.Step, int64(l.Ttl.Seconds()))
	cur := cast.ToInt64(v)
	if err != nil {
		return cur, exception.ERR_VERIFYCODE_SENDLIMIT.Wrap(err)
	}
	// 判断自增计算结果
	if cur == -1 {
		return 0, l.LimitError()
	}
	return cur, nil
}

func (l Limiter) Clear(ctx context.Context) (int64, exception.Exception) {
	if l.Store == nil {
		return 0, ExStoreNil
	}
	// 获取当前值
	key := l.BuldKeyFunc(l.Code)
	v, err := l.Store.DelContext(ctx, key)
	cur := cast.ToInt64(v)
	if err != nil {
		return cur, exception.ERR_VERIFYCODE_SENDLIMIT.Wrap(err)
	}
	return cur, nil
}
