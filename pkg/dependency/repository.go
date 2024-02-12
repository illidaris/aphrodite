package dependency

import (
	"context"
)

// IMQProducerRepository
type IMQProducerRepository[T IEventMessage] interface {
	IRepository[T]
	InsertAction(ctx context.Context, db string, t *T) func(context.Context) error
	WaitExecWithLock(ctx context.Context, t T, batch int) (string, int64, error)
	FindLockeds(ctx context.Context, locker string) ([]T, error)
}

// action
type DbAction func(ctx context.Context) error

// IUnitOfWork  trans
type IUnitOfWork interface {
	Execute(ctx context.Context, fs ...DbAction) (e error)
}

// IRepository repo
type IRepository[T IEntity] interface {
	BaseCreate(ctx context.Context, ps []*T, opts ...BaseOptionFunc) (int64, error)
	BaseSave(ctx context.Context, ps []*T, opts ...BaseOptionFunc) (int64, error)
	BaseUpdate(ctx context.Context, p *T, opts ...BaseOptionFunc) (int64, error)
	BaseGet(ctx context.Context, opts ...BaseOptionFunc) (*T, error)
	BaseDelete(ctx context.Context, p *T, opts ...BaseOptionFunc) (int64, error)
	BaseCount(ctx context.Context, opts ...BaseOptionFunc) (int64, error)
	BaseQuery(ctx context.Context, opts ...BaseOptionFunc) ([]T, error)
	BaseQueryWithCount(ctx context.Context, opts ...BaseOptionFunc) ([]T, int64, error)
}
