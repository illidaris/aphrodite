# Aphrodite 项目代码审查报告

## 1. 结论摘要

本次审查覆盖仓库中的 303 个 Go 源文件、77 个测试文件，重点检查了 HTTP 输入、OAuth2、加密、数据库事务、组件注册、后台任务、二维码处理以及 CI 质量门禁。

当前没有发现可以仅凭仓库代码确认的 P0（立即造成远程代码执行或完整数据泄露）问题；确认 6 项 P1、5 项 P2 和 2 项 P3 问题。建议先处理并发组件注册、二维码不可信输入、解密边界和事件投递可靠性，再恢复测试及 CI 门禁。

验证基线：

- `go test ./... -count=1`：失败，涉及 `cache`、`component/gormex`、`pkg/convert`。
- `GOCACHE=/tmp/aphrodite-go-cache go vet ./...`：通过。
- `go test -race ./component/embedded ./idgenerate/idsnow -count=1`：现有测试通过，但现有用例没有覆盖并发注册/读取。
- 本机未安装 `golangci-lint`、`govulncheck`，因此未取得这两项扫描结果。

## 2. 问题清单

### P1-1：全局组件注册表并发不安全，可能导致数据竞争或进程崩溃

**证据**：`component/embedded/component.go:32-42` 使用非并发安全的包级 `rand.Rand`；`component/embedded/component.go:55-87` 直接读写公开的 `Writers`、`Readers` map 和负载函数，没有任何同步。MySQL、Mongo、Elastic 等组件均以该类型作为包级全局实例。

**影响**：服务启动后若动态注册节点、配置热更新或多个 goroutine 同时选取节点，会产生 data race；map 并发读写可直接触发 `fatal error: concurrent map read and map write`。`rand.Rand` 的并发调用也不安全。

**建议方案**：

1. 将 map、切片和负载函数改为非导出字段，禁止调用方绕过同步。
2. 使用 `sync.RWMutex` 保护注册、读取和负载策略更新；读取时复制节点切片后再执行用户提供的 balance 函数，避免持锁调用外部代码。
3. 使用并发安全的包级 `math/rand` 函数，或为自建 `rand.Rand` 单独加锁；更推荐轮询/加权轮询以获得可预测性。
4. 增加并发注册、读取、替换策略的测试，并执行 `go test -race ./component/embedded ./component/...`。

**验收标准**：持续并发注册和读取 1 万次时 race detector 无报告，且不存在 map panic。

### P1-2：二维码输入可触发 SSRF、任意本地文件读取和无界内存占用

**证据**：`pkg/qrcodes/parse.go:17-22` 将调用方传入的字符串直接交给 `ReadFile`；`pkg/qrcodes/qrcodes.go:73-85` 自动把字符串解释为 URL、本地路径或 Base64；`pkg/qrcodes/qrcodes.go:97-103` 使用无超时的 `http.Get` 并通过 `io.ReadAll` 无限制读取响应，也未校验 HTTP 状态码。

**影响**：如果 HTTP API 将用户输入传给 `ParseQrCode`，攻击者可以访问云元数据、内网管理端点，探测本机文件是否存在并尝试读取，还可通过慢响应或超大响应耗尽连接、goroutine 和内存。

**建议方案**：拆分为明确的 `ParseQRCodeBytes`、`ParseQRCodeURL` 和受控文件 API，不再自动猜测输入类型。URL 模式使用注入的 `http.Client`、请求上下文、连接/总超时、响应体上限；拒绝 loopback、link-local、私网和非 HTTP(S) 目标，并在重定向后重新校验地址。文件模式使用 Go 1.24 `os.Root` 将路径限制在业务目录。图片解码前后均限制字节数和像素数。

**验收标准**：针对 `127.0.0.1`、`169.254.169.254`、私网 DNS、重定向绕过、超大响应、慢响应和 `../` 路径均有拒绝测试。

### P1-3：图片合成参数可导致除零 panic 或超大内存分配

**证据**：`pkg/qrcodes/qrcodes.go:19-21` 明确允许 `logoP == 0`，但 `pkg/qrcodes/qrcodes.go:48` 随后执行除法；`zoom` 未校验，`pkg/qrcodes/qrcodes.go:32-33` 会把负数或过大尺寸转换为 `uint` 后交给缩放库。

**影响**：不可信参数可稳定触发 panic；异常尺寸可能造成 CPU/内存型拒绝服务。该问题位于通用包中，不能假定所有上游均有校验。

**建议方案**：要求 `logoP` 在 `[1, 10]`，为 `zoom` 设置明确的小范围；使用 checked arithmetic 计算宽高并限制最大像素数。返回参数错误而不是依赖恢复中间件。增加 `logoP=0`、负数、极大值及超大图片测试。

### P1-4：解密 API 对畸形密文处理不完整，可被输入触发 panic

