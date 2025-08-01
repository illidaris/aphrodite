package classic

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/ollama/biz/base"

	group "github.com/illidaris/aphrodite/pkg/group/v2"
)

func ClassicFunc(host string, template string, opts ...base.Option) func(ctx context.Context, categories []string, labels []Label, srcs ...*base.Entry) error {
	g := base.Generater(host)
	return func(ctx context.Context, categories []string, labels []Label, srcs ...*base.Entry) error {
		for _, v := range srcs {
			v.Prompt = buildPrompt(template, labels, categories, v.Content)
		}
		_, _ = group.GroupFunc(func(vs ...*base.Entry) (int64, error) {
			for _, v := range vs {
				_ = g(ctx, v, opts...)
			}
			return 0, nil
		}, srcs)
		return nil
	}
}
