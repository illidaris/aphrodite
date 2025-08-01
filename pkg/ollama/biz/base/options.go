package base

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/ollama/api"
)

type options struct {
	Model     string                              `json:"mode"`
	Suffix    string                              `json:"suffix"`   // 模型响应后的后缀文本
	System    string                              `json:"system"`   // 系统提示，覆盖Modefile中的配置
	Template  string                              `json:"template"` // 模板
	Context   []int                               `json:"context,omitempty"`
	Stream    *bool                               `json:"stream,omitempty"`
	Raw       bool                                `json:"raw,omitempty"`
	Options   map[string]any                      `json:"options"`
	Think     *bool                               `json:"think,omitempty"`
	Handle    func(context.Context, *Entry) error `json:"-"`
	RawHandle func(api.GenerateResponse) error    `json:"-"`
}

func (o options) NewGenerate() *api.GenerateRequest {
	return &api.GenerateRequest{
		Model:    o.Model,
		Suffix:   o.Suffix,
		Template: o.Template,
		Context:  o.Context,
		System:   o.System,
		Stream:   o.Stream,
		Think:    o.Think,
	}
}

func newOptions(opts ...Option) options {
	o := options{
		Model:  "deepseek-r1:1.5b",
		Stream: new(bool),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

type Option func(o *options)

func WithModel(model string) Option {
	return func(o *options) {
		o.Model = model
	}
}

func WithThink(think bool) Option {
	return func(o *options) {
		o.Think = &think
	}
}

func WithSuffix(suffix string) Option {
	return func(o *options) {
		o.Suffix = suffix
	}
}

func WithStream(stream bool) Option {
	return func(o *options) {
		o.Stream = &stream
	}
}

func WithHandle(handle func(context.Context, *Entry) error) Option {
	return func(o *options) {
		o.Handle = handle
	}
}
