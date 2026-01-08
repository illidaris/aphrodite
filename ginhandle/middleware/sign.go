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
			secret, ok := opts.Secrets[app]
			if !ok {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_APP.New("应用访问被拒绝")))
				return
			}
			if !opts.AllowVer(ver) {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_VER.New("签名版本失效")))
				return
			}
			now := time.Now()
			beg := now.Add(-1 * opts.Timeout)
			end := now.Add(opts.Timeout)
			if ts := cast.ToInt64(tsStr); beg.Unix() > ts || end.Unix() < ts {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_SIGN_EXPIRED.New("签名已过期")))
				return
			}
			realOpts := []signature.OptionFunc{
				signature.WithHmacFunc(func(s string, as ...string) string {
					return signature.HashMac(sha256.New, s, as...)
				}),
				signature.WithExpire(opts.Timeout),
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
		AllowHosts:  map[string]struct{}{},
		AllowVers:   map[string]struct{}{},
		Timeout:     time.Minute * 3,
		Secrets:     map[string]string{},
		RestOptions: []signature.OptionFunc{},
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

func WithSecret(app, secret string) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.Secrets[app] = secret
	}
}

func WithSecrets(secrets map[string]string) WebsignOption {
	return func(opts *WebsignOptions) {
		for k, v := range secrets {
			opts.Secrets[k] = v
		}
	}
}

func WithExpire(timeout time.Duration) WebsignOption {
	return func(opts *WebsignOptions) {
		opts.Timeout = timeout
	}
}

func WithHosts(hosts ...string) WebsignOption {
	return func(opts *WebsignOptions) {
		for _, v := range hosts {
			opts.AllowHosts[v] = struct{}{}
		}
	}
}

func WithVers(vers ...string) WebsignOption {
	return func(opts *WebsignOptions) {
		for _, v := range vers {
			opts.AllowVers[v] = struct{}{}
		}
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
	Secrets     map[string]string      // 签名所用的密钥
	Timeout     time.Duration          // 请求超时时间
	AllowHosts  map[string]struct{}    // 允许签名的主机名集合
	AllowVers   map[string]struct{}    // 允许签名的版本集合
	RestOptions []signature.OptionFunc // 框架参数
	SkipFunc    func(context.Context) bool
}

func (opts *WebsignOptions) Skip(ctx context.Context) bool {
	if opts.SkipFunc == nil {
		return false
	}
	return opts.SkipFunc(ctx)
}

// AllowHost 判断给定的主机名是否在允许的主机名集合中。
// v: 要判断的主机名
// 返回值表示是否允许该主机名进行签名。
func (opts *WebsignOptions) AllowHost(v string) bool {
	if opts.AllowHosts == nil {
		return false
	}
	for k := range opts.AllowHosts {
		if k == "*" {
			return true
		}
	}
	if len(v) == 0 {
		return false
	}
	_, ok := opts.AllowHosts[v]
	return ok
}

// AllowVer 判断给定的版本是否在允许的版本集合中。
// v: 要判断的版本号
// 返回值表示是否允许该版本进行签名。
func (opts *WebsignOptions) AllowVer(v string) bool {
	if len(v) == 0 {
		return false
	}
	if opts.AllowVers == nil {
		return false
	}
	_, ok := opts.AllowVers[v]
	return ok
}
