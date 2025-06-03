package v2

// Count 计算切片元素按指定批次参数划分后的总批次数量
//   - srcs : 源数据切片，支持泛型类型
//   - opts : 可变长度配置选项，用于自定义批次参数
//
// 返回值：计算得出的批次总数
func Count[T any](srcs []T, opts ...Option) int {
	return count(srcs, newOptions(opts...))
}

// count 内部批次计算核心逻辑，根据配置参数计算实际批次数量
//   - srcs : 源数据切片
//   - opt  : 经过初始化的配置对象指针
//
// 返回值：最终计算得出的批次数量
func count[T any](srcs []T, opt *options) int {
	// 基础校验：空切片直接返回0批次
	total := len(srcs)
	if total == 0 {
		return 0
	}

	// 动态批次计算逻辑：考虑整除和余数两种情况
	batch := opt.GetRealBatch(total)
	count := 0
	if total%batch == 0 {
		count = total / batch
	} else {
		count = total/batch + 1 // 有余数时增加一个批次
	}
	return count
}

// batchByCount 计算批次划分后的实际批次参数
//   - srcs : 源数据切片
//   - opt  : 配置对象指针
//
// 返回值：返回两个整数值，分别表示实际批次数量和每批元素数量
func batchByCount[T any](srcs []T, opt *options) (int, int) {
	gCount := count(srcs, opt)
	// 通过配置对象获取最终批次参数
	return opt.GetRealCount(len(srcs), gCount)
}
