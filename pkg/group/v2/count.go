package v2

func Count[T any](srcs []T, opts ...Option) int {
	return count(srcs, newOptions(opts...))
}

func count[T any](srcs []T, opt *options) int {
	total := len(srcs) // 获取切片p的长度
	if total == 0 {    // 判断切片p是否为空
		return 0 // 返回0表示无法组成批次
	}
	batch := opt.GetRealBatch(total)
	count := 0
	if int(total)%batch == 0 { // 判断切片p的元素数量能否整除batch
		count = int(total) / batch // 返回能够整除的结果，即批次数量
	} else {
		count = int(total)/batch + 1 // 返回无法整除的结果，即批次数量加一
	}
	return count
}

func batchByCount[T any](srcs []T, opt *options) (int, int) {
	gCount := count(srcs, opt)
	return opt.GetRealCount(len(srcs), gCount)
}
