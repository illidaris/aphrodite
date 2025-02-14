package mongoex

import (
	"context"
	"errors"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/convert"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/group"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		count := int64(0)
		if len(v) > 0 {
			t = v[0]
		}
		args := []interface{}{}
		for _, item := range v {
			args = append(args, item)
		}
		opt := dependency.NewBaseOption(opts...)
		finalErr := r.BuildFrmOption(ctx, t, opt, func(colls *mongo.Collection) error {
			result, err := colls.InsertMany(ctx, args)
			if result != nil {
				count = int64(len(result.InsertedIDs))
			}
			return err
		})
		return count, finalErr
	}, opt, ps...)
}

// BaseSave
func (r *BaseRepository[T]) BaseSave(ctx context.Context, ps []*T, opts ...dependency.BaseOptionFunc) (int64, error) {
	panic("no impl")
}

// BaseUpdate
func (r *BaseRepository[T]) BaseUpdate(ctx context.Context, p *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	count := int64(0)
	opt := dependency.NewBaseOption(opts...)
	finalErr := r.BuildFrmOption(ctx, nil, opt, func(colls *mongo.Collection) error {
		updated := bson.E{Key: "$set", Value: p}
		if opt.UpdatedMap != nil {
			updated.Value = opt.UpdatedMap
		}
		res, err := colls.UpdateOne(ctx, QueryConds(opt), bson.D{updated})
		if res != nil {
			count = res.ModifiedCount
		}
		return err
	})

	return count, finalErr
}

// BaseGet
func (r *BaseRepository[T]) BaseGet(ctx context.Context, opts ...dependency.BaseOptionFunc) (*T, error) {
	var t T
	opt := dependency.NewBaseOption(opts...)
	finalErr := r.BuildFrmOption(ctx, nil, opt, func(colls *mongo.Collection) error {
		filter := QueryConds(opt)
		res := colls.FindOne(ctx, filter)
		err := res.Decode(&t)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil
			}
			return err
		}
		return nil
	})
	return &t, finalErr
}

// BaseDelete
func (r *BaseRepository[T]) BaseDelete(ctx context.Context, p *T, opts ...dependency.BaseOptionFunc) (int64, error) {
	count := int64(0)
	opt := dependency.NewBaseOption(opts...)
	finalErr := r.BuildFrmOption(ctx, nil, opt, func(colls *mongo.Collection) error {
		res, err := colls.DeleteMany(ctx, QueryConds(opt))
		if res != nil {
			count = res.DeletedCount
		}
		return err
	})
	return count, finalErr
}

// BaseCount
func (r *BaseRepository[T]) BaseCount(ctx context.Context, opts ...dependency.BaseOptionFunc) (int64, error) {
	count := int64(0)
	opt := dependency.NewBaseOption(opts...)
	err := r.BuildFrmOption(ctx, nil, opt, func(colls *mongo.Collection) error {
		total, docError := colls.CountDocuments(ctx, QueryConds(opt))
		count = total
		return docError
	})
	return count, err
}

