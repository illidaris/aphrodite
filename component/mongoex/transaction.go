package mongoex

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ = dependency.IUnitOfWork(&MongoTransactionImpl{})

type UOWOptions struct {
	IsolationLevel sql.IsolationLevel
}

type UOWOptionFunc func(*UOWOptions)

func WithIsolationLevel(level sql.IsolationLevel) UOWOptionFunc {
	return func(o *UOWOptions) {
		o.IsolationLevel = level
	}
}

func NewUnitOfWork(id string, opts ...UOWOptionFunc) dependency.IUnitOfWork {
	return NewMongoUnitOfWork(id, opts...)
}

func NewMongoUnitOfWork(id string, opts ...UOWOptionFunc) *MongoTransactionImpl {
	opt := &UOWOptions{}
	for _, o := range opts {
		o(opt)
	}
	impl := &MongoTransactionImpl{
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

func GetTransactionDb(id string) *mongo.Client {
	return MongoComponent.GetWriter(id)
}

type MongoTransactionImpl struct {
	id     string
	db     *mongo.Client
	txOpts *sql.TxOptions
}

func (t *MongoTransactionImpl) WithTxOpts(txOpts *sql.TxOptions) dependency.IUnitOfWork {
	t.txOpts = txOpts
	return t
}

// Execute An execution function is passed in, and transactions are executed within the function.
func (t *MongoTransactionImpl) Execute(ctx context.Context, fs ...dependency.DbAction) (e error) {
	return t.db.UseSessionWithOptions(ctx, options.Session(), func(sc mongo.SessionContext) error {
		_, err := sc.WithTransaction(ctx, func(tx mongo.SessionContext) (res interface{}, e error) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						e = err
					} else {
						e = fmt.Errorf("unkonw %v", r)
					}
				}
			}()
			for _, f := range fs {
				if err := f(tx); err != nil {
					e = err
					return
				}
			}
			return nil, nil
		})
		return err
	})
}
