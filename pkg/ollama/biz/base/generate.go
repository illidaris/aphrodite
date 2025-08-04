package base

import (
	"context"
	"net/http"
	"net/url"

	"github.com/illidaris/aphrodite/pkg/ollama/api"
)

func Generater(raw string) func(context.Context, *Entry, ...Option) error {
	urlV, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil
	}
	client := api.NewClient(urlV, http.DefaultClient)
	return func(ctx context.Context, entry *Entry, opts ...Option) error {
		options := newOptions(opts...)
		req := options.NewGenerate()
		req.Prompt = entry.Prompt
		// 如果存在原始委托函数
		if options.RawHandle != nil {
			return client.Generate(ctx, req, options.RawHandle)
		}
		respFunc := func(resp api.GenerateResponse) error {
			entry.Duration = resp.TotalDuration.Milliseconds()
			entry.Result = resp.Response
			if options.Handle == nil {
				return nil
			}
			return options.Handle(ctx, entry)
		}
		err = client.Generate(ctx, req, respFunc)
		if err != nil {
			return err
		}
		return nil
	}
}