// BaseQuery
func (r *BaseRepository[T]) BaseQuery(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, error) {
	result := []T{}
	opt := dependency.NewBaseOption(opts...)
	err := r.BuildFrmOption(ctx, nil, opt, func(colls *mongo.Collection) error {
		cur, findErr := colls.Find(ctx, QueryConds(opt), Option2Page(opt))
		if findErr != nil {
			return findErr
		}
		if cur == nil {
			return ErrCurNil
		}
		findErr = cur.All(ctx, &result)
		return findErr
	})
	return result, err
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
func (r *BaseRepository[T]) BuildConds(ctx context.Context, t *T, opt *dependency.BaseOption, dbcallback func(*mongo.Database) error) error {
	if t == nil {
		t = new(T)
	}
	if sharding, ok := any(t).(dependency.IDbSharding); ok {
		opt.DataBase = sharding.DbSharding(opt.DbShardingKey...)
	}
	if opt.DataBase == "" {
		opt.DataBase = any(t).(dependency.IPo).Database()
	}
	if opt.DataBase == "" {
		opt.DataBase = getKey(ctx)
	}
	// 转化真实的数据库
	realDb := GetMongoNameByKey(opt.DataBase)
	session := mongo.SessionFromContext(ctx)
	if session == nil {
		c := MongoComponent.GetWriter(opt.DataBase)
		return c.UseSessionWithOptions(ctx, options.Session(), func(sc mongo.SessionContext) error {
			return dbcallback(sc.Client().Database(realDb))
		})
	}
	return dbcallback(session.Client().Database(realDb))
}

// BuildFrmOption
func (r *BaseRepository[T]) BuildFrmOption(ctx context.Context, t *T, opt *dependency.BaseOption, colcallback func(*mongo.Collection) error) error {
	return r.BuildConds(ctx, t, opt, func(db *mongo.Database) error {
		if t == nil {
			t = new(T)
		}
		var (
			colls *mongo.Collection
		)
		if len(opt.TableName) == 0 {
			opt.TableName = any(t).(dependency.IEntity).TableName()
		}
		if sharding, ok := any(t).(dependency.ITableSharding); ok {
			opt.TableName = sharding.TableSharding(opt.TbShardingKey...)
		}
		if opt.TableName != "" {
			colls = db.Collection(opt.TableName)
		}
		return colcallback(colls)
	})
}

// BuildFrmOptions
func (r *BaseRepository[T]) BuildFrmOptions(ctx context.Context, t *T, colcallback func(*mongo.Collection) error, opts ...dependency.BaseOptionFunc) error {
	opt := dependency.NewBaseOption(opts...)
	return r.BuildFrmOption(ctx, t, opt, colcallback)
}

// Option2Page
func Option2Page(opt *dependency.BaseOption) *options.FindOptions {
	mongoFindOpts := options.Find()
	if opt.Page != nil {
		sorts := bson.D{}
		for _, f := range opt.Page.GetSorts() {
			key, _ := convert.FieldFilter(f.GetField(), convert.FieldFilterLevelDefault)
			if key == "" {
				continue
			}
			asc := 1
			if f.GetIsDesc() {
				asc = -1
			}
			sorts = append(sorts, bson.E{Key: key, Value: asc})
		}
		mongoFindOpts.SetSort(sorts)
		mongoFindOpts.SetSkip((opt.Page.GetPageIndex() - 1) * opt.Page.GetPageSize())
		mongoFindOpts.SetLimit(opt.Page.GetPageSize())
		if d := QueryFields(opt); len(d) > 0 {
			mongoFindOpts.SetProjection(d)
		}
	} else {
		if opt.ReadOnly && opt.BatchSize > 0 {
			mongoFindOpts.SetLimit(opt.BatchSize)
		}
	}
	return mongoFindOpts
}

func QueryFields(opt *dependency.BaseOption) bson.D {
	d := bson.D{}
	for _, sel := range opt.Selects {
		d = append(d, bson.E{Key: sel, Value: 1})
	}
	for _, sel := range opt.Omits {
		d = append(d, bson.E{Key: sel, Value: 0})
	}
	return d
}

func QueryConds(opt *dependency.BaseOption) bson.D {
	l := len(opt.Conds)
	switch l {
	case 0:
		return bson.D{}
	case 1:
		d, ok := opt.Conds[0].(bson.D)
		if !ok {
			return bson.D{}
		}
		return d
	default:
		keys := []string{}
		values := []interface{}{}
		d := bson.D{}
		for k, v := range opt.Conds {
			if k%2 == 0 {
				keys = append(keys, v.(string))
			} else {
				values = append(values, v)
			}
		}
		for index, v := range keys {
			e := bson.E{Key: v}
			if index < len(values) {
				e.Value = values[index]
			}
			d = append(d, e)
		}
		return d
	}
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
