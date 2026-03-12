package redq

import (
	"context"

	"github.com/hibiken/asynq"
)

type SubscribeHandler func(context.Context, []byte) error
type RedqOption func(*RedqOptions)

type RedqOptions struct {
	RedisConfig asynq.RedisClientOpt
	QueueConfig asynq.Config
	HandleMap   map[string]SubscribeHandler
}

func WithRedisConfig(cfg asynq.RedisClientOpt) RedqOption {
	return func(o *RedqOptions) {
		o.RedisConfig = cfg
	}
}

func WithQueueConfig(cfg asynq.Config) RedqOption {
	return func(o *RedqOptions) {
		o.QueueConfig = cfg
	}
}

func WithHandle(topic string, handler SubscribeHandler) RedqOption {
	return func(o *RedqOptions) {
		if o.HandleMap == nil {
			o.HandleMap = map[string]SubscribeHandler{}
		}
		o.HandleMap[topic] = handler
	}
}

func NewRedqOptions(opts ...RedqOption) *RedqOptions {
	opt := &RedqOptions{}
	for _, f := range opts {
		f(opt)
	}
	return opt
}
