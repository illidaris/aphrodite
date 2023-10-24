package gormex

import (
	"context"
	"fmt"

	"github.com/IvanWhisper/aphrodite/component/dependency"

	"gorm.io/gorm"
)

var _ = dependency.IUnitOfWork(&GormTransactionImpl{})

func GetDbTX(id string) ContextKey {
	return ContextKey(fmt.Sprintf("%s_%s", DbTXPrefix, id))
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
	tx := t.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
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
			return err
		}
	}
	return tx.Commit().Error
}
