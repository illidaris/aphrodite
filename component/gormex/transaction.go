package gormex

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"gorm.io/gorm"
)

var _ = dependency.IUnitOfWork(&GormTransactionImpl{})

type UOWOptions struct {
	IsolationLevel sql.IsolationLevel
}

type UOWOptionFunc func(*UOWOptions)

func WithIsolationLevel(level sql.IsolationLevel) UOWOptionFunc {
	return func(o *UOWOptions) {
		o.IsolationLevel = level
	}
}

func GetDbTX(id string) contextex.ContextKey {
	return contextex.DbTxID.ID(id)
}
func NewUnitOfWork(id string, opts ...UOWOptionFunc) dependency.IUnitOfWork {
	return NewGormUnitOfWork(id, opts...)
}

func NewGormUnitOfWork(id string, opts ...UOWOptionFunc) *GormTransactionImpl {
	opt := &UOWOptions{}
	for _, o := range opts {
		o(opt)
	}
	impl := &GormTransactionImpl{
		id: id,
		db: GetTransactionDb(id),
	}
	// 自定义事务隔离等级
	if opt.IsolationLevel != sql.LevelDefault {
		impl.txOpts = &sql.TxOptions{
			Isolation: opt.IsolationLevel,
		}
	}
	return impl
}

func GetTransactionDb(id string) *gorm.DB {
	return MySqlComponent.GetWriter(id)
}

type GormTransactionImpl struct {
	id     string
	db     *gorm.DB
	txOpts *sql.TxOptions
}

func (t *GormTransactionImpl) WithTxOpts(txOpts *sql.TxOptions) dependency.IUnitOfWork {
	t.txOpts = txOpts
	return t
}

// Execute An execution function is passed in, and transactions are executed within the function.
func (t *GormTransactionImpl) Execute(ctx context.Context, fs ...dependency.DbAction) (e error) {
	var tx *gorm.DB
	if t.txOpts != nil {
		tx = t.db.WithContext(ctx).Begin(t.txOpts)
	} else {
		tx = t.db.WithContext(ctx).Begin()
	}
	tx.Logger.Info(ctx, "transaction begin")
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			tx.Logger.Warn(ctx, fmt.Sprintf("transaction panic rollback %v", r))
			if err, ok := r.(error); ok {
				e = err
			} else {
				e = fmt.Errorf("unkonw %v", r)
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
			tx.Logger.Warn(ctx, fmt.Sprintf("transaction rollback %v", err))
			return err
		}
	}
	e = tx.Commit().Error
	tx.Logger.Info(ctx, "transaction commit")
	return
}
