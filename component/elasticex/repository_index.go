package elasticex

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
)

// BaseCreate
func (r *BaseRepository[T]) BaseCreate(ctx context.Context, ts []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	switch len(ts) {
	case 0:
		return 0, nil
	case 1:
		return r.BaseCreateOne(ctx, ts[0], opts...)
	default:
		return r.BaseCreates(ctx, ts, opts...)
	}
}

// BaseCreate
func (r *BaseRepository[T]) BaseCreateOne(ctx context.Context, t *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	srv := r.GetIndexService(ctx, opts...)
	p, ok := any(t).(dependency.IEntity)
	if ok && p.ID() != nil {
		srv = srv.Id(cast.ToString(p.ID()))
	}
	_, err := srv.
		BodyJson(t).
		Do(ctx)
	if err != nil {
		// Handle error
		return 0, err
	}
	// println(res.Shards.Total)
	return 1, err
}

// BaseCreates
func (r *BaseRepository[T]) BaseCreates(ctx context.Context, ts []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	bulkRequest := db.Bulk()
	for _, t := range ts {
		p, ok := any(t).(dependency.IEntity)
		if !ok {
			continue
		}
		srv := elastic.NewBulkIndexRequest().Index(opt.GetTableName(p))
		if p.ID() != nil {
			srv = srv.Id(cast.ToString(p.ID()))
		}
		req := srv.Doc(t)
		bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(ctx)
	if err != nil {
		// Handle error
		return 0, err
	}
	// println(res.Shards.Total)
	return int64(len(res.Succeeded())), err
}

// BaseSave
func (r *BaseRepository[T]) BaseSave(ctx context.Context, ts []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	return r.BaseCreate(ctx, ts, opts...)
}

// BaseUpdate
func (r *BaseRepository[T]) BaseUpdate(ctx context.Context, t *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	srv := r.GetIndexService(ctx, opts...)
	p, ok := any(t).(dependency.IEntity)
	if ok && p.ID() != nil {
		srv = srv.Id(cast.ToString(p.ID()))
	} else {
		return 0, nil
	}
	_, err := srv.
		BodyJson(t).
		Do(ctx)
	if err != nil {
		// Handle error
		return 0, err
	}
	// println(res.Shards.Total)
	return 1, err
}

// BaseDelete
func (r *BaseRepository[T]) BaseDelete(ctx context.Context, t *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	srv := r.GetDeleteService(ctx, opts...)
	p, ok := any(t).(dependency.IEntity)
	if ok && p.ID() != nil {
		srv = srv.Id(cast.ToString(p.ID()))
	} else {
		return 0, nil
	}
	_, err := srv.Do(ctx)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// GetIndexService
func (r *BaseRepository[T]) GetIndexService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.IndexService {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.Index().Index(opt.GetTableName(t))
	return srv
}

// GetDeleteService
func (r *BaseRepository[T]) GetDeleteService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.DeleteService {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.Delete().Index(opt.GetTableName(t))
	return srv
}
