# aphrodite

`aphrodite` 是一个基于 Go 1.24 的服务端基础框架，提供了组件抽象（DB / MQ / ES / Mongo）、业务编排（CRUD / OAuth2 / 事件）、HTTP / RPC 接入（Gin / Dubbo）、ID 生成、缓存、加解密、以及一组通用工具包。框架以「显式读写分离 + 选项模式 + 泛型 Repository」为核心思想，目标是让上层业务可以用极少的样板代码完成常见后端需求。

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

通过 `embedded.NewComponent[T]()` 即可创建一个具备读写分离、可插拔负载均衡的组件实例，框架内已内置 MySQL / Mongo / Elastic / Kafka 等组件。

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

特殊场景

1. 批量插入，忽略已经存在的，只新增未入库的（覆盖导入建议使用删除，导入形式）。当插入的数量大于`opt.BatchSize`时会触发分组协程操作，默认使用多协程插入。

```go
repo := &BaseRepository[testStructPo]{}
affect, err := repo.BaseCreate(ctx, pos,
    dependency.WithIgnore(true), // 忽略插入失败，否则失败则返回错误并且不执行下去
	dependency.WithBatchSize(1000), // 1000条为一组执行，默认1000
)
```

### kafkaex sarama 扩展组件

基于 `IBM/sarama` 的 Kafka 客户端封装，提供统一的 `Manager` 管理生产者、消费者组、确认（Receipt）以及主题元数据，与 `event` 包配合可以做到「发送失败 → 落库补偿 → 重发」的最终一致投递。

```go
mgr, err := kafkaex.InitDefaultManager(opts...)
// 同步发送
_, _, err = mgr.SyncSend(ctx, "topic", "key", []byte("payload"))
// 启动消费者组
go mgr.ConsumeGroup(ctx, "group-id", []string{"topic"}, handler)
```

### mongoex mongo 扩展组件

基于 `go.mongodb.org/mongo-driver` 的 Mongo 封装，提供 `Repository[T]` 通用 CRUD、事务（Session）、变更追踪（VAO）等能力，接口风格与 `gormex` 保持一致，便于在 SQL / NoSQL 仓储之间平滑切换。

### elasticex elasticsearch 扩展组件

基于 `olivere/elastic/v7` 的 ES 客户端封装，包含索引管理（`repository_index*`）、查询（`repository_query`）与限流搜索（`limit_search`），方便在业务层使用统一的 Repository 风格调用 ES。

### dbex 原生 SQL 辅助

`QueryScan2Map` 等工具方法，用于将 `sql.Rows` 直接扫描为 `[]map[string]any`，适合不便建模的临时查询或脚本类任务。

## biz 业务层

业务层封装了一些与具体业务无关、但在大多数后台系统都会重复出现的能力，避免在每个项目里都重写一遍。

### biz/crud 通用 CRUD

将仓储层的 `BaseCreate / BaseUpdate / BaseQuery ...` 与 `dto.Page` 等组合成可直接挂在路由上的函数，支持入参 / 出参的类型转换：

```go
createFn := crud.Create(repo, nil)
affected, err := createFn(ctx, []*Entity{entity})

detailFn := crud.DetailByIdFunc(repo)
po, err := detailFn(ctx, []any{"db"}, []any{"table"}, int64(123))
```

### biz/oauth2 OAuth2 客户端

通用 OAuth2 授权码 + PKCE 流程，包含 AES 加密的 state、回调处理、token 缓存等。

```go
url, err := oauth2.GetAuthorizeURl(ctx, opts...)
token, user, err := oauth2.OAuthCallback(ctx, code, state, opts...)
```

### biz/ginoauth2 Gin OAuth2 路由

在 `biz/oauth2` 之上提供 Gin 路由处理器：登录跳转、登录链接返回、回调处理，可一行挂载到 Gin Engine 上。

### biz/api 业务请求客户端

对 HTTP API 调用做了 Get / Post / Put / Delete 的封装，统一了选项模式、签名与上下文透传。

### biz/bizlog 业务日志

为业务流水提供结构化日志接口，支持选项扩展。

## event 事件 / 工作单元

事件发布管理器，默认通过 Kafka 投递，失败时落库到 `po.MqMessage` 形成补偿队列；`event_db_uow` 提供数据库事务与事件发布的一体化协调（本地消息表）。

```go
event.InitDefault(producer, repo) // 初始化默认事件管理器
event.Publish(ctx, eventArgs)     // 业务侧只调用 Publish
```

## ginhandle Gin HTTP 处理框架

为常用的 HTTP 处理范式抽出模板函数（`GinExHandler / BizGinExHandler / GinOneHandler` 等），自动完成请求绑定、参数校验、异常包装与统一响应；`crud.go` 进一步将 CRUD 仓储接入 Gin 路由。

### ginhandle/middleware 中间件

- `log.go` — 访问日志中间件，记录方法、路径、状态码、耗时
- `param.go` — 请求 / 响应体记录，支持大小限制
- `sign.go` — Web 签名验证（应用签名 + 时间戳 + 版本号）
- `recover.go` — panic 恢复
- `prometheus/` — Prometheus 指标采集

## micro/dubboex Dubbo 接入

对 `dubbo-go/v3` 的封装，支持从 Nacos 配置中心动态拉取 Dubbo 配置后再启动实例：

```go
dubboex.InitFrmDubboNacos() // 从 Nacos 拉取配置并初始化 Dubbo
dubboex.NewInstance(opts...)
```

## idgenerate ID 生成

提供两种生产可用的全局唯一 ID 方案，统一实现 `dep.IIDGenerate` 接口，便于上层按需切换。

### idsegment 号段模式

依赖数据库 segment 表分配号段，本地内存预取下一段，适合对趋势递增有要求的场景。

