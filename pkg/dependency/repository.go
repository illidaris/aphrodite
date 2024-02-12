package dependency

import (
	"context"
	"time"
)

// IMQProducerRepository
type IMQProducerRepository interface {
	InsertAction(ctx context.Context, message IEventMessage) (func(context.Context) error, string)
	WaitExecWithLock(ctx context.Context, bizId, category, batch int, name string, timeout time.Duration) (string, int64, error)
	FindLockeds(ctx context.Context, locker string) ([]ITask, error)
	Clear(ctx context.Context, id string) (int64, error)
	ClearByLocker(ctx context.Context, locker string) (int64, error)
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
