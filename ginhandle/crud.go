package ginhandle

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/biz/crud"
	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/jinzhu/copier"
)

func ListHandler[Req dependency.ICondPage, T dependency.IEntity]() func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (*dto.RecordPtrPager[T], exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		return crud.PagingListFunc[T, T](repo, nil)(ctx, *r)
	})
}

func CreateManyHandler[Req any, T dependency.IEntity](f func(Req) []*T) func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (int64, exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		return crud.Create(repo, nil)(ctx, f(*r))
	})
}

func CreateHandler[Req any, T dependency.IEntity]() func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (int64, exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		t := new(T)
		_ = copier.Copy(t, r)
		return crud.Create(repo, nil)(ctx, []*T{t})
	})
}

func UpdateHandler[Req dependency.ICond, T dependency.IEntity]() func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (int64, exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		t := new(T)
		_ = copier.Copy(t, r)
		return crud.Update(repo, nil)(ctx, t, (*r).GetConds()...)
	})
}

func DeleteHandler[Req dependency.ICond, T dependency.IEntity]() func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (int64, exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		t := new(T)
		_ = copier.Copy(t, r)
		return crud.Delete(repo, nil)(ctx, r, (*r).GetConds()...)
	})
}

func DetailHandler[Req dependency.ICond, T dependency.IEntity]() func(c *gin.Context) {
	return GinOneHandler(func(ctx context.Context, r *Req) (*T, exception.Exception) {
		repo := &gormex.BaseRepository[T]{}
		return crud.DetailFunc(repo, nil)(ctx, *r)
	})
}
