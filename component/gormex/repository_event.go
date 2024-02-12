package gormex

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.IMQProducerRepository[dependency.IEventMessage](EventRepository[dependency.IEventMessage]{})

type EventRepository[T dependency.IEventMessage] struct {
	TaskQueueRepository[T]
}

func (r EventRepository[T]) InsertAction(ctx context.Context, t *T) (func(context.Context) error, any) {
	return func(ctx context.Context) error {
		_, err := r.BaseCreate(ctx, []*T{t})
		return err
	}, (*t).ID()
}
