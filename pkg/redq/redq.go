package redq

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

func InitRedqSrv(opts ...RedqOption) error {
	o := NewRedqOptions(opts...)
	srv := asynq.NewServer(
		o.RedisConfig,
		o.QueueConfig,
	)
	mux := asynq.NewServeMux()
	for topic, sub := range o.HandleMap {
		mux.HandleFunc(topic, func(ctx context.Context, t *asynq.Task) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("[%v]panic is %v", topic, r)
				}
			}()
			if t == nil {
				err = fmt.Errorf("[%v]asynq.Task is nil", topic)
				return
			}
			err = sub(ctx, t.Payload())
			return
		})
	}
	return srv.Run(mux)
}

func SendFunc(opts ...RedqOption) func(ctx context.Context, topic string, val []byte, asyncopts ...asynq.Option) (*asynq.TaskInfo, error) {
	o := NewRedqOptions(opts...)
	client := asynq.NewClient(o.RedisConfig)
	return func(ctx context.Context, topic string, val []byte, asyncopts ...asynq.Option) (*asynq.TaskInfo, error) {
		t := asynq.NewTask(topic, val, asyncopts...)
		return client.EnqueueContext(ctx, t)
	}
}
