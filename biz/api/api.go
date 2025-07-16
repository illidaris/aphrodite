package api

import (
	"context"
	"time"

	"github.com/illidaris/aphrodite/pkg/exception"
	restSender "github.com/illidaris/rest/sender"
	"github.com/illidaris/rest/signature"
)

func POST[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewPostAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func FORM[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewFormAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func PUT[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewPutAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func GET[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewGetAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func DELETE[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (Out, exception.Exception) {
	return func(ctx context.Context, req In) (Out, exception.Exception) {
		r := NewDeleteAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return r.Response.Data, exception.ERR_BUSI.Wrap(err)
		}
		return r.Response.Data, r.Response.ToException()
	}
}

func invoke(ctx context.Context, req restSender.IRequest, host, secret string, timeout time.Duration) error {
	opts := []restSender.Option{
		restSender.WithTimeout(timeout),
		restSender.WithHost(host),
	}
	if secret != "" {
		opts = append(opts, restSender.WithSignSetMode(signature.SignSetlInURL, secret, signature.Generate))
	}
	s := restSender.NewSender(opts...)
	_, err := s.Invoke(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
