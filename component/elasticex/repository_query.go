package elasticex

import (
	"context"
	"encoding/json"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
)

// BaseGet
func (r *BaseRepository[T]) BaseGet(ctx context.Context, opts ...dependency.BaseOptionFunc) (*T, error) {
	t := new(T)
	res, err := r.GetGetService(ctx, opts...).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, nil
	}
	bs, err := res.Source.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// BaseCount
func (r *BaseRepository[T]) BaseCount(ctx context.Context, opts ...dependency.BaseOptionFunc) (int64, error) {
	return r.GetCountService(ctx, opts...).Do(ctx)
}

// BaseQuery
func (r *BaseRepository[T]) BaseQuery(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, error) {
	res, _, err := r.BaseQueryWithCount(ctx, opts...)
	return res, err
}

// BaseQuery
func (r *BaseRepository[T]) BaseQueryWithCount(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, int64, error) {
	var (
		t  T
		ts = []T{}
	)
	// normal query or search after
	opt := dependency.NewBaseOption(opts...)
	if opt.Page == nil {
		return r.BaseSearch(ctx, opts...)
	}
	cursor := (opt.Page.GetPageIndex() - 1) * opt.Page.GetPageSize()
	size := opt.Page.GetSize()
	// form + size window search in max size
	if cursor+size <= MAX_WINDOW_SIZE {
		return r.BaseSearch(ctx, opts...)
	}
	// query conditions
	query := WithQuery(opt.Conds...)
	// query sorts
	sorts := WithSort(opt.Page.GetSorts()...)
	// search after
	ls := &LimitSearch{
		Limit:  int64(MAX_WINDOW_SIZE),
		Offset: int64(cursor),
		Total:  int64(size),
		Spans:  []*BatchSpan{},
	}
	sortValues := []any{}
	_, total, err := ls.SearchByStep(func(index, subCursor, subSize int64) ([]any, int64, error) {
		res := []any{}
		db := CoreFrmCtx(ctx, opt.GetDataBase(t))
		srv := db.Search().Index(opt.GetTableName(t))
		result, err := srv.SearchAfter(sortValues...).
			Query(query).
			SortBy(sorts...).
			Size(int(subSize)).
			Do(ctx)
		for _, row := range result.Hits.Hits {
			res = append(res, row)
			sortValues = row.Sort
			tPtr := ParseHit[T](row)
			ts = append(ts, *tPtr)
			opt.Iterate(row)
		}
		return res, result.TotalHits(), err
	})
	return ts, total, err
}

func (r *BaseRepository[T]) BaseSearch(ctx context.Context, opts ...dependency.BaseOptionFunc) ([]T, int64, error) {
	result := []T{}
	opt := dependency.NewBaseOption(opts...)
	srv := r.GetSearchServiceFrmOpt(ctx, opt)
	res, err := srv.Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	for _, hit := range res.Hits.Hits {
		var t T
		if err := json.Unmarshal(hit.Source, &t); err == nil {
			result = append(result, t)
		}
		opt.Iterate(hit)
	}
	return result, res.TotalHits(), nil
}

// GetCountService
func (r *BaseRepository[T]) GetCountService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.CountService {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.Count(opt.GetTableName(t))
	if len(opt.Conds) > 0 {
		qs := []elastic.Query{}
		for _, cond := range opt.Conds {
			if v, ok := cond.(elastic.Query); ok {
				qs = append(qs, v)
			}
		}
		srv.Query(elastic.NewBoolQuery().Must(qs...))
	}
	return srv
}

// GetService
func (r *BaseRepository[T]) GetGetService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.GetService {
	var (
		t T
	)
	opt := dependency.NewBaseOption(opts...)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.Get().Index(opt.GetTableName(t))
	if len(opt.Conds) > 0 {
		srv = srv.Id(cast.ToString(opt.Conds[0]))
	}
	return srv
}

// GetSearchService
func (r *BaseRepository[T]) GetSearchService(ctx context.Context, opts ...dependency.BaseOptionFunc) *elastic.SearchService {
	opt := dependency.NewBaseOption(opts...)
	return r.GetSearchServiceFrmOpt(ctx, opt)
}

// GetSearchService
func (r *BaseRepository[T]) GetSearchServiceFrmOpt(ctx context.Context, opt *dependency.BaseOption) *elastic.SearchService {
	var (
		t T
	)
	db := CoreFrmCtx(ctx, opt.GetDataBase(t))
	srv := db.Search().Index(opt.GetTableName(t))

	if len(opt.Conds) > 0 {
		// query conditions
		query := WithQuery(opt.Conds...)
		srv.Query(query)
	}
	if opt.Page != nil {
		// query sorts
		sorts := WithSort(opt.Page.GetSorts()...)
		srv = srv.SortBy(sorts...)
		srv = srv.From(int((opt.Page.GetPageIndex() - 1) * opt.Page.GetPageSize())).Size(int(opt.Page.GetPageSize()))
	} else if opt.BatchSize > 0 {
		if opt.BatchSize > MAX_WINDOW_SIZE {
			opt.BatchSize = MAX_WINDOW_SIZE
		}
		srv = srv.From(0).Size(int(opt.BatchSize))
	}
	if sa := opt.SearchAfter; sa != nil {
		for _, f := range opt.Page.GetSorts() {
			srv = srv.Sort(f.GetField(), !f.GetIsDesc())
		}
		srv = srv.SearchAfter(sa.GetSortValues()...)
	}
	return srv
}
