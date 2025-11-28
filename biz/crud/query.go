package crud

import (
	"context"

	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"github.com/illidaris/aphrodite/pkg/exception"
)

type Option func(*Options)
type Options struct {
	RepoOptions []dependency.BaseOptionFunc
}

func WithRepoOptins(vs ...dependency.BaseOptionFunc) Option {
	return func(o *Options) {
		o.RepoOptions = append(o.RepoOptions, vs...)
	}
}

func PagingListFunc[T dependency.IEntity, K dto.IRow](repo dependency.IRepository[T], iterater func(T) *K, opts ...Option) func(ctx context.Context, req dependency.ICondPage) (*dto.RecordPtrPager[K], exception.Exception) {
	return func(ctx context.Context, req dependency.ICondPage) (*dto.RecordPtrPager[K], exception.Exception) {
		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{
				dependency.WithPage(req),
				dependency.WithReadOnly(true),
			},
		}

		if keys := req.GetDbShardingKeys(); len(keys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithDbShardingKey(keys...))
		}
		if keys := req.GetTbShardingKeys(); len(keys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithTbShardingKey(keys...))
		}
		if keys := req.GetConds(); len(keys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithConds(keys...))
		}

		for _, opt := range opts {
			opt(option)
		}

		result := new(dto.RecordPtrPager[K])
		result.PageIndex = req.GetPageIndex()
		result.PageSize = req.GetPageSize()
		result.Data = make([]*K, 0)

		ps, total, err := repo.BaseQueryWithCount(ctx, option.RepoOptions...)
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

func ListFunc[T dependency.IEntity, K dto.IRow](repo dependency.IRepository[T], iterater func(T) *K, opts ...Option) func(ctx context.Context, req dependency.ICond) ([]*K, exception.Exception) {
	return func(ctx context.Context, req dependency.ICond) ([]*K, exception.Exception) {
		result := []*K{}
		f := EntitiesFunc(repo, opts...)
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

func IDMapFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (map[any]*T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (map[any]*T, exception.Exception) {
		result := map[any]*T{}
		f := EntitiesFunc(repo, opts...)
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

func EntitiesFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) ([]T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) ([]T, exception.Exception) {
		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{
				dependency.WithReadOnly(true),
			},
		}

		if len(dbKeys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithDbShardingKey(dbKeys...))
		}
		if len(tbKeys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithTbShardingKey(tbKeys...))
		}
		if len(conds) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithConds(conds...))
		}

		for _, opt := range opts {
			opt(option)
		}

		ps, err := repo.BaseQuery(ctx, option.RepoOptions...)
		if err != nil {
			return nil, exception.ERR_BUSI.Wrap(err)
		}
		return ps, nil
	}
}

func DetailByIdFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, dbKeys, tbKeys []any, id any) (*T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, id any) (*T, exception.Exception) {
		f := DetailFunc(repo, opts...)
		return f(ctx, dbKeys, tbKeys, "`id` = ?", id)
	}
}

func DetailFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (*T, exception.Exception) {
	return func(ctx context.Context, dbKeys, tbKeys []any, conds ...any) (*T, exception.Exception) {
		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{
				dependency.WithReadOnly(true),
			},
		}

		if len(dbKeys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithDbShardingKey(dbKeys...))
		}
		if len(tbKeys) > 0 {
			option.RepoOptions = append(option.RepoOptions, dependency.WithTbShardingKey(tbKeys...))
		}
		option.RepoOptions = append(option.RepoOptions, dependency.WithConds(conds...))

		for _, opt := range opts {
			opt(option)
		}

		p, err := repo.BaseGet(ctx, option.RepoOptions...)
		if err != nil {
			return p, exception.ERR_BUSI_NOFOUND.Wrap(err)
		}
		return p, nil
	}
}