```go
gen := &idsegment.IdSegment{Batch: 1000, Cache: cache, Repo: repo}
id, err := gen.NewID(ctx, "user_id")
```

### idsnow Snowflake 模式

基于 `pkg/snowflake` 的多实例负载均衡，支持自定义机器位 / 时间位 / 业务基因位。

```go
gen := idsnow.NewIdGenerate()
gen.Run(ctx, "tmp", 4) // 启动 4 个 snowflake worker
id, err := gen.NewID(ctx, "user_id")
```

## cache 缓存

- `limit.go` — 基于 Redis + Lua 的滑动窗口限流器，配置最大请求数与窗口时长
- `shell.go` — 缓存外壳，统一处理键生成、过期时间、回源、强制刷新

```go
opts := cache.NewLimitOptions(
    cache.WithLimitCache(redisClient),
    cache.WithLimitMax(100),
    cache.WithLimitDur(time.Second),
)
allow, err := cache.Limit(ctx, "key", opts)
```

## dto 通用 DTO

- `request_base / response_base` — 统一响应壳 `Response[T]`（code / message / data）
- `request_page / response_page` — `Page`（page / pageSize / sorts）+ `Pager[T]`（含 totalRecord / totalPage）
- `request_biz` — 含 bizId 的业务请求基类
- `request_range` — `Range`（Beg / End）范围查询
- `request_schedule` — `ScheduleBase`（batch / timeout）调度参数
- `request_ip` — 客户端 IP 提取

## po 通用持久化对象

预定义的 GORM 嵌入字段与基础消息模型：

- `section.IDAutoSection` — 自增主键
- `section.RawBizSection` / `BizSection` — 业务 ID（带分片）
- `section.OperationSection` — 状态 / 创建者 / 修改时间
- `section.LockSection` — 分布式锁字段（locker / expire / timeout / retries）
- `TaskQueueMessage` — 异步任务 / 补偿消息持久化模型
- `MqMessage` — Kafka 发送失败时的补偿消息
- `EventArgs` — 事件参数（id / 主题 / key / value）

## pkg 工具集

### pkg/app 应用初始化

启动时打印 Go 运行时信息并通过 `automaxprocs` 适配容器 CPU 配额。

### pkg/dependency 接口契约

框架最核心的接口集合：`IEntity / IPo / IRepository / ICache / ILog / IPage / IMessage` 等；`BaseOption` 与 `WithXxx` 选项函数也在这里定义。

### pkg/contextex 上下文扩展

定义跨组件透传的 Context Key（`ElasticID / DbTxID / MongoID` 等），用于在事务、追踪、日志间共享上下文。

### pkg/exception 异常

统一异常类型 `Exception(Code, SubCode, Msg)`，支持 `Wrap` 链式封装，便于跨层透传错误码。

```go
ex := exception.ERR_BUSI_CREATE.Wrap(err)
```

### pkg/group 并发分组

对大批量数据按指定 size 切片并发执行函数，自动聚合结果与错误。

```go
total, errs := group.GroupFunc(func(items ...*User) (int64, error) {
    return repo.BaseCreate(ctx, items)
}, 100, users...)
```

### pkg/imex 导入导出

基于分页拉取 + 流式写出的导入导出框架，支持 CSV / Excel，适合大数据量报表导出。

### pkg/encrypter KMS

KMS（密钥管理）抽象：`IKmsAdapter / IKmsStore / IKmsCache`，内置嵌入式与腾讯云 KMS 适配；支持 DEK 生成、加密、缓存与流式加解密。

### pkg/encryptex 加解密算法

通用对称 / 非对称算法集合：AES（ECB / CBC + 各种 padding）、3DES、DES、RSA、MD5、SHA。

```go
ciphertext, err := encryptex.AesCBCEncrypt(data, key, iv, "PKCS5")
```

### pkg/snowflake Snowflake

可配置位长的 Snowflake 实现，提供 `NextIdFunc` 返回一个线程安全的取号闭包。

### pkg/check 业务校验

文件上传校验（扩展名 + 文件头）、手机号、邮箱等常用校验。

### pkg/redq Redis 队列

基于 `hibiken/asynq` 的 Redis 任务队列，封装了服务端启动与发送端的简化调用。

```go
_ = redq.InitRedqSrv(opts...)
send := redq.SendFunc(opts...)
_, err := send(ctx, "topic", []byte("payload"))
```

### pkg/logex 日志

`go.uber.org/zap` 的封装，支持自动从 Context 中提取链路字段。

### pkg/structure 数据结构

通用数据结构工具：唯一过滤器 `NewUnqueFilter`、ID 切片 `IIDSection`、`ItemMap`、对象 diff / compare 等。

### pkg/convert 类型转换

`any2struct / ConvertToStructByJson / ConvertToStructByRef`、身份证号解析、时间转换、UUID、字段过滤等。

### pkg/proxy HTTP 反向代理

可注入请求 / 响应钩子的 HTTP 反向代理。

```go
p, _ := proxy.NewProxy("http://upstream:8080",
    proxy.WithRequestHooks(rewriteReq),
    proxy.WithResponseHook(rewriteResp))
```

### pkg/smtp 邮件

支持 TLS 与腾讯云特殊握手的 SMTP 发送函数。

### pkg/qrcodes 二维码

二维码生成与中心 Logo 合成。

### pkg/ollama Ollama 客户端

Ollama LLM 客户端：模型管理、推理 / 流式响应、鉴权与环境配置。

### pkg/netex 网络工具

获取空闲端口、IP 段、地区行政码等。

### pkg/backup 本地备份

将对象序列化为 JSON 落盘 / 加载（`DiskSave` / `DiskLoad`），用于轻量级状态持久化。

## License

详见 [LICENSE](LICENSE)。
