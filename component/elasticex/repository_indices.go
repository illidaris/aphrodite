package elasticex

import (
	"context"
	"errors"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/olivere/elastic/v7"
)

// BaseTableCreate
func (r *BaseRepository[T]) BaseTableCreate(ctx context.Context, opts ...dependency.BaseOptionFunc) (int64, error) {
	var (
		t T
	)
	mapping, ok := any(t).(dependency.IMapping)
	if !ok {
		return 0, nil
	}
	res, err := r.GetIndicesCreateService(ctx, opts...).BodyString(mapping.GetMapping()).Do(ctx)
	if err != nil {
		// Handle error
		return 0, err
	}
	if !res.Acknowledged {
		// Not acknowledged
		return 0, errors.New("no acknowledged")
	}
	return 1, err
}

// BaseTableCreate
func (r *BaseRepository[T]) BaseTableExists(ctx context.Context, opts ...dependency.BaseOptionFunc) (bool, error) {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	// Check if the index called "twitter" exists
	return db.IndexExists(t.TableName()).Do(ctx)
}

// GetIndicesCreateService
func (r *BaseRepository[T]) GetIndicesCreateService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.IndicesCreateService {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.CreateIndex(opt.GetTableName(t))
	return srv
}
