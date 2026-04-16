package crud

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/exception"
)

func Create[T dependency.IEntity](repo dependency.IRepository[T], iterater func(*T) *T, opts ...Option) func(ctx context.Context, ts []*T) (int64, exception.Exception) {
	return func(ctx context.Context, ts []*T) (int64, exception.Exception) {
		if len(ts) == 0 {
			return 0, nil
		}

		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{},
		}

		first := ts[0]
		option.RepoOptions = append(option.RepoOptions, dependency.ShardingOptions(first)...)

		for _, opt := range opts {
			opt(option)
		}
		if iterater != nil {
			for _, t := range ts {
				t = iterater(t)
			}
		}
		affect, err := repo.BaseCreate(ctx, ts, option.RepoOptions...)
		if err != nil {
			return affect, exception.ERR_BUSI_CREATE.Wrap(err)
		}
		return affect, nil
	}
}

func DetailByIdFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, req any, id int64) (*T, exception.Exception) {
	return func(ctx context.Context, req any, id int64) (*T, exception.Exception) {
		return BaseDetailFunc(repo, opts...)(ctx, req, "`id` = ?", id)
	}
}

func DetailByCodeFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, req any, code string) (*T, exception.Exception) {
	return func(ctx context.Context, req any, code string) (*T, exception.Exception) {
		return BaseDetailFunc(repo, opts...)(ctx, req, "`code` = ?", code)
	}
}

func DetailFunc[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, req dependency.ICond) (*T, exception.Exception) {
	return func(ctx context.Context, req dependency.ICond) (*T, exception.Exception) {
		return BaseDetailFunc(repo, opts...)(ctx, req, req.GetConds()...)
	}
}

func Update[T dependency.IEntity](repo dependency.IRepository[T], iterater func(*T) *T, opts ...Option) func(ctx context.Context, t *T, conds ...any) (int64, exception.Exception) {
	return func(ctx context.Context, t *T, conds ...any) (int64, exception.Exception) {
		if t == nil {
			return 0, nil
		}
		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{},
		}
		option.RepoOptions = append(option.RepoOptions, dependency.ShardingOptions(t)...)
		option.RepoOptions = append(option.RepoOptions, dependency.WithConds(conds...))
		for _, opt := range opts {
			opt(option)
		}
		if iterater != nil {
			t = iterater(t)
		}
		affect, err := repo.BaseUpdate(ctx, t, option.RepoOptions...)
		if err != nil {
			return affect, exception.ERR_BUSI_UPDATE.Wrap(err)
		}
		return affect, nil
	}
}

func Delete[T dependency.IEntity](repo dependency.IRepository[T], opts ...Option) func(ctx context.Context, req any, conds ...any) (int64, exception.Exception) {
	return func(ctx context.Context, req any, conds ...any) (int64, exception.Exception) {
		option := &Options{
			RepoOptions: []dependency.BaseOptionFunc{},
		}
		option.RepoOptions = append(option.RepoOptions, dependency.ShardingOptions(req)...)
		option.RepoOptions = append(option.RepoOptions, dependency.WithConds(conds...))
		for _, opt := range opts {
			opt(option)
		}
		affect, err := repo.BaseDelete(ctx, new(T), option.RepoOptions...)
		if err != nil {
			return affect, exception.ERR_BUSI_DELETE.Wrap(err)
		}
		return affect, nil
	}
}
