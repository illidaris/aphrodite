package api

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/exception"
	restSender "github.com/illidaris/rest/sender"
	"github.com/illidaris/rest/signature"
)

func POST[In IRequest, Out any](host string, opts ...Option) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewPostAPI[In, Out](req)
		err := invoke(ctx, r, host, opts...)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func FORM[In IRequest, Out any](host string, opts ...Option) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewFormAPI[In, Out](req)
		err := invoke(ctx, r, host, opts...)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func PUT[In IRequest, Out any](host string, opts ...Option) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewPutAPI[In, Out](req)
		err := invoke(ctx, r, host, opts...)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func GET[In IRequest, Out any](host string, opts ...Option) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewGetAPI[In, Out](req)
		err := invoke(ctx, r, host, opts...)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func DELETE[In IRequest, Out any](host string, opts ...Option) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewDeleteAPI[In, Out](req)
		err := invoke(ctx, r, host, opts...)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func invoke(ctx context.Context, req restSender.IRequest, host string, opts ...Option) error {
	o := newOptions(opts...)
	rsOpts := []restSender.Option{
		restSender.WithTimeout(o.Timeout),
		restSender.WithHost(host),
	}

	if o.Secret != "" {
		rsOpts = append(rsOpts, restSender.WithSignSetMode(signature.SignSetlInURL, o.Secret, signature.Generate))
	}

	if len(o.RsOptions) > 0 {
		rsOpts = append(rsOpts, o.RsOptions...)
	}
	s := restSender.NewSender(rsOpts...)

	if o.RequestFunc != nil {
		o.RequestFunc(ctx, req)
	}

	resp, err := s.Invoke(ctx, req)

	if o.ResponseFunc != nil {
		o.ResponseFunc(ctx, resp, err)
	}

	if err != nil {
		return err
	}
	return nil
}