**证据**：`pkg/encryptex/cbc.go:46-58` 在调用 `CryptBlocks` 前未校验密文长度是否为块大小的整数倍；`pkg/encryptex/ecb.go` 的 `CryptBlocks` 对非整块输入直接 panic；`pkg/encryptex/padding.go:52-61` 只检查最后一个 padding 数值是否越界，不验证其为正及所有填充字节一致；`ZerosUnPadding` 在空输入或全零输入上会越界（`pkg/encryptex/padding.go:73-78`）。

**影响**：只要密文来自请求、消息或数据库中的非可信数据，攻击者或损坏数据即可使服务 panic；宽松的 PKCS#7 校验还会静默接受损坏明文。

**建议方案**：所有解密入口先验证 block 非 nil、密文非空且 `len(src)%blockSize == 0`；严格验证 PKCS#7 的 padding 范围和每个尾字节；安全处理空/全零输入。长期将新业务迁移到 AES-GCM 等 AEAD，ECB、DES、3DES、裸 CBC 仅保留为明确标注的兼容 API。

**验收标准**：对任意字节输入做 fuzz 测试时无 panic；被篡改密文必须返回认证或格式错误，不能返回明文。

### P1-5：OAuth2 无缓存回退 state 仅加密不认证，且编码错误被吞掉

**证据**：`biz/oauth2/aes.go:13-19` 使用 SM4-CTR；`pkg/encrypter/stream.go:79-103` 只生成 IV 并流加密，没有 MAC/AEAD tag；`biz/oauth2/dto.go:26-28` 忽略编码错误并返回空字符串。`biz/oauth2/oauth2.go:41-45` 支持直接从回调 state 解出 verifier、BizId 和过期时间。

**影响**：CTR 提供机密性但不提供完整性，state 被修改时不能在解密层可靠拒绝。当前 PKCE 和业务校验降低了直接接管风险，但安全边界依赖 JSON 恰好解析失败或后续字段校验，不适合作为鉴权状态保护；空 secret 还是默认配置。

**建议方案**：首选只在服务端保存随机 state，回调时原子读取并删除；确需无状态 state 时使用 AES-GCM/标准 AEAD，将版本、用途和租户作为附加认证数据，并强制非空、足够长度的密钥。把 `AuthorizeParam.Encode` 改为 `(string, error)`，任何加密/缓存错误均 fail closed。增加篡改、重放、过期、错租户和空密钥测试。

### P1-6：事件发布失败后的补偿写入异步且忽略错误，可能永久丢失事件

**证据**：`event/event_db_uow.go:39-43` 提交本地消息事务后同步发布；`event/event_db_uow.go:44-66` 在后台 goroutine 中更新失败状态或删除消息，但两次数据库操作均忽略错误；函数最后始终返回 nil。goroutine 没有生命周期管理，panic 只输出到 stdout。

**影响**：发布失败且补偿状态更新也失败时，调用方仍收到成功；消息可能停留在无法正确重试的状态。发布成功但删除失败会造成重复投递。进程退出时后台更新也可能尚未完成。

**建议方案**：采用标准 transactional outbox：事务内只写待发送记录，独立、受生命周期管理的 worker 拉取并发布，通过带条件的状态迁移记录重试、幂等键和最后错误；只有状态持久化成功后才确认处理。至少在现结构中同步检查 `BaseUpdate/BaseDelete` 错误并记录指标/告警，不要用裸 goroutine。

**验收标准**：注入 Kafka 失败、数据库失败、进程重启及重复消费后，事件最终可重试且下游以幂等键去重，不出现静默丢失。

### P2-1：全量测试当前不通过，且 GORM 测试受全局状态污染

**证据**：本次 `go test ./... -count=1` 中 `cache/TestNewsInfoForCache`、多个 `component/gormex` 用例以及 `pkg/convert/TestTruncateMySQLVarchar` 失败。GORM 日志显示后续测试取得了已关闭数据库；测试持续向包级 `MySqlComponent` 追加同名 writer，而组件没有清空、替换或移除能力。

**影响**：回归结果不可信，真实缺陷容易被失败噪声掩盖；全局注册状态也使测试顺序和进程生命周期相互耦合。

**建议方案**：业务组件通过构造函数注入，包级默认实例仅作为兼容层；注册 API 明确支持 replace/remove，并为测试创建独立 component。Redis mock 需匹配 Lua 参数。先把三组失败逐一固定为可重复的单包测试，再恢复全量绿灯。

### P2-2：`TruncateMySQLVarchar` 未按 MySQL 字节长度截断

**证据**：`pkg/convert/string.go:17-26` 的注释和测试要求按存储字节长度截断，但实现按 rune 数量截断；本次测试中的中文、emoji 和混合字符串均失败。

**影响**：多字节文本可能超过按字节定义的外部限制，导致写库失败或上游/下游长度不一致。

**建议方案**：先明确目标是“UTF-8 字节预算”还是“MySQL `VARCHAR(n)` 字符数”。若是字节预算，逐 rune 累加 `utf8.RuneLen`，绝不截断 UTF-8 序列；若是 MySQL `VARCHAR(n)`，函数和测试应改为字符语义，并另设行字节上限检查。

