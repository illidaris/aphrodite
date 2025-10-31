package sender

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/spf13/cast"
)

func NewSenderOptions(opts ...SenderOption) *SenderOptions {
	o := &SenderOptions{
		App:            globalApp,
		Busi:           "def",
		Name:           "name",
		Limiters:       []Limiter{},
		CodeLen:        6,
		CodeTtl:        time.Minute * 5,
		CodeGenateFunc: RandVerifyCode,
		SendFunc: func(ctx context.Context, a any, s string) error {
			return nil
		},
		ArgsFmtFunc: func(v any) string {
			return base64.RawURLEncoding.EncodeToString([]byte(cast.ToString(v)))
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type SenderOption func(*SenderOptions)

type SenderOptions struct {
	App            string           // app
	Busi           string           // 业务类型标识
	Name           string           // 名称
	Limiters       []Limiter        // 限制器
	CodeLimiter    Limiter          // 验证码限制器
	CodeLen        uint             // 随机码长度
	CodeTtl        time.Duration    // 随机码有效期
	CodeGenateFunc CodeGenateHandle // 函数：生成指定的随机码
	ArgsFmtFunc    ArgsFmtHandle    // 函数：Args格式化
	SendFunc       NotifyMsgHandle  // 函数：发送消息
	Store          IStore           // 持久化
}

func (i SenderOptions) CheckExistLocked(ctx context.Context, id any) (int64, exception.Exception) {
	if i.Store == nil {
		return 0, ExStoreNil
	}
	bizId := contextex.GetBizId(ctx) // 获取业务标识
	key := i.BuildKeyLocked(bizId, id)
	ttl, err := i.Store.TTLContext(ctx, key)
	if err != nil {
		return int64(ttl.Seconds()), exception.ERR_VERIFYCODE_SENDFAIL.Wrap(err)
	}
	return int64(ttl.Seconds()), nil
}

func (i SenderOptions) SetLocked(ctx context.Context, id any, code string) (int64, bool, exception.Exception) {
	if i.Store == nil {
		return 0, false, ExStoreNil
	}
	bizId := contextex.GetBizId(ctx) // 获取业务标识
	// 设置间隔时间
	lockDur := i.CodeLimiter.Interval
	lockedKey := i.BuildKeyLocked(bizId, id)
	b, err := i.Store.SetNXContext(ctx, lockedKey, code, lockDur)
	if err != nil {
		return 0, false, exception.ERR_VERIFYCODE_SENDFAIL.Wrap(err)
	}
	// 在发送间隔内，存在已经发送消息
	if !b {
		t, err := i.Store.TTLContext(ctx, lockedKey)
		if err != nil {
			return int64(t.Seconds()), true, exception.ERR_VERIFYCODE_HASEXPIRED.Wrap(err)
		}
		return int64(t.Seconds()), true, nil
	}
	return int64(lockDur.Seconds()), false, nil
}

func (i SenderOptions) GetCode(ctx context.Context, id any) (string, exception.Exception) {
	if i.Store == nil {
		return "", ExStoreNil
	}
	bizId := contextex.GetBizId(ctx) // 获取业务标识
	key := i.BuildKeyCode(bizId, id)
	v, err := i.Store.GetContext(ctx, key)
	if err != nil {
		return cast.ToString(v), exception.ERR_BUSI.Wrap(err)
	}
	return cast.ToString(v), nil
}

func (i SenderOptions) SetCode(ctx context.Context, id any, code string) (int64, exception.Exception) {
	if i.Store == nil {
		return 0, ExStoreNil
	}
	bizId := contextex.GetBizId(ctx) // 获取业务标识
	dur := i.CodeLimiter.Ttl
	key := i.BuildKeyCode(bizId, id)
	v, err := i.Store.SetContext(ctx, key, code, dur)
	if err != nil {
		return cast.ToInt64(v), exception.ERR_BUSI.Wrap(err)
	}
	return cast.ToInt64(v), nil
}

// ============================ Build Key =============================
func (i SenderOptions) BuildKeyCode(bizId int64, id any) string {
	return i.BuildKey(bizId, KEY_SECTION_CODE, id)
}

// 构建间隔时间的KEY
func (i SenderOptions) BuildKeyLocked(bizId int64, id any) string {
	return i.BuildKey(bizId, KEY_SECTION_LOCKED, id)
}
func (i SenderOptions) BuildKey(bizId int64, section string, id any) string {
	keys := []string{i.App, cast.ToString(bizId), i.Busi, section}
	if i.ArgsFmtFunc != nil {
		keys = append(keys, i.ArgsFmtFunc(id))
	} else {
		keys = append(keys, cast.ToString(id))
	}
	return strings.Join(keys, ":")
}
