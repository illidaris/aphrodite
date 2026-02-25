package oauth2

import (
	"context"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type Handle func(context.Context) string
type Option func(*options)
type options struct {
	GetAuthUrlHandle             Handle
	GetTokenUrlHandle            Handle
	GetClientIdHandle            Handle
	GetClientSecretHandle        Handle
	GetRedirectUrlHandle         Handle
	GetScopesHandle              func(context.Context) []string
	GetBusiSecretHandle          Handle
	GetCodeChallengeExpireHandle func(context.Context) time.Duration
	Cache                        ICache
}

func (opt options) GetOAuth2Config(ctx context.Context) *oauth2.Config {
	var oauthConfig = &oauth2.Config{
		ClientID:     opt.GetClientIdHandle(ctx),
		ClientSecret: opt.GetClientSecretHandle(ctx),
		RedirectURL:  opt.GetRedirectUrlHandle(ctx),
		Endpoint: oauth2.Endpoint{
			AuthURL:  opt.GetAuthUrlHandle(ctx),
			TokenURL: opt.GetTokenUrlHandle(ctx),
		},
		Scopes: opt.GetScopesHandle(ctx),
	}
	return oauthConfig
}

func NewOptions(opts ...Option) *options {
	o := &options{
		GetAuthUrlHandle:      func(ctx context.Context) string { return "" },
		GetTokenUrlHandle:     func(ctx context.Context) string { return "" },
		GetClientIdHandle:     func(context.Context) string { return "" },
		GetClientSecretHandle: func(context.Context) string { return "" },
		GetRedirectUrlHandle:  func(context.Context) string { return "" },
		GetScopesHandle:       func(context.Context) []string { return []string{} },
		GetBusiSecretHandle:   func(context.Context) string { return "" },
		GetCodeChallengeExpireHandle: func(context.Context) time.Duration {
			return 5 * time.Minute
		},
		Cache: defaultCache{&sync.Map{}},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithGetAuthUrlHandle(handle Handle) Option {
	return func(o *options) {
		o.GetAuthUrlHandle = handle
	}
}

func WithGetTokenUrlHandle(handle Handle) Option {
	return func(o *options) {
		o.GetTokenUrlHandle = handle
	}
}

func WithGetClientIdHandle(handle Handle) Option {
	return func(o *options) {
		o.GetClientIdHandle = handle
	}
}

func WithGetClientSecretHandle(handle Handle) Option {
	return func(o *options) {
		o.GetClientSecretHandle = handle
	}
}

func WithGetRedirectUrlHandle(handle Handle) Option {
	return func(o *options) {
		o.GetRedirectUrlHandle = handle
	}
}

func WithGetScopesHandle(handle func(context.Context) []string) Option {
	return func(o *options) {
		o.GetScopesHandle = handle
	}
}

func WithGetBusiSecretHandle(handle Handle) Option {
	return func(o *options) {
		o.GetBusiSecretHandle = handle
	}
}

func WithGetCodeChallengeExpireHandle(handle func(context.Context) time.Duration) Option {
	return func(o *options) {
		o.GetCodeChallengeExpireHandle = handle
	}
}

func WithCache(cache ICache) Option {
	return func(o *options) {
		o.Cache = cache
	}
}
