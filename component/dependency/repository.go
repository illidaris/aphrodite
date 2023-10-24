package dependency

import (
	"context"
)

// action
type DbAction func(ctx context.Context) error

// IUnitOfWork  trans
type IUnitOfWork interface {
	Execute(ctx context.Context, fs ...DbAction) (e error)
}

// BaseOption base repo exec
type BaseOption struct {
	Ignore    bool     `json:"ignore"`    // ignore if exist
	Lock      bool     `json:"lock"`      // lock row
	ReadOnly  bool     `json:"readOnly"`  // read only
	Selects   []string `json:"selects"`   // select fields
	Omits     []string `json:"omits"`     // omit fields select omit
	Conds     []any    `json:"conds"`     // conds where
	Page      IPage    `json:"page"`      // page
	BatchSize int64    `json:"batchSize"` // exec by batch
}

// IRepository repo
type IRepository[T IEntity] interface {
	BaseCreate(ctx context.Context, opt BaseOption, p ...T) (int64, error)
	BaseSave(ctx context.Context, opt BaseOption, p ...T) (int64, error)
	BaseUpdate(ctx context.Context, opt BaseOption, p T) (int64, error)
	BaseGet(ctx context.Context, opt BaseOption, p T) (int64, error)
	BaseDelete(ctx context.Context, opt BaseOption, p T) (int64, error)
	BaseCount(ctx context.Context, opt BaseOption, p T) (int64, error)
	BaseQuery(ctx context.Context, opt BaseOption, p T) ([]T, error)
}
