# aphrodite

## component 组件

针对于各种基础组件的封装，对于每一种组件，定义一个`component`。
组件默认采用显式读写分离，负载采用随机负载模式，也可以根据需要自定义。

```go
type IComponent[T IItem] interface {
	NewWriter(id string, items ...T) // 新增写库
	GetWriter(id string) T // 获取写库
	NewReader(id string, items ...T) // 新增读库
	GetReader(id string) T // 获取读库
	SetWriterBalance(f func(ts ...IInstance[T]) IInstance[T]) // 设置负载均衡
	SetReaderBalance(f func(ts ...IInstance[T]) IInstance[T]) // 设置负载均衡
}
```

### gormex gorm 扩展组件

下图为内置的`mysql`组件，采用 gorm 框架，使用本项目中`gormex`包

```go
var MySqlComponent = embedded.NewComponent[*gorm.DB]()
```

主要是为了方便对于数据库操作，以`gorm`为基础，将来可以扩展到其他组件，并实现以下方法：

```go
// IRepository repo
type IRepository[T IEntity] interface {
	BaseCreate(ctx context.Context, ps []*T, opts ...BaseOptionFunc) (int64, error)
	BaseSave(ctx context.Context, ps []*T, opts ...BaseOptionFunc) (int64, error)
	BaseUpdate(ctx context.Context, p *T, opts ...BaseOptionFunc) (int64, error)
	BaseGet(ctx context.Context, opts ...BaseOptionFunc) (*T, error)
	BaseDelete(ctx context.Context, p *T, opts ...BaseOptionFunc) (int64, error)
	BaseCount(ctx context.Context, opts ...BaseOptionFunc) (int64, error)
	BaseQuery(ctx context.Context, opts ...BaseOptionFunc) ([]T, error)
}
```

// 内置数仓结构已经实现了上述方法，也可以根据需要重写

```go
type BaseRepository[T dependency.IEntity] struct{} // base repository
```

操作配置，使用配置可以根据需要调整数据操作，使用选项模式使用

```go
// BaseOption base repo exec
type BaseOption struct {
	Ignore      bool                          `json:"ignore"`      // ignore if exist
	Lock        bool                          `json:"lock"`        // lock row
	ReadOnly    bool                          `json:"readOnly"`    // read only
	Selects     []string                      `json:"selects"`     // select fields
	Omits       []string                      `json:"omits"`       // omit fields select omit
	Conds       []any                         `json:"conds"`       // conds where
	Page        IPage                         `json:"page"`        // page
	BatchSize   int64                         `json:"batchSize"`   // exec by batch
	TableName   string                        `json:"tableName"`   // table name
	DataBase    string                        `json:"dataBase"`    // db name
	ShardingKey []any                         `json:"shardingKey"` // sharding key
	IDGenerate  func(ctx context.Context) any `json:"-"`           // id generate func
}
```

使用配置

- `dependency.WithReadOnly(true)` 使用读节点，只能用于单纯查询
- `dependency.WithConds(2)` 查询条件一般为`[]interface{}{"code = ? AND age > ?","123",22}`,单纯使用主键时可以忽略字段名
- `dependency.WithBatchSize(3)` 默认批量为 1000, 限定查询数量，防止全量查询的风险，在使用分页时不生效

```go
ps, err := repo.BaseQuery(ctx, dependency.WithReadOnly(true), dependency.WithConds(2), dependency.WithBatchSize(3))
```

特殊配置

1. 行锁 `dependency.WithConds(2)` 与 `dependency.WithLock(true)` cond 的参数必须为索引或者主键，防止出现表锁问题
2. 指定库表 `dependency.WithTableName("table")` 与 `dependency.WithDataBase("db")`
3. 分库分表 `dependency.WithShardingKey([]interface{}{1,2})`，同时实体结构需要实现如下方法：

```go
// ITableSharding split table by keys
type ITableSharding interface {
	TableSharding(keys ...any) string
}

// IDbSharding split database by keys
type IDbSharding interface {
	DbSharding(keys ...any) string
}

// 实现
func (s testStructShardingPo) TableSharding(keys ...any) string {
	if len(keys) == 0 {
		return s.TableName()
	}
	return fmt.Sprintf("%s_%v", s.TableName(), keys[0])
}
func (s testStructShardingPo) DbSharding(keys ...any) string {
	if len(keys) < 2 {
		return s.TableName()
	}
	return fmt.Sprintf("%s_%v", s.Database(), keys[1])
}

```

4. 自定义主键生成 `dependency.WithIDGenerate(f)`,同时实体结构需要实现如下方法：

```go

// IGenerateID customer id generate
type IGenerateID interface {
	SetID(id any)
}
// 实现
func (s *testStructIdGeneratePo) SetID(id any) {
	if v, ok := id.(int64); ok {
		s.Id = v
	}
}
```

使用事务

```go
uok := NewUnitOfWork("db") // 初始化一个事务对象，它将存于context中往下传递
repo := &BaseRepository[testStructPo]{}
// 执行操作1 返回err panic 则回滚
f1 := func(subCtx context.Context) error {
	_, _ = repo.BaseUpdate(subCtx, &testStructPo{
		Id:     1,
		Code:   "122",
		Status: 2,
	})
	return right
}
// 执行操作2 返回err panic 则回滚
f2 := func(subCtx context.Context) error {
	_, err := repo.BaseCreate(subCtx, []*testStructPo{
		{
			Code: "1221",
		},
	})
	return err
}
// 执行操作
err := uok.Execute(ctx, f1, f2)
```
