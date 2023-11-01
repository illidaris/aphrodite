package dependency

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
	Ignore      bool     `json:"ignore"`      // ignore if exist
	Lock        bool     `json:"lock"`        // lock row
	ReadOnly    bool     `json:"readOnly"`    // read only
	Selects     []string `json:"selects"`     // select fields
	Omits       []string `json:"omits"`       // omit fields select omit
	Conds       []any    `json:"conds"`       // conds where
	Page        IPage    `json:"page"`        // page
	BatchSize   int64    `json:"batchSize"`   // exec by batch
	TableName   string   `json:"tableName"`   // table name
	DataBase    string   `json:"dataBase"`    // db name
	ShardingKey []any    `json:"shardingKey"` // sharding key
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

// WithShardingKey
func WithShardingKey(v ...any) BaseOptionFunc {
	return func(o *BaseOption) {
		o.ShardingKey = v
	}
}
