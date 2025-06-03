package v2

// DEF_MAX_GROUP_COUNT 定义默认的最大并行组数
const (
	DEF_MAX_GROUP_COUNT = 10000
)

// options 包含并行处理配置参数
// ParallelismMax: 允许的最大并行处理组数
// Batch: 每个处理组的批次大小基准值
type options struct {
	ParallelismMax int
	Batch          int
}

// GetRealBatch 计算实际的批次大小
// srcsLen: 源数据总量长度
// 返回值: 修正后的有效批次大小(保证至少为1)
func (o options) GetRealBatch(srcsLen int) int {
	if o.Batch < 1 {
		return 1
	}
	return o.Batch
}

// GetRealCount 计算实际的批次数量和并行度
// srcsLen: 源数据总量长度
// count: 请求的原始并行度
// 返回值1: 调整后的实际每批处理量
// 返回值2: 调整后的实际最大并行组数
func (o options) GetRealCount(srcsLen, count int) (int, int) {
	realBatch := o.GetRealBatch(srcsLen)

	// 当请求的并行度超过最大值时，重新计算批次分割方案
	if count > o.ParallelismMax {
		// 根据总量和最大并行度计算每批处理量
		if srcsLen%o.ParallelismMax == 0 {
			realBatch = srcsLen / o.ParallelismMax
		} else {
			realBatch = srcsLen/o.ParallelismMax + 1
		}
		return realBatch, o.ParallelismMax
	}
	return realBatch, count
}

// newOptions 创建配置选项实例
// opts: 可变长度配置选项函数集
// 返回值: 初始化后的配置对象指针
func newOptions(opts ...Option) *options {
	// 初始化默认配置值
	opt := &options{
		ParallelismMax: DEF_MAX_GROUP_COUNT,
		Batch:          1,
	}
	// 应用所有配置选项函数
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// Option 定义配置选项的函数类型
type Option func(*options)

// WithParallelismMax 设置最大并行度配置选项
// parallelism: 最大允许的并行处理组数
func WithParallelismMax(parallelism int) Option {
	return func(o *options) {
		o.ParallelismMax = parallelism
	}
}

// WithBatch 设置批次大小配置选项
// batch: 每个处理组的基准批次数量
func WithBatch(batch int) Option {
	return func(o *options) {
		o.Batch = batch
	}
}
