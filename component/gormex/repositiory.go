package gormex

import (
	"context"
	"errors"
	"fmt"

	"github.com/IvanWhisper/aphrodite/component/dependency"
	"github.com/IvanWhisper/aphrodite/pkg/group"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	BATCH_SIZE = 1000 // default batch size
)

var _ = dependency.IRepository[dependency.IEntity](&BaseRepository[dependency.IEntity]{}) // impl check

type BaseRepository[T dependency.IEntity] struct{} // base repository

// BaseCreate
func (r *BaseRepository[T]) BaseCreate(ctx context.Context, opt dependency.BaseOption, p ...T) (int64, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return BaseGroup(func(v ...T) (int64, error) {
		db := r.BuildFrmOption(ctx, p[0], opt)
		result := db.Create(v)
		return result.RowsAffected, result.Error
	}, opt, p...)
}

// BaseSave
func (r *BaseRepository[T]) BaseSave(ctx context.Context, opt dependency.BaseOption, p ...T) (int64, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return BaseGroup(func(v ...T) (int64, error) {
		db := r.BuildFrmOption(ctx, p[0], opt)
		result := db.Save(v)
		return result.RowsAffected, result.Error
	}, opt, p...)
}

// BaseUpdate
func (r *BaseRepository[T]) BaseUpdate(ctx context.Context, opt dependency.BaseOption, p T) (int64, error) {
	result := r.BuildFrmOption(ctx, p, opt).Updates(p)
	return result.RowsAffected, result.Error
}

// BaseGet
func (r *BaseRepository[T]) BaseGet(ctx context.Context, opt dependency.BaseOption, p T) (int64, error) {
	db := r.BuildFrmOption(ctx, p, opt)
	res := db.First(p)
	return res.RowsAffected, res.Error
}

// BaseDelete
func (r *BaseRepository[T]) BaseDelete(ctx context.Context, opt dependency.BaseOption, p T) (int64, error) {
	result := r.BuildFrmOption(ctx, p, opt).Delete(p)
	return result.RowsAffected, result.Error
}

// BaseCount
func (r *BaseRepository[T]) BaseCount(ctx context.Context, opt dependency.BaseOption, p T) (int64, error) {
	var count int64
	db := r.BuildConds(ctx, p, true, opt.Conds...)
	res := db.Count(&count)
	return count, res.Error
}

// BaseQuery
func (r *BaseRepository[T]) BaseQuery(ctx context.Context, opt dependency.BaseOption, p T) ([]T, error) {
	result := []T{}
	db := r.BuildFrmOption(ctx, p, opt)
	res := db.Find(&result)
	return result, res.Error
}

func (r *BaseRepository[T]) BuildConds(ctx context.Context, p T, readOnly bool, conds ...any) *gorm.DB {
	var db *gorm.DB
	if readOnly {
		db = ReadOnly(ctx, p.Database())
	} else {
		db = CoreFrmCtx(ctx, p.Database())
	}
	db = db.Model(p)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}
	return db
}

func (r *BaseRepository[T]) BuildFrmOption(ctx context.Context, p T, opt dependency.BaseOption) *gorm.DB {
	db := r.BuildConds(ctx, p, opt.ReadOnly, opt.Conds...)
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
		if opt.BatchSize == 0 {
			opt.BatchSize = BATCH_SIZE
		}
		db = db.Limit(int(opt.BatchSize))
	}
	return db
}

func CoreFrmCtx(ctx context.Context, id string) *gorm.DB {
	return WithContext(ctx, id)
}
func ReadOnly(ctx context.Context, id string) *gorm.DB {
	return MySqlComponent.GetReader(id).Session(&gorm.Session{
		QueryFields: !disableQueryFields,
		Context:     ctx,
	})
}
func BaseGroup[T dependency.IEntity](f func(v ...T) (int64, error), opt dependency.BaseOption, p ...T) (int64, error) {
	if opt.BatchSize == 0 {
		opt.BatchSize = BATCH_SIZE
	}
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
