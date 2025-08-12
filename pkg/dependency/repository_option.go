package dependency

import "context"

const (
	BATCH_SIZE = 1000 // default batch size
)

// BaseOptionFunc base option func
type BaseOptionFunc func(o *BaseOption)

// NewBaseOption with func
func NewBaseOption(opts ...BaseOptionFunc) *BaseOption {
	def := &BaseOption{
		BatchSize: BATCH_SIZE,
	}
	for _, opt := range opts {
		opt(def)
	}
	return def
}

// BaseOption base repo exec
type BaseOption struct {
	Ignore         bool                          `json:"ignore"`        // ignore if exist
	Lock           bool                          `json:"lock"`          // lock row
	ReadOnly       bool                          `json:"readOnly"`      // read only
	Selects        []string                      `json:"selects"`       // select fields
	Omits          []string                      `json:"omits"`         // omit fields select omit
	Conds          []any                         `json:"conds"`         // conds where
	Page           IPage                         `json:"page"`          // page
	DeepPage       IDeepPage                     `json:"deepPage"`      // deep page
	SearchAfter    ISearchAfter                  `json:"searchAfter"`   // search after
	BatchSize      int64                         `json:"batchSize"`     // exec by batch
	TableName      string                        `json:"tableName"`     // table name
	DataBase       string                        `json:"dataBase"`      // db name
	DbShardingKey  []any                         `json:"dbShardingKey"` // db sharding key
	TbShardingKey  []any                         `json:"tbShardingKey"` // table sharding key
	UpdatedMap     map[string]any                `json:"-"`             // updated map
	IDGenerate     func(ctx context.Context) any `json:"-"`             // id generate func
	IterativeFuncs []func(any)                   `json:"-"`             // iterative func
}

// GetDataBase
func (opt BaseOption) GetDataBase(t IEntity) string {
	if len(opt.DataBase) > 0 {
		return opt.DataBase
	}
	opt.DataBase = t.Database()
	if sharding, ok := any(t).(IDbSharding); ok {
		opt.DataBase = sharding.DbSharding(opt.DbShardingKey...)
	}
	return opt.DataBase
}

// GetTableName
func (opt BaseOption) GetTableName(t IEntity) string {
	if len(opt.TableName) > 0 {
		return opt.TableName
	}
	opt.TableName = t.TableName()
	if sharding, ok := any(t).(ITableSharding); ok {
		opt.TableName = sharding.TableSharding(opt.TbShardingKey...)
	}
	return opt.TableName
}

// Iterate
func (opt BaseOption) Iterate(v any) {
	for _, iterate := range opt.IterativeFuncs {
		iterate(v)
	}
}

// WithIgnore
func WithIgnore(v bool) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Ignore = v
	}
}

// WithLock
func WithLock(v bool) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Lock = v
	}
}

// WithReadOnly
func WithReadOnly(v bool) BaseOptionFunc {
	return func(o *BaseOption) {
		o.ReadOnly = v
	}
}

// WithSelects
func WithSelects(vs ...string) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Selects = vs
	}
}

// WithOmits
func WithOmits(vs ...string) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Omits = vs
	}
}

// WithConds
func WithConds(vs ...any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Conds = vs
	}
}

// WithPage
func WithPage(v IPage) BaseOptionFunc {
	return func(o *BaseOption) {
		o.Page = v
	}
}

// WithDeepPage
func WithDeepPage(v IDeepPage) BaseOptionFunc {
	return func(o *BaseOption) {
		o.DeepPage = v
	}
}

// WithSearchAfter
func WithSearchAfter(v ISearchAfter) BaseOptionFunc {
	return func(o *BaseOption) {
		o.SearchAfter = v
	}
}

// WithBatchSize
func WithBatchSize(v int64) BaseOptionFunc {
	return func(o *BaseOption) {
		o.BatchSize = v
	}
}

// WithTableName
func WithTableName(v string) BaseOptionFunc {
	return func(o *BaseOption) {
		o.TableName = v
	}
}

// WithDataBase
func WithDataBase(v string) BaseOptionFunc {
	return func(o *BaseOption) {
		o.DataBase = v
	}
}

// WithDbShardingKey
func WithDbShardingKey(v ...any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.DbShardingKey = v
	}
}

// WithTbShardingKey
func WithTbShardingKey(v ...any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.TbShardingKey = v
	}
}

// WithUpdatedMap
func WithUpdatedMap(v map[string]any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.UpdatedMap = v
	}
}

// WithIDGenerate
func WithIDGenerate(v func(context.Context) any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.IDGenerate = v
	}
}

// WithIterativeFunc
func WithIterativeFunc(v ...func(any)) BaseOptionFunc {
	return func(o *BaseOption) {
		if o.IterativeFuncs == nil {
			o.IterativeFuncs = []func(any){}
		}
		o.IterativeFuncs = append(o.IterativeFuncs, v...)
	}
}
