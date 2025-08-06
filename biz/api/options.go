package api

import (
	"context"
	"time"

	restSender "github.com/illidaris/rest/sender"
)

func newOptions(opts ...Option) *options {
	opt := &options{
		Timeout:   5 * time.Second,
		RsOptions: []restSender.Option{},
	}
	for _, f := range opts {
		f(opt)
	}
	return opt
}

type options struct {
	Timeout      time.Duration
	Secret       string
	RequestFunc  func(context.Context, restSender.IRequest)
	ResponseFunc func(context.Context, any, error)
	RsOptions    []restSender.Option
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.Timeout = timeout
	}
}

func WithSecret(secret string) Option {
	return func(o *options) {
		o.Secret = secret
	}
}

func WithRequestFunc(f func(context.Context, restSender.IRequest)) Option {
	return func(o *options) {
		o.RequestFunc = f
	}
}

func WithResponseFunc(f func(context.Context, any, error)) Option {
	return func(o *options) {
		o.ResponseFunc = f
	}
}

func WithRsOptions(opts ...restSender.Option) Option {
	return func(o *options) {
		o.RsOptions = append(o.RsOptions, opts...)
	}
}
