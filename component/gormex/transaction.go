package gormex

import (
	"context"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"gorm.io/gorm"
)

var _ = dependency.IUnitOfWork(&GormTransactionImpl{})

func GetDbTX(id string) contextex.ContextKey {
	return contextex.DbTxID.ID(id)
}
func NewUnitOfWork(id string) dependency.IUnitOfWork {
	return &GormTransactionImpl{
		id: id,
		db: GetTransactionDb(id),
	}
}

func GetTransactionDb(id string) *gorm.DB {
	return MySqlComponent.GetWriter(id)
}

type GormTransactionImpl struct {
	id string
	db *gorm.DB
}

// Execute An execution function is passed in, and transactions are executed within the function.
func (t *GormTransactionImpl) Execute(ctx context.Context, fs ...dependency.DbAction) (e error) {
	tx := t.db.WithContext(ctx).Begin()
	tx.Logger.Info(ctx, "transaction begin")
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			tx.Logger.Warn(ctx, fmt.Sprintf("transaction panic rollback %s", r))
			if err, ok := r.(error); ok {
				e = err
			} else {
				e = fmt.Errorf("unkonw %s", r)
			}
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	ctx = NewContext(ctx, t.id, tx)
	for _, f := range fs {
		if err := f(ctx); err != nil {
			tx.Rollback()
			tx.Logger.Warn(ctx, "transaction rollback", err.Error())
			return err
		}
	}
	e = tx.Commit().Error
	tx.Logger.Info(ctx, "transaction commit")
	return
}
