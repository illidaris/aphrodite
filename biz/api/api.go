package api

import (
	"context"
	"fmt"
	"time"

	restSender "github.com/illidaris/rest/sender"
	"github.com/illidaris/rest/signature"
)

func POST[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (*Out, error) {
	return func(ctx context.Context, req In) (*Out, error) {
		r := NewPostAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, timeout)
		if err != nil {
			return nil, err
		}
		if r.Response == nil {
			return nil, fmt.Errorf("[POST]%v resp is nil", req.GetPath())
		}
		return r.Response.Data, nil
	}
}

func PUT[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (*Out, error) {
	return func(ctx context.Context, req In) (*Out, error) {
		r := NewPutAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, time.Second*5)
		if err != nil {
			return nil, err
		}
		if r.Response == nil {
			return nil, fmt.Errorf("[PUT]%v resp is nil", req.GetPath())
		}
		return r.Response.Data, nil
	}
}

func GET[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (*Out, error) {
	return func(ctx context.Context, req In) (*Out, error) {
		r := NewGetAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, time.Second*5)
		if err != nil {
			return nil, err
		}
		if r.Response == nil {
			return nil, fmt.Errorf("[GET]%v resp is nil", req.GetPath())
		}
		return r.Response.Data, nil
	}
}

func DELETE[In IRequest, Out any](host, secret string, timeout time.Duration) func(ctx context.Context, req In) (*Out, error) {
	return func(ctx context.Context, req In) (*Out, error) {
		r := NewDeleteAPI[In, Out](req)
		err := invoke(ctx, r, host, secret, time.Second*5)
		if err != nil {
			return nil, err
		}
		if r.Response == nil {
			return nil, fmt.Errorf("[DELETE]%v resp is nil", req.GetPath())
		}
		return r.Response.Data, nil
	}
}

func invoke(ctx context.Context, req restSender.IRequest, host, secret string, timeout time.Duration) error {
	opts := []restSender.Option{
		restSender.WithTimeout(timeout),
		restSender.WithTimeConsume(true),
		restSender.WithHost(host),
	}
	if secret == "" {
		opts = append(opts, restSender.WithSignSetMode(signature.SignSetlInURL, secret, signature.Generate))
	}
	s := restSender.NewSender(opts...)
	_, err := s.Invoke(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
