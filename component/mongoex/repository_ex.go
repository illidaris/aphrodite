package mongoex

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"go.mongodb.org/mongo-driver/bson"
)

func (r BaseRepository[T]) EntitiesByIDs(ctx context.Context, ids ...any) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}
	opts := []dependency.BaseOptionFunc{
		dependency.WithConds(bson.D{bson.E{Key: "_id", Value: bson.M{"$in": ids}}}),
	}
	if bizId := contextex.GetBizId(ctx); bizId != 0 {
		opts = append(opts, dependency.WithDbShardingKey(bizId))
	}
	return r.BaseQuery(ctx, opts...)
}

func (r BaseRepository[T]) EntityMapByIDs(ctx context.Context, ids ...any) (map[any]T, error) {
	res := map[any]T{}
	es, err := r.EntitiesByIDs(ctx, ids...)
	if err != nil {
		return res, nil
	}
	for _, v := range es {
		res[v.ID()] = v
	}
	return res, nil
}
