package elasticex

import (
	"context"
	"encoding/json"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/olivere/elastic/v7"
)

const (
	MAX_WINDOW_SIZE = 10000
)

// ParseHit
func ParseHit[T dependency.IEntity](h *elastic.SearchHit) *T {
	t := new(T)
	bs, err := h.Source.MarshalJSON()
	if err != nil {
		return t
	}
	err = json.Unmarshal(bs, t)
	if err != nil {
		return t
	}
	return t
}

// CoreFrmCtx
func CoreFrmCtx(ctx context.Context, id string) *elastic.Client {
	es, _ := elastic.NewClient()
	return es
}

func WithQuery(conds ...any) elastic.Query {
	qs := []elastic.Query{}
	if len(conds) > 0 {
		for _, cond := range conds {
			if v, ok := cond.(elastic.Query); ok {
				qs = append(qs, v)
			}
		}
	}
	return elastic.NewBoolQuery().Filter(qs...)
}

func WithSort(sorts ...dependency.ISortField) []elastic.Sorter {
	esSorts := []elastic.Sorter{}
	for _, sort := range sorts {
		if sort != nil {
			s := elastic.NewFieldSort(sort.GetField())
			if sort.GetIsDesc() {
				s.Desc()
			} else {
				s.Asc()
			}
			esSorts = append(esSorts, s)
		}
	}
	// 排序方式默认
	if len(esSorts) == 0 {
		esSorts = append(esSorts, elastic.NewFieldSort("_id").Asc())
	}
	return esSorts
}
