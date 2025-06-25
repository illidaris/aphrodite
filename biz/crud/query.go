package crud

import (
	"context"

	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"github.com/illidaris/aphrodite/pkg/exception"
)

func PagingListFunc[T dependency.IEntity, K dto.IRow](repo dependency.IRepository[T], iterater func(T) *K) func(ctx context.Context, req dependency.ICondPage) (*dto.RecordPtrPager[K], exception.Exception) {
	return func(ctx context.Context, req dependency.ICondPage) (*dto.RecordPtrPager[K], exception.Exception) {
		result := new(dto.RecordPtrPager[K])
		result.PageIndex = req.GetPageIndex()
		result.PageSize = req.GetPageSize()
		result.Data = make([]*K, 0)

		opts := []dependency.BaseOptionFunc{}
		if keys := req.GetDbShardingKeys(); len(keys) > 0 {
			opts = append(opts, dependency.WithDbShardingKey(keys...))
		}
		if keys := req.GetTbShardingKeys(); len(keys) > 0 {
			opts = append(opts, dependency.WithTbShardingKey(keys...))
		}
		if keys := req.GetConds(); len(keys) > 0 {
			opts = append(opts, dependency.WithConds(keys...))
		}
		opts = append(opts,
			dependency.WithPage(req),
			dependency.WithReadOnly(true),
		)

		ps, total, err := repo.BaseQueryWithCount(ctx, opts...)
		if err != nil {
			return result, exception.ERR_BUSI.Wrap(err)
		}
		result.TotalRecord = total
		result.Paginator()

		for _, v := range ps {
			ptr := iterater(v)
			result.Data = append(result.Data, ptr)
		}
		return result, nil
	}
}

func ListFunc[T dependency.IEntity, K dto.IRow](repo dependency.IRepository[T], iterater func(T) *K) func(ctx context.Context, req dependency.ICond) ([]*K, exception.Exception) {
	f := EntitiesFunc(repo)
	return func(ctx context.Context, req dependency.ICond) ([]*K, exception.Exception) {
		result := []*K{}
		ps, ex := f(ctx, req.GetDbShardingKeys(), req.GetTbShardingKeys(), req.GetConds()...)
		if ex != nil {
			return result, ex
		}
		for _, v := range ps {
			result = append(result, iterater(v))
		}
		return result, nil
	}
}

func IDMapFunc[T dependency.IEntity](repo dependency.IRepository[T]) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (map[any]*T, exception.Exception) {
	f := EntitiesFunc(repo)
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (map[any]*T, exception.Exception) {
		result := map[any]*T{}
		ps, ex := f(ctx, dbKeys, tbKeys, conds...)
		if ex != nil {
			return result, ex
		}
		for _, v := range ps {
			result[v.ID()] = &v
		}
		return result, nil
	}
}

func EntitiesFunc[T dependency.IEntity](repo dependency.IRepository[T]) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) ([]T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) ([]T, exception.Exception) {
		opts := []dependency.BaseOptionFunc{}
		if len(dbKeys) > 0 {
			opts = append(opts, dependency.WithDbShardingKey(dbKeys...))
		}
		if len(tbKeys) > 0 {
			opts = append(opts, dependency.WithTbShardingKey(tbKeys...))
		}
		if len(conds) > 0 {
			opts = append(opts, dependency.WithConds(conds...))
		}
		opts = append(opts,
			dependency.WithReadOnly(true),
		)
		ps, err := repo.BaseQuery(ctx, opts...)
		if err != nil {
			return nil, exception.ERR_BUSI.Wrap(err)
		}
		return ps, nil
	}
}

func DetailByIdFunc[T dependency.IEntity](repo dependency.IRepository[T]) func(ctx context.Context, dbKeys, tbKeys []any, id any) (*T, exception.Exception) {
	f := DetailFunc(repo)
	return func(ctx context.Context, dbKeys, tbKeys []any, id any) (*T, exception.Exception) {
		return f(ctx, dbKeys, tbKeys, "`id` = ?", id)
	}
}

func DetailFunc[T dependency.IEntity](repo dependency.IRepository[T]) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (*T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (*T, exception.Exception) {
		opts := []dependency.BaseOptionFunc{}
		if len(dbKeys) > 0 {
			opts = append(opts, dependency.WithDbShardingKey(dbKeys...))
		}
		if len(tbKeys) > 0 {
			opts = append(opts, dependency.WithTbShardingKey(tbKeys...))
		}
		opts = append(opts, dependency.WithConds(conds...))
		opts = append(opts,
			dependency.WithReadOnly(true),
		)
		p, err := repo.BaseGet(ctx, opts...)
		if err != nil {
			return p, exception.ERR_BUSI_NOFOUND.Wrap(err)
		}
		return p, nil
	}
}
