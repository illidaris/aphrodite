package gormex

import (
	"context"

	"gorm.io/gorm"
)

func NewContext(ctx context.Context, id string, newdb *gorm.DB) context.Context {
	if newdb != nil {
		return context.WithValue(ctx, GetDbTX(id), newdb)
	}
	return context.WithValue(ctx, GetDbTX(id), MySqlComponent.GetWriter(id))
}

func WithContext(ctx context.Context, id string) *gorm.DB {
	v := ctx.Value(GetDbTX(id))
	if d, ok := v.(*gorm.DB); ok {
		return d
	}
	db := MySqlComponent.GetWriter(id)
	if db == nil {
		return nil
	}
	return db.Session(&gorm.Session{
		QueryFields: !disableQueryFields,
		Context:     ctx,
	})
}
