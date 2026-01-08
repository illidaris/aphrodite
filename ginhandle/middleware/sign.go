package middleware

import (
	"context"
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/illidaris/rest/signature"

	"github.com/spf13/cast"
)

// 用于前端签名使用，最好配合mojito_bg.wasm
func WebSignMiddleware(sopts ...WebsignOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		opts := NewWebsignOptions(sopts...)
		app := c.Query(signature.SignAppID)
		ver := c.Query("ver")
		tsStr := c.Query(signature.SignKeyTimestamp)
		// 校验签名
		if !opts.Skip(ctx) {
			secret := opts.SecretFunc(ctx, app)
			if secret == "" {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_APP.New("应用访问被拒绝")))
				return
			}
			if !opts.AllowVerFunc(ctx, ver) {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_VER.New("签名版本失效")))
				return
			}
			timeout := opts.TimeoutFunc(ctx)
			now := time.Now()
			beg := now.Add(-1 * timeout)
			end := now.Add(timeout)
			if ts := cast.ToInt64(tsStr); beg.Unix() > ts || end.Unix() < ts {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_EXPIRED.New("签名已过期")))
				return
			}
			realOpts := []signature.OptionFunc{
				signature.WithHmacFunc(func(s string, as ...string) string {
					return signature.HashMac(sha256.New, s, as...)
				}),
				signature.WithExpire(timeout),
				signature.WithSecret(secret),
			}
			realOpts = append(realOpts, opts.RestOptions...)
			err := signature.VerifySign(c.Request, realOpts...)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
				return
			}
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

type WebsignOption func(*WebsignOptions)

// NewWebsignOptions 创建并返回一个新的WebsignOptions实例，初始化允许的主机和版本为空，超时时间为2分钟，签名不可用。
func NewWebsignOptions(opts ...WebsignOption) *WebsignOptions {
	o := &WebsignOptions{
		AllowHostFunc: func(ctx context.Context, s string) bool { return true },
		AllowVerFunc:  func(ctx context.Context, s string) bool { return true },
		TimeoutFunc:   func(ctx context.Context) time.Duration { return time.Minute * 3 },
		SecretFunc:    func(ctx context.Context, s string) string { return "" },
		RestOptions:   []signature.OptionFunc{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithSkipFunc(v func(context.Context) bool) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.SkipFunc = v
	}
}

func WithSecretFunc(v func(ctx context.Context, s string) string) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.SecretFunc = v
	}
}

func WithExpireFunc(v func(ctx context.Context) time.Duration) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.TimeoutFunc = v
	}
}

func WithHostFunc(v func(ctx context.Context, s string) bool) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.AllowHostFunc = v
	}
}

func WithVerFunc(v func(ctx context.Context, s string) bool) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.AllowVerFunc = v
	}
}

func WithRestOptions(vs ...signature.OptionFunc) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.RestOptions = append(opts.RestOptions, vs...)
	}
}

// WebsignOptions 定义了Web签名选项的结构体。
// 其中包括是否启用签名、签名密钥、超时时间、允许的主机和允许的版本。
type WebsignOptions struct {
	SecretFunc    func(context.Context, string) string // 签名所用的密钥
	TimeoutFunc   func(context.Context) time.Duration  // 请求超时时间
	AllowHostFunc func(context.Context, string) bool   // 允许签名的主机名集合
	AllowVerFunc  func(context.Context, string) bool   // 允许签名的版本集合
	RestOptions   []signature.OptionFunc               // 框架参数
	SkipFunc      func(context.Context) bool
}

func (opts *WebsignOptions) Skip(ctx context.Context) bool {
	if opts.SkipFunc == nil {
		return false
	}
	return opts.SkipFunc(ctx)
}
