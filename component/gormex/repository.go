package gormex

import (
	"context"
	"errors"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/group"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ = dependency.IRepository[dependency.IEntity](&BaseRepository[dependency.IEntity]{}) // impl check

type BaseRepository[T dependency.IEntity] struct{} // base repository

// BaseCreate
func (r *BaseRepository[T]) BaseCreate(ctx context.Context, ps []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	if len(ps) == 0 {
		return 0, nil
	}
	opt := dependency.NewBaseOption(opts...)
	if idgen, ok := any(ps[0]).(dependency.IGenerateID); ok && opt.IDGenerate != nil {
		idgen.SetID(opt.IDGenerate(ctx))
	}
	return BaseGroup(func(v ...*T) (int64, error) {
		var t *T
		if len(v) > 0 {
			t = v[0]
		}
		db := r.BuildFrmOption(ctx, t, opt)
		result := db.Create(v)
		return result.RowsAffected, result.Error
	}, opt, ps...)
}

// BaseSave
func (r *BaseRepository[T]) BaseSave(ctx context.Context, ps []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	if len(ps) == 0 {
		return 0, nil
	}
	opt := dependency.NewBaseOption(opts...)
	if idgen, ok := any(ps[0]).(dependency.IGenerateID); ok && opt.IDGenerate != nil {
		idgen.SetID(opt.IDGenerate(ctx))
	}
	return BaseGroup(func(v ...*T) (int64, error) {
		var t *T
		if len(v) > 0 {
			t = v[0]
		}
		db := r.BuildFrmOption(ctx, t, opt)
		result := db.Save(v)
		return result.RowsAffected, result.Error
	}, opt, ps...)
}

// BaseUpdate
func (r *BaseRepository[T]) BaseUpdate(ctx context.Context, p *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	result := r.BuildFrmOptions(ctx, p, opts...).Updates(p)
	return result.RowsAffected, result.Error
}

// BaseGet
func (r *BaseRepository[T]) BaseGet(ctx context.Context, opts ...dependency.BaseOptionFunc) (*T, error) {
	var t T
	db := r.BuildFrmOptions(ctx, nil, opts...)
	res := db.First(&t)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &t, res.Error
}

// BaseDelete
func (r *BaseRepository[T]) BaseDelete(ctx context.Context, p *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	result := r.BuildFrmOptions(ctx, p, opts...).Delete(p)
	return result.RowsAffected, result.Error
}

// BaseCount
func (r *BaseRepository[T]) BaseCount(ctx context.Context, opts ...dependency.BaseOptionFunc) (int64, error) {
	var count int64
	opt := dependency.NewBaseOption(opts...)
	db := r.BuildConds(ctx, nil, opt)
	res := db.Count(&count)
	return count, res.Error
}

// BaseQuery
func (r *BaseRepository[T]) BaseQuery(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, error) {
	result := []T{}
	db := r.BuildFrmOptions(ctx, nil, opts...)
	res := db.Find(&result)
	return result, res.Error
}

// BaseQueryWithCount
func (r *BaseRepository[T]) BaseQueryWithCount(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, int64, error) {
	count, err := r.BaseCount(ctx, opts...)
	if err != nil {
		return nil, count, err
	}
	ts, err := r.BaseQuery(ctx, opts...)
	if err != nil {
		return ts, count, err
	}
	return ts, count, err
}

// BuildConds
func (r *BaseRepository[T]) BuildConds(ctx context.Context, t *T, opt *dependency.BaseOption) *gorm.DB {
	var (
		db *gorm.DB
	)
	if t == nil {
		t = new(T)
	}
	if sharding, ok := any(t).(dependency.IDbSharding); ok {
		opt.DataBase = sharding.DbSharding(opt.DbShardingKey...)
	}
	if opt.DataBase == "" {
		opt.DataBase = any(t).(dependency.IPo).Database()
	}
	if opt != nil && opt.ReadOnly {
		db = ReadOnly(ctx, opt.DataBase)
	} else {
		db = CoreFrmCtx(ctx, opt.DataBase)
	}
	db = db.Model(&t)
	if sharding, ok := any(t).(dependency.ITableSharding); ok {
		opt.TableName = sharding.TableSharding(opt.TbShardingKey...)
	}
	if opt.TableName != "" {
		db = db.Table(opt.TableName)
	}
	if opt != nil && len(opt.Conds) > 0 {
		db = db.Where(opt.Conds[0], opt.Conds[1:]...)
	}
	return db
}

// BuildFrmOption
func (r *BaseRepository[T]) BuildFrmOption(ctx context.Context, t *T, opt *dependency.BaseOption) *gorm.DB {
	db := r.BuildConds(ctx, t, opt)
	if opt.Ignore {
		db = db.Clauses(clause.Insert{Modifier: "IGNORE"})
	}
	if opt.Lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if len(opt.Selects) > 0 {
		db = db.Select(opt.Selects)
	}
	if len(opt.Omits) > 0 {
		db = db.Omit(opt.Omits...)
	}
	db = Option2Page(db, opt)
	return db
}

// BuildFrmOptions
func (r *BaseRepository[T]) BuildFrmOptions(ctx context.Context, t *T, opts ...dependency.BaseOptionFunc) *gorm.DB {
	opt := dependency.NewBaseOption(opts...)
	db := r.BuildFrmOption(ctx, t, opt)
	return db
}

// Option2Page
func Option2Page(db *gorm.DB, opt *dependency.BaseOption) *gorm.DB {
	if opt.Page != nil {
		for _, f := range opt.Page.GetSorts() {
			key := f.GetField()
			if f.GetIsDesc() {
				key = fmt.Sprintf("%s %s", key, "desc")
			}
			db = db.Order(key)
		}
		db = db.Offset(int((opt.Page.GetPageIndex() - 1) * opt.Page.GetPageSize())).Limit(int(opt.Page.GetPageSize()))
	} else {
		if opt.ReadOnly && opt.BatchSize > 0 {
			db = db.Limit(int(opt.BatchSize))
		}
	}
	return db
}

// CoreFrmCtx
func CoreFrmCtx(ctx context.Context, id string) *gorm.DB {
	return WithContext(ctx, id)
}

// ReadOnly
func ReadOnly(ctx context.Context, id string) *gorm.DB {
	return MySqlComponent.GetReader(id).Session(&gorm.Session{
		QueryFields: !disableQueryFields,
		Context:     ctx,
	})
}

// BaseGroup
func BaseGroup[T dependency.IEntity](f func(v ...*T) (int64, error), opt *dependency.BaseOption, p ...*T) (int64, error) {
	if opt.BatchSize >= int64(len(p)) {
		return f(p...)
	}
	affect, errM := group.GroupFunc[T](f, int(opt.BatchSize), p...)
	if l := len(errM); l > 0 {
		errMsg := fmt.Sprintf("%d err ", l)
		for _, v := range errM {
			errMsg += v.Error()
		}
		return affect, errors.New(errMsg)
	}
	return affect, nil
}