### P2-3：数据库结果遍历结束后未检查 `rows.Err()`

**证据**：`component/dbex/query_scan2map.go:44-69` 在 `rows.Next()` 结束后直接返回结果。

**影响**：网络中断、驱动解码或游标错误可能被当作正常 EOF，调用方收到不完整结果且 error 为 nil，形成静默数据截断。

**建议方案**：循环后检查并返回 `rows.Err()`；同时检查 `stmt.Close`/`rows.Close` 中对业务有意义的错误。用自定义 rows/mock 添加“读取若干行后失败”的回归测试。

### P2-4：Prometheus 推送后台任务无法停止，HTTP 请求没有超时

**证据**：`ginhandle/middleware/prometheus/promethus.go:291-310` 使用无超时 `http.Get`；`330-341` 创建无超时客户端；`343-349` 的 ticker 不会 `Stop`，goroutine 没有 context 或关闭方法。重复调用 `SetPushGateway` 会重复启动永久 worker。

**影响**：目标端挂起时 goroutine/连接长期占用；配置重载会泄漏 ticker 和 worker，关闭服务时也无法优雅等待。

**建议方案**：让组件实现 `Start(ctx)`/`Close()`，保存 ticker 和 cancel；注入带超时及连接池限制的 `http.Client`，使用 `NewRequestWithContext`，校验非 2xx 响应。保证多次 Start 幂等。

### P2-5：CI/Makefile 将关键安全检查降级为“永远成功”

**证据**：`Makefile:4-5` 和 `Makefile:11-12` 对 `golangci-lint`、`govulncheck` 使用 `|| true`；`.golangci.yml:13` 排除测试文件，`:16-17` 固定不存在业务含义的 `mytag` build tag，可能改变实际分析文件集合。

**影响**：即使静态检查或依赖漏洞扫描发现问题，流水线仍可通过；测试代码中的资源泄漏、竞态和错误断言不会被 lint。

**建议方案**：质量门禁去掉 `|| true`，将工具版本固定在 CI；删除无依据的 build tag，启用 tests；遗留问题使用基线或 `new-from-rev` 渐进收敛，而不是吞掉退出码。增加 `go test -race ./...` 和覆盖率报告。

### P3-1：恢复中间件的 broken-pipe 判断失效，并存在二次 panic 风险

**证据**：`ginhandle/middleware/recover.go:28-34` 用一个 nil 的 `*os.SyscallError` 调用 `errors.Is`，无法获得实际错误对象；`:41` 又直接执行 `err.(error)`，而 Go 允许 panic 任意值。

**影响**：连接断开会被记录为完整 panic，增加噪声；若未来修正分支判断但 panic 值不是 error，会在恢复逻辑中再次 panic。

**建议方案**：使用 `errors.As(ne.Err, &se)`，并对 panic 值通过类型 switch 转为 error；删除 `fmt.Println`，统一结构化日志并避免记录敏感请求信息。

### P3-2：部分 API 以 panic 或静默错误表示“未实现/失败”

**证据**：如 `idgenerate/idsnow/id_generate.go:27-29` 吞掉 ID 生成错误，`:40-42` 对公开方法直接 `panic("no impl")`；仓库其他包也存在同类占位实现。

**影响**：调用方无法区分合法零值和生成失败；公开接口在运行期才崩溃，接口实现断言给出了错误的完整性信号。

**建议方案**：删除未实现方法对应的接口声明，或返回明确的 `ErrNotImplemented`；废弃不返回 error 的 `NewIDX`，业务路径统一使用 `NewID`。对所有 `panic("no impl")` 建清单并在发布前消除。

## 3. 建议实施顺序

1. **立即（1-3 天）**：修复二维码参数 panic、解密畸形输入 panic；限制 URL/文件读取；让现有全量测试恢复通过。
2. **当前迭代**：重构组件注册同步与测试隔离；OAuth2 state 改为服务端一次性存储或 AEAD；补齐 race/fuzz/篡改测试。
3. **下一迭代**：落地 transactional outbox worker；为 Prometheus 后台任务增加生命周期；修复 `rows.Err()`。
4. **持续治理**：恢复 lint、漏洞扫描、race 门禁，逐步废弃 ECB/DES/3DES/裸 CBC 和吞错 API。

## 4. 推荐新增测试

- `component/embedded`：并发 `NewWriter/GetWriter/SetWriterBalance` race 测试。
- `pkg/qrcodes`：参数边界、SSRF 地址、重定向、超时、响应体/像素上限测试。
- `pkg/encryptex`：CBC/ECB/padding fuzz 测试，约束“任意输入不 panic”。
- `biz/oauth2`：state 单字节篡改、重放、过期、错租户、空密钥测试。
- `event`：发布失败 + 状态更新失败、进程重启、重复投递的集成测试。
- `component/dbex`：`rows.Next` 中途失败时必须返回 error。

